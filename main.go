package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	"html/template"
	"log"
	_ "log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var err error

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
	// Открываем (или создаём) файл базы данных
	db, err = sql.Open("sqlite3", "casino.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе:", err)
	}

	fmt.Println("✅ Подключение к базе успешно!")

	// Создаём таблицу пользователей (если нет)
	createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE,
        password TEXT,
        coins INTEGER
    );`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}

	fmt.Println("✅ Таблица 'users' готова.")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/play", playHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/profile", profileHandler)
	println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
