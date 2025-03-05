package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type APIServer struct {
	listenAddr      string
	store           Storage
	documentService *DocumentService
	aiService       *AIService
	config          *Configuration
}

func NewAPIServer(listenAddr string, store Storage, config *Configuration) (*APIServer, error) {
	// Initialize document service
	documentService, err := NewDocumentService(store, config.UploadPath)
	if err != nil {
		return nil, err
	}
	
	// Initialize AI service
	aiService := NewAIService(config.AIServiceURL, store)
	
	return &APIServer{
		listenAddr:      listenAddr,
		store:           store,
		documentService: documentService,
		aiService:       aiService,
		config:          config,
	}, nil
}

func (s *APIServer) Run() error {
	logger := GetLogger()
	logger.Info("Setting up API routes...")
	
	mux := http.NewServeMux()

	// Expert routes
	mux.HandleFunc("GET /api/experts", makeHTTPHandleFunc(s.handleGetExperts))
	mux.HandleFunc("POST /api/experts", makeHTTPHandleFunc(requireAdmin(s.handleCreateExpert)))
	mux.HandleFunc("GET /api/experts/{id}", makeHTTPHandleFunc(s.handleGetExpert))
	mux.HandleFunc("PUT /api/experts/{id}", makeHTTPHandleFunc(requireAdmin(s.handleUpdateExpert)))
	mux.HandleFunc("DELETE /api/experts/{id}", makeHTTPHandleFunc(requireAdmin(s.handleDeleteExpert)))
	
	// Expert request routes
	mux.HandleFunc("POST /api/expert-requests", makeHTTPHandleFunc(requireAuth(s.handleCreateExpertRequest)))
	mux.HandleFunc("GET /api/expert-requests", makeHTTPHandleFunc(requireAuth(s.handleGetExpertRequests)))
	mux.HandleFunc("GET /api/expert-requests/{id}", makeHTTPHandleFunc(requireAuth(s.handleGetExpertRequest)))
	mux.HandleFunc("PUT /api/expert-requests/{id}", makeHTTPHandleFunc(requireAdmin(s.handleUpdateExpertRequest)))

	// ISCED reference data routes
	mux.HandleFunc("GET /api/isced/levels", makeHTTPHandleFunc(s.handleGetISCEDLevels))
	mux.HandleFunc("GET /api/isced/fields", makeHTTPHandleFunc(s.handleGetISCEDFields))
	
	// Document routes
	mux.HandleFunc("POST /api/documents", makeHTTPHandleFunc(requireAuth(s.handleUploadDocument)))
	mux.HandleFunc("GET /api/documents/{id}", makeHTTPHandleFunc(requireAuth(s.handleGetDocument)))
	mux.HandleFunc("DELETE /api/documents/{id}", makeHTTPHandleFunc(requireAdmin(s.handleDeleteDocument)))
	mux.HandleFunc("GET /api/experts/{id}/documents", makeHTTPHandleFunc(requireAuth(s.handleGetExpertDocuments)))
	
	// Engagement routes
	mux.HandleFunc("POST /api/engagements", makeHTTPHandleFunc(requireAdmin(s.handleCreateEngagement)))
	mux.HandleFunc("GET /api/engagements/{id}", makeHTTPHandleFunc(requireAuth(s.handleGetEngagement)))
	mux.HandleFunc("PUT /api/engagements/{id}", makeHTTPHandleFunc(requireAdmin(s.handleUpdateEngagement)))
	mux.HandleFunc("DELETE /api/engagements/{id}", makeHTTPHandleFunc(requireAdmin(s.handleDeleteEngagement)))
	mux.HandleFunc("GET /api/experts/{id}/engagements", makeHTTPHandleFunc(requireAuth(s.handleGetExpertEngagements)))
	
	// AI integration routes
	mux.HandleFunc("POST /api/ai/generate-profile", makeHTTPHandleFunc(requireAdmin(s.handleGenerateProfile)))
	mux.HandleFunc("POST /api/ai/suggest-isced", makeHTTPHandleFunc(requireAdmin(s.handleSuggestISCED)))
	mux.HandleFunc("POST /api/ai/extract-skills", makeHTTPHandleFunc(requireAdmin(s.handleExtractSkills)))
	mux.HandleFunc("POST /api/ai/suggest-panel", makeHTTPHandleFunc(requireAuth(s.handleSuggestPanel)))
	
	// Statistics routes
	mux.HandleFunc("GET /api/statistics", makeHTTPHandleFunc(requireAuth(s.handleGetStatistics)))
	mux.HandleFunc("GET /api/statistics/nationality", makeHTTPHandleFunc(requireAuth(s.handleGetNationalityStats)))
	mux.HandleFunc("GET /api/statistics/isced", makeHTTPHandleFunc(requireAuth(s.handleGetISCEDStats)))
	mux.HandleFunc("GET /api/statistics/engagements", makeHTTPHandleFunc(requireAuth(s.handleGetEngagementStats)))
	mux.HandleFunc("GET /api/statistics/growth", makeHTTPHandleFunc(requireAuth(s.handleGetGrowthStats)))
	
	// User management routes
	mux.HandleFunc("POST /api/auth/login", makeHTTPHandleFunc(s.handleLogin))
	mux.HandleFunc("POST /api/users", makeHTTPHandleFunc(requireAdmin(s.handleCreateUser)))
	mux.HandleFunc("GET /api/users", makeHTTPHandleFunc(requireAdmin(s.handleGetUsers)))
	mux.HandleFunc("GET /api/users/{id}", makeHTTPHandleFunc(requireAuth(s.handleGetUser)))
	mux.HandleFunc("PUT /api/users/{id}", makeHTTPHandleFunc(requireAdmin(s.handleUpdateUser)))
	mux.HandleFunc("DELETE /api/users/{id}", makeHTTPHandleFunc(requireAdmin(s.handleDeleteUser)))
	
	// Add CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", s.config.CORSAllowOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			// Proceed with the request
			next.ServeHTTP(w, r)
		})
	}
	
	// Get logger for request logging
	apiLogger := GetLogger()
	
	// Apply middleware chain: request logging -> CORS
	handler := apiLogger.RequestLoggerMiddleware(corsMiddleware(mux))

	// Start server
	apiLogger.Info("API server listening on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, handler)
}

func (s *APIServer) handleGetExperts(w http.ResponseWriter, r *http.Request) error {
	// Parse query parameters for filtering
	queryParams := r.URL.Query()
	filters := make(map[string]interface{})

	// Add supported filters
	if name := queryParams.Get("name"); name != "" {
		filters["name"] = name
	}
	if area := queryParams.Get("area"); area != "" {
		filters["area"] = area
	}
	if available := queryParams.Get("is_available"); available != "" {
		if available == "true" {
			filters["is_available"] = true
		} else if available == "false" {
			filters["is_available"] = false
		}
	}
	if role := queryParams.Get("role"); role != "" {
		filters["role"] = role
	}
	if iscedLevelID := queryParams.Get("isced_level_id"); iscedLevelID != "" {
		if id, err := strconv.ParseInt(iscedLevelID, 10, 64); err == nil {
			filters["isced_level_id"] = id
		}
	}
	if iscedFieldID := queryParams.Get("isced_field_id"); iscedFieldID != "" {
		if id, err := strconv.ParseInt(iscedFieldID, 10, 64); err == nil {
			filters["isced_field_id"] = id
		}
	}
	
	// Add rating filter
	if minRating := queryParams.Get("min_rating"); minRating != "" {
		if rating, err := strconv.ParseFloat(minRating, 64); err == nil {
			filters["min_rating"] = rating
		}
	}

	// Parse pagination parameters
	limit := 10 // Default to 10 per page
	offset := 0
	if limitParam := queryParams.Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if offsetParam := queryParams.Get("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	// Parse sorting parameters
	sortBy := "name"
	sortOrder := "asc"
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

	// Get total count for pagination
	
	// Create a copy of filters without pagination
	countFilters := make(map[string]interface{})
	for k, v := range filters {
		if k != "sort_by" && k != "sort_order" {
			countFilters[k] = v
		}
	}
	
	// Get total count
	totalCount, err := s.store.CountExperts(countFilters)
	if err != nil {
		return err
	}
	
	// Get filtered experts
	experts, err := s.store.ListExperts(filters, limit, offset)
	if err != nil {
		return err
	}
	
	// Set total count header for pagination
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", totalCount))

	return WriteJson(w, http.StatusOK, experts)
}

func (s *APIServer) handleGetExpert(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from path using native Go 1.22+ pattern matching
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}

	expert, err := s.store.GetExpert(id)
	if err != nil {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "Expert not found"})
	}

	return WriteJson(w, http.StatusOK, expert)
}

func (s *APIServer) handleCreateExpert(w http.ResponseWriter, r *http.Request) error {
	createReq := new(CreateExpertRequest)
	if err := json.NewDecoder(r.Body).Decode(createReq); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}

	// Validate request
	if err := ValidateCreateExpertRequest(createReq); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	// Create expert
	expert := NewExpert(*createReq)
	id, err := s.store.CreateExpert(expert)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Failed to create expert"})
	}

	return WriteJson(w, http.StatusCreated, CreateExpertResponse{
		ID:      id,
		Success: true,
		Message: "Expert created successfully",
	})
}

func (s *APIServer) handleUpdateExpert(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from path using native Go 1.22+ pattern matching
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}

	// Check if expert exists
	_, err = s.store.GetExpert(id)
	if err != nil {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "Expert not found"})
	}

	// Parse update request
	var updateExpert Expert
	if err := json.NewDecoder(r.Body).Decode(&updateExpert); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}

	// Ensure ID matches
	updateExpert.ID = id

	// Update expert
	if err := s.store.UpdateExpert(&updateExpert); err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Failed to update expert"})
	}

	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Expert updated successfully",
	})
}

func (s *APIServer) handleDeleteExpert(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from path using native Go 1.22+ pattern matching
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}

	if err := s.store.DeleteExpert(id); err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Failed to delete expert"})
	}

	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Expert deleted successfully",
	})
}

func (s *APIServer) handleGetISCEDLevels(w http.ResponseWriter, r *http.Request) error {
	// Query ISCED levels from database
	levels, err := s.store.GetISCEDLevels()
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Failed to fetch ISCED levels"})
	}

	return WriteJson(w, http.StatusOK, levels)
}

func (s *APIServer) handleGetISCEDFields(w http.ResponseWriter, r *http.Request) error {
	// Query ISCED fields from database
	fields, err := s.store.GetISCEDFields()
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Failed to fetch ISCED fields"})
	}

	return WriteJson(w, http.StatusOK, fields)
}

// Document handlers

func (s *APIServer) handleUploadDocument(w http.ResponseWriter, r *http.Request) error {
	// Parse the multipart form data, 10 MB max
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Failed to parse form"})
	}
	
	// Get expert ID
	expertIDStr := r.FormValue("expertId")
	if expertIDStr == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Expert ID is required"})
	}
	
	expertID, err := strconv.ParseInt(expertIDStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}
	
	// Get document type
	docType := r.FormValue("documentType")
	if docType == "" {
		docType = "cv" // Default type
	}
	
	// Get the file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "No file provided"})
	}
	defer file.Close()
	
	// Upload and store the document
	doc, err := s.documentService.CreateDocument(expertID, file, header, docType)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to upload document: %v", err)})
	}
	
	return WriteJson(w, http.StatusCreated, doc)
}

func (s *APIServer) handleGetDocument(w http.ResponseWriter, r *http.Request) error {
	// Extract document ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid document ID"})
	}
	
	// Get the document
	doc, err := s.documentService.GetDocument(id)
	if err != nil {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "Document not found"})
	}
	
	return WriteJson(w, http.StatusOK, doc)
}

func (s *APIServer) handleDeleteDocument(w http.ResponseWriter, r *http.Request) error {
	// Extract document ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid document ID"})
	}
	
	// Delete the document
	if err := s.documentService.DeleteDocument(id); err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to delete document: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Document deleted successfully",
	})
}

func (s *APIServer) handleGetExpertDocuments(w http.ResponseWriter, r *http.Request) error {
	// Extract expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}
	
	// Get the expert's documents
	docs, err := s.documentService.GetDocumentsByExpertID(id)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve documents: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, docs)
}

// Engagement handlers

func (s *APIServer) handleCreateEngagement(w http.ResponseWriter, r *http.Request) error {
	var engagement Engagement
	if err := json.NewDecoder(r.Body).Decode(&engagement); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Set creation time
	engagement.CreatedAt = time.Now()
	
	// Validate required fields
	if engagement.ExpertID == 0 {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Expert ID is required"})
	}
	if engagement.EngagementType == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Engagement type is required"})
	}
	if engagement.StartDate.IsZero() {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Start date is required"})
	}
	if engagement.Status == "" {
		engagement.Status = "pending" // Default status
	}
	
	// Create the engagement
	id, err := s.store.CreateEngagement(&engagement)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to create engagement: %v", err)})
	}
	
	// Set the ID in the response
	engagement.ID = id
	
	return WriteJson(w, http.StatusCreated, engagement)
}

func (s *APIServer) handleGetEngagement(w http.ResponseWriter, r *http.Request) error {
	// Extract engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid engagement ID"})
	}
	
	// Get the engagement
	engagement, err := s.store.GetEngagement(id)
	if err != nil {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "Engagement not found"})
	}
	
	return WriteJson(w, http.StatusOK, engagement)
}

func (s *APIServer) handleUpdateEngagement(w http.ResponseWriter, r *http.Request) error {
	// Extract engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid engagement ID"})
	}
	
	// Get existing engagement
	existing, err := s.store.GetEngagement(id)
	if err != nil {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "Engagement not found"})
	}
	
	// Parse update request
	var updateEngagement Engagement
	if err := json.NewDecoder(r.Body).Decode(&updateEngagement); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Ensure ID matches
	updateEngagement.ID = id
	
	// Update only the fields that are provided in the request
	// Maintaining original values for fields not in the update
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
	
	// Update the engagement
	if err := s.store.UpdateEngagement(&updateEngagement); err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to update engagement: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Engagement updated successfully",
	})
}

func (s *APIServer) handleDeleteEngagement(w http.ResponseWriter, r *http.Request) error {
	// Extract engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid engagement ID"})
	}
	
	// Delete the engagement
	if err := s.store.DeleteEngagement(id); err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to delete engagement: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Engagement deleted successfully",
	})
}

func (s *APIServer) handleGetExpertEngagements(w http.ResponseWriter, r *http.Request) error {
	// Extract expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}
	
	// Get the expert's engagements
	engagements, err := s.store.GetEngagementsByExpertID(id)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve engagements: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, engagements)
}

// AI Integration handlers

func (s *APIServer) handleGenerateProfile(w http.ResponseWriter, r *http.Request) error {
	var request struct {
		ExpertID int64 `json:"expertId"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Check if expert exists
	_, expert_err := s.store.GetExpert(request.ExpertID)
	if expert_err != nil {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "Expert not found"})
	}
	
	// Generate the profile using the Storage interface
	result, err := s.store.GenerateProfile(request.ExpertID)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to generate profile: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, result)
}

func (s *APIServer) handleSuggestISCED(w http.ResponseWriter, r *http.Request) error {
	var request struct {
		ExpertID        int64  `json:"expertId"`
		GeneralArea     string `json:"generalArea"`
		SpecializedArea string `json:"specializedArea"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Validate required fields
	if request.GeneralArea == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "General area is required"})
	}
	
	// Combine general and specialized areas for input
	input := request.GeneralArea
	if request.SpecializedArea != "" {
		input += " - " + request.SpecializedArea
	}
	
	// Generate the ISCED suggestion via the Storage interface
	result, err := s.store.SuggestISCED(request.ExpertID, input)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to suggest ISCED classification: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, result)
}

func (s *APIServer) handleExtractSkills(w http.ResponseWriter, r *http.Request) error {
	// Parse the multipart form data, 10 MB max
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Failed to parse form"})
	}
	
	// Get expert ID
	expertIDStr := r.FormValue("expertId")
	if expertIDStr == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Expert ID is required"})
	}
	
	expertID, err := strconv.ParseInt(expertIDStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}
	
	// Get the file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "No file provided"})
	}
	defer file.Close()
	
	// Upload and store the document
	doc, err := s.documentService.CreateDocument(expertID, file, header, "skills_extraction")
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to upload document: %v", err)})
	}
	
	// Extract text from the document
	text, err := s.documentService.ExtractTextFromDocument(doc.FilePath)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to extract text: %v", err)})
	}
	
	// Extract skills using the Storage interface
	result, err := s.store.ExtractSkills(expertID, text)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to extract skills: %v", err)})
	}
	
	// Set document ID in result
	result.DocumentID = doc.ID
	
	return WriteJson(w, http.StatusOK, result)
}

// handleSuggestPanel suggests an expert panel for a project
func (s *APIServer) handleSuggestPanel(w http.ResponseWriter, r *http.Request) error {
	var request struct {
		ProjectName string `json:"projectName"`
		ISCEDFieldID int64 `json:"iscedFieldId"`
		NumExperts int `json:"numExperts"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Validate request
	if strings.TrimSpace(request.ProjectName) == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Project name is required"})
	}
	
	if request.NumExperts <= 0 {
		request.NumExperts = 3 // Default to 3 experts if not specified
	}
	
	// Call storage interface to suggest expert panel
	experts, err := s.store.SuggestExpertPanel(request.ProjectName, request.NumExperts)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to suggest expert panel: %v", err)})
	}
	
	// Prepare response
	result := struct {
		Experts []Expert `json:"experts"`
		Count   int      `json:"count"`
	}{
		Experts: experts,
		Count:   len(experts),
	}
	
	return WriteJson(w, http.StatusOK, result)
}

// Statistics handlers

func (s *APIServer) handleGetStatistics(w http.ResponseWriter, r *http.Request) error {
	// Get all statistics
	stats, err := s.store.GetStatistics()
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve statistics: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, stats)
}

func (s *APIServer) handleGetNationalityStats(w http.ResponseWriter, r *http.Request) error {
	// Get nationality statistics
	bahrainiCount, nonBahrainiCount, err := s.store.GetExpertsByNationality()
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve nationality statistics: %v", err)})
	}
	
	// Calculate total
	total := bahrainiCount + nonBahrainiCount
	
	// Calculate percentages
	var bahrainiPercentage, nonBahrainiPercentage float64
	if total > 0 {
		bahrainiPercentage = float64(bahrainiCount) / float64(total) * 100
		nonBahrainiPercentage = float64(nonBahrainiCount) / float64(total) * 100
	}
	
	// Create stats array in the format expected by the frontend
	stats := []AreaStat{
		{Name: "Bahraini", Count: bahrainiCount, Percentage: bahrainiPercentage},
		{Name: "Non-Bahraini", Count: nonBahrainiCount, Percentage: nonBahrainiPercentage},
	}
	
	// Prepare response
	result := map[string]interface{}{
		"total": total,
		"stats": stats,
	}
	
	return WriteJson(w, http.StatusOK, result)
}

func (s *APIServer) handleGetISCEDStats(w http.ResponseWriter, r *http.Request) error {
	// Get ISCED statistics
	stats, err := s.store.GetExpertsByISCEDField()
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve ISCED statistics: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, stats)
}

func (s *APIServer) handleGetEngagementStats(w http.ResponseWriter, r *http.Request) error {
	// Get engagement statistics
	stats, err := s.store.GetEngagementStatistics()
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve engagement statistics: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, stats)
}

func (s *APIServer) handleGetGrowthStats(w http.ResponseWriter, r *http.Request) error {
	// Parse months parameter
	months := 12 // Default to 12 months
	
	monthsParam := r.URL.Query().Get("months")
	if monthsParam != "" {
		parsedMonths, err := strconv.Atoi(monthsParam)
		if err == nil && parsedMonths > 0 {
			months = parsedMonths
		}
	}
	
	// Get growth statistics
	stats, err := s.store.GetExpertGrowthByMonth(months)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve growth statistics: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, stats)
}

// Expert Request Handler Functions

func (s *APIServer) handleCreateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	var request ExpertRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Validate required fields
	if strings.TrimSpace(request.Name) == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Name is required"})
	}
	if strings.TrimSpace(request.Institution) == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Institution is required"})
	}
	if strings.TrimSpace(request.Designation) == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Designation is required"})
	}
	if strings.TrimSpace(request.Email) == "" && strings.TrimSpace(request.Phone) == "" {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "At least one contact method (email or phone) is required"})
	}
	
	// Set default values
	request.Status = "pending"
	request.CreatedAt = time.Now()
	
	// Save the request to the database
	id, err := s.store.CreateExpertRequest(&request)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to create expert request: %v", err)})
	}
	
	// Set the ID in the response
	request.ID = id
	
	return WriteJson(w, http.StatusCreated, request)
}

func (s *APIServer) handleGetExpertRequests(w http.ResponseWriter, r *http.Request) error {
	// Get query parameters for filtering
	status := r.URL.Query().Get("status")
	
	// Parse pagination parameters
	limit := 100
	offset := 0
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	// Build filters
	filters := make(map[string]interface{})
	if status != "" {
		filters["status"] = status
	}
	
	// Get requests with filters
	requests, err := s.store.ListExpertRequests(filters, limit, offset)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve expert requests: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, requests)
}

func (s *APIServer) handleGetExpertRequest(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request ID"})
	}
	
	// Get the expert request
	request, err := s.store.GetExpertRequest(id)
	if err != nil {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Expert request not found: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, request)
}

func (s *APIServer) handleUpdateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request ID"})
	}
	
	// Get existing request
	existingRequest, err := s.store.GetExpertRequest(id)
	if err != nil {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Expert request not found: %v", err)})
	}
	
	// Parse update data
	var updateRequest ExpertRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Ensure ID matches
	updateRequest.ID = id
	
	// Handle status changes - if status is changing to "approved", create an expert record
	if existingRequest.Status != "approved" && updateRequest.Status == "approved" {
		// Create a new expert from the request data
		expert := &Expert{
			ExpertID:        updateRequest.ExpertID,
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
			CreatedAt:       time.Now(),
		}
		
		// Create the expert record
		expertID, err := s.store.CreateExpert(expert)
		if err != nil {
			return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to create expert from request: %v", err)})
		}
		
		// Set the reviewed timestamp and reviewer info
		updateRequest.ReviewedAt = time.Now()
		// Note: In a real app, we'd get the reviewer ID from the authenticated user
		
		// Update the expert request with the expert ID and other review info
		updateRequest.ExpertID = fmt.Sprintf("EXP-%d", expertID)
	}
	
	// Update the expert request
	if err := s.store.UpdateExpertRequest(&updateRequest); err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to update expert request: %v", err)})
	}
	
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Expert request updated successfully",
	})
}

// Helper functions

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := GetLogger()
		
		if err := f(w, r); err != nil {
			// Get request details for logging
			path := r.URL.Path
			method := r.Method
			
			// Handle specific errors
			switch {
			case err == ErrUnauthorized:
				logger.Warn("Unauthorized access: %s %s from %s", method, path, r.RemoteAddr)
				WriteJson(w, http.StatusUnauthorized, ApiError{Error: "Unauthorized: Authentication required"})
				
			case err == ErrForbidden:
				logger.Warn("Forbidden access: %s %s from %s", method, path, r.RemoteAddr)
				WriteJson(w, http.StatusForbidden, ApiError{Error: "Forbidden: Insufficient permissions"})
				
			case err == ErrInvalidCredentials:
				logger.Warn("Invalid credentials attempt: %s %s from %s", method, path, r.RemoteAddr)
				WriteJson(w, http.StatusUnauthorized, ApiError{Error: "Invalid email or password"})
				
			case err == ErrNotFound:
				logger.Info("Resource not found: %s %s", method, path)
				WriteJson(w, http.StatusNotFound, ApiError{Error: "Resource not found"})
				
			default:
				// Generic error handler
				logger.Error("Handler error: %s %s - %v", method, path, err)
				WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
			}
		}
	}
}