# Expert Management API Reference

**Date**: January 22, 2025  
**Version**: 1.1 (Updated for Document Management System)  
**Context**: This document provides comprehensive API documentation for all Expert Management endpoints in the ExpertDB system, extracted from the main API_REFERENCE.md file. Updated to reflect the migration from file path-based document storage to a centralized document management system using foreign key relationships.

## Table of Contents

1. [Overview](#overview)
2. [Data Model](#data-model)
3. [Expert Management Endpoints](#expert-management-endpoints)
   - [GET /api/experts](#get-apiexperts)
   - [GET /api/experts/{id}](#get-apiexpertsid)
   - [POST /api/experts](#post-apiexperts)
   - [PUT /api/experts/{id}](#put-apiexpertsid)
   - [DELETE /api/experts/{id}](#delete-apiexpertsid)
4. [Expert Areas Endpoints](#expert-areas-endpoints)
   - [GET /api/expert/areas](#get-apiexpertareas)
   - [POST /api/expert/areas](#post-apiexpertareas)
   - [PUT /api/expert/areas/{id}](#put-apiexpertareasid)
5. [Specialized Areas Endpoints](#specialized-areas-endpoints)
   - [GET /api/specialized-areas](#get-apispecialized-areas)

## Overview

The Expert Management system in ExpertDB provides comprehensive functionality for managing expert profiles, their specialization areas, and related metadata. The system supports:

- Full CRUD operations for expert profiles
- Advanced filtering and sorting capabilities with multi-value support
- Normalized specialized areas management
- File upload support for CVs and approval documents
- Integration with the expert request workflow

### Key Features

- **Multi-Value Filtering**: Support for comma-separated values in filter parameters (v1.5 update)
- **Normalized Data Model**: Specialized areas are stored as IDs with separate lookup table
- **Dual Content-Type Support**: Expert updates support both JSON and multipart/form-data
- **Comprehensive Metadata**: Includes professional experience, education, and documents

### Authentication & Authorization

- All endpoints require JWT authentication
- Expert creation/modification requires admin role
- Expert viewing is available to all authenticated users
- Area management requires admin role

## Data Model

### Expert Entity

The expert entity includes the following fields:

```json
{
  "id": int,                      // Auto-generated database ID
  "expertId": "string",           // Business ID (e.g., "EXP-0001")
  "name": "string",               // Expert's full name
  "designation": "string",        // Professional title
  "affiliation": "string",        // Organization/affiliation
  "isBahraini": boolean,          // Nationality flag
  "isAvailable": boolean,         // Availability status
  "rating": int,                  // Performance rating (1-5)
  "role": "string",               // Expert role (validator/evaluator)
  "employmentType": "string",     // Employment category
  "generalAreaName": "string",    // General specialization area name
  "specializedAreaNames": "string", // Comma-separated specialized area names
  "isTrained": boolean,           // BQA training status
  "cvDocumentId": int,            // CV document reference ID
  "phone": "string",              // Contact phone
  "email": "string",              // Contact email
  "isPublished": boolean,         // Publication status
  "approvalDocumentId": int,      // Approval document reference ID
  "experienceEntries": [...],     // Professional experience array
  "educationEntries": [...],      // Educational background array
  "createdAt": "string",          // Creation timestamp
  "updatedAt": "string"           // Last update timestamp
}
```

### Professional Experience Structure

```json
{
  "id": number,
  "organization": "string",
  "position": "string", 
  "startDate": "string",
  "endDate": "string",
  "isCurrent": boolean,
  "country": "string",
  "description": "string",
  "createdAt": "string",
  "updatedAt": "string"
}
```

### Educational Background Structure

```json
{
  "id": number,
  "institution": "string",
  "degree": "string",
  "fieldOfStudy": "string", 
  "graduationYear": "string",
  "country": "string",
  "description": "string",
  "createdAt": "string",
  "updatedAt": "string"
}
```

### Specialized Areas

The system uses a normalized approach for specialized areas:
- Expert records store specialized area IDs as comma-separated values (e.g., "1,4,6")
- A separate `specialized_areas` table maintains the ID-to-name mapping
- This allows for efficient searching and consistent area naming

## Expert Management Endpoints

### GET /api/experts

Retrieves a paginated list of experts with enhanced multi-value filtering and sorting capabilities.

#### Request

- **Method**: GET
- **Path**: `/api/experts`
- **Headers**: 
  - `Authorization: Bearer <token>`

#### Query Parameters

**Pagination & Sorting:**
- `limit` - Number of results per page (default: 100)
- `offset` - Number of results to skip (default: 0)
- `sort_by` - Sort field options:
  - `name`, `institution`, `role`, `created_at`, `updated_at`
  - `rating`, `general_area`, `designation`, `employment_type`
  - `specialized_area`, `is_bahraini`, `is_available`, `is_published`
- `sort_order` - Sort direction: `asc` or `desc` (default: `asc`)

**Multi-Value Filtering (supports comma-separated values):**
- `general_area` - General area ID(s) (e.g., `3` or `3,5,12`)
- `affiliation` - Institution/affiliation text search (e.g., `University` or `University,Polytechnic`)
- `role` - Expert role(s) (e.g., `validator` or `validator,evaluator`)
- `employment_type` - Employment type(s) (e.g., `Academic` or `Academic,Employer`)
- `specialized_area` - Specialized area text search (supports multiple comma-separated values)

**Boolean Filters (single value only):**
- `is_available` - Availability status (`true` or `false`)
- `is_bahraini` - Nationality filter (`true` or `false`)
- `is_published` - Publication status (`true` or `false`)

#### Filter Logic

- **Within same parameter**: OR logic (e.g., `role=validator,evaluator` finds experts who are validators OR evaluators)
- **Between different parameters**: AND logic (e.g., `role=validator&general_area=3` finds experts who are validators AND in area 3)

#### Examples

```bash
# Single value filters
GET /api/experts?general_area=3
GET /api/experts?role=validator
GET /api/experts?affiliation=University

# Multi-value filters (OR within field)
GET /api/experts?role=validator,evaluator
GET /api/experts?general_area=3,5,12
GET /api/experts?employment_type=Academic,Employer

# Combined filters (AND between fields)
GET /api/experts?role=validator&general_area=3
GET /api/experts?role=validator,evaluator&affiliation=University&is_available=true

# With sorting and pagination
GET /api/experts?role=validator&sort_by=rating&sort_order=desc&limit=20&offset=0
```

#### Response Headers

- `X-Total-Count` - Total number of experts matching filters
- `X-Total-Pages` - Total number of pages
- `X-Current-Page` - Current page number
- `X-Page-Size` - Number of items per page
- `X-Has-Next-Page` - Boolean indicating if next page exists
- `X-Has-Prev-Page` - Boolean indicating if previous page exists

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "data": {
    "experts": [
      {
        "id": 1,
        "expertId": "EXP-0001",
        "name": "Dr. John Smith",
        "designation": "Professor",
        "affiliation": "University of Bahrain",
        "isBahraini": true,
        "isAvailable": true,
        "rating": 5,
        "role": "validator",
        "employmentType": "Academic",
        "generalAreaName": "Business - Management & Marketing",
        "specializedAreaNames": "Software Engineering, Database Design",
        "isTrained": true,
        "cvDocumentId": 123,
        "phone": "+973 12345678",
        "email": "john.smith@example.com",
        "isPublished": true,
        "experienceEntries": [
          {
            "id": 1,
            "organization": "Tech Company",
            "position": "Senior Software Engineer",
            "startDate": "2020-01",
            "endDate": "Present",
            "isCurrent": true,
            "country": "Bahrain",
            "description": "Led development team",
            "createdAt": "2025-01-20T10:00:00Z",
            "updatedAt": "2025-01-20T10:00:00Z"
          }
        ],
        "educationEntries": [
          {
            "id": 1,
            "institution": "MIT",
            "degree": "PhD",
            "fieldOfStudy": "Computer Science",
            "graduationYear": "2015",
            "country": "USA",
            "description": "Focus on AI/ML",
            "createdAt": "2025-01-20T10:00:00Z",
            "updatedAt": "2025-01-20T10:00:00Z"
          }
        ],
        "approvalDocumentId": 124,
        "createdAt": "2025-01-20T10:00:00Z",
        "updatedAt": "2025-01-21T15:30:00Z"
      }
    ],
    "pagination": {
      "totalCount": 441,
      "totalPages": 5,
      "currentPage": 1,
      "pageSize": 100,
      "hasNextPage": true,
      "hasPrevPage": false,
      "hasMore": true
    }
  }
}
```

**Error (400 Bad Request):**
```json
{
  "error": "Invalid query parameters"
}
```

#### Implementation Notes

- **Performance**: Single SQL query with IN clauses maintains optimal performance (1-2 second response times)
- **Database**: Uses indexed columns for filtered fields
- **Multi-Value Processing**: Helper functions `parseMultiValue()`, `buildInClause()`, `buildLikeClause()`
- **Breaking Changes (v1.5)**: Legacy parameter names (`by_role`, `by_general_area`, etc.) no longer supported

### GET /api/experts/{id}

Retrieves detailed information for a specific expert.

#### Request

- **Method**: GET
- **Path**: `/api/experts/{id}`
- **Headers**: 
  - `Authorization: Bearer <token>`
- **Parameters**:
  - `id` (path) - Expert ID

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "expertId": "EXP-0001",
    "name": "Dr. John Smith",
    "designation": "Professor",
    "affiliation": "University of Bahrain",
    "isBahraini": true,
    "isAvailable": true,
    "rating": 5,
    "role": "validator",
    "employmentType": "Academic",
    "generalAreaName": "Business - Management & Marketing",
    "specializedAreaNames": "Software Engineering, Database Design",
    "isTrained": true,
    "cvDocumentId": 123,
    "phone": "+973 12345678",
    "email": "john.smith@example.com",
    "isPublished": true,
    "approvalDocumentId": 124,
    "createdAt": "2025-01-20T10:00:00Z",
    "updatedAt": "2025-01-21T15:30:00Z"
  }
}
```

**Error (404 Not Found):**
```json
{
  "error": "Expert not found"
}
```

### POST /api/experts

Creates a new expert profile.

#### Request

- **Method**: POST
- **Path**: `/api/experts`
- **Headers**: 
  - `Authorization: Bearer <admin_token>`
  - `Content-Type: application/json`

#### Request Payload

```json
{
  "name": "Dr. Jane Doe",              // Required
  "affiliation": "Tech University",     // Required
  "email": "jane.doe@example.com",     // Required
  "designation": "Associate Professor", // Required
  "isBahraini": false,                 // Required
  "isAvailable": true,                 // Required
  "rating": "4",                       // Required
  "role": "evaluator",                 // Required
  "employmentType": "Academic",        // Required
  "generalArea": 3,                    // Required (area ID)
  "specializedArea": "1,4,6",          // Required (comma-separated IDs)
  "isTrained": true,                   // Required
  "cvDocumentId": 123,                 // Required (document reference ID)
  "phone": "+973 87654321",            // Required
  "isPublished": false,                // Required
  "experienceEntries": [               // Optional
    {
      "organization": "Previous University",
      "position": "Assistant Professor",
      "startDate": "2015-09",
      "endDate": "2020-08",
      "isCurrent": false,
      "country": "UK",
      "description": "Teaching and research in computer science"
    }
  ],
  "educationEntries": [                // Optional
    {
      "institution": "Oxford University",
      "degree": "PhD",
      "fieldOfStudy": "Computer Science",
      "graduationYear": "2015",
      "country": "UK",
      "description": "Research in distributed systems"
    }
  ],
  "approvalDocumentId": 124           // Required (document reference ID)
}
```

#### Response Payload

**Success (201 Created):**
```json
{
  "success": true,
  "message": "Expert created successfully",
  "data": {
    "id": 442,
    "expertId": "EXP-0442"
  }
}
```

**Error (400 Bad Request):**
```json
{
  "errors": ["name is required", "invalid general_area"]
}
```

**Error (409 Conflict):**
```json
{
  "error": "Expert ID already exists"
}
```

#### Business Rules

- Generates unique expert ID (e.g., `EXP-0001`)
- Validates all required fields
- Ensures general area ID exists
- Validates specialized area IDs if provided

### PUT /api/experts/{id}

Updates an expert profile with support for both JSON and file uploads.

#### Request

- **Method**: PUT
- **Path**: `/api/experts/{id}`
- **Headers**: 
  - `Authorization: Bearer <admin_token>`
- **Parameters**:
  - `id` (path) - Expert ID

#### Content-Type Support

The endpoint supports two content types:

1. **JSON Updates** - For data-only updates
   - `Content-Type: application/json`
   
2. **Multipart Updates** - For updates with file attachments
   - `Content-Type: multipart/form-data`

#### JSON Request (Data Only)

**Headers:**
- `Content-Type: application/json`

**Request Payload:**
```json
{
  "name": "Dr. Jane Smith",
  "affiliation": "New University",
  "email": "jane.smith@newuni.com",
  "designation": "Full Professor",
  "isBahraini": false,
  "isAvailable": true,
  "rating": 5,
  "role": "validator",
  "employmentType": "Academic",
  "generalArea": 5,
  "specializedArea": "2,5,8",
  "isTrained": true,
  "cvDocumentId": 123,
  "phone": "+973 87654321",
  "isPublished": true,
  "approvalDocumentId": 124
}
```

#### Multipart Request (With File Uploads)

**Headers:**
- `Content-Type: multipart/form-data`

**Form Fields:**
- `data` - JSON string containing expert data (same structure as JSON request)
- `cvFile` - (optional) New CV file (PDF format) - creates new document record
- `approvalDocument` - (optional) New approval document file - creates new document record

**Example:**
```bash
curl -X PUT http://api.example.com/api/experts/442 \
  -H "Authorization: Bearer <token>" \
  -F "data={\"name\":\"Dr. Jane Smith\",\"rating\":5}" \
  -F "cvFile=@new_cv.pdf" \
  -F "approvalDocument=@new_approval.pdf"
```

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "message": "Expert updated successfully",
  "data": {
    "id": 442
  }
}
```

**Error (400 Bad Request):**
```json
{
  "error": "Invalid request payload"
}
```

**Error (404 Not Found):**
```json
{
  "error": "Expert not found"
}
```

#### Implementation Notes

- **Smart Content-Type Detection**: Automatically detects JSON vs multipart requests
- **File Handling**: New files automatically processed via document service and document IDs updated
- **Document References**: Uses document IDs instead of file paths for data integrity
- **Partial Updates**: Only provided fields are updated
- **Type Conversion**: Handles rating as int
- **Backward Compatibility**: JSON-only clients continue to work unchanged

### DELETE /api/experts/{id}

Deletes an expert and all associated documents.

#### Request

- **Method**: DELETE
- **Path**: `/api/experts/{id}`
- **Headers**: 
  - `Authorization: Bearer <admin_token>`
- **Parameters**:
  - `id` (path) - Expert ID

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "message": "Expert deleted successfully"
}
```

**Error (404 Not Found):**
```json
{
  "error": "Expert not found"
}
```

#### Business Rules

- Cascades deletion to all associated documents
- Removes expert from any phase applications
- Cannot be undone

## Expert Areas Endpoints

### GET /api/expert/areas

Retrieves all general specialization areas.

#### Request

- **Method**: GET
- **Path**: `/api/expert/areas`
- **Headers**: 
  - `Authorization: Bearer <token>`

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "data": [
    { "id": 1, "name": "Accounting" },
    { "id": 2, "name": "Applied Sciences - Health and Sports" },
    { "id": 3, "name": "Business - Management & Marketing" },
    { "id": 4, "name": "Computing - Graphic Design" },
    { "id": 5, "name": "Computing - Information Communication Technology" }
  ]
}
```

**Error (401 Unauthorized):**
```json
{
  "error": "Unauthorized"
}
```

#### Access Control

- Available to all authenticated users
- Read-only for non-admin users

### POST /api/expert/areas

Creates a new general specialization area.

#### Request

- **Method**: POST
- **Path**: `/api/expert/areas`
- **Headers**: 
  - `Authorization: Bearer <admin_token>`
  - `Content-Type: application/json`

#### Request Payload

```json
{
  "name": "Quantum Computing"  // Required
}
```

#### Response Payload

**Success (201 Created):**
```json
{
  "success": true,
  "message": "Area created successfully",
  "data": {
    "id": 35,
    "name": "Quantum Computing"
  }
}
```

**Error (400 Bad Request):**
```json
{
  "error": "Name is required"
}
```

**Error (409 Conflict):**
```json
{
  "error": "Area name already exists"
}
```

### PUT /api/expert/areas/{id}

Renames an existing general specialization area.

#### Request

- **Method**: PUT
- **Path**: `/api/expert/areas/{id}`
- **Headers**: 
  - `Authorization: Bearer <admin_token>`
  - `Content-Type: application/json`
- **Parameters**:
  - `id` (path) - Area ID

#### Request Payload

```json
{
  "name": "Quantum & Classical Computing"  // Required
}
```

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "message": "Area updated successfully",
  "data": {
    "id": 35,
    "name": "Quantum & Classical Computing"
  }
}
```

**Error (400 Bad Request):**
```json
{
  "error": "Name is required"
}
```

**Error (404 Not Found):**
```json
{
  "error": "Area not found"
}
```

#### Business Rules

- Updates cascade to all experts and expert requests
- Uses database transaction for consistency
- Area name must be unique

## Specialized Areas Endpoints

### GET /api/specialized-areas

Retrieves all specialized areas for search and selection functionality.

#### Request

- **Method**: GET
- **Path**: `/api/specialized-areas`
- **Headers**: 
  - `Authorization: Bearer <token>`

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "data": [
    { "id": 1, "name": "Software Engineering", "createdAt": "2025-01-15T10:00:00Z" },
    { "id": 2, "name": "Database Design", "createdAt": "2025-01-15T10:00:00Z" },
    { "id": 3, "name": "Network Security", "createdAt": "2025-01-15T10:00:00Z" },
    { "id": 4, "name": "Machine Learning", "createdAt": "2025-01-15T10:00:00Z" },
    { "id": 5, "name": "Cloud Computing", "createdAt": "2025-01-15T10:00:00Z" }
  ]
}
```

**Error (401 Unauthorized):**
```json
{
  "error": "Unauthorized"
}
```

#### Implementation Notes

- Returns all available specialized areas from normalized table
- Expert records reference these areas by ID
- Used for populating selection dropdowns in UI
- Part of the specialized areas normalization system

## Migration Guide (v1.4 to v1.5)

### Breaking Changes

The following parameter names are no longer supported:

```bash
# OLD PARAMETERS (v1.4) - NO LONGER SUPPORTED
GET /api/experts?by_role=validator&by_general_area=3&name=University

# NEW PARAMETERS (v1.5) - REQUIRED FORMAT  
GET /api/experts?role=validator&general_area=3&affiliation=University
```

### Parameter Mapping

| Old Parameter | New Parameter | Notes |
|--------------|---------------|-------|
| `by_role` | `role` | Now supports multi-value |
| `by_general_area` | `general_area` | Now supports multi-value |
| `by_employment_type` | `employment_type` | Now supports multi-value |
| `name` | `affiliation` | Fixed: now searches affiliation field |

### New Multi-Value Capabilities

```bash
# Single values still work
GET /api/experts?role=validator

# But now you can use comma-separated values
GET /api/experts?role=validator,evaluator
GET /api/experts?general_area=3,5,12
GET /api/experts?affiliation=University,College,Institute
```

## Performance Considerations

- **Response Times**: Typical response times are 1-2 seconds for filtered queries
- **Indexing**: The following columns should be indexed for optimal performance:
  - `role`, `general_area`, `employment_type`, `is_bahraini`
  - `is_available`, `is_published`, `institution` (partial text index)
- **Multi-Value Queries**: Use IN clauses rather than multiple OR conditions
- **Pagination**: Always use pagination for large result sets

## Error Handling

All endpoints follow a consistent error response pattern:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common HTTP status codes:
- `200 OK` - Successful GET, PUT, DELETE
- `201 Created` - Successful POST
- `400 Bad Request` - Invalid request data or parameters
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate resource
- `500 Internal Server Error` - Unexpected server error

## Best Practices

1. **Always use pagination** for listing endpoints to avoid performance issues
2. **Use multi-value filters** efficiently - combine filters to narrow results
3. **Include only necessary fields** in update requests (partial updates)
4. **Handle file uploads** separately from data updates when possible
5. **Cache area lists** on the client side as they change infrequently
6. **Validate specialized area IDs** before submitting expert create/update requests