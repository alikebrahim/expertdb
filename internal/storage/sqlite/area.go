package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	
	"expertdb/internal/domain"
	"expertdb/internal/logger"
)

// ListAreas retrieves all expert areas
func (s *SQLiteStore) ListAreas() ([]*domain.Area, error) {
	query := "SELECT id, name FROM expert_areas ORDER BY name"
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expert areas: %w", err)
	}
	defer rows.Close()
	
	var areas []*domain.Area
	for rows.Next() {
		var area domain.Area
		if err := rows.Scan(&area.ID, &area.Name); err != nil {
			return nil, fmt.Errorf("failed to scan area row: %w", err)
		}
		areas = append(areas, &area)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating area rows: %w", err)
	}
	
	return areas, nil
}

// GetArea retrieves a specific area by its ID
func (s *SQLiteStore) GetArea(id int64) (*domain.Area, error) {
	var area domain.Area
	err := s.db.QueryRow("SELECT id, name FROM expert_areas WHERE id = ?", id).Scan(&area.ID, &area.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get expert area: %w", err)
	}
	return &area, nil
}

// CreateArea creates a new expert area with the given name
func (s *SQLiteStore) CreateArea(name string) (int64, error) {
	log := logger.Get()
	
	// Check for empty name
	if strings.TrimSpace(name) == "" {
		return 0, fmt.Errorf("area name cannot be empty")
	}
	
	// Check if area with same name already exists (case insensitive)
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM expert_areas WHERE LOWER(name) = LOWER(?)", name).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to check for duplicate area: %w", err)
	}
	
	if count > 0 {
		return 0, fmt.Errorf("area with name '%s' already exists", name)
	}
	
	// Insert new area
	result, err := s.db.Exec("INSERT INTO expert_areas (name) VALUES (?)", name)
	if err != nil {
		log.Error("Failed to insert new area: %v", err)
		return 0, fmt.Errorf("failed to create expert area: %w", err)
	}
	
	// Get the ID of the new area
	id, err := result.LastInsertId()
	if err != nil {
		log.Error("Failed to get last insert ID: %v", err)
		return 0, fmt.Errorf("failed to retrieve new area ID: %w", err)
	}
	
	log.Info("Created new expert area: %s (ID: %d)", name, id)
	return id, nil
}

// UpdateArea renames an existing expert area
func (s *SQLiteStore) UpdateArea(id int64, name string) error {
	log := logger.Get()
	
	// Check for empty name
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("area name cannot be empty")
	}
	
	// Check if area exists
	var existing domain.Area
	err := s.db.QueryRow("SELECT id, name FROM expert_areas WHERE id = ?", id).Scan(&existing.ID, &existing.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to check if area exists: %w", err)
	}
	
	// Check if another area with the same name already exists (case insensitive)
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM expert_areas WHERE LOWER(name) = LOWER(?) AND id != ?", name, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate area: %w", err)
	}
	
	if count > 0 {
		return fmt.Errorf("another area with name '%s' already exists", name)
	}
	
	// Use a transaction to update area name and cascade changes
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Rename the area
	_, err = tx.Exec("UPDATE expert_areas SET name = ? WHERE id = ?", name, id)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to update area name: %v", err)
		return fmt.Errorf("failed to update area: %w", err)
	}
	
	// Update experts and expert_requests that reference this area
	// Note: We don't update the database directly here as the foreign key is just
	// referencing the ID. The updated name will be fetched when needed.
	
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error("Failed to commit area update transaction: %v", err)
		return fmt.Errorf("failed to commit area update: %w", err)
	}
	
	log.Info("Updated expert area ID %d from '%s' to '%s'", id, existing.Name, name)
	return nil
}