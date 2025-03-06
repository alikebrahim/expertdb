# ExpertDB API Documentation

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

### Login
- **URL**: `/api/auth/login`
- **Method**: `POST`
- **Authentication**: None
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
- **Notes**: Returns a JWT token that should be included in subsequent requests in the Authorization header as `Bearer {token}`.

## Experts

### List Experts
- **URL**: `/api/experts`
- **Method**: `GET`
- **Authentication**: None
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
  - `sort_by`: Field to sort by (name, institution, role, created_at, rating, general_area)
  - `sort_order`: Sort order (asc/desc)
- **Response**: Array of Expert objects
  ```json
  [
    {
      "id": "integer",
      "expertId": "string",
      "name": "string",
      "designation": "string",
      "institution": "string",
      "isBahraini": "boolean",
      "nationality": "string",
      "isAvailable": "boolean",
      "rating": "string",
      "role": "string",
      "employmentType": "string",
      "generalArea": "string",
      "specializedArea": "string",
      "isTrained": "boolean",
      "cvPath": "string",
      "phone": "string",
      "email": "string",
      "isPublished": "boolean",
      "iscedLevel": { "id": "integer", "code": "string", "name": "string", "description": "string" },
      "iscedField": { "id": "integer", "broadCode": "string", "broadName": "string", ... },
      "areas": [{ "id": "integer", "name": "string" }, ...],
      "biography": "string",
      "createdAt": "timestamp",
      "updatedAt": "timestamp"
    },
    ...
  ]
  ```
- **Notes**: Response includes `X-Total-Count` header with total number of experts matching the filters.

### Get Expert
- **URL**: `/api/experts/{id}`
- **Method**: `GET`
- **Authentication**: None
- **Response**: Expert object

### Create Expert
- **URL**: `/api/experts`
- **Method**: `POST`
- **Authentication**: Admin only
- **Request Body**:
  ```json
  {
    "name": "string",
    "affiliation": "string",
    "primaryContact": "string",
    "contactType": "string", 
    "skills": ["string", ...],
    "role": "string",
    "employmentType": "string",
    "generalArea": "string",
    "cvPath": "string",
    "biography": "string",
    "isBahraini": "boolean",
    "availability": "string"
  }
  ```
- **Response**:
  ```json
  {
    "id": "integer",
    "success": "boolean",
    "message": "string"
  }
  ```
- **Notes**: `contactType` must be "email" or "phone". `role` must be one of: evaluator, validator, consultant, trainer, expert. `employmentType` must be one of: academic, employer, freelance, government, other.

### Update Expert
- **URL**: `/api/experts/{id}`
- **Method**: `PUT`
- **Authentication**: Admin only
- **Request Body**: Expert object
- **Response**:
  ```json
  {
    "success": "boolean",
    "message": "string"
  }
  ```

### Delete Expert
- **URL**: `/api/experts/{id}`
- **Method**: `DELETE`
- **Authentication**: Admin only
- **Response**:
  ```json
  {
    "success": "boolean",
    "message": "string"
  }
  ```

## Expert Requests

### Create Expert Request
- **URL**: `/api/expert-requests`
- **Method**: `POST`
- **Authentication**: Authenticated user
- **Request Body**:
  ```json
  {
    "name": "string",
    "designation": "string",
    "institution": "string",
    "isBahraini": "boolean",
    "isAvailable": "boolean",
    "rating": "string",
    "role": "string",
    "employmentType": "string",
    "generalArea": "string",
    "specializedArea": "string",
    "isTrained": "boolean",
    "cvPath": "string",
    "phone": "string",
    "email": "string",
    "isPublished": "boolean"
  }
  ```
- **Response**: ExpertRequest object including the created ID

### List Expert Requests
- **URL**: `/api/expert-requests`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Query Parameters**:
  - `status`: Filter by status (pending, approved, rejected)
  - `limit`: Number of results per page (default: 100)
  - `offset`: Number of results to skip (default: 0)
- **Response**: Array of ExpertRequest objects

### Get Expert Request
- **URL**: `/api/expert-requests/{id}`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: ExpertRequest object

### Update Expert Request
- **URL**: `/api/expert-requests/{id}`
- **Method**: `PUT`
- **Authentication**: Admin only
- **Request Body**: ExpertRequest object
- **Response**:
  ```json
  {
    "success": "boolean",
    "message": "string"
  }
  ```
- **Notes**: If the status is changed to "approved", a new expert record will be automatically created.

## Documents

### Upload Document
- **URL**: `/api/documents`
- **Method**: `POST`
- **Authentication**: Authenticated user
- **Content-Type**: `multipart/form-data`
- **Form Fields**:
  - `expertId`: ID of the expert (required)
  - `documentType`: Type of document (default: "cv")
  - `file`: File to upload (required)
- **Response**: Document object
- **Notes**: Maximum file size is 10MB.

### Get Document
- **URL**: `/api/documents/{id}`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: Document object

### Delete Document
- **URL**: `/api/documents/{id}`
- **Method**: `DELETE`
- **Authentication**: Admin only
- **Response**:
  ```json
  {
    "success": "boolean",
    "message": "string"
  }
  ```

### Get Expert Documents
- **URL**: `/api/experts/{id}/documents`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: Array of Document objects

## Engagements

### Create Engagement
- **URL**: `/api/engagements`
- **Method**: `POST`
- **Authentication**: Admin only
- **Request Body**:
  ```json
  {
    "expertId": "integer",
    "engagementType": "string",
    "startDate": "timestamp",
    "endDate": "timestamp",
    "projectName": "string",
    "status": "string",
    "feedbackScore": "integer",
    "notes": "string"
  }
  ```
- **Response**: Engagement object
- **Notes**: Required fields are `expertId`, `engagementType`, and `startDate`.

### Get Engagement
- **URL**: `/api/engagements/{id}`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: Engagement object

### Update Engagement
- **URL**: `/api/engagements/{id}`
- **Method**: `PUT`
- **Authentication**: Admin only
- **Request Body**: Engagement object
- **Response**:
  ```json
  {
    "success": "boolean",
    "message": "string"
  }
  ```

### Delete Engagement
- **URL**: `/api/engagements/{id}`
- **Method**: `DELETE`
- **Authentication**: Admin only
- **Response**:
  ```json
  {
    "success": "boolean",
    "message": "string"
  }
  ```

### Get Expert Engagements
- **URL**: `/api/experts/{id}/engagements`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: Array of Engagement objects

## AI Integration

### Generate Profile
- **URL**: `/api/ai/generate-profile`
- **Method**: `POST`
- **Authentication**: Admin only
- **Request Body**:
  ```json
  {
    "expertId": "integer"
  }
  ```
- **Response**: AIAnalysisResult object

### Suggest ISCED Classification
- **URL**: `/api/ai/suggest-isced`
- **Method**: `POST`
- **Authentication**: Admin only
- **Request Body**:
  ```json
  {
    "expertId": "integer",
    "generalArea": "string",
    "specializedArea": "string"
  }
  ```
- **Response**: AIAnalysisResult object
- **Notes**: `generalArea` is required.

### Extract Skills
- **URL**: `/api/ai/extract-skills`
- **Method**: `POST`
- **Authentication**: Admin only
- **Content-Type**: `multipart/form-data`
- **Form Fields**:
  - `expertId`: ID of the expert (required)
  - `file`: Document to analyze (required)
- **Response**: AIAnalysisResult object

### Suggest Expert Panel
- **URL**: `/api/ai/suggest-panel`
- **Method**: `POST`
- **Authentication**: Authenticated user
- **Request Body**:
  ```json
  {
    "projectName": "string",
    "iscedFieldId": "integer",
    "numExperts": "integer"
  }
  ```
- **Response**:
  ```json
  {
    "experts": [Expert, ...],
    "count": "integer"
  }
  ```
- **Notes**: `projectName` is required. If `numExperts` is not provided, defaults to 3.

## ISCED Reference Data

### Get ISCED Levels
- **URL**: `/api/isced/levels`
- **Method**: `GET`
- **Authentication**: None
- **Response**: Array of ISCEDLevel objects

### Get ISCED Fields
- **URL**: `/api/isced/fields`
- **Method**: `GET`
- **Authentication**: None
- **Response**: Array of ISCEDField objects

## Statistics

### Get All Statistics
- **URL**: `/api/statistics`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: Statistics object

### Get Nationality Statistics
- **URL**: `/api/statistics/nationality`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**:
  ```json
  {
    "total": "integer",
    "stats": [
      { "name": "Bahraini", "count": "integer", "percentage": "number" },
      { "name": "Non-Bahraini", "count": "integer", "percentage": "number" }
    ]
  }
  ```

### Get ISCED Statistics
- **URL**: `/api/statistics/isced`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: Array of AreaStat objects

### Get Engagement Statistics
- **URL**: `/api/statistics/engagements`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: Array of AreaStat objects

### Get Growth Statistics
- **URL**: `/api/statistics/growth`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Query Parameters**:
  - `months`: Number of months to include (default: 12)
- **Response**: Array of GrowthStat objects

## User Management

### Create User
- **URL**: `/api/users`
- **Method**: `POST`
- **Authentication**: Admin only
- **Request Body**:
  ```json
  {
    "name": "string",
    "email": "string",
    "password": "string",
    "role": "string",
    "isActive": "boolean"
  }
  ```
- **Response**:
  ```json
  {
    "id": "integer",
    "success": "boolean",
    "message": "string"
  }
  ```
- **Notes**: `role` must be "admin" or "user". If not provided, defaults to "user".

### List Users
- **URL**: `/api/users`
- **Method**: `GET`
- **Authentication**: Admin only
- **Query Parameters**:
  - `limit`: Number of results per page (default: 10)
  - `offset`: Number of results to skip (default: 0)
- **Response**: Array of User objects

### Get User
- **URL**: `/api/users/{id}`
- **Method**: `GET`
- **Authentication**: Authenticated user
- **Response**: User object

### Update User
- **URL**: `/api/users/{id}`
- **Method**: `PUT`
- **Authentication**: Admin only
- **Request Body**:
  ```json
  {
    "name": "string",
    "email": "string",
    "password": "string",
    "role": "string",
    "isActive": "boolean"
  }
  ```
- **Response**:
  ```json
  {
    "id": "integer",
    "success": "boolean",
    "message": "string"
  }
  ```
- **Notes**: Only provided fields will be updated.

### Delete User
- **URL**: `/api/users/{id}`
- **Method**: `DELETE`
- **Authentication**: Admin only
- **Response**:
  ```json
  {
    "success": "boolean",
    "message": "string"
  }
  ```

## Authentication and Authorization

All protected endpoints require a JWT token in the `Authorization` header: `Bearer <token>`.

- **User Authentication**: Endpoints marked with "Authenticated user" require a valid JWT token.
- **Admin Authentication**: Endpoints marked with "Admin only" require a valid JWT token with a role of "admin".
