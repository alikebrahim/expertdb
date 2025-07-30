// Package domain contains the core business entities for the ExpertDB application
package domain

import (
	"errors"
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

// Database-specific structs for structured expert profiles
type ExpertExperienceEntry struct {
	ID           int64     `json:"id" db:"id"`
	ExpertID     int64     `json:"expertId" db:"expert_id"`
	Organization string    `json:"organization" db:"organization"`
	Position     string    `json:"position" db:"position"`
	StartDate    string    `json:"startDate" db:"start_date"`
	EndDate      string    `json:"endDate" db:"end_date"`
	IsCurrent    bool      `json:"isCurrent" db:"is_current"`
	Country      string    `json:"country" db:"country"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type ExpertEducationEntry struct {
	ID             int64     `json:"id" db:"id"`
	ExpertID       int64     `json:"expertId" db:"expert_id"`
	Institution    string    `json:"institution" db:"institution"`
	Degree         string    `json:"degree" db:"degree"`
	FieldOfStudy   string    `json:"fieldOfStudy" db:"field_of_study"`
	GraduationYear string    `json:"graduationYear" db:"graduation_year"`
	Country        string    `json:"country" db:"country"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type ExpertRequestExperienceEntry struct {
	ID              int64  `json:"id" db:"id"`
	ExpertRequestID int64  `json:"expertRequestId" db:"expert_request_id"`
	Organization    string `json:"organization" db:"organization"`
	Position        string `json:"position" db:"position"`
	StartDate       string `json:"startDate" db:"start_date"`
	EndDate         string `json:"endDate" db:"end_date"`
	IsCurrent       bool   `json:"isCurrent" db:"is_current"`
	Country         string `json:"country" db:"country"`
	Description     string `json:"description" db:"description"`
}

type ExpertRequestEducationEntry struct {
	ID              int64  `json:"id" db:"id"`
	ExpertRequestID int64  `json:"expertRequestId" db:"expert_request_id"`
	Institution     string `json:"institution" db:"institution"`
	Degree          string `json:"degree" db:"degree"`
	FieldOfStudy    string `json:"fieldOfStudy" db:"field_of_study"`
	GraduationYear  string `json:"graduationYear" db:"graduation_year"`
	Country         string `json:"country" db:"country"`
	Description     string `json:"description" db:"description"`
}

type CreateExpertRequest struct {
	Name                      string                         `json:"name"`                      // Full name of the expert
	Designation               string                         `json:"designation"`               // Professional title: Prof., Dr., Mr., Ms., Mrs., Miss, Eng.
	Affiliation               string                         `json:"affiliation"`               // Organization or institution the expert is affiliated with
	Phone                     string                         `json:"phone"`                     // Contact phone number
	Email                     string                         `json:"email"`                     // Contact email address
	Role                      string                         `json:"role"`                      // Expert's role: "evaluator", "validator", or "evaluator/validator"
	EmploymentType            string                         `json:"employmentType"`            // Type of employment: "academic" or "employer"
	GeneralArea               int64                          `json:"generalArea"`               // ID referencing expert_areas table
	SpecializedAreaIds        []int64                        `json:"specializedAreaIds"`        // Selected existing area IDs
	SuggestedSpecializedAreas []string                       `json:"suggestedSpecializedAreas"` // User-suggested area names
	CVDocumentID              *int64                         `json:"cvDocumentId,omitempty"`    // Reference to CV document
	ExperienceEntries         []ExpertRequestExperienceEntry `json:"experienceEntries"`         // Professional experience entries
	EducationEntries          []ExpertRequestEducationEntry  `json:"educationEntries"`          // Educational background entries
	IsBahraini                bool                           `json:"isBahraini"`                // Flag indicating if expert is Bahraini citizen
	IsAvailable               bool                           `json:"isAvailable"`               // Current availability status for assignments
	Rating                    int                            `json:"rating"`                    // Performance rating (1-5 scale)
	IsTrained                 bool                           `json:"isTrained"`                 // Indicates if expert has completed required training
	IsPublished               bool                           `json:"isPublished"`               // Indicates if expert has published work
}

type CreateExpertResponse struct {
	ID      int64  `json:"id"`                // ID of the newly created expert
	Success bool   `json:"success"`           // Indicates if the creation was successful
	Message string `json:"message,omitempty"` // Optional message providing additional details
}

// Expert represents a domain expert in the system
type Expert struct {
	ID                       int64                   `json:"id"`                                          // Primary key identifier
	Name                     string                  `json:"name"`                                        // Full name of the expert
	Designation              string                  `json:"designation"`                                 // Professional title or position
	Affiliation              string                  `json:"affiliation"`                                 // Organization or institution affiliation
	IsBahraini               bool                    `json:"isBahraini"`                                  // Flag indicating if expert is Bahraini citizen
	IsAvailable              bool                    `json:"isAvailable"`                                 // Current availability status for assignments
	Rating                   int                     `json:"rating"`                                      // Performance rating (1-5 scale)
	Role                     string                  `json:"role"`                                        // Expert's role: "evaluator", "validator", or "evaluator/validator"
	EmploymentType           string                  `json:"employmentType"`                              // Type of employment: "academic" or "employer"
	GeneralArea              int64                   `json:"-" db:"general_area"`                         // ID referencing expert_areas table - internal use only
	GeneralAreaName          string                  `json:"generalAreaName"`                             // Name of the general area (from expert_areas table)
	SpecializedArea          string                  `json:"-" db:"specialized_area"`                     // Comma-separated specialized area IDs (e.g., "1,4,6") - internal use only
	SpecializedAreaNames     string                  `json:"specializedAreaNames"`                        // Comma-separated specialized area names (e.g., "Software Engineering, Database Design")
	SpecializedAreasResolved []*SpecializedArea      `json:"specialized_areas_resolved,omitempty" db:"-"` // Resolved specialized area names
	IsTrained                bool                    `json:"isTrained"`                                   // Indicates if expert has completed required training
	CVDocumentID             *int64                  `json:"cvDocumentId,omitempty"`                      // Reference to CV document
	ApprovalDocumentID       *int64                  `json:"approvalDocumentId,omitempty"`                // Reference to approval document
	CVDocument               *Document               `json:"cvDocument,omitempty"`                        // Resolved CV document
	ApprovalDocument         *Document               `json:"approvalDocument,omitempty"`                  // Resolved approval document
	Phone                    string                  `json:"phone"`                                       // Contact phone number
	Email                    string                  `json:"email"`                                       // Contact email address
	IsPublished              bool                    `json:"isPublished"`                                 // Indicates if expert profile should be publicly visible
	ExperienceEntries        []ExpertExperienceEntry `json:"experienceEntries,omitempty"`                 // Professional experience entries
	EducationEntries         []ExpertEducationEntry  `json:"educationEntries,omitempty"`                  // Educational background entries
	Documents                []Document              `json:"documents,omitempty"`                         // Associated documents
	Engagements              []Engagement            `json:"engagements,omitempty"`                       // Associated engagements
	CreatedAt                time.Time               `json:"createdAt"`                                   // Timestamp when expert was created
	UpdatedAt                time.Time               `json:"updatedAt"`                                   // Timestamp when expert was last updated
	OriginalRequestID        int64                   `json:"originalRequestId,omitempty"`                 // Reference to the request that created this expert
	LastEditedBy             *int64                  `json:"lastEditedBy,omitempty" db:"last_edited_by"`  // ID of user who last edited this expert
	LastEditedAt             *time.Time              `json:"lastEditedAt,omitempty" db:"last_edited_at"`  // Timestamp when expert was last edited
}

// ExpertEditHistoryEntry represents a single edit made to an expert profile
type ExpertEditHistoryEntry struct {
	ID            int64     `json:"id" db:"id"`                                // Primary key identifier
	ExpertID      int64     `json:"expertId" db:"expert_id"`                   // ID of the expert that was edited
	EditedBy      int64     `json:"editedBy" db:"edited_by"`                   // ID of user who made the edit
	EditedAt      time.Time `json:"editedAt" db:"edited_at"`                   // Timestamp when the edit was made
	FieldsChanged string    `json:"fieldsChanged" db:"fields_changed"`         // JSON array of field names that were changed
	OldValues     *string   `json:"oldValues,omitempty" db:"old_values"`       // JSON object of previous field values
	NewValues     *string   `json:"newValues,omitempty" db:"new_values"`       // JSON object of new field values
	ChangeReason  *string   `json:"changeReason,omitempty" db:"change_reason"` // Optional reason for the change
	EditorName    *string   `json:"editorName,omitempty" db:"-"`               // Name of the user who made the edit (resolved)
}

// Area represents an expert specialization area
type Area struct {
	ID   int64  `json:"id"`   // Unique identifier for the area
	Name string `json:"name"` // Name of the specialization area
}

// SpecializedArea represents a specialized area for experts
type SpecializedArea struct {
	ID        int64     `json:"id" db:"id"`                 // Unique identifier for the specialized area
	Name      string    `json:"name" db:"name"`             // Name of the specialized area
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Timestamp when specialized area was created
}

// ExpertRequest represents a request to add a new expert
type ExpertRequest struct {
	ID                        int64                          `json:"id"`                           // Primary key identifier
	Name                      string                         `json:"name"`                         // Full name of the expert
	Designation               string                         `json:"designation"`                  // Professional title: Prof., Dr., Mr., Ms., Mrs., Miss, Eng.
	Affiliation               string                         `json:"affiliation"`                  // Organization or institution the expert is affiliated with
	Phone                     string                         `json:"phone"`                        // Contact phone number
	Email                     string                         `json:"email"`                        // Contact email address
	IsBahraini                bool                           `json:"isBahraini"`                   // Flag indicating if expert is Bahraini citizen
	IsAvailable               bool                           `json:"isAvailable"`                  // Current availability status for assignments
	Role                      string                         `json:"role"`                         // Expert's role: "evaluator", "validator", or "evaluator/validator"
	EmploymentType            string                         `json:"employmentType"`               // Type of employment: "academic" or "employer"
	GeneralArea               int64                          `json:"generalArea"`                  // ID referencing expert_areas table
	SpecializedArea           string                         `json:"specializedArea"`              // Comma-separated existing area IDs
	SuggestedSpecializedAreas []string                       `json:"suggestedSpecializedAreas"`    // User-suggested area names
	IsTrained                 bool                           `json:"isTrained"`                    // Indicates if expert has completed required training
	IsPublished               bool                           `json:"isPublished"`                  // Indicates if expert profile should be publicly visible
	CVDocumentID              *int64                         `json:"cvDocumentId,omitempty"`       // Reference to CV document
	ApprovalDocumentID        *int64                         `json:"approvalDocumentId,omitempty"` // Reference to approval document
	CVDocument                *Document                      `json:"cvDocument,omitempty"`         // Resolved CV document
	ApprovalDocument          *Document                      `json:"approvalDocument,omitempty"`   // Resolved approval document
	ExperienceEntries         []ExpertRequestExperienceEntry `json:"experienceEntries"`            // Professional experience entries
	EducationEntries          []ExpertRequestEducationEntry  `json:"educationEntries"`             // Educational background entries
	Status                    string                         `json:"status"`                       // Request status: "pending", "approved", "rejected"
	RejectionReason           string                         `json:"rejectionReason,omitempty"`    // Reason for rejection if status is "rejected"
	CreatedAt                 time.Time                      `json:"createdAt"`                    // Timestamp when request was submitted
	ReviewedAt                time.Time                      `json:"reviewedAt,omitempty"`         // Timestamp when request was reviewed
	ReviewedBy                int64                          `json:"reviewedBy,omitempty"`         // ID of admin who reviewed the request
	CreatedBy                 int64                          `json:"createdBy,omitempty"`          // ID of user who created the request
}

// Document resolution methods for Expert
func (e *Expert) ResolveCVDocument(getDocument func(int64) (*Document, error)) error {
	if e.CVDocumentID != nil {
		doc, err := getDocument(*e.CVDocumentID)
		if err == nil {
			e.CVDocument = doc
		}
	}
	return nil
}

func (e *Expert) ResolveApprovalDocument(getDocument func(int64) (*Document, error)) error {
	if e.ApprovalDocumentID != nil {
		doc, err := getDocument(*e.ApprovalDocumentID)
		if err == nil {
			e.ApprovalDocument = doc
		}
	}
	return nil
}

// Document resolution methods for ExpertRequest
func (er *ExpertRequest) ResolveCVDocument(getDocument func(int64) (*Document, error)) error {
	if er.CVDocumentID != nil {
		doc, err := getDocument(*er.CVDocumentID)
		if err == nil {
			er.CVDocument = doc
		}
	}
	return nil
}

func (er *ExpertRequest) ResolveApprovalDocument(getDocument func(int64) (*Document, error)) error {
	if er.ApprovalDocumentID != nil {
		doc, err := getDocument(*er.ApprovalDocumentID)
		if err == nil {
			er.ApprovalDocument = doc
		}
	}
	return nil
}

// User represents a system user
type User struct {
	ID           int64     `json:"id"`                  // Primary key identifier
	Name         string    `json:"name"`                // Full name of the user
	Email        string    `json:"email"`               // Email address (used for login)
	PasswordHash string    `json:"-"`                   // Hashed password (never exposed in JSON)
	Role         string    `json:"role"`                // User role: "super_user", "admin", "planner", or "user"
	IsActive     bool      `json:"isActive"`            // Account status (active/inactive)
	CreatedAt    time.Time `json:"createdAt"`           // Timestamp when user was created
	LastLogin    time.Time `json:"lastLogin,omitempty"` // Timestamp of last successful login
}

// Document represents an uploaded document for an expert
type Document struct {
	ID           int64     `json:"id"`           // Primary key identifier
	ExpertID     int64     `json:"expertId"`     // Foreign key reference to expert
	DocumentType string    `json:"documentType"` // Type of document: "cv" or "approval"
	Filename     string    `json:"filename"`     // Original filename as uploaded
	FilePath     string    `json:"filePath"`     // Path where file is stored on server
	ContentType  string    `json:"contentType"`  // MIME type of the document
	FileSize     int64     `json:"fileSize"`     // Size of document in bytes
	UploadDate   time.Time `json:"uploadDate"`   // Timestamp when document was uploaded
}

// Engagement represents expert assignment to projects/activities
type Engagement struct {
	ID             int64     `json:"id"`                      // Primary key identifier
	ExpertID       int64     `json:"expertId"`                // Foreign key reference to expert
	EngagementType string    `json:"engagementType"`          // Type of work: "evaluation", "consultation", "project", etc.
	StartDate      time.Time `json:"startDate"`               // Date when engagement begins
	EndDate        time.Time `json:"endDate,omitempty"`       // Date when engagement ends
	ProjectName    string    `json:"projectName,omitempty"`   // Name of the project or activity
	Status         string    `json:"status"`                  // Current status: "pending", "active", "completed", "cancelled"
	FeedbackScore  int       `json:"feedbackScore,omitempty"` // Performance rating (1-5 scale)
	Notes          string    `json:"notes,omitempty"`         // Additional comments or observations
	CreatedAt      time.Time `json:"createdAt"`               // Timestamp when record was created
}

// Statistics represents system-wide statistics
type Statistics struct {
	TotalExperts         int           `json:"totalExperts"`         // Total number of experts in the system
	ActiveCount          int           `json:"activeCount"`          // Number of experts marked as available
	BahrainiPercentage   float64       `json:"bahrainiPercentage"`   // Percentage of experts who are Bahraini nationals
	PublishedCount       int           `json:"publishedCount"`       // Number of experts marked as published
	PublishedRatio       float64       `json:"publishedRatio"`       // Percentage of experts who are published
	TopAreas             []AreaStat    `json:"topAreas"`             // Most common expertise areas
	EngagementsByType    []AreaStat    `json:"engagementsByType"`    // Distribution of engagements by type
	EngagementsByStatus  []AreaStat    `json:"engagementsByStatus"`  // Distribution of engagements by status (active/completed)
	NationalityStats     []AreaStat    `json:"nationalityStats"`     // Bahraini vs Non-Bahraini breakdown with counts
	SpecializedAreas     AreaBreakdown `json:"specializedAreas"`     // Top and bottom specialized areas
	YearlyGrowth         []GrowthStat  `json:"yearlyGrowth"`         // Yearly growth in expert count
	MostRequestedExperts []ExpertStat  `json:"mostRequestedExperts"` // Most frequently requested experts
	LastUpdated          time.Time     `json:"lastUpdated"`          // Timestamp when statistics were last calculated
}

// AreaStat represents statistics for a specific area/category
type AreaStat struct {
	Name       string  `json:"name"`       // Name of the area or category
	Count      int     `json:"count"`      // Number of items in this area
	Percentage float64 `json:"percentage"` // Percentage of total this area represents
}

// AreaBreakdown represents top and bottom specialized areas
type AreaBreakdown struct {
	Top    []AreaStat `json:"top"`    // Top 5 specialized areas
	Bottom []AreaStat `json:"bottom"` // Bottom 5 specialized areas
}

// GrowthStat represents growth statistics over time
type GrowthStat struct {
	Period     string  `json:"period"`     // Time period identifier: "2023-01", "2023-Q1", etc.
	Count      int     `json:"count"`      // Number of items in this period
	GrowthRate float64 `json:"growthRate"` // Percentage growth from previous period
}

// ExpertStat represents statistics for a specific expert
type ExpertStat struct {
	ExpertID int64  `json:"expertId"` // Database ID of the expert
	Name     string `json:"name"`     // Expert's name
	Count    int    `json:"count"`    // Number of requests/engagements for this expert
}

// DocumentUploadRequest represents a request to upload a document
type DocumentUploadRequest struct {
	ExpertID     int64  `json:"expertId"`     // ID of the expert to associate the document with
	DocumentType string `json:"documentType"` // Type of document: "cv" or "approval"
}

// Phase represents a collection of qualification applications to be processed
type Phase struct {
	ID                int64              `json:"id"`                     // Primary key identifier
	PhaseID           string             `json:"phaseId"`                // Business identifier (e.g., "PH-2025-001")
	Title             string             `json:"title"`                  // Title/name of the phase
	AssignedPlannerID int64              `json:"assignedPlannerId"`      // ID of planner user assigned to this phase
	PlannerName       string             `json:"plannerName,omitempty"`  // Name of assigned planner (not stored in DB)
	Status            string             `json:"status"`                 // Status: "draft", "in_progress", "completed", "cancelled"
	Applications      []PhaseApplication `json:"applications,omitempty"` // List of applications in this phase
	CreatedAt         time.Time          `json:"createdAt"`              // When the phase was created
	UpdatedAt         time.Time          `json:"updatedAt"`              // When the phase was last updated
}

// PhaseApplication represents an application for a qualification requiring expert review
type PhaseApplication struct {
	ID                int64     `json:"id"`                       // Primary key identifier
	PhaseID           int64     `json:"phaseId"`                  // Foreign key reference to phases table
	Type              string    `json:"type"`                     // Type: "QP" (Qualification Placement) or "IL" (Institutional Listing)
	InstitutionName   string    `json:"institutionName"`          // Name of the institution
	QualificationName string    `json:"qualificationName"`        // Name of the qualification being reviewed
	Expert1           int64     `json:"expert1"`                  // First expert ID
	Expert1Name       string    `json:"expert1Name,omitempty"`    // First expert name (not stored in DB)
	Expert2           int64     `json:"expert2"`                  // Second expert ID
	Expert2Name       string    `json:"expert2Name,omitempty"`    // Second expert name (not stored in DB)
	Status            string    `json:"status"`                   // Status: "pending", "assigned", "approved", "rejected"
	RejectionNotes    string    `json:"rejectionNotes,omitempty"` // Notes for rejection (if status is "rejected")
	CreatedAt         time.Time `json:"createdAt"`                // When the application was created
	UpdatedAt         time.Time `json:"updatedAt"`                // When the application was last updated
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
	ID      int64  `json:"id"`                // ID of the newly created user
	Success bool   `json:"success"`           // Indicates if the creation was successful
	Message string `json:"message,omitempty"` // Optional message providing additional details
}
