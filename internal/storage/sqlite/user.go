package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	
	"expertdb/internal/domain"
)

// CreateUser creates a new user in the database with role management rules
func (s *SQLiteStore) CreateUser(user *domain.User) (int64, error) {
	// Validate required fields
	if user.Name == "" || user.Email == "" || user.PasswordHash == "" {
		return 0, domain.ErrValidation
	}

	// Check if email already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("failed to check for existing email: %w", err)
	}
	
	if exists {
		return 0, errors.New("email already exists")
	}

	// Validate role assignment based on the current user's role (if available)
	if user.Role == "" {
		// Default to regular user role
		user.Role = "user"
	} else {
		// Extract creator user ID and role from context (if available)
		// This logic assumes the user creator context is stored in the current request
		// In a real implementation, this would be passed through an additional param
		
		// Apply role creation constraints:
		// 1. If creating a super_user, verify the user is self-bootstrapping (like server initialization)
		// 2. For creating admin and below, verify creator is a super_user
		// 3. For creating planner and below, verify creator is admin or super_user
		
		// For the initial implementation, we'll simply check if the role is valid
		validRoles := []string{"super_user", "admin", "planner", "user"}
		roleValid := false
		for _, role := range validRoles {
			if user.Role == role {
				roleValid = true
				break
			}
		}
		
		if !roleValid {
			return 0, fmt.Errorf("invalid role: %s", user.Role)
		}
	}

	// Initialize with current time if not set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	
	// Handle NULL last_login
	var lastLogin interface{} = nil
	if !user.LastLogin.IsZero() {
		lastLogin = user.LastLogin
	}

	query := `
		INSERT INTO users (
			name, email, password_hash, role, is_active, created_at, last_login
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		user.Name, user.Email, user.PasswordHash, user.Role,
		user.IsActive, user.CreatedAt, lastLogin,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	user.ID = id
	return id, nil
}

// GetUser retrieves a user by ID
func (s *SQLiteStore) GetUser(id int64) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, last_login
		FROM users
		WHERE id = ?
	`

	var user domain.User
	var nullableLastLogin sql.NullTime
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &nullableLastLogin,
	)
	
	if nullableLastLogin.Valid {
		user.LastLogin = nullableLastLogin.Time
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *SQLiteStore) GetUserByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, last_login
		FROM users
		WHERE email = ?
	`

	var user domain.User
	var nullableLastLogin sql.NullTime
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &nullableLastLogin,
	)
	
	if nullableLastLogin.Valid {
		user.LastLogin = nullableLastLogin.Time
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// ListUsers retrieves all users with pagination
func (s *SQLiteStore) ListUsers(limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, last_login
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var nullableLastLogin sql.NullTime
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Role, &user.IsActive, &user.CreatedAt, &nullableLastLogin,
		)
		
		if nullableLastLogin.Valid {
			user.LastLogin = nullableLastLogin.Time
		}
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (s *SQLiteStore) UpdateUser(user *domain.User) error {
	// First get the current user to preserve values that aren't explicitly changed
	current, err := s.GetUser(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get current user state: %w", err)
	}
	
	// Preserve existing values for empty fields
	if user.Name == "" {
		user.Name = current.Name
	}
	
	if user.Email == "" {
		user.Email = current.Email
	}
	
	if user.PasswordHash == "" {
		user.PasswordHash = current.PasswordHash
	}
	
	if user.Role == "" {
		user.Role = current.Role
	}
	
	// Handle NULL for last_login
	var lastLogin interface{} = nil
	if !user.LastLogin.IsZero() {
		lastLogin = user.LastLogin
	} else if !current.LastLogin.IsZero() {
		// Preserve the existing last login time
		lastLogin = current.LastLogin
	}

	query := `
		UPDATE users
		SET name = ?, email = ?, password_hash = ?, role = ?, 
			is_active = ?, last_login = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(
		query,
		user.Name, user.Email, user.PasswordHash, user.Role,
		user.IsActive, lastLogin, user.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	// Verify the update affected a row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for user update: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// UpdateUserLastLogin updates the last login timestamp for a user
func (s *SQLiteStore) UpdateUserLastLogin(id int64) error {
	query := "UPDATE users SET last_login = ? WHERE id = ?"
	result, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update user last login: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for last login update: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	return nil
}

// EnsureSuperUserExists ensures that a super user with the given credentials exists
func (s *SQLiteStore) EnsureSuperUserExists(email, name, passwordHash string) error {
	// Check if super user with given email already exists
	user, err := s.GetUserByEmail(email)
	if err == nil && user != nil {
		// User already exists - check if role needs to be upgraded to super_user
		if user.Role != "super_user" {
			// Update the existing user to a super_user
			user.Role = "super_user"
			err = s.UpdateUser(user)
			if err != nil {
				return fmt.Errorf("failed to upgrade user to super_user: %w", err)
			}
		}
		return nil
	}
	
	// Super user doesn't exist, create a new one
	superUser := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "super_user",
		IsActive:     true,
		CreatedAt:    time.Now(),
		LastLogin:    time.Now(),
	}
	
	_, err = s.CreateUser(superUser)
	if err != nil {
		return fmt.Errorf("failed to create super user: %w", err)
	}
	
	return nil
}

// CreateUserWithRoleCheck creates a new user while enforcing role hierarchy rules
func (s *SQLiteStore) CreateUserWithRoleCheck(user *domain.User, creatorRole string) (int64, error) {
	// First, check if the user can be created by the creator based on role hierarchy
	if !canCreateUserWithRole(creatorRole, user.Role) {
		return 0, fmt.Errorf("creator with role '%s' cannot create a user with role '%s'", 
			creatorRole, user.Role)
	}
	
	// If allowed, proceed with regular user creation
	return s.CreateUser(user)
}

// Helper function to check if a user with a given role can create a user with a target role
func canCreateUserWithRole(creatorRole, targetRole string) bool {
	switch creatorRole {
	case "super_user":
		// Super user can create admin, planner, and regular users
		return targetRole == "admin" || targetRole == "planner" || targetRole == "user"
	case "admin":
		// Admin can create planner and regular users
		return targetRole == "planner" || targetRole == "user"
	default:
		// No other roles can create users
		return false
	}
}

// DeleteUser deletes a user by ID
func (s *SQLiteStore) DeleteUser(id int64) error {
	// First check if this is a protected user (like the first super_user)
	var role string
	err := s.db.QueryRow("SELECT role FROM users WHERE id = ?", id).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to check user role: %w", err)
	}
	
	// If deleting a super_user, ensure it's not the last one
	if role == "super_user" {
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'super_user'").Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to count super users: %w", err)
		}
		
		if count <= 1 {
			return fmt.Errorf("cannot delete the last super user")
		}
	}
	
	// If deleting a planner, handle cascade deletion of planner assignments
	if role == "planner" {
		// In a real implementation, you would delete/reassign any planner assignments here
		// For example:
		_, err := s.db.Exec("UPDATE phases SET assigned_planner_id = NULL WHERE assigned_planner_id = ?", id)
		if err != nil {
		    return fmt.Errorf("failed to clear planner assignments: %w", err)
		}
	}
	
	// Now proceed with deleting the user
	result, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	return nil
}