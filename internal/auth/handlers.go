package auth

import (
	"encoding/json"
	"net/http"
)

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

// RegisterHandler creates a new user account
func RegisterHandler(users map[string]*User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
			return
		}

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		var registerReq RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
			return
		}

		// Validate input
		if registerReq.Username == "" || registerReq.Email == "" || registerReq.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "username, email, and password are required"})
			return
		}

		// Check if user already exists
		if _, exists := users[registerReq.Username]; exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "username already exists"})
			return
		}

		// Hash password
		hashedPassword, err := HashPassword(registerReq.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to hash password"})
			return
		}

		// Create new user
		newUserID := len(users) + 1
		newUser := &User{
			ID:       newUserID,
			Username: registerReq.Username,
			Email:    registerReq.Email,
			Password: hashedPassword,
		}

		// Store user
		users[registerReq.Username] = newUser

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(RegisterResponse{
			ID:       newUser.ID,
			Username: newUser.Username,
			Email:    newUser.Email,
			Message:  "User registered successfully",
		})
	}
}

// LoginHandler authenticates a user and returns JWT tokens
func LoginHandler(jm *JWTManager, users map[string]*User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
			return
		}

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		var loginReq LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
			return
		}

		// Validate credentials
		user, exists := users[loginReq.Username]
		if !exists || !VerifyPassword(user.Password, loginReq.Password) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid username or password"})
			return
		}

		// Generate tokens
		tokens, err := jm.GenerateTokens(user.ID, user.Username, user.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to generate tokens"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokens)
	}
}

// RefreshHandler generates a new access token from a refresh token
func RefreshHandler(jm *JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
			return
		}

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		var refreshReq RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
			return
		}

		// Generate new access token
		newAccessToken, err := jm.RefreshAccessToken(refreshReq.RefreshToken)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid refresh token"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": newAccessToken,
			"expires_in":   900,
		})
	}
}

// ProtectedHandler is an example protected endpoint
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	// Access user info from headers set by middleware
	userID := r.Header.Get("X-User-ID")
	username := r.Header.Get("X-Username")
	email := r.Header.Get("X-Email")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "This is a protected resource",
		"user": map[string]string{
			"id":       userID,
			"username": username,
			"email":    email,
		},
	})
}

// HealthHandler checks if the server is running
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
