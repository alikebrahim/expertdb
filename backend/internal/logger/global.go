package logger

import "log"

// Global logger instance
var globalLogger *Logger

// Init initializes the global logger
func Init(level LogLevel, logDir string, useColors bool) error {
	logger, err := New(level, logDir, useColors)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// Get returns the global logger instance
func Get() *Logger {
	if globalLogger == nil {
		// Create a default logger if not initialized
		logger, err := New(LevelInfo, "./logs", true)
		if err != nil {
			log.Fatalf("Failed to create default logger: %v", err)
		}
		globalLogger = logger
	}
	return globalLogger
}