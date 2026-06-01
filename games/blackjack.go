package games

import (
	"math/rand"
	"strconv"
	"time"
)

// Масти: H(Hearts/Черви-Красн), D(Diamonds/Бубны-Красн), C(Clubs/Трефы-Черн), S(Spades/Пики-Черн)
// Ранги: 2-9, T(10), J, Q, K, A

func generateDeck() []string {
	suits := []string{"H", "D", "C", "S"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
	deck := make([]string, 0, 52)
	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, suit+rank)
		}
	}
	// Перемешиваем
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func CalculateScore(hand []string) int {
	score := 0
	aces := 0
	for _, card := range hand {
		if len(card) < 2 {
			continue
		}
		rank := card[1:]
		if rank == "A" {
			aces++
			score += 11
		} else if rank == "K" || rank == "Q" || rank == "J" || rank == "T" {
			score += 10
		} else {
			val, _ := strconv.Atoi(rank)
			score += val
		}
	}
	// Корректируем тузы, если перебор
	for score > 21 && aces > 0 {
		score -= 10
		aces--
	}
	return score
}

// FormatCard возвращает масть (символ), ранг и цвет для отображения в шаблоне
func FormatCard(cardStr string) (string, string, string) {
	if len(cardStr) < 2 {
		return "", "", ""
	}
	suitChar := cardStr[0:1]
	rankChar := cardStr[1:]

	suitSymbol := ""
	color := "black" // Для нашего CSS: red или black

	switch suitChar {
	case "H":
		suitSymbol = "♥"
		color = "red"
	case "D":
		suitSymbol = "♦"
		color = "red"
	case "C":
		suitSymbol = "♣"
		color = "black"
	case "S":
		suitSymbol = "♠"
		color = "black"
	}

	if rankChar == "T" {
		rankChar = "10"
	}

	return rankChar, suitSymbol, color
}

// InitialDeal возвращает начальную руку игрока, дилера и оставшуюся колоду
func InitialDeal() ([]string, []string, []string) {
	deck := generateDeck()
	player := []string{deck[0], deck[2]}
	dealer := []string{deck[1], deck[3]}
	return player, dealer, deck[4:]
}

func DrawCard(hand []string, deck []string) ([]string, []string) {
	if len(deck) == 0 {
		return hand, deck
	}
	return append(hand, deck[0]), deck[1:]
}
