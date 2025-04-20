# ExpertDB API Endpoints Documentation

**Date**: April 17, 2025  
**Version**: 1.0  
**Context**:  
ExpertDB is a small internal tool for managing a database of experts, designed for a department with 10-12 users and a maximum of 1200 database entries over 5 years. The tool is not exposed to the internet, and security is handled organizationally, so the focus is on simplicity, maintainability, and clear error messaging rather than high scalability or robust security measures. The backend is built in Go, uses SQLite as the database, and provides a RESTful API with JSON payloads, JWT authentication, and permissive CORS settings (`*`).

**Purpose**:  
This document provides a detailed reference for all API endpoints, including their functionality, request/response structures, and implementation notes. It serves as a guide for developers and users within the department to interact with the ExpertDB system effectively.

## Table of Contents
1. [Overview](#overview)
2. [General Notes](#general-notes)
3. [Authentication Endpoints](#authentication-endpoints)
   - [POST /api/auth/login](#post-apiauthlogin)
4. [User Management Endpoints](#user-management-endpoints)
   - [POST /api/users](#post-apiusers)
   - [DELETE /api/users/{id}](#delete-apiusersid)
5. [Expert Management Endpoints](#expert-management-endpoints)
   - [POST /api/experts](#post-apiexperts)
   - [GET /api/experts](#get-apiexperts)
   - [GET /api/experts/{id}](#get-apiexpertsid)
   - [PUT /api/experts/{id}](#put-apiexpertsid)
   - [DELETE /api/experts/{id}](#delete-apiexpertsid)
   - [GET /api/expert/areas](#get-apiexpertareas)
6. [Expert Request Management Endpoints](#expert-request-management-endpoints)
   - [POST /api/expert-requests](#post-apiexpert-requests)
   - [GET /api/expert-requests](#get-apiexpert-requests)
   - [GET /api/expert-requests/{id}](#get-apiexpert-requestsid)
   - [PUT /api/expert-requests/{id}](#put-apiexpert-requestsid)
7. [Document Management Endpoints](#document-management-endpoints)
   - [POST /api/documents](#post-apidocuments)
   - [GET /api/experts/{id}/documents](#get-apiexpertsiddocuments)
   - [GET /api/documents/{id}](#get-apidocumentsid)
   - [DELETE /api/documents/{id}](#delete-apidocumentsid)
8. [Engagement Management Endpoints](#engagement-management-endpoints)
   - [GET /api/expert-engagements](#get-apiexpert-engagements)
9. [Statistics Endpoints](#statistics-endpoints)
   - [GET /api/statistics](#get-apistatistics)
   - [GET /api/statistics/growth](#get-apistatisticsgrowth)
   - [GET /api/statistics/nationality](#get-apistatisticsnationality)
   - [GET /api/statistics/engagements](#get-apistatisticsengagements)

## Overview
The ExpertDB backend provides a RESTful API for managing expert profiles, expert requests, user accounts, documents, engagements, and system statistics. The API is implemented in Go, with endpoints defined in the `internal/api` directory, primarily in `server.go` and the `handlers` subpackage. It uses SQLite (`expertdb.sqlite`) for data storage, JWT for authentication, and the `internal/logger` package for logging requests and responses to `./logs`.

The API is designed for:
- **Small Scale**: Supports 10-12 users and up to 1200 expert entries over 5 years.
- **Internal Use**: Not exposed to the internet, with security managed organizationally.
- **Simplicity**: Prioritizes clear error messages and straightforward CRUD operations over complex optimizations.
- **Modularity**: Follows a layered architecture (Domain, Storage, Service, API) as outlined in `backend/README.md`.

Endpoints are grouped into categories for authentication, user management, expert management, expert requests, documents, engagements, and statistics. Most endpoints require authentication via a JWT token, with admin-only endpoints restricted to users with the `admin` role.

## General Notes
- **Authentication**: All endpoints (except `/api/auth/login` and health checks) require a JWT token in the `Authorization: Bearer <token>` header, obtained via `/api/auth/login`. Admin-only endpoints are protected by middleware in `internal/auth/middleware.go`.
- **Access Levels**:
  - **Public**: Only `/api/auth/login` and health check endpoints are accessible without authentication
  - **User**: Authenticated users have read-only access to experts, expert areas, documents, and engagements
  - **Scheduler**: Can manage engagements (create, update, delete)
  - **Admin**: Can manage experts, expert requests, documents, and users
  - **Super User**: Has full system access including statistics and user deletion
- **CORS**: Configured to allow all origins (`*`), suitable for internal use but may need adjustment if exposed externally.
- **Error Handling**: Errors return JSON with an `error` field. Improvements suggested in `ERRORS.md` include specific messages and aggregated validation errors for clarity.
- **Database**: Uses SQLite with schema defined in `db/migrations/sqlite`. The small scale ensures SQLite is sufficient.
- **Logging**: Requests and responses are logged to `./logs` with details like HTTP status, headers, and payloads.
- **Testing**: The `test_api.sh` script validates endpoints, covering happy paths and edge cases (e.g., invalid payloads).
- **Payload Validation**: Required fields are enforced, with defaults (e.g., `pending` status for expert requests) applied where applicable.
- **HTTP Status Codes**:
  - `200 OK`: Successful GET, PUT, DELETE.
  - `201 Created`: Successful POST.
  - `400 Bad Request`: Invalid payload or parameters.
  - `401 Unauthorized`: Invalid or missing token.
  - `403 Forbidden`: Insufficient permissions (e.g., non-admin access).
  - `404 Not Found`: Resource not found.
  - `409 Conflict`: Duplicate resource (e.g., email or expert ID).
  - `500 Internal Server Error`: Unexpected server errors.

## Authentication Endpoints

### POST /api/auth/login
- **Purpose**: Authenticates a user and returns a JWT token for session management.
- **Method**: POST
- **Path**: `/api/auth/login`
- **Request Payload**:
  ```json
  {
    "email": "string",    // Required: User's email (e.g., "admin@expertdb.com")
    "password": "string"  // Required: User's password (e.g., "adminpassword")
  }
  ```
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "user": {
        "id": int,           // User ID (e.g., 1)
        "name": "string",    // User name (e.g., "Admin User")
        "email": "string",   // User email (e.g., "admin@expertdb.com")
        "role": "string",    // Role ("admin" or "user")
        "isActive": boolean, // Active status (e.g., true)
        "createdAt": "string", // ISO 8601 timestamp (e.g., "2025-04-10T10:04:59.744473095+03:00")
        "lastLogin": "string"  // ISO 8601 timestamp (e.g., "2025-04-17T10:13:09.703248012+03:00")
      },
      "token": "string"      // JWT token (e.g., "eyJhbGciOiJIUzI1NiIs...")
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"     // e.g., "Invalid request payload"
    }
    ```
  - **Error (401 Unauthorized)**:
    ```json
    {
      "error": "string"     // e.g., "Invalid credentials"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/auth.go:HandleLogin`
  - Validates email/password, checks user in `users` table, verifies password hash using `golang.org/x/crypto`.
  - Generates JWT via `internal/auth/jwt.go`.
- **Notes**:
  - Used by admin and regular users.
  - Logs successful logins (e.g., "User logged in successfully: admin@expertdb.com").
  - Token is required for all other endpoints except this one.

## User Management Endpoints

### POST /api/users
- **Purpose**: Creates a new user account (admin-only).
- **Method**: POST
- **Path**: `/api/users`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:
  ```json
  {
    "name": "string",     // Required: User name (e.g., "Test User 1744873989")
    "email": "string",    // Required: User email (e.g., "testuser1744873989@example.com")
    "password": "string", // Required: Password (e.g., "password123")
    "role": "string",     // Required: Role ("admin" or "user")
    "isActive": boolean   // Required: Active status (e.g., true)
  }
  ```
- **Response Payload**:
  - **Success (201 Created)**:
    ```json
    {
      "id": int,           // New user ID (e.g., 18)
      "success": boolean,  // true
      "message": "string"  // e.g., "User created successfully"
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"    // e.g., "Invalid request payload"
    }
    ```
  - **Error (409 Conflict)**:
    ```json
    {
      "error": "string"    // e.g., "Email already exists"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/user.go`
  - Validates input, checks email uniqueness, hashes password using `internal/auth/password.go`.
  - Inserts into `users` table via `internal/storage/sqlite/user.go`.
- **Notes**:
  - Requires admin role, enforced by middleware.
  - Logs creation (e.g., "New user created: testuser1744873989@example.com").

### DELETE /api/users/{id}
- **Purpose**: Deletes a user by ID (admin-only).
- **Method**: DELETE
- **Path**: `/api/users/{id}` (e.g., `/api/users/18`)
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "success": boolean,  // true
      "message": "string"  // e.g., "User deleted successfully"
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"    // e.g., "User not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/user.go`
  - Deletes from `users` table via `internal/storage/sqlite/user.go`.
- **Notes**:
  - Requires admin role.
  - Logs deletion (e.g., "User deleted: ID 18, Email: testuser1744873989@example.com").

## Expert Management Endpoints

### POST /api/experts
- **Purpose**: Creates a new expert profile (admin-only).
- **Method**: POST
- **Path**: `/api/experts`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:
  ```json
  {
    "name": "string",           // Required: Name (e.g., "Test Expert 1744873989")
    "institution": "string",    // Optional: Institution (e.g., "Test University 1744873989")
    "primaryContact": "string", // Required: Contact (e.g., "expert1744873989@example.com")
    "contactType": "string",    // Required: Contact type (e.g., "email")
    "designation": "string",    // Optional: Designation (e.g., "Professor")
    "isBahraini": boolean,      // Optional: Bahraini status (e.g., true)
    "availability": "string",   // Optional: Availability ("yes" or "no")
    "rating": "string",         // Optional: Rating (e.g., "5")
    "role": "string",           // Required: Role (e.g., "evaluator")
    "employmentType": "string", // Optional: Employment type (e.g., "academic")
    "generalArea": int,         // Required: General area ID (e.g., 1)
    "specializedArea": "string",// Optional: Specialized area (e.g., "Software Engineering")
    "isTrained": boolean,       // Optional: Trained status (e.g., true)
    "isPublished": boolean,     // Optional: Published status (e.g., true)
    "biography": "string",      // Optional: Biography (e.g., "Expert created for testing.")
    "skills": ["string"]        // Optional: Skills (e.g., ["Go", "Testing"])
  }
  ```
- **Response Payload**:
  - **Success (201 Created)**:
    ```json
    {
      "id": int,           // Expert ID (e.g., 459)
      "success": boolean,  // true
      "message": "string"  // e.g., "Expert created successfully"
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"    // e.g., "role is required"
    }
    ```
  - **Error (409 Conflict)**:
    ```json
    {
      "error": "string"    // e.g., "Expert ID already exists"
    }
    ```
  - **Error (500 Internal Server Error)**:
    ```json
    {
      "error": "string"    // e.g., "Failed to create expert: UNIQUE constraint failed: experts.expert_id"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go:HandleCreateExpert`
  - Validates required fields, generates `expert_id` if not provided (format: `EXP-<request_id>-<timestamp>`).
  - Inserts into `experts` table via `internal/storage/sqlite/expert.go`.
- **Notes**:
  - Requires admin role.
  - Logs creation attempts and errors (e.g., "Creating expert: Test Expert 1744873989").
  - Validation failures (e.g., invalid email) return specific errors, but improvements are suggested in `ERRORS.md`.

### GET /api/experts
- **Purpose**: Retrieves a paginated list of experts with optional filters, enhanced sorting, and pagination metadata.
- **Method**: GET
- **Path**: `/api/experts`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Query Parameters**:
  - `limit`: Integer (e.g., 5) – Max number of experts.
  - `offset`: Integer (e.g., 0) – Pagination offset.
  - `sort_by`: String – Sort field with expanded options:
    - `name` (default), `institution`, `role`, `created_at`, `updated_at`, `rating`, `general_area` 
    - New options: `expert_id`, `designation`, `employment_type`, `nationality`, `specialized_area`, `is_bahraini`, `is_available`, `is_published`
    - Also accepts camelCase versions (e.g., `expertId`, `specializedArea`)
  - `sort_order`: String (e.g., "asc") – Sort direction ("asc" or "desc").
  - `name`: String – Filter by name (partial match).
  - `is_available`: String ("true"/"false") – Filter by availability.
  - `role`: String – Filter by role (exact match).
  - `generalArea`: Integer – Filter by general area ID.
  - `by_nationality`: String ("Bahraini"/"non-Bahraini") – Filter by nationality.
  - `by_general_area`: Integer – Filter by general area ID (alternative parameter).
  - `by_specialized_area`: String – Filter by specialized area (partial match).
  - `by_employment_type`: String – Filter by employment type (e.g., "academic").
  - `by_role`: String – Filter by role (alternative parameter).
- **Request Payload**: None
- **Response Headers**:
  - `X-Total-Count`: Total number of experts matching filters
  - `X-Total-Pages`: Total number of pages available
  - `X-Current-Page`: Current page number
  - `X-Page-Size`: Number of items per page
  - `X-Has-Next-Page`: Boolean indicating if there's a next page
  - `X-Has-Prev-Page`: Boolean indicating if there's a previous page
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "experts": [
        {
          "id": int,               // Expert ID (e.g., 440)
          "expertId": "string",    // Unique ID (e.g., "EXP-5-1744706079")
          "name": "string",        // Name
          "designation": "string", // Designation
          "institution": "string", // Institution
          "isBahraini": boolean,   // Bahraini status
          "nationality": "string", // Nationality
          "isAvailable": boolean,  // Availability
          "rating": "string",      // Rating
          "role": "string",        // Role
          "employmentType": "string", // Employment type
          "generalArea": int,      // General area ID
          "generalAreaName": "string", // General area name
          "specializedArea": "string", // Specialized area
          "isTrained": boolean,    // Trained status
          "cvPath": "string",      // CV path
          "phone": "string",       // Phone
          "email": "string",       // Email
          "isPublished": boolean,  // Published status
          "biography": "string",   // Biography
          "createdAt": "string",   // ISO 8601 timestamp
          "updatedAt": "string"    // ISO 8601 timestamp
        }
      ],
      "pagination": {
        "totalCount": int,        // Total number of experts matching filters
        "totalPages": int,        // Total number of pages
        "currentPage": int,       // Current page number
        "pageSize": int,          // Number of items per page
        "hasNextPage": boolean,   // Indicates if there's a next page
        "hasPrevPage": boolean,   // Indicates if there's a previous page
        "hasMore": boolean        // Indicates if there are more results
      }
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"        // e.g., "Invalid query parameters"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Queries `experts` table with pagination, sorting, and filtering.
- **Notes**:
  - Accessible to authenticated users.
  - Enhanced in Phase 3B with improved sorting and pagination metadata.
  - Headers continue to be provided for API clients that rely on them.
  - Adds detailed pagination metadata in the response body.
  - Fields can be sorted in many different ways (e.g., by name, nationality, availability).
  - Combines multiple filters with AND logic.
  - Supports a variety of filters for nationality, area, role, etc.
  - Logs query details (e.g., "Retrieving experts with filters: map[by_nationality:Bahraini sort_by:name sort_order:asc]").

### GET /api/experts/{id}
- **Purpose**: Retrieves details of a specific expert by ID.
- **Method**: GET
- **Path**: `/api/experts/{id}` (e.g., `/api/experts/440`)
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "id": int,               // Expert ID
      "expertId": "string",    // Unique ID
      "name": "string",        // Name
      "designation": "string", // Designation
      "institution": "string", // Institution
      "isBahraini": boolean,   // Bahraini status
      "nationality": "string", // Nationality
      "isAvailable": boolean,  // Availability
      "rating": "string",      // Rating
      "role": "string",        // Role
      "employmentType": "string", // Employment type
      "generalArea": int,      // General area ID
      "generalAreaName": "string", // General area name
      "specializedArea": "string", // Specialized area
      "isTrained": boolean,    // Trained status
      "cvPath": "string",      // CV path
      "phone": "string",       // Phone
      "email": "string",       // Email
      "isPublished": boolean,  // Published status
      "biography": "string",   // Biography
      "createdAt": "string",   // ISO 8601 timestamp
      "updatedAt": "string"    // ISO 8601 timestamp
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"        // e.g., "Expert not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Fetches from `experts` table.
- **Notes**:
  - Accessible to authenticated users.
  - Logs retrieval (e.g., "Successfully retrieved expert: ID: 440").

### PUT /api/experts/{id}
- **Purpose**: Updates an existing expert profile (admin-only).
- **Method**: PUT
- **Path**: `/api/experts/{id}` (e.g., `/api/experts/440`)
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:
  ```json
  {
    "name": "string",           // Optional: Name
    "institution": "string",    // Optional: Institution
    "primaryContact": "string", // Optional: Contact
    "contactType": "string",    // Optional: Contact type
    "designation": "string",    // Optional: Designation
    "isBahraini": boolean,      // Optional: Bahraini status
    "availability": "string",   // Optional: Availability
    "rating": "string",         // Optional: Rating
    "role": "string",           // Optional: Role
    "employmentType": "string", // Optional: Employment type
    "generalArea": int,         // Optional: General area ID
    "specializedArea": "string",// Optional: Specialized area
    "isTrained": boolean,       // Optional: Trained status
    "isPublished": boolean,     // Optional: Published status
    "biography": "string",      // Optional: Biography
    "skills": ["string"]        // Optional: Skills
  }
  ```
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "success": boolean,  // true
      "message": "string"  // e.g., "Expert updated successfully"
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"    // e.g., "Invalid request payload"
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"    // e.g., "Expert not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Updates specified fields in `experts` table.
- **Notes**:
  - Requires admin role.
  - Only provided fields are updated.
  - Logs updates (e.g., "Expert updated successfully: ID: 440").

### DELETE /api/experts/{id}
- **Purpose**: Deletes an expert by ID (admin-only).
- **Method**: DELETE
- **Path**: `/api/experts/{id}` (e.g., `/api/experts/440`)
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "success": boolean,  // true
      "message": "string"  // e.g., "Expert deleted successfully"
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"    // e.g., "Expert not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Deletes from `experts` table, cascading to `expert_documents`.
- **Notes**:
  - Requires admin role.
  - Logs deletion (e.g., "Expert deleted successfully: ID: 440").

### GET /api/expert/areas
- **Purpose**: Retrieves a list of available expert areas.
- **Method**: GET
- **Path**: `/api/expert/areas`
- **Request Headers**:
  - `Authorization: Bearer <token>` (required)
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    [
      {
        "id": int,         // Area ID (e.g., 1)
        "name": "string"   // Area name (e.g., "Art and Design")
      }
    ]
    ```
  - **Error (401 Unauthorized)**:
    ```json
    {
      "error": "string"    // e.g., "Unauthorized"
    }
    ```
  - **Error (500 Internal Server Error)**:
    ```json
    {
      "error": "string"    // e.g., "Failed to retrieve expert areas"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Fetches from `expert_areas` table.
- **Notes**:
  - Requires authentication (previously was public).
  - Returns 34 areas (e.g., "Art and Design", "Information Technology").
  - Logs retrieval (e.g., "Returning 34 expert areas").

## Expert Request Management Endpoints

### POST /api/expert-requests
- **Purpose**: Submits a new expert request for review.
- **Method**: POST
- **Path**: `/api/expert-requests`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Payload**:
  ```json
  {
    "name": "string",           // Required: Name (e.g., "Request Expert 1744873989")
    "designation": "string",    // Optional: Designation (e.g., "Researcher")
    "institution": "string",    // Optional: Institution (e.g., "Request University 1744873989")
    "isBahraini": boolean,      // Optional: Bahraini status (e.g., false)
    "isAvailable": boolean,     // Optional: Availability (e.g., true)
    "rating": "string",         // Optional: Rating (e.g., "4")
    "role": "string",           // Optional: Role (e.g., "reviewer")
    "employmentType": "string", // Optional: Employment type (e.g., "freelance")
    "generalArea": int,         // Required: General area ID (e.g., 1)
    "specializedArea": "string",// Optional: Specialized area (e.g., "Quantum Physics")
    "isTrained": boolean,       // Optional: Trained status (e.g., false)
    "phone": "string",          // Optional: Phone (e.g., "+97311111744873989")
    "email": "string",          // Optional: Email (e.g., "request1744873989@example.com")
    "isPublished": boolean,     // Optional: Published status (e.g., false)
    "biography": "string"       // Optional: Biography (e.g., "Researcher requesting addition.")
  }
  ```
- **Response Payload**:
  - **Success (201 Created)**:
    ```json
    {
      "id": int,               // Request ID (e.g., 26)
      "name": "string",        // Name
      "designation": "string", // Designation
      "institution": "string", // Institution
      "isBahraini": boolean,   // Bahraini status
      "isAvailable": boolean,  // Availability
      "rating": "string",      // Rating
      "role": "string",        // Role
      "employmentType": "string", // Employment type
      "generalArea": int,      // General area ID
      "specializedArea": "string", // Specialized area
      "isTrained": boolean,    // Trained status
      "cvPath": "string",      // CV path
      "phone": "string",       // Phone
      "email": "string",       // Email
      "isPublished": boolean,  // Published status
      "status": "string",      // Status (e.g., "pending")
      "biography": "string",   // Biography
      "createdAt": "string",   // ISO 8601 timestamp
      "reviewedAt": "string"   // ISO 8601 timestamp
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"        // e.g., "name is required"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go:HandleCreateExpertRequest`
  - Validates required fields, sets default `status` to "pending".
  - Inserts into `expert_requests` table.
- **Notes**:
  - Accessible to authenticated users.
  - Logs creation (e.g., "Expert request created successfully: ID: 26").
  - Validation improvements suggested in `ERRORS.md`.

### GET /api/expert-requests
- **Purpose**: Retrieves a paginated list of expert requests.
- **Method**: GET
- **Path**: `/api/expert-requests`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Query Parameters**:
  - `limit`: Integer (e.g., 5)
  - `offset`: Integer (e.g., 0)
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    [
      {
        "id": int,               // Request ID
        "name": "string",        // Name
        "designation": "string", // Designation
        "institution": "string", // Institution
        "isBahraini": boolean,   // Bahraini status
        "isAvailable": boolean,  // Availability
        "rating": "string",      // Rating
        "role": "string",        // Role
        "employmentType": "string", // Employment type
        "generalArea": int,      // General area ID
        "specializedArea": "string", // Specialized area
        "isTrained": boolean,    // Trained status
        "cvPath": "string",      // CV path
        "phone": "string",       // Phone
        "email": "string",       // Email
        "isPublished": boolean,  // Published status
        "status": "string",      // Status
        "biography": "string",   // Biography
        "createdAt": "string",   // ISO 8601 timestamp
        "reviewedAt": "string"   // ISO 8601 timestamp
      }
    ]
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"        // e.g., "Invalid query parameters"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Queries `expert_requests` table with pagination.
- **Notes**:
  - Requires admin role.
  - Logs retrieval (e.g., "Returning 5 expert requests").

### GET /api/expert-requests/{id}
- **Purpose**: Retrieves details of a specific expert request.
- **Method**: GET
- **Path**: `/api/expert-requests/{id}` (e.g., `/api/expert-requests/26`)
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "id": int,               // Request ID
      "name": "string",        // Name
      "designation": "string", // Designation
      "institution": "string", // Institution
      "isBahraini": boolean,   // Bahraini status
      "isAvailable": boolean,  // Availability
      "rating": "string",      // Rating
      "role": "string",        // Role
      "employmentType": "string", // Employment type
      "generalArea": int,      // General area ID
      "specializedArea": "string", // Specialized area
      "isTrained": boolean,    // Trained status
      "cvPath": "string",      // CV path
      "phone": "string",       // Phone
      "email": "string",       // Email
      "isPublished": boolean,  // Published status
      "status": "string",      // Status
      "biography": "string",   // Biography
      "createdAt": "string",   // ISO 8601 timestamp
      "reviewedAt": "string"   // ISO 8601 timestamp
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"        // e.g., "Expert request not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Fetches from `expert_requests` table.
- **Notes**:
  - Requires admin role.
  - Logs retrieval (e.g., "Successfully retrieved expert request: ID: 26").

### PUT /api/expert-requests/{id}
- **Purpose**: Updates an expert request, typically to approve or reject it (admin-only).
- **Method**: PUT
- **Path**: `/api/expert-requests/{id}` (e.g., `/api/expert-requests/26`)
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:
  ```json
  {
    "status": "string",          // Required: Status ("approved" or "rejected")
    "rejectionReason": "string"  // Optional: Reason for rejection (e.g., "Test rejection")
  }
  ```
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "success": boolean,  // true
      "message": "string"  // e.g., "Expert request updated successfully"
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"    // e.g., "Invalid status"
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"    // e.g., "Expert request not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Updates `status` and `reviewedAt` in `expert_requests` table.
  - On approval, creates an expert in `experts` table.
- **Notes**:
  - Requires admin role.
  - Logs updates (e.g., "Expert request updated successfully: ID: 26, Status: approved").
  - Approval generates a unique `expert_id` (e.g., "EXP-26-1744873990").

## Document Management Endpoints

### POST /api/documents
- **Purpose**: Uploads a document (e.g., CV) for an expert (admin-only).
- **Method**: POST
- **Path**: `/api/documents`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: Form-data
  ```text
  file: file           // Required: File (e.g., sample_cv.txt)
  documentType: string // Required: Type (e.g., "cv")
  expertId: int        // Required: Expert ID (e.g., 440)
  ```
- **Response Payload**:
  - **Success (201 Created)**:
    ```json
    {
      "id": int,           // Document ID
      "success": boolean,  // true
      "message": "string"  // e.g., "Document uploaded successfully"
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"    // e.g., "Missing file"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/documents/document_handler.go`
  - Stores file in `UPLOAD_PATH` (default: `./data/documents`).
  - Inserts metadata into `expert_documents` table.
- **Notes**:
  - Requires admin role.
  - Logs upload attempts (e.g., via `test_api.sh`).

### GET /api/experts/{id}/documents
- **Purpose**: Retrieves a list of documents for an expert.
- **Method**: GET
- **Path**: `/api/experts/{id}/documents` (e.g., `/api/experts/440/documents`)
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    [
      {
        "id": int,         // Document ID
        "expertId": int,   // Expert ID
        "documentType": "string", // Type
        "filePath": "string",     // File path
        "createdAt": "string"     // ISO 8601 timestamp
      }
    ]
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"    // e.g., "Expert not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/documents/document_handler.go`
  - Queries `expert_documents` table.
- **Notes**:
  - Accessible to authenticated users.
  - Logs retrieval (e.g., via `test_api.sh`).

### GET /api/documents/{id}
- **Purpose**: Retrieves details of a specific document.
- **Method**: GET
- **Path**: `/api/documents/{id}` (e.g., `/api/documents/1`)
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "id": int,         // Document ID
      "expertId": int,   // Expert ID
      "documentType": "string", // Type
      "filePath": "string",     // File path
      "createdAt": "string"     // ISO 8601 timestamp
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"    // e.g., "Document not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/documents/document_handler.go`
  - Fetches from `expert_documents` table.
- **Notes**:
  - Accessible to authenticated users.
  - Logs retrieval (e.g., via `test_api.sh`).

### DELETE /api/documents/{id}
- **Purpose**: Deletes a document by ID (admin-only).
- **Method**: DELETE
- **Path**: `/api/documents/{id}` (e.g., `/api/documents/1`)
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "success": boolean,  // true
      "message": "string"  // e.g., "Document deleted successfully"
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "error": "string"    // e.g., "Document not found"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/documents/document_handler.go`
  - Deletes from `expert_documents` table and removes file.
- **Notes**:
  - Requires admin role.
  - Logs deletion (e.g., via `test_api.sh`).

## Engagement Management Endpoints

### GET /api/expert-engagements
- **Purpose**: Retrieves a list of expert engagements.
- **Method**: GET
- **Path**: `/api/expert-engagements`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Query Parameters**:
  - `limit`: Integer (e.g., 5)
  - `offset`: Integer (e.g., 0)
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    [
      {
        "id": int,               // Engagement ID
        "expertId": int,         // Expert ID
        "type": "string",        // Type (e.g., "evaluation")
        "description": "string", // Description
        "startDate": "string",   // ISO 8601 timestamp
        "endDate": "string",     // ISO 8601 timestamp
        "createdAt": "string"    // ISO 8601 timestamp
      }
    ]
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"        // e.g., "Invalid query parameters"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/engagements/engagement_handler.go`
  - Queries `expert_engagements` table.
- **Notes**:
  - Accessible to authenticated users.
  - Logs retrieval (e.g., via `test_api.sh`).

## Statistics Endpoints

### GET /api/statistics
- **Purpose**: Retrieves overall system statistics.
- **Method**: GET
- **Path**: `/api/statistics`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "totalExperts": int,        // Total experts (e.g., 459)
      "activeCount": int,         // Active experts (e.g., 379)
      "bahrainiPercentage": float,// Bahraini percentage (e.g., 59.25925925925925)
      "topAreas": [
        {
          "name": "string",      // Area ID (e.g., "10")
          "count": int,          // Count (e.g., 67)
          "percentage": float     // Percentage (e.g., 14.596949891067537)
        }
      ],
      "engagementsByType": [
        {
          "name": "string",      // Type (e.g., "evaluation")
          "count": int,          // Count (e.g., 100)
          "percentage": float     // Percentage (e.g., 25)
        }
      ],
      "monthlyGrowth": [
        {
          "period": "string",    // Year-month (e.g., "2025-03")
          "count": int,          // Count (e.g., 436)
          "growthRate": float     // Rate (e.g., 0)
        }
      ],
      "mostRequestedExperts": [
        {
          "expertId": "string",  // ID (e.g., "E001")
          "name": "string",      // Name (e.g., "Ammar Jreisat")
          "count": int           // Count (e.g., 160)
        }
      ],
      "lastUpdated": "string"     // ISO 8601 timestamp
    }
    ```
  - **Error (500 Internal Server Error)**:
    ```json
    {
      "error": "string"           // e.g., "Failed to retrieve statistics"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
  - Aggregates data from multiple tables.
- **Notes**:
  - Accessible to authenticated users.
  - Logs retrieval (e.g., "Successfully retrieved system statistics").

### GET /api/statistics/growth
- **Purpose**: Retrieves expert growth statistics over a period.
- **Method**: GET
- **Path**: `/api/statistics/growth`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Query Parameters**:
  - `months`: Integer (e.g., 6) – Number of months.
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    [
      {
        "period": "string",    // Year-month (e.g., "2025-03")
        "count": int,          // Count (e.g., 436)
        "growthRate": float     // Rate (e.g., 0)
      }
    ]
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "error": "string"        // e.g., "Invalid months parameter"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
  - Queries `experts` table, grouped by month.
- **Notes**:
  - Accessible to authenticated users.
  - Logs retrieval (e.g., "Retrieved growth statistics for 6 months").

### GET /api/statistics/nationality
- **Purpose**: Retrieves nationality distribution statistics.
- **Method**: GET
- **Path**: `/api/statistics/nationality`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    {
      "stats": [
        {
          "name": "string",    // Nationality (e.g., "Bahraini")
          "count": int,        // Count (e.g., 272)
          "percentage": float   // Percentage (e.g., 62.96296296296296)
        }
      ],
      "total": int             // Total experts (e.g., 432)
    }
    ```
  - **Error (500 Internal Server Error)**:
    ```json
    {
      "error": "string"        // e.g., "Failed to retrieve nationality statistics"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
  - Aggregates from `experts` table.
- **Notes**:
  - Accessible to authenticated users.
  - Logs retrieval (e.g., "Total experts: 432 (Bahraini: 272, Non-Bahraini: 160)").

### GET /api/statistics/engagements
- **Purpose**: Retrieves engagement type statistics.
- **Method**: GET
- **Path**: `/api/statistics/engagements`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:
    ```json
    [
      {
        "name": "string",    // Type (e.g., "evaluation")
        "count": int,        // Count (e.g., 100)
        "percentage": float   // Percentage (e.g., 25)
      }
    ]
    ```
  - **Error (500 Internal Server Error)**:
    ```json
    {
      "error": "string"        // e.g., "Failed to retrieve engagement statistics"
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
  - Queries `expert_engagements` table.
- **Notes**:
  - Accessible to authenticated users.
  - Logs retrieval (e.g., "Successfully retrieved engagement statistics").

## Conclusion
This artifact provides a complete and detailed reference for all ExpertDB API endpoints, covering their purpose, request/response structures, implementation details, and usage notes. It is designed to support the small team managing the tool, ensuring clarity and ease of use for development and maintenance. The endpoints are tested via `test_api.sh`, which validates functionality and edge cases, as seen in the logs (`api_test_run_20250417_101309.log`). For further improvements, refer to `ERRORS.md` for enhanced error messaging suggestions.