# Expert Request Data Issue Investigation Report

## Executive Summary

The Expert Request Management system displays incomplete data for approved requests while pending requests show correctly. This investigation identified a critical SQL schema mismatch bug in the batch approval process that causes data corruption for approved expert requests.

## Issue Description

### Problem Statement
- **Approved requests** return empty strings for critical fields: `name`, `designation`, `affiliation`, `phone`, `email`
- **Pending requests** return complete data with all fields properly populated
- Frontend correctly displays whatever data the API returns
- All UI tabs use the same API endpoint (`/api/expert-requests`) with different query parameters

### API Endpoint Behavior
**Single Endpoint Used:** `GET /api/expert-requests`

**Query Parameters by Tab:**
1. **All Requests Tab:** `?limit=100&offset=0` (no status filter)
2. **Awaiting Action Tab:** Combines results from:
   - `?status=pending&limit=100&offset=0`
   - `?status=rejected&limit=100&offset=0`
3. **Processed Tab:** `?status=approved&limit=100&offset=0`

### Data Comparison

#### Working Response (Pending Request)
```json
{
  "id": 4,
  "name": "TestResponse",              ✓ Has data
  "designation": "Prof.",              ✓ Has data
  "affiliation": "TestResponse",       ✓ Has data
  "phone": "+973 12345678",            ✓ Has data
  "email": "test@response.com",        ✓ Has data
  "generalArea": 5,                    ✓ Actual area ID
  "specializedArea": "1,2,3",          ✓ Has data
  "status": "pending",
  "createdAt": "2025-07-28T11:30:00.123456789+03:00"
}
```

#### Broken Response (Approved Requests)
```json
{
  "id": 3,
  "name": "",                          ✗ Empty string
  "designation": "",                   ✗ Empty string
  "affiliation": "",                   ✗ Empty string
  "phone": "",                         ✗ Empty string
  "email": "",                         ✗ Empty string
  "generalArea": 0,                    ✗ Default value
  "specializedArea": "",               ✗ Empty string
  "status": "approved",
  "createdAt": "2025-07-27T14:32:06.211080523+03:00",
  "reviewedAt": "0001-01-01T00:00:00Z"  ✗ Not updated
}
```

## Root Cause Analysis

### Primary Issue Identified
**Location:** `/home/alikebrahim/dev/expertdb_backend/internal/storage/sqlite/expert_request.go:742`

**Problem:** SQL schema mismatch in the `BatchApproveExpertRequests` method

### Technical Details

#### 1. SQL Query Error
The problematic query attempts to SELECT a `rating` field from the `expert_requests` table:

```sql
SELECT id, name, designation, affiliation, is_bahraini, 
       is_available, rating, role, employment_type, general_area, 
       specialized_area, is_trained, cv_document_id, phone, email, 
       is_published, status, created_by
FROM expert_requests
WHERE id = ?
```

#### 2. Schema Mismatch
- **expert_requests table:** Does NOT have a `rating` column
- **experts table:** DOES have a `rating` column
- **SQL Error Generated:** `"no such column: rating"`

#### 3. Database Schema Verification

**expert_requests table structure:**
```sql
CREATE TABLE IF NOT EXISTS "expert_requests" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    designation TEXT,
    affiliation TEXT,
    -- ... other fields ...
    -- NOTE: NO rating field
    status TEXT DEFAULT 'pending'
);
```

**experts table structure:**
```sql
CREATE TABLE IF NOT EXISTS "experts" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    designation TEXT,
    affiliation TEXT,
    rating INTEGER DEFAULT 0 CHECK (rating >= 0 AND rating <= 5),
    -- ... other fields ...
);
```

#### 4. Code Analysis
**File:** `internal/storage/sqlite/expert_request.go`
**Lines:** 742-755

```go
// PROBLEMATIC QUERY (line 742)
query := `
    SELECT 
        id, name, designation, affiliation, is_bahraini, 
        is_available, rating, role, employment_type, general_area,  // ← rating doesn't exist
        specialized_area, is_trained, cv_document_id, phone, email, 
        is_published, status, created_by
    FROM expert_requests
    WHERE id = ?
`

// SCAN OPERATION MISSING rating VARIABLE (line 749-755)
err = tx.QueryRow(query, id).Scan(
    &req.ID, &req.Name, &req.Designation, &req.Affiliation, 
    &req.IsBahraini, &req.IsAvailable, &req.Role,  // ← Missing rating variable
    &req.EmploymentType, &req.GeneralArea, &req.SpecializedArea, 
    &req.IsTrained, &cvDocumentID, &req.Phone, &req.Email, 
    &req.IsPublished, &req.Status, &req.CreatedBy,
)
```

### Impact Assessment

#### Database Evidence
Query results showing the corruption:
```bash
sqlite> SELECT id, name, designation, affiliation, phone, email, status FROM expert_requests;
1||||||approved    ← Empty core fields
2||||||approved    ← Empty core fields  
3||||||approved    ← Empty core fields
4|TestResponse|Prof.|TestResponse|+97311111111|TestResponse@db.com|pending  ← Complete data
```

#### Workflow Analysis
1. **Pending requests** work correctly because they use different code paths
2. **Approved requests** fail during the batch approval process
3. **Data loss occurs** when the SQL error prevents proper data retrieval
4. **Frontend displays empty data** because the API returns default/empty values

## Solution Requirements

### Immediate Fix Required

**File:** `/home/alikebrahim/dev/expertdb_backend/internal/storage/sqlite/expert_request.go`

**Action:** Remove `rating` from the SELECT query in `BatchApproveExpertRequests` method

#### Before (Broken):
```sql
SELECT id, name, designation, affiliation, is_bahraini, 
       is_available, rating, role, employment_type, general_area,
       -- ... rest of fields
```

#### After (Fixed):
```sql
SELECT id, name, designation, affiliation, is_bahraini, 
       is_available, role, employment_type, general_area,
       -- ... rest of fields
```

### Additional Considerations

1. **Default Rating Value:** When creating expert records from approved requests, set `rating` to default value (0)
2. **Data Recovery:** The existing corrupted records (IDs 1, 2, 3) may need manual data recovery if the original request data is available elsewhere
3. **Testing:** Verify the approval workflow works end-to-end after the fix

## Verification Steps

### Pre-Fix Verification
1. ✅ Confirmed approved requests have empty core fields in database
2. ✅ Confirmed pending requests have complete data
3. ✅ Identified exact SQL error in code
4. ✅ Verified schema mismatch between tables

### Post-Fix Verification Required
1. ⏳ Test batch approval process with valid request
2. ⏳ Verify expert record creation with proper data transfer
3. ⏳ Confirm API returns complete data for all request statuses
4. ⏳ Validate frontend displays correct information

## Technical Context

- **Backend:** Go with SQLite database
- **API Pattern:** RESTful with status-based filtering  
- **Authentication:** JWT-based (working correctly)
- **Database:** SQLite with foreign key relationships
- **Frontend:** React + TypeScript (functioning correctly)

## Priority Classification

**Severity:** High  
**Impact:** Critical data loss for approved requests  
**Effort:** Low (single line SQL query fix)  
**Risk:** Low (well-isolated change)

---

**Report Generated:** 2025-07-28  
**Investigation by:** Claude Code  
**Status:** Root cause identified, solution ready for implementation