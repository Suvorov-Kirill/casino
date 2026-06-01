package main

import (
	"casino/db"
	"casino/games"
	"errors"
	"fmt"
	"html/template"
	"log"
	_ "log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type slotsPageData struct {
	Balance int
	Result  string
	Win     bool
	Message string
}

type crapsPageData struct {
	Balance    int
	Point      int
	InProgress bool
	Message    string
	Win        bool
}

type roulettePageData struct {
	Balance     int
	Message     string
	Win         bool
	HasResult   bool
	ResultText  string
	ResultClass string
	BetKind     string
	BetChoice   string
	BetAmount   int
}
type CardDisplay struct {
	Rank  string
	Suit  string
	Color string
}

type blackjackPageData struct {
	Balance     int
	Message     string
	Status      string // "playing", "won", "lost", "push", "blackjack"
	PlayerCards []CardDisplay
	PlayerScore int
	DealerCards []CardDisplay
	DealerScore int
	InProgress  bool
	Bet         int
}

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

	err = db.DB.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		return 0, 0, err
	}

	return userID, balance, nil
}

func updateBalance(userID, newBalance int) error {
	_, err := db.DB.Exec("UPDATE users SET coins = $1 WHERE id = $2", newBalance, userID)
	return err
}

func renderSlotsPage(w http.ResponseWriter, data slotsPageData) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/play_slots.html")
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		log.Println("Template parse error:", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		log.Println("Template execute error:", err)
	}
}

func renderCrapsPage(w http.ResponseWriter, data crapsPageData) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/play_craps.html")
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		log.Println("Template parse error:", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		log.Println("Template execute error:", err)
	}
}

func renderRoulettePage(w http.ResponseWriter, data roulettePageData) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/roulette.html")
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		log.Println("Template parse error:", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		log.Println("Template execute error:", err)
	}
}

func saveBet(userID, bet int, game string, win bool) error {
	_, err := db.DB.Exec(
		"INSERT INTO bets (user_id, amount, game, result) VALUES ($1, $2, $3, $4)",
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
	userID, balance, err := getUserIDAndBalance(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := slotsPageData{Balance: balance}

	if r.Method != http.MethodPost {
		renderSlotsPage(w, data)
		return
	}

	bet, parseErr := strconv.Atoi(r.FormValue("bet"))
	if parseErr != nil || bet <= 0 {
		data.Message = "Введите корректную ставку"
		renderSlotsPage(w, data)
		return
	}

	if bet > balance {
		data.Message = "Недостаточно монет"
		renderSlotsPage(w, data)
		return
	}

	resultText, win := games.SlotsGameLogic()
	newBalance := balance - bet
	if win {
		newBalance += bet * 2
	}

	if err := updateBalance(userID, newBalance); err != nil {
		data.Message = "Ошибка обновления баланса"
		renderSlotsPage(w, data)
		return
	}

	if err := saveBet(userID, bet, "Slots", win); err != nil {
		data.Message = "Ошибка записи ставки"
		renderSlotsPage(w, data)
		return
	}

	_, currentBalance, err := getUserIDAndBalance(r)
	if err == nil {
		data.Balance = currentBalance
	} else {
		data.Balance = newBalance
	}
	data.Result = fmt.Sprintf("%s", resultText)
	data.Win = win
	if win {
		data.Message = fmt.Sprintf("Выигрыш! Ставка %d. Баланс обновлен.", bet)
	} else {
		data.Message = fmt.Sprintf("Проигрыш. Ставка %d.", bet)
	}

	renderSlotsPage(w, data)
}

func deductBet(userID, balance, bet int) error {
	if bet > balance {
		return errors.New("недостаточно монет")
	}
	_, err := db.DB.Exec("UPDATE users SET coins = $1 WHERE id = $2", balance-bet, userID)
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

	data := crapsPageData{
		Balance:    balance,
		Point:      game.Point,
		InProgress: game.InProgress,
	}

	switch r.Method {
	case http.MethodGet:
		renderCrapsPage(w, data)

	case http.MethodPost:
		if !game.InProgress {
			bet, parseErr := strconv.Atoi(r.FormValue("bet"))
			if parseErr != nil || bet <= 0 {
				data.Message = "Введите корректную ставку"
				renderCrapsPage(w, data)
				return
			}

			if err := deductBet(userID, balance, bet); err != nil {
				data.Message = err.Error()
				renderCrapsPage(w, data)
				return
			}

			game.Bet = bet
			data.Balance = balance - bet
		} else {
			data.Balance = balance
		}

		d1, d2, sum := game.RollDice()

		var msg string
		if !game.InProgress {
			switch game.ProcessComeOutRoll(sum) {
			case "win":
				games.FinishCrapsGame(userID, true)
				msg = fmt.Sprintf("Первый бросок: %d + %d = %d → ВЫИГРЫШ!", d1, d2, sum)
			case "lose":
				games.FinishCrapsGame(userID, false)
				msg = fmt.Sprintf("Первый бросок: %d + %d = %d → ПРОИГРЫШ!", d1, d2, sum)
			case "point":
				msg = fmt.Sprintf("Point установлен: %d. Бросайте снова!", sum)
				data.Point = game.Point
				data.InProgress = true
			}
		} else {
			switch game.ProcessPointRoll(sum) {
			case "win":
				games.FinishCrapsGame(userID, true)
				msg = "Попал в POINT! Победа!"
			case "lose":
				games.FinishCrapsGame(userID, false)
				msg = "Выпало 7! Проигрыш"
			case "continue":
				msg = fmt.Sprintf("Выпало: %d + %d = %d. Point: %d. Бросайте снова!", d1, d2, sum, game.Point)
				data.Point = game.Point
				data.InProgress = true
			}
		}

		_, currentBalance, err := getUserIDAndBalance(r)
		if err == nil {
			data.Balance = currentBalance
		}
		data.Message = msg
		data.Point = game.Point
		data.InProgress = game.InProgress
		renderCrapsPage(w, data)
	}
}

func renderBlackjackPage(w http.ResponseWriter, data blackjackPageData) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/play_blackjack.html")
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		log.Println("Template parse error:", err)
		return
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func parseHand(handStr string) []string {
	if handStr == "" {
		return []string{}
	}
	return strings.Split(handStr, ",")
}

func makeDisplayHand(cards []string, hideDealerSecondCard bool) []CardDisplay {
	display := make([]CardDisplay, len(cards))
	for i, card := range cards {
		if i == 1 && hideDealerSecondCard {
			display[i] = CardDisplay{Rank: "?", Suit: "?", Color: "hidden"}
		} else {
			rank, suit, color := games.FormatCard(card)
			display[i] = CardDisplay{Rank: rank, Suit: suit, Color: color}
		}
	}
	return display
}

func playBlackjackHandler(w http.ResponseWriter, r *http.Request) {
	userID, balance, err := getUserIDAndBalance(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := blackjackPageData{Balance: balance}

	// 1. Проверяем, есть ли активная игра в БД
	var betAmount int
	var playerHandStr, dealerHandStr, deckStr, status string
	err = db.DB.QueryRow("SELECT bet_amount, player_hand, dealer_hand, deck, status FROM blackjack_games WHERE user_id = $1", userID).
		Scan(&betAmount, &playerHandStr, &dealerHandStr, &deckStr, &status)

	hasActiveGame := err == nil

	if r.Method == http.MethodGet {
		if hasActiveGame {
			playerHand := parseHand(playerHandStr)
			dealerHand := parseHand(dealerHandStr)
			data.InProgress = true
			data.Bet = betAmount
			data.PlayerScore = games.CalculateScore(playerHand)
			data.PlayerCards = makeDisplayHand(playerHand, false)
			// Скрываем вторую карту дилера, если игра еще идет
			data.DealerCards = makeDisplayHand(dealerHand, true)
		}
		renderBlackjackPage(w, data)
		return
	}

	// Обработка POST-запросов
	action := r.FormValue("action") // "start", "hit", "stand"

	if action == "start" && !hasActiveGame {
		bet, parseErr := strconv.Atoi(r.FormValue("bet"))
		if parseErr != nil || bet <= 0 {
			data.Message = "Введите корректную ставку"
			renderBlackjackPage(w, data)
			return
		}
		if bet > balance {
			data.Message = "Недостаточно монет"
			renderBlackjackPage(w, data)
			return
		}

		// Снимаем ставку
		newBalance := balance - bet
		updateBalance(userID, newBalance)
		data.Balance = newBalance

		playerHand, dealerHand, deck := games.InitialDeal()
		playerScore := games.CalculateScore(playerHand)

		// Проверка на натуральный блэкджек со старта
		gameStatus := "playing"
		if playerScore == 21 {
			gameStatus = "blackjack"
			// Выплачиваем сразу 3 к 2
			winAmount := bet + (bet * 3 / 2)
			updateBalance(userID, newBalance+winAmount)
			data.Balance = newBalance + winAmount
			saveBet(userID, bet, "Blackjack", true)
		}

		// Сохраняем в БД
		_, err = db.DB.Exec("INSERT INTO blackjack_games (user_id, bet_amount, player_hand, dealer_hand, deck, status) VALUES ($1, $2, $3, $4, $5, $6)",
			userID, bet, strings.Join(playerHand, ","), strings.Join(dealerHand, ","), strings.Join(deck, ","), gameStatus)

		if gameStatus == "playing" {
			data.InProgress = true
			data.Message = "Игра начата!"
			data.DealerCards = makeDisplayHand(dealerHand, true) // Скрываем
		} else {
			data.Status = gameStatus
			data.Message = "Блэкджек! Вы выиграли!"
			data.DealerCards = makeDisplayHand(dealerHand, false) // Открываем
			data.DealerScore = games.CalculateScore(dealerHand)
			db.DB.Exec("DELETE FROM blackjack_games WHERE user_id = $1", userID) // Удаляем завершенную игру
		}

		data.PlayerCards = makeDisplayHand(playerHand, false)
		data.PlayerScore = playerScore
		data.Bet = bet
		renderBlackjackPage(w, data)
		return
	}

	if hasActiveGame && data.Status != "playing" {
		playerHand := parseHand(playerHandStr)
		dealerHand := parseHand(dealerHandStr)
		deck := parseHand(deckStr)

		if action == "hit" {
			playerHand, deck = games.DrawCard(playerHand, deck)
			playerScore := games.CalculateScore(playerHand)

			if playerScore > 21 {
				status = "lost"
				saveBet(userID, betAmount, "Blackjack", false)
				db.DB.Exec("DELETE FROM blackjack_games WHERE user_id = $1", userID)
				data.Message = "Перебор! Вы проиграли."
			} else {
				// Обновляем состояние
				db.DB.Exec("UPDATE blackjack_games SET player_hand = $1, deck = $2 WHERE user_id = $3",
					strings.Join(playerHand, ","), strings.Join(deck, ","), userID)
				data.InProgress = true
			}

			data.PlayerCards = makeDisplayHand(playerHand, false)
			data.PlayerScore = playerScore
			data.DealerCards = makeDisplayHand(dealerHand, status != "lost") // Если не проиграли, карта скрыта
			data.Status = status
			data.Bet = betAmount

		} else if action == "stand" {
			// Ход дилера
			playerScore := games.CalculateScore(playerHand)
			dealerScore := games.CalculateScore(dealerHand)

			// Дилер берет до 17
			for dealerScore < 17 {
				dealerHand, deck = games.DrawCard(dealerHand, deck)
				dealerScore = games.CalculateScore(dealerHand)
			}

			win := false
			if dealerScore > 21 {
				status = "won"
				win = true
				data.Message = "Перебор у дилера! Вы выиграли."
			} else if playerScore > dealerScore {
				status = "won"
				win = true
				data.Message = "Вы победили!"
			} else if dealerScore > playerScore {
				status = "lost"
				data.Message = "Дилер победил."
			} else {
				status = "push"
				data.Message = "Ничья. Ставка возвращена."
			}

			if win {
				updateBalance(userID, balance+(betAmount*2))
				data.Balance = balance + (betAmount * 2)
				saveBet(userID, betAmount, "Blackjack", true)
			} else if status == "push" {
				updateBalance(userID, balance+betAmount)
				data.Balance = balance + betAmount
			} else {
				saveBet(userID, betAmount, "Blackjack", false)
			}

			db.DB.Exec("DELETE FROM blackjack_games WHERE user_id = $1", userID)

			data.PlayerCards = makeDisplayHand(playerHand, false)
			data.PlayerScore = playerScore
			data.DealerCards = makeDisplayHand(dealerHand, false) // Открываем все карты
			data.DealerScore = dealerScore
			data.Status = status
			data.Bet = betAmount
		}

		renderBlackjackPage(w, data)
		return
	}

	// По умолчанию, если ничего не подошло
	http.Redirect(w, r, "/play/blackjack", http.StatusSeeOther)
}
func playRouletteHandler(w http.ResponseWriter, r *http.Request) {
	userID, balance, err := getUserIDAndBalance(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := roulettePageData{Balance: balance}

	if r.Method != http.MethodPost {
		renderRoulettePage(w, data)
		return
	}

	bet, parseErr := strconv.Atoi(r.FormValue("bet"))
	if parseErr != nil || bet <= 0 {
		data.Message = "Введите корректную ставку"
		renderRoulettePage(w, data)
		return
	}

	if bet > balance {
		data.Message = "Недостаточно монет"
		renderRoulettePage(w, data)
		return
	}

	betKind := strings.ToLower(strings.TrimSpace(r.FormValue("bet_kind")))
	result := games.SpinRoulette()
	data.HasResult = true
	data.ResultText = fmt.Sprintf("%d / %s", result.Number, strings.ToUpper(result.Color))
	data.ResultClass = result.Color
	data.BetAmount = bet
	data.BetKind = betKind

	newBalance := balance - bet
	win := false
	message := ""

	switch betKind {
	case "number":
		betChoice, choiceErr := strconv.Atoi(r.FormValue("number"))
		if choiceErr != nil || betChoice < 0 || betChoice > 36 {
			data.Message = "Введите число от 0 до 36"
			renderRoulettePage(w, data)
			return
		}

		data.BetChoice = strconv.Itoa(betChoice)
		if betChoice == result.Number {
			win = true
			newBalance += bet * 36
			message = fmt.Sprintf("Число %d сыграло. Выигрыш x36.", betChoice)
		} else {
			message = fmt.Sprintf("Выпало %d. Ставка на %d не зашла.", result.Number, betChoice)
		}
	case "color":
		betChoice := strings.ToLower(strings.TrimSpace(r.FormValue("color")))
		if betChoice != "red" && betChoice != "black" {
			data.Message = "Выберите красное или чёрное"
			renderRoulettePage(w, data)
			return
		}

		data.BetChoice = betChoice
		if betChoice == result.Color {
			win = true
			newBalance += bet * 2
			message = fmt.Sprintf("%s сыграло. Выигрыш x2.", strings.ToUpper(betChoice))
		} else {
			message = fmt.Sprintf("Выпало %s. Ставка на %s не зашла.", result.Color, betChoice)
		}
	default:
		data.Message = "Выберите тип ставки"
		renderRoulettePage(w, data)
		return
	}

	if err := updateBalance(userID, newBalance); err != nil {
		data.Message = "Ошибка обновления баланса"
		renderRoulettePage(w, data)
		return
	}

	if err := saveBet(userID, bet, "Roulette", win); err != nil {
		data.Message = "Ошибка записи ставки"
		renderRoulettePage(w, data)
		return
	}

	_, currentBalance, err := getUserIDAndBalance(r)
	if err == nil {
		data.Balance = currentBalance
	} else {
		data.Balance = newBalance
	}

	data.Win = win
	data.Message = message
	renderRoulettePage(w, data)
}
