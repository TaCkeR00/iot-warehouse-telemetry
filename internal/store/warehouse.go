package store

import "iot-warehouse-telemetry/pkg/warehouse"

func (s *Store) AddNewWarehouse(warehouse warehouse.Warehouse) (int, error) {
	stmt := "INSERT INTO warehouses (state, city, address) VALUES (?, ?, ?)"

	result, err := s.DB.Exec(stmt, warehouse.State, warehouse.City, warehouse.Address)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
