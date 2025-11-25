package models

import("time")

type Form struct {
	Id int
	Nom string
	Prenom string
	Date_naissance time.Time
	Ville_naissance string
	Niveau_diplome string
	Mail string
	// Addresse
	Adresse string
	Complement string
	Code_postal string
	Ville string
	Pays string
	Num_secu string
	Telephone string
}

