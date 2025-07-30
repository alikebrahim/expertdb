package utils

import (
	"net/http"
	"expertdb/internal/api/response"
	"expertdb/internal/validation"
)

// RespondWithValidationErrors responds with validation errors using the standard format
func RespondWithValidationErrors(w http.ResponseWriter, validationResult *validation.ValidationResult) error {
	if !validationResult.HasErrors() {
		return nil
	}
	
	return response.ValidationError(w, validationResult.Errors())
}

// RespondWithValidationErrorsMap responds with validation errors as a field->message map
func RespondWithValidationErrorsMap(w http.ResponseWriter, validationResult *validation.ValidationResult) error {
	if !validationResult.HasErrors() {
		return nil
	}
	
	// Convert to the format expected by the frontend
	errorMap := validationResult.ErrorsMap()
	
	return response.JSON(w, http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error":   "Validation failed",
		"errors":  errorMap,
	})
}

// RespondWithPaginatedData responds with paginated data using standard format
func RespondWithPaginatedData(w http.ResponseWriter, data interface{}, totalCount int, pagination Pagination) error {
	responseData := BuildPaginationResponse(data, totalCount, pagination)
	return response.Success(w, http.StatusOK, "", responseData)
}

// RespondWithSimplePaginatedData responds with simple paginated data (no total counts)
func RespondWithSimplePaginatedData(w http.ResponseWriter, data interface{}, hasMore bool, pagination Pagination) error {
	responseData := BuildSimplePaginationResponse(data, hasMore, pagination)
	return response.Success(w, http.StatusOK, "", responseData)
}

// RespondWithCreated responds with a created resource
func RespondWithCreated(w http.ResponseWriter, id int64, message string) error {
	return response.Created(w, id, message)
}

// RespondWithSuccess responds with a success message and optional data
func RespondWithSuccess(w http.ResponseWriter, message string, data interface{}) error {
	if data == nil {
		return response.Success(w, http.StatusOK, message, map[string]interface{}{})
	}
	return response.Success(w, http.StatusOK, message, data)
}

// RespondWithError responds with an error using the standard error handling
func RespondWithError(w http.ResponseWriter, err error) error {
	return response.Error(w, err)
}

// RespondWithNoContent responds with 204 No Content
func RespondWithNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// RespondWithBadRequest responds with a bad request error
func RespondWithBadRequest(w http.ResponseWriter, message string) error {
	return response.BadRequest(w, message)
}

// RespondWithNotFound responds with a not found error
func RespondWithNotFound(w http.ResponseWriter, message string) error {
	return response.NotFound(w, message)
}

// RespondWithJSON responds with raw JSON data
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) error {
	return response.JSON(w, status, data)
}

// RespondWithValidationErrorStrings responds with validation errors from string array
func RespondWithValidationErrorStrings(w http.ResponseWriter, validationErrors []string) error {
	return response.ValidationError(w, validationErrors)
}

// RespondWithCustomError responds with a custom error message and optional details
func RespondWithCustomError(w http.ResponseWriter, status int, message string, details map[string]interface{}) error {
	errorData := map[string]interface{}{
		"error": message,
	}
	
	// Add any additional details
	for key, value := range details {
		errorData[key] = value
	}
	
	return response.JSON(w, status, errorData)
}