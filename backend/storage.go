package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	_ "github.com/mattn/go-sqlite3"
)

// Storage defines the interface for database operations
type Storage interface {
	// Expert methods
	ListExperts(filters map[string]interface{}, limit, offset int) ([]*Expert, error)
	CountExperts(filters map[string]interface{}) (int, error)
	GetExpert(id int64) (*Expert, error)
	CreateExpert(expert *Expert) (int64, error)
	UpdateExpert(expert *Expert) error
	DeleteExpert(id int64) error
	
	// Expert request methods
	CreateExpertRequest(request *ExpertRequest) (int64, error)
	ListExpertRequests(filters map[string]interface{}, limit, offset int) ([]*ExpertRequest, error)
	GetExpertRequest(id int64) (*ExpertRequest, error)
	UpdateExpertRequest(request *ExpertRequest) error
	
	// Document methods
	CreateDocument(doc *Document) (int64, error)
	GetDocument(id int64) (*Document, error)
	DeleteDocument(id int64) error
	GetDocumentsByExpertID(expertID int64) ([]*Document, error)
	
	// Engagement methods
	CreateEngagement(engagement *Engagement) (int64, error)
	GetEngagement(id int64) (*Engagement, error)
	UpdateEngagement(engagement *Engagement) error
	DeleteEngagement(id int64) error
	GetEngagementsByExpertID(expertID int64) ([]*Engagement, error)
	
	// Statistics methods
	GetStatistics() (*Statistics, error)
	GetExpertsByNationality() (int, int, error) // Returns (bahrainiCount, nonBahrainiCount, error)
	GetExpertsByISCEDField() ([]AreaStat, error)
	GetEngagementStatistics() ([]AreaStat, error)
	GetExpertGrowthByMonth(months int) ([]GrowthStat, error)
	
	// ISCED methods
	GetISCEDLevels() ([]ISCEDLevel, error)
	GetISCEDFields() ([]ISCEDField, error)
	
	// AI analysis methods
	StoreAIAnalysisResult(analysis *AIAnalysisResult) error
	SuggestISCED(expertID int64, input string) (*AIAnalysisResult, error)
	ExtractSkills(expertID int64, input string) (*AIAnalysisResult, error)
	GenerateProfile(expertID int64) (*AIAnalysisResult, error)
	SuggestExpertPanel(request string, count int) ([]Expert, error)
	
	// User methods
	CreateUser(user *User) error
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	ListUsers(limit, offset int) ([]*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
	CountUsers() (int, error)
	EnsureAdminExists(adminEmail, adminName, adminPasswordHash string) error
	
	// Utility methods
	Close() error
}

// SQLiteStore implements the Storage interface with SQLite backend
type SQLiteStore struct {
	db *sql.DB
}

// Verify that SQLiteStore implements the Storage interface at compile time
var _ Storage = (*SQLiteStore)(nil)

// NewSQLiteStore creates a new SQLite database connection
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}
	
	// Connect to the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	// Create the store
	store := &SQLiteStore{
		db: db,
	}
	
	// Initialize the database schema
	if err := store.initSchema(); err != nil {
		return nil, err
	}
	
	return store, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// initSchema initializes the database schema
func (s *SQLiteStore) initSchema() error {
	logger := GetLogger()
	logger.Info("Initializing database schema")
	
	// Check if we're using an in-memory database
	isMemoryDB := false
	
	// Try to determine if it's a memory database
	var name, file string
	var seq int
	err := s.db.QueryRow("PRAGMA database_list").Scan(&seq, &name, &file)
	if err != nil {
		logger.Warn("Failed to determine database type: %v", err)
		// Continue anyway, assume it's a file database
	} else {
		isMemoryDB = file == "" || file == ":memory:"
		logger.Info("Database type: %s (memory=%v)", file, isMemoryDB)
	}
	
	// For in-memory database, use a simplified schema
	if isMemoryDB {
		logger.Info("Using in-memory database with simplified schema")
		
		// Begin transaction
		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		
		// Define basic schema for testing
		schema := []string{
			// Create experts table
			`CREATE TABLE IF NOT EXISTS "experts" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				expert_id TEXT UNIQUE,
				name TEXT NOT NULL,
				designation TEXT,
				institution TEXT,
				is_bahraini BOOLEAN,
				nationality TEXT,
				is_available BOOLEAN,
				rating TEXT,
				role TEXT,
				employment_type TEXT,
				general_area TEXT,
				specialized_area TEXT,
				is_trained BOOLEAN,
				cv_path TEXT,
				phone TEXT,
				email TEXT,
				is_published BOOLEAN,
				isced_level_id INTEGER,
				isced_field_id INTEGER,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP
			)`,
			
			// Create expert_requests table
			`CREATE TABLE IF NOT EXISTS "expert_requests" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				expert_id TEXT,
				name TEXT NOT NULL,
				designation TEXT,
				institution TEXT,
				is_bahraini BOOLEAN,
				is_available BOOLEAN,
				rating TEXT,
				role TEXT,
				employment_type TEXT,
				general_area TEXT,
				specialized_area TEXT,
				is_trained BOOLEAN,
				cv_path TEXT,
				phone TEXT,
				email TEXT,
				is_published BOOLEAN,
				status TEXT DEFAULT 'pending',
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				reviewed_at TIMESTAMP,
				reviewed_by INTEGER
			)`,
			
			// Create users table
			`CREATE TABLE IF NOT EXISTS "users" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				email TEXT UNIQUE NOT NULL,
				password_hash TEXT NOT NULL,
				role TEXT NOT NULL,
				is_active BOOLEAN NOT NULL DEFAULT 1,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				last_login TIMESTAMP
			)`,
			
			// Create expert_documents table
			`CREATE TABLE IF NOT EXISTS "expert_documents" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				expert_id INTEGER NOT NULL,
				document_type TEXT NOT NULL,
				filename TEXT NOT NULL,
				file_path TEXT NOT NULL,
				content_type TEXT,
				file_size INTEGER,
				upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
			
			// Create expert_engagements table
			`CREATE TABLE IF NOT EXISTS "expert_engagements" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				expert_id INTEGER NOT NULL,
				engagement_type TEXT NOT NULL,
				start_date TIMESTAMP NOT NULL,
				end_date TIMESTAMP,
				project_name TEXT,
				status TEXT DEFAULT 'pending',
				feedback_score INTEGER,
				notes TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`,
			
			// Create ai_analysis table
			`CREATE TABLE IF NOT EXISTS "ai_analysis" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				expert_id INTEGER,
				document_id INTEGER,
				analysis_type TEXT NOT NULL,
				analysis_result TEXT NOT NULL,
				confidence_score REAL,
				model_used TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP
			)`,
			
			// Create isced_levels table
			`CREATE TABLE IF NOT EXISTS "isced_levels" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				code TEXT UNIQUE NOT NULL,
				name TEXT NOT NULL,
				description TEXT
			)`,
			
			// Create isced_fields table
			`CREATE TABLE IF NOT EXISTS "isced_fields" (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				broad_code TEXT NOT NULL,
				broad_name TEXT NOT NULL,
				narrow_code TEXT,
				narrow_name TEXT,
				detailed_code TEXT,
				detailed_name TEXT,
				description TEXT
			)`,
		}
		
		// Execute each statement
		for _, stmt := range schema {
			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute schema statement: %w", err)
			}
		}
		
		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		
	} else {
		// For file-based database, apply migrations from files
		logger.Info("Using file-based database with migrations")
		
		// Apply migration files from the db/migrations/sqlite directory
		migrations, err := filepath.Glob("db/migrations/sqlite/*.sql")
		if err != nil {
			return fmt.Errorf("failed to list migration files: %w", err)
		}
		
		// Sort migrations by filename (they should be prefixed with numbers)
		logger.Info("Found %d migration files", len(migrations))
		
		// Begin transaction
		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		
		// Apply migrations
		for _, migration := range migrations {
			logger.Info("Applying migration: %s", filepath.Base(migration))
			
			// Read the migration file
			content, err := os.ReadFile(migration)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to read migration file %s: %w", migration, err)
			}
			
			// Execute the migration
			if _, err := tx.Exec(string(content)); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %s: %w", migration, err)
			}
		}
		
		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}
	
	logger.Info("Database schema initialized successfully")
	return nil
}

// ISCED Methods

// GetISCEDLevels returns all ISCED levels from the database
func (s *SQLiteStore) GetISCEDLevels() ([]ISCEDLevel, error) {
	query := `SELECT id, code, name, description FROM isced_levels ORDER BY id`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var levels []ISCEDLevel
	for rows.Next() {
		var level ISCEDLevel
		if err := rows.Scan(&level.ID, &level.Code, &level.Name, &level.Description); err != nil {
			return nil, err
		}
		levels = append(levels, level)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return levels, nil
}

// GetISCEDFields returns all ISCED fields from the database
func (s *SQLiteStore) GetISCEDFields() ([]ISCEDField, error) {
	query := `
		SELECT id, broad_code, broad_name, narrow_code, narrow_name, 
		       detailed_code, detailed_name, description 
		FROM isced_fields 
		ORDER BY broad_code, narrow_code, detailed_code
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var fields []ISCEDField
	for rows.Next() {
		var field ISCEDField
		if err := rows.Scan(
			&field.ID, &field.BroadCode, &field.BroadName,
			&field.NarrowCode, &field.NarrowName,
			&field.DetailedCode, &field.DetailedName,
			&field.Description,
		); err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return fields, nil
}

// AI Analysis Methods

// StoreAIAnalysisResult stores an AI analysis result in the database
func (s *SQLiteStore) StoreAIAnalysisResult(analysis *AIAnalysisResult) error {
	query := `
		INSERT INTO ai_analysis (
			expert_id, document_id, analysis_type, analysis_result, confidence_score, 
			model_used, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		analysis.ExpertID, analysis.DocumentID, analysis.AnalysisType,
		analysis.AnalysisResult, analysis.ConfidenceScore, analysis.ModelUsed,
		analysis.CreatedAt, analysis.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	// Set the ID in the result
	analysis.ID = id
	
	return nil
}

// SuggestISCED suggests ISCED classification for expert data
func (s *SQLiteStore) SuggestISCED(expertID int64, input string) (*AIAnalysisResult, error) {
	// This would normally call the AI service, but for now we'll create a placeholder
	result := &AIAnalysisResult{
		ExpertID:        expertID,
		AnalysisType:    "isced_suggestion",
		AnalysisResult:  `{"iscedLevel":"7","iscedField":"0111"}`,
		ConfidenceScore: 0.85,
		ModelUsed:       "placeholder",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Store the result
	err := s.StoreAIAnalysisResult(result)
	return result, err
}

// ExtractSkills extracts skills from expert data
func (s *SQLiteStore) ExtractSkills(expertID int64, input string) (*AIAnalysisResult, error) {
	// This would normally call the AI service, but for now we'll create a placeholder
	result := &AIAnalysisResult{
		ExpertID:        expertID,
		AnalysisType:    "skills_extraction",
		AnalysisResult:  `["Project Management","Research Methods","Data Analysis"]`,
		ConfidenceScore: 0.9,
		ModelUsed:       "placeholder",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Store the result
	err := s.StoreAIAnalysisResult(result)
	return result, err
}

// GenerateProfile generates a profile for an expert
func (s *SQLiteStore) GenerateProfile(expertID int64) (*AIAnalysisResult, error) {
	// This would normally call the AI service, but for now we'll create a placeholder
	result := &AIAnalysisResult{
		ExpertID:        expertID,
		AnalysisType:    "profile_generation",
		AnalysisResult:  `{"summary":"Experienced professional with background in...","highlights":["Key achievement 1","Key achievement 2"]}`,
		ConfidenceScore: 0.8,
		ModelUsed:       "placeholder",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Store the result
	err := s.StoreAIAnalysisResult(result)
	return result, err
}

// SuggestExpertPanel suggests a panel of experts based on a request
func (s *SQLiteStore) SuggestExpertPanel(request string, count int) ([]Expert, error) {
	// For now, just return a sample of experts
	query := `
		SELECT id, expert_id, name, designation, institution, is_bahraini, 
		       is_available, rating, role, employment_type, general_area, 
		       specialized_area, is_trained, cv_path, phone, email, is_published, 
		       isced_level_id, isced_field_id, created_at, updated_at
		FROM experts
		WHERE is_available = 1 AND is_published = 1
		ORDER BY RANDOM()
		LIMIT ?
	`
	
	rows, err := s.db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var experts []Expert
	for rows.Next() {
		var expert Expert
		if err := rows.Scan(
			&expert.ID, &expert.ExpertID, &expert.Name, &expert.Designation,
			&expert.Institution, &expert.IsBahraini, &expert.IsAvailable,
			&expert.Rating, &expert.Role, &expert.EmploymentType,
			&expert.GeneralArea, &expert.SpecializedArea, &expert.IsTrained,
			&expert.CVPath, &expert.Phone, &expert.Email, &expert.IsPublished,
			&expert.ISCEDLevel, &expert.ISCEDField, &expert.CreatedAt, &expert.UpdatedAt,
		); err != nil {
			return nil, err
		}
		experts = append(experts, expert)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return experts, nil
}

// ListExperts retrieves experts based on filters with pagination and sorting
func (s *SQLiteStore) ListExperts(filters map[string]interface{}, limit, offset int) ([]*Expert, error) {
	// Default limit if not specified
	if limit <= 0 {
		limit = 10
	}

	// Build the query
	query := `
		SELECT e.id, e.expert_id, e.name, e.designation, e.institution, e.is_bahraini, 
		       e.is_available, e.rating, e.role, e.employment_type, e.general_area, 
		       e.specialized_area, e.is_trained, e.cv_path, e.phone, e.email, e.is_published, 
		       e.isced_level_id, e.isced_field_id, e.created_at, e.updated_at
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
				conditions = append(conditions, "e.general_area LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			case "is_available":
				conditions = append(conditions, "e.is_available = ?")
				args = append(args, value)
			case "role":
				conditions = append(conditions, "e.role LIKE ?")
				args = append(args, fmt.Sprintf("%%%s%%", value))
			case "isced_level_id":
				conditions = append(conditions, "e.isced_level_id = ?")
				args = append(args, value)
			case "isced_field_id":
				conditions = append(conditions, "e.isced_field_id = ?")
				args = append(args, value)
			case "min_rating":
				// Handle minimum rating filter (convert to numeric for comparison)
				conditions = append(conditions, "CAST(e.rating AS REAL) >= ?")
				args = append(args, value)
			case "sort_by":
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

	// Parse the results
	var experts []*Expert
	for rows.Next() {
		var expert Expert
		var createdAt, updatedAt string
		var iscedLevelID, iscedFieldID sql.NullInt64

		err := rows.Scan(
			&expert.ID, &expert.ExpertID, &expert.Name, &expert.Designation, &expert.Institution,
			&expert.IsBahraini, &expert.IsAvailable, &expert.Rating, &expert.Role, &expert.EmploymentType,
			&expert.GeneralArea, &expert.SpecializedArea, &expert.IsTrained, &expert.CVPath,
			&expert.Phone, &expert.Email, &expert.IsPublished, &iscedLevelID, &iscedFieldID,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expert row: %w", err)
		}

		// Parse timestamps
		expert.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		if updatedAt != "" {
			expert.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		}

		// Add ISCED data (basic info only in list view)
		if iscedLevelID.Valid {
			var level ISCEDLevel
			err := s.db.QueryRow(
				"SELECT id, code, name FROM isced_levels WHERE id = ?",
				iscedLevelID.Int64,
			).Scan(&level.ID, &level.Code, &level.Name)
			if err == nil {
				expert.ISCEDLevel = &level
			}
		}

		if iscedFieldID.Valid {
			var field ISCEDField
			err := s.db.QueryRow(
				"SELECT id, broad_code, broad_name FROM isced_fields WHERE id = ?",
				iscedFieldID.Int64,
			).Scan(&field.ID, &field.BroadCode, &field.BroadName)
			if err == nil {
				expert.ISCEDField = &field
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
			case "isced_level_id":
				conditions = append(conditions, "e.isced_level_id = ?")
				args = append(args, value)
			case "isced_field_id":
				conditions = append(conditions, "e.isced_field_id = ?")
				args = append(args, value)
			case "min_rating":
				conditions = append(conditions, "CAST(e.rating AS REAL) >= ?")
				args = append(args, value)
			}
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count experts: %w", err)
	}

	return count, nil
}

// GetStatistics retrieves system-wide statistics
func (s *SQLiteStore) GetStatistics() (*Statistics, error) {
    stats := &Statistics{
        LastUpdated: time.Now(),
    }
    
    // Get total experts count
    var totalExperts int
    err := s.db.QueryRow("SELECT COUNT(*) FROM experts").Scan(&totalExperts)
    if err != nil {
        return nil, fmt.Errorf("failed to count experts: %w", err)
    }
    stats.TotalExperts = totalExperts
    
    // Get Bahraini percentage
    bahrainiCount, _, err := s.GetExpertsByNationality()
    if err != nil {
        return nil, fmt.Errorf("failed to count experts by nationality: %w", err)
    }
    
    if totalExperts > 0 {
        stats.BahrainiPercentage = float64(bahrainiCount) / float64(totalExperts) * 100
    }
    
    // Get top areas
    rows, err := s.db.Query(`
        SELECT general_area, COUNT(*) as count
        FROM experts
        WHERE general_area != ''
        GROUP BY general_area
        ORDER BY count DESC
        LIMIT 10
    `)
    if err != nil {
        return nil, fmt.Errorf("failed to query top areas: %w", err)
    }
    defer rows.Close()
    
    var topAreas []AreaStat
    for rows.Next() {
        var area AreaStat
        var count int
        if err := rows.Scan(&area.Name, &count); err != nil {
            return nil, fmt.Errorf("failed to scan area row: %w", err)
        }
        area.Count = count
        if totalExperts > 0 {
            area.Percentage = float64(count) / float64(totalExperts) * 100
        }
        topAreas = append(topAreas, area)
    }
    stats.TopAreas = topAreas
    
    // Get ISCED field distribution
    iscedStats, err := s.GetExpertsByISCEDField()
    if err != nil {
        return nil, err
    }
    stats.ExpertsByISCEDField = iscedStats
    
    // Get engagement statistics
    engagementStats, err := s.GetEngagementStatistics()
    if err != nil {
        return nil, err
    }
    stats.EngagementsByType = engagementStats
    
    // Get monthly growth
    growthStats, err := s.GetExpertGrowthByMonth(12) // Last 12 months
    if err != nil {
        return nil, err
    }
    stats.MonthlyGrowth = growthStats
    
    // Get most requested experts
    rows, err = s.db.Query(`
        SELECT e.expert_id, e.name, COUNT(eng.id) as request_count
        FROM experts e
        JOIN expert_engagements eng ON e.id = eng.expert_id
        GROUP BY e.id
        ORDER BY request_count DESC
        LIMIT 10
    `)
    if err != nil {
        return nil, fmt.Errorf("failed to query most requested experts: %w", err)
    }
    defer rows.Close()
    
    var mostRequested []ExpertStat
    for rows.Next() {
        var stat ExpertStat
        if err := rows.Scan(&stat.ExpertID, &stat.Name, &stat.Count); err != nil {
            return nil, fmt.Errorf("failed to scan expert stat row: %w", err)
        }
        mostRequested = append(mostRequested, stat)
    }
    stats.MostRequestedExperts = mostRequested
    
    return stats, nil
}

// GetExpertsByNationality retrieves counts of experts by nationality (Bahraini vs non-Bahraini)
func (s *SQLiteStore) GetExpertsByNationality() (int, int, error) {
    var bahrainiCount, nonBahrainiCount int
    
    // Count Bahraini experts
    err := s.db.QueryRow("SELECT COUNT(*) FROM experts WHERE is_bahraini = 1").Scan(&bahrainiCount)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to count Bahraini experts: %w", err)
    }
    
    // Count non-Bahraini experts
    err = s.db.QueryRow("SELECT COUNT(*) FROM experts WHERE is_bahraini = 0").Scan(&nonBahrainiCount)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to count non-Bahraini experts: %w", err)
    }
    
    return bahrainiCount, nonBahrainiCount, nil
}

// GetExpertsByISCEDField retrieves counts of experts by ISCED field
func (s *SQLiteStore) GetExpertsByISCEDField() ([]AreaStat, error) {
    rows, err := s.db.Query(`
        SELECT f.broad_name, COUNT(e.id) as count
        FROM experts e
        JOIN isced_fields f ON e.isced_field_id = f.id
        GROUP BY f.broad_name
        ORDER BY count DESC
    `)
    if err != nil {
        return nil, fmt.Errorf("failed to query experts by ISCED field: %w", err)
    }
    defer rows.Close()
    
    var stats []AreaStat
    var totalInCategorizedFields int
    
    // First, collect all counts
    for rows.Next() {
        var stat AreaStat
        if err := rows.Scan(&stat.Name, &stat.Count); err != nil {
            return nil, fmt.Errorf("failed to scan ISCED field row: %w", err)
        }
        totalInCategorizedFields += stat.Count
        stats = append(stats, stat)
    }
    
    // Calculate percentages
    if totalInCategorizedFields > 0 {
        for i := range stats {
            stats[i].Percentage = float64(stats[i].Count) / float64(totalInCategorizedFields) * 100
        }
    }
    
    return stats, nil
}

// GetEngagementStatistics retrieves statistics about expert engagements
func (s *SQLiteStore) GetEngagementStatistics() ([]AreaStat, error) {
    rows, err := s.db.Query(`
        SELECT engagement_type, COUNT(*) as count, 
               SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count
        FROM expert_engagements
        GROUP BY engagement_type
        ORDER BY count DESC
    `)
    if err != nil {
        return nil, fmt.Errorf("failed to query engagement statistics: %w", err)
    }
    defer rows.Close()
    
    var stats []AreaStat
    var totalEngagements int
    
    // First, collect all counts
    for rows.Next() {
        var stat AreaStat
        var completedCount int
        if err := rows.Scan(&stat.Name, &stat.Count, &completedCount); err != nil {
            return nil, fmt.Errorf("failed to scan engagement row: %w", err)
        }
        totalEngagements += stat.Count
        stats = append(stats, stat)
    }
    
    // Calculate percentages
    if totalEngagements > 0 {
        for i := range stats {
            stats[i].Percentage = float64(stats[i].Count) / float64(totalEngagements) * 100
        }
    }
    
    return stats, nil
}

// GetExpertGrowthByMonth retrieves statistics about expert growth by month
func (s *SQLiteStore) GetExpertGrowthByMonth(months int) ([]GrowthStat, error) {
    // Default to 12 months if not specified
    if months <= 0 {
        months = 12
    }
    
    // Query to get expert count by month
    rows, err := s.db.Query(`
        SELECT 
            strftime('%Y-%m', created_at) as month,
            COUNT(*) as count
        FROM experts
        WHERE created_at >= date('now', '-' || ? || ' months')
        GROUP BY month
        ORDER BY month
    `, months)
    if err != nil {
        return nil, fmt.Errorf("failed to query expert growth: %w", err)
    }
    defer rows.Close()
    
    var stats []GrowthStat
    var prevCount int
    
    // Process each month
    for rows.Next() {
        var stat GrowthStat
        var monthStr string
        var count int
        if err := rows.Scan(&monthStr, &count); err != nil {
            return nil, fmt.Errorf("failed to scan growth stats row: %w", err)
        }
        
        stat.Period = monthStr
        stat.Count = count
        
        // Calculate growth rate (except for first month)
        if len(stats) > 0 && prevCount > 0 {
            stat.GrowthRate = (float64(count) - float64(prevCount)) / float64(prevCount) * 100
        }
        
        prevCount = count
        stats = append(stats, stat)
    }
    
    // If no data for some months in the range, fill with zeroes for continuity
    if len(stats) < months {
        // Generate a complete list of months
        endDate := time.Now()
        startDate := endDate.AddDate(0, -months, 0)
        
        filledStats := make([]GrowthStat, 0, months)
        
        // Create a map of existing stats for lookup
        existingStats := make(map[string]GrowthStat)
        for _, stat := range stats {
            existingStats[stat.Period] = stat
        }
        
        // Fill in all months
        for m := 0; m < months; m++ {
            currDate := startDate.AddDate(0, m, 0)
            monthStr := fmt.Sprintf("%04d-%02d", currDate.Year(), currDate.Month())
            
            if stat, exists := existingStats[monthStr]; exists {
                filledStats = append(filledStats, stat)
            } else {
                // Add empty stat
                filledStats = append(filledStats, GrowthStat{
                    Period: monthStr,
                    Count:  0,
                })
            }
        }
        
        // Recalculate growth rates with filled data
        for i := 1; i < len(filledStats); i++ {
            prevCount := filledStats[i-1].Count
            if prevCount > 0 {
                filledStats[i].GrowthRate = (float64(filledStats[i].Count) - float64(prevCount)) / float64(prevCount) * 100
            }
        }
        
        stats = filledStats
    }
    
    return stats, nil
}

// Document Methods

// CreateDocument creates a new document record in the database
func (s *SQLiteStore) CreateDocument(doc *Document) (int64, error) {
    query := `
        INSERT INTO expert_documents (
            expert_id, document_type, filename, file_path,
            content_type, file_size, upload_date
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `
    
    result, err := s.db.Exec(
        query,
        doc.ExpertID, doc.DocumentType, doc.Filename, doc.FilePath,
        doc.ContentType, doc.FileSize, doc.UploadDate,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to create document: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get document ID: %w", err)
    }
    
    return id, nil
}

// GetDocument retrieves a document by ID
func (s *SQLiteStore) GetDocument(id int64) (*Document, error) {
    query := `
        SELECT id, expert_id, document_type, filename, file_path,
               content_type, file_size, upload_date
        FROM expert_documents
        WHERE id = ?
    `
    
    var doc Document
    err := s.db.QueryRow(query, id).Scan(
        &doc.ID, &doc.ExpertID, &doc.DocumentType, &doc.Filename,
        &doc.FilePath, &doc.ContentType, &doc.FileSize, &doc.UploadDate,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("failed to get document: %w", err)
    }
    
    return &doc, nil
}

// DeleteDocument deletes a document by ID
func (s *SQLiteStore) DeleteDocument(id int64) error {
    _, err := s.db.Exec("DELETE FROM expert_documents WHERE id = ?", id)
    if err != nil {
        return fmt.Errorf("failed to delete document: %w", err)
    }
    
    return nil
}

// GetDocumentsByExpertID retrieves all documents for an expert
func (s *SQLiteStore) GetDocumentsByExpertID(expertID int64) ([]*Document, error) {
    query := `
        SELECT id, expert_id, document_type, filename, file_path,
               content_type, file_size, upload_date
        FROM expert_documents
        WHERE expert_id = ?
    `
    
    rows, err := s.db.Query(query, expertID)
    if err != nil {
        return nil, fmt.Errorf("failed to get expert documents: %w", err)
    }
    defer rows.Close()
    
    var docs []*Document
    for rows.Next() {
        var doc Document
        err := rows.Scan(
            &doc.ID, &doc.ExpertID, &doc.DocumentType, &doc.Filename,
            &doc.FilePath, &doc.ContentType, &doc.FileSize, &doc.UploadDate,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan document row: %w", err)
        }
        
        docs = append(docs, &doc)
    }
    
    return docs, nil
}

// Engagement Methods

// CreateEngagement creates a new engagement record
func (s *SQLiteStore) CreateEngagement(engagement *Engagement) (int64, error) {
    query := `
        INSERT INTO expert_engagements (
            expert_id, engagement_type, start_date, end_date,
            project_name, status, feedback_score, notes, created_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    result, err := s.db.Exec(
        query,
        engagement.ExpertID, engagement.EngagementType, engagement.StartDate,
        engagement.EndDate, engagement.ProjectName, engagement.Status,
        engagement.FeedbackScore, engagement.Notes, engagement.CreatedAt,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to create engagement: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get engagement ID: %w", err)
    }
    
    return id, nil
}

// GetEngagement retrieves an engagement by ID
func (s *SQLiteStore) GetEngagement(id int64) (*Engagement, error) {
    query := `
        SELECT id, expert_id, engagement_type, start_date, end_date,
               project_name, status, feedback_score, notes, created_at
        FROM expert_engagements
        WHERE id = ?
    `
    
    var engagement Engagement
    err := s.db.QueryRow(query, id).Scan(
        &engagement.ID, &engagement.ExpertID, &engagement.EngagementType,
        &engagement.StartDate, &engagement.EndDate, &engagement.ProjectName,
        &engagement.Status, &engagement.FeedbackScore, &engagement.Notes,
        &engagement.CreatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("failed to get engagement: %w", err)
    }
    
    return &engagement, nil
}

// UpdateEngagement updates an existing engagement
func (s *SQLiteStore) UpdateEngagement(engagement *Engagement) error {
    query := `
        UPDATE expert_engagements
        SET expert_id = ?, engagement_type = ?, start_date = ?, end_date = ?,
            project_name = ?, status = ?, feedback_score = ?, notes = ?
        WHERE id = ?
    `
    
    _, err := s.db.Exec(
        query,
        engagement.ExpertID, engagement.EngagementType, engagement.StartDate,
        engagement.EndDate, engagement.ProjectName, engagement.Status,
        engagement.FeedbackScore, engagement.Notes, engagement.ID,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update engagement: %w", err)
    }
    
    return nil
}

// DeleteEngagement deletes an engagement by ID
func (s *SQLiteStore) DeleteEngagement(id int64) error {
    _, err := s.db.Exec("DELETE FROM expert_engagements WHERE id = ?", id)
    if err != nil {
        return fmt.Errorf("failed to delete engagement: %w", err)
    }
    
    return nil
}

// GetEngagementsByExpertID retrieves all engagements for an expert
func (s *SQLiteStore) GetEngagementsByExpertID(expertID int64) ([]*Engagement, error) {
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
    
    var engagements []*Engagement
    for rows.Next() {
        var engagement Engagement
        err := rows.Scan(
            &engagement.ID, &engagement.ExpertID, &engagement.EngagementType,
            &engagement.StartDate, &engagement.EndDate, &engagement.ProjectName,
            &engagement.Status, &engagement.FeedbackScore, &engagement.Notes,
            &engagement.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan engagement row: %w", err)
        }
        
        engagements = append(engagements, &engagement)
    }
    
    return engagements, nil
}