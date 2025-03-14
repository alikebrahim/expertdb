# UI/UX Guidelines for ExpertDB

This document outlines React UI components used in the ExpertDB frontend and their current implementation status. It serves as a reference for maintaining consistent UI/UX across the application and tracking component development.

## Component Library

ExpertDB uses the [shadcn/ui](https://ui.shadcn.com/) component library, which provides accessible, customizable UI components built on Radix UI and styled with Tailwind CSS.

## Current Component Status

### Authentication & User Management

#### Login Page ⚠️
- **Components**:
  - Form (shadcn/ui) - Needs styling
  - Button (shadcn/ui) - Needs styling
  - Input (shadcn/ui) - Needs styling
  - Card (shadcn/ui) - Needs implementation
- **Behavior**:
  - Form submits credentials to `/api/auth/login`
  - JWT stored in localStorage
  - Redirects to dashboard on success
- **Issues**:
  - No loading spinner during authentication
  - Error messages not styled
  - Form validation incomplete
  - Layout and styling not finalized
  - No "Remember Me" functionality

#### User Management (Admin) ⚠️
- **Components**:
  - Table (shadcn/ui) - Partial implementation
  - Dialog (shadcn/ui) - Not implemented
  - Form (shadcn/ui) - Not implemented
- **Status**: Basic user listing implemented, but creation/editing UI incomplete

### Expert Search & Management

#### Search Page ⚠️
- **Components**:
  - Combobox (shadcn/ui) - Not implemented
  - Select (shadcn/ui) - Partial implementation
  - Table (shadcn/ui) - Partial implementation
  - Pagination (shadcn/ui) - Not implemented
  - Card (shadcn/ui) - Not implemented
- **Behavior**:
  - Fetches experts from `/api/experts`
  - Allows filtering by various criteria
- **Issues**:
  - Advanced filtering UI incomplete
  - Results display inconsistent
  - Pagination not implemented
  - Sort controls not fully functional

#### Expert Profile ❌
- **Components**: Not implemented
- **Status**: Backend model exists but frontend page not started

#### Expert Request Form ❌
- **Components**: Not implemented
- **Status**: Backend endpoints ready but frontend form not started

### Dashboard & Statistics

#### Statistics Dashboard ❌
- **Components**: Not implemented
- **Status**: Backend endpoints ready but frontend visualizations not started

### Document Management

#### Document Upload ❌
- **Components**: Not implemented
- **Status**: Backend endpoints ready but frontend UI not started

## Global UI Elements

### Navigation

#### Sidebar ⚠️
- **Components**:
  - Sheet (shadcn/ui) - Partial implementation
  - Button (shadcn/ui) - Implemented
- **Status**: Basic navigation implemented but needs styling and role-based rendering

#### Header ⚠️
- **Components**: 
  - DropdownMenu (shadcn/ui) - Partial implementation
- **Status**: Basic header implemented but needs styling

### Common UI Patterns

#### Loading States ❌
- **Components**:
  - Skeleton (shadcn/ui) - Not implemented
  - Spinner (custom) - Not implemented
- **Status**: No consistent loading indicators implemented

#### Error Handling ❌
- **Components**:
  - Alert (shadcn/ui) - Not implemented
  - Toast (shadcn/ui) - Not implemented
- **Status**: Error feedback inconsistent across application

#### Forms
- **Components**:
  - Form (shadcn/ui)
  - Input (shadcn/ui)
  - Select (shadcn/ui)
  - Checkbox (shadcn/ui)
  - RadioGroup (shadcn/ui)
  - Switch (shadcn/ui)
- **Status**: Components available but not consistently used

## Design System

### Colors
- **Primary**: #0f172a (Slate 900)
- **Secondary**: #6366f1 (Indigo 500)
- **Accent**: #f97316 (Orange 500)
- **Background**: #f8fafc (Slate 50)
- **Text**: #334155 (Slate 700)
- **Error**: #ef4444 (Red 500)
- **Success**: #22c55e (Green 500)
- **Warning**: #f59e0b (Amber 500)

### Typography
- **Font Family**: Inter (system fallback: sans-serif)
- **Base Size**: 16px
- **Headings**:
  - H1: 2rem, 700 weight
  - H2: 1.5rem, 700 weight
  - H3: 1.25rem, 600 weight
  - H4: 1rem, 600 weight
- **Body**: 1rem, 400 weight
- **Small**: 0.875rem, 400 weight

### Spacing
- **Base Unit**: 0.25rem (4px)
- **Common Spacing**:
  - xs: 0.5rem (8px)
  - sm: 0.75rem (12px)
  - md: 1rem (16px)
  - lg: 1.5rem (24px)
  - xl: 2rem (32px)
  - 2xl: 3rem (48px)

### Responsive Breakpoints
- **sm**: 640px
- **md**: 768px
- **lg**: 1024px
- **xl**: 1280px
- **2xl**: 1536px

## Next Steps

### Immediate Priority: Login Page
1. Complete form styling using shadcn/ui Card component
2. Add loading spinner during authentication
3. Implement proper error handling with styled messages
4. Add form validation with error states
5. Test JWT persistence and protected route redirection

### Secondary Priorities
1. Complete expert search filtering UI
2. Implement expert profile page
3. Create expert request form
4. Build statistics dashboard visualizations
5. Complete user management admin UI