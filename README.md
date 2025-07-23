# ExpertDB Backend

This is the backend service for ExpertDB, a system for managing expert profiles, requests, and engagements.

## Project Structure

The project follows a clean, modular structure using Go best practices:

```
├── cmd/
│   └── server/       # Application entry point
├── internal/         # Private application code
│   ├── api/          # HTTP API handlers
│   ├── auth/         # Authentication and authorization
│   ├── config/       # Configuration management
│   ├── domain/       # Core business entities
│   ├── documents/    # Document handling service
│   ├── logger/       # Logging functionality
│   └── storage/      # Database operations
│       └── sqlite/   # SQLite implementation
└── db/               # Database migrations and schema
    └── migrations/   # Migration files
```

## Architecture

The application follows a layered architecture:

1. **Domain Layer** (`internal/domain`)
   - Core entities and validation
   - Error definitions

2. **Storage Layer** (`internal/storage`)
   - Database interface
   - SQLite implementation

3. **Service Layer** (`internal/documents`, etc.)
   - Logic for documents, etc.

4. **API Layer** (`internal/api`)
   - HTTP handlers and routing
   - Request/response handling

5. **Cross-cutting Concerns**
   - `internal/auth`: Authentication and authorization
   - `internal/config`: Configuration management
   - `internal/logger`: Logging

## Running the Application

```bash
# Set environment variables (or use defaults)
export PORT=8080
export DB_PATH=./db/sqlite/main.db
export UPLOAD_PATH=./data/documents
export CORS_ALLOWED_ORIGINS=*
export LOG_LEVEL=info

# Run the application
go run cmd/server/main.go
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `DB_PATH` | Path to SQLite database | `./db/sqlite/main.db` |
| `UPLOAD_PATH` | Directory for document uploads | `./data/documents` |
| `CORS_ALLOWED_ORIGINS` | CORS allowed origins | `*` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `LOG_DIR` | Directory for log files | `./logs` |
| `ADMIN_EMAIL` | Default admin email | `admin@expertdb.com` |
| `ADMIN_NAME` | Default admin name | `Admin User` |
| `ADMIN_PASSWORD` | Default admin password | `adminpassword` |
