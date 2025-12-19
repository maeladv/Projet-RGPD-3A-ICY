package models

import (
	"time"
)

type Form struct {
	ID             int       `json:"id"`
	Nom            string    `json:"nom"`
	Prenom         string    `json:"prenom"`
	DateNaissance  time.Time `json:"date_naissance"`
	VilleNaissance string    `json:"ville_naissance"`
	NiveauDiplome  string    `json:"niveau_diplome"`
	Mail           string    `json:"mail"`
	Adresse        string    `json:"adresse"`
	Complement     string    `json:"complement"`
	CodePostal     string    `json:"code_postal"`
	Ville          string    `json:"ville"`
	Pays           string    `json:"pays"`
	NumSecu        string    `json:"num_secu_sociale"`
	Telephone      string    `json:"num_telephone"`
}

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

type Session struct {
	ID        int
	UserID    int
	Token     string
	ExpiresAt time.Time
}
