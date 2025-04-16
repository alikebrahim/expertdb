#!/bin/bash

# ==============================================================================
# ExpertDB API Test Script (Comprehensive Version)
#
# Description:
#   Tests the core API endpoints of the ExpertDB application with enhanced coverage.
#   Focuses on expert creation with edge cases, validation, and robustness.
#   Provides detailed logging of requests and responses.
#   Uses unique data for creation tests to ensure re-runnability.
#   Handles dependencies and known backend issues gracefully.
#
# Usage:
#   ./test_api.sh
#
# Requirements:
#   - curl: For making HTTP requests.
#   - jq: For parsing JSON responses.
#   - A running instance of the ExpertDB backend API (default: http://localhost:8080).
# ==============================================================================

# --- Configuration ---
BASE_URL="http://localhost:8080"
ADMIN_EMAIL="admin@expertdb.com"
ADMIN_PASSWORD="adminpassword" # Default password set in server.go
TIMESTAMP=$(date +%s)
USER_EMAIL="testuser${TIMESTAMP}@example.com"
USER_PASSWORD="password123"
USER_NAME="Test User ${TIMESTAMP}"
EXPERT_NAME="Test Expert ${TIMESTAMP}"
EXPERT_EMAIL="expert${TIMESTAMP}@example.com"
EXPERT_PHONE="+9731111${TIMESTAMP: -4}"
REQUEST_NAME="Request Expert ${TIMESTAMP}"
REQUEST_EMAIL="request${TIMESTAMP}@example.com"

# --- Log Files ---
LOG_DIR="logs"
RUN_LOG_FILE="${LOG_DIR}/api_test_run_$(date +%Y%m%d_%H%M%S).log"
REQUEST_LOG_FILE="/tmp/api_request_details.log" # Temporary log for current request details

# --- State Variables ---
ADMIN_TOKEN=""
USER_TOKEN=""
USER_ID=""
EXPERT_ID=""
EXPERT_INTERNAL_ID=""
EXPERT_REQUEST_ID=""
DOCUMENT_ID=""
GENERAL_AREA_ID="1" # Default, will be updated dynamically

# --- Test Counters ---
declare -A test_stats
test_stats["total"]=0
test_stats["passed"]=0
test_stats["failed"]=0
test_stats["skipped"]=0
declare -A test_details # Store details for summary

# --- Colors ---
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;94m'
NC='\033[0m' # No Color

# --- Utility Functions ---

# Function to print messages to console and log file
log_message() {
    local level="$1"
    local message="$2"
    local color="$NC"
    local log_prefix=""

    case "$level" in
        HEADER) color="$BLUE"; log_prefix="[HEADER] ";;
        INFO) color="$NC"; log_prefix="[INFO]   ";;
        STEP) color="$YELLOW"; log_prefix="[STEP]   ";;
        SUCCESS) color="$GREEN"; log_prefix="[SUCCESS]";;
        ERROR) color="$RED"; log_prefix="[ERROR]  ";;
        SKIP) color="$YELLOW"; log_prefix="[SKIP]   ";;
        DETAIL) color="$NC"; log_prefix="[DETAIL] ";; # Only log to file
        *) color="$NC"; log_prefix="[LOG]    ";;
    esac

    # Log to file always
    echo "$(date '+%Y-%m-%d %H:%M:%S') ${log_prefix} ${message}" >> "$RUN_LOG_FILE"

    # Print to console unless it's DETAIL
    if [ "$level" != "DETAIL" ]; then
        echo -e "${color}${message}${NC}"
    fi
}

print_header() {
    log_message "HEADER" "\n===== $1 ====="
}

# Function to record test results
record_test_result() {
    local test_name=$1
    local status=$2 # passed, failed, skipped
    local detail=${3:-}

    test_stats["total"]=$((test_stats["total"] + 1))
    test_stats["$status"]=$((test_stats["$status"] + 1))
    test_details["$test_name"]="$status - $detail"
}

# Function to make API requests with detailed logging
# Usage: api_request <METHOD> <ENDPOINT> [JSON_DATA] [TOKEN] [EXPECTED_STATUS]
api_request() {
    local method=$1
    local endpoint=$2
    local data=$3 # Optional JSON data string
    local token=$4 # Optional JWT token
    local expected_status=$5 # Optional expected HTTP status

    local test_name="$method $endpoint"
    local output_file="/tmp/api_response.json"
    local status_code_file="/tmp/api_status_code.txt"
    local headers_file="/tmp/api_response_headers.txt"

    # Clear temporary request detail log
    > "$REQUEST_LOG_FILE"

    log_message "STEP" "Executing: $test_name"

    local curl_cmd="curl -s -w '%{http_code}' -o '$output_file' --dump-header '$headers_file'"
    curl_cmd+=" -X $method"

    # Add Authorization header if token is provided
    if [ -n "$token" ]; then
        curl_cmd+=" -H 'Authorization: Bearer $token'"
        echo "Authorization: Bearer ${token:0:10}..." >> "$REQUEST_LOG_FILE" # Log redacted token
    fi

    # Add Content-Type and data if provided
    if [ -n "$data" ]; then
        curl_cmd+=" -H 'Content-Type: application/json'"
        curl_cmd+=" -d '$data'"
        echo "Request Body:" >> "$REQUEST_LOG_FILE"
        echo "$data" | jq '.' >> "$REQUEST_LOG_FILE" 2>/dev/null || echo "$data" >> "$REQUEST_LOG_FILE"
    else
        echo "Request Body: (None)" >> "$REQUEST_LOG_FILE"
    fi

    # Construct full URL
    local url="$BASE_URL$endpoint"
    curl_cmd+=" '$url'"

    # Log the command (with token redacted for security)
    log_message "DETAIL" "URL: $url"
    log_message "DETAIL" "Command: $curl_cmd"
    echo "--- Request ---" >> "$RUN_LOG_FILE"
    cat "$REQUEST_LOG_FILE" >> "$RUN_LOG_FILE"
    echo "--- End Request ---" >> "$RUN_LOG_FILE"

    # Execute the command
    http_status=$(eval "$curl_cmd")
    local curl_exit_code=$?

    # Log response details
    log_message "DETAIL" "Curl Exit Code: $curl_exit_code"
    log_message "DETAIL" "HTTP Status: $http_status"
    echo "--- Response ---" >> "$RUN_LOG_FILE"
    echo "HTTP Status: $http_status" >> "$RUN_LOG_FILE"
    echo "Response Headers:" >> "$RUN_LOG_FILE"
    cat "$headers_file" >> "$RUN_LOG_FILE"
    echo "Response Body:" >> "$RUN_LOG_FILE"
    jq '.' "$output_file" >> "$RUN_LOG_FILE" 2>/dev/null || cat "$output_file" >> "$RUN_LOG_FILE"
    echo "" >> "$RUN_LOG_FILE"
    echo "--- End Response ---" >> "$RUN_LOG_FILE"

    # --- Check for errors ---

    # Check curl execution error
    if [ $curl_exit_code -ne 0 ]; then
        log_message "ERROR" "curl command failed with exit code $curl_exit_code."
        record_test_result "$test_name" "failed" "Curl command error"
        return 1
    fi

    # Check HTTP status code
    if [ -n "$expected_status" ] && [ "$http_status" -ne "$expected_status" ]; then
        log_message "ERROR" "Expected status $expected_status but received $http_status."
        jq '.' "$output_file" # Print response body to console on unexpected status
        record_test_result "$test_name" "failed" "Expected status $expected_status, got $http_status"
        return 1
    elif [ "$http_status" -ge 400 ]; then
        log_message "ERROR" "API request failed with HTTP status $http_status."
        jq '.' "$output_file" # Print response body to console on error status
        record_test_result "$test_name" "failed" "HTTP status $http_status"
        return 1
    fi

    # Check for API-level error in JSON response
    if jq -e . "$output_file" > /dev/null 2>&1; then
        if jq -e '.error' "$output_file" > /dev/null 2>&1; then
            error_msg=$(jq -r '.error' "$output_file")
            log_message "ERROR" "API returned error: $error_msg"
            record_test_result "$test_name" "failed" "API Error: $error_msg"
            return 1
        fi
    elif [ -s "$output_file" ]; then
        if [ "$http_status" -lt 300 ]; then
            log_message "DETAIL" "Received non-JSON response with success status $http_status."
        else
            log_message "ERROR" "Received non-JSON response with error status $http_status."
            record_test_result "$test_name" "failed" "Non-JSON response with status $http_status"
            return 1
        fi
    fi

    log_message "SUCCESS" "$test_name completed successfully (Status: $http_status)."
    record_test_result "$test_name" "passed" "Status $http_status"
    return 0
}

# --- Test Sections ---

# 0. Setup: Fetch Valid General Area
test_setup() {
    print_header "TESTING SETUP"

    # Fetch expert areas to get a valid generalArea ID
    api_request "GET" "/api/expert/areas" "" "$ADMIN_TOKEN" 200
    if [ $? -eq 0 ]; then
        GENERAL_AREA_ID=$(jq -r '.[0].id // 1' /tmp/api_response.json)
        local areas_count=$(jq '. | length' /tmp/api_response.json)
        log_message "INFO" "Retrieved $areas_count expert areas. Using generalArea ID: $GENERAL_AREA_ID"
    else
        log_message "WARN" "Failed to retrieve expert areas. Using default generalArea ID: $GENERAL_AREA_ID"
    fi
}

# 1. Test Authentication
test_auth() {
    print_header "TESTING AUTHENTICATION"

    # Retry login up to 3 times
    local retries=3
    local attempt=1
    while [ $attempt -le $retries ]; do
        log_message "INFO" "Attempting admin login ($attempt/$retries)..."
        api_request "POST" "/api/auth/login" '{"email":"'"$ADMIN_EMAIL"'","password":"'"$ADMIN_PASSWORD"'"}' "" 200
        if [ $? -eq 0 ]; then
            ADMIN_TOKEN=$(jq -r '.token' /tmp/api_response.json)
            if [ -n "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "null" ]; then
                log_message "INFO" "Admin login successful."
                break
            fi
            log_message "ERROR" "Failed to extract admin token."
        fi
        attempt=$((attempt + 1))
        sleep 1
    done

    if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
        log_message "ERROR" "Admin login failed after $retries attempts. Cannot proceed."
        exit 1
    fi

    # Create a test user
    local user_payload='{"name":"'"$USER_NAME"'","email":"'"$USER_EMAIL"'","password":"'"$USER_PASSWORD"'","role":"user","isActive":true}'
    api_request "POST" "/api/users" "$user_payload" "$ADMIN_TOKEN" 201
    if [ $? -eq 0 ]; then
        USER_ID=$(jq -r '.id // empty' /tmp/api_response.json)
        log_message "INFO" "Test user created successfully. ID: ${USER_ID:-N/A}"
    else
        if grep -q "email already exists" "$RUN_LOG_FILE"; then
            log_message "INFO" "Test user already exists: $USER_EMAIL. Attempting to retrieve ID..."
            api_request "GET" "/api/users?email=$USER_EMAIL" "" "$ADMIN_TOKEN" 200
            if [ $? -eq 0 ]; then
                USER_ID=$(jq -r --arg email "$USER_EMAIL" '.[] | select(.email == $email) | .id // empty' /tmp/api_response.json | head -n 1)
                if [ -n "$USER_ID" ]; then
                    log_message "INFO" "Retrieved existing user ID: $USER_ID"
                else
                    log_message "ERROR" "Failed to find existing user ID for $USER_EMAIL."
                fi
            fi
        else
            log_message "ERROR" "Failed to create test user."
        fi
    fi

    # Login as test user
    if [ -n "$USER_ID" ]; then
        api_request "POST" "/api/auth/login" '{"email":"'"$USER_EMAIL"'","password":"'"$USER_PASSWORD"'"}' "" 200
        if [ $? -eq 0 ]; then
            USER_TOKEN=$(jq -r '.token' /tmp/api_response.json)
            if [ -n "$USER_TOKEN" ] && [ "$USER_TOKEN" != "null" ]; then
                log_message "INFO" "Test user login successful."
            else
                log_message "ERROR" "Failed to extract user token."
            fi
        else
            log_message "ERROR" "Test user login failed."
        fi
    else
        log_message "WARN" "Skipping test user login (USER_ID not available)."
    fi
}

# 2. Test Expert Management (Admin)
test_experts() {
    print_header "TESTING EXPERT MANAGEMENT (Admin)"

    # --- Test 1: Create Expert (Happy Path) ---
    local expert_payload='{
        "name":"'"$EXPERT_NAME"'",
        "affiliation":"Test University '""$TIMESTAMP""'",
        "primaryContact":"'"$EXPERT_EMAIL"'",
        "contactType":"email",
        "designation":"Professor",
        "isBahraini":true,
        "availability":"yes",
        "rating":"5",
        "role":"evaluator",
        "employmentType":"academic",
        "generalArea":'"$GENERAL_AREA_ID"',
        "specializedArea":"Software Engineering",
        "isTrained":true,
        "isPublished":true,
        "biography":"An expert created during automated testing.",
        "skills":["Go","React","Testing"]
    }'
    api_request "POST" "/api/experts" "$expert_payload" "$ADMIN_TOKEN" 201
    if [ $? -eq 0 ]; then
        EXPERT_INTERNAL_ID=$(jq -r '.id // empty' /tmp/api_response.json)
        if [ -n "$EXPERT_INTERNAL_ID" ]; then
            api_request "GET" "/api/experts/$EXPERT_INTERNAL_ID" "" "$ADMIN_TOKEN" 200
            if [ $? -eq 0 ]; then
                EXPERT_ID=$(jq -r '.expertId // empty' /tmp/api_response.json)
                log_message "INFO" "Expert created successfully. Internal ID: $EXPERT_INTERNAL_ID, Expert ID: ${EXPERT_ID:-N/A}"
            else
                log_message "ERROR" "Failed to retrieve created expert details."
            fi
        else
            log_message "ERROR" "Expert creation succeeded but no internal ID returned."
        fi
    else
        if grep -q "expert ID already exists" "$RUN_LOG_FILE"; then
            log_message "INFO" "Expert ID conflict detected. Attempting to retrieve existing expert..."
            api_request "GET" "/api/experts?email=$EXPERT_EMAIL" "" "$ADMIN_TOKEN" 200
            if [ $? -eq 0 ]; then
                EXPERT_INTERNAL_ID=$(jq -r '.[0].id // empty' /tmp/api_response.json)
                EXPERT_ID=$(jq -r '.[0].expertId // empty' /tmp/api_response.json)
                if [ -n "$EXPERT_INTERNAL_ID" ]; then
                    log_message "INFO" "Retrieved existing expert. Internal ID: $EXPERT_INTERNAL_ID, Expert ID: $EXPERT_ID"
                else
                    log_message "ERROR" "Failed to retrieve existing expert."
                fi
            fi
        elif grep -q "email already exists" "$RUN_LOG_FILE"; then
            log_message "ERROR" "Failed to create expert - Email conflict: $EXPERT_EMAIL"
        else
            log_message "ERROR" "Failed to create expert."
        fi
    fi

    # --- Test 2: Create Expert without ExpertID (Backend Generation) ---
    local no_id_payload='{
        "name":"NoID Expert '"$TIMESTAMP"'",
        "affiliation":"Test University '""$TIMESTAMP""'",
        "primaryContact":"noid${TIMESTAMP}@example.com",
        "contactType":"email",
        "generalArea":'"$GENERAL_AREA_ID"'"
    }'
    api_request "POST" "/api/experts" "$no_id_payload" "$ADMIN_TOKEN" 201
    if [ $? -eq 0 ]; then
        local noid_expert_id=$(jq -r '.id // empty' /tmp/api_response.json)
        if [ -n "$noid_expert_id" ]; then
            api_request "GET" "/api/experts/$noid_expert_id" "" "$ADMIN_TOKEN" 200
            if [ $? -eq 0 ]; then
                local generated_expert_id=$(jq -r '.expertId // empty' /tmp/api_response.json)
                if [[ "$generated_expert_id" =~ ^EXP- ]]; then
                    log_message "INFO" "Successfully created expert with backend-generated ID: $generated_expert_id"
                else
                    log_message "ERROR" "Backend-generated expert ID is invalid: $generated_expert_id"
                    record_test_result "POST /api/experts (No ExpertID)" "failed" "Invalid generated expert ID"
                fi
            fi
        else
            log_message "ERROR" "Expert creation (no ID) succeeded but no internal ID returned."
            record_test_result "POST /api/experts (No ExpertID)" "failed" "No internal ID returned"
        fi
    else
        log_message "ERROR" "Failed to create expert without expertId."
    fi

    # --- Test 3: Create Expert with Missing Required Fields ---
    local invalid_payload='{
        "affiliation":"Test University",
        "generalArea":'"$GENERAL_AREA_ID"'"
    }' # Missing name, primaryContact
    api_request "POST" "/api/experts" "$invalid_payload" "$ADMIN_TOKEN" 400
    if [ $? -eq 0 ]; then
        log_message "INFO" "Correctly rejected expert creation with missing required fields."
    else
        log_message "ERROR" "Failed to reject invalid expert payload."
    fi

    # --- Test 4: Create Expert with Invalid General Area ---
    local invalid_area_payload='{
        "name":"Invalid Area Expert '"$TIMESTAMP"'",
        "affiliation":"Test University",
        "primaryContact":"invalidarea${TIMESTAMP}@example.com",
        "contactType":"email",
        "generalArea":-1
    }'
    api_request "POST" "/api/experts" "$invalid_area_payload" "$ADMIN_TOKEN" 400
    if [ $? -eq 0 ]; then
        log_message "INFO" "Correctly rejected expert creation with invalid generalArea."
    else
        log_message "ERROR" "Failed to reject invalid generalArea."
    fi

    # --- Test 5: Create Expert with Invalid Email ---
    local invalid_email_payload='{
        "name":"Invalid Email Expert '"$TIMESTAMP"'",
        "affiliation":"Test University",
        "primaryContact":"invalid-email",
        "contactType":"email",
        "generalArea":'"$GENERAL_AREA_ID"'"
    }'
    api_request "POST" "/api/experts" "$invalid_email_payload" "$ADMIN_TOKEN" 400
    if [ $? -eq 0 ]; then
        log_message "INFO" "Correctly rejected expert creation with invalid email."
    else
        log_message "ERROR" "Failed to reject invalid email."
    fi

    # --- Test 6: List Experts ---
    api_request "GET" "/api/experts?limit=5" "" "$ADMIN_TOKEN" 200
    if [ $? -eq 0 ]; then
        local experts_count=$(jq '. | length' /tmp/api_response.json)
        log_message "INFO" "Listed $experts_count experts."
    else
        if grep -q "failed to scan expert row" "$RUN_LOG_FILE"; then
            log_message "ERROR" "Failed to list experts - Possible nullable field issue."
        else
            log_message "ERROR" "Failed to list experts."
        fi
    fi

    # --- Test 7: Get Expert Details ---
    if [ -n "$EXPERT_INTERNAL_ID" ]; then
        api_request "GET" "/api/experts/$EXPERT_INTERNAL_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Get Expert Details (EXPERT_INTERNAL_ID not available)."
        record_test_result "GET /api/experts/{id}" "skipped" "EXPERT_INTERNAL_ID not available"
    fi

    # --- Test 8: Update Expert ---
    if [ -n "$EXPERT_INTERNAL_ID" ]; then
        local update_payload='{"designation":"Senior Professor","availability":"no"}'
        api_request "PUT" "/api/experts/$EXPERT_INTERNAL_ID" "$update_payload" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Update Expert (EXPERT_INTERNAL_ID not available)."
        record_test_result "PUT /api/experts/{id}" "skipped" "EXPERT_INTERNAL_ID not available"
    fi
}

# 3. Test Expert Request Management
test_expert_requests() {
    print_header "TESTING EXPERT REQUEST MANAGEMENT"

    local request_payload='{
        "name":"'"$REQUEST_NAME"'",
        "designation":"Researcher",
        "institution":"Request University '""$TIMESTAMP""'",
        "isBahraini":false,
        "isAvailable":true,
        "rating":"4",
        "role":"reviewer",
        "employmentType":"freelance",
        "generalArea":'"$GENERAL_AREA_ID"',
        "specializedArea":"Quantum Physics",
        "isTrained":false,
        "phone":"'"$EXPERT_PHONE"'",
        "email":"'"$REQUEST_EMAIL"'",
        "isPublished":false,
        "biography":"A researcher requesting addition."
    }'

    # Create request
    local creator_token="$USER_TOKEN"
    if [ -z "$creator_token" ]; then
        log_message "WARN" "User token not available, creating expert request as admin."
        creator_token="$ADMIN_TOKEN"
    fi

    api_request "POST" "/api/expert-requests" "$request_payload" "$creator_token" 201
    if [ $? -eq 0 ]; then
        EXPERT_REQUEST_ID=$(jq -r '.id // empty' /tmp/api_response.json)
        log_message "INFO" "Expert request created successfully. ID: ${EXPERT_REQUEST_ID:-N/A}"
    else
        log_message "ERROR" "Failed to create expert request."
    fi

    # List requests
    api_request "GET" "/api/expert-requests?limit=5" "" "$ADMIN_TOKEN" 200

    # Get request details
    if [ -n "$EXPERT_REQUEST_ID" ]; then
        api_request "GET" "/api/expert-requests/$EXPERT_REQUEST_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Get Expert Request Details (EXPERT_REQUEST_ID not available)."
        record_test_result "GET /api/expert-requests/{id}" "skipped" "EXPERT_REQUEST_ID not available"
    fi

    # Approve request
    if [ -n "$EXPERT_REQUEST_ID" ]; then
        local approve_payload='{"status":"approved"}'
        api_request "PUT" "/api/expert-requests/$EXPERT_REQUEST_ID" "$approve_payload" "$ADMIN_TOKEN" 200
        if [ $? -ne 0 ]; then
            if grep -q "expert ID already exists" "$RUN_LOG_FILE"; then
                log_message "ERROR" "Failed to approve expert request - Expert ID conflict."
            else
                log_message "ERROR" "Failed to approve expert request."
            fi
        fi
    else
        log_message "SKIP" "Skipping Approve Expert Request (EXPERT_REQUEST_ID not available)."
        record_test_result "PUT /api/expert-requests/{id} (Approve)" "skipped" "EXPERT_REQUEST_ID not available"
    fi

    # Reject request (create a new one to test rejection)
    local reject_request_payload='{
        "name":"Reject '"$REQUEST_NAME"'",
        "designation":"Researcher",
        "institution":"Reject University '""$TIMESTAMP""'",
        "email":"reject${TIMESTAMP}@example.com",
        "generalArea":'"$GENERAL_AREA_ID"'"
    }'
    api_request "POST" "/api/expert-requests" "$reject_request_payload" "$creator_token" 201
    if [ $? -eq 0 ]; then
        local reject_request_id=$(jq -r '.id // empty' /tmp/api_response.json)
        if [ -n "$reject_request_id" ]; then
            local reject_payload='{"status":"rejected","rejectionReason":"Test rejection"}'
            api_request "PUT" "/api/expert-requests/$reject_request_id" "$reject_payload" "$ADMIN_TOKEN" 200
            if [ $? -eq 0 ]; then
                log_message "INFO" "Successfully rejected expert request ID: $reject_request_id"
            else
                log_message "ERROR" "Failed to reject expert request."
            fi
        fi
    else
        log_message "ERROR" "Failed to create expert request for rejection test."
    fi
}

# 4. Test Document Management
test_documents() {
    print_header "TESTING DOCUMENT MANAGEMENT"

    if [ -z "$EXPERT_INTERNAL_ID" ]; then
        log_message "SKIP" "Skipping Document Management tests (EXPERT_INTERNAL_ID not available)."
        record_test_result "Document Management" "skipped" "EXPERT_INTERNAL_ID not available"
        return
    fi

    # Create a sample document file
    local doc_file="/tmp/sample_cv_${TIMESTAMP}.txt"
    echo "Sample CV for expert $EXPERT_NAME ($EXPERT_ID)" > "$doc_file"
    log_message "INFO" "Created sample document file: $doc_file"

    # Upload document
    log_message "STEP" "Executing: POST /api/documents (Upload)"
    local upload_output_file="/tmp/upload_response.json"
    local upload_status=$(curl -s -w '%{http_code}' -o "$upload_output_file" \
        -X POST \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -F "file=@${doc_file}" \
        -F "documentType=cv" \
        -F "expertId=${EXPERT_INTERNAL_ID}" \
        "${BASE_URL}/api/documents")
    local curl_exit_code=$?

    # Log upload details
    log_message "DETAIL" "Upload Command: curl ... -F file=@${doc_file} -F documentType=cv -F expertId=${EXPERT_INTERNAL_ID} ${BASE_URL}/api/documents"
    echo "--- Upload Response ---" >> "$RUN_LOG_FILE"
    echo "HTTP Status: $upload_status" >> "$RUN_LOG_FILE"
    echo "Response Body:" >> "$RUN_LOG_FILE"
    jq '.' "$upload_output_file" >> "$RUN_LOG_FILE" 2>/dev/null || cat "$upload_output_file" >> "$RUN_LOG_FILE"
    echo "--- End Upload Response ---" >> "$RUN_LOG_FILE"

    if [ $curl_exit_code -ne 0 ] || [ "$upload_status" -ne 201 ]; then
        log_message "ERROR" "Failed to upload document. Status: $upload_status, Curl Exit: $curl_exit_code"
        jq '.' "$upload_output_file"
        record_test_result "POST /api/documents" "failed" "Status $upload_status"
    else
        DOCUMENT_ID=$(jq -r '.id // empty' "$upload_output_file")
        log_message "SUCCESS" "Document uploaded successfully. ID: ${DOCUMENT_ID:-N/A}"
        record_test_result "POST /api/documents" "passed" "Status $upload_status"
    fi

    # Clean up sample file
    rm -f "$doc_file"

    # List documents
    api_request "GET" "/api/experts/$EXPERT_INTERNAL_ID/documents" "" "$ADMIN_TOKEN" 200

    # Get document details
    if [ -n "$DOCUMENT_ID" ]; then
        api_request "GET" "/api/documents/$DOCUMENT_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Get Document Details (DOCUMENT_ID not available)."
        record_test_result "GET /api/documents/{id}" "skipped" "DOCUMENT_ID not available"
    fi

    # Delete document
    if [ -n "$DOCUMENT_ID" ]; then
        api_request "DELETE" "/api/documents/$DOCUMENT_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Delete Document (DOCUMENT_ID not available)."
        record_test_result "DELETE /api/documents/{id}" "skipped" "DOCUMENT_ID not available"
    fi
}

# 5. Test Statistics
test_statistics() {
    print_header "TESTING STATISTICS"

    api_request "GET" "/api/statistics" "" "$ADMIN_TOKEN" 200
    api_request "GET" "/api/statistics/growth?months=6" "" "$ADMIN_TOKEN" 200
    api_request "GET" "/api/statistics/nationality" "" "$ADMIN_TOKEN" 200
    api_request "GET" "/api/statistics/engagements" "" "$ADMIN_TOKEN" 200
}

# 6. Test Cleanup
test_cleanup() {
    print_header "TESTING CLEANUP"

    # Delete expert
    if [ -n "$EXPERT_INTERNAL_ID" ]; then
        api_request "DELETE" "/api/experts/$EXPERT_INTERNAL_ID" "" "$ADMIN_TOKEN" 200
        if [ $? -eq 0 ]; then
            # Verify cascade deletion of documents
            api_request "GET" "/api/experts/$EXPERT_INTERNAL_ID/documents" "" "$ADMIN_TOKEN" 200
            if [ $? -eq 0 ]; then
                local doc_count=$(jq '. | length' /tmp/api_response.json)
                if [ "$doc_count" -eq 0 ]; then
                    log_message "INFO" "Confirmed cascade deletion of expert documents."
                else
                    log_message "ERROR" "Documents not cascade-deleted for expert $EXPERT_INTERNAL_ID."
                    record_test_result "Cascade Deletion (Documents)" "failed" "Found $doc_count documents"
                fi
            fi
        else
            log_message "WARN" "Failed to cleanup expert ID $EXPERT_INTERNAL_ID."
        fi
    else
        log_message "SKIP" "Skipping Expert cleanup (EXPERT_INTERNAL_ID not available)."
    fi

    # Delete user
    if [ -n "$USER_ID" ]; then
        api_request "DELETE" "/api/users/$USER_ID" "" "$ADMIN_TOKEN" 200
        if [ $? -ne 0 ]; then
            log_message "WARN" "Failed to cleanup user ID $USER_ID."
        fi
    else
        log_message "SKIP" "Skipping User cleanup (USER_ID not available)."
    fi
}

# Function to print test summary
print_test_summary() {
    log_message "HEADER" "===== TEST SUMMARY ====="
    log_message "INFO" "Total tests run: ${test_stats["total"]}"
    log_message "SUCCESS" "Passed: ${test_stats["passed"]}"
    log_message "ERROR" "Failed: ${test_stats["failed"]}"
    log_message "SKIP" "Skipped: ${test_stats["skipped"]}"
    echo ""

    if [ ${test_stats["failed"]} -gt 0 ]; then
        log_message "ERROR" "Failed Tests:"
        for test_name in "${!test_details[@]}"; do
            if [[ "${test_details[$test_name]}" == failed* ]]; then
                log_message "ERROR" "- $test_name: ${test_details[$test_name]}"
            fi
        done
        echo ""
    fi

    if [ ${test_stats["skipped"]} -gt 0 ]; then
        log_message "SKIP" "Skipped Tests:"
        for test_name in "${!test_details[@]}"; do
            if [[ "${test_details[$test_name]}" == skipped* ]]; then
                log_message "SKIP" "- $test_name: ${test_details[$test_name]}"
            fi
        done
        echo ""
    fi

    log_message "INFO" "Detailed logs available in: $RUN_LOG_FILE"
}

# --- Main Execution ---
main() {
    # Create log directory
    mkdir -p "$LOG_DIR"
    echo "===== ExpertDB API Test Run Started: $(date) =====" > "$RUN_LOG_FILE"

    log_message "INFO" "Starting ExpertDB API tests against $BASE_URL..."
    log_message "INFO" "Detailed logs will be saved to $RUN_LOG_FILE"

    test_auth
    test_setup
    test_experts
    test_expert_requests
    test_documents
    test_statistics
    test_cleanup

    print_test_summary

    if [ ${test_stats["failed"]} -gt 0 ]; then
        log_message "ERROR" "===== Test Run Finished with Failures ====="
        exit 1
    else
        log_message "SUCCESS" "===== Test Run Finished Successfully ====="
        exit 0
    fi
}

main