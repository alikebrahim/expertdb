package sqlite

import (
	"expertdb/internal/domain"
	"fmt"
	"strings"
)

// ListSpecializedAreas retrieves all specialized areas
func (s *SQLiteStore) ListSpecializedAreas() ([]*domain.SpecializedArea, error) {
	query := `
		SELECT id, name, created_at
		FROM specialized_areas
		ORDER BY name ASC
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query specialized areas: %w", err)
	}
	defer rows.Close()
	
	var areas []*domain.SpecializedArea
	for rows.Next() {
		var area domain.SpecializedArea
		if err := rows.Scan(&area.ID, &area.Name, &area.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan specialized area: %w", err)
		}
		areas = append(areas, &area)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating specialized areas: %w", err)
	}
	
	return areas, nil
}

// GetSpecializedAreasByIds retrieves specialized areas by their IDs
func (s *SQLiteStore) GetSpecializedAreasByIds(ids []int64) ([]*domain.SpecializedArea, error) {
	if len(ids) == 0 {
		return []*domain.SpecializedArea{}, nil
	}
	
	// Create placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	
	query := fmt.Sprintf(`
		SELECT id, name, created_at
		FROM specialized_areas
		WHERE id IN (%s)
		ORDER BY name ASC
	`, strings.Join(placeholders, ","))
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query specialized areas by IDs: %w", err)
	}
	defer rows.Close()
	
	var areas []*domain.SpecializedArea
	for rows.Next() {
		var area domain.SpecializedArea
		if err := rows.Scan(&area.ID, &area.Name, &area.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan specialized area: %w", err)
		}
		areas = append(areas, &area)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating specialized areas: %w", err)
	}
	
	return areas, nil
}

// CreateSpecializedArea creates a new specialized area
func (s *SQLiteStore) CreateSpecializedArea(area *domain.SpecializedArea) (int64, error) {
	query := `
		INSERT INTO specialized_areas (name, created_at)
		VALUES (?, CURRENT_TIMESTAMP)
	`
	
	result, err := s.db.Exec(query, area.Name)
	if err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, fmt.Errorf("specialized area '%s' already exists", area.Name)
		}
		return 0, fmt.Errorf("failed to create specialized area: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}
	
	return id, nil
}