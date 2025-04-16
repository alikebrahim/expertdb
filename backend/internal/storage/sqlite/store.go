// Package sqlite provides a SQLite implementation of the storage interface
package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	
	_ "github.com/mattn/go-sqlite3"
	"expertdb/internal/storage"
)

// SQLiteStore implements the Storage interface with SQLite backend
type SQLiteStore struct {
	db *sql.DB
}

// Verify that SQLiteStore implements the Storage interface at compile time
var _ storage.Storage = (*SQLiteStore)(nil)

// New creates a new SQLite database connection
func New(dbPath string) (*SQLiteStore, error) {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}
	
	// Connect to the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	// Create the store
	store := &SQLiteStore{
		db: db,
	}
	
	return store, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// InitDB initializes the database schema if it doesn't exist
func (s *SQLiteStore) InitDB() error {
	// Check if tables exist by querying a key table
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='experts'").Scan(&count)
	if err != nil || count == 0 {
		return fmt.Errorf("database schema not initialized: %w", err)
	}
	
	return nil
}