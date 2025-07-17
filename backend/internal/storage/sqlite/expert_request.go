package sqlite

import (
	"fmt"
	"time"
	"database/sql"
	
	"expertdb/internal/domain"
)

// CreateExpertRequest creates a new expert request in the database
func (s *SQLiteStore) CreateExpertRequest(req *domain.ExpertRequest) (int64, error) {
	query := `
		INSERT INTO expert_requests (
			name, designation, institution, is_bahraini, is_available,
			rating, role, employment_type, general_area, specialized_area,
			is_trained, cv_path, approval_document_path, phone, email, is_published, biography,
			status, created_at, created_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	// Set default values if not provided
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}
	if req.Status == "" {
		req.Status = "pending"
	}
	
	// Handle nullable fields or use empty string defaults for non-nullable text fields
	designation := req.Designation
	if designation == "" {
		designation = "" // Not NULL but empty string
	}
	
	affiliation := req.Affiliation
	if affiliation == "" {
		affiliation = "" // Not NULL but empty string
	}
	
	// Rating can be NULL
	var rating interface{} = nil
	if req.Rating != 0 {
		rating = req.Rating
	}
	
	// For specialized area: can be NULL
	var specializedArea interface{} = nil
	if req.SpecializedArea != "" {
		specializedArea = req.SpecializedArea
	}
	
	// CV path can be NULL
	var cvPath interface{} = nil
	if req.CVPath != "" {
		cvPath = req.CVPath
	}
	
	// Approval document path can be NULL
	var approvalDocPath interface{} = nil
	if req.ApprovalDocumentPath != "" {
		approvalDocPath = req.ApprovalDocumentPath
	}
	
	// Biography can be NULL
	var biography interface{} = nil
	if req.Biography != "" {
		biography = req.Biography
	}
	
	result, err := s.db.Exec(
		query,
		req.Name, designation, affiliation,
		req.IsBahraini, req.IsAvailable, rating,
		req.Role, req.EmploymentType, req.GeneralArea,
		specializedArea, req.IsTrained, cvPath,
		approvalDocPath, req.Phone, req.Email, req.IsPublished, biography,
		req.Status, req.CreatedAt, req.CreatedBy,
	)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create expert request: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get expert request ID: %w", err)
	}
	
	// Set the request ID
	req.ID = id
	
	return id, nil
}

// GetExpertRequest retrieves an expert request by ID
func (s *SQLiteStore) GetExpertRequest(id int64) (*domain.ExpertRequest, error) {
	query := `
		SELECT 
			id, expert_id, name, designation, institution, is_bahraini, 
			is_available, rating, role, employment_type, general_area, 
			specialized_area, is_trained, cv_path, approval_document_path, phone, email, 
			is_published, biography, status, rejection_reason, 
			created_at, reviewed_at, reviewed_by, created_by
		FROM expert_requests
		WHERE id = ?
	`
	
	var req domain.ExpertRequest
	var expertID sql.NullString
	var reviewedAt sql.NullTime
	var reviewedBy sql.NullInt64
	var createdBy sql.NullInt64
	
	err := s.db.QueryRow(query, id).Scan(
		&req.ID, &expertID, &req.Name, &req.Designation, &req.Affiliation, 
		&req.IsBahraini, &req.IsAvailable, &req.Rating, &req.Role, 
		&req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
		&req.IsTrained, &req.CVPath, &req.ApprovalDocumentPath, &req.Phone, &req.Email, 
		&req.IsPublished, &req.Biography, &req.Status, &req.RejectionReason, 
		&req.CreatedAt, &reviewedAt, &reviewedBy, &createdBy,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get expert request: %w", err)
	}
	
	// Set nullable fields
	if expertID.Valid {
		req.ExpertID = expertID.String
	}
	
	if reviewedAt.Valid {
		req.ReviewedAt = reviewedAt.Time
	}
	
	if reviewedBy.Valid {
		req.ReviewedBy = reviewedBy.Int64
	}
	
	if createdBy.Valid {
		req.CreatedBy = createdBy.Int64
	}
	
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
				id, expert_id, name, designation, institution, is_bahraini, 
				is_available, rating, role, employment_type, general_area, 
				specialized_area, is_trained, cv_path, approval_document_path, phone, email, 
				is_published, biography, status, rejection_reason, 
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
				id, expert_id, name, designation, institution, is_bahraini, 
				is_available, rating, role, employment_type, general_area, 
				specialized_area, is_trained, cv_path, approval_document_path, phone, email, 
				is_published, biography, status, rejection_reason, 
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
		var expertID sql.NullString
		var reviewedAt sql.NullTime
		var reviewedBy sql.NullInt64
		var createdBy sql.NullInt64
		
		err := rows.Scan(
			&req.ID, &expertID, &req.Name, &req.Designation, &req.Affiliation, 
			&req.IsBahraini, &req.IsAvailable, &req.Rating, &req.Role, 
			&req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
			&req.IsTrained, &req.CVPath, &req.ApprovalDocumentPath, &req.Phone, &req.Email, 
			&req.IsPublished, &req.Biography, &req.Status, &req.RejectionReason, 
			&req.CreatedAt, &reviewedAt, &reviewedBy, &createdBy,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan expert request row: %w", err)
		}
		
		// Set nullable fields
		if expertID.Valid {
			req.ExpertID = expertID.String
		}
		
		if reviewedAt.Valid {
			req.ReviewedAt = reviewedAt.Time
		}
		
		if reviewedBy.Valid {
			req.ReviewedBy = reviewedBy.Int64
		}
		
		if createdBy.Valid {
			req.CreatedBy = createdBy.Int64
		}
		
		requests = append(requests, &req)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expert request rows: %w", err)
	}
	
	return requests, nil
}

// UpdateExpertRequestStatus updates the status of an expert request
func (s *SQLiteStore) UpdateExpertRequestStatus(id int64, status, rejectionReason string, reviewedBy int64) error {
	query := `
		UPDATE expert_requests
		SET status = ?, rejection_reason = ?, reviewed_at = ?, reviewed_by = ?
		WHERE id = ?
	`
	
	now := time.Now().UTC()
	
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
	
	// If approved, automatically create expert
	if status == "approved" {
		// Retrieve the request
		req, err := s.GetExpertRequest(id)
		if err != nil {
			return fmt.Errorf("failed to retrieve request for expert creation: %w", err)
		}
		
		// Create expert
		expert := &domain.Expert{
			// Don't set ExpertID - let the CreateExpert function generate it using GenerateUniqueExpertID
			Name:                req.Name,
			Email:               req.Email,
			Phone:               req.Phone,
			Biography:           req.Biography,
			CVPath:              req.CVPath,
			ApprovalDocumentPath: req.ApprovalDocumentPath,
			Designation:         req.Designation,
			Affiliation:         req.Affiliation,
			IsBahraini:          req.IsBahraini,
			IsAvailable:         req.IsAvailable,
			Rating:              req.Rating,
			Role:                req.Role,
			EmploymentType:      req.EmploymentType,
			GeneralArea:         req.GeneralArea,
			SpecializedArea:     req.SpecializedArea,
			IsPublished:         req.IsPublished,
			IsTrained:           req.IsTrained,
			CreatedAt:           now,
			UpdatedAt:           now,
			OriginalRequestID:   id, // Set the reference to the original request
		}
		
		// Create the expert using the sequential ID generator
		expertID, err := s.CreateExpert(expert)
		if err != nil {
			return fmt.Errorf("failed to create expert on approval: %w", err)
		}
		
		// Update the request with the generated expert ID
		createdExpert, err := s.GetExpert(expertID)
		if err != nil {
			return fmt.Errorf("failed to retrieve created expert: %w", err)
		}
		
		// Set the expert_id in the request record
		_, err = s.db.Exec("UPDATE expert_requests SET expert_id = ? WHERE id = ?", createdExpert.ExpertID, id)
		if err != nil {
			return fmt.Errorf("failed to update expert request with expert_id: %w", err)
		}
	}
	
	return nil
}

// UpdateExpertRequest updates an expert request with new data
func (s *SQLiteStore) UpdateExpertRequest(req *domain.ExpertRequest) error {
	query := `
		UPDATE expert_requests
		SET name = ?, designation = ?, institution = ?, is_bahraini = ?,
			is_available = ?, rating = ?, role = ?, employment_type = ?,
			general_area = ?, specialized_area = ?, is_trained = ?,
			cv_path = ?, approval_document_path = ?, phone = ?, email = ?, is_published = ?,
			biography = ?, status = ?, rejection_reason = ?,
			expert_id = ?, reviewed_at = ?, reviewed_by = ?, created_by = ?
		WHERE id = ?
	`
	
	// Handle nullable fields
	var rating, specializedArea, cvPath, approvalDocPath, biography, rejectionReason, expertID interface{} = nil, nil, nil, nil, nil, nil, nil
	
	if req.Rating != 0 {
		rating = req.Rating
	}
	if req.SpecializedArea != "" {
		specializedArea = req.SpecializedArea
	}
	if req.CVPath != "" {
		cvPath = req.CVPath
	}
	if req.ApprovalDocumentPath != "" {
		approvalDocPath = req.ApprovalDocumentPath
	}
	if req.Biography != "" {
		biography = req.Biography
	}
	if req.RejectionReason != "" {
		rejectionReason = req.RejectionReason
	}
	if req.ExpertID != "" {
		expertID = req.ExpertID
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
	
	// Execute update
	result, err := s.db.Exec(
		query,
		req.Name, req.Designation, req.Affiliation, req.IsBahraini,
		req.IsAvailable, rating, req.Role, req.EmploymentType,
		req.GeneralArea, specializedArea, req.IsTrained,
		cvPath, approvalDocPath, req.Phone, req.Email, req.IsPublished,
		biography, req.Status, rejectionReason,
		expertID, reviewedAt, reviewedBy, createdBy,
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

// BatchApproveExpertRequests approves multiple expert requests in a single transaction
// Returns a list of successfully approved request IDs and a map of errors for failed approvals
func (s *SQLiteStore) BatchApproveExpertRequests(requestIDs []int64, approvalDocumentPath string, reviewedBy int64) ([]int64, map[int64]error) {
	successIDs := []int64{}
	errors := make(map[int64]error)
	
	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		// If we can't even start a transaction, return error for all IDs
		for _, id := range requestIDs {
			errors[id] = fmt.Errorf("failed to begin transaction: %w", err)
		}
		return successIDs, errors
	}
	
	// Defer rollback in case of error - this will be a no-op if we commit
	defer tx.Rollback()
	
	// Prepare update statement
	now := time.Now().UTC()
	updateStmt, err := tx.Prepare(`
		UPDATE expert_requests
		SET status = ?, rejection_reason = ?, reviewed_at = ?, reviewed_by = ?, approval_document_path = ?
		WHERE id = ? AND status = 'pending'
	`)
	if err != nil {
		for _, id := range requestIDs {
			errors[id] = fmt.Errorf("failed to prepare statement: %w", err)
		}
		return successIDs, errors
	}
	defer updateStmt.Close()
	
	// Process each request
	for _, id := range requestIDs {
		// Update the status
		result, err := updateStmt.Exec("approved", "", now, reviewedBy, approvalDocumentPath, id)
		if err != nil {
			errors[id] = fmt.Errorf("failed to update request status: %w", err)
			continue
		}
		
		// Check if a row was affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			errors[id] = fmt.Errorf("failed to get rows affected: %w", err)
			continue
		}
		
		if rowsAffected == 0 {
			// This could be because the request doesn't exist or was not in pending status
			// Get the request to find out
			var status string
			err := tx.QueryRow("SELECT status FROM expert_requests WHERE id = ?", id).Scan(&status)
			
			if err != nil {
				if err == sql.ErrNoRows {
					errors[id] = domain.ErrNotFound
				} else {
					errors[id] = fmt.Errorf("failed to check request status: %w", err)
				}
				continue
			}
			
			if status != "pending" {
				errors[id] = fmt.Errorf("request is not in pending status (current: %s)", status)
				continue
			}
			
			errors[id] = fmt.Errorf("request not updated for unknown reason")
			continue
		}
		
		// Get the request data for expert creation
		var req domain.ExpertRequest
		query := `
			SELECT 
				id, expert_id, name, designation, institution, is_bahraini, 
				is_available, rating, role, employment_type, general_area, 
				specialized_area, is_trained, cv_path, phone, email, 
				is_published, biography, status, created_by
			FROM expert_requests
			WHERE id = ?
		`
		
		err = tx.QueryRow(query, id).Scan(
			&req.ID, &req.ExpertID, &req.Name, &req.Designation, &req.Affiliation, 
			&req.IsBahraini, &req.IsAvailable, &req.Rating, &req.Role, 
			&req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
			&req.IsTrained, &req.CVPath, &req.Phone, &req.Email, 
			&req.IsPublished, &req.Biography, &req.Status, &req.CreatedBy,
		)
		
		if err != nil {
			errors[id] = fmt.Errorf("failed to retrieve request data: %w", err)
			continue
		}
		
		// Create expert
		expert := &domain.Expert{
			// Don't set ExpertID - will be generated by CreateExpert
			Name:                req.Name,
			Email:               req.Email,
			Phone:               req.Phone,
			Biography:           req.Biography,
			CVPath:              req.CVPath,
			ApprovalDocumentPath: approvalDocumentPath,
			Designation:         req.Designation,
			Affiliation:         req.Affiliation,
			IsBahraini:          req.IsBahraini,
			IsAvailable:         req.IsAvailable,
			Rating:              req.Rating,
			Role:                req.Role,
			EmploymentType:      req.EmploymentType,
			GeneralArea:         req.GeneralArea,
			SpecializedArea:     req.SpecializedArea,
			IsPublished:         req.IsPublished,
			IsTrained:           req.IsTrained,
			CreatedAt:           now,
			UpdatedAt:           now,
			OriginalRequestID:   id,
		}
		
		// Generate unique expert ID
		expertIDStmt, err := tx.Prepare(`
			SELECT COALESCE(MAX(CAST(SUBSTR(expert_id, 5) AS INTEGER)), 0) + 1
			FROM experts
			WHERE expert_id LIKE 'EXP-%'
		`)
		if err != nil {
			errors[id] = fmt.Errorf("failed to prepare expert ID statement: %w", err)
			continue
		}
		
		var nextID int
		err = expertIDStmt.QueryRow().Scan(&nextID)
		expertIDStmt.Close()
		
		if err != nil {
			errors[id] = fmt.Errorf("failed to generate expert ID: %w", err)
			continue
		}
		
		expertID := fmt.Sprintf("EXP-%04d", nextID)
		expert.ExpertID = expertID
		
		// Insert the expert
		expertStmt, err := tx.Prepare(`
			INSERT INTO experts (
				expert_id, name, designation, institution, is_bahraini, is_available,
				rating, role, employment_type, general_area, specialized_area,
				is_trained, cv_path, approval_document_path, phone, email, is_published, biography,
				created_at, updated_at, original_request_id
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			errors[id] = fmt.Errorf("failed to prepare expert insert statement: %w", err)
			continue
		}
		
		result, err = expertStmt.Exec(
			expert.ExpertID, expert.Name, expert.Designation, expert.Affiliation,
			expert.IsBahraini, expert.IsAvailable, expert.Rating, expert.Role,
			expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
			expert.IsTrained, expert.CVPath, expert.ApprovalDocumentPath,
			expert.Phone, expert.Email, expert.IsPublished, expert.Biography,
			expert.CreatedAt, expert.UpdatedAt, expert.OriginalRequestID,
		)
		expertStmt.Close()
		
		if err != nil {
			errors[id] = fmt.Errorf("failed to insert expert: %w", err)
			continue
		}
		
		// Update the request with the expert ID
		_, err = tx.Exec("UPDATE expert_requests SET expert_id = ? WHERE id = ?", expertID, id)
		if err != nil {
			errors[id] = fmt.Errorf("failed to update request with expert ID: %w", err)
			continue
		}
		
		// This request was successful
		successIDs = append(successIDs, id)
	}
	
	// If we have at least one success, commit the transaction
	if len(successIDs) > 0 {
		if err := tx.Commit(); err != nil {
			// If commit fails, all operations fail
			for _, id := range successIDs {
				errors[id] = fmt.Errorf("failed to commit transaction: %w", err)
			}
			return []int64{}, errors
		}
	} else {
		// No successful operations, so return all errors
		return []int64{}, errors
	}
	
	return successIDs, errors
}