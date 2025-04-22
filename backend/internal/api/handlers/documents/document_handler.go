// Package documents provides handlers for document-related API endpoints
package documents

import (
	"fmt"
	"net/http"
	"strconv"
	
	"expertdb/internal/api/response"
	"expertdb/internal/documents"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// Handler manages document-related HTTP endpoints
type Handler struct {
	store           storage.Storage
	documentService *documents.Service
}

// NewHandler creates a new document handler
func NewHandler(store storage.Storage, documentService *documents.Service) *Handler {
	return &Handler{
		store:           store,
		documentService: documentService,
	}
}

// HandleUploadDocument handles POST /api/documents requests
func (h *Handler) HandleUploadDocument(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing POST /api/documents request")
	
	// Parse multipart form data (10MB max)
	maxSize := int64(10 << 20)
	if err := r.ParseMultipartForm(maxSize); err != nil {
		log.Warn("Failed to parse multipart form: %v", err)
		return fmt.Errorf("failed to parse form - file may be too large: %w", err)
	}
	
	// Extract and validate expert ID
	expertIDStr := r.FormValue("expertId")
	if expertIDStr == "" {
		log.Warn("Missing expert ID in document upload request")
		return fmt.Errorf("expert ID is required")
	}
	
	expertID, err := strconv.ParseInt(expertIDStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert ID in document upload request: %s", expertIDStr)
		return fmt.Errorf("invalid expert ID: %w", err)
	}
	
	// Get document type or use default
	docType := r.FormValue("documentType")
	if docType == "" {
		log.Debug("No document type specified, using default type: cv")
		docType = "cv" // Default type
	}
	
	// Get the file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Warn("No file provided in document upload request: %v", err)
		return fmt.Errorf("no file provided: %w", err)
	}
	defer file.Close()
	
	// Upload and store the document
	log.Debug("Uploading document for expert ID: %d, type: %s, filename: %s",
		expertID, docType, header.Filename)
	doc, err := h.documentService.CreateDocument(expertID, file, header, docType)
	if err != nil {
		log.Error("Failed to upload document: %v", err)
		return fmt.Errorf("failed to upload document: %w", err)
	}
	
	// Return document information with standardized response
	log.Info("Document uploaded successfully: ID: %d, Type: %s, Expert: %d", doc.ID, doc.Type, doc.ExpertID)
	return response.Success(w, http.StatusCreated, "Document uploaded successfully", doc)
}

// HandleGetDocument handles GET /api/documents/{id} requests
func (h *Handler) HandleGetDocument(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate document ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid document ID provided: %s", idStr)
		return fmt.Errorf("invalid document ID: %w", err)
	}
	
	// Retrieve document from document service
	log.Debug("Retrieving document with ID: %d", id)
	doc, err := h.documentService.GetDocument(id)
	if err != nil {
		log.Warn("Document not found with ID: %d - %v", id, err)
		return fmt.Errorf("document not found: %w", err)
	}
	
	// Return document information with standardized response
	log.Debug("Returning document: ID: %d, Type: %s, Expert: %d", doc.ID, doc.Type, doc.ExpertID)
	return response.Success(w, http.StatusOK, "", doc)
}

// HandleDeleteDocument handles DELETE /api/documents/{id} requests
func (h *Handler) HandleDeleteDocument(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate document ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid document ID provided for deletion: %s", idStr)
		return fmt.Errorf("invalid document ID: %w", err)
	}
	
	// Delete the document
	log.Debug("Deleting document with ID: %d", id)
	if err := h.documentService.DeleteDocument(id); err != nil {
		log.Error("Failed to delete document: %v", err)
		return fmt.Errorf("failed to delete document: %w", err)
	}
	
	// Return success response
	log.Info("Document deleted successfully: ID: %d", id)
	return response.Success(w, http.StatusOK, "Document deleted successfully", nil)
}

// HandleGetExpertDocuments handles GET /api/experts/{id}/documents requests
func (h *Handler) HandleGetExpertDocuments(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate expert ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid expert ID provided for document retrieval: %s", idStr)
		return fmt.Errorf("invalid expert ID: %w", err)
	}
	
	// Retrieve the expert's documents
	log.Debug("Retrieving documents for expert with ID: %d", id)
	docs, err := h.documentService.ListDocuments(id)
	if err != nil {
		log.Error("Failed to retrieve documents for expert %d: %v", id, err)
		return fmt.Errorf("failed to retrieve documents: %w", err)
	}
	
	// Return documents with standardized response
	log.Debug("Returning %d documents for expert ID: %d", len(docs), id)
	responseData := map[string]interface{}{
		"documents": docs,
		"count":     len(docs),
		"expertId":  id,
	}
	return response.Success(w, http.StatusOK, "", responseData)
}