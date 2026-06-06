package routes

import (
	"casino/app"
	"casino/handlers"
	"casino/views"
	"net/http"
)

func RegisterRoutes(app *app.CasinoApp) *http.ServeMux {
	mux := http.NewServeMux()
	http.HandleFunc("/", views.IndexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(app, w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(app, w, r)
	})
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProfileHandler(app, w, r)
	})
	http.HandleFunc("/admin/users", handlers.RequireAdmin(app, func(w http.ResponseWriter, r *http.Request) { handlers.AdminUsersHandler(app, w, r) }))
	http.HandleFunc("/admin/bets", handlers.RequireAdmin(app, func(w http.ResponseWriter, r *http.Request) { handlers.AdminBetsHandler(app, w, r) }))
	http.HandleFunc("/admin/users/edit", handlers.RequireAdmin(app, func(w http.ResponseWriter, r *http.Request) { handlers.AdminEditUser(app, w, r) }))
	http.HandleFunc("/admin/users/delete", handlers.RequireAdmin(app, func(w http.ResponseWriter, r *http.Request) { handlers.AdminDeleteUser(app, w, r) }))
	http.HandleFunc("/play/slots", func(w http.ResponseWriter, r *http.Request) {
		handlers.PlaySlotsHandler(app, w, r)
	})
	http.HandleFunc("/play/craps", func(w http.ResponseWriter, r *http.Request) {
		handlers.PlayCrapsHandler(app, w, r)
	})
	http.HandleFunc("/play/roulette", func(w http.ResponseWriter, r *http.Request) {
		handlers.PlayRouletteHandler(app, w, r)
	})
	http.HandleFunc("/play/blackjack", func(w http.ResponseWriter, r *http.Request) {
		handlers.PlayBlackjackHandler(app, w, r)
	})

	return mux
}
