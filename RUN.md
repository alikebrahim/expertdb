Backend log:
~/dev/expertdb/backend  master !72 ?14 ❯ air                                                                                                                                         8m 0s ▼

  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \_ v1.61.7, built with Go go1.24.2

watching .
watching data
watching data/documents
watching db
watching db/migrations
watching db/migrations/sqlite
watching db/sqlite
watching http
watching logs
!exclude tmp
building...
# expertdb
./api.go:410:86: syntax error: unexpected newline in composite literal; possibly missing comma or }
failed to build, error: exit status 1
running...
2025/04/15 11:13:38 [INFO] server.go:114: Starting ExpertDB initialization...
2025/04/15 11:13:38 [INFO] server.go:118: Configuration loaded successfully
2025/04/15 11:13:38 [INFO] server.go:125: Database directory created: db/sqlite
2025/04/15 11:13:38 [INFO] server.go:131: Upload directory created: ./data/documents
2025/04/15 11:13:38 [INFO] server.go:134: Connecting to database at ./db/sqlite/expertdb.sqlite
2025/04/15 11:13:38 [INFO] server.go:142: Database connection established successfully
2025/04/15 11:13:38 [INFO] server.go:44: Verifying database schema...
2025/04/15 11:13:38 [INFO] server.go:85: Database schema verification completed successfully
2025/04/15 11:13:38 [INFO] server.go:151: Initializing JWT secret...
2025/04/15 11:13:38 [INFO] server.go:155: JWT secret initialized successfully
2025/04/15 11:13:38 [INFO] server.go:174: Checking for admin user with email: admin@expertdb.com
2025/04/15 11:13:38 [INFO] server.go:201: Admin user already exists, skipping creation
2025/04/15 11:13:38 [INFO] server.go:205: Creating API server on port 8080
2025/04/15 11:13:38 [INFO] server.go:211: Starting ExpertDB with configuration:
2025/04/15 11:13:38 [INFO] server.go:212: - Port: 8080
2025/04/15 11:13:38 [INFO] server.go:213: - Database: ./db/sqlite/expertdb.sqlite
2025/04/15 11:13:38 [INFO] server.go:214: - Upload Path: ./data/documents
2025/04/15 11:13:38 [INFO] server.go:215: - CORS: *
2025/04/15 11:13:38 [INFO] server.go:216: - Log Level: DEBUG
2025/04/15 11:13:38 [INFO] server.go:217: - Log Directory: ./logs
2025/04/15 11:13:38 [INFO] server.go:219: Server starting, press Ctrl+C to stop
2025/04/15 11:13:38 [INFO] api.go:75: Setting up API routes...
2025/04/15 11:13:38 [DEBUG] api.go:84: Registering expert management routes
2025/04/15 11:13:38 [DEBUG] api.go:92: Registering expert request routes
2025/04/15 11:13:38 [DEBUG] api.go:102: Registering document management routes
2025/04/15 11:13:38 [DEBUG] api.go:109: Registering expert engagement routes
2025/04/15 11:13:38 [DEBUG] api.go:117: Registering statistics routes
2025/04/15 11:13:38 [DEBUG] api.go:125: Registering user management routes
2025/04/15 11:13:38 [DEBUG] api.go:136: Setting up CORS middleware
2025/04/15 11:13:38 [DEBUG] api.go:160: Applying middleware chain
2025/04/15 11:13:38 [INFO] api.go:164: API server listening on :8080
2025/04/15 11:13:44 [INFO] auth.go:387: User logged in successfully: admin@expertdb.com (ID: 1, Role: admin)
2025/04/15 11:13:44 [INFO] logger.go:209: HTTP POST /api/auth/login from 127.0.0.1:59496 - 200 (OK) - 185.67274ms
2025/04/15 11:13:44 [INFO] auth.go:467: User creation failed - duplicate email: testuser3@example.com
2025/04/15 11:13:44 [ERROR] api.go:1672: Handler error: POST /api/users - email already exists
2025/04/15 11:13:44 [ERROR] logger.go:209: HTTP POST /api/users from 127.0.0.1:59506 - 500 (Internal Server Error) - 180.222574ms
2025/04/15 11:13:44 [INFO] auth.go:387: User logged in successfully: testuser3@example.com (ID: 2, Role: user)
2025/04/15 11:13:44 [INFO] logger.go:209: HTTP POST /api/auth/login from 127.0.0.1:59516 - 200 (OK) - 185.564964ms
2025/04/15 11:13:44 [DEBUG] auth.go:537: User list retrieved: 2 users returned (limit: 10, offset: 0)
2025/04/15 11:13:44 [INFO] logger.go:209: HTTP GET /api/users from 127.0.0.1:59526 - 200 (OK) - 309.231µs
2025/04/15 11:13:44 [DEBUG] api.go:351: Processing POST /api/experts request
2025/04/15 11:13:44 [DEBUG] api.go:401: Creating expert: Test Expert, Institution: Test University
2025/04/15 11:13:44 [ERROR] api.go:404: Failed to create expert in database: failed to create expert: UNIQUE constraint failed: experts.expert_id
2025/04/15 11:13:44 [ERROR] logger.go:209: HTTP POST /api/experts from 127.0.0.1:59536 - 500 (Internal Server Error) - 356.811µs
2025/04/15 11:13:44 [DEBUG] api.go:189: Processing GET /api/experts request
2025/04/15 11:13:44 [DEBUG] api.go:274: Retrieving experts with filters: map[sort_by:name sort_order:asc], limit: 10, offset: 0
2025/04/15 11:13:44 [ERROR] api.go:277: Failed to list experts: failed to scan expert row: sql: Scan error on column index 11, name "specialized_area": converting NULL to string is unsupported
2025/04/15 11:13:44 [ERROR] api.go:1672: Handler error: GET /api/experts - failed to retrieve experts: failed to scan expert row: sql: Scan error on column index 11, name "specialized_area": converting NULL to string is unsupported
2025/04/15 11:13:44 [ERROR] logger.go:209: HTTP GET /api/experts from 127.0.0.1:59540 - 500 (Internal Server Error) - 537.426µs
2025/04/15 11:13:44 [DEBUG] api.go:1264: Processing POST /api/expert-requests request
2025/04/15 11:13:44 [DEBUG] api.go:1306: Validating expert request fields
2025/04/15 11:13:44 [DEBUG] api.go:1333: Setting default values for expert request
2025/04/15 11:13:44 [DEBUG] api.go:1338: Creating expert request in database: Request Test Expert, Institution: Request University
2025/04/15 11:13:44 [INFO] api.go:1348: Expert request created successfully: ID: 4, Name: Request Test Expert
2025/04/15 11:13:44 [INFO] logger.go:209: HTTP POST /api/expert-requests from 127.0.0.1:59546 - 201 (Created) - 2.876723ms
2025/04/15 11:13:44 [DEBUG] api.go:1376: Processing GET /api/expert-requests request
2025/04/15 11:13:44 [DEBUG] api.go:1387: Using pagination: limit=100, offset=0
2025/04/15 11:13:44 [DEBUG] api.go:1396: Retrieving expert requests with filters: map[]
2025/04/15 11:13:44 [DEBUG] api.go:1404: Returning 4 expert requests
2025/04/15 11:13:44 [INFO] logger.go:209: HTTP GET /api/expert-requests from 127.0.0.1:59560 - 200 (OK) - 500.906µs
2025/04/15 11:13:44 [DEBUG] api.go:1436: Retrieving expert request with ID: 4
2025/04/15 11:13:44 [DEBUG] api.go:1444: Successfully retrieved expert request: ID: 4, Name: Request Test Expert
2025/04/15 11:13:44 [INFO] logger.go:209: HTTP GET /api/expert-requests/4 from 127.0.0.1:59568 - 200 (OK) - 246.922µs
2025/04/15 11:13:45 [DEBUG] api.go:1480: Checking if expert request exists with ID: 4
2025/04/15 11:13:45 [DEBUG] api.go:1498: Processing request update, current status: pending, new status: approved
2025/04/15 11:13:45 [INFO] api.go:1502: Expert request being approved, creating expert record from request data
2025/04/15 11:13:45 [DEBUG] api.go:1509: Generated unique expert ID: EXP-4-1744704825
2025/04/15 11:13:45 [DEBUG] api.go:1535: Creating expert record: Request Test Expert, Institution: Request University
2025/04/15 11:13:45 [INFO] api.go:1548: Expert created successfully from request: Expert ID: 439
2025/04/15 11:13:45 [DEBUG] api.go:1552: Updating expert request ID: 4, Status: approved
2025/04/15 11:13:45 [INFO] api.go:1559: Expert request updated successfully: ID: 4, Status: approved
2025/04/15 11:13:45 [INFO] logger.go:209: HTTP PUT /api/expert-requests/4 from 127.0.0.1:59580 - 200 (OK) - 5.12782ms
2025/04/15 11:13:45 [DEBUG] api.go:568: Processing GET /api/expert/areas request
2025/04/15 11:13:45 [DEBUG] api.go:578: Returning 34 expert areas
2025/04/15 11:13:45 [INFO] logger.go:209: HTTP GET /api/expert/areas from 127.0.0.1:59592 - 200 (OK) - 260.778µs
2025/04/15 11:13:45 [DEBUG] api.go:1085: Processing GET /api/statistics request
2025/04/15 11:13:45 [DEBUG] api.go:1088: Retrieving overall system statistics
2025/04/15 11:13:45 [DEBUG] api.go:1096: Successfully retrieved system statistics
2025/04/15 11:13:45 [INFO] logger.go:209: HTTP GET /api/statistics from 127.0.0.1:59608 - 200 (OK) - 1.163381ms
2025/04/15 11:13:45 [DEBUG] api.go:1085: Processing GET /api/statistics request
2025/04/15 11:13:45 [DEBUG] api.go:1088: Retrieving overall system statistics
2025/04/15 11:13:45 [DEBUG] api.go:1096: Successfully retrieved system statistics
2025/04/15 11:13:45 [INFO] logger.go:209: HTTP GET /api/statistics from 127.0.0.1:59622 - 200 (OK) - 618.641µs

api_test.sh log:
~/dev/expertdb  master !72 ?14 ❯ ./test_api.sh                                                                                                                                        3m 56s
Starting ExpertDB API tests...

===== TESTING AUTHENTICATION =====

Logging in as admin...
Executing: POST http://localhost:8080/api/auth/login
\SUCCESS: Admin login successful. Token: eyJhbGciOiJIUzI...
Creating test user...
Executing: POST http://localhost:8080/api/users
ERROR: API error: email already exists
ERROR: Failed to create test user
Logging in as test user...
Executing: POST http://localhost:8080/api/auth/login
SUCCESS: User login successful. Token: eyJhbGciOiJIUzI...

===== TESTING USER MANAGEMENT =====

Listing all users...
Executing: GET http://localhost:8080/api/users
SUCCESS: Listed 2 users
Getting user details for user ...
Skipping user details as no user ID is available
SUCCESS: User details skipped
Updating user ...
Skipping user update as no user ID is available
SUCCESS: User update skipped

===== TESTING EXPERT MANAGEMENT =====

Creating a new expert...
Executing: POST http://localhost:8080/api/experts
ERROR: API error: Failed to create expert: failed to create expert: UNIQUE constraint failed: experts.expert_id
ERROR: Failed to create expert
Listing all experts...
Executing: GET http://localhost:8080/api/experts
ERROR: API error: failed to retrieve experts: failed to scan expert row: sql: Scan error on column index 11, name "specialized_area": converting NULL to string is unsupported
ERROR: Failed to list experts
Getting expert details for ID ...
Skipping expert details as no expert was created
SUCCESS: Expert details skipped
Updating expert ...
Skipping expert update as no expert was created
SUCCESS: Expert updated successfully

===== TESTING EXPERT REQUEST MANAGEMENT =====

Creating a new expert request...
Executing: POST http://localhost:8080/api/expert-requests
SUCCESS: Expert request created successfully with ID: 4
Listing all expert requests...
Executing: GET http://localhost:8080/api/expert-requests
SUCCESS: Listed 4 expert requests
Getting expert request details for ID 4...
Executing: GET http://localhost:8080/api/expert-requests/4
SUCCESS: Got details for expert request: Request Test Expert
Approving expert request 4...
Executing: PUT http://localhost:8080/api/expert-requests/4
SUCCESS: Expert request approved successfully

===== TESTING DOCUMENT MANAGEMENT =====

Creating a sample document for upload...
Uploading document... (note: this requires a real server with file upload support)
curl -s -X POST -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGV4cGVydGRiLmNvbSIsImV4cCI6MTc0NDc5MTIyNCwibmFtZSI6IkFkbWluIFVzZXIiLCJyb2xlIjoiYWRtaW4iLCJzdWIiOiIxIn0.6jiTEvZQGPwPz3blDkTf1zH9EXQJQ0L4KAeP2yklvVA' -F 'file=@/tmp/sample_cv.txt' -F 'documentType=cv' -F 'expertId=' http://localhost:8080/api/documents
Simulating successful document upload...
SUCCESS: Document uploaded with ID: 1 (simulated)
Listing documents for expert ...
Skipping document listing as no expert was created
SUCCESS: Document listing skipped

===== TESTING EXPERT AREAS =====

Getting expert areas...
Executing: GET http://localhost:8080/api/expert/areas
SUCCESS: Got 34 expert areas

===== TESTING STATISTICS =====

Getting overall statistics...
Executing: GET http://localhost:8080/api/statistics
SUCCESS: Successfully retrieved overall statistics
Getting filtered statistics...
Executing: GET http://localhost:8080/api/statistics?period=month&count=6
SUCCESS: Successfully retrieved filtered statistics

All tests completed!
