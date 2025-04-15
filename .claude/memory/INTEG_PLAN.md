# API Integration Implementation Plan

This document outlines a phased approach to address the integration issues identified in the [API_INTEG_DESC.md](API_INTEG_DESC.md) assessment. Each phase has its own summary document tracking implementation details.

**Current Progress**:
- Phase 1 (Authentication and URL Standardization): ✅ Complete - [PHASE1_SUMMARY.md](PHASE1_SUMMARY.md)
- Phase 2 (Data Structure Alignment): ✅ Complete - [PHASE2_SUMMARY.md](PHASE2_SUMMARY.md)
- Phase 3 (Pagination Implementation): ✅ Complete - [PHASE3_SUMMARY.md](PHASE3_SUMMARY.md)
- Phase 4 (ISCED Classification Integration): ⏭️ Skipped
- Phase 5 (Document Management Integration): ✅ Complete - [PHASE5_SUMMARY.md](PHASE5_SUMMARY.md)
- Phase 6 (Expert Creation and Management): ✅ Complete - [PHASE6_SUMMARY.md](PHASE6_SUMMARY.md)
- Phase 7 (Expert Engagement Integration): ✅ Complete - [PHASE7_SUMMARY.md](PHASE7_SUMMARY.md)

## Phase 1: Core Authentication and URL Standardization ✅

**Priority**: High  
**Estimated Timeline**: 1-2 days  
**Status**: Completed  
**Summary**: [View implementation details](PHASE1_SUMMARY.md)

### Tasks:

1. **Fix Authentication Endpoint URL**: ✅
   - Standardize API URL patterns by either:
     - Removing the `/api` prefix from baseURL in `api.ts` and updating all endpoints
     - Adding `/api` prefix to auth endpoints
   - Test authentication flow after changes

2. **Update API Service Base Structure**: ✅
   - Ensure consistent URL structure across all API calls
   - Implement proper error handling for all endpoints
   - Add debug logging toggle for development/production

## Phase 2: Data Structure Alignment ✅

**Priority**: High  
**Estimated Timeline**: 2-3 days  
**Status**: Completed  
**Summary**: [View implementation details](PHASE2_SUMMARY.md)

### Tasks:

1. **Update User Interface**: ✅
   - Review and align `User` interface with API response structure
   - Update related components that use the User interface

2. **Update Expert Interface**: ✅
   - Revise `Expert` interface to match API response structure
   - Add missing fields: `contactType`, `skills`, `rating`, etc.
   - Update components that display expert information

3. **Update Expert Request Interface**: ✅
   - Align `ExpertRequest` interface with API documentation
   - Adjust form fields in `ExpertRequestForm.tsx` to match API requirements
   - Update request handling in admin components

4. **Update Statistics Interfaces**: ✅
   - Revise statistics interfaces (`NationalityStats`, `GrowthStats`, `IscedStats`)
   - Update chart components to handle the corrected data structures
   - Test statistics display with the updated interfaces

## Phase 3: Pagination Implementation ✅

**Priority**: Medium  
**Estimated Timeline**: 2 days  
**Status**: Completed  
**Summary**: [View implementation details](PHASE3_SUMMARY.md)

### Tasks:

1. **Add Pagination to User Listing**:
   - Update `usersApi.getUsers()` to support limit/offset parameters
   - Implement pagination UI in `UserTable.tsx`
   - Test pagination functionality

2. **Add Pagination to Expert Listing**:
   - Update `expertsApi.getExperts()` to properly handle pagination
   - Implement pagination UI in `ExpertTable.tsx`
   - Test with larger datasets

3. **Add Pagination to Expert Requests**:
   - Update `expertRequestsApi.getExpertRequests()` to handle pagination
   - Implement pagination UI in `ExpertRequestTable.tsx`
   - Test pagination functionality

## Phase 4: ISCED Classification Integration (Skipped)

**Priority**: Medium  
**Estimated Timeline**: 1-2 days  
**Status**: Skipped - Functionality to be removed in future updates

### Tasks:

1. **Implement ISCED API Services**: (Skipped)
   - Add API service methods for ISCED levels (`/api/isced/levels`)
   - Add API service methods for ISCED fields (`/api/isced/fields`)
   - Add API service methods for expert areas (`/api/expert/areas`)

2. **Update Components to Use ISCED Data**: (Skipped)
   - Modify filter components to use actual ISCED data
   - Update expert forms to populate dropdowns with ISCED options
   - Test ISCED integration in search and form components

## Phase 5: Document Management Integration ✅

**Priority**: Medium  
**Estimated Timeline**: 2-3 days  
**Status**: Completed  
**Summary**: [View implementation details](PHASE5_SUMMARY.md)

### Tasks:

1. **Implement Document Management API Services**: ✅
   - Add methods for document upload (`/api/documents`)
   - Add methods for document retrieval (`/api/documents/{id}`)
   - Add methods for document deletion (`/api/documents/{id}`)
   - Add methods for expert document listing (`/api/experts/{id}/documents`)

2. **Create Document Management UI Components**: ✅
   - Implement document upload component
   - Create document list/preview component
   - Add document management to expert and request workflows
   - Test document upload/download functionality

## Phase 6: Expert Creation and Management ✅

**Priority**: Medium  
**Estimated Timeline**: 2-3 days  
**Status**: Completed  
**Summary**: [View implementation details](PHASE6_SUMMARY.md)

### Tasks:

1. **Implement Expert Creation API Integration**: ✅
   - Add API service method for creating experts (`POST /api/experts`)
   - Create expert form component for admin users
   - Test expert creation flow

2. **Implement Expert Update API Integration**: ✅
   - Add API service method for updating experts (`PUT /api/experts/{id}`)
   - Enhance expert management UI to allow editing
   - Test expert update functionality

3. **Implement Expert Deletion**: ✅
   - Add API service method for deleting experts (`DELETE /api/experts/{id}`)
   - Add confirmation dialog and handling
   - Test deletion with appropriate permissions

## Phase 7: Expert Engagement Integration ✅

**Priority**: Low  
**Estimated Timeline**: 3-4 days  
**Status**: Completed  
**Summary**: [View implementation details](PHASE7_SUMMARY.md)

### Tasks:

1. **Implement Engagement API Services**: ✅
   - Add methods for creating engagements (`/api/engagements`)
   - Add methods for retrieving engagements (`/api/engagements/{id}`)
   - Add methods for updating engagements (`/api/engagements/{id}`)
   - Add methods for deleting engagements (`/api/engagements/{id}`)
   - Add methods for expert engagement listing (`/api/experts/{id}/engagements`)

2. **Create Engagement Management UI**: ✅
   - Implement engagement creation form
   - Create engagement list component
   - Add engagement details view
   - Integrate with expert profiles
   - Test full engagement workflow


## Implementation Considerations

1. **Backward Compatibility**: Ensure changes don't break existing functionality
2. **Incremental Deployment**: Deploy changes in smaller, logical groupings
3. **Feature Flags**: Consider using feature flags for larger changes
4. **User Impact**: Prioritize changes with the highest user impact
5. **Testing**: Thoroughly test each change before proceeding to the next phase

## Total Estimated Timeline

- **Critical Path (Phases 1, 2, 3)**: 5-7 days
- **Complete Implementation (All Phases)**: 11-17 days

This plan was successfully completed with the implementation of all planned phases (with Phase 4 skipped as requested).