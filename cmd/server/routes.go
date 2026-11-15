package main

import "net/http"

func (app *application) setRoutes() {
	// app.mux.HandleFunc("POST /records/new", app.newRecordHandler)
	app.mux.HandleFunc("POST /users/login", app.logInHandler)
	app.mux.Handle("POST /users/create", app.ValidateAdmin(http.HandlerFunc(app.signNewAccountHandler)))
	app.mux.Handle("POST /sensors/add", app.ValidateAdmin(http.HandlerFunc(app.addNewSensorHandler)))
	app.mux.Handle("POST /warehouses/add", app.ValidateAdmin(http.HandlerFunc(app.addNewWarehouseHandler)))
	app.mux.Handle("GET /ws", app.validateJWT(http.HandlerFunc(app.wsHandler)))
}
