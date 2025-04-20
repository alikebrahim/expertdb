# ExpertDB API Endpoints Documentation

**Date**: April 20, 2025\
**Version**: 1.1\
**Context**:\
ExpertDB is a lightweight internal tool for managing a database of experts, designed for a department with 10-12 users and a maximum of 1200 database entries over 5 years. The tool operates on an intranet, with security handled organizationally, prioritizing simplicity, maintainability, and clear error messaging over complex scalability or security measures. The backend is built in Go, uses SQLite as the database, and provides a RESTful API with JSON payloads, JWT authentication, and permissive CORS settings (`*`).

**Purpose**:\
This document provides an updated reference for all API endpoints, incorporating new features and changes implemented through Phase 12 of the ExpertDB Implementation Plan. It details endpoint functionality, request/response structures, and implementation notes to guide developers and users within the department.

## Table of Contents

 1. Overview
 2. General Notes
 3. Authentication Endpoints
    - POST /api/auth/login
 4. User Management Endpoints
    - POST /api/users
    - DELETE /api/users/{id}
 5. Expert Management Endpoints
    - POST /api/experts
    - GET /api/experts
    - GET /api/experts/{id}
    - PUT /api/experts/{id}
    - DELETE /api/experts/{id}
    - GET /api/expert/areas
    - POST /api/expert/areas
    - PUT /api/expert/areas/{id}
 6. Expert Request Management Endpoints
    - POST /api/expert-requests
    - GET /api/expert-requests
    - GET /api/expert-requests/{id}
    - PUT /api/expert-requests/{id}
    - PUT /api/expert-requests/{id}/edit
    - POST /api/expert-requests/batch-approve
 7. Document Management Endpoints
    - POST /api/documents
    - GET /api/experts/{id}/documents
    - GET /api/documents/{id}
    - DELETE /api/documents/{id}
 8. Engagement Management Endpoints
    - GET /api/expert-engagements
    - POST /api/engagements/import
 9. Phase Planning Endpoints
    - POST /api/phases
    - GET /api/phases
    - PUT /api/phases/{id}/applications/{app_id}
    - PUT /api/phases/{id}/applications/{app_id}/review
10. Statistics Endpoints
    - GET /api/statistics
    - GET /api/statistics/growth
    - GET /api/statistics/nationality
    - GET /api/statistics/engagements
    - GET /api/statistics/areas
11. Backup Endpoints
    - GET /api/backup

## Overview

The ExpertDB backend provides a RESTful API for managing expert profiles, requests, users, documents, engagements, phase planning, statistics, and backups. Implemented in Go, it uses SQLite (`expertdb.sqlite`), JWT authentication, and logs requests/responses to `./logs` via `internal/logger`. The API supports a small user base with a focus on simplicity and clear error messaging, as outlined in `SRS.md` and `ExpertDB Implementation Plan.markdown`.

Recent updates (Phases 2-12) include:

- Enhanced user roles (`super_user`, `scheduler`).
- Approval document support and batch approvals.
- Extended access for expert, document, and area endpoints.
- Phase planning with application and engagement management.
- Improved statistics (published experts, yearly growth, area stats).
- Specialization area creation/renaming.
- CSV backup functionality.
- Engagement filtering and import.

Endpoints are grouped by functionality, with most requiring JWT authentication and role-based access control enforced via `internal/auth/middleware.go`.

## General Notes

- **Authentication**: All endpoints except `/api/auth/login` require a JWT token in the `Authorization: Bearer <token>` header. Role-based permissions (`super_user`, `admin`, `scheduler`, `regular`) are enforced.
- **Access Levels**:
  - **Super User**: Full access, including admin creation and deletion.
  - **Admin**: Manages users, experts, requests, documents, areas, phases, and backups.
  - **Scheduler**: Submits requests, proposes experts for phase plans, views experts/documents.
  - **Regular**: Submits requests, views experts/documents.
- **CORS**: Allows all origins (`*`), suitable for intranet use.
- **Error Handling**: Returns JSON with specific `error` messages, improved per `ERRORS.md` recommendations (e.g., aggregated validation errors).
- **Database**: SQLite with schema in `db/migrations/sqlite`. Indexes added for filters (`nationality`, `general_area`, etc.).
- **Logging**: Logs to `./logs` with HTTP status, headers, and payloads.
- **Testing**: `test_api.sh` validates endpoints, covering new features and edge cases.
- **Payload Validation**: Enforces required fields, applies defaults (e.g., `pending` status).
- **HTTP Status Codes**:
  - `200 OK`: Successful GET, PUT, DELETE.
  - `201 Created`: Successful POST.
  - `400 Bad Request`: Invalid payload/parameters.
  - `401 Unauthorized`: Invalid/missing token.
  - `403 Forbidden`: Insufficient permissions.
  - `404 Not Found`: Resource not found.
  - `409 Conflict`: Duplicate resource.
  - `500 Internal Server Error`: Unexpected errors.

## Authentication Endpoints

### POST /api/auth/login

- **Purpose**: Authenticates a user and returns a JWT token.
- **Method**: POST
- **Path**: `/api/auth/login`
- **Request Payload**:

  ```json
  {
    "email": "string",    // Required: e.g., "admin@expertdb.com"
    "password": "string"  // Required: e.g., "adminpassword"
  }
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "user": {
        "id": int,           // e.g., 1
        "name": "string",    // e.g., "Admin User"
        "email": "string",   // e.g., "admin@expertdb.com"
        "role": "string",    // e.g., "super_user"
        "isActive": boolean, // e.g., true
        "createdAt": "string", // e.g., "2025-04-10T10:04:59Z"
        "lastLogin": "string"  // e.g., "2025-04-20T10:13:09Z"
      },
      "token": "string"      // JWT token
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid request payload" }
    ```
  - **Error (401 Unauthorized)**:

    ```json
    { "error": "Invalid credentials" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/auth.go:HandleLogin`
  - Validates credentials, checks `users` table, uses `golang.org/x/crypto` for password verification, generates JWT via `internal/auth/jwt.go`.
- **Notes**:
  - Supports `super_user`, `admin`, `scheduler`, `regular` roles.
  - Logs successful logins.

## User Management Endpoints

### POST /api/users

- **Purpose**: Creates a new user account.
- **Method**: POST
- **Path**: `/api/users`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token|admin_token>`
- **Request Payload**:

  ```json
  {
    "name": "string",     // Required: e.g., "Test User"
    "email": "string",    // Required: e.g., "test@example.com"
    "password": "string", // Required: e.g., "password123"
    "role": "string",     // Required: e.g., "scheduler"
    "isActive": boolean   // Required: e.g., true
  }
  ```
- **Response Payload**:
  - **Success (201 Created)**:

    ```json
    {
      "id": int,           // e.g., 18
      "success": true,
      "message": "User created successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid request payload" }
    ```
  - **Error (403 Forbidden)**:

    ```json
    { "error": "Creator cannot create user with role 'admin'" }
    ```
  - **Error (409 Conflict)**:

    ```json
    { "error": "Email already exists" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/user.go`
  - Enforces role hierarchy via `internal/storage/sqlite/user.go:CreateUserWithRoleCheck`.
  - Super users create admins; admins create schedulers/regular users.
- **Notes**:
  - Logs creation (e.g., "New user created: test@example.com").
  - Updated in Phase 2 to support `super_user` and `scheduler` roles.

### DELETE /api/users/{id}

- **Purpose**: Deletes a user by ID.
- **Method**: DELETE
- **Path**: `/api/users/{id}`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token|admin_token>`
- **Request Payload**: None
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "User deleted successfully"
    }
    ```
  - **Error (403 Forbidden)**:

    ```json
    { "error": "Only super users can delete admin accounts" }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "User not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/user.go`
  - Super users delete admins; admins delete schedulers/regular users.
  - Cascades deletion of scheduler assignments.
- **Notes**:
  - Logs deletion (e.g., "User deleted: ID 18").
  - Updated in Phase 2 for role-based deletion restrictions.

## Expert Management Endpoints

### POST /api/experts

- **Purpose**: Creates a new expert profile.
- **Method**: POST
- **Path**: `/api/experts`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:

  ```json
  {
    "name": "string",           // Required
    "institution": "string",    // Required
    "email": "string",          // Required
    "designation": "string",    // Required
    "isBahraini": boolean,      // Required
    "isAvailable": boolean,     // Required
    "rating": "string",         // Required
    "role": "string",           // Required
    "employmentType": "string", // Required
    "generalArea": int,         // Required
    "specializedArea": "string",// Required
    "isTrained": boolean,       // Required
    "cvPath": "string",         // Required
    "phone": "string",          // Required
    "isPublished": boolean,     // Required
    "biography": "string",      // Required
    "skills": ["string"],       // Required
    "approvalDocumentPath": "string" // Required
  }
  ```
- **Response Payload**:
  - **Success (201 Created)**:

    ```json
    {
      "id": int,
      "success": true,
      "message": "Expert created successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "errors": ["name is required", "invalid general_area"] }
    ```
  - **Error (409 Conflict)**:

    ```json
    { "error": "Expert ID already exists" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go:HandleCreateExpert`
  - Generates unique `expert_id` (e.g., `EXP-0001`) via `internal/storage/sqlite/expert.go`.
  - Validates fields, stores in `experts` table.
- **Notes**:
  - Updated in Phase 1 to fix `UNIQUE constraint` issue.
  - Phase 5 added `approvalDocumentPath`.
  - Logs creation (e.g., "Creating expert: Test Expert").

### GET /api/experts

- **Purpose**: Retrieves a paginated list of experts with filters and sorting.
- **Method**: GET
- **Path**: `/api/experts`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Query Parameters**:
  - `limit`, `offset`: Pagination.
  - `sort_by`: e.g., `name`, `rating`, `expert_id`, `specialized_area`.
  - `sort_order`: `asc` or `desc`.
  - `by_nationality`: `Bahraini` or `non-Bahraini`.
  - `by_general_area`: Area ID.
  - `by_specialized_area`: Text search.
  - `by_employment_type`: e.g., `academic`.
  - `by_role`: e.g., `evaluator`.
- **Response Headers**:
  - `X-Total-Count`, `X-Total-Pages`, `X-Current-Page`, `X-Page-Size`, `X-Has-Next-Page`, `X-Has-Prev-Page`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "experts": [
        {
          "id": int,
          "expertId": "string",
          "name": "string",
          "designation": "string",
          "institution": "string",
          "isBahraini": boolean,
          "isAvailable": boolean,
          "rating": "string",
          "role": "string",
          "employmentType": "string",
          "generalArea": int,
          "generalAreaName": "string",
          "specializedArea": "string",
          "isTrained": boolean,
          "cvPath": "string",
          "phone": "string",
          "email": "string",
          "isPublished": boolean,
          "biography": "string",
          "approvalDocumentPath": "string",
          "createdAt": "string",
          "updatedAt": "string"
        }
      ],
      "pagination": {
        "totalCount": int,
        "totalPages": int,
        "currentPage": int,
        "pageSize": int,
        "hasNextPage": boolean,
        "hasPrevPage": boolean,
        "hasMore": boolean
      }
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid query parameters" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Supports extended filters and sorting (Phase 3).
- **Notes**:
  - Accessible to all authenticated users (Phase 2D).
  - Logs queries (e.g., "Retrieving experts with filters").

### GET /api/experts/{id}

- **Purpose**: Retrieves a specific expert’s details.
- **Method**: GET
- **Path**: `/api/experts/{id}`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "id": int,
      "expertId": "string",
      "name": "string",
      "designation": "string",
      "institution": "string",
      "isBahraini": boolean,
      "isAvailable": boolean,
      "rating": "string",
      "role": "string",
      "employmentType": "string",
      "generalArea": int,
      "generalAreaName": "string",
      "specializedArea": "string",
      "isTrained": boolean,
      "cvPath": "string",
      "phone": "string",
      "email": "string",
      "isPublished": boolean,
      "biography": "string",
      "approvalDocumentPath": "string",
      "createdAt": "string",
      "updatedAt": "string"
    }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "Expert not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Includes `approvalDocumentPath` (Phase 3C).
- **Notes**:
  - Accessible to all authenticated users (Phase 2D).

### PUT /api/experts/{id}

- **Purpose**: Updates an expert profile.
- **Method**: PUT
- **Path**: `/api/experts/{id}`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:

  ```json
  {
    "name": "string",
    "institution": "string",
    "email": "string",
    "designation": "string",
    "isBahraini": boolean,
    "isAvailable": boolean,
    "rating": "string",
    "role": "string",
    "employmentType": "string",
    "generalArea": int,
    "specializedArea": "string",
    "isTrained": boolean,
    "cvPath": "string",
    "phone": "string",
    "isPublished": boolean,
    "biography": "string",
    "skills": ["string"],
    "approvalDocumentPath": "string"
  }
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Expert updated successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid request payload" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Updates specified fields, including `approvalDocumentPath`.
- **Notes**:
  - Logs updates (e.g., "Expert updated: ID 440").

### DELETE /api/experts/{id}

- **Purpose**: Deletes an expert.
- **Method**: DELETE
- **Path**: `/api/experts/{id}`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Expert deleted successfully"
    }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "Expert not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Cascades to documents, including approval documents (Phase 6C).
- **Notes**:
  - Logs deletion (e.g., "Expert deleted: ID 440").

### GET /api/expert/areas

- **Purpose**: Retrieves all specialization areas.
- **Method**: GET
- **Path**: `/api/expert/areas`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    [
      { "id": int, "name": "string" }
    ]
    ```
  - **Error (401 Unauthorized)**:

    ```json
    { "error": "Unauthorized" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Extended to all authenticated users (Phase 8A).
- **Notes**:
  - Logs retrieval (e.g., "Returning 34 expert areas").

### POST /api/expert/areas

- **Purpose**: Creates a new specialization area.
- **Method**: POST
- **Path**: `/api/expert/areas`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:

  ```json
  { "name": "string" } // Required
  ```
- **Response Payload**:
  - **Success (201 Created)**:

    ```json
    {
      "id": int,
      "success": true,
      "message": "Area created successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Name is required" }
    ```
  - **Error (409 Conflict)**:

    ```json
    { "error": "Area name already exists" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Inserts into `expert_areas` table (Phase 8B).
- **Notes**:
  - Logs creation (e.g., "Area created: New Area").

### PUT /api/expert/areas/{id}

- **Purpose**: Renames a specialization area.
- **Method**: PUT
- **Path**: `/api/expert/areas/{id}`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:

  ```json
  { "name": "string" } // Required
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Area updated successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Name is required" }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "Area not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert.go`
  - Updates `expert_areas`, cascades to `experts` and `expert_requests` (Phase 8C).
- **Notes**:
  - Uses transactions for integrity.
  - Logs update (e.g., "Area renamed: ID 1").

## Expert Request Management Endpoints

### POST /api/expert-requests

- **Purpose**: Submits an expert request with CV upload.
- **Method**: POST
- **Path**: `/api/expert-requests`
- **Request Headers**:
  - `Authorization: Bearer <scheduler_token|regular_token>`
- **Request Payload**: Form-data

  ```text
  name: string           // Required
  designation: string     // Required
  institution: string     // Required
  isBahraini: boolean    // Required
  isAvailable: boolean   // Required
  rating: string         // Required
  role: string           // Required
  employmentType: string // Required
  generalArea: int       // Required
  specializedArea: string// Required
  isTrained: boolean     // Required
  phone: string          // Required
  email: string          // Required
  biography: string      // Required
  skills: string         // Required: JSON array
  isPublished: boolean   // Optional, defaults to false
  cv: file               // Required: PDF
  ```
- **Response Payload**:
  - **Success (201 Created)**:

    ```json
    {
      "id": int,
      "name": "string",
      "status": "pending",
      "cvPath": "string",
      ...
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "errors": ["name is required", "cv file missing"] }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Stores CV via `internal/documents/service.go` (Phase 4A).
- **Notes**:
  - Logs creation (e.g., "Expert request created: ID 26").
  - Improved validation (Phase 4A).

### GET /api/expert-requests

- **Purpose**: Retrieves paginated expert requests with status filtering.
- **Method**: GET
- **Path**: `/api/expert-requests`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Query Parameters**:
  - `limit`, `offset`: Pagination.
  - `status`: `pending`, `approved`, `rejected`.
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    [
      {
        "id": int,
        "name": "string",
        "status": "string",
        "cvPath": "string",
        "approvalDocumentPath": "string",
        ...
      }
    ]
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid status parameter" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Filters by status (Phase 4B).
- **Notes**:
  - Logs retrieval (e.g., "Returning 5 requests").

### GET /api/expert-requests/{id}

- **Purpose**: Retrieves a specific expert request.
- **Method**: GET
- **Path**: `/api/expert-requests/{id}`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "id": int,
      "name": "string",
      "status": "string",
      "cvPath": "string",
      "approvalDocumentPath": "string",
      ...
    }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "Expert request not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Includes `cvPath` (Phase 4B).
- **Notes**:
  - Logs retrieval (e.g., "Retrieved request: ID 26").

### PUT /api/expert-requests/{id}

- **Purpose**: Approves or rejects an expert request with approval document.
- **Method**: PUT
- **Path**: `/api/expert-requests/{id}`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: Form-data

  ```text
  status: string          // Required: "approved" or "rejected"
  rejectionReason: string // Optional
  approvalDocument: file  // Required for approval
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Expert request updated successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Approval document required" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Stores approval document, creates expert on approval (Phase 5B).
- **Notes**:
  - Logs update (e.g., "Request approved: ID 26").

### PUT /api/expert-requests/{id}/edit

- **Purpose**: Edits an expert request before approval.
- **Method**: PUT
- **Path**: `/api/expert-requests/{id}/edit`
- **Request Headers**:
  - `Authorization: Bearer <admin_token|owner_token>`
- **Request Payload**: Form-data

  ```text
  name: string
  designation: string
  institution: string
  ...
  cv: file
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Expert request updated successfully"
    }
    ```
  - **Error (403 Forbidden)**:

    ```json
    { "error": "Only admins or request owner can edit" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Admins edit any request; users edit their rejected requests (Phase 4C).
- **Notes**:
  - Logs update (e.g., "Request edited: ID 26").

### POST /api/expert-requests/batch-approve

- **Purpose**: Approves multiple expert requests with one approval document.
- **Method**: POST
- **Path**: `/api/expert-requests/batch-approve`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: Form-data

  ```text
  requestIds: string     // JSON array of IDs
  approvalDocument: file // Required
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Batch approval completed",
      "results": [
        { "id": int, "status": "success" | "failed", "error": "string" }
      ]
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Missing approval document" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/expert_request.go`
  - Uses transactions for consistency (Phase 5C).
- **Notes**:
  - Logs results (e.g., "Batch approved 3 requests").

## Document Management Endpoints

### POST /api/documents

- **Purpose**: Uploads a document for an expert.
- **Method**: POST
- **Path**: `/api/documents`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: Form-data

  ```text
  file: file           // Required
  documentType: string // Required: e.g., "cv", "approval"
  expertId: int        // Required
  ```
- **Response Payload**:
  - **Success (201 Created)**:

    ```json
    {
      "id": int,
      "success": true,
      "message": "Document uploaded successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid document type" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/documents/document_handler.go`
  - Supports `cv`, `approval`, `certificate`, `publication`, `other` (Phase 6B).
- **Notes**:
  - Logs upload (e.g., "Document uploaded for expert: ID 440").

### GET /api/experts/{id}/documents

- **Purpose**: Lists documents for an expert.
- **Method**: GET
- **Path**: `/api/experts/{id}/documents`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    [
      {
        "id": int,
        "expertId": int,
        "documentType": "string",
        "filePath": "string",
        "createdAt": "string"
      }
    ]
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "Expert not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/documents/document_handler.go`
  - Accessible to all users (Phase 6A).
- **Notes**:
  - Includes CVs and approval documents.

### GET /api/documents/{id}

- **Purpose**: Retrieves a specific document.
- **Method**: GET
- **Path**: `/api/documents/{id}`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "id": int,
      "expertId": int,
      "documentType": "string",
      "filePath": "string",
      "createdAt": "string"
    }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "Document not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/documents/document_handler.go`
  - Accessible to all users (Phase 6A).
- **Notes**:
  - Logs retrieval (e.g., "Retrieved document: ID 1").

### DELETE /api/documents/{id}

- **Purpose**: Deletes a document.
- **Method**: DELETE
- **Path**: `/api/documents/{id}`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Document deleted successfully"
    }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "Document not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/documents/document_handler.go`
  - Cascades with expert deletion (Phase 6C).
- **Notes**:
  - Logs deletion (e.g., "Document deleted: ID 1").

## Engagement Management Endpoints

### GET /api/expert-engagements

- **Purpose**: Lists engagements with filters.
- **Method**: GET
- **Path**: `/api/expert-engagements`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Query Parameters**:
  - `limit`, `offset`: Pagination.
  - `expert_id`: Filter by expert ID.
  - `type`: `validator` or `evaluator`.
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    [
      {
        "id": int,
        "expertId": int,
        "engagementType": "string",
        "startDate": "string",
        "projectName": "string",
        "status": "string",
        "notes": "string",
        "createdAt": "string"
      }
    ]
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid query parameters" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/engagements/engagement_handler.go`
  - Filters added in Phase 11A; types restricted to `validator`, `evaluator` (Phase 11B).
- **Notes**:
  - Logs retrieval (e.g., "Retrieved engagements").

### POST /api/engagements/import

- **Purpose**: Imports past engagements from CSV/JSON.
- **Method**: POST
- **Path**: `/api/engagements/import`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**: Form-data

  ```text
  format: string      // Required: "csv" or "json"
  file: file          // Required
  ```
  - CSV Example:

    ```csv
    expert_id,type,date,details
    1,validator,2025-01-01,Project X
    ```
  - JSON Example:

    ```json
    [
      {
        "expert_id": 1,
        "type": "validator",
        "date": "2025-01-01",
        "details": "Project X"
      }
    ]
    ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Engagements imported successfully",
      "imported": int,
      "failed": int
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid format" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/engagements/engagement_handler.go`
  - Validates and deduplicates via `internal/storage/sqlite/engagement.go` (Phase 11C).
- **Notes**:
  - Logs import results (e.g., "Imported 10 engagements").

## Phase Planning Endpoints

### POST /api/phases

- **Purpose**: Creates a phase plan with applications.
- **Method**: POST
- **Path**: `/api/phases`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:

  ```json
  {
    "title": "string",
    "assignedSchedulerId": int,
    "status": "string",
    "applications": [
      {
        "type": "string",       // "QP" or "IL"
        "institutionName": "string",
        "qualificationName": "string",
        "expert1": int,
        "expert2": int,
        "status": "string"
      }
    ]
  }
  ```
- **Response Payload**:
  - **Success (201 Created)**:

    ```json
    {
      "id": int,
      "success": true,
      "message": "Phase created successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Title is required" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/phase/phase_handler.go`
  - Stores in `phases` and `phase_applications` tables (Phase 10B).
- **Notes**:
  - Logs creation (e.g., "Phase created: ID 1").

### GET /api/phases

- **Purpose**: Lists phase plans with filters.
- **Method**: GET
- **Path**: `/api/phases`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Query Parameters**:
  - `limit`, `offset`: Pagination.
  - `status`: Phase status.
  - `scheduler_id`: Assigned scheduler ID.
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    [
      {
        "id": int,
        "phaseId": "string",
        "title": "string",
        "assignedSchedulerId": int,
        "status": "string",
        "createdAt": "string",
        "updatedAt": "string",
        "applications": [
          {
            "id": int,
            "phaseId": int,
            "type": "string",
            "institutionName": "string",
            "qualificationName": "string",
            "expert1": int,
            "expert1Name": "string",
            "expert2": int,
            "expert2Name": "string",
            "status": "string",
            "rejectionNotes": "string",
            "createdAt": "string",
            "updatedAt": "string"
          }
        ]
      }
    ]
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid query parameters" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/phase/phase_handler.go`
  - Supports filtering (Phase 10E).
- **Notes**:
  - Logs retrieval (e.g., "Retrieved phases").

### PUT /api/phases/{id}/applications/{app_id}

- **Purpose**: Proposes experts for a phase application.
- **Method**: PUT
- **Path**: `/api/phases/{id}/applications/{app_id}`
- **Request Headers**:
  - `Authorization: Bearer <scheduler_token>`
- **Request Payload**:

  ```json
  {
    "expert1": int,
    "expert2": int
  }
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Experts proposed successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid expert ID" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/phase/phase_handler.go`
  - Validates expert IDs (Phase 10C).
- **Notes**:
  - Logs update (e.g., "Experts proposed for application: ID 1").

### PUT /api/phases/{id}/applications/{app_id}/review

- **Purpose**: Approves or rejects a phase application, creating engagements.
- **Method**: PUT
- **Path**: `/api/phases/{id}/applications/{app_id}/review`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Request Payload**:

  ```json
  {
    "status": "string",       // "approved", "rejected", "pending"
    "rejectionNotes": "string"
  }
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "Application reviewed successfully"
    }
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Rejection notes required for rejection" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/phase/phase_handler.go`
  - Creates `validator`/`evaluator` engagements on approval (Phase 10D).
- **Notes**:
  - Uses transactions; logs review (e.g., "Application approved: ID 1").

## Statistics Endpoints

### GET /api/statistics

- **Purpose**: Retrieves overall system statistics.
- **Method**: GET
- **Path**: `/api/statistics`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "totalExperts": int,
      "activeCount": int,
      "bahrainiPercentage": float,
      "publishedCount": int,
      "publishedRatio": float,
      "topAreas": [
        { "name": "string", "count": int, "percentage": float }
      ],
      "engagementsByType": [
        { "name": "string", "count": int, "percentage": float }
      ],
      "yearlyGrowth": [
        { "period": "string", "count": int, "growthRate": float }
      ],
      "mostRequestedExperts": [
        { "expertId": "string", "name": "string", "count": int }
      ],
      "lastUpdated": "string"
    }
    ```
  - **Error (500 Internal Server Error)**:

    ```json
    { "error": "Failed to retrieve statistics" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
  - Added `publishedCount`, `publishedRatio` (Phase 7A).
- **Notes**:
  - Logs retrieval (e.g., "Retrieved statistics").

### GET /api/statistics/growth

- **Purpose**: Retrieves yearly expert growth.
- **Method**: GET
- **Path**: `/api/statistics/growth`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token>`
- **Query Parameters**:
  - `years`: e.g., 5
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    [
      { "period": "string", "count": int, "growthRate": float }
    ]
    ```
  - **Error (400 Bad Request)**:

    ```json
    { "error": "Invalid years parameter" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
  - Switched to yearly from monthly (Phase 7B).
- **Notes**:
  - Logs retrieval (e.g., "Retrieved growth for 5 years").

### GET /api/statistics/nationality

- **Purpose**: Retrieves nationality distribution.
- **Method**: GET
- **Path**: `/api/statistics/nationality`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "stats": [
        { "name": "string", "count": int, "percentage": float }
      ],
      "total": int
    }
    ```
  - **Error (500 Internal Server Error)**:

    ```json
    { "error": "Failed to retrieve nationality statistics" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
- **Notes**:
  - Logs retrieval (e.g., "Retrieved nationality stats").

### GET /api/statistics/engagements

- **Purpose**: Retrieves engagement type statistics.
- **Method**: GET
- **Path**: `/api/statistics/engagements`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    [
      { "name": "string", "count": int, "percentage": float }
    ]
    ```
  - **Error (500 Internal Server Error)**:

    ```json
    { "error": "Failed to retrieve engagement statistics" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
  - Restricted to `validator`, `evaluator` (Phase 7C).
- **Notes**:
  - Logs retrieval (e.g., "Retrieved engagement stats").

### GET /api/statistics/areas

- **Purpose**: Retrieves general and specialized area statistics.
- **Method**: GET
- **Path**: `/api/statistics/areas`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "generalAreas": [
        { "name": "string", "count": int, "percentage": float }
      ],
      "topSpecializedAreas": [
        { "name": "string", "count": int, "percentage": float }
      ],
      "bottomSpecializedAreas": [
        { "name": "string", "count": int, "percentage": float }
      ]
    }
    ```
  - **Error (500 Internal Server Error)**:

    ```json
    { "error": "Failed to retrieve area statistics" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/statistics/statistics_handler.go`
  - Added in Phase 7D for top/bottom 5 specialized areas.
- **Notes**:
  - Logs retrieval (e.g., "Retrieved area stats").

## Backup Endpoints

### GET /api/backup

- **Purpose**: Generates a ZIP file with CSV exports of database tables.
- **Method**: GET
- **Path**: `/api/backup`
- **Request Headers**:
  - `Authorization: Bearer <admin_token>`
- **Response Payload**:
  - **Success (200 OK)**:
    - Content-Type: `application/zip`
    - Content-Disposition: `attachment; filename=expertdb_backup.zip`
    - ZIP file containing CSVs for `experts`, `expert_requests`, `expert_engagements`, `expert_documents`, `expert_areas`.
  - **Error (500 Internal Server Error)**:

    ```json
    { "error": "Failed to generate backup" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/backup/backup_handler.go`
  - Uses `archive/zip` (Phase 9A).
- **Notes**:
  - Logs generation (e.g., "Backup generated successfully").

## Conclusion

This updated API reference reflects enhancements through Phase 12, including new endpoints for phase planning, engagement imports, area management, and backups. It supports the department’s needs with clear documentation, tested via `test_api.sh`, and aligns with `SRS.md` requirements. For further error handling improvements, refer to `ERRORS.md`.