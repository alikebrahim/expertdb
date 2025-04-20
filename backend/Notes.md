- TASK: Analyze current db to determine how to create conversion/ import mechanism to sql in order to :
    - assign general aread and specialised areas based on research being done (grok convo https://grok.com/chat/2a8a04c0-6ddb-4ed0-a630-bf9cbb7121da)
    - Ex. General area: Business, Special Areas: Accounting and Auditing, Banking and Finance...etc
    - Add field for skills to replace spesialized area e.g. Education Quality Assurance, System Programming...etc. Use current specialised areas as a reference
- TASK: Sketch a flow chart for both Expert_Creation and Phase_Planning workflows and generate a mermaid.js to:
    - aid the SRS
    - help improve script to simulate workflow

# Implementation Progress Notes

## 2025-04-21: Phase 7 Implementation (Statistics and Reporting Enhancements)

### Completed Work:

#### Phase 7A: Published Expert Statistics ✅
- Added `PublishedCount` and `PublishedRatio` fields to the Statistics struct
- Implemented `GetPublishedExpertStats()` method to calculate published expert statistics
- Updated the main statistics endpoint to include published expert data
- Added test cases to verify published statistics

#### Phase 7B: Growth Statistics Enhancement ✅
- Converted monthly growth statistics to yearly statistics
- Created new `GetExpertGrowthByYear(years int)` method in the repository
- Updated growth statistics endpoint to use years parameter instead of months
- Improved growth calculation algorithm with better year formatting
- Added test cases for yearly growth in test_api.sh

#### Phase 7C: Engagement Type Statistics ✅
- Updated engagement statistics query to filter by "validator" and "evaluator" types only
- Enhanced statistics handler to limit results to these specific engagement types
- Added test cases to verify engagement type restrictions

#### Phase 7D: Area Statistics Implementation ✅
- Created new `/api/statistics/areas` endpoint
- Implemented `GetAreaStatistics()` method to calculate:
  - General area statistics (all areas)
  - Top 5 specialized areas by expert count
  - Bottom 5 specialized areas by expert count
- Added test cases to verify area statistics
- Ensured proper permission controls (super_user access)

#### Files Modified:
- `/internal/domain/types.go` - Added new statistics fields for published experts and yearly growth
- `/internal/storage/interface.go` - Added new methods for the statistics repository
- `/internal/storage/sqlite/statistics.go` - Implemented all the new statistics methods
- `/internal/api/handlers/statistics/statistics_handler.go` - Added area statistics handler and updated growth statistics
- `/internal/api/server.go` - Added new statistics endpoint for area statistics
- `/test_api.sh` - Enhanced tests for all statistics endpoints

## 2025-04-20: Phase 3B, 3C, 4A, 4B, 4C, 5A, 5B, 5C, and 6 Implementation

### Completed Work:

#### Phase 3B: Sorting and Pagination Improvements ✅
- Enhanced sorting options for the expert listing endpoint
- Implemented safe column name validation for SQL injection prevention
- Improved pagination responses with metadata
- Added pagination headers for client-side UI improvements

#### Phase 3C: Expert Detail Access ✅
- Added approval_document_path to Expert struct
- Updated database schema in migrations
- Updated expert creation, retrieval, and update methods to include approval document
- Ensured endpoints for expert details remain accessible to all authenticated users

#### Phase 4A: Expert Request Creation with CV Upload ✅
- Made expert request creation endpoint accessible to all authenticated users (not just admins)
- Added support for approval_document_path in ExpertRequest struct and database schema
- Updated handlers to accept CV and approval document file uploads
- Updated database operations to store document paths properly

#### Phase 4B: Request Listing with Status Filtering ✅
- Added status filter parameter to expert requests listing endpoint
- Implemented repository filter logic in ListExpertRequests method
- Added comprehensive tests for status filtering (pending/approved/rejected)
- Ensured proper validation and error handling for status filters

#### Phase 4C: Request Editing Before Approval ✅
- Enhanced expert request update endpoint to support multiple use cases:
  - Admins can edit any pending or rejected request
  - Users can edit their own rejected requests for corrections
- Added support for multipart form uploads during edits (CV and approval documents)
- Implemented proper permission checks to enforce access rules
- Added tests to verify permission handling and edit functionality

#### Phase 5A: Schema Updates for Approval Documents ✅
- Added approval_document_path to Expert and ExpertRequest structs
- Updated database schema in both expert and expert_request tables
- Ensured proper handling of the field in all database operations

#### Phase 5B: Single Request Approval with Document ✅
- Added validation to require an approval document when approving a request
- Ensured approval document path is copied to the expert when creating from a request
- Added tests to verify the approval document requirement

#### Phase 5C: Batch Approval Implementation ✅
- Created new batch approval endpoint at POST /api/expert-requests/batch-approve
- Implemented transactional batch approval to handle multiple requests at once
- Enhanced error handling to track success/failure of each request in the batch
- Added a shared approval document for all batch-approved experts

#### Phase 6A: Document Access Extension ✅
- Verified document endpoints already allow access to all authenticated users
- Document endpoints include GET /api/documents/{id} and GET /api/experts/{id}/documents
- Access control was already properly implemented in server.go

#### Phase 6B: Document Type Handling ✅
- Enhanced document type validation to include "cv", "approval", "certificate", "publication", and "other"
- Added clear error messages for invalid document types
- Improved validation in the document service

#### Phase 6C: Document Cascade Deletion ✅
- Implemented document and file deletion when deleting an expert
- Added file existence checks before deletion attempts
- Used transactions to ensure atomic operations
- Added proper error handling and logging

#### Files Modified:
- `/internal/domain/types.go` - Added approval_document_path field to both Expert and ExpertRequest structs
- `/db/migrations/sqlite/0002_create_expert-request_table.sql` - Added approval_document_path column
- `/db/migrations/sqlite/0004_create_expert_table_up.sql` - Updated schema with approval_document_path
- `/internal/storage/sqlite/expert.go` - Updated database operations for experts and added document cascade deletion
- `/internal/storage/sqlite/expert_request.go` - Updated database operations for requests and batch approval
- `/internal/api/handlers/expert_request.go` - Updated to handle file uploads, status filtering, and approval document validation
- `/internal/documents/service.go` - Enhanced document type validation
- `/internal/api/server.go` - Added batch approval endpoint and confirmed document access permissions
- `/internal/storage/interface.go` - Added BatchApproveExpertRequests method to the interface
- `/test_api.sh` - Added tests for approval document requirement and status filtering
- `/internal/api/server.go` - Updated permissions for expert request creation

## 2025-04-17: Phase 2D and 3A Implementation

### Completed Work:

#### Phase 2D: Resource Access Expansion ✅
- Updated `GET /api/expert/areas` endpoint to require authentication (previously public)
- Verified all other endpoints already had proper authentication
- Updated API documentation to reflect the authentication requirement
- Enhanced general notes section with access level descriptions

#### Phase 3A: Add Expert Filtering ✅
- Added new filter parameters to `/api/experts` endpoint:
  - `by_nationality` - Filter by Bahraini/non-Bahraini status
  - `by_general_area` - Filter by general area ID
  - `by_specialized_area` - Filter by specialized area (text search)
  - `by_employment_type` - Filter by employment type
  - `by_role` - Filter by role
- Fixed SQL NULL handling for datetime fields in expert queries
- Created test script (`test_filters.sh`) to verify filter functionality
- Added new test cases to `test_api.sh`
- Updated API documentation with new filter parameters

#### Files Modified:
- `/internal/api/server.go`
- `/internal/api/handlers/expert.go`
- `/internal/storage/sqlite/expert.go`
- `/internal/domain/types.go`
- `/db/migrations/sqlite/0004_create_expert_table_up.sql`
- `/ExpertDB API Endpoints Documentation.markdown`
- `/test_api.sh`
- Created `/test_filters.sh`
- Updated `/ExpertDB Implementation Plan.markdown`
