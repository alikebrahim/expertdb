package sqlite

import (
	"database/sql"
	"errors"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"fmt"
	"strings"
	"time"
)

// ListPhases retrieves a list of phases with optional filtering
func (s *SQLiteStore) ListPhases(status string, plannerID int64, limit, offset int) ([]*domain.Phase, error) {
	log := logger.Get()
	
	query := `SELECT p.id, p.phase_id, p.title, p.assigned_planner_id, u.name AS planner_name,
	          p.status, p.created_at, p.updated_at
	          FROM phases p
	          LEFT JOIN users u ON p.assigned_planner_id = u.id
	          WHERE 1=1`
	
	var args []interface{}
	
	// Apply status filter if provided
	if status != "" && status != "all" {
		query += " AND p.status = ?"
		args = append(args, status)
	}
	
	// Apply planner filter if provided
	if plannerID > 0 {
		query += " AND p.assigned_planner_id = ?"
		args = append(args, plannerID)
	}
	
	// Add order by and limits
	query += " ORDER BY p.created_at DESC"
	
	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}
	
	// Execute query
	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Error("Failed to query phases: %v", err)
		return nil, fmt.Errorf("failed to query phases: %w", err)
	}
	defer rows.Close()
	
	// Parse results
	phases := []*domain.Phase{}
	for rows.Next() {
		var phase domain.Phase
		var plannerName sql.NullString
		var createdAt, updatedAt string
		
		err := rows.Scan(
			&phase.ID,
			&phase.PhaseID,
			&phase.Title,
			&phase.AssignedPlannerID,
			&plannerName,
			&phase.Status,
			&createdAt,
			&updatedAt,
		)
		
		if err != nil {
			log.Error("Failed to scan phase row: %v", err)
			return nil, fmt.Errorf("failed to scan phase row: %w", err)
		}
		
		// Convert string timestamps to time.Time
		phase.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		phase.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		
		// Set planner name if available
		if plannerName.Valid {
			phase.PlannerName = plannerName.String
		}
		
		// Fetch applications for this phase
		applications, err := s.ListPhaseApplications(phase.ID)
		if err != nil {
			log.Error("Failed to fetch applications for phase %d: %v", phase.ID, err)
			return nil, fmt.Errorf("failed to fetch applications for phase: %w", err)
		}
		
		phase.Applications = applications
		phases = append(phases, &phase)
	}
	
	if err = rows.Err(); err != nil {
		log.Error("Error iterating phase rows: %v", err)
		return nil, fmt.Errorf("error iterating phase rows: %w", err)
	}
	
	return phases, nil
}

// GetPhase retrieves a phase by ID
func (s *SQLiteStore) GetPhase(id int64) (*domain.Phase, error) {
	log := logger.Get()
	
	query := `SELECT p.id, p.phase_id, p.title, p.assigned_planner_id, u.name AS planner_name,
	          p.status, p.created_at, p.updated_at
	          FROM phases p
	          LEFT JOIN users u ON p.assigned_planner_id = u.id
	          WHERE p.id = ?`
	
	var phase domain.Phase
	var plannerName sql.NullString
	var createdAt, updatedAt string
	
	err := s.db.QueryRow(query, id).Scan(
		&phase.ID,
		&phase.PhaseID,
		&phase.Title,
		&phase.AssignedPlannerID,
		&plannerName,
		&phase.Status,
		&createdAt,
		&updatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Error("Failed to get phase %d: %v", id, err)
		return nil, fmt.Errorf("failed to get phase: %w", err)
	}
	
	// Convert string timestamps to time.Time
	phase.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	phase.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	
	// Set planner name if available
	if plannerName.Valid {
		phase.PlannerName = plannerName.String
	}
	
	// Fetch applications for this phase
	applications, err := s.ListPhaseApplications(phase.ID)
	if err != nil {
		log.Error("Failed to fetch applications for phase %d: %v", phase.ID, err)
		return nil, fmt.Errorf("failed to fetch applications for phase: %w", err)
	}
	
	phase.Applications = applications
	
	return &phase, nil
}

// GetPhaseByPhaseID retrieves a phase by its business identifier
func (s *SQLiteStore) GetPhaseByPhaseID(phaseID string) (*domain.Phase, error) {
	log := logger.Get()
	
	query := `SELECT p.id, p.phase_id, p.title, p.assigned_planner_id, u.name AS planner_name,
	          p.status, p.created_at, p.updated_at
	          FROM phases p
	          LEFT JOIN users u ON p.assigned_planner_id = u.id
	          WHERE p.phase_id = ?`
	
	var phase domain.Phase
	var plannerName sql.NullString
	var createdAt, updatedAt string
	
	err := s.db.QueryRow(query, phaseID).Scan(
		&phase.ID,
		&phase.PhaseID,
		&phase.Title,
		&phase.AssignedPlannerID,
		&plannerName,
		&phase.Status,
		&createdAt,
		&updatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Error("Failed to get phase by ID %s: %v", phaseID, err)
		return nil, fmt.Errorf("failed to get phase by ID: %w", err)
	}
	
	// Convert string timestamps to time.Time
	phase.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	phase.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	
	// Set planner name if available
	if plannerName.Valid {
		phase.PlannerName = plannerName.String
	}
	
	// Fetch applications for this phase
	applications, err := s.ListPhaseApplications(phase.ID)
	if err != nil {
		log.Error("Failed to fetch applications for phase %d: %v", phase.ID, err)
		return nil, fmt.Errorf("failed to fetch applications for phase: %w", err)
	}
	
	phase.Applications = applications
	
	return &phase, nil
}

// CreatePhase creates a new phase with applications
func (s *SQLiteStore) CreatePhase(phase *domain.Phase) (int64, error) {
	log := logger.Get()
	
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Failed to begin transaction: %v", err)
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Error("Transaction rolled back: %v", err)
		}
	}()
	
	// Validate phase fields
	if strings.TrimSpace(phase.Title) == "" {
		return 0, fmt.Errorf("phase title cannot be empty")
	}
	
	if strings.TrimSpace(phase.Status) == "" {
		// Set default status
		phase.Status = "draft"
	}
	
	// Generate a unique phase ID if not provided
	if strings.TrimSpace(phase.PhaseID) == "" {
		phase.PhaseID, err = s.GenerateUniquePhaseID()
		if err != nil {
			return 0, fmt.Errorf("failed to generate phase ID: %w", err)
		}
	}
	
	// Set timestamps
	now := time.Now().UTC()
	phase.CreatedAt = now
	phase.UpdatedAt = now
	
	// Insert the phase
	result, err := tx.Exec(
		"INSERT INTO phases (phase_id, title, assigned_planner_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		phase.PhaseID,
		phase.Title,
		phase.AssignedPlannerID,
		phase.Status,
		phase.CreatedAt.Format(time.RFC3339),
		phase.UpdatedAt.Format(time.RFC3339),
	)
	
	if err != nil {
		log.Error("Failed to insert phase: %v", err)
		return 0, fmt.Errorf("failed to insert phase: %w", err)
	}
	
	// Get the phase ID
	phaseID, err := result.LastInsertId()
	if err != nil {
		log.Error("Failed to get phase ID: %v", err)
		return 0, fmt.Errorf("failed to get phase ID: %w", err)
	}
	
	phase.ID = phaseID
	
	// Insert phase applications if provided
	if len(phase.Applications) > 0 {
		for i := range phase.Applications {
			app := &phase.Applications[i]
			app.PhaseID = phaseID
			app.CreatedAt = now
			app.UpdatedAt = now
			
			// Set default status if not provided
			if strings.TrimSpace(app.Status) == "" {
				app.Status = "pending"
			}
			
			// Insert the application
			appResult, err := tx.Exec(
				`INSERT INTO phase_applications 
				(phase_id, type, institution_name, qualification_name, expert_1, expert_2, status, rejection_notes, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				app.PhaseID,
				app.Type,
				app.InstitutionName,
				app.QualificationName,
				nullableInt64(app.Expert1),
				nullableInt64(app.Expert2),
				app.Status,
				app.RejectionNotes,
				app.CreatedAt.Format(time.RFC3339),
				app.UpdatedAt.Format(time.RFC3339),
			)
			
			if err != nil {
				log.Error("Failed to insert phase application: %v", err)
				return 0, fmt.Errorf("failed to insert phase application: %w", err)
			}
			
			// Get the application ID
			appID, err := appResult.LastInsertId()
			if err != nil {
				log.Error("Failed to get application ID: %v", err)
				return 0, fmt.Errorf("failed to get application ID: %w", err)
			}
			
			app.ID = appID
		}
	}
	
	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Error("Failed to commit transaction: %v", err)
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	log.Info("Created new phase: %s (ID: %d) with %d applications", phase.PhaseID, phaseID, len(phase.Applications))
	return phaseID, nil
}

// UpdatePhase updates an existing phase
func (s *SQLiteStore) UpdatePhase(phase *domain.Phase) error {
	log := logger.Get()
	
	// Validate phase fields
	if strings.TrimSpace(phase.Title) == "" {
		return fmt.Errorf("phase title cannot be empty")
	}
	
	if strings.TrimSpace(phase.Status) == "" {
		return fmt.Errorf("phase status cannot be empty")
	}
	
	// Check if phase exists
	exists := false
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM phases WHERE id = ?)", phase.ID).Scan(&exists)
	if err != nil {
		log.Error("Failed to check if phase exists: %v", err)
		return fmt.Errorf("failed to check if phase exists: %w", err)
	}
	
	if !exists {
		return domain.ErrNotFound
	}
	
	// Update phase
	phase.UpdatedAt = time.Now().UTC()
	_, err = s.db.Exec(
		"UPDATE phases SET title = ?, assigned_planner_id = ?, status = ?, updated_at = ? WHERE id = ?",
		phase.Title,
		phase.AssignedPlannerID,
		phase.Status,
		phase.UpdatedAt.Format(time.RFC3339),
		phase.ID,
	)
	
	if err != nil {
		log.Error("Failed to update phase: %v", err)
		return fmt.Errorf("failed to update phase: %w", err)
	}
	
	log.Info("Updated phase: %s (ID: %d)", phase.PhaseID, phase.ID)
	return nil
}

// GenerateUniquePhaseID generates a unique business identifier for phases
func (s *SQLiteStore) GenerateUniquePhaseID() (string, error) {
	log := logger.Get()
	
	// Get the current year
	currentYear := time.Now().Year()
	
	// Count existing phases for this year
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM phases WHERE phase_id LIKE ?", fmt.Sprintf("PH-%d-%%", currentYear)).Scan(&count)
	if err != nil {
		log.Error("Failed to count phases: %v", err)
		return "", fmt.Errorf("failed to count phases: %w", err)
	}
	
	// Generate new ID
	newID := fmt.Sprintf("PH-%d-%03d", currentYear, count+1)
	
	// Verify it's unique
	exists := false
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM phases WHERE phase_id = ?)", newID).Scan(&exists)
	if err != nil {
		log.Error("Failed to check if phase ID exists: %v", err)
		return "", fmt.Errorf("failed to check if phase ID exists: %w", err)
	}
	
	if exists {
		// This shouldn't happen if our counting is correct, but just in case
		// Let's try with a higher number
		for i := count + 2; i < 1000; i++ {
			newID = fmt.Sprintf("PH-%d-%03d", currentYear, i)
			exists = false
			err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM phases WHERE phase_id = ?)", newID).Scan(&exists)
			if err != nil {
				log.Error("Failed to check if phase ID exists: %v", err)
				return "", fmt.Errorf("failed to check if phase ID exists: %w", err)
			}
			
			if !exists {
				break
			}
			
			// If we've tried 1000 times, something is wrong
			if i == 999 {
				log.Error("Failed to generate unique phase ID after 1000 attempts")
				return "", fmt.Errorf("failed to generate unique phase ID after 1000 attempts")
			}
		}
	}
	
	log.Debug("Generated unique phase ID: %s", newID)
	return newID, nil
}

// nullableInt64 returns a sql.NullInt64 based on the given value
func nullableInt64(val int64) sql.NullInt64 {
	if val <= 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Valid: true, Int64: val}
}

// ListPhaseApplications retrieves all applications for a phase
func (s *SQLiteStore) ListPhaseApplications(phaseID int64) ([]domain.PhaseApplication, error) {
	log := logger.Get()
	
	query := `SELECT 
	          a.id, a.phase_id, a.type, a.institution_name, a.qualification_name, 
	          a.expert_1, e1.name as expert_1_name, 
	          a.expert_2, e2.name as expert_2_name, 
	          a.status, a.rejection_notes, a.created_at, a.updated_at
	          FROM phase_applications a
	          LEFT JOIN experts e1 ON a.expert_1 = e1.id
	          LEFT JOIN experts e2 ON a.expert_2 = e2.id
	          WHERE a.phase_id = ?
	          ORDER BY a.id ASC`
	
	rows, err := s.db.Query(query, phaseID)
	if err != nil {
		log.Error("Failed to query phase applications: %v", err)
		return nil, fmt.Errorf("failed to query phase applications: %w", err)
	}
	defer rows.Close()
	
	applications := []domain.PhaseApplication{}
	for rows.Next() {
		var app domain.PhaseApplication
		var expert1Name, expert2Name sql.NullString
		var expert1, expert2 sql.NullInt64
		var createdAt, updatedAt string
		var rejectionNotes sql.NullString
		
		err := rows.Scan(
			&app.ID,
			&app.PhaseID,
			&app.Type,
			&app.InstitutionName,
			&app.QualificationName,
			&expert1,
			&expert1Name,
			&expert2,
			&expert2Name,
			&app.Status,
			&rejectionNotes,
			&createdAt,
			&updatedAt,
		)
		
		if err != nil {
			log.Error("Failed to scan application row: %v", err)
			return nil, fmt.Errorf("failed to scan application row: %w", err)
		}
		
		// Convert string timestamps to time.Time
		app.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		app.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		
		// Set nullable fields
		if expert1.Valid {
			app.Expert1 = expert1.Int64
		}
		
		if expert2.Valid {
			app.Expert2 = expert2.Int64
		}
		
		if expert1Name.Valid {
			app.Expert1Name = expert1Name.String
		}
		
		if expert2Name.Valid {
			app.Expert2Name = expert2Name.String
		}
		
		if rejectionNotes.Valid {
			app.RejectionNotes = rejectionNotes.String
		}
		
		applications = append(applications, app)
	}
	
	if err = rows.Err(); err != nil {
		log.Error("Error iterating application rows: %v", err)
		return nil, fmt.Errorf("error iterating application rows: %w", err)
	}
	
	return applications, nil
}

// GetPhaseApplication retrieves a single application by ID
func (s *SQLiteStore) GetPhaseApplication(id int64) (*domain.PhaseApplication, error) {
	log := logger.Get()
	
	query := `SELECT 
	          a.id, a.phase_id, a.type, a.institution_name, a.qualification_name, 
	          a.expert_1, e1.name as expert_1_name, 
	          a.expert_2, e2.name as expert_2_name, 
	          a.status, a.rejection_notes, a.created_at, a.updated_at
	          FROM phase_applications a
	          LEFT JOIN experts e1 ON a.expert_1 = e1.id
	          LEFT JOIN experts e2 ON a.expert_2 = e2.id
	          WHERE a.id = ?`
	
	var app domain.PhaseApplication
	var expert1Name, expert2Name sql.NullString
	var expert1, expert2 sql.NullInt64
	var createdAt, updatedAt string
	var rejectionNotes sql.NullString
	
	err := s.db.QueryRow(query, id).Scan(
		&app.ID,
		&app.PhaseID,
		&app.Type,
		&app.InstitutionName,
		&app.QualificationName,
		&expert1,
		&expert1Name,
		&expert2,
		&expert2Name,
		&app.Status,
		&rejectionNotes,
		&createdAt,
		&updatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		log.Error("Failed to get application %d: %v", id, err)
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	
	// Convert string timestamps to time.Time
	app.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	app.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	
	// Set nullable fields
	if expert1.Valid {
		app.Expert1 = expert1.Int64
	}
	
	if expert2.Valid {
		app.Expert2 = expert2.Int64
	}
	
	if expert1Name.Valid {
		app.Expert1Name = expert1Name.String
	}
	
	if expert2Name.Valid {
		app.Expert2Name = expert2Name.String
	}
	
	if rejectionNotes.Valid {
		app.RejectionNotes = rejectionNotes.String
	}
	
	return &app, nil
}

// CreatePhaseApplication creates a new phase application
func (s *SQLiteStore) CreatePhaseApplication(app *domain.PhaseApplication) (int64, error) {
	log := logger.Get()
	
	// Validate required fields
	if app.PhaseID <= 0 {
		return 0, fmt.Errorf("phase ID is required")
	}
	
	if strings.TrimSpace(app.Type) == "" {
		return 0, fmt.Errorf("application type is required")
	}
	
	if strings.TrimSpace(app.InstitutionName) == "" {
		return 0, fmt.Errorf("institution name is required")
	}
	
	if strings.TrimSpace(app.QualificationName) == "" {
		return 0, fmt.Errorf("qualification name is required")
	}
	
	// Verify phase exists
	exists := false
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM phases WHERE id = ?)", app.PhaseID).Scan(&exists)
	if err != nil {
		log.Error("Failed to check if phase exists: %v", err)
		return 0, fmt.Errorf("failed to check if phase exists: %w", err)
	}
	
	if !exists {
		return 0, fmt.Errorf("phase with ID %d does not exist", app.PhaseID)
	}
	
	// Set defaults for nullable fields
	if strings.TrimSpace(app.Status) == "" {
		app.Status = "pending"
	}
	
	// Set timestamps
	now := time.Now().UTC()
	app.CreatedAt = now
	app.UpdatedAt = now
	
	// Insert the application
	result, err := s.db.Exec(
		`INSERT INTO phase_applications 
		(phase_id, type, institution_name, qualification_name, expert_1, expert_2, status, rejection_notes, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		app.PhaseID,
		app.Type,
		app.InstitutionName,
		app.QualificationName,
		nullableInt64(app.Expert1),
		nullableInt64(app.Expert2),
		app.Status,
		app.RejectionNotes,
		app.CreatedAt.Format(time.RFC3339),
		app.UpdatedAt.Format(time.RFC3339),
	)
	
	if err != nil {
		log.Error("Failed to insert phase application: %v", err)
		return 0, fmt.Errorf("failed to insert phase application: %w", err)
	}
	
	// Get the application ID
	appID, err := result.LastInsertId()
	if err != nil {
		log.Error("Failed to get application ID: %v", err)
		return 0, fmt.Errorf("failed to get application ID: %w", err)
	}
	
	app.ID = appID
	
	log.Info("Created new phase application (ID: %d) for phase %d", appID, app.PhaseID)
	return appID, nil
}

// UpdatePhaseApplication updates an existing phase application
func (s *SQLiteStore) UpdatePhaseApplication(app *domain.PhaseApplication) error {
	log := logger.Get()
	
	// Validate required fields
	if app.ID <= 0 {
		return fmt.Errorf("application ID is required")
	}
	
	if strings.TrimSpace(app.Type) == "" {
		return fmt.Errorf("application type is required")
	}
	
	if strings.TrimSpace(app.InstitutionName) == "" {
		return fmt.Errorf("institution name is required")
	}
	
	if strings.TrimSpace(app.QualificationName) == "" {
		return fmt.Errorf("qualification name is required")
	}
	
	if strings.TrimSpace(app.Status) == "" {
		return fmt.Errorf("application status is required")
	}
	
	// Check if application exists
	exists := false
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM phase_applications WHERE id = ?)", app.ID).Scan(&exists)
	if err != nil {
		log.Error("Failed to check if application exists: %v", err)
		return fmt.Errorf("failed to check if application exists: %w", err)
	}
	
	if !exists {
		return domain.ErrNotFound
	}
	
	// Update the application
	app.UpdatedAt = time.Now().UTC()
	_, err = s.db.Exec(
		`UPDATE phase_applications 
		SET type = ?, institution_name = ?, qualification_name = ?, 
		expert_1 = ?, expert_2 = ?, status = ?, rejection_notes = ?, updated_at = ? 
		WHERE id = ?`,
		app.Type,
		app.InstitutionName,
		app.QualificationName,
		nullableInt64(app.Expert1),
		nullableInt64(app.Expert2),
		app.Status,
		app.RejectionNotes,
		app.UpdatedAt.Format(time.RFC3339),
		app.ID,
	)
	
	if err != nil {
		log.Error("Failed to update phase application: %v", err)
		return fmt.Errorf("failed to update phase application: %w", err)
	}
	
	log.Info("Updated phase application (ID: %d)", app.ID)
	return nil
}

// UpdatePhaseApplicationExperts updates the experts assigned to an application
func (s *SQLiteStore) UpdatePhaseApplicationExperts(id int64, expert1ID, expert2ID int64) error {
	log := logger.Get()
	
	// Check if application exists
	_, err := s.GetPhaseApplication(id)
	if err != nil {
		return err // Error already logged in GetPhaseApplication
	}
	
	// Verify experts exist (if specified)
	if expert1ID > 0 {
		exists := false
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM experts WHERE id = ?)", expert1ID).Scan(&exists)
		if err != nil {
			log.Error("Failed to check if expert1 exists: %v", err)
			return fmt.Errorf("failed to check if expert1 exists: %w", err)
		}
		
		if !exists {
			return fmt.Errorf("expert with ID %d does not exist", expert1ID)
		}
	}
	
	if expert2ID > 0 {
		exists := false
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM experts WHERE id = ?)", expert2ID).Scan(&exists)
		if err != nil {
			log.Error("Failed to check if expert2 exists: %v", err)
			return fmt.Errorf("failed to check if expert2 exists: %w", err)
		}
		
		if !exists {
			return fmt.Errorf("expert with ID %d does not exist", expert2ID)
		}
	}
	
	// Update the experts
	now := time.Now().UTC()
	_, err = s.db.Exec(
		"UPDATE phase_applications SET expert_1 = ?, expert_2 = ?, status = ?, updated_at = ? WHERE id = ?",
		nullableInt64(expert1ID),
		nullableInt64(expert2ID),
		"assigned", // Update status to assigned when experts are set
		now.Format(time.RFC3339),
		id,
	)
	
	if err != nil {
		log.Error("Failed to update application experts: %v", err)
		return fmt.Errorf("failed to update application experts: %w", err)
	}
	
	log.Info("Updated experts for application %d", id)
	return nil
}

// UpdatePhaseApplicationStatus updates the status of an application
func (s *SQLiteStore) UpdatePhaseApplicationStatus(id int64, status, rejectionNotes string) error {
	log := logger.Get()
	
	// Check if application exists
	app, err := s.GetPhaseApplication(id)
	if err != nil {
		return err // Error already logged in GetPhaseApplication
	}
	
	// Validate status
	if strings.TrimSpace(status) == "" {
		return fmt.Errorf("status cannot be empty")
	}
	
	validStatuses := []string{"pending", "assigned", "approved", "rejected"}
	if !containsString(validStatuses, status) {
		return fmt.Errorf("invalid status: %s", status)
	}
	
	// Check for rejection notes if rejecting
	if status == "rejected" && strings.TrimSpace(rejectionNotes) == "" {
		return fmt.Errorf("rejection notes are required when rejecting an application")
	}
	
	// Start a transaction for potential engagement creation
	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Error("Transaction rolled back: %v", err)
		}
	}()
	
	// Update the status
	now := time.Now().UTC()
	result, err := tx.Exec(
		"UPDATE phase_applications SET status = ?, rejection_notes = ?, updated_at = ? WHERE id = ?",
		status,
		rejectionNotes,
		now.Format(time.RFC3339),
		id,
	)
	
	if err != nil {
		log.Error("Failed to update application status: %v", err)
		return fmt.Errorf("failed to update application status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Failed to get rows affected: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	// If approved, create engagements for the assigned experts
	if status == "approved" {
		// Check if experts are assigned
		if app.Expert1 <= 0 && app.Expert2 <= 0 {
			return fmt.Errorf("cannot approve application without assigned experts")
		}
		
		// Determine engagement type based on application type
		// Mapping QP (Qualification Placement) to validator and IL (Institutional Listing) to evaluator
		var engagementType string
		if app.Type == "validation" || app.Type == "QP" {
			engagementType = "validator"  // QP (Qualification Placement) maps to validator
		} else if app.Type == "evaluation" || app.Type == "IL" {
			engagementType = "evaluator"  // IL (Institutional Listing) maps to evaluator
		} else {
			return fmt.Errorf("invalid application type: %s", app.Type)
		}
		
		// Create engagement for expert 1 if assigned
		if app.Expert1 > 0 {
			engagement := &domain.Engagement{
				ExpertID:       app.Expert1,
				EngagementType: engagementType,
				StartDate:      now,
				ProjectName:    fmt.Sprintf("%s - %s", app.InstitutionName, app.QualificationName),
				Status:         "active",
				Notes:          fmt.Sprintf("Automatically created from phase application ID %d", app.ID),
				CreatedAt:      now,
			}
			
			_, err = tx.Exec(
				`INSERT INTO expert_engagements 
				(expert_id, engagement_type, start_date, project_name, status, notes, created_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				engagement.ExpertID,
				engagement.EngagementType,
				engagement.StartDate.Format(time.RFC3339),
				engagement.ProjectName,
				engagement.Status,
				engagement.Notes,
				engagement.CreatedAt.Format(time.RFC3339),
			)
			
			if err != nil {
				log.Error("Failed to create engagement for expert 1: %v", err)
				return fmt.Errorf("failed to create engagement for expert 1: %w", err)
			}
			
			log.Info("Created engagement for expert %d", app.Expert1)
		}
		
		// Create engagement for expert 2 if assigned
		if app.Expert2 > 0 {
			engagement := &domain.Engagement{
				ExpertID:       app.Expert2,
				EngagementType: engagementType,
				StartDate:      now,
				ProjectName:    fmt.Sprintf("%s - %s", app.InstitutionName, app.QualificationName),
				Status:         "active",
				Notes:          fmt.Sprintf("Automatically created from phase application ID %d", app.ID),
				CreatedAt:      now,
			}
			
			_, err = tx.Exec(
				`INSERT INTO expert_engagements 
				(expert_id, engagement_type, start_date, project_name, status, notes, created_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				engagement.ExpertID,
				engagement.EngagementType,
				engagement.StartDate.Format(time.RFC3339),
				engagement.ProjectName,
				engagement.Status,
				engagement.Notes,
				engagement.CreatedAt.Format(time.RFC3339),
			)
			
			if err != nil {
				log.Error("Failed to create engagement for expert 2: %v", err)
				return fmt.Errorf("failed to create engagement for expert 2: %w", err)
			}
			
			log.Info("Created engagement for expert %d", app.Expert2)
		}
	}
	
	// Commit the transaction
	if err = tx.Commit(); err != nil {
		log.Error("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	log.Info("Updated status to %s for application %d", status, id)
	return nil
}

// Helper function to check if a string is in a slice
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}