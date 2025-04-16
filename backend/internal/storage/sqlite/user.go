package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	
	"expertdb/internal/domain"
)

// CreateUser creates a new user in the database
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
	result, err := s.db.Exec(query, time.Now().UTC(), id)
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

// EnsureAdminExists ensures that an admin user with the given credentials exists
func (s *SQLiteStore) EnsureAdminExists(adminEmail, adminName, adminPasswordHash string) error {
	// Check if admin with given email already exists
	user, err := s.GetUserByEmail(adminEmail)
	if err == nil && user != nil {
		// Admin already exists, no need to create
		return nil
	}
	
	// Admin doesn't exist, create a new one
	admin := &domain.User{
		Name:         adminName,
		Email:        adminEmail,
		PasswordHash: adminPasswordHash,
		Role:         "admin",
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		LastLogin:    time.Now().UTC(),
	}
	
	_, err = s.CreateUser(admin)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	
	return nil
}