# Multi-stage Docker build for ExpertDB
# This dockerfile builds both frontend and backend, serving a complete application

# Stage 1: Build Frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Install dependencies (including dev dependencies for build)
RUN npm ci

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.22-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app/backend

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Stage 3: Final Runtime Image
FROM alpine:3.19

# Install runtime dependencies including Go for goose
RUN apk add --no-cache \
    sqlite \
    ca-certificates \
    tzdata \
    go \
    && rm -rf /var/cache/apk/*

# Set Go environment variables
ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=$PATH:$GOBIN

# Create go directory and install goose (compatible with Go 1.21)
RUN mkdir -p /go && \
    go install github.com/pressly/goose/v3/cmd/goose@v3.20.0

# Create non-root user
RUN addgroup -g 1001 -S expertdb && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G expertdb -g expertdb expertdb

WORKDIR /app

# Copy built backend binary
COPY --from=backend-builder /app/backend/main ./

# Copy built frontend assets
COPY --from=frontend-builder /app/frontend/dist ./public

# Copy database migrations (if any)
COPY backend/db/ ./db/

# Copy entrypoint script
COPY docker-entrypoint.sh ./
RUN chmod +x ./docker-entrypoint.sh

# Create necessary directories and set permissions
RUN mkdir -p /app/data/documents /app/logs /app/db/sqlite && \
    chown -R expertdb:expertdb /app /go

# Switch to non-root user
USER expertdb

# Expose port
EXPOSE 8080

# Environment variables with defaults
ENV PORT=8080
ENV DB_PATH=/app/db/sqlite/main.db
ENV UPLOAD_PATH=/app/data/documents
ENV LOG_DIR=/app/logs
ENV LOG_LEVEL=info
ENV CORS_ALLOWED_ORIGINS=*
ENV ADMIN_EMAIL=admin@expertdb.com
ENV ADMIN_NAME="Admin User"
# ADMIN_PASSWORD should be set via docker run -e or docker-compose.yml
# Default is only for development - change this in production
ARG ADMIN_PASSWORD=changeme
ENV ADMIN_PASSWORD=${ADMIN_PASSWORD}

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/api/health || exit 1

# Create volumes for persistent data
VOLUME ["/app/data", "/app/logs", "/app/db"]

# Start the application with migrations
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["./main"]
