package main

import (
	"backend/auth"
	"backend/db"
	"backend/routes"
	"backend/sftp"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	auth.InitJWTSecret()

	/*---[Connection Database]---*/

	hostDB := "db"
	portDB := "5432"
	userDB := os.Getenv("POSTGRES_USER")
	passDB := os.Getenv("POSTGRES_PASSWORD")
	nameDB := userDB
	sslmodeDB := "disable"
	database, err := db.InitDB(hostDB, portDB, userDB, passDB, nameDB, sslmodeDB)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	/*---[Connection SFTP]---*/

	hostSFTP := "sftp:2222"
	usernameSFTP := os.Getenv("USER_NAME")
	privateKey, err := loadPrivateKey("/keys/id_ed25519")
	if err != nil {
		log.Fatal(err)
	}
	sftp, err := sftp.InitSFTP(hostSFTP, usernameSFTP, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	/*---[Routing]---*/

	routes.SetupRoutes(database, sftp)

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
