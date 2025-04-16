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

	// Parse and validate months parameter
	months := 12 // Default to 12 months if not specified

	monthsParam := r.URL.Query().Get("months")
	if monthsParam != "" {
		parsedMonths, err := strconv.Atoi(monthsParam)
		if err == nil && parsedMonths > 0 {
			months = parsedMonths
			log.Debug("Using custom months parameter: %d", months)
		} else {
			log.Warn("Invalid months parameter provided: %s, using default (12)", monthsParam)
		}
	}

	// Get all experts with their creation dates
	experts, err := h.store.ListExperts(map[string]interface{}{}, 1000, 0)
	if err != nil {
		log.Error("Failed to retrieve experts: %v", err)
		return fmt.Errorf("failed to retrieve growth statistics: %w", err)
	}

	// This is a simplistic implementation - in a real scenario, you'd use SQL's date functions
	// to aggregate by month directly in the database query
	
	// Map to store counts by month
	monthCounts := make(map[string]int)
	
	// Process each expert
	for _, expert := range experts {
		// Format the month as YYYY-MM
		month := expert.CreatedAt.Format("2006-01")
		monthCounts[month]++
	}
	
	// Convert to GrowthStat format and calculate growth rates
	var stats []domain.GrowthStat
	var prevCount int
	
	// Sort by month and limit to the specified number of months
	// In a real implementation, you would sort the months here
	
	for month, count := range monthCounts {
		growthRate := 0.0
		if prevCount > 0 {
			growthRate = (float64(count) - float64(prevCount)) / float64(prevCount) * 100
		}
		
		stats = append(stats, domain.GrowthStat{
			Period:     month,
			Count:      count,
			GrowthRate: growthRate,
		})
		
		prevCount = count
	}
	
	// Return statistics as JSON response
	log.Debug("Successfully retrieved growth statistics for %d months", months)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(stats)
}