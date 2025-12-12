// Package db: initialisation de la connection à la db
package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(host, port, user, password, dbname, sslmode string) (*sql.DB, error) {
	log.Printf("[i] Tentative de connexion à PostgreSQL: host=%s port=%s user=%s dbname=%s", host, port, user, dbname)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var db *sql.DB
	var err error

	maxRetries := 5
	for i := range maxRetries {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("[!] Tentative %d/%d: erreur lors de la connection à PostgreSQL: %v", i+1, maxRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("[i] Connexion à la base de données réussie")
			return db, nil
		}

		log.Printf("[!] Tentative %d/%d: erreur ping PostgreSQL: %v", i+1, maxRetries, err)
		db.Close()
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("[!] Impossible de se connecter après %d tentatives: %v", maxRetries, err)
}
