package middleware

import (
	"net/http"
	"product-api/internal/jwt"
	"strings"
)

// JWTMiddleware validates the JWT token from the Authorization header.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			http.Error(w, "Malformed Authorization header", http.StatusUnauthorized)
			return
		}

		// Validate the token
		_, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
