package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	
	"expertdb/internal/auth"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// UserHandler handles user-related API endpoints
type UserHandler struct {
	store storage.Storage
}

// NewUserHandler creates a new user handler
func NewUserHandler(store storage.Storage) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// HandleCreateUser creates a new user account (admin only)
func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Parse request body into CreateUserRequest
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("Failed to parse user creation request: %v", err)
		return fmt.Errorf("invalid request body: %w", err)
	}
	
	// Validate required fields
	if req.Email == "" || req.Name == "" || req.Password == "" {
		log.Debug("User creation failed - missing required fields")
		return domain.ErrValidation
	}
	
	// Validate role against allowed values
	validRoles := []string{auth.RoleSuperUser, auth.RoleAdmin, auth.RolePlanner, auth.RoleUser}
	isValidRole := false
	for _, role := range validRoles {
		if req.Role == role {
			isValidRole = true
			break
		}
	}
	
	if !isValidRole {
		// Default to regular user role for security
		log.Debug("Invalid role specified (%s), defaulting to user role", req.Role)
		req.Role = auth.RoleUser
	}
	
	// Get creator's role from context
	creatorRole, ok := auth.GetUserRoleFromContext(r.Context())
	if !ok {
		log.Error("Failed to get creator role from context")
		return fmt.Errorf("authentication error: missing user role")
	}
	
	// Check if creator can create a user with the requested role
	if !auth.CanManageRole(creatorRole, req.Role) {
		log.Warn("User with role %s attempted to create user with role %s", creatorRole, req.Role)
		return fmt.Errorf("insufficient privileges to create user with role '%s'", req.Role)
	}
	
	// Validate email format
	if !strings.Contains(req.Email, "@") || len(req.Email) < 5 {
		log.Debug("User creation failed - invalid email format: %s", req.Email)
		return fmt.Errorf("invalid email format")
	}
	
	// Generate secure password hash
	passwordHash, err := auth.GeneratePasswordHash(req.Password)
	if err != nil {
		log.Error("Failed to hash password for new user: %v", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create user object
	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
		IsActive:     req.IsActive,
		CreatedAt:    time.Now(),
	}
	
	// Attempt to create user in database with role check
	id, err := h.store.CreateUserWithRoleCheck(user, creatorRole)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Info("User creation failed - duplicate email: %s", req.Email)
			return fmt.Errorf("email already exists")
		}
		if strings.Contains(err.Error(), "cannot create") {
			log.Warn("Role-based access control prevented user creation: %v", err)
			return fmt.Errorf("insufficient privileges: %s", err.Error())
		}
		log.Error("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	// Prepare success response
	resp := domain.CreateUserResponse{
		ID:      id,
		Success: true,
		Message: "User created successfully",
	}
	
	// Log successful user creation
	log.Info("New user created: %s (ID: %d, Role: %s)", req.Email, id, req.Role)
	
	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(resp)
}

// HandleGetUsers retrieves a paginated list of users (admin only)
func (h *UserHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Parse pagination parameters
	const DefaultLimit = 10
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = DefaultLimit // default page size
	}
	
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // default starting position
	}
	
	// Retrieve users from database
	users, err := h.store.ListUsers(limit, offset)
	if err != nil {
		log.Error("Failed to list users: %v", err)
		return fmt.Errorf("failed to retrieve users: %w", err)
	}
	
	// Remove sensitive data (password hashes) for security
	for _, user := range users {
		user.PasswordHash = ""
	}
	
	// Log user list request
	log.Debug("User list retrieved: %d users returned (limit: %d, offset: %d)", 
		len(users), limit, offset)
	
	// Return user list as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(users)
}

// HandleGetUser retrieves details for a specific user
func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate the user ID from the URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Debug("Invalid user ID format: %s", idStr)
		return fmt.Errorf("invalid user ID: must be a number")
	}
	
	// Retrieve the user from the database
	user, err := h.store.GetUser(id)
	if err != nil {
		if err == domain.ErrNotFound {
			log.Info("User not found with ID %d", id)
			return domain.ErrNotFound
		}
		log.Error("Failed to get user with ID %d: %v", id, err)
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// Remove sensitive data before returning the response
	user.PasswordHash = ""
	
	// Log user retrieval
	log.Debug("User retrieved: ID %d", id)
	
	// Return user details as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(user)
}

// HandleUpdateUser updates an existing user's information
func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate the user ID from the URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Debug("Invalid user ID format in update request: %s", idStr)
		return fmt.Errorf("invalid user ID: must be a number")
	}
	
	// Parse the update request from JSON body
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("Failed to parse user update request: %v", err)
		return fmt.Errorf("invalid request format: %w", err)
	}
	
	// Retrieve the existing user from database
	user, err := h.store.GetUser(id)
	if err != nil {
		if err == domain.ErrNotFound {
			log.Info("User not found with ID %d for update", id)
			return domain.ErrNotFound
		}
		log.Error("Failed to get user with ID %d for update: %v", id, err)
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// Update fields selectively (only if provided in request)
	
	// Update name if provided
	if req.Name != "" {
		log.Debug("Updating name for user ID %d: %s -> %s", id, user.Name, req.Name)
		user.Name = req.Name
	}
	
	// Update email if provided and different from current
	if req.Email != "" && req.Email != user.Email {
		// Validate email uniqueness
		existingUser, err := h.store.GetUserByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			log.Info("Email already in use during user update: %s", req.Email)
			return fmt.Errorf("email already in use by another user")
		}
		
		// Email is unique, update it
		log.Debug("Updating email for user ID %d: %s -> %s", id, user.Email, req.Email)
		user.Email = req.Email
	}
	
	// Update password if provided
	if req.Password != "" {
		log.Debug("Updating password for user ID %d", id)
		passwordHash, err := auth.GeneratePasswordHash(req.Password)
		if err != nil {
			log.Error("Failed to hash password during user update: %v", err)
			return fmt.Errorf("failed to update password: %w", err)
		}
		user.PasswordHash = passwordHash
	}
	
	// Update role if provided and valid
	if req.Role == auth.RoleUser || req.Role == auth.RoleAdmin {
		log.Debug("Updating role for user ID %d: %s -> %s", id, user.Role, req.Role)
		user.Role = req.Role
	}
	
	// Update active status
	if user.IsActive != req.IsActive {
		log.Debug("Updating active status for user ID %d: %v -> %v", id, user.IsActive, req.IsActive)
		user.IsActive = req.IsActive
	}
	
	// Save the updated user to the database
	if err := h.store.UpdateUser(user); err != nil {
		log.Error("Failed to update user ID %d: %v", id, err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	// Prepare and return success response
	resp := domain.CreateUserResponse{
		ID:      user.ID,
		Success: true,
		Message: "User updated successfully",
	}
	
	// Log successful update
	log.Info("User updated successfully: ID %d, Email: %s, Role: %s", user.ID, user.Email, user.Role)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

// HandleDeleteUser permanently deletes a user account
func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate the user ID from the URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Debug("Invalid user ID format in delete request: %s", idStr)
		return fmt.Errorf("invalid user ID: must be a number")
	}
	
	// Get user details for logging before deletion
	user, err := h.store.GetUser(id)
	if err != nil {
		if err == domain.ErrNotFound {
			log.Info("User not found with ID %d for deletion", id)
			return domain.ErrNotFound
		}
		log.Error("Failed to get user with ID %d for deletion: %v", id, err)
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// Safety check: don't delete the last admin account
	if user.Role == auth.RoleAdmin {
		// Count total admins
		adminCount := 0
		users, err := h.store.ListUsers(100, 0) // Get up to 100 users
		if err == nil {
			for _, u := range users {
				if u.Role == auth.RoleAdmin && u.IsActive {
					adminCount++
				}
			}
		}
		
		if adminCount <= 1 {
			log.Warn("Attempt to delete the last admin account (ID: %d)", id)
			return fmt.Errorf("cannot delete the last admin account")
		}
	}
	
	// Delete the user from the database
	if err := h.store.DeleteUser(id); err != nil {
		log.Error("Failed to delete user ID %d: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	// Log successful deletion
	log.Info("User deleted: ID %d, Email: %s, Role: %s", id, user.Email, user.Role)
	
	// Prepare and return success response
	resp := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "User deleted successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}