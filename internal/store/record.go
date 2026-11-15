package store

import (
	"database/sql"
	"time"

	"iot-warehouse-telemetry/pkg/sensor"

	_ "github.com/go-sql-driver/mysql"
	// influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Record struct {
	ID          int       `json:"id"`
	SensorID    int       `json:"sensor_id"`
	Temperature int       `json:"temperature"`
	Humidity    int       `json:"humidity"`
	Time        time.Time `json:"time-stamp"`
}

type Sensor struct {
	ID    int
	Model string
	Site  string
}

type Store struct {
	DB *sql.DB
}

func (s *Store) InsertRecord(record sensor.SensorRecord) (int, error) {
	stmt := `INSERT INTO records (sensor_id, temperature, humidity, time)
	VALUES (?, ?, ?, ?)`

	result, err := s.DB.Exec(stmt, record.SensorID, record.Temperature, record.Humidity, record.TimeStamp)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (s *Store) GetRecord(sensorID int, site string) (Record, error) {
	stmt := `SELECT records.id, sensor_id, temperature, humidity, time FROM records JOIN sensors ON records.sensor_id = sensors.id
	WHERE sensor_id=? and site=? ORDER BY time DESC LIMIT 1`
	var record Record

	row := s.DB.QueryRow(stmt, sensorID, site)

	if err := row.Scan(&record.ID, &record.SensorID, &record.Temperature, &record.Humidity, &record.Time); err != nil {
		return Record{}, err
	}

	return record, nil
}
