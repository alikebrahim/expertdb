FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go.mod and go.sum
COPY go.mod go.sum ./
# Force Go to respect our go.mod version
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled for SQLite support
RUN CGO_ENABLED=1 GOOS=linux go build -o expertdb

# Final stage
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates sqlite tzdata

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/expertdb /app/

# Copy migrations and other necessary files
COPY --from=builder /app/db/migrations /app/db/migrations

# Create directories for data persistence
RUN mkdir -p /app/data/documents

# Default environment variables
ENV PORT=8080
ENV DB_PATH=/app/data/expertdb.sqlite
ENV UPLOAD_PATH=/app/data/documents
ENV CORS_ALLOWED_ORIGINS=http://localhost:3000

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["/app/expertdb"]