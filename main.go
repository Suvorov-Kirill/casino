package main

import (
	_ "database/sql"
	"fmt"
	"html/template"
	_ "log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type PageData struct {
	Message string
}

// обработчик главной страницы
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

// обработчик кнопки "Играть"

func playHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	result := "Проигрыш 😢"
	if rand.Intn(2) == 0 {
		result = "Победа 🎉"
	}
	fmt.Fprintln(w, result)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/play", playHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
