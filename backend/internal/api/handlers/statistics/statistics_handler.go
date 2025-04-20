// Package statistics provides handlers for statistics-related API endpoints
package statistics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// Handler manages statistics-related HTTP endpoints
type Handler struct {
	store storage.Storage
}

// NewHandler creates a new statistics handler
func NewHandler(store storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}

// HandleGetStatistics handles GET /api/statistics requests
func (h *Handler) HandleGetStatistics(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/statistics request")

	// Retrieve all statistics from database
	log.Debug("Retrieving overall system statistics")
	stats, err := h.store.GetStatistics()
	if err != nil {
		log.Error("Failed to retrieve statistics: %v", err)
		return fmt.Errorf("failed to retrieve statistics: %w", err)
	}

	// Return statistics as JSON response
	log.Debug("Successfully retrieved system statistics")
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(stats)
}

// HandleGetNationalityStats handles GET /api/statistics/nationality requests
func (h *Handler) HandleGetNationalityStats(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/statistics/nationality request")

	// Query for experts with Bahraini status
	filters := make(map[string]interface{})
	filters["isBahraini"] = true
	bahrainiCount, err := h.store.CountExperts(filters)
	if err != nil {
		log.Error("Failed to count Bahraini experts: %v", err)
		return fmt.Errorf("failed to count Bahraini experts: %w", err)
	}

	// Get total count
	totalCount, err := h.store.CountExperts(map[string]interface{}{})
	if err != nil {
		log.Error("Failed to count total experts: %v", err)
		return fmt.Errorf("failed to count total experts: %w", err)
	}

	// Calculate non-Bahraini count
	nonBahrainiCount := totalCount - bahrainiCount

	// Calculate percentages, avoiding division by zero
	var bahrainiPercentage, nonBahrainiPercentage float64
	if totalCount > 0 {
		bahrainiPercentage = float64(bahrainiCount) / float64(totalCount) * 100
		nonBahrainiPercentage = float64(nonBahrainiCount) / float64(totalCount) * 100
	}

	// Create stats array in the format expected by the frontend
	stats := []domain.AreaStat{
		{Name: "Bahraini", Count: bahrainiCount, Percentage: bahrainiPercentage},
		{Name: "Non-Bahraini", Count: nonBahrainiCount, Percentage: nonBahrainiPercentage},
	}

	// Prepare response with total and detailed stats
	result := map[string]interface{}{
		"total": totalCount,
		"stats": stats,
	}

	// Return statistics as JSON response
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(result)
}

// HandleGetEngagementStats handles GET /api/statistics/engagements requests
func (h *Handler) HandleGetEngagementStats(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/statistics/engagements request")

	// Get all engagements (we'll process them in memory since SQLite doesn't support complex aggregations)
	allEngagements, err := h.store.ListEngagements(0) // 0 means all engagements
	if err != nil {
		log.Error("Failed to retrieve engagements: %v", err)
		return fmt.Errorf("failed to retrieve engagement statistics: %w", err)
	}

	// Count by engagement type
	typeCount := make(map[string]int)
	statusCount := make(map[string]int)
	total := len(allEngagements)

	for _, engagement := range allEngagements {
		// Count by type
		typeCount[engagement.EngagementType]++
		
		// Count by status
		statusCount[engagement.Status]++
	}

	// Convert to the AreaStat format
	var typeStats []domain.AreaStat
	for typeName, count := range typeCount {
		percentage := 0.0
		if total > 0 {
			percentage = float64(count) / float64(total) * 100
		}
		typeStats = append(typeStats, domain.AreaStat{
			Name:       typeName,
			Count:      count,
			Percentage: percentage,
		})
	}

	var statusStats []domain.AreaStat
	for statusName, count := range statusCount {
		percentage := 0.0
		if total > 0 {
			percentage = float64(count) / float64(total) * 100
		}
		statusStats = append(statusStats, domain.AreaStat{
			Name:       statusName,
			Count:      count,
			Percentage: percentage,
		})
	}

	// Prepare the response
	result := map[string]interface{}{
		"total":    total,
		"byType":   typeStats,
		"byStatus": statusStats,
	}

	// Return statistics as JSON response
	log.Debug("Successfully retrieved engagement statistics")
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(result)
}

// HandleGetGrowthStats handles GET /api/statistics/growth requests
func (h *Handler) HandleGetGrowthStats(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/statistics/growth request")

	// Parse and validate years parameter (changed from months)
	years := 5 // Default to 5 years if not specified

	yearsParam := r.URL.Query().Get("years")
	if yearsParam != "" {
		parsedYears, err := strconv.Atoi(yearsParam)
		if err == nil && parsedYears > 0 {
			years = parsedYears
			log.Debug("Using custom years parameter: %d", years)
		} else {
			log.Warn("Invalid years parameter provided: %s, using default (5)", yearsParam)
		}
	}

	// Get yearly growth statistics using the repository method
	stats, err := h.store.GetExpertGrowthByYear(years)
	if err != nil {
		log.Error("Failed to retrieve growth statistics: %v", err)
		return fmt.Errorf("failed to retrieve growth statistics: %w", err)
	}
	
	// Return statistics as JSON response
	log.Debug("Successfully retrieved growth statistics for %d years", years)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(stats)
}

// HandleGetAreaStats handles GET /api/statistics/areas requests
func (h *Handler) HandleGetAreaStats(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/statistics/areas request")

	// Retrieve area statistics from repository
	areaStats, err := h.store.GetAreaStatistics()
	if err != nil {
		log.Error("Failed to retrieve area statistics: %v", err)
		return fmt.Errorf("failed to retrieve area statistics: %w", err)
	}
	
	// Return statistics as JSON response
	log.Debug("Successfully retrieved area statistics")
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(areaStats)
}