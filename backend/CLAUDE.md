# ExpertDB Go Project Guidelines

## Build Commands
```bash
# Build the application
go build -o expertdb

# Run the application
./expertdb

# Run tests
go test ./...

# Run a specific test
go test -v -run TestName

# Format code
go fmt ./...

# Lint code
go vet ./...
```

## Code Style Guidelines
- **Imports**: Standard library first, then third-party packages (separate with newline)
- **Naming**: Use camelCase for variables, PascalCase for exported functions/types
- **Error Handling**: Always check errors, use descriptive error messages
- **Types**: Favor strong typing, use interfaces for flexibility
- **Functions**: Keep functions focused on a single responsibility
- **Comments**: Document exported functions with proper godoc format
- **SQL**: Use parameterized queries to prevent SQL injection
- **HTTP**: Use route-specific handler functions with proper error handling

## Project Structure
- SQL migrations in db/migrations/sqlite/
- Main data types defined in types.go
- Authentication handling in auth.go and context.go
- User management in user_storage.go

## Authentication System
The system implements JWT-based authentication with two roles:
- **Admin**: Full access to all features, including user management
- **User**: Limited access to view data and submit requests

Default admin credentials (configurable via environment variables):
- Email: admin@expertdb.com
- Password: adminpassword
- Name: Admin User

Environment variables for admin configuration:
- ADMIN_EMAIL
- ADMIN_PASSWORD
- ADMIN_NAME

## User Management
All user management is handled by admins:
- Admins can create, update, and delete users
- No public registration is available
- Only admins can approve expert requests

## API Security
All sensitive endpoints are protected with authentication middleware:
- requireAuth: Ensures the user is authenticated
- requireAdmin: Ensures the user is an admin

## Current Project Status
We've been implementing enhancements for the ExpertDB backend system with these key activities:

### Authentication & User Management
- Added JWT-based authentication system
- Implemented role-based access control (admin/user)
- Added user management endpoints (create, update, delete)
- Protected sensitive endpoints with authentication middleware
- Added automatic admin user creation on startup

### Database Schema Extension
- Created new tables: `expert_documents`, `expert_engagements`, `ai_analysis_results`, and `system_statistics`
- Added nationality tracking to the experts table

### Code Implementation
- Updated `types.go` with new data structures for documents, engagements, AI analysis, statistics, and authentication
- Expanded the `Storage` interface in `storage.go` with methods for the new functionality
- Implemented storage methods for document handling, engagement tracking, and statistics generation
- Created service layer files:
  - `document_service.go` for file uploads and management
  - `ai_service.go` for AI integration (placeholder for now)
  - `auth.go` for authentication and authorization
  - `user_storage.go` for user management

### API Endpoints
- Added new REST endpoints in `api.go` for documents, engagements, AI features, statistics, and user management
- Implemented CORS middleware for cross-service communication
- Added authentication middleware for protected endpoints

### CSV Import Enhancement
- Modified `import_experts.go` with ISCED classification mapping
- Added nationality detection based on "BH" field

### Docker Setup
- Created `Dockerfile` and Docker Compose files
- Added placeholder implementations for frontend and AI services

### Current Focus
- Ensuring proper database initialization and admin user creation
- Proper authentication and authorization throughout the system
- Building comprehensive admin dashboard for user management

### Next Steps
1. Test the authentication system and user management
2. Verify proper database initialization
3. Test the expert request workflow with role-based permissions
4. Finalize the Docker configuration for integration
