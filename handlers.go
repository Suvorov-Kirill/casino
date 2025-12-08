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

// обработчик главной страницы
func indexHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Шаблон главной страницы не найден", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения", http.StatusInternalServerError)
		return
	}
}

// обработчик кнопки "Играть"

func playHandler(w http.ResponseWriter, _ *http.Request) {
	rand.Seed(time.Now().UnixNano())
	result := "Проигрыш 😢"
	if rand.Intn(2) == 0 {
		result = "Победа 🎉"
	}
	_, err := fmt.Fprintln(w, result)
	if err != nil {
		http.Error(w, "Ошибка вывода ответа", http.StatusInternalServerError)
		return
	}
}
