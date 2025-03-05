~/dev/expertdb_new  master !8 ?10 ❯ ./run.sh
Starting ExpertDB Initialization
Creating necessary directories...
Database already exists.
Found experts.csv file.
CSV contains approximately 441 records.
Admin user will be created with email: admin@expertdb.com
You can change admin settings by setting environment variables ADMIN_EMAIL, ADMIN_NAME, and ADMIN_PASSWORD
Logging configured with level INFO to directory ./logs
You can change logging settings with LOG_LEVEL and LOG_DIR environment variables
Checking database migrations...
Found 12 migration files that will be applied automatically
CSV import will be performed if the database is new
Importing from ./backend/experts.csv
Docker not detected or Docker Compose not available. Starting services manually...
Starting frontend and backend services in separate processes...
Starting backend server...
Frontend dependencies not found. Installing dependencies...
./run.sh: line 99: cd: frontend: No such file or directory
Starting frontend server...
Services started successfully!
Backend running at http://localhost:8080
Frontend running at http://localhost:3000
Please access the application at http://localhost:3000
Press Ctrl+C to stop all services
./run.sh: line 105: cd: frontend: No such file or directory
2025/03/04 12:03:11 [INFO] server.go:66: Starting ExpertDB initialization...
2025/03/04 12:03:11 [INFO] server.go:70: Configuration loaded successfully
2025/03/04 12:03:11 [INFO] server.go:77: Database directory created: db/sqlite
2025/03/04 12:03:11 [INFO] server.go:83: Upload directory created: ./data/documents
2025/03/04 12:03:11 [INFO] server.go:86: Connecting to database at ./db/sqlite/expertdb.sqlite
2025/03/04 12:03:11 [INFO] storage.go:123: Initializing database schema
2025/03/04 12:03:11 [INFO] storage.go:137: Database type:  (memory=true)
2025/03/04 12:03:11 [INFO] storage.go:142: Using in-memory database with simplified schema
2025/03/04 12:03:11 [INFO] storage.go:331: Database schema initialized successfully
2025/03/04 12:03:11 [INFO] server.go:99: Database connection established successfully
2025/03/04 12:03:11 [INFO] server.go:102: Initializing JWT secret...
2025/03/04 12:03:11 [INFO] server.go:106: JWT secret initialized successfully
2025/03/04 12:03:11 [INFO] server.go:125: Checking for admin user with email: admin@expertdb.com
2025/03/04 12:03:11 [INFO] server.go:128: Admin user not found, creating...
2025/03/04 12:03:11 [INFO] server.go:150: Created default admin user with email: admin@expertdb.com
2025/03/04 12:03:11 [INFO] server.go:156: Creating API server on port 8080
2025/03/04 12:03:11 [INFO] server.go:162: Starting ExpertDB with configuration:
2025/03/04 12:03:11 [INFO] server.go:163: - Port: 8080
2025/03/04 12:03:11 [INFO] server.go:164: - Database: ./db/sqlite/expertdb.sqlite
2025/03/04 12:03:11 [INFO] server.go:165: - Upload Path: ./data/documents
2025/03/04 12:03:11 [INFO] server.go:166: - CORS: *
2025/03/04 12:03:11 [INFO] server.go:167: - AI Service: http://localhost:9000
2025/03/04 12:03:11 [INFO] server.go:168: - Log Level: INFO
2025/03/04 12:03:11 [INFO] server.go:169: - Log Directory: ./logs
2025/03/04 12:03:11 [INFO] server.go:171: Server starting, press Ctrl+C to stop
2025/03/04 12:03:11 [INFO] api.go:41: Setting up API routes...
2025/03/04 12:03:11 [INFO] api.go:122: API server listening on :8080
