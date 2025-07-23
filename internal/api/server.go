// Package api provides the HTTP API server for the ExpertDB application
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"expertdb/internal/api/handlers"
	"expertdb/internal/api/handlers/backup"
	"expertdb/internal/api/handlers/documents"
	"expertdb/internal/api/handlers/engagements"
	"expertdb/internal/api/handlers/phase"
	"expertdb/internal/api/handlers/statistics"
	"expertdb/internal/auth"
	"expertdb/internal/config"
	docsvc "expertdb/internal/documents"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// Server represents the HTTP API server for the ExpertDB application
type Server struct {
	listenAddr      string
	store           storage.Storage
	documentService *docsvc.Service
	config          *config.Configuration
	mux             *http.ServeMux
}

// writeJSON is a helper function to write JSON responses
func writeJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// errorResponse represents a standard error response
type errorResponse struct {
	Error string `json:"error"`
}

// handleError is a helper function to handle API errors
func handleError(w http.ResponseWriter, err error) {
	log := logger.Get()
	
	// Determine appropriate status code based on error type
	var statusCode int
	
	switch {
	case err == domain.ErrNotFound:
		statusCode = http.StatusNotFound
	case err == domain.ErrUnauthorized:
		statusCode = http.StatusUnauthorized
	case err == domain.ErrForbidden:
		statusCode = http.StatusForbidden
	case err == domain.ErrInvalidCredentials:
		statusCode = http.StatusUnauthorized
	case err == domain.ErrValidation:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}
	
	// Log the error
	if statusCode >= 500 {
		log.Error("Server error: %v", err)
	} else {
		log.Debug("Client error: %v", err)
	}
	
	// Write error response
	resp := errorResponse{Error: err.Error()}
	writeJSON(w, statusCode, resp)
}

// NewServer creates a new API server
func NewServer(listenAddr string, store storage.Storage, docService *docsvc.Service, cfg *config.Configuration) (*Server, error) {
	server := &Server{
		listenAddr:      listenAddr,
		store:           store,
		documentService: docService,
		config:          cfg,
		mux:             http.NewServeMux(),
	}
	
	// Register routes
	server.registerRoutes()
	
	return server, nil
}

// registerRoutes sets up the API routes
func (s *Server) registerRoutes() {
	log := logger.Get()
	
	// Define the middleware for handling CORS and logging
	corsAndLogMiddleware := func(next http.Handler) http.Handler {
		return log.RequestLoggerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add CORS headers
			w.Header().Set("Access-Control-Allow-Origin", s.config.CORSAllowOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			// Call the next handler
			next.ServeHTTP(w, r)
		}))
	}
	
	// Register routes with middleware
	// Create handlers
	expertHandler := handlers.NewExpertHandler(s.store, s.documentService)
	expertRequestHandler := handlers.NewExpertRequestHandler(s.store, s.documentService)
	expertEditRequestHandler := handlers.NewExpertEditRequestHandler(s.store, s.documentService)
	documentHandler := documents.NewHandler(s.store, s.documentService)
	engagementHandler := engagements.NewHandler(s.store)
	statisticsHandler := statistics.NewHandler(s.store)
	backupHandler := backup.NewHandler(s.store)
	userHandler := handlers.NewUserHandler(s.store)
	authHandler := handlers.NewAuthHandler(s.store)
	phaseHandler := phase.NewHandler(s.store)
	roleAssignmentHandler := handlers.NewRoleAssignmentHandler(s.store)
	specializedAreasHandler := handlers.NewSpecializedAreasHandler(s.store)
	
	// Define a generic error handler wrapper for converting HandlerFunc to http.HandlerFunc
	errorHandler := func(h auth.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if err := h(w, r); err != nil {
				handleError(w, err)
			}
		}
	}
	
	//
	// PUBLIC ENDPOINTS (No auth required)
	//
	
	// Global OPTIONS handler for CORS preflight requests
	s.mux.Handle("OPTIONS /", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers are already set by corsAndLogMiddleware
		// Just return 200 OK for all preflight requests
		w.WriteHeader(http.StatusOK)
	})))
	
	// Health check endpoint - public access
	s.mux.Handle("GET /api/health", corsAndLogMiddleware(http.HandlerFunc(s.handleHealth)))
	
	// User authentication endpoints
	s.mux.Handle("POST /api/auth/login", corsAndLogMiddleware(errorHandler(func(w http.ResponseWriter, r *http.Request) error {
		return authHandler.HandleLogin(w, r)
	})))
	
	// Get expert areas - authenticated user access (Phase 8A: Area Access Extension)
	s.mux.Handle("GET /api/expert/areas", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertHandler.HandleGetExpertAreas(w, r)
	}))))
	
	// Get specialized areas - authenticated user access (for search functionality)
	s.mux.Handle("GET /api/specialized-areas", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return specializedAreasHandler.HandleListSpecializedAreas(w, r)
	}))))
	
	// Create expert area - admin access (Phase 8B: Area Creation)
	s.mux.Handle("POST /api/expert/areas", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return expertHandler.HandleCreateArea(w, r)
	}))))
	
	// Update expert area - admin access (Phase 8C: Area Renaming)
	s.mux.Handle("PUT /api/expert/areas/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return expertHandler.HandleUpdateArea(w, r)
	}))))
	
	//
	// USER ACCESS (All authenticated users)
	//
	
	// Read-only expert endpoints
	s.mux.Handle("GET /api/experts", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertHandler.HandleGetExperts(w, r)
	}))))
	
	s.mux.Handle("GET /api/experts/{id}", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertHandler.HandleGetExpert(w, r)
	}))))
	
	// Read-only document endpoints
	s.mux.Handle("GET /api/documents/{id}", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return documentHandler.HandleGetDocument(w, r)
	}))))
	
	
	// Read-only engagement endpoints
	s.mux.Handle("GET /api/engagements/{id}", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return engagementHandler.HandleGetEngagement(w, r)
	}))))
	
	// Phase 11A: Global engagement listing with filtering capability
	s.mux.Handle("GET /api/engagements", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return engagementHandler.HandleListEngagements(w, r)
	}))))
	
	s.mux.Handle("GET /api/experts/{id}/engagements", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return engagementHandler.HandleGetExpertEngagements(w, r)
	}))))
	
	//
	// PLANNER ACCESS
	//
	
	// Engagement management endpoints (create, update, delete) - legacy system, admin only
	s.mux.Handle("POST /api/engagements", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return engagementHandler.HandleCreateEngagement(w, r)
	}))))
	
	s.mux.Handle("PUT /api/engagements/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return engagementHandler.HandleUpdateEngagement(w, r)
	}))))
	
	s.mux.Handle("DELETE /api/engagements/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return engagementHandler.HandleDeleteEngagement(w, r)
	}))))
	
	// Phase 11C: Engagement import endpoint - restricted to admin/developer role
	s.mux.Handle("POST /api/engagements/import", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return engagementHandler.HandleImportEngagements(w, r)
	}))))
	
	//
	// ADMIN ACCESS
	//
	
	// Expert management (create, update, delete)
	s.mux.Handle("POST /api/experts", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return expertHandler.HandleCreateExpert(w, r)
	}))))
	
	s.mux.Handle("PUT /api/experts/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return expertHandler.HandleUpdateExpert(w, r)
	}))))
	
	s.mux.Handle("DELETE /api/experts/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return expertHandler.HandleDeleteExpert(w, r)
	}))))
	
	// Expert request management
	s.mux.Handle("POST /api/expert-requests", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertRequestHandler.HandleCreateExpertRequest(w, r)
	}))))
	
	s.mux.Handle("GET /api/expert-requests", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertRequestHandler.HandleGetExpertRequests(w, r)
	}))))
	
	s.mux.Handle("GET /api/expert-requests/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return expertRequestHandler.HandleGetExpertRequest(w, r)
	}))))
	
	s.mux.Handle("PUT /api/expert-requests/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return expertRequestHandler.HandleUpdateExpertRequest(w, r)
	}))))
	
	s.mux.Handle("PUT /api/expert-requests/{id}/edit", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertRequestHandler.HandleEditExpertRequest(w, r)
	}))))
	
	// Batch approval endpoint for multiple expert requests
	s.mux.Handle("POST /api/expert-requests/batch-approve", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return expertRequestHandler.HandleBatchApproveExpertRequests(w, r)
	}))))

	// Expert profile edit system (for requesting changes to existing expert profiles)
	// Create edit proposal - authenticated users can request edits to existing experts
	s.mux.Handle("POST /api/experts/{id}/edit", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertEditRequestHandler.HandleCreateExpertEditRequest(w, r)
	}))))

	// List expert edit proposals - users see their own, admins see all
	s.mux.Handle("GET /api/expert-edits", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertEditRequestHandler.HandleGetExpertEditRequests(w, r)
	}))))

	// Get specific expert edit proposal - users see their own, admins see all
	s.mux.Handle("GET /api/expert-edits/{id}", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertEditRequestHandler.HandleGetExpertEditRequest(w, r)
	}))))

	// Update expert edit proposal status - admin approve/reject (auto-applies if approved)
	s.mux.Handle("PUT /api/expert-edits/{id}/status", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return expertEditRequestHandler.HandleUpdateExpertEditRequestStatus(w, r)
	}))))
	
	// Document management (upload, delete)
	s.mux.Handle("POST /api/documents", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return documentHandler.HandleUploadDocument(w, r)
	}))))
	
	s.mux.Handle("DELETE /api/documents/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return documentHandler.HandleDeleteDocument(w, r)
	}))))
	
	//
	// SUPER USER ACCESS
	//
	
	// Statistics endpoint (accessible to all authenticated users)
	s.mux.Handle("GET /api/statistics", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return statisticsHandler.HandleGetStatistics(w, r)
	}))))
	
	// Backup endpoints (restricted to admins)
	s.mux.Handle("GET /api/backup", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return backupHandler.HandleBackupCSV(w, r)
	}))))
	
	// Role assignment endpoints (admin access only)
	s.mux.Handle("POST /api/users/{id}/planner-assignments", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		roleAssignmentHandler.AssignPlannerApplications(w, r)
		return nil
	}))))
	s.mux.Handle("POST /api/users/{id}/manager-assignments", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		roleAssignmentHandler.AssignManagerApplications(w, r)
		return nil
	}))))
	s.mux.Handle("DELETE /api/users/{id}/planner-assignments", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		roleAssignmentHandler.RemovePlannerAssignments(w, r)
		return nil
	}))))
	s.mux.Handle("DELETE /api/users/{id}/manager-assignments", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		roleAssignmentHandler.RemoveManagerAssignments(w, r)
		return nil
	}))))
	s.mux.Handle("GET /api/users/{id}/assignments", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		roleAssignmentHandler.GetUserAssignments(w, r)
		return nil
	}))))
	
	// Phase planning endpoints
	// Phase listing - all authenticated users can view phases
	s.mux.Handle("GET /api/phases", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return phaseHandler.HandleListPhases(w, r)
	}))))
	
	// Get specific phase - all authenticated users can view phases
	s.mux.Handle("GET /api/phases/{id}", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return phaseHandler.HandleGetPhase(w, r)
	}))))
	
	// Create phase - admin access
	s.mux.Handle("POST /api/phases", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return phaseHandler.HandleCreatePhase(w, r)
	}))))
	
	// Update phase - admin access
	s.mux.Handle("PUT /api/phases/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return phaseHandler.HandleUpdatePhase(w, r)
	}))))
	
	// Update application experts - planner access (context-aware)
	s.mux.Handle("PUT /api/phases/{id}/applications/{app_id}", corsAndLogMiddleware(errorHandler(auth.RequirePlannerForApplication(s.store, phaseHandler.HandleUpdateApplicationExperts))))
	
	// Review application - admin access
	s.mux.Handle("PUT /api/phases/{id}/applications/{app_id}/review", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return phaseHandler.HandleReviewApplication(w, r)
	}))))
	
	// User management endpoints (most restricted to super_user)
	// Get own user profile - any authenticated user can access their own profile
	s.mux.Handle("GET /api/users/me", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		// Get user ID from context
		userID, ok := auth.GetUserIDFromContext(r.Context())
		if !ok {
			return domain.ErrUnauthorized
		}
		
		// Set the ID in the URL path
		r = r.WithContext(r.Context())
		r.URL.Path = fmt.Sprintf("/api/users/%d", userID)
		
		return userHandler.HandleGetUser(w, r)
	}))))
	
	// User list - admin access
	s.mux.Handle("GET /api/users", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return userHandler.HandleGetUsers(w, r)
	}))))
	
	// Get specific user - admin access
	s.mux.Handle("GET /api/users/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return userHandler.HandleGetUser(w, r)
	}))))
	
	// Create user - super user and admin access
	s.mux.Handle("POST /api/users", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return userHandler.HandleCreateUser(w, r)
	}))))
	
	// Update user - admin access with role-based restrictions
	s.mux.Handle("PUT /api/users/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return userHandler.HandleUpdateUser(w, r)
	}))))
	
	// Delete user - super user and admin access
	s.mux.Handle("DELETE /api/users/{id}", corsAndLogMiddleware(errorHandler(auth.RequireRole(auth.RoleAdmin, func(w http.ResponseWriter, r *http.Request) error {
		return userHandler.HandleDeleteUser(w, r)
	}))))
	
	// Manager-specific endpoints
	// Rate experts - manager access (context-aware)
	s.mux.Handle("POST /api/phases/{id}/applications/{app_id}/ratings", corsAndLogMiddleware(errorHandler(auth.RequireManagerForApplication(s.store, phaseHandler.HandleRateExperts))))
	
	// Get manager tasks - any authenticated user can see their own tasks
	s.mux.Handle("GET /api/users/me/manager-tasks", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return phaseHandler.HandleGetManagerTasks(w, r)
	}))))
	
	// Application listing with filtering - all authenticated users can view
	s.mux.Handle("GET /api/applications", corsAndLogMiddleware(errorHandler(auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) error {
		return phaseHandler.HandleListApplications(w, r)
	}))))
}

// Run starts the HTTP server
func (s *Server) Run() error {
	log := logger.Get()
	log.Info("API server listening on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.mux)
}

// Handler for health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"status": "ok",
		"message": "ExpertDB API is running",
	}
	writeJSON(w, http.StatusOK, resp)
}