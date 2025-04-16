package auth

import (
	"fmt"
	
	"golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash creates a secure bcrypt hash from a plaintext password
func GeneratePasswordHash(password string) (string, error) {
	// Hash the password using bcrypt with the configured cost factor
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Return the hash as a string, suitable for database storage
	return string(hash), nil
}

// VerifyPassword checks if a plaintext password matches a previously hashed password
func VerifyPassword(password, hash string) bool {
	// Compare the provided password against the hash
	// bcrypt.CompareHashAndPassword handles all the work of extracting the salt and cost
	// parameters from the hash and performing the comparison securely
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}