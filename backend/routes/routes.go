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

    // Routes Admin
    	http.HandleFunc("/admin", middleware.RequireRole(database, []string{"admin"}, func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "/usr/share/nginx/html/pages/admin.html")
    }))

	// Routes protégées (nécessitent JWT)
    http.HandleFunc("/api/forms", middleware.RequireJWT(database, handlers.GetAllForms(database)))
    http.HandleFunc("/api/form", middleware.RequireJWT(database, handlers.GetForm(database)))
    http.HandleFunc("/api/form/add", middleware.RequireJWT(database, handlers.AjoutForm(database)))
    http.HandleFunc("/api/users", middleware.RequireJWT(database, handlers.GetAllUsersHandler(database)))

    // Routes publiques (pas de JWT)
    http.HandleFunc("/api/user/add", handlers.AjoutUser(database))
    http.HandleFunc("/api/login", handlers.LoginHandler(database))
}