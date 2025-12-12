package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, username, password, role string) error {
	_, err := db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", username, password, role)
	return err
}

func AjoutUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		var u models.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		if u.Username == "" || u.Password == "" || u.Role == "" {
			http.Error(w, "Champs requis manquants", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
			return
}

		existingUser, err := GetUserByUsername(db, u.Username)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		if existingUser != nil {
			http.Error(w, "Nom d'utilisateur déjà pris", http.StatusConflict)
			return
		}

		if err := CreateUser(db, u.Username, string(hashedPassword), u.Role); err != nil {
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Utilisateur créé avec succès"))
	}
}