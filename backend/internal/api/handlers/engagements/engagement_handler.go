// Package engagements provides handlers for engagement-related API endpoints
package engagements

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"expertdb/internal/api/response"
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

	// Set the ID in the response and return with standardized format
	engagement.ID = id
	log.Info("Engagement created successfully: ID: %d, Type: %s, Expert: %d",
		id, engagement.EngagementType, engagement.ExpertID)
	
	return response.Success(w, http.StatusCreated, "Engagement created successfully", engagement)
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

	// Return engagement data with standardized response
	log.Debug("Successfully retrieved engagement: ID: %d, Type: %s", engagement.ID, engagement.EngagementType)
	return response.Success(w, http.StatusOK, "", engagement)
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

	// Return success response with standardized format
	log.Info("Engagement updated successfully: ID: %d", id)
	return response.Success(w, http.StatusOK, "Engagement updated successfully", map[string]interface{}{
		"id": id,
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

	// Return success response with standardized format
	log.Info("Engagement deleted successfully: ID: %d", id)
	return response.Success(w, http.StatusOK, "Engagement deleted successfully", nil)
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

	// Extract query parameters for filtering
	engagementType := r.URL.Query().Get("type")
	
	// Parse pagination parameters
	limit, offset := 50, 0 // Default values
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Retrieve the expert's engagements from database with filters
	log.Debug("Retrieving engagements for expert with ID: %d, type: %s", id, engagementType)
	engagements, err := h.store.ListEngagements(id, engagementType, limit, offset)
	if err != nil {
		log.Error("Failed to retrieve engagements for expert %d: %v", id, err)
		return fmt.Errorf("failed to retrieve engagements: %w", err)
	}

	// Return engagements with standardized response and pagination metadata
	log.Debug("Returning %d engagements for expert ID: %d", len(engagements), id)
	responseData := map[string]interface{}{
		"engagements": engagements,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(engagements),
		},
		"expertId": id,
	}
	return response.Success(w, http.StatusOK, "", responseData)
}

// HandleListEngagements handles GET /api/engagements requests with filtering capabilities
func (h *Handler) HandleListEngagements(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/engagements request")

	// Extract query parameters for filtering
	expertIDStr := r.URL.Query().Get("expert_id")
	engagementType := r.URL.Query().Get("type")

	// Parse expert_id if provided
	var expertID int64 = 0 // Default to 0 (all experts)
	if expertIDStr != "" {
		if parsed, err := strconv.ParseInt(expertIDStr, 10, 64); err == nil && parsed > 0 {
			expertID = parsed
		} else {
			log.Warn("Invalid expert_id parameter: %s", expertIDStr)
			return fmt.Errorf("invalid expert_id parameter")
		}
	}

	// Parse pagination parameters
	limit, offset := 50, 0 // Default values
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Validate engagement type if provided (Phase 11B)
	if engagementType != "" && engagementType != "validator" && engagementType != "evaluator" {
		log.Warn("Invalid engagement type parameter: %s", engagementType)
		return fmt.Errorf("engagement type must be 'validator' or 'evaluator'")
	}

	// Retrieve engagements with filters
	log.Debug("Retrieving engagements with filters - expert_id: %d, type: %s", expertID, engagementType)
	engagements, err := h.store.ListEngagements(expertID, engagementType, limit, offset)
	if err != nil {
		log.Error("Failed to retrieve engagements: %v", err)
		return fmt.Errorf("failed to retrieve engagements: %w", err)
	}

	// Return engagements with standardized response and pagination metadata
	log.Debug("Returning %d engagements", len(engagements))
	responseData := map[string]interface{}{
		"engagements": engagements,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(engagements),
		},
		"filters": map[string]interface{}{
			"expertId":       expertID,
			"engagementType": engagementType,
		},
	}
	return response.Success(w, http.StatusOK, "", responseData)
}

// HandleImportEngagements handles POST /api/engagements/import requests
func (h *Handler) HandleImportEngagements(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/engagements/import request")

	// Determine content type
	contentType := r.Header.Get("Content-Type")
	isJSON := strings.Contains(contentType, "application/json")

	// Prepare for engagement parsing
	var engagements []*domain.Engagement
	var err error

	if isJSON {
		// Parse JSON payload
		err = json.NewDecoder(r.Body).Decode(&engagements)
		if err != nil {
			log.Warn("Failed to parse JSON engagement import: %v", err)
			return fmt.Errorf("invalid JSON format: %w", err)
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Parse multipart form data (CSV upload)
		err = r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			log.Warn("Failed to parse multipart form: %v", err)
			return fmt.Errorf("invalid form data: %w", err)
		}

		// Get the file from the form data
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Warn("Failed to get file from form: %v", err)
			return fmt.Errorf("no file uploaded: %w", err)
		}
		defer file.Close()

		// Read and parse CSV
		engagements, err = parseCSVEngagements(file)
		if err != nil {
			log.Warn("Failed to parse CSV data: %v", err)
			return fmt.Errorf("failed to parse CSV: %w", err)
		}
	} else {
		log.Warn("Unsupported content type for import: %s", contentType)
		return fmt.Errorf("unsupported content type: expected application/json or multipart/form-data (CSV)")
	}

	// Validate engagements list
	if len(engagements) == 0 {
		log.Warn("Empty engagement list in import request")
		return fmt.Errorf("no engagements provided for import")
	}

	// Import engagements
	log.Debug("Importing %d engagements", len(engagements))
	successCount, errors := h.store.ImportEngagements(engagements)

	// Prepare response
	type ImportResponse struct {
		Success        bool              `json:"success"`
		SuccessCount   int               `json:"successCount"`
		FailureCount   int               `json:"failureCount"`
		TotalCount     int               `json:"totalCount"`
		Errors         map[string]string `json:"errors,omitempty"` // Index-error mapping for failed imports
	}

	// Convert error map to string map for JSON
	errorMap := make(map[string]string)
	for index, err := range errors {
		errorMap[fmt.Sprintf("%d", index)] = err.Error()
	}

	respData := ImportResponse{
		Success:      successCount > 0,
		SuccessCount: successCount,
		FailureCount: len(errors),
		TotalCount:   len(engagements),
		Errors:       errorMap,
	}

	// Return response with standardized format
	log.Info("Engagement import completed: %d successful, %d failed", successCount, len(errors))
	
	message := fmt.Sprintf("Import completed: %d successful, %d failed out of %d total", 
		successCount, len(errors), len(engagements))
		
	return response.Success(w, http.StatusOK, message, respData)
}

// parseCSVEngagements parses CSV data into engagement records
func parseCSVEngagements(file io.Reader) ([]*domain.Engagement, error) {
	// Create a new CSV reader
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header row
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create header index map
	headerIndex := make(map[string]int)
	for i, col := range header {
		headerIndex[strings.ToLower(strings.TrimSpace(col))] = i
	}

	// Verify required columns
	requiredColumns := []string{"expert_id", "engagement_type", "start_date"}
	for _, col := range requiredColumns {
		if _, ok := headerIndex[col]; !ok {
			return nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	// Read and parse engagement rows
	var engagements []*domain.Engagement
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row: %w", err)
		}

		// Parse expert ID
		expertID, err := strconv.ParseInt(row[headerIndex["expert_id"]], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid expert_id in row: %w", err)
		}

		// Parse engagement type
		engagementType := strings.TrimSpace(row[headerIndex["engagement_type"]])
		if engagementType != "validator" && engagementType != "evaluator" {
			return nil, fmt.Errorf("invalid engagement_type '%s': must be 'validator' or 'evaluator'", engagementType)
		}

		// Parse start date
		startDate, err := time.Parse("2006-01-02", strings.TrimSpace(row[headerIndex["start_date"]]))
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format (expected YYYY-MM-DD): %w", err)
		}

		// Create engagement object
		engagement := &domain.Engagement{
			ExpertID:       expertID,
			EngagementType: engagementType,
			StartDate:      startDate,
			Status:         "active", // Default status for imports
			CreatedAt:      time.Now().UTC(),
		}

		// Optional fields
		if idx, ok := headerIndex["end_date"]; ok && idx < len(row) && row[idx] != "" {
			endDate, err := time.Parse("2006-01-02", strings.TrimSpace(row[idx]))
			if err == nil {
				engagement.EndDate = endDate
			}
		}

		if idx, ok := headerIndex["project_name"]; ok && idx < len(row) {
			engagement.ProjectName = strings.TrimSpace(row[idx])
		}

		if idx, ok := headerIndex["status"]; ok && idx < len(row) && row[idx] != "" {
			engagement.Status = strings.TrimSpace(row[idx])
		}

		if idx, ok := headerIndex["feedback_score"]; ok && idx < len(row) && row[idx] != "" {
			if score, err := strconv.Atoi(row[idx]); err == nil {
				engagement.FeedbackScore = score
			}
		}

		if idx, ok := headerIndex["notes"]; ok && idx < len(row) {
			engagement.Notes = strings.TrimSpace(row[idx])
		}

		engagements = append(engagements, engagement)
	}

	return engagements, nil
}