package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"iot-warehouse-telemetry/pkg/sensor"
	"iot-warehouse-telemetry/pkg/user"
	"iot-warehouse-telemetry/pkg/warehouse"
	"iot-warehouse-telemetry/pkg/ws"
)

type idResponse struct {
	ID int `json:"id"`
}

func (app *application) newRecordHandler(w http.ResponseWriter, r *http.Request) {
	var sensorRecord sensor.SensorRecord

	err := json.NewDecoder(r.Body).Decode(&sensorRecord)
	if err != nil {
		app.httpError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	id, err := app.store.InsertRecord(sensorRecord)
	if err != nil {
		app.httpError(w, http.StatusInternalServerError, err)
		return
	}

	json.NewEncoder(w).Encode(idResponse{
		ID: id,
	})
}

func (app *application) logInHandler(w http.ResponseWriter, r *http.Request) {
	var user user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		app.httpError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	isUserExists, isAdmin := app.store.IsUserExistsAndIsAdmin(user)
	if !isUserExists {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	tokenString, err := app.auth.GenerateJWT(user.UserName, isAdmin)
	if err != nil {
		app.httpError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func (app *application) signNewAccountHandler(w http.ResponseWriter, r *http.Request) {
	var user user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		app.httpError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	err := app.store.SignUser(user)
	if err != nil {
		app.httpError(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Fprint(w, "{\"status\": \"ok\"}")
}

func (app *application) addNewSensorHandler(w http.ResponseWriter, r *http.Request) {
	var sensor sensor.Sensor
	if err := json.NewDecoder(r.Body).Decode(&sensor); err != nil {
		app.httpError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	id, err := app.store.AddNewSensor(sensor)
	if err != nil {
		app.httpError(w, http.StatusInternalServerError, err)
		return
	}

	json.NewEncoder(w).Encode(idResponse{
		ID: id,
	})
}

func (app *application) addNewWarehouseHandler(w http.ResponseWriter, r *http.Request) {
	var warehouse warehouse.Warehouse
	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		app.httpError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	id, err := app.store.AddNewWarehouse(warehouse)
	if err != nil {
		app.httpError(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Fprintf(w, "{\"id\": %d}", id)
}

func (app *application) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.httpError(w, http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	client := ws.NewClient(conn)
	ws.Mu.Lock()
	ws.Clients[client] = true
	ws.Mu.Unlock()

	app.logger.Info("Client Connected")

	for {
		_, _, err := conn.NextReader()
		if err != nil {
			ws.Mu.Lock()
			delete(ws.Clients, client)
			ws.Mu.Unlock()
			conn.Close()
			app.logger.Info("Client Disconnected")
			break
		}
	}
}
