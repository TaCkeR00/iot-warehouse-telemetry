package sensor

import "time"

type SensorRecord struct {
	SensorID    int       `json:"sensor_id"`
	Temperature int       `json:"temperature"`
	Humidity    int       `json:"humidity"`
	TimeStamp   time.Time `json:"time_stamp"`
}

type Sensor struct {
	SensorID   int    `json:"sensor_id"`
	Model      string `json:"model"`
	WarehousId int    `json:"warehouse_id"`
	MaxTemp    int    `json:"max_temperature"`
	MinTemp    int    `json:"min_temperature"`
	MaxHum     int    `json:"max_humidity"`
	MinHum     int    `json:"min_humidity"`
}
