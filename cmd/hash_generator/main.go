package main

import (
	"fmt"

	"github.com/vladislavkn/texterstellungsuebungsmaschine/internal/auth"
)

func main() {
	hash, _ := auth.HashPassword("password123")
	fmt.Println("Hashed password for 'password123':")
	fmt.Println(hash)
}
