package handlers

import (
	"backend/auth"
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}
		if req.Username == "" || req.Password == "" {
			http.Error(w, "Champs requis manquants", http.StatusBadRequest)
			return
		}

		// Récupère utilisateur avec hash
		var u models.User
		var hashed string
		err := db.QueryRow(
			"SELECT id, username, role, password FROM users WHERE username = $1",
			req.Username,
		).Scan(&u.ID, &u.Username, &u.Role, &hashed)
		if err == sql.ErrNoRows {
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		// Compare le mot de passe
		if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(req.Password)); err != nil {
			http.Error(w, "Identifiants invalides", http.StatusUnauthorized)
			return
		}

		// Génère le token
		token, err := auth.GenerateJWT(u.ID, u.Username, u.Role)
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // mets true si HTTPS
			SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{Token: token, User: u})
	}
}

