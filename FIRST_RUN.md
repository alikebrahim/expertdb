# ExpertDB - First Run Guide

This document provides instructions for the first run and initialization of the ExpertDB system. It covers the setup, initialization, and testing process.

## System Requirements

- Go 1.18+ installed
- SQLite 3
- Bash shell

## Initialization Process

The system has been set up with an automated initialization process through the `run.sh` script. This script performs the following actions:

1. Creates necessary directories for:
   - Database files
   - Document storage
   - Log files

2. Validates the presence of required files:
   - Experts CSV file for initial data import
   - Database migration scripts

3. Sets up environment variables with defaults:
   - Admin user credentials
   - Logging configuration
   - System paths

4. Initializes the database:
   - Creates SQLite database if not exists
   - Applies migration scripts automatically
   - Imports expert data from CSV file

5. Creates default admin user if not exists

## First Run Instructions

1. **Navigate to the project root directory**:
   ```bash
   cd /path/to/expertdb_new
   ```

2. **Run the initialization script**:
   ```bash
   ./run.sh
   ```

3. **Optional: Configure admin user**
   You can set environment variables before running to customize the admin account:
   ```bash
   export ADMIN_EMAIL="custom.admin@example.com"
   export ADMIN_NAME="Custom Admin Name"
   export ADMIN_PASSWORD="securepassword"
   ./run.sh
   ```

4. **Optional: Configure logging level**
   You can adjust the logging verbosity:
   ```bash
   export LOG_LEVEL="DEBUG" # Options: DEBUG, INFO, WARN, ERROR
   export LOG_DIR="./custom_logs"
   ./run.sh
   ```

## Post-Initialization

After initialization, the system will have:

1. A running backend API server
2. A configured database with:
   - User tables with default admin user
   - Expert tables populated from CSV
   - All required schema objects

3. Log files stored in the `logs` directory:
   - Application logs with timestamps
   - Request/response logs
   - Error logs

## Testing the Setup

To verify the system is properly initialized:

1. **Test Admin Login**:
   ```bash
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@expertdb.com","password":"adminpassword"}'
   ```
   This should return a JSON response with user details and a JWT token.

2. **View Expert Data**:
   ```bash
   curl http://localhost:8080/api/experts?limit=5
   ```
   This should return the first 5 experts from the database.

3. **Check Statistics**:
   ```bash
   # Get JWT token first from login response above
   curl http://localhost:8080/api/statistics \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```
   This should return system statistics.

## Common Issues

1. **Database Initialization Failure**:
   - Check if SQLite is installed
   - Verify permissions on the db directory
   - Check logs for specific errors

2. **Admin User Creation Failure**:
   - Check logs for password hashing errors
   - Verify database permissions

3. **CSV Import Failure**:
   - Ensure the CSV file is in the correct location (./backend/experts.csv)
   - Check the CSV format matches expected schema

4. **Permission Issues**:
   - Ensure the run.sh script is executable (`chmod +x run.sh`)
   - Verify write permissions in all created directories

## Implemented Features

1. **Authentication System**:
   - JWT-based authentication
   - Role-based access control (admin/user)
   - Secure password hashing with bcrypt
   - Token validation and middleware

2. **Logging System**:
   - Multi-level logging (DEBUG, INFO, WARN, ERROR)
   - File and console outputs
   - Colored console logs
   - Request/response logging
   - Structured log format with timestamps and source location

3. **User Management**:
   - Admin user auto-creation
   - User CRUD operations (admin only)
   - Role management
   - No public registration (admin-controlled)

4. **Database Initialization**:
   - Automatic directory creation
   - Schema migrations
   - Initial data import from CSV
   - Data validation

5. **Error Handling**:
   - Comprehensive error logging
   - HTTP status code mapping
   - Detailed error messages
   - Client-safe error responses

## Next Steps

After the initial setup:

1. Complete the frontend implementation
2. Run integration tests between frontend and backend
3. Set up production deployment configuration
4. Implement database backup procedures
5. Configure SSL for production