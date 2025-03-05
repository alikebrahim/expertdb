package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type CreateExpertRequest struct {
	Name           string   `json:"name"`
	Affiliation    string   `json:"affiliation"`
	PrimaryContact string   `json:"primaryContact"`
	ContactType    string   `json:"contactType"` // "email" or "phone"
	Skills         []string `json:"skills"`
	Biography      string   `json:"biography"`
	Availability   string   `json:"availability"` // "full-time", "part-time", "weekends"
}

type CreateExpertResponse struct {
	ID      int64  `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ISCEDLevel represents an educational level according to ISCED classification
type ISCEDLevel struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ISCEDField represents a field of education according to ISCED classification
type ISCEDField struct {
	ID          int64  `json:"id"`
	BroadCode   string `json:"broadCode"`
	BroadName   string `json:"broadName"`
	NarrowCode  string `json:"narrowCode,omitempty"`
	NarrowName  string `json:"narrowName,omitempty"`
	DetailedCode string `json:"detailedCode,omitempty"`
	DetailedName string `json:"detailedName,omitempty"`
	Description string `json:"description,omitempty"`
}

type Expert struct {
	ID              int64       `json:"id"`
	ExpertID        string      `json:"expertId"`
	Name            string      `json:"name"`
	Designation     string      `json:"designation"`
	Institution     string      `json:"institution"`
	IsBahraini      bool        `json:"isBahraini"`
	Nationality     string      `json:"nationality"`
	IsAvailable     bool        `json:"isAvailable"`
	Rating          string      `json:"rating"`
	Role            string      `json:"role"`
	EmploymentType  string      `json:"employmentType"`
	GeneralArea     string      `json:"generalArea"`
	SpecializedArea string      `json:"specializedArea"`
	IsTrained       bool        `json:"isTrained"`
	CVPath          string      `json:"cvPath"`
	Phone           string      `json:"phone"`
	Email           string      `json:"email"`
	IsPublished     bool        `json:"isPublished"`
	ISCEDLevel      *ISCEDLevel `json:"iscedLevel,omitempty"`
	ISCEDField      *ISCEDField `json:"iscedField,omitempty"`
	Areas           []Area      `json:"areas,omitempty"`
	Documents       []Document  `json:"documents,omitempty"`
	Engagements     []Engagement `json:"engagements,omitempty"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt,omitempty"`
}

// Area represents an expert specialization area
type Area struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ExpertRequest represents a request to add a new expert
type ExpertRequest struct {
	ID              int64     `json:"id"`
	ExpertID        string    `json:"expertId,omitempty"`
	Name            string    `json:"name"`
	Designation     string    `json:"designation"`
	Institution     string    `json:"institution"`
	IsBahraini      bool      `json:"isBahraini"`
	IsAvailable     bool      `json:"isAvailable"`
	Rating          string    `json:"rating"`
	Role            string    `json:"role"`
	EmploymentType  string    `json:"employmentType"`
	GeneralArea     string    `json:"generalArea"`
	SpecializedArea string    `json:"specializedArea"`
	IsTrained       bool      `json:"isTrained"`
	CVPath          string    `json:"cvPath"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	IsPublished     bool      `json:"isPublished"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	ReviewedAt      time.Time `json:"reviewedAt,omitempty"`
	ReviewedBy      int64     `json:"reviewedBy,omitempty"`
}

// User represents a system user
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	Role         string    `json:"role"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	LastLogin    time.Time `json:"lastLogin,omitempty"`
}

// Document represents an uploaded document for an expert
type Document struct {
	ID           int64     `json:"id"`
	ExpertID     int64     `json:"expertId"`
	DocumentType string    `json:"documentType"` // "cv", "certificate", "publication", etc.
	Filename     string    `json:"filename"`
	FilePath     string    `json:"filePath"`
	ContentType  string    `json:"contentType"` // MIME type
	FileSize     int64     `json:"fileSize"`    // In bytes
	UploadDate   time.Time `json:"uploadDate"`
}

// Engagement represents expert assignment to projects/activities
type Engagement struct {
	ID             int64     `json:"id"`
	ExpertID       int64     `json:"expertId"`
	EngagementType string    `json:"engagementType"` // "evaluation", "consultation", "project", etc.
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate,omitempty"`
	ProjectName    string    `json:"projectName,omitempty"`
	Status         string    `json:"status"` // "pending", "active", "completed", "cancelled"
	FeedbackScore  int       `json:"feedbackScore,omitempty"` // 1-5 rating
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

// AIAnalysisResult represents AI-generated content/analysis
type AIAnalysisResult struct {
	ID              int64     `json:"id"`
	ExpertID        int64     `json:"expertId,omitempty"`
	DocumentID      int64     `json:"documentId,omitempty"`
	AnalysisType    string    `json:"analysisType"` // "profile", "isced_suggestion", "skills_extraction"
	AnalysisResult  string    `json:"analysisResult"`   // JSON or text data
	ResultData      string    `json:"resultData"`   // JSON or text data (alias for AnalysisResult)
	ConfidenceScore float64   `json:"confidenceScore,omitempty"`
	ModelUsed       string    `json:"modelUsed,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// Statistics represents system-wide statistics
type Statistics struct {
	TotalExperts          int       `json:"totalExperts"`
	BahrainiPercentage    float64   `json:"bahrainiPercentage"`
	TopAreas              []AreaStat `json:"topAreas"`
	ExpertsByISCEDField   []AreaStat `json:"expertsByISCEDField"`
	EngagementsByType     []AreaStat `json:"engagementsByType"`
	MonthlyGrowth         []GrowthStat `json:"monthlyGrowth"`
	MostRequestedExperts  []ExpertStat `json:"mostRequestedExperts"`
	LastUpdated           time.Time    `json:"lastUpdated"`
}

// AreaStat represents statistics for a specific area/category
type AreaStat struct {
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// GrowthStat represents growth statistics over time
type GrowthStat struct {
	Period     string  `json:"period"` // "2023-01", "2023-Q1", etc.
	Count      int     `json:"count"`
	GrowthRate float64 `json:"growthRate"` // Percentage growth from previous period
}

// ExpertStat represents statistics for a specific expert
type ExpertStat struct {
	ExpertID string `json:"expertId"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

// DocumentUploadRequest represents a request to upload a document
type DocumentUploadRequest struct {
	ExpertID     int64  `json:"expertId"`
	DocumentType string `json:"documentType"`
}

// AIAnalysisRequest represents a request for AI analysis
type AIAnalysisRequest struct {
	ExpertID    int64  `json:"expertId,omitempty"`
	DocumentID  int64  `json:"documentId,omitempty"`
	AnalysisType string `json:"analysisType"`
	InputData   string `json:"inputData,omitempty"` // Optional additional input data
}

// Authentication types

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a user login response
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}

// CreateUserResponse represents a response to creating a new user
type CreateUserResponse struct {
	ID      int64  `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// Configuration represents application configuration
type Configuration struct {
	Port             string `json:"port"`
	DBPath           string `json:"dbPath"`
	UploadPath       string `json:"uploadPath"`
	CORSAllowOrigins string `json:"corsAllowOrigins"`
	AIServiceURL     string `json:"aiServiceURL"`
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
		Name:         req.Name,
		Institution:  req.Affiliation,
		IsAvailable:  req.Availability == "yes" || req.Availability == "full-time",
		Email:        email,
		Phone:        phone,
		CreatedAt:    time.Now().UTC(),
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

	// Limit biography length
	if len(req.Biography) > 1000 {
		return errors.New("biography exceeds maximum length of 1000 characters")
	}
	
	return nil
}
