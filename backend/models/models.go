package models

import("time")

type Form struct {
	Id int
	Nom string
	Prenom string
	Date_naissance time.Time
	Niveau_diplome string
	// Addresse
	Adresse string
	Complement string
	Code_postal int
	Pays string
	Num_secu string
	Telephone string
}

