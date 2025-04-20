# ExpertDB Implementation Plan

## Overview
This plan outlines a phased approach to resolve critical bugs and implement new features for the ExpertDB system, a lightweight internal tool for managing expert profiles. The implementation prioritizes simplicity, maintainability, and clear error messaging, as security is handled organizationally and high load is not expected. The plan breaks down requirements from the System Requirements Specification (SRS) into smaller, actionable tasks with clear dependencies and validation steps.

- **Context**: Small department (10-12 users), max 1200 expert entries over 5 years, internal use, Go backend with SQLite, JWT authentication.
- **Goals**:
  - Fix critical bugs (e.g., expert creation UNIQUE constraint issue).
  - Implement high-priority features (e.g., approval documents, batch approvals, user role enhancements).
  - Enhance error messaging as per `ERRORS.md`.
  - Support new features like phase planning and CSV backups.
  - Ensure minimal dependencies and maintainable code.

## Phase 1: Critical Bug Fixes and Foundation
**Objective**: Resolve critical bugs blocking core functionality and improve error handling for better usability.

### 1A: Fix Expert Creation and ID Generation
**Tasks**:
1. **Fix expert_id generation in `sqlite/expert.go`**:
   - Modify `GenerateUniqueExpertID` to use sequential numbering (e.g., `EXP-0001`, `EXP-0002`).
   - Add database check to prevent ID collisions.
   - Test with multiple sequential insertions.
   - Validation: Verify unique IDs are created consistently.
   - Files: `internal/storage/sqlite/expert.go`

2. **Improve error handling for duplicate IDs**:
   - Update `CreateExpert` to catch and convert SQLite UNIQUE constraint errors to user-friendly messages.
   - Return HTTP 409 Conflict for duplicate ID scenarios.
   - Files: `internal/storage/sqlite/expert.go`, `internal/api/handlers/expert.go`

3. **Update test script for expert creation**:
   - Modify `test_api.sh` to include test cases for proper ID generation.
   - Add conflict handling test case.
   - Files: `test_api.sh`

### 1B: Error Messaging Improvements
**Tasks**:
1. **Implement validation error aggregation**:
   - Update request validation in `HandleCreateExpert` and `HandleCreateExpertRequest` to collect all validation errors.
   - Return them as an array in JSON response (e.g., `{"errors": ["name is required", "invalid general_area"]}`).
   - Files: `internal/api/handlers/expert.go`, `internal/api/handlers/expert_request.go`

2. **Create error parsing utility for SQLite errors**:
   - Create a helper function `ParseSQLiteError` in `internal/errors/errors.go`.
   - Map common SQLite errors to user-friendly messages.
   - Handle both UNIQUE constraint and foreign key violations with specific messaging.
   - Files: `internal/errors/errors.go`

3. **Integrate error parsing into handlers**:
   - Use the error parser in all handler functions dealing with database operations.
   - Log the original error but return the user-friendly version.
   - Files: All handler files in `internal/api/handlers/`

### 1C: Database Performance Enhancements
**Tasks**:
1. **Create migration for missing indexes**:
   - Add migration file `db/migrations/sqlite/0009_add_indexes.sql` with indexes for:
     - `idx_experts_nationality` on `experts(nationality)`
     - `idx_experts_general_area` on `experts(general_area)`
     - `idx_experts_specialized_area` on `experts(specialized_area)`
     - `idx_experts_employment_type` on `experts(employment_type)`
     - `idx_experts_role` on `experts(role)`
   - Files: `db/migrations/sqlite/0009_add_indexes.sql`

2. **Update store.go to apply new migrations**:
   - Ensure the migration runner picks up the new migration files.
   - Add unit tests for migration application.
   - Files: `internal/storage/sqlite/store.go`

3. **Validate index performance**:
   - Create benchmark queries using `EXPLAIN QUERY PLAN` for the new indexes.
   - Verify query plans show index usage.
   - Files: `test_api.sh` (new performance test section)

### 1D: Field Validation Adjustments
**Tasks**:
1. **Remove email validation for experts**:
   - Update `domain.Expert` struct in `types.go` to make email optional without format checks.
   - Remove email validation logic from `HandleCreateExpert`.
   - Files: `internal/domain/types.go`, `internal/api/handlers/expert.go`

2. **Clarify required fields for expert creation**:
   - Make `name`, `institution`, `designation`, `is_bahraini`, `is_available`, `rating`, `role`, `employment_type`, `general_area`, `specialized_area`, `is_trained`, `phone`, `biography`, `skills` explicitly required.
   - Update validation logic and error messages accordingly.
   - Files: `internal/domain/types.go`, `internal/api/handlers/expert.go`

## Phase 2: User Management and Access Control
**Objective**: Implement extended user roles and role-based access control as specified in SRS.

### 2A: User Role Structure Updates
**Tasks**:
1. **Update user role enumeration in domain model**:
   - Modify `User.Role` type in `domain/types.go` to include `super_user` and `scheduler` roles.
   - Add clear role hierarchy definition.
   - Files: `internal/domain/types.go`

2. **Create migration for role column updates**:
   - Create migration file `db/migrations/sqlite/0010_update_user_roles.sql` to update any constraints or defaults on the role column.
   - Files: `db/migrations/sqlite/0010_update_user_roles.sql`

3. **Update JWT token generation and parsing**:
   - Modify JWT claims in `auth/jwt.go` to handle new roles.
   - Add token expiration time (e.g., 24 hours).
   - Files: `internal/auth/jwt.go`

### 2B: Initialization and Creation Flow
**Tasks**:
1. **Update server initialization process**:
   - Rename `EnsureAdminExists` to `EnsureSuperUserExists` in `server.go`.
   - Update creation logic to use `super_user` role instead of `admin`.
   - Files: `internal/api/server.go`

2. **Implement role hierarchy for user creation**:
   - Update `CreateUser` in `sqlite/user.go` to enforce role creation rules:
     - `super_user` can create `admin`.
     - `admin` can create `regular`/`scheduler` users.
   - Add appropriate error messages for unauthorized role creation attempts.
   - Files: `internal/storage/sqlite/user.go`

3. **Test super user initialization and role creation**:
   - Add tests in `test_api.sh` for super user creation.
   - Test role restriction violations.
   - Files: `test_api.sh`

### 2C: Role-Based Middleware
**Tasks**:
1. **Refactor middleware for role-based permissions**:
   - Update `auth/middleware.go` to support detailed permission checks.
   - Create a more flexible middleware function allowing specific role checks.
   - Files: `internal/auth/middleware.go`

2. **Restrict admin deletion to super users**:
   - Update `DELETE /api/users/{id}` handler to only allow super users to delete admin accounts.
   - Add role checking to the delete logic.
   - Files: `internal/api/handlers/user.go`

3. **Update scheduler user handling**:
   - Add support for cascade deletion of scheduler assignments when deleting a scheduler user.
   - Files: `internal/storage/sqlite/user.go`

### 2D: Resource Access Expansion ✅ COMPLETED
**Tasks**:
1. **Extend expert listing access to all users** ✅:
   - Update middleware for `GET /api/experts` to allow all authenticated users.
   - Maintain admin-only restriction for modifications.
   - Files: `internal/auth/middleware.go`, `internal/api/server.go`

2. **Extend document access to all users** ✅:
   - Update middleware for `GET /api/experts/{id}/documents` and `GET /api/documents/{id}` endpoints.
   - Files: `internal/auth/middleware.go`, `internal/api/server.go`

3. **Extend specialization area access** ✅:
   - Update middleware for `GET /api/expert/areas` to allow all authenticated users.
   - Files: `internal/auth/middleware.go`, `internal/api/server.go`

## Phase 3: Expert List Enhancements and Filters
**Objective**: Improve expert listing with extended filtering and sorting options.

### 3A: Add Expert Filtering ✅ COMPLETED
**Tasks**:
1. **Add filter parameters to expert listing endpoint** ✅:
   - Update `GET /api/experts` to accept query parameters:
     - `by_nationality` (Bahraini/non-Bahraini)
     - `by_general_area` (area ID)
     - `by_specialized_area` (text search)
     - `by_employment_type` (e.g., academic)
     - `by_role` (e.g., evaluator)
   - Files: `internal/api/handlers/expert.go`

2. **Implement filter logic in repository** ✅:
   - Update `ListExperts` in `sqlite/expert.go` to apply filters to SQL query.
   - Use prepared statements for safe parameter handling.
   - Files: `internal/storage/sqlite/expert.go`

3. **Add tests for filter combinations** ✅:
   - Extend `test_api.sh` to test multiple filter combinations.
   - Created dedicated test script `test_filters.sh` for comprehensive filter testing.
   - Test edge cases (zero results, all results).
   - Files: `test_api.sh`, `test_filters.sh`

### 3B: Sorting and Pagination Improvements
**Tasks**:
1. **Enhance sorting options**:
   - Add support for sorting by more fields (`name`, `rating`, `institution`, etc.).
   - Implement safe column name validation to prevent SQL injection.
   - Files: `internal/api/handlers/expert.go`, `internal/storage/sqlite/expert.go`

2. **Improve pagination responses**:
   - Return total count, page info, and result counts in response headers.
   - Add metadata to simplify client-side pagination UI.
   - Files: `internal/api/handlers/expert.go`

### 3C: Expert Detail Access ✅ COMPLETED
**Tasks**:
1. **Update expert details endpoint access**:
   - Extend `GET /api/experts/{id}` access to all authenticated users.
   - Files: `internal/auth/middleware.go`, `internal/api/server.go`

2. **Include approval document path in responses**:
   - Update expert detail response structure to include `approval_document_path`.
   - Files: `internal/api/handlers/expert.go`

## Phase 4: Expert Request Workflow Improvements
**Objective**: Enhance expert request creation, filtering, and approval workflow.

### 4A: Expert Request Creation with CV Upload ✅ COMPLETED
**Tasks**:
1. **Update request creation endpoint for file uploads** ✅:
   - Modify `POST /api/expert-requests` to accept multipart form data.
   - Include `file` field for CV upload, plus JSON fields for request details.
   - Files: `internal/api/handlers/expert_request.go`

2. **Implement CV storage logic** ✅:
   - Update `documents/service.go` to handle expert request CV uploads.
   - Create directory structure if needed.
   - Generate unique filenames to prevent collisions.
   - Files: `internal/documents/service.go`

3. **Update request validation** ✅:
   - Enforce required fields: `name`, `designation`, `institution`, etc.
   - Make `is_published` optional (default to `false`).
   - Add file type validation for CV uploads.
   - Files: `internal/api/handlers/expert_request.go`, `internal/domain/types.go`

### 4B: Request Listing with Status Filtering ✅ COMPLETED
**Tasks**:
1. **Add status filter to expert request listing** ✅:
   - Update `GET /api/expert-requests` to accept `status` query parameter (`pending`, `approved`, `rejected`).
   - Files: `internal/api/handlers/expert_request.go`

2. **Implement status filtering in repository** ✅:
   - Modify `ListExpertRequests` in `sqlite/expert_request.go` to filter by status.
   - Add efficient indexes for status column if needed.
   - Files: `internal/storage/sqlite/expert_request.go`

3. **Add tests for status filtering** ✅:
   - Extend `test_api.sh` to test status filtering options.
   - Files: `test_api.sh`

### 4C: Request Editing Before Approval ✅ COMPLETED
**Tasks**:
1. **Enhance request update endpoint** ✅:
   - Update the existing `PUT /api/expert-requests/{id}` endpoint to handle different use cases:
     - Allow admins to edit any request (pending or rejected)
     - Allow users to edit their own rejected requests for corrections
   - Support multipart form data for CV and approval document updates.
   - Files: `internal/api/handlers/expert_request.go`

2. **Implement proper permission checks** ✅:
   - Updated permission logic to check user role and request ownership
   - Ensure users can only edit their own rejected requests
   - Ensure admins can edit any request
   - Files: `internal/api/handlers/expert_request.go`

3. **Test request editing** ✅:
   - Add test cases for request editing to `test_api.sh`.
   - Test various permission scenarios
   - Files: `test_api.sh`

## Phase 5: Approval Document Integration
**Objective**: Implement approval document requirement for expert requests as specified in SRS.

### 5A: Schema Updates for Approval Documents ✅ COMPLETED
**Tasks**:
1. **Add approval document path to experts table** ✅:
   - Added `approval_document_path TEXT` column to both experts and expert_requests tables.
   - Files: `db/migrations/sqlite/0004_create_expert_table_up.sql`, `db/migrations/sqlite/0002_create_expert-request_table.sql`

2. **Update domain model** ✅:
   - Added `ApprovalDocumentPath` field to both `Expert` and `ExpertRequest` structs in `domain/types.go`.
   - Updated JSON serialization tags with appropriate omitempty for optional fields.
   - Files: `internal/domain/types.go`

### 5B: Single Request Approval with Document ✅ COMPLETED
**Tasks**:
1. **Modify request approval endpoint** ✅:
   - Updated `PUT /api/expert-requests/{id}` to accept multipart form data.
   - Added validation to require an approval document for approval actions.
   - Files: `internal/api/handlers/expert_request.go`

2. **Implement approval document storage** ✅:
   - Enhanced `documents/service.go` to handle approval document uploads.
   - Used existing document storage mechanisms with unique naming.
   - Files: `internal/documents/service.go`

3. **Update expert creation flow** ✅:
   - Modified request approval process to copy `approval_document_path` when creating expert record.
   - Files: `internal/storage/sqlite/expert_request.go`

4. **Test approval with documents** ✅:
   - Added test case to verify rejection when trying to approve without document.
   - Files: `test_api.sh`

### 5C: Batch Approval Implementation ✅ COMPLETED
**Tasks**:
1. **Create batch approval endpoint** ✅:
   - Added `POST /api/expert-requests/batch-approve` endpoint.
   - Implemented to accept multiple request IDs and a single approval document.
   - Files: `internal/api/handlers/expert_request.go`, `internal/api/server.go`

2. **Implement transactional batch approval** ✅:
   - Added `BatchApproveExpertRequests` method to `sqlite/expert_request.go`.
   - Used SQLite transactions to ensure consistency during batch operations.
   - Set the same `approval_document_path` for all approved experts.
   - Implemented detailed error tracking for individual request failures.
   - Files: `internal/storage/sqlite/expert_request.go`, `internal/storage/interface.go`

3. **Test batch approval** ✅:
   - Added basic test case in `test_api.sh` for batch approval.
   - Note: Full testing would require more complex multipart form handling.
   - Files: `test_api.sh`

## Phase 6: Document Management Enhancements
**Objective**: Improve document handling and access for all user roles.

### 6A: Document Access Extension ✅ COMPLETED
**Tasks**:
1. **Update middleware for document endpoints** ✅:
   - Verified that access to `GET /api/experts/{id}/documents` is already available for all authenticated users.
   - Verified that access to `GET /api/documents/{id}` is already available for all authenticated users.
   - Files: `internal/api/server.go`

2. **Include CV and approval documents in responses** ✅:
   - Confirmed document listings already include approval documents.
   - Files: `internal/api/handlers/documents/document_handler.go`

### 6B: Document Type Handling ✅ COMPLETED
**Tasks**:
1. **Enhance document type validation** ✅:
   - Added explicit validation for document types to include `cv`, `approval`, `certificate`, `publication`, and `other`.
   - Added clear error messages for invalid document types.
   - Files: `internal/documents/service.go`

2. **Update document service for type-specific handling** ✅:
   - Enhanced document service to properly validate document types.
   - Maintained consistent storage paths across document types.
   - Files: `internal/documents/service.go`

### 6C: Document Cascade Deletion ✅ COMPLETED
**Tasks**:
1. **Ensure document deletion on expert deletion** ✅:
   - Enhanced `DeleteExpert` in `sqlite/expert.go` to handle document deletion.
   - Implemented deletion of both files and database records.
   - Used transaction to ensure atomic operations.
   - Files: `internal/storage/sqlite/expert.go`

2. **Implement document existence check** ✅:
   - Added file existence checks before deletion attempts.
   - Added graceful handling for missing files with appropriate logging.
   - Ensured database records are deleted even if files are missing.
   - Files: `internal/storage/sqlite/expert.go`

## Phase 7: Statistics and Reporting Enhancements ✅ COMPLETED
**Objective**: Implement additional statistics features as specified in SRS.

### 7A: Published Expert Statistics ✅
**Tasks**:
1. **Add published statistics to main statistics endpoint** ✅:
   - Updated `GET /api/statistics` to include `published_count` and `published_ratio`.
   - Files: `internal/api/handlers/statistics/statistics_handler.go`

2. **Implement published stats calculation** ✅:
   - Added `GetPublishedExpertStats()` method to calculate published expert statistics.
   - Files: `internal/storage/sqlite/statistics.go`

3. **Test published statistics** ✅:
   - Added test cases for published statistics to `test_api.sh`.
   - Files: `test_api.sh`

### 7B: Growth Statistics Enhancement ✅
**Tasks**:
1. **Convert growth statistics to yearly** ✅:
   - Updated `GET /api/statistics/growth` to accept `years` parameter instead of `months`.
   - Files: `internal/api/handlers/statistics/statistics_handler.go`

2. **Implement yearly growth calculation** ✅:
   - Created new `GetExpertGrowthByYear(years int)` method with improved year formatting.
   - Files: `internal/storage/sqlite/statistics.go`

3. **Test yearly growth** ✅:
   - Added test cases for yearly growth to `test_api.sh`.
   - Files: `test_api.sh`

### 7C: Engagement Type Statistics ✅
**Tasks**:
1. **Update engagement type statistics** ✅:
   - Restricted `GET /api/statistics/engagements` to `validator` and `evaluator` types.
   - Files: `internal/api/handlers/statistics/statistics_handler.go`

2. **Update type validation** ✅:
   - Modified engagement type query to filter by the correct types.
   - Files: `internal/storage/sqlite/statistics.go`

### 7D: Area Statistics Implementation ✅
**Tasks**:
1. **Create area statistics endpoint** ✅:
   - Added `GET /api/statistics/areas` endpoint for super_user role.
   - Implemented return of general area and specialized area counts.
   - Included top 5 and bottom 5 areas by expert count.
   - Files: `internal/api/handlers/statistics/statistics_handler.go`, `internal/api/server.go`

2. **Implement area statistics calculation** ✅:
   - Added `GetAreaStatistics()` method to calculate area statistics.
   - Used efficient queries with proper result categorization.
   - Files: `internal/storage/sqlite/statistics.go`

3. **Test area statistics** ✅:
   - Added test cases for area statistics to `test_api.sh`.
   - Files: `test_api.sh`

## Phase 8: Specialization Area Management ✅ COMPLETED
**Objective**: Implement specialization area creation and renaming as specified in SRS.

### 8A: Area Access Extension ✅
**Tasks**:
1. **Update middleware for area endpoint** ✅:
   - Extend access to `GET /api/expert/areas` to all authenticated users.
   - Files: `internal/auth/middleware.go`, `internal/api/server.go`

### 8B: Area Creation ✅
**Tasks**:
1. **Create area creation endpoint** ✅:
   - Add `POST /api/expert/areas` endpoint for admins.
   - Require unique area name.
   - Files: `internal/api/handlers/expert.go`

2. **Implement area creation in repository** ✅:
   - Add area creation method to `sqlite/area.go`.
   - Enforce uniqueness of area names.
   - Files: `internal/storage/sqlite/area.go`

3. **Test area creation** ✅:
   - Add test cases for area creation to `test_api.sh`.
   - Test duplicate name handling.
   - Files: `test_api.sh`

### 8C: Area Renaming ✅
**Tasks**:
1. **Create area rename endpoint** ✅:
   - Add `PUT /api/expert/areas/{id}` endpoint for admins.
   - Files: `internal/api/handlers/expert.go`

2. **Implement transactional area renaming** ✅:
   - Add area renaming method to `sqlite/area.go`.
   - Use transaction to cascade updates to `experts` and `expert_requests`.
   - Files: `internal/storage/sqlite/area.go`

3. **Test area renaming** ✅:
   - Add test cases for area renaming to `test_api.sh`.
   - Verify cascade updates.
   - Files: `test_api.sh`

## Phase 9: CSV Backup Implementation ✅ COMPLETED
**Objective**: Implement CSV backup functionality as specified in SRS.

### 9A: CSV Backup Endpoint ✅
**Tasks**:
1. **Create backup endpoint** ✅:
   - Add `GET /api/backup` endpoint for admins.
   - Return ZIP file with CSV exports.
   - Files: `internal/api/handlers/backup/backup_handler.go`

2. **Implement CSV generation** ✅:
   - Added CSV export methods in backup handler.
   - Export tables: `experts`, `expert_requests`, `expert_engagements`, `expert_documents`, `expert_areas`.
   - Include fields like `expert_id`, `cv_path`, `approval_document_path`.
   - Files: `internal/api/handlers/backup/backup_handler.go`

3. **Implement ZIP creation** ✅:
   - Used `archive/zip` from Go standard library.
   - Created temporary files for CSVs.
   - Compressed into a single ZIP archive.
   - Files: `internal/api/handlers/backup/backup_handler.go`

4. **Test backup generation** ✅:
   - Added test case for backup to `test_api.sh`.
   - Verified ZIP content structure.
   - Files: `test_api.sh`

## Phase 10: Phase Planning and Engagement System ✅ COMPLETED
**Objective**: Implement phase planning functionality as specified in SRS.

### 10A: Phase Planning Schema ✅
**Tasks**:
1. **Create phase planning tables** ✅:
   - Create migration file `db/migrations/sqlite/0012_create_phases.sql`.
   - Add `phases` table with fields: `id`, `phase_id`, `title`, `assigned_scheduler_id`, `status`, `created_at`.
   - Add `phase_applications` table with fields: `id`, `phase_id`, `type`, `institution_name`, `qualification_name`, `expert_1`, `expert_2`, `status`, `rejection_notes`.
   - Add appropriate indexes and foreign keys.
   - Files: `db/migrations/sqlite/0012_create_phases.sql`

2. **Update domain model** ✅:
   - Add `Phase` and `PhaseApplication` structs to `domain/types.go`.
   - Add validation rules.
   - Files: `internal/domain/types.go`

### 10B: Phase Creation ✅
**Tasks**:
1. **Create phase creation endpoint** ✅:
   - Add `POST /api/phases` endpoint for admins.
   - Accept phase details with list of applications.
   - Files: `internal/api/handlers/phase/phase_handler.go`

2. **Implement phase creation in repository** ✅:
   - Add phase creation methods to new `sqlite/phase.go`.
   - Use transaction to create phase and applications together.
   - Files: `internal/storage/sqlite/phase.go`

3. **Test phase creation** ✅:
   - Add test cases for phase creation to `test_api.sh`.
   - Files: `test_api.sh`

### 10C: Expert Proposal for Applications ✅
**Tasks**:
1. **Create expert proposal endpoint** ✅:
   - Add `PUT /api/phases/{id}/applications/{app_id}` endpoint for scheduler users.
   - Allow updating `expert_1` and `expert_2` fields.
   - Files: `internal/api/handlers/phase/phase_handler.go`

2. **Implement proposal validation** ✅:
   - Add expert ID validation to proposal logic.
   - Verify experts exist and are valid for the application.
   - Files: `internal/storage/sqlite/phase.go`

3. **Test expert proposals** ✅:
   - Add test cases for expert proposals to `test_api.sh`.
   - Test invalid expert IDs.
   - Files: `test_api.sh`

### 10D: Application Review and Engagement Creation ✅
**Tasks**:
1. **Create application review endpoint** ✅:
   - Add `PUT /api/phases/{id}/applications/{app_id}/review` endpoint for admins.
   - Support approve/reject actions with notes for rejections.
   - Support reopening applications for modification.
   - Files: `internal/api/handlers/phase/phase_handler.go`

2. **Implement automatic engagement creation** ✅:
   - Add logic to create engagements on application approval.
   - Create `validator` or `evaluator` engagements based on application type.
   - Handle reopening by updating/removing engagements.
   - Files: `internal/storage/sqlite/phase.go`

3. **Test application review** ✅:
   - Add test cases for application approvals/rejections to `test_api.sh`.
   - Verify engagement creation.
   - Test reopening flow.
   - Files: `test_api.sh`

### 10E: Phase Listing ✅
**Tasks**:
1. **Create phase listing endpoint** ✅:
   - Add `GET /api/phases` endpoint for admins.
   - Support filters for `status` and `scheduler_id`.
   - Files: `internal/api/handlers/phase/phase_handler.go`

2. **Implement phase listing in repository** ✅:
   - Add phase listing method to `sqlite/phase.go`.
   - Apply filters to SQL query.
   - Files: `internal/storage/sqlite/phase.go`

3. **Test phase listing** ✅:
   - Add test cases for phase listing to `test_api.sh`.
   - Test filters.
   - Files: `test_api.sh`

## Phase 11: Engagement Management ✅ COMPLETED
**Objective**: Enhance engagement management with filtering and import functionality.

### 11A: Engagement Filtering ✅
**Tasks**:
1. **Add filters to engagement listing endpoint** ✅:
   - Updated `GET /api/engagements` to accept filters: `expert_id`, `type`.
   - Added support for pagination with `limit` and `offset` parameters.
   - Files: `internal/api/handlers/engagements/engagement_handler.go`

2. **Implement filter logic in repository** ✅:
   - Updated `ListEngagements` method in `sqlite/engagement.go` to apply filters.
   - Added query builder for dynamic filter composition.
   - Files: `internal/storage/sqlite/engagement.go`

3. **Test engagement filtering** ✅:
   - Added test cases for engagement filtering to `test_api.sh`.
   - Tested various filter combinations.
   - Files: `test_api.sh`

### 11B: Engagement Type Restriction ✅
**Tasks**:
1. **Update engagement type validation** ✅:
   - Restricted engagement `type` to `validator` or `evaluator`.
   - Updated validation logic in create and update operations.
   - Added explicit type checking with appropriate error messages.
   - Files: `internal/storage/sqlite/engagement.go`

### 11C: Engagement Import ✅
**Tasks**:
1. **Create engagement import endpoint** ✅:
   - Added `POST /api/engagements/import` endpoint for administrators.
   - Implemented support for both CSV and JSON data formats.
   - Added validation for required fields: `expert_id`, `type`, `date`.
   - Files: `internal/api/handlers/engagements/engagement_handler.go`

2. **Implement import logic** ✅:
   - Added `ImportEngagements` method to `sqlite/engagement.go`.
   - Implemented validation for expert existence and engagement types.
   - Added deduplication checks based on expert, type, date, and project.
   - Used transactions for atomic batch operations.
   - Files: `internal/storage/sqlite/engagement.go`

3. **Test engagement import** ✅:
   - Added test cases for both CSV and JSON imports to `test_api.sh`.
   - Verified successful import and proper error handling.
   - Files: `test_api.sh`

## Phase 12: Testing and Documentation
**Objective**: Comprehensive testing and documentation updates.

### 12A: Integration Testing
**Tasks**:
1. **Update test_api.sh for all new endpoints**:
   - Add test cases for all new endpoints.
   - Cover edge cases and error handling.
   - Files: `test_api.sh`

2. **Create test data generation script**:
   - Add script to generate test data for all features.
   - Useful for quick setup and testing.
   - Files: `scripts/generate_test_data.sh`

### 12B: Documentation Updates
**Tasks**:
1. **Update API documentation**:
   - Update `ExpertDB API Endpoints Documentation.markdown` with all new endpoints.
   - Add detailed request/response examples.
   - Files: `ExpertDB API Endpoints Documentation.markdown`

2. **Update SRS documentation**:
   - Update `ExpertDB System Requirements Specification.markdown` to reflect implemented features.
   - Files: `ExpertDB System Requirements Specification.markdown`

3. **Update README**:
   - Update `README.md` with new setup instructions if needed.
   - Document any new environment variables.
   - Files: `README.md`

### 12C: Deployment Preparation
**Tasks**:
1. **Create deployment checklist**:
   - Document migration application steps.
   - List environment variables and configuration options.
   - Files: `docs/deployment.md`

2. **Build and test deployment package**:
   - Create release build script.
   - Test on staging environment.
   - Files: `scripts/build_release.sh`

## Dependencies and Timeline
- **Dependencies**:
  - Phase 1 (Critical Bug Fixes) is a prerequisite for all other phases.
  - User Role updates (Phase 2) should precede access control changes.
  - Schema changes (approval documents, phases) should precede related feature implementation.
  - Testing should be integrated throughout all phases.

- **Estimated Timeline**:
  - Phase 1 (Critical Bug Fixes): 1-2 weeks
  - Phase 2 (User Management): 1-2 weeks
  - Phase 3 (Expert Listing): 1 week
  - Phase 4 (Expert Requests): 1-2 weeks
  - Phase 5 (Approval Documents): 1-2 weeks
  - Phase 6 (Document Management): 1 week
  - Phase 7 (Statistics): 1-2 weeks
  - Phase 8 (Areas): 1 week
  - Phase 9 (Backup): 1 week
  - Phase 10 (Phase Planning): 2-3 weeks
  - Phase 11 (Engagements): 1-2 weeks
  - Phase 12 (Testing & Documentation): 1-2 weeks
  - **Total**: 14-22 weeks

## Implementation Approach
1. **Start with foundation fixes** - Address critical bugs before adding new features
2. **Work in small, testable increments** - Each task should be individually testable
3. **Maintain backward compatibility** - Avoid breaking existing clients when possible
4. **Test as you go** - Update test script with each feature addition
5. **Keep documentation current** - Update docs alongside code changes

## Conclusion
This detailed plan breaks down the ExpertDB implementation into manageable tasks with clear dependencies and validation steps. By following this incremental approach, we can systematically address bugs, implement new features, and maintain a simple, maintainable codebase for the department's expert management needs.