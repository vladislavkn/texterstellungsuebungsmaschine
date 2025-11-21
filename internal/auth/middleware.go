package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// AuthMiddleware validates JWT tokens from the Authorization header
func AuthMiddleware(jm *JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "missing authorization header"})
			return
		}

		// Extract the token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// Validate the token
		claims, err := jm.ValidateToken(tokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid token: " + err.Error()})
			return
		}

		// Store claims in request context for use in handlers
		r.Header.Set("X-User-ID", string(rune(claims.UserID)))
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-Email", claims.Email)

		next.ServeHTTP(w, r)
	})
}
