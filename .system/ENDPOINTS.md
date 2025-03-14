# ExpertDB API Endpoints

This document maps all Go API endpoints in the ExpertDB system, reflecting their current status and implementation. It serves as a reference for frontend-backend integration, helping to identify and fix issues with API communication.

## Endpoint Status Key
- ✅ **Working**: Endpoint functions correctly
- ⚠️ **Partial**: Endpoint works with some issues
- ❌ **Failing**: Frontend fetch fails or backend errors occur
- 🧪 **Experimental**: Not for production use

## Table of Contents
1. [Authentication](#authentication)
2. [Experts](#experts)
3. [Expert Requests](#expert-requests)
4. [Documents](#documents)
5. [Engagements](#engagements)
6. [AI Integration](#ai-integration)
7. [ISCED Reference Data](#isced-reference-data)
8. [Statistics](#statistics)
9. [User Management](#user-management)

## Authentication

### Login ⚠️
- **URL**: `/api/auth/login`
- **Method**: `POST`
- **Authentication**: None
- **Frontend Status**: Issues with JWT persistence fixed but needs validation
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response**:
  ```json
  {
    "user": {
      "id": "integer",
      "name": "string",
      "email": "string",
      "role": "string",
      "isActive": "boolean",
      "createdAt": "timestamp",
      "lastLogin": "timestamp"
    },
    "token": "string"
  }
  ```
- **Integration Issues**: 
  - Token storage and persistence needs validation
  - Login page styling incomplete

## Experts

### List Experts ⚠️
- **URL**: `/api/experts`
- **Method**: `GET`
- **Authentication**: None
- **Frontend Status**: Fetch working but display issues
- **Query Parameters**:
  - `name`: Filter by expert name
  - `area`: Filter by area of expertise
  - `is_available`: Filter by availability (true/false)
  - `role`: Filter by role
  - `isced_level_id`: Filter by ISCED level ID
  - `isced_field_id`: Filter by ISCED field ID
  - `min_rating`: Filter by minimum rating
  - `limit`: Number of results per page (default: 10)
  - `offset`: Number of results to skip (default: 0)
  - `sort_by`: Field to sort by
  - `sort_order`: Sort order (asc/desc)
- **Response**: Array of Expert objects
- **Integration Issues**:
  - Expert data fetch works but display in UI inconsistent
  - Advanced filtering not fully implemented in frontend

### Get Expert ✅
- **URL**: `/api/experts/{id}`
- **Method**: `GET`
- **Authentication**: None
- **Frontend Status**: Working

### Create Expert ✅
- **URL**: `/api/experts`
- **Method**: `POST`
- **Authentication**: Admin only
- **Frontend Status**: Not fully implemented in UI

### Update Expert ✅
- **URL**: `/api/experts/{id}`
- **Method**: `PUT`
- **Authentication**: Admin only
- **Frontend Status**: Not fully implemented in UI

### Delete Expert ✅
- **URL**: `/api/experts/{id}`
- **Method**: `DELETE`
- **Authentication**: Admin only
- **Frontend Status**: Not implemented in UI

## Expert Requests

### Create Expert Request ❌
- **URL**: `/api/expert-requests`
- **Method**: `POST`
- **Authentication**: Authenticated user
- **Frontend Status**: Form not implemented
- **Integration Issues**: Expert request page not started

### List Expert Requests ❌
- **URL**: `/api/expert-requests`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Frontend Status**: Not implemented
- **Integration Issues**: Admin request approval UI not implemented

### Get/Update/Delete Expert Request ❌
- **Status**: Backend implemented but frontend integration pending

## Documents

### Document Management ❌
- **Status**: Backend implemented but frontend integration not started
- **Integration Issues**: Document upload/view UI not implemented

## Engagements

### Engagement Management ❌
- **Status**: Backend implemented but frontend integration not started
- **Integration Issues**: Engagement tracking UI not implemented

## AI Integration

### AI Features 🧪
- **Status**: Experimental placeholders only
- **Note**: All AI integration endpoints are currently placeholders and should not be used

## ISCED Reference Data

### Get ISCED Levels ✅
- **URL**: `/api/isced/levels`
- **Method**: `GET`
- **Authentication**: None
- **Frontend Status**: Working but UI integration partial

### Get ISCED Fields ✅
- **URL**: `/api/isced/fields`
- **Method**: `GET`
- **Authentication**: None
- **Frontend Status**: Working but UI integration partial

## Statistics

### Statistics Endpoints ❌
- **Status**: Backend implemented but frontend dashboard not started
- **Integration Issues**: Statistics visualization components not implemented

## User Management

### User Management ⚠️
- **Status**: Backend implemented but admin UI partially implemented
- **Integration Issues**: User listing works but creation/editing UI incomplete

## Authentication and Authorization

### JWT Implementation ⚠️
- **Token Lifetime**: 24 hours
- **Storage**: LocalStorage with persistence
- **Integration Issues**: 
  - Recent fixes for token persistence need validation
  - Role-based UI rendering needs testing
  - Protected routes implementation complete but needs testing

### Error Handling
- **Status**: Backend error responses standardized but frontend error handling inconsistent
- **Integration Issues**: Improved error feedback needed in UI