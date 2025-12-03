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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedHash string
	var id int

	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).
		Scan(&id, &storedHash)

	if err != nil {
		fmt.Fprintln(w, "Неверный логин или пароль")
		return
	}

	// Сравниваем пароль
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		fmt.Fprintln(w, "Неверный логин или пароль")
		return
	}

	// Устанавливаем cookie (сессию)
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: fmt.Sprint(id),
		Path:  "/",
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
