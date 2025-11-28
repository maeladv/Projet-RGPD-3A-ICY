package main

import (
    "backend/auth"
    "backend/db"
    "backend/handlers"
    "log"
    "net/http"
)

func main() {
    auth.InitJWTSecret()

    database, err := db.InitDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.CloseDB()

    http.HandleFunc("/api/forms", handlers.RequireJWT(database, handlers.GetAllForms(database)))
    http.HandleFunc("/api/form", handlers.RequireJWT(database, handlers.GetForm(database)))
    http.HandleFunc("/api/form/add", handlers.RequireJWT(database, handlers.AjoutForm(database)))

    http.HandleFunc("/api/user/add", handlers.AjoutUser(database))       // inscription publique
    http.HandleFunc("/api/login", handlers.LoginHandler(database))       // login pour obtenir JWT
    http.HandleFunc("/api/users", handlers.RequireJWT(database, handlers.GetAllUsersHandler(database)))

    log.Println("[i] Serveur démarré sur le port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}