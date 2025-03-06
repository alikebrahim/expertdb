# ExpertDB Frontend Testing Guide

This document provides comprehensive instructions for testing the ExpertDB frontend application, with specific guidance for each implemented phase.

## Prerequisites

- Node.js 18.x or newer
- npm 9.x or newer
- A modern web browser (Chrome, Firefox, Edge, or Safari)

## Setup Instructions

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/expertdb.git
   cd expertdb/frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```
   The application will be available at http://localhost:5173/ (or another port if 5173 is in use)

## Phase 1: Authentication Testing

### Testing Authentication Flow

1. **Visit the Login Page**
   - Navigate to http://localhost:5173/login
   - Verify the login form is displayed with email and password fields
   - Check that the page is styled properly with Tailwind CSS

2. **Testing Form Validation**
   - Try submitting with empty fields - form should prevent submission
   - Enter an invalid email format - validate the error message appears
   - Enter a very short password - validate any minimum length requirements

3. **Testing Authentication API Integration**
   - If backend is running:
     - Enter valid credentials and submit the form
     - Verify successful login redirects to the correct page based on role
   - If backend is not running or for development testing:
     - Open browser console and run: `localStorage.setItem("token", "mock-token")`
     - Also set a mock user: 
       ```javascript
       localStorage.setItem("user", JSON.stringify({
         id: "1", 
         name: "Test User", 
         email: "test@example.com", 
         role: "admin"  // Try with both "admin" and "user" roles
       }))
       ```
     - Refresh the page and verify redirection works

4. **Testing Protected Routes**
   - Try accessing a protected route directly without authentication:
     - Visit http://localhost:5173/search
     - Verify redirection to login page
   - With a mock regular user token:
     - Try accessing http://localhost:5173/admin
     - Verify redirection to /search
   - With a mock admin token:
     - Verify access to http://localhost:5173/admin is allowed

5. **Testing Logout Functionality**
   - Currently logout must be tested via console:
     - Run `localStorage.removeItem("token")`
     - Run `localStorage.removeItem("user")`
     - Refresh the page and verify redirection to login

### Common Issues and Troubleshooting

1. **CORS Issues with Backend**
   - If using the API and encountering CORS errors, ensure the backend is properly configured
   - For local development, the backend should allow requests from http://localhost:5173

2. **Authentication State Persistence**
   - If login state is not persisting between refreshes:
     - Check if localStorage is being properly set/retrieved
     - Verify that AuthContext is correctly initialized with token from localStorage

3. **Mock Authentication for Development**
   - If the backend is unavailable, use the localStorage method to test:
     ```javascript
     // Set token and user for testing
     localStorage.setItem("token", "mock-token")
     localStorage.setItem("user", JSON.stringify({
       id: "1", 
       name: "Test User", 
       email: "test@example.com", 
       role: "admin"
     }))
     
     // For regular user testing
     localStorage.setItem("user", JSON.stringify({
       id: "2", 
       name: "Regular User", 
       email: "user@example.com", 
       role: "user"
     }))
     
     // To simulate logout
     localStorage.removeItem("token")
     localStorage.removeItem("user")
     ```

## Upcoming Tests (Future Phases)

As more phases are implemented, testing instructions for each will be added here:

- Phase 2: Expert Database Searching (coming soon)
- Phase 3: Request Submission (pending)
- Phase 4: Statistics Dashboard (pending)
- Phase 5: Admin Panel (pending)
- Phase 6: Polish and Deployment (pending)

## Automated Testing

Automated tests will be implemented in a future phase. Currently, all testing is manual.