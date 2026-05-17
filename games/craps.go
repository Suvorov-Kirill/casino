package games

import (
	"casino/db"
	"math/rand"
)

type CrapsGame struct {
	Point      int
	InProgress bool
	Bet        int
}

var CrapsGames = map[int]*CrapsGame{}

func (g *CrapsGame) RollDice() (int, int, int) {
	d1 := rand.Intn(6) + 1
	d2 := rand.Intn(6) + 1
	return d1, d2, d1 + d2
}

func (g *CrapsGame) ProcessComeOutRoll(sum int) string {
	switch sum {
	case 7, 11:
		return "win"
	case 2, 3, 12:
		return "lose"
	default:
		g.Point = sum
		g.InProgress = true
		return "point"
	}
}

func (g *CrapsGame) ProcessPointRoll(sum int) string {
	if sum == g.Point {
		return "win"
	}
	if sum == 7 {
		return "lose"
	}
	return "continue"
}

func FinishCrapsGame(userID int, win bool) {
	game := CrapsGames[userID]
	game.InProgress = false
	game.Point = 0
	bet := game.Bet
	game.Bet = 0

	if win {
		db.DB.Exec("UPDATE users SET coins = coins + ? WHERE id = ?", bet*2, userID)
	}

	db.DB.Exec(
		"INSERT INTO bets (user_id, amount, game, result) VALUES (?, ?, 'Craps', ?)",
		userID, bet, win)
}
