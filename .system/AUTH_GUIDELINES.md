# Authentication Guidelines

This document details the JWT authentication and role-based access control implementation in ExpertDB. It serves as a reference for both backend and frontend authentication flow and logic.

## Authentication Flow

### Login Process
1. User submits credentials (email/password) to `/api/auth/login`
2. Server validates credentials against database
3. On success, server generates JWT token containing user details and permissions
4. Server responds with user object and token
5. Frontend stores token in localStorage
6. Frontend updates authentication context with user data
7. User is redirected to appropriate page based on role

### Token Management
- **Storage**: JWT stored in browser's localStorage
- **Expiration**: Tokens valid for 24 hours
- **Persistence**: Recent fixes implemented to maintain authentication between sessions/page refreshes
- **Status**: Persistence fixes need validation

### Protected Routes
- Protected routes redirect unauthenticated users to login page
- Admin-only routes redirect non-admin users to dashboard
- Implementation uses React Router and custom ProtectedRoute component

## Role-Based Access Control

### User Roles
1. **Admin**
   - Full access to all system features
   - Can manage users, approve expert requests
   - Can create/edit/delete experts directly
   - Can view all statistics and reports

2. **User**
   - Can search and view experts
   - Can submit expert requests (cannot directly create experts)
   - Limited access to statistics
   - Cannot manage other users

### Backend Implementation
- Middleware verifies token and role for protected endpoints
- Two middleware functions:
  - `authMiddleware`: Verifies token validity
  - `roleMiddleware`: Checks user role

### Frontend Implementation
- AuthContext provides authentication state and functions
- ProtectedRoute component handles route protection
- Role-based UI rendering shows/hides elements based on user role
- Admin features conditionally rendered using `isAdmin` check from auth context

## JWT Implementation Details

### Token Structure
- **Header**: Standard JWT header with algorithm information
- **Payload**:
  - `sub`: User ID
  - `name`: User's name
  - `email`: User's email
  - `role`: User's role ("admin" or "user")
  - `exp`: Token expiration timestamp
- **Signature**: HMAC SHA-256 signature with server secret

### Backend JWT Functions
- `generateToken(user *User) (string, error)`: Creates new JWT for user
- `parseToken(tokenString string) (*jwt.Token, error)`: Validates and parses JWT
- `getUserFromToken(tokenString string) (*User, error)`: Extracts user data from JWT

### Frontend JWT Handling
- Token stored in localStorage with key `expertdb_token`
- Token added to all API requests via Authorization header
- Token format: `Bearer <token>`
- AuthContext checks token on initial load to restore session

## Current Authentication Issues

### Backend Issues
- ✅ **JWT Generation**: Working correctly
- ✅ **Middleware Protection**: Working correctly
- ⚠️ **Default Admin Creation**: Works but needs better first-run detection

### Frontend Issues
- ⚠️ **Token Persistence**: Recent fixes need validation
- ⚠️ **Login Form**: Needs styling and improved error handling
- ⚠️ **AuthContext**: Working but needs validation with protected routes
- ❌ **Loading States**: No loading indicators during authentication
- ⚠️ **Error Handling**: Basic error display but needs improvement

## Authentication Testing Plan

1. **Login Validation**
   - Test valid credentials (both admin and user roles)
   - Test invalid credentials (wrong password, non-existent user)
   - Test empty form submission
   - Verify error messages

2. **Token Persistence**
   - Verify authentication persists after page refresh
   - Verify authentication persists after browser restart
   - Test token expiration handling

3. **Protected Routes**
   - Verify unauthenticated users are redirected to login
   - Verify users can't access admin-only routes
   - Verify authenticated users can access appropriate routes

4. **Role-Based UI**
   - Verify admin users see admin-only features
   - Verify regular users don't see admin features
   - Test conditional rendering of UI elements

## Next Authentication Tasks

1. **Login Page Improvements**
   - Complete login form styling with shadcn/ui
   - Add loading spinner during authentication
   - Improve error message display
   - Add form validation with error states

2. **Persistence Validation**
   - Test and validate recent auth persistence fixes
   - Implement proper session management
   - Add token refresh mechanism (if needed)

3. **UI Enhancement**
   - Implement role-based sidebar navigation
   - Add user profile/settings page
   - Improve login/logout UX
   - Add "Remember Me" functionality