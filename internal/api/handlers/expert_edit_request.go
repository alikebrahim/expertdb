package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"expertdb/internal/api/response"
	"expertdb/internal/auth"
	"expertdb/internal/documents"
	"expertdb/internal/domain"
	errs "expertdb/internal/errors"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// ExpertEditRequestHandler handles expert edit request-related API endpoints
type ExpertEditRequestHandler struct {
	store           storage.Storage
	documentService *documents.Service
}

// NewExpertEditRequestHandler creates a new expert edit request handler
func NewExpertEditRequestHandler(store storage.Storage, documentService *documents.Service) *ExpertEditRequestHandler {
	return &ExpertEditRequestHandler{
		store:           store,
		documentService: documentService,
	}
}

// HandleCreateExpertEditRequest handles POST /api/experts/{id}/edit-requests requests
func (h *ExpertEditRequestHandler) HandleCreateExpertEditRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract expert ID from URL path
	expertIDStr := r.PathValue("id")
	expertID, err := strconv.ParseInt(expertIDStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert ID provided: %s", expertIDStr)
		return fmt.Errorf("invalid expert ID: %w", err)
	}

	log.Debug("Processing POST /api/experts/%d/edit-requests request", expertID)

	// Verify expert exists
	existingExpert, err := h.store.GetExpert(expertID)
	if err != nil {
		if err == domain.ErrNotFound {
			log.Warn("Expert not found with ID: %d", expertID)
			return domain.ErrNotFound
		}
		log.Error("Failed to get expert: %v", err)
		return fmt.Errorf("failed to get expert: %w", err)
	}

	// Get user ID from JWT context
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		log.Warn("Failed to get user claims from context")
		return domain.ErrUnauthorized
	}

	var userID int64
	if sub, ok := claims["sub"].(string); ok {
		parsedID, err := strconv.ParseInt(sub, 10, 64)
		if err == nil {
			userID = parsedID
		} else {
			log.Warn("Failed to parse user ID from claims: %v", err)
			return domain.ErrUnauthorized
		}
	} else {
		log.Warn("Failed to get user ID from claims")
		return domain.ErrUnauthorized
	}

	// Parse the request
	var req domain.CreateExpertEditRequest

	// Check if this is multipart form data (for file uploads) or JSON
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form (max 10MB for files)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Warn("Failed to parse multipart form: %v", err)
			return fmt.Errorf("failed to parse form: %w", err)
		}

		// Parse JSON part for the request data
		jsonData := r.FormValue("data")
		if jsonData == "" {
			log.Warn("Missing JSON data in form")
			return fmt.Errorf("missing request data")
		}

		if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
			log.Warn("Failed to parse JSON data: %v", err)
			return fmt.Errorf("invalid JSON data: %w", err)
		}
	} else {
		// Parse JSON request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Warn("Failed to parse expert edit request: %v", err)
			return fmt.Errorf("invalid request body: %w", err)
		}
	}

	// Ensure expert ID matches URL parameter
	req.ExpertID = expertID

	// Validate the request
	if err := domain.ValidateCreateExpertEditRequest(&req); err != nil {
		log.Warn("Expert edit request validation failed: %v", err)
		return response.ValidationError(w, []string{err.Error()})
	}

	// Check if user can request edits for this expert
	role, _ := claims["role"].(string)
	if role != auth.RoleAdmin && role != auth.RoleSuperUser {
		// Regular users can only request edits for their own submissions or rejected requests
		// This is a simplified permission model - you might want to enhance this based on business rules
		log.Debug("Regular user requesting edit for expert ID: %d", expertID)
	}

	// Create the edit request domain object
	editRequest := &domain.ExpertEditRequest{
		ExpertID:                  req.ExpertID,
		Name:                     req.Name,
		Designation:              req.Designation,
		Institution:              req.Institution,
		Phone:                    req.Phone,
		Email:                    req.Email,
		IsBahraini:               req.IsBahraini,
		IsAvailable:              req.IsAvailable,
		Rating:                   req.Rating,
		Role:                     req.Role,
		EmploymentType:           req.EmploymentType,
		GeneralArea:              req.GeneralArea,
		IsTrained:                req.IsTrained,
		IsPublished:              req.IsPublished,
		Biography:                req.Biography,
		SuggestedSpecializedAreas: req.SuggestedSpecializedAreas,
		RemoveCV:                 req.RemoveCV,
		RemoveApprovalDocument:   req.RemoveApprovalDocument,
		ExperienceChanges:        req.ExperienceChanges,
		EducationChanges:         req.EducationChanges,
		ChangeSummary:            req.ChangeSummary,
		ChangeReason:             req.ChangeReason,
		Status:                   "pending",
		CreatedAt:                time.Now(),
		CreatedBy:                userID,
	}

	// Handle specialized areas
	if len(req.SpecializedAreaIds) > 0 {
		idStrings := make([]string, len(req.SpecializedAreaIds))
		for i, id := range req.SpecializedAreaIds {
			idStrings[i] = strconv.FormatInt(id, 10)
		}
		specializedAreaStr := strings.Join(idStrings, ",")
		editRequest.SpecializedArea = &specializedAreaStr
	}

	// Calculate which fields are being changed
	fieldsChanged := h.calculateChangedFields(existingExpert, editRequest)
	editRequest.FieldsChanged = fieldsChanged

	// Handle document uploads if multipart form
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle CV file upload
		cvFile, cvFileHeader, err := r.FormFile("cv")
		if err == nil {
			defer cvFile.Close()
			
			// Upload CV using document service  
			cvDoc, err := h.documentService.CreateDocument(expertID, cvFile, cvFileHeader, "cv")
			if err != nil {
				log.Error("Failed to upload CV: %v", err)
				return fmt.Errorf("failed to upload CV: %w", err)
			}
			editRequest.NewCVPath = &cvDoc.FilePath
			log.Debug("CV uploaded successfully for edit request: %s", cvDoc.FilePath)
		}

		// Handle approval document upload
		approvalFile, approvalFileHeader, err := r.FormFile("approval_document")
		if err == nil {
			defer approvalFile.Close()
			
			approvalDoc, err := h.documentService.CreateDocument(expertID, approvalFile, approvalFileHeader, "approval")
			if err != nil {
				log.Error("Failed to upload approval document: %v", err)
				return fmt.Errorf("failed to upload approval document: %w", err)
			}
			editRequest.NewApprovalDocumentPath = &approvalDoc.FilePath
			log.Debug("Approval document uploaded successfully for edit request: %s", approvalDoc.FilePath)
		}
	}

	// Create the edit request in the database
	log.Debug("Creating expert edit request for expert %s", existingExpert.Name)
	id, err := h.store.CreateExpertEditRequest(editRequest)
	if err != nil {
		log.Error("Failed to create expert edit request: %v", err)
		
		// Use the error parser for user-friendly errors
		userErr := errs.ParseSQLiteError(err, "expert edit request")
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error":      userErr.Error(),
			"suggestion": "Please check your input and try again",
		})
	}

	// Return success response
	log.Info("Expert edit request created successfully with ID: %d", id)
	responseData := map[string]interface{}{
		"id": id,
	}
	return response.Success(w, http.StatusCreated, "Expert edit request created successfully", responseData)
}

// HandleGetExpertEditRequests handles GET /api/expert-edit-requests requests
func (h *ExpertEditRequestHandler) HandleGetExpertEditRequests(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/expert-edit-requests request")

	// Parse query parameters for filtering
	filters := make(map[string]interface{})

	if expertID := r.URL.Query().Get("expertId"); expertID != "" {
		if id, err := strconv.ParseInt(expertID, 10, 64); err == nil {
			filters["expertId"] = id
		}
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
		log.Debug("Filtering expert edit requests by status: %s", status)
	}

	if createdBy := r.URL.Query().Get("createdBy"); createdBy != "" {
		if id, err := strconv.ParseInt(createdBy, 10, 64); err == nil {
			filters["createdBy"] = id
		}
	}

	// Parse pagination parameters
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 100 // Default limit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	// Get user context for permission filtering
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		return domain.ErrUnauthorized
	}

	role, _ := claims["role"].(string)
	if role != auth.RoleAdmin && role != auth.RoleSuperUser {
		// Regular users can only see their own requests
		if userIDStr, ok := claims["sub"].(string); ok {
			if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
				filters["createdBy"] = userID
			}
		}
	}

	// Retrieve edit requests from database
	log.Debug("Retrieving expert edit requests with filters: %v", filters)
	requests, err := h.store.ListExpertEditRequests(filters, limit, offset)
	if err != nil {
		log.Error("Failed to retrieve expert edit requests: %v", err)
		return fmt.Errorf("failed to retrieve expert edit requests: %w", err)
	}

	// Get total count for pagination
	totalCount, err := h.store.CountExpertEditRequests(filters)
	if err != nil {
		log.Error("Failed to count expert edit requests: %v", err)
		return fmt.Errorf("failed to count expert edit requests: %w", err)
	}

	// Create response with pagination metadata
	responseData := map[string]interface{}{
		"requests": requests,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(requests),
			"total":  totalCount,
		},
	}

	log.Debug("Returning %d expert edit requests", len(requests))
	return response.Success(w, http.StatusOK, "", responseData)
}

// HandleGetExpertEditRequest handles GET /api/expert-edit-requests/{id} requests
func (h *ExpertEditRequestHandler) HandleGetExpertEditRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate edit request ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert edit request ID provided: %s", idStr)
		return fmt.Errorf("invalid request ID: %w", err)
	}

	// Retrieve edit request from database
	log.Debug("Retrieving expert edit request with ID: %d", id)
	request, err := h.store.GetExpertEditRequest(id)
	if err != nil {
		if err == domain.ErrNotFound {
			log.Warn("Expert edit request not found with ID: %d", id)
			return domain.ErrNotFound
		}

		log.Error("Failed to get expert edit request: %v", err)
		return fmt.Errorf("failed to retrieve expert edit request: %w", err)
	}

	// Check permissions
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		return domain.ErrUnauthorized
	}

	role, _ := claims["role"].(string)
	if role != auth.RoleAdmin && role != auth.RoleSuperUser {
		// Regular users can only view their own requests
		userIDStr, _ := claims["sub"].(string)
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)
		if request.CreatedBy != userID {
			log.Warn("User %d attempted to view edit request %d without permission", userID, id)
			return domain.ErrForbidden
		}
	}

	// Return edit request data
	log.Debug("Successfully retrieved expert edit request: ID: %d, Expert: %s", request.ID, request.ExpertName)
	return response.Success(w, http.StatusOK, "", request)
}

// HandleUpdateExpertEditRequestStatus handles PUT /api/expert-edit-requests/{id}/status requests
func (h *ExpertEditRequestHandler) HandleUpdateExpertEditRequestStatus(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate edit request ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert edit request ID provided: %s", idStr)
		return fmt.Errorf("invalid request ID: %w", err)
	}

	// Parse status update request
	var statusUpdate struct {
		Status          string `json:"status"`          // "approved", "rejected", "cancelled"
		RejectionReason string `json:"rejectionReason,omitempty"` // Required if status is "rejected"
		AdminNotes      string `json:"adminNotes,omitempty"`     // Optional admin notes
	}

	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		log.Warn("Failed to parse status update request: %v", err)
		return fmt.Errorf("invalid request body: %w", err)
	}

	// Validate status
	validStatuses := []string{"approved", "rejected", "cancelled"}
	if !contains(validStatuses, statusUpdate.Status) {
		log.Warn("Invalid status provided: %s", statusUpdate.Status)
		return fmt.Errorf("invalid status: %s", statusUpdate.Status)
	}

	// Get user context
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		return domain.ErrUnauthorized
	}

	userIDStr, _ := claims["sub"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return domain.ErrUnauthorized
	}

	role, _ := claims["role"].(string)

	// Check permissions
	if statusUpdate.Status == "cancelled" {
		// Users can cancel their own pending requests
		editRequest, err := h.store.GetExpertEditRequest(id)
		if err != nil {
			return err
		}

		if editRequest.CreatedBy != userID || editRequest.Status != "pending" {
			return domain.ErrForbidden
		}

		return h.store.CancelExpertEditRequest(id, userID)
	}

	// Only admins can approve/reject requests
	if role != auth.RoleAdmin && role != auth.RoleSuperUser {
		log.Warn("Non-admin user %d attempted to update edit request status", userID)
		return domain.ErrForbidden
	}

	// Validate rejection reason
	if statusUpdate.Status == "rejected" && strings.TrimSpace(statusUpdate.RejectionReason) == "" {
		return fmt.Errorf("rejection reason is required when rejecting a request")
	}

	// Update status in database
	log.Debug("Updating expert edit request status: ID: %d, Status: %s", id, statusUpdate.Status)
	err = h.store.UpdateExpertEditRequestStatus(id, statusUpdate.Status, statusUpdate.RejectionReason, statusUpdate.AdminNotes, userID)
	if err != nil {
		log.Error("Failed to update expert edit request status: %v", err)
		return fmt.Errorf("failed to update status: %w", err)
	}

	log.Info("Expert edit request status updated: ID: %d, Status: %s", id, statusUpdate.Status)
	return response.Success(w, http.StatusOK, "Status updated successfully", nil)
}

// HandleApplyExpertEditRequest handles POST /api/expert-edit-requests/{id}/apply requests
func (h *ExpertEditRequestHandler) HandleApplyExpertEditRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

	// Extract and validate edit request ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert edit request ID provided: %s", idStr)
		return fmt.Errorf("invalid request ID: %w", err)
	}

	// Get user context
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		return domain.ErrUnauthorized
	}

	role, _ := claims["role"].(string)
	if role != auth.RoleAdmin && role != auth.RoleSuperUser {
		log.Warn("Non-admin user attempted to apply edit request")
		return domain.ErrForbidden
	}

	userIDStr, _ := claims["sub"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return domain.ErrUnauthorized
	}

	// Apply the edit request
	log.Debug("Applying expert edit request: ID: %d", id)
	err = h.store.ApplyExpertEditRequest(id, userID)
	if err != nil {
		log.Error("Failed to apply expert edit request: %v", err)
		return fmt.Errorf("failed to apply edit request: %w", err)
	}

	log.Info("Expert edit request applied successfully: ID: %d", id)
	return response.Success(w, http.StatusOK, "Edit request applied successfully", nil)
}

// Helper functions

func (h *ExpertEditRequestHandler) calculateChangedFields(existing *domain.Expert, editRequest *domain.ExpertEditRequest) []string {
	var changedFields []string

	if editRequest.Name != nil && *editRequest.Name != existing.Name {
		changedFields = append(changedFields, "name")
	}
	if editRequest.Designation != nil && *editRequest.Designation != existing.Designation {
		changedFields = append(changedFields, "designation")
	}
	if editRequest.Institution != nil && *editRequest.Institution != existing.Affiliation {
		changedFields = append(changedFields, "institution")
	}
	if editRequest.Phone != nil && *editRequest.Phone != existing.Phone {
		changedFields = append(changedFields, "phone")
	}
	if editRequest.Email != nil && *editRequest.Email != existing.Email {
		changedFields = append(changedFields, "email")
	}
	if editRequest.IsBahraini != nil && *editRequest.IsBahraini != existing.IsBahraini {
		changedFields = append(changedFields, "isBahraini")
	}
	if editRequest.IsAvailable != nil && *editRequest.IsAvailable != existing.IsAvailable {
		changedFields = append(changedFields, "isAvailable")
	}
	if editRequest.Rating != nil && *editRequest.Rating != existing.Rating {
		changedFields = append(changedFields, "rating")
	}
	if editRequest.Role != nil && *editRequest.Role != existing.Role {
		changedFields = append(changedFields, "role")
	}
	if editRequest.EmploymentType != nil && *editRequest.EmploymentType != existing.EmploymentType {
		changedFields = append(changedFields, "employmentType")
	}
	if editRequest.GeneralArea != nil && *editRequest.GeneralArea != existing.GeneralArea {
		changedFields = append(changedFields, "generalArea")
	}
	if editRequest.SpecializedArea != nil && *editRequest.SpecializedArea != existing.SpecializedArea {
		changedFields = append(changedFields, "specializedArea")
	}
	if editRequest.IsTrained != nil && *editRequest.IsTrained != existing.IsTrained {
		changedFields = append(changedFields, "isTrained")
	}
	if editRequest.IsPublished != nil && *editRequest.IsPublished != existing.IsPublished {
		changedFields = append(changedFields, "isPublished")
	}

	if editRequest.NewCVPath != nil {
		changedFields = append(changedFields, "cvPath")
	}
	if editRequest.RemoveCV {
		changedFields = append(changedFields, "removeCv")
	}
	if editRequest.NewApprovalDocumentPath != nil {
		changedFields = append(changedFields, "approvalDocumentPath")
	}
	if editRequest.RemoveApprovalDocument {
		changedFields = append(changedFields, "removeApprovalDocument")
	}

	if len(editRequest.ExperienceChanges) > 0 {
		changedFields = append(changedFields, "experience")
	}
	if len(editRequest.EducationChanges) > 0 {
		changedFields = append(changedFields, "education")
	}
	if len(editRequest.SuggestedSpecializedAreas) > 0 {
		changedFields = append(changedFields, "suggestedSpecializedAreas")
	}

	return changedFields
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}