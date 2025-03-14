package main

import (
	"fmt"
	"time"
)

// CreateExpertRequest creates a new expert request in the database
func (s *SQLiteStore) CreateExpertRequest(request *ExpertRequest) (int64, error) {
	query := `
		INSERT INTO expert_requests (
			name, designation, institution, is_bahraini, is_available,
			rating, role, employment_type, general_area, specialized_area,
			is_trained, cv_path, phone, email, is_published, biography,
			status, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	// Set default values if not provided
	if request.CreatedAt.IsZero() {
		request.CreatedAt = time.Now()
	}
	if request.Status == "" {
		request.Status = "pending"
	}
	
	result, err := s.db.Exec(
		query,
		request.Name, request.Designation, request.Institution,
		request.IsBahraini, request.IsAvailable, request.Rating,
		request.Role, request.EmploymentType, request.GeneralArea,
		request.SpecializedArea, request.IsTrained, request.CVPath,
		request.Phone, request.Email, request.IsPublished, request.Biography,
		request.Status, request.CreatedAt,
	)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create expert request: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get expert request ID: %w", err)
	}
	
	// Set the request ID
	request.ID = id
	
	return id, nil
}