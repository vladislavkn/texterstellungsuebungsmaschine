package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vladislavkn/texterstellungsuebungsmaschine/internal/auth"
)

// In-memory user database for demo purposes
var users = map[string]*auth.User{
	"testuser": {
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "$2a$10$P9iTmndOegEkd9OZ0DZxNOvWzemcb/bGGguvbWVBqGyCyX36vm77q", // hashed "password123"
	},
}

func main() {
	// Initialize JWT Manager
	jwtManager := auth.NewJWTManager(
		"your-secret-key-min-32-chars-long!",
		"your-refresh-secret-key-min-32-chars!",
		15*time.Minute,
		7*24*time.Hour,
	)

	// Routes
	http.HandleFunc("/health", auth.HealthHandler)
	http.HandleFunc("/auth/register", auth.RegisterHandler(users))
	http.HandleFunc("/auth/login", auth.LoginHandler(jwtManager, users))
	http.HandleFunc("/auth/refresh", auth.RefreshHandler(jwtManager))
	http.Handle("/protected", auth.AuthMiddleware(jwtManager, http.HandlerFunc(auth.ProtectedHandler)))

	port := ":8080"
	fmt.Printf("Server starting on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
