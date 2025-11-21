package main

import (
	"fmt"

	"github.com/vladislavkn/texterstellungsuebungsmaschine/internal/auth"
)

func main() {
	// Test password hash stored in main.go
	storedHash := "$2a$10$P9iTmndOegEkd9OZ0DZxNOvWzemcb/bGGguvbWVBqGyCyX36vm77q"
	testPassword := "password123"

	result := auth.VerifyPassword(storedHash, testPassword)
	fmt.Printf("Stored hash: %s\n", storedHash)
	fmt.Printf("Test password: %s\n", testPassword)
	fmt.Printf("Verification result: %v\n", result)

	// Generate new hash and verify
	fmt.Println("\n--- Generating new hash ---")
	newHash, _ := auth.HashPassword(testPassword)
	fmt.Printf("New hash: %s\n", newHash)
	result2 := auth.VerifyPassword(newHash, testPassword)
	fmt.Printf("Verification result: %v\n", result2)
}
