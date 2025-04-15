# ExpertDB API Documentation

This document provides comprehensive documentation for the ExpertDB API, which serves as the backend for the expert database management system.

## Base URL

The API base URL is `/api` in the standard configuration.

## Authentication

### Authentication Endpoints

#### POST /api/auth/login

- **Description**: Authenticates a user and returns a JWT token
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
    "success": true,
    "message": "Login successful",
    "data": {
      "token": "string",
      "user": {
        "id": 1,
        "email": "string",
        "name": "string",
        "role": "string",
        "isActive": true,
        "createdAt": "string",
        "lastLogin": "string"
      }
    }
  }
  ```

### Authorization

All authenticated endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer {token}
```

Token validity is set to 24 hours.

### Authorization Levels

- **Unauthenticated**: Only login endpoint is accessible
- **User Role**: Most read operations and some write operations
- **Admin Role**: Full system access, including user management and administrative functions

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "message": "string",
  "data": any
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "data": null
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "string",
  "data": {
    "data": [...],
    "total": 100,
    "page": 1,
    "limit": 10,
    "totalPages": 10
  }
}
```

## User Management Endpoints

### GET /api/users

- **Description**: List users with pagination
- **Authentication**: Required (Admin only)
- **Query Parameters**:
  - `limit` (number, optional): Results per page, default: 10
  - `offset` (number, optional): Starting position, default: 0
  - `role` (string, optional): Filter by user role
  - `isActive` (boolean, optional): Filter by active status
- **Response**: List of users (password hash excluded)

### GET /api/users/:id

- **Description**: Get detailed information about a user
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): User ID
- **Response**: Detailed user data (password hash excluded)

### POST /api/users

- **Description**: Create a new user
- **Authentication**: Required (Admin only)
- **Request Body**:
  ```json
  {
    "email": "string",
    "name": "string",
    "role": "string",
    "password": "string",
    "isActive": boolean
  }
  ```
- **Response**: Created user data with ID (password excluded)

### PUT /api/users/:id

- **Description**: Update an existing user
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): User ID
- **Request Body**: User details to update (all fields optional)
- **Response**: Success confirmation

### DELETE /api/users/:id

- **Description**: Delete a user
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): User ID
- **Response**: Success confirmation

## Expert Endpoints

### GET /api/experts

- **Description**: List experts with pagination and filtering
- **Authentication**: Required
- **Query Parameters**:
  - `page` (number, optional): Page number, default: 1
  - `limit` (number, optional): Results per page, default: 10
  - `name` (string, optional): Filter by expert name
  - `area` (string, optional): Filter by expert area
  - `is_available` (boolean, optional): Filter by availability
  - `role` (string, optional): Filter by role
  - `min_rating` (number, optional): Filter by minimum rating
  - `sort_by` (string, optional): Field to sort by (name, institution, role, created_at, rating, general_area)
  - `sort_order` (string, optional): "asc" or "desc", default: "asc"
- **Response**: Paginated list of experts

### GET /api/experts/:id

- **Description**: Get detailed information about an expert
- **Authentication**: Required
- **Path Parameters**:
  - `id` (number): Expert ID
- **Response**: Detailed expert data

### POST /api/experts

- **Description**: Create a new expert
- **Authentication**: Required (Admin only)
- **Request Body**: FormData containing expert details and optional CV file
  - `name` (string): Expert name
  - `affiliation` (string): Current affiliation
  - `primaryContact` (string): Primary contact information
  - `contactType` (string): Type of contact
  - `skills` (string[]): Array of skills
  - `role` (string): Expert role
  - `employmentType` (string): Employment type
  - `generalArea` (number): General expertise area ID
  - `biography` (string): Expert bio
  - `isBahraini` (boolean): Nationality indicator
  - `availability` (string): Availability status
  - `cvFile` (file, optional): CV document
- **Response**: Created expert data

### PUT /api/experts/:id

- **Description**: Update an existing expert
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): Expert ID
- **Request Body**: FormData containing expert details to update
- **Response**: Updated expert data

### DELETE /api/experts/:id

- **Description**: Delete an expert
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): Expert ID
- **Response**: Success confirmation

### GET /api/experts/:id/approval-pdf

- **Description**: Generate and download an expert approval PDF
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): Expert ID
- **Response**: PDF file (binary)

### GET /api/experts/:id/engagements

- **Description**: Get a list of engagements for an expert
- **Authentication**: Required
- **Path Parameters**:
  - `id` (number): Expert ID
- **Query Parameters**:
  - `page` (number, optional): Page number, default: 1
  - `limit` (number, optional): Results per page, default: 10
- **Response**: Paginated list of engagements

### GET /api/experts/:id/documents

- **Description**: Get a list of documents for an expert
- **Authentication**: Required
- **Path Parameters**:
  - `id` (number): Expert ID
- **Response**: List of documents

## Expert Request Endpoints

### GET /api/expert-requests

- **Description**: List expert requests with pagination and filtering
- **Authentication**: Required
- **Query Parameters**:
  - `status` (string, optional): Filter by request status
  - `limit` (number, optional): Results per page, default: 100
  - `offset` (number, optional): Starting position, default: 0
- **Response**: Paginated list of expert requests

### GET /api/expert-requests/:id

- **Description**: Get detailed information about an expert request
- **Authentication**: Required
- **Path Parameters**:
  - `id` (number): Request ID
- **Response**: Detailed request data

### POST /api/expert-requests

- **Description**: Create a new expert request
- **Authentication**: Required
- **Request Body**:
  ```json
  {
    "requestorName": "string",
    "requestorEmail": "string",
    "organizationName": "string",
    "projectName": "string",
    "projectDescription": "string",
    "expertiseRequired": "string",
    "timeframe": "string"
  }
  ```
- **Response**: Created request data

### PUT /api/expert-requests/:id

- **Description**: Update an existing expert request
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): Request ID
- **Request Body**: Request details to update
  - `status` (string): When changed to "approved", creates a new expert
- **Response**: Updated request data

### DELETE /api/expert-requests/:id

- **Description**: Delete an expert request
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): Request ID
- **Response**: Success confirmation

## Document Endpoints

### POST /api/documents

- **Description**: Upload a document
- **Authentication**: Required
- **Request Body**: FormData
  - `file` (file): Document file
  - `expertId` (number): Associated expert ID
  - `documentType` (string): Type of document
- **Response**: Document metadata

### GET /api/documents/:id

- **Description**: Get document metadata
- **Authentication**: Required
- **Path Parameters**:
  - `id` (number): Document ID
- **Response**: Document metadata

### DELETE /api/documents/:id

- **Description**: Delete a document
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): Document ID
- **Response**: Success confirmation

## Engagement Endpoints

### GET /api/engagements

- **Description**: List engagements with pagination and filtering
- **Authentication**: Required
- **Query Parameters**:
  - `page` (number, optional): Page number, default: 1
  - `limit` (number, optional): Results per page, default: 10
  - `status` (string, optional): Filter by status
  - `engagementType` (string, optional): Filter by type
- **Response**: Paginated list of engagements

### GET /api/engagements/:id

- **Description**: Get detailed information about an engagement
- **Authentication**: Required
- **Path Parameters**:
  - `id` (number): Engagement ID
- **Response**: Detailed engagement data

### POST /api/engagements

- **Description**: Create a new engagement
- **Authentication**: Required (Admin only)
- **Request Body**:
  ```json
  {
    "expertId": number,
    "requestId": number,
    "title": "string",
    "description": "string",
    "engagementType": "string",
    "status": "string",
    "startDate": "string",
    "endDate": "string",
    "contactPerson": "string",
    "contactEmail": "string",
    "organizationName": "string",
    "notes": "string"
  }
  ```
- **Response**: Created engagement data

### PUT /api/engagements/:id

- **Description**: Update an existing engagement
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): Engagement ID
- **Request Body**: Engagement details to update
- **Response**: Updated engagement data

### DELETE /api/engagements/:id

- **Description**: Delete an engagement
- **Authentication**: Required (Admin only)
- **Path Parameters**:
  - `id` (number): Engagement ID
- **Response**: Success confirmation

## Expert Areas Endpoint

### GET /api/expert/areas

- **Description**: Get a list of expert areas
- **Authentication**: Optional
- **Response**: List of expert area objects with ID, name, and description

## Statistics Endpoints

### GET /api/statistics

- **Description**: Get overall system statistics
- **Authentication**: Required
- **Response**:
  ```json
  {
    "success": true,
    "data": {
      "totalExperts": 100,
      "totalBahraini": 60,
      "totalInternational": 40,
      "totalEngagements": 150,
      "byEmploymentType": {
        "Full-time": 50,
        "Part-time": 30,
        "Consultant": 20
      },
      "byAvailability": {
        "Available": 70,
        "Limited": 20,
        "Unavailable": 10
      }
    }
  }
  ```

### GET /api/statistics/nationality

- **Description**: Get nationality distribution statistics
- **Authentication**: Required
- **Response**:
  ```json
  {
    "success": true,
    "data": {
      "total": 100,
      "stats": [
        {"name": "Bahraini", "count": 60, "percentage": 60},
        {"name": "Non-Bahraini", "count": 40, "percentage": 40}
      ]
    }
  }
  ```

### GET /api/statistics/growth

- **Description**: Get expert growth statistics over time
- **Authentication**: Required
- **Query Parameters**:
  - `months` (number, optional): Number of months to include, default: 12
- **Response**: Array of monthly data points with new and total experts

### GET /api/statistics/engagements

- **Description**: Get engagement statistics
- **Authentication**: Required
- **Response**:
  ```json
  {
    "success": true,
    "data": {
      "total": 150,
      "byStatus": {
        "Active": 50,
        "Completed": 80,
        "Cancelled": 20
      },
      "byType": {
        "Consultation": 70,
        "Project": 50,
        "Workshop": 30
      }
    }
  }
  ```

## Error Codes

The API uses standard HTTP status codes:

- 200: OK
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## Data Types

### User
```json
{
  "id": number,
  "email": "string",
  "name": "string",
  "role": "string",
  "isActive": boolean,
  "createdAt": "string",
  "lastLogin": "string"
}
```

### Expert
```json
{
  "id": number,
  "name": "string",
  "affiliation": "string",
  "primaryContact": "string",
  "contactType": "string",
  "skills": ["string"],
  "role": "string",
  "employmentType": "string",
  "generalArea": number,
  "cvPath": "string",
  "biography": "string",
  "isBahraini": boolean,
  "availability": "string",
  "rating": number,
  "created_at": "string",
  "updated_at": "string"
}
```

### ExpertRequest
```json
{
  "id": number,
  "requestorId": number,
  "requestorName": "string",
  "requestorEmail": "string",
  "organizationName": "string",
  "projectName": "string",
  "projectDescription": "string",
  "expertiseRequired": "string",
  "timeframe": "string",
  "status": "string",
  "notes": "string",
  "createdAt": "string",
  "updatedAt": "string"
}
```

### Document
```json
{
  "id": number,
  "expertId": number,
  "filename": "string",
  "originalFilename": "string",
  "documentType": "string",
  "contentType": "string",
  "size": number,
  "uploadedBy": number,
  "uploadedAt": "string"
}
```

### Engagement
```json
{
  "id": number,
  "expertId": number,
  "requestId": number,
  "title": "string",
  "description": "string",
  "engagementType": "string",
  "status": "string",
  "startDate": "string",
  "endDate": "string",
  "contactPerson": "string",
  "contactEmail": "string",
  "organizationName": "string",
  "notes": "string",
  "createdAt": "string",
  "updatedAt": "string"
}
```