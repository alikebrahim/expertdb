// Package documents provides document management functionality
package documents

import (
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

// CreateDocument handles file upload and database registration
func (s *Service) CreateDocument(expertID int64, file multipart.File, header *multipart.FileHeader, docType string) (*domain.Document, error) {
	log := logger.Get()
	log.Debug("Document service: CreateDocument called with expertID=%d, docType='%s', filename='%s'", 
		expertID, docType, header.Filename)
	
	// Validate document type
	validTypes := map[string]bool{
		"cv":       true,
		"approval": true,
	}
	
	if !validTypes[docType] {
		log.Debug("Document service: Invalid document type '%s'", docType)
		return nil, fmt.Errorf("document type '%s' is not allowed; must be one of: cv, approval", docType)
	}
	log.Debug("Document service: Document type '%s' is valid", docType)
	
	// Validate file size
	log.Debug("Document service: Validating file size: %d bytes (max: %d bytes)", header.Size, s.maxSize)
	if header.Size > s.maxSize {
		log.Debug("Document service: File size %d exceeds maximum %d", header.Size, s.maxSize)
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxSize)
	}
	log.Debug("Document service: File size validation passed")

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	log.Debug("Document service: Validating content type: '%s'", contentType)
	log.Debug("Document service: Allowed types: %v", s.allowedTypes)
	if !s.allowedTypes[contentType] {
		log.Debug("Document service: Content type '%s' not allowed", contentType)
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}
	log.Debug("Document service: Content type validation passed")

	// Determine directory and filename based on document type
	var targetDir, filename string
	timestamp := time.Now().Format("20060102_150405")
	extension := filepath.Ext(header.Filename)
	
	switch docType {
	case "cv":
		if expertID < 0 {
			// This is an expert request CV
			targetDir = filepath.Join(s.uploadDir, "expert_requests")
			filename = fmt.Sprintf("expert_request_%d_%s%s", -expertID, timestamp, extension)
		} else {
			// This is an approved expert CV
			targetDir = filepath.Join(s.uploadDir, "experts")
			filename = fmt.Sprintf("cv_%d_%s%s", expertID, timestamp, extension)
		}
	case "approval":
		// Approval documents go in approvals directory
		targetDir = filepath.Join(s.uploadDir, "approvals")
		// For approvals, expertID should contain the expert IDs (could be formatted string)
		filename = fmt.Sprintf("approval_%d_%s%s", expertID, timestamp, extension)
	default:
		// Other document types use expert-specific directories
		targetDir = filepath.Join(s.uploadDir, fmt.Sprintf("expert_%d", expertID))
		filename = fmt.Sprintf("%d_%s%s", expertID, timestamp, extension)
	}
	
	log.Debug("Document service: Target directory: %s", targetDir)
	log.Debug("Document service: Generated filename: %s", filename)
	
	// Create the target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Debug("Document service: Failed to create directory: %v", err)
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	log.Debug("Document service: Directory created successfully")

	filePath := filepath.Join(targetDir, filename)
	log.Debug("Document service: Full file path: %s", filePath)

	// Create the file
	log.Debug("Document service: Creating file at: %s", filePath)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Debug("Document service: Failed to create file: %v", err)
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()
	log.Debug("Document service: File created successfully")

	// Copy the file data
	log.Debug("Document service: Copying file data...")
	bytesWritten, err := io.Copy(dst, file)
	if err != nil {
		log.Debug("Document service: Failed to copy file data: %v", err)
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	log.Debug("Document service: File data copied successfully (%d bytes written)", bytesWritten)

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
	log.Debug("Document service: Created document record: %+v", doc)

	// Store in database
	log.Debug("Document service: Storing document in database...")
	docID, err := s.store.CreateDocument(doc)
	if err != nil {
		log.Debug("Document service: Failed to store in database: %v", err)
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to store document in database: %w", err)
	}
	log.Debug("Document service: Document stored in database with ID: %d", docID)

	doc.ID = docID
	log.Debug("Document service: CreateDocument completed successfully - returning doc ID: %d", docID)
	return doc, nil
}

// CreateDocumentForExpert creates a document and automatically updates expert references
func (s *Service) CreateDocumentForExpert(expertID int64, file multipart.File, header *multipart.FileHeader, docType string) (*domain.Document, error) {
	log := logger.Get()
	log.Debug("Creating document for expert %d, type: %s", expertID, docType)
	
	// Create document record
	doc, err := s.CreateDocument(expertID, file, header, docType)
	if err != nil {
		return nil, err
	}

	// Automatically update expert table with document reference
	switch docType {
	case "cv":
		err = s.store.UpdateExpertCVDocument(expertID, doc.ID)
	case "approval":
		err = s.store.UpdateExpertApprovalDocument(expertID, doc.ID)
	}

	if err != nil {
		// Rollback document creation
		s.store.DeleteDocument(doc.ID)
		os.Remove(doc.FilePath)
		return nil, fmt.Errorf("failed to update expert reference: %w", err)
	}

	log.Debug("Document created and expert updated successfully")
	return doc, nil
}

// CreateDocumentForExpertRequest handles CV upload for expert requests using request ID
func (s *Service) CreateDocumentForExpertRequest(requestID int64, file multipart.File, header *multipart.FileHeader) (*domain.Document, error) {
	log := logger.Get()
	log.Debug("Creating document for expert request %d", requestID)
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
	filename := fmt.Sprintf("expert_request_%d_%s%s", requestID, timestamp, extension)
	filePath := filepath.Join(targetDir, filename)

	// Create the file
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

	// Use negative request ID to distinguish from expert documents
	doc, err := s.CreateDocument(-requestID, file, header, "cv")
	if err != nil {
		return nil, err
	}

	// Update expert_requests table with document reference
	err = s.store.UpdateExpertRequestCVDocument(requestID, doc.ID)
	if err != nil {
		s.store.DeleteDocument(doc.ID)
		os.Remove(doc.FilePath)
		return nil, fmt.Errorf("failed to update request reference: %w", err)
	}

	log.Debug("Document created and expert request updated successfully")
	return doc, nil
}

// MoveExpertRequestCVToExpert moves CV file from expert_requests to experts directory and returns new path
func (s *Service) MoveExpertRequestCVToExpert(requestID, expertID int64) (string, error) {
	log := logger.Get()
	log.Debug("Moving CV from request %d to expert %d", requestID, expertID)
	
	// Find the request CV file
	requestDir := filepath.Join(s.uploadDir, "expert_requests")
	pattern := fmt.Sprintf("expert_request_%d_*.pdf", requestID)
	matches, err := filepath.Glob(filepath.Join(requestDir, pattern))
	if err != nil {
		log.Error("Failed to find request CV files: %v", err)
		return "", fmt.Errorf("failed to find request CV files: %w", err)
	}
	
	if len(matches) == 0 {
		log.Warn("No CV file found for request %d", requestID)
		return "", fmt.Errorf("no CV file found for request %d", requestID)
	}
	
	if len(matches) > 1 {
		log.Warn("Multiple CV files found for request %d, using first one", requestID)
	}
	
	oldPath := matches[0]
	log.Debug("Found request CV file: %s", oldPath)
	
	// Create target directory
	expertDir := filepath.Join(s.uploadDir, "experts")
	if err := os.MkdirAll(expertDir, 0755); err != nil {
		log.Error("Failed to create experts directory: %v", err)
		return "", fmt.Errorf("failed to create experts directory: %w", err)
	}
	
	// Generate new filename
	timestamp := time.Now().Format("20060102_150405")
	extension := filepath.Ext(oldPath)
	newFilename := fmt.Sprintf("cv_%d_%s%s", expertID, timestamp, extension)
	newPath := filepath.Join(expertDir, newFilename)
	
	log.Debug("Moving file from %s to %s", oldPath, newPath)
	
	// Move the file
	err = os.Rename(oldPath, newPath)
	if err != nil {
		log.Error("Failed to move file from %s to %s: %v", oldPath, newPath, err)
		return "", fmt.Errorf("failed to move file: %w", err)
	}
	
	log.Debug("File moved successfully to: %s", newPath)
	return newPath, nil
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