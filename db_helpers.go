package main

import (
	"casino/db"
	"errors"
	"net/http"
	"strconv"
)

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

func saveBet(userID, bet int, game string, win bool) error {
	_, err := db.DB.Exec(
		"INSERT INTO bets (user_id, amount, game, result) VALUES ($1, $2, $3, $4)",
		userID, bet, game, win,
	)
	return err
}

func deductBet(userID, balance, bet int) error {
	if bet > balance {
		return errors.New("недостаточно монет")
	}
	_, err := db.DB.Exec("UPDATE users SET coins = $1 WHERE id = $2", balance-bet, userID)
	return err
}
