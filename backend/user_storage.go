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

	query := `
		INSERT INTO users (
			name, email, password_hash, role, is_active, created_at, last_login
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		user.Name, user.Email, user.PasswordHash, user.Role,
		user.IsActive, user.CreatedAt, user.LastLogin,
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
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
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
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
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
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user
func (s *SQLiteStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, password_hash = ?, role = ?, 
			is_active = ?, last_login = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		user.Name, user.Email, user.PasswordHash, user.Role,
		user.IsActive, user.LastLogin, user.ID,
	)

	return err
}

// DeleteUser deletes a user by ID
func (s *SQLiteStore) DeleteUser(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := s.db.Exec(query, id)
	return err
}

// CountUsers returns the total number of users
func (s *SQLiteStore) CountUsers() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM users"
	err := s.db.QueryRow(query).Scan(&count)
	return count, err
}

// EnsureAdminExists checks if an admin user exists and creates one if not
func (s *SQLiteStore) EnsureAdminExists(adminEmail, adminName, adminPasswordHash string) error {
	// Check if any admin exists
	query := "SELECT COUNT(*) FROM users WHERE role = 'admin'"
	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return err
	}

	// If no admin exists, create one
	if count == 0 {
		now := time.Now()
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