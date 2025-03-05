# ExpertDB Project Overview

## Project Documentation
This project has detailed documentation for implementation status and guidelines:

- [Main Implementation Plan](/IMPLEMENTATION.md) - Project-wide integration tracking across the full stack
- [Backend Implementation](/backend/IMPLEMENTATION.md) - Backend-specific implementation details
- [Frontend Implementation](/frontend/IMPLEMENTATION.md) - Frontend-specific implementation details
- [Backend Guidelines](/backend/CLAUDE.md) - Backend development guidelines
- [Frontend Guidelines](/frontend/CLAUDE.md) - Frontend development guidelines and UI/UX guidance
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
