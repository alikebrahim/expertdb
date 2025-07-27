# ExpertDB Request Management API Documentation

**Date**: July 22, 2025  
**Version**: 1.0  
**Component**: Request Management Subsystem

## Overview

The ExpertDB Request Management system provides two primary workflows for managing expert data:

1. **Expert Request Management**: Allows users to submit requests for creating new expert profiles, which admins review and approve/reject.
2. **Expert Edit Request Management** (Planned): Enables users to propose changes to existing expert profiles with admin oversight.

Both systems follow a similar workflow pattern:
- Users submit requests with supporting documentation
- Admins review and approve/reject with appropriate documentation
- Approved requests trigger automatic updates to the expert database

## Table of Contents

1. [Request Management Overview](#request-management-overview)
   - [Expert Request Workflow](#expert-request-workflow)
   - [Expert Edit Request Workflow](#expert-edit-request-workflow)
2. [Expert Request Management Endpoints](#expert-request-management-endpoints)
   - [POST /api/expert-requests](#post-apiexpert-requests)
   - [GET /api/expert-requests](#get-apiexpert-requests)
   - [GET /api/expert-requests/{id}](#get-apiexpert-requestsid)
   - [PUT /api/expert-requests/{id}](#put-apiexpert-requestsid)
   - [PUT /api/expert-requests/{id}/edit](#put-apiexpert-requestsidedit)
   - [POST /api/expert-requests/batch-approve](#post-apiexpert-requestsbatch-approve)
3. [Expert Profile Edit System](#expert-profile-edit-system)
   - [POST /api/experts/{id}/edit](#post-apiexpertsidedit)
   - [GET /api/expert-edits](#get-apiexpert-edits)
   - [GET /api/expert-edits/{id}](#get-apiexpert-editsid)
   - [PUT /api/expert-edits/{id}/status](#put-apiexpert-editsidstatus)
4. [Business Rules and Validation](#business-rules-and-validation)
5. [Status Transitions](#status-transitions)
6. [Implementation Notes](#implementation-notes)

## Request Management Overview

### Expert Request Workflow

The expert request workflow handles the creation of new expert profiles:

```
User Submits Request → Admin Reviews → Approves with Document → Expert Profile Created
                                   ↘ Rejects → User Can Edit & Resubmit
```

**Key Features:**
- Form-based submission with CV upload
- Structured professional background (experience and education entries)
- Specialized areas selection with suggestion capability
- Batch approval support for multiple requests
- Automatic expert ID generation upon approval

### Expert Edit Request Workflow (Planned)

The expert edit request workflow manages changes to existing expert profiles:

```
User Proposes Changes → Admin Reviews → Approves → Changes Applied Automatically
                                     ↘ Rejects → Request Closed
```

**Key Features:**
- Track changes for all expert fields
- File upload support for new CV/documents
- Change summary and reasoning
- Audit trail of all modifications
- Automatic application of approved changes

## Expert Request Management Endpoints

### POST /api/expert-requests

**Purpose**: Submits an expert request with CV upload and structured professional background.

**Method**: POST  
**Path**: `/api/expert-requests`  
**Access Control**: Any authenticated user  
**Content-Type**: `multipart/form-data`

#### Request Payload

```text
name: string                       // Required: Expert's full name (min 2 chars)
designation: string                // Required: Professional title
                                  // Options: "Prof.", "Dr.", "Mr.", "Ms.", "Mrs.", "Miss", "Eng."
affiliation: string                // Required: Organization/institution (min 2 chars)
phone: string                      // Required: Contact phone (min 8 chars, format validated)
email: string                      // Required: Contact email (email format validated)
isBahraini: boolean               // Required: Bahraini citizenship status
isAvailable: boolean              // Required: Current availability for assignments
role: string                      // Required: Expert role
                                  // Options: "evaluator", "validator", "evaluator/validator"
employmentType: string            // Required: Employment type
                                  // Options: "academic", "employer"
generalArea: int                  // Required: ID from expert_areas table
specializedAreaIds: string        // Optional: JSON array of existing area IDs
                                  // Example: "[1,4,6]"
suggestedSpecializedAreas: string // Optional: JSON array of new area suggestions
                                  // Example: ["Machine Learning", "Blockchain"]
isTrained: boolean                // Required: BQA training completion status
isPublished: boolean              // Optional: Publication status (defaults to false)
experienceEntries: string         // Optional: JSON array of experience entries
educationEntries: string          // Optional: JSON array of education entries
cv: file                          // Required: CV document (PDF format, max 5MB)
```

#### Professional Background Structures

**Experience Entries Format:**
```json
[
  {
    "organization": "Tech Company",
    "position": "Senior Software Engineer", 
    "startDate": "2020-01",
    "endDate": "Present",
    "isCurrent": true,
    "country": "United States",
    "description": "Led development team and managed project delivery"
  }
]
```

**Education Entries Format:**
```json
[
  {
    "institution": "University Name",
    "degree": "Bachelor of Computer Science",
    "fieldOfStudy": "Computer Science", 
    "graduationYear": "2020",
    "country": "United States",
    "description": "Focused on software engineering and algorithms"
  }
]
```

#### Response Payloads

**Success (201 Created):**
```json
{
  "success": true,
  "message": "Expert request created successfully",
  "data": {
    "id": 26
  }
}
```

**Error (400 Bad Request):**
```json
{
  "errors": [
    "name is required",
    "cv file missing",
    "invalid email format"
  ]
}
```

#### Specialized Areas Workflow

Users can work with specialized areas in two ways:

1. **Select Existing Areas**: Use `specializedAreaIds` with IDs from `/api/specialized-areas`
2. **Suggest New Areas**: Use `suggestedSpecializedAreas` for areas that don't exist yet

**Example Combined Usage:**
```javascript
formData.append('specializedAreaIds', '[1, 4, 6]');
formData.append('suggestedSpecializedAreas', '["Machine Learning", "Blockchain Technology"]');
```

**Admin Notification**: Admins are notified of suggested areas and can create them before approving the request.

### GET /api/expert-requests

**Purpose**: Retrieves paginated expert requests with status filtering.

**Method**: GET  
**Path**: `/api/expert-requests`  
**Access Control**: 
- **Admin/Super User**: Can view all expert requests
- **Regular User**: Can view only their own submitted expert requests

#### Query Parameters

- `limit`: Number of results per page (default: 20)
- `offset`: Number of results to skip (default: 0)
- `status`: Filter by status - `pending`, `approved`, `rejected`

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "data": {
    "requests": [
      {
        "id": 26,
        "name": "Dr. John Smith",
        "status": "pending",
        "cvDocumentId": 345,
        "approvalDocumentId": null,
        "designation": "Dr.",
        "affiliation": "University of Bahrain",
        "isBahraini": true,
        "isAvailable": true,
        "role": "evaluator",
        "employmentType": "academic",
        "generalArea": 3,
        "specializedArea": "1,4,6",
        "suggestedSpecializedAreas": ["Machine Learning", "AI Ethics"],
        "isTrained": true,
        "phone": "+973 17123456",
        "email": "john.smith@uob.edu.bh",
        "experienceEntries": [
          {
            "id": 1,
            "organization": "University of Bahrain",
            "position": "Associate Professor",
            "startDate": "2018-09",
            "endDate": null,
            "isCurrent": true,
            "country": "Bahrain",
            "description": "Teaching computer science courses",
            "createdAt": "2025-07-22T10:00:00Z",
            "updatedAt": "2025-07-22T10:00:00Z"
          }
        ],
        "educationEntries": [
          {
            "id": 1,
            "institution": "MIT",
            "degree": "PhD in Computer Science",
            "fieldOfStudy": "Computer Science",
            "graduationYear": "2015",
            "country": "United States",
            "description": "Research in machine learning",
            "createdAt": "2025-07-22T10:00:00Z",
            "updatedAt": "2025-07-22T10:00:00Z"
          }
        ],
        "isPublished": false,
        "createdAt": "2025-07-22T10:00:00Z",
        "updatedAt": "2025-07-22T10:00:00Z"
      }
    ],
    "pagination": {
      "limit": 20,
      "offset": 0,
      "count": 5
    }
  }
}
```

### GET /api/expert-requests/{id}

**Purpose**: Retrieves a specific expert request with full details.

**Method**: GET  
**Path**: `/api/expert-requests/{id}`  
**Access Control**: Admin only

#### Response Payload

Returns a single expert request object with the same structure as the list endpoint, including full professional background details.

### PUT /api/expert-requests/{id}

**Purpose**: Approves or rejects an expert request with approval document.

**Method**: PUT  
**Path**: `/api/expert-requests/{id}`  
**Access Control**: Admin only  
**Content-Type**: `multipart/form-data`

#### Request Payload

```text
status: string                    // Required: "approved" or "rejected"
rejectionReason: string           // Required if status is "rejected"
approvalDocument: file            // Required if status is "approved"
```

#### Response Payloads

**Success (200 OK) - Approved:**
```json
{
  "success": true,
  "message": "Expert request approved and expert profile created",
  "data": {
    "expertId": 440,
    "expertBusinessId": "EXP-0440"
  }
}
```

**Success (200 OK) - Rejected:**
```json
{
  "success": true,
  "message": "Expert request rejected"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "Approval document required for approval"
}
```

#### Approval Process

When a request is approved:
1. Approval document is stored in the document system
2. Expert profile is created with all submitted data
3. Professional background entries are created
4. Unique expert ID is generated (e.g., "EXP-0440")
5. Request status is updated to "approved"

### PUT /api/expert-requests/{id}/edit

**Purpose**: Edits an expert request before approval.

**Method**: PUT  
**Path**: `/api/expert-requests/{id}/edit`  
**Access Control**: 
- Admin: Can edit any pending request
- User: Can edit their own rejected requests only  
**Content-Type**: `multipart/form-data`

#### Request Payload

Same structure as POST /api/expert-requests, with all fields optional. Only provided fields will be updated.

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "message": "Expert request updated successfully"
}
```

**Error (403 Forbidden):**
```json
{
  "error": "Only admins or request owner can edit"
}
```

### POST /api/expert-requests/batch-approve

**Purpose**: Approves multiple expert requests with one approval document.

**Method**: POST  
**Path**: `/api/expert-requests/batch-approve`  
**Access Control**: Admin only  
**Content-Type**: `multipart/form-data`

#### Request Payload

```text
data: string                      // JSON array of request IDs
                                 // Example: "[26, 27, 28]"
approvalDocument: file           // Required: Single approval document for all
```

#### Response Payload

**Success (200 OK):**
```json
{
  "success": true,
  "message": "Approved 3 of 4 requests",
  "data": {
    "totalRequests": 4,
    "approvedCount": 3,
    "approvedIds": [26, 27, 28],
    "errors": {
      "29": "Request already approved"
    },
    "errorCount": 1
  }
}
```

**Error (400 Bad Request):**
```json
{
  "error": "Missing approval document"
}
```

#### Batch Processing Rules

- All requests must be in "pending" status
- Same approval document is attached to all approved requests
- Transaction ensures all-or-nothing for database operations
- Individual errors don't affect other approvals
- Each approved request creates a separate expert profile

## Expert Profile Edit System

### POST /api/experts/{id}/edit

**Purpose**: Creates an expert profile edit request for an existing expert.

**Method**: POST  
**Path**: `/api/experts/{id}/edit`  
**Access Control**: Any authenticated user  
**Content-Type**: `multipart/form-data`

#### Request Payload

```text
data: string                      // JSON string with proposed changes
cv: file                          // Optional: New CV file
approval_document: file           // Optional: New approval document
```

**JSON Data Structure:**
```json
{
  "proposedName": "Dr. John Smith Jr.",
  "proposedDesignation": "Prof.",
  "proposedEmail": "john.smith.jr@uob.edu.bh",
  "proposedGeneralArea": 5,
  "proposedSpecializedArea": "1,4,6,12",
  "experienceChanges": {
    "add": [...],
    "update": [...],
    "delete": [...]
  },
  "educationChanges": {
    "add": [...],
    "update": [...],
    "delete": [...]
  },
  "removeCv": false,
  "removeApprovalDocument": false,
  "changeSummary": "Updated name, email, and added new experience entry",
  "changeReason": "Recent promotion and department change"
}
```

### GET /api/expert-edits

**Purpose**: Lists expert profile edit proposals with filtering.

**Method**: GET  
**Path**: `/api/expert-edits`  
**Access Control**: 
- Regular users: See only their own requests
- Admin/Super_user: See all requests

#### Query Parameters

- `status`: Filter by status (`pending`, `approved`, `rejected`, `cancelled`)
- `expert_id`: Filter by expert ID
- `limit`, `offset`: Pagination

### GET /api/expert-edits/{id}

**Purpose**: Retrieves a specific expert profile edit proposal with full change details.

**Method**: GET  
**Path**: `/api/expert-edits/{id}`  
**Access Control**: Requester can view own requests, admin/super_user can view all

### PUT /api/expert-edits/{id}/status

**Purpose**: Updates the status of an expert profile edit proposal (auto-applies if approved).

**Method**: PUT  
**Path**: `/api/expert-edits/{id}/status`  
**Access Control**: 
- Users: Can cancel their own pending requests
- Admin/Super_user: Can approve/reject any request

#### Request Payload

```json
{
  "status": "approved",
  "rejectionReason": null
}
```

#### Auto-Application Behavior

When an admin approves an edit proposal (status = "approved"), the system automatically applies the changes to the expert profile. No separate apply step is required.

**Success Response (with auto-apply):**
```json
{
  "success": true,
  "message": "Edit proposal approved and changes applied successfully",
  "data": {
    "expertId": 440,
    "changesApplied": 7
  }
}
```

## Business Rules and Validation

### Expert Request Validation

1. **Required Fields**:
   - All fields marked as required must be provided
   - CV file must be PDF format and under 5MB
   - Email must be valid format
   - Phone must be at least 8 characters

2. **Designation Options**:
   - Must be one of: "Prof.", "Dr.", "Mr.", "Ms.", "Mrs.", "Miss", "Eng."

3. **Role Options**:
   - Must be one of: "evaluator", "validator", "evaluator/validator"

4. **Employment Type Options**:
   - Must be one of: "academic", "employer"

5. **General Area**:
   - Must be valid ID from expert_areas table


### Edit Request Validation

1. **Change Tracking**:
   - At least one field must be changed
   - changeSummary and changeReason are required

2. **File Operations**:
   - Cannot remove and upload same file type in single request
   - File uploads follow same rules as expert requests

3. **Status Constraints**:
   - Can only edit pending requests
   - Approved requests require new edit request

## Status Transitions

### Expert Request Status Flow

```
pending → approved (with approval document)
        ↘ rejected (with reason)
          ↗ pending (after edit by user/admin)
```

### Expert Edit Request Status Flow

```
pending → approved → applied (automatic)
        ↘ rejected (with reason)
        ↘ cancelled (by requester)
```

## Implementation Notes

### File Storage

- CV files stored in: `/storage/documents/cv/`
- Approval documents stored in: `/storage/documents/approval/`
- File naming: `{type}_{requestId}_{timestamp}.pdf`

### Database Transactions

- Approval process uses transactions to ensure consistency
- Batch approvals use single transaction for all operations
- Edit request application uses transaction for all field updates

### Logging

All operations are logged with:
- User performing action
- Timestamp
- Request ID
- Action type
- Result status

Example log entries:
```
Expert request created: ID 26 by user@example.com
Request approved: ID 26 by admin@expertdb.com
Batch approved 3 requests by admin@expertdb.com
Edit request created for expert 440 by user@example.com
```

### Performance Considerations

- Batch approval limited to 100 requests per operation
- File uploads processed asynchronously after validation
- Professional background entries use bulk insert operations
- Edit request changes tracked using JSON diff algorithm

### Security

- File uploads scanned for malicious content
- SQL injection prevention through parameterized queries
- XSS prevention through input sanitization
- CSRF protection via token validation
- Rate limiting on request submissions (5 per hour per user)

---

This documentation provides a comprehensive reference for the ExpertDB request management subsystem. For general API information, refer to the main API_REFERENCE.md document.