package main

import (
	"fmt"
	"strings"
	"time"
	"database/sql"
)

// ListExpertRequests retrieves expert requests based on filters with pagination
func (s *SQLiteStore) ListExpertRequests(filters map[string]interface{}, limit, offset int) ([]*ExpertRequest, error) {
	if limit <= 0 {
		limit = 10
	}

	// Build a dynamic query based on available columns
	columns, err := s.getExpertRequestColumns()
	if err != nil {
		return nil, fmt.Errorf("failed to get expert request columns: %w", err)
	}
	
	// Join columns for SELECT statement
	query := fmt.Sprintf(`
		SELECT %s
		FROM expert_requests
	`, strings.Join(columns, ", "))

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

	// Check columns to handle schema discrepancies
	var columnNames []string
	columnNames, err = rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	
	// Create column map for dynamic scanning
	colMap := make(map[string]int)
	for i, col := range columnNames {
		colMap[col] = i
	}
	
	// Parse the results
	var requests []*ExpertRequest
	for rows.Next() {
		var req ExpertRequest
		var expertID, reviewedAt, reviewedBy interface{}
		
		// Create values slice with the right number of interface{} values
		values := make([]interface{}, len(columnNames))
		for i := range values {
			values[i] = new(interface{})
		}
		
		if err := rows.Scan(values...); err != nil {
			return nil, fmt.Errorf("failed to scan expert request row: %w", err)
		}
		
		// Map column values to struct fields
		for col, idx := range colMap {
			v := *(values[idx].(*interface{}))
			if v == nil {
				continue
			}
			
			switch col {
			case "id":
				if id, ok := v.(int64); ok {
					req.ID = id
				}
			case "expert_id":
				expertID = v
			case "name":
				if s, ok := v.(string); ok {
					req.Name = s
				}
			case "designation":
				if s, ok := v.(string); ok {
					req.Designation = s
				}
			case "institution":
				if s, ok := v.(string); ok {
					req.Institution = s
				}
			case "is_bahraini":
				if b, ok := v.(bool); ok {
					req.IsBahraini = b
				} else if i, ok := v.(int64); ok {
					req.IsBahraini = i > 0
				}
			case "is_available":
				if b, ok := v.(bool); ok {
					req.IsAvailable = b
				} else if i, ok := v.(int64); ok {
					req.IsAvailable = i > 0
				}
			case "rating":
				if s, ok := v.(string); ok {
					req.Rating = s
				}
			case "role":
				if s, ok := v.(string); ok {
					req.Role = s
				}
			case "employment_type":
				if s, ok := v.(string); ok {
					req.EmploymentType = s
				}
			case "general_area":
				if n, ok := v.(int64); ok {
					req.GeneralArea = n
				} else if s, ok := v.(string); ok {
					// Try to convert string to int64
					var areaID int64
					if _, err := fmt.Sscanf(s, "%d", &areaID); err == nil {
						req.GeneralArea = areaID
					}
				}
			case "specialized_area":
				if s, ok := v.(string); ok {
					req.SpecializedArea = s
				}
			case "is_trained":
				if b, ok := v.(bool); ok {
					req.IsTrained = b
				} else if i, ok := v.(int64); ok {
					req.IsTrained = i > 0
				}
			case "cv_path":
				if s, ok := v.(string); ok {
					req.CVPath = s
				}
			case "phone":
				if s, ok := v.(string); ok {
					req.Phone = s
				}
			case "email":
				if s, ok := v.(string); ok {
					req.Email = s
				}
			case "is_published":
				if b, ok := v.(bool); ok {
					req.IsPublished = b
				} else if i, ok := v.(int64); ok {
					req.IsPublished = i > 0
				}
			case "biography":
				if s, ok := v.(string); ok {
					req.Biography = s
				}
			case "status":
				if s, ok := v.(string); ok {
					req.Status = s
				}
			case "rejection_reason":
				if s, ok := v.(string); ok {
					req.RejectionReason = s
				}
			case "created_at":
				if s, ok := v.(string); ok && s != "" {
					parsedTime, err := time.Parse(time.RFC3339, s)
					if err == nil {
						req.CreatedAt = parsedTime
					} else {
						// Try alternate formats if RFC3339 fails
						parsedTime, err = time.Parse("2006-01-02 15:04:05", s)
						if err == nil {
							req.CreatedAt = parsedTime
						}
					}
				}
			case "reviewed_at":
				reviewedAt = v
			case "reviewed_by":
				reviewedBy = v
			}
		}
		
		// Handle nullable fields
		if expertID != nil {
			req.ExpertID = fmt.Sprintf("%v", expertID)
		}
		
		if reviewedAt != nil {
			if timeStr, ok := reviewedAt.(string); ok && timeStr != "" {
				parsedTime, err := time.Parse(time.RFC3339, timeStr)
				if err == nil {
					req.ReviewedAt = parsedTime
				} else {
					// Try alternate formats if RFC3339 fails
					parsedTime, err = time.Parse("2006-01-02 15:04:05", timeStr)
					if err == nil {
						req.ReviewedAt = parsedTime
					}
				}
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
	// First, get the list of columns from the table
	columnQuery := `PRAGMA table_info(expert_requests)`
	rows, err := s.db.Query(columnQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get table info: %w", err)
	}
	
	// Get column names
	var columns []string
	for rows.Next() {
		var cid, notnull, pk int
		var name, dataType string
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &dataType, &notnull, &dfltValue, &pk); err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}
		columns = append(columns, name)
	}
	rows.Close()
	
	// Build a dynamic query based on available columns
	selectCols := []string{"id"}
	columnMap := make(map[string]bool)
	
	for _, col := range columns {
		columnMap[col] = true
		selectCols = append(selectCols, col)
	}
	
	query := fmt.Sprintf(`SELECT %s FROM expert_requests WHERE id = ?`, strings.Join(selectCols, ", "))
	
	// Execute the query
	row := s.db.QueryRow(query, id)
	
	// Create a slice to hold the values
	values := make([]interface{}, len(selectCols))
	for i := range values {
		values[i] = new(interface{})
	}
	
	if err := row.Scan(values...); err != nil {
		return nil, fmt.Errorf("failed to scan expert request: %w", err)
	}
	
	// Create the request object with default values
	var req ExpertRequest
	req.ID = id // We know this is valid
	var expertID, reviewedAt, reviewedBy interface{}
	
	// Map values to struct fields based on column names
	for i, col := range selectCols {
		v := *(values[i].(*interface{}))
		if v == nil {
			continue
		}
		
		switch col {
		case "expert_id":
			expertID = v
		case "name":
			if s, ok := v.(string); ok {
				req.Name = s
			}
		case "designation":
			if s, ok := v.(string); ok {
				req.Designation = s
			}
		case "institution":
			if s, ok := v.(string); ok {
				req.Institution = s
			}
		case "is_bahraini":
			if b, ok := v.(bool); ok {
				req.IsBahraini = b
			} else if i, ok := v.(int64); ok {
				req.IsBahraini = i > 0
			}
		case "is_available":
			if b, ok := v.(bool); ok {
				req.IsAvailable = b
			} else if i, ok := v.(int64); ok {
				req.IsAvailable = i > 0
			}
		case "rating":
			if s, ok := v.(string); ok {
				req.Rating = s
			}
		case "role":
			if s, ok := v.(string); ok {
				req.Role = s
			}
		case "employment_type":
			if s, ok := v.(string); ok {
				req.EmploymentType = s
			}
		case "general_area":
			if n, ok := v.(int64); ok {
				req.GeneralArea = n
			} else if s, ok := v.(string); ok {
				// Try to convert string to int64
				var areaID int64
				if _, err := fmt.Sscanf(s, "%d", &areaID); err == nil {
					req.GeneralArea = areaID
				}
			}
		case "specialized_area":
			if s, ok := v.(string); ok {
				req.SpecializedArea = s
			}
		case "is_trained":
			if b, ok := v.(bool); ok {
				req.IsTrained = b
			} else if i, ok := v.(int64); ok {
				req.IsTrained = i > 0
			}
		case "cv_path":
			if s, ok := v.(string); ok {
				req.CVPath = s
			}
		case "phone":
			if s, ok := v.(string); ok {
				req.Phone = s
			}
		case "email":
			if s, ok := v.(string); ok {
				req.Email = s
			}
		case "is_published":
			if b, ok := v.(bool); ok {
				req.IsPublished = b
			} else if i, ok := v.(int64); ok {
				req.IsPublished = i > 0
			}
		case "biography":
			if s, ok := v.(string); ok {
				req.Biography = s
			}
		case "status":
			if s, ok := v.(string); ok {
				req.Status = s
			}
		case "rejection_reason":
			if s, ok := v.(string); ok {
				req.RejectionReason = s
			}
		case "created_at":
			if s, ok := v.(string); ok && s != "" {
				parsedTime, err := time.Parse(time.RFC3339, s)
				if err == nil {
					req.CreatedAt = parsedTime
				} else {
					// Try alternate formats if RFC3339 fails
					parsedTime, err = time.Parse("2006-01-02 15:04:05", s)
					if err == nil {
						req.CreatedAt = parsedTime
					}
				}
			}
		case "reviewed_at":
			reviewedAt = v
		case "reviewed_by":
			reviewedBy = v
		}
	}
	
	// Handle nullable fields
	if expertID != nil {
		req.ExpertID = fmt.Sprintf("%v", expertID)
	}
	
	if reviewedAt != nil {
		if timeStr, ok := reviewedAt.(string); ok && timeStr != "" {
			parsedTime, err := time.Parse(time.RFC3339, timeStr)
			if err == nil {
				req.ReviewedAt = parsedTime
			} else {
				// Try alternate formats if RFC3339 fails
				parsedTime, err = time.Parse("2006-01-02 15:04:05", timeStr)
				if err == nil {
					req.ReviewedAt = parsedTime
				}
			}
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
	// First, get the list of columns from the table
	columnQuery := `PRAGMA table_info(expert_requests)`
	rows, err := s.db.Query(columnQuery)
	if err != nil {
		return fmt.Errorf("failed to get table info: %w", err)
	}
	
	// Get column names and create a map
	columnMap := make(map[string]bool)
	for rows.Next() {
		var cid, notnull, pk int
		var name, dataType string
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &dataType, &notnull, &dfltValue, &pk); err != nil {
			rows.Close()
			return fmt.Errorf("failed to scan column info: %w", err)
		}
		columnMap[name] = true
	}
	rows.Close()
	
	// Build dynamic update query based on available columns
	var setClauses []string
	var args []interface{}
	
	// Add fields only if the column exists in the table
	if columnMap["expert_id"] {
		setClauses = append(setClauses, "expert_id = ?")
		args = append(args, request.ExpertID)
	}
	
	if columnMap["name"] {
		setClauses = append(setClauses, "name = ?")
		args = append(args, request.Name)
	}
	
	if columnMap["designation"] {
		setClauses = append(setClauses, "designation = ?")
		args = append(args, request.Designation)
	}
	
	if columnMap["institution"] {
		setClauses = append(setClauses, "institution = ?")
		args = append(args, request.Institution)
	}
	
	if columnMap["is_bahraini"] {
		setClauses = append(setClauses, "is_bahraini = ?")
		args = append(args, request.IsBahraini)
	}
	
	if columnMap["is_available"] {
		setClauses = append(setClauses, "is_available = ?")
		args = append(args, request.IsAvailable)
	}
	
	if columnMap["rating"] {
		setClauses = append(setClauses, "rating = ?")
		args = append(args, request.Rating)
	}
	
	if columnMap["role"] {
		setClauses = append(setClauses, "role = ?")
		args = append(args, request.Role)
	}
	
	if columnMap["employment_type"] {
		setClauses = append(setClauses, "employment_type = ?")
		args = append(args, request.EmploymentType)
	}
	
	if columnMap["general_area"] {
		setClauses = append(setClauses, "general_area = ?")
		args = append(args, request.GeneralArea)
	}
	
	if columnMap["specialized_area"] {
		setClauses = append(setClauses, "specialized_area = ?")
		args = append(args, request.SpecializedArea)
	}
	
	if columnMap["is_trained"] {
		setClauses = append(setClauses, "is_trained = ?")
		args = append(args, request.IsTrained)
	}
	
	if columnMap["cv_path"] {
		setClauses = append(setClauses, "cv_path = ?")
		args = append(args, request.CVPath)
	}
	
	if columnMap["phone"] {
		setClauses = append(setClauses, "phone = ?")
		args = append(args, request.Phone)
	}
	
	if columnMap["email"] {
		setClauses = append(setClauses, "email = ?")
		args = append(args, request.Email)
	}
	
	if columnMap["is_published"] {
		setClauses = append(setClauses, "is_published = ?")
		args = append(args, request.IsPublished)
	}
	
	if columnMap["biography"] {
		setClauses = append(setClauses, "biography = ?")
		args = append(args, request.Biography)
	}
	
	if columnMap["status"] {
		setClauses = append(setClauses, "status = ?")
		args = append(args, request.Status)
	}
	
	if columnMap["rejection_reason"] {
		setClauses = append(setClauses, "rejection_reason = ?")
		args = append(args, request.RejectionReason)
	}
	
	// Set reviewed date if not set and status is changing
	if request.Status == "approved" || request.Status == "rejected" {
		if request.ReviewedAt.IsZero() {
			request.ReviewedAt = time.Now()
		}
	}
	
	if columnMap["reviewed_at"] {
		setClauses = append(setClauses, "reviewed_at = ?")
		args = append(args, request.ReviewedAt)
	}
	
	if columnMap["reviewed_by"] {
		setClauses = append(setClauses, "reviewed_by = ?")
		args = append(args, request.ReviewedBy)
	}
	
	// Add ID as the last argument for the WHERE clause
	args = append(args, request.ID)
	
	// Build the final query
	query := fmt.Sprintf("UPDATE expert_requests SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update expert request: %w", err)
	}
	
	return nil
}

// Helper method to get column names from expert_requests table
func (s *SQLiteStore) getExpertRequestColumns() ([]string, error) {
	// Query table schema
	rows, err := s.db.Query("PRAGMA table_info(expert_requests)")
	if err != nil {
		return nil, fmt.Errorf("failed to get table info: %w", err)
	}
	defer rows.Close()
	
	// Extract column names
	var columns []string
	for rows.Next() {
		var cid, notnull, pk int
		var name, dataType string
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &dataType, &notnull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}
		columns = append(columns, name)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating columns: %w", err)
	}
	
	if len(columns) == 0 {
		return nil, fmt.Errorf("no columns found in expert_requests table")
	}
	
	return columns, nil
}