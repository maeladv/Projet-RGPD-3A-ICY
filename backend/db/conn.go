package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := "db" // nom docker db
	port := "5432"
	user := "postgres"
	dbname := "postgres"
	sslmode := "disable"

	log.Printf("[i] Tentative de connexion à PostgreSQL: host=%s port=%s user=%s dbname=%s", host, port, user, dbname)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, sslmode)

	var db *sql.DB
	var err error

	// Retry logic - attendre que PostgreSQL soit prêt
	maxRetries := 10
	for i := range maxRetries {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("[!] Tentative %d/%d: erreur lors de l'ouverture: %v", i+1, maxRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("[i] Connexion à la base de données réussie")
			DB = db
			return db, nil
		}

		log.Printf("[!] Tentative %d/%d: erreur ping: %v", i+1, maxRetries, err)
		db.Close()
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("[!] Impossible de se connecter après %d tentatives: %v", maxRetries, err)
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("[i] Connexion à la base de données fermée")
	}
}

