package main

import (
	"flag"
	"net/http"
	"os"
	"sync/atomic"

	"iot-warehouse-telemetry/internal/auth"
	"iot-warehouse-telemetry/internal/store"
	"iot-warehouse-telemetry/pkg/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

// TODO: add influxdb client
type application struct {
	mux                 *http.ServeMux
	logger              *logger.Logger
	store               *store.Store
	mqttClient          mqtt.Client
	upgrader            websocket.Upgrader
	activeSensorsNumber atomic.Int64
	auth                *auth.Service
}

func main() {
	addr := flag.String("addr", "localhost:1883", "MQTT Broker address")
	flag.Parse()

	app, err := newApp(nil, *addr)
	if err != nil {
		panic(err)
	}
	app.setRoutes()
	defer app.store.DB.Close()

	if err := app.start(":8080"); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
