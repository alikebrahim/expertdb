// Package engagements provides handlers for engagement-related API endpoints
package engagements

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// Handler manages engagement-related HTTP endpoints
type Handler struct {
	store storage.Storage
}

// NewHandler creates a new engagement handler
func NewHandler(store storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}

// HandleCreateEngagement handles POST /api/engagements requests
func (h *Handler) HandleCreateEngagement(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/engagements request")

	// Parse request body
	var engagement domain.Engagement
	if err := json.NewDecoder(r.Body).Decode(&engagement); err != nil {
		log.Warn("Failed to parse engagement creation request: %v", err)
		return fmt.Errorf("invalid request payload: %w", err)
	}

	// Set default values and validate
	// Set creation time
	engagement.CreatedAt = time.Now()

	// Validate required fields
	if engagement.ExpertID == 0 {
		log.Warn("Missing expert ID in engagement creation request")
		return fmt.Errorf("expert ID is required")
	}
	if engagement.EngagementType == "" {
		log.Warn("Missing engagement type in creation request")
		return fmt.Errorf("engagement type is required")
	}
	if engagement.StartDate.IsZero() {
		log.Warn("Missing start date in engagement creation request")
		return fmt.Errorf("start date is required")
	}

	// Set default status if not provided
	if engagement.Status == "" {
		log.Debug("No status specified, using default status: pending")
		engagement.Status = "pending" // Default status
	}

	// Create the engagement in database
	log.Debug("Creating engagement for expert ID: %d, type: %s",
		engagement.ExpertID, engagement.EngagementType)
	id, err := h.store.CreateEngagement(&engagement)
	if err != nil {
		log.Error("Failed to create engagement in database: %v", err)
		return fmt.Errorf("failed to create engagement: %w", err)
	}

	// Set the ID in the response and return
	engagement.ID = id
	log.Info("Engagement created successfully: ID: %d, Type: %s, Expert: %d",
		id, engagement.EngagementType, engagement.ExpertID)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(engagement)
}

// HandleGetEngagement handles GET /api/engagements/{id} requests
func (h *Handler) HandleGetEngagement(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid engagement ID provided: %s", idStr)
		return fmt.Errorf("invalid engagement ID: %w", err)
	}

	// Retrieve engagement from database
	log.Debug("Retrieving engagement with ID: %d", id)
	engagement, err := h.store.GetEngagement(id)
	if err != nil {
		log.Warn("Engagement not found with ID: %d - %v", id, err)
		return fmt.Errorf("engagement not found: %w", err)
	}

	// Return engagement data
	log.Debug("Successfully retrieved engagement: ID: %d, Type: %s", engagement.ID, engagement.EngagementType)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(engagement)
}

// HandleUpdateEngagement handles PUT /api/engagements/{id} requests
func (h *Handler) HandleUpdateEngagement(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid engagement ID provided for update: %s", idStr)
		return fmt.Errorf("invalid engagement ID: %w", err)
	}

	// Retrieve existing engagement from database
	log.Debug("Checking if engagement exists with ID: %d", id)
	existing, err := h.store.GetEngagement(id)
	if err != nil {
		log.Warn("Engagement not found for update ID: %d - %v", id, err)
		return fmt.Errorf("engagement not found: %w", err)
	}

	// Parse update request
	var updateEngagement domain.Engagement
	if err := json.NewDecoder(r.Body).Decode(&updateEngagement); err != nil {
		log.Warn("Failed to parse engagement update request: %v", err)
		return fmt.Errorf("invalid request payload: %w", err)
	}

	// Ensure ID matches path parameter
	updateEngagement.ID = id

	// Merge with existing engagement data - use existing data for empty fields
	// This approach maintains data integrity by preserving fields not included in the update
	if updateEngagement.ExpertID == 0 {
		updateEngagement.ExpertID = existing.ExpertID
	}
	if updateEngagement.EngagementType == "" {
		updateEngagement.EngagementType = existing.EngagementType
	}
	if updateEngagement.StartDate.IsZero() {
		updateEngagement.StartDate = existing.StartDate
	}
	if updateEngagement.Status == "" {
		updateEngagement.Status = existing.Status
	}
	// Preserve creation time
	if updateEngagement.CreatedAt.IsZero() {
		updateEngagement.CreatedAt = existing.CreatedAt
	}
	// Preserve additional fields if they exist
	if updateEngagement.EndDate.IsZero() && !existing.EndDate.IsZero() {
		updateEngagement.EndDate = existing.EndDate
	}
	if updateEngagement.Notes == "" && existing.Notes != "" {
		updateEngagement.Notes = existing.Notes
	}

	// Update the engagement in database
	log.Debug("Updating engagement ID: %d, Type: %s", id, updateEngagement.EngagementType)
	if err := h.store.UpdateEngagement(&updateEngagement); err != nil {
		log.Error("Failed to update engagement in database: %v", err)
		return fmt.Errorf("failed to update engagement: %w", err)
	}

	// Return success response
	log.Info("Engagement updated successfully: ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Engagement updated successfully",
	})
}

// HandleDeleteEngagement handles DELETE /api/engagements/{id} requests
func (h *Handler) HandleDeleteEngagement(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid engagement ID provided for deletion: %s", idStr)
		return fmt.Errorf("invalid engagement ID: %w", err)
	}

	// Delete the engagement from database
	log.Debug("Deleting engagement with ID: %d", id)
	if err := h.store.DeleteEngagement(id); err != nil {
		log.Error("Failed to delete engagement: %v", err)
		return fmt.Errorf("failed to delete engagement: %w", err)
	}

	// Return success response
	log.Info("Engagement deleted successfully: ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Engagement deleted successfully",
	})
}

// HandleGetExpertEngagements handles GET /api/experts/{id}/engagements requests
func (h *Handler) HandleGetExpertEngagements(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert ID provided for engagement retrieval: %s", idStr)
		return fmt.Errorf("invalid expert ID: %w", err)
	}

	// Retrieve the expert's engagements from database
	log.Debug("Retrieving engagements for expert with ID: %d", id)
	engagements, err := h.store.ListEngagements(id)
	if err != nil {
		log.Error("Failed to retrieve engagements for expert %d: %v", id, err)
		return fmt.Errorf("failed to retrieve engagements: %w", err)
	}

	// Return engagements
	log.Debug("Returning %d engagements for expert ID: %d", len(engagements), id)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(engagements)
}