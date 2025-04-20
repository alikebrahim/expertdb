# ExpertDB System Requirements Specification

## 1. Introduction

### 1.1 Purpose

This System Requirements Specification (SRS) defines the functional and non-functional requirements for the ExpertDB tool, an internal web application for managing a database of experts. The tool supports a small department (10-12 users) with a modest database (up to 1200 expert entries over 5 years). The SRS guides further development by clarifying existing functionality, addressing gaps, and specifying new features based on stakeholder requirements and known implementation issues.

### 1.2 Scope

ExpertDB is a lightweight tool for managing expert profiles, requests, engagements, documents, and statistics. It supports:

- User management (super_user, admin, regular, scheduler roles).
- Expert profile creation, management, and requests with mandatory approval documents.
- Document uploads (CVs, approval documents).
- Engagement tracking and phase planning.
- Statistics and reporting (including published status and area breakdowns).
- Specialization area management.
- Data import and CSV backup.

The tool uses Go, SQLite, and JWT authentication, emphasizing simplicity, internal use, and minimal dependencies. Security is handled organizationally, and high load is not expected.

### 1.3 Definitions

- **Super User**: A privileged user created during initialization, responsible for creating admin users.
- **Admin**: A user with full access to manage users, experts, requests, engagements, areas, and phase planning.
- **Regular User**: A user who can submit expert requests and view expert data/documents.
- **Scheduler User**: A user with regular privileges plus the ability to propose experts for phase planning.
- **Expert**: A professional with a profile (e.g., name, institution, skills) in the database.
- **Expert Request**: A user-submitted proposal to add a new expert, pending admin review.
- **Engagement**: An expert’s assignment to a task (validator or evaluator).
- **Specialization Area**: A category for classifying experts (e.g., Business, Engineering).
- **Phase Plan**: A planning period with a list of applications requiring expert assignments.

## 2. Overall Description

### 2.1 User Needs

- **Super Users** need to:
  - Create admin users during system setup.
- **Admins** need to:
  - Manage users, experts, and specialization areas.
  - Review, edit, and approve/reject expert requests with mandatory approval documents, including batch approvals.
  - Create phase plans and approve scheduler-proposed experts.
  - View statistics (published experts, general/specialized areas).
  - Generate CSV backups.
- **Regular Users** need to:
  - Submit expert requests with CVs.
  - View and filter expert profiles/documents.
- **Scheduler Users** need to:
  - Submit expert requests with CVs.
  - Propose experts for phase plan applications.
  - View and filter expert profiles/documents.
- The system must support data import, CSV backups, and annual growth statistics.

### 2.2 Assumptions and Constraints

- **Assumptions**:
  - The tool is internal, with organizational security measures.
  - Users are trained to use the API, requiring minimal UI complexity.
  - SQLite is sufficient for the database size and load.
  - Approval documents are PDFs provided by admins.
  - Phase plan applications are predefined by admins.
- **Constraints**:
  - Maximum 10-12 concurrent users.
  - Database growth capped at 1200 expert entries over 5 years.
  - No internet exposure; no advanced security features needed.
  - Minimal new dependencies to maintain simplicity.

## 3. Functional Requirements

### 3.1 User Management

#### FR1.1: Create User

- **Description**: A super user shall be created during system initialization (`super_user` role, default credentials: `admin@expertdb.com`, `adminpassword`). Super users shall create admin users. Admins shall create regular and scheduler users with fields: `name`, `email`, `password`, `role` (regular, scheduler), `is_active`. Scheduler users inherit regular user privileges.
- **Current Implementation**: A default admin user is created during initialization (see `server.go:174-201`). Only `admin` and `user` roles exist. Known issue: No `super_user` or `scheduler` roles, requiring schema changes in `domain/types.go`.
- **Requirement**: Add `super_user` role for initialization and `scheduler` role for user creation. Update `sqlite/user.go` to enforce role hierarchy (super_user &gt; admin &gt; regular/scheduler). Modify `server.go` to create a `super_user` instead of `admin`.
- **Priority**: High (new role hierarchy requirement).

#### FR1.2: Authenticate User

- **Description**: Users (super_user, admin, regular, scheduler) shall log in with email and password, receiving a JWT token for authorized requests.
- **Current Implementation**: Supported via `POST /api/auth/login` (see `auth.go`). No issues reported.
- **Requirement**: No changes needed; ensure `super_user` and `scheduler` roles are recognized in JWT claims.
- **Priority**: High (core functionality).

#### FR1.3: Delete User

- **Description**: Admins shall delete regular and scheduler users by ID. Super users shall delete admin users.
- **Current Implementation**: Supported via `DELETE /api/users/{id}` (see `auth.go`). No support for super_user deletion rights.
- **Requirement**: Update `auth/middleware.go` to restrict admin deletion to super users. Ensure cascade deletion of scheduler assignments.
- **Priority**: Medium (administrative function).

### 3.2 Expert Management

#### FR2.1: Create Expert

- **Description**: Admins shall create expert profiles from approved requests, with fields: `expert_id` (auto-generated incrementally), `name` (required), `institution` (required), `email` (required, no validation), `designation` (required), `is_bahraini` (required), `is_available` (required), `rating` (required), `role` (required), `employment_type` (required), `general_area` (required, valid ID), `specialized_area` (required), `is_trained` (required), `cv_path` (required), `phone` (required), `is_published` (required, defaults to false), `biography` (required), `skills` (required), `approval_document_path` (required). Admins can edit request details before approval.
- **Current Implementation**: Supported via `POST /api/experts`, but fails due to `UNIQUE constraint failed: experts.expert_id` (see `api.go`). Validation exists for `name`, `role`, `email`, `general_area`. Known issue: Non-unique `expert_id` generation in `sqlite/expert.go`; vague error messages (per `ERRORS.md`); no `approval_document_path` or edit-before-approval support.
- **Requirement**: Modify `sqlite/expert.go` to use incremental `expert_id` (e.g., `EXP-<sequence>` with database check). Remove `email` validation. Add `approval_document_path` to `experts` table. Add `PUT /api/expert-requests/{id}/edit` endpoint for admins to edit request details before approval. Implement `ERRORS.md` suggestions for detailed validation errors.
- **Priority**: High (core functionality, bug fix, new requirements).

#### FR2.2: List Experts

- **Description**: All users (admin, regular, scheduler) shall retrieve a paginated list of experts with filters: `by_nationality` (Bahraini/non-Bahraini), `by_general_area` (area ID), `by_specialized_area` (text), `by_employment_type` (e.g., academic), `by_role` (e.g., evaluator). Sorting shall be supported (e.g., by `name`, `rating`).
- **Current Implementation**: Supported via `GET /api/experts` for admins with pagination and limited filters (`sort_by`, `sort_order`). Known issue: Restricted to admins; limited filter options.
- **Requirement**: Extend access to all users in `auth/middleware.go`. Add filters (`by_nationality`, `by_general_area`, etc.) and sorting in `sqlite/expert.go`. Update API to accept query parameters (e.g., `/api/experts?by_general_area=1&sort_by=name`).
- **Priority**: High (new access and filter requirements).

#### FR2.3: Retrieve Expert Details

- **Description**: All users shall retrieve a specific expert’s details by ID.
- **Current Implementation**: Supported via `GET /api/experts/{id}` for admins (implemented in `expert.go`).
- **Requirement**: Extend access to all users. Include `approval_document_path` in response.
- **Priority**: Medium (expanded access).

#### FR2.4: Update Expert

- **Description**: Admins shall update expert fields by ID, including `is_published`.
- **Current Implementation**: Supported via `PUT /api/experts/{id}` (implemented in `expert.go`).
- **Requirement**: No changes needed; ensure all fields (including `approval_document_path`) are updatable.
- **Priority**: Medium (administrative function).

#### FR2.5: Delete Expert

- **Description**: Admins shall delete an expert by ID, cascading to associated documents and engagements.
- **Current Implementation**: Supported via `DELETE /api/experts/{id}` (tested in `test_api.sh`).
- **Requirement**: No changes needed; ensure cascade deletion includes approval documents.
- **Priority**: Medium (administrative function).

### 3.3 Expert Request Management

#### FR3.1: Create Expert Request

- **Description**: Regular and scheduler users shall submit expert requests with required fields: `name`, `designation`, `institution`, `is_bahraini`, `is_available`, `rating`, `role`, `employment_type`, `general_area` (valid ID), `specialized_area`, `is_trained`, `phone`, `email`, `biography`, `skills`, `cv_path` (file upload). `is_published` is optional, defaulting to `false`. Status defaults to `pending`.
- **Current Implementation**: Supported via `POST /api/expert-requests` (see `api.go`). Known issue: No CV upload support; not all fields are required; vague validation errors (per `ERRORS.md`).
- **Requirement**: Add CV file upload (multipart form data) to `POST /api/expert-requests`. Enforce all fields as required except `is_published` (default `false`). Update `expert_request.go` and `documents/service.go`. Implement detailed validation errors.
- **Priority**: High (new field requirements, CV upload).

#### FR3.2: List Expert Requests

- **Description**: Admins shall retrieve a paginated list of expert requests, filtered by status (`pending`, `approved`, `rejected`).
- **Current Implementation**: Supported via `GET /api/expert-requests` with pagination. Known issue: No status filtering.
- **Requirement**: Add `status` query parameter (e.g., `/api/expert-requests?status=pending`). Update `sqlite/expert_request.go`.
- **Priority**: High (tabbed UI support).

#### FR3.3: Retrieve Expert Request Details

- **Description**: Admins shall retrieve a specific request’s details by ID.
- **Current Implementation**: Supported via `GET /api/expert-requests/{id}` (see `api.go`).
- **Requirement**: No changes needed; include `cv_path` in response.
- **Priority**: Medium (administrative function).

#### FR3.4: Approve/Reject Expert Request

- **Description**: Admins shall approve (`status: approved`) or reject (`status: rejected`, with optional `rejection_reason`) individual requests, attaching a mandatory approval document for approvals. Approved requests create an expert record with `cv_path` and `approval_document_path`.
- **Current Implementation**: Supported via `PUT /api/expert-requests/{id}`. Known issue: No approval document support; vague database errors (per `ERRORS.md`).
- **Requirement**: Add mandatory file upload for approval documents (multipart form data). Store in `documents` table and set `approval_document_path`. Update `documents/service.go` and `sqlite/expert_request.go`. Improve error handling.
- **Priority**: High (mandatory approval document).

#### FR3.5: Batch Approve Expert Requests

- **Description**: Admins shall approve multiple requests with a single mandatory approval document, creating expert records.
- **Current Implementation**: Not supported.
- **Requirement**: Add `POST /api/expert-requests/batch-approve` endpoint accepting request IDs and an approval document. Create expert records with shared `approval_document_path`. Ensure transactional integrity. Update `documents/service.go` and `sqlite/expert_request.go`.
- **Priority**: High (batch approval requirement).

### 3.4 Document Management

#### FR4.1: Upload Document

- **Description**: Regular and scheduler users shall upload CVs during expert request creation (FR3.1). Admins shall upload approval documents during request approval (FR3.4, FR3.5) and additional documents for experts.
- **Current Implementation**: Admin uploads supported via `POST /api/documents` (see `document_handler.go`). Known issue: No user CV or approval document support.
- **Requirement**: Integrate CV upload into `POST /api/expert-requests` and approval document upload into `PUT /api/expert-requests/{id}` and `POST /api/expert-requests/batch-approve`. Maintain admin upload capability.
- **Priority**: High (CV and approval document requirements).

#### FR4.2: List Documents

- **Description**: All users (admin, regular, scheduler) shall list documents (CVs, approval documents) for a specific expert.
- **Current Implementation**: Supported via `GET /api/experts/{id}/documents` for admins. Known issue: Restricted to admins.
- **Requirement**: Extend access to all users in `auth/middleware.go`. Ensure CVs and approval documents are included.
- **Priority**: High (new access requirement).

#### FR4.3: Retrieve Document

- **Description**: All users shall retrieve a specific document by ID.
- **Current Implementation**: Supported via `GET /api/documents/{id}` for admins.
- **Requirement**: Extend access to all users. Ensure access to approval documents.
- **Priority**: Medium (expanded access).

#### FR4.4: Delete Document

- **Description**: Admins shall delete a document by ID.
- **Current Implementation**: Supported via `DELETE /api/documents/{id}`.
- **Requirement**: No changes needed; ensure cascade deletion with expert removal.
- **Priority**: Medium (administrative function).

### 3.5 Engagement Management

#### FR5.1: Create Engagement

- **Description**: Engagements (type: `validator`, `evaluator`) shall be created automatically when admins approve phase plan applications (FR6.3). Admins can approve/reject applications individually, reopening applications for modification, automatically adding/removing engagements.
- **Current Implementation**: Likely supported via `POST /api/expert-engagements`. Known issue: No automatic creation or dynamic modification.
- **Requirement**: Integrate engagement creation into `PUT /api/phases/{id}/applications/{app_id}/review`. Support reopening applications to update `expert_1`/`expert_2`, modifying `expert_engagements` table. Restrict `type` to `validator`/`evaluator`. Update `engagement_handler.go` and `sqlite/engagement.go`.
- **Priority**: High (automatic engagement requirement).

#### FR5.2: List Engagements

- **Description**: Admins shall list engagements, filtered by expert or type (`validator`, `evaluator`).
- **Current Implementation**: Likely supported via `GET /api/expert-engagements`.
- **Requirement**: Confirm endpoint and add filters (e.g., `expert_id`, `type`).
- **Priority**: Medium (administrative function).

#### FR5.3: Update Past Engagements

- **Description**: Developers shall update past expert engagements using cleaned data from honorarium documents at a later stage.
- **Current Implementation**: Not supported.
- **Requirement**: Provide a script or endpoint (e.g., `POST /api/engagements/import`) for developers to import engagement data (CSV/JSON, fields: `expert_id`, `type`, `date`, `details`). Update `expert_engagements` table. Ensure validation and deduplication.
- **Priority**: Medium (developer task).

### 3.6 Phase Planning

#### FR6.1: Create Phase Plan

- **Description**: Admins shall create a phase plan with a list of applications, each detailing: `type` (QP for qualification placement, IL for institutional listing), `institution_name`, `qualification_name` (for QPs), `expert_1` (scheduler-proposed expert ID), `expert_2` (scheduler-proposed expert ID), `status` (pending, approved, rejected). A scheduler user is assigned to propose experts.
- **Current Implementation**: Not supported. No phase planning tables or endpoints exist.
- **Requirement**: Add `POST /api/phases` endpoint to create phases with fields: `phase_id` (auto-generated), `title`, `applications` (list of `{type, institution_name, qualification_name, expert_1, expert_2, status}`), `assigned_scheduler_id`, `status`, `created_at`. Store in `phases` and `phase_applications` tables.
- **Priority**: Medium (new feature).

#### FR6.2: Propose Experts for Phase Plan

- **Description**: Assigned scheduler users shall propose expert IDs for `expert_1` and `expert_2` in each application.
- **Current Implementation**: Not supported.
- **Requirement**: Add `PUT /api/phases/{id}/applications/{app_id}` endpoint for schedulers to update `expert_1` and `expert_2`. Validate expert IDs.
- **Priority**: Medium (new feature).

#### FR6.3: Approve/Reject Phase Plan Applications

- **Description**: Admins shall approve or reject proposed experts per application, with notes for rejected proposals. Approved applications create engagements (`validator` or `evaluator`). Admins can reopen applications for modification, updating engagements.
- **Current Implementation**: Not supported.
- **Requirement**: Add `PUT /api/phases/{id}/applications/{app_id}/review` endpoint to set `status` (approved/rejected) and `rejection_notes`. Create engagements on approval. Support reopening via `status: pending`. Ensure transactional integrity.
- **Priority**: Medium (new feature).

#### FR6.4: List Phase Plans

- **Description**: Admins shall list phase plans, filtered by status or assigned scheduler.
- **Current Implementation**: Not supported.
- **Requirement**: Add `GET /api/phases` endpoint with filters (e.g., `status`, `scheduler_id`).
- **Priority**: Medium (new feature).

### 3.7 Statistics and Reporting

#### FR7.1: Overall Statistics

- **Description**: Admins shall view system statistics: total experts, active count, Bahraini percentage, published count, published ratio (published/total), top areas, engagement types (`validator`, `evaluator`), most requested experts.
- **Current Implementation**: Supported via `GET /api/statistics` (see `api.go`). Known issue: No `is_published` statistics.
- **Requirement**: Add `published_count` and `published_ratio` to response. Update `sqlite/statistics.go`.
- **Priority**: High (new statistics requirement).

#### FR7.2: Annual Growth Statistics

- **Description**: Admins shall view expert growth by year.
- **Current Implementation**: `GET /api/statistics/growth?months=6` shows monthly growth. Known issue: Monthly, not yearly.
- **Requirement**: Replace `months` with `years` (e.g., `/api/statistics/growth?years=5`). Update `sqlite/statistics.go`.
- **Priority**: Medium (stakeholder requirement).

#### FR7.3: Nationality Statistics

- **Description**: Admins shall view Bahraini vs. non-Bahraini distribution.
- **Current Implementation**: Supported via `GET /api/statistics/nationality` (see `api.go`).
- **Requirement**: No changes needed.
- **Priority**: Medium (administrative function).

#### FR7.4: Engagement Statistics

- **Description**: Admins shall view engagement counts by type (`validator`, `evaluator`).
- **Current Implementation**: Supported via `GET /api/statistics/engagements`. Known issue: Lists outdated types (e.g., evaluation).
- **Requirement**: Restrict to `validator` and `evaluator`. Update `sqlite/statistics.go`.
- **Priority**: High (engagement type requirement).

#### FR7.5: Area Statistics

- **Description**: Admins shall view statistics for general and specialized areas, including top 5 and bottom 5 areas by expert count.
- **Current Implementation**: Not supported. Top general areas are partially reported in `GET /api/statistics`.
- **Requirement**: Add `GET /api/statistics/areas` endpoint to report general area counts and specialized area counts, with top/bottom 5 for each. Update `sqlite/statistics.go` to query `general_area` and `specialized_area`.
- **Priority**: Medium (new statistics requirement).

### 3.8 Specialization Area Management

#### FR8.1: List Specialization Areas

- **Description**: All users shall retrieve all specialization areas.
- **Current Implementation**: Supported via `GET /api/expert/areas` for admins (see `api.go`).
- **Requirement**: Extend access to all users in `auth/middleware.go`.
- **Priority**: Medium (expanded access).

#### FR8.2: Create Specialization Area

- **Description**: Admins shall create new specialization areas with a unique name.
- **Current Implementation**: Not supported; areas are static (see `py_import.py`).
- **Requirement**: Add `POST /api/expert/areas` endpoint to insert area (`name` required, unique). Update `sqlite/area.go`.
- **Priority**: Medium (stakeholder requirement).

#### FR8.3: Rename Specialization Area

- **Description**: Admins shall rename specialization areas, updating associated expert and request records.
- **Current Implementation**: Not supported.
- **Requirement**: Add `PUT /api/expert/areas/{id}` endpoint to update `name`, cascading to `experts` and `expert_requests`. Ensure transactional integrity.
- **Priority**: Medium (stakeholder requirement).

### 3.9 Data Import and Backup

#### FR9.1: Import Expert Data

- **Description**: Developers shall import expert data from CSV files, mapping fields to the `experts` table.
- **Current Implementation**: Supported via `py_import.py`.
- **Requirement**: No changes needed; consider API endpoint (e.g., `POST /api/import`).
- **Priority**: Medium (existing functionality).

#### FR9.2: Generate CSV Backup

- **Description**: Admins shall generate a CSV backup of the database (`experts`, `expert_requests`, `expert_engagements`, `expert_documents`, `expert_areas`).
- **Current Implementation**: Not supported.
- **Requirement**: Add `GET /api/backup` endpoint to export tables as CSV files (zipped). Include fields like `expert_id`, `cv_path`, `approval_document_path`. Update `sqlite/store.go`.
- **Priority**: Medium (stakeholder requirement).

## 4. Non-Functional Requirements

### NFR1: Performance

- **Description**: The system shall handle up to 12 concurrent users with response times under 2 seconds for API requests under normal load (1200 expert entries).
- **Current Implementation**: SQLite and Go are adequate (see `CLAUDE.md`).
- **Requirement**: Monitor performance for batch approvals, backups, and area statistics.
- **Priority**: Medium.

### NFR2: Scalability

- **Description**: The system shall support up to 1200 expert entries without performance degradation.
- **Current Implementation**: SQLite supports current scale. Known issue: Ensure indexes for new filters (FR2.2).
- **Requirement**: Add indexes on `nationality`, `general_area`, `specialized_area`, `employment_type`, `role`, `phase_id`.
- **Priority**: Medium.

### NFR3: Usability

- **Description**: The API shall provide clear error messages for invalid inputs.
- **Current Implementation**: Known issue: Vague error messages (e.g., `role is required`, per `ERRORS.md`).
- **Requirement**: Implement `ERRORS.md` suggestions (list all validation errors, specific database errors, e.g., `expert_id already exists`).
- **Priority**: High (improves user experience).

### NFR4: Maintainability

- **Description**: The codebase shall remain simple with minimal dependencies.
- **Current Implementation**: Uses Go, SQLite, few dependencies (see `go.mod`).
- **Requirement**: Avoid new dependencies unless critical (e.g., file handling). Update documentation for new features.
- **Priority**: High (core design principle).

### NFR5: Security

- **Description**: The system shall rely on organizational security, using JWT for authentication and role-based access control (super_user, admin, regular, scheduler).
- **Current Implementation**: JWT implemented (see `auth/jwt.go`). Known issue: Role permissions need clarification.
- **Requirement**: Define permissions: super_user (create admins), admin (full access except super_user tasks), regular/scheduler (requests, view experts/documents). Ensure token expiration.
- **Priority**: Medium (organizational security assumed).

## 5. Future Considerations

- **Frontend Development**: Status-filtered requests (FR3.2) and area statistics (FR7.5) support UI tabs. Ensure JSON-friendly API responses.
- **Phase Planning**: Prioritize after core fixes (expert creation, approval documents, user access).
- **Data Import/Backup API**: Enhance `py_import.py` and CSV backup as API endpoints.
- **Testing**: Update `test_api.sh` to cover new endpoints (area statistics, phase planning, backups) and fix expert creation bug (FR2.1).

## 6. Assumptions

- Users have basic API training (no immediate frontend required).
- CSV imports and backups occur infrequently.
- Honorarium documents are pre-cleaned (CSV/JSON).
- Approval documents are PDFs provided by admins.
- Applications in phase plans are predefined by admins.