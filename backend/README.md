# ExpertDB Backend

A Go-based backend for the ExpertDB system that manages expert profiles, documents, engagements, and provides statistics. It's designed to integrate with separate frontend and AI services via Docker Compose.

## Features

- **Expert Management**: CRUD operations for expert profiles
- **Document Handling**: Upload, store, and manage expert documents (CVs, certificates, etc.)
- **Engagement Tracking**: Record and monitor expert assignments and activities
- **ISCED Classification**: Support for international standard classification of education
- **AI Integration**: APIs for AI-generated profiles, ISCED suggestions, and skills extraction
- **Statistics**: Comprehensive metrics on experts, engagements, and system usage

## Architecture

The system is designed with a microservice architecture:

- **Backend (this repository)**: Handles data storage, business logic, and API serving
- **Frontend** (separate repository): UI for user interaction
- **AI Service** (separate repository): Provides AI-powered features for profile generation and document analysis

## Getting Started

### Prerequisites

- Go 1.22+
- Docker and Docker Compose (for containerized deployment)
- SQLite (for local development)

### Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database path (default: ./db/sqlite/expertdb.sqlite)
- `UPLOAD_PATH`: Document storage path (default: ./data/documents)
- `CORS_ALLOWED_ORIGINS`: CORS configuration (default: *)
- `AI_SERVICE_URL`: URL for AI service (default: http://localhost:9000)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/expertdb-backend.git
   cd expertdb-backend
   ```

2. Initialize the database:
   ```bash
   mkdir -p db/sqlite
   goose -dir db/migrations/sqlite sqlite3 ./db/sqlite/expertdb.sqlite up
   ```

3. Import expert data from CSV (optional):
   ```bash
   go build -o import_csv ./cmd/import_csv
   ./import_csv -csv path/to/your/experts.csv -db ./db/sqlite/expertdb.sqlite
   ```

4. Run the server:
   ```bash
   go run *.go
   ```

### CSV Import Instructions

The CSV import tool is a standalone utility for importing expert data. It requires that you have already run the database migrations using goose before importing:

1. Build the import utility:
   ```
   go build -o import_csv ./cmd/import_csv
   ```

2. Run the import utility:
   ```
   ./import_csv -csv path/to/your/data.csv -db ./db/sqlite/expertdb.sqlite
   ```

The CSV file should have the following columns:
- ID - Unique identifier
- Name - Expert's full name
- Designation - Professional title
- Institution - Affiliated institution
- BH - Whether the expert is Bahraini (Yes/No)
- Available - Availability status (Yes/No)
- Rating - Expert rating
- Validator/ Evaluator - Role specification
- Academic/Employer - Employment type
- General Area - Main area of expertise
- Specialised Area - Specific specialization
- Trained - Training status (Yes/No)
- Phone - Contact number
- Email - Contact email
- Published - Whether profile is published (Yes/No)

### Docker Deployment

1. Build and run using Docker Compose:
   ```bash
   docker-compose up -d
   ```

This will start the backend service and create necessary volumes for data persistence.

## API Endpoints

### Expert Management
- `GET /api/experts`: List experts with filtering
- `POST /api/experts`: Create a new expert
- `GET /api/experts/{id}`: Get expert details
- `PUT /api/experts/{id}`: Update expert
- `DELETE /api/experts/{id}`: Delete expert

### Document Management
- `POST /api/documents`: Upload a document
- `GET /api/documents/{id}`: Get document details
- `DELETE /api/documents/{id}`: Delete document
- `GET /api/experts/{id}/documents`: Get expert's documents

### Engagement Management
- `POST /api/engagements`: Create an engagement
- `GET /api/engagements/{id}`: Get engagement details
- `PUT /api/engagements/{id}`: Update engagement
- `DELETE /api/engagements/{id}`: Delete engagement
- `GET /api/experts/{id}/engagements`: Get expert's engagements

### AI Integration
- `POST /api/ai/generate-profile`: Generate an expert profile
- `POST /api/ai/suggest-isced`: Suggest ISCED classification
- `POST /api/ai/extract-skills`: Extract skills from document

### Statistics
- `GET /api/statistics`: Get all statistics
- `GET /api/statistics/nationality`: Get nationality distribution
- `GET /api/statistics/isced`: Get ISCED field distribution
- `GET /api/statistics/engagements`: Get engagement statistics
- `GET /api/statistics/growth`: Get growth statistics

## Database Schema

The system uses SQLite with the following key tables:

- `experts`: Expert profiles and details
- `expert_documents`: Documents uploaded for experts
- `expert_engagements`: Records of expert assignments
- `ai_analysis_results`: Results from AI processing
- `isced_levels`: Education levels (ISCED classification)
- `isced_fields`: Education fields (ISCED classification)
- `system_statistics`: Cached statistics for performance

## License

This project is licensed under the MIT License.