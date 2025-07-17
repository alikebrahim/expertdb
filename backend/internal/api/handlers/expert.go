package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"expertdb/internal/api/response"
	"expertdb/internal/documents"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// ExpertHandler handles expert-related API endpoints
type ExpertHandler struct {
	store           storage.Storage
	documentService *documents.Service
}

// NewExpertHandler creates a new expert handler
func NewExpertHandler(store storage.Storage, documentService *documents.Service) *ExpertHandler {
	return &ExpertHandler{
		store:           store,
		documentService: documentService,
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
	
	// Phase 3A: Additional Expert Filtering
	
	// Filter by nationality (Bahraini/non-Bahraini)
	if nationality := queryParams.Get("by_nationality"); nationality != "" {
		filters["by_nationality"] = nationality
	}
	
	// Filter by general area
	if generalArea := queryParams.Get("by_general_area"); generalArea != "" {
		if area, err := strconv.ParseInt(generalArea, 10, 64); err == nil {
			filters["by_general_area"] = area
		}
	}
	
	// Filter by specialized area
	if specializedArea := queryParams.Get("by_specialized_area"); specializedArea != "" {
		filters["by_specialized_area"] = specializedArea
	}
	
	// Filter by employment type
	if employmentType := queryParams.Get("by_employment_type"); employmentType != "" {
		filters["by_employment_type"] = employmentType
	}
	
	// Filter by role
	if role := queryParams.Get("by_role"); role != "" {
		filters["by_role"] = role
	}

	// Process sorting parameters
	sortBy := "name"   // Default sort field
	sortOrder := "asc" // Default sort order

	if sortParam := queryParams.Get("sort_by"); sortParam != "" {
		// Enhanced sort field validation - more options as per Phase 3B
		allowedSortFields := map[string]bool{
			"name": true, 
			"institution": true, 
			"role": true,
			"created_at": true, 
			"updated_at": true,
			"rating": true, 
			"general_area": true,
			"expert_id": true,  // Add ability to sort by expert ID
			"designation": true, // Add ability to sort by designation
			"employment_type": true, // Add ability to sort by employment type
			"nationality": true, // Add ability to sort by nationality
			"specialized_area": true, // Add ability to sort by specialized area
			"is_bahraini": true, // Add ability to sort by Bahraini status
			"is_available": true, // Add ability to sort by availability
			"is_published": true, // Add ability to sort by published status
		}
		// Convert to database column name format if needed (e.g., camelCase to snake_case)
		dbFieldName := sortParam
		if sortParam == "expertId" {
			dbFieldName = "expert_id"
		} else if sortParam == "specializedArea" {
			dbFieldName = "specialized_area"
		} else if sortParam == "employmentType" {
			dbFieldName = "employment_type"
		} else if sortParam == "generalArea" {
			dbFieldName = "general_area"
		} else if sortParam == "isBahraini" {
			dbFieldName = "is_bahraini"
		} else if sortParam == "isAvailable" {
			dbFieldName = "is_available"
		} else if sortParam == "isPublished" {
			dbFieldName = "is_published"
		} else if sortParam == "createdAt" {
			dbFieldName = "created_at"
		} else if sortParam == "updatedAt" {
			dbFieldName = "updated_at"
		}
		
		if allowedSortFields[dbFieldName] {
			sortBy = dbFieldName // Use the validated field name
		} else {
			log.Warn("Invalid sort field requested: %s. Using default: name", sortParam)
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

	// Enhanced pagination metadata for Phase 3B
	// Calculate pagination information
	totalPages := (totalCount + limit - 1) / limit // Ceiling division
	currentPage := (offset / limit) + 1
	hasMore := offset+len(experts) < totalCount
	hasNext := currentPage < totalPages
	hasPrev := currentPage > 1
	
	// Set pagination headers for client convenience
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", totalCount))
	w.Header().Set("X-Total-Pages", fmt.Sprintf("%d", totalPages))
	w.Header().Set("X-Current-Page", fmt.Sprintf("%d", currentPage))
	w.Header().Set("X-Page-Size", fmt.Sprintf("%d", limit))
	w.Header().Set("X-Has-Next-Page", fmt.Sprintf("%t", hasNext))
	w.Header().Set("X-Has-Prev-Page", fmt.Sprintf("%t", hasPrev))
	
	// Create a response object that includes both experts and metadata
	responseData := map[string]interface{}{
		"experts": experts,
		"pagination": map[string]interface{}{
			"totalCount": totalCount,
			"totalPages": totalPages,
			"currentPage": currentPage,
			"pageSize": limit,
			"hasNextPage": hasNext,
			"hasPrevPage": hasPrev,
			"hasMore": hasMore,
		},
	}

	// Return results
	log.Debug("Returning %d experts (page %d/%d, total count: %d)", 
		len(experts), currentPage, totalPages, totalCount)
	
	// Use the standardized response format
	return response.Success(w, http.StatusOK, "", responseData)
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
	
	// Use the standardized response format
	return response.Success(w, http.StatusOK, "", expert)
}

// HandleCreateExpert handles POST /api/experts requests
func (h *ExpertHandler) HandleCreateExpert(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/experts request")

	// Parse request body
	var expert domain.Expert
	if err := json.NewDecoder(r.Body).Decode(&expert); err != nil {
		log.Warn("Failed to parse expert creation request: %v", err)
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid JSON format",
			"details": err.Error(),
			"suggestion": "Check the request syntax and ensure all fields have proper types",
		})
	}

	// Validate required fields - collect all validation errors
	errors := []string{}
	
	// The following fields are required per SRS
	if expert.Name == "" {
		errors = append(errors, "name is required")
	}
	
	if expert.Affiliation == "" {
		errors = append(errors, "institution is required")
	}
	
	if expert.Designation == "" {
		errors = append(errors, "designation is required")
	}
	
	if expert.Role == "" {
		errors = append(errors, "role is required")
	} else {
		// Validate role values
		validRoles := []string{"evaluator", "validator", "evaluator/validator"}
		if !containsString(validRoles, strings.ToLower(expert.Role)) {
			errors = append(errors, "role must be one of: evaluator, validator, evaluator/validator")
		}
	}
	
	if expert.EmploymentType == "" {
		errors = append(errors, "employmentType is required")
	} else {
		// Validate employment type values
		validEmploymentTypes := []string{"academic", "employer"}
		if !containsString(validEmploymentTypes, strings.ToLower(expert.EmploymentType)) {
			errors = append(errors, "employmentType must be one of: academic, employer")
		}
	}
	
	if expert.GeneralArea <= 0 {
		errors = append(errors, "generalArea must be a positive number")
	}
	
	if expert.SpecializedArea == "" {
		errors = append(errors, "specializedArea is required")
	}
	
	if expert.Phone == "" {
		errors = append(errors, "phone is required")
	}
	
	if expert.Biography == "" {
		errors = append(errors, "biography is required")
	}
	
	// Skip email validation per SRS requirement - no validation for emails
	
	if len(errors) > 0 {
		log.Warn("Expert validation failed: %v", errors)
		return response.ValidationError(w, errors)
	}

	// Set creation time and default values if not provided
	if expert.CreatedAt.IsZero() {
		expert.CreatedAt = time.Now()
		expert.UpdatedAt = expert.CreatedAt
	}
	
	// Default values
	if !expert.IsPublished {
		expert.IsPublished = false // Explicitly set to false if not provided
	}

	// Create expert in database
	log.Debug("Creating expert: %s, Institution: %s", expert.Name, expert.Affiliation)
	id, err := h.store.CreateExpert(&expert)
	if err != nil {
		log.Error("Failed to create expert in database: %v", err)

		// Check for different types of errors and return appropriate status codes
		if strings.Contains(err.Error(), "expert ID already exists") {
			return writeJSON(w, http.StatusConflict, map[string]interface{}{
				"error": err.Error(),
				"suggestion": "Let the system generate a unique ID automatically by omitting the expertId field",
			})
		}
		
		if strings.Contains(err.Error(), "email already exists") {
			return writeJSON(w, http.StatusConflict, map[string]interface{}{
				"error": err.Error(),
				"suggestion": "Either use a different email or update the existing expert record",
			})
		}
		
		if strings.Contains(err.Error(), "invalid general area") || 
		   strings.Contains(err.Error(), "referenced resource does not exist") {
			return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
				"suggestion": "Use GET /api/expert/areas to see the list of valid general areas",
			})
		}
		
		if strings.Contains(err.Error(), "required field") {
			return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		
		// Generic database error
		return writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Database error creating expert",
			"details": err.Error(),
		})
	}

	// Return success response
	log.Info("Expert created successfully with ID: %d", id)
	return response.Success(w, http.StatusCreated, "Expert created successfully", map[string]interface{}{
		"id":      id,
		"expertId": expert.ExpertID,
	})
}

// Helper function to write JSON responses
func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	// Use the standardized response format
	if resp, ok := data.(map[string]interface{}); ok {
		// If data already has success and message fields, use response.JSON directly
		if _, hasSuccess := resp["success"]; hasSuccess {
			return response.JSON(w, status, resp)
		}
		
		// If it's an error response
		if errMsg, hasError := resp["error"]; hasError && status >= 400 {
			errorResp := response.ErrorResponse{
				Error: errMsg.(string),
			}
			
			// Include errors array if present
			if errDetails, hasErrors := resp["errors"]; hasErrors {
				if errMap, isMap := errDetails.(map[string]string); isMap {
					errorResp.Errors = errMap
				}
			}
			
			return response.JSON(w, status, errorResp)
		}
	}
	
	// Otherwise, use Success with default parameters
	return response.Success(w, status, "", data)
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
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

	// Check if this is a multipart form or JSON update
	contentType := r.Header.Get("Content-Type")
	var updateExpert domain.Expert
	
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// This is a file upload with form data
		log.Debug("Processing multipart form update for expert ID: %d", id)
		
		// Parse multipart form (max 10MB for file)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Warn("Failed to parse multipart form: %v", err)
			return fmt.Errorf("failed to parse form: %w", err)
		}
		
		// Parse JSON part for the expert data
		jsonData := r.FormValue("data")
		if jsonData == "" {
			log.Warn("Missing JSON data in form")
			return fmt.Errorf("missing expert data")
		}
		
		if err := json.Unmarshal([]byte(jsonData), &updateExpert); err != nil {
			log.Warn("Failed to parse JSON data: %v", err)
			return fmt.Errorf("invalid JSON data: %w", err)
		}
		
		// Process CV file if provided
		cvFile, cvFileHeader, err := r.FormFile("cvFile")
		if err == nil {
			// CV file was provided, upload it
			defer cvFile.Close()
			
			cvDoc, err := h.documentService.CreateDocument(id, cvFile, cvFileHeader, "cv")
			if err != nil {
				log.Error("Failed to upload CV file: %v", err)
				return fmt.Errorf("failed to upload CV: %w", err)
			}
			
			// Update the expert's CV path
			updateExpert.CVPath = cvDoc.FilePath
			log.Debug("Updated CV path for expert ID: %d to %s", id, cvDoc.FilePath)
		} else if err != http.ErrMissingFile {
			log.Warn("Error accessing CV file: %v", err)
			return fmt.Errorf("error processing CV file: %w", err)
		}
		
		// Process approval document file if provided
		approvalFile, approvalFileHeader, err := r.FormFile("approvalDocument")
		if err == nil {
			// Approval document file was provided, upload it
			defer approvalFile.Close()
			
			approvalDoc, err := h.documentService.CreateDocument(id, approvalFile, approvalFileHeader, "approval")
			if err != nil {
				log.Error("Failed to upload approval document: %v", err)
				return fmt.Errorf("failed to upload approval document: %w", err)
			}
			
			// Update the expert's approval document path
			updateExpert.ApprovalDocumentPath = approvalDoc.FilePath
			log.Debug("Updated approval document path for expert ID: %d to %s", id, approvalDoc.FilePath)
		} else if err != http.ErrMissingFile {
			log.Warn("Error accessing approval document file: %v", err)
			return fmt.Errorf("error processing approval document file: %w", err)
		}
		
	} else {
		// This is a regular JSON update
		log.Debug("Processing JSON update for expert ID: %d", id)
		
		if err := json.NewDecoder(r.Body).Decode(&updateExpert); err != nil {
			log.Warn("Failed to parse expert update request: %v", err)
			return fmt.Errorf("invalid request body: %w", err)
		}
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
		if updateExpert.Affiliation == "" {
			updateExpert.Affiliation = existingExpert.Affiliation
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
	return response.Success(w, http.StatusOK, "Expert updated successfully", map[string]interface{}{
		"id": id,
	})
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
	return response.Success(w, http.StatusOK, "Expert deleted successfully", nil)
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
	
	// Use the standardized response format
	return response.Success(w, http.StatusOK, "", areas)
}

// AreaRequest represents a request to create or update an area
type AreaRequest struct {
	Name string `json:"name"`
}

// HandleCreateArea handles POST /api/expert/areas requests
func (h *ExpertHandler) HandleCreateArea(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/expert/areas request")

	// Parse the request body
	var req AreaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("Failed to parse area creation request: %v", err)
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid JSON format",
			"details": err.Error(),
		})
	}

	// Validate area name
	if strings.TrimSpace(req.Name) == "" {
		log.Warn("Area name validation failed: empty name")
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": "Area name cannot be empty",
		})
	}

	// Create area in database
	id, err := h.store.CreateArea(req.Name)
	if err != nil {
		log.Error("Failed to create area: %v", err)
		
		// Check for duplicate name error
		if strings.Contains(err.Error(), "already exists") {
			return writeJSON(w, http.StatusConflict, map[string]interface{}{
				"error": err.Error(),
				"suggestion": "Use a different area name",
			})
		}
		
		return writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create area",
			"details": err.Error(),
		})
	}

	// Return success response
	log.Info("Area created successfully with ID: %d", id)
	return response.Success(w, http.StatusCreated, "Area created successfully", map[string]interface{}{
		"id": id,
		"name": req.Name,
	})
}

// HandleUpdateArea handles PUT /api/expert/areas/{id} requests
func (h *ExpertHandler) HandleUpdateArea(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate area ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid area ID provided: %s", idStr)
		return fmt.Errorf("invalid area ID: %w", err)
	}
	
	// Parse the request body
	var req AreaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("Failed to parse area update request: %v", err)
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid JSON format",
			"details": err.Error(),
		})
	}
	
	// Validate area name
	if strings.TrimSpace(req.Name) == "" {
		log.Warn("Area name validation failed: empty name")
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": "Area name cannot be empty",
		})
	}
	
	// Update area in database
	err = h.store.UpdateArea(id, req.Name)
	if err != nil {
		log.Error("Failed to update area: %v", err)
		
		// Check for specific errors
		if err == domain.ErrNotFound {
			return response.NotFound(w, fmt.Sprintf("No area exists with ID: %d", id))
		}
		
		if strings.Contains(err.Error(), "already exists") {
			return writeJSON(w, http.StatusConflict, map[string]interface{}{
				"error": err.Error(),
				"suggestion": "Use a different area name",
			})
		}
		
		return writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update area",
			"details": err.Error(),
		})
	}
	
	// Return success response
	log.Info("Area updated successfully: ID: %d", id)
	return response.Success(w, http.StatusOK, "Area updated successfully", map[string]interface{}{
		"id": id,
		"name": req.Name,
	})
}
