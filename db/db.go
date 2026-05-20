package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error

	// Открываем (или создаём) файл базы данных
	DB, err = sql.Open("sqlite3", "casino.db")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Подключение к базе успешно!")

	sqlBytes, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal("Не могу прочитать schema.sql:", err)
	}

	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Не могу выполнить SQL из schema.sql:", err)
	}

	fmt.Println("Структура базы данных загружена.")
}

func Close() {
	if DB == nil {
		return
	}

	if err := DB.Close(); err != nil {
		log.Println("Ошибка при закрытии базы:", err)
	}
}
