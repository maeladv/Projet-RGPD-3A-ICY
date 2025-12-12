package middleware

import (
	"database/sql"
	"net/http"
)

func RequireRole(db *sql.DB, role string, next http.HandlerFunc) http.HandlerFunc {
    return RequireJWT(db, func(w http.ResponseWriter, r *http.Request) {
        claims := GetClaims(r)
        if claims == nil || claims.Role != role {
            http.Error(w, "Accès interdit", http.StatusForbidden)
            return
        }
        next(w, r)
    })
}