package main

import "net/http"

func registerRoutes() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/play", playHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/admin/users", requireAdmin(adminUsersHandler))
	http.HandleFunc("/admin/users/edit", requireAdmin(adminEditUser))
	http.HandleFunc("/admin/users/delete", requireAdmin(adminDeleteUser))
	http.HandleFunc("/play/slots", playSlotsHandler)
	http.HandleFunc("/play/craps", playCrapsHandler)
}
