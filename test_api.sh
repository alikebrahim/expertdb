#!/bin/bash

# ==============================================================================
# ExpertDB API Test Script (Enhanced Version)
#
# Description:
#   Tests the core API endpoints of the ExpertDB application.
#   Provides detailed logging of requests and responses.
#   Uses unique data for creation tests to improve re-runnability.
#   Skips dependent tests if prerequisites fail.
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
EXPERT_PHONE="+9731111${TIMESTAMP: -4}" # Keep generating phone, might be useful later
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
EXPERT_INTERNAL_ID="" # The auto-incrementing ID
EXPERT_REQUEST_ID=""
DOCUMENT_ID=""

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

    # Log the command (with token potentially redacted if needed for security)
    log_message "DETAIL" "URL: $url"
    log_message "DETAIL" "Command: $curl_cmd" # Be cautious logging full token in production logs
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
    echo "" >> "$RUN_LOG_FILE" # Add a newline for readability
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

    # Check for API-level error in JSON response (if response is JSON)
    if jq -e . "$output_file" > /dev/null 2>&1; then # Check if valid JSON
        if jq -e '.error' "$output_file" > /dev/null 2>&1; then # Check if .error field exists
            error_msg=$(jq -r '.error' "$output_file")
            log_message "ERROR" "API returned error: $error_msg"
            record_test_result "$test_name" "failed" "API Error: $error_msg"
            return 1
        fi
    elif [ -s "$output_file" ]; then # If not JSON but response file has content
        log_message "WARN" "Response is not valid JSON."
        # Decide if this is an error based on context or expected status
        if [ "$http_status" -lt 300 ]; then # Treat non-JSON as warning for success statuses
             log_message "DETAIL" "Received non-JSON response with success status $http_status."
        else # Treat non-JSON as error for error statuses
            log_message "ERROR" "Received non-JSON response with error status $http_status."
            record_test_result "$test_name" "failed" "Non-JSON response with status $http_status"
            return 1
        fi
    fi

    # If we reached here, the request is considered successful for this step
    log_message "SUCCESS" "$test_name completed successfully (Status: $http_status)."
    record_test_result "$test_name" "passed" "Status $http_status"
    return 0
}


# --- Test Sections ---

# 1. Test Authentication
test_auth() {
    print_header "TESTING AUTHENTICATION"

    # Login as admin
    api_request "POST" "/api/auth/login" '{"email":"'"$ADMIN_EMAIL"'","password":"'"$ADMIN_PASSWORD"'"}' "" 200
    if [ $? -ne 0 ]; then
        log_message "ERROR" "Admin login failed. Cannot proceed with tests."
        exit 1
    fi
    ADMIN_TOKEN=$(jq -r '.token' /tmp/api_response.json)
    if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
        log_message "ERROR" "Failed to extract admin token. Cannot proceed."
        exit 1
    fi
    log_message "INFO" "Admin login successful."

    # Create a test user
    local user_payload='{"name":"'"$USER_NAME"'","email":"'"$USER_EMAIL"'","password":"'"$USER_PASSWORD"'","role":"user","isActive":true}'
    api_request "POST" "/api/users" "$user_payload" "$ADMIN_TOKEN" 201
    local create_user_status=$?
    if [ $create_user_status -eq 0 ]; then
        USER_ID=$(jq -r '.id // empty' /tmp/api_response.json) # Use default if ID is missing
        log_message "INFO" "Test user created successfully. ID: ${USER_ID:-N/A}"
    else
        # Check if it failed because the user already exists (common on reruns)
        if grep -q "email already exists" "$RUN_LOG_FILE"; then
             log_message "WARN" "Test user creation failed (likely already exists): $USER_EMAIL"
             # Try to retrieve the existing user's ID
             # Note: This assumes the GET /api/users endpoint supports filtering by email.
             # If not, this part will fail or return all users.
             api_request "GET" "/api/users?email=$USER_EMAIL" "" "$ADMIN_TOKEN" 200
             if [ $? -eq 0 ]; then
                 # Attempt to find the user in the list
                 USER_ID=$(jq -r --arg email "$USER_EMAIL" '.[] | select(.email == $email) | .id // empty' /tmp/api_response.json | head -n 1)
                 if [ -n "$USER_ID" ]; then
                    log_message "INFO" "Retrieved existing user ID: $USER_ID"
                 else
                    log_message "ERROR" "Failed to find existing user ID for $USER_EMAIL in the list."
                 fi
             else
                 log_message "ERROR" "Failed to list users to find existing user ID."
             fi
        else
             log_message "ERROR" "Failed to create test user for reasons other than duplication."
             # USER_ID remains empty
        fi
    fi


    # Login as test user
    api_request "POST" "/api/auth/login" '{"email":"'"$USER_EMAIL"'","password":"'"$USER_PASSWORD"'"}' "" 200
    if [ $? -ne 0 ]; then
        log_message "ERROR" "Test user login failed."
        # Decide if this is critical - maybe continue with admin token?
        # For now, we'll continue but USER_TOKEN will be empty
    else
        USER_TOKEN=$(jq -r '.token' /tmp/api_response.json)
        if [ -z "$USER_TOKEN" ] || [ "$USER_TOKEN" = "null" ]; then
            log_message "ERROR" "Failed to extract user token."
        else
            log_message "INFO" "Test user login successful."
        fi
    fi
}

# 2. Test User Management (Admin)
test_users() {
    print_header "TESTING USER MANAGEMENT (Admin)"

    # List all users
    api_request "GET" "/api/users" "" "$ADMIN_TOKEN" 200
    if [ $? -eq 0 ]; then
        local users_count=$(jq '. | length' /tmp/api_response.json)
        log_message "INFO" "Listed $users_count users."
    fi

    # Get details for the created test user
    if [ -n "$USER_ID" ]; then
        api_request "GET" "/api/users/$USER_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Get User Details (USER_ID not available)."
        record_test_result "GET /api/users/{id}" "skipped" "USER_ID not available"
    fi

    # Update the created test user
    if [ -n "$USER_ID" ]; then
        local update_payload='{"name":"Updated '"$USER_NAME"'","role":"user","isActive":false}' # Keep email same
        api_request "PUT" "/api/users/$USER_ID" "$update_payload" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Update User (USER_ID not available)."
        record_test_result "PUT /api/users/{id}" "skipped" "USER_ID not available"
    fi
}

# 3. Test Expert Management (Admin)
test_experts() {
    print_header "TESTING EXPERT MANAGEMENT (Admin)"

    # --- Corrected Payload ---
    # Create a new expert payload including primaryContact and contactType
    # Use integer ID for generalArea based on migration 0005
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
      "generalArea": 1,
      "specializedArea":"Software Engineering",
      "isTrained":true,
      "isPublished":true,
      "biography":"An expert created during automated testing.",
      "skills": ["Go", "React", "Testing"]
    }'
    # Note: The backend CreateExpertRequest uses 'affiliation', 'availability', 'skills'
    # The backend Expert struct uses 'institution', 'isAvailable'
    # The NewExpert function maps these. Ensure the payload matches CreateExpertRequest.

    api_request "POST" "/api/experts" "$expert_payload" "$ADMIN_TOKEN" 201
    if [ $? -eq 0 ]; then
        EXPERT_INTERNAL_ID=$(jq -r '.id // empty' /tmp/api_response.json)
        # We need the *actual* expert record to get the business ID if generated by backend
        if [ -n "$EXPERT_INTERNAL_ID" ]; then
             api_request "GET" "/api/experts/$EXPERT_INTERNAL_ID" "" "$ADMIN_TOKEN" 200
             if [ $? -eq 0 ]; then
                  EXPERT_ID=$(jq -r '.expertId // empty' /tmp/api_response.json)
                  log_message "INFO" "Expert created successfully. Internal ID: $EXPERT_INTERNAL_ID, Expert ID: ${EXPERT_ID:-N/A}"
             else
                 log_message "WARN" "Expert created (Internal ID: $EXPERT_INTERNAL_ID), but failed to retrieve details to get Expert ID."
             fi
        else
             log_message "ERROR" "Expert creation reported success, but no internal ID found in response."
        fi
    else
        # Check for known unique constraint error
        if grep -q "UNIQUE constraint failed: experts.expert_id" "$RUN_LOG_FILE" || grep -q "expert ID already exists" "$RUN_LOG_FILE"; then
            log_message "ERROR" "Failed to create expert - KNOWN BACKEND ISSUE (Unique constraint on expert_id)."
        elif grep -q "email already exists" "$RUN_LOG_FILE"; then
             log_message "ERROR" "Failed to create expert - Email conflict: $EXPERT_EMAIL"
        # Check if the error is the one we just fixed
        elif grep -q "primary contact is required" "$RUN_LOG_FILE"; then
             log_message "ERROR" "Failed to create expert - Payload issue persists (primary contact required)."
        else
            log_message "ERROR" "Failed to create expert for other reasons."
        fi
        # EXPERT_ID and EXPERT_INTERNAL_ID remain empty
    fi

    # List all experts
    api_request "GET" "/api/experts?limit=5" "" "$ADMIN_TOKEN" 200
    if [ $? -ne 0 ]; then
         # Check for known scan error
         if grep -q "converting NULL to string is unsupported" "$RUN_LOG_FILE" || grep -q "failed to scan expert row" "$RUN_LOG_FILE"; then
              log_message "ERROR" "Failed to list experts - KNOWN BACKEND ISSUE (Scan error on nullable field)."
         else
              log_message "ERROR" "Failed to list experts for other reasons."
         fi
    fi

    # Get expert details
    if [ -n "$EXPERT_INTERNAL_ID" ]; then
        api_request "GET" "/api/experts/$EXPERT_INTERNAL_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Get Expert Details (EXPERT_INTERNAL_ID not available)."
        record_test_result "GET /api/experts/{id}" "skipped" "EXPERT_INTERNAL_ID not available"
    fi

    # Update expert
    if [ -n "$EXPERT_INTERNAL_ID" ]; then
        local update_payload='{"designation":"Senior Professor","isAvailable":false}'
        api_request "PUT" "/api/experts/$EXPERT_INTERNAL_ID" "$update_payload" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Update Expert (EXPERT_INTERNAL_ID not available)."
        record_test_result "PUT /api/experts/{id}" "skipped" "EXPERT_INTERNAL_ID not available"
    fi
}

# 4. Test Expert Request Management (User creates, Admin manages)
test_expert_requests() {
    print_header "TESTING EXPERT REQUEST MANAGEMENT"

    # Note: Expert Request payload seems okay based on previous logs, uses email/phone directly
    local request_payload='{
      "name":"'"$REQUEST_NAME"'",
      "designation":"Researcher",
      "institution":"Request University '""$TIMESTAMP""'",
      "isBahraini":false,
      "isAvailable":true,
      "rating":"4",
      "role":"reviewer",
      "employmentType":"freelance",
      "generalArea": 2,
      "specializedArea":"Quantum Physics",
      "isTrained":false,
      "phone":"+9732222${TIMESTAMP: -4}",
      "email":"'"$REQUEST_EMAIL"'",
      "isPublished":false,
      "biography":"A researcher requesting addition to the database."
    }'

    # Create request (as regular user, if token available)
    local creator_token="$USER_TOKEN"
    if [ -z "$creator_token" ]; then
        log_message "WARN" "User token not available, creating expert request as admin instead."
        creator_token="$ADMIN_TOKEN"
    fi

    api_request "POST" "/api/expert-requests" "$request_payload" "$creator_token" 201
    if [ $? -eq 0 ]; then
        EXPERT_REQUEST_ID=$(jq -r '.id // empty' /tmp/api_response.json)
        log_message "INFO" "Expert request created successfully. ID: ${EXPERT_REQUEST_ID:-N/A}"
    else
        log_message "ERROR" "Failed to create expert request."
        # EXPERT_REQUEST_ID remains empty
    fi

    # List requests (as admin)
    api_request "GET" "/api/expert-requests?limit=5" "" "$ADMIN_TOKEN" 200

    # Get request details (as admin)
    if [ -n "$EXPERT_REQUEST_ID" ]; then
        api_request "GET" "/api/expert-requests/$EXPERT_REQUEST_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Get Expert Request Details (EXPERT_REQUEST_ID not available)."
        record_test_result "GET /api/expert-requests/{id}" "skipped" "EXPERT_REQUEST_ID not available"
    fi

    # Approve request (as admin)
    if [ -n "$EXPERT_REQUEST_ID" ]; then
        # The backend automatically populates reviewedBy from the token (implicitly admin ID 1 here)
        # The backend also generates the expert record upon approval.
        local approve_payload='{"status":"approved"}' # Minimal payload, backend fills details
        api_request "PUT" "/api/expert-requests/$EXPERT_REQUEST_ID" "$approve_payload" "$ADMIN_TOKEN" 200
        if [ $? -ne 0 ]; then
             # Check for potential conflict if the expert ID generated already exists (related to expert creation bug)
             if grep -q "expert ID already exists" "$RUN_LOG_FILE"; then
                  log_message "ERROR" "Failed to approve expert request - KNOWN BACKEND ISSUE (Expert ID conflict during approval)."
             else
                  log_message "ERROR" "Failed to approve expert request."
             fi
        fi
    else
        log_message "SKIP" "Skipping Approve Expert Request (EXPERT_REQUEST_ID not available)."
        record_test_result "PUT /api/expert-requests/{id} (Approve)" "skipped" "EXPERT_REQUEST_ID not available"
    fi

    # TODO: Add tests for rejecting requests
}

# 5. Test Document Management (Admin)
test_documents() {
    print_header "TESTING DOCUMENT MANAGEMENT (Admin)"

    if [ -z "$EXPERT_INTERNAL_ID" ]; then
        log_message "SKIP" "Skipping Document Management tests (EXPERT_INTERNAL_ID not available)."
        record_test_result "Document Management" "skipped" "EXPERT_INTERNAL_ID not available"
        return
    fi

    # Create a sample document file
    local doc_file="/tmp/sample_cv_${TIMESTAMP}.txt"
    echo "Sample CV content for expert $EXPERT_NAME ($EXPERT_ID)" > "$doc_file"
    log_message "INFO" "Created sample document file: $doc_file"

    # Upload document using multipart/form-data
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
    echo "" >> "$RUN_LOG_FILE"
    echo "--- End Upload Response ---" >> "$RUN_LOG_FILE"

    if [ $curl_exit_code -ne 0 ] || [ "$upload_status" -ne 201 ]; then
        log_message "ERROR" "Failed to upload document. Status: $upload_status, Curl Exit: $curl_exit_code"
        jq '.' "$upload_output_file" # Print response to console
        record_test_result "POST /api/documents" "failed" "Status $upload_status"
    else
        DOCUMENT_ID=$(jq -r '.id // empty' "$upload_output_file")
        log_message "SUCCESS" "Document uploaded successfully. ID: ${DOCUMENT_ID:-N/A}"
        record_test_result "POST /api/documents" "passed" "Status $upload_status"
    fi

    # Clean up sample file
    rm -f "$doc_file"

    # List documents for the expert
    api_request "GET" "/api/experts/$EXPERT_INTERNAL_ID/documents" "" "$ADMIN_TOKEN" 200

    # Get the uploaded document details
    if [ -n "$DOCUMENT_ID" ]; then
        api_request "GET" "/api/documents/$DOCUMENT_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Get Document Details (DOCUMENT_ID not available)."
        record_test_result "GET /api/documents/{id}" "skipped" "DOCUMENT_ID not available"
    fi

    # Delete the uploaded document
    if [ -n "$DOCUMENT_ID" ]; then
        api_request "DELETE" "/api/documents/$DOCUMENT_ID" "" "$ADMIN_TOKEN" 200
    else
        log_message "SKIP" "Skipping Delete Document (DOCUMENT_ID not available)."
        record_test_result "DELETE /api/documents/{id}" "skipped" "DOCUMENT_ID not available"
    fi
}

# 6. Test Expert Areas (Authenticated)
test_expert_areas() {
    print_header "TESTING EXPERT AREAS"
    api_request "GET" "/api/expert/areas" "" "$ADMIN_TOKEN" 200
    if [ $? -eq 0 ]; then
        local areas_count=$(jq '. | length' /tmp/api_response.json)
        log_message "INFO" "Retrieved $areas_count expert areas."
    fi
}

# 7. Test Statistics (Authenticated)
test_statistics() {
    print_header "TESTING STATISTICS"

    # Get overall statistics
    api_request "GET" "/api/statistics" "" "$ADMIN_TOKEN" 200

    # Get filtered statistics (example: growth over last 6 months)
    api_request "GET" "/api/statistics/growth?months=6" "" "$ADMIN_TOKEN" 200

    # Get nationality statistics
    api_request "GET" "/api/statistics/nationality" "" "$ADMIN_TOKEN" 200

    # Get engagement statistics
    api_request "GET" "/api/statistics/engagements" "" "$ADMIN_TOKEN" 200
}

# 8. Test Cleanup (Optional - Best Effort)
test_cleanup() {
    print_header "TESTING CLEANUP (Attempting)"

    # Delete the created expert
    if [ -n "$EXPERT_INTERNAL_ID" ]; then
        api_request "DELETE" "/api/experts/$EXPERT_INTERNAL_ID" "" "$ADMIN_TOKEN" 200 || log_message "WARN" "Failed to cleanup expert ID $EXPERT_INTERNAL_ID (might have been deleted or failed previously)."
    else
        log_message "SKIP" "Skipping Expert cleanup (EXPERT_INTERNAL_ID not available)."
    fi

    # Delete the created user
    if [ -n "$USER_ID" ]; then
        api_request "DELETE" "/api/users/$USER_ID" "" "$ADMIN_TOKEN" 200 || log_message "WARN" "Failed to cleanup user ID $USER_ID (might have been deleted or failed previously)."
    else
        log_message "SKIP" "Skipping User cleanup (USER_ID not available)."
    fi

    # Note: Expert Requests are not deleted via API in this version.
    # Note: Documents associated with the expert should be cascade deleted by the backend if implemented correctly.
}

# Function to print test summary
print_test_summary() {
    log_message "HEADER" "===== TEST SUMMARY ====="
    log_message "INFO" "Total tests run: ${test_stats["total"]}"
    log_message "SUCCESS" "Passed: ${test_stats["passed"]}"
    log_message "ERROR" "Failed: ${test_stats["failed"]}"
    log_message "SKIP" "Skipped: ${test_stats["skipped"]}"
    echo "" # Add a newline

    if [ ${test_stats["failed"]} -gt 0 ]; then
        log_message "ERROR" "Failed Tests:"
        for test_name in "${!test_details[@]}"; do
            if [[ "${test_details[$test_name]}" == failed* ]]; then
                log_message "ERROR" "- $test_name: ${test_details[$test_name]}"
            fi
        done
        echo "" # Add a newline
    fi

    if [ ${test_stats["skipped"]} -gt 0 ]; then
        log_message "SKIP" "Skipped Tests:"
        for test_name in "${!test_details[@]}"; do
            if [[ "${test_details[$test_name]}" == skipped* ]]; then
                log_message "SKIP" "- $test_name: ${test_details[$test_name]}"
            fi
        done
         echo "" # Add a newline
    fi

    log_message "INFO" "Detailed logs available in: $RUN_LOG_FILE"
}

# --- Main Execution ---
main() {
    # Create log directory
    mkdir -p "$LOG_DIR"
    # Initialize run log file
    echo "===== ExpertDB API Test Run Started: $(date) =====" > "$RUN_LOG_FILE"

    log_message "INFO" "Starting ExpertDB API tests against $BASE_URL..."
    log_message "INFO" "Detailed logs will be saved to $RUN_LOG_FILE"

    # Run test sections sequentially
    test_auth
    test_users
    test_experts
    test_expert_requests
    test_documents
    test_expert_areas
    test_statistics
    # test_cleanup # Optional: uncomment to attempt cleanup

    # Print summary
    print_test_summary

    # Exit with non-zero code if any tests failed
    if [ ${test_stats["failed"]} -gt 0 ]; then
        log_message "ERROR" "===== Test Run Finished with Failures ====="
        exit 1
    else
        log_message "SUCCESS" "===== Test Run Finished Successfully ====="
        exit 0
    fi
}

# Execute main function
main
