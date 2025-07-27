package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError holds a field-specific validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult holds the collection of validation errors
type ValidationResult struct {
	errors []ValidationError
}

// New creates a new validation result
func New() *ValidationResult {
	return &ValidationResult{errors: make([]ValidationError, 0)}
}

// Required validates that a string field is not empty
func (v *ValidationResult) Required(field string, value string, fieldName string) *ValidationResult {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s is required", fieldName),
		})
	}
	return v
}

// MinLength validates minimum length of a string field
func (v *ValidationResult) MinLength(field string, value string, minLen int, fieldName string) *ValidationResult {
	if len(strings.TrimSpace(value)) < minLen && value != "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be at least %d characters long", fieldName, minLen),
		})
	}
	return v
}

// MaxLength validates maximum length of a string field
func (v *ValidationResult) MaxLength(field string, value string, maxLen int, fieldName string) *ValidationResult {
	if len(value) > maxLen {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be no more than %d characters long", fieldName, maxLen),
		})
	}
	return v
}

// Email validates email format
func (v *ValidationResult) Email(field string, value string, fieldName string) *ValidationResult {
	if value != "" {
		// Basic email regex pattern
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			v.errors = append(v.errors, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s must be a valid email address", fieldName),
			})
		}
	}
	return v
}

// OneOf validates that value is one of the allowed options (case-insensitive)
func (v *ValidationResult) OneOf(field string, value string, options []string, fieldName string) *ValidationResult {
	if value != "" {
		found := false
		lowerValue := strings.ToLower(strings.TrimSpace(value))
		for _, option := range options {
			if lowerValue == strings.ToLower(option) {
				found = true
				break
			}
		}
		if !found {
			v.errors = append(v.errors, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s must be one of: %s", fieldName, strings.Join(options, ", ")),
			})
		}
	}
	return v
}

// Range validates that a numeric value is within a specified range
func (v *ValidationResult) Range(field string, value int, min, max int, fieldName string) *ValidationResult {
	if value < min || value > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be between %d and %d", fieldName, min, max),
		})
	}
	return v
}

// FloatRange validates that a float value is within a specified range
func (v *ValidationResult) FloatRange(field string, value float64, min, max float64, fieldName string) *ValidationResult {
	if value < min || value > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s must be between %.2f and %.2f", fieldName, min, max),
		})
	}
	return v
}

// Phone validates basic phone number format (digits, spaces, hyphens, parentheses, plus)
func (v *ValidationResult) Phone(field string, value string, fieldName string) *ValidationResult {
	if value != "" {
		// Allow digits, spaces, hyphens, parentheses, and plus sign
		phoneRegex := regexp.MustCompile(`^[\d\s\-\(\)\+]+$`)
		if !phoneRegex.MatchString(value) || len(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", ""), "(", ""), ")", ""), "+", "")) < 7 {
			v.errors = append(v.errors, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s must be a valid phone number", fieldName),
			})
		}
	}
	return v
}

// URL validates basic URL format
func (v *ValidationResult) URL(field string, value string, fieldName string) *ValidationResult {
	if value != "" {
		urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
		if !urlRegex.MatchString(value) {
			v.errors = append(v.errors, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("%s must be a valid URL", fieldName),
			})
		}
	}
	return v
}

// Custom allows for custom validation logic
func (v *ValidationResult) Custom(field string, isValid bool, message string) *ValidationResult {
	if !isValid {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: message,
		})
	}
	return v
}

// Conditional applies validation only if condition is true
func (v *ValidationResult) Conditional(condition bool, validationFunc func(*ValidationResult) *ValidationResult) *ValidationResult {
	if condition {
		return validationFunc(v)
	}
	return v
}

// IsValid returns true if no validation errors occurred
func (v *ValidationResult) IsValid() bool {
	return len(v.errors) == 0
}

// HasErrors returns true if validation errors occurred
func (v *ValidationResult) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns all validation errors as a slice of strings
func (v *ValidationResult) Errors() []string {
	messages := make([]string, len(v.errors))
	for i, err := range v.errors {
		messages[i] = err.Message
	}
	return messages
}

// ErrorsMap returns validation errors as a map of field -> error message
func (v *ValidationResult) ErrorsMap() map[string]string {
	errorMap := make(map[string]string)
	for _, err := range v.errors {
		errorMap[err.Field] = err.Message
	}
	return errorMap
}

// DetailedErrors returns the full ValidationError objects
func (v *ValidationResult) DetailedErrors() []ValidationError {
	return v.errors
}

// First returns the first validation error message, or empty string if no errors
func (v *ValidationResult) First() string {
	if len(v.errors) > 0 {
		return v.errors[0].Message
	}
	return ""
}

// AddError manually adds a validation error
func (v *ValidationResult) AddError(field, message string) *ValidationResult {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
	return v
}

// Merge combines validation results from multiple validators
func (v *ValidationResult) Merge(other *ValidationResult) *ValidationResult {
	v.errors = append(v.errors, other.errors...)
	return v
}

// ValidateString is a helper for common string validations
func ValidateString(field, value, fieldName string, required bool, minLen, maxLen int) *ValidationResult {
	validator := New()
	
	if required {
		validator.Required(field, value, fieldName)
	}
	
	if value != "" {
		if minLen > 0 {
			validator.MinLength(field, value, minLen, fieldName)
		}
		if maxLen > 0 {
			validator.MaxLength(field, value, maxLen, fieldName)
		}
	}
	
	return validator
}

// ValidateEmail is a helper for email validation with optional requirement
func ValidateEmail(field, value, fieldName string, required bool) *ValidationResult {
	validator := New()
	
	if required {
		validator.Required(field, value, fieldName)
	}
	
	validator.Email(field, value, fieldName)
	return validator
}

// ValidateChoice is a helper for choice validation with optional requirement
func ValidateChoice(field, value, fieldName string, options []string, required bool) *ValidationResult {
	validator := New()
	
	if required {
		validator.Required(field, value, fieldName)
	}
	
	validator.OneOf(field, value, options, fieldName)
	return validator
}

// ParseInt safely parses a string to int with validation
func ParseInt(field, value, fieldName string) (int, *ValidationResult) {
	validator := New()
	
	if value == "" {
		return 0, validator
	}
	
	parsed, err := strconv.Atoi(value)
	if err != nil {
		validator.AddError(field, fmt.Sprintf("%s must be a valid number", fieldName))
		return 0, validator
	}
	
	return parsed, validator
}

// ParseFloat safely parses a string to float64 with validation
func ParseFloat(field, value, fieldName string) (float64, *ValidationResult) {
	validator := New()
	
	if value == "" {
		return 0, validator
	}
	
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		validator.AddError(field, fmt.Sprintf("%s must be a valid number", fieldName))
		return 0, validator
	}
	
	return parsed, validator
}