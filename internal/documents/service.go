// Package documents provides document management functionality
package documents

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	
	"expertdb/internal/domain"
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
	// Validate document type
	validTypes := map[string]bool{
		"cv":       true,
		"approval": true,
		"certificate": true,
		"publication": true,
		"other":    true,
	}
	
	if !validTypes[docType] {
		return nil, fmt.Errorf("document type '%s' is not allowed; must be one of: cv, approval, certificate, publication, other", docType)
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

	// Create expert-specific directory
	expertDir := filepath.Join(s.uploadDir, fmt.Sprintf("expert_%d", expertID))
	if err := os.MkdirAll(expertDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create expert directory: %w", err)
	}

	// Generate a unique filename
	timestamp := time.Now().Format("20060102_150405")
	extension := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%s%s", expertID, timestamp, extension)
	filePath := filepath.Join(expertDir, filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy the file data
	if _, err = io.Copy(dst, file); err != nil {
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