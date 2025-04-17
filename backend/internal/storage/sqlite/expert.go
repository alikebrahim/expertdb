package sqlite

import (
	"database/sql"
	"expertdb/internal/domain"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUniqueExpertID generates a unique ID for an expert
func (s *SQLiteStore) GenerateUniqueExpertID() (string, error) {
	id := "EXP-" + uuid.New().String()
	exists, err := s.ExpertIDExists(id)
	if err != nil {
		return "", fmt.Errorf("failed to check expert ID: %w", err)
	}
	if exists {
		return s.GenerateUniqueExpertID() // Recursive retry
	}
	return id, nil
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
			cv_path, phone, email, is_published, biography, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		expert.IsTrained, expert.CVPath, expert.Phone, expert.Email, expert.IsPublished,
		expert.Biography, expert.CreatedAt, expert.UpdatedAt,
	)

	if err != nil {
		// Parse SQLite error to provide more specific error messages
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "expert_id") {
				return 0, fmt.Errorf("expert ID already exists: %s", expert.ExpertID)
			} else if strings.Contains(err.Error(), "email") {
				return 0, fmt.Errorf("email already exists: %s", expert.Email)
			} else {
				return 0, fmt.Errorf("unique constraint violation: %w", err)
			}
		} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return 0, fmt.Errorf("referenced resource does not exist (check generalArea): %w", err)
		}
		
		return 0, fmt.Errorf("failed to create expert: %w", err)
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
		       e.specialized_area, e.is_trained, e.cv_path, e.phone, e.email, 
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
		&expert.SpecializedArea, &expert.IsTrained, &expert.CVPath, &expert.Phone, &expert.Email,
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
	expert.Documents = documents

	engagements, err := s.ListEngagements(expert.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expert engagements: %w", err)
	}
	expert.Engagements = engagements

	return &expert, nil
}

// GetExpertByEmail retrieves an expert by their email address
func (s *SQLiteStore) GetExpertByEmail(email string) (*domain.Expert, error) {
	query := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, 
		       e.is_bahraini, e.nationality, e.is_available, e.rating, e.role, 
		       e.employment_type, e.general_area, ea.name as general_area_name, 
		       e.specialized_area, e.is_trained, e.cv_path, e.phone, e.email, 
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
		&expert.SpecializedArea, &expert.IsTrained, &expert.CVPath, &expert.Phone, &expert.Email,
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
			is_trained = ?, cv_path = ?, phone = ?, email = ?,
			is_published = ?, biography = ?, updated_at = ?
		WHERE id = ?
	`

	_, err = s.db.Exec(
		query,
		expert.Name, expert.Designation, expert.Institution, expert.IsBahraini,
		expert.Nationality, expert.IsAvailable, expert.Rating, expert.Role,
		expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
		expert.IsTrained, expert.CVPath, expert.Phone, expert.Email,
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

// DeleteExpert deletes an expert by ID
func (s *SQLiteStore) DeleteExpert(id int64) error {
	result, err := s.db.Exec("DELETE FROM experts WHERE id = ?", id)
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

	return nil
}

// ListExperts retrieves a paginated list of experts with filters
func (s *SQLiteStore) ListExperts(filters map[string]interface{}, limit, offset int) ([]*domain.Expert, error) {
	// Build the query with filters
	queryBase := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, 
		       e.is_bahraini, e.nationality, e.is_available, e.rating, e.role, 
		       e.employment_type, e.general_area, ea.name as general_area_name, 
		       e.specialized_area, e.is_trained, e.cv_path, e.phone, e.email, 
		       e.is_published, e.biography, e.created_at, e.updated_at
		FROM experts e
		LEFT JOIN expert_areas ea ON e.general_area = ea.id
	`

	// Add WHERE clause and parameters if filters are provided
	whereClause, params := buildWhereClauseForExpertFilters(filters)
	if whereClause != "" {
		queryBase += " WHERE " + whereClause
	}

	// Add ORDER BY and LIMIT
	queryBase += " ORDER BY e.updated_at DESC LIMIT ? OFFSET ?"
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
		var phone sql.NullString
		var email sql.NullString
		var biography sql.NullString

		err := rows.Scan(
			&expert.ID, &expertID, &name, &designation, &institution,
			&expert.IsBahraini, &nationality, &expert.IsAvailable, &rating, &role,
			&employmentType, &expert.GeneralArea, &generalAreaName,
			&specializedArea, &expert.IsTrained, &cvPath, &phone, &email,
			&expert.IsPublished, &biography, &expert.CreatedAt, &expert.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan expert row: %w", err)
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

// ListAreas retrieves all expert areas
func (s *SQLiteStore) ListAreas() ([]*domain.Area, error) {
	query := "SELECT id, name FROM expert_areas ORDER BY name"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expert areas: %w", err)
	}
	defer rows.Close()

	var areas []*domain.Area
	for rows.Next() {
		var area domain.Area
		if err := rows.Scan(&area.ID, &area.Name); err != nil {
			return nil, fmt.Errorf("failed to scan area row: %w", err)
		}
		areas = append(areas, &area)
	}

	return areas, nil
}

// GetArea retrieves a specific area by its ID
func (s *SQLiteStore) GetArea(id int64) (*domain.Area, error) {
	var area domain.Area
	err := s.db.QueryRow("SELECT id, name FROM expert_areas WHERE id = ?", id).Scan(&area.ID, &area.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get expert area: %w", err)
	}
	return &area, nil
}
