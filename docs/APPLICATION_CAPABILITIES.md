# ExpertDB Application Capabilities

This document provides a comprehensive overview of the ExpertDB system's capabilities, workflows, and processes. It is intended to help developers understand the full range of functionality supported by the application.

## Table of Contents

1. [System Overview](#system-overview)
2. [User Roles and Permissions](#user-roles-and-permissions)
3. [Core Workflows](#core-workflows)
   - [Expert Creation Workflow](#expert-creation-workflow)
   - [Phase Planning Workflow](#phase-planning-workflow)
4. [Key Functional Areas](#key-functional-areas)
   - [Expert Management](#expert-management)
   - [Document Management](#document-management)
   - [Engagement Tracking](#engagement-tracking)
   - [Statistics and Reporting](#statistics-and-reporting)
5. [Operational Capabilities](#operational-capabilities)
   - [System Administration](#system-administration)
   - [Data Management](#data-management)

## System Overview

ExpertDB is a lightweight internal web application for managing a database of experts with the following characteristics:
- Small scale: Designed for 10-12 users, supporting up to 2000 expert entries over 5 years
- Internal use: Not exposed to the internet, deployed within the organization's security boundaries
- Performance: Optimized for response times under 2 seconds with modest server requirements
- Technology: Go backend with SQLite database

The system enables the organization to:
- Maintain a centralized database of qualified experts
- Manage the expert creation request and approval process
- Assign experts to qualification reviews (QP or IL applications)
- Track expert engagements and performance
- Generate statistics and reports on expert utilization

## User Roles and Permissions

The system implements a role-based access control model with three base roles and contextual elevations:

### Super User
- **Capabilities**: Full system access, including all administrative functions
- **Special Powers**: Can modify user roles, access system configuration, create admin users
- **Statistics Access**: Complete access to all statistics endpoints (**Enhanced: Now shared with Admin**)
- **Typical Tasks**: System setup, user provisioning, data recovery

### Admin
- **Capabilities**: Review and approve expert requests, manage phases, assign users to applications
- **Primary Workflows**: Expert approval, phase creation, user elevation assignments, final approval of expert assignments
- **Access Level**: Can view and edit all experts, requests, and phases. Has inherent access to all phase applications.
- **Statistics Access**: (**New Enhancement**) Full access to statistics dashboard and reports, including new request tracking statistics

### User
- **Base Capabilities**: Submit expert requests, view approved experts
- **Primary Workflows**: Expert request submission, viewing expert profiles
- **Access Level**: Limited to submitting requests and viewing approved expert profiles

### Contextual Elevations

Users can be elevated to have additional privileges for specific applications within phases:

#### Planner Elevation
- **Additional Capabilities**: Propose experts for assigned applications
- **Scope**: Limited to specific applications within assigned phases
- **Access**: Can update expert assignments for applications where they have planner elevation
- **Assignment**: Managed by admin via `/api/users/{id}/planner-assignments`

#### Manager Elevation
- **Additional Capabilities**: Rate experts for assigned applications when requested by admin
- **Scope**: Limited to specific applications within assigned phases
- **Access**: Can provide expert ratings and evaluations for applications where they have manager elevation
- **Assignment**: Managed by admin via `/api/users/{id}/manager-assignments`

## Core Workflows

### Expert Creation Workflow

The expert creation workflow manages how new experts are proposed, reviewed, and added to the system.

#### Process Flow
1. **Request Initiation**
   - User submits a new expert request (`/api/expert-requests` POST)
   - Required information: expert details, CV document upload
   - **Specialized Areas Selection**: Users can search and select from existing specialized areas, and if no suitable area exists, they can suggest new area names for admin review
   - System assigns "pending" status to the request

2. **Request Review**
   - Admin reviews pending requests (`/api/expert-requests` GET)
   - Admin examines details and attached CV
   - **Suggested Areas Review**: Admin reviews any suggested specialized areas from users and can create new areas in the system if appropriate

3. **Request Decision**
   - Admin approves or rejects the request (`/api/expert-requests/{id}` PUT)
   - For approval: Admin uploads approval document and can assign newly created specialized areas to the expert
   - For rejection: Admin provides rejection reason

4. **Expert Creation**
   - Upon approval, system automatically:
     - Creates new expert record
     - Generates unique expert ID (e.g., "EXP-0001")
     - Links original request ID to expert record
     - Associates CV and approval documents with expert

5. **Notification**
   - Request submitter notified of decision
   - For rejections, user can edit and resubmit the rejected request:
     - Users can edit their rejected requests using `/api/expert-requests/{id}/edit` endpoint
     - This allows for correction of issues that led to rejection
     - The resubmitted request will return to "pending" status for admin review

#### Key Features
- **Structured Professional Background System**: Database-driven professional background collection
  - Separate Experience and Education entry management with relational storage
  - Individual entry CRUD operations with proper foreign key relationships
  - Standardized data format stored in dedicated database tables
  - Enhanced data integrity and query capabilities
- **Enhanced Form Validation**: Comprehensive field validation with standardized options
  - Designation dropdown with professional titles: Prof., Dr., Mr., Ms., Mrs., Miss, Eng.
  - Performance rating is set to 0 by default when expert is created from approved request
  - Skills management via tag-based input system
- **Specialized Areas Management**: Flexible area assignment with user-driven expansion
  - Users can search and select from existing specialized areas during request creation
  - Users can suggest new specialized area names when suitable options don't exist
  - Suggested areas are stored with expert requests for admin review and approval
  - Admin can create new specialized areas based on user suggestions and assign them to approved experts
- **File Management**: Robust document handling with validation
  - PDF-only CV upload with 5MB size limit
  - Drag-and-drop interface with progress indication
- **Quality Assurance Features**:
  - Batch approval capability for processing multiple requests
  - Document validation for required CV format
  - Traceability from expert back to original request
  - Rejection reason tracking for quality improvement

### Phase Planning Workflow

The phase planning workflow manages the process of creating review phases, defining applications, assigning experts, and tracking reviews. This workflow is central to the system's purpose, allowing users with planner elevations to match qualified experts to appropriate applications.

#### Process Flow
1. **Phase Creation**
   - Admin creates new phase (`/api/phases` POST) 
   - Provides phase title and system automatically generates a unique phase ID
   - Sets initial phase status (typically "pending")

2. **Application Definition**
   - Admin adds applications to the phase as part of phase creation:
     - Application type must be either "QP" (Qualification Placement) or "IL" (Institutional Listing)
     - Requires institution name and qualification details
     - Applications are created with "pending" status by default
   - Each phase can contain multiple applications

3. **User Elevation Assignment**
   - Admin assigns users to specific applications within phases:
     - Planner elevation: Users can propose experts for assigned applications (`/api/users/{id}/planner-assignments`)
     - Manager elevation: Users can rate experts for assigned applications (`/api/users/{id}/manager-assignments`)
   - Users only have elevated privileges for applications they are specifically assigned to

4. **Expert Assignment**
   - Users with planner elevation view applications they're assigned to
   - For each assigned application, elevated user proposes two experts:
     - Expert 1 (Primary): Required for all applications
     - Expert 2 (Secondary): Optional but recommended for quality assurance
   - User submits expert selections (`/api/phases/{id}/applications/{app_id}` PUT)
   - Application status automatically changes to "assigned" upon submission

5. **Assignment Review**
   - Admin reviews expert assignments (`/api/phases/{id}` GET)
   - Admin approves or rejects each application's expert assignments (`/api/phases/{id}/applications/{app_id}/review` PUT)
   - For rejection: Admin must provide rejection notes explaining why the proposed experts are unsuitable
   - Rejected applications are returned to the elevated user for new expert proposals

6. **Engagement Creation**
   - Upon assignment approval, system automatically:
     - Creates engagement records for both assigned experts
     - Sets engagement type strictly based on application type:
       - QP applications ALWAYS create "validator" engagements
       - IL applications ALWAYS create "evaluator" engagements
     - Sets initial engagement status to "pending"
     - Links experts to the specific project/qualification via the projectName field
     - Records start date for the engagement

#### Key Features
- Role-based workflow with clear separation of responsibilities:
  - Admins create phases, assign user elevations, and make final approvals
  - Users with planner elevation assign appropriate experts based on expertise for their assigned applications
- Application types directly determine engagement types:
  - QP (Qualification Placement) → validator engagements
  - IL (Institutional Listing) → evaluator engagements
- Status tracking throughout the entire workflow:
  - Phases track overall status
  - Applications track assignment status
  - Engagements track completion status
- Comprehensive validation ensuring:
  - Experts are qualified for their assigned applications
  - Experts are available during the required timeframe
  - Proposed experts meet the criteria for the application type

## Key Functional Areas

### Expert Management

Expert management is the central capability of the system, providing comprehensive tools for maintaining the expert database.

#### Features
- **Expert Profiles**
  - Unique business identifiers (EXP-0001, EXP-0002, etc.)
  - Comprehensive profile information:
    - Professional details (name, designation, institution)
    - Qualifications (specialized area, is_trained)
    - Contact information (phone, email)
    - Nationality tracking (is_bahraini, nationality)
    - Professional background (experience and education entries stored in relational tables)

- **Expert Discovery**
  - Advanced search and filtering:
    - By name, designation, institution
    - By specialization area (ID-based normalized system)
    - By employment type 
    - By nationality (Bahraini status)
    - By availability and training status
  - Multi-field sorting capabilities
  - Pagination for large result sets
  - Specialized areas endpoint (`/api/specialized-areas`) provides searchable master list

- **Expert Classification**
  - Hierarchical specialization areas with normalized data model:
    - General areas: Broad categorization stored in `expert_areas` table
    - Specialized areas: Detailed expertise stored in `specialized_areas` table with ID-based references
    - Expert records store specialized areas as comma-separated IDs (e.g., "1,4,6") for data consistency
    - 327 normalized specialized areas derived from original CSV data normalization
  - Role categorization: Limited to only three options:
    - "evaluator" - For experts who evaluate IL applications
    - "validator" - For experts who validate QP applications
    - "evaluator/validator" - For experts who can perform both roles
  - Employment type tracking: Limited to only two options:
    - "academic" - For experts from academic institutions
    - "employer" - For experts from industry/employer organizations

- **Expert Lifecycle Management**
  - **Direct Expert Editing**: Any authenticated user can edit expert profiles directly
  - **Comprehensive Audit Trails**: All changes automatically tracked with user ID, timestamps, and field-level changes
  - Availability status updates
  - Performance rating tracking (defaults to 0 for new experts, can be updated by any authenticated user)
  - Publication status control (is_published) //NOTE: is published pertains to expert profile being published on another website. Ensure the implementation does not handle this as being published in the database
  - Record update history (created_at, updated_at, last_edited_by, last_edited_at)
  - **Edit History Viewing**: Complete audit trail accessible via `/api/experts/{id}/edit-history`

### Document Management

Document management enables secure storage and retrieval of expert-related documents within the system.

#### Features
- **Document Types**
  - CV documents: Expert qualifications and history
  - Approval documents: Official system approval

- **Document Operations**
  - Upload: Secure document storage with metadata
  - Retrieval: Access to stored documents
  - Deletion: Removal of obsolete documents

- **Document Metadata**
  - Original filename preservation
  - Content type (MIME) tracking
  - File size information
  - Upload date timestamps

- **Expert Association**
  - Multiple documents per expert
  - Document type categorization
  - Direct linkage between experts and documents

### Engagement Tracking

Engagement tracking monitors expert assignments to specific qualification reviews and projects, with a strict mapping between application types and engagement types.

#### Features
- **Engagement Types**
  - Strict mapping to application types:
    - "validator" engagements - Created ONLY from QP (Qualification Placement) applications
    - "evaluator" engagements - Created ONLY from IL (Institutional Listing) applications
  - No other engagement types are supported in the current implementation
  - Engagement type is automatically determined by application type and cannot be modified manually

- **Engagement Lifecycle**
  - Status tracking (pending, active, completed, cancelled)
  - Date range management (start_date, end_date)
  - Project association
  - Outcome recording (feedback_score)

- **Engagement Reporting**
  - Expert-specific engagement history
  - Project-based engagement tracking
  - Status-based filtering
  - Date range filtering

- **Engagement Operations**
  - Manual creation for ad-hoc assignments
  - Automatic creation from approved phase applications
  - Bulk import capability via CSV

### Statistics and Reporting

Statistics and reporting provide analytical insights into the expert database and engagement patterns.

#### Features
- **General Statistics**
  - Total expert count
  - Active (available) vs inactive experts
  - Trained vs untrained ratios
  - Published vs unpublished counts
  - Accessible via `/api/statistics` endpoint

- **Nationality Statistics**
  - Bahraini vs non-Bahraini expert ratios
  - Nationality distribution breakdown
  - Accessible via `/api/statistics/nationality` endpoint

- **Engagement Statistics**
  - Categorization strictly by type:
    - "validator" engagements (from QP applications)
    - "evaluator" engagements (from IL applications)
  - Status distribution of engagements
  - Accessible via `/api/statistics/engagements` endpoint
  - Also shows number of engagements per expert in profile page

- **Growth Statistics**
  - Expert database yearly growth rate
  - Time-based growth visualization
  - New experts by year
  - Accessible via `/api/statistics/growth` endpoint
  - Important: Growth statistics are aggregated by year (2023, 2024, etc.)
  - Experts with creation dates from 2022 or earlier are grouped as "Before 2023"

- **Specialization Statistics**
  - Expert distribution by general area
  - Top 5 most represented specialized areas
  - Bottom 5 least represented specialized areas
  - Accessible via `/api/statistics/areas` endpoint

## Operational Capabilities

### System Administration

The system provides administrative capabilities for efficient operation and maintenance.

#### Features
- **User Management**
  - User creation and role assignment
  - Account activation/deactivation
  - Password management
  - Login auditing (last_login tracking)

- **Access Control**
  - JWT-based authentication
  - Role-based permission enforcement
  - Session management
  - API security (protected endpoints)

- **System Monitoring**
  - Error logging and tracking
  - Activity auditing
  - Performance metrics

### Data Management

Data management capabilities ensure data integrity, backup, and maintenance.

#### Features
- **Backup and Restore**
  - CSV backup generation
  - Data export for reporting
  - External backup integration

- **Data Validation**
  - Input validation across all endpoints
  - Data integrity checks
  - Referential integrity enforcement

- **Data Lifecycle**
  - Soft deletion support
  - History tracking
  - Audit fields (created_at, updated_at, created_by)


---

This document provides a high-level overview of the capabilities and workflows supported by the ExpertDB application, including planned enhancements. Developers should refer to the API Reference, ENHANCEMENTS.md, and codebase for implementation details.
