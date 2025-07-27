# ExpertDB System Requirements Specification

## 1. Introduction

### 1.1 Purpose

This System Requirements Specification (SRS) defines the functional and non-functional requirements for the ExpertDB tool, an internal web application for managing a database of experts. The tool supports a small department (10-12 users) with a modest database (up to 2000 expert entries over 5 years). The SRS guides development by clarifying existing functionality, addressing gaps, and specifying new features based on stakeholder requirements and known implementation issues.

### 1.2 Scope

ExpertDB is a lightweight tool for managing expert profiles, requests, engagements, documents, phase planning, and statistics. It supports:

- User management with roles: `super_user`, `admin`, `user` with contextual elevations.
- Expert profile creation, management, and requests with mandatory approval documents.
- Document uploads (CVs, approval documents).
- Engagement tracking and phase planning for Qualification Placement (QP) and Institutional Listing (IL) applications.
- Statistics (annual growth, nationality representation, engagement counts by QP/IL).
- Specialization area management.
- Data import and CSV backup.

The tool uses Go, SQLite, and JWT authentication, emphasizing simplicity, internal use, and minimal dependencies. Security is handled organizationally, and high load is not expected.

### 1.3 Definitions

- **Super User**: A privileged user created during initialization, responsible for creating admin users.
- **Admin**: A user with full access to manage users, experts, requests, engagements, areas, and phase planning.
- **User**: A user who can submit expert requests, view expert data/documents, and view phases. Can be elevated to planner/manager for specific applications.
- **Planner Elevation**: Contextual privilege allowing a user to propose experts for specific applications within phases.
- **Manager Elevation**: Contextual privilege allowing a user to provide expert ratings for specific applications when requested by admin.
- **Expert**: A professional with a profile (e.g., name, institution, skills) in the database.
- **Expert Request**: A user-submitted proposal to add a new expert, pending admin review (statuses: `pending`, `rejected`, `approved`).
- **Engagement**: An expert’s assignment to a task (validator or evaluator) tied to QP or IL applications.
- **Specialization Area**: A category for classifying experts (e.g., Business, Engineering).
- **Phase Plan**: A planning period with applications (QP or IL) requiring expert assignments.
- **Application**: A task within a phase plan, either Qualification Placement (QP) or Institutional Listing (IL).

## 2. Overall Description

### 2.1 User Needs

- **Super Users** need to:
  - Create admin users during system setup.
- **Admins** need to:
  - Manage users, experts, and specialization areas.
  - Review, edit, approve/reject expert requests with mandatory approval documents, including batch approvals.
  - Create phase plans, assign applications to planners, and approve/reject planner-proposed experts.
  - View statistics (annual growth, nationality, engagements by QP/IL).
  - Generate CSV backups.
- **Users (Regular)** need to:
  - Submit expert requests with CVs.
  - Browse and filter expert profiles/documents.
- **Planner Users** need to:
  - Submit expert requests with CVs.
  - Propose experts (Expert-1, Expert-2) for phase plan applications.
  - Browse and filter expert profiles/documents.
- The system must support data import, CSV backups, and detailed statistics.

### 2.2 Assumptions and Constraints

- **Assumptions**:
  - The tool is internal, with organizational security measures.
  - SQLite is sufficient for the database size and load.
  - Accepted CV and approval documents are in PDF format.
  - Phase plan applications are predefined by admins as QP or IL.
  - Dependencies are minimal: `golang-jwt`, `google/uuid`, `mattn/go-sqlite3`, `golang.org/x/crypto`.
- **Constraints**:
  - Maximum 10-12 concurrent users.
  - Database growth capped at 2000 expert entries over 5 years.
  - No internet exposure; no advanced security features needed.
  - Minimal new dependencies to maintain simplicity.

## 3. Functional Requirements

### 3.1 User Management

#### FR1.1: Create User

- **Description**: A `super_user` shall be created during system initialization (`super_user` role, default credentials: `admin@expertdb.com`, `adminpassword`). Super users shall create admin users. Admins shall create regular and planner users with fields: `name`, `email`, `password`, `role` (regular, planner), `is_active`. Planner users inherit regular user privileges.
- **Current Implementation**: Implemented in `server.go:EnsureSuperUserExists` and `sqlite/user.go:CreateUserWithRoleCheck`. Role hierarchy (`super_user > admin > regular/planner`) is enforced.
- **Requirement**: No further changes needed.
- **Priority**: High (core functionality).

#### FR1.2: Authenticate User

- **Description**: Users (`super_user`, `admin`, `planner`, `regular`) shall log in with email and password, receiving a JWT token for authorized requests.
- **Current Implementation**: Supported via `POST /api/auth/login` (`auth.go`). All roles recognized in JWT claims.
- **Requirement**: No changes needed.
- **Priority**: High (core functionality).

#### FR1.3: Delete User

- **Description**: Admins shall delete regular and planner users by ID. Super users shall delete admin users. Deletion of planner users sets `phases.assigned_scheduler_id` to NULL.
- **Current Implementation**: Supported via `DELETE /api/users/{id}` (`user.go`). Super user deletion restricted to prevent deleting the last `super_user`. Cascade deletion for planner assignments partially implemented.
- **Requirement**: Validate cascade deletion for planner assignments in `sqlite/user.go`.
- **Priority**: Medium (administrative function).

### 3.2 Expert Management

#### FR2.1: Create Expert

- **Description**: Admins shall create expert profiles from approved requests, with fields: `expert_id` (auto-generated, `EXP-<sequence>`), `name` (required), `affiliation` (required), `email` (required, no validation), `designation` (required), `is_bahraini` (required), `is_available` (required), `rating` (set to 0 by default), `role` (required), `employment_type` (required), `general_area` (required, valid ID), `specialized_area` (required), `is_trained` (required), `cv_document_id` (required, references expert_documents table), `phone` (required), `is_published` (required, defaults to false), `experience_entries` (optional, stored in dedicated table), `education_entries` (optional, stored in dedicated table), `skills` (required), `approval_document_id` (required, references expert_documents table). Admins can edit requests before approval via `PUT /api/expert-requests/{id}/edit`.
- **Current Implementation**: Supported via `POST /api/experts`. Bug fixed for `UNIQUE constraint failed` in `sqlite/expert.go:GenerateUniqueExpertID`. Document management system implemented with `expert_documents` table for centralized file management. Edit endpoint implemented.
- **Requirement**: No changes needed.
- **Priority**: High (core functionality).

#### FR2.2: List Experts

- **Description**: All users shall retrieve a paginated list of experts with filters: `by_nationality` (Bahraini/non-Bahraini), `by_general_area` (area ID), `by_specialized_area` (text search against normalized specialized areas), `by_employment_type` (e.g., academic), `by_role` (e.g., evaluator). Sorting supported (e.g., `name`, `rating`, `institution`, `expert_id`, `is_bahraini`, `is_published`). Specialized areas are stored as comma-separated IDs and searched via normalized lookup.
- **Current Implementation**: Supported via `GET /api/experts` for all users, with filters and sorting in `sqlite/expert.go`. Pagination metadata included in response and headers (`X-Total-Count`, etc.). ID-based specialized areas normalization implemented with 327 specialized areas. Tested via `test_filters.sh`.
- **Requirement**: No changes needed.
- **Priority**: High (core functionality).

#### FR2.3: Retrieve Expert Details

- **Description**: All users shall retrieve a specific expert's details by ID, including `approval_document_id` for document reference.
- **Current Implementation**: Supported via `GET /api/experts/{id}` for all users (`expert.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (user function).

#### FR2.4: Update Expert

- **Description**: Admins shall update expert fields by ID, including `is_published`.
- **Current Implementation**: Supported via `PUT /api/experts/{id}` (`expert.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (administrative function).

#### FR2.5: Delete Expert

- **Description**: Admins shall delete an expert by ID, cascading to associated documents (`cv`, `approval_document`).
- **Current Implementation**: Supported via `DELETE /api/experts/{id}` (`test_api.sh`).
- **Requirement**: No changes needed.
- **Priority**: Medium (administrative function).

### 3.3 Expert Request Management

#### FR3.1: Create Expert Request

- **Description**: Regular and planner users shall submit expert requests with required fields: `name`, `designation`, `affiliation`, `is_bahraini`, `is_available`, `role`, `employment_type`, `general_area` (valid ID), `specialized_area`, `is_trained`, `phone`, `email`, `experience_entries` (optional array), `education_entries` (optional array), `skills`, `cv_document_id` (document reference after file upload). `is_published` is optional (default: `false`). Status defaults to `pending`. Users can select from existing specialized areas or suggest new area names when suitable options don't exist. Rating is set to 0 by default when expert is created from approved request.
- **Current Implementation**: Supported via `POST /api/expert-requests` (`api.go`). CV upload implemented (`documents/service.go`). Specialized area suggestions stored as JSON array in `suggested_specialized_areas` column.
- **Requirement**: No changes needed.
- **Priority**: High (core functionality).

#### FR3.2: List Expert Requests

- **Description**: Admins shall retrieve a paginated list of expert requests, filtered by status (`pending`, `approved`, `rejected`).
- **Current Implementation**: Supported via `GET /api/expert-requests` with status filtering (`sqlite/expert_request.go`).
- **Requirement**: No changes needed.
- **Priority**: High (administrative function).

#### FR3.3: Retrieve Expert Request Details

- **Description**: Admins shall retrieve a specific request's details by ID, including `cv_document_id` (document reference) and `suggested_specialized_areas` for review of user-proposed areas.
- **Current Implementation**: Supported via `GET /api/expert-requests/{id}` (`api.go`). Returns `suggested_specialized_areas` as JSON array for admin review.
- **Requirement**: No changes needed.
- **Priority**: Medium (administrative function).

#### FR3.4: Approve/Reject Expert Request

- **Description**: Admins shall approve (`status: approved`) or reject (`status: rejected`, with optional `rejection_reason`) individual requests, attaching a mandatory approval document for approvals. Approved requests create an expert record with `cv_document_id` and `approval_document_id` referencing documents in the expert_documents table. During approval, admin can review `suggested_specialized_areas` and create new specialized areas in the system if appropriate, then assign them to the expert.
- **Current Implementation**: Supported via `PUT /api/expert-requests/{id}` with document upload (`handlers/expert_request.go`). Suggested areas stored in expert request for admin review.
- **Requirement**: No changes needed.
- **Priority**: High (core functionality).

#### FR3.5: Batch Approve Expert Requests

- **Description**: Admins shall approve multiple requests with a single mandatory approval document, creating expert records.
- **Current Implementation**: Supported via `POST /api/expert-requests/batch-approve` with transactions (`sqlite/expert_request.go`).
- **Requirement**: No changes needed.
- **Priority**: High (administrative function).

### 3.4 Document Management

#### FR4.1: Upload Document

- **Description**: Regular and planner users shall upload CVs during expert request creation. Admins shall upload approval documents during request approval and additional CVs/approval documents for experts. Only `cv` and `approval_document` types are supported.
- **Current Implementation**: Supported via `POST /api/documents` and integrated into `POST /api/expert-requests` and `PUT /api/expert-requests/{id}` (`document_handler.go`).
- **Requirement**: No changes needed.
- **Priority**: High (core functionality).

#### FR4.2: List Documents

- **Description**: All users shall list documents (`cv`, `approval_document`) for a specific expert.
- **Current Implementation**: Supported via `GET /api/experts/{id}/documents` for all users (`auth/middleware.go`).
- **Requirement**: No changes needed.
- **Priority**: High (user function).

#### FR4.3: Retrieve Document

- **Description**: All users shall retrieve a specific document by ID.
- **Current Implementation**: Supported via `GET /api/documents/{id}` for all users.
- **Requirement**: No changes needed.
- **Priority**: Medium (user function).

#### FR4.4: Delete Document

- **Description**: Admins shall delete a document by ID, cascading with expert removal.
- **Current Implementation**: Supported via `DELETE /api/documents/{id}`.
- **Requirement**: No changes needed.
- **Priority**: Medium (administrative function).

### 3.5 Engagement Management

#### FR5.1: Create Engagement

- **Description**: Engagements (type: `validator`, `evaluator`) shall be created automatically when admins approve phase plan applications. Engagements are `validator` for QP applications, `evaluator` for IL applications. Admins can approve/reject applications individually, reopening for modification, updating/removing engagements.
- **Current Implementation**: Supported via `PUT /api/phases/{id}/applications/{app_id}/review` (`sqlite/phase.go`, `engagement_handler.go`).
- **Requirement**: No changes needed.
- **Priority**: High (core functionality).

#### FR5.2: List Engagements

- **Description**: Admins shall list engagements, filtered by expert or type (`validator`, `evaluator`).
- **Current Implementation**: Supported via `GET /api/expert-engagements` with filters (`engagement_handler.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (administrative function).

#### FR5.3: Update Past Engagements

- **Description**: Developers shall import past engagement data via `POST /api/engagements/import` using CSV/JSON, with fields: `expert_id`, `type`, `date`, `details`. Imports validate expert existence and deduplicate entries.
- **Current Implementation**: Supported via `POST /api/engagements/import` (`engagement_handler.go`, `sqlite/engagement.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (developer task).

### 3.6 Phase Planning

#### FR6.1: Create Phase Plan

- **Description**: Admins shall create a phase plan with applications (QP or IL), each detailing: `type` (QP, IL), `institution_name`, `qualification_name` (for QP), `expert_1` (planner-proposed expert ID), `expert_2` (planner-proposed expert ID), `status` (pending, approved, rejected). A planner user is assigned to propose experts.
- **Current Implementation**: Supported via `POST /api/phases` (`sqlite/phase.go`, `phase_handler.go`, `0012_create_phases.sql`).
- **Requirement**: No changes needed.
- **Priority**: Medium (new feature).

#### FR6.2: Propose Experts for Phase Plan

- **Description**: Assigned planner users shall propose expert IDs for `expert_1` and `expert_2` in each application.
- **Current Implementation**: Supported via `PUT /api/phases/{id}/applications/{app_id}` (`phase_handler.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (new feature).

#### FR6.3: Approve/Reject Phase Plan Applications

- **Description**: Admins shall approve or reject proposed experts per application, with notes for rejections. Approved applications create engagements (`validator` for QP, `evaluator` for IL). Admins can reopen applications for modification, updating engagements.
- **Current Implementation**: Supported via `PUT /api/phases/{id}/applications/{app_id}/review` with transactions (`phase_handler.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (new feature).

#### FR6.4: List Phase Plans

- **Description**: Admins shall list phase plans, filtered by status or assigned planner.
- **Current Implementation**: Supported via `GET /api/phases` with filters (`phase_handler.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (new feature).

### 3.7 Statistics and Reporting

#### FR7.1: Overall Statistics

- **Description**: Admins shall view system statistics: total experts, active count, Bahraini percentage, published count, published ratio, top areas, engagement types (`validator`, `evaluator`), most requested experts.
- **Current Implementation**: Supported via `GET /api/statistics` (`sqlite/statistics.go`). Includes `published_count` and `published_ratio`.
- **Requirement**: No changes needed.
- **Priority**: High (stakeholder requirement).

#### FR7.2: Annual Growth Statistics

- **Description**: Admins shall view expert growth by year (year-over-year, since last year).
- **Current Implementation**: Supported via `GET /api/statistics/growth?years=<int>` (`sqlite/statistics.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (stakeholder requirement).

#### FR7.3: Nationality Statistics

- **Description**: Admins shall view Bahraini vs. non-Bahraini distribution.
- **Current Implementation**: Supported via `GET /api/statistics/nationality` (`api.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (administrative function).

#### FR7.4: Engagement Statistics

- **Description**: Admins shall view engagement counts per expert by type (QP: `validator`, IL: `evaluator`).
- **Current Implementation**: Supported via `GET /api/statistics/engagements`, restricted to `validator` and `evaluator` (`sqlite/statistics.go`).
- **Requirement**: Update `sqlite/statistics.go` to map QP to `validator`, IL to `evaluator` in statistics.
- **Priority**: High (stakeholder requirement).

#### FR7.5: Area Statistics

- **Description**: Admins shall view statistics for general and specialized areas, including top 5 and bottom 5 areas by expert count.
- **Current Implementation**: Supported via `GET /api/statistics/areas` (`sqlite/statistics.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (stakeholder requirement).

### 3.8 Specialization Area Management

#### FR8.1: List General Areas

- **Description**: All users shall retrieve all general specialization areas.
- **Current Implementation**: Supported via `GET /api/expert/areas` for all users (`auth/middleware.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (user function).

#### FR8.2: List Specialized Areas

- **Description**: All users shall retrieve all specialized areas for search functionality.
- **Current Implementation**: Supported via `GET /api/specialized-areas` for all users. Returns 327 normalized specialized areas with ID-based references.
- **Requirement**: No changes needed.
- **Priority**: High (search functionality).

#### FR8.3: Create General Area

- **Description**: Admins shall create new general specialization areas with a unique name.
- **Current Implementation**: Supported via `POST /api/expert/areas` (`sqlite/area.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (stakeholder requirement).

#### FR8.4: Rename General Area

- **Description**: Admins shall rename general specialization areas, updating associated expert and request records.
- **Current Implementation**: Supported via `PUT /api/expert/areas/{id}` with transactions (`sqlite/area.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (stakeholder requirement).

#### FR8.5: Specialized Areas Data Model

- **Description**: The system shall use normalized specialized areas with ID-based storage for data consistency. Expert records store specialized areas as comma-separated IDs (e.g., "1,4,6") referencing the specialized_areas table.
- **Current Implementation**: Implemented via specialized_areas table with 327 entries, normalization via py_import.py, ID-based storage in experts table.
- **Requirement**: No changes needed.
- **Priority**: High (data integrity).

### 3.9 Data Import and Backup

#### FR9.1: Import Expert Data

- **Description**: Developers shall import expert data from CSV files, mapping fields to the `experts` table.
- **Current Implementation**: Supported via `py_import.py`.
- **Requirement**: No changes needed.
- **Priority**: Medium (existing functionality).

#### FR9.2: Generate CSV Backup

- **Description**: Admins shall generate a CSV backup of the database (`experts`, `expert_requests`, `expert_engagements`, `expert_documents`, `expert_areas`) as a ZIP file.
- **Current Implementation**: Supported via `GET /api/backup` (`backup_handler.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (stakeholder requirement).

### 3.10 Workflows

- **Description**: The system supports structured workflows for expert creation and phase planning to streamline operations.
- **Workflow 1: Expert Creation**:
  1. **User Submits Request**: A regular or planner user submits an `expert_request` via `POST /api/expert-requests`, filling a form with required fields, uploading a CV (stored via document service with `cv_document_id` reference), selecting from existing specialized areas, and optionally suggesting new specialized area names if suitable options don't exist.
  2. **Admin Reviews Request**: Admin receives the request (`GET /api/expert-requests`) and reviews both expert details and any suggested specialized areas:
     - **Approve**: Sets `status: approved` via `PUT /api/expert-requests/{id}`, uploads `approval_document`, reviews suggested areas and creates new specialized areas if appropriate, creates expert in `experts` table.
     - **Reject**: Sets `status: rejected` with `rejection_reason`, returns to user for amendment.
     - **Update and Approve**: Edits request via `PUT /api/expert-requests/{id}/edit`, then approves.
- **Workflow 2: Phase Planning**:
  1. **Admin Creates Phase**: Admin creates a phase via `POST /api/phases` with a title and applications (QP or IL).
  2. **Admin Creates Applications**: Applications are defined with `type` (QP, IL), `institution_name`, `qualification_name` (for QP).
  3. **Admin Assigns Applications**: Admin assigns applications (single or batch) to a planner via `POST /api/phases` or `PUT /api/phases/{id}/applications/{app_id}`.
  4. **Planner Proposes Experts**: Planner assigns `expert_1` and `expert_2` via `PUT /api/phases/{id}/applications/{app_id}`, submits for review.
  5. **Admin Reviews Proposals**:
     - **Approve**: Sets `status: approved` via `PUT /api/phases/{id}/applications/{app_id}/review`, creates engagements (`validator` for QP, `evaluator` for IL).
     - **Reject**: Sets `status: rejected` with `rejection_notes`, returns to planner for amendment or admin amends and approves.
- **Implementation**: Workflows are supported by existing endpoints (`expert_request.go`, `phase_handler.go`) and enforced via role-based access (`auth/middleware.go`).
- **Priority**: High (core functionality).

## 4. Non-Functional Requirements

### NFR1: Performance

- **Description**: The system shall handle up to 12 concurrent users with response times under 2 seconds for API requests under normal load (1200 expert entries).
- **Current Implementation**: SQLite and Go are adequate (`CLAUDE.md`). Indexes added in `0009_add_indexes.sql` for `nationality`, `general_area`, `specialized_area`, `employment_type`, `role`. Performance validated for batch approvals and statistics.
- **Requirement**: No changes needed.
- **Priority**: Medium.

### NFR2: Scalability

- **Description**: The system shall support up to 1200 expert entries without performance degradation.
- **Current Implementation**: SQLite supports scale with indexes (`0009_add_indexes.sql`).
- **Requirement**: No changes needed.
- **Priority**: Medium.

### NFR3: Usability

- **Description**: The API shall provide clear error messages for invalid inputs.
- **Current Implementation**: Improved in `errors/errors.go:ParseSQLiteError`. Handlers return aggregated validation errors per `ERRORS.md`.
- **Requirement**: No changes needed.
- **Priority**: High (user experience).

### NFR4: Maintainability

- **Description**: The codebase shall remain simple with minimal dependencies.
- **Current Implementation**: Uses Go, SQLite, and dependencies: `golang-jwt`, `google/uuid`, `mattn/go-sqlite3`, `golang.org/x/crypto` (`go.mod`).
- **Requirement**: No changes needed.
- **Priority**: High (core design principle).

### NFR5: Security

- **Description**: The system shall rely on organizational security, using JWT for authentication and role-based access control (`super_user`, `admin`, `planner`, `user`). JWT tokens expire after 24 hours. Permissions: `super_user` (create/delete admins), `admin` (manage users/experts/requests/phases), `planner` (propose experts, view/submit requests), `user` (view/submit requests).
- **Current Implementation**: JWT implemented with role checks (`auth/middleware.go`, `jwt.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (organizational security).

## 5. Future Considerations

- **API Development**: Status-filtered requests (FR3.2) and area statistics (FR7.5) support client applications. Ensure JSON-friendly API responses.
- **Phase Planning**: Validate cascade deletion for planner assignments (FR1.3).
- **Data Import/Backup API**: Enhance `py_import.py` as an API endpoint.
- **Testing**: Update `test_api.sh` to cover new endpoints (phase planning, engagement import, area statistics) and workflows.

## 6. Assumptions

- Users have basic API training.
- CSV imports and backups occur infrequently.
- Honorarium documents are pre-cleaned (CSV/JSON).
- Approval documents are validated as type `approval_document`, with flexibility for formats (e.g., PDF, DOC).
- Applications in phase plans are predefined by admins as QP or IL.
- Dependencies are minimal: `golang-jwt`, `google/uuid`, `mattn/go-sqlite3`, `golang.org/x/crypto`.