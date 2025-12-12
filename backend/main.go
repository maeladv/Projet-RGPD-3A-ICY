package main

import (
    "backend/auth"
    "backend/db"
	"backend/routes"
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

	routes.SetupRoutes(database)

    log.Println("[i] Serveur démarré sur le port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}