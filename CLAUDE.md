# ExpertDB Project Overview

## Persona
I am a meticulous craftsman refining ExpertDB's Go backend and React frontend integration.

## Guidelines
My primary focus is analyzing the codebase, prioritizing frontend-backend integration, documenting gaps, and completing unfinished features like the login page. I will use the structured context files in the `.system/` directory to maintain comprehensive understanding of the project.

## Structured Context Files
The following files provide detailed context about different aspects of the project:
- [API Endpoints Map](/.system/ENDPOINTS.md) - Complete mapping of backend API endpoints
- [UI/UX Guidelines](/.system/UI_UX_GUIDELINES.md) - React component standards and current status
- [Function Signatures](/.system/FUNCTION_SIGNATURES.md) - Index of key Go/React functions
- [Authentication Guidelines](/.system/AUTH_GUIDELINES.md) - JWT auth flow and role-based logic
- [Implementation Status](/.system/IMPLEMENTATION.md) - Current progress and next steps
- [Master Record](/.system/MASTER_RECORD.md) - Timeline of actions and changes
- [Issue Log](/.system/ISSUE_LOG.md) - Tracking of bugs and their resolutions

## Project Documentation
This project has detailed documentation for implementation status and guidelines:

- [Implementation Status](/.system/IMPLEMENTATION.md) - Current progress and implementation details  
- [Issue Management System](/ISSUES.md) - Guidelines for tracking and resolving issues
- [Git Strategy](/GIT_STRATEGY.md) - Branching model and commit conventions
- [Frontend Issues Repository](/frontend/issues/) - Archive of resolved frontend issues
- [Backend Issues Repository](/backend/issues/) - Archive of resolved backend issues

## Project Architecture
- **Backend**: Go-based REST API with SQLite database and JWT authentication
- **Frontend**: Next.js with TypeScript, Vite for build management, and shadcn/ui component library
- **AI Integration**: To be implemented with Python using langchain
- **Authentication**: JWT-based authentication with role-based access control

## Key Features
- Expert database management
- Expert request submission and review (admin approval required)
- Document management for expert profiles
- AI-assisted profile generation from PDFs
- AI-suggested ISCED classifications and specialized areas
- Expert search with advanced filtering
- Statistics and reporting
- User management (admin-only)
- Role-based access control (admin/user)

## Authentication and Access Control
- **Two User Roles**:
  - Admin: Full access to all features, including user management and request approval
  - User: Limited access to view data and submit requests
- **No Public Registration**: All users are created by admin
- **Default Admin**: System automatically creates a default admin on first startup
- **JWT Authentication**: Secure token-based authentication
- **Protected Endpoints**: Role-based middleware secures sensitive endpoints

## AI Integration Plan
- **PDF Analysis**: Extract expert information from uploaded CVs/documents
  - Generate profile suggestions with human-in-the-loop approval
  - Extract skills, certifications, and experience
- **Classification Suggestions**:
  - Suggest appropriate ISCED classifications based on expert profile
  - Identify specialized areas based on employment history and certifications
- **Implementation**: Python service with langchain, integrated via REST API

## Backend Structure
- Go-based REST API with JWT authentication
- SQLite database with migrations
- Service-oriented architecture
- CSV import functionality for initial data population
- Automatic database initialization and admin user creation

## Frontend Stack
- **Framework**: Next.js with TypeScript and App Router
- **Build Tool**: Vite for fast development and optimized builds
- **UI Library**: shadcn/ui for accessible, customizable components
- **State Management**: React Context API and React Query for data fetching
- **Authentication**: JWT integration with secure storage
- **Key Features**:
  - Expert profile management
  - Expert request submission form
  - Search interface with advanced filtering
  - Document upload
  - AI-assisted profile creation
  - Admin dashboard with user management
  - Role-based UI components
  - Responsive design for desktop

## Current Status
- Backend implementation complete with authentication
- Database schema established with automatic initialization
- API endpoints defined with role-based security
- CSV import functionality finalized
- Frontend authentication and admin components developed
- Docker configuration ready for deployment

## Current Development Priority
Following a page-by-page systematic approach to ensure all main user-facing features are properly designed and functional:

1. Login page styling and functionality
2. Expert search with advanced filtering and expert selection
3. Expert request page and workflow
4. Statistics dashboard page

## Next Steps
1. Complete login page styling and functionality (first point of contact for users)
2. Implement expert search with advanced filtering and intelligent search
3. Prepare expert bio page implementation (backend model and frontend page)
4. Complete expert request page and functionality
5. Implement statistics dashboard page
6. Complete user management UI components
7. Test authentication and role-based access
8. Run integration tests for all main features
9. Deploy the complete system

See directory-specific CLAUDE.md files for detailed implementation guidelines.
