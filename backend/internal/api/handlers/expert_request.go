package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// ExpertRequestHandler handles expert request-related API endpoints
type ExpertRequestHandler struct {
	store storage.Storage
}

// NewExpertRequestHandler creates a new expert request handler
func NewExpertRequestHandler(store storage.Storage) *ExpertRequestHandler {
	return &ExpertRequestHandler{
		store: store,
	}
}

// HandleCreateExpertRequest handles POST /api/expert-requests requests
func (h *ExpertRequestHandler) HandleCreateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/expert-requests request")
	
	// Parse request body
	var request domain.ExpertRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Warn("Failed to parse expert request: %v", err)
		return fmt.Errorf("invalid request body: %w", err)
	}

	// Validate required fields
	log.Debug("Validating expert request fields")
	
	// Validate name
	if strings.TrimSpace(request.Name) == "" {
		log.Warn("Missing name in expert request")
		return fmt.Errorf("name is required")
	}
	
	// Validate institution
	if strings.TrimSpace(request.Institution) == "" {
		log.Warn("Missing institution in expert request")
		return fmt.Errorf("institution is required")
	}
	
	// Validate designation
	if strings.TrimSpace(request.Designation) == "" {
		log.Warn("Missing designation in expert request")
		return fmt.Errorf("designation is required")
	}
	
	// Validate contact information (email or phone)
	if strings.TrimSpace(request.Email) == "" && strings.TrimSpace(request.Phone) == "" {
		log.Warn("Missing contact information in expert request")
		return fmt.Errorf("at least one contact method (email or phone) is required")
	}

	// Set default values
	log.Debug("Setting default values for expert request")
	request.Status = "pending" // Default status for new requests
	request.CreatedAt = time.Now()

	// Create the request in database
	log.Debug("Creating expert request in database: %s, Institution: %s", 
		request.Name, request.Institution)
	id, err := h.store.CreateExpertRequest(&request)
	if err != nil {
		log.Error("Failed to create expert request in database: %v", err)
		return fmt.Errorf("failed to create expert request: %w", err)
	}

	// Set the ID in the response
	request.ID = id
	log.Info("Expert request created successfully: ID: %d, Name: %s", id, request.Name)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(request)
}

// HandleGetExpertRequests handles GET /api/expert-requests requests
func (h *ExpertRequestHandler) HandleGetExpertRequests(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/expert-requests request")
	
	// Parse query parameters for filtering
	status := r.URL.Query().Get("status")
	if status != "" {
		log.Debug("Filtering expert requests by status: %s", status)
	}
	
	// Parse pagination parameters
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 100 // Default limit for requests
	}
	
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}
	
	// Build filters map
	filters := make(map[string]interface{})
	if status != "" {
		filters["status"] = status
	}
	
	// Retrieve requests from database
	log.Debug("Retrieving expert requests with filters: %v", filters)
	requests, err := h.store.ListExpertRequests(status, limit, offset)
	if err != nil {
		log.Error("Failed to retrieve expert requests: %v", err)
		return fmt.Errorf("failed to retrieve expert requests: %w", err)
	}
	
	// Return requests as JSON
	log.Debug("Returning %d expert requests", len(requests))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(requests)
}

// HandleGetExpertRequest handles GET /api/expert-requests/{id} requests
func (h *ExpertRequestHandler) HandleGetExpertRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate expert request ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert request ID provided: %s", idStr)
		return fmt.Errorf("invalid request ID: %w", err)
	}
	
	// Retrieve expert request from database
	log.Debug("Retrieving expert request with ID: %d", id)
	request, err := h.store.GetExpertRequest(id)
	if err != nil {
		if err == domain.ErrNotFound {
			log.Warn("Expert request not found with ID: %d", id)
			return domain.ErrNotFound
		}
		
		log.Error("Failed to get expert request: %v", err)
		return fmt.Errorf("failed to retrieve expert request: %w", err)
	}
	
	// Return expert request data
	log.Debug("Successfully retrieved expert request: ID: %d, Name: %s", request.ID, request.Name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(request)
}

// HandleUpdateExpertRequest handles PUT /api/expert-requests/{id} requests
func (h *ExpertRequestHandler) HandleUpdateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate expert request ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert request ID provided for update: %s", idStr)
		return fmt.Errorf("invalid request ID: %w", err)
	}
	
	// Retrieve existing expert request from database
	log.Debug("Checking if expert request exists with ID: %d", id)
	existingRequest, err := h.store.GetExpertRequest(id)
	if err != nil {
		if err == domain.ErrNotFound {
			log.Warn("Expert request not found for update ID: %d", id)
			return domain.ErrNotFound
		}
		
		log.Error("Failed to get expert request: %v", err)
		return fmt.Errorf("failed to retrieve expert request: %w", err)
	}
	
	// Parse update data
	var updateRequest domain.ExpertRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		log.Warn("Failed to parse expert request update: %v", err)
		return fmt.Errorf("invalid request body: %w", err)
	}
	
	// Ensure ID matches path parameter
	updateRequest.ID = id
	
	// Handle status changes - if status is changing to "approved", create an expert record
	log.Debug("Processing request update, current status: %s, new status: %s", 
		existingRequest.Status, updateRequest.Status)
	
	if existingRequest.Status != "approved" && updateRequest.Status == "approved" {
		log.Info("Expert request being approved, creating expert record from request data")
		
		// Create a new expert from the request data
		// Generate a unique expert ID if not provided
		expertIDStr := updateRequest.ExpertID
		if expertIDStr == "" || len(expertIDStr) < 3 {
			var genErr error
			expertIDStr, genErr = h.store.GenerateUniqueExpertID()
			if genErr != nil {
				log.Error("Failed to generate unique expert ID: %v", genErr)
				return fmt.Errorf("failed to generate unique expert ID: %w", genErr)
			}
			log.Debug("Generated unique expert ID: %s", expertIDStr)
		}
		
		// Create expert record with data from the request
		expert := &domain.Expert{
			ExpertID:        expertIDStr,
			Name:            updateRequest.Name,
			Designation:     updateRequest.Designation,
			Institution:     updateRequest.Institution,
			IsBahraini:      updateRequest.IsBahraini,
			IsAvailable:     updateRequest.IsAvailable,
			Rating:          updateRequest.Rating,
			Role:            updateRequest.Role,
			EmploymentType:  updateRequest.EmploymentType,
			GeneralArea:     updateRequest.GeneralArea,
			SpecializedArea: updateRequest.SpecializedArea,
			IsTrained:       updateRequest.IsTrained,
			CVPath:          updateRequest.CVPath,
			Phone:           updateRequest.Phone,
			Email:           updateRequest.Email,
			IsPublished:     updateRequest.IsPublished,
			Biography:       updateRequest.Biography,
			CreatedAt:       time.Now(),
		}
		
		// Create the expert record in database
		log.Debug("Creating expert record: %s, Institution: %s", expert.Name, expert.Institution)
		createdID, err := h.store.CreateExpert(expert)
		if err != nil {
			log.Error("Failed to create expert from request: %v", err)
			
			// Check if this is a uniqueness constraint error
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return fmt.Errorf("an expert with this ID already exists: %w", err)
			}
			
			return fmt.Errorf("failed to create expert from request: %w", err)
		}
		
		// Set the reviewed timestamp
		updateRequest.ReviewedAt = time.Now()
		
		// Update the expert request with the expert ID
		updateRequest.ExpertID = fmt.Sprintf("EXP-%d", createdID)
		log.Info("Expert created successfully from request: Expert ID: %d", createdID)
	}
	
	// Perform update to the expert request
	// Use a specific method for status updates
	if updateRequest.Status != "" && updateRequest.Status != existingRequest.Status {
		log.Debug("Updating expert request ID: %d, Status: %s", id, updateRequest.Status)
		if err := h.store.UpdateExpertRequestStatus(id, updateRequest.Status, updateRequest.RejectionReason, 0); err != nil {
			log.Error("Failed to update expert request status: %v", err)
			return fmt.Errorf("failed to update expert request: %w", err)
		}
	}
	
	// Return success response
	log.Info("Expert request updated successfully: ID: %d, Status: %s", id, updateRequest.Status)
	resp := map[string]interface{}{
		"success": true,
		"message": "Expert request updated successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}