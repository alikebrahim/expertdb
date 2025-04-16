package sqlite

import (
	"fmt"
	"time"
	"strings"
	"database/sql"
	
	"expertdb/internal/domain"
)

// CreateExpertRequest creates a new expert request in the database
func (s *SQLiteStore) CreateExpertRequest(req *domain.ExpertRequest) (int64, error) {
	query := `
		INSERT INTO expert_requests (
			name, designation, institution, is_bahraini, is_available,
			rating, role, employment_type, general_area, specialized_area,
			is_trained, cv_path, phone, email, is_published, biography,
			status, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
	
	institution := req.Institution
	if institution == "" {
		institution = "" // Not NULL but empty string
	}
	
	// Rating can be NULL
	var rating interface{} = nil
	if req.Rating != "" {
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
	
	// Biography can be NULL
	var biography interface{} = nil
	if req.Biography != "" {
		biography = req.Biography
	}
	
	result, err := s.db.Exec(
		query,
		req.Name, designation, institution,
		req.IsBahraini, req.IsAvailable, rating,
		req.Role, req.EmploymentType, req.GeneralArea,
		specializedArea, req.IsTrained, cvPath,
		req.Phone, req.Email, req.IsPublished, biography,
		req.Status, req.CreatedAt,
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
			specialized_area, is_trained, cv_path, phone, email, 
			is_published, biography, status, rejection_reason, 
			created_at, reviewed_at, reviewed_by
		FROM expert_requests
		WHERE id = ?
	`
	
	var req domain.ExpertRequest
	var expertID sql.NullString
	var reviewedAt sql.NullTime
	var reviewedBy sql.NullInt64
	
	err := s.db.QueryRow(query, id).Scan(
		&req.ID, &expertID, &req.Name, &req.Designation, &req.Institution, 
		&req.IsBahraini, &req.IsAvailable, &req.Rating, &req.Role, 
		&req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
		&req.IsTrained, &req.CVPath, &req.Phone, &req.Email, 
		&req.IsPublished, &req.Biography, &req.Status, &req.RejectionReason, 
		&req.CreatedAt, &reviewedAt, &reviewedBy,
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
	
	return &req, nil
}

// ListExpertRequests retrieves a list of expert requests with the given status
func (s *SQLiteStore) ListExpertRequests(status string, limit, offset int) ([]*domain.ExpertRequest, error) {
	if limit <= 0 {
		limit = 10
	}
	
	var query string
	var args []interface{}
	
	if status != "" {
		query = `
			SELECT 
				id, expert_id, name, designation, institution, is_bahraini, 
				is_available, rating, role, employment_type, general_area, 
				specialized_area, is_trained, cv_path, phone, email, 
				is_published, biography, status, rejection_reason, 
				created_at, reviewed_at, reviewed_by
			FROM expert_requests
			WHERE status = ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{status, limit, offset}
	} else {
		query = `
			SELECT 
				id, expert_id, name, designation, institution, is_bahraini, 
				is_available, rating, role, employment_type, general_area, 
				specialized_area, is_trained, cv_path, phone, email, 
				is_published, biography, status, rejection_reason, 
				created_at, reviewed_at, reviewed_by
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
		
		err := rows.Scan(
			&req.ID, &expertID, &req.Name, &req.Designation, &req.Institution, 
			&req.IsBahraini, &req.IsAvailable, &req.Rating, &req.Role, 
			&req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
			&req.IsTrained, &req.CVPath, &req.Phone, &req.Email, 
			&req.IsPublished, &req.Biography, &req.Status, &req.RejectionReason, 
			&req.CreatedAt, &reviewedAt, &reviewedBy,
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
	result, err := s.db.Exec(
		query,
		status,
		rejectionReason,
		now,
		reviewedBy,
		id,
	)
	
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
	
	return nil
}

// UpdateExpertRequest updates an expert request with new data
func (s *SQLiteStore) UpdateExpertRequest(req *domain.ExpertRequest) error {
	query := `
		UPDATE expert_requests
		SET name = ?, designation = ?, institution = ?, is_bahraini = ?,
			is_available = ?, rating = ?, role = ?, employment_type = ?,
			general_area = ?, specialized_area = ?, is_trained = ?,
			cv_path = ?, phone = ?, email = ?, is_published = ?,
			biography = ?, status = ?, rejection_reason = ?,
			expert_id = ?, reviewed_at = ?, reviewed_by = ?
		WHERE id = ?
	`
	
	// Handle nullable fields
	var rating, specializedArea, cvPath, biography, rejectionReason, expertID interface{} = nil, nil, nil, nil, nil, nil
	
	if req.Rating != "" {
		rating = req.Rating
	}
	if req.SpecializedArea != "" {
		specializedArea = req.SpecializedArea
	}
	if req.CVPath != "" {
		cvPath = req.CVPath
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
	
	// Execute update
	result, err := s.db.Exec(
		query,
		req.Name, req.Designation, req.Institution, req.IsBahraini,
		req.IsAvailable, rating, req.Role, req.EmploymentType,
		req.GeneralArea, specializedArea, req.IsTrained,
		cvPath, req.Phone, req.Email, req.IsPublished,
		biography, req.Status, rejectionReason,
		expertID, reviewedAt, reviewedBy,
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