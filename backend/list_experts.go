package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// ListExperts retrieves experts based on filters with pagination and sorting
func (s *SQLiteStore) ListExperts(filters map[string]interface{}, limit, offset int) ([]*Expert, error) {
	// Default limit if not specified
	if limit <= 0 {
		limit = 10
	}

	// Build the query to select all expert fields with optional filtering and sorting
	// This query retrieves complete expert records from the experts table
	// with support for filtering by various criteria and sorting options
	query := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, e.is_bahraini, 
		       e.is_available, e.rating, e.role, e.employment_type, e.general_area, 
		       e.specialized_area, e.is_trained, e.cv_path, e.phone, e.email, e.is_published,
		       e.created_at, e.updated_at, e.biography
		FROM experts e
	`

	// Apply filters
	var conditions []string
	var args []interface{}
	var sortBy string = "e.name"
	var sortOrder string = "ASC"
	
	if len(filters) > 0 {
		for key, value := range filters {
			switch key {
			case "name":
				conditions = append(conditions, "e.name LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			case "area":
				conditions = append(conditions, "e.general_area = ?")
				// Try to convert to integer if it's a string
				if strVal, ok := value.(string); ok {
					if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
						args = append(args, intVal)
					} else {
						args = append(args, value)
					}
				} else {
					args = append(args, value)
				}
			case "is_available":
				conditions = append(conditions, "e.is_available = ?")
				args = append(args, value)
			case "role":
				conditions = append(conditions, "e.role LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			case "min_rating":
				// Handle minimum rating filter (convert to numeric for comparison)
				conditions = append(conditions, "CAST(e.rating AS REAL) >= ?")
				args = append(args, value)
			case "nationality":
				if value == "bahraini" {
					conditions = append(conditions, "e.is_bahraini = 1")
				} else if value == "international" {
					conditions = append(conditions, "e.is_bahraini = 0")
				}
			case "sort":
				// Handle sorting column
				switch value {
				case "name":
					sortBy = "e.name"
				case "institution":
					sortBy = "e.institution"
				case "role":
					sortBy = "e.role"
				case "created_at":
					sortBy = "e.created_at"
				case "rating":
					sortBy = "e.rating"
				case "general_area":
					sortBy = "e.general_area"
				default:
					sortBy = "e.name"
				}
			case "sort_order":
				// Handle sort order
				orderVal, isString := value.(string)
				if isString {
					if strings.ToLower(orderVal) == "desc" {
						sortOrder = "DESC"
					}
				}
			}
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add sorting and pagination
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT ? OFFSET ?", sortBy, sortOrder)
	args = append(args, limit, offset)

	// Execute the query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query experts: %w", err)
	}
	defer rows.Close()

	var experts []*Expert
	for rows.Next() {
		var expert Expert
		var createdAt string
		var nullableExpertID, nullableCVPath, nullablePhone, nullableUpdatedAt, nullableSpecializedArea, nullableBiography sql.NullString
		var nullableIsAvailable sql.NullBool

		err := rows.Scan(
			&expert.ID, &nullableExpertID, &expert.Name, &expert.Designation, &expert.Institution,
			&expert.IsBahraini, &nullableIsAvailable, &expert.Rating, &expert.Role, &expert.EmploymentType,
			&expert.GeneralArea, &nullableSpecializedArea, &expert.IsTrained, &nullableCVPath,
			&nullablePhone, &expert.Email, &expert.IsPublished,
			&createdAt, &nullableUpdatedAt, &nullableBiography,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expert row: %w", err)
		}

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

		// Convert time strings to Time objects
		if createdAt != "" {
			expert.CreatedAt, err = parseTime(createdAt)
			if err != nil {
				// Don't fail the whole query if one timestamp is invalid
				// Just default to current time
				expert.CreatedAt = time.Now()
			}
		}

		if nullableUpdatedAt.Valid && nullableUpdatedAt.String != "" {
			expert.UpdatedAt, err = parseTime(nullableUpdatedAt.String)
			if err != nil {
				// Don't fail the whole query if one timestamp is invalid
				expert.UpdatedAt = time.Time{}
			}
		}

		// Load the expert area name if possible
		if expert.GeneralArea > 0 {
			area, err := s.GetExpertAreaByID(expert.GeneralArea)
			if err == nil && area != nil {
				expert.GeneralAreaName = area.Name
			}
		}

		experts = append(experts, &expert)
	}

	return experts, nil
}

// CountExperts counts the total number of experts that match the given filters
func (s *SQLiteStore) CountExperts(filters map[string]interface{}) (int, error) {
	query := "SELECT COUNT(*) FROM experts e"

	// Apply filters
	var conditions []string
	var args []interface{}
	
	if len(filters) > 0 {
		for key, value := range filters {
			switch key {
			case "name":
				conditions = append(conditions, "e.name LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			case "area":
				conditions = append(conditions, "e.general_area LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			case "is_available":
				conditions = append(conditions, "e.is_available = ?")
				args = append(args, value)
			case "role":
				conditions = append(conditions, "e.role LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			case "min_rating":
				conditions = append(conditions, "CAST(e.rating AS REAL) >= ?")
				args = append(args, value)
			case "nationality":
				if value == "bahraini" {
					conditions = append(conditions, "e.is_bahraini = 1")
				} else if value == "international" {
					conditions = append(conditions, "e.is_bahraini = 0")
				}
			}
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Execute count query
	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count experts: %w", err)
	}

	return count, nil
}

// Helper function to parse time strings
func parseTime(timeStr string) (time.Time, error) {
	// Try RFC3339 format first (standard for API)
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t, nil
	}
	
	// Try RFC3339Nano format
	if t, err := time.Parse(time.RFC3339Nano, timeStr); err == nil {
		return t, nil
	}
	
	// Try common SQLite ISO format
	if t, err := time.Parse("2006-01-02 15:04:05", timeStr); err == nil {
		return t, nil
	}
	
	// Try date-only formats
	if t, err := time.Parse("2006-01-02", timeStr); err == nil {
		return t, nil
	}
	
	// If all else fails, try to parse original format
	return time.Parse(time.RFC3339, timeStr)
}