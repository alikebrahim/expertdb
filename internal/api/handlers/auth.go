package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"expertdb/internal/api/utils"
	"expertdb/internal/auth"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// AuthHandler handles authentication-related API endpoints
type AuthHandler struct {
	store storage.Storage
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(store storage.Storage) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

// HandleLogin processes user login requests and issues JWT tokens
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Parse login request from JSON body
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("Failed to parse login request: %v", err)
		return fmt.Errorf("invalid request format: %w", err)
	}
	
	// Validate required fields
	if req.Email == "" || req.Password == "" {
		log.Debug("Login attempt with missing credentials")
		return fmt.Errorf("email and password required")
	}
	
	// Retrieve user by email
	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		// Use generic error message to prevent user enumeration
		log.Info("Login failed - email not found: %s", req.Email)
		return domain.ErrInvalidCredentials
	}
	
	// Verify password
	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		log.Info("Login failed - invalid password for user: %s", req.Email)
		return domain.ErrInvalidCredentials
	}
	
	// Check if user is active
	if !user.IsActive {
		log.Info("Login denied - inactive account: %s", req.Email)
		return fmt.Errorf("account is inactive, please contact administrator")
	}
	
	// Generate JWT token
	token, err := auth.GenerateJWT(user)
	if err != nil {
		log.Error("Failed to generate token for user %s: %v", req.Email, err)
		return fmt.Errorf("failed to generate auth token: %w", err)
	}
	
	// Update last login time
	if err := h.store.UpdateUserLastLogin(user.ID); err != nil {
		// Non-fatal error - log but continue
		log.Warn("Failed to update last login time for user %s: %v", req.Email, err)
	}
	
	// Prepare response (mask password hash for security)
	user.PasswordHash = ""
	responseData := map[string]interface{}{
		"user":  user,
		"token": token,
	}
	
	// Log successful login
	log.Info("User logged in successfully: %s (ID: %d, Role: %s)", 
		user.Email, user.ID, user.Role)
	
	// Return standardized success response
	return utils.RespondWithSuccess(w, "Login successful", responseData)
}