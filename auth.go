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

type authPageData struct {
	Username  string
	Message   string
	IsSuccess bool
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderLoginPage := func(data authPageData) {
		tmpl, err := template.ParseFiles("templates/layout.html", "templates/login.html")
		if err != nil {
			http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
		}
	}

	if r.Method == http.MethodGet {
		renderLoginPage(authPageData{})
		return
	}

	// POST
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		renderLoginPage(authPageData{
			Username: username,
			Message:  "Введите логин и пароль",
		})
		return
	}

	var storedHash string
	var id int
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&id, &storedHash)
	if errors.Is(err, sql.ErrNoRows) {
		renderLoginPage(authPageData{
			Username: username,
			Message:  "Неверный логин или пароль",
		})
		return
	} else if err != nil {
		log.Println("DB error:", err)
		renderLoginPage(authPageData{
			Username: username,
			Message:  "Ошибка сервера",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		renderLoginPage(authPageData{
			Username: username,
			Message:  "Неверный логин или пароль",
		})
		return
	}

	// Устанавливаем cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: fmt.Sprint(id),
		Path:  "/",
	})

	renderLoginPage(authPageData{
		Username:  username,
		Message:   "Login successful",
		IsSuccess: true,
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	renderRegisterPage := func(data authPageData) {
		tmpl, err := template.ParseFiles("templates/layout.html", "templates/register.html")
		if err != nil {
			http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
			log.Println("Template error:", err)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
		}
	}

	if r.Method == http.MethodGet {
		renderRegisterPage(authPageData{})
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			renderRegisterPage(authPageData{
				Username: username,
				Message:  "Введите логин и пароль",
			})
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
				renderRegisterPage(authPageData{
					Username: username,
					Message:  "Пользователь с таким именем уже существует",
				})
			} else {
				log.Println("DB insert error:", err)
				renderRegisterPage(authPageData{
					Username: username,
					Message:  "Ошибка сервера",
				})
			}
			return
		}

		renderRegisterPage(authPageData{
			Username:  username,
			Message:   "Registration successful",
			IsSuccess: true,
		})
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

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/profile.html")
	if err != nil {
		http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := struct {
		Username string
		Coins    int
	}{
		Username: username,
		Coins:    coins,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
		log.Println("Template execute error:", err)
	}
}
