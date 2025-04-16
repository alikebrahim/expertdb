package sqlite

import (
	"database/sql"
	"fmt"
	
	"expertdb/internal/domain"
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