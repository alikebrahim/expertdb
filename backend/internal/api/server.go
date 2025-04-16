// Package api provides the HTTP API server for the ExpertDB application
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"expertdb/internal/api/handlers"
	"expertdb/internal/api/handlers/documents"
	"expertdb/internal/api/handlers/engagements"
	"expertdb/internal/api/handlers/statistics"
	"expertdb/internal/config"
	"expertdb/internal/documents"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// Server represents the HTTP API server for the ExpertDB application
type Server struct {
	listenAddr      string
	store           storage.Storage
	documentService *documents.Service
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
func NewServer(listenAddr string, store storage.Storage, docService *documents.Service, cfg *config.Configuration) (*Server, error) {
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
	// Health check endpoint
	s.mux.Handle("GET /api/health", corsAndLogMiddleware(http.HandlerFunc(s.handleHealth)))
	
	// Create handlers
	userHandler := handlers.NewUserHandler(s.store)
	expertHandler := handlers.NewExpertHandler(s.store)
	expertRequestHandler := handlers.NewExpertRequestHandler(s.store, s.documentService)
	documentHandler := documents.NewHandler(s.store, s.documentService)
	engagementHandler := engagements.NewHandler(s.store)
	statisticsHandler := statistics.NewHandler(s.store)
	
	// Expert endpoints
	s.mux.Handle("GET /api/experts", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertHandler.HandleGetExperts(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("POST /api/experts", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertHandler.HandleCreateExpert(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/experts/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertHandler.HandleGetExpert(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("PUT /api/experts/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertHandler.HandleUpdateExpert(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("DELETE /api/experts/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertHandler.HandleDeleteExpert(w, r); err != nil {
			handleError(w, err)
		}
	})))
	
	// Expert request endpoints
	s.mux.Handle("POST /api/expert-requests", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertRequestHandler.HandleCreateExpertRequest(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/expert-requests", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertRequestHandler.HandleGetExpertRequests(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/expert-requests/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertRequestHandler.HandleGetExpertRequest(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("PUT /api/expert-requests/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertRequestHandler.HandleUpdateExpertRequest(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/expert/areas", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := expertHandler.HandleGetExpertAreas(w, r); err != nil {
			handleError(w, err)
		}
	})))

	// Document endpoints
	s.mux.Handle("POST /api/documents", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := documentHandler.HandleUploadDocument(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/documents/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := documentHandler.HandleGetDocument(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("DELETE /api/documents/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := documentHandler.HandleDeleteDocument(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/experts/{id}/documents", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := documentHandler.HandleGetExpertDocuments(w, r); err != nil {
			handleError(w, err)
		}
	})))

	// Engagement endpoints
	s.mux.Handle("POST /api/engagements", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := engagementHandler.HandleCreateEngagement(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/engagements/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := engagementHandler.HandleGetEngagement(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("PUT /api/engagements/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := engagementHandler.HandleUpdateEngagement(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("DELETE /api/engagements/{id}", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := engagementHandler.HandleDeleteEngagement(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/experts/{id}/engagements", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := engagementHandler.HandleGetExpertEngagements(w, r); err != nil {
			handleError(w, err)
		}
	})))

	// Statistics endpoints
	s.mux.Handle("GET /api/statistics", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := statisticsHandler.HandleGetStatistics(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/statistics/nationality", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := statisticsHandler.HandleGetNationalityStats(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/statistics/engagements", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := statisticsHandler.HandleGetEngagementStats(w, r); err != nil {
			handleError(w, err)
		}
	})))
	s.mux.Handle("GET /api/statistics/growth", corsAndLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := statisticsHandler.HandleGetGrowthStats(w, r); err != nil {
			handleError(w, err)
		}
	})))
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