# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Summary
<memory>
ExpertDB is a production-ready REST API for managing a database of experts with the following characteristics:
- Architecture: Backend-only REST API service (frontend removed from codebase)
- Scale: 10-12 users, up to 2000 expert entries over 5 years  
- User roles: super_user, admin, user (with contextual elevations to planner/manager)
- Core features: expert profile management, request workflows, document management, engagement tracking, phase planning, statistics, expert profile editing
- Backend: Go with SQLite database, JWT authentication, minimal dependencies
- Status: All planned integration stages completed + additional enterprise features

Key workflows:
1. Expert Creation: User submits request with CV → Admin reviews → Approves with document → Expert profile created
2. Expert Profile Editing: Users propose changes to existing experts → Admin reviews → Auto-applies approved changes
3. Phase Planning: Admin creates phases → Assigns users to applications → Users elevated to planner/manager for specific applications → Admin approves/rejects

Role System:
- super_user: Complete system access, can create admins
- admin: Full system access, can create users and manage all phases/applications  
- user: Limited access, can be elevated for specific applications within phases:
  - Planner elevation: Can propose experts for assigned applications
  - Manager elevation: Can rate experts for assigned applications (when requested by admin)

Technical characteristics:
- Production-ready enterprise-grade REST API
- Comprehensive security with contextual role elevation
- Advanced filtering, sorting, and search capabilities
- Complete audit trails and change tracking
- Batch operations and data import/export
- Internal use only with organizational security
</memory>

## Commands
- Backend: Use `go run cmd/server/main.go` (from root directory)
- Testing: API tests available in `api-tests/` directory using Hurl
- Database: SQLite database at `db/sqlite/main.db`
- Logs: Application logs in `logs/` directory

## Code Style
- Backend (Go): Use modular architecture with domain, storage, service and API layers
- Error handling: Use custom error types from `internal/errors`
- API responses: Follow `ApiResponse<T>` pattern
- Imports: Group imports by source (stdlib, external, internal)
- Naming: snake_case for database fields

## Tool Usage
- Always use `bash -c "command"` syntax with Bash tool to avoid zoxide integration issues
- For authentication testing, use these test credentials:
  - Admin: admin@expertdb.com / adminpassword
  - Regular User: user@expertdb.com / userpassword
  - Note: Users can be elevated to planner/manager for specific applications via admin assignment

## Core Entities and Relationships
<memory>
**Primary Entities:**
- User: Authentication with roles (super_user > admin > user) and contextual elevations
- Expert: Professional profile with specialization areas, experience, education, and documents
- Expert Request: User proposal to add expert (pending → approved/rejected)
- Expert Edit Request: User proposal to modify existing expert (pending → approved/rejected)
- Document: Files (CV, approval documents) attached to experts
- Engagement: Expert's assignment to tasks (validator or evaluator) - legacy system
- Phase Plan: Planning period with applications requiring expert assignments
- Application: Task within phase plan (QP or IL) with contextual user assignments
- Expert Areas: General categorization system for experts
- Specialized Areas: Detailed specialization categories with search capabilities
- Expert Experience: Professional work history entries
- Expert Education: Educational background entries
- Role Assignments: Contextual elevations (planner/manager) for users on specific applications

**Key Workflows:**
- Expert creation workflow: requests → approval → expert profiles
- Expert editing workflow: edit proposals → admin review → auto-application
- Phase planning workflow: phases → applications → contextual user elevations → expert assignments
- Role elevation system: temporary planner/manager privileges for specific applications
- Document management: upload, attachment, and retrieval system
- Batch operations: bulk approval, import/export capabilities
</memory>

## API Structure
<memory>
The backend API is production-ready with 50+ endpoints across 12 feature areas:

**Core APIs:**
- Auth: `/api/auth/login` - JWT authentication
- Users: `/api/users` - User management with role-based access
- Experts: `/api/experts` - Expert profile CRUD with advanced filtering/sorting
- Expert Requests: `/api/expert-requests` - Request workflow with batch operations
- Expert Edit Requests: `/api/experts/{id}/edit` - Profile editing workflow
- Documents: `/api/documents` - Document management and file handling
- Engagements: `/api/engagements` - Engagement tracking with CSV import
- Phases: `/api/phases` - Phase planning workflow with contextual assignments
- Role Assignments: `/api/users/{id}/{planner|manager}-assignments` - Contextual elevations
- Expert Areas: `/api/expert/areas` - General area management
- Specialized Areas: `/api/specialized-areas` - Detailed specializations with search
- Statistics: `/api/statistics` - Comprehensive system statistics
- Backup: `/api/backup` - Complete system backup to CSV/ZIP

**Advanced Features:**
- Multi-criteria filtering and sorting on all major endpoints
- Batch operations for bulk processing
- Contextual role-based access control with middleware
- Complete audit trails and change tracking
- Search capabilities with fuzzy matching
- CSV import/export functionality
- Document upload/download with secure file handling

API responses follow the `ApiResponse<T>` pattern with comprehensive error handling.
</memory>

## Role Elevation System
<memory>
The application implements a three-tier role system with contextual elevations:

**Base Roles:**
- `super_user`: Complete system access, can create admin users
- `admin`: Full system access, can create regular users and manage all phases/applications
- `user`: Can submit expert requests, view expert data/documents, and view all phases. Can be elevated for specific applications to propose experts (planner) or provide ratings upon admin request (manager)

**Contextual Elevations:**
Regular users can be elevated to have special privileges for specific applications within phases:

- **Planner Elevation**: Allows user to propose experts for assigned applications
  - Scoped to specific applications within a phase
  - Managed via `/api/users/{id}/planner-assignments` endpoints
  - Uses `RequirePlannerForApplication` middleware for access control

- **Manager Elevation**: Allows user to provide expert ratings for assigned applications when requested by admin
  - Scoped to specific applications within a phase
  - Managed via `/api/users/{id}/manager-assignments` endpoints
  - Uses `RequireManagerForApplication` middleware for access control

**Implementation:**
- Database tables: `application_planners` and `application_managers`
- Storage methods: `IsUserPlannerForApplication`, `IsUserManagerForApplication`
- Middleware: Admin/super_user bypass elevation checks and have inherent access
- API endpoints use application-specific access control rather than global role checks

This design provides fine-grained access control while maintaining system simplicity.
</memory>

## Project Structure
- Documented `docs/` directory structure for comprehensive project documentation
- Backend code organized in modular architecture (domain, storage, service, API layers)
- API tests available in `api-tests/` directory
- Database migrations in `backend/db/migrations/sqlite/`
- Application logs stored in `backend/logs/` directory

## Project Status
All integration stages completed. Production-ready REST API with 50+ endpoints, 18+ database tables, comprehensive security, and enterprise-grade features including expert profile editing, specialized areas management, experience/education tracking, advanced role assignments, document management, batch operations, audit trails, and CSV import/export.