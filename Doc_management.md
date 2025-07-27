# Document Management System Analysis & Fix Plan

## **Critical Assessment: Current vs Intended Architecture**

### **Database Schema Analysis**
The system has **TWO separate document storage approaches** that are inconsistent:

1. **Modern Architecture**: `expert_documents` table with proper foreign key relationships
2. **Legacy Architecture**: File paths stored directly in `experts` and `expert_requests` tables

### **CRITICAL DISCOVERY: Extensive Query Dependencies on Path Fields**

**Investigation Results:**
- **28+ SQL queries** across 3 major files directly reference `cv_path` and `approval_document_path` fields
- **Domain structs** (Expert, ExpertRequest) have path fields as core properties with JSON serialization
- **All CRUD operations** (CREATE, READ, UPDATE, SELECT) depend on path-based fields
- **API responses** expose path fields directly to clients
- **File operations** in approval workflows directly update path fields in database

**Path Field Dependencies:**
1. **expert.go**: 8 references across CREATE, SELECT, UPDATE operations
2. **expert_request.go**: 15+ references in complex approval workflows
3. **expert_edit_request.go**: 5+ references in edit request processing
4. **Domain types**: Expert and ExpertRequest structs expose paths as primary JSON fields

### **Current Document Flow Issues**

#### **Expert Request Creation:**
- ✅ CV documents stored in `expert_documents` with `expert_id = -requestID` 
- ✅ File stored in `data/documents/expert_requests/`
- ❌ ExpertRequest.CVPath also stores path (redundant)

#### **Expert Request Approval:**
- ❌ **CRITICAL**: Approval documents created with hardcoded `expert_id = -1`
- ❌ CV documents in `expert_documents` table NEVER updated with real expert ID
- ✅ Files moved correctly on filesystem
- ❌ Database remains inconsistent

#### **Single vs Batch Approval Discrepancy:**
- **Single approval**: Uses old path-based approach, copies paths to expert record
- **Batch approval**: Creates experts with empty paths, then updates via file move

## **Root Cause: Deep Legacy Integration**

The system has **hybrid architecture** with:
- **Path-based storage** in entity tables (legacy) - **DEEPLY INTEGRATED**
- **Proper document table** with foreign keys (modern) - **PARTIALLY IMPLEMENTED**

**CRITICAL INSIGHT: Schema Migration Complexity**
Changing from path fields to document ID references requires:
1. **28+ query modifications** across storage layer
2. **Domain struct changes** affecting JSON API contracts
3. **Handler updates** for file operations and API responses
4. **Client compatibility** - API consumers expect path fields
5. **Migration scripts** to transition existing data

## **Key Considerations for Fix**
- **File path transitions**: Handle `data/documents/expert_requests/` → `data/documents/cvs/` moves
- **Migration file editing**: Modify existing migration files, no new ones
- **Atomic operations**: Maintain transaction integrity for expert request approval

## **REVISED Strategy: Gradual Migration Approach**

**DECISION: Maintain Path Fields During Transition**
Given the extensive integration of path fields, we'll implement a **dual-storage approach** that maintains backward compatibility while gradually migrating to document-centric architecture.

### **Phase 1: Fix Critical Issues (Immediate)**

#### **1. Fix Approval Document Creation**
**Problem**: Approval documents created with hardcoded `expert_id = -1`
**Solution**: 
- During approval process, create approval documents with `expert_id = -requestID` (consistent with CV approach)
- After expert creation, update both CV and approval documents to use real expert ID
- **Keep existing path-based workflow intact**

#### **2. Complete Document Transition Logic**  
**Current Issue**: Only file moves happen, database `expert_documents` table not updated
**Solution**:
- Extend `moveCVFileToExpertDirectory()` to also update `expert_documents` table
- Create `moveApprovalDocumentToExpertDirectory()` for approval documents
- **Maintain path field updates for backward compatibility**
- Both functions handle: file move + database update + path field update atomically

#### **3. Database Schema Updates (Add Document References)**
**Approach**: Add document ID fields alongside existing path fields
**File**: `db/migrations/sqlite/0004_create_expert_table.sql`
- Add new fields: `cv_document_id`, `approval_document_id` (nullable)
- **Keep existing**: `cv_path`, `approval_document_path` fields
- Add comment explaining dual-storage approach

**File**: `db/migrations/sqlite/0002_create_expert-request_table.sql` 
- Add new fields: `cv_document_id`, `approval_document_id` (nullable)
- **Keep existing**: `cv_path`, `approval_document_path` fields
- Add comment explaining dual-storage approach

## **Implementation Details**

### **Enhanced Document Transition Service (Dual-Storage)**
```go
// Enhanced function signatures maintaining backward compatibility:
func (s *SQLiteStore) moveCVFileToExpertDirectory(oldPath string, expertID int64) error {
    // 1. Move file from expert_requests/ to cvs/
    // 2. Update expert_documents SET expert_id = expertID WHERE expert_id = -requestID
    // 3. Update expert_documents SET file_path = newPath WHERE expert_id = expertID
    // 4. Update experts.cv_path (MAINTAIN for backward compatibility)
    // 5. Update experts.cv_document_id = documentID (NEW for future migration)
}

func (s *SQLiteStore) moveApprovalDocumentToExpertDirectory(oldPath string, expertID int64) error {
    // 1. Move file from approvals/ to approvals/ (rename with expert ID)
    // 2. Update expert_documents SET expert_id = expertID WHERE expert_id = -1 AND file_path = oldPath
    // 3. Update expert_documents SET file_path = newPath WHERE expert_id = expertID AND document_type = 'approval'
    // 4. Update experts.approval_document_path (MAINTAIN for backward compatibility)
    // 5. Update experts.approval_document_id = documentID (NEW for future migration)
}
```

**Key Changes:**
- **Dual updates**: Both path fields AND document ID fields are maintained
- **Backward compatibility**: All existing queries continue to work
- **Forward compatibility**: New document ID fields prepare for future migration
- **Atomic operations**: All updates within single transaction

### **Atomic Transaction Flow**
```go
// In UpdateExpertRequestStatus:
func (s *SQLiteStore) UpdateExpertRequestStatus(...) error {
    tx, err := s.db.Begin()
    defer tx.Rollback()
    
    // 1. Update expert_requests status
    // 2. If approved: Create expert
    // 3. Move CV file + update documents table
    // 4. Move approval file + update documents table  
    // 5. Update expert with final paths
    
    tx.Commit()
}
```

## **Current State Analysis**

### **Database Tables**
- `expert_documents`: Modern approach with proper foreign keys
- `experts`: Has path fields (cv_path, approval_document_path) - legacy
- `expert_requests`: Has path fields (cv_path, approval_document_path) - legacy

### **Document Flow Problems**
1. **CV Documents**: Created correctly in `expert_documents` with `-requestID`, but never updated to real expert ID
2. **Approval Documents**: Created with hardcoded `-1` expert ID
3. **File Moves**: Work correctly but database not synchronized
4. **Path Storage**: Redundant storage in both tables and document records

## **REVISED Migration Strategy: Compatibility-First Approach**

### **Phase 1: Critical Fixes (Current Priority)**
1. Fix approval document creation to use proper expert ID tracking
2. Add document table update logic during expert creation  
3. Ensure both CV and approval documents get proper expert IDs
4. **Add document ID fields alongside existing path fields**

### **Phase 2: Dual-Storage Implementation**
1. Edit existing migration files to add document ID references (nullable)
2. Update domain models to include document IDs (optional fields)
3. **Maintain all existing path field functionality**
4. Implement document ID population in all file operations

### **Phase 3: Gradual API Enhancement**
1. Add optional document ID fields to API responses
2. Maintain existing path-based file operations
3. Add document-centric endpoints as alternatives
4. **No breaking changes to existing APIs**

### **Phase 4: Optional Long-term Migration**
1. Monitor usage of document-centric vs path-based approaches
2. Evaluate client adoption of new document ID fields
3. Plan deprecation timeline based on actual usage
4. **Only remove path fields after confirmed client migration**

**Benefits of This Approach:**
- **Zero disruption** to existing functionality
- **Gradual migration** path with full backward compatibility
- **Risk mitigation** - can rollback at any phase
- **Production stability** maintained throughout

## **Files to Modify**

### **Phase 1: Critical Fixes (Immediate)**
1. `internal/storage/sqlite/expert_request.go` - Enhanced move functions with dual-storage
2. `internal/api/handlers/expert_request.go` - Fix approval document creation logic
3. `internal/storage/sqlite/document.go` - Add document update methods

### **Phase 2: Schema Updates (Dual-Storage)**
1. `db/migrations/sqlite/0004_create_expert_table.sql` - Add document ID fields (nullable)
2. `db/migrations/sqlite/0002_create_expert-request_table.sql` - Add document ID fields (nullable)  
3. `internal/domain/types.go` - Add optional document ID fields to structs

### **Phase 3: Service Layer (Backward Compatible)**
1. `internal/documents/service.go` - Enhanced document management with dual-storage
2. `internal/storage/interface.go` - Add new method signatures (non-breaking)
3. API handlers - Add document ID support without breaking existing path-based operations

**CRITICAL CHANGE**: All modifications maintain existing functionality while adding new capabilities

## **Benefits of Complete Migration**
- **Eliminates data inconsistency** between file paths and document table
- **Proper foreign key relationships** with cascade deletes
- **Centralized document management** through document service
- **Audit trail** of all document operations
- **Scalable architecture** for multiple document types per expert

## **Testing Strategy**
- Test single expert approval with both CV and approval documents
- Test batch approval with mixed document scenarios
- Verify atomic rollback on failures
- Validate file system consistency with database state
- Test document retrieval through both legacy and modern approaches

## **Risk Mitigation**
- Maintain backward compatibility during transition
- Implement proper error handling and rollback mechanisms
- Ensure atomic operations for all document transitions
- Validate file system state matches database state
- Add comprehensive logging for document operations

## **Success Criteria**
1. **Immediate**: Approval documents have correct expert IDs
2. **Short-term**: All documents properly transitioned during expert creation
3. **Long-term**: Complete migration to document-centric architecture
4. **Always**: Atomic operations and data consistency maintained