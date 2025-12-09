package main

import (
	"casino/db"
	"database/sql"
	_ "database/sql"
	"fmt"
	"log"
	_ "log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var err error

func main() {
	// Открываем (или создаём) файл базы данных
	db.DB, err = sql.Open("sqlite3", "casino.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Println("Ошибка при закрытии базы:", err)
		}
	}(db.DB)

	// Проверяем подключение
	err = db.DB.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе:", err)
	}

	fmt.Println("Подключение к базе успешно!")

	sqlBytes, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal("Не могу прочитать schema.sql:", err)
	}

	_, err = db.DB.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Не могу выполнить SQL из schema.sql:", err)
	}

	fmt.Println("Структура базы данных загружена.")

	fmt.Println("Таблица 'users' готова.")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/play", playHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/admin/users", requireAdmin(adminUsersHandler))
	http.HandleFunc("/admin/users/edit", requireAdmin(adminEditUser))
	http.HandleFunc("/admin/users/delete", requireAdmin(adminDeleteUser))
	println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Сервер не запустился:", err)
	}

}
