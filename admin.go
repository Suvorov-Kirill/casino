package main

import (
	"casino/db"
	"casino/models"
	"database/sql"
	_ "database/sql"
	"errors"
	"html/template"
	"log"
	_ "log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Получаем cookie с user_id
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Получаем роль пользователя из базы
		var userRole string
		err = db.DB.QueryRow("SELECT user_role FROM users WHERE id = ?", cookie.Value).Scan(&userRole)
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

func adminUsersHandler(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.DB.Query("SELECT id, username, coins, user_role FROM users")
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

	tmpl, err := template.ParseFiles("templates/admin_users.html")
	tmpl, err = template.ParseFiles("templates/layout.html", "templates/admin_users.html")
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

func adminDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func adminEditUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if r.Method == http.MethodPost {
		coins := r.FormValue("coins")
		userRole := r.FormValue("user_role")

		_, err := db.DB.Exec("UPDATE users SET coins = ?, user_role = ? WHERE id = ?", coins, userRole, id)
		if err != nil {
			http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}

	var u models.User
	err := db.DB.QueryRow("SELECT id, username, coins, user_role FROM users WHERE id = ?", id).Scan(
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
