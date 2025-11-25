package handlers

import (
    "backend/models"
    "database/sql"
    "encoding/json"
    "net/http"
    "regexp"
    "strings"
    "time"
)

// Regex pour valider l'email
var mailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func AjoutForm(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Vérifier la méthode HTTP
        if r.Method != http.MethodPost {
            http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
            return
        }

        // 2. Parser les données du formulaire
        var f models.Form

        // Parser les données
        var data map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "Données invalides", http.StatusBadRequest)
            return
        }

        // Traiter chaque champ individuellement
        // Nom (obligatoire, nettoyage)
        if nom, ok := data["nom"].(string); ok {
            f.Nom = strings.TrimSpace(nom)
        } else {
            http.Error(w, "Nom requis", http.StatusBadRequest)
            return
        }

        // Prenom (obligatoire, nettoyage)
        if prenom, ok := data["prenom"].(string); ok {
            f.Prenom = strings.TrimSpace(prenom)
        } else {
            http.Error(w, "Prénom requis", http.StatusBadRequest)
            return
        }

        // Date de naissance (obligatoire)
        if dateStr, ok := data["date_naissance"].(string); ok {
            parsedDate, err := time.Parse("2006-01-02", dateStr)
            if err != nil {
                http.Error(w, "Format de date invalide (YYYY-MM-DD)", http.StatusBadRequest)
                return
            }
            f.Date_naissance = parsedDate // Assigner le time.Time
        } else {
            http.Error(w, "Date de naissance requise", http.StatusBadRequest)
            return
        }

        // Ville de naissance (Obligatoire)
        if ville_naissance, ok := data["ville_naissance"].(string); ok {
            f.Ville_naissance = strings.TrimSpace(ville_naissance)
        } else {
            http.Error(w, "Ville de naissance requise", http.StatusBadRequest)
            return
        }

        // Niveau de diplôme (obligatoire)
        if niveau, ok := data["niveau_diplome"].(string); ok {
            f.Niveau_diplome = strings.TrimSpace(niveau)
        } else {
            http.Error(w, "Niveau de diplôme requis", http.StatusBadRequest)
            return
        }

        // Adresse (obligatoire)
        if adresse, ok := data["adresse"].(string); ok {
            f.Adresse = strings.TrimSpace(adresse)
        } else {
            http.Error(w, "Adresse requise", http.StatusBadRequest)
            return
        }

        // Complément d'adresse (optionnel)
        if complement, ok := data["complement"].(string); ok {
            f.Complement = strings.TrimSpace(complement)
        }

        // Code postal (obligatoire)
        if codePostal, ok := data["code_postal"].(string); ok {
            f.Code_postal = strings.TrimSpace(codePostal)
        } else {
            http.Error(w, "Code postal requis", http.StatusBadRequest)
            return
        }

        // Ville (Obligatoire)
        if ville, ok := data["ville"].(string); ok {
            f.Ville = strings.TrimSpace(ville)
        } else {
            http.Error(w, "Ville requise", http.StatusBadRequest)
            return
        }

        // Pays (obligatoire)
        if pays, ok := data["pays"].(string); ok {
            f.Pays = strings.TrimSpace(pays)
        } else {
            http.Error(w, "Pays requis", http.StatusBadRequest)
            return
        }

        // Numéro de sécurité sociale (obligatoire)
        if numSecu, ok := data["num_secu"].(string); ok {
            f.Num_secu = strings.TrimSpace(numSecu)
        } else {
            http.Error(w, "Numéro de sécurité sociale requis", http.StatusBadRequest)
            return
        }

        // Téléphone (obligatoire)
        if telephone, ok := data["telephone"].(string); ok {
            f.Telephone = strings.TrimSpace(telephone)
        } else {
            http.Error(w, "Téléphone requis", http.StatusBadRequest)
            return
        }

        // Email (obligatoire avec validation)
        if mail, ok := data["mail"].(string); ok {
            mail = strings.TrimSpace(strings.ToLower(mail))
            if !mailRegex.MatchString(mail) {
                http.Error(w, "Format d'email invalide", http.StatusBadRequest)
                return
            }
            f.Mail = mail
        } else {
            http.Error(w, "Mail requis", http.StatusBadRequest)
            return
        }

        // Enregistrer dans la DB
        query := `INSERT INTO answers (nom, prenom, date_naissance, ville_naissance, niveau_diplome, mail, 
                  adresse, complement_adresse, code_postal, ville, pays, num_secu_sociale, num_telephone)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

        err := db.QueryRow(query, f.Nom, f.Prenom, f.Date_naissance, f.Ville_naissance,
            f.Niveau_diplome, f.Adresse, f.Complement, f.Code_postal, f.Ville, f.Pays,
            f.Num_secu, f.Telephone, f.Mail).Scan(&f.Id)

        if err != nil {
            http.Error(w, "Erreur lors de l'enregistrement: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Réponse
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "Formulaire enregistré avec succès",
            "id":      f.Id,
        })
    }
}