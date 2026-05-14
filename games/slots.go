package games

import (
	"fmt"
	"math/rand"
)

func SlotsGameLogic() (string, bool) {
	symbols := []string{"🍒", "🍋", "⭐", "🔔"}
	reel1 := symbols[rand.Intn(len(symbols))]
	reel2 := symbols[rand.Intn(len(symbols))]
	reel3 := symbols[rand.Intn(len(symbols))]

	win := false
	if reel1 == reel2 && reel2 == reel3 {
		win = true
	}

	result := fmt.Sprintf("%s | %s | %s", reel1, reel2, reel3)
	return result, win
}
