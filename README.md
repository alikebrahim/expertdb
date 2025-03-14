# ExpertDB Project

ExpertDB is a comprehensive system for managing a database of experts, allowing users to search for experts based on various criteria and submit expert requests for review.

## Documentation

### Project Guidelines
- [Project Overview](/CLAUDE.md) - Main project architecture and guidelines
- [Implementation Plan](/IMPLEMENTATION.md) - Current implementation status and roadmap
- [Issue Tracking Guidelines](/ISSUES.md) - How to track and document issues
- [Git Strategy](/GIT_STRATEGY.md) - Branching model and commit conventions

### Issue Management
- [Frontend Issues](/frontend/issues.md) - Current frontend issues and status
- [Frontend Issue Archives](/frontend/issues/) - Detailed documentation of resolved issues
- [Backend Issue Archives](/backend/issues/) - Detailed documentation of resolved backend issues

## Project Structure

- **`backend/`**: Go-based REST API backend with SQLite database
- **`frontend/`**: Next.js frontend with TypeScript and shadcn/ui
- **`ai/`**: Future AI integration using Python with langchain

## Key Features

- Expert database management with ISCED classification
- Expert request submission and review workflow
- Expert search with advanced filtering
- Document management for expert profiles
- AI-assisted profile generation
- Statistics and reporting

## Getting Started

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Build the application:
   ```bash
   go build -o expertdb
   ```

3. Run the application:
   ```bash
   ./expertdb
   ```

See `backend/README.md` for more details.

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

See `frontend/README.md` for more details.

## Implemented Functionality

### Backend
- RESTful API endpoints for experts, expert requests, documents, and statistics
- SQLite database with migrations
- Expert data import from CSV
- ISCED classification integration

### Frontend
- Modern UI built with Next.js and shadcn/ui components
- Expert request submission form
- Expert search with filters
- Expert detail view

### Future Plans
- AI integration using langchain for:
  - PDF analysis and profile generation
  - ISCED classification suggestions
  - Specialized area suggestions

## System Architecture
ExpertDB follows a microservices architecture:
- Backend service: Go API handling business logic and data storage
- Frontend service: Next.js providing the user interface
- AI service: Python-based service for intelligent data processing (future)

Components communicate via REST APIs, with data flowing through well-defined interfaces.

## Development Workflow

This project follows a structured Git workflow:

1. Development happens on feature branches created from `develop`
   ```bash
   git checkout develop
   git checkout -b feature/new-feature
   ```

2. After testing, features are merged into `develop`
   ```bash
   git checkout develop
   git merge feature/new-feature
   ```
   
3. Releases are prepared on the `release` branch
   ```bash
   git checkout -b release/v1.0.0
   # Final testing and preparation
   ```
   
4. Production code is maintained on `main`
   ```bash
   git checkout main
   git merge release/v1.0.0
   git tag v1.0.0
   ```

See the [Git Strategy](/GIT_STRATEGY.md) for detailed guidelines on branching and commit conventions.
