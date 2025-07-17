package sqlite

import (
	"fmt"
	"strings"
)

// IsUserPlannerForApplication checks if a user has planner privileges for a specific application
func (s *SQLiteStore) IsUserPlannerForApplication(userID int, applicationID int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM application_planners 
		WHERE user_id = ? AND application_id = ?
	`
	
	var count int
	err := s.db.QueryRow(query, userID, applicationID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check planner permissions: %w", err)
	}
	
	return count > 0, nil
}

// IsUserManagerForApplication checks if a user has manager privileges for a specific application
func (s *SQLiteStore) IsUserManagerForApplication(userID int, applicationID int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM application_managers 
		WHERE user_id = ? AND application_id = ?
	`
	
	var count int
	err := s.db.QueryRow(query, userID, applicationID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check manager permissions: %w", err)
	}
	
	return count > 0, nil
}

// AssignUserToPlannerApplications assigns a user as planner to multiple applications in batch
func (s *SQLiteStore) AssignUserToPlannerApplications(userID int, applicationIDs []int) error {
	if len(applicationIDs) == 0 {
		return nil
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Build placeholders for the IN clause
	placeholders := strings.Repeat("?,", len(applicationIDs))
	placeholders = placeholders[:len(placeholders)-1] // Remove last comma
	
	// First, remove existing assignments for these applications for this user
	deleteQuery := fmt.Sprintf(`
		DELETE FROM application_planners 
		WHERE user_id = ? AND application_id IN (%s)
	`, placeholders)
	
	deleteArgs := make([]interface{}, len(applicationIDs)+1)
	deleteArgs[0] = userID
	for i, appID := range applicationIDs {
		deleteArgs[i+1] = appID
	}
	
	_, err = tx.Exec(deleteQuery, deleteArgs...)
	if err != nil {
		return fmt.Errorf("failed to remove existing planner assignments: %w", err)
	}
	
	// Insert new assignments
	insertQuery := fmt.Sprintf(`
		INSERT INTO application_planners (user_id, application_id) 
		VALUES %s
	`, strings.Repeat("(?,?),", len(applicationIDs))[:len(applicationIDs)*5-1])
	
	insertArgs := make([]interface{}, len(applicationIDs)*2)
	for i, appID := range applicationIDs {
		insertArgs[i*2] = userID
		insertArgs[i*2+1] = appID
	}
	
	_, err = tx.Exec(insertQuery, insertArgs...)
	if err != nil {
		return fmt.Errorf("failed to insert planner assignments: %w", err)
	}
	
	return tx.Commit()
}

// AssignUserToManagerApplications assigns a user as manager to multiple applications in batch
func (s *SQLiteStore) AssignUserToManagerApplications(userID int, applicationIDs []int) error {
	if len(applicationIDs) == 0 {
		return nil
	}
	
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Build placeholders for the IN clause
	placeholders := strings.Repeat("?,", len(applicationIDs))
	placeholders = placeholders[:len(placeholders)-1] // Remove last comma
	
	// First, remove existing assignments for these applications for this user
	deleteQuery := fmt.Sprintf(`
		DELETE FROM application_managers 
		WHERE user_id = ? AND application_id IN (%s)
	`, placeholders)
	
	deleteArgs := make([]interface{}, len(applicationIDs)+1)
	deleteArgs[0] = userID
	for i, appID := range applicationIDs {
		deleteArgs[i+1] = appID
	}
	
	_, err = tx.Exec(deleteQuery, deleteArgs...)
	if err != nil {
		return fmt.Errorf("failed to remove existing manager assignments: %w", err)
	}
	
	// Insert new assignments
	insertQuery := fmt.Sprintf(`
		INSERT INTO application_managers (user_id, application_id) 
		VALUES %s
	`, strings.Repeat("(?,?),", len(applicationIDs))[:len(applicationIDs)*5-1])
	
	insertArgs := make([]interface{}, len(applicationIDs)*2)
	for i, appID := range applicationIDs {
		insertArgs[i*2] = userID
		insertArgs[i*2+1] = appID
	}
	
	_, err = tx.Exec(insertQuery, insertArgs...)
	if err != nil {
		return fmt.Errorf("failed to insert manager assignments: %w", err)
	}
	
	return tx.Commit()
}

// RemoveUserPlannerAssignments removes planner assignments for a user from specific applications
func (s *SQLiteStore) RemoveUserPlannerAssignments(userID int, applicationIDs []int) error {
	if len(applicationIDs) == 0 {
		return nil
	}
	
	placeholders := strings.Repeat("?,", len(applicationIDs))
	placeholders = placeholders[:len(placeholders)-1] // Remove last comma
	
	query := fmt.Sprintf(`
		DELETE FROM application_planners 
		WHERE user_id = ? AND application_id IN (%s)
	`, placeholders)
	
	args := make([]interface{}, len(applicationIDs)+1)
	args[0] = userID
	for i, appID := range applicationIDs {
		args[i+1] = appID
	}
	
	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove planner assignments: %w", err)
	}
	
	return nil
}

// RemoveUserManagerAssignments removes manager assignments for a user from specific applications
func (s *SQLiteStore) RemoveUserManagerAssignments(userID int, applicationIDs []int) error {
	if len(applicationIDs) == 0 {
		return nil
	}
	
	placeholders := strings.Repeat("?,", len(applicationIDs))
	placeholders = placeholders[:len(placeholders)-1] // Remove last comma
	
	query := fmt.Sprintf(`
		DELETE FROM application_managers 
		WHERE user_id = ? AND application_id IN (%s)
	`, placeholders)
	
	args := make([]interface{}, len(applicationIDs)+1)
	args[0] = userID
	for i, appID := range applicationIDs {
		args[i+1] = appID
	}
	
	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove manager assignments: %w", err)
	}
	
	return nil
}

// GetUserPlannerApplications returns all application IDs that a user has planner access to
func (s *SQLiteStore) GetUserPlannerApplications(userID int) ([]int, error) {
	query := `
		SELECT application_id 
		FROM application_planners 
		WHERE user_id = ?
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user planner applications: %w", err)
	}
	defer rows.Close()
	
	var applicationIDs []int
	for rows.Next() {
		var appID int
		if err := rows.Scan(&appID); err != nil {
			return nil, fmt.Errorf("failed to scan application ID: %w", err)
		}
		applicationIDs = append(applicationIDs, appID)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	
	return applicationIDs, nil
}

// GetUserManagerApplications returns all application IDs that a user has manager access to
func (s *SQLiteStore) GetUserManagerApplications(userID int) ([]int, error) {
	query := `
		SELECT application_id 
		FROM application_managers 
		WHERE user_id = ?
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user manager applications: %w", err)
	}
	defer rows.Close()
	
	var applicationIDs []int
	for rows.Next() {
		var appID int
		if err := rows.Scan(&appID); err != nil {
			return nil, fmt.Errorf("failed to scan application ID: %w", err)
		}
		applicationIDs = append(applicationIDs, appID)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	
	return applicationIDs, nil
}