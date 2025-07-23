# Debugging Lines Added

This file tracks temporary debugging output added to help diagnose issues. Remove these lines when debugging is complete.

## Issue: Expert Requests Endpoint Returning 0 Results When 1 Should Exist

User reports that the `/api/expert-requests` endpoint is returning 0 expert requests although one exists in the expert_requests table by the same user (user ID: 2).

## File: internal/api/handlers/expert_request.go

### Location 1: Lines 266-271 (in HandleGetExpertRequests function)
**Added after userRole assignment**
```go
	// DEBUG: Log all claims
	log.Debug("=== USER CLAIMS DEBUG ===")
	for key, value := range claims {
		log.Debug("Claim '%s': %v", key, value)
	}
	log.Debug("=== END USER CLAIMS ===")
```

### Location 2: Lines 299-302 (in HandleGetExpertRequests function) 
**Enhanced existing debug log for regular user**
```go
	log.Debug("Regular user (ID: %d) retrieving own expert requests with status: '%s', limit: %d, offset: %d", userID, status, limit, offset)
	requests, err = h.store.ListExpertRequestsByUser(userID, status, limit, offset)
	
	// DEBUG: Log the storage method call result
	log.Debug("Storage method returned %d requests, error: %v", len(requests), err)
```

## File: internal/storage/sqlite/expert_request.go

### Location 1: Lines 250-253 (in ListExpertRequestsByUser function)
**Added after parameter validation**
```go
	// DEBUG: Log input parameters
	log := logger.Get()
	log.Debug("=== ListExpertRequestsByUser DEBUG ===")
	log.Debug("Input parameters: userID=%d, status='%s', limit=%d, offset=%d", userID, status, limit, offset)
```

### Location 2: Lines 284-296 (before executing the query)
**Added database validation checks**
```go
	// DEBUG: First check total count in expert_requests table
	var totalCount int
	s.db.QueryRow("SELECT COUNT(*) FROM expert_requests").Scan(&totalCount)
	log.Debug("Total expert_requests in database: %d", totalCount)
	
	// DEBUG: Check count for this specific user
	var userCount int
	s.db.QueryRow("SELECT COUNT(*) FROM expert_requests WHERE created_by = ?", userID).Scan(&userCount)
	log.Debug("Expert requests for user %d: %d", userID, userCount)
	
	// DEBUG: Log the query and arguments
	log.Debug("Query: %s", query)
	log.Debug("Query args: %v", args)
```

### Location 3: Lines 342-348 (before returning results)
**Added result logging**
```go
	// DEBUG: Log final result
	log.Debug("Returning %d expert requests", len(requests))
	for i, req := range requests {
		log.Debug("Request %d: ID=%d, Name=%s, CreatedBy=%d, Status=%s", i, req.ID, req.Name, req.CreatedBy, req.Status)
	}
	log.Debug("=== END ListExpertRequestsByUser DEBUG ===")
```

## Purpose
- Diagnose why ListExpertRequestsByUser returns 0 results when there should be 1 request
- Verify user authentication and claims are correct
- Check database content and query parameters
- Trace the full flow from API handler to storage layer

## Expected Debug Output Flow
1. User claims logging will show the authenticated user details
2. Handler parameters will show the query filters being applied
3. Storage method will show database counts and query construction
4. Final results will show what's actually returned

## Removal Instructions
1. Remove all debug blocks marked with "DEBUG:" comments
2. Restore the logger import in expert_request.go if not needed elsewhere
3. Delete this DEBUG.md file