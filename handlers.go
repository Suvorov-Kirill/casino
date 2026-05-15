package main

import (
	"casino/db"
	"casino/games"
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
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		http.Error(w, "Шаблон главной страницы не найден", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
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

func getUserIDAndBalance(r *http.Request) (userID int, balance int, err error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return 0, 0, err
	}

	userID, err = strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, 0, err
	}

	err = db.DB.QueryRow("SELECT coins FROM users WHERE id = ?", userID).Scan(&balance)
	if err != nil {
		return 0, 0, err
	}

	return userID, balance, nil
}

func updateBalance(userID, newBalance int) error {
	_, err := db.DB.Exec("UPDATE users SET coins = ? WHERE id = ?", newBalance, userID)
	return err
}

func saveBet(userID, bet int, game string, win bool) error {
	_, err := db.DB.Exec(
		"INSERT INTO bets (user_id, amount, game, result) VALUES (?, ?, ?, ?)",
		userID, bet, game, win,
	)
	return err
}

func writeResponse(w http.ResponseWriter, msg string) {
	_, err := fmt.Fprintln(w, msg)
	if err != nil {
		http.Error(w, "Ошибка при отправке ответа клиенту", http.StatusInternalServerError)
	}
}

func playSlotsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/play_slots.html"))
		err = tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
		}
		return
	}

	userID, balance, err := getUserIDAndBalance(r)
	if err != nil {
		http.Error(w, "Ошибка получения пользователя", http.StatusBadRequest)
		return
	}

	bet, _ := strconv.Atoi(r.FormValue("bet"))
	if bet > balance {
		http.Error(w, "Недостаточно монет", http.StatusBadRequest)
		return
	}

	resultText, win := games.SlotsGameLogic()
	newBalance := balance - bet
	if win {
		newBalance += bet * 2
	}

	if err := updateBalance(userID, newBalance); err != nil {
		http.Error(w, "Ошибка обновления баланса", http.StatusInternalServerError)
		return
	}

	if err := saveBet(userID, bet, "Slots", win); err != nil {
		http.Error(w, "Ошибка записи ставки", http.StatusInternalServerError)
		return
	}

	writeResponse(w, fmt.Sprintf("Результат игры: %s. Ваш баланс: %d", resultText, newBalance))
}

func deductBet(userID, balance, bet int) error {
	if bet > balance {
		return errors.New("недостаточно монет")
	}
	_, err := db.DB.Exec("UPDATE users SET coins = ? WHERE id = ?", balance-bet, userID)
	return err
}

func playCrapsHandler(w http.ResponseWriter, r *http.Request) {
	userID, balance, err := getUserIDAndBalance(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	game, ok := games.CrapsGames[userID]
	if !ok {
		game = &games.CrapsGame{}
		games.CrapsGames[userID] = game
	}

	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/play_craps.html"))
		err = tmpl.ExecuteTemplate(w, "base", game)
		if err != nil {
			http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
			return
		}

	case http.MethodPost:
		bet, _ := strconv.Atoi(r.FormValue("bet"))

		if !game.InProgress {
			if err := deductBet(userID, balance, bet); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		d1, d2, sum := game.RollDice()

		var msg string
		if !game.InProgress {
			switch game.ProcessComeOutRoll(sum) {
			case "win":
				games.FinishCrapsGame(userID, bet, true)
				msg = fmt.Sprintf("Первый бросок: %d + %d = %d → ВЫИГРЫШ!", d1, d2, sum)
			case "lose":
				games.FinishCrapsGame(userID, bet, false)
				msg = fmt.Sprintf("Первый бросок: %d + %d = %d → ПРОИГРЫШ!", d1, d2, sum)
			case "point":
				msg = fmt.Sprintf("Point установлен: %d. Бросайте снова!", sum)
			}
		} else {
			switch game.ProcessPointRoll(sum) {
			case "win":
				games.FinishCrapsGame(userID, bet, true)
				msg = fmt.Sprintf("Попал в POINT! Победа!")
			case "lose":
				games.FinishCrapsGame(userID, bet, false)
				msg = "Выпало 7! Проигрыш"
			case "continue":
				msg = fmt.Sprintf("Выпало: %d + %d = %d. Point: %d. Бросайте снова!", d1, d2, sum, game.Point)
			}
		}

		writeResponse(w, msg)
	}
}
