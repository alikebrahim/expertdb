// Package domain contains the core business entities for the ExpertDB application
package domain

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

// Domain errors
var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("access forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrValidation         = errors.New("validation error")
	ErrBadRequest         = errors.New("bad request")
	ErrInternalServer     = errors.New("internal server error")
)

// Biography structures for structured expert profiles
type Biography struct {
	Experience []ExperienceEntry `json:"experience"` // Professional experience entries
	Education  []EducationEntry  `json:"education"`  // Educational background entries
}

type ExperienceEntry struct {
	StartDate    string `json:"start_date"`    // Start date of experience
	EndDate      string `json:"end_date"`      // End date of experience (can be "Present")
	Title        string `json:"title"`         // Job title or position
	Organization string `json:"organization"`  // Organization or company name
	Description  string `json:"description"`   // Description of responsibilities and achievements
}

type EducationEntry struct {
	StartDate   string `json:"start_date"`   // Start date of education
	EndDate     string `json:"end_date"`     // End date of education
	Title       string `json:"title"`        // Degree, qualification, or program title
	Institution string `json:"institution"`  // Educational institution name
}

type CreateExpertRequest struct {
	Name            string     `json:"name"`            // Full name of the expert
	Designation     string     `json:"designation"`     // Professional title: Prof., Dr., Mr., Ms., Mrs., Miss, Eng.
	Affiliation     string     `json:"affiliation"`     // Organization or institution the expert is affiliated with
	Phone           string     `json:"phone"`           // Contact phone number
	Email           string     `json:"email"`           // Contact email address
	Skills          []string   `json:"skills"`          // List of expert's skills and competencies
	Role            string     `json:"role"`            // Expert's role: "evaluator", "validator", or "evaluator/validator"
	EmploymentType  string     `json:"employmentType"`  // Type of employment: "academic" or "employer"
	GeneralArea     int64      `json:"generalArea"`     // ID referencing expert_areas table
	SpecializedArea string     `json:"specializedArea"` // Specific field of specialization
	CVPath          string     `json:"cvPath"`          // Path to the expert's CV file
	Biography       Biography  `json:"biography"`       // Structured biography with experience and education
	IsBahraini      bool       `json:"isBahraini"`      // Flag indicating if expert is Bahraini citizen
	IsAvailable     bool       `json:"isAvailable"`     // Current availability status for assignments
	Rating          int        `json:"rating"`          // Performance rating (1-5 scale)
	IsTrained       bool       `json:"isTrained"`       // Indicates if expert has completed required training
	IsPublished     bool       `json:"isPublished"`     // Indicates if expert has published work
}

type CreateExpertResponse struct {
	ID      int64  `json:"id"`      // ID of the newly created expert
	Success bool   `json:"success"` // Indicates if the creation was successful
	Message string `json:"message,omitempty"` // Optional message providing additional details
}

// ISCED classification types have been removed as part of schema simplification

// Expert represents a domain expert in the system
type Expert struct {
	ID                  int64     `json:"id"`              // Primary key identifier
	ExpertID            string    `json:"expertId,omitempty"` // Business identifier
	Name                string    `json:"name"`            // Full name of the expert
	Designation         string    `json:"designation"`     // Professional title or position
	Affiliation         string    `json:"affiliation"`     // Organization or institution affiliation
	IsBahraini          bool      `json:"isBahraini"`      // Flag indicating if expert is Bahraini citizen
	Nationality         string    `json:"nationality"`     // Expert's nationality 
	IsAvailable         bool      `json:"isAvailable"`     // Current availability status for assignments
	Rating              int       `json:"rating"`          // Performance rating (1-5 scale)
	Role                string    `json:"role"`            // Expert's role: "evaluator", "validator", or "evaluator/validator"
	EmploymentType      string    `json:"employmentType"`  // Type of employment: "academic" or "employer"
	GeneralArea         int64     `json:"generalArea"`     // ID referencing expert_areas table 
	GeneralAreaName     string    `json:"generalAreaName"` // Name of the general area (from expert_areas table)
	SpecializedArea     string    `json:"specializedArea"` // Specific field of specialization
	IsTrained           bool      `json:"isTrained"`       // Indicates if expert has completed required training
	CVPath              string    `json:"cvPath"`          // Path to the expert's CV file
	ApprovalDocumentPath string    `json:"approvalDocumentPath,omitempty"` // Path to the approval document
	Phone               string    `json:"phone"`           // Contact phone number
	Email               string    `json:"email"`           // Contact email address
	IsPublished         bool      `json:"isPublished"`     // Indicates if expert profile should be publicly visible
	Biography           string    `json:"biography"`       // Professional summary or background
	Documents           []Document `json:"documents,omitempty"` // Associated documents
	Engagements         []Engagement `json:"engagements,omitempty"` // Associated engagements
	CreatedAt           time.Time `json:"createdAt"`       // Timestamp when expert was created
	UpdatedAt           time.Time `json:"updatedAt"`       // Timestamp when expert was last updated
	OriginalRequestID   int64     `json:"originalRequestId,omitempty"` // Reference to the request that created this expert
}

// Area represents an expert specialization area
type Area struct {
	ID   int64  `json:"id"`   // Unique identifier for the area
	Name string `json:"name"` // Name of the specialization area
}

// ExpertRequest represents a request to add a new expert
type ExpertRequest struct {
	ID                   int64     `json:"id"`                   // Primary key identifier
	ExpertID             string    `json:"expertId,omitempty"`   // Business identifier (assigned after approval)
	Name                 string    `json:"name"`                 // Full name of the expert
	Designation          string    `json:"designation"`          // Professional title: Prof., Dr., Mr., Ms., Mrs., Miss, Eng.
	Affiliation          string    `json:"affiliation"`          // Organization or institution the expert is affiliated with
	Phone                string    `json:"phone"`                // Contact phone number
	Email                string    `json:"email"`                // Contact email address
	IsBahraini           bool      `json:"isBahraini"`           // Flag indicating if expert is Bahraini citizen
	IsAvailable          bool      `json:"isAvailable"`          // Current availability status for assignments
	Rating               int       `json:"rating"`               // Performance rating (1-5 scale)
	Role                 string    `json:"role"`                 // Expert's role: "evaluator", "validator", or "evaluator/validator"
	EmploymentType       string    `json:"employmentType"`       // Type of employment: "academic" or "employer"
	GeneralArea          int64     `json:"generalArea"`          // ID referencing expert_areas table
	SpecializedArea      string    `json:"specializedArea"`      // Specific field of specialization
	IsTrained            bool      `json:"isTrained"`            // Indicates if expert has completed required training
	IsPublished          bool      `json:"isPublished"`          // Indicates if expert profile should be publicly visible
	CVPath               string    `json:"cvPath"`               // Path to the expert's CV file
	ApprovalDocumentPath string    `json:"approvalDocumentPath,omitempty"` // Path to the approval document
	Biography            string    `json:"biography"`            // Structured biography as JSON string
	Status               string    `json:"status"`               // Request status: "pending", "approved", "rejected"
	RejectionReason      string    `json:"rejectionReason,omitempty"` // Reason for rejection if status is "rejected"
	CreatedAt            time.Time `json:"createdAt"`            // Timestamp when request was submitted
	ReviewedAt           time.Time `json:"reviewedAt,omitempty"` // Timestamp when request was reviewed
	ReviewedBy           int64     `json:"reviewedBy,omitempty"` // ID of admin who reviewed the request
	CreatedBy            int64     `json:"createdBy,omitempty"`  // ID of user who created the request
}

// User represents a system user
type User struct {
	ID           int64     `json:"id"`           // Primary key identifier
	Name         string    `json:"name"`         // Full name of the user
	Email        string    `json:"email"`        // Email address (used for login)
	PasswordHash string    `json:"-"`            // Hashed password (never exposed in JSON)
	Role         string    `json:"role"`         // User role: "super_user", "admin", "planner", or "user"
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
	PublishedCount       int          `json:"publishedCount"`        // Number of experts marked as published
	PublishedRatio       float64      `json:"publishedRatio"`        // Percentage of experts who are published
	TopAreas             []AreaStat   `json:"topAreas"`             // Most common expertise areas
	EngagementsByType    []AreaStat   `json:"engagementsByType"`    // Distribution of engagements by type
	YearlyGrowth         []GrowthStat `json:"yearlyGrowth"`         // Yearly growth in expert count
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

// Phase represents a collection of qualification applications to be processed
type Phase struct {
	ID                int64             `json:"id"`                // Primary key identifier
	PhaseID           string            `json:"phaseId"`           // Business identifier (e.g., "PH-2025-001")
	Title             string            `json:"title"`             // Title/name of the phase
	AssignedPlannerID int64           `json:"assignedPlannerId"` // ID of planner user assigned to this phase
	PlannerName     string            `json:"plannerName,omitempty"` // Name of assigned planner (not stored in DB)
	Status            string            `json:"status"`            // Status: "draft", "in_progress", "completed", "cancelled"
	Applications      []PhaseApplication `json:"applications,omitempty"` // List of applications in this phase
	CreatedAt         time.Time         `json:"createdAt"`         // When the phase was created
	UpdatedAt         time.Time         `json:"updatedAt"`         // When the phase was last updated
}

// PhaseApplication represents an application for a qualification requiring expert review
type PhaseApplication struct {
	ID              int64     `json:"id"`              // Primary key identifier
	PhaseID         int64     `json:"phaseId"`         // Foreign key reference to phases table
	Type            string    `json:"type"`            // Type: "QP" (Qualification Placement) or "IL" (Institutional Listing)
	InstitutionName string    `json:"institutionName"` // Name of the institution
	QualificationName string  `json:"qualificationName"` // Name of the qualification being reviewed
	Expert1         int64     `json:"expert1"`         // First expert ID
	Expert1Name     string    `json:"expert1Name,omitempty"` // First expert name (not stored in DB)
	Expert2         int64     `json:"expert2"`         // Second expert ID
	Expert2Name     string    `json:"expert2Name,omitempty"` // Second expert name (not stored in DB)
	Status          string    `json:"status"`          // Status: "pending", "assigned", "approved", "rejected"
	RejectionNotes  string    `json:"rejectionNotes,omitempty"` // Notes for rejection (if status is "rejected")
	CreatedAt       time.Time `json:"createdAt"`       // When the application was created
	UpdatedAt       time.Time `json:"updatedAt"`       // When the application was last updated
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
	Role     string `json:"role"`     // User role: "super_user", "admin", "planner", or "user"
	IsActive bool   `json:"isActive"` // Initial account status
}

// CreateUserResponse represents a response to creating a new user
type CreateUserResponse struct {
	ID      int64  `json:"id"`      // ID of the newly created user
	Success bool   `json:"success"` // Indicates if the creation was successful
	Message string `json:"message,omitempty"` // Optional message providing additional details
}

// NewExpert creates a new Expert from a CreateExpertRequest
func NewExpert(req CreateExpertRequest) *Expert {
	// Convert Biography struct to JSON string for storage
	biographyJSON := ""
	if len(req.Biography.Experience) > 0 || len(req.Biography.Education) > 0 {
		if data, err := json.Marshal(req.Biography); err == nil {
			biographyJSON = string(data)
		}
	}

	return &Expert{
		Name:            req.Name,
		Designation:     req.Designation,
		Affiliation:     req.Affiliation,
		IsAvailable:     req.IsAvailable,
		Email:           req.Email,
		Phone:           req.Phone,
		Role:            req.Role,
		EmploymentType:  req.EmploymentType,
		GeneralArea:     req.GeneralArea,
		SpecializedArea: req.SpecializedArea,
		CVPath:          req.CVPath,
		Biography:       biographyJSON,
		IsBahraini:      req.IsBahraini,
		IsTrained:       req.IsTrained,
		IsPublished:     req.IsPublished,
		Rating:          req.Rating,
		CreatedAt:       time.Now().UTC(),
	}
}

// ValidateCreateExpertRequest validates the expert request fields
func ValidateCreateExpertRequest(req *CreateExpertRequest) error {
	// Required fields
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(req.Designation) == "" {
		return errors.New("designation is required")
	}

	if strings.TrimSpace(req.Affiliation) == "" {
		return errors.New("affiliation is required")
	}

	if strings.TrimSpace(req.Phone) == "" {
		return errors.New("phone is required")
	}

	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(req.SpecializedArea) == "" {
		return errors.New("specialized area is required")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	// Validate phone format
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	if !phoneRegex.MatchString(req.Phone) {
		return errors.New("invalid phone number format")
	}

	// Validate designation
	validDesignations := []string{"Prof.", "Dr.", "Mr.", "Ms.", "Mrs.", "Miss", "Eng."}
	if !containsString(validDesignations, req.Designation) {
		return errors.New("designation must be one of: Prof., Dr., Mr., Ms., Mrs., Miss, Eng.")
	}

	// Validate role values
	if strings.TrimSpace(req.Role) == "" {
		return errors.New("role is required")
	}
	validRoles := []string{"evaluator", "validator", "evaluator/validator"}
	if !containsString(validRoles, strings.ToLower(req.Role)) {
		return errors.New("role must be one of: evaluator, validator, evaluator/validator")
	}

	// Validate employment type values
	if strings.TrimSpace(req.EmploymentType) == "" {
		return errors.New("employment type is required")
	}
	validEmploymentTypes := []string{"academic", "employer"}
	if !containsString(validEmploymentTypes, strings.ToLower(req.EmploymentType)) {
		return errors.New("employment type must be one of: academic, employer")
	}

	// Validate general area
	if req.GeneralArea == 0 {
		return errors.New("general area is required")
	}

	// Validate rating (1-5 scale)
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// Validate skills
	if len(req.Skills) == 0 {
		return errors.New("at least one skill is required")
	}

	// Validate biography structure
	if len(req.Biography.Experience) == 0 && len(req.Biography.Education) == 0 {
		return errors.New("biography must contain at least one experience or education entry")
	}

	// Validate experience entries
	for i, exp := range req.Biography.Experience {
		if strings.TrimSpace(exp.Title) == "" {
			return errors.New("experience title is required for entry " + string(rune(i+1)))
		}
		if strings.TrimSpace(exp.Organization) == "" {
			return errors.New("experience organization is required for entry " + string(rune(i+1)))
		}
		if strings.TrimSpace(exp.StartDate) == "" {
			return errors.New("experience start date is required for entry " + string(rune(i+1)))
		}
	}

	// Validate education entries
	for i, edu := range req.Biography.Education {
		if strings.TrimSpace(edu.Title) == "" {
			return errors.New("education title is required for entry " + string(rune(i+1)))
		}
		if strings.TrimSpace(edu.Institution) == "" {
			return errors.New("education institution is required for entry " + string(rune(i+1)))
		}
		if strings.TrimSpace(edu.StartDate) == "" {
			return errors.New("education start date is required for entry " + string(rune(i+1)))
		}
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