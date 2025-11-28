package main

import (
	"backend/db"
	"backend/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	host := "db"
	port := "5432"
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := "postgres"
	sslmode := "disable"
	database, err := db.InitDB(host, port, user, pass, dbname, sslmode)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	http.HandleFunc("/api/forms", handlers.GetAllForms(database))
	http.HandleFunc("/api/form", handlers.GetForm(database))
	http.HandleFunc("/api/form/add", handlers.AjoutForm(database))

	log.Println("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
