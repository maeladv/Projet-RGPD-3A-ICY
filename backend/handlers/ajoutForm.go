package handlers

import (
	"backend/models"
	"net/http"
	"encoding/json"
	"strings"
	"time"
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

	// Prenom (obligatoire, nettoyage)
    if prenom, ok := data["prenom"].(string); ok {
        f.Prenom = strings.TrimSpace(prenom)
        if f.Prenom == "" {
            http.Error(w, "Le prénom est obligatoire", http.StatusBadRequest)
            return
        }
    } else {
        http.Error(w, "Le prénom est obligatoire", http.StatusBadRequest)
        return
    }

    // Date de naissance (obligatoire)
    if dateStr, ok := data["date_naissance"].(string); ok {
        parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
        if err != nil {
            http.Error(w, "Format de date invalide (YYYY-MM-DD)", http.StatusBadRequest)
            return
        }
        f.Date_naissance = parsedDate
    } else {
        http.Error(w, "La date de naissance est obligatoire", http.StatusBadRequest)
        return
    }

    // Niveau de diplôme (obligatoire)
    if niveau, ok := data["niveau_diplome"].(string); ok {
        f.Niveau_diplome = strings.TrimSpace(niveau)
        if f.Niveau_diplome == "" {
            http.Error(w, "Le niveau de diplôme est obligatoire", http.StatusBadRequest)
            return
        }
    } else {
        http.Error(w, "Le niveau de diplôme est obligatoire", http.StatusBadRequest)
        return
    }

    // Adresse (obligatoire)
    if adresse, ok := data["adresse"].(string); ok {
        f.Adresse = strings.TrimSpace(adresse)
        if f.Adresse == "" {
            http.Error(w, "L'adresse est obligatoire", http.StatusBadRequest)
            return
        }
    } else {
        http.Error(w, "L'adresse est obligatoire", http.StatusBadRequest)
        return
    }

    // Complément d'adresse (optionnel)
    if complement, ok := data["complement"].(string); ok {
        f.Complement = strings.TrimSpace(complement)
    }

    // Code postal (obligatoire)
    if codePostal, ok := data["code_postal"].(float64); ok {
        f.Code_postal = int(codePostal)
        if f.Code_postal <= 0 {
            http.Error(w, "Code postal invalide", http.StatusBadRequest)
            return
        }
    } else {
        http.Error(w, "Le code postal est obligatoire", http.StatusBadRequest)
        return
    }

    // Pays (obligatoire)
    if pays, ok := data["pays"].(string); ok {
        f.Pays = strings.TrimSpace(pays)
        if f.Pays == "" {
            http.Error(w, "Le pays est obligatoire", http.StatusBadRequest)
            return
        }
    } else {
        http.Error(w, "Le pays est obligatoire", http.StatusBadRequest)
        return
    }

    // Numéro de sécurité sociale (obligatoire)
    if numSecu, ok := data["num_secu"].(string); ok {
        f.Num_secu = strings.TrimSpace(numSecu)
        if f.Num_secu == "" {
            http.Error(w, "Le numéro de sécurité sociale est obligatoire", http.StatusBadRequest)
            return
        }
        // Validation basique (15 chiffres en France)
        if len(f.Num_secu) != 15 {
            http.Error(w, "Numéro de sécurité sociale invalide (15 chiffres attendus)", http.StatusBadRequest)
            return
        }
    } else {
        http.Error(w, "Le numéro de sécurité sociale est obligatoire", http.StatusBadRequest)
        return
    }

    // Téléphone (obligatoire)
    if telephone, ok := data["telephone"].(string); ok && telephone != "" {
        f.Telephone = strings.TrimSpace(telephone)
		if f.Telephone == "" {
			http.Error(w, "Le numéro de télephone est obligatoire", http.StatusBadRequest)
            return
		} else if len(f.Telephone) < 10 {
            http.Error(w, "Numéro de téléphone invalide", http.StatusBadRequest)
            return
        }
    }

	// Enregistrer dans la DB

	// Reponse
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Formulaire enregistré avec succès",
        "id": f.Id,
    })
}