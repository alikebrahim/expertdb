// Package statistics provides handlers for statistics-related API endpoints
package statistics

import (
	"fmt"
	"net/http"
	"strconv"

	"expertdb/internal/api/response"
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

	// Parse and validate years parameter
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

	// Retrieve comprehensive statistics from database
	log.Debug("Retrieving comprehensive system statistics for %d years", years)
	stats, err := h.store.GetStatistics(years)
	if err != nil {
		log.Error("Failed to retrieve statistics: %v", err)
		return fmt.Errorf("failed to retrieve statistics: %w", err)
	}

	// Return statistics with standardized response
	log.Debug("Successfully retrieved comprehensive system statistics")
	return response.Success(w, http.StatusOK, "", stats)
}

