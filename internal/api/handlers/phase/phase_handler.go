package phase

import (
	"encoding/json"
	"expertdb/internal/api/response"
	"expertdb/internal/auth"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handler represents the phase planning handler
type Handler struct {
	store storage.Storage
}

// NewHandler creates a new phase handler
func NewHandler(store storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}

// HandleListPhases handles GET /api/phases requests
func (h *Handler) HandleListPhases(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Parse query parameters
	queryParams := r.URL.Query()
	status := queryParams.Get("status")
	limit := 100 // Default limit
	offset := 0  // Default offset
	
	// Parse planner ID
	var plannerID int64
	plannerIDParam := queryParams.Get("planner_id")
	if plannerIDParam != "" {
		var err error
		plannerID, err = strconv.ParseInt(plannerIDParam, 10, 64)
		if err != nil {
			log.Debug("Invalid planner_id parameter: %v", err)
			return fmt.Errorf("invalid planner_id parameter: %v", err)
		}
	}
	
	// Parse pagination parameters
	if limitParam := queryParams.Get("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	if offsetParam := queryParams.Get("offset"); offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	// Get phases from store
	phases, err := h.store.ListPhases(status, plannerID, limit, offset)
	if err != nil {
		log.Error("Failed to list phases: %v", err)
		return fmt.Errorf("failed to list phases: %w", err)
	}
	
	// Write standardized response
	responseData := map[string]interface{}{
		"phases": phases,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(phases),
		},
		"filters": map[string]interface{}{
			"status":    status,
			"plannerId": plannerID,
		},
	}
	return response.Success(w, http.StatusOK, "", responseData)
}

// HandleGetPhase handles GET /api/phases/{id} requests
func (h *Handler) HandleGetPhase(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract phase ID from URL
	idStr := r.PathValue("id")
	if idStr == "" {
		return fmt.Errorf("phase ID is required")
	}
	
	// Check if it's a numeric ID or a business ID
	if strings.HasPrefix(idStr, "PH-") {
		// It's a business ID (phase_id)
		phase, err := h.store.GetPhaseByPhaseID(idStr)
		if err != nil {
			if err == domain.ErrNotFound {
				return domain.ErrNotFound
			}
			log.Error("Failed to get phase by business ID: %v", err)
			return fmt.Errorf("failed to get phase: %w", err)
		}
		
		// Write standardized response
		return response.Success(w, http.StatusOK, "", phase)
	}
	
	// It's a numeric ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Debug("Invalid phase ID: %v", err)
		return fmt.Errorf("invalid phase ID: %v", err)
	}
	
	// Get phase from store
	phase, err := h.store.GetPhase(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		log.Error("Failed to get phase: %v", err)
		return fmt.Errorf("failed to get phase: %w", err)
	}
	
	// Write standardized response
	return response.Success(w, http.StatusOK, "", phase)
}

// createPhaseRequest represents the request to create a new phase
type createPhaseRequest struct {
	Title             string                 `json:"title"`
	AssignedPlannerID int64                 `json:"assignedPlannerId"`
	Status            string                 `json:"status"`
	Applications      []createApplicationRequest `json:"applications"`
}

// createApplicationRequest represents the request to create a new application
type createApplicationRequest struct {
	Type              string `json:"type"`
	InstitutionName   string `json:"institutionName"`
	QualificationName string `json:"qualificationName"`
	Expert1           int64  `json:"expert1,omitempty"`
	Expert2           int64  `json:"expert2,omitempty"`
	Status            string `json:"status,omitempty"`
}

// HandleCreatePhase handles POST /api/phases requests
func (h *Handler) HandleCreatePhase(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Parse request body
	var req createPhaseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug("Failed to decode request body: %v", err)
		return fmt.Errorf("invalid request body: %v", err)
	}
	
	// Validate required fields
	var validationErrors []string
	
	if strings.TrimSpace(req.Title) == "" {
		validationErrors = append(validationErrors, "title is required")
	}
	
	if req.AssignedPlannerID <= 0 {
		validationErrors = append(validationErrors, "assigned planner ID is required")
	} else {
		// Verify assigned user exists and can be elevated to planner
		user, err := h.store.GetUser(req.AssignedPlannerID)
		if err != nil {
			if err == domain.ErrNotFound {
				validationErrors = append(validationErrors, fmt.Sprintf("user with ID %d does not exist", req.AssignedPlannerID))
			} else {
				log.Error("Failed to get assigned user: %v", err)
				return fmt.Errorf("failed to verify assigned user: %w", err)
			}
		} else if user.Role == "super_user" || user.Role == "admin" {
			// Admin and super_user roles don't need elevation assignments - they have inherent access
			log.Info("Assigning admin/super_user to phase - they have inherent planner access")
		} else if user.Role != "user" {
			validationErrors = append(validationErrors, fmt.Sprintf("user with ID %d has invalid role for elevation assignment", req.AssignedPlannerID))
		}
	}
	
	// Validate status if provided
	if req.Status != "" {
		validStatuses := []string{"draft", "in_progress", "completed", "cancelled"}
		valid := false
		for _, s := range validStatuses {
			if req.Status == s {
				valid = true
				break
			}
		}
		if !valid {
			validationErrors = append(validationErrors, "status must be one of: draft, in_progress, completed, cancelled")
		}
	}
	
	// Validate applications if provided
	for i, app := range req.Applications {
		if strings.TrimSpace(app.Type) == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("application %d: type is required", i+1))
		} else {
			// Only support QP and IL types
			validTypes := []string{"QP", "IL"}
			valid := false
			for _, t := range validTypes {
				if app.Type == t {
					valid = true
					break
				}
			}
			if !valid {
				validationErrors = append(validationErrors, fmt.Sprintf("application %d: type must be one of: QP (Qualification Placement), IL (Institutional Listing)", i+1))
			}
		}
		
		if strings.TrimSpace(app.InstitutionName) == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("application %d: institution name is required", i+1))
		}
		
		if strings.TrimSpace(app.QualificationName) == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("application %d: qualification name is required", i+1))
		}
		
		// Validate status if provided
		if app.Status != "" {
			validStatuses := []string{"pending", "assigned", "approved", "rejected"}
			valid := false
			for _, s := range validStatuses {
				if app.Status == s {
					valid = true
					break
				}
			}
			if !valid {
				validationErrors = append(validationErrors, fmt.Sprintf("application %d: status must be one of: pending, assigned, approved, rejected", i+1))
			}
		}
		
		// If experts are provided, verify they exist
		if app.Expert1 > 0 {
			exists, err := expertExists(h.store, app.Expert1)
			if err != nil {
				log.Error("Failed to check if expert exists: %v", err)
				return fmt.Errorf("failed to verify expert: %w", err)
			}
			if !exists {
				validationErrors = append(validationErrors, fmt.Sprintf("application %d: expert 1 with ID %d does not exist", i+1, app.Expert1))
			}
		}
		
		if app.Expert2 > 0 {
			exists, err := expertExists(h.store, app.Expert2)
			if err != nil {
				log.Error("Failed to check if expert exists: %v", err)
				return fmt.Errorf("failed to verify expert: %w", err)
			}
			if !exists {
				validationErrors = append(validationErrors, fmt.Sprintf("application %d: expert 2 with ID %d does not exist", i+1, app.Expert2))
			}
		}
	}
	
	// Return validation errors if any
	if len(validationErrors) > 0 {
		log.Debug("Validation errors: %v", validationErrors)
		return respondWithValidationErrors(w, validationErrors)
	}
	
	// Create phase object
	phase := &domain.Phase{
		Title:             req.Title,
		AssignedPlannerID: req.AssignedPlannerID,
		Status:            req.Status,
		Applications:      make([]domain.PhaseApplication, len(req.Applications)),
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
	
	// Set default status if not provided
	if phase.Status == "" {
		phase.Status = "draft"
	}
	
	// Generate phase ID
	phaseID, err := h.store.GenerateUniquePhaseID()
	if err != nil {
		log.Error("Failed to generate phase ID: %v", err)
		return fmt.Errorf("failed to generate phase ID: %w", err)
	}
	phase.PhaseID = phaseID
	
	// Convert application requests to domain objects
	for i, appReq := range req.Applications {
		app := domain.PhaseApplication{
			Type:              appReq.Type,
			InstitutionName:   appReq.InstitutionName,
			QualificationName: appReq.QualificationName,
			Expert1:           appReq.Expert1,
			Expert2:           appReq.Expert2,
			Status:            appReq.Status,
			CreatedAt:         time.Now().UTC(),
			UpdatedAt:         time.Now().UTC(),
		}
		
		// Set default status if not provided
		if app.Status == "" {
			app.Status = "pending"
		}
		
		phase.Applications[i] = app
	}
	
	// Create phase in store
	phaseIDInt, err := h.store.CreatePhase(phase)
	if err != nil {
		log.Error("Failed to create phase: %v", err)
		return fmt.Errorf("failed to create phase: %w", err)
	}
	
	// Get created phase to return
	createdPhase, err := h.store.GetPhase(phaseIDInt)
	if err != nil {
		log.Error("Failed to get created phase: %v", err)
		return fmt.Errorf("failed to get created phase: %w", err)
	}
	
	// Write standardized response
	return response.Success(w, http.StatusCreated, "Phase created successfully", createdPhase)
}

// updatePhaseRequest represents the request to update a phase
type updatePhaseRequest struct {
	Title             string `json:"title"`
	AssignedPlannerID int64 `json:"assignedPlannerId"`
	Status            string `json:"status"`
}

// HandleUpdatePhase handles PUT /api/phases/{id} requests
func (h *Handler) HandleUpdatePhase(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract phase ID from URL
	idStr := r.PathValue("id")
	if idStr == "" {
		return fmt.Errorf("phase ID is required")
	}
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Debug("Invalid phase ID: %v", err)
		return fmt.Errorf("invalid phase ID: %v", err)
	}
	
	// Get existing phase
	phase, err := h.store.GetPhase(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		log.Error("Failed to get phase: %v", err)
		return fmt.Errorf("failed to get phase: %w", err)
	}
	
	// Parse request body
	var req updatePhaseRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug("Failed to decode request body: %v", err)
		return fmt.Errorf("invalid request body: %v", err)
	}
	
	// Validate and update fields
	var validationErrors []string
	
	// Update title if provided
	if strings.TrimSpace(req.Title) != "" {
		phase.Title = req.Title
	}
	
	// Update assigned planner if provided
	if req.AssignedPlannerID > 0 && req.AssignedPlannerID != phase.AssignedPlannerID {
		// Verify assigned user exists and can be elevated to planner
		user, err := h.store.GetUser(req.AssignedPlannerID)
		if err != nil {
			if err == domain.ErrNotFound {
				validationErrors = append(validationErrors, fmt.Sprintf("user with ID %d does not exist", req.AssignedPlannerID))
			} else {
				log.Error("Failed to get assigned user: %v", err)
				return fmt.Errorf("failed to verify assigned user: %w", err)
			}
		} else if user.Role == "super_user" || user.Role == "admin" {
			// Admin and super_user roles don't need elevation assignments - they have inherent access
			log.Info("Assigning admin/super_user to phase - they have inherent planner access")
		} else if user.Role != "user" {
			validationErrors = append(validationErrors, fmt.Sprintf("user with ID %d has invalid role for elevation assignment", req.AssignedPlannerID))
		} else {
			phase.AssignedPlannerID = req.AssignedPlannerID
			phase.PlannerName = user.Name
		}
	}
	
	// Update status if provided
	if req.Status != "" {
		validStatuses := []string{"draft", "in_progress", "completed", "cancelled"}
		valid := false
		for _, s := range validStatuses {
			if req.Status == s {
				valid = true
				break
			}
		}
		if !valid {
			validationErrors = append(validationErrors, "status must be one of: draft, in_progress, completed, cancelled")
		} else {
			phase.Status = req.Status
		}
	}
	
	// Return validation errors if any
	if len(validationErrors) > 0 {
		log.Debug("Validation errors: %v", validationErrors)
		return respondWithValidationErrors(w, validationErrors)
	}
	
	// Update timestamp
	phase.UpdatedAt = time.Now().UTC()
	
	// Update phase in store
	err = h.store.UpdatePhase(phase)
	if err != nil {
		log.Error("Failed to update phase: %v", err)
		return fmt.Errorf("failed to update phase: %w", err)
	}
	
	// Get updated phase to return
	updatedPhase, err := h.store.GetPhase(id)
	if err != nil {
		log.Error("Failed to get updated phase: %v", err)
		return fmt.Errorf("failed to get updated phase: %w", err)
	}
	
	// Write standardized response
	return response.Success(w, http.StatusOK, "Phase updated successfully", updatedPhase)
}

// updateExpertsRequest represents the request to update the experts assigned to an application
type updateExpertsRequest struct {
	Expert1 int64 `json:"expert1"`
	Expert2 int64 `json:"expert2"`
}

// HandleUpdateApplicationExperts handles PUT /api/phases/{id}/applications/{app_id} requests
func (h *Handler) HandleUpdateApplicationExperts(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract phase and application IDs from URL
	phaseIDStr := r.PathValue("id")
	appIDStr := r.PathValue("app_id")
	
	if phaseIDStr == "" {
		return fmt.Errorf("phase ID is required")
	}
	
	if appIDStr == "" {
		return fmt.Errorf("application ID is required")
	}
	
	// Parse IDs
	phaseID, err := strconv.ParseInt(phaseIDStr, 10, 64)
	if err != nil {
		log.Debug("Invalid phase ID: %v", err)
		return fmt.Errorf("invalid phase ID: %v", err)
	}
	
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		log.Debug("Invalid application ID: %v", err)
		return fmt.Errorf("invalid application ID: %v", err)
	}
	
	// Verify phase exists
	_, err = h.store.GetPhase(phaseID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		log.Error("Failed to get phase: %v", err)
		return fmt.Errorf("failed to get phase: %w", err)
	}
	
	// Verify application exists and belongs to the phase
	app, err := h.store.GetPhaseApplication(appID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		log.Error("Failed to get application: %v", err)
		return fmt.Errorf("failed to get application: %w", err)
	}
	
	if app.PhaseID != phaseID {
		log.Debug("Application does not belong to phase")
		return fmt.Errorf("application does not belong to specified phase")
	}
	
	// Parse request body
	var req updateExpertsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug("Failed to decode request body: %v", err)
		return fmt.Errorf("invalid request body: %v", err)
	}
	
	// Validate experts
	var validationErrors []string
	
	// At least one expert must be provided
	if req.Expert1 <= 0 && req.Expert2 <= 0 {
		validationErrors = append(validationErrors, "at least one expert must be provided")
	}
	
	// If experts are provided, verify they exist
	if req.Expert1 > 0 {
		exists, err := expertExists(h.store, req.Expert1)
		if err != nil {
			log.Error("Failed to check if expert exists: %v", err)
			return fmt.Errorf("failed to verify expert: %w", err)
		}
		if !exists {
			validationErrors = append(validationErrors, fmt.Sprintf("expert 1 with ID %d does not exist", req.Expert1))
		}
	}
	
	if req.Expert2 > 0 {
		exists, err := expertExists(h.store, req.Expert2)
		if err != nil {
			log.Error("Failed to check if expert exists: %v", err)
			return fmt.Errorf("failed to verify expert: %w", err)
		}
		if !exists {
			validationErrors = append(validationErrors, fmt.Sprintf("expert 2 with ID %d does not exist", req.Expert2))
		}
	}
	
	// Return validation errors if any
	if len(validationErrors) > 0 {
		log.Debug("Validation errors: %v", validationErrors)
		return respondWithValidationErrors(w, validationErrors)
	}
	
	// Update experts in store
	err = h.store.UpdatePhaseApplicationExperts(appID, req.Expert1, req.Expert2)
	if err != nil {
		log.Error("Failed to update application experts: %v", err)
		return fmt.Errorf("failed to update application experts: %w", err)
	}
	
	// Get updated application to return
	updatedApp, err := h.store.GetPhaseApplication(appID)
	if err != nil {
		log.Error("Failed to get updated application: %v", err)
		return fmt.Errorf("failed to get updated application: %w", err)
	}
	
	// Write standardized response
	return response.Success(w, http.StatusOK, "Application experts updated successfully", updatedApp)
}

// applicationReviewRequest represents the request to review an application
type applicationReviewRequest struct {
	Action         string `json:"action"` // "approve" or "reject"
	RejectionNotes string `json:"rejectionNotes,omitempty"`
}

// HandleReviewApplication handles PUT /api/phases/{id}/applications/{app_id}/review requests
func (h *Handler) HandleReviewApplication(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract phase and application IDs from URL
	phaseIDStr := r.PathValue("id")
	appIDStr := r.PathValue("app_id")
	
	if phaseIDStr == "" {
		return fmt.Errorf("phase ID is required")
	}
	
	if appIDStr == "" {
		return fmt.Errorf("application ID is required")
	}
	
	// Parse IDs
	phaseID, err := strconv.ParseInt(phaseIDStr, 10, 64)
	if err != nil {
		log.Debug("Invalid phase ID: %v", err)
		return fmt.Errorf("invalid phase ID: %v", err)
	}
	
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		log.Debug("Invalid application ID: %v", err)
		return fmt.Errorf("invalid application ID: %v", err)
	}
	
	// Verify phase exists
	_, err = h.store.GetPhase(phaseID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		log.Error("Failed to get phase: %v", err)
		return fmt.Errorf("failed to get phase: %w", err)
	}
	
	// Verify application exists and belongs to the phase
	app, err := h.store.GetPhaseApplication(appID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		log.Error("Failed to get application: %v", err)
		return fmt.Errorf("failed to get application: %w", err)
	}
	
	if app.PhaseID != phaseID {
		log.Debug("Application does not belong to phase")
		return fmt.Errorf("application does not belong to specified phase")
	}
	
	// Parse request body
	var req applicationReviewRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug("Failed to decode request body: %v", err)
		return fmt.Errorf("invalid request body: %v", err)
	}
	
	// Validate action
	var validationErrors []string
	
	if req.Action != "approve" && req.Action != "reject" {
		validationErrors = append(validationErrors, "action must be either 'approve' or 'reject'")
	}
	
	// Application must be in "assigned" status to be reviewed
	if app.Status != "assigned" {
		validationErrors = append(validationErrors, "application must be in 'assigned' status to be reviewed")
	}
	
	// If rejecting, rejection notes are required
	if req.Action == "reject" && strings.TrimSpace(req.RejectionNotes) == "" {
		validationErrors = append(validationErrors, "rejection notes are required when rejecting an application")
	}
	
	// If approving, experts must be assigned
	if req.Action == "approve" && (app.Expert1 <= 0 && app.Expert2 <= 0) {
		validationErrors = append(validationErrors, "at least one expert must be assigned before approving")
	}
	
	// Return validation errors if any
	if len(validationErrors) > 0 {
		log.Debug("Validation errors: %v", validationErrors)
		return respondWithValidationErrors(w, validationErrors)
	}
	
	// Set status based on action
	var status string
	if req.Action == "approve" {
		status = "approved"
	} else {
		status = "rejected"
	}
	
	// Update application status in store
	err = h.store.UpdatePhaseApplicationStatus(appID, status, req.RejectionNotes)
	if err != nil {
		log.Error("Failed to update application status: %v", err)
		return fmt.Errorf("failed to update application status: %w", err)
	}
	
	// Get updated application to return
	updatedApp, err := h.store.GetPhaseApplication(appID)
	if err != nil {
		log.Error("Failed to get updated application: %v", err)
		return fmt.Errorf("failed to get updated application: %w", err)
	}
	
	// Write standardized response with dynamic message
	message := "Application " + req.Action + "ed successfully"
	return response.Success(w, http.StatusOK, message, updatedApp)
}

// Helper function to check if an expert exists
func expertExists(store storage.Storage, expertID int64) (bool, error) {
	expert, err := store.GetExpert(expertID)
	if err != nil {
		if err == domain.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return expert != nil, nil
}

// expertRatingRequest represents the request to rate experts in an application
type expertRatingRequest struct {
	ExpertID int    `json:"expertId"`
	Rating   int    `json:"rating"`   // 1-5 scale
	Comments string `json:"comments"` // Optional feedback
}

// HandleRateExperts handles POST /api/phases/{id}/applications/{app_id}/ratings requests (Manager access)
func (h *Handler) HandleRateExperts(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract phase and application IDs from URL
	phaseIDStr := r.PathValue("id")
	appIDStr := r.PathValue("app_id")
	
	if phaseIDStr == "" {
		return fmt.Errorf("phase ID is required")
	}
	
	if appIDStr == "" {
		return fmt.Errorf("application ID is required")
	}
	
	// Parse IDs
	phaseID, err := strconv.ParseInt(phaseIDStr, 10, 64)
	if err != nil {
		log.Debug("Invalid phase ID: %v", err)
		return fmt.Errorf("invalid phase ID: %v", err)
	}
	
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		log.Debug("Invalid application ID: %v", err)
		return fmt.Errorf("invalid application ID: %v", err)
	}
	
	// Verify application exists and belongs to the phase
	app, err := h.store.GetPhaseApplication(appID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrNotFound
		}
		log.Error("Failed to get application: %v", err)
		return fmt.Errorf("failed to get application: %w", err)
	}
	
	if app.PhaseID != phaseID {
		log.Debug("Application does not belong to phase")
		return fmt.Errorf("application does not belong to specified phase")
	}
	
	// Parse request body
	var req expertRatingRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug("Failed to decode request body: %v", err)
		return fmt.Errorf("invalid request body: %v", err)
	}
	
	// Validate rating
	var validationErrors []string
	
	if req.ExpertID <= 0 {
		validationErrors = append(validationErrors, "expert ID is required")
	}
	
	if req.Rating < 1 || req.Rating > 5 {
		validationErrors = append(validationErrors, "rating must be between 1 and 5")
	}
	
	// Verify expert exists and is assigned to this application
	expertAssigned := false
	if req.ExpertID == int(app.Expert1) || req.ExpertID == int(app.Expert2) {
		exists, err := expertExists(h.store, int64(req.ExpertID))
		if err != nil {
			log.Error("Failed to verify expert: %v", err)
			return fmt.Errorf("failed to verify expert: %w", err)
		}
		expertAssigned = exists
	}
	
	if !expertAssigned {
		validationErrors = append(validationErrors, "expert is not assigned to this application")
	}
	
	// Return validation errors if any
	if len(validationErrors) > 0 {
		log.Debug("Validation errors: %v", validationErrors)
		return respondWithValidationErrors(w, validationErrors)
	}
	
	// TODO: Store rating in application_ratings table (when implemented)
	// For now, just return success
	
	log.Info("Expert rated successfully: expert=%d, rating=%d, app=%d", req.ExpertID, req.Rating, appID)
	
	responseData := map[string]interface{}{
		"expertId": req.ExpertID,
		"rating":   req.Rating,
		"appId":    appID,
		"message":  "Rating will be implemented with application_ratings table",
	}
	
	// Write standardized response
	return response.Success(w, http.StatusOK, "Expert rating recorded successfully", responseData)
}

// managerTasksResponse represents the response for manager's pending tasks
type managerTasksResponse struct {
	ApplicationID       int64  `json:"applicationId"`
	PhaseID            int64  `json:"phaseId"`
	PhaseTitle         string `json:"phaseTitle"`
	InstitutionName    string `json:"institutionName"`
	QualificationName  string `json:"qualificationName"`
	Expert1ID          int64  `json:"expert1Id"`
	Expert1Name        string `json:"expert1Name"`
	Expert2ID          int64  `json:"expert2Id"`
	Expert2Name        string `json:"expert2Name"`
	Status             string `json:"status"`
	RatingRequested    bool   `json:"ratingRequested"`
}

// HandleGetManagerTasks handles GET /api/users/me/manager-tasks requests
func (h *Handler) HandleGetManagerTasks(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		return domain.ErrUnauthorized
	}
	
	// Get applications where user is assigned as manager
	managerApps, err := h.store.GetUserManagerApplications(int(userID))
	if err != nil {
		log.Error("Failed to get manager applications: %v", err)
		return fmt.Errorf("failed to get manager applications: %w", err)
	}
	
	// For each application, get the details
	tasks := make([]managerTasksResponse, 0, len(managerApps))
	
	for _, appID := range managerApps {
		app, err := h.store.GetPhaseApplication(int64(appID))
		if err != nil {
			log.Error("Failed to get application details: %v", err)
			continue // Skip this application but continue with others
		}
		
		// Get phase details
		phase, err := h.store.GetPhase(app.PhaseID)
		if err != nil {
			log.Error("Failed to get phase details: %v", err)
			continue
		}
		
		// TODO: Get expert names and rating request status from database
		// For now, create basic task info
		task := managerTasksResponse{
			ApplicationID:      app.ID,
			PhaseID:           app.PhaseID,
			PhaseTitle:        phase.Title,
			InstitutionName:   app.InstitutionName,
			QualificationName: app.QualificationName,
			Expert1ID:         app.Expert1,
			Expert1Name:       app.Expert1Name, // This should be resolved from DB
			Expert2ID:         app.Expert2,
			Expert2Name:       app.Expert2Name, // This should be resolved from DB
			Status:            app.Status,
			RatingRequested:   app.Status == "approved", // Simplification for now
		}
		
		tasks = append(tasks, task)
	}
	
	responseData := map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	}
	
	// Write standardized response
	return response.Success(w, http.StatusOK, "", responseData)
}

// HandleListApplications handles GET /api/applications requests with filtering
func (h *Handler) HandleListApplications(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Parse query parameters
	queryParams := r.URL.Query()
	phaseIDParam := queryParams.Get("phase_id")
	status := queryParams.Get("status")
	appType := queryParams.Get("type")
	expertIDParam := queryParams.Get("expert_id")
	
	limit := 100 // Default limit
	offset := 0  // Default offset
	
	// Parse pagination parameters
	if limitParam := queryParams.Get("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	if offsetParam := queryParams.Get("offset"); offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	var phaseID int64
	var expertID int64
	var err error
	
	// Parse phase ID if provided
	if phaseIDParam != "" {
		phaseID, err = strconv.ParseInt(phaseIDParam, 10, 64)
		if err != nil {
			log.Debug("Invalid phase_id parameter: %v", err)
			return fmt.Errorf("invalid phase_id parameter: %v", err)
		}
	}
	
	// Parse expert ID if provided
	if expertIDParam != "" {
		expertID, err = strconv.ParseInt(expertIDParam, 10, 64)
		if err != nil {
			log.Debug("Invalid expert_id parameter: %v", err)
			return fmt.Errorf("invalid expert_id parameter: %v", err)
		}
	}
	
	// For now, we'll get applications by phase ID
	// TODO: Implement more sophisticated filtering in storage layer
	
	var applications []domain.PhaseApplication
	
	if phaseID > 0 {
		apps, err := h.store.ListPhaseApplications(phaseID)
		if err != nil {
			log.Error("Failed to list applications: %v", err)
			return fmt.Errorf("failed to list applications: %w", err)
		}
		applications = apps
	} else {
		// If no phase ID specified, this would require a new storage method
		// For now, return empty results with appropriate message
		applications = []domain.PhaseApplication{}
		log.Debug("Application listing without phase_id not yet implemented")
	}
	
	// Apply additional filtering
	filteredApps := make([]domain.PhaseApplication, 0)
	for _, app := range applications {
		// Filter by status
		if status != "" && app.Status != status {
			continue
		}
		
		// Filter by type
		if appType != "" && app.Type != appType {
			continue
		}
		
		// Filter by expert ID
		if expertID > 0 && app.Expert1 != expertID && app.Expert2 != expertID {
			continue
		}
		
		filteredApps = append(filteredApps, app)
	}
	
	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(filteredApps) {
		start = len(filteredApps)
	}
	if end > len(filteredApps) {
		end = len(filteredApps)
	}
	
	paginatedApps := filteredApps[start:end]
	
	// Write standardized response
	responseData := map[string]interface{}{
		"applications": paginatedApps,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(paginatedApps),
			"total":  len(filteredApps),
		},
		"filters": map[string]interface{}{
			"phase_id":  phaseID,
			"status":    status,
			"type":      appType,
			"expert_id": expertID,
		},
	}
	
	return response.Success(w, http.StatusOK, "", responseData)
}

// Helper function to respond with validation errors
func respondWithValidationErrors(w http.ResponseWriter, errors []string) error {
	return response.ValidationError(w, errors)
}