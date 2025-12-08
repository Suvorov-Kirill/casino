package main

import (
	"casino/db"
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
		}
		return
	}

	// POST
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Введите логин и пароль", http.StatusBadRequest)
		return
	}

	var storedHash string
	var id int
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&id, &storedHash)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		log.Println("DB error:", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// Устанавливаем cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: fmt.Sprint(id),
		Path:  "/",
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
			log.Println("Template error:", err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
		}
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Введите логин и пароль", http.StatusBadRequest)
			return
		}

		// Хэшируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			log.Println("Password hash error:", err)
			return
		}

		// Добавляем пользователя
		_, err = db.DB.Exec("INSERT INTO users (username, password, coins) VALUES (?, ?, ?)",
			username, string(hashedPassword), 100, // начальный баланс
		)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				http.Error(w, "Пользователь с таким именем уже существует", http.StatusConflict)
			} else {
				http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
				log.Println("DB insert error:", err)
			}
			return
		}

		_, err = fmt.Fprintln(w, "Регистрация успешна! Теперь можете войти.")
		if err != nil {
			log.Println("Ошибка вывода:", err)
		}
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID := cookie.Value

	// Получаем данные пользователя из базы
	var username string
	var coins int
	err = db.DB.QueryRow("SELECT username, coins FROM users WHERE id = ?", userID).Scan(&username, &coins)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		log.Println("DB error:", err)
		return
	}

	// Выводим страницу
	_, err = fmt.Fprintf(w, "Привет, %s! Ваш баланс: %d монет.", username, coins)
	if err != nil {
		log.Println("Ошибка вывода:", err)
	}
}
