// Package response provides utility functions for handling HTTP responses
package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"expertdb/internal/errors"
	"expertdb/internal/logger"
)

// Standard response structures

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error  string            `json:"error"`
	Errors map[string]string `json:"errors,omitempty"`
}

// ValidationErrorResponse represents a response with validation errors
type ValidationErrorResponse struct {
	Error  string   `json:"error"`
	Errors []string `json:"errors"`
}

// JSON writes a JSON response with the given status code and data
func JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Success writes a success response with the given status code, message, and optional data
func Success(w http.ResponseWriter, status int, message string, data interface{}) error {
	resp := SuccessResponse{
		Success: true,
		Message: message,
	}
	
	if data != nil {
		resp.Data = data
	}
	
	return JSON(w, status, resp)
}

// Created writes a success response for resource creation with ID
func Created(w http.ResponseWriter, id int64, message string) error {
	resp := map[string]interface{}{
		"id":      id,
		"success": true,
		"message": message,
	}
	
	return JSON(w, http.StatusCreated, resp)
}

// Error writes an error response based on an error object
func Error(w http.ResponseWriter, err error) error {
	log := logger.Get()
	
	// Convert the error to an API error
	apiErr := errors.HTTPErrorFromError(err)
	
	// Log the error appropriately based on status code
	if apiErr.StatusCode >= 500 {
		log.Error("Server error: %v", err)
	} else {
		log.Debug("Client error: %v", err)
	}
	
	// Write the error response
	return apiErr.WriteJSON(w)
}

// BadRequest writes a bad request error response
func BadRequest(w http.ResponseWriter, message string) error {
	return JSON(w, http.StatusBadRequest, ErrorResponse{
		Error: message,
	})
}

// ValidationError writes a validation error response with multiple error messages
func ValidationError(w http.ResponseWriter, validationErrors []string) error {
	return JSON(w, http.StatusBadRequest, ValidationErrorResponse{
		Error:  "Validation failed",
		Errors: validationErrors,
	})
}

// NotFound writes a not found error response
func NotFound(w http.ResponseWriter, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return JSON(w, http.StatusNotFound, ErrorResponse{
		Error: message,
	})
}

// InternalError writes an internal server error response
func InternalError(w http.ResponseWriter, err error) error {
	log := logger.Get()
	log.Error("Internal server error: %v", err)
	
	return JSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: "Internal server error",
	})
}

// Conflict writes a conflict error response
func Conflict(w http.ResponseWriter, message string) error {
	return JSON(w, http.StatusConflict, ErrorResponse{
		Error: message,
	})
}

// Unauthorized writes an unauthorized error response
func Unauthorized(w http.ResponseWriter, message string) error {
	if message == "" {
		message = "Unauthorized access"
	}
	return JSON(w, http.StatusUnauthorized, ErrorResponse{
		Error: message,
	})
}

// Forbidden writes a forbidden error response
func Forbidden(w http.ResponseWriter, message string) error {
	if message == "" {
		message = "Access forbidden"
	}
	return JSON(w, http.StatusForbidden, ErrorResponse{
		Error: message,
	})
}

// ParseJSONError parses the common JSON error patterns and returns a user-friendly message
func ParseJSONError(err error) string {
	msg := err.Error()
	
	if strings.Contains(msg, "unexpected EOF") {
		return "Incomplete JSON data provided"
	} else if strings.Contains(msg, "cannot unmarshal") {
		return "Invalid data type in JSON"
	} else if strings.Contains(msg, "invalid character") {
		return "Invalid JSON format"
	}
	
	return "Error parsing request body: " + msg
}

// HandleJSONParsingError handles JSON parsing errors with appropriate responses
func HandleJSONParsingError(w http.ResponseWriter, err error) error {
	message := ParseJSONError(err)
	logger.Get().Warn("JSON parsing error: %v", err)
	return BadRequest(w, message)
}

// DatabaseError writes a database error response with appropriate details
func DatabaseError(w http.ResponseWriter, err *errors.DatabaseError) error {
	log := logger.Get()
	log.Error("Database error: %v", err)
	
	switch err.Code {
	case "unique_constraint":
		field := err.Field
		if field == "" {
			field = "resource"
		}
		return Conflict(w, field+" already exists")
		
	case "foreign_key":
		return BadRequest(w, "Referenced resource does not exist")
		
	default:
		return InternalError(w, err)
	}
}