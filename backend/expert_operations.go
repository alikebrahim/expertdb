package main

import (
	"database/sql"
	"fmt"
	"time"
)

// GetExpert retrieves an expert by ID
func (s *SQLiteStore) GetExpert(id int64) (*Expert, error) {
	logger := GetLogger()
	logger.Debug("Getting expert with ID: %d", id)
	
	// Step 1: Retrieve base expert data
	// This query gets all core expert fields from the experts table
	query := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, e.is_bahraini, 
			   e.is_available, e.rating, e.role, e.employment_type, e.general_area, 
			   e.specialized_area, e.is_trained, e.cv_path, e.phone, e.email, e.is_published, 
			   e.biography, e.created_at, e.updated_at
		FROM experts e
		WHERE e.id = ?
	`
	
	var expert Expert
	var createdAt string
	var nullableExpertID, nullableCVPath, nullablePhone, nullableUpdatedAt, nullableSpecializedArea, nullableBiography sql.NullString
	var nullableIsAvailable sql.NullBool
	
	// Execute the query and scan results into expert struct
	err := s.db.QueryRow(query, id).Scan(
		&expert.ID, &nullableExpertID, &expert.Name, &expert.Designation, &expert.Institution,
		&expert.IsBahraini, &nullableIsAvailable, &expert.Rating, &expert.Role, &expert.EmploymentType,
		&expert.GeneralArea, &nullableSpecializedArea, &expert.IsTrained, &nullableCVPath,
		&nullablePhone, &expert.Email, &expert.IsPublished, &nullableBiography,
		&createdAt, &nullableUpdatedAt,
	)
	
	// Handle nullable fields
	if nullableExpertID.Valid {
		expert.ExpertID = nullableExpertID.String
	}
	
	if nullableCVPath.Valid {
		expert.CVPath = nullableCVPath.String
	}
	
	if nullablePhone.Valid {
		expert.Phone = nullablePhone.String
	}
	
	if nullableIsAvailable.Valid {
		expert.IsAvailable = nullableIsAvailable.Bool
	}
	
	if nullableSpecializedArea.Valid {
		expert.SpecializedArea = nullableSpecializedArea.String
	}
	
	if nullableBiography.Valid {
		expert.Biography = nullableBiography.String
	}
	
	// Load the area name if needed
	if expert.GeneralArea > 0 {
		area, err := s.GetExpertAreaByID(expert.GeneralArea)
		if err == nil && area != nil {
			expert.GeneralAreaName = area.Name
		}
	}
	
	// Handle query errors
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("Expert not found with ID: %d", id)
			return nil, ErrNotFound
		}
		logger.Error("Database error retrieving expert ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get expert: %w", err)
	}
	
	// Step 2: Parse timestamp strings into time.Time objects
	expert.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	if nullableUpdatedAt.Valid {
		expert.UpdatedAt, _ = time.Parse(time.RFC3339, nullableUpdatedAt.String)
	}
	
	// Step 3: Load associated documents
	docs, err := s.GetDocumentsByExpertID(expert.ID)
	if err != nil {
		logger.Debug("Failed to load documents for expert %d: %v", id, err)
	} else if docs != nil {
		for _, doc := range docs {
			expert.Documents = append(expert.Documents, *doc)
		}
	}
	
	// Step 4: Load engagement history
	expertEngagements, err := s.GetEngagementsByExpertID(expert.ID)
	if err != nil {
		logger.Debug("Failed to load engagements for expert %d: %v", id, err)
	} else {
		for _, eng := range expertEngagements {
			expert.Engagements = append(expert.Engagements, *eng)
		}
	}
	
	logger.Debug("Successfully retrieved expert ID %d: %s", id, expert.Name)
	return &expert, nil
}

// UpdateExpert updates an existing expert in the database
func (s *SQLiteStore) UpdateExpert(expert *Expert) error {
	logger := GetLogger()
	logger.Debug("Updating expert with ID: %d", expert.ID)
	
	// Step 1: Validate parameters
	if expert.ID <= 0 {
		return fmt.Errorf("invalid expert ID: %d", expert.ID)
	}
	
	// Step 2: Update timestamps
	expert.UpdatedAt = time.Now().UTC()
	
	// Build update query for all fields
	query := `
		UPDATE experts
		SET expert_id = ?, name = ?, designation = ?, institution = ?,
			is_bahraini = ?, is_available = ?, rating = ?,
			role = ?, employment_type = ?, general_area = ?, specialized_area = ?,
			is_trained = ?, cv_path = ?, phone = ?, email = ?, is_published = ?,
			biography = ?, updated_at = ?
		WHERE id = ?
	`
	
	// Step 3: Execute the update query
	result, err := s.db.Exec(
		query,
		expert.ExpertID, expert.Name, expert.Designation, expert.Institution,
		expert.IsBahraini, expert.IsAvailable, expert.Rating,
		expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.Phone, expert.Email, expert.IsPublished,
		expert.Biography, expert.UpdatedAt, expert.ID,
	)
	
	if err != nil {
		logger.Error("Failed to update expert ID %d: %v", expert.ID, err)
		return fmt.Errorf("failed to update expert: %w", err)
	}
	
	// Step 4: Verify that the update affected a row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Failed to get rows affected for expert update ID %d: %v", expert.ID, err)
		return fmt.Errorf("failed to verify expert update: %w", err)
	}
	
	if rowsAffected == 0 {
		logger.Warn("No rows affected when updating expert ID %d", expert.ID)
		return ErrNotFound
	}
	
	logger.Info("Successfully updated expert ID %d: %s", expert.ID, expert.Name)
	return nil
}

// DeleteExpert deletes an expert by ID
func (s *SQLiteStore) DeleteExpert(id int64) error {
	logger := GetLogger()
	logger.Debug("Deleting expert with ID: %d", id)
	
	// Use a transaction to ensure data integrity
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("Failed to start transaction for expert deletion: %v", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	
	// First delete associated documents (cascaded delete)
	_, err = tx.Exec("DELETE FROM expert_documents WHERE expert_id = ?", id)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to delete documents for expert ID %d: %v", id, err)
		return fmt.Errorf("failed to delete expert documents: %w", err)
	}
	
	// Then delete associated engagements (cascaded delete)
	_, err = tx.Exec("DELETE FROM expert_engagements WHERE expert_id = ?", id)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to delete engagements for expert ID %d: %v", id, err)
		return fmt.Errorf("failed to delete expert engagements: %w", err)
	}
	
	// Finally delete the expert
	result, err := tx.Exec("DELETE FROM experts WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to delete expert ID %d: %v", id, err)
		return fmt.Errorf("failed to delete expert: %w", err)
	}
	
	// Verify that a row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		logger.Error("Failed to get rows affected for expert deletion ID %d: %v", id, err)
		return fmt.Errorf("failed to verify expert deletion: %w", err)
	}
	
	if rowsAffected == 0 {
		tx.Rollback()
		logger.Warn("No rows affected when deleting expert ID %d", id)
		return ErrNotFound
	}
	
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction for expert deletion: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	logger.Info("Successfully deleted expert ID %d", id)
	return nil
}