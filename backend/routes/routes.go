package routes

import (
    "backend/handlers"
    "backend/middleware"
    "database/sql"
    "net/http"
)

func SetupRoutes(database *sql.DB) {
	// Routes RH et admin
	http.HandleFunc("/rh", middleware.RequireRole(database, []string{"rh", "admin"}, func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "/usr/share/nginx/html/pages/rh.html")
    }))
    http.HandleFunc("/api/forms", middleware.RequireRole(database, []string{"rh", "admin"}, handlers.GetAllForms(database)))
    http.HandleFunc("/api/form", middleware.RequireRole(database, []string{"rh", "admin"}, handlers.GetForm(database)))
    http.HandleFunc("/api/form/add", middleware.RequireRole(database, []string{"rh", "admin"}, handlers.AjoutForm(database)))
    http.HandleFunc("/api/form/delete", middleware.RequireRole(database, []string{"rh", "admin"}, handlers.SuppForm(database)))
    http.HandleFunc("/api/form/modify", middleware.RequireRole(database, []string{"rh", "admin"}, handlers.ModifForm(database)))

    // Routes Admin
    http.HandleFunc("/api/users", middleware.RequireRole(database, []string{"admin"}, handlers.GetAllUsersHandler(database)))


    // Routes publiques (pas de JWT)
    http.HandleFunc("/api/user/add", handlers.AjoutUser(database))
    http.HandleFunc("/api/login", handlers.LoginHandler(database))
}