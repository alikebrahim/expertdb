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

// Use existing writeJSON function from expert.go

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

	// Create CreateExpertRequest from form data
	var req domain.CreateExpertRequest
	
	// Parse individual form fields
	req.Name = r.FormValue("name")
	req.Designation = r.FormValue("designation")
	req.Affiliation = r.FormValue("affiliation")
	req.Phone = r.FormValue("phone")
	req.Email = r.FormValue("email")
	req.SpecializedArea = r.FormValue("specializedArea")
	req.Role = r.FormValue("role")
	req.EmploymentType = r.FormValue("employmentType")
	
	// Parse boolean fields
	req.IsBahraini, _ = strconv.ParseBool(r.FormValue("isBahraini"))
	req.IsAvailable, _ = strconv.ParseBool(r.FormValue("isAvailable"))
	req.IsTrained, _ = strconv.ParseBool(r.FormValue("isTrained"))
	req.IsPublished, _ = strconv.ParseBool(r.FormValue("isPublished"))
	
	// Parse numeric fields
	if generalAreaStr := r.FormValue("generalArea"); generalAreaStr != "" {
		req.GeneralArea, _ = strconv.ParseInt(generalAreaStr, 10, 64)
	}
	if ratingStr := r.FormValue("rating"); ratingStr != "" {
		req.Rating, _ = strconv.Atoi(ratingStr)
	}
	
	// Parse skills array from JSON
	if skillsJSON := r.FormValue("skills"); skillsJSON != "" {
		if err := json.Unmarshal([]byte(skillsJSON), &req.Skills); err != nil {
			log.Warn("Failed to parse skills JSON: %v", err)
			return fmt.Errorf("invalid skills data: %w", err)
		}
	}
	
	// Parse biography object from JSON
	if biographyJSON := r.FormValue("biography"); biographyJSON != "" {
		if err := json.Unmarshal([]byte(biographyJSON), &req.Biography); err != nil {
			log.Warn("Failed to parse biography JSON: %v", err)
			return fmt.Errorf("invalid biography data: %w", err)
		}
	}

	// Validate using the domain validation function
	if err := domain.ValidateCreateExpertRequest(&req); err != nil {
		log.Warn("Expert request validation failed: %v", err)
		return response.ValidationError(w, []string{err.Error()})
	}

	// Handle CV file - required
	cvFile, cvFileHeader, err := r.FormFile("cv")
	if err != nil {
		log.Warn("Failed to get CV file: %v", err)
		return fmt.Errorf("CV file is required: %w", err)
	}
	defer cvFile.Close()

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

	// Use a temporary negative ID to indicate this is for a request
	tempExpertID := int64(-1) 
	
	// Upload CV using document service
	cvDoc, err := h.documentService.CreateDocument(tempExpertID, cvFile, cvFileHeader, "cv")
	if err != nil {
		log.Error("Failed to upload CV: %v", err)
		return fmt.Errorf("failed to upload CV: %w", err)
	}
	req.CVPath = cvDoc.FilePath
	
	// Handle optional approval document file
	approvalFile, approvalFileHeader, err := r.FormFile("approval_document")
	var approvalDocPath string
	if err == nil {
		// Approval document was provided, upload it
		defer approvalFile.Close()
		
		approvalDoc, err := h.documentService.CreateDocument(tempExpertID, approvalFile, approvalFileHeader, "approval")
		if err != nil {
			log.Error("Failed to upload approval document: %v", err)
			return fmt.Errorf("failed to upload approval document: %w", err)
		}
		approvalDocPath = approvalDoc.FilePath
		log.Debug("Approval document uploaded successfully: %s", approvalDocPath)
	} else {
		log.Debug("No approval document provided (optional)")
	}

	// Convert biography struct to JSON string for storage
	biographyJSON := ""
	if len(req.Biography.Experience) > 0 || len(req.Biography.Education) > 0 {
		if data, err := json.Marshal(req.Biography); err == nil {
			biographyJSON = string(data)
		}
	}

	// Create ExpertRequest for storage (with additional fields)
	// Note: Temporarily mapping Affiliation to Institution until storage layer is updated
	expertRequest := &domain.ExpertRequest{
		Name:                 req.Name,
		Designation:          req.Designation,
		Affiliation:          req.Affiliation, // This maps to Institution in storage
		Phone:                req.Phone,
		Email:                req.Email,
		IsBahraini:           req.IsBahraini,
		IsAvailable:          req.IsAvailable,
		Rating:               req.Rating,
		Role:                 req.Role,
		EmploymentType:       req.EmploymentType,
		GeneralArea:          req.GeneralArea,
		SpecializedArea:      req.SpecializedArea,
		IsTrained:            req.IsTrained,
		IsPublished:          req.IsPublished,
		CVPath:               req.CVPath,
		ApprovalDocumentPath: approvalDocPath,
		Biography:            biographyJSON,
		Status:               "pending",
		CreatedAt:            time.Now(),
		CreatedBy:            userID,
	}

	// Create request in database
	log.Debug("Creating expert request for %s", req.Name)
	id, err := h.store.CreateExpertRequest(expertRequest)
	if err != nil {
		log.Error("Failed to create expert request: %v", err)
		
		// Use the new error parser for user-friendly errors
		userErr := errs.ParseSQLiteError(err, "expert request")
		return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": userErr.Error(),
			"suggestion": "Please check your input and try again",
		})
	}

	// Return success response
	log.Info("Expert request created successfully with ID: %d", id)
	responseData := map[string]interface{}{
		"id": id,
	}
	return response.Success(w, http.StatusCreated, "Expert request created successfully", responseData)
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
	
	// Create response with pagination metadata
	responseData := map[string]interface{}{
		"requests": requests,
		"pagination": map[string]interface{}{
			"limit": limit,
			"offset": offset,
			"count": len(requests),
		},
	}
	
	return response.Success(w, http.StatusOK, "", responseData)
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
	return response.Success(w, http.StatusOK, "", request)
}

// HandleUpdateExpertRequest handles PUT /api/expert-requests/{id} requests
func (h *ExpertRequestHandler) HandleUpdateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Get user claims for authentication
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		log.Warn("Failed to get user claims from context")
		return domain.ErrUnauthorized
	}
	
	// Extract user ID from claims
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
	
	// Get user role
	role, ok := claims["role"].(string)
	if !ok {
		log.Warn("Failed to get user role from context")
		return domain.ErrUnauthorized
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
	
	// Check permissions:
	// 1. Admins can edit any request
	// 2. Regular users can edit only their own rejected requests
	isAdmin := role == auth.RoleAdmin
	isOwner := existingRequest.CreatedBy == userID
	isRejected := existingRequest.Status == "rejected"
	
	if !isAdmin && !(isOwner && isRejected) {
		log.Warn("User %d attempted to update request %d without permission. Admin: %v, Owner: %v, Rejected: %v", 
			userID, id, isAdmin, isOwner, isRejected)
		return domain.ErrForbidden
	}
	
	// Check if this is a multipart form or JSON update
	contentType := r.Header.Get("Content-Type")
	var updateRequest domain.ExpertRequest
	
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// This is a file upload with form data
		log.Debug("Processing multipart form update for request ID: %d", id)
		
		// Parse multipart form (max 10MB for file)
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
		
		if err := json.Unmarshal([]byte(jsonData), &updateRequest); err != nil {
			log.Warn("Failed to parse JSON data: %v", err)
			return fmt.Errorf("invalid JSON data: %w", err)
		}
		
		// Process CV file if provided
		cvFile, cvFileHeader, err := r.FormFile("cv")
		if err == nil {
			// CV file was provided, upload it
			defer cvFile.Close()
			
			// Use a temporary negative ID to indicate this is for a request
			tempExpertID := int64(-1)
			
			cvDoc, err := h.documentService.CreateDocument(tempExpertID, cvFile, cvFileHeader, "cv")
			if err != nil {
				log.Error("Failed to upload updated CV: %v", err)
				return fmt.Errorf("failed to upload CV: %w", err)
			}
			updateRequest.CVPath = cvDoc.FilePath
			log.Debug("CV updated successfully for request ID %d: %s", id, updateRequest.CVPath)
		}
		
		// Process approval document if provided
		approvalFile, approvalFileHeader, err := r.FormFile("approval_document")
		if err == nil {
			// Approval document was provided, upload it
			defer approvalFile.Close()
			
			tempExpertID := int64(-1)
			approvalDoc, err := h.documentService.CreateDocument(tempExpertID, approvalFile, approvalFileHeader, "approval")
			if err != nil {
				log.Error("Failed to upload approval document: %v", err)
				return fmt.Errorf("failed to upload approval document: %w", err)
			}
			updateRequest.ApprovalDocumentPath = approvalDoc.FilePath
			log.Debug("Approval document updated successfully for request ID %d: %s", id, updateRequest.ApprovalDocumentPath)
		}
	} else {
		// This is a regular JSON update
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			log.Warn("Failed to parse expert request update: %v", err)
			return fmt.Errorf("invalid request body: %w", err)
		}
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
	
	// Perform status update if it's changing and user is admin
	if isAdmin && updateRequest.Status != "" && updateRequest.Status != existingRequest.Status {
		log.Debug("Admin updating expert request ID: %d, Status: %s", id, updateRequest.Status)
		
		// If approving the request, require an approval document
		if updateRequest.Status == "approved" {
			// Check if the request has an approval document
			hasApprovalDoc := existingRequest.ApprovalDocumentPath != "" || updateRequest.ApprovalDocumentPath != ""
			
			if !hasApprovalDoc {
				log.Warn("Attempted to approve request without approval document: %d", id)
				return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
					"error": "Approval document is required",
					"details": "An approval document must be uploaded before approving a request",
					"suggestion": "Please upload an approval document and try again",
				})
			}
			
			log.Debug("Approval document verified for request ID: %d", id)
		}
		
		if err := h.store.UpdateExpertRequestStatus(id, updateRequest.Status, updateRequest.RejectionReason, userID); err != nil {
			log.Error("Failed to update expert request status: %v", err)
			
			// Use the new error parser for user-friendly errors
			userErr := errs.ParseSQLiteError(err, "expert request")
			return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"error": userErr.Error(),
				"details": "There was an issue updating the expert request status",
				"suggestion": "Please check your input and try again",
			})
		}
	} else {
		// If not a status update or if user is not admin, update the request fields
		// Important: Preserve critical fields from the existing request
		updateRequest.CreatedBy = existingRequest.CreatedBy
		
		// If CV or approval document wasn't updated, keep the existing one
		if updateRequest.CVPath == "" {
			updateRequest.CVPath = existingRequest.CVPath
		}
		if updateRequest.ApprovalDocumentPath == "" {
			updateRequest.ApprovalDocumentPath = existingRequest.ApprovalDocumentPath
		}
		
		// Regular users shouldn't be able to change status
		if !isAdmin {
			updateRequest.Status = existingRequest.Status
			updateRequest.ReviewedAt = existingRequest.ReviewedAt
			updateRequest.ReviewedBy = existingRequest.ReviewedBy
		}
		
		log.Debug("Updating expert request ID: %d fields", id)
		if err := h.store.UpdateExpertRequest(&updateRequest); err != nil {
			log.Error("Failed to update expert request: %v", err)
			
			// Use the new error parser for user-friendly errors
			userErr := errs.ParseSQLiteError(err, "expert request")
			return writeJSON(w, http.StatusBadRequest, map[string]interface{}{
				"error": userErr.Error(),
				"details": "There was an issue updating the expert request",
				"suggestion": "Please check your input and try again",
			})
		}
	}
	
	// Return success response
	log.Info("Expert request updated successfully: ID: %d, Status: %s", id, updateRequest.Status)
	return response.Success(w, http.StatusOK, "Expert request updated successfully", nil)
}

// BatchApprovalRequest represents a request to approve multiple expert requests at once
type BatchApprovalRequest struct {
	RequestIDs []int64 `json:"requestIds"` // Array of expert request IDs to approve
}

// HandleBatchApproveExpertRequests handles POST /api/expert-requests/batch-approve requests
// This endpoint allows admins to approve multiple expert requests at once with a single approval document
func (h *ExpertRequestHandler) HandleBatchApproveExpertRequests(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/expert-requests/batch-approve request")
	
	// Get user claims for authentication
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		log.Warn("Failed to get user claims from context")
		return domain.ErrUnauthorized
	}
	
	// Only admins can perform batch approvals
	role, ok := claims["role"].(string)
	if !ok || role != auth.RoleAdmin {
		log.Warn("Non-admin attempted to perform batch approval")
		return domain.ErrForbidden
	}
	
	// Extract user ID from claims
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
	
	// Parse multipart form (max 10MB for file)
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
	
	var batchRequest BatchApprovalRequest
	if err := json.Unmarshal([]byte(jsonData), &batchRequest); err != nil {
		log.Warn("Failed to parse JSON data: %v", err)
		return fmt.Errorf("invalid JSON data: %w", err)
	}
	
	// Validate that at least one request ID is provided
	if len(batchRequest.RequestIDs) == 0 {
		log.Warn("No request IDs provided for batch approval")
		return fmt.Errorf("at least one request ID is required")
	}
	
	// Process approval document (required)
	approvalFile, approvalFileHeader, err := r.FormFile("approval_document")
	if err != nil {
		log.Warn("Failed to get approval document: %v", err)
		return fmt.Errorf("approval document is required: %w", err)
	}
	defer approvalFile.Close()
	
	// Upload the approval document
	tempExpertID := int64(-1) // Use a temporary negative ID
	approvalDoc, err := h.documentService.CreateDocument(tempExpertID, approvalFile, approvalFileHeader, "approval")
	if err != nil {
		log.Error("Failed to upload approval document: %v", err)
		return fmt.Errorf("failed to upload approval document: %w", err)
	}
	
	// Call the storage method for batch approval
	log.Debug("Batch approving %d expert requests", len(batchRequest.RequestIDs))
	approved, errors := h.store.BatchApproveExpertRequests(batchRequest.RequestIDs, approvalDoc.FilePath, userID)
	
	// Prepare response data
	responseData := map[string]interface{}{
		"totalRequests": len(batchRequest.RequestIDs),
		"approvedCount": len(approved),
		"approvedIds": approved,
	}
	
	if len(errors) > 0 {
		errorMessages := make(map[int64]string)
		for id, err := range errors {
			errorMessages[id] = err.Error()
		}
		responseData["errors"] = errorMessages
		responseData["errorCount"] = len(errors)
	}
	
	return response.Success(w, http.StatusOK, fmt.Sprintf("Approved %d of %d requests", len(approved), len(batchRequest.RequestIDs)), responseData)
}

// Use containsString from expert.go