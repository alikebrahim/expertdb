# ExpertDB Frontend Debugging Notes

This document contains detailed technical notes on debugging issues and their resolutions. It complements [issues.md](./issues.md), which tracks overall project issues and implementation requirements.

## Documentation Structure

The ExpertDB project uses the following documentation system:

1. **Project-wide Guidelines**: [CLAUDE.md](/CLAUDE.md) in the project root
2. **Issue Management System**: [ISSUES.md](/ISSUES.md) in the project root
3. **Frontend-specific Guidelines**: [frontend/CLAUDE.md](./CLAUDE.md)
4. **Active Issues Tracker**: [issues.md](./issues.md)
5. **Technical Solutions**: This document (debugging_notes.md)
6. **Archived Issues**: Individual files in [frontend/issues/](./issues/) directory

This file is intended for:
- Detailed technical notes on issues encountered
- Code-level solutions for resolving bugs
- Reference patterns for common problems
- Implementation details for feature enhancements

For a complete history of resolved issues with requirements and solutions, see the [issues directory](./issues/).

## Active Debugging Tasks

### Expert Listing Enhancements
**Issue**: Expert listing table view needs improvements
**Location**: 
- app/search/page.tsx (needs UI enhancements)
- app/search/expert-search.tsx (needs filtering improvements)

**Implementation Approach**:
1. Enhance expert listing with better filtering options
2. Improve responsiveness of the expert list view
3. Add pagination or infinite scrolling
4. Implement loading states during data fetching

### Loading State Improvements
**Issue**: Loading states needed for better UX during data fetching
**Location**: 
- Various components that make API calls

**Implementation Approach**:
1. Add consistent loading indicators across components
2. Implement skeleton loading patterns 
3. Ensure error states are handled gracefully
4. Show appropriate loading feedback during authentication

## Resolved Issues

### 1. Login Page UI (Fixed ✅)
**Issue**: Login page incorrectly showing navbar
**Root Cause**: Navbar component being rendered unconditionally
**Solution**: 
```tsx
// Added conditional rendering for navbar
{!isLoginPage && <Navbar />}
```
**Files Modified**:
- app/layout.tsx (Added path check for navbar display)
- components/layout/navbar.tsx (Fixed styling)

### 2. Auto-Login Behavior (Fixed ✅)
**Issue**: App attempting login unnecessarily
**Root Cause**: API client interceptor checking auth on all requests
**Solution**:
```typescript
// Added path-based auth checking
if (!path.includes('/login') && !token) {
  // Redirect to login only for protected routes
  router.push('/login');
}
```
**Files Modified**:
- lib/api.ts (Updated request interceptor logic)

### 3. RequireAuth Implementation (Fixed ✅)
**Issue**: Protected routes accessible without authentication
**Solution**: Created HOC for route protection
```tsx
// Created auth wrapper component
export function RequireAuth({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();
  
  if (isLoading) {
    return <div>Loading...</div>;
  }
  
  if (!isAuthenticated) {
    redirect('/login');
    return null;
  }
  
  return <>{children}</>;
}
```
**Files Modified**:
- components/auth/require-auth.tsx (Created component)
- app/panel/page.tsx (Applied component)

### 4. Build Error in panel/page.tsx (Fixed ✅)
**Issue**: JSX syntax error in panel component
**Error Message**: Unexpected token, expected "}"
**Root Cause**: Improperly nested JSX elements and indentation issues
**Solution**:
- Fixed indentation throughout the file
- Properly closed all JSX tags
- Ensured consistent formatting

**Files Modified**:
- app/panel/page.tsx

## Technical Notes

### Authentication Flow
The authentication flow uses JWT tokens stored in localStorage with the following pattern:
1. User submits credentials to `/api/auth/login`
2. On success, token is stored in localStorage
3. Axios interceptor adds Authorization header to subsequent requests
4. RequireAuth component checks token validity for protected routes

### API Client Structure
The API client (lib/api.ts) has these key components:
- Base axios instance with interceptors
- Login/logout functions
- Error handling middleware
- Request/response transformation

### Common Error Patterns
- 401 errors usually indicate authentication issues
- Empty response objects often point to CORS or preflight issues
- JSX errors typically show as unexpected token errors

## Resolved Issues (March 4, 2025)

### 10. Protected Routes Implementation (Fixed ✅)
**Issue**: Not all routes protected with RequireAuth component
**Root Cause**: Missing RequireAuth wrapper on statistics and request pages
**Solution**: Added RequireAuth component to all protected pages
```jsx
// Added RequireAuth wrapper
<RequireAuth>
  <>
    <Navbar />
    {/* Page content */}
  </>
</RequireAuth>
```
**Files Modified**:
- app/statistics/page.tsx (Added RequireAuth wrapper)
- app/request/page.tsx (Added RequireAuth wrapper)
- app/panel/page.tsx (Already had RequireAuth)

### 11. AI References Removal (Fixed ✅)
**Issue**: Panel page showing AI integration references that aren't ready
**Root Cause**: Early implementation of AI features before they're ready
**Solution**: Replaced AI panel suggestion UI with admin dashboard content
**Files Modified**:
- app/panel/page.tsx (Completely refactored to show admin panel content)

### 12. Navigation Flow Correction (Fixed ✅)
**Issue**: Inconsistent page flow between authentication states
**Root Cause**: Root page (/) not properly managing authenticated vs. unauthenticated states
**Solution**: Updated root page to always redirect appropriately
```typescript
// Root page implementation
useEffect(() => {
  const checkAuth = () => {
    // If authenticated, redirect to search page
    if (authAPI.isAuthenticated()) {
      router.push('/search');
    } else {
      // If not authenticated, redirect to login
      router.push('/login');
    }
  };
  
  checkAuth();
}, [router]);
```
**Files Modified**:
- app/page.tsx (Updated redirect logic)

### 5. Login API Error (Fixed ✅)
**Issue**: Empty API response during login `Error: API Error Response: {}`
**Root Cause**: Inadequate error handling in login form and API client
**Solution**: 
```typescript
// Improved error handling with specific messages
if (error.response) {
  // Server responded with an error
  setError(error.response.data?.error || `Server error: ${error.response.status}`);
} else if (error.request) {
  // Request was made but no response received
  setError('Server not responding. Please try again later.');
} else {
  // Request setup error
  setError('Login failed. Please check your credentials and try again.');
}
```
**Files Modified**:
- app/login/page.tsx (Enhanced error handling)

### 6. Page Flow Correction (Fixed ✅)
**Issue**: Incorrect page flow after login (redirects to panel page)
**Root Cause**: Hardcoded redirect to panel page in login component
**Solution**: Updated redirect destination after successful login to point to search page
**Files Modified**:
- app/login/page.tsx (Changed redirect target)

### 7. Landing Page (Fixed ✅)
**Issue**: Landing page not redirecting to login for unauthenticated users
**Solution**: Added authentication check in useEffect hook on landing page
**Files Modified**:
- app/page.tsx (Added auth check and redirect)

### 8. Protected Search Page (Fixed ✅)
**Issue**: Search page accessible without authentication
**Solution**: Wrapped search page content with RequireAuth component
**Files Modified**:
- app/search/page.tsx (Added RequireAuth wrapper)

### 9. Select Item Value Error (Fixed ✅)
**Issue**: Error with SelectItem needing a non-empty string value
**Error Message**: 
```
Error: A <Select.Item /> must have a value prop that is not an empty string.
```
**Root Cause**: ISCED fields data may have missing or null id values, causing SelectItem to have empty values
**Solution**: 
```tsx
// Added fallback value mechanism
<SelectItem 
  key={field.id} 
  value={field.id ? field.id.toString() : `field-${field.broadCode || Math.random()}`}
>
  {field.broadName}
</SelectItem>
```
**Files Modified**:
- app/search/expert-search.tsx (Added fallback value handling for SelectItem components)

**Technical Details**:
- The radix-ui Select component requires all SelectItem components to have a non-empty string value
- The ISCED fields fetched from the API might have missing or null id values
- The solution ensures each SelectItem always gets a valid string value by:
  1. Using the field.id if available
  2. Falling back to a constructed value using broadCode if available
  3. Using a random number as a last resort if all else fails
- This ensures the Select component won't throw an error even with inconsistent data

## Next Steps
1. ✅ Apply RequireAuth to all remaining protected pages (statistics, request, panel, etc.) - COMPLETED
2. ✅ Remove AI integration references from panel page - COMPLETED 
3. ✅ Implement expert listing with search/filtering enhancements - COMPLETED
4. ✅ Add loading states for data fetching and authentication checks - COMPLETED
5. Improve expert profile page with detailed information display
6. Test the application with different user roles

## Reference Code Snippets

### Proper Error Handling Example
```typescript
// Recommended error handling in API client
try {
  const response = await axios.post('/endpoint', data);
  return response.data;
} catch (error) {
  if (axios.isAxiosError(error)) {
    // Log detailed error information
    console.error('API Error:', {
      status: error.response?.status,
      data: error.response?.data,
      headers: error.response?.headers
    });
    
    // Structured error response
    throw new Error(`API Error: ${error.response?.status || 'Unknown'} - ${JSON.stringify(error.response?.data || {})}`);
  }
  throw error;
}
```

### Defensive Component Props Example
```tsx
// Ensuring component props are always valid
// For components that require strict prop formats
const SafeComponent = ({ id, name, ...props }) => {
  // Ensure id is always a valid string
  const safeId = id ? String(id) : `fallback-${Date.now()}`;
  
  // Ensure name is never undefined
  const safeName = name || 'Default Name';
  
  return <BaseComponent id={safeId} name={safeName} {...props} />;
};
```