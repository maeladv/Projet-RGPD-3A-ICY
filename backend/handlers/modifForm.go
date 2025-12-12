package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

func ModifForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		var form models.Form
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		query := `UPDATE answers SET nom=$1, prenom=$2, date_naissance=$3, ville_naissance=$4,
				  niveau_diplome=$5, mail=$6, adresse=$7, complement_adresse=$8,
				  code_postal=$9, ville=$10, pays=$11, num_secu_sociale=$12, num_telephone=$13
				  WHERE id=$14`

		_, err := db.Exec(query,
			form.Nom, form.Prenom, form.DateNaissance, form.VilleNaissance,
			form.NiveauDiplome, form.Mail, form.Adresse, form.Complement,
			form.CodePostal, form.Ville, form.Pays, form.NumSecu, form.Telephone,
			form.ID,
		)
		if err != nil {
			http.Error(w, "Erreur lors de la mise à jour du formulaire", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Formulaire mis à jour avec succès"))
	}
}