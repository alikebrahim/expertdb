// Package config provides configuration management for the ExpertDB application
package config

import "os"

// Configuration represents application configuration
type Configuration struct {
	Port             string `json:"port"`             // HTTP server port
	DBPath           string `json:"dbPath"`           // Path to SQLite database file
	UploadPath       string `json:"uploadPath"`       // Directory for uploaded documents
	CORSAllowOrigins string `json:"corsAllowOrigins"` // CORS allowed origins (comma-separated)
	AdminEmail       string `json:"-"`                // Default admin email
	AdminName        string `json:"-"`                // Default admin name
	AdminPassword    string `json:"-"`                // Default admin password
	LogDir           string `json:"-"`                // Directory for log files
	LogLevel         string `json:"-"`                // Log level (debug, info, warn, error)
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Configuration {
	config := &Configuration{
		Port:             os.Getenv("PORT"),
		DBPath:           os.Getenv("DB_PATH"),
		UploadPath:       os.Getenv("UPLOAD_PATH"),
		CORSAllowOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
		AdminEmail:       os.Getenv("ADMIN_EMAIL"),
		AdminName:        os.Getenv("ADMIN_NAME"),
		AdminPassword:    os.Getenv("ADMIN_PASSWORD"),
		LogDir:           os.Getenv("LOG_DIR"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
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
	if config.AdminEmail == "" {
		config.AdminEmail = "admin@expertdb.com"
	}
	if config.AdminName == "" {
		config.AdminName = "Admin User"
	}
	if config.AdminPassword == "" {
		config.AdminPassword = "adminpassword"
	}
	if config.LogDir == "" {
		config.LogDir = "./logs"
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return config
}