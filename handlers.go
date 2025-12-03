package main

import (
	_ "database/sql"
	"fmt"
	"html/template"
	_ "log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/register.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Хэшируем пароль
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		// Добавляем пользователя
		_, err := db.Exec("INSERT INTO users (username, password, coins) VALUES (?, ?, ?)",
			username, string(hashedPassword), 100, // начальный баланс
		)

		if err != nil {
			fmt.Fprintln(w, "Ошибка: пользователь с таким именем уже существует")
			return
		}

		fmt.Fprintln(w, "Регистрация успешна! Теперь можете войти.")
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := cookie.Value

	var username string
	var coins int
	db.QueryRow("SELECT username, coins FROM users WHERE id = ?", userID).
		Scan(&username, &coins)

	fmt.Fprintf(w, "Привет, %s! Ваш баланс: %d монет.", username, coins)
}
