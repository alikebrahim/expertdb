# ExpertDB User Management API Reference

**Date**: July 21, 2025  
**Version**: 1.0  
**Purpose**: This document provides a focused reference for user management and role assignment endpoints in the ExpertDB API.

## Table of Contents

1. [Overview](#overview)
2. [User Management Endpoints](#user-management-endpoints)
   - [GET /api/users](#get-apiusers)
   - [GET /api/users/{id}](#get-apiusersid)
   - [GET /api/users/me](#get-apiusersme)
   - [POST /api/users](#post-apiusers)
   - [PUT /api/users/{id}](#put-apiusersid)
   - [DELETE /api/users/{id}](#delete-apiusersid)
3. [Role Assignment Endpoints](#role-assignment-endpoints)
   - [POST /api/users/{id}/planner-assignments](#post-apiusersidplanner-assignments)
   - [POST /api/users/{id}/manager-assignments](#post-apiusersidmanager-assignments)
   - [DELETE /api/users/{id}/planner-assignments](#delete-apiusersidplanner-assignments)
   - [DELETE /api/users/{id}/manager-assignments](#delete-apiusersidmanager-assignments)
   - [GET /api/users/{id}/assignments](#get-apiusersidassignments)

## Overview

The ExpertDB user management system implements a three-tier role system with contextual elevations:

**Base Roles:**
- `super_user`: Complete system access, can create admin users
- `admin`: Full system access, can create regular users and manage all phases/applications
- `user`: Can submit expert requests, view expert data/documents, and view all phases. Can be elevated for specific applications to propose experts (planner) or provide ratings upon admin request (manager)

**Contextual Elevations:**
Regular users can be elevated to have special privileges for specific applications within phases:

- **Planner Elevation**: Allows user to propose experts for assigned applications
  - Scoped to specific applications within a phase
  - Managed via `/api/users/{id}/planner-assignments` endpoints
  - Uses `RequirePlannerForApplication` middleware for access control

- **Manager Elevation**: Allows user to provide expert ratings for assigned applications when requested by admin
  - Scoped to specific applications within a phase
  - Managed via `/api/users/{id}/manager-assignments` endpoints
  - Uses `RequireManagerForApplication` middleware for access control

**Implementation Details:**
- Database tables: `application_planners` and `application_managers`
- Storage methods: `IsUserPlannerForApplication`, `IsUserManagerForApplication`
- Middleware: Admin/super_user bypass elevation checks and have inherent access
- API endpoints use application-specific access control rather than global role checks

## User Management Endpoints

### GET /api/users

- **Purpose**: Retrieves a paginated list of users.
- **Method**: GET
- **Path**: `/api/users`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token|admin_token>`
- **Query Parameters**:
  - `limit`, `offset`: Pagination.
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "data": {
        "users": [
          {
            "id": int,
            "name": "string",
            "email": "string",
            "role": "string",
            "isActive": boolean,
            "createdAt": "string",
            "lastLogin": "string"
          }
        ],
        "pagination": {
          "limit": int,
          "offset": int,
          "count": int
        }
      }
    }
    ```
- **Implementation**:
  - File: `internal/api/handlers/user.go`
  - Lists users with pagination, excluding password data.
- **Notes**:
  - Super users see all users; admins only see non-admin users.

### GET /api/users/{id}

- **Purpose**: Retrieves a specific user's details.
- **Method**: GET
- **Path**: `/api/users/{id}`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token|admin_token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "data": {
        "id": int,
        "name": "string",
        "email": "string",
        "role": "string",
        "isActive": boolean,
        "createdAt": "string",
        "lastLogin": "string"
      }
    }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "User not found" }
    ```
- **Implementation**:
  - File: `internal/api/handlers/user.go`
  - Super users can view any user; admins can only view non-admin users.
- **Notes**:
  - Password data is never returned.

### GET /api/users/me

- **Purpose**: Retrieves the current authenticated user's profile.
- **Method**: GET
- **Path**: `/api/users/me`
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "data": {
        "id": int,
        "name": "string",
        "email": "string",
        "role": "string",
        "isActive": boolean,
        "createdAt": "string",
        "lastLogin": "string"
      }
    }
    ```
- **Implementation**:
  - User ID extracted from JWT token claims
  - Available to all authenticated users
- **Notes**:
  - Users can always access their own profile regardless of role

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
    "role": "string",     // Required: e.g., "user"
    "isActive": boolean   // Required: e.g., true
  }
  ```
- **Response Payload**:
  - **Success (201 Created)**:

    ```json
    {
      "success": true,
      "message": "User created successfully",
      "data": {
        "id": int           // e.g., 18
      }
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
  - Super users create admins; admins create users.
- **Notes**:
  - Logs creation (e.g., "New user created: test@example.com").
  - Updated for new role system with contextual elevations.

### PUT /api/users/{id}

- **Purpose**: Updates an existing user's information.
- **Method**: PUT
- **Path**: `/api/users/{id}`
- **Request Headers**:
  - `Authorization: Bearer <super_user_token|admin_token>`
- **Request Payload**:

  ```json
  {
    "name": "string",     // Optional
    "email": "string",    // Optional
    "password": "string", // Optional
    "role": "string",     // Optional
    "isActive": boolean   // Optional
  }
  ```
- **Response Payload**:
  - **Success (200 OK)**:

    ```json
    {
      "success": true,
      "message": "User updated successfully"
    }
    ```
  - **Error (403 Forbidden)**:

    ```json
    { "error": "Only super users can modify admin accounts" }
    ```
  - **Error (404 Not Found)**:

    ```json
    { "error": "User not found" }
    ```
- **Implementation**:
  - Updates only provided fields
  - Role hierarchy enforced (admins cannot modify other admins)
  - Password is hashed if provided
- **Notes**:
  - Contextual role assignments are managed separately through role assignment endpoints

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
  - Super users delete admins; admins delete users.
  - Cascades deletion of contextual role assignments.
- **Notes**:
  - Logs deletion (e.g., "User deleted: ID 18").
  - Updated in Phase 2 for role-based deletion restrictions.

## Role Assignment Endpoints

### POST /api/users/{id}/planner-assignments

- **Purpose**: Assigns a user as planner to multiple applications within phases.
- **Method**: POST
- **Path**: `/api/users/{id}/planner-assignments`
- **Request Headers**:
  - `Authorization: Bearer <JWT_TOKEN>`
  - `Content-Type: application/json`
- **Access Control**: Admin only
- **Request Payload**:
  ```json
  {
    "application_ids": [1, 2, 3]
  }
  ```
- **Response Payload (Success)**:
  ```json
  {
    "success": true,
    "data": {
      "message": "Planner assignments updated successfully",
      "user_id": 123,
      "assigned_applications": 3
    }
  }
  ```
- **Implementation**:
  - File: `internal/api/handlers/role_assignments.go`
  - Replaces existing planner assignments for the user with the new list
  - Uses batch operations within a database transaction
- **Notes**:
  - User must exist in the system
  - Application IDs must be valid existing applications
  - Assignment is contextual - limited to specific applications within phases

### POST /api/users/{id}/manager-assignments

- **Purpose**: Assigns a user as manager to multiple applications within phases.
- **Method**: POST
- **Path**: `/api/users/{id}/manager-assignments`
- **Request Headers**:
  - `Authorization: Bearer <JWT_TOKEN>`
  - `Content-Type: application/json`
- **Access Control**: Admin only
- **Request Payload**:
  ```json
  {
    "application_ids": [1, 2, 3]
  }
  ```
- **Response Payload (Success)**:
  ```json
  {
    "success": true,
    "data": {
      "message": "Manager assignments updated successfully",
      "user_id": 123,
      "assigned_applications": 3
    }
  }
  ```
- **Implementation**:
  - File: `internal/api/handlers/role_assignments.go`
  - Replaces existing manager assignments for the user with the new list
  - Uses batch operations within a database transaction
- **Notes**:
  - Manager role provides rating privileges for assigned applications
  - User can receive requests and provide expert ratings only for assigned applications

### DELETE /api/users/{id}/planner-assignments

- **Purpose**: Removes planner assignments for a user from specific applications.
- **Method**: DELETE
- **Path**: `/api/users/{id}/planner-assignments`
- **Request Headers**:
  - `Authorization: Bearer <JWT_TOKEN>`
  - `Content-Type: application/json`
- **Access Control**: Admin only
- **Request Payload**:
  ```json
  {
    "application_ids": [1, 2]
  }
  ```
- **Response Payload (Success)**:
  ```json
  {
    "success": true,
    "data": {
      "message": "Planner assignments removed successfully",
      "user_id": 123,
      "removed_applications": 2
    }
  }
  ```

### DELETE /api/users/{id}/manager-assignments

- **Purpose**: Removes manager assignments for a user from specific applications.
- **Method**: DELETE
- **Path**: `/api/users/{id}/manager-assignments`
- **Request Headers**:
  - `Authorization: Bearer <JWT_TOKEN>`
  - `Content-Type: application/json`
- **Access Control**: Admin only
- **Request Payload**:
  ```json
  {
    "application_ids": [1, 2]
  }
  ```
- **Response Payload (Success)**:
  ```json
  {
    "success": true,
    "data": {
      "message": "Manager assignments removed successfully",
      "user_id": 123,
      "removed_applications": 2
    }
  }
  ```

### GET /api/users/{id}/assignments

- **Purpose**: Retrieves all planner and manager assignments for a specific user.
- **Method**: GET
- **Path**: `/api/users/{id}/assignments`
- **Request Headers**:
  - `Authorization: Bearer <JWT_TOKEN>`
- **Access Control**: Admin only
- **Response Payload (Success)**:
  ```json
  {
    "success": true,
    "data": {
      "user_id": 123,
      "planner_applications": [1, 3, 5],
      "manager_applications": [2, 4, 6]
    }
  }
  ```
- **Implementation**:
  - File: `internal/api/handlers/role_assignments.go`
  - Returns array of application IDs for each role type
- **Notes**:
  - Used by admin UI to display current assignments when managing user roles
  - Applications can be cross-referenced with phase data to show context

## Standard Response Structure

All endpoints use a standard response structure:
```json
{
  "success": boolean,    // true for successful requests
  "message": "string",   // optional success message
  "data": <object>       // optional response data
}
```

## HTTP Status Codes

- `200 OK`: Successful GET, PUT, DELETE.
- `201 Created`: Successful POST.
- `400 Bad Request`: Invalid payload/parameters.
- `401 Unauthorized`: Invalid/missing token.
- `403 Forbidden`: Insufficient permissions.
- `404 Not Found`: Resource not found.
- `409 Conflict`: Duplicate resource.
- `500 Internal Server Error`: Unexpected errors.

## Authentication

All endpoints except `/api/auth/login` require a JWT token in the `Authorization: Bearer <token>` header. Role-based permissions (`super_user`, `admin`, `user`) are enforced.

## Notes

- Password data is never returned in responses
- Role hierarchy is strictly enforced (admins cannot create/modify super_users)
- Contextual elevations (planner/manager) are managed separately from base roles
- All user operations are logged for audit purposes