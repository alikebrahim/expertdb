# Issue Log

This document tracks bugs, issues, and their resolutions in the ExpertDB project. It serves as a reference for ongoing problems and their fixes to avoid repeating solutions.

## Open Issues

### Authentication

#### AUTH-001: Login Page Styling Incomplete
- **Description**: The login page lacks proper styling and visual feedback
- **Steps to Reproduce**: Visit the login page
- **Expected Behavior**: Professionally styled login form with proper layout
- **Actual Behavior**: Basic unstyled form with minimal visual feedback
- **Priority**: High
- **Status**: Open
- **Notes**: Need to implement shadcn/ui Card component and improve form styling

#### AUTH-002: Loading State Missing During Authentication
- **Description**: No loading indicator shown during login API request
- **Steps to Reproduce**: Submit login form
- **Expected Behavior**: Loading spinner or indicator while request processes
- **Actual Behavior**: UI remains static until response received
- **Priority**: Medium
- **Status**: Open
- **Notes**: Need to add loading state in login form component

#### AUTH-003: Authentication Persistence Needs Validation
- **Description**: Recent fixes for token persistence between sessions need validation
- **Steps to Reproduce**: Login, close browser, reopen application
- **Expected Behavior**: Session should be maintained
- **Actual Behavior**: Need to verify if fix works consistently
- **Priority**: High
- **Status**: Open
- **Notes**: Fixed in recent commit but requires testing

### Expert Management

#### EXP-001: Expert Search Filtering UI Incomplete
- **Description**: Advanced filtering UI in expert search is not fully implemented
- **Steps to Reproduce**: Go to expert search page, try to use advanced filters
- **Expected Behavior**: Full set of filter controls for all filterable properties
- **Actual Behavior**: Basic filtering only, advanced options missing or non-functional
- **Priority**: Medium
- **Status**: Open
- **Notes**: Backend supports filtering but frontend UI incomplete

#### EXP-002: Expert Profile Page Not Implemented
- **Description**: No dedicated page to view full expert details
- **Steps to Reproduce**: Try to view detailed expert information
- **Expected Behavior**: Dedicated page with complete expert profile
- **Actual Behavior**: Only table listing with limited information
- **Priority**: Medium
- **Status**: Open
- **Notes**: Backend endpoint exists but frontend page not created

### User Interface

#### UI-001: Inconsistent Error Handling
- **Description**: Error handling and display varies across the application
- **Steps to Reproduce**: Trigger errors in different parts of the application
- **Expected Behavior**: Consistent error message styling and handling
- **Actual Behavior**: Different error display methods, some errors not displayed
- **Priority**: Medium
- **Status**: Open
- **Notes**: Need standardized error handling approach

#### UI-002: Missing Loading States
- **Description**: Loading indicators inconsistent or missing across the application
- **Steps to Reproduce**: Perform actions that trigger API requests
- **Expected Behavior**: Consistent loading indicators for all async operations
- **Actual Behavior**: Some operations show no loading state
- **Priority**: Medium
- **Status**: Open
- **Notes**: Implement shadcn/ui Skeleton components for loading states

## Recently Fixed Issues

### Authentication

#### AUTH-000: JWT Token Persistence Between Sessions
- **Description**: Authentication state was lost on page refresh
- **Steps to Reproduce**: Login, refresh page
- **Expected Behavior**: Stay logged in after refresh
- **Actual Behavior**: User was logged out and redirected to login page
- **Resolution**: Fixed localStorage usage in AuthContext and added proper initialization
- **Fixed In**: Commit 3dc02ebc
- **Status**: Fixed (needs validation)
- **Notes**: Fix implemented but needs testing across browsers

### Expert Management

#### EXP-000: Basic Expert Search Not Working
- **Description**: Expert search API integration not functioning
- **Steps to Reproduce**: Go to search page
- **Expected Behavior**: List of experts displayed
- **Actual Behavior**: No experts shown or errors in console
- **Resolution**: Fixed API endpoint URL and error handling in fetch function
- **Fixed In**: Commit d942ce5c
- **Status**: Fixed
- **Notes**: Basic functionality working but advanced filtering still needed