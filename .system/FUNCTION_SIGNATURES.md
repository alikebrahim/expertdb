# Function Signatures Reference

This document provides a reference for key Go and React function signatures in the ExpertDB codebase. Use this for tracing code paths and understanding integration points between frontend and backend.

## Table of Contents
1. [Backend (Go)](#backend-go)
   - [Authentication](#authentication)
   - [Expert Management](#expert-management)
   - [Expert Requests](#expert-requests)
   - [Document Management](#document-management)
   - [User Management](#user-management)
   - [Utilities](#utilities)
2. [Frontend (React/TypeScript)](#frontend-reacttypescript)
   - [API Functions](#api-functions)
   - [Authentication](#authentication-1)
   - [Components](#components)
   - [Utilities](#utilities-1)

## Backend (Go)

### Authentication

#### `handleLogin(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles user login, validates credentials, and generates JWT token
- **Parameters**: HTTP request and response writer
- **Returns**: JSON response with user data and token on success, error on failure

#### `generateToken(user *User) (string, error)`
- **File**: `backend/auth.go`
- **Description**: Generates JWT token for authenticated user
- **Parameters**: User object
- **Returns**: JWT token string and error

#### `authMiddleware(next http.Handler) http.Handler`
- **File**: `backend/auth.go`
- **Description**: Middleware to verify JWT authentication for protected routes
- **Parameters**: HTTP handler to wrap
- **Returns**: HTTP handler

#### `roleMiddleware(role string) func(http.Handler) http.Handler`
- **File**: `backend/auth.go`
- **Description**: Middleware to verify user role for role-restricted routes
- **Parameters**: Required role string
- **Returns**: Middleware function

### Expert Management

#### `getExperts(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles GET requests to list experts with filtering and pagination
- **Parameters**: HTTP request and response writer
- **Returns**: JSON array of expert objects

#### `getExpert(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles GET requests to retrieve a specific expert
- **Parameters**: HTTP request and response writer
- **Returns**: JSON expert object

#### `createExpert(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles POST requests to create a new expert
- **Parameters**: HTTP request and response writer
- **Returns**: JSON response with success status

#### `findExperts(ctx context.Context, filters ExpertFilters) ([]Expert, int, error)`
- **File**: `backend/expert_operations.go`
- **Description**: Finds experts matching the provided filters
- **Parameters**: Context, filter criteria
- **Returns**: Array of experts, total count, and error

### Expert Requests

#### `createExpertRequest(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles POST requests to create expert request
- **Parameters**: HTTP request and response writer
- **Returns**: JSON response with created request data

#### `updateExpertRequest(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles PUT requests to update expert request status
- **Parameters**: HTTP request and response writer
- **Returns**: JSON response with success status

#### `approveExpertRequest(ctx context.Context, requestID int) error`
- **File**: `backend/expert_request_operations.go`
- **Description**: Approves expert request and creates expert record
- **Parameters**: Context, request ID
- **Returns**: Error if operation fails

### Document Management

#### `uploadDocument(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles document upload via multipart form
- **Parameters**: HTTP request and response writer
- **Returns**: JSON response with document data

#### `getExpertDocuments(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Retrieves documents for a specific expert
- **Parameters**: HTTP request and response writer
- **Returns**: JSON array of document objects

### User Management

#### `createUser(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles POST requests to create a new user
- **Parameters**: HTTP request and response writer
- **Returns**: JSON response with created user ID

#### `getUsers(w http.ResponseWriter, r *http.Request)`
- **File**: `backend/api.go`
- **Description**: Handles GET requests to list users with pagination
- **Parameters**: HTTP request and response writer
- **Returns**: JSON array of user objects

### Utilities

#### `respondWithJSON(w http.ResponseWriter, code int, payload interface{})`
- **File**: `backend/api.go`
- **Description**: Helper to send JSON responses
- **Parameters**: Response writer, HTTP status code, payload object
- **Returns**: None (writes to response)

#### `respondWithError(w http.ResponseWriter, code int, message string)`
- **File**: `backend/api.go`
- **Description**: Helper to send error responses
- **Parameters**: Response writer, HTTP status code, error message
- **Returns**: None (writes to response)

## Frontend (React/TypeScript)

### API Functions

#### `login(email: string, password: string): Promise<LoginResponse>`
- **File**: `frontend/src/api/api.ts`
- **Description**: Authenticates user and retrieves JWT token
- **Parameters**: Email and password
- **Returns**: Promise resolving to user data and token

#### `fetchExperts(filters?: ExpertFilters): Promise<Expert[]>`
- **File**: `frontend/src/api/api.ts`
- **Description**: Retrieves experts with optional filtering
- **Parameters**: Optional filter object
- **Returns**: Promise resolving to array of experts

#### `createExpertRequest(request: ExpertRequestCreate): Promise<ExpertRequest>`
- **File**: `frontend/src/api/api.ts`
- **Description**: Creates new expert request
- **Parameters**: Expert request data
- **Returns**: Promise resolving to created request

#### `fetchUsers(): Promise<User[]>`
- **File**: `frontend/src/api/api.ts`
- **Description**: Retrieves list of users (admin only)
- **Parameters**: None
- **Returns**: Promise resolving to array of users

### Authentication

#### `useAuth()`
- **File**: `frontend/src/context/AuthContext.tsx`
- **Description**: Custom hook providing authentication context
- **Returns**: {
  - `user`: Current user or null
  - `login`: Login function
  - `logout`: Logout function
  - `isAuthenticated`: Boolean indicating auth status
  - `isAdmin`: Boolean indicating admin status
  }

#### `AuthProvider({ children })`
- **File**: `frontend/src/context/AuthContext.tsx`
- **Description**: Context provider for authentication state
- **Parameters**: Child components
- **Returns**: Auth context provider component

#### `ProtectedRoute({ children, adminOnly })`
- **File**: `frontend/src/components/ProtectedRoute.tsx`
- **Description**: Route wrapper that redirects unauthenticated users
- **Parameters**: Child components, boolean indicating admin-only access
- **Returns**: Component that renders children or redirects

### Components

#### `LoginForm({ onSuccess })`
- **File**: `frontend/src/pages/Login.tsx`
- **Description**: Login form component
- **Parameters**: Success callback
- **Returns**: Form component

#### `ExpertSearchFilters({ onFilterChange })`
- **File**: `frontend/src/pages/Search.tsx`
- **Description**: Search filters component
- **Parameters**: Filter change callback
- **Returns**: Filters component

#### `ExpertTable({ experts, loading })`
- **File**: `frontend/src/pages/Search.tsx`
- **Description**: Table displaying expert data
- **Parameters**: Expert array, loading state
- **Returns**: Table component

### Utilities

#### `formatDate(date: string): string`
- **File**: `frontend/src/lib/utils.ts`
- **Description**: Formats date string for display
- **Parameters**: ISO date string
- **Returns**: Formatted date string

#### `cn(...inputs: ClassValue[]): string`
- **File**: `frontend/src/lib/utils.ts`
- **Description**: Utility for conditional class name joining
- **Parameters**: Class values
- **Returns**: Combined class string