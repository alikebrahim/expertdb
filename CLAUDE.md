# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Summary
<memory>
ExpertDB is a lightweight internal web application for managing a database of experts with the following characteristics:
- Small scale: 10-12 users, up to 2000 expert entries over 5 years
- User roles: super_user, admin, planner, regular user
- Core features: expert profile management, request workflows, document management, engagement tracking, phase planning, statistics
- Backend: Go with SQLite database, JWT authentication, minimal dependencies
- Frontend: React, TypeScript, Tailwind CSS

Key workflows:
1. Expert Creation: User submits request with CV → Admin reviews → Approves with document → Expert profile created
2. Phase Planning: Admin creates phases → Assigns to planners → Planners propose experts → Admin approves/rejects

Technical constraints:
- Internal use only with organizational security
- Simplicity and maintainability are priorities
- No internet exposure needed
- Modest performance expectations (response times under 2 seconds)
</memory>

## Commands
- Backend: Use `go run cmd/server/main.go`
- Frontend: `npm run dev` (development), `npm run build` (production), `npm run lint` (lint)

## Code Style
- Backend (Go): Use modular architecture with domain, storage, service and API layers
- Error handling: Use custom error types from `internal/errors`
- Frontend (TypeScript): Strong typing with interfaces in `src/types`
- React components: Functional components with hooks
- CSS: Tailwind for styling
- API responses: Follow `ApiResponse<T>` pattern
- Imports: Group imports by source (stdlib, external, internal)
- Naming: camelCase for JS/TS, snake_case for database fields, PascalCase for React components

## Tool Usage
- Always use `bash -c "command"` syntax with Bash tool to avoid zoxide integration issues
- Use Playwright MCP tools for debugging and UI testing:
  - Navigate with `mcp__playwright__browser_navigate`
  - Take screenshots with `mcp__playwright__browser_take_screenshot`
  - Monitor console with `mcp__playwright__browser_console_messages`
  - View page UI with `mcp__playwright__browser_snapshot`
  - Interact with UI using `mcp__playwright__browser_click`, `mcp__playwright__browser_type`
  - Save all screenshots and logs to `./frontend_debugging/` directory
- For authentication testing, use these test credentials:
  - Admin: admin@expertdb.com / adminpassword
  - Planner: planner@expertdb.com / plannerpassword
  - Regular: user@expertdb.com / userpassword

## Core Entities and Relationships
<memory>
- User: Authentication with roles (super_user > admin > planner/regular)
- Expert: Professional profile with specialization areas and documents
- Expert Request: User proposal to add expert (pending → approved/rejected)
- Document: Files (CV, approval documents) attached to experts
- Engagement: Expert's assignment to tasks (validator or evaluator)
- Phase Plan: Planning period with applications requiring expert assignments
- Application: Task within phase plan (QP or IL)
- Specialization Area: Category for classifying experts

Workflows connect these entities:
- Expert creation workflow links requests → approval → expert profiles
- Phase planning workflow links phases → applications → expert assignments → engagements
</memory>

## API Structure
<memory>
The backend API follows RESTful conventions with these key endpoints:
- Auth: `/api/auth/login` - JWT authentication
- Users: `/api/users` - User management
- Experts: `/api/experts` - Expert profile CRUD, filtering, sorting
- Expert Requests: `/api/expert-requests` - Request workflow
- Documents: `/api/documents` - Document management
- Engagements: `/api/expert-engagements` - Engagement tracking
- Phases: `/api/phases` - Phase planning workflow
- Statistics: `/api/statistics` - System statistics and reports
- Specialization Areas: `/api/expert/areas` - Area management
- Backup: `/api/backup` - CSV backup generation

API responses follow the `ApiResponse<T>` pattern with consistent error handling.
</memory>

## Integration Plan Status
- Stage 1 (Search & Filter Foundation): Pending
- Stage 2 (Expert Creation Workflow): Pending
- Stage 3 (Phase Planning Workflow): Pending
- Stage 4 (Statistics & Reporting): Pending

When a stage is completed, update this file to mark the stage as "Completed" and create a corresponding STAGE#.complete.md file with implementation details and lessons learned.