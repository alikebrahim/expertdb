package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Custom error types for different categories of errors
var (
	ErrNotFound        = fmt.Errorf("resource not found")
	ErrInvalidInput    = fmt.Errorf("invalid input data")
	ErrConflict        = fmt.Errorf("resource conflict")
	ErrDatabaseError   = fmt.Errorf("database error")
	ErrForbidden       = fmt.Errorf("access forbidden")
	ErrUnauthorized    = fmt.Errorf("unauthorized access")
	ErrInternalError   = fmt.Errorf("internal server error")
)

// ValidationError represents an error with multiple validation issues
type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

// NewValidationError creates a new validation error
func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors: make(map[string]string),
	}
}

// Add adds a field validation error
func (v *ValidationError) Add(field, message string) {
	v.Errors[field] = message
}

// HasErrors checks if there are any validation errors
func (v *ValidationError) HasErrors() bool {
	return len(v.Errors) > 0
}

// Error implements the error interface
func (v *ValidationError) Error() string {
	if len(v.Errors) == 0 {
		return "no validation errors"
	}

	var errMsgs []string
	for field, msg := range v.Errors {
		errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", field, msg))
	}
	return strings.Join(errMsgs, "; ")
}

// JSONParsingError represents a JSON parsing error
type JSONParsingError struct {
	OriginalError error
}

// NewJSONParsingError creates a new JSON parsing error
func NewJSONParsingError(err error) *JSONParsingError {
	return &JSONParsingError{
		OriginalError: err,
	}
}

// Error implements the error interface
func (e *JSONParsingError) Error() string {
	return fmt.Sprintf("invalid JSON format: %v", e.OriginalError)
}

// Unwrap returns the original error
func (e *JSONParsingError) Unwrap() error {
	return e.OriginalError
}

// DatabaseError represents a database error with details about the specific issue
type DatabaseError struct {
	OriginalError error
	Code          string // Optional error code
	Field         string // Optional field related to the error
	Operation     string // The database operation (create, update, delete, etc.)
}

// NewDatabaseError creates a new database error
func NewDatabaseError(err error, operation string) *DatabaseError {
	dbErr := &DatabaseError{
		OriginalError: err,
		Operation:     operation,
	}

	// Parse common SQLite errors
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "UNIQUE constraint failed") {
			dbErr.Code = "unique_constraint"
			// Extract field name from error message
			parts := strings.Split(errMsg, "UNIQUE constraint failed: ")
			if len(parts) > 1 {
				fieldParts := strings.Split(parts[1], ".")
				if len(fieldParts) > 1 {
					dbErr.Field = fieldParts[1]
				}
			}
		} else if strings.Contains(errMsg, "FOREIGN KEY constraint failed") {
			dbErr.Code = "foreign_key"
		} else if strings.Contains(errMsg, "no such table") {
			dbErr.Code = "missing_table"
		} else if strings.Contains(errMsg, "NOT NULL constraint failed") {
			dbErr.Code = "required_field"
			// Extract field name from error message
			parts := strings.Split(errMsg, "NOT NULL constraint failed: ")
			if len(parts) > 1 {
				fieldParts := strings.Split(parts[1], ".")
				if len(fieldParts) > 1 {
					dbErr.Field = fieldParts[1]
				}
			}
		}
	}

	return dbErr
}

// ParseSQLiteError parses a SQLite error into a user-friendly message
// This function is designed to be used by all handlers that interact with the database
func ParseSQLiteError(err error, entityName string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	
	// Handle UNIQUE constraint violations
	if strings.Contains(errMsg, "UNIQUE constraint failed") {
		// Extract the field name from the error message
		field := extractConstraintField(errMsg, "UNIQUE constraint failed")
		fieldName := formatFieldName(field)
		
		// Check for specific field handling
		switch fieldName {
		case "email":
			return fmt.Errorf("email already exists: an %s with this email is already registered", entityName)
		default:
			return fmt.Errorf("unique constraint violation on %s: duplicate value not allowed", fieldName)
		}
	}
	
	// Handle FOREIGN KEY constraint violations
	if strings.Contains(errMsg, "FOREIGN KEY constraint failed") {
		// Special case for general_area (common foreign key)
		if strings.Contains(errMsg, "general_area") {
			return fmt.Errorf("invalid general area ID: this area does not exist in the system")
		}
		
		return fmt.Errorf("referenced resource does not exist: %v", err)
	}
	
	// Handle NOT NULL constraint violations
	if strings.Contains(errMsg, "NOT NULL constraint failed") {
		field := extractConstraintField(errMsg, "NOT NULL constraint failed")
		return fmt.Errorf("required field missing: %s cannot be empty", formatFieldName(field))
	}
	
	// Handle CHECK constraint violations
	if strings.Contains(errMsg, "CHECK constraint failed") {
		field := extractConstraintField(errMsg, "CHECK constraint failed")
		return fmt.Errorf("invalid value for %s: value does not meet requirements", formatFieldName(field))
	}
	
	// Default error handling for other database errors
	return fmt.Errorf("database error: %v", err)
}

// Helper function to extract the field name from a constraint error message
func extractConstraintField(errMsg string, constraintType string) string {
	parts := strings.Split(errMsg, constraintType+": ")
	if len(parts) < 2 {
		return "unknown field"
	}
	
	fieldPart := strings.TrimSpace(parts[1])
	fieldParts := strings.Split(fieldPart, ".")
	if len(fieldParts) < 2 {
		return fieldPart
	}
	
	return fieldParts[1]
}

// Helper function to format a field name for user-friendly messages
// Converts snake_case to space-separated words
func formatFieldName(field string) string {
	// Replace underscores with spaces
	formatted := strings.ReplaceAll(field, "_", " ")
	return formatted
}

// Error implements the error interface
func (e *DatabaseError) Error() string {
	var details string
	if e.Field != "" {
		details = fmt.Sprintf(" (field: %s)", e.Field)
	}

	switch e.Code {
	case "unique_constraint":
		return fmt.Sprintf("resource already exists%s", details)
	case "foreign_key":
		return fmt.Sprintf("referenced resource does not exist%s", details)
	case "missing_table":
		return "database schema error: table not found"
	default:
		return fmt.Sprintf("database error during %s: %v", e.Operation, e.OriginalError)
	}
}

// Unwrap returns the original error
func (e *DatabaseError) Unwrap() error {
	return e.OriginalError
}

// IsUniqueConstraintError checks if the error is a unique constraint violation
func IsUniqueConstraintError(err error) bool {
	var dbErr *DatabaseError
	return As(err, &dbErr) && dbErr.Code == "unique_constraint"
}

// IsForeignKeyError checks if the error is a foreign key constraint violation
func IsForeignKeyError(err error) bool {
	var dbErr *DatabaseError
	return As(err, &dbErr) && dbErr.Code == "foreign_key" 
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	var valErr *ValidationError
	return As(err, &valErr)
}

// APIError represents an error response to be sent to API clients
type APIError struct {
	StatusCode int               `json:"-"`
	Message    string            `json:"message,omitempty"`
	Errors     map[string]string `json:"errors,omitempty"`
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Errors:     make(map[string]string),
	}
}

// AddError adds a field error
func (e *APIError) AddError(field, message string) {
	e.Errors[field] = message
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// WriteJSON writes the error as JSON to the http response
func (e *APIError) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	
	resp := map[string]interface{}{
		"error": e.Message,
	}
	if len(e.Errors) > 0 {
		resp["errors"] = e.Errors
	}
	
	return json.NewEncoder(w).Encode(resp)
}

// Helper function to parse JSON and handle errors consistently
func ParseJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return NewJSONParsingError(err)
	}
	return nil
}

// HTTPErrorFromError converts a Go error to an appropriate HTTP error
func HTTPErrorFromError(err error) *APIError {
	// Check for specific error types in order of specificity
	switch {
	case err == nil:
		return NewAPIError(http.StatusOK, "")
		
	case err == ErrNotFound:
		return NewAPIError(http.StatusNotFound, "Resource not found")
		
	case IsValidationError(err):
		var valErr *ValidationError
		As(err, &valErr)
		apiErr := NewAPIError(http.StatusBadRequest, "Validation failed")
		for field, msg := range valErr.Errors {
			apiErr.AddError(field, msg)
		}
		return apiErr
		
	case IsUniqueConstraintError(err):
		var dbErr *DatabaseError
		As(err, &dbErr)
		apiErr := NewAPIError(http.StatusConflict, "Resource already exists")
		if dbErr.Field != "" {
			apiErr.AddError(dbErr.Field, "already exists")
		}
		return apiErr
		
	case IsForeignKeyError(err):
		return NewAPIError(http.StatusBadRequest, "Referenced resource does not exist")
		
	case err == ErrUnauthorized:
		return NewAPIError(http.StatusUnauthorized, "Unauthorized access")
		
	case err == ErrForbidden:
		return NewAPIError(http.StatusForbidden, "Access forbidden")
	
	default:
		// For unexpected errors, log them but don't expose details
		return NewAPIError(http.StatusInternalServerError, "Internal server error")
	}
}

// For compatibility with standard errors package
func As(err error, target interface{}) bool {
	return fmt.Errorf("%w", err).(interface{ As(interface{}) bool }).As(target)
}

func Is(err, target error) bool {
	return fmt.Errorf("%w", err).(interface{ Is(error) bool }).Is(target)
}

func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}