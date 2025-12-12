package main

import (
	"backend/db"
	"backend/handlers"
	"backend/sftp"
	"log"
	"net/http"
	"os"
)

func main() {
	hostDB := "db"
	portDB := "5432"
	userDB := os.Getenv("POSTGRES_USER")
	passDB := os.Getenv("POSTGRES_PASSWORD")
	nameDB := "postgres"
	sslmodeDB := "disable"
	database, err := db.InitDB(hostDB, portDB, userDB, passDB, nameDB, sslmodeDB)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	hostSFTP := "sftp"
	portSFTP := "22"
	usernameSFTP := os.Getenv("SFTP_USERNAME")
	passwordSFTP := os.Getenv("SFTP_PASSWORD")
	sftp, err := sftp.InitSFTP(hostSFTP, portSFTP, usernameSFTP, passwordSFTP)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	http.HandleFunc("/api/forms", handlers.GetAllForms(database))
	http.HandleFunc("/api/form", handlers.GetForm(database))
	http.HandleFunc("/api/form/add", handlers.AjoutForm(database))

	log.Println("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
