package main

import (
	"fmt"
	"strings"
	"time"
)

// ListExpertRequests retrieves expert requests based on filters with pagination
func (s *SQLiteStore) ListExpertRequests(filters map[string]interface{}, limit, offset int) ([]*ExpertRequest, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, expert_id, name, designation, institution, is_bahraini, 
		       is_available, rating, role, employment_type, general_area, 
		       specialized_area, is_trained, cv_path, phone, email, is_published,
		       status, created_at, reviewed_at, reviewed_by
		FROM expert_requests
	`

	// Apply filters
	var conditions []string
	var args []interface{}
	
	if len(filters) > 0 {
		for key, value := range filters {
			switch key {
			case "status":
				conditions = append(conditions, "status = ?")
				args = append(args, value)
			case "name":
				conditions = append(conditions, "name LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			case "institution":
				conditions = append(conditions, "institution LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			}
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add sorting and pagination
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute the query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query expert requests: %w", err)
	}
	defer rows.Close()

	// Parse the results
	var requests []*ExpertRequest
	for rows.Next() {
		var req ExpertRequest
		var expertID, reviewedAt, reviewedBy interface{}
		
		err := rows.Scan(
			&req.ID, &expertID, &req.Name, &req.Designation, &req.Institution,
			&req.IsBahraini, &req.IsAvailable, &req.Rating, &req.Role, &req.EmploymentType,
			&req.GeneralArea, &req.SpecializedArea, &req.IsTrained, &req.CVPath,
			&req.Phone, &req.Email, &req.IsPublished, &req.Status, &req.CreatedAt,
			&reviewedAt, &reviewedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expert request row: %w", err)
		}
		
		// Handle nullable fields
		if expertID != nil {
			req.ExpertID = fmt.Sprintf("%v", expertID)
		}
		
		if reviewedAt != nil {
			if timeStr, ok := reviewedAt.(string); ok {
				req.ReviewedAt, _ = time.Parse(time.RFC3339, timeStr)
			}
		}
		
		if reviewedBy != nil {
			if idNum, ok := reviewedBy.(int64); ok {
				req.ReviewedBy = idNum
			}
		}
		
		requests = append(requests, &req)
	}

	return requests, nil
}

// GetExpertRequest retrieves an expert request by ID
func (s *SQLiteStore) GetExpertRequest(id int64) (*ExpertRequest, error) {
	query := `
		SELECT id, expert_id, name, designation, institution, is_bahraini, 
		       is_available, rating, role, employment_type, general_area, 
		       specialized_area, is_trained, cv_path, phone, email, is_published,
		       status, created_at, reviewed_at, reviewed_by
		FROM expert_requests
		WHERE id = ?
	`
	
	var req ExpertRequest
	var expertID, reviewedAt, reviewedBy interface{}
	
	err := s.db.QueryRow(query, id).Scan(
		&req.ID, &expertID, &req.Name, &req.Designation, &req.Institution,
		&req.IsBahraini, &req.IsAvailable, &req.Rating, &req.Role, &req.EmploymentType,
		&req.GeneralArea, &req.SpecializedArea, &req.IsTrained, &req.CVPath,
		&req.Phone, &req.Email, &req.IsPublished, &req.Status, &req.CreatedAt,
		&reviewedAt, &reviewedBy,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get expert request: %w", err)
	}
	
	// Handle nullable fields
	if expertID != nil {
		req.ExpertID = fmt.Sprintf("%v", expertID)
	}
	
	if reviewedAt != nil {
		if timeStr, ok := reviewedAt.(string); ok {
			req.ReviewedAt, _ = time.Parse(time.RFC3339, timeStr)
		}
	}
	
	if reviewedBy != nil {
		if idNum, ok := reviewedBy.(int64); ok {
			req.ReviewedBy = idNum
		}
	}
	
	return &req, nil
}

// UpdateExpertRequest updates an existing expert request
func (s *SQLiteStore) UpdateExpertRequest(request *ExpertRequest) error {
	query := `
		UPDATE expert_requests
		SET expert_id = ?, name = ?, designation = ?, institution = ?,
		    is_bahraini = ?, is_available = ?, rating = ?, role = ?,
		    employment_type = ?, general_area = ?, specialized_area = ?,
		    is_trained = ?, cv_path = ?, phone = ?, email = ?, is_published = ?,
		    status = ?, reviewed_at = ?, reviewed_by = ?
		WHERE id = ?
	`
	
	// Set reviewed date if not set and status is changing
	if request.Status == "approved" || request.Status == "rejected" {
		if request.ReviewedAt.IsZero() {
			request.ReviewedAt = time.Now()
		}
	}
	
	_, err := s.db.Exec(
		query,
		request.ExpertID, request.Name, request.Designation, request.Institution,
		request.IsBahraini, request.IsAvailable, request.Rating, request.Role,
		request.EmploymentType, request.GeneralArea, request.SpecializedArea,
		request.IsTrained, request.CVPath, request.Phone, request.Email, request.IsPublished,
		request.Status, request.ReviewedAt, request.ReviewedBy, request.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update expert request: %w", err)
	}
	
	return nil
}