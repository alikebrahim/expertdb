package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	
	"expertdb/internal/api/response"
	"expertdb/internal/api/utils"
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


// HandleCreateExpertRequest handles POST /api/expert-requests requests
func (h *ExpertRequestHandler) HandleCreateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()

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
	
	
	// Parse experience entries from JSON
	experienceJSON := r.FormValue("experienceEntries")
	if experienceJSON != "" {
		if err := json.Unmarshal([]byte(experienceJSON), &req.ExperienceEntries); err != nil {
			log.Warn("Failed to parse experience entries JSON: %v", err)
			return fmt.Errorf("invalid experience entries data: %w", err)
		}
	}
	
	// Parse education entries from JSON
	educationJSON := r.FormValue("educationEntries")
	if educationJSON != "" {
		if err := json.Unmarshal([]byte(educationJSON), &req.EducationEntries); err != nil {
			log.Warn("Failed to parse education entries JSON: %v", err)
			return fmt.Errorf("invalid education entries data: %w", err)
		}
	}
	
	// Parse specialized area IDs from JSON array
	if specAreaIdsJSON := r.FormValue("specializedAreaIds"); specAreaIdsJSON != "" {
		if err := json.Unmarshal([]byte(specAreaIdsJSON), &req.SpecializedAreaIds); err != nil {
			log.Warn("Failed to parse specialized area IDs JSON: %v", err)
			return fmt.Errorf("invalid specialized area IDs data: %w", err)
		}
	}
	
	// Parse suggested specialized areas from JSON array
	if suggestedAreasJSON := r.FormValue("suggestedSpecializedAreas"); suggestedAreasJSON != "" {
		if err := json.Unmarshal([]byte(suggestedAreasJSON), &req.SuggestedSpecializedAreas); err != nil {
			log.Warn("Failed to parse suggested specialized areas JSON: %v", err)
			return fmt.Errorf("invalid suggested specialized areas data: %w", err)
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

	// Get user ID from JWT context using the new utility
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		log.Warn("Failed to get user ID from request context")
		return err
	}

	// Convert specialized area IDs to comma-separated string for storage compatibility
	specializedAreaStr := ""
	if len(req.SpecializedAreaIds) > 0 {
		idStrings := make([]string, len(req.SpecializedAreaIds))
		for i, id := range req.SpecializedAreaIds {
			idStrings[i] = strconv.FormatInt(id, 10)
		}
		specializedAreaStr = strings.Join(idStrings, ",")
	}

	// Create ExpertRequest for storage (without file paths initially)
	expertRequest := &domain.ExpertRequest{
		Name:                      req.Name,
		Designation:               req.Designation,
		Affiliation:               req.Affiliation,
		Phone:                     req.Phone,
		Email:                     req.Email,
		IsBahraini:                req.IsBahraini,
		IsAvailable:               req.IsAvailable,
		Role:                      req.Role,
		EmploymentType:            req.EmploymentType,
		GeneralArea:               req.GeneralArea,
		SpecializedArea:           specializedAreaStr,
		SuggestedSpecializedAreas: req.SuggestedSpecializedAreas,
		IsTrained:                 req.IsTrained,
		IsPublished:               req.IsPublished,
		ExperienceEntries:         req.ExperienceEntries,
		EducationEntries:          req.EducationEntries,
		Status:                    "pending",
		CreatedAt:                 time.Now(),
		CreatedBy:                 userID,
	}

	// Transaction-based approach: Create request -> Upload CV -> Update paths
	
	// Step 1: Create expert request record and get ID
	requestID, err := h.store.CreateExpertRequestWithoutPaths(expertRequest)
	if err != nil {
		log.Error("Failed to create expert request: %v", err)
		userErr := errs.ParseSQLiteError(err, "expert request")
		return utils.RespondWithError(w, userErr)
	}

	// Step 2: Upload CV file using the actual request ID
	cvDoc, err := h.documentService.CreateDocumentForExpertRequest(requestID, cvFile, cvFileHeader)
	if err != nil {
		log.Error("Failed to upload CV for request %d: %v", requestID, err)
		// TODO: Consider rolling back the request creation here
		return fmt.Errorf("failed to upload CV: %w", err)
	}

	// Step 3: Document reference already updated by CreateDocumentForExpertRequest
	log.Debug("Expert request %d created successfully with CV document %d", requestID, cvDoc.ID)

	// Return success response
	return utils.RespondWithCreated(w, requestID, "Expert request created successfully")
}

// HandleGetExpertRequests handles GET /api/expert-requests requests
func (h *ExpertRequestHandler) HandleGetExpertRequests(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Get user context for access control using the new utility
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		log.Warn("Failed to get user ID from request context")
		return err
	}
	
	userRole, err := auth.GetUserRoleFromRequest(r)
	if err != nil {
		log.Warn("Failed to get user role from request context")
		return err
	}
	
	// Check if user is admin for access control
	isAdmin := userRole == "admin" || userRole == "super_user"
	
	// Parse query parameters for filtering
	status := r.URL.Query().Get("status")
	if status != "" {
		log.Debug("Filtering expert requests by status: %s", status)
	}
	
	// Parse pagination parameters using the new utility
	pagination := utils.ParsePaginationParams(r, 100) // Default limit of 100 for requests
	
	var requests []*domain.ExpertRequest
	
	if isAdmin {
		// Admin users can see all requests
		requests, err = h.store.ListExpertRequests(status, pagination.Limit, pagination.Offset)
	} else {
		// Regular users can only see their own requests
		requests, err = h.store.ListExpertRequestsByUser(userID, status, pagination.Limit, pagination.Offset)
	}
	
	if err != nil {
		log.Error("Failed to retrieve expert requests: %v", err)
		return fmt.Errorf("failed to retrieve expert requests: %w", err)
	}
	
	// Create response with pagination metadata
	responseData := map[string]interface{}{
		"requests": requests,
		"pagination": map[string]interface{}{
			"limit": pagination.Limit,
			"offset": pagination.Offset,
			"count": len(requests),
		},
	}
	
	return response.Success(w, http.StatusOK, "", responseData)
}

// HandleGetExpertRequest handles GET /api/expert-requests/{id} requests
func (h *ExpertRequestHandler) HandleGetExpertRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate expert request ID from path using the new utility
	id, err := utils.ExtractIDFromPath(r, "id", "expert request")
	if err != nil {
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
	return utils.RespondWithSuccess(w, "", request)
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
	// 1. Admins and super users can edit any request
	// 2. Regular users can edit only their own rejected requests
	isAdmin := role == auth.RoleAdmin || role == auth.RoleSuperUser
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
		
		// Parse multipart form (max 10MB for file)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Warn("Failed to parse multipart form: %v", err)
			return fmt.Errorf("failed to parse form: %w", err)
		}
		
		
		// Parse JSON part for the request data
		jsonData := r.FormValue("data")
		if jsonData == "" {
			// No JSON data provided, try to handle as simple form fields
			// Check for status field directly
			if status := r.FormValue("status"); status != "" {
				updateRequest.Status = status
			}
			if rejectionReason := r.FormValue("rejection_reason"); rejectionReason != "" {
				updateRequest.RejectionReason = rejectionReason
			}
			
			// If still no data, return error
			if updateRequest.Status == "" {
				log.Warn("Missing JSON data and status in form")
				return fmt.Errorf("missing request data")
			}
		} else {
			if err := json.Unmarshal([]byte(jsonData), &updateRequest); err != nil {
				log.Warn("Failed to parse JSON data: %v", err)
				return fmt.Errorf("invalid JSON data: %w", err)
			}
		}
		
		// Process CV file if provided
		cvFile, cvFileHeader, err := r.FormFile("cv")
		if err == nil {
			// CV file was provided, upload it using the request-specific method
			defer cvFile.Close()
			
			// Create CV document for the request (will be moved during approval)
			_, err := h.documentService.CreateDocumentForExpertRequest(id, cvFile, cvFileHeader)
			if err != nil {
				log.Error("Failed to upload updated CV: %v", err)
				return fmt.Errorf("failed to upload CV: %w", err)
			}
			// Document reference already updated by CreateDocumentForExpertRequest
		}
		
		// Process approval document if provided - store it for later use during approval
		approvalFile, approvalFileHeader, err := r.FormFile("approval_document")
		if err == nil {
			// Approval document was provided, store it properly as approval document
			defer approvalFile.Close()
			
			// Create approval document for the request
			doc, err := h.documentService.CreateApprovalDocumentForExpertRequest(id, approvalFile, approvalFileHeader)
			if err != nil {
				log.Error("Failed to upload approval document: %v", err)
				return fmt.Errorf("failed to upload approval document: %w", err)
			}
			
			// Update the request with the document ID
			updateRequest.ApprovalDocumentID = &doc.ID
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
		log.Info("Admin %d updating request %d status from '%s' to '%s'", userID, id, existingRequest.Status, updateRequest.Status)
		
		// If approving the request, require an approval document
		if updateRequest.Status == "approved" {
			// Check if the request has an approval document
			hasApprovalDoc := existingRequest.ApprovalDocumentID != nil || updateRequest.ApprovalDocumentID != nil
			
			if !hasApprovalDoc {
				log.Warn("Approval rejected for request %d: no approval document", id)
				return response.BadRequest(w, "Approval document is required before approving a request")
			}
		}
		
		
		if updateRequest.Status == "approved" {
			// Use the new method that handles approval document with proper expert ID
			expertID, err := h.store.ApproveExpertRequestWithDocument(id, userID, h.documentService)
			if err != nil {
				log.Error("Failed to approve expert request %d: %v", id, err)
				
				// Use the new error parser for user-friendly errors
				userErr := errs.ParseSQLiteError(err, "expert request")
				return utils.RespondWithError(w, userErr)
			}
			log.Info("Expert request %d approved - created expert %d", id, expertID)
		} else {
			// For rejection, use the old method
			if err := h.store.UpdateExpertRequestStatus(id, updateRequest.Status, updateRequest.RejectionReason, userID); err != nil {
				log.Error("Failed to update expert request %d status: %v", id, err)
				
				// Use the new error parser for user-friendly errors
				userErr := errs.ParseSQLiteError(err, "expert request")
				return utils.RespondWithError(w, userErr)
			}
		}
	} else {
		// If not a status update or if user is not admin, update the request fields
		
		// CRITICAL FIX: Prevent double request corruption
		// If all core fields are empty, this is likely an empty form submission from the second request
		// Reject it to prevent data corruption
		if updateRequest.Name == "" && updateRequest.Email == "" && updateRequest.Phone == "" && 
		   updateRequest.Designation == "" && updateRequest.Affiliation == "" && updateRequest.Status == "" {
			log.Warn("Rejected empty form submission for request ID: %d to prevent data corruption", id)
			return response.BadRequest(w, "Cannot process empty request data")
		}
		
		// Important: Preserve critical fields from the existing request
		updateRequest.CreatedBy = existingRequest.CreatedBy
		
		// If CV or approval document wasn't updated, keep the existing one
		if updateRequest.CVDocumentID == nil {
			updateRequest.CVDocumentID = existingRequest.CVDocumentID
		}
		if updateRequest.ApprovalDocumentID == nil {
			updateRequest.ApprovalDocumentID = existingRequest.ApprovalDocumentID
		}
		
		// Regular users shouldn't be able to change status
		if !isAdmin {
			updateRequest.Status = existingRequest.Status
			updateRequest.ReviewedAt = existingRequest.ReviewedAt
			updateRequest.ReviewedBy = existingRequest.ReviewedBy
		}
		
		if err := h.store.UpdateExpertRequest(&updateRequest); err != nil {
			log.Error("Failed to update expert request %d: %v", id, err)
			
			// Use the new error parser for user-friendly errors
			userErr := errs.ParseSQLiteError(err, "expert request")
			return utils.RespondWithError(w, userErr)
		}
	}
	
	// Return success response
	log.Info("Expert request %d updated successfully", id)
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
	
	// Call the storage method for batch approval with file moving
	log.Debug("Batch approving %d expert requests with file moving", len(batchRequest.RequestIDs))
	log.Debug("DEBUG: Request IDs to approve: %v", batchRequest.RequestIDs)
	log.Debug("DEBUG: Reviewer user ID: %d", userID)
	approved, expertIDs, errors := h.store.BatchApproveExpertRequestsWithFileMove(batchRequest.RequestIDs, userID, h.documentService)
	log.Debug("DEBUG: Batch approval completed - approved: %v, expertIDs: %v, errors: %v", approved, expertIDs, errors)
	if err != nil {
		log.Error("Failed during batch approval: %v", err)
		return fmt.Errorf("failed during batch approval: %w", err)
	}
	
	// Upload the approval document for the successfully approved experts
	var approvalDoc *domain.Document
	if len(expertIDs) > 0 {
		log.Debug("Creating approval document for %d experts", len(expertIDs))
		log.Debug("DEBUG: Expert IDs for approval document: %v", expertIDs)
		approvalDoc, err = h.documentService.CreateApprovalDocument(expertIDs, approvalFile, approvalFileHeader)
		log.Debug("DEBUG: Approval document created successfully: %+v", approvalDoc)
		if err != nil {
			log.Error("Failed to upload approval document: %v", err)
			return fmt.Errorf("failed to upload approval document: %w", err)
		}
		log.Debug("Approval document created: %s", approvalDoc.FilePath)
		
		// Update all approved experts with the approval document path
		log.Debug("DEBUG: Updating experts with approval document path: %s", approvalDoc.FilePath)
		err = h.store.UpdateExpertsApprovalPath(expertIDs, approvalDoc.FilePath)
		log.Debug("DEBUG: Experts approval path update completed with error: %v", err)
		if err != nil {
			log.Error("Failed to update experts with approval path: %v", err)
			return fmt.Errorf("failed to update experts with approval path: %w", err)
		}
	}
	
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

// HandleEditExpertRequest handles PUT /api/expert-requests/{id}/edit requests
func (h *ExpertRequestHandler) HandleEditExpertRequest(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Get user claims for authentication
	claims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		log.Warn("Failed to get user claims from context")
		return domain.ErrUnauthorized
	}
	
	userRole, _ := claims["role"].(string)
	isAdmin := userRole == "admin" || userRole == "super_user"
	
	// Extract request ID from path
	idStr := r.PathValue("id")
	requestID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid request ID provided: %s", idStr)
		return fmt.Errorf("invalid request ID: %w", err)
	}
	
	// Get existing request to check ownership and status
	existingRequest, err := h.store.GetExpertRequest(requestID)
	if err != nil {
		if err == domain.ErrNotFound {
			log.Warn("Expert request not found with ID: %d", requestID)
			return domain.ErrNotFound
		}
		log.Error("Failed to get expert request: %v", err)
		return fmt.Errorf("failed to retrieve expert request: %w", err)
	}
	
	// Check access permissions
	if !isAdmin {
		// Regular users can only edit their own rejected requests
		userID, _ := strconv.ParseInt(claims["sub"].(string), 10, 64)
		if existingRequest.CreatedBy != userID {
			log.Warn("User %d attempted to edit request %d not owned by them", userID, requestID)
			return domain.ErrForbidden
		}
		if existingRequest.Status != "rejected" {
			log.Warn("User %d attempted to edit request %d with status %s (only rejected allowed)", userID, requestID, existingRequest.Status)
			return fmt.Errorf("only rejected requests can be edited by users")
		}
	} else {
		// Admins can edit any pending request
		if existingRequest.Status != "pending" {
			log.Warn("Admin attempted to edit request %d with status %s (only pending allowed)", requestID, existingRequest.Status)
			return fmt.Errorf("only pending requests can be edited by admins")
		}
	}
	
	// Parse multipart form
	err = r.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		log.Warn("Failed to parse multipart form: %v", err)
		return fmt.Errorf("failed to parse form data: %w", err)
	}
	
	// Create updated request object based on existing data
	updatedRequest := *existingRequest
	
	// Update fields that are provided in the form
	if name := r.FormValue("name"); name != "" {
		updatedRequest.Name = name
	}
	if designation := r.FormValue("designation"); designation != "" {
		updatedRequest.Designation = designation
	}
	if affiliation := r.FormValue("affiliation"); affiliation != "" {
		updatedRequest.Affiliation = affiliation
	}
	if phone := r.FormValue("phone"); phone != "" {
		updatedRequest.Phone = phone
	}
	if email := r.FormValue("email"); email != "" {
		updatedRequest.Email = email
	}
	if isBahrainiStr := r.FormValue("isBahraini"); isBahrainiStr != "" {
		updatedRequest.IsBahraini = isBahrainiStr == "true"
	}
	if isAvailableStr := r.FormValue("isAvailable"); isAvailableStr != "" {
		updatedRequest.IsAvailable = isAvailableStr == "true"
	}
	if role := r.FormValue("role"); role != "" {
		updatedRequest.Role = role
	}
	if employmentType := r.FormValue("employmentType"); employmentType != "" {
		updatedRequest.EmploymentType = employmentType
	}
	if generalAreaStr := r.FormValue("generalArea"); generalAreaStr != "" {
		if generalArea, err := strconv.ParseInt(generalAreaStr, 10, 64); err == nil {
			updatedRequest.GeneralArea = generalArea
		}
	}
	if specializedAreaIds := r.FormValue("specializedAreaIds"); specializedAreaIds != "" {
		updatedRequest.SpecializedArea = specializedAreaIds
	}
	if suggestedAreasStr := r.FormValue("suggestedSpecializedAreas"); suggestedAreasStr != "" {
		var suggestedAreas []string
		if err := json.Unmarshal([]byte(suggestedAreasStr), &suggestedAreas); err == nil {
			updatedRequest.SuggestedSpecializedAreas = suggestedAreas
		}
	}
	if isTrainedStr := r.FormValue("isTrained"); isTrainedStr != "" {
		updatedRequest.IsTrained = isTrainedStr == "true"
	}
	
	// Handle CV file upload if provided
	if cvFile, cvHeader, err := r.FormFile("cv"); err == nil {
		defer cvFile.Close()
		
		// Upload new CV for expert
		_, err := h.documentService.CreateDocumentForExpert(requestID, cvFile, cvHeader, "cv")
		if err != nil {
			log.Error("Failed to upload CV: %v", err)
			return fmt.Errorf("failed to upload CV: %w", err)
		}
		// Document reference already updated by CreateDocumentForExpert
	}
	
	// Reset status to pending if user edited a rejected request
	if !isAdmin && existingRequest.Status == "rejected" {
		updatedRequest.Status = "pending"
		updatedRequest.RejectionReason = ""
		updatedRequest.ReviewedAt = time.Time{}
		updatedRequest.ReviewedBy = 0
	}
	
	// Update the request in storage
	err = h.store.UpdateExpertRequest(&updatedRequest)
	if err != nil {
		log.Error("Failed to update expert request: %v", err)
		return fmt.Errorf("failed to update expert request: %w", err)
	}
	
	log.Info("Expert request updated successfully: ID %d by user %s", requestID, claims["sub"])
	return response.Success(w, http.StatusOK, "Expert request updated successfully", nil)
}

// Use containsString from expert.go