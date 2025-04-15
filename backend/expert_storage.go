package main

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateExpert creates a new expert in the database
func (s *SQLiteStore) CreateExpert(expert *Expert) (int64, error) {
	// Generate a unique expert_id if not provided or empty
	if expert.ExpertID == "" {
		var err error
		expert.ExpertID, err = s.GenerateUniqueExpertID()
		if err != nil {
			return 0, fmt.Errorf("failed to generate unique expert ID: %w", err)
		}
	} else {
		// Check if the provided expert_id already exists
		exists, err := s.ExpertIDExists(expert.ExpertID)
		if err != nil {
			return 0, fmt.Errorf("failed to check if expert ID exists: %w", err)
		}
		if exists {
			return 0, fmt.Errorf("expert ID already exists: %s", expert.ExpertID)
		}
	}
	
	query := `
		INSERT INTO experts (
			expert_id, name, designation, institution, is_bahraini, is_available, rating,
			role, employment_type, general_area, specialized_area, is_trained,
			cv_path, phone, email, is_published, biography, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	if expert.CreatedAt.IsZero() {
		expert.CreatedAt = time.Now().UTC()
		expert.UpdatedAt = expert.CreatedAt
	}
	
	result, err := s.db.Exec(
		query,
		expert.ExpertID, expert.Name, expert.Designation, expert.Institution,
		expert.IsBahraini, expert.IsAvailable, expert.Rating,
		expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.Phone, expert.Email, expert.IsPublished,
		expert.Biography, expert.CreatedAt, expert.UpdatedAt,
	)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create expert: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get expert ID: %w", err)
	}
	
	return id, nil
}

// ExpertIDExists checks if an expert_id already exists in the database
func (s *SQLiteStore) ExpertIDExists(expertID string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM experts WHERE expert_id = ?"
	err := s.db.QueryRow(query, expertID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if expert ID exists: %w", err)
	}
	return count > 0, nil
}

// GenerateUniqueExpertID generates a unique expert ID
func (s *SQLiteStore) GenerateUniqueExpertID() (string, error) {
	// Generate a base ID
	baseID := fmt.Sprintf("EXP-%d-%d", time.Now().Unix(), time.Now().UnixNano()%1000)
	
	// Check if the generated ID already exists
	exists, err := s.ExpertIDExists(baseID)
	if err != nil {
		return "", fmt.Errorf("failed to check if expert ID exists: %w", err)
	}
	
	// If the ID already exists, add a random component and check again
	if exists {
		// In the unlikely case of a collision, add a random suffix
		randomSuffix := fmt.Sprintf("%04d", time.Now().UnixNano()%10000)
		baseID = fmt.Sprintf("%s-%s", baseID, randomSuffix)
		
		// Check one more time
		exists, err = s.ExpertIDExists(baseID)
		if err != nil {
			return "", fmt.Errorf("failed to check if expert ID exists: %w", err)
		}
		if exists {
			return "", fmt.Errorf("failed to generate a unique expert ID after multiple attempts")
		}
	}
	
	return baseID, nil
}

// GetExpertAreas retrieves all expert areas from the database
func (s *SQLiteStore) GetExpertAreas() ([]Area, error) {
	query := "SELECT id, name FROM expert_areas ORDER BY name"
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expert areas: %w", err)
	}
	defer rows.Close()
	
	var areas []Area
	for rows.Next() {
		var area Area
		if err := rows.Scan(&area.ID, &area.Name); err != nil {
			return nil, fmt.Errorf("failed to scan area row: %w", err)
		}
		areas = append(areas, area)
	}
	
	return areas, nil
}

// GetExpertAreaByID retrieves a specific area by its ID
func (s *SQLiteStore) GetExpertAreaByID(id int64) (*Area, error) {
	var area Area
	err := s.db.QueryRow("SELECT id, name FROM expert_areas WHERE id = ?", id).Scan(&area.ID, &area.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get expert area: %w", err)
	}
	return &area, nil
}