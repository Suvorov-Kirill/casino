package handlers

import (
	"casino/app"
	"casino/models"
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
)

func RequireAdmin(app *app.CasinoApp, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Получаем cookie с user_id
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Получаем роль пользователя из базы
		var userRole string
		err = app.DB.QueryRow("SELECT user_role FROM users WHERE id = $1::int", cookie.Value).Scan(&userRole)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Проверка,что пользователь админ
		if userRole != "admin" {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func AdminUsersHandler(app *app.CasinoApp, w http.ResponseWriter, _ *http.Request) {
	rows, err := app.DB.Query("SELECT id, username, coins, user_role FROM users")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("rows.Close() error: %v", cerr)
		}
	}()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Email, &u.Balance, &u.Role)
		if err != nil {
			http.Error(w, "Database scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	// Проверяем, есть ли ошибки после завершения цикла
	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Disable caching for admin pages so browsers always fetch fresh data
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/admin_users.html")
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", users)
	if err != nil {
		http.Error(w, "Template execute error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func AdminBetsHandler(app *app.CasinoApp, w http.ResponseWriter, _ *http.Request) {
	// Prevent browser caching so admin sees latest bets
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	rows, err := app.DB.Query("SELECT b.id, b.user_id, u.username, b.amount, b.game, b.result, b.created_at FROM bets b LEFT JOIN users u ON u.id = b.user_id ORDER BY b.created_at DESC LIMIT 15")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("rows.Close() error: %v", cerr)
		}
	}()

	var bets []models.Bet
	for rows.Next() {
		var b models.Bet
		err := rows.Scan(&b.ID, &b.UserID, &b.Username, &b.Amount, &b.Game, &b.Result, &b.CreatedAt)
		if err != nil {
			http.Error(w, "Database scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		bets = append(bets, b)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/admin_bets.html")
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", bets)
	if err != nil {
		http.Error(w, "Template execute error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func AdminDeleteUser(app *app.CasinoApp, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := app.DB.Exec("DELETE FROM users WHERE id = $1::int", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func AdminEditUser(app *app.CasinoApp, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if r.Method == http.MethodPost {
		coins := r.FormValue("coins")
		userRole := r.FormValue("user_role")

		_, err := app.DB.Exec("UPDATE users SET coins = $1, user_role = $2 WHERE id = $3::int", coins, userRole, id)
		if err != nil {
			http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	var u models.User
	err := app.DB.QueryRow("SELECT id, username, coins, user_role FROM users WHERE id = $1::int", id).Scan(
		&u.ID, &u.Email, &u.Balance, &u.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/admin_edit_user.html")
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", u)
	if err != nil {
		http.Error(w, "Template execute error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
