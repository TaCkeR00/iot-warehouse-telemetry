package store

import "iot-warehouse-telemetry/pkg/sensor"

func (s *Store) AddNewSensor(sensor sensor.Sensor) (int, error) {
	stmt := `INSERT INTO sensors(model, max_temperature, min_temperature, max_humidity, min_humidity, warehouse_id) VALUES(?, ?, ?, ?, ?, ?)`

	result, err := s.DB.Exec(stmt, sensor.Model, sensor.MaxTemp, sensor.MinTemp, sensor.MaxHum, sensor.MinHum, sensor.WarehousId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) GetSensorsNumber() (int, error) {
	stmt := `SELECT COUNT(*) FROM sensors`
	var count int

	err := s.DB.QueryRow(stmt).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
