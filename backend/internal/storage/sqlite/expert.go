package sqlite

import (
	"database/sql"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"fmt"
	"os"
	"strings"
	"time"
)

// GenerateUniqueExpertID generates a unique sequential ID for an expert
func (s *SQLiteStore) GenerateUniqueExpertID() (string, error) {
	// Use a transaction to ensure atomicity when incrementing the sequence
	tx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get and increment the sequence in a single step
	var nextVal int
	err = tx.QueryRow("UPDATE expert_id_sequence SET next_val = next_val + 1 WHERE id = 1 RETURNING next_val").Scan(&nextVal)
	if err != nil {
		return "", fmt.Errorf("failed to get next expert ID sequence: %w", err)
	}

	// Format the ID with leading zeros to ensure consistent formatting (EXP-0001, EXP-0002, etc.)
	expertID := fmt.Sprintf("EXP-%04d", nextVal)

	// Verify the ID doesn't already exist (safeguard against potential issues)
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM experts WHERE expert_id = ?", expertID).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("failed to check if expert ID exists: %w", err)
	}

	// In the unlikely case of a collision, try the next ID
	if count > 0 {
		// Release current transaction
		tx.Rollback()
		// Recursive call should be rare and only happens in case of a conflict
		return s.GenerateUniqueExpertID()
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return expertID, nil
}

// ExpertIDExists checks if an expert ID already exists in the database
func (s *SQLiteStore) ExpertIDExists(expertID string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM experts WHERE expert_id = ?"

	err := s.db.QueryRow(query, expertID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if expert ID exists: %w", err)
	}

	return count > 0, nil
}

// CreateExpert creates a new expert in the database
func (s *SQLiteStore) CreateExpert(expert *domain.Expert) (int64, error) {
	// Generate a unique expert_id if not provided or empty
	if expert.ExpertID == "" {
		var err error
		expert.ExpertID, err = s.GenerateUniqueExpertID()
		if err != nil {
			return 0, fmt.Errorf("failed to generate unique expert ID: %w", err)
		}
	} else {
		// Check if the provided expert_id already exists
		exists, err := s.ExpertIDExists(expert.ExpertID)
		if err != nil {
			return 0, fmt.Errorf("failed to check if expert ID exists: %w", err)
		}
		if exists {
			return 0, fmt.Errorf("expert ID already exists: %s", expert.ExpertID)
		}
	}

	query := `
		INSERT INTO experts (
			expert_id, name, designation, institution, is_bahraini, is_available, rating,
			role, employment_type, general_area, specialized_area, is_trained,
			cv_path, approval_document_path, phone, email, is_published, biography, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if expert.CreatedAt.IsZero() {
		expert.CreatedAt = time.Now().UTC()
		expert.UpdatedAt = expert.CreatedAt
	}

	result, err := s.db.Exec(
		query,
		expert.ExpertID, expert.Name, expert.Designation, expert.Institution,
		expert.IsBahraini, expert.IsAvailable, expert.Rating,
		expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.ApprovalDocumentPath, expert.Phone, expert.Email, expert.IsPublished,
		expert.Biography, expert.CreatedAt, expert.UpdatedAt,
	)

	if err != nil {
		// Parse SQLite error to provide more specific error messages
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "expert_id") {
				return 0, fmt.Errorf("expert ID already exists: %s (use a different ID or let the system generate one)", expert.ExpertID)
			} else if strings.Contains(err.Error(), "email") {
				return 0, fmt.Errorf("email already exists: %s (an expert with this email is already registered)", expert.Email)
			} else {
				// Identify other unique constraint violations
				constraintName := extractConstraintName(err.Error())
				return 0, fmt.Errorf("unique constraint violation on %s: duplicate value not allowed", constraintName)
			}
		} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			// Identify which foreign key failed
			if strings.Contains(err.Error(), "general_area") {
				return 0, fmt.Errorf("invalid general area ID %d: this area does not exist in the system", expert.GeneralArea)
			} else {
				return 0, fmt.Errorf("referenced resource does not exist: %w", err)
			}
		} else if strings.Contains(err.Error(), "NOT NULL constraint failed") {
			// Extract column name from error message
			colName := extractColumnName(err.Error())
			return 0, fmt.Errorf("required field missing: %s cannot be empty", colName)
		}
		
		// Log the full error for debugging but return a cleaner message to the user
		logger.Get().Error("Database error creating expert: %v", err)
		return 0, fmt.Errorf("failed to create expert: database error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get expert ID: %w", err)
	}

	return id, nil
}

// GetExpert retrieves an expert by their ID
func (s *SQLiteStore) GetExpert(id int64) (*domain.Expert, error) {
	query := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, 
		       e.is_bahraini, e.nationality, e.is_available, e.rating, e.role, 
		       e.employment_type, e.general_area, ea.name as general_area_name, 
		       e.specialized_area, e.is_trained, e.cv_path, e.approval_document_path, e.phone, e.email, 
		       e.is_published, e.biography, e.created_at, e.updated_at
		FROM experts e
		LEFT JOIN expert_areas ea ON e.general_area = ea.id
		WHERE e.id = ?
	`

	var expert domain.Expert
	var generalAreaName sql.NullString
	var nationality sql.NullString

	err := s.db.QueryRow(query, id).Scan(
		&expert.ID, &expert.ExpertID, &expert.Name, &expert.Designation, &expert.Institution,
		&expert.IsBahraini, &nationality, &expert.IsAvailable, &expert.Rating, &expert.Role,
		&expert.EmploymentType, &expert.GeneralArea, &generalAreaName,
		&expert.SpecializedArea, &expert.IsTrained, &expert.CVPath, &expert.ApprovalDocumentPath, &expert.Phone, &expert.Email,
		&expert.IsPublished, &expert.Biography, &expert.CreatedAt, &expert.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get expert: %w", err)
	}

	if generalAreaName.Valid {
		expert.GeneralAreaName = generalAreaName.String
	}

	if nationality.Valid {
		expert.Nationality = nationality.String
	}

	// Fetch documents and engagements
	documents, err := s.ListDocuments(expert.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expert documents: %w", err)
	}
	// Convert []*Document to []Document
	if len(documents) > 0 {
		docSlice := make([]domain.Document, len(documents))
		for i, doc := range documents {
			docSlice[i] = *doc
		}
		expert.Documents = docSlice
	}

	engagements, err := s.ListEngagements(expert.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expert engagements: %w", err)
	}
	// Convert []*Engagement to []Engagement
	if len(engagements) > 0 {
		engSlice := make([]domain.Engagement, len(engagements))
		for i, eng := range engagements {
			engSlice[i] = *eng
		}
		expert.Engagements = engSlice
	}

	return &expert, nil
}

// GetExpertByEmail retrieves an expert by their email address
func (s *SQLiteStore) GetExpertByEmail(email string) (*domain.Expert, error) {
	query := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, 
		       e.is_bahraini, e.nationality, e.is_available, e.rating, e.role, 
		       e.employment_type, e.general_area, ea.name as general_area_name, 
		       e.specialized_area, e.is_trained, e.cv_path, e.approval_document_path, e.phone, e.email, 
		       e.is_published, e.biography, e.created_at, e.updated_at
		FROM experts e
		LEFT JOIN expert_areas ea ON e.general_area = ea.id
		WHERE e.email = ?
	`

	var expert domain.Expert
	var generalAreaName sql.NullString
	var nationality sql.NullString

	err := s.db.QueryRow(query, email).Scan(
		&expert.ID, &expert.ExpertID, &expert.Name, &expert.Designation, &expert.Institution,
		&expert.IsBahraini, &nationality, &expert.IsAvailable, &expert.Rating, &expert.Role,
		&expert.EmploymentType, &expert.GeneralArea, &generalAreaName,
		&expert.SpecializedArea, &expert.IsTrained, &expert.CVPath, &expert.ApprovalDocumentPath, &expert.Phone, &expert.Email,
		&expert.IsPublished, &expert.Biography, &expert.CreatedAt, &expert.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get expert by email: %w", err)
	}

	if generalAreaName.Valid {
		expert.GeneralAreaName = generalAreaName.String
	}

	if nationality.Valid {
		expert.Nationality = nationality.String
	}

	return &expert, nil
}

// UpdateExpert updates an existing expert in the database
func (s *SQLiteStore) UpdateExpert(expert *domain.Expert) error {
	// Get the current expert to avoid overwriting fields with empty values
	currentExpert, err := s.GetExpert(expert.ID)
	if err != nil {
		return fmt.Errorf("failed to get current expert data: %w", err)
	}

	// Only update fields that are explicitly set
	if expert.Name == "" {
		expert.Name = currentExpert.Name
	}
	if expert.Designation == "" {
		expert.Designation = currentExpert.Designation
	}
	if expert.Institution == "" {
		expert.Institution = currentExpert.Institution
	}
	if expert.Nationality == "" {
		expert.Nationality = currentExpert.Nationality
	}
	if expert.Rating == "" {
		expert.Rating = currentExpert.Rating
	}
	if expert.Role == "" {
		expert.Role = currentExpert.Role
	}
	if expert.EmploymentType == "" {
		expert.EmploymentType = currentExpert.EmploymentType
	}
	if expert.GeneralArea == 0 {
		expert.GeneralArea = currentExpert.GeneralArea
	}
	if expert.SpecializedArea == "" {
		expert.SpecializedArea = currentExpert.SpecializedArea
	}
	if expert.CVPath == "" {
		expert.CVPath = currentExpert.CVPath
	}
	if expert.ApprovalDocumentPath == "" {
		expert.ApprovalDocumentPath = currentExpert.ApprovalDocumentPath
	}
	if expert.Phone == "" {
		expert.Phone = currentExpert.Phone
	}
	if expert.Email == "" {
		expert.Email = currentExpert.Email
	}
	if expert.Biography == "" {
		expert.Biography = currentExpert.Biography
	}

	expert.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE experts SET
			name = ?, designation = ?, institution = ?, is_bahraini = ?,
			nationality = ?, is_available = ?, rating = ?, role = ?,
			employment_type = ?, general_area = ?, specialized_area = ?,
			is_trained = ?, cv_path = ?, approval_document_path = ?, phone = ?, email = ?,
			is_published = ?, biography = ?, updated_at = ?
		WHERE id = ?
	`

	_, err = s.db.Exec(
		query,
		expert.Name, expert.Designation, expert.Institution, expert.IsBahraini,
		expert.Nationality, expert.IsAvailable, expert.Rating, expert.Role,
		expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.ApprovalDocumentPath, expert.Phone, expert.Email,
		expert.IsPublished, expert.Biography, expert.UpdatedAt,
		expert.ID,
	)

	if err != nil {
		// Parse SQLite error to provide more specific error messages
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "email") {
				return fmt.Errorf("email already exists: %s", expert.Email)
			} else {
				return fmt.Errorf("unique constraint violation: %w", err)
			}
		} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return fmt.Errorf("referenced resource does not exist (check generalArea): %w", err)
		}
		
		return fmt.Errorf("failed to update expert: %w", err)
	}

	return nil
}

// DeleteExpert deletes an expert by ID and their associated documents
func (s *SQLiteStore) DeleteExpert(id int64) error {
	// Start a transaction to ensure all operations are atomic
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback will be no-op if transaction is committed
	
	// First, get the expert to verify existence and collect CV path and documents
	expert, err := s.GetExpert(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to get expert information: %w", err)
	}
	
	// Get list of document IDs associated with this expert
	var documentIDs []int64
	err = tx.QueryRow("SELECT COUNT(*) FROM expert_documents WHERE expert_id = ?", id).Scan(&documentIDs)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check for expert documents: %w", err)
	}
	
	// If there are associated documents, delete them first
	if len(expert.Documents) > 0 {
		// Delete the document files
		for _, doc := range expert.Documents {
			// First check if file exists
			filePath := doc.FilePath
			if filePath != "" {
				if _, err := os.Stat(filePath); err == nil {
					// File exists, delete it
					if err := os.Remove(filePath); err != nil {
						logger.Get().Warn("Failed to delete document file: %s - %v", filePath, err)
						// Log but continue - we still want to delete the database records
					} else {
						logger.Get().Debug("Deleted document file: %s", filePath)
					}
				}
			}
		}
		
		// Delete document records from database
		_, err = tx.Exec("DELETE FROM expert_documents WHERE expert_id = ?", id)
		if err != nil {
			return fmt.Errorf("failed to delete expert documents: %w", err)
		}
	}
	
	// Delete CV file if exists
	if expert.CVPath != "" {
		if _, err := os.Stat(expert.CVPath); err == nil {
			if err := os.Remove(expert.CVPath); err != nil {
				logger.Get().Warn("Failed to delete CV file: %s - %v", expert.CVPath, err)
			} else {
				logger.Get().Debug("Deleted CV file: %s", expert.CVPath)
			}
		}
	}
	
	// Delete approval document if exists
	if expert.ApprovalDocumentPath != "" {
		if _, err := os.Stat(expert.ApprovalDocumentPath); err == nil {
			if err := os.Remove(expert.ApprovalDocumentPath); err != nil {
				logger.Get().Warn("Failed to delete approval document: %s - %v", expert.ApprovalDocumentPath, err)
			} else {
				logger.Get().Debug("Deleted approval document: %s", expert.ApprovalDocumentPath)
			}
		}
	}
	
	// Now delete the expert record
	result, err := tx.Exec("DELETE FROM experts WHERE id = ?", id)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return fmt.Errorf("cannot delete expert: it is referenced by other records: %w", err)
		}
		return fmt.Errorf("failed to delete expert: %w", err)
	}
	
	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ListExperts retrieves a paginated list of experts with filters
func (s *SQLiteStore) ListExperts(filters map[string]interface{}, limit, offset int) ([]*domain.Expert, error) {
	// Build the query with filters
	queryBase := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, 
		       e.is_bahraini, e.nationality, e.is_available, e.rating, e.role, 
		       e.employment_type, e.general_area, ea.name as general_area_name, 
		       e.specialized_area, e.is_trained, e.cv_path, e.approval_document_path, e.phone, e.email, 
		       e.is_published, e.biography, e.created_at, e.updated_at
		FROM experts e
		LEFT JOIN expert_areas ea ON e.general_area = ea.id
	`

	// Add WHERE clause and parameters if filters are provided
	whereClause, params := buildWhereClauseForExpertFilters(filters)
	if whereClause != "" {
		queryBase += " WHERE " + whereClause
	}

	// Add dynamic ORDER BY based on sort_by and sort_order parameters
	// Default to updated_at DESC if not specified
	sortBy := "e.updated_at"
	sortOrder := "DESC"

	if val, ok := filters["sort_by"]; ok && val != "" {
		// To prevent SQL injection, validate against a whitelist of column names
		sortByStr := val.(string)
		
		// Mapping of allowed sort fields to their actual database column expressions
		allowedSortFields := map[string]string{
			"name":            "e.name",
			"expert_id":       "e.expert_id",
			"institution":     "e.institution",
			"designation":     "e.designation",
			"role":            "e.role",
			"employment_type": "e.employment_type",
			"nationality":     "e.nationality",
			"specialized_area": "e.specialized_area",
			"general_area":    "e.general_area",
			"rating":          "e.rating",
			"created_at":      "e.created_at",
			"updated_at":      "e.updated_at",
			"is_bahraini":     "e.is_bahraini",
			"is_available":    "e.is_available",
			"is_published":    "e.is_published",
		}
		
		if columnExpr, exists := allowedSortFields[sortByStr]; exists {
			sortBy = columnExpr
		}
	}

	if val, ok := filters["sort_order"]; ok && val != "" {
		orderStr := strings.ToUpper(val.(string))
		if orderStr == "ASC" || orderStr == "DESC" {
			sortOrder = orderStr
		}
	}

	queryBase += " ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	params = append(params, limit, offset)

	// Execute query
	rows, err := s.db.Query(queryBase, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query experts: %w", err)
	}
	defer rows.Close()

	var experts []*domain.Expert
	for rows.Next() {
		var expert domain.Expert
		var generalAreaName sql.NullString
		var nationality sql.NullString
		var expertID sql.NullString
		var name sql.NullString
		var designation sql.NullString
		var institution sql.NullString
		var rating sql.NullString
		var role sql.NullString
		var employmentType sql.NullString
		var specializedArea sql.NullString
		var cvPath sql.NullString
		var approvalDocumentPath sql.NullString
		var phone sql.NullString
		var email sql.NullString
		var biography sql.NullString
		var createdAt sql.NullTime
		var updatedAt sql.NullTime

		err := rows.Scan(
			&expert.ID, &expertID, &name, &designation, &institution,
			&expert.IsBahraini, &nationality, &expert.IsAvailable, &rating, &role,
			&employmentType, &expert.GeneralArea, &generalAreaName,
			&specializedArea, &expert.IsTrained, &cvPath, &approvalDocumentPath, &phone, &email,
			&expert.IsPublished, &biography, &createdAt, &updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan expert row: %w", err)
		}
		
		// Set createdAt and updatedAt from NullTime
		if createdAt.Valid {
			expert.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			expert.UpdatedAt = updatedAt.Time
		}

		// Assign NULL-safe values to expert struct
		if expertID.Valid {
			expert.ExpertID = expertID.String
		}
		if name.Valid {
			expert.Name = name.String
		}
		if designation.Valid {
			expert.Designation = designation.String
		}
		if institution.Valid {
			expert.Institution = institution.String
		}
		if rating.Valid {
			expert.Rating = rating.String
		}
		if role.Valid {
			expert.Role = role.String
		}
		if employmentType.Valid {
			expert.EmploymentType = employmentType.String
		}
		if specializedArea.Valid {
			expert.SpecializedArea = specializedArea.String
		}
		if cvPath.Valid {
			expert.CVPath = cvPath.String
		}
		if approvalDocumentPath.Valid {
			expert.ApprovalDocumentPath = approvalDocumentPath.String
		}
		if phone.Valid {
			expert.Phone = phone.String
		}
		if email.Valid {
			expert.Email = email.String
		}
		if biography.Valid {
			expert.Biography = biography.String
		}
		if generalAreaName.Valid {
			expert.GeneralAreaName = generalAreaName.String
		}
		if nationality.Valid {
			expert.Nationality = nationality.String
		}

		experts = append(experts, &expert)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over expert rows: %w", err)
	}

	return experts, nil
}

// CountExperts counts the total number of experts matching the given filters
func (s *SQLiteStore) CountExperts(filters map[string]interface{}) (int, error) {
	queryBase := "SELECT COUNT(*) FROM experts e"

	// Add WHERE clause if filters are provided
	whereClause, params := buildWhereClauseForExpertFilters(filters)
	if whereClause != "" {
		queryBase += " WHERE " + whereClause
	}

	var count int
	err := s.db.QueryRow(queryBase, params...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count experts: %w", err)
	}

	return count, nil
}

// Helper function to extract constraint name from SQLite error message
func extractConstraintName(errMsg string) string {
	// Example: "UNIQUE constraint failed: experts.expert_id"
	parts := strings.Split(errMsg, ":")
	if len(parts) < 2 {
		return "unknown field"
	}
	
	fieldPart := strings.TrimSpace(parts[len(parts)-1])
	fieldParts := strings.Split(fieldPart, ".")
	if len(fieldParts) < 2 {
		return fieldPart
	}
	
	// Convert snake_case to readable format
	field := fieldParts[1]
	field = strings.ReplaceAll(field, "_", " ")
	return field
}

// Helper function to extract column name from SQLite error message
func extractColumnName(errMsg string) string {
	// Example: "NOT NULL constraint failed: experts.name"
	parts := strings.Split(errMsg, ":")
	if len(parts) < 2 {
		return "unknown field"
	}
	
	fieldPart := strings.TrimSpace(parts[len(parts)-1])
	fieldParts := strings.Split(fieldPart, ".")
	if len(fieldParts) < 2 {
		return fieldPart
	}
	
	// Convert snake_case to readable format
	field := fieldParts[1]
	field = strings.ReplaceAll(field, "_", " ")
	return field
}

// Helper function to build WHERE clause for expert filters
func buildWhereClauseForExpertFilters(filters map[string]interface{}) (string, []interface{}) {
	var conditions []string
	var params []interface{}

	// Add conditions based on filters
	if val, ok := filters["name"]; ok && val != "" {
		conditions = append(conditions, "e.name LIKE ?")
		params = append(params, "%"+val.(string)+"%")
	}

	if val, ok := filters["institution"]; ok && val != "" {
		conditions = append(conditions, "e.institution LIKE ?")
		params = append(params, "%"+val.(string)+"%")
	}

	if val, ok := filters["role"]; ok && val != "" {
		conditions = append(conditions, "e.role = ?")
		params = append(params, val)
	}

	if val, ok := filters["generalArea"]; ok && val != 0 {
		conditions = append(conditions, "e.general_area = ?")
		params = append(params, val)
	}

	// Add new filters for nationality, specialized area, and employment type
	if val, ok := filters["by_nationality"]; ok {
		conditions = append(conditions, "e.is_bahraini = ?")
		isBahraini := val == "Bahraini" || val == "bahraini" || val == true
		params = append(params, isBahraini)
	}

	if val, ok := filters["by_specialized_area"]; ok && val != "" {
		conditions = append(conditions, "e.specialized_area LIKE ?")
		params = append(params, "%"+val.(string)+"%")
	}

	if val, ok := filters["by_employment_type"]; ok && val != "" {
		conditions = append(conditions, "e.employment_type = ?")
		params = append(params, val)
	}

	if val, ok := filters["by_role"]; ok && val != "" {
		conditions = append(conditions, "e.role = ?")
		params = append(params, val)
	}

	if val, ok := filters["by_general_area"]; ok && val != 0 {
		conditions = append(conditions, "e.general_area = ?")
		params = append(params, val)
	}

	// Original filters (for backward compatibility)
	if val, ok := filters["isBahraini"]; ok {
		conditions = append(conditions, "e.is_bahraini = ?")
		params = append(params, val)
	}

	if val, ok := filters["isAvailable"]; ok {
		conditions = append(conditions, "e.is_available = ?")
		params = append(params, val)
	}

	if val, ok := filters["isPublished"]; ok {
		conditions = append(conditions, "e.is_published = ?")
		params = append(params, val)
	}

	// Combine conditions with AND
	whereClause := ""
	if len(conditions) > 0 {
		for i, condition := range conditions {
			if i == 0 {
				whereClause = condition
			} else {
				whereClause += " AND " + condition
			}
		}
	}

	return whereClause, params
}

// NOTE: ListAreas and GetArea implementations are in area.go
