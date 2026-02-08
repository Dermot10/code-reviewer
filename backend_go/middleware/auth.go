package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization format", http.StatusUnauthorized)
				return
			}

			// split header up and xtract token, and validate JWT

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			idClaim, ok := claims["user_id"]
			if !ok {
				http.Error(w, "user_id missing in token", http.StatusUnauthorized)
				return
			}

			idFloat, ok := idClaim.(float64) // jwt numeric claims - float64
			if !ok {
				http.Error(w, "user_id invalid type", http.StatusUnauthorized)
				return
			}

			userID := uint(idFloat)
			// add user_id to context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// next handler
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
