package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Init() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/casino?sslmode=disable"
	}

	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к PostgreSQL:", err)
	}

	fmt.Println("Подключение к PostgreSQL успешно!")

	sqlBytes, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatal("Не могу прочитать schema.sql:", err)
	}

	for _, statement := range strings.Split(string(sqlBytes), ";") {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		if _, err := DB.Exec(statement); err != nil {
			log.Fatal("Не могу выполнить SQL из schema.sql:", err)
		}
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
