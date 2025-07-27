# AGGRESSIVE DOCUMENT MANAGEMENT OVERHAUL PLAN

**Date**: January 25, 2025  
**Context**: Complete architectural restructure for document management  
**Approach**: NO BACKWARD COMPATIBILITY - Clean Architecture Implementation  
**Strategy**: Edit existing files, modify migration files directly

---

## 🎯 STRATEGIC APPROACH: COMPLETE ARCHITECTURAL OVERHAUL

**DECISION**: Eliminate hybrid architecture completely. Move to pure document-centric approach with foreign key relationships. All path fields will be removed and replaced with document references.

---

## 📋 IMPLEMENTATION PHASES

### **PHASE 1: DATABASE SCHEMA OVERHAUL**

#### **1.1 Edit Migration Files Directly**

**File: `db/migrations/sqlite/0002_create_expert-request_table.sql`**
```sql
-- REMOVE path fields, ADD document references
-- BEFORE (lines 14-15):
cv_path TEXT,
approval_document_path TEXT,

-- AFTER (replace with):
cv_document_id INTEGER REFERENCES expert_documents(id),
approval_document_id INTEGER REFERENCES expert_documents(id),
```

**File: `db/migrations/sqlite/0004_create_expert_table.sql`**
```sql
-- REMOVE path fields, ADD document references
-- BEFORE (lines 15-16):
cv_path TEXT,           -- Path to the CV file NOTE: This is better be replaced with expert_documents(id)
approval_document_path TEXT, -- Path to the approval document

-- AFTER (replace with):
cv_document_id INTEGER REFERENCES expert_documents(id),
approval_document_id INTEGER REFERENCES expert_documents(id),
```

**File: `db/migrations/sqlite/0012_create_expert_edit_requests_table.sql`**
```sql
-- REMOVE path fields, ADD document references
-- BEFORE (lines 26-27):
new_cv_path TEXT,
new_approval_document_path TEXT,

-- AFTER (replace with):
new_cv_document_id INTEGER REFERENCES expert_documents(id),
new_approval_document_id INTEGER REFERENCES expert_documents(id),
```

#### **1.2 Add Proper Indexes**
```sql
-- Add to migration files
CREATE INDEX idx_experts_cv_document ON experts(cv_document_id);
CREATE INDEX idx_experts_approval_document ON experts(approval_document_id);
CREATE INDEX idx_expert_requests_cv_document ON expert_requests(cv_document_id);
CREATE INDEX idx_expert_requests_approval_document ON expert_requests(approval_document_id);
```

### **PHASE 2: DOMAIN MODEL OVERHAUL**

#### **2.1 Edit `internal/domain/types.go`**

**Expert struct changes:**
```go
type Expert struct {
    ID                  int64     `json:"id"`
    Name                string    `json:"name"`
    // ... other fields ...
    
    // REMOVE these fields completely:
    // CVPath              string    `json:"cvPath"`
    // ApprovalDocumentPath string    `json:"approvalDocumentPath,omitempty"`
    
    // REPLACE with document references:
    CVDocumentID         *int64    `json:"cvDocumentId,omitempty"`
    ApprovalDocumentID   *int64    `json:"approvalDocumentId,omitempty"`
    
    // Enhanced document access:
    CVDocument          *Document `json:"cvDocument,omitempty"`
    ApprovalDocument    *Document `json:"approvalDocument,omitempty"`
    Documents           []Document `json:"documents,omitempty"`
    
    // ... rest unchanged
}
```

**ExpertRequest struct changes:**
```go
type ExpertRequest struct {
    ID                   int64     `json:"id"`
    // ... other fields ...
    
    // REMOVE these fields:
    // CVPath               string    `json:"cvPath"`
    // ApprovalDocumentPath string    `json:"approvalDocumentPath,omitempty"`
    
    // REPLACE with document references:
    CVDocumentID         *int64    `json:"cvDocumentId,omitempty"`
    ApprovalDocumentID   *int64    `json:"approvalDocumentId,omitempty"`
    
    // Enhanced document access:
    CVDocument          *Document `json:"cvDocument,omitempty"`
    ApprovalDocument    *Document `json:"approvalDocument,omitempty"`
    
    // ... rest unchanged
}
```

#### **2.2 Add Document Resolution Methods**
```go
// Add to types.go
func (e *Expert) ResolveCVDocument(store Storage) error {
    if e.CVDocumentID != nil {
        doc, err := store.GetDocument(*e.CVDocumentID)
        if err == nil {
            e.CVDocument = doc
        }
    }
    return nil
}

func (e *Expert) ResolveApprovalDocument(store Storage) error {
    if e.ApprovalDocumentID != nil {
        doc, err := store.GetDocument(*e.ApprovalDocumentID)
        if err == nil {
            e.ApprovalDocument = doc
        }
    }
    return nil
}
```

### **PHASE 3: STORAGE LAYER COMPLETE REWRITE**

#### **3.1 Edit `internal/storage/sqlite/expert.go`**

**CreateExpert method:**
```go
func (s *SQLiteStore) CreateExpert(expert *domain.Expert) (int64, error) {
    tx, err := s.db.Begin()
    if err != nil {
        return 0, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    query := `
        INSERT INTO experts (
            name, designation, affiliation, is_bahraini, is_available, rating,
            role, employment_type, general_area, specialized_area, is_trained,
            cv_document_id, approval_document_id, phone, email, is_published, 
            created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

    if expert.CreatedAt.IsZero() {
        expert.CreatedAt = time.Now()
        expert.UpdatedAt = expert.CreatedAt
    }

    result, err := tx.Exec(
        query,
        expert.Name, expert.Designation, expert.Affiliation,
        expert.IsBahraini, expert.IsAvailable, expert.Rating,
        expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
        expert.IsTrained, expert.CVDocumentID, expert.ApprovalDocumentID, 
        expert.Phone, expert.Email, expert.IsPublished,
        expert.CreatedAt, expert.UpdatedAt,
    )
    // ... rest of transaction logic
}
```

**GetExpert method with document resolution:**
```go
func (s *SQLiteStore) GetExpert(id int64) (*domain.Expert, error) {
    query := `
        SELECT e.id, e.name, e.designation, e.affiliation, 
               e.is_bahraini, e.is_available, e.rating, e.role, 
               e.employment_type, e.general_area, ea.name as general_area_name, 
               e.specialized_area, e.is_trained, e.cv_document_id, e.approval_document_id, 
               e.phone, e.email, e.is_published, e.created_at, e.updated_at,
               COALESCE(
                   (SELECT GROUP_CONCAT(sa.name, ', ')
                   FROM specialized_areas sa
                   WHERE ',' || e.specialized_area || ',' LIKE '%,' || sa.id || ',%'
                   AND e.specialized_area IS NOT NULL 
                   AND e.specialized_area != ''),
                   ''
               ) as specialized_area_names
        FROM experts e
        LEFT JOIN expert_areas ea ON e.general_area = ea.id
        WHERE e.id = ?
    `

    var expert domain.Expert
    var generalAreaName sql.NullString
    var specializedAreaNames sql.NullString
    var cvDocumentID sql.NullInt64
    var approvalDocumentID sql.NullInt64
    var createdAt sql.NullTime
    var updatedAt sql.NullTime

    err := s.db.QueryRow(query, id).Scan(
        &expert.ID, &expert.Name, &expert.Designation, &expert.Affiliation,
        &expert.IsBahraini, &expert.IsAvailable, &expert.Rating, &expert.Role,
        &expert.EmploymentType, &expert.GeneralArea, &generalAreaName,
        &expert.SpecializedArea, &expert.IsTrained, &cvDocumentID, &approvalDocumentID, 
        &expert.Phone, &expert.Email, &expert.IsPublished, &createdAt, &updatedAt,
        &specializedAreaNames,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("failed to get expert: %w", err)
    }

    // Set document IDs
    if cvDocumentID.Valid {
        expert.CVDocumentID = &cvDocumentID.Int64
    }
    if approvalDocumentID.Valid {
        expert.ApprovalDocumentID = &approvalDocumentID.Int64
    }

    // Resolve document objects
    expert.ResolveCVDocument(s)
    expert.ResolveApprovalDocument(s)

    // ... rest of method (bio data, documents, engagements)
    return &expert, nil
}
```

#### **3.2 Edit `internal/storage/sqlite/expert_request.go`**

**Complete rewrite of document handling methods:**
```go
// REMOVE moveCVFileToExpertDirectory - replace with document-centric approach
func (s *SQLiteStore) TransferExpertRequestToExpert(requestID, expertID int64) error {
    tx, err := s.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    // 1. Update expert_documents: expert_id = -requestID → expertID
    _, err = tx.Exec(`
        UPDATE expert_documents 
        SET expert_id = ? 
        WHERE expert_id = ?
    `, expertID, -requestID)
    if err != nil {
        return fmt.Errorf("failed to update document expert_id: %w", err)
    }

    // 2. Get document IDs for the transferred documents
    var cvDocID, approvalDocID sql.NullInt64
    err = tx.QueryRow(`
        SELECT 
            (SELECT id FROM expert_documents WHERE expert_id = ? AND document_type = 'cv' LIMIT 1),
            (SELECT id FROM expert_documents WHERE expert_id = ? AND document_type = 'approval' LIMIT 1)
    `, expertID, expertID).Scan(&cvDocID, &approvalDocID)
    
    // 3. Update expert record with document references
    _, err = tx.Exec(`
        UPDATE experts 
        SET cv_document_id = ?, approval_document_id = ?
        WHERE id = ?
    `, cvDocID, approvalDocID, expertID)
    if err != nil {
        return fmt.Errorf("failed to update expert document references: %w", err)
    }

    return tx.Commit()
}
```

#### **3.3 Edit `internal/storage/sqlite/expert_edit_request.go`**

**Update all methods to use document references instead of paths:**
```go
// Replace all cv_path/approval_document_path references with document IDs
func (s *SQLiteStore) CreateExpertEditRequest(req *domain.ExpertEditRequest) (int64, error) {
    query := `
        INSERT INTO expert_edit_requests (
            expert_id, name, designation, institution, phone, email,
            is_bahraini, is_available, rating, role, employment_type,
            general_area, specialized_area, is_trained, is_published, rating,
            suggested_specialized_areas, new_cv_document_id, new_approval_document_id,
            remove_cv, remove_approval_document, change_summary, change_reason,
            fields_changed, status, created_at, created_by
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    // ... implementation with document ID fields
}
```

### **PHASE 4: DOCUMENT SERVICE OVERHAUL**

#### **4.1 Edit `internal/documents/service.go`**

**Enhanced CreateDocument with automatic entity updates:**
```go
func (s *Service) CreateDocumentForExpert(expertID int64, file multipart.File, header *multipart.FileHeader, docType string) (*domain.Document, error) {
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

    return doc, nil
}

func (s *Service) CreateDocumentForExpertRequest(requestID int64, file multipart.File, header *multipart.FileHeader) (*domain.Document, error) {
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

    return doc, nil
}
```

**Document path resolution methods:**
```go
func (s *Service) GetDocumentPath(documentID int64) (string, error) {
    doc, err := s.store.GetDocument(documentID)
    if err != nil {
        return "", err
    }
    return doc.FilePath, nil
}

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
```

### **PHASE 5: API HANDLER OVERHAUL**

#### **5.1 Edit `internal/api/handlers/expert.go`**

**Remove all path-based file handling, use document service:**
```go
func (h *Handler) HandleCreateExpert(w http.ResponseWriter, r *http.Request) error {
    // Parse expert data
    expert, err := h.parseExpertFromRequest(r)
    if err != nil {
        return err
    }

    // Handle CV upload if present
    if cvFile, cvHeader, err := r.FormFile("cv"); err == nil {
        defer cvFile.Close()
        
        // Create expert first to get ID
        expertID, err := h.store.CreateExpert(expert)
        if err != nil {
            return err
        }
        
        // Create CV document
        doc, err := h.documentService.CreateDocumentForExpert(expertID, cvFile, cvHeader, "cv")
        if err != nil {
            return err
        }
        
        expert.CVDocumentID = &doc.ID
        expert.CVDocument = doc
    }

    // Similar handling for approval document...
    
    return response.Success(w, http.StatusCreated, "Expert created successfully", expert)
}
```

#### **5.2 Edit `internal/api/handlers/expert_request.go`**

**Replace all path-based operations with document references:**
```go
func (h *Handler) HandleCreateExpertRequest(w http.ResponseWriter, r *http.Request) error {
    // Parse request data
    expertReq, err := h.parseExpertRequestFromForm(r)
    if err != nil {
        return err
    }

    // Create expert request first
    requestID, err := h.store.CreateExpertRequest(expertReq)
    if err != nil {
        return err
    }

    // Handle CV upload
    if cvFile, cvHeader, err := r.FormFile("cv"); err == nil {
        defer cvFile.Close()
        
        doc, err := h.documentService.CreateDocumentForExpertRequest(requestID, cvFile, cvHeader)
        if err != nil {
            return err
        }
        
        expertReq.CVDocumentID = &doc.ID
        expertReq.CVDocument = doc
    }

    return response.Success(w, http.StatusCreated, "Expert request created successfully", expertReq)
}

func (h *Handler) HandleUpdateExpertRequestStatus(w http.ResponseWriter, r *http.Request) error {
    // ... existing validation logic ...

    if status == "approved" {
        // Create expert from request
        expert := h.convertRequestToExpert(expertReq)
        expertID, err := h.store.CreateExpert(expert)
        if err != nil {
            return err
        }

        // Transfer documents from request to expert
        err = h.store.TransferExpertRequestToExpert(requestID, expertID)
        if err != nil {
            return err
        }
    }

    return response.Success(w, http.StatusOK, "Expert request updated successfully", expertReq)
}
```

#### **5.3 Create New Storage Interface Methods**

**Edit `internal/storage/interface.go`:**
```go
type Storage interface {
    // Existing methods...
    
    // New document reference methods
    UpdateExpertCVDocument(expertID, documentID int64) error
    UpdateExpertApprovalDocument(expertID, documentID int64) error
    UpdateExpertRequestCVDocument(requestID, documentID int64) error
    UpdateExpertRequestApprovalDocument(requestID, documentID int64) error
    
    TransferExpertRequestToExpert(requestID, expertID int64) error
    
    // Document resolution methods
    GetExpertWithDocuments(expertID int64) (*domain.Expert, error)
    GetExpertRequestWithDocuments(requestID int64) (*domain.ExpertRequest, error)
}
```

### **PHASE 6: FILE SYSTEM ORGANIZATION**

#### **6.1 Simplify Directory Structure**
```
data/documents/
├── experts/           # All expert documents by ID
│   ├── cv_<expertId>_<timestamp>.pdf
│   └── approval_<expertId>_<timestamp>.pdf
└── requests/          # Temporary request documents
    └── cv_request_<requestId>_<timestamp>.pdf
```

#### **6.2 Update Document Service File Handling**
```go
func (s *Service) generateFilePath(expertID int64, docType, extension string) string {
    timestamp := time.Now().Format("20060102_150405")
    
    if expertID < 0 {
        // Request document
        return filepath.Join(s.uploadDir, "requests", 
            fmt.Sprintf("cv_request_%d_%s%s", -expertID, timestamp, extension))
    } else {
        // Expert document
        return filepath.Join(s.uploadDir, "experts", 
            fmt.Sprintf("%s_%d_%s%s", docType, expertID, timestamp, extension))
    }
}
```

---

## 🔧 IMPLEMENTATION ORDER

### **Week 1: Database & Domain**
1. Edit migration files to replace path fields with document references
2. Update domain structs to use document IDs
3. Add document resolution methods

### **Week 2: Storage Layer**
1. Rewrite all storage methods to use document references
2. Remove all path-based operations
3. Implement document transfer logic

### **Week 3: Services & APIs**
1. Update document service for automatic entity updates
2. Rewrite API handlers to use document-centric approach
3. Remove all path-based file handling

### **Week 4: Testing & Cleanup**
1. Update all tests to use new document-centric approach
2. Clean up unused path-based methods
3. Verify data consistency and file system organization

---

## 💡 ARCHITECTURAL BENEFITS

1. **Clean Architecture**: Single source of truth for document management
2. **Proper Foreign Keys**: Automatic cascade deletes and referential integrity
3. **Simplified File Operations**: Centralized through document service
4. **Better Performance**: Proper indexing on document references
5. **Audit Trail**: Complete document lifecycle tracking
6. **Scalability**: Easy to add new document types and operations

---

## 🚨 BREAKING CHANGES

**API Response Changes:**
- `cvPath` field removed from all responses
- `approvalDocumentPath` field removed from all responses
- New fields: `cvDocumentId`, `approvalDocumentId`, `cvDocument`, `approvalDocument`

**Database Schema Changes:**
- All path columns removed from tables
- Document reference columns added
- Foreign key constraints added

**File System Changes:**
- Simplified directory structure
- Consistent naming conventions
- Document-centric file operations

---

**CONCLUSION**: This aggressive overhaul completely eliminates the hybrid architecture and establishes a clean, document-centric approach. All path-based operations are removed and replaced with proper foreign key relationships to the document table.