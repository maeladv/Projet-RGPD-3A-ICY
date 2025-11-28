package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formID := r.URL.Query().Get("id")

		var form models.Form
		query := `SELECT id, nom, prenom, date_naissance, ville_naissance, niveau_diplome, 
                  adresse, complement_adresse, code_postal, ville, pays, num_secu_sociale, num_telephone 
                  FROM answers WHERE id = $1`

		err := db.QueryRow(query, formID).Scan(
			&form.ID, &form.Nom, &form.Prenom, &form.DateNaissance,
			&form.VilleNaissance, &form.NiveauDiplome, &form.Adresse,
			&form.Complement, &form.CodePostal, &form.Ville, &form.Pays,
			&form.NumSecu, &form.Telephone,
		)

		if err == sql.ErrNoRows {
			http.Error(w, "Form not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(form)
	}
}

func GetAllForms(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `SELECT id, nom, prenom, date_naissance, ville_naissance, niveau_diplome,
                  adresse, complement_adresse, code_postal, ville, pays, num_secu_sociale, num_telephone 
                  FROM answers`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var forms []models.Form
		for rows.Next() {
			var form models.Form
			if err := rows.Scan(
				&form.ID, &form.Nom, &form.Prenom, &form.DateNaissance,
				&form.VilleNaissance, &form.NiveauDiplome, &form.Adresse,
				&form.Complement, &form.CodePostal, &form.Ville, &form.Pays,
				&form.NumSecu, &form.Telephone,
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			forms = append(forms, form)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(forms)
	}
}
