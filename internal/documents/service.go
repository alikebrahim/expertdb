// Package documents provides document management functionality
package documents

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// Service manages document uploads and storage
type Service struct {
	store       storage.Storage
	uploadDir   string
	maxSize     int64
	allowedTypes map[string]bool
}

// New creates a new Service instance
func New(store storage.Storage, uploadDir string) (*Service, error) {
	// Create the upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &Service{
		store:     store,
		uploadDir: uploadDir,
		maxSize:   10 * 1024 * 1024, // 10 MB default limit
		allowedTypes: map[string]bool{
			"application/pdf":                                       true,
			"application/msword":                                    true,
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
			"image/jpeg":                                            true,
			"image/png":                                             true,
		},
	}, nil
}

// CreateDocument handles core file upload and database registration only
func (s *Service) CreateDocument(expertID int64, file multipart.File, header *multipart.FileHeader, docType string) (*domain.Document, error) {
	log := logger.Get()
	log.Debug("Creating document: expertID=%d, docType='%s', filename='%s'", expertID, docType, header.Filename)
	
	// Validate document type
	validTypes := map[string]bool{
		"cv":       true,
		"approval": true,
	}
	
	if !validTypes[docType] {
		return nil, fmt.Errorf("document type '%s' is not allowed; must be one of: cv, approval", docType)
	}
	
	// Validate file size
	if header.Size > s.maxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxSize)
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if !s.allowedTypes[contentType] {
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}

	// Determine directory and filename based on document type
	var targetDir, filename string
	timestamp := time.Now().Format("20060102_150405")
	extension := filepath.Ext(header.Filename)
	
	switch docType {
	case "cv":
		targetDir = filepath.Join(s.uploadDir, "experts")
		filename = fmt.Sprintf("cv_%d_%s%s", expertID, timestamp, extension)
	case "approval":
		targetDir = filepath.Join(s.uploadDir, "approvals")
		filename = fmt.Sprintf("approval_%d_%s%s", expertID, timestamp, extension)
	}
	
	// Create the target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(targetDir, filename)

	// Create and write file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy the file data
	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create document record
	doc := &domain.Document{
		ExpertID:     expertID,
		DocumentType: docType,
		Filename:     header.Filename,
		FilePath:     filePath,
		ContentType:  contentType,
		FileSize:     header.Size,
		UploadDate:   time.Now(),
	}

	// Store in database
	docID, err := s.store.CreateDocument(doc)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to store document in database: %w", err)
	}

	doc.ID = docID
	log.Debug("Document created successfully with ID: %d", docID)
	return doc, nil
}

// LinkDocumentToExpert updates expert table references to point to a document
func (s *Service) LinkDocumentToExpert(expertID, docID int64, docType string) error {
	log := logger.Get()
	log.Debug("Linking document %d to expert %d as %s", docID, expertID, docType)
	
	var err error
	switch docType {
	case "cv":
		err = s.store.UpdateExpertCVDocument(expertID, docID)
	case "approval":
		err = s.store.UpdateExpertApprovalDocument(expertID, docID)
	default:
		return fmt.Errorf("invalid document type for expert linking: %s", docType)
	}
	
	if err != nil {
		return fmt.Errorf("failed to link %s document to expert: %w", docType, err)
	}
	
	log.Debug("Successfully linked document %d to expert %d", docID, expertID)
	return nil
}

// GetExpertDocument retrieves the current document of specified type for an expert
func (s *Service) GetExpertDocument(expertID int64, docType string) (*domain.Document, error) {
	log := logger.Get()
	log.Debug("Getting %s document for expert %d", docType, expertID)
	
	// Query for the current document based on expert table reference
	var query string
	switch docType {
	case "cv":
		query = `
			SELECT d.id, d.expert_id, d.document_type, d.filename, d.file_path, 
			       d.content_type, d.file_size, d.upload_date
			FROM expert_documents d
			INNER JOIN experts e ON d.id = e.cv_document_id
			WHERE e.id = ? AND d.document_type = 'cv'
		`
	case "approval":
		query = `
			SELECT d.id, d.expert_id, d.document_type, d.filename, d.file_path, 
			       d.content_type, d.file_size, d.upload_date
			FROM expert_documents d
			INNER JOIN experts e ON d.id = e.approval_document_id
			WHERE e.id = ? AND d.document_type = 'approval'
		`
	default:
		return nil, fmt.Errorf("invalid document type: %s", docType)
	}
	
	var doc domain.Document
	db := s.store.GetDB().(*sql.DB)
	err := db.QueryRow(query, expertID).Scan(
		&doc.ID, &doc.ExpertID, &doc.DocumentType, &doc.Filename, &doc.FilePath,
		&doc.ContentType, &doc.FileSize, &doc.UploadDate,
	)
	
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Debug("No %s document found for expert %d", docType, expertID)
			return nil, nil // No document found, not an error
		}
		return nil, fmt.Errorf("failed to get expert document: %w", err)
	}
	
	log.Debug("Found %s document with ID %d for expert %d", docType, doc.ID, expertID)
	return &doc, nil
}

// ReplaceExpertDocument replaces an expert's document with cleanup of old document
func (s *Service) ReplaceExpertDocument(expertID int64, file multipart.File, header *multipart.FileHeader, docType string) (*domain.Document, error) {
	log := logger.Get()
	log.Debug("Replacing %s document for expert %d", docType, expertID)
	
	// Get current document for cleanup (if exists)
	oldDoc, err := s.GetExpertDocument(expertID, docType)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing document: %w", err)
	}
	
	// Create new document
	newDoc, err := s.CreateDocument(expertID, file, header, docType)
	if err != nil {
		return nil, fmt.Errorf("failed to create new document: %w", err)
	}
	
	// Link new document to expert
	if err := s.LinkDocumentToExpert(expertID, newDoc.ID, docType); err != nil {
		// Rollback: delete the new document we just created
		s.store.DeleteDocument(newDoc.ID)
		if newDoc.FilePath != "" {
			os.Remove(newDoc.FilePath)
		}
		return nil, fmt.Errorf("failed to link new document: %w", err)
	}
	
	// Clean up old document if it exists
	if oldDoc != nil {
		log.Debug("Cleaning up old %s document with ID %d", docType, oldDoc.ID)
		
		// Delete from database
		if err := s.store.DeleteDocument(oldDoc.ID); err != nil {
			log.Warn("Failed to delete old document from database (ID: %d): %v", oldDoc.ID, err)
			// Continue anyway - new document is already linked
		}
		
		// Delete file from filesystem
		if oldDoc.FilePath != "" {
			if err := os.Remove(oldDoc.FilePath); err != nil {
				log.Warn("Failed to delete old document file (%s): %v", oldDoc.FilePath, err)
				// Continue anyway - file might already be missing
			} else {
				log.Debug("Deleted old document file: %s", oldDoc.FilePath)
			}
		}
	}
	
	log.Debug("Successfully replaced %s document for expert %d (new ID: %d)", docType, expertID, newDoc.ID)
	return newDoc, nil
}


// CreateDocumentForRequest handles document upload for expert requests using request ID
func (s *Service) CreateDocumentForRequest(requestID int64, file multipart.File, header *multipart.FileHeader, docType string) (*domain.Document, error) {
	log := logger.Get()
	log.Debug("Creating %s document for expert request %d", docType, requestID)
	
	// Validate document type
	if docType != "cv" && docType != "approval" {
		return nil, fmt.Errorf("invalid document type '%s' for request; must be 'cv' or 'approval'", docType)
	}
	
	// Validate file size
	if header.Size > s.maxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxSize)
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if !s.allowedTypes[contentType] {
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}

	// Create expert_requests directory
	targetDir := filepath.Join(s.uploadDir, "expert_requests")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate filename for expert request
	timestamp := time.Now().Format("20060102_150405")
	extension := filepath.Ext(header.Filename)
	var filename string
	if docType == "approval" {
		filename = fmt.Sprintf("expert_request_%d_approval_%s%s", requestID, timestamp, extension)
	} else {
		filename = fmt.Sprintf("expert_request_%d_%s%s", requestID, timestamp, extension)
	}
	filePath := filepath.Join(targetDir, filename)

	// Create and write file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy the file data
	bytesWritten, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	log.Debug("File written successfully: %d bytes to %s", bytesWritten, filePath)

	// Create document record
	doc := &domain.Document{
		ExpertID:     requestID, // Use expert_request id during request creation
		DocumentType: docType,
		Filename:     header.Filename,
		FilePath:     filePath,
		ContentType:  contentType,
		FileSize:     header.Size,
		UploadDate:   time.Now(),
	}

	// Store in database
	docID, err := s.store.CreateDocument(doc)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to store document in database: %w", err)
	}
	doc.ID = docID

	// Update expert_requests table with document reference
	if docType == "cv" {
		err = s.store.UpdateExpertRequestCVDocument(requestID, doc.ID)
	} else {
		err = s.store.UpdateExpertRequestApprovalDocument(requestID, doc.ID)
	}
	if err != nil {
		s.store.DeleteDocument(doc.ID)
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to update request reference: %w", err)
	}

	log.Debug("Document created and expert request updated successfully: ID %d", docID)
	return doc, nil
}

// CreateDocumentForExpertRequest maintains compatibility - delegates to CreateDocumentForRequest
func (s *Service) CreateDocumentForExpertRequest(requestID int64, file multipart.File, header *multipart.FileHeader) (*domain.Document, error) {
	return s.CreateDocumentForRequest(requestID, file, header, "cv")
}

// CreateApprovalDocumentForExpertRequest maintains compatibility - delegates to CreateDocumentForRequest
func (s *Service) CreateApprovalDocumentForExpertRequest(requestID int64, file multipart.File, header *multipart.FileHeader) (*domain.Document, error) {
	return s.CreateDocumentForRequest(requestID, file, header, "approval")
}

// MoveRequestDocumentToExpert moves a document from expert_requests to expert directories during approval
func (s *Service) MoveRequestDocumentToExpert(documentID, expertID int64) error {
	log := logger.Get()
	log.Debug("Moving document %d from request to expert %d", documentID, expertID)
	
	// Get the current document
	doc, err := s.store.GetDocument(documentID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}
	
	// Generate new path and filename
	timestamp := time.Now().Format("20060102_150405")
	extension := filepath.Ext(doc.FilePath)
	var newDir, newFilename string
	
	switch doc.DocumentType {
	case "cv":
		newDir = filepath.Join(s.uploadDir, "experts")
		newFilename = fmt.Sprintf("cv_%d_%s%s", expertID, timestamp, extension)
	case "approval":
		newDir = filepath.Join(s.uploadDir, "approvals")
		newFilename = fmt.Sprintf("approval_%d_%s%s", expertID, timestamp, extension)
	default:
		return fmt.Errorf("unsupported document type for migration: %s", doc.DocumentType)
	}
	
	// Create target directory
	if err := os.MkdirAll(newDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}
	
	newPath := filepath.Join(newDir, newFilename)
	
	// Move the file
	err = os.Rename(doc.FilePath, newPath)
	if err != nil {
		log.Error("Failed to move file from %s to %s: %v", doc.FilePath, newPath, err)
		return fmt.Errorf("failed to move file: %w", err)
	}
	
	// Update document record
	doc.FilePath = newPath
	doc.ExpertID = expertID // Update to use expert ID instead of request ID
	err = s.store.UpdateDocument(doc)
	if err != nil {
		// Try to move file back on database update failure
		os.Rename(newPath, doc.FilePath)
		return fmt.Errorf("failed to update document record: %w", err)
	}
	
	log.Debug("Successfully moved document %d to expert %d: %s", documentID, expertID, newPath)
	return nil
}

// CreateApprovalDocument creates an approval document for multiple expert IDs
func (s *Service) CreateApprovalDocument(expertIDs []int64, file multipart.File, header *multipart.FileHeader) (*domain.Document, error) {
	log := logger.Get()
	
	// Create approval filename with all expert IDs
	idStrs := make([]string, len(expertIDs))
	for i, id := range expertIDs {
		idStrs[i] = strconv.FormatInt(id, 10)
	}
	expertIDsStr := strings.Join(idStrs, "-")
	
	log.Debug("Creating approval document for expert IDs: %s, filename: %s", expertIDsStr, header.Filename)
	
	// Validate file size
	if header.Size > s.maxSize {
		log.Debug("Document service: File size %d exceeds maximum %d", header.Size, s.maxSize)
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxSize)
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	log.Debug("Document service: Validating content type: '%s'", contentType)
	if !s.allowedTypes[contentType] {
		log.Debug("Document service: Content type '%s' not allowed", contentType)
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}

	// Create approvals directory
	targetDir := filepath.Join(s.uploadDir, "approvals")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Debug("Document service: Failed to create directory: %v", err)
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate filename for approval
	timestamp := time.Now().Format("20060102_150405")
	extension := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("approval_%s_%s%s", expertIDsStr, timestamp, extension)
	filePath := filepath.Join(targetDir, filename)
	log.Debug("Document service: Full file path: %s", filePath)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		log.Debug("Document service: Failed to create file: %v", err)
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy the file data
	bytesWritten, err := io.Copy(dst, file)
	if err != nil {
		log.Debug("Document service: Failed to copy file data: %v", err)
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	log.Debug("Document service: File data copied successfully (%d bytes written)", bytesWritten)

	// Create document record - use first expert ID for database record
	doc := &domain.Document{
		ExpertID:     expertIDs[0], // Use first expert ID
		DocumentType: "approval",
		Filename:     header.Filename,
		FilePath:     filePath,
		ContentType:  contentType,
		FileSize:     header.Size,
		UploadDate:   time.Now(),
	}

	// Store in database
	docID, err := s.store.CreateDocument(doc)
	if err != nil {
		log.Debug("Document service: Failed to store in database: %v", err)
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to store document in database: %w", err)
	}

	doc.ID = docID
	log.Debug("Approval document created successfully - doc ID: %d, path: %s", docID, filePath)
	return doc, nil
}

// GetDocument retrieves a document by ID
func (s *Service) GetDocument(id int64) (*domain.Document, error) {
	return s.store.GetDocument(id)
}

// ListDocuments retrieves all documents for an expert
func (s *Service) ListDocuments(expertID int64) ([]*domain.Document, error) {
	return s.store.ListDocuments(expertID)
}

// DeleteDocument removes a document and its file
func (s *Service) DeleteDocument(id int64) error {
	// Get the document first to find the file path
	doc, err := s.store.GetDocument(id)
	if err != nil {
		return err
	}

	// Delete from database first
	if err := s.store.DeleteDocument(id); err != nil {
		return err
	}

	// Delete the file
	if err := os.Remove(doc.FilePath); err != nil {
		// Log but don't fail if file is missing
		// The database record is already deleted
		fmt.Printf("Warning: Could not delete file %s: %v\n", doc.FilePath, err)
	}

	return nil
}

// GetDocumentPath retrieves the file path for a document by ID
func (s *Service) GetDocumentPath(documentID int64) (string, error) {
	doc, err := s.store.GetDocument(documentID)
	if err != nil {
		return "", err
	}
	return doc.FilePath, nil
}

// GetExpertCVPath retrieves the CV file path for an expert
func (s *Service) GetExpertCVPath(expertID int64) (string, error) {
	expert, err := s.store.GetExpert(expertID)
	if err != nil {
		return "", err
	}
	
	if expert.CVDocumentID == nil {
		return "", fmt.Errorf("no CV document found for expert %d", expertID)
	}
	
	return s.GetDocumentPath(*expert.CVDocumentID)
}

// GetExpertApprovalPath retrieves the approval document path for an expert
func (s *Service) GetExpertApprovalPath(expertID int64) (string, error) {
	expert, err := s.store.GetExpert(expertID)
	if err != nil {
		return "", err
	}
	
	if expert.ApprovalDocumentID == nil {
		return "", fmt.Errorf("no approval document found for expert %d", expertID)
	}
	
	return s.GetDocumentPath(*expert.ApprovalDocumentID)
}

// MoveApprovalDocumentToExpert moves an approval document to use the correct expert ID in filename
func (s *Service) MoveApprovalDocumentToExpert(documentID, expertID int64) error {
	log := logger.Get()
	
	// Get the current document
	doc, err := s.store.GetDocument(documentID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}
	
	// Check if it's already properly named (in case it was already moved)
	currentPath := doc.FilePath
	
	// Generate the new filename with proper expert ID
	timestamp := time.Now().Format("20060102_150405")
	extension := filepath.Ext(currentPath)
	newFilename := fmt.Sprintf("approval_%d_%s%s", expertID, timestamp, extension)
	
	// Create target directory (approvals)
	targetDir := filepath.Join(s.uploadDir, "approvals")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create approvals directory: %w", err)
	}
	
	newPath := filepath.Join(targetDir, newFilename)
	
	// Move the file if paths are different
	if currentPath != newPath {
		err = os.Rename(currentPath, newPath)
		if err != nil {
			log.Error("Failed to move approval document from %s to %s: %v", currentPath, newPath, err)
			return fmt.Errorf("failed to move file: %w", err)
		}
		
		// Update the document record with new path
		doc.FilePath = newPath
		err = s.store.UpdateDocument(doc)
		if err != nil {
			// Try to move the file back on database update failure
			os.Rename(newPath, currentPath)
			return fmt.Errorf("failed to update document path in database: %w", err)
		}
		
		log.Info("Approval document renamed for expert %d: %s", expertID, newFilename)
	}
	
	return nil
}