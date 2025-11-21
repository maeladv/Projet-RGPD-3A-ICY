package handlers

import (
	"backend/models"
	"net/http"
	"encoding/json"
	"strings"
)

func ajoutForm(w http.ResponseWriter, r *http.Request)  {
	// 1. Vérifier la méthode HTTP
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

	// 2. Parser les données du formulaire
	var f models.Form

	 //Parser les données
    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Erreur de décodage JSON", http.StatusBadRequest)
        return
    }

	// Traiter chaque champ individuellement
    // Nom (obligatoire, nettoyage)
    if nom, ok := data["nom"].(string); ok {
        f.Nom = strings.TrimSpace(nom)
        if f.Nom == "" {
            http.Error(w, "Le nom est obligatoire", http.StatusBadRequest)
            return
        }
    } else {
        http.Error(w, "Le nom est obligatoire", http.StatusBadRequest)
        return
    }
}