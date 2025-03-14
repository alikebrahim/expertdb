// Package main provides the backend functionality for the ExpertDB application
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

// Authentication related errors
var (
	// JWTSecretKey is the key used to sign and verify JWT tokens
	// In production, this should be loaded from environment or configuration
	// This key is generated randomly at application startup
	JWTSecretKey []byte

	// ErrInvalidCredentials is returned when login credentials are invalid
	// This error is intentionally generic to prevent user enumeration
	ErrInvalidCredentials = errors.New("invalid email or password")
	
	// ErrUnauthorized is returned when a user is not authenticated
	// This indicates missing or invalid JWT token in the request
	ErrUnauthorized = errors.New("unauthorized access")
	
	// ErrForbidden is returned when an authenticated user lacks sufficient permissions
	// This indicates the user doesn't have the required role (typically admin) for an operation
	ErrForbidden = errors.New("forbidden: insufficient permissions")
)

// InitJWTSecret initializes the JWT secret key used for token signing and verification
// 
// It generates a secure random 32-byte key that will be used for the lifetime of the application.
// A new key is generated each time the application restarts, invalidating previous tokens.
// 
// Returns:
//   - error: If random number generation fails
func InitJWTSecret() error {
	// Step 1: Generate a cryptographically secure random 32-byte key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return fmt.Errorf("failed to generate random JWT secret: %w", err)
	}
	
	// Step 2: Store the key in the global variable
	JWTSecretKey = key
	
	return nil
}

// GeneratePasswordHash creates a secure bcrypt hash from a plaintext password
//
// The generated hash includes the salt and cost factor, making it self-contained
// for future password verification.
//
// Inputs:
//   - password (string): The plaintext password to hash
//
// Returns:
//   - string: The password hash as a string
//   - error: If hashing fails due to algorithm constraints or system issues
func GeneratePasswordHash(password string) (string, error) {
	// Step 1: Hash the password using bcrypt with the configured cost factor
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Step 2: Return the hash as a string, suitable for database storage
	return string(hash), nil
}

// VerifyPassword checks if a plaintext password matches a previously hashed password
//
// This function uses bcrypt's comparison function which is resistant to timing attacks.
//
// Inputs:
//   - password (string): The plaintext password to check
//   - hash (string): The stored password hash to compare against
//
// Returns:
//   - bool: True if the password matches the hash, false otherwise
func VerifyPassword(password, hash string) bool {
	// Compare the provided password against the hash
	// bcrypt.CompareHashAndPassword handles all the work of extracting the salt and cost
	// parameters from the hash and performing the comparison securely
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user with standard claims
//
// The token includes user identity and role information as claims,
// along with an expiration time to ensure tokens don't remain valid indefinitely.
//
// Inputs:
//   - user (*User): The user for whom to generate the token
//
// Returns:
//   - string: The signed JWT token string
//   - error: If token signing fails
func GenerateJWT(user *User) (string, error) {
	// Step 1: Calculate token expiration time
	expiration := time.Now().Add(jwtExpiration)
	
	// Step 2: Create claims map with user information
	claims := jwt.MapClaims{
		"sub":   strconv.FormatInt(user.ID, 10), // Subject (user ID)
		"name":  user.Name,                      // User's full name
		"email": user.Email,                     // User's email address
		"role":  user.Role,                      // User's role (admin/user)
		"exp":   expiration.Unix(),              // Expiration timestamp
	}
	
	// Step 3: Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Step 4: Sign the token with the secret key
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}
	
	// Step 5: Return the signed token string
	return tokenString, nil
}

// VerifyJWT verifies and parses a JWT token
//
// This function validates the token signature, ensures it was created with the expected
// signing method, and checks that it hasn't expired.
//
// Inputs:
//   - tokenString (string): The JWT token string to verify
//
// Returns:
//   - *jwt.Token: The parsed token object
//   - jwt.MapClaims: The token claims if valid
//   - error: If token validation fails
func VerifyJWT(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	// Step 1: Initialize an empty claims map
	claims := jwt.MapClaims{}
	
	// Step 2: Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Step 2.1: Validate the signing algorithm is as expected (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		// Step 2.2: Return the secret key for signature verification
		return JWTSecretKey, nil
	})
	
	// Step 3: Handle parsing errors (includes expiration checks)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}
	
	// Step 4: Verify token is valid (this check is redundant as ParseWithClaims
	// already checks validity, but kept for clarity and safety)
	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}
	
	// Step 5: Return the validated token and claims
	return token, claims, nil
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
//
// This helper function parses the Authorization header and expects a "Bearer" token.
//
// Inputs:
//   - r (*http.Request): The HTTP request containing the Authorization header
//
// Returns:
//   - string: The extracted token
//   - error: If the header is missing or malformed
func extractTokenFromHeader(r *http.Request) (string, error) {
	// Step 1: Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	
	// Step 2: Split the header into parts and validate format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	
	// Step 3: Return the token part
	return parts[1], nil
}

// Authentication and authorization middleware
// These functions provide role-based access control for API endpoints

// requireAuth is middleware that verifies a user is authenticated before allowing access
//
// This middleware wraps an API handler function and ensures the request contains a valid
// JWT token. If valid, the user claims are added to the request context for use in handlers.
//
// Inputs:
//   - next (apiFunc): The handler function to wrap with authentication
//
// Returns:
//   - apiFunc: A wrapped handler function that first checks authentication
func requireAuth(next apiFunc) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Step 1: Extract the token from the Authorization header
		token, err := extractTokenFromHeader(r)
		if err != nil {
			logger := GetLogger()
			logger.Debug("Authentication failed: %v", err)
			return ErrUnauthorized
		}
		
		// Step 2: Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			logger := GetLogger()
			logger.Debug("JWT verification failed: %v", err)
			return ErrUnauthorized
		}
		
		// Step 3: Add user claims to request context for downstream handlers
		ctx := setUserContext(r.Context(), claims)
		
		// Step 4: Pass control to the next handler with the updated context
		return next(w, r.WithContext(ctx))
	}
}

// requireAdmin is middleware that ensures only admin users can access protected endpoints
//
// This middleware extends the basic authentication check by also verifying the user
// has the admin role. It's used for administrative endpoints that modify system data.
//
// Inputs:
//   - next (apiFunc): The handler function to wrap with admin authorization
//
// Returns:
//   - apiFunc: A wrapped handler function that first checks admin authorization
func requireAdmin(next apiFunc) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := GetLogger()
		
		// Step 1: Extract the token from the Authorization header
		token, err := extractTokenFromHeader(r)
		if err != nil {
			logger.Debug("Admin check failed - missing token: %v", err)
			return ErrUnauthorized
		}
		
		// Step 2: Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			logger.Debug("Admin check failed - invalid token: %v", err)
			return ErrUnauthorized
		}
		
		// Step 3: Check if the user has admin role
		role, ok := claims["role"].(string)
		if !ok || role != RoleAdmin {
			// User is authenticated but not an admin
			logger.Info("Forbidden access attempt by non-admin user (ID: %v) to %s", 
				claims["sub"], r.URL.Path)
			return ErrForbidden
		}
		
		// Step 4: Add user claims to request context for downstream handlers
		ctx := setUserContext(r.Context(), claims)
		
		// Step 5: Pass control to the next handler with the updated context
		return next(w, r.WithContext(ctx))
	}
}

// API handlers for authentication and user management
// These handlers process user authentication and account management operations

// handleLogin processes user login requests and issues JWT tokens
//
// This handler authenticates users based on email and password, generating a JWT token
// for successful logins that can be used for subsequent authenticated requests.
//
// HTTP Method: POST
// Endpoint: /api/auth/login
// 
// Inputs:
//   - w (http.ResponseWriter): Response writer for sending the JSON response
//   - r (*http.Request): Request containing login credentials in JSON body
//
// Flow:
//   1. Parse and validate login credentials from request body
//   2. Retrieve user from database by email
//   3. Verify password against stored hash
//   4. Generate and sign a JWT token
//   5. Update user's last login timestamp
//   6. Return user information and token
//
// Returns:
//   - error: An authentication error or internal error if the login process fails
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Parse login request from JSON body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug("Failed to parse login request: %v", err)
		return fmt.Errorf("invalid request format: %w", err)
	}
	
	// Step 2: Validate required fields
	if req.Email == "" || req.Password == "" {
		logger.Debug("Login attempt with missing credentials")
		return fmt.Errorf("email and password required")
	}
	
	// Step 3: Retrieve user by email
	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		// Use generic error message to prevent user enumeration
		logger.Info("Login failed - email not found: %s", req.Email)
		return ErrInvalidCredentials
	}
	
	// Step 4: Verify password
	if !VerifyPassword(req.Password, user.PasswordHash) {
		logger.Info("Login failed - invalid password for user: %s", req.Email)
		return ErrInvalidCredentials
	}
	
	// Step 5: Check if user is active
	if !user.IsActive {
		logger.Info("Login denied - inactive account: %s", req.Email)
		return fmt.Errorf("account is inactive, please contact administrator")
	}
	
	// Step 6: Generate JWT token
	token, err := GenerateJWT(user)
	if err != nil {
		logger.Error("Failed to generate token for user %s: %v", req.Email, err)
		return fmt.Errorf("failed to generate auth token: %w", err)
	}
	
	// Step 7: Update last login time
	user.LastLogin = time.Now()
	if err := s.store.UpdateUser(user); err != nil {
		// Non-fatal error - log but continue
		logger.Warn("Failed to update last login time for user %s: %v", req.Email, err)
	}
	
	// Step 8: Prepare response (mask password hash for security)
	user.PasswordHash = ""
	resp := LoginResponse{
		User:  *user,
		Token: token,
	}
	
	// Log successful login
	logger.Info("User logged in successfully: %s (ID: %d, Role: %s)", 
		user.Email, user.ID, user.Role)
	
	// Step 9: Return success response
	return WriteJson(w, http.StatusOK, resp)
}

// handleCreateUser creates a new user account (admin only)
//
// This handler allows administrators to create new user accounts with specified
// permissions and attributes. New users can be created with either admin or regular
// user roles.
//
// HTTP Method: POST
// Endpoint: /api/users
// Access: Admin only (via requireAdmin middleware)
//
// Inputs:
//   - w (http.ResponseWriter): Response writer for sending JSON response
//   - r (*http.Request): Request containing user details in JSON body
//
// Flow:
//   1. Parse and validate user creation request
//   2. Validate role and required fields
//   3. Generate secure password hash
//   4. Create user record in database
//   5. Return success response with new user ID
//
// Returns:
//   - error: If validation fails or user creation fails
func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Parse request body into CreateUserRequest
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug("Failed to parse user creation request: %v", err)
		return fmt.Errorf("invalid request body: %w", err)
	}
	
	// Step 2: Validate required fields
	if req.Email == "" || req.Name == "" || req.Password == "" {
		logger.Debug("User creation failed - missing required fields")
		return ErrInvalidData
	}
	
	// Step 3: Validate and normalize role
	if req.Role != RoleUser && req.Role != RoleAdmin {
		// Default to regular user role for security
		logger.Debug("Invalid role specified (%s), defaulting to user role", req.Role)
		req.Role = RoleUser
	}
	
	// Step 4: Validate email format
	// NOTE: A more comprehensive validation could be extracted to a helper function
	if !strings.Contains(req.Email, "@") || len(req.Email) < 5 {
		logger.Debug("User creation failed - invalid email format: %s", req.Email)
		return fmt.Errorf("invalid email format")
	}
	
	// Step 5: Generate secure password hash
	passwordHash, err := GeneratePasswordHash(req.Password)
	if err != nil {
		logger.Error("Failed to hash password for new user: %v", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Step 6: Create user object
	user := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
		IsActive:     req.IsActive,
		CreatedAt:    time.Now(),
	}
	
	// Step 7: Attempt to create user in database
	if err := s.store.CreateUser(user); err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			logger.Info("User creation failed - duplicate email: %s", req.Email)
			return ErrDuplicateEmail
		}
		logger.Error("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	// Step 8: Prepare success response
	resp := CreateUserResponse{
		ID:      user.ID,
		Success: true,
		Message: "User created successfully",
	}
	
	// Log successful user creation
	logger.Info("New user created: %s (ID: %d, Role: %s)", req.Email, user.ID, req.Role)
	
	// Step 9: Return success response
	return WriteJson(w, http.StatusCreated, resp)
}

// handleGetUsers retrieves a paginated list of users (admin only)
//
// This handler returns a list of all users in the system with pagination support.
// For security, password hashes are removed from the response.
//
// HTTP Method: GET
// Endpoint: /api/users
// Access: Admin only (via requireAdmin middleware)
//
// Inputs:
//   - w (http.ResponseWriter): Response writer for sending JSON response
//   - r (*http.Request): Request containing optional pagination parameters
//
// Flow:
//   1. Parse pagination parameters from query string
//   2. Retrieve users from database with pagination
//   3. Remove sensitive data (password hashes)
//   4. Return user list as JSON response
//
// Returns:
//   - error: If database retrieval fails
func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Parse pagination parameters
	const DefaultLimit = 10
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = DefaultLimit // default page size
	}
	
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // default starting position
	}
	
	// Step 2: Retrieve users from database
	users, err := s.store.ListUsers(limit, offset)
	if err != nil {
		logger.Error("Failed to list users: %v", err)
		return fmt.Errorf("failed to retrieve users: %w", err)
	}
	
	// Step 3: Remove sensitive data (password hashes) for security
	for _, user := range users {
		user.PasswordHash = ""
	}
	
	// Log user list request
	logger.Debug("User list retrieved: %d users returned (limit: %d, offset: %d)", 
		len(users), limit, offset)
	
	// Step 4: Return user list as JSON response
	return WriteJson(w, http.StatusOK, users)
}

// handleGetUser retrieves details for a specific user
//
// This handler retrieves a single user by their ID. For security,
// the password hash is removed from the response.
//
// HTTP Method: GET
// Endpoint: /api/users/{id}
// Access: Any authenticated user (via requireAuth middleware)
//
// Inputs:
//   - w (http.ResponseWriter): Response writer for sending JSON response
//   - r (*http.Request): Request containing the user ID in the URL path
//
// Flow:
//   1. Extract and validate user ID from URL path
//   2. Retrieve user from database
//   3. Remove sensitive data (password hash)
//   4. Return user details as JSON response
//
// Returns:
//   - error: If user doesn't exist or ID is invalid
func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate the user ID from the URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Debug("Invalid user ID format: %s", idStr)
		return fmt.Errorf("invalid user ID: must be a number")
	}
	
	// Step 2: Retrieve the user from the database
	user, err := s.store.GetUserByID(id)
	if err != nil {
		logger.Info("Failed to get user with ID %d: %v", id, err)
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	// Step 3: Remove sensitive data before returning the response
	user.PasswordHash = ""
	
	// Log user retrieval
	logger.Debug("User retrieved: ID %d", id)
	
	// Step 4: Return user details as JSON response
	return WriteJson(w, http.StatusOK, user)
}

// handleUpdateUser updates an existing user's information (admin only)
//
// This handler allows administrators to update user details including name,
// email, password, role, and active status.
//
// HTTP Method: PUT
// Endpoint: /api/users/{id}
// Access: Admin only (via requireAdmin middleware)
//
// Inputs:
//   - w (http.ResponseWriter): Response writer for sending JSON response
//   - r (*http.Request): Request containing user ID in path and updated details in body
//
// Flow:
//   1. Extract and validate user ID from URL path
//   2. Parse update request from JSON body
//   3. Retrieve existing user from database
//   4. Update fields selectively based on request
//   5. Validate email uniqueness if changed
//   6. Hash new password if provided
//   7. Save updated user to database
//   8. Return success response
//
// Returns:
//   - error: If user doesn't exist, validation fails, or update fails
func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate the user ID from the URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Debug("Invalid user ID format in update request: %s", idStr)
		return fmt.Errorf("invalid user ID: must be a number")
	}
	
	// Step 2: Parse the update request from JSON body
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug("Failed to parse user update request: %v", err)
		return fmt.Errorf("invalid request format: %w", err)
	}
	
	// Step 3: Retrieve the existing user from database
	user, err := s.store.GetUserByID(id)
	if err != nil {
		logger.Info("Failed to get user with ID %d for update: %v", id, err)
		return fmt.Errorf("user not found: %w", err)
	}
	
	// Step 4: Update fields selectively (only if provided in request)
	
	// Update name if provided
	if req.Name != "" {
		logger.Debug("Updating name for user ID %d: %s -> %s", id, user.Name, req.Name)
		user.Name = req.Name
	}
	
	// Update email if provided and different from current
	if req.Email != "" && req.Email != user.Email {
		// Step 5: Validate email uniqueness
		existingUser, err := s.store.GetUserByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			logger.Info("Email already in use during user update: %s", req.Email)
			return fmt.Errorf("email already in use by another user")
		}
		
		// Email is unique, update it
		logger.Debug("Updating email for user ID %d: %s -> %s", id, user.Email, req.Email)
		user.Email = req.Email
	}
	
	// Step 6: Update password if provided
	if req.Password != "" {
		logger.Debug("Updating password for user ID %d", id)
		passwordHash, err := GeneratePasswordHash(req.Password)
		if err != nil {
			logger.Error("Failed to hash password during user update: %v", err)
			return fmt.Errorf("failed to update password: %w", err)
		}
		user.PasswordHash = passwordHash
	}
	
	// Update role if provided and valid
	if req.Role == RoleUser || req.Role == RoleAdmin {
		logger.Debug("Updating role for user ID %d: %s -> %s", id, user.Role, req.Role)
		user.Role = req.Role
	}
	
	// Update active status
	if user.IsActive != req.IsActive {
		logger.Debug("Updating active status for user ID %d: %v -> %v", id, user.IsActive, req.IsActive)
		user.IsActive = req.IsActive
	}
	
	// Step 7: Save the updated user to the database
	if err := s.store.UpdateUser(user); err != nil {
		logger.Error("Failed to update user ID %d: %v", id, err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	// Step 8: Prepare and return success response
	resp := CreateUserResponse{
		ID:      user.ID,
		Success: true,
		Message: "User updated successfully",
	}
	
	// Log successful update
	logger.Info("User updated successfully: ID %d, Email: %s, Role: %s", user.ID, user.Email, user.Role)
	
	return WriteJson(w, http.StatusOK, resp)
}

// handleDeleteUser permanently deletes a user account (admin only)
//
// This handler allows administrators to delete user accounts from the system.
// The operation is irreversible and removes all access for the deleted user.
//
// HTTP Method: DELETE
// Endpoint: /api/users/{id}
// Access: Admin only (via requireAdmin middleware)
//
// Inputs:
//   - w (http.ResponseWriter): Response writer for sending JSON response
//   - r (*http.Request): Request containing the user ID to delete in the URL path
//
// Flow:
//   1. Extract and validate user ID from URL path
//   2. Delete user from database
//   3. Return success response
//
// Returns:
//   - error: If user doesn't exist or deletion fails
func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate the user ID from the URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Debug("Invalid user ID format in delete request: %s", idStr)
		return fmt.Errorf("invalid user ID: must be a number")
	}
	
	// Get user details for logging before deletion
	user, err := s.store.GetUserByID(id)
	if err != nil {
		logger.Info("Failed to get user with ID %d for deletion: %v", id, err)
		return fmt.Errorf("user not found: %w", err)
	}
	
	// Safety check: don't delete the last admin account
	if user.Role == RoleAdmin {
		// Count total admins
		adminCount := 0
		users, err := s.store.ListUsers(100, 0) // Get up to 100 users
		if err == nil {
			for _, u := range users {
				if u.Role == RoleAdmin && u.IsActive {
					adminCount++
				}
			}
		}
		
		if adminCount <= 1 {
			logger.Warn("Attempt to delete the last admin account (ID: %d)", id)
			return fmt.Errorf("cannot delete the last admin account")
		}
	}
	
	// Step 2: Delete the user from the database
	if err := s.store.DeleteUser(id); err != nil {
		logger.Error("Failed to delete user ID %d: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	// Log successful deletion
	logger.Info("User deleted: ID %d, Email: %s, Role: %s", id, user.Email, user.Role)
	
	// Step 3: Prepare and return success response
	resp := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: true,
		Message: "User deleted successfully",
	}
	
	return WriteJson(w, http.StatusOK, resp)
}

// NOTE: Added safety check to prevent deletion of the last admin account, which would lock
// out all administrative access to the system. Consider adding additional confirmation
// mechanisms for destructive operations.