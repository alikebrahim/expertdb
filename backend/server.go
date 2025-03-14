package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// loadConfig loads configuration from environment variables
func loadConfig() *Configuration {
	config := &Configuration{
		Port: os.Getenv("PORT"),
		DBPath: os.Getenv("DB_PATH"),
		UploadPath: os.Getenv("UPLOAD_PATH"),
		CORSAllowOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
		AIServiceURL: os.Getenv("AI_SERVICE_URL"),
	}

	// Set defaults for empty values
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.DBPath == "" {
		config.DBPath = "./db/sqlite/expertdb.sqlite"
	}
	if config.UploadPath == "" {
		config.UploadPath = "./data/documents"
	}
	if config.CORSAllowOrigins == "" {
		config.CORSAllowOrigins = "*"
	}
	if config.AIServiceURL == "" {
		config.AIServiceURL = "http://localhost:9000"
	}

	return config
}

// verifyDatabaseSchema checks if the database schema has been properly initialized
func verifyDatabaseSchema(store Storage) error {
	logger := GetLogger()
	logger.Info("Verifying database schema...")
	
	// Get the SQLite store to access the db directly
	sqliteStore, ok := store.(*SQLiteStore)
	if !ok {
		return fmt.Errorf("store is not a SQLiteStore")
	}
	
	// List of required tables that must exist for the application to function
	requiredTables := []string{
		"experts",
		"expert_requests",
		"users",
		"expert_documents",
		"expert_engagements",
		"ai_analysis",
		"isced_levels",
		"isced_fields",
	}
	
	// Query to get list of all tables
	rows, err := sqliteStore.db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()
	
	// Build a map of existing tables
	existingTables := make(map[string]bool)
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		existingTables[tableName] = true
	}
	
	// Check if all required tables exist
	for _, table := range requiredTables {
		if !existingTables[table] {
			return fmt.Errorf("required table '%s' not found - database schema is incomplete", table)
		}
	}
	
	logger.Info("Database schema verification completed successfully")
	return nil
}

func main() {
	// Initialize logging system
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "./logs"
	}
	
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel := LevelInfo // Default level
	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		logLevel = LevelDebug
	case "INFO":
		logLevel = LevelInfo
	case "WARN":
		logLevel = LevelWarn
	case "ERROR":
		logLevel = LevelError
	}
	
	if err := InitLogger(logLevel, logDir, true); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	
	logger := GetLogger()
	logger.Info("Starting ExpertDB initialization...")
	
	// Load configuration
	config := loadConfig()
	logger.Info("Configuration loaded successfully")

	// Create the DB directory if it doesn't exist
	dbDir := filepath.Dir(config.DBPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		logger.Fatal("Failed to create database directory: %v", err)
	}
	logger.Info("Database directory created: %s", dbDir)

	// Create the upload directory if it doesn't exist
	if err := os.MkdirAll(config.UploadPath, 0755); err != nil {
		logger.Fatal("Failed to create upload directory: %v", err)
	}
	logger.Info("Upload directory created: %s", config.UploadPath)

	// Initialize database connection
	logger.Info("Connecting to database at %s", config.DBPath)
	
	// For testing, create an in-memory database
	dbPath := ":memory:"
	if os.Getenv("USE_FILE_DB") == "true" {
		dbPath = config.DBPath
	}
	
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer store.Close()
	logger.Info("Database connection established successfully")
	
	// For file-based databases, verify that the schema has been properly initialized
	if dbPath != ":memory:" {
		if err := verifyDatabaseSchema(store); err != nil {
			logger.Fatal("Database schema verification failed: %v", err)
			logger.Fatal("Please run migrations using: goose -dir db/migrations/sqlite sqlite3 %s up", dbPath)
		}
	}

	// Initialize JWT secret
	logger.Info("Initializing JWT secret...")
	if err := InitJWTSecret(); err != nil {
		logger.Fatal("Failed to initialize JWT secret: %v", err)
	}
	logger.Info("JWT secret initialized successfully")
	
	// Check if admin user exists, create if not
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminName := os.Getenv("ADMIN_NAME")
	
	// Set defaults if not provided
	if adminEmail == "" {
		adminEmail = "admin@expertdb.com"
	}
	if adminPassword == "" {
		adminPassword = "adminpassword"
	}
	if adminName == "" {
		adminName = "Admin User"
	}
	
	// Create admin user if it doesn't exist
	logger.Info("Checking for admin user with email: %s", adminEmail)
	_, err = store.GetUserByEmail(adminEmail)
	if err != nil {
		logger.Info("Admin user not found, creating...")
		
		// Create admin user
		passwordHash, err := GeneratePasswordHash(adminPassword)
		if err != nil {
			logger.Fatal("Failed to hash admin password: %v", err)
		}
		
		admin := &User{
			Name:         adminName,
			Email:        adminEmail,
			PasswordHash: passwordHash,
			Role:         RoleAdmin,
			IsActive:     true,
			CreatedAt:    time.Now(),
			LastLogin:    time.Now(),
		}
		
		if err := store.CreateUser(admin); err != nil {
			logger.Fatal("Failed to create admin user: %v", err)
		}
		
		logger.Info("Created default admin user with email: %s", adminEmail)
	} else {
		logger.Info("Admin user already exists, skipping creation")
	}

	// Create and start API server
	logger.Info("Creating API server on port %s", config.Port)
	server, err := NewAPIServer(":"+config.Port, store, config)
	if err != nil {
		logger.Fatal("Failed to create API server: %v", err)
	}
	
	logger.Info("Starting ExpertDB with configuration:")
	logger.Info("- Port: %s", config.Port)
	logger.Info("- Database: %s", config.DBPath)
	logger.Info("- Upload Path: %s", config.UploadPath)
	logger.Info("- CORS: %s", config.CORSAllowOrigins)
	logger.Info("- AI Service: %s", config.AIServiceURL)
	logger.Info("- Log Level: %s", logLevel.String())
	logger.Info("- Log Directory: %s", logDir)
	
	logger.Info("Server starting, press Ctrl+C to stop")
	if err := server.Run(); err != nil {
		logger.Fatal("Server error: %v", err)
	}
}