package main

import (
	"casino/db"
	"log"
	"net/http"
)

func main() {
	db.Init()
	defer db.Close()

	registerRoutes()

	println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Сервер не запустился:", err)
	}

}
