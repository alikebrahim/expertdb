package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Cost for bcrypt password hashing
	bcryptCost = 12
	
	// JWT expiration time
	jwtExpiration = time.Hour * 24 // 24 hours
	
	// Admin role
	RoleAdmin = "admin"
	RoleUser  = "user"
)

var (
	// JWTSecretKey is the key used to sign JWT tokens
	// In production, this should be loaded from environment or configuration
	JWTSecretKey []byte

	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid email or password")
	
	// ErrUnauthorized is returned when a user is not authorized
	ErrUnauthorized = errors.New("unauthorized access")
	
	// ErrForbidden is returned when a user is forbidden from accessing a resource
	ErrForbidden = errors.New("forbidden: insufficient permissions")
)

// InitJWTSecret initializes the JWT secret key
func InitJWTSecret() error {
	// Generate a random 32-byte key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}
	
	JWTSecretKey = key
	return nil
}

// GeneratePasswordHash generates a bcrypt hash from a password
func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	
	return string(hash), nil
}

// VerifyPassword checks if a password matches a hash
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(user *User) (string, error) {
	expiration := time.Now().Add(jwtExpiration)
	
	claims := jwt.MapClaims{
		"sub":  strconv.FormatInt(user.ID, 10),
		"name": user.Name,
		"email": user.Email,
		"role": user.Role,
		"exp":  expiration.Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// VerifyJWT verifies and parses a JWT token
func VerifyJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		return JWTSecretKey, nil
	})
	
	if err != nil {
		return nil, nil, err
	}
	
	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}
	
	return token, claims, nil
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	
	return parts[1], nil
}

// Authentication and authorization middleware

// requireAuth is middleware that wraps an apiFunc to check if a user is authenticated
func requireAuth(next apiFunc) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := extractTokenFromHeader(r)
		if err != nil {
			return ErrUnauthorized
		}
		
		_, claims, err := VerifyJWT(token)
		if err != nil {
			return ErrUnauthorized
		}
		
		// Add claims to request context for handlers to use
		ctx := setUserContext(r.Context(), claims)
		return next(w, r.WithContext(ctx))
	}
}

// requireAdmin is middleware that wraps an apiFunc to check if a user is an admin
func requireAdmin(next apiFunc) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := extractTokenFromHeader(r)
		if err != nil {
			return ErrUnauthorized
		}
		
		_, claims, err := VerifyJWT(token)
		if err != nil {
			return ErrUnauthorized
		}
		
		// Check if the user is an admin
		role, ok := claims["role"].(string)
		if !ok || role != RoleAdmin {
			return ErrForbidden
		}
		
		// Add claims to request context for handlers to use
		ctx := setUserContext(r.Context(), claims)
		return next(w, r.WithContext(ctx))
	}
}

// API handlers for authentication and user management

// handleLogin handles user login authentication
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	
	// Validate input
	if req.Email == "" || req.Password == "" {
		return fmt.Errorf("email and password required")
	}
	
	// Get user by email
	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		// Use same error message regardless of whether user exists to prevent enumeration
		return ErrInvalidCredentials
	}
	
	// Verify password
	if !VerifyPassword(req.Password, user.PasswordHash) {
		return ErrInvalidCredentials
	}
	
	// Generate JWT token
	token, err := GenerateJWT(user)
	if err != nil {
		return err
	}
	
	// Update last login time
	user.LastLogin = time.Now()
	if err := s.store.UpdateUser(user); err != nil {
		return err
	}
	
	// Mask password hash before sending response
	user.PasswordHash = ""
	
	resp := LoginResponse{
		User:  *user,
		Token: token,
	}
	
	return WriteJson(w, http.StatusOK, resp)
}

// handleCreateUser handles creating a new user (admin only)
func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	var req CreateUserRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}
	
	// Validate input
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return ErrInvalidData
	}
	
	// Validate role
	if req.Role != RoleUser && req.Role != RoleAdmin {
		req.Role = RoleUser // Default to user role
	}
	
	// Generate password hash
	passwordHash, err := GeneratePasswordHash(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create user
	user := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
		IsActive:     req.IsActive,
		CreatedAt:    time.Now(),
	}
	
	// Attempt to create user (CreateUser already handles duplicate email check)
	if err := s.store.CreateUser(user); err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return ErrDuplicateEmail
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	resp := CreateUserResponse{
		ID:      user.ID,
		Success: true,
		Message: "User created successfully",
	}
	
	return WriteJson(w, http.StatusCreated, resp)
}

// handleGetUsers handles getting a list of users (admin only)
func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	// Parse pagination params
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // default
	}
	
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // default
	}
	
	users, err := s.store.ListUsers(limit, offset)
	if err != nil {
		return err
	}
	
	// Remove password hashes
	for _, user := range users {
		user.PasswordHash = ""
	}
	
	return WriteJson(w, http.StatusOK, users)
}

// handleGetUser handles getting a single user
func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}
	
	user, err := s.store.GetUserByID(id)
	if err != nil {
		return err
	}
	
	// Clear password hash before returning
	user.PasswordHash = ""
	
	return WriteJson(w, http.StatusOK, user)
}

// handleUpdateUser handles updating a user (admin only)
func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}
	
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	
	// Get existing user
	user, err := s.store.GetUserByID(id)
	if err != nil {
		return err
	}
	
	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	
	if req.Email != "" && req.Email != user.Email {
		// Check if email is already in use
		existingUser, err := s.store.GetUserByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			return fmt.Errorf("email already in use")
		}
		user.Email = req.Email
	}
	
	// Update password if provided
	if req.Password != "" {
		passwordHash, err := GeneratePasswordHash(req.Password)
		if err != nil {
			return err
		}
		user.PasswordHash = passwordHash
	}
	
	// Update role if provided
	if req.Role == RoleUser || req.Role == RoleAdmin {
		user.Role = req.Role
	}
	
	// Update active status
	user.IsActive = req.IsActive
	
	if err := s.store.UpdateUser(user); err != nil {
		return err
	}
	
	resp := CreateUserResponse{
		ID:      user.ID,
		Success: true,
		Message: "User updated successfully",
	}
	
	return WriteJson(w, http.StatusOK, resp)
}

// handleDeleteUser handles deleting a user (admin only)
func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}
	
	if err := s.store.DeleteUser(id); err != nil {
		return err
	}
	
	resp := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "User deleted successfully",
	}
	
	return WriteJson(w, http.StatusOK, resp)
}