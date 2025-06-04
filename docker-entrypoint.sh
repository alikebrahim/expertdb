#!/bin/sh
set -e

echo "Starting ExpertDB container..."

# Set GOPATH and GOBIN
export GOPATH=/go
export GOBIN=/go/bin
export PATH=$PATH:$GOBIN

# Create go directory if it doesn't exist
mkdir -p /go

# Run database migrations
echo "Running database migrations..."
cd /app/db/migrations/sqlite
/go/bin/goose sqlite /app/db/sqlite/main.db up

# Return to app directory
cd /app

echo "Migrations completed. Starting application..."

# Start the main application
exec "$@"