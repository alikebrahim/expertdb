// Package main provides the backend functionality for the ExpertDB application
package main

import (
	"database/sql"
	"fmt"
	"time"
)

// GetExpert retrieves a complete expert profile by ID including all related data
//
// This function retrieves an expert by ID and loads all associated data including
// general area, documents, and engagement history.
//
// Inputs:
//   - id (int64): The unique identifier of the expert to retrieve
//
// Returns:
//   - *Expert: A complete expert record with all related data
//   - error: ErrNotFound if the expert doesn't exist, or any database error
//
// Flow:
//   1. Retrieve base expert data from experts table with joined general area name
//   2. Parse and convert timestamp fields
//   3. Load associated documents
//   4. Load engagement history records
func (s *SQLiteStore) GetExpert(id int64) (*Expert, error) {
	logger := GetLogger()
	
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
	var createdAt, updatedAt string
	var generalAreaName string
	
	// Execute the query and scan results into expert struct
	var generalAreaName sql.NullString
	err := s.db.QueryRow(query, id).Scan(
		&expert.ID, &expert.ExpertID, &expert.Name, &expert.Designation, &expert.Institution,
		&expert.IsBahraini, &expert.IsAvailable, &expert.Rating, &expert.Role, &expert.EmploymentType,
		&expert.GeneralArea, &generalAreaName, &expert.SpecializedArea, &expert.IsTrained, &expert.CVPath,
		&expert.Phone, &expert.Email, &expert.IsPublished, &expert.Biography,
		&createdAt, &updatedAt,
	)
	
	// Store area name if available
	if generalAreaName.Valid {
		expert.GeneralAreaName = generalAreaName.String
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
	if updatedAt != "" {
		expert.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
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

// UpdateExpert updates an expert's record in the database
//
// This function updates all fields of an expert record with the provided values.
// It sets the updated_at timestamp to the current time.
//
// Inputs:
//   - expert (*Expert): The expert object with updated field values to save
//
// Returns:
//   - error: Any database error that occurs during the update operation
//
// Note: This method updates only the core expert data. Associated data like
// documents and engagements must be updated separately using their
// respective methods.
func (s *SQLiteStore) UpdateExpert(expert *Expert) error {
	logger := GetLogger()
	
	// Step 1: Set the updated timestamp to the current time
	expert.UpdatedAt = time.Now()
	
	// Step 2: Prepare the SQL update query
	query := `
		UPDATE experts
		SET expert_id = ?, name = ?, designation = ?, institution = ?,
			is_bahraini = ?, nationality = ?, is_available = ?, rating = ?, role = ?,
			employment_type = ?, general_area = ?, specialized_area = ?,
			is_trained = ?, cv_path = ?, phone = ?, email = ?, is_published = ?,
			biography = ?, updated_at = ?
		WHERE id = ?
	`
	
	// Step 3: Execute the update query
	result, err := s.db.Exec(
		query,
		expert.ExpertID, expert.Name, expert.Designation, expert.Institution,
		expert.IsBahraini, expert.Nationality, expert.IsAvailable, expert.Rating,
		expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.Phone, expert.Email, expert.IsPublished,
		expert.Biography, expert.UpdatedAt, expert.ID,
	)
	
	// Step 4: Handle database errors
	if err != nil {
		logger.Error("Failed to update expert ID %d: %v", expert.ID, err)
		return fmt.Errorf("failed to update expert: %w", err)
	}
	
	// Step 5: Check if any rows were affected (optional validation)
	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		// This indicates the expert ID doesn't exist
		logger.Warn("Update for expert ID %d affected 0 rows - record may not exist", expert.ID)
	} else {
		logger.Debug("Successfully updated expert ID %d: %s", expert.ID, expert.Name)
	}
	
	return nil
}

// DeleteExpert deletes an expert and all related records in a transaction
//
// This function removes an expert and all associated data including area mappings,
// engagements, and documents in a single atomic transaction. If any part of the
// deletion fails, the entire operation is rolled back.
//
// Inputs:
//   - id (int64): The unique identifier of the expert to delete
//
// Returns:
//   - error: Any database error that occurs during the deletion transaction
//
// Flow:
//   1. Begin database transaction
//   2. Delete related records in correct order to maintain referential integrity
//   3. Delete the expert record
//   4. Commit the transaction or roll back if any step fails
//
// NOTE: This is a destructive operation that permanently removes all expert data
// and cannot be undone. Consider implementing an archive/soft delete mechanism
// for production applications.
func (s *SQLiteStore) DeleteExpert(id int64) error {
	logger := GetLogger()
	
	// Step 1: Begin transaction to ensure atomic operation
	// This ensures that either all deletions succeed or none do
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("Failed to begin transaction for deleting expert ID %d: %v", id, err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Transaction rollback helper function to avoid repetition
	rollback := func(err error, message string) error {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logger.Error("Failed to rollback transaction: %v", rollbackErr)
		}
		logger.Error("%s for expert ID %d: %v", message, id, err)
		return fmt.Errorf("%s: %w", message, err)
	}
	
	// Step 2: Delete related records in order that maintains referential integrity
	
	// Note: expert_area_map junction table has been removed in schema simplification
	
	// Step 2.2: Delete expert engagements
	// These track participation of experts in various activities
	_, err = tx.Exec("DELETE FROM expert_engagements WHERE expert_id = ?", id)
	if err != nil {
		return rollback(err, "failed to delete expert engagements")
	}
	
	// Step 2.3: Delete expert documents
	// These are files like CVs and certificates associated with the expert
	_, err = tx.Exec("DELETE FROM expert_documents WHERE expert_id = ?", id)
	if err != nil {
		return rollback(err, "failed to delete expert documents")
	}
	
	// Step 3: Delete the expert record itself
	// This must be done after all dependent records are removed
	result, err := tx.Exec("DELETE FROM experts WHERE id = ?", id)
	if err != nil {
		return rollback(err, "failed to delete expert")
	}
	
	// Verify that the expert was actually deleted
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Expert ID didn't exist, but we don't consider this an error
		// since the end state is as requested (expert doesn't exist)
		logger.Warn("Expert ID %d was not found for deletion", id)
	}
	
	// Step 4: Commit the transaction
	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction for deleting expert ID %d: %v", id, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	logger.Info("Successfully deleted expert ID %d and all related records", id)
	return nil
}

// NOTE: Consider implementing a confirmation mechanism or soft delete feature
// to prevent accidental data loss. The current implementation permanently deletes
// all expert data with no recovery option.

// GetExpertAreaByID retrieves a single expert area by ID
func (s *SQLiteStore) GetExpertAreaByID(id int64) (*Area, error) {
	logger := GetLogger()
	
	var area Area
	err := s.db.QueryRow("SELECT id, name FROM expert_areas WHERE id = ?", id).Scan(&area.ID, &area.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("Expert area not found with ID: %d", id)
			return nil, ErrNotFound
		}
		logger.Error("Database error retrieving expert area ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get expert area: %w", err)
	}
	
	return &area, nil
}

// GetExpertAreas retrieves all expert areas
func (s *SQLiteStore) GetExpertAreas() ([]Area, error) {
	logger := GetLogger()
	
	rows, err := s.db.Query("SELECT id, name FROM expert_areas ORDER BY name")
	if err != nil {
		logger.Error("Database error retrieving expert areas: %v", err)
		return nil, fmt.Errorf("failed to get expert areas: %w", err)
	}
	defer rows.Close()
	
	var areas []Area
	for rows.Next() {
		var area Area
		if err := rows.Scan(&area.ID, &area.Name); err != nil {
			logger.Error("Error scanning expert area row: %v", err)
			return nil, fmt.Errorf("failed to scan expert area: %w", err)
		}
		areas = append(areas, area)
	}
	
	if err = rows.Err(); err != nil {
		logger.Error("Error iterating expert area rows: %v", err)
		return nil, fmt.Errorf("failed to iterate expert areas: %w", err)
	}
	
	logger.Debug("Successfully retrieved %d expert areas", len(areas))
	return areas, nil
}