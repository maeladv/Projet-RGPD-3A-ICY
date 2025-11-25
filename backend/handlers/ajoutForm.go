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

        var f models.Form
        contentType := r.Header.Get("Content-Type")

        // Parser selon le Content-Type
        if strings.Contains(contentType, "multipart/form-data") {
            // Parser multipart/form-data
            err := r.ParseMultipartForm(10 << 20) // 10 MB max
            if err != nil {
                http.Error(w, "Erreur lors du parsing du formulaire", http.StatusBadRequest)
                return
            }

            // Mapper les champs du formulaire HTML vers le modèle
            f.Nom = strings.TrimSpace(r.FormValue("nom"))
            f.Prenom = strings.TrimSpace(r.FormValue("prenom"))
            f.Ville_naissance = strings.TrimSpace(r.FormValue("lieu_naissance"))
            f.Niveau_diplome = strings.TrimSpace(r.FormValue("diplome"))
            f.Adresse = strings.TrimSpace(r.FormValue("adresse_ligne1"))
            f.Complement = strings.TrimSpace(r.FormValue("adresse_ligne2"))
            f.Code_postal = strings.TrimSpace(r.FormValue("code_postal"))
            f.Ville = strings.TrimSpace(r.FormValue("ville"))
            f.Pays = strings.TrimSpace(r.FormValue("pays"))
            f.Mail = strings.TrimSpace(strings.ToLower(r.FormValue("email")))
            f.Num_secu = strings.TrimSpace(r.FormValue("securite_sociale"))
            f.Telephone = strings.TrimSpace(r.FormValue("telephone"))

            // Parser la date
            dateStr := r.FormValue("date_naissance")
            if dateStr != "" {
                parsedDate, err := time.Parse("2006-01-02", dateStr)
                if err != nil {
                    http.Error(w, "Format de date invalide (YYYY-MM-DD)", http.StatusBadRequest)
                    return
                }
                f.Date_naissance = parsedDate
            }

        } else if strings.Contains(contentType, "application/json") {
            // Parser JSON (pour compatibilité avec API)
            var data map[string]interface{}
            if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
                http.Error(w, "Données invalides", http.StatusBadRequest)
                return
            }

            // Mapper les données JSON
            if nom, ok := data["nom"].(string); ok {
                f.Nom = strings.TrimSpace(nom)
            }
            if prenom, ok := data["prenom"].(string); ok {
                f.Prenom = strings.TrimSpace(prenom)
            }
            if dateStr, ok := data["date_naissance"].(string); ok {
                parsedDate, err := time.Parse("2006-01-02", dateStr)
                if err != nil {
                    http.Error(w, "Format de date invalide (YYYY-MM-DD)", http.StatusBadRequest)
                    return
                }
                f.Date_naissance = parsedDate
            }
            if ville_naissance, ok := data["ville_naissance"].(string); ok {
                f.Ville_naissance = strings.TrimSpace(ville_naissance)
            }
            if niveau, ok := data["niveau_diplome"].(string); ok {
                f.Niveau_diplome = strings.TrimSpace(niveau)
            }
            if adresse, ok := data["adresse"].(string); ok {
                f.Adresse = strings.TrimSpace(adresse)
            }
            if complement, ok := data["complement"].(string); ok {
                f.Complement = strings.TrimSpace(complement)
            }
            if codePostal, ok := data["code_postal"].(string); ok {
                f.Code_postal = strings.TrimSpace(codePostal)
            }
            if ville, ok := data["ville"].(string); ok {
                f.Ville = strings.TrimSpace(ville)
            }
            if pays, ok := data["pays"].(string); ok {
                f.Pays = strings.TrimSpace(pays)
            }
            if numSecu, ok := data["num_secu"].(string); ok {
                f.Num_secu = strings.TrimSpace(numSecu)
            }
            if telephone, ok := data["telephone"].(string); ok {
                f.Telephone = strings.TrimSpace(telephone)
            }
            if mail, ok := data["mail"].(string); ok {
                f.Mail = strings.TrimSpace(strings.ToLower(mail))
            }
        } else {
            http.Error(w, "Content-Type non supporté", http.StatusUnsupportedMediaType)
            return
        }

        // Validations
        if f.Nom == "" {
            http.Error(w, "Nom requis", http.StatusBadRequest)
            return
        }
        if f.Prenom == "" {
            http.Error(w, "Prénom requis", http.StatusBadRequest)
            return
        }
        if f.Date_naissance.IsZero() {
            http.Error(w, "Date de naissance requise", http.StatusBadRequest)
            return
        }
        if f.Ville_naissance == "" {
            http.Error(w, "Ville de naissance requise", http.StatusBadRequest)
            return
        }
        if f.Niveau_diplome == "" {
            http.Error(w, "Niveau de diplôme requis", http.StatusBadRequest)
            return
        }
        if f.Adresse == "" {
            http.Error(w, "Adresse requise", http.StatusBadRequest)
            return
        }
        if f.Code_postal == "" {
            http.Error(w, "Code postal requis", http.StatusBadRequest)
            return
        }
        if f.Ville == "" {
            http.Error(w, "Ville requise", http.StatusBadRequest)
            return
        }
        if f.Pays == "" {
            http.Error(w, "Pays requis", http.StatusBadRequest)
            return
        }
        if f.Num_secu == "" {
            http.Error(w, "Numéro de sécurité sociale requis", http.StatusBadRequest)
            return
        }
        if f.Telephone == "" {
            http.Error(w, "Téléphone requis", http.StatusBadRequest)
            return
        }
        if f.Mail == "" {
            http.Error(w, "Email requis", http.StatusBadRequest)
            return
        }
        if !mailRegex.MatchString(f.Mail) {
            http.Error(w, "Format d'email invalide", http.StatusBadRequest)
            return
        }

        // Enregistrer dans la DB
        query := `INSERT INTO answers (nom, prenom, date_naissance, ville_naissance, niveau_diplome, mail, 
                  adresse, complement_adresse, code_postal, ville, pays, num_secu_sociale, num_telephone)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

        err := db.QueryRow(query, f.Nom, f.Prenom, f.Date_naissance, f.Ville_naissance,
            f.Niveau_diplome, f.Mail, f.Adresse, f.Complement, f.Code_postal, f.Ville, f.Pays,
            f.Num_secu, f.Telephone).Scan(&f.Id)

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