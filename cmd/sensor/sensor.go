package main

import (
	"math/rand/v2"
	"time"

	"iot-warehouse-telemetry/pkg/sensor"
)

func (app *application) getSensorsVals(id int) sensor.SensorRecord {
	return sensor.SensorRecord{
		SensorID:    id,
		Temperature: 35 + int(rand.Float64()*10), //20
		Humidity:    60 + int(rand.Float64()*20), //40
		TimeStamp:   time.Now(),
	}
}
