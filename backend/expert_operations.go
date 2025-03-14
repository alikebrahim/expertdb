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
// education classification, specialization areas, documents, and engagement history.
//
// Inputs:
//   - id (int64): The unique identifier of the expert to retrieve
//
// Returns:
//   - *Expert: A complete expert record with all related data
//   - error: ErrNotFound if the expert doesn't exist, or any database error
//
// Flow:
//   1. Retrieve base expert data from experts table
//   2. Parse and convert timestamp fields
//   3. Load related ISCED education classification data
//   4. Load expert specialization areas
//   5. Load associated documents
//   6. Load engagement history records
func (s *SQLiteStore) GetExpert(id int64) (*Expert, error) {
	logger := GetLogger()
	
	// Step 1: Retrieve base expert data with main query
	// This query gets all core expert fields from the experts table
	query := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, e.is_bahraini, 
			   e.is_available, e.rating, e.role, e.employment_type, e.general_area, 
			   e.specialized_area, e.is_trained, e.cv_path, e.phone, e.email, e.is_published, 
			   e.biography, e.isced_level_id, e.isced_field_id, e.created_at, e.updated_at
		FROM experts e
		WHERE e.id = ?
	`
	
	var expert Expert
	var createdAt, updatedAt string
	var iscedLevelID, iscedFieldID sql.NullInt64
	
	// Execute the query and scan results into expert struct
	err := s.db.QueryRow(query, id).Scan(
		&expert.ID, &expert.ExpertID, &expert.Name, &expert.Designation, &expert.Institution,
		&expert.IsBahraini, &expert.IsAvailable, &expert.Rating, &expert.Role, &expert.EmploymentType,
		&expert.GeneralArea, &expert.SpecializedArea, &expert.IsTrained, &expert.CVPath,
		&expert.Phone, &expert.Email, &expert.IsPublished, &expert.Biography, &iscedLevelID, &iscedFieldID,
		&createdAt, &updatedAt,
	)
	
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
	
	// Step 3: Load related ISCED education classification data if available
	
	// Load ISCED level if the expert has one assigned
	if iscedLevelID.Valid {
		var level ISCEDLevel
		// Query ISCED level details
		err := s.db.QueryRow(
			"SELECT id, code, name, description FROM isced_levels WHERE id = ?",
			iscedLevelID.Int64,
		).Scan(&level.ID, &level.Code, &level.Name, &level.Description)
		if err == nil {
			expert.ISCEDLevel = &level
		} else {
			logger.Debug("Failed to load ISCED level %d for expert %d: %v", 
				iscedLevelID.Int64, id, err)
		}
	}
	
	// Load ISCED field if the expert has one assigned
	if iscedFieldID.Valid {
		var field ISCEDField
		// Query ISCED field details
		err := s.db.QueryRow(
			"SELECT id, broad_code, broad_name, narrow_code, narrow_name, detailed_code, detailed_name, description FROM isced_fields WHERE id = ?",
			iscedFieldID.Int64,
		).Scan(&field.ID, &field.BroadCode, &field.BroadName, &field.NarrowCode, &field.NarrowName, &field.DetailedCode, &field.DetailedName, &field.Description)
		if err == nil {
			expert.ISCEDField = &field
		} else {
			logger.Debug("Failed to load ISCED field %d for expert %d: %v", 
				iscedFieldID.Int64, id, err)
		}
	}
	
	// Step 4: Load expert specialization areas
	// Query to get all areas associated with this expert through the mapping table
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
	} else {
		logger.Debug("Failed to load areas for expert %d: %v", id, err)
	}
	
	// Step 5: Load associated documents
	docs, err := s.GetDocumentsByExpertID(expert.ID)
	if err != nil {
		logger.Debug("Failed to load documents for expert %d: %v", id, err)
	} else if docs != nil {
		for _, doc := range docs {
			expert.Documents = append(expert.Documents, *doc)
		}
	}
	
	// Step 6: Load engagement history
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
// It sets the updated_at timestamp to the current time and handles nullable fields
// like ISCED classification IDs.
//
// Inputs:
//   - expert (*Expert): The expert object with updated field values to save
//
// Returns:
//   - error: Any database error that occurs during the update operation
//
// Note: This method updates only the core expert data. Associated data like
// documents, areas, and engagements must be updated separately using their
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
			biography = ?, isced_level_id = ?, isced_field_id = ?, updated_at = ?
		WHERE id = ?
	`
	
	// Step 3: Handle nullable foreign key fields (ISCED classifications)
	// Convert pointers to sql.NullInt64 for nullable database fields
	var iscedLevelID, iscedFieldID sql.NullInt64
	
	// Set ISCED level ID if present
	if expert.ISCEDLevel != nil {
		iscedLevelID.Int64 = expert.ISCEDLevel.ID
		iscedLevelID.Valid = true
	}
	
	// Set ISCED field ID if present
	if expert.ISCEDField != nil {
		iscedFieldID.Int64 = expert.ISCEDField.ID
		iscedFieldID.Valid = true
	}
	
	// Step 4: Execute the update query
	result, err := s.db.Exec(
		query,
		expert.ExpertID, expert.Name, expert.Designation, expert.Institution,
		expert.IsBahraini, expert.Nationality, expert.IsAvailable, expert.Rating,
		expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.Phone, expert.Email, expert.IsPublished,
		expert.Biography, iscedLevelID, iscedFieldID, expert.UpdatedAt, expert.ID,
	)
	
	// Step 5: Handle database errors
	if err != nil {
		logger.Error("Failed to update expert ID %d: %v", expert.ID, err)
		return fmt.Errorf("failed to update expert: %w", err)
	}
	
	// Step 6: Check if any rows were affected (optional validation)
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
	
	// Step 2.1: Delete expert-area mappings first
	// These are junction table records linking experts to their specialization areas
	_, err = tx.Exec("DELETE FROM expert_area_map WHERE expert_id = ?", id)
	if err != nil {
		return rollback(err, "failed to delete expert area mappings")
	}
	
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