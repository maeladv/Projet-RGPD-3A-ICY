package auth

import (
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWTSecret() {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "change-me"
    }
    jwtSecret = []byte(secret)
}

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateJWT(userID int, username, role string) (string, error) {
    claims := &Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseJWT(tokenStr string) (*Claims, error) {
    tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := tok.Claims.(*Claims); ok && tok.Valid {
        return claims, nil
    }
    return nil, jwt.ErrTokenInvalidClaims
}