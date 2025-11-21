package auth

import "golang.org/x/crypto/bcrypt"

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // Never expose password
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// VerifyPassword compares a hashed password with a plaintext password
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
