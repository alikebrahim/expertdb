package phase

import (
	"encoding/json"
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
	
	// Parse scheduler ID
	var schedulerID int64
	schedulerIDParam := queryParams.Get("scheduler_id")
	if schedulerIDParam != "" {
		var err error
		schedulerID, err = strconv.ParseInt(schedulerIDParam, 10, 64)
		if err != nil {
			log.Debug("Invalid scheduler_id parameter: %v", err)
			return fmt.Errorf("invalid scheduler_id parameter: %v", err)
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
	phases, err := h.store.ListPhases(status, schedulerID, limit, offset)
	if err != nil {
		log.Error("Failed to list phases: %v", err)
		return fmt.Errorf("failed to list phases: %w", err)
	}
	
	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(phases)
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
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return json.NewEncoder(w).Encode(phase)
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
	
	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(phase)
}

// createPhaseRequest represents the request to create a new phase
type createPhaseRequest struct {
	Title             string                 `json:"title"`
	AssignedSchedulerID int64               `json:"assignedSchedulerId"`
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
	
	if req.AssignedSchedulerID <= 0 {
		validationErrors = append(validationErrors, "assigned scheduler ID is required")
	} else {
		// Verify scheduler exists and has scheduler role
		user, err := h.store.GetUser(req.AssignedSchedulerID)
		if err != nil {
			if err == domain.ErrNotFound {
				validationErrors = append(validationErrors, fmt.Sprintf("scheduler with ID %d does not exist", req.AssignedSchedulerID))
			} else {
				log.Error("Failed to get scheduler user: %v", err)
				return fmt.Errorf("failed to verify scheduler: %w", err)
			}
		} else if user.Role != "scheduler" {
			validationErrors = append(validationErrors, fmt.Sprintf("user with ID %d is not a scheduler", req.AssignedSchedulerID))
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
			validTypes := []string{"validation", "evaluation"}
			valid := false
			for _, t := range validTypes {
				if app.Type == t {
					valid = true
					break
				}
			}
			if !valid {
				validationErrors = append(validationErrors, fmt.Sprintf("application %d: type must be one of: validation, evaluation", i+1))
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
		AssignedSchedulerID: req.AssignedSchedulerID,
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
	phaseID, err = h.store.CreatePhase(phase)
	if err != nil {
		log.Error("Failed to create phase: %v", err)
		return fmt.Errorf("failed to create phase: %w", err)
	}
	
	// Get created phase to return
	createdPhase, err := h.store.GetPhase(phaseID)
	if err != nil {
		log.Error("Failed to get created phase: %v", err)
		return fmt.Errorf("failed to get created phase: %w", err)
	}
	
	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(createdPhase)
}

// updatePhaseRequest represents the request to update a phase
type updatePhaseRequest struct {
	Title             string `json:"title"`
	AssignedSchedulerID int64 `json:"assignedSchedulerId"`
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
	
	// Update scheduler if provided
	if req.AssignedSchedulerID > 0 && req.AssignedSchedulerID != phase.AssignedSchedulerID {
		// Verify scheduler exists and has scheduler role
		user, err := h.store.GetUser(req.AssignedSchedulerID)
		if err != nil {
			if err == domain.ErrNotFound {
				validationErrors = append(validationErrors, fmt.Sprintf("scheduler with ID %d does not exist", req.AssignedSchedulerID))
			} else {
				log.Error("Failed to get scheduler user: %v", err)
				return fmt.Errorf("failed to verify scheduler: %w", err)
			}
		} else if user.Role != "scheduler" {
			validationErrors = append(validationErrors, fmt.Sprintf("user with ID %d is not a scheduler", req.AssignedSchedulerID))
		} else {
			phase.AssignedSchedulerID = req.AssignedSchedulerID
			phase.SchedulerName = user.Name
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
	
	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(updatedPhase)
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
	
	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(updatedApp)
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
	
	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(updatedApp)
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

// Helper function to respond with validation errors
func respondWithValidationErrors(w http.ResponseWriter, errors []string) error {
	response := map[string]interface{}{
		"success": false,
		"errors":  errors,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(response)
}