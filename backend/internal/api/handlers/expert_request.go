package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	
	"expertdb/internal/auth"
	"expertdb/internal/documents"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// ExpertRequestHandler handles expert request-related API endpoints
type ExpertRequestHandler struct {
	store           storage.Storage
	documentService *documents.Service
}

// NewExpertRequestHandler creates a new expert request handler
func NewExpertRequestHandler(store storage.Storage, documentService *documents.Service) *ExpertRequestHandler {
	return &ExpertRequestHandler{
		store:           store,
		documentService: documentService,
	}
}

// writeJSON is a helper function to write JSON responses
func writeJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// HandleCreateExpertRequest handles POST /api/expert-requests requests
func (h *ExpertRequestHandler) HandleCreateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/expert-requests request")

	// Parse multipart form (max 10MB for file)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Warn("Failed to parse multipart form: %v", err)
		return fmt.Errorf("failed to parse form: %w", err)
	}

	// Parse JSON part (expert request data)
	var req domain.ExpertRequest
	jsonData := r.FormValue("data")
	if jsonData == "" {
		log.Warn("Missing JSON data in form")
		return fmt.Errorf("missing request data")
	}
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		log.Warn("Failed to parse JSON data: %v", err)
		return fmt.Errorf("invalid JSON data: %w", err)
	}

	// Validate mandatory fields
	if req.Name == "" || req.Email == "" || req.Phone == "" || req.Biography == "" ||
		req.Designation == "" || req.Institution == "" || req.GeneralArea == 0 ||
		req.Role == "" || req.EmploymentType == "" || req.Rating == "" {
		log.Warn("Missing required fields in expert request")
		return fmt.Errorf("all fields (name, email, phone, biography, designation, institution, general_area, role, employment_type, rating) are required")
	}

	// Handle CV file
	file, fileHeader, err := r.FormFile("cv")
	if err != nil {
		log.Warn("Failed to get CV file: %v", err)
		return fmt.Errorf("CV file is required: %w", err)
	}
	defer file.Close()

	// Set created_by from JWT context
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		log.Warn("Failed to get user claims from context")
		return domain.ErrUnauthorized
	}
	
	// Extract user ID from claims
	if sub, ok := claims["sub"].(string); ok {
		userID, err := strconv.ParseInt(sub, 10, 64)
		if err == nil {
			req.CreatedBy = userID
		} else {
			log.Warn("Failed to parse user ID from claims: %v", err)
			return domain.ErrUnauthorized
		}
	} else {
		log.Warn("Failed to get user ID from claims")
		return domain.ErrUnauthorized
	}

	// Set default values
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}
	if req.Status == "" {
		req.Status = "pending"
	}

	// Upload CV using document service
	// Since the expert hasn't been created yet, we'll use a temporary ID
	// We'll update it later once we have the actual expert ID
	tempExpertID := int64(-1) // Use a temporary negative ID to indicate this is for a request
	
	// Upload document using the document service's method
	doc, err := h.documentService.CreateDocument(tempExpertID, file, fileHeader, "cv")
	if err != nil {
		log.Error("Failed to upload CV: %v", err)
		return fmt.Errorf("failed to upload CV: %w", err)
	}
	req.CVPath = doc.FilePath

	// Create request in database
	log.Debug("Creating expert request for %s", req.Name)
	id, err := h.store.CreateExpertRequest(&req)
	if err != nil {
		log.Error("Failed to create expert request: %v", err)
		return fmt.Errorf("failed to create expert request: %w", err)
	}

	// Return success response
	log.Info("Expert request created successfully with ID: %d", id)
	resp := map[string]interface{}{
		"id": id,
	}
	return writeJSON(w, http.StatusCreated, resp)
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
	
	// Check if user is admin for status updates
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		log.Warn("Failed to get user claims from context")
		return domain.ErrUnauthorized
	}
	
	// Check if user has admin role
	role, ok := claims["role"].(string)
	if !ok || role != auth.RoleAdmin {
		log.Warn("Non-admin attempted to update request status")
		return domain.ErrForbidden
	}
	
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
	
	// Validate status
	if updateRequest.Status != "" && 
	   updateRequest.Status != "approved" && 
	   updateRequest.Status != "rejected" && 
	   updateRequest.Status != "pending" {
		log.Warn("Invalid status provided: %s", updateRequest.Status)
		return fmt.Errorf("invalid status: %s", updateRequest.Status)
	}
	
	// Get user ID from context for reviewer (claims already validated above)
	var userID int64 = 0
	if sub, ok := claims["sub"].(string); ok {
		parsedID, err := strconv.ParseInt(sub, 10, 64)
		if err == nil {
			userID = parsedID
		} else {
			log.Warn("Failed to parse user ID from claims: %v", err)
			return domain.ErrUnauthorized
		}
	}
	
	// Perform update to the expert request status
	// This will now trigger expert creation automatically in the storage layer
	if updateRequest.Status != "" && updateRequest.Status != existingRequest.Status {
		log.Debug("Updating expert request ID: %d, Status: %s", id, updateRequest.Status)
		if err := h.store.UpdateExpertRequestStatus(id, updateRequest.Status, updateRequest.RejectionReason, userID); err != nil {
			log.Error("Failed to update expert request status: %v", err)
			return fmt.Errorf("failed to update expert request: %w", err)
		}
	} else {
		// If not a status update, update the entire request
		// This preserves created_by field
		updateRequest.CreatedBy = existingRequest.CreatedBy
		if err := h.store.UpdateExpertRequest(&updateRequest); err != nil {
			log.Error("Failed to update expert request: %v", err)
			return fmt.Errorf("failed to update expert request: %w", err)
		}
	}
	
	// Return success response
	log.Info("Expert request updated successfully: ID: %d, Status: %s", id, updateRequest.Status)
	return writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Expert request updated successfully",
	})
}