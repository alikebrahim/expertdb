# ExpertDB System API Reference

**Date**: July 22, 2025  
**Version**: 1.0  
**Purpose**: This document provides comprehensive documentation for ExpertDB system endpoints including statistics, backup, and health check functionality.

## Overview

ExpertDB provides system-level endpoints for monitoring, statistics, and data backup. These endpoints support administrative functions and system health monitoring with a focus on simplicity and maintainability.

## Table of Contents

1. [Statistics Endpoint](#statistics-endpoint)
   - GET /api/statistics
2. [Backup Endpoint](#backup-endpoint)
   - GET /api/backup
3. [Health Check Endpoint](#health-check-endpoint)
   - GET /api/health

## Statistics Endpoint

### GET /api/statistics

**Purpose**: Retrieves comprehensive system statistics in a single consolidated response. This endpoint provides all statistical data needed for dashboards and reports.

**Method**: GET  
**Path**: `/api/statistics`  
**Access Control**: All authenticated users (previously restricted to super users only)

#### Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| years | integer | No | 5 | Number of years for growth statistics (must be > 0) |

#### Request Headers

```
Authorization: Bearer <JWT_TOKEN>
```

#### Response Structure

**Success Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "totalExperts": 441,
    "activeCount": 393,
    "bahrainiPercentage": 60.77,
    "publishedCount": 32,
    "publishedRatio": 7.26,
    "topAreas": [
      {
        "name": "Business - Management & Marketing",
        "count": 67,
        "percentage": 15.19
      },
      {
        "name": "Information Technology",
        "count": 49,
        "percentage": 11.11
      }
    ],
    "engagementsByType": [
      {
        "name": "QP (Qualification Placement)",
        "count": 45,
        "percentage": 50.6
      },
      {
        "name": "IL (Institutional Listing)",
        "count": 44,
        "percentage": 49.4
      }
    ],
    "engagementsByStatus": [
      {
        "name": "completed",
        "count": 67,
        "percentage": 75.3
      },
      {
        "name": "active",
        "count": 22,
        "percentage": 24.7
      }
    ],
    "nationalityStats": [
      {
        "name": "Bahraini",
        "count": 268,
        "percentage": 60.77
      },
      {
        "name": "Non-Bahraini",
        "count": 173,
        "percentage": 39.23
      }
    ],
    "specializedAreas": {
      "top": [
        {
          "name": "Software Engineering",
          "count": 12,
          "percentage": 2.72
        },
        {
          "name": "Digital Marketing",
          "count": 8,
          "percentage": 1.81
        }
      ],
      "bottom": [
        {
          "name": "Quantum Computing",
          "count": 1,
          "percentage": 0.23
        },
        {
          "name": "Blockchain Technology",
          "count": 1,
          "percentage": 0.23
        }
      ]
    },
    "yearlyGrowth": [
      {
        "period": "2023",
        "count": 0,
        "growthRate": 0
      },
      {
        "period": "2024",
        "count": 0,
        "growthRate": 0
      },
      {
        "period": "2025",
        "count": 441,
        "growthRate": 0
      }
    ],
    "mostRequestedExperts": [
      {
        "expertId": 1,
        "name": "Dr. John Smith",
        "count": 8
      },
      {
        "expertId": 2,
        "name": "Prof. Sarah Johnson",
        "count": 6
      }
    ],
    "lastUpdated": "2025-07-22T17:43:21.677449311+03:00"
  }
}
```

**Error Response (500 Internal Server Error)**:

```json
{
  "error": "Failed to retrieve statistics"
}
```

#### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| totalExperts | integer | Total number of experts in the system |
| activeCount | integer | Number of active experts |
| bahrainiPercentage | float | Percentage of Bahraini experts |
| publishedCount | integer | Number of published experts |
| publishedRatio | float | Percentage of published experts |
| topAreas | array | Top general areas by expert count |
| engagementsByType | array | Distribution of engagements by type (QP/IL) |
| engagementsByStatus | array | Distribution of engagements by status |
| nationalityStats | array | Breakdown by nationality with counts and percentages |
| specializedAreas | object | Top 5 and bottom 5 specialized areas |
| yearlyGrowth | array | Year-by-year expert growth statistics |
| mostRequestedExperts | array | Experts with most engagements |
| lastUpdated | string | ISO timestamp of statistics generation |

#### Implementation Details

- **Handler**: `internal/api/handlers/statistics/statistics_handler.go`
- **Storage**: `internal/storage/sqlite/statistics.go`
- **Performance**: Typical response time 200-500ms
- **Caching**: No caching implemented; real-time data
- **Arrays**: All arrays guaranteed to be empty arrays (never null) for consistent frontend handling

#### Usage Examples

**Basic Request**:
```bash
curl -X GET "https://api.example.com/api/statistics" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Custom Time Range**:
```bash
curl -X GET "https://api.example.com/api/statistics?years=3" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### Notes

- Consolidates data from previously separate endpoints for simplification
- All authenticated users can access (previously super_user only)
- Growth statistics dynamically calculated based on years parameter
- Percentages calculated server-side for consistency

## Backup Endpoint

### GET /api/backup

**Purpose**: Generates a comprehensive backup of the database as a ZIP file containing CSV exports of all major tables.

**Method**: GET  
**Path**: `/api/backup`  
**Access Control**: Admin and super_user only

#### Request Headers

```
Authorization: Bearer <JWT_TOKEN>
```

#### Response

**Success Response (200 OK)**:

- **Content-Type**: `application/zip`
- **Content-Disposition**: `attachment; filename=expertdb_backup_YYYYMMDD_HHMMSS.zip`
- **Body**: Binary ZIP file data

**Error Response (500 Internal Server Error)**:

```json
{
  "error": "Failed to generate backup"
}
```

#### ZIP File Structure

The generated ZIP file contains the following CSV files:

1. **experts.csv**
   - All expert profiles with complete data
   - Columns: ID, Name, Designation, Affiliation, IsBahraini, IsAvailable, Rating, Role, EmploymentType, GeneralArea, GeneralAreaName, SpecializedArea, IsTrained, CVPath, ApprovalDocumentPath, Phone, Email, IsPublished, CreatedAt, UpdatedAt, OriginalRequestID

2. **expert_requests.csv**
   - All expert requests (pending, approved, rejected)
   - Columns: ID, Name, Designation, Affiliation, IsBahraini, IsAvailable, Rating, Role, EmploymentType, GeneralArea, SpecializedArea, IsTrained, CVPath, ApprovalDocumentPath, Phone, Email, IsPublished, Status, RejectionReason, CreatedAt, ReviewedAt, ReviewedBy, CreatedBy

3. **engagements.csv**
   - All expert engagements
   - Columns: ID, ExpertID, EngagementType, StartDate, EndDate, ProjectName, Status, FeedbackScore, Notes, CreatedAt

4. **documents.csv**
   - All uploaded documents metadata
   - Columns: ID, ExpertID, DocumentType, Filename, FilePath, ContentType, FileSize, UploadDate

5. **expert_areas.csv**
   - All general specialization areas
   - Columns: ID, Name

#### Implementation Details

- **Handler**: `internal/api/handlers/backup/backup_handler.go`
- **Method**: Uses Go's `archive/zip` package
- **Temporary Files**: Creates temp directory, cleaned up after generation
- **Memory**: Streams data to avoid loading entire database into memory
- **Filename**: Includes timestamp for uniqueness (e.g., `expertdb_backup_20250722_143021.zip`)

#### Usage Example

```bash
curl -X GET "https://api.example.com/api/backup" \
  -H "Authorization: Bearer <ADMIN_JWT_TOKEN>" \
  -o expertdb_backup.zip
```

#### Notes

- Backup includes all data regardless of status or visibility
- CSV format chosen for compatibility with spreadsheet applications
- Timestamps formatted in RFC3339 format
- Zero values for timestamps displayed as empty strings
- File paths in CSVs are relative paths as stored in database

## Health Check Endpoint

### GET /api/health

**Purpose**: Provides a simple health check endpoint to verify the API is running and accessible.

**Method**: GET  
**Path**: `/api/health`  
**Access Control**: Public (no authentication required)

#### Response

**Success Response (200 OK)**:

```json
{
  "status": "ok",
  "message": "ExpertDB API is running"
}
```

#### Implementation Details

- **Handler**: `internal/api/server.go` (handleHealth method)
- **Middleware**: Bypasses authentication middleware
- **CORS**: Allows all origins
- **Logging**: Minimal logging to avoid noise

#### Usage Example

```bash
curl -X GET "https://api.example.com/api/health"
```

#### Use Cases

1. **Load Balancer Health Checks**: Configure load balancers to monitor this endpoint
2. **Uptime Monitoring**: External monitoring services can poll this endpoint
3. **Deployment Verification**: Verify API is accessible after deployment
4. **Network Debugging**: Test basic connectivity without authentication

#### Notes

- Intentionally simple to ensure fast response times
- Does not check database connectivity (by design for simplicity)
- Always returns 200 OK if the server is running
- No rate limiting applied to this endpoint

## Error Handling

All system endpoints follow the standard ExpertDB error response format:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common HTTP status codes:
- `200 OK`: Successful request
- `400 Bad Request`: Invalid parameters
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `500 Internal Server Error`: Server-side error

## Performance Considerations

1. **Statistics Endpoint**:
   - No caching implemented; queries run in real-time
   - Consider adding caching if performance becomes an issue
   - Database indexes exist on filtered columns

2. **Backup Endpoint**:
   - Generates files on-demand; no pre-generated backups
   - Uses streaming to minimize memory usage
   - Large databases may take several seconds

3. **Health Check**:
   - Minimal processing for fast response
   - Typical response time < 10ms
   - No database queries performed

## Security Notes

1. **Authentication**: Statistics and backup require valid JWT tokens
2. **Authorization**: Backup restricted to admin/super_user roles
3. **CORS**: All origins allowed (suitable for intranet use)
4. **Rate Limiting**: Not implemented (internal tool assumption)

## Future Enhancements

1. **Statistics**:
   - Add caching with configurable TTL
   - Include phase planning statistics

2. **Backup**:
   - Add incremental backup support
   - Include phase planning data
   - Support for encrypted backups

3. **Health Check**:
   - Optional database connectivity check
   - Include version information
   - Add detailed system status

## Conclusion

These system endpoints provide essential administrative and monitoring capabilities for ExpertDB. The consolidated statistics endpoint simplifies frontend development while the backup functionality ensures data can be exported for analysis or archival purposes. The health check enables reliable deployment and monitoring practices.