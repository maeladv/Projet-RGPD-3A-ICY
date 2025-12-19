package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/sftp"
)

var mailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func AjoutForm(db *sql.DB, sftp *sftp.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		var f models.Form
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			err := r.ParseMultipartForm(10 << 20) // 10 MB max
			if err != nil {
				http.Error(w, "Erreur lors du parsing du formulaire", http.StatusBadRequest)
				return
			}

			f.Nom = strings.TrimSpace(r.FormValue("nom"))
			f.Prenom = strings.TrimSpace(r.FormValue("prenom"))
			f.VilleNaissance = strings.TrimSpace(r.FormValue("lieu_naissance"))
			f.NiveauDiplome = strings.TrimSpace(r.FormValue("diplome"))
			f.Adresse = strings.TrimSpace(r.FormValue("adresse_ligne1"))
			f.Complement = strings.TrimSpace(r.FormValue("adresse_ligne2"))
			f.CodePostal = strings.TrimSpace(r.FormValue("code_postal"))
			f.Ville = strings.TrimSpace(r.FormValue("ville"))
			f.Pays = strings.TrimSpace(r.FormValue("pays"))
			f.Mail = strings.TrimSpace(strings.ToLower(r.FormValue("email")))
			f.NumSecu = strings.TrimSpace(r.FormValue("securite_sociale"))
			f.Telephone = strings.TrimSpace(r.FormValue("telephone"))

			dateStr := r.FormValue("date_naissance")
			if dateStr != "" {
				parsedDate, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					http.Error(w, "Format de date invalide (YYYY-MM-DD)", http.StatusBadRequest)
					return
				}
				f.DateNaissance = parsedDate
			}

			cv, _, err := r.FormFile("cv")
			if err != nil && err != http.ErrMissingFile {
				http.Error(w, "Erreur lors du parsing du cv", http.StatusBadRequest)
				return
			} else if err == nil {
				defer cv.Close()

				f.CvPath = uuid.NewString() + ".pdf"

				remoteCv, err := sftp.Create("/data" + f.CvPath)
				if err != nil {
					fmt.Printf("erreur : %s\n", err)
					http.Error(w, "Erreur lors de la sauvegarde du cv", http.StatusInternalServerError)
					return
				}
				defer remoteCv.Close()

				_, err = io.Copy(remoteCv, cv)
				if err != nil {
					http.Error(w, "Erreur lors de la sauvegarde du cv", http.StatusInternalServerError)
					return
				}
			}
		} else if strings.Contains(contentType, "application/json") {
			var data map[string]any
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, "Données invalides", http.StatusBadRequest)
				return
			}

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
				f.DateNaissance = parsedDate
			}
			if villeNaissance, ok := data["ville_naissance"].(string); ok {
				f.VilleNaissance = strings.TrimSpace(villeNaissance)
			}
			if niveau, ok := data["niveau_diplome"].(string); ok {
				f.NiveauDiplome = strings.TrimSpace(niveau)
			}
			if adresse, ok := data["adresse"].(string); ok {
				f.Adresse = strings.TrimSpace(adresse)
			}
			if complement, ok := data["complement"].(string); ok {
				f.Complement = strings.TrimSpace(complement)
			}
			if codePostal, ok := data["code_postal"].(string); ok {
				f.CodePostal = strings.TrimSpace(codePostal)
			}
			if ville, ok := data["ville"].(string); ok {
				f.Ville = strings.TrimSpace(ville)
			}
			if pays, ok := data["pays"].(string); ok {
				f.Pays = strings.TrimSpace(pays)
			}
			if numSecu, ok := data["num_secu"].(string); ok {
				f.NumSecu = strings.TrimSpace(numSecu)
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

		// for i := 0; i < f.NumField(); i++ {
		// 	fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		// }
		if f.Nom == "" {
			http.Error(w, "Nom requis", http.StatusBadRequest)
			return
		}
		if f.Prenom == "" {
			http.Error(w, "Prénom requis", http.StatusBadRequest)
			return
		}
		if f.DateNaissance.IsZero() {
			http.Error(w, "Date de naissance requise", http.StatusBadRequest)
			return
		}
		if f.VilleNaissance == "" {
			http.Error(w, "Ville de naissance requise", http.StatusBadRequest)
			return
		}
		if f.NiveauDiplome == "" {
			http.Error(w, "Niveau de diplôme requis", http.StatusBadRequest)
			return
		}
		if f.Adresse == "" {
			http.Error(w, "Adresse requise", http.StatusBadRequest)
			return
		}
		if f.CodePostal == "" {
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
		if f.NumSecu == "" {
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

		query := `INSERT INTO answers (nom, prenom, date_naissance, ville_naissance, niveau_diplome, mail, 
                  adresse, complement_adresse, code_postal, ville, pays, num_secu_sociale, num_telephone)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

		err := db.QueryRow(query, f.Nom, f.Prenom, f.DateNaissance, f.VilleNaissance,
			f.NiveauDiplome, f.Mail, f.Adresse, f.Complement, f.CodePostal, f.Ville, f.Pays,
			f.NumSecu, f.Telephone).Scan(&f.ID)
		if err != nil {
			http.Error(w, "Erreur lors de l'enregistrement: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Formulaire enregistré avec succès",
			"id":      f.ID,
		})
		log.Printf("[+] Formulaire de %v enregistré", r.RemoteAddr)
	}
}
