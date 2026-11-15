package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"iot-warehouse-telemetry/internal/auth"
	"iot-warehouse-telemetry/internal/store"
	"iot-warehouse-telemetry/pkg/logger"
	"iot-warehouse-telemetry/pkg/sensor"
	"iot-warehouse-telemetry/pkg/ws"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

func newApp(ho *slog.HandlerOptions, addr string) (*application, error) {
	logger := logger.NewLogger(ho)

	authSvc, err := auth.New(os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, err
	}

	db, err := OpenDB()
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	app := &application{
		mux:    http.NewServeMux(),
		logger: logger,
		store: &store.Store{
			DB: db,
		},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		auth: authSvc,
	}

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s", addr)).
		SetClientID(fmt.Sprintf("go-server-%d", os.Getpid())).
		SetConnectTimeout(5 * time.Second).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(c mqtt.Client) {
			if err := app.SubscribeMqttAndInser(); err != nil {
				logger.Error("mqtt subscribe failed", "err", err)
			}
		}).
		SetConnectionNotificationHandler(func(c mqtt.Client, cn mqtt.ConnectionNotification) {
			switch v := cn.(type) {
			case mqtt.ConnectionNotificationConnected:
				logger.Info("mqtt connected", "broker", addr)
			case mqtt.ConnectionNotificationLost:
				logger.Error("mqtt connection lost", "err", v.Reason)
			case mqtt.ConnectionNotificationConnecting:
				if v.IsReconnect {
					logger.Warn("reconnecting to mqtt", "broker", addr)
				}
			case mqtt.ConnectionNotificationFailed:
				logger.Error("mqtt connect/reconnect failed", "err", v.Reason)
			case mqtt.ConnectionNotificationBrokerFailed:
				logger.Error("mqtt broker attempt failed", "broker", v.Broker, "err", v.Reason)
			}
		})

	app.mqttClient = mqtt.NewClient(opts)

	if token := app.mqttClient.Connect(); !token.WaitTimeout(5*time.Second) || token.Error() != nil {
		db.Close()
		if token.Error() != nil {
			return nil, fmt.Errorf("connect to mqtt: %w", token.Error())
		}
		return nil, fmt.Errorf("connect to mqtt: timeout")
	}

	return app, nil
}

func (app *application) subscribeTopic(topic string, handler mqtt.MessageHandler) error {
	token := app.mqttClient.Subscribe(topic, 0, handler)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("subscribe to %s: timeout", topic)
	}
	if token.Error() != nil {
		return fmt.Errorf("subscribe to %s: %w", topic, token.Error())
	}
	return nil
}

func (app *application) SubscribeMqttAndInser() error {
	return app.subscribeTopic(sensor.DATA_TOPIC, app.handleSensorData)
}

func (app *application) handleSensorData(_ mqtt.Client, msg mqtt.Message) {
	var record sensor.SensorRecord
	if err := json.Unmarshal(msg.Payload(), &record); err != nil {
		app.logger.Error("unmarshaling sensor record", "err", err)
		return
	}

	if _, err := app.store.InsertRecord(record); err != nil {
		app.logger.Error("inserting sensor record", "err", err)
		return
	}

	payload := struct {
		sensor.SensorRecord
		MessageType int `json:"message_type"`
	}{
		MessageType:  ws.SENSOR_RECORD_MSG,
		SensorRecord: record,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		app.logger.Error("marshaling payload", "err", err)
		return
	}

	if err := ws.Broadcast(data); err != nil {
		app.logger.Error("broadcast sensor record to websocket", "err", err)
		return
	}

	var alertErr error
	switch {
	case sensor.MAXTEMP <= record.Temperature:
		alertErr = app.broadcastAlert(ws.SENSOR_ALERT_TEMP_MSG, "max temperature reached")
	case sensor.MINTEMP >= record.Temperature:
		alertErr = app.broadcastAlert(ws.SENSOR_ALERT_TEMP_MSG, "min temperature reached")
	}
	if alertErr != nil {
		app.logger.Error("broadcasting temperature alert", "err", alertErr)
	}

	alertErr = nil
	switch {
	case sensor.MAXHUM <= record.Humidity:
		alertErr = app.broadcastAlert(ws.SENSOR_ALERT_HUM_MSG, "max humidity reached")
	case sensor.MINHUM >= record.Humidity:
		alertErr = app.broadcastAlert(ws.SENSOR_ALERT_HUM_MSG, "min humidity reached")
	}
	if alertErr != nil {
		app.logger.Error("broadcasting humidity alert", "err", alertErr)
	}
}

func (app *application) SubscribeMqttForSensorsNumber() error {
	return app.subscribeTopic(sensor.SENSORS_NUMBER_TOPIC, func(_ mqtt.Client, msg mqtt.Message) {
		n, err := strconv.ParseInt(string(msg.Payload()), 10, 64)
		if err != nil {
			app.logger.Error("parsing active sensors number", "err", err)
			return
		}
		app.activeSensorsNumber.Store(n)
	})
}

func (app *application) start(addr string) error {
	app.logger.Info("starting app", "addr", addr)
	return http.ListenAndServe(addr, app.mux)
}

func OpenDB() (*sql.DB, error) {
	dsn, ok := os.LookupEnv("MYSQL_DSN")
	if !ok || dsn == "" {
		return nil, fmt.Errorf("MYSQL_DSN is unset or empty")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (app *application) httpError(w http.ResponseWriter, status int, err error) {
	app.logger.Error("http request error", "status", status, "err", err)
	http.Error(w, http.StatusText(status), status)
}

func closeWSClient(client *ws.Client) {
	ws.Mu.Lock()
	defer ws.Mu.Unlock()
	client.Conn.Close()
	delete(ws.Clients, client)
}

func (app *application) sendSensorNumber(client *ws.Client) {
	count, err := app.store.GetSensorsNumber()
	if err != nil {
		app.logger.Error("getting sensors number", "err", err)
		closeWSClient(client)
		return
	}

	sn := struct {
		SensorsNumber int `json:"sensors-number"`
	}{
		SensorsNumber: count,
	}

	data, err := json.Marshal(sn)
	if err != nil {
		app.logger.Error("marshaling sensors number", "err", err)
		return
	}

	ws.Mu.Lock()
	err = client.Conn.WriteMessage(websocket.TextMessage, data)
	ws.Mu.Unlock()
	if err != nil {
		closeWSClient(client)
		app.logger.Error("writing sensors number to websocket", "err", err)
		return
	}
}

func (app *application) broadcastAlert(msgType int, msg string) error {
	dataSensorAlert := struct {
		MessageType int    `json:"message_type"`
		Message     string `json:"message"`
	}{
		MessageType: msgType,
		Message:     msg,
	}

	data, err := json.Marshal(dataSensorAlert)
	if err != nil {
		return fmt.Errorf("marshal alert: %w", err)
	}

	if err := ws.Broadcast(data); err != nil {
		return fmt.Errorf("broadcast alert: %w", err)
	}

	return nil
}
