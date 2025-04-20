# ExpertDB Frontend Development Phased Plan

This document outlines a phased approach to implementing the remaining features of the ExpertDB frontend, focusing on backend integration and feature completion.

## Phase 1: Foundation & Core API Integration (2 weeks)

### Week 1: API Refactoring & Structure Improvements

#### 1.1 API Module Separation (3 days)
- [x] Split `api.ts` into separate modules according to the planned structure:
  - [x] Create `src/api/` directory
  - [x] Implement `auth.ts` for authentication
  - [x] Implement `experts.ts` for expert management
  - [x] Implement `requests.ts` for expert requests
  - [x] Implement `documents.ts` for document management
  - [x] Implement `engagements.ts` for engagement management
  - [x] Implement `phases.ts` for phase planning
  - [x] Implement `statistics.ts` for statistics
  - [x] Implement `areas.ts` for specialization areas
- [x] Create base API client with interceptors for reuse
- [x] Add token refresh mechanism

#### 1.2 Component Reorganization (2 days)
- [x] Create proper component folder structure:
  - [x] Move form components to `/components/forms/`
  - [x] Move table components to `/components/tables/`
  - [x] Move chart components to `/components/charts/`
  - [x] Move modals to `/components/modals/`
- [x] Update imports throughout the application

### Week 2: Utils & Context Implementation

#### 2.1 Utility Functions (2 days)
- [x] Create `/src/utils/` directory
- [x] Implement `formatters.ts` with functions for:
  - [x] Date formatting
  - [x] Currency formatting
  - [x] Name formatting
  - [x] Status badge formatting
- [x] Implement `validators.ts` for form validations
- [x] Implement `permissions.ts` for role-based access control

#### 2.2 Context Implementation (3 days)
- [x] Enhance existing `AuthContext` with better error handling
- [x] Add token refresh to `AuthContext`
- [x] Create `UIContext.tsx` for managing UI state:
  - [x] Sidebar collapsed state
  - [x] Notification system
  - [x] Theme preferences (if applicable)
- [x] Implement hooks for accessing context

## Phase 2: Feature Completion (3 weeks)

### Week 3: Expert Management & Documents

#### 3.1 Expert Management Enhancement (3 days)
- [x] Complete expert filtering functionality
- [x] Implement batch operations
- [x] Add expert area management interface
- [x] Enhance expert detail view
- [x] Add export functionality

#### 3.2 Document Management (2 days)
- [x] Implement document preview feature
- [x] Add document versioning 
- [x] Enhance document listing with filters
- [x] Add document metadata editing

### Week 4: Phase Planning & Engagement

#### 4.1 Phase Planning Interface (3 days)
- [x] Create phase planning creation interface
- [x] Implement expert assignment workflow for phases
- [x] Build admin review interface for applications
- [x] Add phase status tracking

#### 4.2 Engagement Management (2 days)
- [x] Complete engagement listing with filters
- [x] Implement engagement import functionality with validation
- [x] Add engagement editing capabilities
- [x] Create engagement statistics view

### Week 5: Statistics & Advanced Features

#### 5.1 Statistics Dashboard (3 days)
- [x] Implement complete dashboard with key metrics
- [x] Create yearly growth chart
- [x] Build nationality distribution charts
- [x] Add area utilization statistics
- [x] Implement export functionality for reports

#### 5.2 Advanced Features (2 days)
- [x] Add global search functionality
- [x] Implement user preferences saving
- [x] Create system backup and restore interface
- [x] Add batch operations for common tasks

## Phase 3: Polish & Testing (2 weeks) ✅

### Week 6: Form Handling & Loading States ✅

#### 6.1 Form Handling Improvements (3 days) ✅
- [x] Implement React Hook Form consistently across all forms
- [x] Add Zod for validation schemas
- [x] Create reusable form components
- [x] Implement optimistic updates for better UX

#### 6.2 Loading States & Feedback (2 days) ✅
- [x] Add skeleton loaders for data fetching
- [x] Implement progress indicators for file uploads
- [x] Create toast notifications for user feedback
- [x] Add transition animations for state changes

### Week 7: Testing & Bug Fixing

#### 7.1 Testing Implementation (3 days)
- [ ] Write unit tests for critical components
- [ ] Create integration tests for core workflows
- [ ] Test role-based access control
- [ ] Verify API error handling

#### 7.2 Bug Fixing & Final Polish (2 days)
- [ ] Address any identified bugs
- [ ] Improve responsive design for all screen sizes
- [ ] Optimize bundle size with code splitting
- [ ] Conduct final review of all features

## Implementation Details for Phase 3: Form Handling & Loading States

### Form Handling Improvements

The Phase 3 implementation has successfully standardized form handling across the application using React Hook Form with Zod validation:

```typescript
// Form schema definition using Zod
export const expertSchema = z.object({
  id: z.number().optional(),
  name: z.string().min(2, 'Name is required and must be at least 2 characters'),
  affiliation: z.string().min(2, 'Affiliation is required'),
  primaryContact: z.string().min(5, 'Primary contact is required'),
  contactType: z.enum(['email', 'phone', 'linkedin']),
  skills: z.string().min(3, 'Skills are required'),
  role: z.string().min(2, 'Role is required'),
  employmentType: z.enum(['full-time', 'part-time', 'consultant', 'retired', 'other']),
  generalArea: z.string().min(1, 'General area is required'),
  biography: z.string().optional(),
  isBahraini: z.boolean().default(false),
  availability: z.enum(['Available', 'Limited', 'Unavailable']),
  cvFile: z.any().optional(),
});

// Custom hook for form handling with notifications
export function useFormWithNotifications<T extends FieldValues>({
  schema,
  onSuccess,
  ...formProps
}: UseZodFormProps<T> & { onSuccess?: (data: T) => void }) {
  const form = useZodForm<T>({ schema, ...formProps });
  const { addNotification } = useUI();

  const handleSubmitWithNotifications = (
    callback: (data: T) => Promise<{ success: boolean; message?: string }>
  ) => {
    return form.handleSubmit(async (data) => {
      try {
        const result = await callback(data);
        
        if (result.success) {
          addNotification({
            type: 'success',
            message: result.message || 'Operation completed successfully',
            duration: 5000,
          });
          
          if (onSuccess) {
            onSuccess(data);
          }
        } else {
          addNotification({
            type: 'error',
            message: result.message || 'Operation failed',
            duration: 5000,
          });
        }
      } catch (error) {
        addNotification({
          type: 'error',
          message: error instanceof Error ? error.message : 'An unexpected error occurred',
          duration: 5000,
        });
      }
    });
  };

  return {
    ...form,
    handleSubmitWithNotifications,
  };
}
```

### Loading States & UI Feedback

Phase 3 also implemented comprehensive loading state components for better user feedback:

```tsx
// LoadingOverlay component example
export const LoadingOverlay: React.FC<{
  isLoading: boolean;
  children: React.ReactNode;
  spinner?: React.ReactNode;
  label?: string;
  className?: string;
}> = ({ isLoading, children, spinner, label, className = '' }) => {
  if (!isLoading) return <>{children}</>;

  return (
    <div className={`relative ${className}`}>
      <div className="absolute inset-0 flex items-center justify-center bg-white/70 dark:bg-gray-800/70 z-10 rounded-lg">
        {spinner || <LoadingSpinner label={label} />}
      </div>
      <div className="opacity-50 pointer-events-none">{children}</div>
    </div>
  );
};

// Implementation example with the Form component
<LoadingOverlay 
  isLoading={isSubmitting}
  className="w-full"
  label="Submitting request..."
>
  <Form
    form={form}
    onSubmit={form.handleSubmitWithNotifications(onSubmit)}
    className="space-y-4"
    showResetButton={true}
    resetText="Reset"
    onReset={handleFormReset}
    submitText="Submit Request"
  >
    {/* Form fields */}
  </Form>
</LoadingOverlay>
```

### Optimistic Updates

Phase 3 added optimistic updates for a more responsive user experience:

```tsx
// Using optimistic collection for expert management
const {
  items: experts,
  isLoading,
  error,
  addItem,
  updateItem,
  deleteItem
} = useOptimisticCollection<Expert>(
  () => fetchExperts().then(data => data.experts),
  {
    deps: [page, limit, filters],
  }
);

// Example of optimistic update with error handling
const handleEditExpert = async (expert: Expert) => {
  try {
    await updateItem(
      expert,
      (data) => expertsApi.updateExpert(data.id.toString(), data).then(res => res.data),
      {
        successMessage: 'Expert updated successfully',
        errorMessage: 'Failed to update expert'
      }
    );
    setIsEditModalOpen(false);
    setSelectedExpert(null);
  } catch (error) {
    // Error is handled by the hook
    console.error('Error updating expert:', error);
  }
};
```

All core forms in the application now use this standardized approach, providing consistent user experience, validation, and feedback mechanisms.

## Next Steps for Phase 4

1. Complete testing:
   - Add unit tests for validation schemas
   - Test form submission with mocked API responses
   - Verify optimistic updates work correctly

2. Add documentation:
   - Document the new hooks and components
   - Create usage examples for developers

3. Performance optimization:
   - Measure and optimize form rendering performance
   - Implement lazy loading for larger forms