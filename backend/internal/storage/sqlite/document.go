package sqlite

import (
	"database/sql"
	"fmt"
	"time"
	
	"expertdb/internal/domain"
)

// ListDocuments retrieves all documents for an expert
func (s *SQLiteStore) ListDocuments(expertID int64) ([]*domain.Document, error) {
	query := `
		SELECT id, expert_id, document_type, filename, file_path,
				content_type, file_size, upload_date
		FROM expert_documents
		WHERE expert_id = ?
	`
	
	rows, err := s.db.Query(query, expertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expert documents: %w", err)
	}
	defer rows.Close()
	
	var docs []*domain.Document
	for rows.Next() {
		var doc domain.Document
		err := rows.Scan(
			&doc.ID, &doc.ExpertID, &doc.DocumentType, &doc.Filename,
			&doc.FilePath, &doc.ContentType, &doc.FileSize, &doc.UploadDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document row: %w", err)
		}
		
		// Set Type field as alias of DocumentType for API compatibility
		doc.Type = doc.DocumentType
		
		docs = append(docs, &doc)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating document rows: %w", err)
	}
	
	return docs, nil
}

// GetDocument retrieves a document by ID
func (s *SQLiteStore) GetDocument(id int64) (*domain.Document, error) {
	query := `
		SELECT id, expert_id, document_type, filename, file_path,
				content_type, file_size, upload_date
		FROM expert_documents
		WHERE id = ?
	`
	
	var doc domain.Document
	err := s.db.QueryRow(query, id).Scan(
		&doc.ID, &doc.ExpertID, &doc.DocumentType, &doc.Filename,
		&doc.FilePath, &doc.ContentType, &doc.FileSize, &doc.UploadDate,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	
	// Set Type field as alias of DocumentType for API compatibility
	doc.Type = doc.DocumentType
	
	return &doc, nil
}

// CreateDocument creates a new document in the database
func (s *SQLiteStore) CreateDocument(doc *domain.Document) (int64, error) {
	query := `
		INSERT INTO expert_documents (
			expert_id, document_type, filename, file_path,
			content_type, file_size, upload_date
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	// Handle potentially nullable fields
	var contentType interface{} = nil
	if doc.ContentType != "" {
		contentType = doc.ContentType
	}
	
	// Set default upload date if not provided
	if doc.UploadDate.IsZero() {
		doc.UploadDate = time.Now()
	}
	
	result, err := s.db.Exec(
		query,
		doc.ExpertID, doc.DocumentType, doc.Filename, doc.FilePath,
		contentType, doc.FileSize, doc.UploadDate,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create document: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get document ID: %w", err)
	}
	
	doc.ID = id
	return id, nil
}

// DeleteDocument deletes a document by ID
func (s *SQLiteStore) DeleteDocument(id int64) error {
	result, err := s.db.Exec("DELETE FROM expert_documents WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	
	return nil
}