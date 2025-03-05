# Docker Setup for ExpertDB

This document provides instructions for running the ExpertDB system using Docker and Docker Compose.

## Prerequisites

- Docker Engine (version 19.03.0+)
- Docker Compose (version 1.27.0+)

## Getting Started

### 1. Initialize Data Directory Structure

Before starting the containers, create the necessary directory structure:

```bash
./scripts/init_data_dirs.sh
```

### 2. Backend Only Setup

To run just the backend service (for development):

```bash
docker-compose up -d
```

This will start the backend service on port 8080.

### 3. Full System Setup

To run the full system with placeholders for frontend and AI services:

```bash
docker-compose -f docker-compose.full.yml up -d
```

This will start:
- Backend service on port 8080
- Frontend placeholder on port 3000
- AI service placeholder on port 9000

### 4. Database Migration

To initialize or update the database schema, run:

```bash
docker exec expertdb-backend sh -c "cd /app && goose -dir db/migrations/sqlite sqlite3 /app/data/expertdb.sqlite up"
```

### 5. Importing Test Data

To import sample data for testing:

```bash
docker exec expertdb-backend sh -c "cd /app && ./import_csv -csv /app/experts.csv -db /app/data/expertdb.sqlite"
```

## Monitoring and Logs

View logs for all services:

```bash
docker-compose logs -f
```

View logs for a specific service:

```bash
docker-compose logs -f expertdb-backend
```

## Development Workflow

When developing the backend service:

1. Make changes to the code
2. Rebuild the container:
   ```bash
   docker-compose up -d --build
   ```

## Production Considerations

For production deployment:

1. Update the CORS settings to allow only specific origins
2. Enable TLS/SSL using a reverse proxy like Nginx
3. Set appropriate environment variables for security
4. Implement proper authentication
5. Use a more robust database like PostgreSQL instead of SQLite
6. Set up proper logging and monitoring

## Volume Management

The system uses Docker volumes for data persistence:

- Database files: `/app/data/expertdb.sqlite`
- Document storage: `/app/data/documents/`

To backup these volumes:

```bash
docker run --rm -v expertdb_expertdb-data:/source -v $(pwd)/backup:/backup alpine tar -czf /backup/expertdb-data-$(date +%Y%m%d).tar.gz -C /source .
```

## Networking

The services communicate over a Docker network called `expertdb-network`. The backend service expects the AI service to be available at `http://expertdb-ai:9000`.

## Troubleshooting

If you encounter issues:

1. Check container status:
   ```bash
   docker-compose ps
   ```

2. Check container logs:
   ```bash
   docker-compose logs
   ```

3. Inspect network:
   ```bash
   docker network inspect expertdb-network
   ```

4. Restart services:
   ```bash
   docker-compose restart
   ```