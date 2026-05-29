package games

import (
	"math/rand"
	"time"
)

type RouletteResult struct {
	Number int
	Color  string
}

func SpinRoulette() RouletteResult {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := r.Intn(37)

	return RouletteResult{
		Number: number,
		Color:  rouletteColor(number),
	}
}

func rouletteColor(number int) string {
	if number == 0 {
		return "green"
	}

	redNumbers := map[int]struct{}{
		1: {}, 3: {}, 5: {}, 7: {}, 9: {}, 12: {}, 14: {}, 16: {}, 18: {},
		19: {}, 21: {}, 23: {}, 25: {}, 27: {}, 30: {}, 32: {}, 34: {}, 36: {},
	}

	if _, ok := redNumbers[number]; ok {
		return "red"
	}

	return "black"
}
