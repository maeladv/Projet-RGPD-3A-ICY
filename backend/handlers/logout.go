package handlers

import (
    "database/sql"
    "net/http"
)

type logoutRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}


func LogoutHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {


        http.SetCookie(w, &http.Cookie{
            Name:     "jwt",
            Value:    "",
            Path:     "/",
            HttpOnly: true,
            Secure:   false, // mets true si HTTPS
            SameSite: http.SameSiteLaxMode,
            MaxAge:   -1, // Supprime le cookie
        })

        // On redirige vers /login
        http.Redirect(w, r, "/login", http.StatusSeeOther)

        

    }
}