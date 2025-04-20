# ExpertDB Frontend Implementation Guide

## Current Implementation Status

Based on an assessment of the existing frontend codebase against the FRONTED_DEV.md specifications, this implementation guide focuses on integrating the frontend with the backend and completing any outstanding development tasks.

## 1. Directory Structure Overview

| Component | Status | Implementation Notes |
|-----------|--------|---------------------|
| **API Integration** | Partially Implemented | Single consolidated `api.ts` file instead of separate modules |
| **Components** | Partially Implemented | Core components exist, but folder structure is flatter than planned |
| **Context** | Partially Implemented | Only `AuthContext` implemented, missing `UIContext` |
| **Hooks** | Implemented | Basic authentication hooks available |
| **Pages** | Partially Implemented | Core pages exist without sub-folder organization |
| **Types** | Partially Implemented | Single `index.ts` file instead of separated type files |
| **Utils** | Missing | No utility functions implemented |

## 2. Backend-Frontend Integration Points

### 2.1 API Integration

The current implementation uses a centralized `api.ts` file that handles all API interactions with organized modules for different entity types:

```javascript
export default {
  auth: authApi,
  experts: expertsApi,
  expertRequests: expertRequestsApi,
  users: usersApi,
  statistics: statisticsApi,
  expertAreas: expertAreasApi,
  documents: documentApi,
  engagements: engagementApi,
  phases: phaseApi,
  backup: backupApi,
};
```

### 2.2 Authentication Flow

Authentication is implemented using JWT tokens:
- Login endpoint (`/api/auth/login`) returns a token
- Token stored in localStorage
- Request interceptor adds token to every API request
- Response interceptor handles 401 errors by logging out

### 2.3 Form Data Handling

- Form submissions use FormData for file uploads
- Content-Type headers automatically set for file uploads
- Error handling in place for form validations

### 2.4 Error Handling

- Centralized error handling in the API service
- Specific error messages based on HTTP status codes
- Debug mode for development with detailed logging

## 3. Implementation Plan

### 3.1 Phase 1: Refactoring Current Structure

1. **API Module Separation**
   - Split current `api.ts` into separate module files as outlined in the plan
   - Create a dedicated API client for reuse
   - Add token refresh mechanism

2. **Component Organization**
   - Reorganize components into proper folder structure:
     - Move form components to `/components/forms/`
     - Move table components to `/components/tables/`
     - Move chart components to `/components/charts/`
     - Move modal components to `/components/modals/`

3. **Add Utility Functions**
   - Create `/src/utils/` directory
   - Implement formatters for dates, currencies, etc.
   - Add validation utilities for forms
   - Create permission utilities for role-based access

### 3.2 Phase 2: Feature Implementation

1. **Expert Management**
   - Complete expert filtering functionality
   - Enhance expert detail view
   - Implement batch operations
   - Add area management

2. **Document Management**
   - Add document preview feature
   - Enhance document listing with filters
   - Implement document versioning

3. **Phase Planning**
   - Implement phase planning interface
   - Build expert assignment workflow
   - Create admin review process

4. **Statistics Enhancement**
   - Complete dashboard implementation
   - Add export functionality
   - Implement detailed reports

### 3.3 Phase 3: State Management Improvements

1. **Context Implementation**
   - Add `UIContext` for managing UI state
   - Implement sidebar collapsible functionality
   - Add notification system

2. **Form Handling**
   - Use React Hook Form consistently
   - Implement Zod for validation schemas
   - Create reusable form components

3. **API Improvements**
   - Add request retries for failed requests
   - Implement token refreshing
   - Add optimistic updates for better UX

## 4. Backend Integration Details

### 4.1 API Endpoints

| Endpoint Category | Base URL | Status | Implementation Notes |
|-------------------|----------|--------|---------------------|
| Authentication | `/api/auth` | Implemented | Login/logout functionality |
| Experts | `/experts` | Implemented | CRUD operations, filtering |
| Expert Requests | `/expert-requests` | Implemented | Creation, approval workflow |
| Users | `/users` | Implemented | User management |
| Statistics | `/statistics` | Implemented | Various statistics endpoints |
| Expert Areas | `/expert/areas` | Implemented | Area management |
| Documents | `/documents` | Implemented | Upload, download, listing |
| Engagements | `/expert-engagements` | Implemented | Listing, importing |
| Phases | `/phases` | Implemented | Creation, management, review |
| Backup | `/backup` | Implemented | System backup |

### 4.2 Data Flow

1. **Authentication Flow**
   ```
   Frontend Login Form → Backend Auth Endpoint → JWT Token → Local Storage → API Interceptors
   ```

2. **Expert Management Flow**
   ```
   Expert Form → FormData → Backend Endpoint → Database → Response → Update UI
   ```

3. **Document Upload Flow**
   ```
   File Input → FormData → Backend Upload Endpoint → File System → Response → Update UI
   ```

4. **Statistics Flow**
   ```
   Stats Page Load → Multiple API Requests → Backend Aggregation → Response → Chart Rendering
   ```

### 4.3 Role-Based Access Control

Access control is implemented through the following mechanisms:
- JWT tokens with role information
- Frontend route protection (`<ProtectedRoute>` component)
- UI conditional rendering based on user role
- Backend permission verification

## 5. Test and Validation Strategy

1. **Component Testing**
   - Implement unit tests for critical components
   - Focus on form validation and submission

2. **API Integration Testing**
   - Test all API endpoints with mock data
   - Verify error handling for all scenarios

3. **End-to-End User Workflows**
   - Test complete user journeys for key features
   - Validate role-based access control

## 6. Implementation Priority

1. **Critical Path (Immediate)**
   - Fix any authentication issues
   - Complete expert management features
   - Implement document management

2. **Secondary Features (Next Phase)**
   - Phase planning
   - Detailed statistics
   - Area management

3. **Polish & Improvements (Final Phase)**
   - UI/UX refinements
   - Performance optimizations
   - Advanced filtering and search

## 7. Best Practices for Backend-Frontend Integration

1. **Type Safety**
   - Ensure consistent types between frontend and backend
   - Define shared interfaces for API requests/responses

2. **Error Handling**
   - Consistent error structure from backend
   - Comprehensive error handling in frontend

3. **Loading States**
   - Implement skeleton loaders for better UX
   - Add progress indicators for file uploads

4. **Data Validation**
   - Client-side validation before submission
   - Graceful handling of server validation errors