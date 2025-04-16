// Package auth provides authentication and authorization functionality for the ExpertDB application
package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	
	"expertdb/internal/domain"
)

// Authentication related constants
const (
	// bcryptCost defines the computational cost for bcrypt password hashing
	// Higher values increase security but require more CPU resources
	bcryptCost = 12
	
	// jwtExpiration defines how long JWT tokens remain valid after issuance
	// Current setting: 24 hours
	jwtExpiration = time.Hour * 24
	
	// User role definitions for access control
	RoleAdmin = "admin" // Admin role has full system access
	RoleUser  = "user"  // User role has limited, read-mostly access
)

// JWTSecretKey is the key used to sign and verify JWT tokens
// This key is generated randomly at application startup
var JWTSecretKey []byte

// InitJWTSecret initializes the JWT secret key used for token signing and verification
func InitJWTSecret() error {
	// Generate a cryptographically secure random 32-byte key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return fmt.Errorf("failed to generate random JWT secret: %w", err)
	}
	
	// Store the key in the global variable
	JWTSecretKey = key
	
	return nil
}

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
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user with standard claims
func GenerateJWT(user *domain.User) (string, error) {
	// Calculate token expiration time
	expiration := time.Now().Add(jwtExpiration)
	
	// Create claims map with user information
	claims := jwt.MapClaims{
		"sub":   strconv.FormatInt(user.ID, 10), // Subject (user ID)
		"name":  user.Name,                      // User's full name
		"email": user.Email,                     // User's email address
		"role":  user.Role,                      // User's role (admin/user)
		"exp":   expiration.Unix(),              // Expiration timestamp
	}
	
	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token with the secret key
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}
	
	return tokenString, nil
}

// VerifyJWT verifies and parses a JWT token
func VerifyJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	// Initialize an empty claims map
	claims := jwt.MapClaims{}
	
	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm is as expected (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		// Return the secret key for signature verification
		return JWTSecretKey, nil
	})
	
	// Handle parsing errors (includes expiration checks)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}
	
	// Verify token is valid
	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}
	
	return token, claims, nil
}