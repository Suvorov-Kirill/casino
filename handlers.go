package main

import (
	"casino/db"
	"casino/games"
	"database/sql"
	_ "database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	_ "log"
	"math/rand"
	"net/http"
	"strconv"
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

func playSlotsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		// Проверка куки
		userCookie, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := strconv.Atoi(userCookie.Value)
		if err != nil {
			http.Error(w, "Неверный формат user_id", http.StatusBadRequest)
			return
		}

		// Получение ставки
		bet, err := strconv.Atoi(r.FormValue("bet"))
		if err != nil || bet <= 0 {
			http.Error(w, "Неверная ставка", http.StatusBadRequest)
			return
		}

		var coins int
		err = db.DB.QueryRow("SELECT coins FROM users WHERE id = ?", userID).Scan(&coins)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Ошибка базы данных (получение баланса)", http.StatusInternalServerError)
			log.Println("DB error:", err)
			return
		}

		if bet > coins {
			http.Error(w, "Недостаточно монет", http.StatusBadRequest)
			return
		}

		resultText, win := games.SlotsGameLogic()

		// Расчёт нового баланса
		newBalance := coins - bet
		if win {
			newBalance += bet * 2
		}

		// Обновление баланса
		_, err = db.DB.Exec("UPDATE users SET coins = ? WHERE id = ?", newBalance, userID)
		if err != nil {
			http.Error(w, "Ошибка обновления баланса", http.StatusInternalServerError)
			log.Println("DB error:", err)
			return
		}

		// Сохранение ставки
		_, err = db.DB.Exec(
			"INSERT INTO bets (user_id, amount, game, result) VALUES (?, ?, ?, ?)",
			userID, bet, "Slots", win,
		)
		if err != nil {
			http.Error(w, "Ошибка записи ставки", http.StatusInternalServerError)
			log.Println("DB error:", err)
			return
		}

		_, err = fmt.Fprintf(w, "Результат игры: %s. Ваш баланс: %d", resultText, newBalance)
		if err != nil {
			log.Println("Ошибка вывода в ResponseWriter:", err)
		}

		return
	}

	tmpl, err := template.ParseFiles("templates/play_slots.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Ошибка выполнения шаблона:", err)
	}
}
