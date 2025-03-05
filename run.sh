#!/bin/bash

# ExpertDB Initialization and Run Script

# Color output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting ExpertDB Initialization${NC}"

# Create required directories
echo -e "${YELLOW}Creating necessary directories...${NC}"
mkdir -p ./backend/db/sqlite
mkdir -p ./backend/data/documents/cv
mkdir -p ./backend/data/documents/certificates
mkdir -p ./backend/data/documents/publications
mkdir -p ./backend/logs

# Check if database exists
if [ ! -f ./backend/db/sqlite/expertdb.sqlite ]; then
    echo -e "${YELLOW}Database does not exist. Will be created on first run.${NC}"
else
    echo -e "${GREEN}Database already exists.${NC}"
fi

# Check if CSV file exists
if [ ! -f ./backend/experts.csv ]; then
    echo -e "${RED}Warning: experts.csv file not found in backend directory. Database import will fail.${NC}"
else
    echo -e "${GREEN}Found experts.csv file.${NC}"
    CSV_COUNT=$(wc -l < ./backend/experts.csv)
    echo -e "${YELLOW}CSV contains approximately $((CSV_COUNT-1)) records.${NC}"
fi

# Set environment variables for the application
export ADMIN_EMAIL=${ADMIN_EMAIL:-"admin@expertdb.com"}
export ADMIN_NAME=${ADMIN_NAME:-"Admin User"}
export ADMIN_PASSWORD=${ADMIN_PASSWORD:-"adminpassword"}

# Set logging configuration
export LOG_LEVEL=${LOG_LEVEL:-"INFO"}  # DEBUG, INFO, WARN, ERROR
export LOG_DIR=${LOG_DIR:-"./logs"}

echo -e "${YELLOW}Admin user will be created with email: ${ADMIN_EMAIL}${NC}"
echo -e "${YELLOW}You can change admin settings by setting environment variables ADMIN_EMAIL, ADMIN_NAME, and ADMIN_PASSWORD${NC}"
echo -e "${YELLOW}Logging configured with level ${LOG_LEVEL} to directory ${LOG_DIR}${NC}"
echo -e "${YELLOW}You can change logging settings with LOG_LEVEL and LOG_DIR environment variables${NC}"

# Database migration check
echo -e "${YELLOW}Checking database migrations...${NC}"
if [ -d ./backend/db/migrations/sqlite ]; then
    MIGRATION_COUNT=$(ls -1 ./backend/db/migrations/sqlite/*.sql 2>/dev/null | wc -l)
    echo -e "${GREEN}Found ${MIGRATION_COUNT} migration files that will be applied automatically${NC}"
else
    echo -e "${RED}Warning: No database migrations found in ./backend/db/migrations/sqlite${NC}"
fi

# Expert import
echo -e "${YELLOW}CSV import will be performed if the database is new${NC}"
echo -e "${YELLOW}Importing from ./backend/experts.csv${NC}"

# Check if Docker is available
if command -v docker >/dev/null 2>&1 && command -v docker-compose >/dev/null 2>&1; then
    echo -e "${GREEN}Docker and Docker Compose detected. Starting services with Docker...${NC}"
    
    # Check if node_modules exists in frontend directory for Docker build efficiency
    if [ ! -d "./frontend/node_modules" ]; then
        echo -e "${YELLOW}Frontend dependencies not found. Installing dependencies before Docker build...${NC}"
        cd frontend && npm install
        cd ..
    fi
    
    echo -e "${GREEN}Starting Docker services...${NC}"
    echo -e "${YELLOW}Frontend will be available at http://localhost:3000${NC}"
    echo -e "${YELLOW}Backend will be available at http://localhost:8080${NC}"
    docker-compose up --build
else
    echo -e "${YELLOW}Docker not detected or Docker Compose not available. Starting services manually...${NC}"
    
    # Check if we need to start the frontend
    if [ -d "./frontend" ] && [ -f "./frontend/package.json" ]; then
        # Check if npm is installed
        if command -v npm >/dev/null 2>&1; then
            echo -e "${YELLOW}Starting frontend and backend services in separate processes...${NC}"
            
            # Start backend in background
            echo -e "${GREEN}Starting backend server...${NC}"
            cd backend && go run . &
            BACKEND_PID=$!
            
            # Go back to root directory
            cd ..
            
            # Check if node_modules exists in frontend directory
            if [ ! -d "./frontend/node_modules" ]; then
                echo -e "${YELLOW}Frontend dependencies not found. Installing dependencies...${NC}"
                cd frontend && npm install
                cd ..
            fi
            
            # Start frontend
            echo -e "${GREEN}Starting frontend server...${NC}"
            cd frontend && npm run dev &
            FRONTEND_PID=$!
            
            echo -e "${GREEN}Services started successfully!${NC}"
            echo -e "${YELLOW}Backend running at http://localhost:8080${NC}"
            echo -e "${YELLOW}Frontend running at http://localhost:3000${NC}"
            echo -e "${YELLOW}Please access the application at http://localhost:3000${NC}"
            echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
            
            # Wait for Ctrl+C
            trap "kill $BACKEND_PID $FRONTEND_PID; exit" INT
            wait
        else
            echo -e "${RED}Error: npm not found. Cannot start frontend.${NC}"
            echo -e "${YELLOW}Starting backend server only...${NC}"
            echo -e "${YELLOW}You'll need to install dependencies and start the frontend manually:${NC}"
            echo -e "${YELLOW}  cd frontend && npm install && npm run dev${NC}"
            echo -e "${YELLOW}Then access the application at http://localhost:3000${NC}"
            
            # Start backend
            cd backend && go run .
        fi
    else
        echo -e "${YELLOW}Frontend directory or package.json not found. Starting backend server only...${NC}"
        echo -e "${YELLOW}You'll need to start the frontend manually if needed.${NC}"
        
        # Start backend
        cd backend && go run .
    fi
fi