package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"expertdb/internal/domain"
	"expertdb/internal/logger"
)

// ListExpertEditRequests retrieves expert edit requests with filtering and pagination
func (s *SQLiteStore) ListExpertEditRequests(filters map[string]interface{}, limit, offset int) ([]*domain.ExpertEditRequest, error) {
	log := logger.Get()
	log.Debug("Listing expert edit requests with filters: %+v, limit: %d, offset: %d", filters, limit, offset)

	// Build dynamic WHERE clause
	whereClause := "1 = 1"
	args := []interface{}{}
	argIndex := 1

	if expertID, ok := filters["expertId"].(int64); ok && expertID > 0 {
		whereClause += fmt.Sprintf(" AND eer.expert_id = ?%d", argIndex)
		args = append(args, expertID)
		argIndex++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		whereClause += fmt.Sprintf(" AND eer.status = ?%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if createdBy, ok := filters["createdBy"].(int64); ok && createdBy > 0 {
		whereClause += fmt.Sprintf(" AND eer.created_by = ?%d", argIndex)
		args = append(args, createdBy)
		argIndex++
	}

	if reviewedBy, ok := filters["reviewedBy"].(int64); ok && reviewedBy > 0 {
		whereClause += fmt.Sprintf(" AND eer.reviewed_by = ?%d", argIndex)
		args = append(args, reviewedBy)
		argIndex++
	}

	// Add limit and offset
	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf(" LIMIT ?%d", argIndex)
		args = append(args, limit)
		argIndex++
		
		if offset > 0 {
			limitClause += fmt.Sprintf(" OFFSET ?%d", argIndex)
			args = append(args, offset)
			argIndex++
		}
	}

	query := `
		SELECT 
			eer.id, eer.expert_id, eer.name, eer.designation, eer.institution,
			eer.phone, eer.email, eer.is_bahraini, eer.is_available, eer.rating,
			eer.role, eer.employment_type, eer.general_area, eer.specialized_area,
			eer.is_trained, eer.is_published, eer.biography, eer.suggested_specialized_areas,
			eer.new_cv_path, eer.new_approval_document_path, eer.remove_cv, eer.remove_approval_document,
			eer.change_summary, eer.change_reason, eer.fields_changed,
			eer.status, eer.rejection_reason, eer.admin_notes,
			eer.created_at, eer.reviewed_at, eer.applied_at, eer.created_by, eer.reviewed_by,
			e.name as expert_name,
			cu.name as created_by_name,
			ru.name as reviewed_by_name
		FROM expert_edit_requests eer
		LEFT JOIN experts e ON eer.expert_id = e.id
		LEFT JOIN users cu ON eer.created_by = cu.id
		LEFT JOIN users ru ON eer.reviewed_by = ru.id
		WHERE ` + whereClause + `
		ORDER BY eer.created_at DESC` + limitClause

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Error("Failed to query expert edit requests: %v", err)
		return nil, fmt.Errorf("failed to query expert edit requests: %w", err)
	}
	defer rows.Close()

	var requests []*domain.ExpertEditRequest
	for rows.Next() {
		req, err := s.scanExpertEditRequest(rows)
		if err != nil {
			log.Error("Failed to scan expert edit request: %v", err)
			return nil, err
		}

		// Load experience and education changes
		if err := s.loadEditRequestChanges(req); err != nil {
			log.Error("Failed to load edit request changes: %v", err)
			return nil, err
		}

		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		log.Error("Error iterating expert edit requests: %v", err)
		return nil, err
	}

	log.Debug("Retrieved %d expert edit requests", len(requests))
	return requests, nil
}

// CountExpertEditRequests counts expert edit requests with filtering
func (s *SQLiteStore) CountExpertEditRequests(filters map[string]interface{}) (int, error) {
	log := logger.Get()

	// Build dynamic WHERE clause
	whereClause := "1 = 1"
	args := []interface{}{}
	argIndex := 1

	if expertID, ok := filters["expertId"].(int64); ok && expertID > 0 {
		whereClause += fmt.Sprintf(" AND expert_id = ?%d", argIndex)
		args = append(args, expertID)
		argIndex++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		whereClause += fmt.Sprintf(" AND status = ?%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if createdBy, ok := filters["createdBy"].(int64); ok && createdBy > 0 {
		whereClause += fmt.Sprintf(" AND created_by = ?%d", argIndex)
		args = append(args, createdBy)
		argIndex++
	}

	query := `SELECT COUNT(*) FROM expert_edit_requests WHERE ` + whereClause

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		log.Error("Failed to count expert edit requests: %v", err)
		return 0, fmt.Errorf("failed to count expert edit requests: %w", err)
	}

	return count, nil
}

// GetExpertEditRequest retrieves a single expert edit request by ID
func (s *SQLiteStore) GetExpertEditRequest(id int64) (*domain.ExpertEditRequest, error) {
	log := logger.Get()
	log.Debug("Getting expert edit request with ID: %d", id)

	query := `
		SELECT 
			eer.id, eer.expert_id, eer.name, eer.designation, eer.institution,
			eer.phone, eer.email, eer.is_bahraini, eer.is_available, eer.rating,
			eer.role, eer.employment_type, eer.general_area, eer.specialized_area,
			eer.is_trained, eer.is_published, eer.biography, eer.suggested_specialized_areas,
			eer.new_cv_path, eer.new_approval_document_path, eer.remove_cv, eer.remove_approval_document,
			eer.change_summary, eer.change_reason, eer.fields_changed,
			eer.status, eer.rejection_reason, eer.admin_notes,
			eer.created_at, eer.reviewed_at, eer.applied_at, eer.created_by, eer.reviewed_by,
			e.name as expert_name,
			cu.name as created_by_name,
			ru.name as reviewed_by_name
		FROM expert_edit_requests eer
		LEFT JOIN experts e ON eer.expert_id = e.id
		LEFT JOIN users cu ON eer.created_by = cu.id
		LEFT JOIN users ru ON eer.reviewed_by = ru.id
		WHERE eer.id = ?`

	row := s.db.QueryRow(query, id)
	req, err := s.scanExpertEditRequest(row)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("Expert edit request not found with ID: %d", id)
			return nil, domain.ErrNotFound
		}
		log.Error("Failed to scan expert edit request: %v", err)
		return nil, err
	}

	// Load experience and education changes
	if err := s.loadEditRequestChanges(req); err != nil {
		log.Error("Failed to load edit request changes: %v", err)
		return nil, err
	}

	log.Debug("Retrieved expert edit request: ID: %d, Expert: %s", req.ID, req.ExpertName)
	return req, nil
}

// CreateExpertEditRequest creates a new expert edit request
func (s *SQLiteStore) CreateExpertEditRequest(req *domain.ExpertEditRequest) (int64, error) {
	log := logger.Get()
	log.Debug("Creating expert edit request for expert ID: %d", req.ExpertID)

	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Serialize fields_changed to JSON
	fieldsChangedJSON, err := json.Marshal(req.FieldsChanged)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal fields changed: %w", err)
	}

	// Serialize suggested_specialized_areas to JSON
	suggestedAreasJSON, err := json.Marshal(req.SuggestedSpecializedAreas)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal suggested areas: %w", err)
	}

	// Insert main expert edit request
	query := `
		INSERT INTO expert_edit_requests (
			expert_id, name, designation, institution, phone, email,
			is_bahraini, is_available, rating, role, employment_type,
			general_area, specialized_area, is_trained, is_published, biography,
			suggested_specialized_areas, new_cv_path, new_approval_document_path,
			remove_cv, remove_approval_document, change_summary, change_reason,
			fields_changed, status, created_at, created_by
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?
		)`

	result, err := tx.Exec(query,
		req.ExpertID,
		stringPtrToInterface(req.Name),
		stringPtrToInterface(req.Designation),
		stringPtrToInterface(req.Institution),
		stringPtrToInterface(req.Phone),
		stringPtrToInterface(req.Email),
		boolPtrToInterface(req.IsBahraini),
		boolPtrToInterface(req.IsAvailable),
		intPtrToInterface(req.Rating),
		stringPtrToInterface(req.Role),
		stringPtrToInterface(req.EmploymentType),
		int64PtrToInterface(req.GeneralArea),
		stringPtrToInterface(req.SpecializedArea),
		boolPtrToInterface(req.IsTrained),
		boolPtrToInterface(req.IsPublished),
		stringPtrToInterface(req.Biography),
		string(suggestedAreasJSON),
		stringPtrToInterface(req.NewCVPath),
		stringPtrToInterface(req.NewApprovalDocumentPath),
		req.RemoveCV,
		req.RemoveApprovalDocument,
		req.ChangeSummary,
		req.ChangeReason,
		string(fieldsChangedJSON),
		req.Status,
		time.Now(),
		req.CreatedBy,
	)
	if err != nil {
		log.Error("Failed to insert expert edit request: %v", err)
		return 0, fmt.Errorf("failed to insert expert edit request: %w", err)
	}

	editRequestID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get edit request ID: %w", err)
	}

	// Insert experience changes
	for _, exp := range req.ExperienceChanges {
		exp.EditRequestID = editRequestID
		if err := s.insertExperienceChange(tx, &exp); err != nil {
			return 0, fmt.Errorf("failed to insert experience change: %w", err)
		}
	}

	// Insert education changes
	for _, edu := range req.EducationChanges {
		edu.EditRequestID = editRequestID
		if err := s.insertEducationChange(tx, &edu); err != nil {
			return 0, fmt.Errorf("failed to insert education change: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info("Expert edit request created successfully with ID: %d", editRequestID)
	return editRequestID, nil
}

// UpdateExpertEditRequestStatus updates the status of an expert edit request
func (s *SQLiteStore) UpdateExpertEditRequestStatus(id int64, status, rejectionReason, adminNotes string, reviewedBy int64) error {
	log := logger.Get()
	log.Debug("Updating expert edit request status: ID: %d, Status: %s", id, status)

	query := `
		UPDATE expert_edit_requests 
		SET status = ?, rejection_reason = ?, admin_notes = ?, reviewed_by = ?, reviewed_at = ?
		WHERE id = ?`

	result, err := s.db.Exec(query, status, rejectionReason, adminNotes, reviewedBy, time.Now(), id)
	if err != nil {
		log.Error("Failed to update expert edit request status: %v", err)
		return fmt.Errorf("failed to update expert edit request status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		log.Debug("No expert edit request found with ID: %d", id)
		return domain.ErrNotFound
	}

	log.Info("Expert edit request status updated: ID: %d, Status: %s", id, status)
	return nil
}

// UpdateExpertEditRequest updates an expert edit request
func (s *SQLiteStore) UpdateExpertEditRequest(req *domain.ExpertEditRequest) error {
	log := logger.Get()
	log.Debug("Updating expert edit request: ID: %d", req.ID)

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Serialize fields_changed to JSON
	fieldsChangedJSON, err := json.Marshal(req.FieldsChanged)
	if err != nil {
		return fmt.Errorf("failed to marshal fields changed: %w", err)
	}

	// Serialize suggested_specialized_areas to JSON
	suggestedAreasJSON, err := json.Marshal(req.SuggestedSpecializedAreas)
	if err != nil {
		return fmt.Errorf("failed to marshal suggested areas: %w", err)
	}

	// Update main expert edit request
	query := `
		UPDATE expert_edit_requests SET
			name = ?, designation = ?, institution = ?, phone = ?, email = ?,
			is_bahraini = ?, is_available = ?, rating = ?, role = ?, employment_type = ?,
			general_area = ?, specialized_area = ?, is_trained = ?, is_published = ?, biography = ?,
			suggested_specialized_areas = ?, new_cv_path = ?, new_approval_document_path = ?,
			remove_cv = ?, remove_approval_document = ?, change_summary = ?, change_reason = ?,
			fields_changed = ?
		WHERE id = ?`

	result, err := tx.Exec(query,
		stringPtrToInterface(req.Name),
		stringPtrToInterface(req.Designation),
		stringPtrToInterface(req.Institution),
		stringPtrToInterface(req.Phone),
		stringPtrToInterface(req.Email),
		boolPtrToInterface(req.IsBahraini),
		boolPtrToInterface(req.IsAvailable),
		intPtrToInterface(req.Rating),
		stringPtrToInterface(req.Role),
		stringPtrToInterface(req.EmploymentType),
		int64PtrToInterface(req.GeneralArea),
		stringPtrToInterface(req.SpecializedArea),
		boolPtrToInterface(req.IsTrained),
		boolPtrToInterface(req.IsPublished),
		stringPtrToInterface(req.Biography),
		string(suggestedAreasJSON),
		stringPtrToInterface(req.NewCVPath),
		stringPtrToInterface(req.NewApprovalDocumentPath),
		req.RemoveCV,
		req.RemoveApprovalDocument,
		req.ChangeSummary,
		req.ChangeReason,
		string(fieldsChangedJSON),
		req.ID,
	)
	if err != nil {
		log.Error("Failed to update expert edit request: %v", err)
		return fmt.Errorf("failed to update expert edit request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		log.Debug("No expert edit request found with ID: %d", req.ID)
		return domain.ErrNotFound
	}

	// Delete existing experience and education changes
	if _, err := tx.Exec("DELETE FROM expert_edit_request_experience WHERE expert_edit_request_id = ?", req.ID); err != nil {
		return fmt.Errorf("failed to delete existing experience changes: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM expert_edit_request_education WHERE expert_edit_request_id = ?", req.ID); err != nil {
		return fmt.Errorf("failed to delete existing education changes: %w", err)
	}

	// Insert updated experience changes
	for _, exp := range req.ExperienceChanges {
		exp.EditRequestID = req.ID
		if err := s.insertExperienceChange(tx, &exp); err != nil {
			return fmt.Errorf("failed to insert experience change: %w", err)
		}
	}

	// Insert updated education changes
	for _, edu := range req.EducationChanges {
		edu.EditRequestID = req.ID
		if err := s.insertEducationChange(tx, &edu); err != nil {
			return fmt.Errorf("failed to insert education change: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info("Expert edit request updated successfully: ID: %d", req.ID)
	return nil
}

// ApplyExpertEditRequest applies the changes from an edit request to the expert profile
func (s *SQLiteStore) ApplyExpertEditRequest(id int64, adminUserID int64) error {
	log := logger.Get()
	log.Debug("Applying expert edit request: ID: %d", id)

	// Get the edit request
	editRequest, err := s.GetExpertEditRequest(id)
	if err != nil {
		return fmt.Errorf("failed to get edit request: %w", err)
	}

	if editRequest.Status != "approved" {
		return fmt.Errorf("edit request must be approved before applying")
	}

	// Get the current expert
	expert, err := s.GetExpert(editRequest.ExpertID)
	if err != nil {
		return fmt.Errorf("failed to get expert: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Apply changes to expert
	s.applyChangesToExpert(expert, editRequest)

	// Update expert in database
	if err := s.updateExpertInTransaction(tx, expert); err != nil {
		return fmt.Errorf("failed to update expert: %w", err)
	}

	// Apply experience changes
	if err := s.applyExperienceChanges(tx, editRequest.ExpertID, editRequest.ExperienceChanges); err != nil {
		return fmt.Errorf("failed to apply experience changes: %w", err)
	}

	// Apply education changes
	if err := s.applyEducationChanges(tx, editRequest.ExpertID, editRequest.EducationChanges); err != nil {
		return fmt.Errorf("failed to apply education changes: %w", err)
	}

	// Mark edit request as applied
	if _, err := tx.Exec("UPDATE expert_edit_requests SET applied_at = ? WHERE id = ?", time.Now(), id); err != nil {
		return fmt.Errorf("failed to mark edit request as applied: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info("Expert edit request applied successfully: ID: %d", id)
	return nil
}

// CancelExpertEditRequest cancels an expert edit request
func (s *SQLiteStore) CancelExpertEditRequest(id int64, userID int64) error {
	log := logger.Get()
	log.Debug("Cancelling expert edit request: ID: %d", id)

	query := `UPDATE expert_edit_requests SET status = 'cancelled' WHERE id = ? AND created_by = ? AND status = 'pending'`

	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		log.Error("Failed to cancel expert edit request: %v", err)
		return fmt.Errorf("failed to cancel expert edit request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("edit request not found or cannot be cancelled")
	}

	log.Info("Expert edit request cancelled: ID: %d", id)
	return nil
}

// Helper functions

func (s *SQLiteStore) scanExpertEditRequest(row interface{ Scan(dest ...interface{}) error }) (*domain.ExpertEditRequest, error) {
	var req domain.ExpertEditRequest
	var fieldsChangedJSON, suggestedAreasJSON string
	var createdAt, reviewedAt, appliedAt sql.NullTime
	
	// Use NullString and NullBool for nullable fields
	var name, designation, institution, phone, email, role, employmentType, specializedArea, biography sql.NullString
	var newCVPath, newApprovalDocumentPath sql.NullString
	var isBahraini, isAvailable, isTrained, isPublished sql.NullBool
	var rating sql.NullInt64
	var generalArea sql.NullInt64

	err := row.Scan(
		&req.ID, &req.ExpertID,
		&name, &designation, &institution, &phone, &email,
		&isBahraini, &isAvailable, &rating,
		&role, &employmentType, &generalArea, &specializedArea,
		&isTrained, &isPublished, &biography,
		&suggestedAreasJSON,
		&newCVPath, &newApprovalDocumentPath,
		&req.RemoveCV, &req.RemoveApprovalDocument,
		&req.ChangeSummary, &req.ChangeReason, &fieldsChangedJSON,
		&req.Status, &req.RejectionReason, &req.AdminNotes,
		&createdAt, &reviewedAt, &appliedAt,
		&req.CreatedBy, &req.ReviewedBy,
		&req.ExpertName, &req.CreatedByName, &req.ReviewedByName,
	)
	if err != nil {
		return nil, err
	}
	
	// Convert nullable fields to pointers
	if name.Valid {
		req.Name = &name.String
	}
	if designation.Valid {
		req.Designation = &designation.String
	}
	if institution.Valid {
		req.Institution = &institution.String
	}
	if phone.Valid {
		req.Phone = &phone.String
	}
	if email.Valid {
		req.Email = &email.String
	}
	if isBahraini.Valid {
		req.IsBahraini = &isBahraini.Bool
	}
	if isAvailable.Valid {
		req.IsAvailable = &isAvailable.Bool
	}
	if rating.Valid {
		ratingInt := int(rating.Int64)
		req.Rating = &ratingInt
	}
	if role.Valid {
		req.Role = &role.String
	}
	if employmentType.Valid {
		req.EmploymentType = &employmentType.String
	}
	if generalArea.Valid {
		req.GeneralArea = &generalArea.Int64
	}
	if specializedArea.Valid {
		req.SpecializedArea = &specializedArea.String
	}
	if isTrained.Valid {
		req.IsTrained = &isTrained.Bool
	}
	if isPublished.Valid {
		req.IsPublished = &isPublished.Bool
	}
	if biography.Valid {
		req.Biography = &biography.String
	}
	if newCVPath.Valid {
		req.NewCVPath = &newCVPath.String
	}
	if newApprovalDocumentPath.Valid {
		req.NewApprovalDocumentPath = &newApprovalDocumentPath.String
	}

	// Handle timestamps
	if createdAt.Valid {
		req.CreatedAt = createdAt.Time
	}
	if reviewedAt.Valid {
		req.ReviewedAt = reviewedAt.Time
	}
	if appliedAt.Valid {
		req.AppliedAt = appliedAt.Time
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal([]byte(fieldsChangedJSON), &req.FieldsChanged); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields changed: %w", err)
	}

	if err := json.Unmarshal([]byte(suggestedAreasJSON), &req.SuggestedSpecializedAreas); err != nil {
		return nil, fmt.Errorf("failed to unmarshal suggested areas: %w", err)
	}

	return &req, nil
}

func (s *SQLiteStore) loadEditRequestChanges(req *domain.ExpertEditRequest) error {
	// Load experience changes
	expQuery := `
		SELECT id, action, experience_id, organization, position, start_date, end_date, is_current, country, description
		FROM expert_edit_request_experience WHERE expert_edit_request_id = ?
		ORDER BY id`

	expRows, err := s.db.Query(expQuery, req.ID)
	if err != nil {
		return fmt.Errorf("failed to query experience changes: %w", err)
	}
	defer expRows.Close()

	req.ExperienceChanges = []domain.ExpertEditRequestExperienceEntry{}
	for expRows.Next() {
		var exp domain.ExpertEditRequestExperienceEntry
		err := expRows.Scan(
			&exp.ID, &exp.Action, &exp.ExperienceID,
			&exp.Organization, &exp.Position, &exp.StartDate, &exp.EndDate,
			&exp.IsCurrent, &exp.Country, &exp.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to scan experience change: %w", err)
		}
		exp.EditRequestID = req.ID
		req.ExperienceChanges = append(req.ExperienceChanges, exp)
	}

	// Load education changes
	eduQuery := `
		SELECT id, action, education_id, institution, degree, field_of_study, graduation_year, country, description
		FROM expert_edit_request_education WHERE expert_edit_request_id = ?
		ORDER BY id`

	eduRows, err := s.db.Query(eduQuery, req.ID)
	if err != nil {
		return fmt.Errorf("failed to query education changes: %w", err)
	}
	defer eduRows.Close()

	req.EducationChanges = []domain.ExpertEditRequestEducationEntry{}
	for eduRows.Next() {
		var edu domain.ExpertEditRequestEducationEntry
		err := eduRows.Scan(
			&edu.ID, &edu.Action, &edu.EducationID,
			&edu.Institution, &edu.Degree, &edu.FieldOfStudy, &edu.GraduationYear,
			&edu.Country, &edu.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to scan education change: %w", err)
		}
		edu.EditRequestID = req.ID
		req.EducationChanges = append(req.EducationChanges, edu)
	}

	return nil
}

func (s *SQLiteStore) insertExperienceChange(tx *sql.Tx, exp *domain.ExpertEditRequestExperienceEntry) error {
	query := `
		INSERT INTO expert_edit_request_experience (
			expert_edit_request_id, action, experience_id, organization, position,
			start_date, end_date, is_current, country, description
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.Exec(query,
		exp.EditRequestID, exp.Action, exp.ExperienceID,
		exp.Organization, exp.Position, exp.StartDate, exp.EndDate,
		exp.IsCurrent, exp.Country, exp.Description,
	)
	return err
}

func (s *SQLiteStore) insertEducationChange(tx *sql.Tx, edu *domain.ExpertEditRequestEducationEntry) error {
	query := `
		INSERT INTO expert_edit_request_education (
			expert_edit_request_id, action, education_id, institution, degree,
			field_of_study, graduation_year, country, description
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.Exec(query,
		edu.EditRequestID, edu.Action, edu.EducationID,
		edu.Institution, edu.Degree, edu.FieldOfStudy, edu.GraduationYear,
		edu.Country, edu.Description,
	)
	return err
}

// Utility functions for handling nullable types
func stringPtrToInterface(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func boolPtrToInterface(b *bool) interface{} {
	if b == nil {
		return nil
	}
	return *b
}

func intPtrToInterface(i *int) interface{} {
	if i == nil {
		return nil
	}
	return *i
}

func int64PtrToInterface(i *int64) interface{} {
	if i == nil {
		return nil
	}
	return *i
}


// Helper functions to apply changes to expert
func (s *SQLiteStore) applyChangesToExpert(expert *domain.Expert, editRequest *domain.ExpertEditRequest) {
	if editRequest.Name != nil {
		expert.Name = *editRequest.Name
	}
	if editRequest.Designation != nil {
		expert.Designation = *editRequest.Designation
	}
	if editRequest.Institution != nil {
		expert.Affiliation = *editRequest.Institution
	}
	if editRequest.Phone != nil {
		expert.Phone = *editRequest.Phone
	}
	if editRequest.Email != nil {
		expert.Email = *editRequest.Email
	}
	if editRequest.IsBahraini != nil {
		expert.IsBahraini = *editRequest.IsBahraini
	}
	if editRequest.IsAvailable != nil {
		expert.IsAvailable = *editRequest.IsAvailable
	}
	if editRequest.Rating != nil {
		expert.Rating = *editRequest.Rating
	}
	if editRequest.Role != nil {
		expert.Role = *editRequest.Role
	}
	if editRequest.EmploymentType != nil {
		expert.EmploymentType = *editRequest.EmploymentType
	}
	if editRequest.GeneralArea != nil {
		expert.GeneralArea = *editRequest.GeneralArea
	}
	if editRequest.SpecializedArea != nil {
		expert.SpecializedArea = *editRequest.SpecializedArea
	}
	if editRequest.IsTrained != nil {
		expert.IsTrained = *editRequest.IsTrained
	}
	if editRequest.IsPublished != nil {
		expert.IsPublished = *editRequest.IsPublished
	}
	if editRequest.Biography != nil {
		// Note: Assuming expert has a Biography field - if not, this would need adjustment
	}

	// Handle document changes
	if editRequest.NewCVPath != nil {
		expert.CVPath = *editRequest.NewCVPath
	} else if editRequest.RemoveCV {
		expert.CVPath = ""
	}

	if editRequest.NewApprovalDocumentPath != nil {
		expert.ApprovalDocumentPath = *editRequest.NewApprovalDocumentPath
	} else if editRequest.RemoveApprovalDocument {
		expert.ApprovalDocumentPath = ""
	}

	expert.UpdatedAt = time.Now()
}

func (s *SQLiteStore) updateExpertInTransaction(tx *sql.Tx, expert *domain.Expert) error {
	query := `
		UPDATE experts SET
			name = ?, designation = ?, institution = ?, phone = ?, email = ?,
			is_bahraini = ?, is_available = ?, rating = ?, role = ?, employment_type = ?,
			general_area = ?, specialized_area = ?, is_trained = ?, is_published = ?,
			cv_path = ?, approval_document_path = ?, updated_at = ?
		WHERE id = ?`

	_, err := tx.Exec(query,
		expert.Name, expert.Designation, expert.Affiliation, expert.Phone, expert.Email,
		expert.IsBahraini, expert.IsAvailable, expert.Rating, expert.Role, expert.EmploymentType,
		expert.GeneralArea, expert.SpecializedArea, expert.IsTrained, expert.IsPublished,
		expert.CVPath, expert.ApprovalDocumentPath, expert.UpdatedAt,
		expert.ID,
	)
	return err
}

func (s *SQLiteStore) applyExperienceChanges(tx *sql.Tx, expertID int64, changes []domain.ExpertEditRequestExperienceEntry) error {
	for _, change := range changes {
		switch change.Action {
		case "add":
			query := `
				INSERT INTO expert_experience (
					expert_id, organization, position, start_date, end_date,
					is_current, country, description, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
			
			now := time.Now()
			_, err := tx.Exec(query,
				expertID, change.Organization, change.Position, change.StartDate, change.EndDate,
				change.IsCurrent, change.Country, change.Description, now, now,
			)
			if err != nil {
				return fmt.Errorf("failed to add experience: %w", err)
			}

		case "update":
			query := `
				UPDATE expert_experience SET
					organization = ?, position = ?, start_date = ?, end_date = ?,
					is_current = ?, country = ?, description = ?, updated_at = ?
				WHERE id = ? AND expert_id = ?`
			
			_, err := tx.Exec(query,
				change.Organization, change.Position, change.StartDate, change.EndDate,
				change.IsCurrent, change.Country, change.Description, time.Now(),
				change.ExperienceID, expertID,
			)
			if err != nil {
				return fmt.Errorf("failed to update experience: %w", err)
			}

		case "delete":
			query := `DELETE FROM expert_experience WHERE id = ? AND expert_id = ?`
			_, err := tx.Exec(query, change.ExperienceID, expertID)
			if err != nil {
				return fmt.Errorf("failed to delete experience: %w", err)
			}
		}
	}
	return nil
}

func (s *SQLiteStore) applyEducationChanges(tx *sql.Tx, expertID int64, changes []domain.ExpertEditRequestEducationEntry) error {
	for _, change := range changes {
		switch change.Action {
		case "add":
			query := `
				INSERT INTO expert_education (
					expert_id, institution, degree, field_of_study, graduation_year,
					country, description, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
			
			now := time.Now()
			_, err := tx.Exec(query,
				expertID, change.Institution, change.Degree, change.FieldOfStudy, change.GraduationYear,
				change.Country, change.Description, now, now,
			)
			if err != nil {
				return fmt.Errorf("failed to add education: %w", err)
			}

		case "update":
			query := `
				UPDATE expert_education SET
					institution = ?, degree = ?, field_of_study = ?, graduation_year = ?,
					country = ?, description = ?, updated_at = ?
				WHERE id = ? AND expert_id = ?`
			
			_, err := tx.Exec(query,
				change.Institution, change.Degree, change.FieldOfStudy, change.GraduationYear,
				change.Country, change.Description, time.Now(),
				change.EducationID, expertID,
			)
			if err != nil {
				return fmt.Errorf("failed to update education: %w", err)
			}

		case "delete":
			query := `DELETE FROM expert_education WHERE id = ? AND expert_id = ?`
			_, err := tx.Exec(query, change.EducationID, expertID)
			if err != nil {
				return fmt.Errorf("failed to delete education: %w", err)
			}
		}
	}
	return nil
}