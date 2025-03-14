// Package main provides the backend functionality for the ExpertDB application
package main

// NOTE: AI integration has been removed as per requirements.
// Previous AI routes and services are no longer available.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Pagination and filtering constants
const (
	// DefaultLimit is the default number of items per page for pagination
	DefaultLimit = 10
	
	// DefaultOffset is the default starting position for pagination
	DefaultOffset = 0
	
	// MaxUploadSize is the maximum file size for document uploads (10 MB)
	MaxUploadSize = 10 << 20
)

// APIServer represents the HTTP API server for the ExpertDB application
// It handles routing, request processing, and coordinating between different services
type APIServer struct {
	listenAddr      string          // HTTP listening address (e.g., ":8080")
	store           Storage         // Database storage interface
	documentService *DocumentService // Service for document upload/management
	config          *Configuration  // Application configuration
}

// NewAPIServer creates a new API server instance with the provided dependencies
//
// This function initializes a new API server with a database connection,
// document service, and configuration settings.
//
// Inputs:
//   - listenAddr (string): The address to listen on (e.g., ":8080")
//   - store (Storage): The database storage implementation
//   - config (*Configuration): Application configuration parameters
//
// Returns:
//   - *APIServer: The initialized API server
//   - error: Any error that occurs during initialization
func NewAPIServer(listenAddr string, store Storage, config *Configuration) (*APIServer, error) {
	// Initialize the document service for handling file uploads
	documentService, err := NewDocumentService(store, config.UploadPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize document service: %w", err)
	}
	
	// Create and return the API server instance
	return &APIServer{
		listenAddr:      listenAddr,
		store:           store,
		documentService: documentService,
		config:          config,
	}, nil
}

// Run starts the HTTP server and listens for API requests
//
// This method initializes all API routes, applies middleware,
// and starts the server on the configured listen address.
//
// Returns:
//   - error: Any error that occurs during server startup or operation
func (s *APIServer) Run() error {
	logger := GetLogger()
	logger.Info("Setting up API routes...")
	
	// Create a new HTTP router
	mux := http.NewServeMux()

	// Step 1: Register all API routes
	// Routes are organized by resource type for clarity
	
	// Step 1.1: Expert management routes
	logger.Debug("Registering expert management routes")
	mux.HandleFunc("GET /api/experts", makeHTTPHandleFunc(s.handleGetExperts))
	mux.HandleFunc("POST /api/experts", makeHTTPHandleFunc(requireAdmin(s.handleCreateExpert)))
	mux.HandleFunc("GET /api/experts/{id}", makeHTTPHandleFunc(s.handleGetExpert))
	mux.HandleFunc("PUT /api/experts/{id}", makeHTTPHandleFunc(requireAdmin(s.handleUpdateExpert)))
	mux.HandleFunc("DELETE /api/experts/{id}", makeHTTPHandleFunc(requireAdmin(s.handleDeleteExpert)))
	
	// Step 1.2: Expert request routes (for requesting addition of new experts)
	logger.Debug("Registering expert request routes")
	mux.HandleFunc("POST /api/expert-requests", makeHTTPHandleFunc(requireAuth(s.handleCreateExpertRequest)))
	mux.HandleFunc("GET /api/expert-requests", makeHTTPHandleFunc(requireAuth(s.handleGetExpertRequests)))
	mux.HandleFunc("GET /api/expert-requests/{id}", makeHTTPHandleFunc(requireAuth(s.handleGetExpertRequest)))
	mux.HandleFunc("PUT /api/expert-requests/{id}", makeHTTPHandleFunc(requireAdmin(s.handleUpdateExpertRequest)))

	// Step 1.3: ISCED classification reference data routes
	logger.Debug("Registering ISCED classification routes")
	mux.HandleFunc("GET /api/isced/levels", makeHTTPHandleFunc(s.handleGetISCEDLevels))
	mux.HandleFunc("GET /api/isced/fields", makeHTTPHandleFunc(s.handleGetISCEDFields))
	
	// Step 1.4: Document management routes
	logger.Debug("Registering document management routes")
	mux.HandleFunc("POST /api/documents", makeHTTPHandleFunc(requireAuth(s.handleUploadDocument)))
	mux.HandleFunc("GET /api/documents/{id}", makeHTTPHandleFunc(requireAuth(s.handleGetDocument)))
	mux.HandleFunc("DELETE /api/documents/{id}", makeHTTPHandleFunc(requireAdmin(s.handleDeleteDocument)))
	mux.HandleFunc("GET /api/experts/{id}/documents", makeHTTPHandleFunc(requireAuth(s.handleGetExpertDocuments)))
	
	// Step 1.5: Expert engagement tracking routes
	logger.Debug("Registering expert engagement routes")
	mux.HandleFunc("POST /api/engagements", makeHTTPHandleFunc(requireAdmin(s.handleCreateEngagement)))
	mux.HandleFunc("GET /api/engagements/{id}", makeHTTPHandleFunc(requireAuth(s.handleGetEngagement)))
	mux.HandleFunc("PUT /api/engagements/{id}", makeHTTPHandleFunc(requireAdmin(s.handleUpdateEngagement)))
	mux.HandleFunc("DELETE /api/engagements/{id}", makeHTTPHandleFunc(requireAdmin(s.handleDeleteEngagement)))
	mux.HandleFunc("GET /api/experts/{id}/engagements", makeHTTPHandleFunc(requireAuth(s.handleGetExpertEngagements)))
	
	// Step 1.6: Statistics and analytics routes
	logger.Debug("Registering statistics routes")
	mux.HandleFunc("GET /api/statistics", makeHTTPHandleFunc(requireAuth(s.handleGetStatistics)))
	mux.HandleFunc("GET /api/statistics/nationality", makeHTTPHandleFunc(requireAuth(s.handleGetNationalityStats)))
	mux.HandleFunc("GET /api/statistics/isced", makeHTTPHandleFunc(requireAuth(s.handleGetISCEDStats)))
	mux.HandleFunc("GET /api/statistics/engagements", makeHTTPHandleFunc(requireAuth(s.handleGetEngagementStats)))
	mux.HandleFunc("GET /api/statistics/growth", makeHTTPHandleFunc(requireAuth(s.handleGetGrowthStats)))
	
	// Step 1.7: User management and authentication routes
	logger.Debug("Registering user management routes")
	mux.HandleFunc("POST /api/auth/login", makeHTTPHandleFunc(s.handleLogin))
	mux.HandleFunc("POST /api/users", makeHTTPHandleFunc(requireAdmin(s.handleCreateUser)))
	mux.HandleFunc("GET /api/users", makeHTTPHandleFunc(requireAdmin(s.handleGetUsers)))
	mux.HandleFunc("GET /api/users/{id}", makeHTTPHandleFunc(requireAuth(s.handleGetUser)))
	mux.HandleFunc("PUT /api/users/{id}", makeHTTPHandleFunc(requireAdmin(s.handleUpdateUser)))
	mux.HandleFunc("DELETE /api/users/{id}", makeHTTPHandleFunc(requireAdmin(s.handleDeleteUser)))
	
	// Step 2: Set up middleware stack
	
	// Step 2.1: CORS middleware for cross-origin requests
	logger.Debug("Setting up CORS middleware")
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set standard CORS headers
			w.Header().Set("Access-Control-Allow-Origin", s.config.CORSAllowOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			// Handle OPTIONS preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			// Proceed with the regular request
			next.ServeHTTP(w, r)
		})
	}
	
	// Get logger for request logging middleware
	apiLogger := GetLogger()
	
	// Step 2.2: Apply middleware chain: request logging -> CORS -> router
	// Order matters: outermost middleware is applied first
	logger.Debug("Applying middleware chain")
	handler := apiLogger.RequestLoggerMiddleware(corsMiddleware(mux))

	// Step 3: Start the HTTP server
	logger.Info("API server listening on %s", s.listenAddr)
	
	// ListenAndServe blocks until the server is shut down or encounters an error
	return http.ListenAndServe(s.listenAddr, handler)
}

// handleGetExperts handles GET /api/experts requests.
//
// This handler retrieves a paginated list of experts, applying optional
// filtering and sorting based on query parameters.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing query parameters
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Parse query parameters for filtering criteria
//   2. Extract pagination and sorting parameters
//   3. Query database with filters and pagination
//   4. Return formatted JSON response with pagination headers
func (s *APIServer) handleGetExperts(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/experts request")

	// Step 1: Parse query parameters for filtering
	queryParams := r.URL.Query()
	filters := make(map[string]interface{})

	// Step 1.1: Process name and area filters
	if name := queryParams.Get("name"); name != "" {
		filters["name"] = name
	}
	if area := queryParams.Get("area"); area != "" {
		filters["area"] = area
	}

	// Step 1.2: Process boolean availability filter
	if available := queryParams.Get("is_available"); available != "" {
		if available == "true" {
			filters["is_available"] = true
		} else if available == "false" {
			filters["is_available"] = false
		}
	}

	// Step 1.3: Process role filter
	if role := queryParams.Get("role"); role != "" {
		filters["role"] = role
	}

	// Step 1.4: Process ISCED classification filters
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
	
	// Step 1.5: Process rating filter
	if minRating := queryParams.Get("min_rating"); minRating != "" {
		if rating, err := strconv.ParseFloat(minRating, 64); err == nil {
			filters["min_rating"] = rating
		}
	}

	// Step 2: Parse pagination parameters using the helper function
	limit, offset := parsePaginationParams(r, DefaultLimit)
	
	// Step 3: Process sorting parameters
	sortBy := "name" // Default sort field
	sortOrder := "asc" // Default sort order
	
	// Step 3.1: Validate and set sort field
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
	
	// Step 3.2: Validate and set sort order
	if orderParam := queryParams.Get("sort_order"); orderParam != "" {
		if orderParam == "desc" {
			sortOrder = "desc"
		}
	}
	
	// Step 3.3: Add sorting to filters
	filters["sort_by"] = sortBy
	filters["sort_order"] = sortOrder

	// Step 4: Get total count for pagination
	
	// Step 4.1: Create a copy of filters without pagination
	countFilters := make(map[string]interface{})
	for k, v := range filters {
		if k != "sort_by" && k != "sort_order" {
			countFilters[k] = v
		}
	}
	
	// Step 4.2: Get total count
	totalCount, err := s.store.CountExperts(countFilters)
	if err != nil {
		logger.Error("Failed to count experts: %v", err)
		return fmt.Errorf("failed to count experts: %w", err)
	}
	
	// Step 5: Retrieve filtered experts with pagination
	logger.Debug("Retrieving experts with filters: %v, limit: %d, offset: %d", filters, limit, offset)
	experts, err := s.store.ListExperts(filters, limit, offset)
	if err != nil {
		logger.Error("Failed to list experts: %v", err)
		return fmt.Errorf("failed to retrieve experts: %w", err)
	}
	
	// Step 6: Set response headers and return results
	
	// Step 6.1: Set total count header for pagination
	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", totalCount))
	
	// Step 6.2: Return JSON response
	logger.Debug("Returning %d experts", len(experts))
	return WriteJson(w, http.StatusOK, experts)
}

// handleGetExpert handles GET /api/experts/{id} requests.
//
// This handler retrieves a single expert by ID. If the expert is not found,
// it returns an empty Expert object with a 200 status code rather than a 404 error.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert ID in the path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the expert ID from the request path
//   2. Retrieve the expert from the database
//   3. Return the expert data in JSON format
func (s *APIServer) handleGetExpert(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid expert ID provided: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}

	// Step 2: Retrieve expert from database
	logger.Debug("Retrieving expert with ID: %d", id)
	expert, err := s.store.GetExpert(id)
	if err != nil {
		// Instead of returning a 404 error, return an empty Expert object
		// This is a business requirement to handle not-found cases gracefully
		logger.Warn("Expert not found for ID: %d - %v", id, err)
		return WriteJson(w, http.StatusOK, &Expert{})
	}

	// Step 3: Return expert data
	logger.Debug("Successfully retrieved expert: %s (ID: %d)", expert.Name, expert.ID)
	return WriteJson(w, http.StatusOK, expert)
}

// handleCreateExpert handles POST /api/experts requests.
//
// This handler creates a new expert in the database. It requires admin privileges,
// which is enforced by the requireAdmin middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert data in the body
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Parse and validate the request body
//   2. Create the expert record in the database
//   3. Return success response with the new expert's ID
func (s *APIServer) handleCreateExpert(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing POST /api/experts request")
	
	// Step 1: Parse request body
	createReq := new(CreateExpertRequest)
	if err := json.NewDecoder(r.Body).Decode(createReq); err != nil {
		logger.Warn("Failed to parse expert creation request: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}

	// Step 2: Validate request data
	if err := ValidateCreateExpertRequest(createReq); err != nil {
		logger.Warn("Expert creation validation failed: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	// Step 3: Create expert record in database
	expert := NewExpert(*createReq)
	logger.Debug("Creating expert: %s, Institution: %s", expert.Name, expert.Institution)
	id, err := s.store.CreateExpert(expert)
	if err != nil {
		logger.Error("Failed to create expert in database: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to create expert: %v", err)})
	}

	// Step 4: Return success response
	logger.Info("Expert created successfully with ID: %d", id)
	return WriteJson(w, http.StatusCreated, CreateExpertResponse{
		ID:      id,
		Success: true,
		Message: "Expert created successfully",
	})
}

// handleUpdateExpert handles PUT /api/experts/{id} requests.
//
// This handler updates an existing expert by ID. It requires admin privileges,
// which is enforced by the requireAdmin middleware. If the expert doesn't exist,
// it creates a new record with the specified ID.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert ID in path and update data in body
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the expert ID from the request path
//   2. Retrieve the existing expert from the database (if exists)
//   3. Parse and validate the update data from the request body
//   4. Merge update data with existing data (maintaining existing values for unspecified fields)
//   5. Update the expert record in the database
//   6. Return success response
func (s *APIServer) handleUpdateExpert(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid expert ID provided for update: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}

	// Step 2: Check if expert exists
	logger.Debug("Checking if expert exists with ID: %d", id)
	existingExpert, err := s.store.GetExpert(id)
	if err != nil {
		// Create default expert with this ID if not exists
		// This approach allows updating non-existent experts
		logger.Warn("Expert not found for update ID: %d - creating empty record", id)
		existingExpert = &Expert{ID: id}
	}

	// Step 3: Parse update request
	var updateExpert Expert
	if err := json.NewDecoder(r.Body).Decode(&updateExpert); err != nil {
		logger.Warn("Failed to parse expert update request: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}

	// Step 4: Ensure ID matches path parameter
	updateExpert.ID = id

	// Step 5: Merge with existing expert data - use existing data for empty fields
	if existingExpert != nil {
		// Only replace fields that are set in the update
		if updateExpert.ExpertID == "" {
			updateExpert.ExpertID = existingExpert.ExpertID
		}
		// NOTE: Additional fields should be added here to preserve existing data
		// for any fields not included in the update request
		if updateExpert.Name == "" {
			updateExpert.Name = existingExpert.Name
		}
		if updateExpert.Institution == "" {
			updateExpert.Institution = existingExpert.Institution
		}
		if updateExpert.Designation == "" {
			updateExpert.Designation = existingExpert.Designation
		}
		// Preserve created date
		if updateExpert.CreatedAt.IsZero() {
			updateExpert.CreatedAt = existingExpert.CreatedAt
		}
	}

	// Step 6: Update expert in database
	logger.Debug("Updating expert ID: %d, Name: %s", id, updateExpert.Name)
	if err := s.store.UpdateExpert(&updateExpert); err != nil {
		logger.Error("Failed to update expert in database: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to update expert: %v", err)})
	}

	// Step 7: Return success response
	logger.Info("Expert updated successfully: ID: %d", id)
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Expert updated successfully",
	})
}

// handleDeleteExpert handles DELETE /api/experts/{id} requests.
//
// This handler deletes an expert by ID. It requires admin privileges,
// which is enforced by the requireAdmin middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert ID in path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the expert ID from the request path
//   2. Delete the expert and related data from the database
//   3. Return success response
func (s *APIServer) handleDeleteExpert(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid expert ID provided for deletion: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}

	// Step 2: Delete expert and related data from database
	logger.Debug("Deleting expert with ID: %d", id)
	if err := s.store.DeleteExpert(id); err != nil {
		logger.Error("Failed to delete expert: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to delete expert: %v", err)})
	}

	// Step 3: Return success response
	logger.Info("Expert deleted successfully: ID: %d", id)
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Expert deleted successfully",
	})
}

// handleGetISCEDLevels handles GET /api/isced/levels requests.
//
// This handler retrieves all ISCED education levels from the database.
// ISCED (International Standard Classification of Education) levels are used
// for categorizing and comparing education qualifications internationally.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Query the database for all ISCED levels
//   2. Return the levels as a JSON response
func (s *APIServer) handleGetISCEDLevels(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/isced/levels request")
	
	// Step 1: Query ISCED levels from database
	levels, err := s.store.GetISCEDLevels()
	if err != nil {
		logger.Error("Failed to fetch ISCED levels: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to fetch ISCED levels: %v", err)})
	}

	// Step 2: Return levels as JSON response
	logger.Debug("Returning %d ISCED levels", len(levels))
	return WriteJson(w, http.StatusOK, levels)
}

// handleGetISCEDFields handles GET /api/isced/fields requests.
//
// This handler retrieves all ISCED fields of education from the database,
// transforms them to match frontend expectations, and returns them as a JSON response.
// ISCED (International Standard Classification of Education) fields categorize
// areas of study and research.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Query the database for all ISCED fields
//   2. Transform the data to match frontend expectations
//   3. Return the simplified fields as a JSON response
func (s *APIServer) handleGetISCEDFields(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/isced/fields request")
	
	// Step 1: Query ISCED fields from database
	fields, err := s.store.GetISCEDFields()
	if err != nil {
		// Return empty array instead of 500 error for better UX
		logger.Error("Failed to fetch ISCED fields: %v", err)
		return WriteJson(w, http.StatusOK, []map[string]interface{}{})
	}

	// Step 2: Transform the data to match frontend expectations
	logger.Debug("Transforming %d ISCED fields to simplified format", len(fields))
	simplifiedFields := make([]map[string]interface{}, 0, len(fields))
	
	// Process each field and create a simplified representation
	for _, field := range fields {
		// Skip any entry with empty ID or name
		if field.ID != 0 && field.BroadName != "" {
			simplifiedFields = append(simplifiedFields, map[string]interface{}{
				"id":   fmt.Sprintf("%d", field.ID), // Convert ID to string for frontend
				"name": field.BroadName,             // Use BroadName as the name field
			})
		}
	}

	// Step 3: Return simplified fields as JSON response
	logger.Debug("Returning %d simplified ISCED fields", len(simplifiedFields))
	return WriteJson(w, http.StatusOK, simplifiedFields)
}

// Document handlers

// handleUploadDocument handles POST /api/documents requests.
//
// This handler uploads and stores a document file (such as a CV, certification,
// or other expert-related document) associated with a specific expert.
// It requires authentication via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the multipart form data
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Parse the multipart form data from the request
//   2. Extract and validate the expert ID and document type
//   3. Get the uploaded file from the request
//   4. Store the document using the document service
//   5. Return the document information as a JSON response
func (s *APIServer) handleUploadDocument(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing POST /api/documents request")
	
	// Step 1: Parse the multipart form data
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		logger.Warn("Failed to parse multipart form: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Failed to parse form - file may be too large"})
	}
	
	// Step 2: Extract and validate expert ID
	expertIDStr := r.FormValue("expertId")
	if expertIDStr == "" {
		logger.Warn("Missing expert ID in document upload request")
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Expert ID is required"})
	}
	
	expertID, err := strconv.ParseInt(expertIDStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid expert ID in document upload request: %s", expertIDStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}
	
	// Step 3: Get document type or use default
	docType := r.FormValue("documentType")
	if docType == "" {
		logger.Debug("No document type specified, using default type: cv")
		docType = "cv" // Default type
	}
	
	// Step 4: Get the file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Warn("No file provided in document upload request: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "No file provided"})
	}
	defer file.Close()
	
	// Step 5: Upload and store the document
	logger.Debug("Uploading document for expert ID: %d, type: %s, filename: %s",
		expertID, docType, header.Filename)
	doc, err := s.documentService.CreateDocument(expertID, file, header, docType)
	if err != nil {
		logger.Error("Failed to upload document: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to upload document: %v", err)})
	}
	
	// Step 6: Return document information
	logger.Info("Document uploaded successfully: ID: %d, Type: %s, Expert: %d", doc.ID, doc.Type, doc.ExpertID)
	return WriteJson(w, http.StatusCreated, doc)
}

// handleGetDocument handles GET /api/documents/{id} requests.
//
// This handler retrieves a document by its ID. It requires authentication
// via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the document ID in the path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the document ID from the request path
//   2. Retrieve the document from the document service
//   3. Return the document information as a JSON response
func (s *APIServer) handleGetDocument(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate document ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid document ID provided: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid document ID"})
	}
	
	// Step 2: Retrieve document from document service
	logger.Debug("Retrieving document with ID: %d", id)
	doc, err := s.documentService.GetDocument(id)
	if err != nil {
		logger.Warn("Document not found with ID: %d - %v", id, err)
		return WriteJson(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Document not found: %v", err)})
	}
	
	// Step 3: Return document information
	logger.Debug("Returning document: ID: %d, Type: %s, Expert: %d", doc.ID, doc.Type, doc.ExpertID)
	return WriteJson(w, http.StatusOK, doc)
}

// handleDeleteDocument handles DELETE /api/documents/{id} requests.
//
// This handler deletes a document by its ID. It requires admin privileges,
// which is enforced by the requireAdmin middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the document ID in the path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the document ID from the request path
//   2. Delete the document using the document service
//   3. Return a success response
func (s *APIServer) handleDeleteDocument(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate document ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid document ID provided for deletion: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid document ID"})
	}
	
	// Step 2: Delete the document
	logger.Debug("Deleting document with ID: %d", id)
	if err := s.documentService.DeleteDocument(id); err != nil {
		logger.Error("Failed to delete document: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to delete document: %v", err)})
	}
	
	// Step 3: Return success response
	logger.Info("Document deleted successfully: ID: %d", id)
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Document deleted successfully",
	})
}

// handleGetExpertDocuments handles GET /api/experts/{id}/documents requests.
//
// This handler retrieves all documents associated with a specific expert.
// It requires authentication via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert ID in the path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the expert ID from the request path
//   2. Retrieve the expert's documents from the document service
//   3. Return the documents as a JSON response
func (s *APIServer) handleGetExpertDocuments(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid expert ID provided for document retrieval: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}
	
	// Step 2: Retrieve the expert's documents
	logger.Debug("Retrieving documents for expert with ID: %d", id)
	docs, err := s.documentService.GetDocumentsByExpertID(id)
	if err != nil {
		logger.Error("Failed to retrieve documents for expert %d: %v", id, err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve documents: %v", err)})
	}
	
	// Step 3: Return documents
	logger.Debug("Returning %d documents for expert ID: %d", len(docs), id)
	return WriteJson(w, http.StatusOK, docs)
}

// Engagement handlers

// handleCreateEngagement handles POST /api/engagements requests.
//
// This handler creates a new expert engagement record in the database.
// It requires admin privileges, which is enforced by the requireAdmin middleware.
// Engagements represent activities involving experts, such as panel participation,
// training delivery, or consulting engagements.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the engagement data in the body
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Parse and validate the request body
//   2. Set default values and validate required fields
//   3. Create the engagement record in the database
//   4. Return the created engagement with its assigned ID
func (s *APIServer) handleCreateEngagement(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing POST /api/engagements request")
	
	// Step 1: Parse request body
	var engagement Engagement
	if err := json.NewDecoder(r.Body).Decode(&engagement); err != nil {
		logger.Warn("Failed to parse engagement creation request: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Step 2: Set default values and validate
	
	// Step 2.1: Set creation time
	engagement.CreatedAt = time.Now()
	
	// Step 2.2: Validate required fields
	if engagement.ExpertID == 0 {
		logger.Warn("Missing expert ID in engagement creation request")
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Expert ID is required"})
	}
	if engagement.EngagementType == "" {
		logger.Warn("Missing engagement type in creation request")
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Engagement type is required"})
	}
	if engagement.StartDate.IsZero() {
		logger.Warn("Missing start date in engagement creation request")
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Start date is required"})
	}
	
	// Step 2.3: Set default status if not provided
	if engagement.Status == "" {
		logger.Debug("No status specified, using default status: pending")
		engagement.Status = "pending" // Default status
	}
	
	// Step 3: Create the engagement in database
	logger.Debug("Creating engagement for expert ID: %d, type: %s", 
		engagement.ExpertID, engagement.EngagementType)
	id, err := s.store.CreateEngagement(&engagement)
	if err != nil {
		logger.Error("Failed to create engagement in database: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to create engagement: %v", err)})
	}
	
	// Step 4: Set the ID in the response and return
	engagement.ID = id
	logger.Info("Engagement created successfully: ID: %d, Type: %s, Expert: %d", 
		id, engagement.EngagementType, engagement.ExpertID)
	return WriteJson(w, http.StatusCreated, engagement)
}

// handleGetEngagement handles GET /api/engagements/{id} requests.
//
// This handler retrieves a single engagement by ID. It requires authentication
// via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the engagement ID in the path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the engagement ID from the request path
//   2. Retrieve the engagement from the database
//   3. Return the engagement data in JSON format
func (s *APIServer) handleGetEngagement(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid engagement ID provided: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid engagement ID"})
	}
	
	// Step 2: Retrieve engagement from database
	logger.Debug("Retrieving engagement with ID: %d", id)
	engagement, err := s.store.GetEngagement(id)
	if err != nil {
		logger.Warn("Engagement not found with ID: %d - %v", id, err)
		return WriteJson(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Engagement not found: %v", err)})
	}
	
	// Step 3: Return engagement data
	logger.Debug("Successfully retrieved engagement: ID: %d, Type: %s", engagement.ID, engagement.EngagementType)
	return WriteJson(w, http.StatusOK, engagement)
}

// handleUpdateEngagement handles PUT /api/engagements/{id} requests.
//
// This handler updates an existing engagement by ID. It requires admin privileges,
// which is enforced by the requireAdmin middleware. The update preserves any
// fields that are not explicitly provided in the update request.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the engagement ID in path and update data in body
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the engagement ID from the request path
//   2. Retrieve the existing engagement from the database
//   3. Parse and validate the update data from the request body
//   4. Merge update data with existing data (maintaining existing values for unspecified fields)
//   5. Update the engagement record in the database
//   6. Return success response
func (s *APIServer) handleUpdateEngagement(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid engagement ID provided for update: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid engagement ID"})
	}
	
	// Step 2: Retrieve existing engagement from database
	logger.Debug("Checking if engagement exists with ID: %d", id)
	existing, err := s.store.GetEngagement(id)
	if err != nil {
		logger.Warn("Engagement not found for update ID: %d - %v", id, err)
		return WriteJson(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Engagement not found: %v", err)})
	}
	
	// Step 3: Parse update request
	var updateEngagement Engagement
	if err := json.NewDecoder(r.Body).Decode(&updateEngagement); err != nil {
		logger.Warn("Failed to parse engagement update request: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Step 4: Ensure ID matches path parameter
	updateEngagement.ID = id
	
	// Step 5: Merge with existing engagement data - use existing data for empty fields
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
	
	// Step 6: Update the engagement in database
	logger.Debug("Updating engagement ID: %d, Type: %s", id, updateEngagement.EngagementType)
	if err := s.store.UpdateEngagement(&updateEngagement); err != nil {
		logger.Error("Failed to update engagement in database: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to update engagement: %v", err)})
	}
	
	// Step 7: Return success response
	logger.Info("Engagement updated successfully: ID: %d", id)
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Engagement updated successfully",
	})
}

// handleDeleteEngagement handles DELETE /api/engagements/{id} requests.
//
// This handler deletes an engagement by ID. It requires admin privileges,
// which is enforced by the requireAdmin middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the engagement ID in the path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the engagement ID from the request path
//   2. Delete the engagement from the database
//   3. Return success response
func (s *APIServer) handleDeleteEngagement(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate engagement ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid engagement ID provided for deletion: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid engagement ID"})
	}
	
	// Step 2: Delete the engagement from database
	logger.Debug("Deleting engagement with ID: %d", id)
	if err := s.store.DeleteEngagement(id); err != nil {
		logger.Error("Failed to delete engagement: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to delete engagement: %v", err)})
	}
	
	// Step 3: Return success response
	logger.Info("Engagement deleted successfully: ID: %d", id)
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Engagement deleted successfully",
	})
}

// handleGetExpertEngagements handles GET /api/experts/{id}/engagements requests.
//
// This handler retrieves all engagements associated with a specific expert.
// It requires authentication via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert ID in the path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the expert ID from the request path
//   2. Retrieve the expert's engagements from the database
//   3. Return the engagements as a JSON response
func (s *APIServer) handleGetExpertEngagements(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid expert ID provided for engagement retrieval: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid expert ID"})
	}
	
	// Step 2: Retrieve the expert's engagements from database
	logger.Debug("Retrieving engagements for expert with ID: %d", id)
	engagements, err := s.store.GetEngagementsByExpertID(id)
	if err != nil {
		logger.Error("Failed to retrieve engagements for expert %d: %v", id, err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve engagements: %v", err)})
	}
	
	// Step 3: Return engagements
	logger.Debug("Returning %d engagements for expert ID: %d", len(engagements), id)
	return WriteJson(w, http.StatusOK, engagements)
}


// Statistics handlers

// handleGetStatistics handles GET /api/statistics requests.
//
// This handler retrieves the overall statistics for the expert database,
// providing a summary of key metrics. It requires authentication via
// the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Retrieve all statistics from the database
//   2. Return the statistics as a JSON response
func (s *APIServer) handleGetStatistics(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/statistics request")
	
	// Step 1: Retrieve all statistics from database
	logger.Debug("Retrieving overall system statistics")
	stats, err := s.store.GetStatistics()
	if err != nil {
		logger.Error("Failed to retrieve statistics: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve statistics: %v", err)})
	}
	
	// Step 2: Return statistics as JSON response
	logger.Debug("Successfully retrieved system statistics")
	return WriteJson(w, http.StatusOK, stats)
}

// handleGetNationalityStats handles GET /api/statistics/nationality requests.
//
// This handler retrieves statistics on expert nationality distribution (Bahraini vs Non-Bahraini).
// It requires authentication via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Retrieve nationality statistics from the database
//   2. Calculate total count and percentages
//   3. Format the data in the structure expected by the frontend
//   4. Return the statistics as a JSON response
func (s *APIServer) handleGetNationalityStats(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/statistics/nationality request")
	
	// Step 1: Retrieve nationality statistics from database
	logger.Debug("Retrieving nationality distribution statistics")
	bahrainiCount, nonBahrainiCount, err := s.store.GetExpertsByNationality()
	if err != nil {
		logger.Error("Failed to retrieve nationality statistics: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve nationality statistics: %v", err)})
	}
	
	// Step 2: Calculate total and percentages
	total := bahrainiCount + nonBahrainiCount
	logger.Debug("Total experts: %d (Bahraini: %d, Non-Bahraini: %d)", total, bahrainiCount, nonBahrainiCount)
	
	// Calculate percentages, avoiding division by zero
	var bahrainiPercentage, nonBahrainiPercentage float64
	if total > 0 {
		bahrainiPercentage = float64(bahrainiCount) / float64(total) * 100
		nonBahrainiPercentage = float64(nonBahrainiCount) / float64(total) * 100
	}
	
	// Step 3: Create stats array in the format expected by the frontend
	stats := []AreaStat{
		{Name: "Bahraini", Count: bahrainiCount, Percentage: bahrainiPercentage},
		{Name: "Non-Bahraini", Count: nonBahrainiCount, Percentage: nonBahrainiPercentage},
	}
	
	// Prepare response with total and detailed stats
	result := map[string]interface{}{
		"total": total,
		"stats": stats,
	}
	
	// Step 4: Return statistics as JSON response
	return WriteJson(w, http.StatusOK, result)
}

// handleGetISCEDStats handles GET /api/statistics/isced requests.
//
// This handler retrieves statistics on expert distribution across ISCED fields (educational areas).
// It requires authentication via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Retrieve ISCED field distribution statistics from the database
//   2. Return the statistics as a JSON response
func (s *APIServer) handleGetISCEDStats(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/statistics/isced request")
	
	// Step 1: Retrieve ISCED field distribution statistics from database
	logger.Debug("Retrieving ISCED field distribution statistics")
	stats, err := s.store.GetExpertsByISCEDField()
	if err != nil {
		logger.Error("Failed to retrieve ISCED statistics: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve ISCED statistics: %v", err)})
	}
	
	// Step 2: Return statistics as JSON response
	logger.Debug("Successfully retrieved ISCED field statistics: %d categories", len(stats))
	return WriteJson(w, http.StatusOK, stats)
}

// handleGetEngagementStats handles GET /api/statistics/engagements requests.
//
// This handler retrieves statistics on expert engagements, showing distribution by type,
// status, or other engagement metrics. It requires authentication via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Retrieve engagement statistics from the database
//   2. Return the statistics as a JSON response
func (s *APIServer) handleGetEngagementStats(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/statistics/engagements request")
	
	// Step 1: Retrieve engagement statistics from database
	logger.Debug("Retrieving engagement statistics")
	stats, err := s.store.GetEngagementStatistics()
	if err != nil {
		logger.Error("Failed to retrieve engagement statistics: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve engagement statistics: %v", err)})
	}
	
	// Step 2: Return statistics as JSON response
	logger.Debug("Successfully retrieved engagement statistics")
	return WriteJson(w, http.StatusOK, stats)
}

// handleGetGrowthStats handles GET /api/statistics/growth requests.
//
// This handler retrieves statistics on the growth of the expert database over time,
// showing expert additions by month for a specified period. It requires authentication
// via the requireAuth middleware.
//
// Request parameters:
//   - months (optional): Number of months to include in the report, defaults to 12
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request with optional query parameters
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Parse and validate the months parameter from query string
//   2. Retrieve growth statistics from the database for the specified period
//   3. Return the statistics as a JSON response
func (s *APIServer) handleGetGrowthStats(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/statistics/growth request")
	
	// Step 1: Parse and validate months parameter
	months := 12 // Default to 12 months if not specified
	
	monthsParam := r.URL.Query().Get("months")
	if monthsParam != "" {
		parsedMonths, err := strconv.Atoi(monthsParam)
		if err == nil && parsedMonths > 0 {
			months = parsedMonths
			logger.Debug("Using custom months parameter: %d", months)
		} else {
			logger.Warn("Invalid months parameter provided: %s, using default (12)", monthsParam)
		}
	}
	
	// Step 2: Retrieve growth statistics for the specified period
	logger.Debug("Retrieving expert growth statistics for past %d months", months)
	stats, err := s.store.GetExpertGrowthByMonth(months)
	if err != nil {
		logger.Error("Failed to retrieve growth statistics: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve growth statistics: %v", err)})
	}
	
	// Step 3: Return statistics as JSON response
	logger.Debug("Successfully retrieved growth statistics for %d months", months)
	return WriteJson(w, http.StatusOK, stats)
}

// Expert Request Handler Functions

// handleCreateExpertRequest handles POST /api/expert-requests requests.
//
// This handler creates a new expert request in the database. Expert requests are
// submissions from users to add a new expert to the database, which require admin approval.
// It requires authentication via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert request data in the body
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Parse and validate the request body
//   2. Validate required fields
//   3. Set default values
//   4. Create the expert request record in the database
//   5. Return the created request with its assigned ID
func (s *APIServer) handleCreateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing POST /api/expert-requests request")
	
	// Step 1: Parse request body
	var request ExpertRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Warn("Failed to parse expert request: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Step 2: Validate required fields
	logger.Debug("Validating expert request fields")
	
	// Step 2.1: Validate name
	if strings.TrimSpace(request.Name) == "" {
		logger.Warn("Missing name in expert request")
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Name is required"})
	}
	
	// Step 2.2: Validate institution
	if strings.TrimSpace(request.Institution) == "" {
		logger.Warn("Missing institution in expert request")
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Institution is required"})
	}
	
	// Step 2.3: Validate designation
	if strings.TrimSpace(request.Designation) == "" {
		logger.Warn("Missing designation in expert request")
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Designation is required"})
	}
	
	// Step 2.4: Validate contact information (email or phone)
	if strings.TrimSpace(request.Email) == "" && strings.TrimSpace(request.Phone) == "" {
		logger.Warn("Missing contact information in expert request")
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "At least one contact method (email or phone) is required"})
	}
	
	// Step 3: Set default values
	logger.Debug("Setting default values for expert request")
	request.Status = "pending" // Default status for new requests
	request.CreatedAt = time.Now()
	
	// Step 4: Save the request to the database
	logger.Debug("Creating expert request in database: %s, Institution: %s", 
		request.Name, request.Institution)
	id, err := s.store.CreateExpertRequest(&request)
	if err != nil {
		logger.Error("Failed to create expert request in database: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to create expert request: %v", err)})
	}
	
	// Step 5: Set the ID in the response and return
	request.ID = id
	logger.Info("Expert request created successfully: ID: %d, Name: %s", id, request.Name)
	return WriteJson(w, http.StatusCreated, request)
}

// handleGetExpertRequests handles GET /api/expert-requests requests.
//
// This handler retrieves a paginated list of expert requests, with optional filtering
// by status. It requires authentication via the requireAuth middleware.
//
// Request parameters:
//   - status (optional): Filter requests by status (e.g., "pending", "approved", "rejected")
//   - limit (optional): Number of items per page (defaults to 100)
//   - offset (optional): Starting position for pagination
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request with optional query parameters
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Parse and extract query parameters for filtering
//   2. Apply pagination parameters
//   3. Retrieve filtered expert requests from the database
//   4. Return the expert requests as a JSON response
func (s *APIServer) handleGetExpertRequests(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	logger.Debug("Processing GET /api/expert-requests request")
	
	// Step 1: Parse query parameters for filtering
	status := r.URL.Query().Get("status")
	if status != "" {
		logger.Debug("Filtering expert requests by status: %s", status)
	}
	
	// Step 2: Parse pagination parameters using the helper function
	// Use custom default limit of 100 for expert requests
	limit, offset := parsePaginationParams(r, 100)
	logger.Debug("Using pagination: limit=%d, offset=%d", limit, offset)
	
	// Step 3: Build filters and retrieve requests
	filters := make(map[string]interface{})
	if status != "" {
		filters["status"] = status
	}
	
	// Get requests with filters from database
	logger.Debug("Retrieving expert requests with filters: %v", filters)
	requests, err := s.store.ListExpertRequests(filters, limit, offset)
	if err != nil {
		logger.Error("Failed to retrieve expert requests: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to retrieve expert requests: %v", err)})
	}
	
	// Step 4: Return expert requests as JSON response
	logger.Debug("Returning %d expert requests", len(requests))
	return WriteJson(w, http.StatusOK, requests)
}

// handleGetExpertRequest handles GET /api/expert-requests/{id} requests.
//
// This handler retrieves a single expert request by ID. It requires authentication
// via the requireAuth middleware.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert request ID in the path
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the expert request ID from the request path
//   2. Retrieve the expert request from the database
//   3. Return the expert request data in JSON format
func (s *APIServer) handleGetExpertRequest(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate expert request ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid expert request ID provided: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request ID"})
	}
	
	// Step 2: Retrieve expert request from database
	logger.Debug("Retrieving expert request with ID: %d", id)
	request, err := s.store.GetExpertRequest(id)
	if err != nil {
		logger.Warn("Expert request not found with ID: %d - %v", id, err)
		return WriteJson(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Expert request not found: %v", err)})
	}
	
	// Step 3: Return expert request data
	logger.Debug("Successfully retrieved expert request: ID: %d, Name: %s", request.ID, request.Name)
	return WriteJson(w, http.StatusOK, request)
}

// handleUpdateExpertRequest handles PUT /api/expert-requests/{id} requests.
//
// This handler updates an existing expert request by ID. It requires admin privileges,
// which is enforced by the requireAdmin middleware. The handler includes special logic
// for handling request approval, which creates a new expert record in the database.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer for returning results
//   - r (*http.Request): The HTTP request containing the expert request ID in path and update data in body
//
// Returns:
//   - error: Any error that occurs during processing
//
// Flow:
//   1. Extract and validate the expert request ID from the request path
//   2. Retrieve the existing expert request from the database
//   3. Parse and validate the update data from the request body
//   4. Handle status changes, particularly approval which creates an expert record
//   5. Update the expert request record in the database
//   6. Return success response
func (s *APIServer) handleUpdateExpertRequest(w http.ResponseWriter, r *http.Request) error {
	logger := GetLogger()
	
	// Step 1: Extract and validate expert request ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn("Invalid expert request ID provided for update: %s", idStr)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request ID"})
	}
	
	// Step 2: Retrieve existing expert request from database
	logger.Debug("Checking if expert request exists with ID: %d", id)
	existingRequest, err := s.store.GetExpertRequest(id)
	if err != nil {
		logger.Warn("Expert request not found for update ID: %d - %v", id, err)
		return WriteJson(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Expert request not found: %v", err)})
	}
	
	// Step 3: Parse update data
	var updateRequest ExpertRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		logger.Warn("Failed to parse expert request update: %v", err)
		return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
	}
	
	// Ensure ID matches path parameter
	updateRequest.ID = id
	
	// Step 4: Handle status changes - if status is changing to "approved", create an expert record
	logger.Debug("Processing request update, current status: %s, new status: %s", 
		existingRequest.Status, updateRequest.Status)
	
	if existingRequest.Status != "approved" && updateRequest.Status == "approved" {
		logger.Info("Expert request being approved, creating expert record from request data")
		
		// Step 4.1: Create a new expert from the request data
		// Generate a unique expert ID if not provided or too short
		expertIDStr := updateRequest.ExpertID
		if expertIDStr == "" || len(expertIDStr) < 3 {
			expertIDStr = fmt.Sprintf("EXP-%d-%d", id, time.Now().Unix())
			logger.Debug("Generated unique expert ID: %s", expertIDStr)
		}
		
		// Step 4.2: Create expert record with data from the request
		expert := &Expert{
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
		
		// Step 4.3: Create the expert record in database
		logger.Debug("Creating expert record: %s, Institution: %s", expert.Name, expert.Institution)
		createdID, err := s.store.CreateExpert(expert)
		if err != nil {
			logger.Error("Failed to create expert from request: %v", err)
			return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to create expert from request: %v", err)})
		}
		
		// Step 4.4: Set the reviewed timestamp and reviewer info
		updateRequest.ReviewedAt = time.Now()
		// Note: In a real app, we'd get the reviewer ID from the authenticated user
		
		// Step 4.5: Update the expert request with the expert ID and other review info
		updateRequest.ExpertID = fmt.Sprintf("EXP-%d", createdID)
		logger.Info("Expert created successfully from request: Expert ID: %d", createdID)
	}
	
	// Step 5: Update the expert request in database
	logger.Debug("Updating expert request ID: %d, Status: %s", id, updateRequest.Status)
	if err := s.store.UpdateExpertRequest(&updateRequest); err != nil {
		logger.Error("Failed to update expert request in database: %v", err)
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: fmt.Sprintf("Failed to update expert request: %v", err)})
	}
	
	// Step 6: Return success response
	logger.Info("Expert request updated successfully: ID: %d, Status: %s", id, updateRequest.Status)
	return WriteJson(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Expert request updated successfully",
	})
}

// Helper functions for API request handling

// apiFunc is a custom HTTP handler function type that can return errors
// This enables more concise error handling in HTTP handlers
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError represents a standardized error response for API clients
type ApiError struct {
	Error string `json:"error"` // Human-readable error message
}

// parsePaginationParams extracts and validates pagination parameters from request query
//
// This helper standardizes pagination parameter handling across different endpoints
// and applies sensible defaults when parameters are missing or invalid.
//
// Inputs:
//   - r (*http.Request): The HTTP request containing query parameters
//   - defaultLimit (int): Optional custom default limit (uses DefaultLimit constant if 0)
//
// Returns:
//   - limit (int): Number of items per page
//   - offset (int): Starting position (number of items to skip)
func parsePaginationParams(r *http.Request, defaultLimit int) (limit, offset int) {
	// Use the provided default limit or fall back to the constant
	if defaultLimit <= 0 {
		defaultLimit = DefaultLimit
	}

	// Parse limit parameter with validation
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}
	
	// Parse offset parameter with validation
	offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = DefaultOffset
	}
	
	return limit, offset
}

// WriteJson writes a JSON response with the appropriate headers
//
// This helper standardizes JSON response formatting across the API.
//
// Inputs:
//   - w (http.ResponseWriter): The response writer to write to
//   - status (int): The HTTP status code to set
//   - v (any): The value to encode as JSON
//
// Returns:
//   - error: Any error that occurs during JSON encoding
func WriteJson(w http.ResponseWriter, status int, v any) error {
	// Set appropriate content type header
	w.Header().Add("Content-Type", "application/json")
	
	// Set the HTTP status code
	w.WriteHeader(status)
	
	// Encode and write the JSON response
	return json.NewEncoder(w).Encode(v)
}

// makeHTTPHandleFunc wraps an apiFunc to handle errors consistently
//
// This adapter converts our custom apiFunc handlers to standard http.HandlerFunc
// and provides centralized error handling logic.
//
// Inputs:
//   - f (apiFunc): The API handler function to wrap
//
// Returns:
//   - http.HandlerFunc: A standard HTTP handler with error handling
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := GetLogger()
		
		// Execute the handler function
		if err := f(w, r); err != nil {
			// Get request details for logging
			path := r.URL.Path
			method := r.Method
			
			// Handle specific error types differently
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
				// Generic error handler for unrecognized errors
				logger.Error("Handler error: %s %s - %v", method, path, err)
				WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
			}
		}
	}
}