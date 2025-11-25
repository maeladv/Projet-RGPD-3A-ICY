package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
    // Charger le fichier .env
    err := godotenv.Load()
    if err != nil {
        log.Println("Erreur lors du chargement du fichier .env, utilisation des variables d'environnement système")
    }

	pass := os.Getenv("POSTGRES_PASSWORD")
	host := "postgres"
	port := "5432"
	user := "postgres"
	dbname := "postgres"
	sslmode := "disable"
	
	// Construire la chaîne de connexion
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, pass, dbname, sslmode)
    
    // Ouvrir la connexion
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("erreur lors de l'ouverture de la connexion: %v", err)
    }
    
    // Vérifier la connexion
    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("erreur lors du ping de la base de données: %v", err)
    }
    
    log.Println("Connexion à la base de données réussie!")
    DB = db
    return db, nil
}

func CloseDB() {
    if DB != nil {
        DB.Close()
        log.Println("Connexion à la base de données fermée")
    }
}