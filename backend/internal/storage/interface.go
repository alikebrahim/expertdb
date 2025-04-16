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
	GenerateUniqueExpertID() (string, error)
	ExpertIDExists(expertID string) (bool, error)
	
	// Expert request methods
	ListExpertRequests(status string, limit, offset int) ([]*domain.ExpertRequest, error)
	GetExpertRequest(id int64) (*domain.ExpertRequest, error)
	CreateExpertRequest(req *domain.ExpertRequest) (int64, error)
	UpdateExpertRequestStatus(id int64, status, rejectionReason string, reviewedBy int64) error
	UpdateExpertRequest(req *domain.ExpertRequest) error
	
	// User methods
	GetUser(id int64) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	CreateUser(user *domain.User) (int64, error)
	UpdateUser(user *domain.User) error
	ListUsers(limit, offset int) ([]*domain.User, error)
	UpdateUserLastLogin(id int64) error
	EnsureAdminExists(adminEmail, adminName, adminPasswordHash string) error
	
	// Area methods
	ListAreas() ([]*domain.Area, error)
	GetArea(id int64) (*domain.Area, error)
	
	// Document methods
	ListDocuments(expertID int64) ([]*domain.Document, error)
	GetDocument(id int64) (*domain.Document, error)
	CreateDocument(doc *domain.Document) (int64, error)
	DeleteDocument(id int64) error
	
	// Engagement methods
	ListEngagements(expertID int64) ([]*domain.Engagement, error)
	GetEngagement(id int64) (*domain.Engagement, error)
	CreateEngagement(engagement *domain.Engagement) (int64, error)
	UpdateEngagement(engagement *domain.Engagement) error
	DeleteEngagement(id int64) error
	
	// Statistics methods
	GetStatistics() (*domain.Statistics, error)
	UpdateStatistics(stats *domain.Statistics) error
	GetExpertsByNationality() (int, int, error)
	GetEngagementStatistics() ([]domain.AreaStat, error)
	GetExpertGrowthByMonth(months int) ([]domain.GrowthStat, error)
	
	// General database methods
	InitDB() error
	Close() error
}