# ExpertDB Engagement Management API Reference

**Date**: July 22, 2025  
**Version**: 1.0  
**Module**: Engagement Management

## Overview

The Engagement Management system in ExpertDB tracks expert assignments to various tasks, specifically as validators or evaluators. This is part of the legacy system that tracks historical engagements while the newer Phase Planning system manages future assignments.

### Key Concepts

- **Engagement Types**: Experts can serve as either "validator" or "evaluator" for projects
- **Status Tracking**: Engagements have status indicators (pending, active, completed)
- **Legacy System**: This is the original engagement tracking system, primarily for historical records
- **Import Support**: Bulk import functionality for migrating past engagement data

### Relationship to Applications

While engagements track expert assignments independently, they conceptually relate to the application types in the Phase Planning system:
- **Validator engagements** typically correspond to validation tasks
- **Evaluator engagements** typically correspond to evaluation tasks
- Future versions may integrate these systems more closely

## Table of Contents

1. [List All Engagements](#list-all-engagements)
2. [Get Expert's Engagements](#get-experts-engagements)
3. [Get Engagement Details](#get-engagement-details)
4. [Create New Engagement](#create-new-engagement)
5. [Update Engagement](#update-engagement)
6. [Delete Engagement](#delete-engagement)
7. [Import Engagements](#import-engagements)

## Authentication

All engagement endpoints require JWT authentication:
```
Authorization: Bearer <JWT_TOKEN>
```

## Access Control

- **View Operations** (GET): All authenticated users
- **Create/Update/Delete Operations** (POST, PUT, DELETE): Admin only
- **Import Operations** (POST): Admin only

## Standard Response Format

All endpoints follow a consistent response structure:

```json
{
  "success": boolean,    // true for successful requests
  "message": "string",   // optional success message
  "data": <object>       // optional response data
}
```

## Endpoints

### List All Engagements

Retrieves a paginated list of engagements with optional filtering.

**Endpoint**: `GET /api/engagements`

**Query Parameters**:
| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `limit` | integer | Number of results per page (default: 20) | `10` |
| `offset` | integer | Number of results to skip | `0` |
| `expert_id` | integer | Filter by expert ID | `123` |
| `type` | string | Filter by engagement type | `validator` or `evaluator` |

**Example Request**:
```bash
curl -X GET "https://api.expertdb.com/api/engagements?type=validator&limit=50" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "engagements": [
      {
        "id": 1,
        "expertId": 123,
        "expertName": "Dr. John Smith",
        "engagementType": "validator",
        "startDate": "2025-01-15",
        "projectName": "University X - Engineering Program",
        "status": "active",
        "notes": "Lead validator for QP review",
        "createdAt": "2025-01-10T08:30:00Z"
      }
    ],
    "pagination": {
      "limit": 50,
      "offset": 0,
      "count": 1
    },
    "filters": {
      "expertId": null,
      "engagementType": "validator"
    }
  }
}
```

### Get Expert's Engagements

Retrieves all engagements for a specific expert.

**Endpoint**: `GET /api/experts/{id}/engagements`

**Path Parameters**:
- `id`: Expert ID

**Query Parameters**:
| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `limit` | integer | Number of results per page | `20` |
| `offset` | integer | Number of results to skip | `0` |
| `type` | string | Filter by engagement type | `evaluator` |

**Example Request**:
```bash
curl -X GET "https://api.expertdb.com/api/experts/123/engagements?type=evaluator" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "engagements": [
      {
        "id": 2,
        "expertId": 123,
        "expertName": "Dr. John Smith",
        "engagementType": "evaluator",
        "startDate": "2025-02-01",
        "projectName": "College Y - Business Studies",
        "status": "pending",
        "notes": "Assigned as evaluator for IL review",
        "createdAt": "2025-01-20T10:00:00Z"
      }
    ],
    "pagination": {
      "limit": 20,
      "offset": 0,
      "count": 1
    },
    "expertId": 123
  }
}
```

### Get Engagement Details

Retrieves detailed information about a specific engagement.

**Endpoint**: `GET /api/engagements/{id}`

**Path Parameters**:
- `id`: Engagement ID

**Example Request**:
```bash
curl -X GET "https://api.expertdb.com/api/engagements/1" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "expertId": 123,
    "engagementType": "validator",
    "startDate": "2025-01-15",
    "endDate": "2025-03-15",
    "projectName": "University X - Engineering Program",
    "status": "active",
    "notes": "Lead validator for QP review. Focus on curriculum assessment.",
    "createdAt": "2025-01-10T08:30:00Z"
  }
}
```

**Error Response** (404 Not Found):
```json
{
  "error": "Engagement not found"
}
```

### Create New Engagement

Creates a new engagement record for an expert.

**Endpoint**: `POST /api/engagements`

**Access**: Admin only

**Request Body**:
```json
{
  "expertId": 123,                    // Required: Expert ID
  "engagementType": "validator",      // Required: "validator" or "evaluator"
  "startDate": "2025-01-15",         // Required: ISO format date
  "endDate": "2025-03-15",           // Optional: ISO format date
  "projectName": "University X",      // Optional: Project/institution name
  "status": "active",                // Optional: defaults to "pending"
  "notes": "Lead validator"          // Optional: Additional notes
}
```

**Validation Rules**:
- `expertId` must reference an existing expert
- `engagementType` must be either "validator" or "evaluator"
- `startDate` is required and must be a valid date
- `status` if provided must be valid (pending, active, completed)

**Example Request**:
```bash
curl -X POST "https://api.expertdb.com/api/engagements" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "expertId": 123,
    "engagementType": "validator",
    "startDate": "2025-01-15",
    "projectName": "University X - Engineering Program"
  }'
```

**Success Response** (201 Created):
```json
{
  "success": true,
  "message": "Engagement created successfully",
  "data": {
    "id": 3
  }
}
```

**Error Response** (400 Bad Request):
```json
{
  "error": "Expert ID is required"
}
```

### Update Engagement

Updates an existing engagement record.

**Endpoint**: `PUT /api/engagements/{id}`

**Access**: Admin only

**Path Parameters**:
- `id`: Engagement ID

**Request Body** (all fields optional):
```json
{
  "engagementType": "evaluator",     // Change type
  "startDate": "2025-02-01",        // Update start date
  "endDate": "2025-04-01",          // Set/update end date
  "projectName": "Updated Project",  // Update project name
  "status": "completed",            // Update status
  "notes": "Updated notes"          // Update notes
}
```

**Example Request**:
```bash
curl -X PUT "https://api.expertdb.com/api/engagements/1" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed",
    "endDate": "2025-03-15",
    "notes": "Successfully completed validation"
  }'
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Engagement updated successfully",
  "data": {
    "id": 1
  }
}
```

### Delete Engagement

Deletes an engagement record.

**Endpoint**: `DELETE /api/engagements/{id}`

**Access**: Admin only

**Path Parameters**:
- `id`: Engagement ID

**Example Request**:
```bash
curl -X DELETE "https://api.expertdb.com/api/engagements/1" \
  -H "Authorization: Bearer <ADMIN_TOKEN>"
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Engagement deleted successfully"
}
```

### Import Engagements

Bulk imports engagement records from CSV or JSON format. Useful for migrating historical data.

**Endpoint**: `POST /api/engagements/import`

**Access**: Admin only

**Content Types Supported**:
- `multipart/form-data` (CSV file upload)
- `application/json` (JSON array)

#### CSV Import

**CSV Format Requirements**:
```csv
expertId,engagementType,startDate,endDate,projectName,status,notes
123,validator,2025-01-15,2025-03-15,University X,active,Lead validator
456,evaluator,2025-02-01,,College Y,pending,Evaluator for IL
```

**Example Request**:
```bash
curl -X POST "https://api.expertdb.com/api/engagements/import" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -F "file=@engagements.csv"
```

#### JSON Import

**JSON Format**:
```json
[
  {
    "expertId": 123,
    "engagementType": "validator",
    "startDate": "2025-01-15",
    "endDate": "2025-03-15",
    "projectName": "University X",
    "status": "active",
    "notes": "Lead validator"
  },
  {
    "expertId": 456,
    "engagementType": "evaluator",
    "startDate": "2025-02-01",
    "projectName": "College Y",
    "status": "pending"
  }
]
```

**Example Request**:
```bash
curl -X POST "https://api.expertdb.com/api/engagements/import" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "expertId": 123,
      "engagementType": "validator",
      "startDate": "2025-01-15",
      "projectName": "University X"
    }
  ]'
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Import completed: 8 successful, 2 failed out of 10 total",
  "data": {
    "success": true,
    "successCount": 8,
    "failureCount": 2,
    "totalCount": 10,
    "errors": {
      "3": "Expert ID 999 not found",
      "7": "Invalid engagement type: reviewer"
    }
  }
}
```

**Import Features**:
- **Validation**: Each record is validated before import
- **Deduplication**: Prevents duplicate engagements for same expert/project/date
- **Error Reporting**: Detailed error messages for failed records
- **Partial Success**: Valid records are imported even if some fail

## Implementation Notes

### Database Schema

Engagements are stored in the `expert_engagements` table with the following structure:
- `id`: Primary key
- `expert_id`: Foreign key to experts table
- `engagement_type`: Either "validator" or "evaluator"
- `start_date`: Engagement start date
- `end_date`: Optional engagement end date
- `project_name`: Optional project/institution name
- `status`: Engagement status (pending, active, completed)
- `notes`: Optional additional notes
- `created_at`: Timestamp of record creation

### Filtering Logic

The engagement filtering system supports:
- **Expert-based filtering**: View all engagements for a specific expert
- **Type-based filtering**: Filter by validator or evaluator roles
- **Combined filtering**: Apply multiple filters simultaneously

### Status Management

Engagement statuses follow this lifecycle:
1. **pending**: Initial state when engagement is planned
2. **active**: Engagement is currently in progress
3. **completed**: Engagement has been successfully completed

### Legacy System Considerations

This engagement system is considered legacy as the newer Phase Planning system provides more structured assignment management. However, it remains important for:
- Tracking historical expert assignments
- Maintaining continuity with existing data
- Supporting simple engagement tracking needs

## Error Handling

All endpoints follow consistent error handling:

| Status Code | Description | Example |
|-------------|-------------|---------|
| 400 | Bad Request - Invalid parameters | Missing required fields |
| 401 | Unauthorized - Invalid/missing token | No authentication token |
| 403 | Forbidden - Insufficient permissions | Non-admin trying to create |
| 404 | Not Found - Resource doesn't exist | Invalid engagement ID |
| 500 | Internal Server Error | Database connection issues |

## Best Practices

1. **Filtering**: Use type filters to improve query performance when looking for specific engagement types
2. **Pagination**: Always use pagination for large datasets to improve response times
3. **Import Format**: Prefer JSON format for imports when programmatically generating data
4. **Status Updates**: Keep engagement statuses current to maintain accurate records
5. **Notes Field**: Use the notes field to provide context about the engagement

## Future Enhancements

The engagement system may be enhanced to:
- Integrate more closely with the Phase Planning system
- Support additional engagement types beyond validator/evaluator
- Include rating and performance tracking
- Provide engagement analytics and reporting
- Support engagement templates for common scenarios