package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)	

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var u models.User
	err := db.QueryRow("SELECT id, username FROM users WHERE username = $1", username).
		Scan(&u.ID, &u.Username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	var u models.User
	err := db.QueryRow("SELECT id, username, role FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Username, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetAllUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query("SELECT id, username, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// HANDLERS HTTP

func GetAllUsersHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        users, err := GetAllUsers(db)
        if err != nil {
            http.Error(w, "Erreur lors de la récupération des utilisateurs", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}

func GetUserByIDHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := r.URL.Query().Get("id")
        if idStr == "" {
            http.Error(w, "ID manquant", http.StatusBadRequest)
            return
        }
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "ID invalide", http.StatusBadRequest)
            return
        }
        user, err := GetUserByID(db, id)
        if err != nil {
            http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}