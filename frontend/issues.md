# ExpertDB Frontend Issues Tracker

This document tracks ongoing issues, implementation requirements, and fixed bugs in the ExpertDB frontend. It serves as a companion to [debugging_notes.md](./debugging_notes.md), which contains more detailed solutions and technical notes.

## Documentation Structure

The ExpertDB project uses the following documentation system:

1. **Project-wide Guidelines**: [CLAUDE.md](/CLAUDE.md) in the project root
2. **Issue Management System**: [ISSUES.md](/ISSUES.md) in the project root
3. **Frontend-specific Guidelines**: [frontend/CLAUDE.md](./CLAUDE.md)
4. **Active Issues Tracker**: This document (frontend/issues.md)
5. **Technical Solutions**: [debugging_notes.md](./debugging_notes.md)
6. **Archived Issues**: Individual files in [frontend/issues/](./issues/) directory

When resolving issues, follow this workflow:
1. Document the issue here in issues.md
2. Add technical solutions in debugging_notes.md
3. Create a detailed write-up in frontend/issues/issue_[REFERENCE-CODE].md
4. Update ISSUES.md with a reference to the completed issue

All issues should follow the reference code format: `EDB-YYYYMMDD-XX` as outlined in the [Issue Management System](/ISSUES.md).

## Priority Implementation Requirements

### Navigation and Page Flow
- [x] **URGENT**: Modify landing page to be login page (root at localhost:3000/) - FIXED
- [x] After login, redirect to expert database table/search page instead of panel - FIXED
- [x] Remove AI Expert Panel Suggestion heading and functionality from panel page - FIXED
- [x] Update navigation flow: login → expert database list (with search/filtering) as main page - FIXED

### Feature Adjustments
- [x] Ensure expert request and statistics pages are accessible via navbar - FIXED
- [x] Remove AI integration references from UI temporarily - FIXED
- [x] Implement expert database table view with search and filtering - COMPLETED

## Active Issues

### Authentication & Authorization
- [x] Fix login API error: "API Error Response: {}" (ERR_1 below) - FIXED
- [x] Implement proper redirect after successful login - FIXED 
- [x] Apply RequireAuth component to search page - FIXED
- [x] Apply RequireAuth component to remaining protected routes - FIXED
- [ ] Test admin-only routes and access control

### UI/UX
- [x] Fix panel page content - remove AI suggestions, replace with appropriate content - FIXED
- [x] Implement loading states for all data fetching operations - COMPLETED
- [x] Create responsive expert listing table view - COMPLETED
- [ ] Improve expert profile page with detailed information display

## Error Log

### ERR_1: Login API Error (FIXED)
```
Error: API Error Response: {}
    at createUnhandledError (webpack-internal:///(app-pages-browser)/./node_modules/next/dist/client/components/errors/console-error.js:27:71)
    at handleClientError (webpack-internal:///(app-pages-browser)/./node_modules/next/dist/client/components/errors/use-error-handler.js:45:56)
    at console.error (webpack-internal:///(app-pages-browser)/./node_modules/next/dist/client/components/globals/intercept-console-error.js:47:56)
    at eval (webpack-internal:///(app-pages-browser)/./lib/api.ts:217:17)
    at async Axios.request (webpack-internal:///(app-pages-browser)/./node_modules/axios/lib/core/Axios.js:52:14)
    at async Object.login (webpack-internal:///(app-pages-browser)/./lib/api.ts:168:34)
    at async handleSubmit (webpack-internal:///(app-pages-browser)/./app/login/page.tsx:40:30)
```
**Investigation Status**: RESOLVED
**Fix Implemented**: 
- Improved error handling in login/page.tsx to properly detect and display different types of errors
- Added more specific error messages based on response type
- Fixed page redirect flow after successful login
- Improved authentication status detection

### ERR_2: Select Item Value Error (FIXED)
```
Error: A <Select.Item /> must have a value prop that is not an empty string. This is because the Select value can be set to an empty string to clear the selection and show the placeholder.
    at SelectItem (webpack-internal:///(app-pages-browser)/./node_modules/@radix-ui/react-select/dist/index.mjs:1072:15)
    at _c8 (webpack-internal:///(app-pages-browser)/./components/ui/select.tsx:160:87)
    at ExpertSearch (webpack-internal:///(app-pages-browser)/./app/search/expert-search.tsx:240:132)
    at SearchPage (webpack-internal:///(app-pages-browser)/./app/search/page.tsx:42:100)
    at ClientPageRoot (webpack-internal:///(app-pages-browser)/./node_modules/next/dist/client/components/client-page.js:20:50)
```
**Investigation Status**: RESOLVED
**Fix Implemented**: 
- Fixed SelectItem values for ISCED fields to ensure they always have a valid, non-empty string value
- Added fallback mechanism to generate a valid value when field.id is missing or null
- Ensured proper value handling in select components

## Resolved Issues
See [debugging_notes.md](./debugging_notes.md) for details on fixed issues including:
- Login page incorrectly showing navbar
- Auto-login behavior problems
- RequireAuth component implementation
- JSX syntax error in panel/page.tsx
- SelectItem value error in expert search component

## Implementation Plan & Progress Tracking

### Current Focus (Priority Order)
1. ✅ Fix login functionality and API error - COMPLETED
2. ✅ Implement proper page flow (login → expert list) - COMPLETED
3. ✅ Fix expert search component - COMPLETED
4. ✅ Apply RequireAuth component to remaining protected routes - COMPLETED
5. ✅ Remove AI references and update panel page - COMPLETED
6. ✅ Implement expert listing with search/filtering improvements - COMPLETED
7. ✅ Add loading states for better user experience - COMPLETED
8. Create expert profile page with detailed information
9. Implement expert request workflow

### Testing Strategy
- Test each feature with both admin and user roles
- Verify all API interactions with backend
- Confirm data fetching and state management work correctly

## Useful Commands
```bash
# Run development server
npm run dev

# Type checking
npm run typecheck  

# Linting
npm run lint

# Build application
npm run build
```
