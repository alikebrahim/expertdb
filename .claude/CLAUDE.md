# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Memory Map

This project follows a structured approach with memory documents tracking different aspects of the system:

1. **API Documentation**: [API_DESC.md](API_DESC.md) - Core API endpoints and specifications
2. **API Integration Analysis**: [API_INTEG_DESC.md](API_INTEG_DESC.md) - Analysis of frontend-backend integration
3. **Implementation Plan**: [INTEG_PLAN.md](INTEG_PLAN.md) - Phase-by-phase plan for fixing integration issues
   - **Phase 1 Summary**: [PHASE1_SUMMARY.md](PHASE1_SUMMARY.md) - Authentication and URL standardization
   - **Phase 2 Summary**: [PHASE2_SUMMARY.md](PHASE2_SUMMARY.md) - Data structure alignment
   - **Phase 3 Summary**: [PHASE3_SUMMARY.md](PHASE3_SUMMARY.md) - Pagination implementation
   - **Phase 5 Summary**: [PHASE5_SUMMARY.md](PHASE5_SUMMARY.md) - Document management integration
   - **Phase 6 Summary**: [PHASE6_SUMMARY.md](PHASE6_SUMMARY.md) - Expert creation and management
   - **Phase 7 Summary**: [PHASE7_SUMMARY.md](PHASE7_SUMMARY.md) - Expert engagement integration
4. **Supplementary Documentation**:
   - **Authentication Report**: [AUTH_REPORT.md](AUTH_REPORT.md) - Authentication system review
   - **Frontend Report**: [FRONTEND_REPORT.md](FRONTEND_REPORT.md) - Frontend architecture analysis
   - **Frontend Debugging**: [FRONTEND_DEBUG.md](FRONTEND_DEBUG.md) - Frontend debugging guidelines

## Current Workflow: Phase-by-Phase Implementation

We followed a structured approach to address API integration issues, which has now been completed:

**All Phases Completed**:
   - ✅ **Phase 1**: Core Authentication and URL Standardization
   - ✅ **Phase 2**: Data Structure Alignment
   - ✅ **Phase 3**: Pagination Implementation
   - ⏭️ **Phase 4**: ISCED Classification Integration (Skipped - ISCED functionality removed)
   - ✅ **Phase 5**: Document Management Integration
   - ✅ **Phase 6**: Expert Creation and Management
   - ✅ **Phase 7**: Expert Engagement Integration

The integration plan has been successfully completed. The API documentation (API_DESC.md) and integration analysis (API_INTEG_DESC.md) have been updated to reflect the current state of the system.

## Memory Management Rules

1. **Document Creation**:
   - When creating a new memory document, add a link to it in CLAUDE.md
   - Each document should begin with a clear purpose statement and reference related documents

2. **Phase Completion**:
   - When completing a phase:
     1. Create a PHASEXX_SUMMARY.md document with details of changes
     2. Update INTEG_PLAN.md to mark the phase as complete with a link to the summary
     3. Update CLAUDE.md to reflect current progress

3. **Document Priority**:
   - INTEG_PLAN.md serves as the roadmap
   - PHASEXX_SUMMARY.md documents contain implementation details
   - API_DESC.md is the source of truth for API specifications

## Build Commands
- Frontend: `cd frontend && npm run dev` - Start frontend development server
- Frontend: `cd frontend && npm run build` - Build frontend for production
- Frontend: `cd frontend && npm run lint` - Run ESLint for frontend
- Backend: `cd backend && go run .` - Run backend server
- Full stack: `./run.sh` - Run both frontend and backend
- API Tests: `./test_api.sh` - Run API integration tests

## Code Style Guidelines

### Frontend (React/TypeScript)
- TypeScript: Strict typing with explicit interfaces
- Components: Functional components with hooks
- Imports: Group by source (React, external, internal)
- Naming: PascalCase for components/interfaces, camelCase for variables/functions
- Error handling: Use try/catch with proper error messages
- Styling: Tailwind with shadcn/ui components

### Backend (Go)
- Error handling: Return explicit errors, no panics
- Logging: Use structured logger (from logger.go)
- Function signatures: Return (data, error) pairs
- Database access: Use prepared statements
- Naming: CamelCase for exported, camelCase for internal
