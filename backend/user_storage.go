package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNotFound is returned when a record is not found
	ErrNotFound = errors.New("record not found")
	// ErrDuplicateEmail is returned when attempting to create a user with an email that already exists
	ErrDuplicateEmail = errors.New("email already exists")
	// ErrInvalidData is returned when attempting to create or update with invalid data
	ErrInvalidData = errors.New("invalid user data")
)

// CreateUser creates a new user in the database
func (s *SQLiteStore) CreateUser(user *User) error {
	// Validate required fields
	if user.Name == "" || user.Email == "" || user.PasswordHash == "" {
		return ErrInvalidData
	}

	// Check if email already exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for existing email: %w", err)
	}
	
	if exists {
		return ErrDuplicateEmail
	}

	// Initialize with current time if not set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
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
		return fmt.Errorf("failed to insert user: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	user.ID = id
	return nil
}

// GetUserByID retrieves a user by ID
func (s *SQLiteStore) GetUserByID(id int64) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, last_login
		FROM users
		WHERE id = ?
	`

	var user User
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
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *SQLiteStore) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active, created_at, last_login
		FROM users
		WHERE email = ?
	`

	var user User
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
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// ListUsers retrieves all users with pagination
func (s *SQLiteStore) ListUsers(limit, offset int) ([]*User, error) {
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

	var users []*User
	for rows.Next() {
		var user User
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
func (s *SQLiteStore) UpdateUser(user *User) error {
	// Handle NULL for last_login
	var lastLogin interface{} = nil
	if !user.LastLogin.IsZero() {
		lastLogin = user.LastLogin
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
		return ErrNotFound
	}

	return nil
}

// DeleteUser deletes a user by ID
func (s *SQLiteStore) DeleteUser(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	// Verify the delete affected a row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for user delete: %w", err)
	}
	
	if rowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}

// CountUsers returns the total number of users
func (s *SQLiteStore) CountUsers() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM users"
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// EnsureAdminExists checks if an admin user exists and creates one if not
func (s *SQLiteStore) EnsureAdminExists(adminEmail, adminName, adminPasswordHash string) error {
	// Check if any admin exists
	query := "SELECT COUNT(*) FROM users WHERE role = 'admin'"
	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for admin user: %w", err)
	}

	// If no admin exists, create one
	if count == 0 {
		now := time.Now().UTC()
		admin := &User{
			Name:         adminName,
			Email:        adminEmail,
			PasswordHash: adminPasswordHash,
			Role:         "admin",
			IsActive:     true,
			CreatedAt:    now,
			LastLogin:    now,
		}
		return s.CreateUser(admin)
	}

	return nil
}