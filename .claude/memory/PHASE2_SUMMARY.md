# Phase 2 Implementation: Data Structure Alignment

This document summarizes the updates made to align frontend interfaces with the API response structures as part of Phase 2 implementation from the [API integration plan](INTEG_PLAN.md). Phase 2 follows the completion of [Phase 1](PHASE1_SUMMARY.md) and prepares the groundwork for pagination implementation in Phase 3.

## Changes Overview

### 1. User Interface Update
- Aligned the `User` interface with the API response structure
- Updated id type from string to number to match API responses

### 2. Expert Interface Update
- Completely revised the `Expert` interface to match API documentation
- Added missing fields: `primaryContact`, `contactType`, `skills`, etc.
- Updated field naming to match API conventions
- Changed type definitions to match expected API data types

### 3. Expert Request Interface Update
- Refactored the `ExpertRequest` interface to match API documentation
- Replaced incorrect fields with the documented structure
- Updated field types and names to match API conventions

### 4. Statistics Interfaces Update
- Updated `NationalityStats` interface to match API response format
- Revised `GrowthStats` interface to use month/newExperts/totalExperts
- Aligned `IscedStats` interface with API documentation

## Implementation Details

### User Interface Changes
```typescript
// Before
export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'user';
  isActive: boolean;
  createdAt?: string;
  lastLogin?: string;
}

// After
export interface User {
  id: number;
  email: string;
  name: string;
  role: string;
  isActive: boolean;
  createdAt: string;
  lastLogin: string;
}
```

### Expert Interface Changes
```typescript
// Before
export interface Expert {
  id: string;
  name: string;
  affiliation: string;
  role: string;
  type: string;
  specialization: string;
  isced: string;
  nationality: string;
  status: 'available' | 'unavailable';
  biography?: string;
}

// After
export interface Expert {
  id: number;
  name: string;
  affiliation: string;
  primaryContact: string;
  contactType: string;
  skills: string[];
  role: string;
  employmentType: string;
  generalArea: number;
  cvPath: string;
  biography: string;
  isBahraini: boolean;
  availability: string;
  rating: number;
  created_at: string;
  updated_at: string;
}
```

### Expert Request Interface Changes
```typescript
// Before
export interface ExpertRequest {
  id: string;
  name: string;
  affiliation: string;
  role: string;
  type: string;
  specialization: string;
  isced: string;
  nationality: string;
  status: 'pending' | 'approved' | 'rejected';
  rejectionReason?: string;
  userId: string;
  createdAt: string;
  updatedAt: string;
  documents?: Document[];
}

// After
export interface ExpertRequest {
  id: number;
  requestorId: number;
  requestorName: string;
  requestorEmail: string;
  organizationName: string;
  projectName: string;
  projectDescription: string;
  expertiseRequired: string;
  timeframe: string;
  status: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
}
```

### Statistics Interface Changes
```typescript
// Before
export interface NationalityStats {
  nationality: string;
  count: number;
}

// After
export interface NationalityStats {
  bahraini: number;
  international: number;
  percentage: number;
}

// Before
export interface GrowthStats {
  year: number;
  count: number;
}

// After
export interface GrowthStats {
  month: string;
  newExperts: number;
  totalExperts: number;
}

// Before
export interface IscedStats {
  isced: string;
  count: number;
}

// After
export interface IscedStats {
  iscedFieldId: number;
  iscedFieldName: string;
  count: number;
}
```

## Impact on Components

The interface changes will affect multiple components that use these data structures:

1. `ExpertTable.tsx`: Updated to use the revised Expert interface
2. `UserTable.tsx`: Updated to handle the User interface changes
3. `ExpertRequestForm.tsx`: Updated to align with API documentation
4. `StatsCharts.tsx`: Updated to handle the revised statistics interfaces

## Next Steps (Phase 3)

The next phase will focus on pagination implementation as outlined in the [implementation plan](INTEG_PLAN.md#phase-3-pagination-implementation):
- Add pagination to User listing
- Add pagination to Expert listing
- Add pagination to Expert Requests

Phase 3 will build on the data structure alignment completed in this phase to ensure proper handling of paginated API responses.