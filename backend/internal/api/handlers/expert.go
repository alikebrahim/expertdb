package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// ExpertHandler handles expert-related API endpoints
type ExpertHandler struct {
	store storage.Storage
}

// NewExpertHandler creates a new expert handler
func NewExpertHandler(store storage.Storage) *ExpertHandler {
	return &ExpertHandler{
		store: store,
	}
}

// HandleGetExperts handles GET /api/experts requests
func (h *ExpertHandler) HandleGetExperts(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/experts request")

	// Parse query parameters for filtering
	queryParams := r.URL.Query()
	filters := make(map[string]interface{})

	// Process name filter
	if name := queryParams.Get("name"); name != "" {
		filters["name"] = name
	}

	// Process boolean availability filter
	if available := queryParams.Get("is_available"); available != "" {
		if available == "true" {
			filters["isAvailable"] = true
		} else if available == "false" {
			filters["isAvailable"] = false
		}
	}

	// Process role filter
	if role := queryParams.Get("role"); role != "" {
		filters["role"] = role
	}

	// Process general area filter
	if generalArea := queryParams.Get("generalArea"); generalArea != "" {
		if area, err := strconv.ParseInt(generalArea, 10, 64); err == nil {
			filters["generalArea"] = area
		}
	}

	// Process sorting parameters
	sortBy := "name"   // Default sort field
	sortOrder := "asc" // Default sort order

	if sortParam := queryParams.Get("sort_by"); sortParam != "" {
		// Validate sort field against allowed fields
		allowedSortFields := map[string]bool{
			"name": true, "institution": true, "role": true,
			"created_at": true, "rating": true, "general_area": true,
		}
		if allowedSortFields[sortParam] {
			sortBy = sortParam
		}
	}

	if orderParam := queryParams.Get("sort_order"); orderParam != "" {
		if orderParam == "desc" {
			sortOrder = "desc"
		}
	}

	// Add sorting to filters
	filters["sort_by"] = sortBy
	filters["sort_order"] = sortOrder

	// Parse pagination parameters
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	// Get total count (without pagination) for headers
	countFilters := make(map[string]interface{})
	for k, v := range filters {
		if k != "sort_by" && k != "sort_order" {
			countFilters[k] = v
		}
	}

	totalCount, err := h.store.CountExperts(countFilters)
	if err != nil {
		log.Error("Failed to count experts: %v", err)
		return fmt.Errorf("failed to count experts: %w", err)
	}

	// Retrieve filtered experts with pagination
	log.Debug("Retrieving experts with filters: %v, limit: %d, offset: %d", filters, limit, offset)
	experts, err := h.store.ListExperts(filters, limit, offset)
	if err != nil {
		log.Error("Failed to list experts: %v", err)
		return fmt.Errorf("failed to retrieve experts: %w", err)
	}

	// Set total count header for pagination
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", totalCount))

	// Return results
	log.Debug("Returning %d experts", len(experts))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(experts)
}

// HandleGetExpert handles GET /api/experts/{id} requests
func (h *ExpertHandler) HandleGetExpert(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert ID provided: %s", idStr)
		return fmt.Errorf("invalid expert ID: %w", err)
	}

	// Retrieve expert from database
	log.Debug("Retrieving expert with ID: %d", id)
	expert, err := h.store.GetExpert(id)
	if err != nil {
		// Return an empty object for not found
		if err == domain.ErrNotFound {
			log.Warn("Expert not found for ID: %d", id)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return json.NewEncoder(w).Encode(&domain.Expert{})
		}

		log.Error("Failed to get expert: %v", err)
		return fmt.Errorf("failed to retrieve expert: %w", err)
	}

	// Return expert data
	log.Debug("Successfully retrieved expert: %s (ID: %d)", expert.Name, expert.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(expert)
}

// HandleCreateExpert handles POST /api/experts requests
func (h *ExpertHandler) HandleCreateExpert(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/experts request")

	// Parse request body
	var expert domain.Expert
	if err := json.NewDecoder(r.Body).Decode(&expert); err != nil {
		log.Warn("Failed to parse expert creation request: %v", err)
		return writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Invalid JSON format: %v", err),
		})
	}

	// Validate required fields
	errors := []string{}
	if expert.Name == "" {
		errors = append(errors, "name is required")
	}
	
	if expert.GeneralArea <= 0 {
		errors = append(errors, "generalArea must be a positive number")
	}
	
	if expert.Email != "" && !isValidEmail(expert.Email) {
		errors = append(errors, fmt.Sprintf("invalid email format: %s", expert.Email))
	}
	
	if len(errors) > 0 {
		log.Warn("Expert validation failed: %v", errors)
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"errors": errors,
		})
	}

	// Set creation time if not provided
	if expert.CreatedAt.IsZero() {
		expert.CreatedAt = time.Now()
		expert.UpdatedAt = expert.CreatedAt
	}

	// Create expert in database
	log.Debug("Creating expert: %s, Institution: %s", expert.Name, expert.Institution)
	id, err := h.store.CreateExpert(&expert)
	if err != nil {
		log.Error("Failed to create expert in database: %v", err)

		// Check specifically for UNIQUE constraint violations on expert_id
		if strings.Contains(err.Error(), "UNIQUE constraint failed") &&
			strings.Contains(err.Error(), "expert_id") {
			return writeJSON(w, http.StatusConflict, map[string]string{
				"error": fmt.Sprintf("Expert ID %s already exists", expert.ExpertID),
			})
		}
		
		// Check for foreign key constraint violations
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Referenced resource does not exist (likely invalid generalArea)",
			})
		}

		return writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Database error creating expert: %v", err),
		})
	}

	// Return success response
	log.Info("Expert created successfully with ID: %d", id)
	return writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"success": true,
		"message": "Expert created successfully",
	})
}

// Helper function to write JSON responses
func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// HandleUpdateExpert handles PUT /api/experts/{id} requests
func (h *ExpertHandler) HandleUpdateExpert(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert ID provided for update: %s", idStr)
		return fmt.Errorf("invalid expert ID: %w", err)
	}

	// Retrieve existing expert (if exists)
	log.Debug("Checking if expert exists with ID: %d", id)
	var existingExpert *domain.Expert
	existingExpert, err = h.store.GetExpert(id)
	if err != nil && err != domain.ErrNotFound {
		log.Error("Error checking for existing expert: %v", err)
		return fmt.Errorf("failed to check existing expert: %w", err)
	}

	// Parse update data
	var updateExpert domain.Expert
	if err := json.NewDecoder(r.Body).Decode(&updateExpert); err != nil {
		log.Warn("Failed to parse expert update request: %v", err)
		return fmt.Errorf("invalid request body: %w", err)
	}

	// Ensure ID matches path parameter
	updateExpert.ID = id

	// If existing expert was found, merge with update data
	if existingExpert != nil {
		// Only replace fields that are set in the update
		if updateExpert.ExpertID == "" {
			updateExpert.ExpertID = existingExpert.ExpertID
		}
		if updateExpert.Name == "" {
			updateExpert.Name = existingExpert.Name
		}
		if updateExpert.Institution == "" {
			updateExpert.Institution = existingExpert.Institution
		}
		if updateExpert.Designation == "" {
			updateExpert.Designation = existingExpert.Designation
		}
		if updateExpert.Nationality == "" {
			updateExpert.Nationality = existingExpert.Nationality
		}
		if updateExpert.Role == "" {
			updateExpert.Role = existingExpert.Role
		}
		if updateExpert.EmploymentType == "" {
			updateExpert.EmploymentType = existingExpert.EmploymentType
		}
		if updateExpert.GeneralArea == 0 {
			updateExpert.GeneralArea = existingExpert.GeneralArea
		}
		if updateExpert.SpecializedArea == "" {
			updateExpert.SpecializedArea = existingExpert.SpecializedArea
		}
		if updateExpert.Phone == "" {
			updateExpert.Phone = existingExpert.Phone
		}
		if updateExpert.Email == "" {
			updateExpert.Email = existingExpert.Email
		}
		if updateExpert.Biography == "" {
			updateExpert.Biography = existingExpert.Biography
		}
		// Preserve created date
		if updateExpert.CreatedAt.IsZero() {
			updateExpert.CreatedAt = existingExpert.CreatedAt
		}
	}

	// Set updated time
	updateExpert.UpdatedAt = time.Now()

	// Update expert in database
	log.Debug("Updating expert ID: %d, Name: %s", id, updateExpert.Name)
	if err := h.store.UpdateExpert(&updateExpert); err != nil {
		log.Error("Failed to update expert in database: %v", err)
		return fmt.Errorf("failed to update expert: %w", err)
	}

	// Return success response
	log.Info("Expert updated successfully: ID: %d", id)
	resp := map[string]interface{}{
		"success": true,
		"message": "Expert updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

// HandleDeleteExpert handles DELETE /api/experts/{id} requests
func (h *ExpertHandler) HandleDeleteExpert(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert ID provided for deletion: %s", idStr)
		return fmt.Errorf("invalid expert ID: %w", err)
	}

	// Delete expert from database
	log.Debug("Deleting expert with ID: %d", id)
	if err := h.store.DeleteExpert(id); err != nil {
		if err == domain.ErrNotFound {
			log.Warn("Expert not found for deletion ID: %d", id)
			return domain.ErrNotFound
		}

		log.Error("Failed to delete expert: %v", err)
		return fmt.Errorf("failed to delete expert: %w", err)
	}

	// Return success response
	log.Info("Expert deleted successfully: ID: %d", id)
	resp := map[string]interface{}{
		"success": true,
		"message": "Expert deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

// HandleGetExpertAreas handles GET /api/expert/areas requests
func (h *ExpertHandler) HandleGetExpertAreas(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/expert/areas request")

	// Retrieve all expert areas from database
	areas, err := h.store.ListAreas()
	if err != nil {
		log.Error("Failed to fetch expert areas: %v", err)
		return fmt.Errorf("failed to fetch expert areas: %w", err)
	}

	// Return areas as JSON
	log.Debug("Returning %d expert areas", len(areas))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(areas)
}
