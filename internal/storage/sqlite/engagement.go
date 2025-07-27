package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	
	"expertdb/internal/domain"
)

// ListEngagements retrieves engagements with optional filtering by expert ID and engagement type
// If expertID is 0, it returns engagements for all experts
// If engagementType is empty, it returns all engagement types
func (s *SQLiteStore) ListEngagements(expertID int64, engagementType string, limit, offset int) ([]*domain.Engagement, error) {
	// Start building the query with filtering support
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT id, expert_id, engagement_type, start_date, end_date,
				project_name, status, feedback_score, notes, created_at
		FROM expert_engagements
		WHERE 1=1
	`)
	
	// Prepare query parameters
	var params []interface{}
	
	// Add expert_id filter if provided
	if expertID > 0 {
		queryBuilder.WriteString(" AND expert_id = ?")
		params = append(params, expertID)
	}
	
	// Add engagement_type filter if provided (Phase 11B: Restrict to validator/evaluator)
	if engagementType != "" {
		queryBuilder.WriteString(" AND engagement_type = ?")
		params = append(params, engagementType)
	}
	
	// Add ordering for consistent results
	queryBuilder.WriteString(" ORDER BY created_at DESC")
	
	// Add pagination
	if limit > 0 {
		queryBuilder.WriteString(" LIMIT ?")
		params = append(params, limit)
		
		if offset > 0 {
			queryBuilder.WriteString(" OFFSET ?")
			params = append(params, offset)
		}
	}
	
	// Execute the query
	query := queryBuilder.String()
	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to get engagements: %w", err)
	}
	defer rows.Close()
	
	var engagements []*domain.Engagement
	for rows.Next() {
		var engagement domain.Engagement
		var endDate sql.NullTime
		var projectName, notes sql.NullString
		var feedbackScore sql.NullInt32
		
		err := rows.Scan(
			&engagement.ID, &engagement.ExpertID, &engagement.EngagementType,
			&engagement.StartDate, &endDate, &projectName,
			&engagement.Status, &feedbackScore, &notes,
			&engagement.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan engagement row: %w", err)
		}
		
		// Set nullable fields
		if endDate.Valid {
			engagement.EndDate = endDate.Time
		}
		
		if projectName.Valid {
			engagement.ProjectName = projectName.String
		}
		
		if feedbackScore.Valid {
			engagement.FeedbackScore = int(feedbackScore.Int32)
		}
		
		if notes.Valid {
			engagement.Notes = notes.String
		}
		
		engagements = append(engagements, &engagement)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating engagement rows: %w", err)
	}
	
	return engagements, nil
}

// GetEngagement retrieves an engagement by ID
func (s *SQLiteStore) GetEngagement(id int64) (*domain.Engagement, error) {
	query := `
		SELECT id, expert_id, engagement_type, start_date, end_date,
				project_name, status, feedback_score, notes, created_at
		FROM expert_engagements
		WHERE id = ?
	`
	
	var engagement domain.Engagement
	var endDate sql.NullTime
	var projectName, notes sql.NullString
	var feedbackScore sql.NullInt32
	
	err := s.db.QueryRow(query, id).Scan(
		&engagement.ID, &engagement.ExpertID, &engagement.EngagementType,
		&engagement.StartDate, &endDate, &projectName,
		&engagement.Status, &feedbackScore, &notes,
		&engagement.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get engagement: %w", err)
	}
	
	// Set nullable fields
	if endDate.Valid {
		engagement.EndDate = endDate.Time
	}
	
	if projectName.Valid {
		engagement.ProjectName = projectName.String
	}
	
	if feedbackScore.Valid {
		engagement.FeedbackScore = int(feedbackScore.Int32)
	}
	
	if notes.Valid {
		engagement.Notes = notes.String
	}
	
	return &engagement, nil
}

// CreateEngagement creates a new engagement record
func (s *SQLiteStore) CreateEngagement(engagement *domain.Engagement) (int64, error) {
	// Phase 11B: Validate engagement type restriction to validator or evaluator
	if engagement.EngagementType != "validator" && engagement.EngagementType != "evaluator" {
		return 0, fmt.Errorf("engagement type must be 'validator' or 'evaluator'")
	}
	
	query := `
		INSERT INTO expert_engagements (
			expert_id, engagement_type, start_date, end_date,
			project_name, status, feedback_score, notes, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	// Set default values
	if engagement.CreatedAt.IsZero() {
		engagement.CreatedAt = time.Now()
	}
	
	if engagement.Status == "" {
		engagement.Status = "pending"
	}
	
	// Handle null values for optional fields
	var endDate interface{} = nil
	if !engagement.EndDate.IsZero() {
		endDate = engagement.EndDate
	}
	
	var projectName interface{} = nil
	if engagement.ProjectName != "" {
		projectName = engagement.ProjectName
	}
	
	var feedbackScore interface{} = nil
	if engagement.FeedbackScore > 0 {
		feedbackScore = engagement.FeedbackScore
	}
	
	var notes interface{} = nil
	if engagement.Notes != "" {
		notes = engagement.Notes
	}
	
	result, err := s.db.Exec(
		query,
		engagement.ExpertID, engagement.EngagementType, engagement.StartDate,
		endDate, projectName, engagement.Status, feedbackScore, notes,
		engagement.CreatedAt,
	)
	
	if err != nil {
		return 0, fmt.Errorf("failed to create engagement: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get engagement ID: %w", err)
	}
	
	engagement.ID = id
	return id, nil
}

// UpdateEngagement updates an existing engagement record
func (s *SQLiteStore) UpdateEngagement(engagement *domain.Engagement) error {
	// Get current engagement to avoid overwriting with empty values
	current, err := s.GetEngagement(engagement.ID)
	if err != nil {
		return fmt.Errorf("failed to get current engagement data: %w", err)
	}
	
	// Only update fields that are set
	if engagement.EngagementType == "" {
		engagement.EngagementType = current.EngagementType
	}
	
	// Phase 11B: Validate engagement type restriction to validator or evaluator
	if engagement.EngagementType != "validator" && engagement.EngagementType != "evaluator" {
		return fmt.Errorf("engagement type must be 'validator' or 'evaluator'")
	}
	
	if engagement.StartDate.IsZero() {
		engagement.StartDate = current.StartDate
	}
	
	// Status defaults to current if not provided
	if engagement.Status == "" {
		engagement.Status = current.Status
	}
	
	query := `
		UPDATE expert_engagements SET
			engagement_type = ?, start_date = ?, end_date = ?,
			project_name = ?, status = ?, feedback_score = ?, notes = ?
		WHERE id = ?
	`
	
	// Handle null values for optional fields
	var endDate interface{} = nil
	if !engagement.EndDate.IsZero() {
		endDate = engagement.EndDate
	} else if !current.EndDate.IsZero() {
		endDate = current.EndDate
	}
	
	var projectName interface{} = nil
	if engagement.ProjectName != "" {
		projectName = engagement.ProjectName
	} else if current.ProjectName != "" {
		projectName = current.ProjectName
	}
	
	var feedbackScore interface{} = nil
	if engagement.FeedbackScore > 0 {
		feedbackScore = engagement.FeedbackScore
	} else if current.FeedbackScore > 0 {
		feedbackScore = current.FeedbackScore
	}
	
	var notes interface{} = nil
	if engagement.Notes != "" {
		notes = engagement.Notes
	} else if current.Notes != "" {
		notes = current.Notes
	}
	
	_, err = s.db.Exec(
		query,
		engagement.EngagementType, engagement.StartDate, endDate,
		projectName, engagement.Status, feedbackScore, notes,
		engagement.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update engagement: %w", err)
	}
	
	return nil
}

// DeleteEngagement deletes an engagement by ID
func (s *SQLiteStore) DeleteEngagement(id int64) error {
	result, err := s.db.Exec("DELETE FROM expert_engagements WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete engagement: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	return nil
}

// ImportEngagements imports multiple engagements at once
// Returns count of successfully imported engagements and a map of errors for failed imports
func (s *SQLiteStore) ImportEngagements(engagements []*domain.Engagement) (int, map[int]error) {
	errors := make(map[int]error)
	successCount := 0
	
	// Start a transaction for the batch operation
	tx, err := s.db.Begin()
	if err != nil {
		errors[-1] = fmt.Errorf("failed to start transaction: %w", err)
		return 0, errors
	}
	
	// Prepare the insert statement
	stmt, err := tx.Prepare(`
		INSERT INTO expert_engagements (
			expert_id, engagement_type, start_date, end_date,
			project_name, status, feedback_score, notes, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		errors[-1] = fmt.Errorf("failed to prepare statement: %w", err)
		return 0, errors
	}
	defer stmt.Close()
	
	// Process each engagement
	for i, engagement := range engagements {
		// Validate expert exists
		var expertExists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM experts WHERE id = ?)", engagement.ExpertID).Scan(&expertExists)
		if err != nil {
			errors[i] = fmt.Errorf("failed to check if expert exists: %w", err)
			continue
		}
		
		if !expertExists {
			errors[i] = fmt.Errorf("expert with ID %d does not exist", engagement.ExpertID)
			continue
		}
		
		// Phase 11B: Validate engagement type restriction to validator or evaluator
		if engagement.EngagementType != "validator" && engagement.EngagementType != "evaluator" {
			errors[i] = fmt.Errorf("engagement type must be 'validator' or 'evaluator'")
			continue
		}
		
		// Set default values
		if engagement.CreatedAt.IsZero() {
			engagement.CreatedAt = time.Now()
		}
		
		if engagement.Status == "" {
			engagement.Status = "pending"
		}
		
		// Handle null values for optional fields
		var endDate interface{} = nil
		if !engagement.EndDate.IsZero() {
			endDate = engagement.EndDate
		}
		
		var projectName interface{} = nil
		if engagement.ProjectName != "" {
			projectName = engagement.ProjectName
		}
		
		var feedbackScore interface{} = nil
		if engagement.FeedbackScore > 0 {
			feedbackScore = engagement.FeedbackScore
		}
		
		var notes interface{} = nil
		if engagement.Notes != "" {
			notes = engagement.Notes
		}
		
		// Check for duplicates (same expert, type, start date, project)
		var duplicateExists bool
		err = s.db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM expert_engagements 
				WHERE expert_id = ? 
				AND engagement_type = ? 
				AND start_date = ? 
				AND (project_name = ? OR (project_name IS NULL AND ? IS NULL))
			)`,
			engagement.ExpertID,
			engagement.EngagementType,
			engagement.StartDate,
			projectName,
			projectName,
		).Scan(&duplicateExists)
		
		if err != nil {
			errors[i] = fmt.Errorf("failed to check for duplicates: %w", err)
			continue
		}
		
		if duplicateExists {
			errors[i] = fmt.Errorf("duplicate engagement found for expert %d with type %s on date %s",
				engagement.ExpertID, engagement.EngagementType, engagement.StartDate.Format("2006-01-02"))
			continue
		}
		
		// Execute the insert
		_, err = stmt.Exec(
			engagement.ExpertID, engagement.EngagementType, engagement.StartDate,
			endDate, projectName, engagement.Status, feedbackScore, notes,
			engagement.CreatedAt,
		)
		
		if err != nil {
			errors[i] = fmt.Errorf("failed to insert engagement: %w", err)
			continue
		}
		
		successCount++
	}
	
	// Commit or rollback the transaction
	if successCount > 0 {
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			errors[-1] = fmt.Errorf("failed to commit transaction: %w", err)
			return 0, errors
		}
	} else {
		tx.Rollback()
	}
	
	return successCount, errors
}