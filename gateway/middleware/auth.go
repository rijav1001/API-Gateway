package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return []byte(secret)
	}
	return []byte("super-secret-key")
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}