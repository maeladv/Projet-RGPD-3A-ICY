package handlers

import (
    "backend/auth"
    "context"
    "database/sql"
    "net/http"
    "strings"
)

type ctxKey string
var userKey ctxKey = "userClaims"

func RequireJWT(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var token string

        // 1) Tente via Authorization: Bearer
        h := r.Header.Get("Authorization")
        if h != "" && strings.HasPrefix(h, "Bearer ") {
            token = strings.TrimPrefix(h, "Bearer ")
        }

        // 2) Fallback: cookie "jwt"
        if token == "" {
            if c, err := r.Cookie("jwt"); err == nil && c.Value != "" {
                token = c.Value
            }
        }

        if token == "" {
            http.Error(w, "Non autorisé", http.StatusUnauthorized)
            return
        }
        claims, err := auth.ParseJWT(token)
        if err != nil {
            http.Error(w, "Token invalide", http.StatusUnauthorized)
            return
        }
        // Option: vérifier que l’utilisateur existe toujours
        // (sécurité supplémentaire)
        ctx := context.WithValue(r.Context(), userKey, claims)
        next(w, r.WithContext(ctx))
    }
}

func GetClaims(r *http.Request) *auth.Claims {
    v := r.Context().Value(userKey)
    if c, ok := v.(*auth.Claims); ok {
        return c
    }
    return nil
}