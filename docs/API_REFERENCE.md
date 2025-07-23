# ExpertDB API Reference - Overview

**Date**: July 22, 2025  
**Version**: 1.5  
**Context**: ExpertDB is a lightweight internal tool for managing a database of experts, designed for a department with 10-12 users and a maximum of 2000 database entries over 5 years. The tool operates on an intranet, with security handled organizationally, prioritizing simplicity, maintainability, and clear error messaging over complex scalability or security measures.

**Backend Technology**: Built in Go, uses SQLite as the database, and provides a RESTful API with JSON payloads, JWT authentication, and permissive CORS settings (`*`).

## Table of Contents

1. [Overview](#overview)
2. [API Architecture](#api-architecture)
3. [Authentication & Authorization](#authentication--authorization)
4. [Standard Response Format](#standard-response-format)
5. [Error Handling](#error-handling)
6. [API Documentation by Category](#api-documentation-by-category)
7. [Recent Updates](#recent-updates)

## Overview

The ExpertDB backend provides a RESTful API for managing expert profiles, requests, users, documents, engagements, phase planning, statistics, and backups. This document serves as the main reference point, with detailed endpoint documentation organized into focused sub-documents for better readability and LLM ingestion.

### Key Features
- JWT-based authentication with 24-hour token expiration
- Role-based access control (super_user, admin, user) with contextual elevations
- RESTful design with consistent JSON responses
- SQLite database with transaction support
- Comprehensive logging and error handling
- CSV data import/export capabilities

## API Architecture

### General Principles
- **RESTful Design**: Standard HTTP methods (GET, POST, PUT, DELETE)
- **JSON Payloads**: All requests and responses use JSON format (except file uploads)
- **Authentication**: Bearer token authentication via JWT
- **Versioning**: API version included in documentation headers
- **CORS**: Permissive settings suitable for intranet deployment

### Database
- **Engine**: SQLite with WAL mode enabled
- **Schema**: Managed via migrations in `db/migrations/sqlite`
- **Indexes**: Applied for common filter fields (nationality, general_area, etc.)
- **Transactions**: Used for data consistency in multi-step operations

### Security
- **Authentication**: JWT tokens with HS256 signing
- **Password Storage**: Bcrypt hashing with cost factor 10
- **File Storage**: Server-side storage with unique naming
- **Access Control**: Role-based with contextual elevations

## Authentication & Authorization

All API endpoints except `/api/auth/login` and `/api/health` require authentication via JWT token in the `Authorization: Bearer <token>` header.

### User Roles
- **super_user**: Complete system access, can create admin users
- **admin**: Full system access, manages all resources and users
- **user**: Limited access, can submit requests and view data

### Contextual Elevations
- **Planner**: Temporary privilege to propose experts for specific applications
- **Manager**: Temporary privilege to rate experts for specific applications

For detailed authentication documentation, see [API_REFERENCE_AUTH.md](./API_REFERENCE_AUTH.md).

## Standard Response Format

All endpoints return responses in a consistent format:

```json
{
  "success": boolean,    // true for successful requests
  "message": "string",   // optional success message  
  "data": <object>       // optional response data
}
```

### HTTP Status Codes
- `200 OK`: Successful GET, PUT, DELETE
- `201 Created`: Successful POST
- `400 Bad Request`: Invalid payload/parameters
- `401 Unauthorized`: Invalid/missing token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate resource
- `500 Internal Server Error`: Unexpected errors

## Error Handling

Error responses follow a consistent format:

```json
{
  "error": "string"    // Error message describing the issue
}
```

Common error scenarios:
- Missing required fields
- Invalid data types or formats
- Authentication failures
- Permission denied
- Resource conflicts
- Database errors

## API Documentation by Category

The API endpoints are organized into logical groups. Click on any category below for detailed endpoint documentation:

### 1. [Authentication & Authorization](./API_REFERENCE_AUTH.md)
- User login and JWT token generation
- Token validation and security

### 2. [User Management & Role Assignments](./API_REFERENCE_USERS.md)
- User CRUD operations
- Role management
- Contextual elevation assignments (planner/manager)
- User profile access

### 3. [Expert Management](./API_REFERENCE_EXPERTS.md)
- Expert CRUD operations
- Advanced filtering and sorting
- General area management
- Specialized area listings
- Multi-value filter support (v1.5)

### 4. [Expert Requests & Edit Requests](./API_REFERENCE_REQUESTS.md)
- Expert creation request workflow
- Request approval/rejection
- Batch approval operations
- Edit request system (planned enhancement)
- Specialized area suggestions

### 5. [Document Management](./API_REFERENCE_DOCUMENTS.md)
- Document upload (CV, approval documents)
- Document retrieval and metadata
- Document deletion
- Expert-document associations

### 6. [Engagement Management](./API_REFERENCE_ENGAGEMENTS.md)
- Engagement CRUD operations
- Expert-engagement tracking
- CSV/JSON import functionality
- Filtering by type and status
- Legacy system for historical records

### 7. [Phase Planning & Applications](./API_REFERENCE_PHASES.md)
- Phase creation and management
- Application assignment workflow
- Expert proposal system
- Review and approval process
- Rating functionality
- Application listing with filters

### 8. [System Endpoints](./API_REFERENCE_SYSTEM.md)
- Statistics dashboard (consolidated endpoint)
- CSV backup generation
- System health check
- Performance monitoring

## Recent Updates

### Version 1.5 (July 21, 2025) - Enhanced Multi-Value Filtering
- **Fixed Institution Search Bug**: Institution filter now properly searches institution/affiliation field
- **Standardized Parameter Names**: Removed `by_*` prefixes for cleaner API
- **Multi-Value Support**: Comma-separated values for filters (e.g., `role=validator,evaluator`)
- **Improved Filter Logic**: OR within parameters, AND between parameters

### Version 1.4 - Statistics Consolidation
- **Single Statistics Endpoint**: Consolidated 5 endpoints into `/api/statistics`
- **All-User Access**: Statistics now available to all authenticated users
- **Performance**: Single query for all statistics improves response time

### Migration Notes
**Breaking Changes in v1.5**:
- Old: `by_role=validator&by_general_area=3`
- New: `role=validator&general_area=3`
- Old: `name=University` (incorrect behavior)
- New: `institution=University` (correct behavior)

## Additional Resources

### Implementation Details
- **Backend Code**: `internal/api/handlers/`
- **Database Schema**: `db/migrations/sqlite/`
- **Authentication**: `internal/auth/`
- **Storage Layer**: `internal/storage/`

### Testing
- API test script: `test_api.sh`
- Test credentials:
  - Admin: `admin@expertdb.com` / `adminpassword`
  - User: `user@expertdb.com` / `userpassword`

### Logging
- Location: `./logs` directory
- Format: JSON with request/response details
- Rotation: Manual cleanup required

---

For questions or issues, consult the specific API category documentation linked above or review the source code implementation.