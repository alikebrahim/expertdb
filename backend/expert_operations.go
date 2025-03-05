package main

import (
	"database/sql"
	"fmt"
	"time"
)

// GetExpert retrieves an expert by ID
func (s *SQLiteStore) GetExpert(id int64) (*Expert, error) {
	query := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, e.is_bahraini, 
			   e.is_available, e.rating, e.role, e.employment_type, e.general_area, 
			   e.specialized_area, e.is_trained, e.cv_path, e.phone, e.email, e.is_published, 
			   e.isced_level_id, e.isced_field_id, e.created_at, e.updated_at
		FROM experts e
		WHERE e.id = ?
	`
	
	var expert Expert
	var createdAt, updatedAt string
	var iscedLevelID, iscedFieldID sql.NullInt64
	
	err := s.db.QueryRow(query, id).Scan(
		&expert.ID, &expert.ExpertID, &expert.Name, &expert.Designation, &expert.Institution,
		&expert.IsBahraini, &expert.IsAvailable, &expert.Rating, &expert.Role, &expert.EmploymentType,
		&expert.GeneralArea, &expert.SpecializedArea, &expert.IsTrained, &expert.CVPath,
		&expert.Phone, &expert.Email, &expert.IsPublished, &iscedLevelID, &iscedFieldID,
		&createdAt, &updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get expert: %w", err)
	}
	
	// Parse timestamps
	expert.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	if updatedAt != "" {
		expert.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	}
	
	// Add ISCED data
	if iscedLevelID.Valid {
		var level ISCEDLevel
		err := s.db.QueryRow(
			"SELECT id, code, name, description FROM isced_levels WHERE id = ?",
			iscedLevelID.Int64,
		).Scan(&level.ID, &level.Code, &level.Name, &level.Description)
		if err == nil {
			expert.ISCEDLevel = &level
		}
	}
	
	if iscedFieldID.Valid {
		var field ISCEDField
		err := s.db.QueryRow(
			"SELECT id, broad_code, broad_name, narrow_code, narrow_name, detailed_code, detailed_name, description FROM isced_fields WHERE id = ?",
			iscedFieldID.Int64,
		).Scan(&field.ID, &field.BroadCode, &field.BroadName, &field.NarrowCode, &field.NarrowName, &field.DetailedCode, &field.DetailedName, &field.Description)
		if err == nil {
			expert.ISCEDField = &field
		}
	}
	
	// Get areas
	areaRows, err := s.db.Query(
		"SELECT a.id, a.name FROM expert_areas a JOIN expert_area_map m ON a.id = m.area_id WHERE m.expert_id = ?",
		expert.ID,
	)
	if err == nil {
		defer areaRows.Close()
		for areaRows.Next() {
			var area Area
			if err := areaRows.Scan(&area.ID, &area.Name); err == nil {
				expert.Areas = append(expert.Areas, area)
			}
		}
	}
	
	// Get documents
	docs, _ := s.GetDocumentsByExpertID(expert.ID)
	if docs != nil {
		for _, doc := range docs {
			expert.Documents = append(expert.Documents, *doc)
		}
	}
	
	// Get engagements
	expertEngagements, _ := s.GetEngagementsByExpertID(expert.ID)
	for _, eng := range expertEngagements {
		expert.Engagements = append(expert.Engagements, *eng)
	}
	
	return &expert, nil
}

// UpdateExpert updates an expert's record
func (s *SQLiteStore) UpdateExpert(expert *Expert) error {
	expert.UpdatedAt = time.Now()
	
	query := `
		UPDATE experts
		SET expert_id = ?, name = ?, designation = ?, institution = ?,
			is_bahraini = ?, nationality = ?, is_available = ?, rating = ?, role = ?,
			employment_type = ?, general_area = ?, specialized_area = ?,
			is_trained = ?, cv_path = ?, phone = ?, email = ?, is_published = ?,
			isced_level_id = ?, isced_field_id = ?, updated_at = ?
		WHERE id = ?
	`
	
	var iscedLevelID, iscedFieldID sql.NullInt64
	if expert.ISCEDLevel != nil {
		iscedLevelID.Int64 = expert.ISCEDLevel.ID
		iscedLevelID.Valid = true
	}
	if expert.ISCEDField != nil {
		iscedFieldID.Int64 = expert.ISCEDField.ID
		iscedFieldID.Valid = true
	}
	
	_, err := s.db.Exec(
		query,
		expert.ExpertID, expert.Name, expert.Designation, expert.Institution,
		expert.IsBahraini, expert.Nationality, expert.IsAvailable, expert.Rating,
		expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.Phone, expert.Email, expert.IsPublished,
		iscedLevelID, iscedFieldID, expert.UpdatedAt, expert.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update expert: %w", err)
	}
	
	return nil
}

// DeleteExpert deletes an expert by ID
func (s *SQLiteStore) DeleteExpert(id int64) error {
	// Begin transaction to delete expert and related records
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Delete expert areas mapping first to maintain referential integrity
	_, err = tx.Exec("DELETE FROM expert_area_map WHERE expert_id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete expert area mappings: %w", err)
	}
	
	// Delete expert engagements
	_, err = tx.Exec("DELETE FROM expert_engagements WHERE expert_id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete expert engagements: %w", err)
	}
	
	// Delete AI analysis results
	_, err = tx.Exec("DELETE FROM ai_analysis WHERE expert_id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete AI analysis results: %w", err)
	}
	
	// Delete expert documents
	_, err = tx.Exec("DELETE FROM expert_documents WHERE expert_id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete expert documents: %w", err)
	}
	
	// Finally, delete the expert
	_, err = tx.Exec("DELETE FROM experts WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete expert: %w", err)
	}
	
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}