package sqlite

import (
	"database/sql"
	"fmt"
	
	"expertdb/internal/domain"
)

// ListEngagements retrieves all engagements for an expert
func (s *SQLiteStore) ListEngagements(expertID int64) ([]*domain.Engagement, error) {
	query := `
		SELECT id, expert_id, engagement_type, start_date, end_date,
				project_name, status, feedback_score, notes, created_at
		FROM expert_engagements
		WHERE expert_id = ?
	`
	
	rows, err := s.db.Query(query, expertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expert engagements: %w", err)
	}
	defer rows.Close()
	
	var engagements []*domain.Engagement
	for rows.Next() {
		var engagement domain.Engagement
		var endDate, projectName sql.NullTime
		var feedbackScore sql.NullInt32
		var notes sql.NullString
		
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
	var endDate, projectName sql.NullTime
	var feedbackScore sql.NullInt32
	var notes sql.NullString
	
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
	query := `
		INSERT INTO expert_engagements (
			expert_id, engagement_type, start_date, end_date,
			project_name, status, feedback_score, notes, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	// Handle nullable fields
	var endDate, projectName interface{} = nil, nil
	var feedbackScore, notes interface{} = nil, nil
	
	// Set end date if provided
	if !engagement.EndDate.IsZero() {
		endDate = engagement.EndDate
	}
	
	// Set project name if provided
	if engagement.ProjectName != "" {
		projectName = engagement.ProjectName
	}
	
	// Set feedback score if provided
	if engagement.FeedbackScore > 0 {
		feedbackScore = engagement.FeedbackScore
	}
	
	// Set notes if provided
	if engagement.Notes != "" {
		notes = engagement.Notes
	}
	
	result, err := s.db.Exec(
		query,
		engagement.ExpertID, engagement.EngagementType, engagement.StartDate,
		endDate, projectName, engagement.Status,
		feedbackScore, notes, engagement.CreatedAt,
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

// UpdateEngagement updates an existing engagement
func (s *SQLiteStore) UpdateEngagement(engagement *domain.Engagement) error {
	// First, get current engagement to preserve values that aren't explicitly updated
	current, err := s.GetEngagement(engagement.ID)
	if err != nil {
		return fmt.Errorf("failed to get current engagement state: %w", err)
	}
	
	query := `
		UPDATE expert_engagements
		SET expert_id = ?, engagement_type = ?, start_date = ?, end_date = ?,
			project_name = ?, status = ?, feedback_score = ?, notes = ?
		WHERE id = ?
	`
	
	// Handle nullable fields - initialize as nil
	var endDate, projectName interface{} = nil, nil
	var feedbackScore, notes interface{} = nil, nil
	
	// Preserve existing values if not explicitly set in the update request
	
	// For end date: use current value if new value is zero, otherwise use new value
	if engagement.EndDate.IsZero() && !current.EndDate.IsZero() {
		endDate = current.EndDate
	} else if !engagement.EndDate.IsZero() {
		endDate = engagement.EndDate
	}
	
	// For project name: use current value if new value is empty, otherwise use new value
	if engagement.ProjectName == "" && current.ProjectName != "" {
		projectName = current.ProjectName
	} else if engagement.ProjectName != "" {
		projectName = engagement.ProjectName
	}
	
	// For feedback score: use current value if new value is zero, otherwise use new value
	if engagement.FeedbackScore == 0 && current.FeedbackScore != 0 {
		feedbackScore = current.FeedbackScore
	} else if engagement.FeedbackScore > 0 {
		feedbackScore = engagement.FeedbackScore
	}
	
	// For notes: use current value if new value is empty, otherwise use new value
	if engagement.Notes == "" && current.Notes != "" {
		notes = current.Notes
	} else if engagement.Notes != "" {
		notes = engagement.Notes
	}
	
	result, err := s.db.Exec(
		query,
		engagement.ExpertID, engagement.EngagementType, engagement.StartDate,
		endDate, projectName, engagement.Status,
		feedbackScore, notes, engagement.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update engagement: %w", err)
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

// DeleteEngagement deletes an engagement by ID
func (s *SQLiteStore) DeleteEngagement(id int64) error {
	result, err := s.db.Exec("DELETE FROM expert_engagements WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete engagement: %w", err)
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