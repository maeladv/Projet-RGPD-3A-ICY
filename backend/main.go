package main

import (
	"backend/db"
	"backend/handlers"
	"backend/sftp"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/ssh"
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
	portSFTP := "2222"
	privateKey, err := loadPrivateKey("/keys/id_ed25519")
	if err != nil {
		log.Fatal(err)
	}
	usernameSFTP := os.Getenv("USER_NAME")
	sftp, err := sftp.InitSFTP(hostSFTP, portSFTP, usernameSFTP, privateKey)
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

func loadPrivateKey(path string) (*ssh.Signer, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[!] echec de la lecture de la clé privée (%s): %v", path, err)
	}

	privateKey, err := ssh.ParsePrivateKey(content)
	if err != nil {
		return nil, fmt.Errorf("[!] echec du parsing de la clé privée: %v", err)
	}

	return &privateKey, nil
}

func loadPublicKey(path string) (*ssh.PublicKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[!] echec de la lecture de la clé publique (%s): %v", path, err)
	}

	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(content)
	if err != nil {
		return nil, fmt.Errorf("[!] echec du parsing de la clé publique: %v", err)
	}

	return &publicKey, nil
}
