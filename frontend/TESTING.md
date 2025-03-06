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
     // For admin testing
     localStorage.setItem("token", "mock-token");
     localStorage.setItem("user", JSON.stringify({
       id: "1", 
       name: "Test Admin", 
       email: "admin@example.com", 
       role: "admin"
     }));
     
     // For regular user testing
     localStorage.setItem("token", "mock-token");
     localStorage.setItem("user", JSON.stringify({
       id: "2", 
       name: "Test User", 
       email: "user@example.com", 
       role: "user"
     }));
     
     // To simulate logout
     localStorage.removeItem("token");
     localStorage.removeItem("user");
     ```

## Phase 2: Expert Database Search Testing

### Testing Search Functionality

1. **Access the Search Page**
   - Log in to the application (see Authentication Testing above)
   - Verify you are redirected to the search page or navigate to http://localhost:5173/search
   - Confirm the page loads with filter controls and an empty expert table
   - Check that shadcn/ui components are properly styled

2. **Testing Filters**
   - **Name Filter**:
     - Enter a name in the search field
     - Verify the API request updates with the name parameter
     - Check that results are filtered by name
   
   - **Affiliation Filter**:
     - Enter an affiliation in the search field
     - Verify the API request updates with the affiliation parameter
     - Check that results are filtered by affiliation
   
   - **ISCED Field Filter**:
     - Select a field from the ISCED dropdown
     - Verify the API request updates with the isced_field_id parameter
     - Check that results are filtered by ISCED field
   
   - **Bahraini Status Filter**:
     - Select "Bahraini" or "Non-Bahraini" from the dropdown
     - Verify the API request updates with the is_bahraini parameter
     - Check that results are filtered by Bahraini status
   
   - **Availability Filter**:
     - Select "Available" or "Not Available" from the dropdown
     - Verify the API request updates with the is_available parameter
     - Check that results are filtered by availability

3. **Testing Sorting**
   - Click the "Name" column header
   - Verify the arrow indicator changes direction
   - Check that experts are sorted alphabetically by name
   - Click again to verify reverse sorting works

4. **Testing Pagination**
   - With more than 10 results, verify the "Next" button is enabled
   - Click "Next" and verify page parameter is incremented
   - Verify results update to show the next page
   - Click "Previous" and verify navigation back to the first page
   - Verify "Previous" button is disabled on the first page
   - Verify "Next" button is disabled when fewer than 10 results are shown

5. **Testing Loading and Error States**
   - Temporarily modify API URL to cause an error
   - Verify the error message is displayed
   - Restore the correct API URL
   - Observe the loading indicator when fetching data

### Testing with Mock Data

If the backend is unavailable, you can modify the Search component temporarily to use mock data:

```javascript
// Add this near the top of the file
const mockExperts: Expert[] = [
  { id: "1", name: "John Doe", affiliation: "University A", is_bahraini: true, isced_field_id: "1", is_available: true, bio: "Expert in field" },
  { id: "2", name: "Jane Smith", affiliation: "University B", is_bahraini: false, isced_field_id: "2", is_available: true, bio: "Researcher" },
  // Add more mock data as needed
];

const mockIscedFields: IscedField[] = [
  { id: "1", name: "Computer Science" },
  { id: "2", name: "Engineering" },
  // Add more mock fields
];

// Then in your useEffect:
useEffect(() => {
  const fetchData = async () => {
    try {
      // Comment out the actual API calls
      // const [expertsData, fieldsData] = await Promise.all([
      //   getExperts(filters),
      //   getIscedFields(),
      // ]);
      
      // Use mock data instead
      setTimeout(() => {
        setExperts(mockExperts.filter(expert => 
          expert.name.toLowerCase().includes(filters.name.toLowerCase())
        ));
        setIscedFields(mockIscedFields);
        setLoading(false);
      }, 500); // Simulate loading delay
    } catch {
      setError("Failed to load data");
      setLoading(false);
    }
  };
  fetchData();
}, [filters]);
```

### Common Issues and Troubleshooting

1. **No Data Loading**
   - Check the network tab to see if API requests are being made
   - Verify the API endpoint URLs are correct
   - Confirm the backend server is running
   - Check for CORS issues in the browser console

2. **Filtering Not Working**
   - Verify filter parameters are correctly included in the API request
   - Check that API parameters match what the backend expects
   - Test each filter individually to isolate issues

3. **Component Styling Issues**
   - Ensure shadcn/ui components are properly imported
   - Check that Tailwind CSS is correctly configured
   - Verify the component classes match shadcn/ui documentation

## Upcoming Tests (Future Phases)

As more phases are implemented, testing instructions for each will be added here:

- Phase 3: Request Submission (pending)
- Phase 4: Statistics Dashboard (pending)
- Phase 5: Admin Panel (pending)
- Phase 6: Polish and Deployment (pending)

## Automated Testing

Automated tests will be implemented in a future phase. Currently, all testing is manual.