package db

import (
	"casino/app"
	"errors"
	"net/http"
	"strconv"
)

func GetUserIDAndBalance(app *app.CasinoApp, r *http.Request) (userID int, balance int, err error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return 0, 0, err
	}

	userID, err = strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, 0, err
	}

	err = app.DB.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		return 0, 0, err
	}

	return userID, balance, nil
}

func UpdateBalance(app *app.CasinoApp, userID, newBalance int) error {
	_, err := app.DB.Exec("UPDATE users SET coins = $1 WHERE id = $2", newBalance, userID)
	return err
}

func SaveBet(app *app.CasinoApp, userID, bet int, game string, win bool) error {
	_, err := app.DB.Exec(
		"INSERT INTO bets (user_id, amount, game, result) VALUES ($1, $2, $3, $4)",
		userID, bet, game, win,
	)
	return err
}

func DeductBet(app *app.CasinoApp, userID, balance, bet int) error {
	if bet > balance {
		return errors.New("недостаточно монет")
	}
	_, err := app.DB.Exec("UPDATE users SET coins = $1 WHERE id = $2", balance-bet, userID)
	return err
}
