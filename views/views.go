package views

import (
	"html/template"
	"log"
	"net/http"
)

type SlotsPageData struct {
	Balance int
	Result  string
	Win     bool
	Message string
}

type CrapsPageData struct {
	Balance    int
	Point      int
	InProgress bool
	Message    string
	Win        bool
}

type RoulettePageData struct {
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

type BlackjackPageData struct {
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
func IndexHandler(w http.ResponseWriter, _ *http.Request) {
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

func RenderSlotsPage(w http.ResponseWriter, data SlotsPageData) {
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

func RenderCrapsPage(w http.ResponseWriter, data CrapsPageData) {
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

func RenderRoulettePage(w http.ResponseWriter, data RoulettePageData) {
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

func RenderBlackjackPage(w http.ResponseWriter, data BlackjackPageData) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/play_blackjack.html")
	if err != nil {
		http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
		log.Println("Template parse error:", err)
		return
	}
	tmpl.ExecuteTemplate(w, "base", data)
}
