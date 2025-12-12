package handlers

import (
    "database/sql"
    "net/http"
    "time"
)


func LogoutHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        // Suppression robuste du cookie
        http.SetCookie(w, &http.Cookie{
            Name:     "jwt",
            Value:    "",
            Path:     "/",
            HttpOnly: true,
            Secure:   false, // mets true si HTTPS
            SameSite: http.SameSiteLaxMode,
            MaxAge:   -1, // Supprime le cookie
            Expires:  time.Unix(0, 0), // Force l'expiration (date dans le passé)
        })

        // Empêcher le cache de la redirection de déconnexion
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")

        // on retroune "deconnexion réussie en json
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message":"Déconnexion réussie"}`))
    }
}