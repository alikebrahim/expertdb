# ExpertDB Document Management API Reference

**Date**: July 22, 2025  
**Version**: 1.0  
**Purpose**: Comprehensive reference for document management endpoints in the ExpertDB system, providing detailed documentation for file upload, retrieval, and management operations.

## Table of Contents

1. [Overview](#overview)
2. [Data Model](#data-model)
3. [Document Types](#document-types)
4. [File Storage](#file-storage)
5. [API Endpoints](#api-endpoints)
   - [POST /api/documents](#post-apidocuments)
   - [GET /api/experts/{id}/documents](#get-apiexpertsiddocuments)
   - [GET /api/documents/{id}](#get-apidocumentsid)
   - [GET /api/documents/{id}/download](#get-apidocumentsiddownload)
   - [DELETE /api/documents/{id}](#delete-apidocumentsid)
6. [Request/Response Examples](#requestresponse-examples)
7. [Security Considerations](#security-considerations)
8. [Implementation Details](#implementation-details)
9. [Error Handling](#error-handling)
10. [Integration with Other Features](#integration-with-other-features)

## Overview

The document management system in ExpertDB provides secure file upload, storage, and retrieval capabilities for expert-related documents. It supports two document types: CVs and approval documents. The system uses multipart form-data for file uploads and stores files locally with metadata tracked in the SQLite database.

### Key Features
- Multipart form-data file upload support
- Two document types (CV and approval)
- Secure file storage with unique naming
- Document metadata tracking
- Integration with expert profiles and requests
- Cascading deletion with expert profiles
- Role-based access control

## Data Model

### Document Structure

```typescript
interface Document {
  id: number;
  expertId: number;
  documentType: string;      // "cv" or "approval"
  filename: string;          // Original filename
  filePath: string;          // Storage path
  contentType: string;       // MIME type
  fileSize: number;          // Size in bytes
  uploadDate: string;        // ISO 8601 timestamp
}
```

### Database Schema

```sql
CREATE TABLE IF NOT EXISTS expert_documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    document_type TEXT NOT NULL,
    filename TEXT NOT NULL,
    file_path TEXT NOT NULL,
    content_type TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE
);
```

## Document Types

The system supports the following document types:

| Type | Description | Common Use Cases |
|------|-------------|------------------|
| `cv` | Curriculum Vitae | Expert's professional resume |
| `approval` | Approval Documents | Administrative approval records |

## File Storage

### Storage Structure
```
uploads/
├── documents/
│   ├── {expertId}/
│   │   ├── cv/
│   │   │   └── {timestamp}_{originalFilename}
│   │   └── approval/
│   │       └── {timestamp}_{originalFilename}
```

### Naming Convention
- Files are stored with timestamp prefix to ensure uniqueness
- Original filename is preserved for user reference
- Directory structure organized by expert ID and document type (cv or approval)

## API Endpoints

### POST /api/documents

**Purpose**: Uploads a document for an expert.

**Method**: POST  
**Path**: `/api/documents`  
**Access Control**: Admin only

#### Request Headers
```http
Authorization: Bearer <admin_token>
Content-Type: multipart/form-data
```

#### Request Payload (Form-data)
```text
file: file                    // Required: The file to upload
documentType: string          // Required: One of "cv" or "approval"
expertId: int                 // Required: ID of the expert
```

#### Response Payload

**Success (201 Created)**:
```json
{
  "success": true,
  "message": "Document uploaded successfully",
  "data": {
    "id": 123,
    "expertId": 456,
    "documentType": "cv",
    "filename": "john_doe_cv.pdf",
    "filePath": "uploads/documents/456/cv/1737553401_john_doe_cv.pdf",
    "contentType": "application/pdf",
    "fileSize": 245678,
    "uploadDate": "2025-07-22T14:30:01Z"
  }
}
```

**Error Responses**:

400 Bad Request:
```json
{
  "error": "Invalid document type"
}
```

400 Bad Request:
```json
{
  "error": "Expert ID is required"
}
```

404 Not Found:
```json
{
  "error": "Expert not found"
}
```

#### Implementation Notes
- File: `internal/api/handlers/documents/document_handler.go`
- Uses `internal/documents/service.go` for file processing
- Validates document type against allowed values
- Creates directory structure if not exists
- Stores file with timestamp prefix for uniqueness

### GET /api/experts/{id}/documents

**Purpose**: Lists all documents for a specific expert.

**Method**: GET  
**Path**: `/api/experts/{id}/documents`  
**Access Control**: All authenticated users

#### Request Headers
```http
Authorization: Bearer <token>
```

#### Path Parameters
- `id`: Expert ID (integer)

#### Response Payload

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "expertId": 456,
    "count": 2,
    "documents": [
      {
        "id": 123,
        "expertId": 456,
        "documentType": "cv",
        "filename": "john_doe_cv.pdf",
        "filePath": "uploads/documents/456/cv/1737553401_john_doe_cv.pdf",
        "contentType": "application/pdf",
        "fileSize": 245678,
        "uploadDate": "2025-07-22T14:30:01Z"
      },
      {
        "id": 124,
        "expertId": 456,
        "documentType": "approval",
        "filename": "approval_letter.pdf",
        "filePath": "uploads/documents/456/approval/1737553402_approval_letter.pdf",
        "contentType": "application/pdf",
        "fileSize": 123456,
        "uploadDate": "2025-07-22T14:31:02Z"
      }
    ]
  }
}
```

**Error Responses**:

404 Not Found:
```json
{
  "error": "Expert not found"
}
```

#### Implementation Notes
- File: `internal/api/handlers/documents/document_handler.go`
- Returns all documents associated with the expert
- Includes both CVs and approval documents
- Accessible to all authenticated users (Phase 6A enhancement)

### GET /api/documents/{id}

**Purpose**: Retrieves metadata for a specific document.

**Method**: GET  
**Path**: `/api/documents/{id}`  
**Access Control**: All authenticated users

#### Request Headers
```http
Authorization: Bearer <token>
```

#### Path Parameters
- `id`: Document ID (integer)

#### Response Payload

**Success (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": 123,
    "expertId": 456,
    "documentType": "cv",
    "filename": "john_doe_cv.pdf",
    "filePath": "uploads/documents/456/cv/1737553401_john_doe_cv.pdf",
    "contentType": "application/pdf",
    "fileSize": 245678,
    "uploadDate": "2025-07-22T14:30:01Z"
  }
}
```

**Error Responses**:

404 Not Found:
```json
{
  "error": "Document not found"
}
```

#### Implementation Notes
- File: `internal/api/handlers/documents/document_handler.go`
- Returns document metadata only (not file content)
- To download the actual file, use the `/api/documents/{id}/download` endpoint
- Accessible to all authenticated users (Phase 6A enhancement)

### GET /api/documents/{id}/download

**Purpose**: Downloads the actual file content for a specific document.

**Method**: GET  
**Path**: `/api/documents/{id}/download`  
**Access Control**: All authenticated users

#### Request Headers
```http
Authorization: Bearer <token>
```

#### Path Parameters
- `id`: Document ID (integer)

#### Response

**Success (200 OK)**:
- **Content-Type**: Original file MIME type (e.g., `application/pdf`)
- **Content-Length**: File size in bytes
- **Content-Disposition**: `attachment; filename="original_filename.pdf"`
- **Cache-Control**: `no-cache, no-store, must-revalidate`
- **Body**: Binary file content

**Error Responses**:

404 Not Found:
```json
{
  "error": "Document not found"
}
```

500 Internal Server Error:
```json
{
  "error": "Failed to open document file"
}
```

#### Implementation Notes
- File: `internal/api/handlers/documents/document_handler.go`
- Streams file content directly to client for efficient memory usage
- Sets appropriate headers for browser download behavior
- Maintains original filename in download
- Security: Only accessible to authenticated users
- Automatically handles file serving with proper MIME types
- Includes cache prevention headers for sensitive documents

#### Usage Examples

**Direct download in browser**:
```javascript
// Get document metadata first
const docResponse = await fetch('/api/documents/123', {
  headers: { 'Authorization': 'Bearer ' + token }
});
const doc = await docResponse.json();

// Download the file
const downloadUrl = '/api/documents/123/download';
window.open(downloadUrl + '?token=' + token, '_blank');
```

**Programmatic download**:
```javascript
const response = await fetch('/api/documents/123/download', {
  headers: { 'Authorization': 'Bearer ' + token }
});

if (response.ok) {
  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'document.pdf'; // Use actual filename from metadata
  a.click();
  window.URL.revokeObjectURL(url);
}
```

### DELETE /api/documents/{id}

**Purpose**: Deletes a document and its associated file.

**Method**: DELETE  
**Path**: `/api/documents/{id}`  
**Access Control**: Admin only

#### Request Headers
```http
Authorization: Bearer <admin_token>
```

#### Path Parameters
- `id`: Document ID (integer)

#### Response Payload

**Success (200 OK)**:
```json
{
  "success": true,
  "message": "Document deleted successfully"
}
```

**Error Responses**:

404 Not Found:
```json
{
  "error": "Document not found"
}
```

#### Implementation Notes
- File: `internal/api/handlers/documents/document_handler.go`
- Deletes both database record and physical file
- Cascades with expert deletion (Phase 6C)
- Admin access only for security

## Request/Response Examples

### Example 1: Upload CV for Expert

**Request**:
```bash
curl -X POST https://api.expertdb.com/api/documents \
  -H "Authorization: Bearer <admin_token>" \
  -F "file=@/path/to/cv.pdf" \
  -F "documentType=cv" \
  -F "expertId=456"
```

**Response**:
```json
{
  "success": true,
  "message": "Document uploaded successfully",
  "data": {
    "id": 123,
    "expertId": 456,
    "documentType": "cv",
    "filename": "cv.pdf",
    "filePath": "uploads/documents/456/cv/1737553401_cv.pdf",
    "contentType": "application/pdf",
    "fileSize": 245678,
    "uploadDate": "2025-07-22T14:30:01Z"
  }
}
```

### Example 2: List All Documents for an Expert

**Request**:
```bash
curl -X GET https://api.expertdb.com/api/experts/456/documents \
  -H "Authorization: Bearer <token>"
```

**Response**:
```json
{
  "success": true,
  "data": {
    "expertId": 456,
    "count": 2,
    "documents": [
      {
        "id": 123,
        "expertId": 456,
        "documentType": "cv",
        "filename": "cv.pdf",
        "filePath": "uploads/documents/456/cv/1737553401_cv.pdf",
        "contentType": "application/pdf",
        "fileSize": 245678,
        "uploadDate": "2025-07-22T14:30:01Z"
      },
      {
        "id": 124,
        "expertId": 456,
        "documentType": "approval",
        "filename": "approval.pdf",
        "filePath": "uploads/documents/456/approval/1737553402_approval.pdf",
        "contentType": "application/pdf",
        "fileSize": 123456,
        "uploadDate": "2025-07-22T14:31:02Z"
      }
    ]
  }
}
```

### Example 3: Download a Document

**Request**:
```bash
curl -X GET https://api.expertdb.com/api/documents/123/download \
  -H "Authorization: Bearer <token>" \
  --output cv.pdf
```

**Response**:
- HTTP Status: 200 OK
- Headers:
  - `Content-Type: application/pdf`
  - `Content-Length: 245678`
  - `Content-Disposition: attachment; filename="cv.pdf"`
- Body: Binary PDF content

**Browser Usage**:
```javascript
// Direct download link (requires authentication)
const downloadUrl = `${API_BASE_URL}/api/documents/123/download`;
fetch(downloadUrl, {
  headers: { 'Authorization': `Bearer ${token}` }
})
.then(response => response.blob())
.then(blob => {
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'cv.pdf';
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  window.URL.revokeObjectURL(url);
});
```

## Security Considerations

### Access Control
- **Upload**: Admin only - ensures only authorized personnel can add documents
- **View**: All authenticated users - allows broader access for viewing
- **Delete**: Admin only - prevents unauthorized document removal

### File Validation
- Document type must be from allowed list
- File size limits enforced by server configuration
- MIME type validation for uploaded files
- Expert ID validation ensures documents linked to existing experts

### Storage Security
- Files stored outside web root to prevent direct access
- Unique naming prevents filename collisions
- Directory structure provides organization and access control

## Implementation Details

### File Upload Service

The document service (`internal/documents/service.go`) handles:
- File validation and processing
- Directory creation and management
- Unique filename generation
- Error handling for upload failures

### Storage Layer

The SQLite storage implementation (`internal/storage/sqlite/document.go`) provides:
- CRUD operations for document metadata
- Transaction support for consistency
- Cascading deletion with expert profiles
- Query optimization for document listing

### Multipart Handling

The API uses Go's standard `multipart` package for handling file uploads:
- Automatic memory/disk storage based on file size
- Stream processing for large files
- Proper cleanup of temporary files

## Error Handling

### Common Error Scenarios

1. **Invalid Document Type**
   - Status: 400 Bad Request
   - Message: "Invalid document type"
   - Solution: Use one of the allowed types: cv, approval

2. **Missing Required Fields**
   - Status: 400 Bad Request
   - Message: "Expert ID is required" or "Document type is required"
   - Solution: Ensure all required fields are provided

3. **Expert Not Found**
   - Status: 404 Not Found
   - Message: "Expert not found"
   - Solution: Verify expert ID exists before uploading documents

4. **File Upload Failure**
   - Status: 500 Internal Server Error
   - Message: "Failed to save document"
   - Solution: Check server logs for disk space or permission issues

### Error Response Format

All errors follow the standard format:
```json
{
  "error": "Error message describing the issue"
}
```

## Integration with Other Features

### Expert Creation Workflow

Documents integrate with the expert creation process:
1. User submits expert request with CV (multipart upload)
2. CV stored via document service during request creation
3. Admin uploads approval document when approving request
4. Both documents linked to created expert profile

### Expert Profile Updates

The expert update endpoint (`PUT /api/experts/{id}`) supports:
- Multipart form-data for file uploads
- Updating CV with `cvFile` field
- Updating approval document with `approvalDocument` field
- JSON-only updates when no files involved

### Batch Approvals

Batch approval endpoint handles:
- Single approval document for multiple expert requests
- Document duplication for each approved expert
- Transaction consistency for all operations

### Expert Deletion

When an expert is deleted:
- All associated documents are automatically deleted (CASCADE)
- Physical files are removed from storage
- Database records are cleaned up

## Best Practices

### Frontend Integration

1. **File Upload Progress**
   - Implement progress tracking for large files
   - Show upload status to users
   - Handle network interruptions gracefully

2. **File Type Validation**
   - Validate file types before upload
   - Check file size limits client-side
   - Provide clear error messages

3. **Document Display**
   - Use document metadata to show file information
   - Implement proper download functionality
   - Handle missing documents gracefully

### Backend Considerations

1. **Storage Management**
   - Monitor disk usage regularly
   - Implement file cleanup for orphaned documents
   - Consider compression for large files

2. **Performance**
   - Use streaming for large file uploads
   - Implement caching for frequently accessed documents
   - Optimize database queries for document listings

3. **Security**
   - Regularly audit file permissions
   - Implement virus scanning for uploads
   - Log all document operations for audit trail

## Conclusion

The document management system provides a robust foundation for handling expert-related files in ExpertDB. With support for multiple document types, secure storage, and comprehensive API endpoints, it enables efficient document workflow management while maintaining security and data integrity. The integration with expert profiles and request workflows ensures seamless operation within the larger system context.