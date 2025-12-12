package handlers

import (
	"database/sql"
	"net/http"
)

func DeleteForm(db *sql.DB, formID string) error {
	_, err := db.Exec("DELETE FROM answers WHERE id = $1", formID)
	return err
}

func SuppForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		formID := r.URL.Query().Get("id")
		if formID == "" {
			http.Error(w, "ID du formulaire requis", http.StatusBadRequest)
			return
		}

		if err := DeleteForm(db, formID); err != nil {
			http.Error(w, "Erreur lors de la suppression du formulaire", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Formulaire supprimé avec succès"))
	}
}