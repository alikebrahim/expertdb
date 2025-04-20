# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Context
ExpertDB is a lightweight internal tool for managing a database of experts and their information. 
- Small user base (10-12 concurrent users)
- Limited data growth (max ~1200 expert entries over 5 years)
- Not exposed to the internet (organizational security measures)
- Role-based authentication (super_user, admin, regular, scheduler)

## System Purpose
ExpertDB supports:
- Expert profile management with approval workflows
- Document uploads (CVs, approval documents)
- User role hierarchy (super_user > admin > regular/scheduler)
- Engagement tracking for validators and evaluators
- Phase planning for applications requiring expert assignments
- Specialization area management
- Statistics and reporting
- Data import and CSV backup

## Build Commands
- Build server: `go build -o ./tmp/main ./cmd/server/main.go`
- Run server: `./tmp/main` or `go run cmd/server/main.go`
- Test API: `./test_api.sh`
- Format code: `go fmt ./...`

## Code Style Guidelines
- **Simplicity First**: Prefer simple, readable solutions over complex optimizations
- **Imports**: Standard library first, then third-party, then local packages
- **Error Handling**: Check errors and log appropriately using the logger; provide specific validation errors
- **Logging**: Use `internal/logger` package, not standard `log`
- **Dependencies**: Avoid adding new dependencies; this is a small internal tool
- **DB Access**: SQLite is sufficient for the scale - no need for complex DB solutions
- **Testing**: Focus on API-level testing via test_api.sh rather than extensive unit tests

## Architecture Guidelines
- Maintain basic layered approach but keep it simple
- Prefer direct CRUD operations over complex abstractions
- Balance maintainability and simplicity over strict architectural purity
- SQLite is perfectly adequate for the application scale
- Enforce role-based access control (super_user, admin, regular, scheduler)
- Support document uploads for CVs and approval documents
- Ensure transactional integrity for multi-table operations
- Use appropriate indexes for performance with filters

## Known Issues
- Need to complete testing and documentation (Phase 12)

## Implementation Progress (as of April 20, 2025)
- ✅ Phase 1A-1C: Completed bug fixes for expert creation, error handling, and database performance
- ✅ Phase 2A-2D: Completed user role structure updates and access control extensions
- ✅ Phase 3A-3C: Added expert filtering capabilities and improved pagination and sorting
- ✅ Phase 4A-4C: Completed expert request workflow improvements
- ✅ Phase 5A-5C: Implemented approval document integration including batch approval
- ✅ Phase 6A-6C: Enhanced document management including document access extension
- ✅ Phase 7A-7D: Added statistics and reporting enhancements
- ✅ Phase 8A-8C: Implemented specialization area management
- ✅ Phase 9A: Implemented CSV backup functionality
- ✅ Phase 10A-10E: Implemented phase planning and engagement system
- ✅ Phase 11A-11C: Enhanced engagement management with filtering and import functionality
- 🔜 NEXT: Phase 12: Testing and Documentation

Refer to the Implementation Plan document and Notes.md for complete details.
We are working through the Implementation Plan systematically - resume with Phase 12 in the next session.