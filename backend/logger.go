package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// Debug level for verbose development logs
	LevelDebug LogLevel = iota
	// Info level for general information
	LevelInfo
	// Warn level for non-critical issues
	LevelWarn
	// Error level for errors that affect functionality
	LevelError
	// Fatal level for critical errors that require shutdown
	LevelFatal
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color returns ANSI color code for a log level
func (l LogLevel) Color() string {
	switch l {
	case LevelDebug:
		return "\033[37m" // White
	case LevelInfo:
		return "\033[32m" // Green
	case LevelWarn:
		return "\033[33m" // Yellow
	case LevelError:
		return "\033[31m" // Red
	case LevelFatal:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// Logger handles logging with different levels and formats
type Logger struct {
	level     LogLevel
	writer    io.Writer
	fileLog   *log.Logger
	consoleLog *log.Logger
	useColors bool
}

// NewLogger creates a new logger with the specified configuration
func NewLogger(level LogLevel, logDir string, useColors bool) (*Logger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file with timestamp in name
	timestamp := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("expertdb_%s.log", timestamp))
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer to log to both file and console
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	return &Logger{
		level:     level,
		writer:    multiWriter,
		fileLog:   log.New(logFile, "", 0),
		consoleLog: log.New(os.Stdout, "", 0),
		useColors: useColors,
	}, nil
}

// formatMessage formats a log message with timestamp, level, and caller info
func (l *Logger) formatMessage(level LogLevel, message string) string {
	// Get caller information
	_, file, line, ok := runtime.Caller(3) // Skip through the logger methods
	if !ok {
		file = "unknown"
		line = 0
	}
	// Get the short file name
	shortFile := filepath.Base(file)

	// Format the log message
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	return fmt.Sprintf("%s [%s] %s:%d: %s", 
		timestamp, 
		level.String(),
		shortFile,
		line,
		message)
}

// log logs a message at the specified level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	message := fmt.Sprintf(format, args...)
	formattedMsg := l.formatMessage(level, message)

	// Log to file without colors
	l.fileLog.Println(formattedMsg)

	// Log to console with colors if enabled
	if l.useColors {
		colorCode := level.Color()
		resetCode := "\033[0m"
		l.consoleLog.Printf("%s%s%s", colorCode, formattedMsg, resetCode)
	} else {
		l.consoleLog.Println(formattedMsg)
	}

	// If fatal, exit the program
	if level == LevelFatal {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LevelFatal, format, args...)
	// The program will exit in the log method
}

// HTTP request logging

// LogRequest logs HTTP request information
func (l *Logger) LogRequest(method, path, ip, userAgent string, statusCode int, duration time.Duration) {
	level := LevelInfo
	if statusCode >= 400 && statusCode < 500 {
		level = LevelWarn
	} else if statusCode >= 500 {
		level = LevelError
	}

	l.log(level, "HTTP %s %s from %s - %d (%s) - %v",
		method, path, ip, statusCode, http.StatusText(statusCode), duration)
}

// RequestLoggerMiddleware returns a middleware that logs HTTP requests
func (l *Logger) RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		
		// Create a response wrapper to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}
		
		// Process the request
		next.ServeHTTP(rw, r)
		
		// Calculate duration
		duration := time.Since(startTime)
		
		// Log the request
		l.LogRequest(
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.UserAgent(),
			rw.statusCode,
			duration,
		)
	})
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Global logger instance
var globalLogger *Logger

// InitLogger initializes the global logger
func InitLogger(level LogLevel, logDir string, useColors bool) error {
	logger, err := NewLogger(level, logDir, useColors)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		// Create a default logger if not initialized
		logger, err := NewLogger(LevelInfo, "./logs", true)
		if err != nil {
			log.Fatalf("Failed to create default logger: %v", err)
		}
		globalLogger = logger
	}
	return globalLogger
}