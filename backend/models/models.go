package models

import (
	"time"
)

type Form struct {
	ID             int
	Nom            string
	Prenom         string
	DateNaissance  time.Time
	VilleNaissance string
	NiveauDiplome  string
	Mail           string
	Adresse        string
	Complement     string
	CodePostal     string
	Ville          string
	Pays           string
	NumSecu        string
	Telephone      string
}
