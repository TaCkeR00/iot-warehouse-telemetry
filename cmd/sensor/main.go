package main

import (
	"flag"
	"time"

	"iot-warehouse-telemetry/pkg/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type application struct {
	logger     *logger.Logger
	mqttClient mqtt.Client
}

func main() {
	sensorsNumber := flag.Int("sn", 1, "sensors number")
	addr := flag.String("addr", "localhost:1883", "address")
	flag.Parse()

	app, err := newApp(nil, *addr)
	if err != nil {
		panic(err)
	}
	defer app.mqttClient.Disconnect(250)

	err = app.publishSensorNumber(*sensorsNumber)
	if err != nil {
		app.logger.Error(err.Error())
	}

	for i := range *sensorsNumber {
		go func(id int) {
			for {
				if err := app.publishSensorsData(id); err != nil {
					app.logger.Error(err.Error())
				}
				time.Sleep(time.Second)
			}
		}(i + 1)
	}

	select {}
}
