Starting ExpertDB API tests...

[0;33m===== TESTING AUTHENTICATION =====[0m

Logging in as admin...
Executing: POST http://localhost:8080/api/auth/login
[0;32mSUCCESS: Admin login successful. Token: eyJhbGciOiJIUzI...[0m
Creating test user...
Executing: POST http://localhost:8080/api/users
[0;32mSUCCESS: Test user created successfully. ID: 2[0m
Logging in as test user...
Executing: POST http://localhost:8080/api/auth/login
[0;32mSUCCESS: User login successful. Token: eyJhbGciOiJIUzI...[0m

[0;33m===== TESTING USER MANAGEMENT =====[0m

Listing all users...
Executing: GET http://localhost:8080/api/users
[0;32mSUCCESS: Listed 2 users[0m
Getting user details for user 2...
Executing: GET http://localhost:8080/api/users/2
[0;32mSUCCESS: Got details for user: Test User[0m
Updating user 2...
Executing: PUT http://localhost:8080/api/users/2
[0;32mSUCCESS: User updated successfully[0m

[0;33m===== TESTING EXPERT MANAGEMENT =====[0m

Creating a new expert...
Executing: POST http://localhost:8080/api/experts
[0;32mSUCCESS: Expert created successfully with ID: 1[0m
Listing all experts...
Executing: GET http://localhost:8080/api/experts
[0;32mSUCCESS: Listed 1 experts[0m
Getting expert details for ID 1...
Executing: GET http://localhost:8080/api/experts/1
[0;31mERROR: API error: Expert not found[0m
[0;31mERROR: Failed to get expert details[0m
Updating expert 1...
Executing: PUT http://localhost:8080/api/experts/1
[0;31mERROR: API error: Expert not found[0m
[0;31mERROR: Failed to update expert[0m

[0;33m===== TESTING EXPERT REQUEST MANAGEMENT =====[0m

Creating a new expert request...
Executing: POST http://localhost:8080/api/expert-requests
[0;32mSUCCESS: Expert request created successfully with ID: 1[0m
Listing all expert requests...
Executing: GET http://localhost:8080/api/expert-requests
[0;31mERROR: API error: Failed to retrieve expert requests: failed to query expert requests: no such column: rejection_reason[0m
[0;31mERROR: Failed to list expert requests[0m
Getting expert request details for ID 1...
Executing: GET http://localhost:8080/api/expert-requests/1
[0;32mSUCCESS: Got details for expert request: Request Test Expert[0m
Approving expert request 1...
Executing: PUT http://localhost:8080/api/expert-requests/1
[0;31mERROR: API error: Failed to create expert from request: failed to create expert: UNIQUE constraint failed: experts.expert_id[0m
[0;31mERROR: Failed to approve expert request[0m

[0;33m===== TESTING DOCUMENT MANAGEMENT =====[0m

Creating a sample document for upload...
Uploading document... (note: this requires a real server with file upload support)
curl -s -X POST -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsaUBlZGIuY29tIiwiZXhwIjoxNzQxODYxNDQ2LCJuYW1lIjoiQWxpIEtoYWxpZCIsInJvbGUiOiJhZG1pbiIsInN1YiI6IjEifQ.QD-bwiNnq0XZ_nee0bkAKnYkIgNnXtZaVUg2LGvWRMU' -F 'file=@/tmp/sample_cv.txt' -F 'documentType=cv' -F 'expertId=1' http://localhost:8080/api/documents
Simulating successful document upload...
[0;32mSUCCESS: Document uploaded with ID: 1 (simulated)[0m
Listing documents for expert 1...
Executing: GET http://localhost:8080/api/experts/1/documents
[0;32mSUCCESS: Listed 0 documents for the expert[0m

[0;33m===== TESTING ISCED CLASSIFICATIONS =====[0m

Getting ISCED levels...
Executing: GET http://localhost:8080/api/isced/levels
[0;32mSUCCESS: Got 0 ISCED levels[0m
Getting ISCED fields...
Executing: GET http://localhost:8080/api/isced/fields
[0;32mSUCCESS: Got 0 ISCED fields[0m

[0;33m===== TESTING STATISTICS =====[0m

Getting overall statistics...
Executing: GET http://localhost:8080/api/statistics
[0;32mSUCCESS: Successfully retrieved overall statistics[0m
Getting filtered statistics...
Executing: GET http://localhost:8080/api/statistics?period=month&count=6
[0;32mSUCCESS: Successfully retrieved filtered statistics[0m

[0;32mAll tests completed![0m
