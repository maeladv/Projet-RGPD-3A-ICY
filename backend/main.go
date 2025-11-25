package main

import (
    "backend/db"
    "backend/handlers"
    "log"
    "net/http"
)

func main() {
    // Initialiser la connexion à la base de données
    database, err := db.InitDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.CloseDB()
    
    // Configurer vos routes
    http.HandleFunc("/api/forms", handlers.GetAllForms(database))
    http.HandleFunc("/api/form", handlers.GetForm(database))
    
    // Démarrer le serveur
    log.Println("Serveur démarré sur le port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
