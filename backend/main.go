package main

import (
	"backend/db"
	"backend/handlers"
	"log"
	"net/http"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	http.HandleFunc("/api/forms", handlers.GetAllForms(database))
	http.HandleFunc("/api/form", handlers.GetForm(database))
	http.HandleFunc("/api/form/add", handlers.AjoutForm(database))
	http.HandleFunc("/api/user/add", handlers.AjoutUser(database))

	http.HandleFunc("/api/users", handlers.GetAllUsersHandler(database))

	log.Println("[i] Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
