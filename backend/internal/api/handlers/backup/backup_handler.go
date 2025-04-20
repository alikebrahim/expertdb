// Package backup provides handlers for backup operations
package backup

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"
)

// Handler handles backup endpoints
type Handler struct {
	store storage.Storage
}

// NewHandler creates a new backup handler
func NewHandler(store storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}

// HandleBackupCSV handles GET /api/backup request to create CSV backup
func (h *Handler) HandleBackupCSV(w http.ResponseWriter, r *http.Request) error {
	log := logger.Get()
	log.Debug("Processing GET /api/backup request")

	// Create a temporary directory to store CSV files
	tempDir, err := os.MkdirTemp("", "expertdb_backup_")
	if err != nil {
		log.Error("Failed to create temporary directory: %v", err)
		return fmt.Errorf("failed to create backup: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temporary directory regardless of outcome

	// Generate CSV files for each table
	timestamp := time.Now().Format("20060102_150405")
	zipFilename := fmt.Sprintf("expertdb_backup_%s.zip", timestamp)

	// Create the ZIP file
	zipPath := filepath.Join(tempDir, zipFilename)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Error("Failed to create ZIP file: %v", err)
		return fmt.Errorf("failed to create backup archive: %w", err)
	}
	defer zipFile.Close()

	// Create a new ZIP archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Generate and add experts CSV
	if err := h.addExpertsToZip(zipWriter, tempDir); err != nil {
		log.Error("Failed to add experts to backup: %v", err)
		return fmt.Errorf("failed to backup experts: %w", err)
	}

	// Generate and add expert requests CSV
	if err := h.addExpertRequestsToZip(zipWriter, tempDir); err != nil {
		log.Error("Failed to add expert requests to backup: %v", err)
		return fmt.Errorf("failed to backup expert requests: %w", err)
	}

	// Generate and add expert engagements CSV
	if err := h.addEngagementsToZip(zipWriter, tempDir); err != nil {
		log.Error("Failed to add engagements to backup: %v", err)
		return fmt.Errorf("failed to backup engagements: %w", err)
	}

	// Generate and add expert documents CSV
	if err := h.addDocumentsToZip(zipWriter, tempDir); err != nil {
		log.Error("Failed to add documents to backup: %v", err)
		return fmt.Errorf("failed to backup documents: %w", err)
	}

	// Generate and add expert areas CSV
	if err := h.addAreasToZip(zipWriter, tempDir); err != nil {
		log.Error("Failed to add areas to backup: %v", err)
		return fmt.Errorf("failed to backup areas: %w", err)
	}

	// Close the ZIP writer before reading the file
	zipWriter.Close()

	// Read the ZIP file
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		log.Error("Failed to read ZIP file: %v", err)
		return fmt.Errorf("failed to read backup archive: %w", err)
	}

	// Set response headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFilename))
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Length", strconv.Itoa(len(zipData)))
	w.WriteHeader(http.StatusOK)

	// Write the ZIP data to the response
	if _, err := w.Write(zipData); err != nil {
		log.Error("Failed to write ZIP data to response: %v", err)
		return fmt.Errorf("failed to send backup: %w", err)
	}

	log.Info("Backup successfully created and sent. Size: %d bytes", len(zipData))
	return nil
}

// addExpertsToZip creates a CSV file with experts data and adds it to the ZIP archive
func (h *Handler) addExpertsToZip(zipWriter *zip.Writer, tempDir string) error {
	// Get all experts
	experts, err := h.store.ListExperts(map[string]interface{}{}, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to retrieve experts: %w", err)
	}

	// Create CSV file
	csvPath := filepath.Join(tempDir, "experts.csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create experts CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{
		"ID", "ExpertID", "Name", "Designation", "Institution", "IsBahraini", 
		"Nationality", "IsAvailable", "Rating", "Role", "EmploymentType",
		"GeneralArea", "GeneralAreaName", "SpecializedArea", "IsTrained",
		"CVPath", "ApprovalDocumentPath", "Phone", "Email", "IsPublished", 
		"Biography", "CreatedAt", "UpdatedAt", "OriginalRequestID",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, expert := range experts {
		row := []string{
			strconv.FormatInt(expert.ID, 10),
			expert.ExpertID,
			expert.Name,
			expert.Designation,
			expert.Institution,
			strconv.FormatBool(expert.IsBahraini),
			expert.Nationality,
			strconv.FormatBool(expert.IsAvailable),
			expert.Rating,
			expert.Role,
			expert.EmploymentType,
			strconv.FormatInt(expert.GeneralArea, 10),
			expert.GeneralAreaName,
			expert.SpecializedArea,
			strconv.FormatBool(expert.IsTrained),
			expert.CVPath,
			expert.ApprovalDocumentPath,
			expert.Phone,
			expert.Email,
			strconv.FormatBool(expert.IsPublished),
			expert.Biography,
			expert.CreatedAt.Format(time.RFC3339),
			expert.UpdatedAt.Format(time.RFC3339),
			strconv.FormatInt(expert.OriginalRequestID, 10),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write expert row: %w", err)
		}
	}
	writer.Flush()

	// Add file to ZIP
	return addFileToZip(zipWriter, csvPath, "experts.csv")
}

// addExpertRequestsToZip creates a CSV file with expert requests data and adds it to the ZIP archive
func (h *Handler) addExpertRequestsToZip(zipWriter *zip.Writer, tempDir string) error {
	// Get all expert requests (both pending, approved, and rejected)
	requests, err := h.store.ListExpertRequests("all", 0, 0)
	if err != nil {
		return fmt.Errorf("failed to retrieve expert requests: %w", err)
	}

	// Create CSV file
	csvPath := filepath.Join(tempDir, "expert_requests.csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create expert requests CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{
		"ID", "ExpertID", "Name", "Designation", "Institution", "IsBahraini", 
		"IsAvailable", "Rating", "Role", "EmploymentType", "GeneralArea",
		"SpecializedArea", "IsTrained", "CVPath", "ApprovalDocumentPath",
		"Phone", "Email", "IsPublished", "Status", "RejectionReason",
		"Biography", "CreatedAt", "ReviewedAt", "ReviewedBy", "CreatedBy",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, req := range requests {
		// Format timestamps, handling zero times
		createdAt := ""
		if !req.CreatedAt.IsZero() {
			createdAt = req.CreatedAt.Format(time.RFC3339)
		}
		
		reviewedAt := ""
		if !req.ReviewedAt.IsZero() {
			reviewedAt = req.ReviewedAt.Format(time.RFC3339)
		}
		
		row := []string{
			strconv.FormatInt(req.ID, 10),
			req.ExpertID,
			req.Name,
			req.Designation,
			req.Institution,
			strconv.FormatBool(req.IsBahraini),
			strconv.FormatBool(req.IsAvailable),
			req.Rating,
			req.Role,
			req.EmploymentType,
			strconv.FormatInt(req.GeneralArea, 10),
			req.SpecializedArea,
			strconv.FormatBool(req.IsTrained),
			req.CVPath,
			req.ApprovalDocumentPath,
			req.Phone,
			req.Email,
			strconv.FormatBool(req.IsPublished),
			req.Status,
			req.RejectionReason,
			req.Biography,
			createdAt,
			reviewedAt,
			strconv.FormatInt(req.ReviewedBy, 10),
			strconv.FormatInt(req.CreatedBy, 10),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write expert request row: %w", err)
		}
	}
	writer.Flush()

	// Add file to ZIP
	return addFileToZip(zipWriter, csvPath, "expert_requests.csv")
}

// addEngagementsToZip creates a CSV file with engagements data and adds it to the ZIP archive
func (h *Handler) addEngagementsToZip(zipWriter *zip.Writer, tempDir string) error {
	// Get all engagements - since there's no global listing, we'll query a large range
	// of experts and combine their engagements
	experts, err := h.store.ListExperts(map[string]interface{}{}, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to retrieve experts: %w", err)
	}

	// Create CSV file
	csvPath := filepath.Join(tempDir, "engagements.csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create engagements CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{
		"ID", "ExpertID", "EngagementType", "StartDate", "EndDate",
		"ProjectName", "Status", "FeedbackScore", "Notes", "CreatedAt",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows for each expert's engagements
	allEngagements := make(map[int64]*domain.Engagement)
	for _, expert := range experts {
		engagements, err := h.store.ListEngagements(expert.ID, "", 1000, 0) // empty string for all types, high limit
		if err != nil {
			return fmt.Errorf("failed to retrieve engagements for expert %d: %w", expert.ID, err)
		}
		
		// Add to map to deduplicate
		for _, engagement := range engagements {
			allEngagements[engagement.ID] = engagement
		}
	}
	
	// Write all unique engagements
	for _, engagement := range allEngagements {
		// Format timestamps, handling zero times
		startDate := ""
		if !engagement.StartDate.IsZero() {
			startDate = engagement.StartDate.Format(time.RFC3339)
		}
		
		endDate := ""
		if !engagement.EndDate.IsZero() {
			endDate = engagement.EndDate.Format(time.RFC3339)
		}
		
		createdAt := ""
		if !engagement.CreatedAt.IsZero() {
			createdAt = engagement.CreatedAt.Format(time.RFC3339)
		}
		
		row := []string{
			strconv.FormatInt(engagement.ID, 10),
			strconv.FormatInt(engagement.ExpertID, 10),
			engagement.EngagementType,
			startDate,
			endDate,
			engagement.ProjectName,
			engagement.Status,
			strconv.Itoa(engagement.FeedbackScore),
			engagement.Notes,
			createdAt,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write engagement row: %w", err)
		}
	}
	writer.Flush()

	// Add file to ZIP
	return addFileToZip(zipWriter, csvPath, "engagements.csv")
}

// addDocumentsToZip creates a CSV file with documents data and adds it to the ZIP archive
func (h *Handler) addDocumentsToZip(zipWriter *zip.Writer, tempDir string) error {
	// Get all experts for their documents
	experts, err := h.store.ListExperts(map[string]interface{}{}, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to retrieve experts: %w", err)
	}

	// Create CSV file
	csvPath := filepath.Join(tempDir, "documents.csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create documents CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{
		"ID", "ExpertID", "DocumentType", "Filename", "FilePath",
		"ContentType", "FileSize", "UploadDate",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows for each expert's documents
	allDocuments := make(map[int64]*domain.Document)
	for _, expert := range experts {
		documents, err := h.store.ListDocuments(expert.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve documents for expert %d: %w", expert.ID, err)
		}
		
		// Add to map to deduplicate
		for _, document := range documents {
			allDocuments[document.ID] = document
		}
	}
	
	// Write all unique documents
	for _, document := range allDocuments {
		// Format timestamp, handling zero time
		uploadDate := ""
		if !document.UploadDate.IsZero() {
			uploadDate = document.UploadDate.Format(time.RFC3339)
		}
		
		row := []string{
			strconv.FormatInt(document.ID, 10),
			strconv.FormatInt(document.ExpertID, 10),
			document.DocumentType,
			document.Filename,
			document.FilePath,
			document.ContentType,
			strconv.FormatInt(document.FileSize, 10),
			uploadDate,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write document row: %w", err)
		}
	}
	writer.Flush()

	// Add file to ZIP
	return addFileToZip(zipWriter, csvPath, "documents.csv")
}

// addAreasToZip creates a CSV file with expert areas data and adds it to the ZIP archive
func (h *Handler) addAreasToZip(zipWriter *zip.Writer, tempDir string) error {
	// Get all areas
	areas, err := h.store.ListAreas()
	if err != nil {
		return fmt.Errorf("failed to retrieve expert areas: %w", err)
	}

	// Create CSV file
	csvPath := filepath.Join(tempDir, "expert_areas.csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create expert areas CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{"ID", "Name"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, area := range areas {
		row := []string{
			strconv.FormatInt(area.ID, 10),
			area.Name,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write area row: %w", err)
		}
	}
	writer.Flush()

	// Add file to ZIP
	return addFileToZip(zipWriter, csvPath, "expert_areas.csv")
}

// addFileToZip adds a file to a ZIP archive
func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	fileWriter, err := zipWriter.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file entry %s: %w", zipPath, err)
	}

	_, err = fileWriter.Write(fileData)
	if err != nil {
		return fmt.Errorf("failed to write data to ZIP entry %s: %w", zipPath, err)
	}

	return nil
}