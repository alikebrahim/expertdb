- Project Overview: What's the project? (e.g., a web app, mobile app, backend API—tech stack, main features). This helps me tailor file contents and focus areas.

ExpertDB is a web application for managing a database of experts with their profiles, skills, and documents. It features a Go-based backend API and React-based frontend. The main features include expert database management, expert request submission and review, document management,  expert search with advanced filtering, statistics and reporting, user management, and role-based access control.

- Current State (60-70% Done?): What's completed (e.g., backend APIs, frontend UI) and what's left (e.g., debugging, new features)? This sets the scope for Claude's context-building.

The project is approximately 60-70% complete. The backend implementation with authentication is complete, database schema established with automatic initialization, API endpoints defined with role-based security, and CSV import functionality finalized. The frontend has authentication and admin components developed. Docker configuration is ready for deployment. What's left: completing the login page styling and functionality, implementing expert search with advanced filtering, preparing expert bio page implementation, completing expert request page and functionality, implementing the statistics dashboard page, completing user management UI components, testing authentication and role-based access, running integration tests, and deploying the complete system.

- Pain Points: What's tricky about the codebase now? (e.g., messy logic, missing docs, bugs). This guides where Claude should dig deepest.

The codebase has some challenges with authentication persistence issues (recently fixed), potentially unclear API integration between frontend and backend which requires to be fixed (data such as users in admin view, experts in /search view are not being fetched and displayed properly in the frontend), and database migration management.  The project structure shows evidence of ongoing reorganization, with some files marked for deletion and others newly added.

- Tech Stack Details: Specific frameworks/libraries (e.g., React, Express, Django)? This refines ref docs like ENDPOINTS.md.

Backend: Go-based REST API with SQLite database, JWT authentication
Frontend: React with TypeScript, Vite for build management, and shadcn/ui component library
Database: SQLite with migrations
Authentication: JWT-based with role-based access control

- Preferred Starting Point: Any specific task for Claude to tackle first (e.g., documenting endpoints, fixing a bug)? This shapes the initial workflow.

Based on the current project state, the following tasks could be good starting points:
1. Fix frontend-backend integrations
2. Complete login page styling and functionality (highest priority according to documentation)
3. Implement expert search with advanced filtering
4. Prepare expert bio page implementation
5. Resolve any authentication persistence issues that might still exist
6. Document the API endpoints for better frontend-backend integration
