package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type CreateExpertRequest struct {
	Name           string   `json:"name"`           // Full name of the expert
	Affiliation    string   `json:"affiliation"`    // Organization or institution the expert is affiliated with
	PrimaryContact string   `json:"primaryContact"` // Main contact information (email or phone)
	ContactType    string   `json:"contactType"`    // Type of contact information: "email" or "phone"
	Skills         []string `json:"skills"`         // List of expert's skills and competencies
	Role           string   `json:"role"`           // Expert's role (evaluator, validator, consultant, etc.)
	EmploymentType string   `json:"employmentType"` // Type of employment (academic, employer, freelance, etc.)
	GeneralArea    int64    `json:"generalArea"`    // ID referencing expert_areas table
	CVPath         string   `json:"cvPath"`         // Path to the expert's CV file
	Biography      string   `json:"biography"`      // Short biography or professional summary
	IsBahraini     bool     `json:"isBahraini"`     // Flag indicating if expert is Bahraini citizen
	Availability   string   `json:"availability"`   // Availability status: "yes"/"full-time" means active
}

type CreateExpertResponse struct {
	ID      int64  `json:"id"`      // ID of the newly created expert
	Success bool   `json:"success"` // Indicates if the creation was successful
	Message string `json:"message,omitempty"` // Optional message providing additional details
}

// ISCED classification types have been removed as part of schema simplification

// Expert represents a domain expert in the system
type Expert struct {
	ID              int64     `json:"id"`              // Primary key identifier
	ExpertID        string    `json:"expertId,omitempty"` // Business identifier
	Name            string    `json:"name"`            // Full name of the expert
	Designation     string    `json:"designation"`     // Professional title or position
	Institution     string    `json:"institution"`     // Organization or institution affiliation
	IsBahraini      bool      `json:"isBahraini"`      // Flag indicating if expert is Bahraini citizen
	Nationality     string    `json:"nationality"`     // Expert's nationality 
	IsAvailable     bool      `json:"isAvailable"`     // Current availability status for assignments
	Rating          string    `json:"rating"`          // Performance rating (if provided)
	Role            string    `json:"role"`            // Expert's role (evaluator, validator, consultant, etc.)
	EmploymentType  string    `json:"employmentType"`  // Type of employment (academic, employer, freelance, etc.)
	GeneralArea     int64     `json:"generalArea"`     // ID referencing expert_areas table 
	GeneralAreaName string    `json:"generalAreaName"` // Name of the general area (from expert_areas table)
	SpecializedArea string    `json:"specializedArea"` // Specific field of specialization
	IsTrained       bool      `json:"isTrained"`       // Indicates if expert has completed required training
	CVPath          string    `json:"cvPath"`          // Path to the expert's CV file
	Phone           string    `json:"phone"`           // Contact phone number
	Email           string    `json:"email"`           // Contact email address
	IsPublished     bool      `json:"isPublished"`     // Indicates if expert profile should be publicly visible
	Biography       string    `json:"biography"`       // Professional summary or background
	Documents       []Document `json:"documents,omitempty"` // Associated documents
	Engagements     []Engagement `json:"engagements,omitempty"` // Associated engagements
	CreatedAt       time.Time `json:"createdAt"`       // Timestamp when expert was created
	UpdatedAt       time.Time `json:"updatedAt"`       // Timestamp when expert was last updated
}

// Area represents an expert specialization area
type Area struct {
	ID   int64  `json:"id"`   // Unique identifier for the area
	Name string `json:"name"` // Name of the specialization area
}

// ExpertRequest represents a request to add a new expert
type ExpertRequest struct {
	ID              int64     `json:"id"`              // Primary key identifier
	ExpertID        string    `json:"expertId,omitempty"` // Business identifier (assigned after approval)
	Name            string    `json:"name"`            // Full name of the expert
	Designation     string    `json:"designation"`     // Professional title or position
	Institution     string    `json:"institution"`     // Organization or institution affiliation
	IsBahraini      bool      `json:"isBahraini"`      // Flag indicating if expert is Bahraini citizen
	IsAvailable     bool      `json:"isAvailable"`     // Current availability status for assignments
	Rating          string    `json:"rating"`          // Performance rating (if provided)
	Role            string    `json:"role"`            // Expert's role (evaluator, validator, consultant, etc.)
	EmploymentType  string    `json:"employmentType"`  // Type of employment (academic, employer, freelance, etc.)
	GeneralArea     int64     `json:"generalArea"`     // ID referencing expert_areas table
	SpecializedArea string    `json:"specializedArea"` // Specific field of specialization
	IsTrained       bool      `json:"isTrained"`       // Indicates if expert has completed required training
	CVPath          string    `json:"cvPath"`          // Path to the expert's CV file
	Phone           string    `json:"phone"`           // Contact phone number
	Email           string    `json:"email"`           // Contact email address
	IsPublished     bool      `json:"isPublished"`     // Indicates if expert profile should be publicly visible
	Status          string    `json:"status"`          // Request status: "pending", "approved", "rejected"
	RejectionReason string    `json:"rejectionReason,omitempty"` // Reason for rejection if status is "rejected"
	Biography       string    `json:"biography"`       // Professional summary or background
	CreatedAt       time.Time `json:"createdAt"`       // Timestamp when request was submitted
	ReviewedAt      time.Time `json:"reviewedAt,omitempty"` // Timestamp when request was reviewed
	ReviewedBy      int64     `json:"reviewedBy,omitempty"` // ID of admin who reviewed the request
}

// User represents a system user
type User struct {
	ID           int64     `json:"id"`           // Primary key identifier
	Name         string    `json:"name"`         // Full name of the user
	Email        string    `json:"email"`        // Email address (used for login)
	PasswordHash string    `json:"-"`            // Hashed password (never exposed in JSON)
	Role         string    `json:"role"`         // User role: "admin" or "user"
	IsActive     bool      `json:"isActive"`     // Account status (active/inactive)
	CreatedAt    time.Time `json:"createdAt"`    // Timestamp when user was created
	LastLogin    time.Time `json:"lastLogin,omitempty"` // Timestamp of last successful login
}

// Document represents an uploaded document for an expert
type Document struct {
	ID           int64     `json:"id"`           // Primary key identifier
	ExpertID     int64     `json:"expertId"`     // Foreign key reference to expert
	DocumentType string    `json:"documentType"` // Type of document: "cv", "certificate", "publication", etc.
	Type         string    `json:"type"`         // Alias for DocumentType for API compatibility
	Filename     string    `json:"filename"`     // Original filename as uploaded
	FilePath     string    `json:"filePath"`     // Path where file is stored on server
	ContentType  string    `json:"contentType"`  // MIME type of the document
	FileSize     int64     `json:"fileSize"`     // Size of document in bytes
	UploadDate   time.Time `json:"uploadDate"`   // Timestamp when document was uploaded
}

// Engagement represents expert assignment to projects/activities
type Engagement struct {
	ID             int64     `json:"id"`             // Primary key identifier
	ExpertID       int64     `json:"expertId"`       // Foreign key reference to expert
	EngagementType string    `json:"engagementType"` // Type of work: "evaluation", "consultation", "project", etc.
	StartDate      time.Time `json:"startDate"`      // Date when engagement begins
	EndDate        time.Time `json:"endDate,omitempty"` // Date when engagement ends
	ProjectName    string    `json:"projectName,omitempty"` // Name of the project or activity
	Status         string    `json:"status"`         // Current status: "pending", "active", "completed", "cancelled"
	FeedbackScore  int       `json:"feedbackScore,omitempty"` // Performance rating (1-5 scale)
	Notes          string    `json:"notes,omitempty"` // Additional comments or observations
	CreatedAt      time.Time `json:"createdAt"`      // Timestamp when record was created
}


// Statistics represents system-wide statistics
type Statistics struct {
	TotalExperts         int          `json:"totalExperts"`         // Total number of experts in the system
	ActiveCount          int          `json:"activeCount"`          // Number of experts marked as available
	BahrainiPercentage   float64      `json:"bahrainiPercentage"`   // Percentage of experts who are Bahraini nationals
	TopAreas             []AreaStat   `json:"topAreas"`             // Most common expertise areas
	EngagementsByType    []AreaStat   `json:"engagementsByType"`    // Distribution of engagements by type
	MonthlyGrowth        []GrowthStat `json:"monthlyGrowth"`        // Monthly growth in expert count
	MostRequestedExperts []ExpertStat `json:"mostRequestedExperts"` // Most frequently requested experts
	LastUpdated          time.Time    `json:"lastUpdated"`          // Timestamp when statistics were last calculated
}

// AreaStat represents statistics for a specific area/category
type AreaStat struct {
	Name       string  `json:"name"`       // Name of the area or category
	Count      int     `json:"count"`      // Number of items in this area
	Percentage float64 `json:"percentage"` // Percentage of total this area represents
}

// GrowthStat represents growth statistics over time
type GrowthStat struct {
	Period     string  `json:"period"`     // Time period identifier: "2023-01", "2023-Q1", etc.
	Count      int     `json:"count"`      // Number of items in this period
	GrowthRate float64 `json:"growthRate"` // Percentage growth from previous period
}

// ExpertStat represents statistics for a specific expert
type ExpertStat struct {
	ExpertID string `json:"expertId"` // Business identifier for the expert
	Name     string `json:"name"`     // Expert's name
	Count    int    `json:"count"`    // Number of requests/engagements for this expert
}

// DocumentUploadRequest represents a request to upload a document
type DocumentUploadRequest struct {
	ExpertID     int64  `json:"expertId"`     // ID of the expert to associate the document with
	DocumentType string `json:"documentType"` // Type of document: "cv", "certificate", "publication", etc.
}


// Authentication types

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email"`    // User's email address for authentication
	Password string `json:"password"` // User's password (plaintext in request only)
}

// LoginResponse represents a user login response
type LoginResponse struct {
	User  User   `json:"user"`  // User information (excluding password)
	Token string `json:"token"` // JWT token for authentication
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Name     string `json:"name"`     // Full name of the user
	Email    string `json:"email"`    // Email address (used for login)
	Password string `json:"password"` // Initial password (plaintext in request only)
	Role     string `json:"role"`     // User role: "admin" or "user"
	IsActive bool   `json:"isActive"` // Initial account status
}

// CreateUserResponse represents a response to creating a new user
type CreateUserResponse struct {
	ID      int64  `json:"id"`      // ID of the newly created user
	Success bool   `json:"success"` // Indicates if the creation was successful
	Message string `json:"message,omitempty"` // Optional message providing additional details
}

// Configuration represents application configuration
type Configuration struct {
	Port             string `json:"port"`             // HTTP server port
	DBPath           string `json:"dbPath"`           // Path to SQLite database file
	UploadPath       string `json:"uploadPath"`       // Directory for uploaded documents
	CORSAllowOrigins string `json:"corsAllowOrigins"` // CORS allowed origins (comma-separated)
}

func NewExpert(req CreateExpertRequest) *Expert {
	var email, phone string
	if req.ContactType == "email" {
		email = req.PrimaryContact
	} else {
		email = ""
	}

	if req.ContactType == "phone" {
		phone = req.PrimaryContact
	} else {
		phone = ""
	}

	return &Expert{
		Name:           req.Name,
		Institution:    req.Affiliation,
		IsAvailable:    req.Availability == "yes" || req.Availability == "full-time",
		Email:          email,
		Phone:          phone,
		Role:           req.Role,
		EmploymentType: req.EmploymentType,
		GeneralArea:    req.GeneralArea,    // Now expecting an int64 ID referencing expert_areas
		CVPath:         req.CVPath,
		Biography:      req.Biography,
		IsBahraini:     req.IsBahraini,
		CreatedAt:      time.Now().UTC(),
	}
}

// ValidateCreateExpertRequest validates the expert request fields
func ValidateCreateExpertRequest(req *CreateExpertRequest) error {
	// Required fields
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(req.PrimaryContact) == "" {
		return errors.New("primary contact is required")
	}

	// Validate contact based on type
	if req.ContactType == "email" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.PrimaryContact) {
			return errors.New("invalid email format")
		}
	} else if req.ContactType == "phone" {
		phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !phoneRegex.MatchString(req.PrimaryContact) {
			return errors.New("invalid phone number format")
		}
	}

	// Set default contact type if not provided
	if req.ContactType == "" {
		req.ContactType = "email"
	}

	// Validate new required fields
	if strings.TrimSpace(req.Role) == "" {
		return errors.New("role is required")
	}

	if strings.TrimSpace(req.EmploymentType) == "" {
		return errors.New("employment type is required")
	}

	if req.GeneralArea == 0 {
		return errors.New("general area is required")
	}

	// Validate role values
	validRoles := []string{"evaluator", "validator", "consultant", "trainer", "expert"}
	if !containsString(validRoles, strings.ToLower(req.Role)) {
		return errors.New("role must be one of: evaluator, validator, consultant, trainer, expert")
	}

	// Validate employment type values
	validEmploymentTypes := []string{"academic", "employer", "freelance", "government", "other"}
	if !containsString(validEmploymentTypes, strings.ToLower(req.EmploymentType)) {
		return errors.New("employment type must be one of: academic, employer, freelance, government, other")
	}

	// Limit biography length
	if len(req.Biography) > 1000 {
		return errors.New("biography exceeds maximum length of 1000 characters")
	}

	return nil
}

// Helper function to check if a string is in a slice
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// NOTE: Added field comments to all structs for enhanced code documentation and clarity.
