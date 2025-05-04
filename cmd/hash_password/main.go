package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	existingHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

	// Generate new hash
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		return
	}

	fmt.Printf("New hash: %s\n", string(newHash))

	// Verify the existing hash
	err = bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(password))
	if err != nil {
		fmt.Printf("Error verifying existing hash: %v\n", err)
		return
	}

	fmt.Println("Existing hash verification successful!")
}
