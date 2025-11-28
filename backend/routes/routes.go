package routes

import (
    "backend/handlers"
    "database/sql"
    "net/http"
)

func SetupRoutes(database *sql.DB) {
    // Routes protégées (nécessitent JWT)
    http.HandleFunc("/api/forms", handlers.RequireJWT(database, handlers.GetAllForms(database)))
    http.HandleFunc("/api/form", handlers.RequireJWT(database, handlers.GetForm(database)))
    http.HandleFunc("/api/form/add", handlers.RequireJWT(database, handlers.AjoutForm(database)))
    http.HandleFunc("/api/users", handlers.RequireJWT(database, handlers.GetAllUsersHandler(database)))

    // Routes publiques (pas de JWT)
    http.HandleFunc("/api/user/add", handlers.AjoutUser(database))
    http.HandleFunc("/api/login", handlers.LoginHandler(database))
}