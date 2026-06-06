package main

import (
	"casino/app"
	"casino/db"
	"casino/routes"
	"log"
	"net/http"
)

func main() {
	database := db.Init()
	defer database.Close()

	app := &app.CasinoApp{
		DB: database,
	}
	routes.RegisterRoutes(app)

	println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Сервер не запустился:", err)
	}

}
