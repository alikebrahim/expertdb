// Package documents provides handlers for document-related API endpoints
package documents

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	
	"expertdb/internal/api/utils"
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
	log.Info("Document uploaded successfully: ID: %d, Type: %s, Expert: %d", doc.ID, doc.DocumentType, doc.ExpertID)
	return utils.RespondWithSuccess(w, "Document uploaded successfully", doc)
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
	log.Debug("Returning document: ID: %d, Type: %s, Expert: %d", doc.ID, doc.DocumentType, doc.ExpertID)
	return utils.RespondWithSuccess(w, "", doc)
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
	return utils.RespondWithSuccess(w, "Document deleted successfully", nil)
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
	return utils.RespondWithSuccess(w, "", responseData)
}

// HandleDownloadDocument handles GET /api/documents/{id}/download requests
func (h *Handler) HandleDownloadDocument(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	
	// Extract and validate document ID from path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Warn("Invalid document ID provided for download: %s", idStr)
		return fmt.Errorf("invalid document ID: %w", err)
	}
	
	// Retrieve document metadata from document service
	log.Debug("Retrieving document metadata for download: ID %d", id)
	doc, err := h.documentService.GetDocument(id)
	if err != nil {
		log.Warn("Document not found for download: ID %d - %v", id, err)
		return fmt.Errorf("document not found: %w", err)
	}
	
	// Open the file for reading
	log.Debug("Opening file for download: %s", doc.FilePath)
	file, err := os.Open(doc.FilePath)
	if err != nil {
		log.Error("Failed to open document file for download: %s - %v", doc.FilePath, err)
		return fmt.Errorf("failed to open document file: %w", err)
	}
	defer file.Close()
	
	// Get file info for size verification
	fileInfo, err := file.Stat()
	if err != nil {
		log.Error("Failed to get file info for download: %s - %v", doc.FilePath, err)
		return fmt.Errorf("failed to get file info: %w", err)
	}
	
	// Set appropriate headers for file download
	w.Header().Set("Content-Type", doc.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", doc.Filename))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	
	// Stream the file content to the response
	log.Debug("Streaming file content for download: %s (%d bytes)", doc.Filename, fileInfo.Size())
	bytesWritten, err := io.Copy(w, file)
	if err != nil {
		log.Error("Failed to stream file content for download: %v", err)
		return fmt.Errorf("failed to stream file content: %w", err)
	}
	
	log.Info("Document downloaded successfully: ID %d, Filename: %s, Bytes: %d", 
		doc.ID, doc.Filename, bytesWritten)
	
	return nil
}