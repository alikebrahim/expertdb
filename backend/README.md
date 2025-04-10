# ExpertDB Backend

## Overview
Refactored to remove AI integration and streamline the codebase.

A Go-based backend for the ExpertDB system that manages expert profiles, documents, engagements, and provides statistics.

## Features

- **Expert Management**: CRUD operations for expert profiles
- **Document Handling**: Upload, store, and manage expert documents (CVs, certificates, etc.)
- **Engagement Tracking**: Record and monitor expert assignments and activities
- **ISCED Classification**: Support for international standard classification of education
- **Statistics**: Comprehensive metrics on experts, engagements, and system usage

## Architecture

The system is designed with a clean architecture:

- **Backend (this repository)**: Handles data storage, business logic, and API serving
- **Frontend** (separate repository): UI for user interaction

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
   python py_import.py
   ```
   Note: This requires a valid experts.csv file in the project directory.

4. Run the server:
   ```bash
   go run *.go
   ```


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
- `isced_levels`: Education levels (ISCED classification)
- `isced_fields`: Education fields (ISCED classification)
- `system_statistics`: Cached statistics for performance

## License

This project is licensed under the MIT License.