~/de/expertdb_grok/backend  master ⇡1 !72 ?48 ❯ air                                       12s

  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \_ v1.52.3, built with Go go1.23.0

watching .
watching cmd
watching cmd/import_csv
watching data
watching data/documents
watching data/documents/certificates
watching data/documents/cv
watching data/documents/publications
watching db
watching db/migrations
watching db/migrations/sqlite
watching db/sqlite
watching http
watching issues
watching logs
watching scripts
watching testdb
!exclude tmp
building...
# expertdb
./expert_request_operations.go:63:15: no new variables on left side of :=
failed to build, error: exit status 1
running...
2025/03/12 13:24:03 [INFO] server.go:66: Starting ExpertDB initialization...
2025/03/12 13:24:03 [INFO] server.go:70: Configuration loaded successfully
2025/03/12 13:24:03 [INFO] server.go:77: Database directory created: db/sqlite
2025/03/12 13:24:03 [INFO] server.go:83: Upload directory created: ./data/documents
2025/03/12 13:24:03 [INFO] server.go:86: Connecting to database at ./db/sqlite/expertdb.sqlite
2025/03/12 13:24:03 [INFO] storage.go:124: Initializing database schema
2025/03/12 13:24:03 [INFO] storage.go:138: Database type:  (memory=true)
2025/03/12 13:24:03 [INFO] storage.go:143: Using in-memory database with simplified schema
2025/03/12 13:24:03 [INFO] storage.go:380: Database schema initialized successfully
2025/03/12 13:24:03 [INFO] server.go:99: Database connection established successfully
2025/03/12 13:24:03 [INFO] server.go:102: Initializing JWT secret...
2025/03/12 13:24:03 [INFO] server.go:106: JWT secret initialized successfully
2025/03/12 13:24:03 [INFO] server.go:125: Checking for admin user with email: ali@edb.com
2025/03/12 13:24:03 [INFO] server.go:128: Admin user not found, creating...
2025/03/12 13:24:03 [INFO] server.go:150: Created default admin user with email: ali@edb.com
2025/03/12 13:24:03 [INFO] server.go:156: Creating API server on port 8080
2025/03/12 13:24:03 [INFO] server.go:162: Starting ExpertDB with configuration:
2025/03/12 13:24:03 [INFO] server.go:163: - Port: 8080
2025/03/12 13:24:03 [INFO] server.go:164: - Database: ./db/sqlite/expertdb.sqlite
2025/03/12 13:24:03 [INFO] server.go:165: - Upload Path: ./data/documents
2025/03/12 13:24:03 [INFO] server.go:166: - CORS: *
2025/03/12 13:24:03 [INFO] server.go:167: - AI Service: http://localhost:9000
2025/03/12 13:24:03 [INFO] server.go:168: - Log Level: INFO
2025/03/12 13:24:03 [INFO] server.go:169: - Log Directory: ./logs
2025/03/12 13:24:03 [INFO] server.go:171: Server starting, press Ctrl+C to stop
2025/03/12 13:24:03 [INFO] api.go:41: Setting up API routes...
2025/03/12 13:24:03 [INFO] api.go:122: API server listening on :8080
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP POST /api/auth/login from [::1]:54612 - 200 (OK) - 181.601992ms
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP POST /api/users from [::1]:54616 - 201 (Created) - 176.273686ms
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP POST /api/auth/login from [::1]:54632 - 200 (OK) - 169.907959ms
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/users from [::1]:54646 - 200 (OK) - 298.48µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/users/2 from [::1]:54662 - 200 (OK) - 146.039µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP PUT /api/users/2 from [::1]:54670 - 200 (OK) - 215.732µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP POST /api/experts from [::1]:54676 - 201 (Created) - 316.574µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/experts from [::1]:54692 - 200 (OK) - 478.063µs
2025/03/12 13:24:06 [WARN] logger.go:209: HTTP GET /api/experts/1 from [::1]:54698 - 404 (Not Found) - 207.186µs
2025/03/12 13:24:06 [WARN] logger.go:209: HTTP PUT /api/experts/1 from [::1]:54710 - 404 (Not Found) - 184.052µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP POST /api/expert-requests from [::1]:54714 - 201 (Created) - 242.613µs
2025/03/12 13:24:06 [ERROR] logger.go:209: HTTP GET /api/expert-requests from [::1]:54718 - 500 (Internal Server Error) - 90.723µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/expert-requests/1 from [::1]:54720 - 200 (OK) - 148.904µs
2025/03/12 13:24:06 [ERROR] logger.go:209: HTTP PUT /api/expert-requests/1 from [::1]:54728 - 500 (Internal Server Error) - 374.215µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/experts/1/documents from [::1]:54740 - 200 (OK) - 126.672µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/isced/levels from [::1]:54750 - 200 (OK) - 91.926µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/isced/fields from [::1]:54764 - 200 (OK) - 69.603µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/statistics from [::1]:54768 - 200 (OK) - 583.625µs
2025/03/12 13:24:06 [INFO] logger.go:209: HTTP GET /api/statistics from [::1]:54776 - 200 (OK) - 401.526µs
