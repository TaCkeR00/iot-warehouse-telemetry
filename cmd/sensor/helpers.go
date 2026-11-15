package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"iot-warehouse-telemetry/pkg/logger"
	"iot-warehouse-telemetry/pkg/sensor"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func newApp(ho *slog.HandlerOptions, addr string) (*application, error) {
	brokerAddr := fmt.Sprintf("tcp://%s", addr)

	opts := mqtt.NewClientOptions().
		AddBroker(brokerAddr).
		SetClientID("go-sensor").
		SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt connect failed: %w", token.Error())
	}

	return &application{
		logger:     logger.NewLogger(ho),
		mqttClient: client,
	}, nil
}

func (app *application) publishData(data []byte, topic string) error {
	token := app.mqttClient.Publish(topic, 0, false, data)
	token.Wait()

	fmt.Println(string(data))

	return token.Error()
}

func (app *application) publishSensorNumber(sensorsNumber int) error {
	sensorsData, err := json.Marshal(map[string]int{
		"sensor-number": sensorsNumber,
	})
	if err != nil {
		return err
	}

	return app.publishData(sensorsData, sensor.SENSORS_NUMBER_TOPIC)
}

func (app *application) publishSensorsData(id int) error {
	data, err := json.Marshal(app.getSensorsVals(id))
	if err != nil {
		return err
	}

	return app.publishData(data, sensor.DATA_TOPIC)
}
