# ExpertDB Issue Management System

## Issue Tracking Guidelines

This document serves as the central reference for issue tracking in the ExpertDB project. When working with Claude, follow these guidelines to maintain organized issue tracking.

### How to Use This System

1. **Issue Identification**
   - Each issue gets a unique reference code: `EDB-YYYYMMDD-XX` where:
     - `EDB` is the project prefix
     - `YYYYMMDD` is the date (e.g., 20250304)
     - `XX` is a sequential number for that day (01, 02, etc.)
   - Example: `EDB-20250304-01` (First issue on March 4, 2025)

2. **Issue Documentation**
   - When Claude identifies an issue, use the following files:
     - `frontend/issue.md` for issues related to the frontend
     - `backend/issue.md` for issues related to the backend
   - After resolving an issue, create detailed documentation:
     - Frontend: `frontend/issues/issue_EDB-YYYYMMDD-XX.md`
     - Backend: `backend/issues/issue_EDB-YYYYMMDD-XX.md`
   - Always include the reference code in all documentation

3. **Working with Claude**
   - Use the appropriate `issue.md` file to provide Claude with current issue details
   - Ask Claude to analyze requirements, propose a plan, implement fixes, and document
   - Reference previous similar issues in the corresponding issues/ directory to help with context

### Issue Lifecycle

1. **Identification**
   - Document issue in the appropriate `issue.md` file with reference code
   - Include error messages, affected components, and reproducible steps

2. **Analysis**
   - Claude reviews the issue and proposes a solution plan
   - Check for similar past issues in `issues/issue_EDB-*.md` files

3. **Implementation**
   - Claude implements the fix across affected files
   - Documents technical details in the appropriate `debugging_notes.md`

4. **Documentation**
   - Create detailed documentation in `issues/issue_EDB-YYYYMMDD-XX.md`
   - Update tracking documents:
     - Frontend: `frontend/issues.md` and `frontend/debugging_notes.md`
     - Backend: `backend/issues.md` and `backend/debugging_notes.md`
   - Cross-reference the issue in `ISSUES.md` for future reference

5. **Verification**
   - Run appropriate test commands to verify the fix
   - Confirm the fix resolves the issue

### History Tracking

This section maintains a chronological list of all issues with their reference codes and links to detailed documentation.

| Reference Code | Title | Status | Component | Documentation |
|---------------|-------|--------|-----------|---------------|
| EDB-20250304-01 | Expert Search Implementation | FIXED | Frontend | [issue_EDB-20250304-01.md](./frontend/issues/issue_EDB-20250304-01.md) |
| EDB-20250304-02 | Login API Error | FIXED | Frontend | See [debugging_notes.md](./frontend/debugging_notes.md) |
| EDB-20250304-03 | Select Item Value Error | FIXED | Frontend | See [debugging_notes.md](./frontend/debugging_notes.md) |
| EDB-20250304-04 | Protected Routes Implementation | FIXED | Frontend | See [debugging_notes.md](./frontend/debugging_notes.md) |
| EDB-20250304-05 | AI References Removal | FIXED | Frontend | See [debugging_notes.md](./frontend/debugging_notes.md) |
| EDB-20250304-06 | Navigation Flow Correction | FIXED | Frontend | See [debugging_notes.md](./frontend/debugging_notes.md) |
| EDB-20250304-07 | UI/UX Improvements and Bug Fixes | FIXED | Frontend | [EDB-20250304-07.md](./frontend/issues/records/EDB-20250304-07.md) |
| EDB-20250304-08 | Root Page Navigation Fix | ACTIVE | Frontend | In progress |

## Current Active Issues

- **EDB-20250304-08**: Root Page Navigation Fix - The root page shows an endless loading spinner instead of properly redirecting to login/search based on authentication state

## Reference Code Lookup

To quickly find issues related to specific areas, use these reference codes:

- **AUTH**: Authentication related issues
- **UI**: User interface issues
- **API**: Backend API integration issues
- **NAV**: Navigation and routing issues
- **DATA**: Data fetching and state management
- **BUILD**: Build and deployment issues

Example search: `AUTH-EDB-20250304-01` for authentication-related issues.