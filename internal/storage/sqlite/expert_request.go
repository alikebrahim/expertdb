package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"expertdb/internal/domain"
	"expertdb/internal/logger"
)


// CreateExpertRequest creates a new expert request in the database
func (s *SQLiteStore) CreateExpertRequest(req *domain.ExpertRequest) (int64, error) {
	log := logger.Get()
	log.Debug("Creating expert request for: %s", req.Name)
	
	query := `
		INSERT INTO expert_requests (
			name, designation, affiliation, is_bahraini, is_available,
			role, employment_type, general_area, specialized_area, is_trained,
			cv_document_id, approval_document_id, phone, email, is_published, 
			suggested_specialized_areas, status, created_at, created_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	// Convert specialized areas to JSON string for storage
	suggestedAreasJSON := s.serializeSuggestedAreas(req.SuggestedSpecializedAreas)
	
	// Set document IDs (initially nil)
	var cvDocumentID, approvalDocumentID *int64
	
	result, err := s.db.Exec(query,
		req.Name, req.Designation, req.Affiliation, req.IsBahraini, req.IsAvailable,
		req.Role, req.EmploymentType, req.GeneralArea, req.SpecializedArea, req.IsTrained,
		cvDocumentID, approvalDocumentID, req.Phone, req.Email, req.IsPublished, 
		suggestedAreasJSON, req.Status, req.CreatedAt, req.CreatedBy,
	)
	if err != nil {
		log.Error("Failed to create expert request: %v", err)
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		log.Error("Failed to get last insert ID: %v", err)
		return 0, err
	}
	
	// Set the request ID
	req.ID = id
	
	// Store experience entries
	for _, entry := range req.ExperienceEntries {
		_, err := s.db.Exec(`
			INSERT INTO expert_request_experience_entries (
				expert_request_id, organization, position, start_date, end_date, is_current, country, description
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			id, entry.Organization, entry.Position, entry.StartDate, entry.EndDate, entry.IsCurrent, entry.Country, entry.Description,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create experience entry: %w", err)
		}
	}
	
	// Store education entries
	for _, entry := range req.EducationEntries {
		_, err := s.db.Exec(`
			INSERT INTO expert_request_education_entries (
				expert_request_id, institution, degree, field_of_study, graduation_year, country, description
			) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			id, entry.Institution, entry.Degree, entry.FieldOfStudy, entry.GraduationYear, entry.Country, entry.Description,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create education entry: %w", err)
		}
	}
	
	log.Debug("Expert request created successfully with ID: %d", id)
	return id, nil
}


// GetExpertRequest retrieves an expert request by ID
func (s *SQLiteStore) GetExpertRequest(id int64) (*domain.ExpertRequest, error) {
	query := `
		SELECT 
			id, name, designation, affiliation, is_bahraini, 
			is_available, role, employment_type, general_area, 
			specialized_area, is_trained, cv_document_id, approval_document_id, phone, email, 
			is_published, suggested_specialized_areas, status, rejection_reason, 
			created_at, reviewed_at, reviewed_by, created_by
		FROM expert_requests
		WHERE id = ?
	`
	
	var req domain.ExpertRequest
	var reviewedAt sql.NullTime
	var reviewedBy sql.NullInt64
	var createdBy sql.NullInt64
	var rejectionReason sql.NullString
	var suggestedAreasJSON string
	var cvDocumentID sql.NullInt64
	var approvalDocumentID sql.NullInt64
	
	err := s.db.QueryRow(query, id).Scan(
		&req.ID, &req.Name, &req.Designation, &req.Affiliation, 
		&req.IsBahraini, &req.IsAvailable, &req.Role, 
		&req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
		&req.IsTrained, &cvDocumentID, &approvalDocumentID, &req.Phone, &req.Email, 
		&req.IsPublished, &suggestedAreasJSON, &req.Status, &rejectionReason, 
		&req.CreatedAt, &reviewedAt, &reviewedBy, &createdBy,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get expert request: %w", err)
	}
	
	// Set nullable fields
	
	if reviewedAt.Valid {
		req.ReviewedAt = reviewedAt.Time
	}
	
	if reviewedBy.Valid {
		req.ReviewedBy = reviewedBy.Int64
	}
	
	if createdBy.Valid {
		req.CreatedBy = createdBy.Int64
	}
	
	if rejectionReason.Valid {
		req.RejectionReason = rejectionReason.String
	}
	
	if cvDocumentID.Valid {
		req.CVDocumentID = &cvDocumentID.Int64
	}
	
	if approvalDocumentID.Valid {
		req.ApprovalDocumentID = &approvalDocumentID.Int64
	}
	
	// Resolve document references
	req.ResolveCVDocument(s.GetDocument)
	req.ResolveApprovalDocument(s.GetDocument)
	
	// Deserialize suggested areas from JSON
	req.SuggestedSpecializedAreas = s.deserializeSuggestedAreas(suggestedAreasJSON)
	
	// Populate experience entries
	experienceEntries, err := s.getExpertRequestExperienceEntries(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experience entries: %w", err)
	}
	req.ExperienceEntries = experienceEntries
	
	// Populate education entries
	educationEntries, err := s.getExpertRequestEducationEntries(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education entries: %w", err)
	}
	req.EducationEntries = educationEntries
	
	return &req, nil
}

// ListExpertRequests retrieves a list of expert requests with the given status
func (s *SQLiteStore) ListExpertRequests(status string, limit, offset int) ([]*domain.ExpertRequest, error) {
	if limit <= 0 {
		limit = 10
	}
	
	var query string
	var args []interface{}
	
	if status != "" && status != "all" {
		query = `
			SELECT 
				id, name, designation, affiliation, is_bahraini, 
				is_available, role, employment_type, general_area, 
				specialized_area, is_trained, cv_document_id, approval_document_id, phone, email, 
				is_published, suggested_specialized_areas, status, rejection_reason, 
				created_at, reviewed_at, reviewed_by, created_by
			FROM expert_requests
			WHERE status = ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{status, limit, offset}
	} else {
		// status is empty or "all" - return all requests
		query = `
			SELECT 
				id, name, designation, affiliation, is_bahraini, 
				is_available, role, employment_type, general_area, 
				specialized_area, is_trained, cv_document_id, approval_document_id, phone, email, 
				is_published, suggested_specialized_areas, status, rejection_reason, 
				created_at, reviewed_at, reviewed_by, created_by
			FROM expert_requests
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{limit, offset}
	}
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query expert requests: %w", err)
	}
	defer rows.Close()
	
	var requests []*domain.ExpertRequest
	for rows.Next() {
		var req domain.ExpertRequest
		var reviewedAt sql.NullTime
		var reviewedBy sql.NullInt64
		var createdBy sql.NullInt64
		var specializedArea sql.NullString
		var cvPath sql.NullInt64
		var approvalDocPath sql.NullInt64
		var rejectionReason sql.NullString
		var suggestedAreasJSON sql.NullString
		
		err := rows.Scan(
			&req.ID, &req.Name, &req.Designation, &req.Affiliation, 
			&req.IsBahraini, &req.IsAvailable, &req.Role, 
			&req.EmploymentType, &req.GeneralArea, &specializedArea, 
			&req.IsTrained, &cvPath, &approvalDocPath, &req.Phone, &req.Email, 
			&req.IsPublished, &suggestedAreasJSON, &req.Status, &rejectionReason, 
			&req.CreatedAt, &reviewedAt, &reviewedBy, &createdBy,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan expert request row: %w", err)
		}
		
		// Handle nullable fields
		if reviewedAt.Valid {
			req.ReviewedAt = reviewedAt.Time
		}
		if reviewedBy.Valid {
			req.ReviewedBy = reviewedBy.Int64
		}
		if createdBy.Valid {
			req.CreatedBy = createdBy.Int64
		}
		if specializedArea.Valid {
			req.SpecializedArea = specializedArea.String
		}
		if cvPath.Valid {
			cvDocID := cvPath.Int64
			req.CVDocumentID = &cvDocID
		}
		if approvalDocPath.Valid {
			approvalDocID := approvalDocPath.Int64
			req.ApprovalDocumentID = &approvalDocID
		}
		if rejectionReason.Valid {
			req.RejectionReason = rejectionReason.String
		}
		
		// Deserialize suggested areas from JSON
		if suggestedAreasJSON.Valid {
			req.SuggestedSpecializedAreas = s.deserializeSuggestedAreas(suggestedAreasJSON.String)
		}
		
		// Populate experience entries
		experienceEntries, err := s.getExpertRequestExperienceEntries(req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get experience entries for request %d: %w", req.ID, err)
		}
		req.ExperienceEntries = experienceEntries
		
		// Populate education entries
		educationEntries, err := s.getExpertRequestEducationEntries(req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get education entries for request %d: %w", req.ID, err)
		}
		req.EducationEntries = educationEntries
		
		requests = append(requests, &req)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expert request rows: %w", err)
	}
	
	return requests, nil
}

// ListExpertRequestsByUser retrieves a list of expert requests created by a specific user
func (s *SQLiteStore) ListExpertRequestsByUser(userID int64, status string, limit, offset int) ([]*domain.ExpertRequest, error) {
	if limit <= 0 {
		limit = 10
	}
	
	var query string
	var args []interface{}
	
	if status != "" && status != "all" {
		query = `
			SELECT 
				id, name, designation, affiliation, is_bahraini, 
				is_available, role, employment_type, general_area, 
				specialized_area, is_trained, cv_document_id, approval_document_id, phone, email, 
				is_published, suggested_specialized_areas, status, rejection_reason, 
				created_at, reviewed_at, reviewed_by, created_by
			FROM expert_requests
			WHERE created_by = ? AND status = ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userID, status, limit, offset}
	} else {
		// status is empty or "all" - return all requests for this user
		query = `
			SELECT 
				id, name, designation, affiliation, is_bahraini, 
				is_available, role, employment_type, general_area, 
				specialized_area, is_trained, cv_document_id, approval_document_id, phone, email, 
				is_published, suggested_specialized_areas, status, rejection_reason, 
				created_at, reviewed_at, reviewed_by, created_by
			FROM expert_requests
			WHERE created_by = ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{userID, limit, offset}
	}
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query user expert requests: %w", err)
	}
	defer rows.Close()
	
	var requests []*domain.ExpertRequest
	for rows.Next() {
		var req domain.ExpertRequest
		var reviewedAt sql.NullTime
		var reviewedBy sql.NullInt64
		var createdBy sql.NullInt64
		var specializedArea sql.NullString
		var cvPath sql.NullInt64
		var approvalDocPath sql.NullInt64
		var rejectionReason sql.NullString
		var suggestedAreasJSON sql.NullString
		
		err := rows.Scan(
			&req.ID, &req.Name, &req.Designation, &req.Affiliation, &req.IsBahraini,
			&req.IsAvailable, &req.Role, &req.EmploymentType, &req.GeneralArea,
			&specializedArea, &req.IsTrained, &cvPath, &approvalDocPath, &req.Phone, &req.Email,
			&req.IsPublished, &suggestedAreasJSON, &req.Status, &rejectionReason,
			&req.CreatedAt, &reviewedAt, &reviewedBy, &createdBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expert request row: %w", err)
		}
		
		// Handle nullable fields
		if reviewedAt.Valid {
			req.ReviewedAt = reviewedAt.Time
		}
		if reviewedBy.Valid {
			req.ReviewedBy = reviewedBy.Int64
		}
		if createdBy.Valid {
			req.CreatedBy = createdBy.Int64
		}
		if specializedArea.Valid {
			req.SpecializedArea = specializedArea.String
		}
		if cvPath.Valid {
			cvDocID := cvPath.Int64
			req.CVDocumentID = &cvDocID
		}
		if approvalDocPath.Valid {
			approvalDocID := approvalDocPath.Int64
			req.ApprovalDocumentID = &approvalDocID
		}
		if rejectionReason.Valid {
			req.RejectionReason = rejectionReason.String
		}
		
		// Deserialize suggested areas from JSON
		if suggestedAreasJSON.Valid {
			req.SuggestedSpecializedAreas = s.deserializeSuggestedAreas(suggestedAreasJSON.String)
		}
		
		// Populate experience entries
		experienceEntries, err := s.getExpertRequestExperienceEntries(req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get experience entries for request %d: %w", req.ID, err)
		}
		req.ExperienceEntries = experienceEntries
		
		// Populate education entries
		educationEntries, err := s.getExpertRequestEducationEntries(req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get education entries for request %d: %w", req.ID, err)
		}
		req.EducationEntries = educationEntries
		
		requests = append(requests, &req)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user expert request rows: %w", err)
	}
	
	return requests, nil
}

// UpdateExpertRequestStatus updates the status of an expert request
func (s *SQLiteStore) UpdateExpertRequestStatus(id int64, status, rejectionReason string, reviewedBy int64) error {
	log := logger.Get()
	log.Debug("Updating expert request status: ID=%d, status='%s'", id, status)
	query := `
		UPDATE expert_requests
		SET status = ?, rejection_reason = ?, reviewed_at = ?, reviewed_by = ?
		WHERE id = ?
	`
	
	now := time.Now()
	
	// Execute the update
	result, err := s.db.Exec(query, status, rejectionReason, now, reviewedBy, id)
	if err != nil {
		return fmt.Errorf("failed to update expert request status: %w", err)
	}
	
	// Check if a row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	// Note: Expert creation for approvals is now handled by ApproveExpertRequestWithDocument method
	// This method is only used for rejections and simple status updates
	
	return nil
}

// UpdateExpertRequest updates an expert request with new data
func (s *SQLiteStore) UpdateExpertRequest(req *domain.ExpertRequest) error {
	log := logger.Get()
	log.Debug("Updating expert request: ID=%d", req.ID)
	
	// CRITICAL FIX: Prevent data corruption from empty form submissions
	// If core fields are empty (except for status updates during approval), reject the update
	if req.Name == "" && req.Email == "" && req.Phone == "" && req.Designation == "" && req.Affiliation == "" {
		// This indicates an empty form submission that would corrupt data
		// Only allow if it's just a status change (status field is not empty)
		if req.Status == "" {
			log.Warn("Rejected UpdateExpertRequest with all empty core fields for ID: %d", req.ID)
			return fmt.Errorf("cannot update expert request with empty core data")
		}
	}
	query := `
		UPDATE expert_requests
		SET name = ?, designation = ?, affiliation = ?, is_bahraini = ?,
			is_available = ?, role = ?, employment_type = ?,
			general_area = ?, specialized_area = ?, is_trained = ?,
			cv_document_id = ?, approval_document_id = ?, phone = ?, email = ?, is_published = ?,
			suggested_specialized_areas = ?, status = ?, rejection_reason = ?,
			reviewed_at = ?, reviewed_by = ?, created_by = ?
		WHERE id = ?
	`
	
	// Handle nullable fields
	var specializedArea, cvDocumentID, approvalDocumentID, rejectionReason interface{} = nil, nil, nil, nil
	if req.SpecializedArea != "" {
		specializedArea = req.SpecializedArea
	}
	if req.CVDocumentID != nil {
		cvDocumentID = *req.CVDocumentID
	}
	if req.ApprovalDocumentID != nil {
		approvalDocumentID = *req.ApprovalDocumentID
	}
	if req.RejectionReason != "" {
		rejectionReason = req.RejectionReason
	}
	
	var reviewedAt interface{} = nil
	if !req.ReviewedAt.IsZero() {
		reviewedAt = req.ReviewedAt
	}
	
	var reviewedBy interface{} = nil
	if req.ReviewedBy != 0 {
		reviewedBy = req.ReviewedBy
	}
	
	var createdBy interface{} = nil
	if req.CreatedBy != 0 {
		createdBy = req.CreatedBy
	}
	
	// Serialize suggested areas to JSON
	suggestedAreasJSON := s.serializeSuggestedAreas(req.SuggestedSpecializedAreas)
	
	// Execute update
	result, err := s.db.Exec(
		query,
		req.Name, req.Designation, req.Affiliation, req.IsBahraini,
		req.IsAvailable, req.Role, req.EmploymentType,
		req.GeneralArea, specializedArea, req.IsTrained,
		cvDocumentID, approvalDocumentID, req.Phone, req.Email, req.IsPublished,
		suggestedAreasJSON, req.Status, rejectionReason,
		reviewedAt, reviewedBy, createdBy,
		req.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update expert request: %w", err)
	}
	
	// Check if the update affected a row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	return nil
}


// BatchApproveExpertRequestsWithFileMove approves multiple expert requests with file moving
func (s *SQLiteStore) BatchApproveExpertRequestsWithFileMove(requestIDs []int64, reviewedBy int64, documentService interface{}) ([]int64, []int64, map[int64]error) {
	log := logger.Get()
	log.Debug("Batch approving %d expert requests with file move", len(requestIDs))
	
	successIDs := []int64{}
	expertIDs := []int64{}
	errors := make(map[int64]error)
	
	// Note: Document service no longer needed for file moving with document-centric approach
	// Documents are managed through foreign key relationships in expert_documents table
	
	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		for _, id := range requestIDs {
			errors[id] = fmt.Errorf("failed to begin transaction: %w", err)
		}
		return successIDs, expertIDs, errors
	}
	defer tx.Rollback()
	
	now := time.Now()
	
	// Process each request
	for _, requestID := range requestIDs {
		log.Debug("Processing request ID: %d", requestID)
		log.Debug("DEBUG: Starting processing for request ID: %d", requestID)
		
		// Step 1: Get the request data
		var req domain.ExpertRequest
		var cvDocumentID, approvalDocumentID sql.NullInt64
		query := `
			SELECT id, name, designation, affiliation, is_bahraini, 
				is_available, role, employment_type, general_area, 
				specialized_area, is_trained, cv_document_id, approval_document_id,
				phone, email, is_published, status, created_by
			FROM expert_requests
			WHERE id = ? AND status = 'pending'
		`
		
		log.Debug("DEBUG: Executing query to get request data for ID: %d", requestID)
		err = tx.QueryRow(query, requestID).Scan(
			&req.ID, &req.Name, &req.Designation, &req.Affiliation, 
			&req.IsBahraini, &req.IsAvailable, &req.Role, 
			&req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
			&req.IsTrained, &cvDocumentID, &approvalDocumentID, &req.Phone, &req.Email, 
			&req.IsPublished, &req.Status, &req.CreatedBy,
		)
		log.Debug("DEBUG: Request data scan completed - Name: '%s', Email: '%s', Phone: '%s', Designation: '%s', Affiliation: '%s'", 
			req.Name, req.Email, req.Phone, req.Designation, req.Affiliation)
		
		if err != nil {
			if err == sql.ErrNoRows {
				errors[requestID] = fmt.Errorf("request not found or not pending")
			} else {
				errors[requestID] = fmt.Errorf("failed to retrieve request data: %w", err)
			}
			continue
		}
		
		// Assign document IDs if valid
		if cvDocumentID.Valid {
			req.CVDocumentID = &cvDocumentID.Int64
		}
		if approvalDocumentID.Valid {
			req.ApprovalDocumentID = &approvalDocumentID.Int64
		}
		
		// Step 2: Create expert record
		log.Debug("DEBUG: Creating expert record from request data for ID: %d", requestID)
		expert := &domain.Expert{
			Name:            req.Name,
			Email:           req.Email,
			Phone:           req.Phone,
			Designation:     req.Designation,
			Affiliation:     req.Affiliation,
			IsBahraini:      req.IsBahraini,
			IsAvailable:     req.IsAvailable,
			Role:            req.Role,
			EmploymentType:  req.EmploymentType,
			GeneralArea:     req.GeneralArea,
			SpecializedArea: req.SpecializedArea,
			IsPublished:     req.IsPublished,
			IsTrained:       req.IsTrained,
			CreatedAt:       now,
			UpdatedAt:       now,
			OriginalRequestID: requestID,
		}
		log.Debug("DEBUG: Expert record created - Name: '%s', Email: '%s', Phone: '%s'", expert.Name, expert.Email, expert.Phone)
		
		// Insert expert
		log.Debug("DEBUG: Inserting expert record into database for request ID: %d", requestID)
		result, err := tx.Exec(`
			INSERT INTO experts (
				name, designation, affiliation, is_bahraini, is_available,
				rating, role, employment_type, general_area, specialized_area,
				is_trained, cv_document_id, approval_document_id, phone, email, is_published,
				created_at, updated_at, original_request_id
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			expert.Name, expert.Designation, expert.Affiliation,
			expert.IsBahraini, expert.IsAvailable, expert.Rating, expert.Role,
			expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
			expert.IsTrained, req.CVDocumentID, req.ApprovalDocumentID, expert.Phone, expert.Email, expert.IsPublished,
			expert.CreatedAt, expert.UpdatedAt, expert.OriginalRequestID,
		)
		log.Debug("DEBUG: Expert insert completed with error: %v", err)
		if err != nil {
			errors[requestID] = fmt.Errorf("failed to insert expert: %w", err)
			continue
		}
		
		expertIDResult, err := result.LastInsertId()
		if err != nil {
			errors[requestID] = fmt.Errorf("failed to get expert ID: %w", err)
			continue
		}
		expertID := expertIDResult
		
		log.Debug("Created expert with ID: %d for request: %d", expertID, requestID)
		
		// Step 3: Copy experience and education entries
		err = s.copyExperienceEntries(tx, requestID, expertID)
		if err != nil {
			errors[requestID] = fmt.Errorf("failed to copy experience entries: %w", err)
			continue
		}

		err = s.copyEducationEntries(tx, requestID, expertID)
		if err != nil {
			errors[requestID] = fmt.Errorf("failed to copy education entries: %w", err)
			continue
		}
		
		// Step 4: Update request status
		log.Debug("DEBUG: Updating request status to approved for ID: %d", requestID)
		_, err = tx.Exec(`
			UPDATE expert_requests
			SET status = ?, reviewed_at = ?, reviewed_by = ?
			WHERE id = ?
		`, "approved", now, reviewedBy, requestID)
		log.Debug("DEBUG: Request status update completed with error: %v", err)
		if err != nil {
			errors[requestID] = fmt.Errorf("failed to update request status: %w", err)
			continue
		}
		
		// This request was successful
		successIDs = append(successIDs, requestID)
		expertIDs = append(expertIDs, expertID)
	}
	
	// Commit transaction if we have any successes
	if len(successIDs) > 0 {
		log.Debug("DEBUG: Committing transaction for %d successful requests", len(successIDs))
		if err := tx.Commit(); err != nil {
			log.Debug("DEBUG: Transaction commit failed: %v", err)
			// If commit fails, all operations fail
			for _, id := range successIDs {
				errors[id] = fmt.Errorf("failed to commit transaction: %w", err)
			}
			return []int64{}, []int64{}, errors
		}
		log.Debug("DEBUG: Transaction committed successfully")
		
		// Handle document movements for successful requests (after transaction committed)
		for i, requestID := range successIDs {
			expertID := expertIDs[i]
			
			// Get the CV document ID for this request
			var cvDocumentID sql.NullInt64
			err := s.db.QueryRow("SELECT cv_document_id FROM expert_requests WHERE id = ?", requestID).Scan(&cvDocumentID)
			if err != nil {
				log.Error("Failed to get CV document ID for request %d: %v", requestID, err)
				continue
			}
			
			// Move CV document using new method
			if cvDocumentID.Valid {
				if docService, ok := documentService.(interface {
					MoveRequestDocumentToExpert(documentID, expertID int64) error
				}); ok {
					err := docService.MoveRequestDocumentToExpert(cvDocumentID.Int64, expertID)
					if err != nil {
						log.Error("Failed to move CV document for request %d: %v", requestID, err)
					} else {
						log.Debug("Successfully moved CV document %d for request %d to expert %d", cvDocumentID.Int64, requestID, expertID)
					}
				}
			}
		}
		
		log.Debug("Document transfer completed for %d successful approvals", len(successIDs))
	} else {
		// No successful operations
		return []int64{}, []int64{}, errors
	}
	
	log.Debug("Batch approval completed: %d successes, %d errors", len(successIDs), len(errors))
	return successIDs, expertIDs, errors
}

// UpdateExpertsApprovalPath updates the approval document path for multiple experts
func (s *SQLiteStore) UpdateExpertsApprovalPath(expertIDs []int64, approvalPath string) error {
	log := logger.Get()
	log.Debug("Updating approval path for %d experts: %s", len(expertIDs), approvalPath)
	
	if len(expertIDs) == 0 {
		return nil
	}
	
	// Create placeholders for the IN clause
	placeholders := make([]string, len(expertIDs))
	args := make([]interface{}, len(expertIDs)+1)
	args[0] = approvalPath
	
	for i, id := range expertIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}
	
	query := fmt.Sprintf(`UPDATE experts SET approval_document_id = ? WHERE id IN (%s)`,
		strings.Join(placeholders, ","))
	
	result, err := s.db.Exec(query, args...)
	if err != nil {
		log.Error("Failed to update experts approval path: %v", err)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Failed to get rows affected: %v", err)
		return err
	}
	
	log.Debug("Updated approval path for %d experts", rowsAffected)
	return nil
}

// serializeSuggestedAreas converts a string slice to JSON for database storage
func (s *SQLiteStore) serializeSuggestedAreas(areas []string) string {
	if len(areas) == 0 {
		return ""
	}
	data, _ := json.Marshal(areas)
	return string(data)
}

// deserializeSuggestedAreas converts JSON string from database to string slice
func (s *SQLiteStore) deserializeSuggestedAreas(jsonData string) []string {
	if jsonData == "" {
		return []string{}
	}
	var areas []string
	json.Unmarshal([]byte(jsonData), &areas)
	return areas
}

// ApproveExpertRequestWithDocument approves a single expert request and creates expert with proper approval document naming
func (s *SQLiteStore) ApproveExpertRequestWithDocument(requestID, reviewedBy int64, documentService interface{}) (int64, error) {
	log := logger.Get()
	
	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	now := time.Now()
	
	// Step 1: Get the request data
	var req domain.ExpertRequest
	var cvDocumentID, approvalDocumentID sql.NullInt64
	query := `
		SELECT id, name, designation, affiliation, is_bahraini, 
			is_available, role, employment_type, general_area, 
			specialized_area, is_trained, cv_document_id, approval_document_id, 
			phone, email, is_published, status, created_by
		FROM expert_requests
		WHERE id = ? AND status = 'pending'
	`
	
	err = tx.QueryRow(query, requestID).Scan(
		&req.ID, &req.Name, &req.Designation, &req.Affiliation, 
		&req.IsBahraini, &req.IsAvailable, &req.Role, 
		&req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
		&req.IsTrained, &cvDocumentID, &approvalDocumentID, 
		&req.Phone, &req.Email, &req.IsPublished, &req.Status, &req.CreatedBy,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("request not found or not pending")
		}
		return 0, fmt.Errorf("failed to retrieve request data: %w", err)
	}
	
	// Assign document IDs if valid
	if cvDocumentID.Valid {
		req.CVDocumentID = &cvDocumentID.Int64
	}
	if approvalDocumentID.Valid {
		req.ApprovalDocumentID = &approvalDocumentID.Int64
	}
	
	// Step 2: Create expert record
	expert := &domain.Expert{
		Name:            req.Name,
		Email:           req.Email,
		Phone:           req.Phone,
		Designation:     req.Designation,
		Affiliation:     req.Affiliation,
		IsBahraini:      req.IsBahraini,
		IsAvailable:     req.IsAvailable,
		Role:            req.Role,
		EmploymentType:  req.EmploymentType,
		GeneralArea:     req.GeneralArea,
		SpecializedArea: req.SpecializedArea,
		IsPublished:     req.IsPublished,
		IsTrained:       req.IsTrained,
		CreatedAt:       now,
		UpdatedAt:       now,
		OriginalRequestID: requestID,
	}
	
	// Insert expert
	result, err := tx.Exec(`
		INSERT INTO experts (
			name, designation, affiliation, is_bahraini, is_available,
			rating, role, employment_type, general_area, specialized_area,
			is_trained, cv_document_id, approval_document_id, phone, email, is_published,
			created_at, updated_at, original_request_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		expert.Name, expert.Designation, expert.Affiliation,
		expert.IsBahraini, expert.IsAvailable, expert.Rating, expert.Role,
		expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, req.CVDocumentID, req.ApprovalDocumentID, expert.Phone, expert.Email, expert.IsPublished,
		expert.CreatedAt, expert.UpdatedAt, expert.OriginalRequestID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert expert: %w", err)
	}
	
	expertIDResult, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get expert ID: %w", err)
	}
	expertID := expertIDResult
	
	// Step 3: Copy experience and education entries
	err = s.copyExperienceEntries(tx, requestID, expertID)
	if err != nil {
		return 0, fmt.Errorf("failed to copy experience entries: %w", err)
	}

	err = s.copyEducationEntries(tx, requestID, expertID)
	if err != nil {
		return 0, fmt.Errorf("failed to copy education entries: %w", err)
	}
	
	// Step 4: Update request status
	_, err = tx.Exec(`
		UPDATE expert_requests
		SET status = ?, reviewed_at = ?, reviewed_by = ?
		WHERE id = ?
	`, "approved", now, reviewedBy, requestID)
	if err != nil {
		return 0, fmt.Errorf("failed to update request status: %w", err)
	}
	
	// Commit transaction
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	// Step 5: Handle document movements (outside transaction)
	
	// Move CV document from expert_requests to experts directory
	if req.CVDocumentID != nil {
		if docService, ok := documentService.(interface {
			MoveRequestDocumentToExpert(documentID, expertID int64) error
		}); ok {
			err := docService.MoveRequestDocumentToExpert(*req.CVDocumentID, expertID)
			if err != nil {
				log.Error("Failed to move CV document: %v", err)
				// Don't fail the entire operation, just log the error
			} else {
				log.Debug("Successfully moved CV document %d for expert %d", *req.CVDocumentID, expertID)
			}
		}
	}
	
	// Handle approval document with proper expert ID
	if req.ApprovalDocumentID != nil {
		if docService, ok := documentService.(interface {
			MoveRequestDocumentToExpert(documentID, expertID int64) error
		}); ok {
			err := docService.MoveRequestDocumentToExpert(*req.ApprovalDocumentID, expertID)
			if err != nil {
				log.Error("Failed to move approval document: %v", err)
				// Don't fail the entire operation, just log the error
			} else {
				log.Debug("Successfully moved approval document %d for expert %d", *req.ApprovalDocumentID, expertID)
			}
		}
	}
	
	return expertID, nil
}

// getExpertRequestExperienceEntries retrieves experience entries for an expert request
func (s *SQLiteStore) getExpertRequestExperienceEntries(requestID int64) ([]domain.ExpertRequestExperienceEntry, error) {
	query := `
		SELECT id, expert_request_id, organization, position, start_date, end_date, is_current, country, description
		FROM expert_request_experience_entries
		WHERE expert_request_id = ?
		ORDER BY start_date DESC
	`
	
	rows, err := s.db.Query(query, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to query experience entries: %w", err)
	}
	defer rows.Close()
	
	var entries []domain.ExpertRequestExperienceEntry
	for rows.Next() {
		var entry domain.ExpertRequestExperienceEntry
		err := rows.Scan(
			&entry.ID, &entry.ExpertRequestID, &entry.Organization, &entry.Position,
			&entry.StartDate, &entry.EndDate, &entry.IsCurrent, &entry.Country, &entry.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan experience entry: %w", err)
		}
		entries = append(entries, entry)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating experience entries: %w", err)
	}
	
	return entries, nil
}

// getExpertRequestEducationEntries retrieves education entries for an expert request
func (s *SQLiteStore) getExpertRequestEducationEntries(requestID int64) ([]domain.ExpertRequestEducationEntry, error) {
	query := `
		SELECT id, expert_request_id, institution, degree, field_of_study, graduation_year, country, description
		FROM expert_request_education_entries
		WHERE expert_request_id = ?
		ORDER BY graduation_year DESC
	`
	
	rows, err := s.db.Query(query, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to query education entries: %w", err)
	}
	defer rows.Close()
	
	var entries []domain.ExpertRequestEducationEntry
	for rows.Next() {
		var entry domain.ExpertRequestEducationEntry
		err := rows.Scan(
			&entry.ID, &entry.ExpertRequestID, &entry.Institution, &entry.Degree,
			&entry.FieldOfStudy, &entry.GraduationYear, &entry.Country, &entry.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan education entry: %w", err)
		}
		entries = append(entries, entry)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating education entries: %w", err)
	}
	
	return entries, nil
}

// moveCVFileToExpertDirectory moves CV file from expert_requests to cvs directory and renames with expert ID
func (s *SQLiteStore) moveCVFileToExpertDirectory(oldPath string, expertID int64) error {
	log := logger.Get()
	
	// Check if the old file exists
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		log.Warn("CV file does not exist at path: %s", oldPath)
		return fmt.Errorf("CV file not found: %s", oldPath)
	}
	
	// Extract file extension from the original path
	ext := filepath.Ext(oldPath)
	
	// Generate new filename with expert ID and timestamp
	timestamp := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("expert_%d_%s%s", expertID, timestamp, ext)
	
	// Construct new path in cvs directory
	newPath := filepath.Join("data", "documents", "cvs", newFileName)
	
	// Ensure the cvs directory exists
	cvsDir := filepath.Dir(newPath)
	if err := os.MkdirAll(cvsDir, 0755); err != nil {
		return fmt.Errorf("failed to create cvs directory: %w", err)
	}
	
	// Move the file
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to move CV file from %s to %s: %w", oldPath, newPath, err)
	}
	
	log.Debug("Successfully moved CV file from %s to %s for expert %d", oldPath, newPath, expertID)
	
	// Note: Database path update no longer needed with document-centric approach
	// Documents are managed through expert_documents table with foreign keys
	
	log.Debug("Successfully updated expert %d CV path to %s", expertID, newPath)
	return nil
}

// copyExperienceEntries copies experience entries from expert request to expert
func (s *SQLiteStore) copyExperienceEntries(tx *sql.Tx, requestID, expertID int64) error {
	query := `
		INSERT INTO expert_experience_entries (
			expert_id, organization, position, start_date, end_date, 
			is_current, country, description, created_at, updated_at
		)
		SELECT ?, organization, position, start_date, end_date, 
			   is_current, country, description, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		FROM expert_request_experience_entries 
		WHERE expert_request_id = ?`
	
	_, err := tx.Exec(query, expertID, requestID)
	if err != nil {
		return fmt.Errorf("failed to copy experience entries from request %d to expert %d: %w", requestID, expertID, err)
	}
	
	return nil
}

// copyEducationEntries copies education entries from expert request to expert
func (s *SQLiteStore) copyEducationEntries(tx *sql.Tx, requestID, expertID int64) error {
	query := `
		INSERT INTO expert_education_entries (
			expert_id, institution, degree, field_of_study, graduation_year,
			country, description, created_at, updated_at
		)
		SELECT ?, institution, degree, field_of_study, graduation_year,
			   country, description, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		FROM expert_request_education_entries 
		WHERE expert_request_id = ?`
	
	_, err := tx.Exec(query, expertID, requestID)
	if err != nil {
		return fmt.Errorf("failed to copy education entries from request %d to expert %d: %w", requestID, expertID, err)
	}
	
	return nil
}

// UpdateExpertRequestCVDocument updates the CV document reference for an expert request
func (s *SQLiteStore) UpdateExpertRequestCVDocument(requestID, documentID int64) error {
	return s.updateDocumentReference("expert_requests", "cv_document_id", requestID, documentID)
}

// UpdateExpertRequestApprovalDocument updates the approval document reference for an expert request
func (s *SQLiteStore) UpdateExpertRequestApprovalDocument(requestID, documentID int64) error {
	return s.updateDocumentReference("expert_requests", "approval_document_id", requestID, documentID)
}

