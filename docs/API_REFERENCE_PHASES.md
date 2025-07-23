# ExpertDB Phase Planning API Reference

**Date**: July 22, 2025  
**Version**: 1.0  
**Purpose**: Comprehensive documentation for ExpertDB phase planning endpoints, workflows, and role-based access control.

## Table of Contents

1. [Overview](#overview)
2. [Phase Planning Workflow](#phase-planning-workflow)
3. [Role-Based Access Control](#role-based-access-control)
4. [Application Types](#application-types)
5. [API Endpoints](#api-endpoints)
   - [POST /api/phases](#post-apiphases)
   - [GET /api/phases](#get-apiphases)
   - [GET /api/phases/{id}](#get-apiphasesid)
   - [PUT /api/phases/{id}](#put-apiphasesid)
   - [PUT /api/phases/{id}/applications/{app_id}](#put-apiphasesidapplicationsapp_id)
   - [PUT /api/phases/{id}/applications/{app_id}/review](#put-apiphasesidapplicationsapp_idreview)
   - [POST /api/phases/{id}/applications/{app_id}/ratings](#post-apiphasesidapplicationsapp_idratings)
   - [GET /api/applications](#get-apiapplications)
6. [Role Assignment Endpoints](#role-assignment-endpoints)
   - [POST /api/users/{id}/planner-assignments](#post-apiusersidplanner-assignments)
   - [POST /api/users/{id}/manager-assignments](#post-apiusersidmanager-assignments)
   - [DELETE /api/users/{id}/planner-assignments](#delete-apiusersidplanner-assignments)
   - [DELETE /api/users/{id}/manager-assignments](#delete-apiusersidmanager-assignments)
   - [GET /api/users/{id}/assignments](#get-apiusersidassignments)
7. [Manager Task Endpoints](#manager-task-endpoints)
   - [GET /api/users/me/manager-tasks](#get-apiusersmemager-tasks)
8. [Implementation Details](#implementation-details)
9. [Examples](#examples)

## Overview

The Phase Planning system in ExpertDB manages the workflow for reviewing qualifications and institutional listings. It implements a sophisticated role-based access control system with contextual elevations, allowing regular users to be temporarily granted planner or manager privileges for specific applications within phases.

**Key Concepts:**
- **Phase Plan**: A planning period containing multiple applications requiring expert assignments
- **Application**: A specific task within a phase (QP or IL) requiring expert review
- **Contextual Elevation**: Temporary privilege assignment for users on specific applications
- **Expert Assignment**: Process of proposing and approving experts for applications

## Phase Planning Workflow

The phase planning workflow follows these steps:

1. **Phase Creation** (Admin)
   - Admin creates a phase with title and initial planner assignment
   - Applications are added with institution and qualification details
   - Phase ID is auto-generated (e.g., "PH-2025-001")

2. **User Assignment** (Admin)
   - Admin assigns users as planners for specific applications
   - Admin can assign managers who will provide expert ratings
   - Assignments are contextual to specific applications

3. **Expert Proposal** (Planner)
   - Users with planner elevation propose experts for their assigned applications
   - Can assign up to 2 experts per application
   - Proposals update application status to "assigned"

4. **Application Review** (Admin)
   - Admin reviews proposed experts
   - Can approve or reject with notes
   - Approval updates application status accordingly

5. **Expert Rating** (Manager)
   - When requested by admin, managers rate expert performance
   - Ratings are scoped to specific applications
   - Future enhancement will store ratings in database

## Role-Based Access Control

The system implements three-tier access control with contextual elevations:

### Base Roles

1. **super_user**
   - Complete system access
   - Can create admin users
   - Inherent access to all phase operations

2. **admin**
   - Full system access
   - Can create regular users
   - Manage all phases and applications
   - Inherent access to all phase operations

3. **user**
   - Limited base access
   - Can view phases and applications
   - Requires elevation for phase operations

### Contextual Elevations

1. **Planner Elevation**
   - Granted for specific applications within phases
   - Allows proposing experts for assigned applications
   - Managed via `/api/users/{id}/planner-assignments`
   - Validated by `RequirePlannerForApplication` middleware

2. **Manager Elevation**
   - Granted for specific applications within phases
   - Allows rating experts when requested by admin
   - Managed via `/api/users/{id}/manager-assignments`
   - Validated by `RequireManagerForApplication` middleware

### Access Control Implementation

```go
// Middleware checks for contextual access
func RequirePlannerForApplication(store storage.Storage) gin.HandlerFunc {
    // Admin/super_user bypass elevation checks
    // Regular users must have planner assignment for the specific application
}

func RequireManagerForApplication(store storage.Storage) gin.HandlerFunc {
    // Admin/super_user bypass elevation checks
    // Regular users must have manager assignment for the specific application
}
```

## Application Types

ExpertDB supports two types of applications within phases:

### QP (Qualification Placement)
- Review of specific qualifications offered by institutions
- Requires expert validation of curriculum and standards
- Typically involves academic program assessment

### IL (Institutional Listing)
- Review of entire institutions for listing status
- Comprehensive institutional assessment
- Covers governance, resources, and quality systems

## API Endpoints

### POST /api/phases

Creates a new phase plan with applications.

**Authorization**: Admin only

**Request Body**:
```json
{
  "title": "Q2 2025 Reviews",
  "assignedPlannerId": 123,
  "status": "draft",
  "applications": [
    {
      "type": "QP",
      "institutionName": "Tech University",
      "qualificationName": "BSc Computer Science",
      "expert1": null,
      "expert2": null,
      "status": "pending"
    }
  ]
}
```

**Response** (201 Created):
```json
{
  "success": true,
  "message": "Phase created successfully",
  "data": {
    "id": 1,
    "phaseId": "PH-2025-001"
  }
}
```

**Validation Rules**:
- `title` is required
- `assignedPlannerId` must be a valid user ID
- Application `type` must be "QP" or "IL"
- `institutionName` and `qualificationName` are required for applications

### GET /api/phases

Lists all phases with filtering and pagination.

**Authorization**: All authenticated users

**Query Parameters**:
- `limit`: Results per page (default: 50)
- `offset`: Number to skip (default: 0)
- `status`: Filter by phase status
- `planner_id`: Filter by assigned planner

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "phases": [
      {
        "id": 1,
        "phaseId": "PH-2025-001",
        "title": "Q2 2025 Reviews",
        "assignedPlannerId": 123,
        "plannerName": "John Doe",
        "status": "in_progress",
        "createdAt": "2025-07-22T10:00:00Z",
        "updatedAt": "2025-07-22T15:30:00Z",
        "applications": [
          {
            "id": 1,
            "phaseId": 1,
            "type": "QP",
            "institutionName": "Tech University",
            "qualificationName": "BSc Computer Science",
            "expert1": 456,
            "expert1Name": "Dr. Smith",
            "expert2": 789,
            "expert2Name": "Prof. Johnson",
            "status": "assigned",
            "rejectionNotes": null,
            "createdAt": "2025-07-22T10:00:00Z",
            "updatedAt": "2025-07-22T14:00:00Z"
          }
        ]
      }
    ],
    "pagination": {
      "totalCount": 25,
      "totalPages": 3,
      "currentPage": 1,
      "pageSize": 10,
      "hasNextPage": true,
      "hasPrevPage": false
    },
    "filters": {
      "status": "in_progress",
      "plannerId": null
    }
  }
}
```

### GET /api/phases/{id}

Retrieves details of a specific phase.

**Authorization**: All authenticated users

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "phaseId": "PH-2025-001",
    "title": "Q2 2025 Reviews",
    "assignedPlannerId": 123,
    "plannerName": "John Doe",
    "status": "in_progress",
    "createdAt": "2025-07-22T10:00:00Z",
    "updatedAt": "2025-07-22T15:30:00Z",
    "applications": [
      {
        "id": 1,
        "phaseId": 1,
        "type": "QP",
        "institutionName": "Tech University",
        "qualificationName": "BSc Computer Science",
        "expert1": 456,
        "expert1Name": "Dr. Smith",
        "expert2": 789,
        "expert2Name": "Prof. Johnson",
        "status": "assigned",
        "rejectionNotes": null,
        "createdAt": "2025-07-22T10:00:00Z",
        "updatedAt": "2025-07-22T14:00:00Z"
      }
    ]
  }
}
```

### PUT /api/phases/{id}

Updates phase details (implementation pending).

**Authorization**: Admin only

**Request Body**:
```json
{
  "title": "Updated Phase Title",
  "status": "completed"
}
```

### PUT /api/phases/{id}/applications/{app_id}

Proposes experts for an application.

**Authorization**: Planner for the specific application or Admin

**Request Body**:
```json
{
  "expert1": 456,
  "expert2": 789
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Experts proposed successfully"
}
```

**Access Control**:
- Uses `RequirePlannerForApplication` middleware
- Admin/super_user have inherent access
- Regular users need planner assignment for the specific application

### PUT /api/phases/{id}/applications/{app_id}/review

Reviews (approves/rejects) an application.

**Authorization**: Admin only

**Request Body**:
```json
{
  "action": "approve",
  "rejectionNotes": ""
}
```

Or for rejection:
```json
{
  "action": "reject",
  "rejectionNotes": "Experts do not meet qualification requirements"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Application approved successfully",
  "data": {
    "id": 1,
    "phaseId": 1,
    "status": "approved",
    "rejectionNotes": null,
    "updatedAt": "2025-07-22T16:00:00Z"
  }
}
```

**Validation**:
- `action` must be "approve" or "reject"
- `rejectionNotes` required when rejecting
- Application must be in "assigned" status
- At least one expert must be assigned before approving

### POST /api/phases/{id}/applications/{app_id}/ratings

Records expert ratings from managers.

**Authorization**: Manager for the specific application or Admin

**Request Body**:
```json
{
  "expertId": 456,
  "rating": 4,
  "comments": "Excellent performance in the review process"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "message": "Expert rating recorded successfully",
  "data": {
    "expertId": 456,
    "rating": 4,
    "appId": 1,
    "message": "Rating will be implemented with application_ratings table"
  }
}
```

**Validation**:
- `expertId` must be assigned to the application
- `rating` must be between 1 and 5
- Expert must exist in the database

**Note**: Full implementation pending - ratings will be stored in `application_ratings` table.

### GET /api/applications

Lists applications with comprehensive filtering.

**Authorization**: All authenticated users

**Query Parameters**:
- `phase_id`: Filter by phase
- `status`: Filter by status ("pending", "assigned", "approved", "rejected")
- `type`: Filter by type ("QP", "IL")
- `expert_id`: Filter by assigned expert
- `limit`: Results per page (default: 100)
- `offset`: Number to skip (default: 0)

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "applications": [
      {
        "id": 1,
        "phaseId": 1,
        "type": "QP",
        "institutionName": "Tech University",
        "qualificationName": "BSc Computer Science",
        "expert1": 456,
        "expert1Name": "Dr. Smith",
        "expert2": 789,
        "expert2Name": "Prof. Johnson",
        "status": "assigned",
        "rejectionNotes": null,
        "createdAt": "2025-07-22T10:00:00Z",
        "updatedAt": "2025-07-22T14:00:00Z"
      }
    ],
    "pagination": {
      "limit": 100,
      "offset": 0,
      "count": 45,
      "total": 150
    },
    "filters": {
      "phase_id": 1,
      "status": "assigned",
      "type": null,
      "expert_id": null
    }
  }
}
```

## Role Assignment Endpoints

### POST /api/users/{id}/planner-assignments

Assigns a user as planner to multiple applications.

**Authorization**: Admin only

**Request Body**:
```json
{
  "application_ids": [1, 2, 3]
}
```

**Response** (200 OK):
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

**Notes**:
- Replaces existing planner assignments
- Uses batch operations within a transaction
- Application IDs must be valid

### POST /api/users/{id}/manager-assignments

Assigns a user as manager to multiple applications.

**Authorization**: Admin only

**Request Body**:
```json
{
  "application_ids": [4, 5, 6]
}
```

**Response** (200 OK):
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

### DELETE /api/users/{id}/planner-assignments

Removes specific planner assignments.

**Authorization**: Admin only

**Request Body**:
```json
{
  "application_ids": [1, 2]
}
```

**Response** (200 OK):
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

Removes specific manager assignments.

**Authorization**: Admin only

**Request Body**:
```json
{
  "application_ids": [4, 5]
}
```

**Response** (200 OK):
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

Retrieves all role assignments for a user.

**Authorization**: Admin only

**Response** (200 OK):
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

## Manager Task Endpoints

### GET /api/users/me/manager-tasks

Retrieves applications where the current user has manager privileges.

**Authorization**: Any authenticated user (sees only their own assignments)

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "applicationId": 4,
        "phaseId": 2,
        "phaseTitle": "Q3 2025 Reviews",
        "institutionName": "Business College",
        "qualificationName": "MBA Program",
        "expert1Id": 234,
        "expert1Name": "Dr. Williams",
        "expert2Id": 567,
        "expert2Name": "Prof. Davis",
        "status": "approved",
        "ratingRequested": true
      }
    ],
    "count": 1
  }
}
```

## Implementation Details

### Database Schema

**phases table**:
```sql
CREATE TABLE phases (
    id INTEGER PRIMARY KEY,
    phase_id TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    assigned_planner_id INTEGER,
    status TEXT DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (assigned_planner_id) REFERENCES users(id)
);
```

**phase_applications table**:
```sql
CREATE TABLE phase_applications (
    id INTEGER PRIMARY KEY,
    phase_id INTEGER NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('QP', 'IL')),
    institution_name TEXT NOT NULL,
    qualification_name TEXT NOT NULL,
    expert_1 INTEGER,
    expert_2 INTEGER,
    status TEXT DEFAULT 'pending',
    rejection_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (phase_id) REFERENCES phases(id),
    FOREIGN KEY (expert_1) REFERENCES experts(id),
    FOREIGN KEY (expert_2) REFERENCES experts(id)
);
```

**application_planners table**:
```sql
CREATE TABLE application_planners (
    application_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (application_id, user_id),
    FOREIGN KEY (application_id) REFERENCES phase_applications(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

**application_managers table**:
```sql
CREATE TABLE application_managers (
    application_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (application_id, user_id),
    FOREIGN KEY (application_id) REFERENCES phase_applications(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Middleware Implementation

The system uses custom middleware for contextual access control:

```go
// RequirePlannerForApplication checks if user has planner access
func RequirePlannerForApplication(store storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(*domain.User)
        
        // Admin/super_user bypass
        if user.Role == "admin" || user.Role == "super_user" {
            c.Next()
            return
        }
        
        // Check contextual elevation
        appID, _ := strconv.Atoi(c.Param("app_id"))
        hasAccess, _ := store.IsUserPlannerForApplication(c, user.ID, appID)
        
        if !hasAccess {
            c.JSON(403, gin.H{"error": "Planner access required"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## Examples

### Complete Phase Planning Workflow Example

1. **Create a Phase** (Admin):
```bash
curl -X POST http://localhost:8080/api/phases \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Q2 2025 Review Cycle",
    "assignedPlannerId": 123,
    "applications": [
      {
        "type": "QP",
        "institutionName": "Tech University",
        "qualificationName": "BSc Computer Science"
      },
      {
        "type": "IL",
        "institutionName": "Business College",
        "qualificationName": "Institutional Review"
      }
    ]
  }'
```

2. **Assign Planners** (Admin):
```bash
curl -X POST http://localhost:8080/api/users/123/planner-assignments \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "application_ids": [1, 2]
  }'
```

3. **Propose Experts** (Planner):
```bash
curl -X PUT http://localhost:8080/api/phases/1/applications/1 \
  -H "Authorization: Bearer $PLANNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "expert1": 456,
    "expert2": 789
  }'
```

4. **Review Application** (Admin):
```bash
curl -X PUT http://localhost:8080/api/phases/1/applications/1/review \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve"
  }'
```

5. **Assign Manager** (Admin):
```bash
curl -X POST http://localhost:8080/api/users/234/manager-assignments \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "application_ids": [1]
  }'
```

6. **Rate Expert** (Manager):
```bash
curl -X POST http://localhost:8080/api/phases/1/applications/1/ratings \
  -H "Authorization: Bearer $MANAGER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "expertId": 456,
    "rating": 5,
    "comments": "Excellent review work"
  }'
```

### Filtering Applications Example

```bash
# Get all QP applications in assigned status
curl -X GET "http://localhost:8080/api/applications?type=QP&status=assigned" \
  -H "Authorization: Bearer $TOKEN"

# Get applications for a specific expert
curl -X GET "http://localhost:8080/api/applications?expert_id=456" \
  -H "Authorization: Bearer $TOKEN"

# Get applications for a specific phase
curl -X GET "http://localhost:8080/api/applications?phase_id=1" \
  -H "Authorization: Bearer $TOKEN"
```

### Managing Role Assignments Example

```bash
# View user's current assignments
curl -X GET http://localhost:8080/api/users/123/assignments \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Remove specific planner assignments
curl -X DELETE http://localhost:8080/api/users/123/planner-assignments \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "application_ids": [1]
  }'

# View manager tasks for current user
curl -X GET http://localhost:8080/api/users/me/manager-tasks \
  -H "Authorization: Bearer $USER_TOKEN"
```

## Error Responses

All endpoints follow a consistent error response format:

```json
{
  "success": false,
  "message": "Error description",
  "errors": ["Detailed error 1", "Detailed error 2"]
}
```

Common HTTP status codes:
- `200 OK`: Successful operation
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Notes

1. **Phase ID Generation**: Phase IDs are auto-generated in the format "PH-YYYY-NNN" (e.g., "PH-2025-001")

2. **Application Status Flow**:
   - `pending` → Initial status
   - `assigned` → Experts proposed
   - `approved` → Admin approved
   - `rejected` → Admin rejected with notes

3. **Contextual Access**: Regular users only have access to applications they're assigned to as planners or managers

4. **Future Enhancements**:
   - Full implementation of rating storage in `application_ratings` table
   - Phase update endpoint implementation
   - Bulk operations for phase management
   - Enhanced reporting for phase completion statistics

5. **Performance Considerations**:
   - All list endpoints support pagination
   - Filtering is performed at the database level
   - Proper indexes on foreign keys for optimal query performance