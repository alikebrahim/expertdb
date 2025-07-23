# ExpertDB API Reference - Authentication & Authorization

This document provides comprehensive API documentation for authentication and authorization endpoints in the ExpertDB system.

## Table of Contents

- [Overview](#overview)
- [Authentication Flow](#authentication-flow)
- [Endpoints](#endpoints)
  - [POST /api/auth/login](#post-apiauthlogin)
- [Security Considerations](#security-considerations)
- [Error Handling](#error-handling)

## Overview

The ExpertDB authentication system uses JWT (JSON Web Tokens) for secure API access. All API endpoints except `/api/auth/login` and `/api/health` require authentication via Bearer tokens.

### Key Features:
- JWT-based authentication with 24-hour token expiration
- Role-based access control (super_user, admin, user)
- Contextual elevations for planner/manager privileges
- Secure password hashing using golang.org/x/crypto
- Automatic last_login tracking

## Authentication Flow

1. **Login Request**: Client sends credentials to `/api/auth/login`
2. **Credential Validation**: Server validates email/password against database
3. **Token Generation**: On success, server generates JWT token with user claims
4. **Token Usage**: Client includes token in `Authorization: Bearer <token>` header for subsequent requests
5. **Token Validation**: Server validates token on each protected endpoint request

## Endpoints

### POST /api/auth/login

Authenticates a user and returns a JWT token for API access.

- **Method**: POST
- **Path**: `/api/auth/login`
- **Authentication**: Not required (public endpoint)

#### Request Payload

```json
{
  "email": "string",    // Required: User's email address
  "password": "string"  // Required: User's password
}
```

#### Response Payload

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": int,               // User's unique identifier
      "name": "string",        // User's full name
      "email": "string",       // User's email address
      "role": "string",        // User's role: "super_user", "admin", or "user"
      "isActive": boolean,     // Account active status
      "createdAt": "string",   // ISO 8601 timestamp of account creation
      "lastLogin": "string"    // ISO 8601 timestamp of last successful login
    },
    "token": "string"          // JWT token for API authentication
  }
}
```

**Error Responses**:

- **400 Bad Request**: Invalid request payload
  ```json
  { "error": "Invalid request payload" }
  ```

- **401 Unauthorized**: Invalid credentials
  ```json
  { "error": "Invalid credentials" }
  ```

#### Example Request

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@expertdb.com",
    "password": "adminpassword"
  }'
```

#### Implementation Details

- **Handler**: `internal/api/handlers/auth.go:HandleLogin`
- **Password Verification**: Uses bcrypt via `golang.org/x/crypto`
- **Token Generation**: JWT created via `internal/auth/jwt.go`
- **Database**: Validates against `users` table
- **Logging**: Successful logins are logged with timestamp

## Security Considerations

### Token Management
- Tokens expire after 24 hours
- No token refresh mechanism (users must re-authenticate)
- Tokens contain user ID, email, and role claims
- Tokens are signed using HS256 algorithm

### Password Security
- Passwords are hashed using bcrypt with cost factor 10
- Plain text passwords are never stored or logged
- Password validation happens server-side only

### Request Security
- All API requests (except login and health) require valid JWT
- Token validation checks signature and expiration
- Invalid tokens result in 401 Unauthorized response

## Error Handling

The authentication system uses consistent error responses:

| Status Code | Error Type | Description |
|------------|------------|-------------|
| 400 | Bad Request | Missing or malformed request data |
| 401 | Unauthorized | Invalid credentials or expired token |
| 403 | Forbidden | Valid token but insufficient permissions |
| 500 | Internal Server Error | Server-side error during authentication |

### Common Authentication Errors

1. **Missing Authorization Header**
   ```json
   { "error": "Authorization header required" }
   ```

2. **Invalid Token Format**
   ```json
   { "error": "Invalid token format" }
   ```

3. **Expired Token**
   ```json
   { "error": "Token has expired" }
   ```

4. **Invalid Signature**
   ```json
   { "error": "Invalid token signature" }
   ```

---

For more information about role-based access control and contextual elevations, see [API_REFERENCE_USERS.md](./API_REFERENCE_USERS.md).