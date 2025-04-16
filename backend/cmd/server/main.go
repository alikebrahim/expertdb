// Package main provides the entry point for the ExpertDB server
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"expertdb/internal/api"
	"expertdb/internal/auth"
	"expertdb/internal/config"
	"expertdb/internal/domain"
	"expertdb/internal/documents"
	"expertdb/internal/logger"
	"expertdb/internal/storage/sqlite"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Initialize logging system
	logLevelStr := cfg.LogLevel
	logLevel := logger.LevelInfo // Default level
	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		logLevel = logger.LevelDebug
	case "INFO":
		logLevel = logger.LevelInfo
	case "WARN":
		logLevel = logger.LevelWarn
	case "ERROR":
		logLevel = logger.LevelError
	}
	
	if err := logger.Init(logLevel, cfg.LogDir, true); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	
	// Get logger
	l := logger.Get()
	l.Info("Starting ExpertDB initialization...")
	
	// Create the DB directory if it doesn't exist
	dbDir := filepath.Dir(cfg.DBPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		l.Fatal("Failed to create database directory: %v", err)
	}
	l.Info("Database directory created: %s", dbDir)
	
	// Create the upload directory if it doesn't exist
	if err := os.MkdirAll(cfg.UploadPath, 0755); err != nil {
		l.Fatal("Failed to create upload directory: %v", err)
	}
	l.Info("Upload directory created: %s", cfg.UploadPath)
	
	// Initialize database connection
	l.Info("Connecting to database at %s", cfg.DBPath)
	
	// Create storage implementation
	store, err := sqlite.New(cfg.DBPath)
	if err != nil {
		l.Fatal("Failed to connect to database: %v", err)
	}
	defer store.Close()
	l.Info("Database connection established successfully")
	
	// Initialize database if needed
	if err := store.InitDB(); err != nil {
		l.Fatal("Failed to initialize database: %v", err)
	}
	
	// Initialize JWT secret
	l.Info("Initializing JWT secret...")
	if err := auth.InitJWTSecret(); err != nil {
		l.Fatal("Failed to initialize JWT secret: %v", err)
	}
	l.Info("JWT secret initialized successfully")
	
	// Create document service
	docService, err := documents.New(store, cfg.UploadPath)
	if err != nil {
		l.Fatal("Failed to create document service: %v", err)
	}
	
	// Create API server
	l.Info("Creating API server on port %s", cfg.Port)
	server, err := api.NewServer(":"+cfg.Port, store, docService, cfg)
	if err != nil {
		l.Fatal("Failed to create API server: %v", err)
	}
	
	// Create admin user if it doesn't exist
	l.Info("Checking for admin user with email: %s", cfg.AdminEmail)
	_, err = store.GetUserByEmail(cfg.AdminEmail)
	if err != nil {
		l.Info("Admin user not found, creating...")
		
		// Create admin user
		passwordHash, err := auth.GeneratePasswordHash(cfg.AdminPassword)
		if err != nil {
			l.Fatal("Failed to hash admin password: %v", err)
		}
		
		admin := &domain.User{
			Name:         cfg.AdminName,
			Email:        cfg.AdminEmail,
			PasswordHash: passwordHash,
			Role:         auth.RoleAdmin,
			IsActive:     true,
			CreatedAt:    time.Now(),
			LastLogin:    time.Now(),
		}
		
		if _, err := store.CreateUser(admin); err != nil {
			l.Fatal("Failed to create admin user: %v", err)
		}
		
		l.Info("Created default admin user with email: %s", cfg.AdminEmail)
	} else {
		l.Info("Admin user already exists, skipping creation")
	}
	
	l.Info("Starting ExpertDB with configuration:")
	l.Info("- Port: %s", cfg.Port)
	l.Info("- Database: %s", cfg.DBPath)
	l.Info("- Upload Path: %s", cfg.UploadPath)
	l.Info("- CORS: %s", cfg.CORSAllowOrigins)
	l.Info("- Log Level: %s", logLevel.String())
	l.Info("- Log Directory: %s", cfg.LogDir)
	
	// For mock data generation, run the populate_mock_data.sh script
	// This keeps the server code clean and focused on its primary responsibility
	
	l.Info("Server starting, press Ctrl+C to stop")
	if err := server.Run(); err != nil {
		l.Fatal("Server error: %v", err)
	}
}