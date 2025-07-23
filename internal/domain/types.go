// Package domain contains the core business entities for the ExpertDB application
package domain

import (
	"errors"
	"regexp"
	"strconv"
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
	Name            string     `json:"name"`            // Full name of the expert
	Designation     string     `json:"designation"`     // Professional title: Prof., Dr., Mr., Ms., Mrs., Miss, Eng.
	Affiliation     string     `json:"affiliation"`     // Organization or institution the expert is affiliated with
	Phone           string     `json:"phone"`           // Contact phone number
	Email           string     `json:"email"`           // Contact email address
	Role            string     `json:"role"`            // Expert's role: "evaluator", "validator", or "evaluator/validator"
	EmploymentType  string     `json:"employmentType"`  // Type of employment: "academic" or "employer"
	GeneralArea               int64    `json:"generalArea"`               // ID referencing expert_areas table
	SpecializedAreaIds        []int64  `json:"specializedAreaIds"`        // Selected existing area IDs
	SuggestedSpecializedAreas []string `json:"suggestedSpecializedAreas"` // User-suggested area names
	CVPath                    string   `json:"cvPath"`                    // Path to the expert's CV file
	ExperienceEntries []ExpertRequestExperienceEntry `json:"experienceEntries"` // Professional experience entries
	EducationEntries  []ExpertRequestEducationEntry  `json:"educationEntries"`  // Educational background entries
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
	Name                string    `json:"name"`            // Full name of the expert
	Designation         string    `json:"designation"`     // Professional title or position
	Affiliation         string    `json:"affiliation"`     // Organization or institution affiliation
	IsBahraini          bool      `json:"isBahraini"`      // Flag indicating if expert is Bahraini citizen
	IsAvailable         bool      `json:"isAvailable"`     // Current availability status for assignments
	Rating              int       `json:"rating"`          // Performance rating (1-5 scale)
	Role                string    `json:"role"`            // Expert's role: "evaluator", "validator", or "evaluator/validator"
	EmploymentType      string    `json:"employmentType"`  // Type of employment: "academic" or "employer"
	GeneralArea         int64     `json:"-" db:"general_area"`     // ID referencing expert_areas table - internal use only
	GeneralAreaName     string    `json:"generalAreaName"` // Name of the general area (from expert_areas table)
	SpecializedArea     string    `json:"-" db:"specialized_area"` // Comma-separated specialized area IDs (e.g., "1,4,6") - internal use only
	SpecializedAreaNames string   `json:"specializedAreaNames"` // Comma-separated specialized area names (e.g., "Software Engineering, Database Design")
	SpecializedAreasResolved []*SpecializedArea `json:"specialized_areas_resolved,omitempty" db:"-"` // Resolved specialized area names
	IsTrained           bool      `json:"isTrained"`       // Indicates if expert has completed required training
	CVPath              string    `json:"cvPath"`          // Path to the expert's CV file
	ApprovalDocumentPath string    `json:"approvalDocumentPath,omitempty"` // Path to the approval document
	Phone               string    `json:"phone"`           // Contact phone number
	Email               string    `json:"email"`           // Contact email address
	IsPublished         bool      `json:"isPublished"`     // Indicates if expert profile should be publicly visible
	ExperienceEntries   []ExpertExperienceEntry `json:"experienceEntries,omitempty"` // Professional experience entries
	EducationEntries    []ExpertEducationEntry  `json:"educationEntries,omitempty"`  // Educational background entries
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

// SpecializedArea represents a specialized area for experts
type SpecializedArea struct {
	ID        int64     `json:"id" db:"id"`               // Unique identifier for the specialized area
	Name      string    `json:"name" db:"name"`           // Name of the specialized area
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Timestamp when specialized area was created
}

// ExpertRequest represents a request to add a new expert
type ExpertRequest struct {
	ID                   int64     `json:"id"`                   // Primary key identifier
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
	GeneralArea               int64    `json:"generalArea"`               // ID referencing expert_areas table
	SpecializedArea           string   `json:"specializedArea"`           // Comma-separated existing area IDs
	SuggestedSpecializedAreas []string `json:"suggestedSpecializedAreas"` // User-suggested area names
	IsTrained                 bool     `json:"isTrained"`                 // Indicates if expert has completed required training
	IsPublished          bool      `json:"isPublished"`          // Indicates if expert profile should be publicly visible
	CVPath               string    `json:"cvPath"`               // Path to the expert's CV file
	ApprovalDocumentPath string    `json:"approvalDocumentPath,omitempty"` // Path to the approval document
	ExperienceEntries    []ExpertRequestExperienceEntry `json:"experienceEntries,omitempty"` // Professional experience entries
	EducationEntries     []ExpertRequestEducationEntry  `json:"educationEntries,omitempty"`  // Educational background entries
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
	EngagementsByStatus  []AreaStat   `json:"engagementsByStatus"`  // Distribution of engagements by status (active/completed)
	NationalityStats     []AreaStat   `json:"nationalityStats"`     // Bahraini vs Non-Bahraini breakdown with counts
	SpecializedAreas     AreaBreakdown `json:"specializedAreas"`    // Top and bottom specialized areas
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

// Expert Edit Request Types

// ExpertEditRequestExperienceEntry represents experience changes in an edit request
type ExpertEditRequestExperienceEntry struct {
	ID              int64  `json:"id" db:"id"`
	EditRequestID   int64  `json:"editRequestId" db:"expert_edit_request_id"`
	Action          string `json:"action" db:"action"` // "add", "update", "delete"
	ExperienceID    int64  `json:"experienceId,omitempty" db:"experience_id"` // For update/delete
	Organization    string `json:"organization" db:"organization"`
	Position        string `json:"position" db:"position"`
	StartDate       string `json:"startDate" db:"start_date"`
	EndDate         string `json:"endDate" db:"end_date"`
	IsCurrent       bool   `json:"isCurrent" db:"is_current"`
	Country         string `json:"country" db:"country"`
	Description     string `json:"description" db:"description"`
}

// ExpertEditRequestEducationEntry represents education changes in an edit request
type ExpertEditRequestEducationEntry struct {
	ID              int64  `json:"id" db:"id"`
	EditRequestID   int64  `json:"editRequestId" db:"expert_edit_request_id"`
	Action          string `json:"action" db:"action"` // "add", "update", "delete"
	EducationID     int64  `json:"educationId,omitempty" db:"education_id"` // For update/delete
	Institution     string `json:"institution" db:"institution"`
	Degree          string `json:"degree" db:"degree"`
	FieldOfStudy    string `json:"fieldOfStudy" db:"field_of_study"`
	GraduationYear  string `json:"graduationYear" db:"graduation_year"`
	Country         string `json:"country" db:"country"`
	Description     string `json:"description" db:"description"`
}

// ExpertEditRequest represents a request to edit an existing expert profile
type ExpertEditRequest struct {
	ID                        int64     `json:"id"`                        // Primary key identifier
	ExpertID                  int64     `json:"expertId"`                  // References experts(id) - the expert being edited
	
	// Core profile fields (nil/null = no change proposed)
	Name                     *string   `json:"name,omitempty"`                     // Full name of the expert
	Designation              *string   `json:"designation,omitempty"`              // Professional title
	Institution              *string   `json:"institution,omitempty"`              // Organization or institution affiliation
	Phone                    *string   `json:"phone,omitempty"`                    // Contact phone number
	Email                    *string   `json:"email,omitempty"`                    // Contact email address
	IsBahraini               *bool     `json:"isBahraini,omitempty"`               // Flag indicating if expert is Bahraini citizen
	IsAvailable              *bool     `json:"isAvailable,omitempty"`              // Current availability status for assignments
	Rating                   *int      `json:"rating,omitempty"`                   // Performance rating (1-5 scale)
	Role                     *string   `json:"role,omitempty"`                     // Expert's role
	EmploymentType           *string   `json:"employmentType,omitempty"`           // Type of employment
	GeneralArea              *int64    `json:"generalArea,omitempty"`              // Reference to expert_areas table
	SpecializedArea          *string   `json:"specializedArea,omitempty"`          // Comma-separated specialized area IDs
	IsTrained                *bool     `json:"isTrained,omitempty"`                // Training status
	IsPublished              *bool     `json:"isPublished,omitempty"`              // Publication status
	Biography                *string   `json:"biography,omitempty"`                // Extended profile information
	SuggestedSpecializedAreas []string `json:"suggestedSpecializedAreas,omitempty"` // User-suggested area names
	
	// Document updates (nil/null = no change proposed)
	NewCVPath                *string   `json:"newCvPath,omitempty"`                // Path to updated CV file
	NewApprovalDocumentPath  *string   `json:"newApprovalDocumentPath,omitempty"`  // Path to updated approval document
	RemoveCV                 bool      `json:"removeCv"`                           // Flag to indicate CV should be removed
	RemoveApprovalDocument   bool      `json:"removeApprovalDocument"`             // Flag to indicate approval document should be removed
	
	// Experience and education changes
	ExperienceChanges        []ExpertEditRequestExperienceEntry `json:"experienceChanges,omitempty"` // Experience modifications
	EducationChanges         []ExpertEditRequestEducationEntry  `json:"educationChanges,omitempty"`  // Education modifications
	
	// Change metadata
	ChangeSummary            string    `json:"changeSummary"`            // User-provided summary of changes
	ChangeReason             string    `json:"changeReason"`             // Reason for requesting the edit
	FieldsChanged            []string  `json:"fieldsChanged"`            // Array of field names being changed
	
	// Status and workflow
	Status                   string    `json:"status"`                   // "pending", "approved", "rejected", "cancelled"
	RejectionReason          string    `json:"rejectionReason,omitempty"` // Reason for rejection if status is 'rejected'
	AdminNotes               string    `json:"adminNotes,omitempty"`     // Internal notes for admin review
	
	// Audit trail
	CreatedAt                time.Time `json:"createdAt"`                // Timestamp when request was created
	ReviewedAt               time.Time `json:"reviewedAt,omitempty"`     // Timestamp when request was reviewed
	AppliedAt                time.Time `json:"appliedAt,omitempty"`      // When the changes were applied to the expert
	CreatedBy                int64     `json:"createdBy"`                // References users(id) - who requested the edit
	ReviewedBy               int64     `json:"reviewedBy,omitempty"`     // References users(id) - who reviewed the request
	
	// Computed fields (not stored in database)
	ExpertName               string    `json:"expertName,omitempty"`     // Name of the expert being edited (for display)
	CreatedByName            string    `json:"createdByName,omitempty"`  // Name of user who created request
	ReviewedByName           string    `json:"reviewedByName,omitempty"` // Name of user who reviewed request
}

// CreateExpertEditRequest represents a request to create an expert edit request
type CreateExpertEditRequest struct {
	ExpertID                  int64     `json:"expertId"`                  // References experts(id) - the expert being edited
	
	// Core profile fields (nil/null = no change proposed)
	Name                     *string   `json:"name,omitempty"`                     // Full name of the expert
	Designation              *string   `json:"designation,omitempty"`              // Professional title
	Institution              *string   `json:"institution,omitempty"`              // Organization or institution affiliation
	Phone                    *string   `json:"phone,omitempty"`                    // Contact phone number
	Email                    *string   `json:"email,omitempty"`                    // Contact email address
	IsBahraini               *bool     `json:"isBahraini,omitempty"`               // Flag indicating if expert is Bahraini citizen
	IsAvailable              *bool     `json:"isAvailable,omitempty"`              // Current availability status for assignments
	Rating                   *int      `json:"rating,omitempty"`                   // Performance rating (1-5 scale)
	Role                     *string   `json:"role,omitempty"`                     // Expert's role
	EmploymentType           *string   `json:"employmentType,omitempty"`           // Type of employment
	GeneralArea              *int64    `json:"generalArea,omitempty"`              // Reference to expert_areas table
	SpecializedAreaIds       []int64   `json:"specializedAreaIds,omitempty"`       // Selected existing area IDs
	IsTrained                *bool     `json:"isTrained,omitempty"`                // Training status
	IsPublished              *bool     `json:"isPublished,omitempty"`              // Publication status
	Biography                *string   `json:"biography,omitempty"`                // Extended profile information
	SuggestedSpecializedAreas []string `json:"suggestedSpecializedAreas,omitempty"` // User-suggested area names
	
	// Document updates
	RemoveCV                 bool      `json:"removeCv"`                           // Flag to indicate CV should be removed
	RemoveApprovalDocument   bool      `json:"removeApprovalDocument"`             // Flag to indicate approval document should be removed
	
	// Experience and education changes
	ExperienceChanges        []ExpertEditRequestExperienceEntry `json:"experienceChanges,omitempty"` // Experience modifications
	EducationChanges         []ExpertEditRequestEducationEntry  `json:"educationChanges,omitempty"`  // Education modifications
	
	// Change metadata - required
	ChangeSummary            string    `json:"changeSummary"`            // User-provided summary of changes
	ChangeReason             string    `json:"changeReason"`             // Reason for requesting the edit
}

// NewExpert creates a new Expert from a CreateExpertRequest
func NewExpert(req CreateExpertRequest) *Expert {
	// Convert SpecializedAreaIds to comma-separated string for storage
	specializedAreaStr := ""
	if len(req.SpecializedAreaIds) > 0 {
		idStrings := make([]string, len(req.SpecializedAreaIds))
		for i, id := range req.SpecializedAreaIds {
			idStrings[i] = strconv.FormatInt(id, 10)
		}
		specializedAreaStr = strings.Join(idStrings, ",")
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
		SpecializedArea: specializedAreaStr,
		CVPath:          req.CVPath,
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

	if len(req.SpecializedAreaIds) == 0 && len(req.SuggestedSpecializedAreas) == 0 {
		return errors.New("at least one specialized area (existing or suggested) is required")
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

	// Validate rating (0-5 scale, where 0 means no rating)
	if req.Rating < 0 || req.Rating > 5 {
		return errors.New("rating must be between 0 and 5 (0 means no rating)")
	}


	// Validate experience and education entries
	if len(req.ExperienceEntries) == 0 && len(req.EducationEntries) == 0 {
		return errors.New("at least one experience or education entry is required")
	}

	// Validate experience entries
	for i, exp := range req.ExperienceEntries {
		if strings.TrimSpace(exp.Position) == "" {
			return errors.New("experience position is required for entry " + strconv.Itoa(i+1))
		}
		if strings.TrimSpace(exp.Organization) == "" {
			return errors.New("experience organization is required for entry " + strconv.Itoa(i+1))
		}
		if strings.TrimSpace(exp.StartDate) == "" {
			return errors.New("experience start date is required for entry " + strconv.Itoa(i+1))
		}
	}

	// Validate education entries
	for i, edu := range req.EducationEntries {
		if strings.TrimSpace(edu.Degree) == "" {
			return errors.New("education degree is required for entry " + strconv.Itoa(i+1))
		}
		if strings.TrimSpace(edu.Institution) == "" {
			return errors.New("education institution is required for entry " + strconv.Itoa(i+1))
		}
		if strings.TrimSpace(edu.GraduationYear) == "" {
			return errors.New("education graduation year is required for entry " + strconv.Itoa(i+1))
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

// ValidateCreateExpertEditRequest validates the expert edit request fields
func ValidateCreateExpertEditRequest(req *CreateExpertEditRequest) error {
	// Required fields
	if req.ExpertID == 0 {
		return errors.New("expert ID is required")
	}
	
	if strings.TrimSpace(req.ChangeSummary) == "" {
		return errors.New("change summary is required")
	}
	
	if strings.TrimSpace(req.ChangeReason) == "" {
		return errors.New("change reason is required")
	}
	
	// At least one field must be changed
	hasChanges := false
	
	// Check if any core fields are being changed
	if req.Name != nil || req.Designation != nil || req.Institution != nil ||
	   req.Phone != nil || req.Email != nil || req.IsBahraini != nil ||
	   req.IsAvailable != nil || req.Rating != nil || req.Role != nil ||
	   req.EmploymentType != nil || req.GeneralArea != nil ||
	   req.IsTrained != nil || req.IsPublished != nil || req.Biography != nil {
		hasChanges = true
	}
	
	// Check specialized areas
	if len(req.SpecializedAreaIds) > 0 || len(req.SuggestedSpecializedAreas) > 0 {
		hasChanges = true
	}
	
	// Check document changes
	if req.RemoveCV || req.RemoveApprovalDocument {
		hasChanges = true
	}
	
	// Check experience/education changes
	if len(req.ExperienceChanges) > 0 || len(req.EducationChanges) > 0 {
		hasChanges = true
	}
	
	if !hasChanges {
		return errors.New("at least one field must be changed")
	}
	
	// Validate changed fields
	if req.Email != nil {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(*req.Email) {
			return errors.New("invalid email format")
		}
	}
	
	if req.Phone != nil {
		phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
		if !phoneRegex.MatchString(*req.Phone) {
			return errors.New("invalid phone number format")
		}
	}
	
	if req.Designation != nil {
		validDesignations := []string{"Prof.", "Dr.", "Mr.", "Ms.", "Mrs.", "Miss", "Eng."}
		if !containsString(validDesignations, *req.Designation) {
			return errors.New("designation must be one of: Prof., Dr., Mr., Ms., Mrs., Miss, Eng.")
		}
	}
	
	if req.Role != nil {
		validRoles := []string{"evaluator", "validator", "evaluator/validator"}
		if !containsString(validRoles, strings.ToLower(*req.Role)) {
			return errors.New("role must be one of: evaluator, validator, evaluator/validator")
		}
	}
	
	if req.EmploymentType != nil {
		validEmploymentTypes := []string{"academic", "employer"}
		if !containsString(validEmploymentTypes, strings.ToLower(*req.EmploymentType)) {
			return errors.New("employment type must be one of: academic, employer")
		}
	}
	
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		return errors.New("rating must be between 1 and 5")
	}
	
	// Validate experience changes
	for i, exp := range req.ExperienceChanges {
		if exp.Action == "" {
			return errors.New("experience action is required for entry " + strconv.Itoa(i+1))
		}
		
		validActions := []string{"add", "update", "delete"}
		if !containsString(validActions, exp.Action) {
			return errors.New("experience action must be one of: add, update, delete for entry " + strconv.Itoa(i+1))
		}
		
		if exp.Action == "update" || exp.Action == "delete" {
			if exp.ExperienceID == 0 {
				return errors.New("experience ID is required for update/delete action in entry " + strconv.Itoa(i+1))
			}
		}
		
		if exp.Action == "add" || exp.Action == "update" {
			if strings.TrimSpace(exp.Position) == "" {
				return errors.New("experience position is required for entry " + strconv.Itoa(i+1))
			}
			if strings.TrimSpace(exp.Organization) == "" {
				return errors.New("experience organization is required for entry " + strconv.Itoa(i+1))
			}
			if strings.TrimSpace(exp.StartDate) == "" {
				return errors.New("experience start date is required for entry " + strconv.Itoa(i+1))
			}
		}
	}
	
	// Validate education changes
	for i, edu := range req.EducationChanges {
		if edu.Action == "" {
			return errors.New("education action is required for entry " + strconv.Itoa(i+1))
		}
		
		validActions := []string{"add", "update", "delete"}
		if !containsString(validActions, edu.Action) {
			return errors.New("education action must be one of: add, update, delete for entry " + strconv.Itoa(i+1))
		}
		
		if edu.Action == "update" || edu.Action == "delete" {
			if edu.EducationID == 0 {
				return errors.New("education ID is required for update/delete action in entry " + strconv.Itoa(i+1))
			}
		}
		
		if edu.Action == "add" || edu.Action == "update" {
			if strings.TrimSpace(edu.Degree) == "" {
				return errors.New("education degree is required for entry " + strconv.Itoa(i+1))
			}
			if strings.TrimSpace(edu.Institution) == "" {
				return errors.New("education institution is required for entry " + strconv.Itoa(i+1))
			}
			if strings.TrimSpace(edu.GraduationYear) == "" {
				return errors.New("education graduation year is required for entry " + strconv.Itoa(i+1))
			}
		}
	}
	
	return nil
}