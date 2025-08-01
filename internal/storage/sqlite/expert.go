package sqlite

import (
	"database/sql"
	"encoding/json"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)


// CreateExpert creates a new expert in the database
func (s *SQLiteStore) CreateExpert(expert *domain.Expert) (int64, error) {
	// Begin transaction for atomic operations
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO experts (
			name, designation, affiliation, is_bahraini, is_available, rating,
			role, employment_type, general_area, specialized_area, is_trained,
			cv_document_id, approval_document_id, phone, email, is_published, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if expert.CreatedAt.IsZero() {
		expert.CreatedAt = time.Now()
		expert.UpdatedAt = expert.CreatedAt
	}

	result, err := tx.Exec(
		query,
		expert.Name, expert.Designation, expert.Affiliation,
		expert.IsBahraini, expert.IsAvailable, expert.Rating,
		expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVDocumentID, expert.ApprovalDocumentID, expert.Phone, expert.Email, expert.IsPublished,
		expert.CreatedAt, expert.UpdatedAt,
	)

	if err != nil {
		// Parse SQLite error to provide more specific error messages
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "email") {
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

	// Insert experience entries
	for _, exp := range expert.ExperienceEntries {
		expQuery := `
			INSERT INTO expert_experience_entries (
				expert_id, organization, position, start_date, end_date, is_current, country, description, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(expQuery, id, exp.Organization, exp.Position, exp.StartDate, exp.EndDate, exp.IsCurrent, exp.Country, exp.Description, time.Now(), time.Now())
		if err != nil {
			return 0, fmt.Errorf("failed to insert experience entry: %w", err)
		}
	}

	// Insert education entries
	for _, edu := range expert.EducationEntries {
		eduQuery := `
			INSERT INTO expert_education_entries (
				expert_id, institution, degree, field_of_study, graduation_year, country, description, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(eduQuery, id, edu.Institution, edu.Degree, edu.FieldOfStudy, edu.GraduationYear, edu.Country, edu.Description, time.Now(), time.Now())
		if err != nil {
			return 0, fmt.Errorf("failed to insert education entry: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

// GetExpert retrieves an expert by their ID
func (s *SQLiteStore) GetExpert(id int64) (*domain.Expert, error) {
	query := `
		SELECT e.id, e.name, e.designation, e.affiliation, 
		       e.is_bahraini, e.is_available, e.rating, e.role, 
		       e.employment_type, e.general_area, ea.name as general_area_name, 
		       e.specialized_area, e.is_trained, e.cv_document_id, e.approval_document_id, e.phone, e.email, 
		       e.is_published, e.created_at, e.updated_at,
		       COALESCE(
		           (SELECT GROUP_CONCAT(sa.name, ', ')
		           FROM specialized_areas sa
		           WHERE ',' || e.specialized_area || ',' LIKE '%,' || sa.id || ',%'
		           AND e.specialized_area IS NOT NULL 
		           AND e.specialized_area != ''),
		           ''
		       ) as specialized_area_names
		FROM experts e
		LEFT JOIN expert_areas ea ON e.general_area = ea.id
		WHERE e.id = ?
	`

	var expert domain.Expert
	var generalAreaName sql.NullString
	var specializedAreaNames sql.NullString
	var cvDocumentID sql.NullInt64
	var approvalDocumentID sql.NullInt64
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	err := s.db.QueryRow(query, id).Scan(
		&expert.ID, &expert.Name, &expert.Designation, &expert.Affiliation,
		&expert.IsBahraini, &expert.IsAvailable, &expert.Rating, &expert.Role,
		&expert.EmploymentType, &expert.GeneralArea, &generalAreaName,
		&expert.SpecializedArea, &expert.IsTrained, &cvDocumentID, &approvalDocumentID, &expert.Phone, &expert.Email,
		&expert.IsPublished, &createdAt, &updatedAt,
		&specializedAreaNames,
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

	if specializedAreaNames.Valid {
		expert.SpecializedAreaNames = specializedAreaNames.String
	}

	if cvDocumentID.Valid {
		expert.CVDocumentID = &cvDocumentID.Int64
	}

	if approvalDocumentID.Valid {
		expert.ApprovalDocumentID = &approvalDocumentID.Int64
	}

	if createdAt.Valid {
		expert.CreatedAt = createdAt.Time
	}

	if updatedAt.Valid {
		expert.UpdatedAt = updatedAt.Time
	}

	// Resolve document references
	expert.ResolveCVDocument(s.GetDocument)
	expert.ResolveApprovalDocument(s.GetDocument)

	// Populate bio data
	err = s.populateBioData(&expert)
	if err != nil {
		return nil, fmt.Errorf("failed to populate bio data: %w", err)
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

	engagements, err := s.ListEngagements(expert.ID, "", 100, 0) // empty string for all types, default limit
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
		SELECT e.id, e.name, e.designation, e.affiliation, 
		       e.is_bahraini, e.is_available, e.rating, e.role, 
		       e.employment_type, e.general_area, ea.name as general_area_name, 
		       e.specialized_area, e.is_trained, e.cv_document_id, e.approval_document_id, e.phone, e.email, 
		       e.is_published, e.created_at, e.updated_at,
		       COALESCE(
		           (SELECT GROUP_CONCAT(sa.name, ', ')
		           FROM specialized_areas sa
		           WHERE ',' || e.specialized_area || ',' LIKE '%,' || sa.id || ',%'
		           AND e.specialized_area IS NOT NULL 
		           AND e.specialized_area != ''),
		           ''
		       ) as specialized_area_names
		FROM experts e
		LEFT JOIN expert_areas ea ON e.general_area = ea.id
		WHERE e.email = ?
	`

	var expert domain.Expert
	var generalAreaName sql.NullString
	var specializedAreaNames sql.NullString
	var cvDocumentID sql.NullInt64
	var approvalDocumentID sql.NullInt64

	err := s.db.QueryRow(query, email).Scan(
		&expert.ID, &expert.Name, &expert.Designation, &expert.Affiliation,
		&expert.IsBahraini, &expert.IsAvailable, &expert.Rating, &expert.Role,
		&expert.EmploymentType, &expert.GeneralArea, &generalAreaName,
		&expert.SpecializedArea, &expert.IsTrained, &cvDocumentID, &approvalDocumentID, &expert.Phone, &expert.Email,
		&expert.IsPublished, &expert.CreatedAt, &expert.UpdatedAt,
		&specializedAreaNames,
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

	if specializedAreaNames.Valid {
		expert.SpecializedAreaNames = specializedAreaNames.String
	}

	if cvDocumentID.Valid {
		expert.CVDocumentID = &cvDocumentID.Int64
	}

	if approvalDocumentID.Valid {
		expert.ApprovalDocumentID = &approvalDocumentID.Int64
	}

	// Resolve document references
	expert.ResolveCVDocument(s.GetDocument)
	expert.ResolveApprovalDocument(s.GetDocument)

	// Fetch experience entries
	expQuery := `
		SELECT id, organization, position, start_date, end_date, is_current, country, description, created_at, updated_at
		FROM expert_experience_entries
		WHERE expert_id = ?
		ORDER BY created_at DESC
	`
	expRows, err := s.db.Query(expQuery, expert.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experience entries: %w", err)
	}
	defer expRows.Close()

	var experienceEntries []domain.ExpertExperienceEntry
	for expRows.Next() {
		var exp domain.ExpertExperienceEntry
		err := expRows.Scan(&exp.ID, &exp.Organization, &exp.Position, &exp.StartDate, &exp.EndDate, &exp.IsCurrent, &exp.Country, &exp.Description, &exp.CreatedAt, &exp.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan experience entry: %w", err)
		}
		exp.ExpertID = expert.ID
		experienceEntries = append(experienceEntries, exp)
	}
	expert.ExperienceEntries = experienceEntries

	// Fetch education entries
	eduQuery := `
		SELECT id, institution, degree, field_of_study, graduation_year, country, description, created_at, updated_at
		FROM expert_education_entries
		WHERE expert_id = ?
		ORDER BY graduation_year DESC
	`
	eduRows, err := s.db.Query(eduQuery, expert.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education entries: %w", err)
	}
	defer eduRows.Close()

	var educationEntries []domain.ExpertEducationEntry
	for eduRows.Next() {
		var edu domain.ExpertEducationEntry
		err := eduRows.Scan(&edu.ID, &edu.Institution, &edu.Degree, &edu.FieldOfStudy, &edu.GraduationYear, &edu.Country, &edu.Description, &edu.CreatedAt, &edu.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan education entry: %w", err)
		}
		edu.ExpertID = expert.ID
		educationEntries = append(educationEntries, edu)
	}
	expert.EducationEntries = educationEntries

	return &expert, nil
}

// UpdateExpert updates an existing expert in the database
func (s *SQLiteStore) UpdateExpert(expert *domain.Expert, editedBy int64) error {
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
	if expert.Affiliation == "" {
		expert.Affiliation = currentExpert.Affiliation
	}
	if expert.Rating == 0 {
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
	if expert.CVDocumentID == nil {
		expert.CVDocumentID = currentExpert.CVDocumentID
	}
	if expert.ApprovalDocumentID == nil {
		expert.ApprovalDocumentID = currentExpert.ApprovalDocumentID
	}
	if expert.Phone == "" {
		expert.Phone = currentExpert.Phone
	}
	if expert.Email == "" {
		expert.Email = currentExpert.Email
	}

	// Calculate changes for audit trail
	changedFields, oldValues, newValues := s.calculateExpertChanges(currentExpert, expert)
	if len(changedFields) == 0 {
		// No changes detected, return early
		return nil
	}

	// Set audit and timestamp fields
	now := time.Now()
	expert.UpdatedAt = now
	expert.LastEditedBy = &editedBy
	expert.LastEditedAt = &now

	// Begin transaction for atomic operations
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE experts SET
			name = ?, designation = ?, affiliation = ?, is_bahraini = ?,
			is_available = ?, rating = ?, role = ?,
			employment_type = ?, general_area = ?, specialized_area = ?,
			is_trained = ?, cv_document_id = ?, approval_document_id = ?, phone = ?, email = ?,
			is_published = ?, updated_at = ?, last_edited_by = ?, last_edited_at = ?
		WHERE id = ?
	`

	_, err = tx.Exec(
		query,
		expert.Name, expert.Designation, expert.Affiliation, expert.IsBahraini,
		expert.IsAvailable, expert.Rating, expert.Role,
		expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVDocumentID, expert.ApprovalDocumentID, expert.Phone, expert.Email,
		expert.IsPublished, expert.UpdatedAt, expert.LastEditedBy, expert.LastEditedAt,
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

	// Create audit history entry
	err = s.createExpertEditHistoryTx(tx, expert.ID, editedBy, changedFields, oldValues, newValues)
	if err != nil {
		return fmt.Errorf("failed to create audit history: %w", err)
	}

	// Update experience entries - delete existing and insert new ones
	if len(expert.ExperienceEntries) > 0 {
		// Delete existing experience entries
		_, err = tx.Exec("DELETE FROM expert_experience_entries WHERE expert_id = ?", expert.ID)
		if err != nil {
			return fmt.Errorf("failed to delete existing experience entries: %w", err)
		}

		// Insert new experience entries
		for _, exp := range expert.ExperienceEntries {
			expQuery := `
				INSERT INTO expert_experience_entries (
					expert_id, organization, position, start_date, end_date, is_current, country, description, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`
			_, err = tx.Exec(expQuery, expert.ID, exp.Organization, exp.Position, exp.StartDate, exp.EndDate, exp.IsCurrent, exp.Country, exp.Description, time.Now(), time.Now())
			if err != nil {
				return fmt.Errorf("failed to insert experience entry: %w", err)
			}
		}
	}

	// Update education entries - delete existing and insert new ones
	if len(expert.EducationEntries) > 0 {
		// Delete existing education entries
		_, err = tx.Exec("DELETE FROM expert_education_entries WHERE expert_id = ?", expert.ID)
		if err != nil {
			return fmt.Errorf("failed to delete existing education entries: %w", err)
		}

		// Insert new education entries
		for _, edu := range expert.EducationEntries {
			eduQuery := `
				INSERT INTO expert_education_entries (
					expert_id, institution, degree, field_of_study, graduation_year, country, description, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`
			_, err = tx.Exec(eduQuery, expert.ID, edu.Institution, edu.Degree, edu.FieldOfStudy, edu.GraduationYear, edu.Country, edu.Description, time.Now(), time.Now())
			if err != nil {
				return fmt.Errorf("failed to insert education entry: %w", err)
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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
	
	// Document files are managed by the document service
	// The document records will be cleaned up when the expert is deleted
	
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
		SELECT e.id, e.name, e.designation, e.affiliation, 
		       e.is_bahraini, e.is_available, e.rating, e.role, 
		       e.employment_type, e.general_area, ea.name as general_area_name, 
		       e.specialized_area, e.is_trained, e.cv_document_id, e.approval_document_id, e.phone, e.email, 
		       e.is_published, e.created_at, e.updated_at,
		       COALESCE(
		           (SELECT GROUP_CONCAT(sa.name, ', ')
		           FROM specialized_areas sa
		           WHERE ',' || e.specialized_area || ',' LIKE '%,' || sa.id || ',%'
		           AND e.specialized_area IS NOT NULL 
		           AND e.specialized_area != ''),
		           ''
		       ) as specialized_area_names
		FROM experts e
		LEFT JOIN expert_areas ea ON e.general_area = ea.id
	`

	// Add WHERE clause and parameters if filters are provided
	whereClause, params := buildWhereClauseForExpertFilters(filters)
	if whereClause != "" {
		queryBase += " WHERE " + whereClause
	}

	// Add dynamic ORDER BY based on sort_by and sort_order parameters
	// Default to id ASC if not specified
	sortBy := "e.id"
	sortOrder := "ASC"

	if val, ok := filters["sort_by"]; ok && val != "" {
		// To prevent SQL injection, validate against a whitelist of column names
		sortByStr := val.(string)
		
		// Mapping of allowed sort fields to their actual database column expressions
		allowedSortFields := map[string]string{
			"name":            "e.name",
			"id":              "e.id",
			"affiliation":     "e.affiliation",
			"designation":     "e.designation",
			"role":            "e.role",
			"employment_type": "e.employment_type",
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
		var specializedAreaNames sql.NullString
		var name sql.NullString
		var designation sql.NullString
		var affiliation sql.NullString
		var rating sql.NullInt32
		var role sql.NullString
		var employmentType sql.NullString
		var specializedArea sql.NullString
		var cvDocumentID sql.NullInt64
		var approvalDocumentID sql.NullInt64
		var phone sql.NullString
		var email sql.NullString
		var createdAt sql.NullTime
		var updatedAt sql.NullTime

		err := rows.Scan(
			&expert.ID, &name, &designation, &affiliation,
			&expert.IsBahraini, &expert.IsAvailable, &rating, &role,
			&employmentType, &expert.GeneralArea, &generalAreaName,
			&specializedArea, &expert.IsTrained, &cvDocumentID, &approvalDocumentID, &phone, &email,
			&expert.IsPublished, &createdAt, &updatedAt,
			&specializedAreaNames,
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
		if name.Valid {
			expert.Name = name.String
		}
		if designation.Valid {
			expert.Designation = designation.String
		}
		if affiliation.Valid {
			expert.Affiliation = affiliation.String
		}
		if rating.Valid {
			expert.Rating = int(rating.Int32)
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
		if cvDocumentID.Valid {
			cvDocID := cvDocumentID.Int64
			expert.CVDocumentID = &cvDocID
		}
		if approvalDocumentID.Valid {
			approvalDocID := approvalDocumentID.Int64
			expert.ApprovalDocumentID = &approvalDocID
		}
		if phone.Valid {
			expert.Phone = phone.String
		}
		if email.Valid {
			expert.Email = email.String
		}
		if generalAreaName.Valid {
			expert.GeneralAreaName = generalAreaName.String
		}
		if specializedAreaNames.Valid {
			expert.SpecializedAreaNames = specializedAreaNames.String
		}

		experts = append(experts, &expert)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over expert rows: %w", err)
	}

	// Populate bio data for each expert
	for _, expert := range experts {
		err = s.populateBioData(expert)
		if err != nil {
			return nil, fmt.Errorf("failed to populate bio data for expert %d: %w", expert.ID, err)
		}
	}

	return experts, nil
}

// populateBioData is a helper function to populate experience and education entries for an expert
func (s *SQLiteStore) populateBioData(expert *domain.Expert) error {
	// Fetch experience entries
	expQuery := `
		SELECT id, organization, position, start_date, end_date, is_current, country, description, created_at, updated_at
		FROM expert_experience_entries
		WHERE expert_id = ?
		ORDER BY created_at DESC
	`
	expRows, err := s.db.Query(expQuery, expert.ID)
	if err != nil {
		return fmt.Errorf("failed to get experience entries: %w", err)
	}
	defer expRows.Close()

	var experienceEntries []domain.ExpertExperienceEntry
	for expRows.Next() {
		var exp domain.ExpertExperienceEntry
		err := expRows.Scan(&exp.ID, &exp.Organization, &exp.Position, &exp.StartDate, &exp.EndDate, &exp.IsCurrent, &exp.Country, &exp.Description, &exp.CreatedAt, &exp.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan experience entry: %w", err)
		}
		exp.ExpertID = expert.ID
		experienceEntries = append(experienceEntries, exp)
	}
	expert.ExperienceEntries = experienceEntries

	// Fetch education entries
	eduQuery := `
		SELECT id, institution, degree, field_of_study, graduation_year, country, description, created_at, updated_at
		FROM expert_education_entries
		WHERE expert_id = ?
		ORDER BY graduation_year DESC
	`
	eduRows, err := s.db.Query(eduQuery, expert.ID)
	if err != nil {
		return fmt.Errorf("failed to get education entries: %w", err)
	}
	defer eduRows.Close()

	var educationEntries []domain.ExpertEducationEntry
	for eduRows.Next() {
		var edu domain.ExpertEducationEntry
		err := eduRows.Scan(&edu.ID, &edu.Institution, &edu.Degree, &edu.FieldOfStudy, &edu.GraduationYear, &edu.Country, &edu.Description, &edu.CreatedAt, &edu.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan education entry: %w", err)
		}
		edu.ExpertID = expert.ID
		educationEntries = append(educationEntries, edu)
	}
	expert.EducationEntries = educationEntries

	return nil
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

// Helper function to parse multiple values from comma-separated string
func parseMultiValue(param string) []string {
	values := strings.Split(param, ",")
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// Helper function to build IN clause for multiple values
func buildInClause(field string, values []string) (string, []interface{}) {
	if len(values) == 0 {
		return "", nil
	}
	
	if len(values) == 1 {
		return fmt.Sprintf("%s = ?", field), []interface{}{values[0]}
	}
	
	placeholders := strings.Repeat("?,", len(values))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma
	
	params := make([]interface{}, len(values))
	for i, value := range values {
		params[i] = value
	}
	
	return fmt.Sprintf("%s IN (%s)", field, placeholders), params
}

// Helper function to build LIKE clause for multiple values (OR logic within field)
func buildLikeClause(field string, values []string) (string, []interface{}) {
	if len(values) == 0 {
		return "", nil
	}
	
	if len(values) == 1 {
		return fmt.Sprintf("%s LIKE ?", field), []interface{}{"%" + values[0] + "%"}
	}
	
	conditions := make([]string, len(values))
	params := make([]interface{}, len(values))
	
	for i, value := range values {
		conditions[i] = fmt.Sprintf("%s LIKE ?", field)
		params[i] = "%" + value + "%"
	}
	
	return fmt.Sprintf("(%s)", strings.Join(conditions, " OR ")), params
}

// Helper function to build WHERE clause for expert filters - Clean implementation
func buildWhereClauseForExpertFilters(filters map[string]interface{}) (string, []interface{}) {
	var conditions []string
	var params []interface{}

	// Multi-value filters with exact matching
	exactMatchFilters := map[string]string{
		"role":            "e.role",
		"employment_type": "e.employment_type",
	}

	for filterKey, dbField := range exactMatchFilters {
		if val, ok := filters[filterKey]; ok && val != "" {
			values := parseMultiValue(val.(string))
			if len(values) > 0 {
				condition, filterParams := buildInClause(dbField, values)
				if condition != "" {
					conditions = append(conditions, condition)
					params = append(params, filterParams...)
				}
			}
		}
	}

	// General area filter (integer values)
	if val, ok := filters["general_area"]; ok && val != "" {
		values := parseMultiValue(val.(string))
		if len(values) > 0 {
			// Convert string values to integers for general_area
			intValues := make([]string, 0, len(values))
			for _, strVal := range values {
				if _, err := strconv.ParseInt(strVal, 10, 64); err == nil {
					intValues = append(intValues, strVal)
				}
			}
			if len(intValues) > 0 {
				condition, filterParams := buildInClause("e.general_area", intValues)
				if condition != "" {
					conditions = append(conditions, condition)
					params = append(params, filterParams...)
				}
			}
		}
	}

	// Text-based filters with LIKE matching
	likeMatchFilters := map[string]string{
		"affiliation":      "e.affiliation",
		"specialized_area": "e.specialized_area",
	}

	for filterKey, dbField := range likeMatchFilters {
		if val, ok := filters[filterKey]; ok && val != "" {
			values := parseMultiValue(val.(string))
			if len(values) > 0 {
				condition, filterParams := buildLikeClause(dbField, values)
				if condition != "" {
					conditions = append(conditions, condition)
					params = append(params, filterParams...)
				}
			}
		}
	}

	// Boolean filters (single value only)
	booleanFilters := map[string]string{
		"is_available": "e.is_available",
		"is_published": "e.is_published",
		"is_bahraini":  "e.is_bahraini",
	}

	for filterKey, dbField := range booleanFilters {
		if val, ok := filters[filterKey]; ok {
			conditions = append(conditions, fmt.Sprintf("%s = ?", dbField))
			params = append(params, val)
		}
	}

	// Combine conditions with AND
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = strings.Join(conditions, " AND ")
	}

	return whereClause, params
}

// UpdateExpertCVDocument updates the CV document reference for an expert
func (s *SQLiteStore) UpdateExpertCVDocument(expertID, documentID int64) error {
	return s.updateDocumentReference("experts", "cv_document_id", expertID, documentID)
}

// UpdateExpertApprovalDocument updates the approval document reference for an expert
func (s *SQLiteStore) UpdateExpertApprovalDocument(expertID, documentID int64) error {
	return s.updateDocumentReference("experts", "approval_document_id", expertID, documentID)
}

// updateDocumentReference is a generic helper for updating document references
func (s *SQLiteStore) updateDocumentReference(table, column string, entityID, documentID int64) error {
	query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", table, column)
	
	result, err := s.db.Exec(query, documentID, entityID)
	if err != nil {
		return fmt.Errorf("failed to update %s document reference: %w", table, err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	return nil
}

// Helper methods for audit trail functionality

// calculateExpertChanges compares old and new expert data to identify changes
func (s *SQLiteStore) calculateExpertChanges(oldExpert, newExpert *domain.Expert) ([]string, map[string]interface{}, map[string]interface{}) {
	var changedFields []string
	oldValues := make(map[string]interface{})
	newValues := make(map[string]interface{})

	// Compare each field and track changes
	if oldExpert.Name != newExpert.Name {
		changedFields = append(changedFields, "name")
		oldValues["name"] = oldExpert.Name
		newValues["name"] = newExpert.Name
	}
	if oldExpert.Designation != newExpert.Designation {
		changedFields = append(changedFields, "designation")
		oldValues["designation"] = oldExpert.Designation
		newValues["designation"] = newExpert.Designation
	}
	if oldExpert.Affiliation != newExpert.Affiliation {
		changedFields = append(changedFields, "affiliation")
		oldValues["affiliation"] = oldExpert.Affiliation
		newValues["affiliation"] = newExpert.Affiliation
	}
	if oldExpert.IsBahraini != newExpert.IsBahraini {
		changedFields = append(changedFields, "isBahraini")
		oldValues["isBahraini"] = oldExpert.IsBahraini
		newValues["isBahraini"] = newExpert.IsBahraini
	}
	if oldExpert.IsAvailable != newExpert.IsAvailable {
		changedFields = append(changedFields, "isAvailable")
		oldValues["isAvailable"] = oldExpert.IsAvailable
		newValues["isAvailable"] = newExpert.IsAvailable
	}
	if oldExpert.Rating != newExpert.Rating {
		changedFields = append(changedFields, "rating")
		oldValues["rating"] = oldExpert.Rating
		newValues["rating"] = newExpert.Rating
	}
	if oldExpert.Role != newExpert.Role {
		changedFields = append(changedFields, "role")
		oldValues["role"] = oldExpert.Role
		newValues["role"] = newExpert.Role
	}
	if oldExpert.EmploymentType != newExpert.EmploymentType {
		changedFields = append(changedFields, "employmentType")
		oldValues["employmentType"] = oldExpert.EmploymentType
		newValues["employmentType"] = newExpert.EmploymentType
	}
	if oldExpert.GeneralArea != newExpert.GeneralArea {
		changedFields = append(changedFields, "generalArea")
		oldValues["generalArea"] = oldExpert.GeneralArea
		newValues["generalArea"] = newExpert.GeneralArea
	}
	if oldExpert.SpecializedArea != newExpert.SpecializedArea {
		changedFields = append(changedFields, "specializedArea")
		oldValues["specializedArea"] = oldExpert.SpecializedArea
		newValues["specializedArea"] = newExpert.SpecializedArea
	}
	if oldExpert.IsTrained != newExpert.IsTrained {
		changedFields = append(changedFields, "isTrained")
		oldValues["isTrained"] = oldExpert.IsTrained
		newValues["isTrained"] = newExpert.IsTrained
	}
	if oldExpert.Phone != newExpert.Phone {
		changedFields = append(changedFields, "phone")
		oldValues["phone"] = oldExpert.Phone
		newValues["phone"] = newExpert.Phone
	}
	if oldExpert.Email != newExpert.Email {
		changedFields = append(changedFields, "email")
		oldValues["email"] = oldExpert.Email
		newValues["email"] = newExpert.Email
	}
	if oldExpert.IsPublished != newExpert.IsPublished {
		changedFields = append(changedFields, "isPublished")
		oldValues["isPublished"] = oldExpert.IsPublished
		newValues["isPublished"] = newExpert.IsPublished
	}

	// Compare document IDs
	oldCVID := int64(0)
	newCVID := int64(0)
	if oldExpert.CVDocumentID != nil {
		oldCVID = *oldExpert.CVDocumentID
	}
	if newExpert.CVDocumentID != nil {
		newCVID = *newExpert.CVDocumentID
	}
	if oldCVID != newCVID {
		changedFields = append(changedFields, "cvDocumentId")
		oldValues["cvDocumentId"] = oldCVID
		newValues["cvDocumentId"] = newCVID
	}

	oldApprovalID := int64(0)
	newApprovalID := int64(0)
	if oldExpert.ApprovalDocumentID != nil {
		oldApprovalID = *oldExpert.ApprovalDocumentID
	}
	if newExpert.ApprovalDocumentID != nil {
		newApprovalID = *newExpert.ApprovalDocumentID
	}
	if oldApprovalID != newApprovalID {
		changedFields = append(changedFields, "approvalDocumentId")
		oldValues["approvalDocumentId"] = oldApprovalID
		newValues["approvalDocumentId"] = newApprovalID
	}

	return changedFields, oldValues, newValues
}

// createExpertEditHistoryTx creates an audit history entry within a transaction
func (s *SQLiteStore) createExpertEditHistoryTx(tx *sql.Tx, expertID, editedBy int64, changedFields []string, oldValues, newValues map[string]interface{}) error {
	
	changedFieldsJSON, err := json.Marshal(changedFields)
	if err != nil {
		return fmt.Errorf("failed to marshal changed fields: %w", err)
	}

	oldValuesJSON, err := json.Marshal(oldValues)
	if err != nil {
		return fmt.Errorf("failed to marshal old values: %w", err)
	}

	newValuesJSON, err := json.Marshal(newValues)
	if err != nil {
		return fmt.Errorf("failed to marshal new values: %w", err)
	}

	query := `
		INSERT INTO expert_edit_history (
			expert_id, edited_by, edited_at, fields_changed, old_values, new_values
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = tx.Exec(query, expertID, editedBy, time.Now(), string(changedFieldsJSON), string(oldValuesJSON), string(newValuesJSON))
	if err != nil {
		return fmt.Errorf("failed to insert edit history: %w", err)
	}

	return nil
}

// GetExpertEditHistory retrieves edit history for an expert
func (s *SQLiteStore) GetExpertEditHistory(expertID int64) ([]*domain.ExpertEditHistoryEntry, error) {
	query := `
		SELECT eh.id, eh.expert_id, eh.edited_by, eh.edited_at, eh.fields_changed, 
		       eh.old_values, eh.new_values, eh.change_reason, u.name
		FROM expert_edit_history eh
		LEFT JOIN users u ON eh.edited_by = u.id
		WHERE eh.expert_id = ?
		ORDER BY eh.edited_at DESC
	`

	rows, err := s.db.Query(query, expertID)
	if err != nil {
		return nil, fmt.Errorf("failed to query edit history: %w", err)
	}
	defer rows.Close()

	var history []*domain.ExpertEditHistoryEntry
	for rows.Next() {
		entry := &domain.ExpertEditHistoryEntry{}
		var editorName sql.NullString
		var oldValues, newValues, changeReason sql.NullString

		err := rows.Scan(
			&entry.ID, &entry.ExpertID, &entry.EditedBy, &entry.EditedAt,
			&entry.FieldsChanged, &oldValues, &newValues, &changeReason, &editorName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan edit history row: %w", err)
		}

		if oldValues.Valid {
			entry.OldValues = &oldValues.String
		}
		if newValues.Valid {
			entry.NewValues = &newValues.String
		}
		if changeReason.Valid {
			entry.ChangeReason = &changeReason.String
		}
		if editorName.Valid {
			entry.EditorName = &editorName.String
		}

		history = append(history, entry)
	}

	return history, nil
}

// CreateExpertEditHistory creates a new edit history entry
func (s *SQLiteStore) CreateExpertEditHistory(entry *domain.ExpertEditHistoryEntry) error {
	query := `
		INSERT INTO expert_edit_history (
			expert_id, edited_by, edited_at, fields_changed, old_values, new_values, change_reason
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query, 
		entry.ExpertID, entry.EditedBy, entry.EditedAt, entry.FieldsChanged,
		entry.OldValues, entry.NewValues, entry.ChangeReason,
	)
	if err != nil {
		return fmt.Errorf("failed to insert edit history: %w", err)
	}

	return nil
}

// NOTE: ListAreas and GetArea implementations are in area.go
