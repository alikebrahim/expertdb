# ExpertDB Docker Deployment Guide

This guide provides instructions for deploying ExpertDB using Docker, making the deployment process straightforward and consistent across different environments.

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone and navigate to the project:**
   ```bash
   git clone <repository-url>
   cd expertdb
   ```

2. **Deploy with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

3. **Access the application:**
   - Open your browser and go to `http://localhost:8080`
   - Default admin credentials:
     - Email: `admin@expertdb.com`
     - Password: `adminpassword`

### Using Docker directly

1. **Build the image:**
   ```bash
   docker build -t expertdb .
   ```

2. **Run the container:**
   ```bash
   docker run -d \
     --name expertdb \
     -p 8080:8080 \
     -v expertdb_data:/app/data \
     -v expertdb_logs:/app/logs \
     -v expertdb_db:/app/db/sqlite \
     expertdb
   ```

## Configuration

### Environment Variables

You can customize the deployment by setting these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `DB_PATH` | SQLite database file path | `/app/db/sqlite/main.db` |
| `UPLOAD_PATH` | Document upload directory | `/app/data/documents` |
| `LOG_DIR` | Log files directory | `/app/logs` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `CORS_ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | `*` |
| `ADMIN_EMAIL` | Default admin email | `admin@expertdb.com` |
| `ADMIN_NAME` | Default admin name | `Admin User` |
| `ADMIN_PASSWORD` | Default admin password | `adminpassword` |

### Custom Configuration Example

Create a `.env` file or modify the `docker-compose.yml`:

```yaml
environment:
  - PORT=3000
  - LOG_LEVEL=debug
  - ADMIN_EMAIL=admin@yourcompany.com
  - ADMIN_PASSWORD=your-secure-password
  - CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourcompany.com
```

## Data Persistence

The Docker setup creates three volumes for data persistence:

- **expertdb_data**: Document uploads and user files
- **expertdb_logs**: Application logs
- **expertdb_db**: SQLite database files

### Backup

To backup your data:

```bash
# Backup database
docker cp expertdb:/app/db/sqlite/main.db ./backup_main.db

# Backup uploads
docker run --rm -v expertdb_data:/data -v $(pwd):/backup alpine tar czf /backup/expertdb_data.tar.gz -C /data .

# Backup logs
docker run --rm -v expertdb_logs:/logs -v $(pwd):/backup alpine tar czf /backup/expertdb_logs.tar.gz -C /logs .
```

### Restore

To restore from backup:

```bash
# Restore database
docker cp ./backup_main.db expertdb:/app/db/sqlite/main.db

# Restore uploads
docker run --rm -v expertdb_data:/data -v $(pwd):/backup alpine tar xzf /backup/expertdb_data.tar.gz -C /data

# Restart container to pick up changes
docker-compose restart
```

## Production Deployment

### Security Considerations

1. **Change default admin password:**
   ```bash
   docker-compose exec expertdb ./main -change-admin-password
   ```

2. **Use environment-specific configurations:**
   - Create environment-specific `.env` files
   - Use secrets management for sensitive data
   - Configure proper CORS origins

3. **Use a reverse proxy:**
   ```nginx
   server {
       listen 80;
       server_name expertdb.yourcompany.com;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```

### Resource Requirements

**Minimum Requirements:**
- CPU: 1 core
- RAM: 512MB
- Storage: 10GB (for database and document uploads)

**Recommended for Production:**
- CPU: 2 cores
- RAM: 1GB
- Storage: 50GB SSD

### Monitoring

The container includes a health check that monitors the application status:

```bash
# Check container health
docker ps
docker-compose ps

# View logs
docker-compose logs -f expertdb

# Check health endpoint
curl http://localhost:8080/api/health
```

## Development

For development with hot reload:

```bash
# Use bind mounts for development
docker run -d \
  --name expertdb-dev \
  -p 8080:8080 \
  -v $(pwd)/backend:/app/backend \
  -v $(pwd)/frontend:/app/frontend \
  expertdb
```

## Troubleshooting

### Common Issues

1. **Port already in use:**
   ```bash
   # Change the port in docker-compose.yml or use a different port
   docker-compose up -d
   ```

2. **Permission issues:**
   ```bash
   # Check container logs
   docker-compose logs expertdb
   
   # Reset permissions
   docker-compose exec expertdb chown -R expertdb:expertdb /app
   ```

3. **Database corruption:**
   ```bash
   # Stop container
   docker-compose down
   
   # Remove database volume
   docker volume rm expertdb_expertdb_db
   
   # Restart (will create fresh database)
   docker-compose up -d
   ```

### Getting Support

- Check application logs: `docker-compose logs expertdb`
- Verify health status: `curl http://localhost:8080/api/health`
- Check resource usage: `docker stats expertdb`

## Scaling

For higher loads, consider:

1. **Load balancing multiple instances:**
   ```yaml
   version: '3.8'
   services:
     expertdb:
       build: .
       deploy:
         replicas: 3
       # ... rest of configuration
   ```

2. **Using external database:**
   - Configure external SQLite or migrate to PostgreSQL
   - Update `DB_PATH` environment variable

3. **Separate file storage:**
   - Use external file storage service
   - Mount shared storage volumes
