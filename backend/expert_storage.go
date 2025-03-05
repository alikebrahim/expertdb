package main

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateExpert creates a new expert in the database
func (s *SQLiteStore) CreateExpert(expert *Expert) (int64, error) {
	query := `
		INSERT INTO experts (
			expert_id, name, designation, institution, is_bahraini, 
			nationality, is_available, rating, role, employment_type, 
			general_area, specialized_area, is_trained, cv_path, 
			phone, email, is_published, isced_level_id, isced_field_id,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	// Set default values if not provided
	if expert.CreatedAt.IsZero() {
		expert.CreatedAt = time.Now()
	}
	if expert.UpdatedAt.IsZero() {
		expert.UpdatedAt = expert.CreatedAt
	}
	
	var iscedLevelID, iscedFieldID sql.NullInt64
	if expert.ISCEDLevel != nil {
		iscedLevelID.Int64 = expert.ISCEDLevel.ID
		iscedLevelID.Valid = true
	}
	if expert.ISCEDField != nil {
		iscedFieldID.Int64 = expert.ISCEDField.ID
		iscedFieldID.Valid = true
	}
	
	result, err := s.db.Exec(
		query,
		expert.ExpertID, expert.Name, expert.Designation, expert.Institution,
		expert.IsBahraini, expert.Nationality, expert.IsAvailable, expert.Rating,
		expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.Phone, expert.Email, expert.IsPublished,
		iscedLevelID, iscedFieldID, expert.CreatedAt, expert.UpdatedAt,
	)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create expert: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get expert ID: %w", err)
	}
	
	// Set the expert ID
	expert.ID = id
	
	return id, nil
}