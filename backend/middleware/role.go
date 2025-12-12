package middleware

import (
	"database/sql"
	"net/http"

	"golang.org/x/exp/slices"
)

func RequireRole(db *sql.DB, roles []string, next http.HandlerFunc) http.HandlerFunc {
	return RequireJWT(db, func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaims(r)
		if claims == nil || !slices.Contains(roles, claims.Role) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
