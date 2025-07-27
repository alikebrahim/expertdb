package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"expertdb/internal/logger"
)

// ExtractIDFromPath extracts and validates an ID from URL path parameter
func ExtractIDFromPath(r *http.Request, paramName string, resourceType string) (int64, error) {
	log := logger.Get()
	idStr := r.PathValue(paramName)
	
	if idStr == "" {
		log.Warn("Missing %s ID in URL path", resourceType)
		return 0, fmt.Errorf("missing %s ID", resourceType)
	}
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid %s ID provided: %s", resourceType, idStr)
		return 0, fmt.Errorf("invalid %s ID: must be a number", resourceType)
	}
	
	if id <= 0 {
		log.Warn("Invalid %s ID provided: %d (must be positive)", resourceType, id)
		return 0, fmt.Errorf("invalid %s ID: must be a positive number", resourceType)
	}
	
	return id, nil
}

// ExtractOptionalIDFromPath extracts an optional ID from URL path parameter
// Returns 0 if the parameter is not present (no error)
func ExtractOptionalIDFromPath(r *http.Request, paramName string, resourceType string) (int64, error) {
	log := logger.Get()
	idStr := r.PathValue(paramName)
	
	if idStr == "" {
		return 0, nil // Optional parameter not provided
	}
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid %s ID provided: %s", resourceType, idStr)
		return 0, fmt.Errorf("invalid %s ID: must be a number", resourceType)
	}
	
	if id <= 0 {
		log.Warn("Invalid %s ID provided: %d (must be positive)", resourceType, id)
		return 0, fmt.Errorf("invalid %s ID: must be a positive number", resourceType)
	}
	
	return id, nil
}

// ExtractIntFromQuery extracts and validates an integer from query parameter
func ExtractIntFromQuery(r *http.Request, paramName string, defaultValue int) int {
	valueStr := r.URL.Query().Get(paramName)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 0 {
		return defaultValue
	}
	
	return value
}

// ExtractStringFromQuery extracts a string from query parameter with trimming
func ExtractStringFromQuery(r *http.Request, paramName string) string {
	value := r.URL.Query().Get(paramName)
	return value // Query parameters are already URL-decoded
}

// ExtractMultiValueFromQuery extracts comma-separated values from query parameter
func ExtractMultiValueFromQuery(r *http.Request, paramName string) []string {
	param := r.URL.Query().Get(paramName)
	if param == "" {
		return []string{}
	}
	
	// Split by comma and filter empty values
	values := make([]string, 0)
	for _, value := range splitAndTrim(param, ",") {
		if value != "" {
			values = append(values, value)
		}
	}
	
	return values
}

// splitAndTrim splits a string by delimiter and trims whitespace from each part
func splitAndTrim(s string, delimiter string) []string {
	if s == "" {
		return []string{}
	}
	
	parts := make([]string, 0)
	for _, part := range splitString(s, delimiter) {
		trimmed := trimWhitespace(part)
		parts = append(parts, trimmed)
	}
	
	return parts
}

// splitString splits a string by delimiter
func splitString(s string, delimiter string) []string {
	if s == "" {
		return []string{}
	}
	
	result := []string{}
	start := 0
	
	for i := 0; i <= len(s)-len(delimiter); i++ {
		if s[i:i+len(delimiter)] == delimiter {
			result = append(result, s[start:i])
			start = i + len(delimiter)
			i += len(delimiter) - 1
		}
	}
	
	// Add the last part
	result = append(result, s[start:])
	return result
}

// trimWhitespace removes leading and trailing whitespace
func trimWhitespace(s string) string {
	start := 0
	end := len(s)
	
	// Find first non-whitespace character
	for start < end && isWhitespace(s[start]) {
		start++
	}
	
	// Find last non-whitespace character
	for end > start && isWhitespace(s[end-1]) {
		end--
	}
	
	return s[start:end]
}

// isWhitespace checks if a character is whitespace
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}