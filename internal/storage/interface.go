// Package storage provides database access layer for the ExpertDB application
package storage

import (
	"expertdb/internal/domain"
)

// Storage defines the interface for database operations
type Storage interface {
	// Expert methods
	ListExperts(filters map[string]interface{}, limit, offset int) ([]*domain.Expert, error)
	CountExperts(filters map[string]interface{}) (int, error)
	GetExpert(id int64) (*domain.Expert, error)
	GetExpertByEmail(email string) (*domain.Expert, error)
	CreateExpert(expert *domain.Expert) (int64, error)
	UpdateExpert(expert *domain.Expert) error
	DeleteExpert(id int64) error
	
	// Expert request methods
	ListExpertRequests(status string, limit, offset int) ([]*domain.ExpertRequest, error)
	ListExpertRequestsByUser(userID int64, status string, limit, offset int) ([]*domain.ExpertRequest, error)
	GetExpertRequest(id int64) (*domain.ExpertRequest, error)
	CreateExpertRequest(req *domain.ExpertRequest) (int64, error)
	CreateExpertRequestWithoutPaths(req *domain.ExpertRequest) (int64, error)
	UpdateExpertRequestStatus(id int64, status, rejectionReason string, reviewedBy int64) error
	UpdateExpertRequest(req *domain.ExpertRequest) error
	ApproveExpertRequestWithDocument(requestID, reviewedBy int64, documentService interface{}) (int64, error)
	BatchApproveExpertRequestsWithFileMove(requestIDs []int64, reviewedBy int64, documentService interface{}) ([]int64, []int64, map[int64]error)
	TransferExpertRequestToExpert(requestID, expertID int64) error
	UpdateExpertsApprovalPath(expertIDs []int64, approvalPath string) error
	
	// Document reference methods
	UpdateExpertCVDocument(expertID, documentID int64) error
	UpdateExpertApprovalDocument(expertID, documentID int64) error
	UpdateExpertRequestCVDocument(requestID, documentID int64) error
	UpdateExpertRequestApprovalDocument(requestID, documentID int64) error
	
	// Expert edit request methods
	ListExpertEditRequests(filters map[string]interface{}, limit, offset int) ([]*domain.ExpertEditRequest, error)
	CountExpertEditRequests(filters map[string]interface{}) (int, error)
	GetExpertEditRequest(id int64) (*domain.ExpertEditRequest, error)
	CreateExpertEditRequest(req *domain.ExpertEditRequest) (int64, error)
	UpdateExpertEditRequestStatus(id int64, status, rejectionReason, adminNotes string, reviewedBy int64) error
	UpdateExpertEditRequest(req *domain.ExpertEditRequest) error
	ApplyExpertEditRequest(id int64, adminUserID int64) error
	CancelExpertEditRequest(id int64, userID int64) error
	
	// User methods
	GetUser(id int64) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	CreateUser(user *domain.User) (int64, error)
	CreateUserWithRoleCheck(user *domain.User, creatorRole string) (int64, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id int64) error
	ListUsers(limit, offset int) ([]*domain.User, error)
	UpdateUserLastLogin(id int64) error
	EnsureSuperUserExists(email, name, passwordHash string) error
	
	// Area methods
	ListAreas() ([]*domain.Area, error)
	GetArea(id int64) (*domain.Area, error)
	CreateArea(name string) (int64, error)
	UpdateArea(id int64, name string) error
	
	// Specialized area methods
	ListSpecializedAreas() ([]*domain.SpecializedArea, error)
	GetSpecializedAreasByIds(ids []int64) ([]*domain.SpecializedArea, error)
	CreateSpecializedArea(area *domain.SpecializedArea) (int64, error)
	
	// Document methods
	ListDocuments(expertID int64) ([]*domain.Document, error)
	GetDocument(id int64) (*domain.Document, error)
	CreateDocument(doc *domain.Document) (int64, error)
	UpdateDocument(doc *domain.Document) error
	DeleteDocument(id int64) error
	
	// Engagement methods
	ListEngagements(expertID int64, engagementType string, limit, offset int) ([]*domain.Engagement, error)
	GetEngagement(id int64) (*domain.Engagement, error)
	CreateEngagement(engagement *domain.Engagement) (int64, error)
	UpdateEngagement(engagement *domain.Engagement) error
	DeleteEngagement(id int64) error
	ImportEngagements(engagements []*domain.Engagement) (int, map[int]error)
	
	// Statistics methods
	GetStatistics(years int) (*domain.Statistics, error)
	UpdateStatistics(stats *domain.Statistics) error
	GetExpertsByNationality() (int, int, error)
	GetEngagementStatistics() ([]domain.AreaStat, error)
	GetExpertGrowthByYear(years int) ([]domain.GrowthStat, error)
	GetPublishedExpertStats() (int, float64, error)
	
	// Phase planning methods
	ListPhases(status string, plannerID int64, limit, offset int) ([]*domain.Phase, error)
	GetPhase(id int64) (*domain.Phase, error)
	GetPhaseByPhaseID(phaseID string) (*domain.Phase, error)
	CreatePhase(phase *domain.Phase) (int64, error)
	UpdatePhase(phase *domain.Phase) error
	GenerateUniquePhaseID() (string, error)
	
	// Phase application methods
	GetPhaseApplication(id int64) (*domain.PhaseApplication, error)
	CreatePhaseApplication(app *domain.PhaseApplication) (int64, error)
	UpdatePhaseApplication(app *domain.PhaseApplication) error
	ListPhaseApplications(phaseID int64) ([]domain.PhaseApplication, error)
	UpdatePhaseApplicationExperts(id int64, expert1ID, expert2ID int64) error
	UpdatePhaseApplicationStatus(id int64, status, rejectionNotes string) error
	
	// Role assignment methods
	IsUserPlannerForApplication(userID int, applicationID int) (bool, error)
	IsUserManagerForApplication(userID int, applicationID int) (bool, error)
	AssignUserToPlannerApplications(userID int, applicationIDs []int) error
	AssignUserToManagerApplications(userID int, applicationIDs []int) error
	RemoveUserPlannerAssignments(userID int, applicationIDs []int) error
	RemoveUserManagerAssignments(userID int, applicationIDs []int) error
	GetUserPlannerApplications(userID int) ([]int, error)
	GetUserManagerApplications(userID int) ([]int, error)
	
	// General database methods
	InitDB() error
	Close() error
}