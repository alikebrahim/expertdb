#!/bin/bash

# test_api.sh - ExpertDB API Test Script
# This script tests the ExpertDB API endpoints for authentication, expert management,
# expert requests, and statistics. It addresses malformed payload issues by ensuring valid
# JSON and replacing placeholders. Assumes a small user base (10-12 users) and modest
# database (up to 1200 entries).

# Configuration
API_BASE_URL="http://localhost:8080"
LOG_DIR="logs"
LOG_FILE="$LOG_DIR/api_test_run_$(date +%Y%m%d_%H%M%S).log"
ADMIN_CREDENTIALS='{"email":"admin@expertdb.com","password":"adminpassword"}'
TIMESTAMP="1744807094" # Fixed for test consistency

# Ensure log directory exists
mkdir -p "$LOG_DIR"

# Redirect all output to log file and stdout
exec > >(tee -a "$LOG_FILE") 2>&1

# Utility functions
log() {
    local level=$1
    shift
    echo "$(date '+%Y-%m-%d %H:%M:%S') [$level] $@"
}

validate_json() {
    local payload=$1
    echo "$payload" | jq . >/dev/null 2>&1
    if [[ $? -ne 0 ]]; then
        log ERROR "Invalid JSON payload: $payload"
        exit 1
    fi
}

execute_curl() {
    local method=$1
    local endpoint=$2
    local payload=$3
    local token=$4
    local expected_status=$5
    local output_file="/tmp/api_response.json"
    local header_file="/tmp/api_response_headers.txt"

    log STEP "Executing: $method $endpoint"
    log DETAIL "URL: $API_BASE_URL$endpoint"

    local curl_cmd="curl -s -w '%{http_code}' -o '$output_file' --dump-header '$header_file' -X $method"
    if [[ -n "$token" ]]; then
        curl_cmd="$curl_cmd -H 'Authorization: Bearer $token'"
    fi
    if [[ "$method" != "GET" && -n "$payload" ]]; then
        validate_json "$payload"
        curl_cmd="$curl_cmd -H 'Content-Type: application/json' -d '$payload'"
    fi
    curl_cmd="$curl_cmd '$API_BASE_URL$endpoint'"

    log DETAIL "Command: $curl_cmd"
    echo "--- Request ---"
    if [[ -n "$token" ]]; then
        echo "Authorization: Bearer ${token:0:10}..."
    fi
    if [[ -n "$payload" ]]; then
        echo "Request Body:"
        echo "$payload" | jq .
    else
        echo "Request Body: (None)"
    fi
    echo "--- End Request ---"

    http_status=$(eval "$curl_cmd")
    curl_exit=$?
    log DETAIL "Curl Exit Code: $curl_exit"
    log DETAIL "HTTP Status: $http_status"

    echo "--- Response ---"
    echo "HTTP Status: $http_status"
    echo "Response Headers:"
    cat "$header_file"
    echo -e "\nResponse Body:"
    jq . "$output_file" 2>/dev/null || cat "$output_file"
    echo -e "\n--- End Response ---"

    if [[ "$http_status" == "$expected_status" ]]; then
        log SUCCESS "$method $endpoint completed successfully (Status: $http_status)."
        return 0
    else
        log ERROR "Expected status $expected_status but received $http_status."
        return 1
    fi
}

# Test counters
total_tests=0
passed_tests=0
failed_tests=0
skipped_tests=0

# Test execution
log INFO "Starting ExpertDB API tests against $API_BASE_URL..."
log INFO "Detailed logs will be saved to $LOG_FILE"

echo -e "\n===== TESTING AUTHENTICATION ====="
((total_tests++))
log INFO "Attempting admin login (1/3)..."
if execute_curl POST "/api/auth/login" "$ADMIN_CREDENTIALS" "" 200; then
    admin_token=$(jq -r '.token' /tmp/api_response.json)
    log INFO "Admin login successful."
    ((passed_tests++))
else
    log ERROR "Admin login failed."
    ((failed_tests++))
    exit 1
fi

((total_tests++))
log STEP "Executing: POST /api/users"
user_payload='{
    "name": "Test User '$TIMESTAMP'",
    "email": "testuser'$TIMESTAMP'@example.com",
    "password": "password123",
    "role": "user",
    "isActive": true
}'
if execute_curl POST "/api/users" "$user_payload" "$admin_token" 201; then
    user_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Test user created successfully. ID: $user_id"
    ((passed_tests++))
else
    log ERROR "Failed to create test user."
    ((failed_tests++))
fi

((total_tests++))
log STEP "Executing: POST /api/auth/login"
user_credentials='{
    "email": "testuser'$TIMESTAMP'@example.com",
    "password": "password123"
}'
if execute_curl POST "/api/auth/login" "$user_credentials" "" 200; then
    user_token=$(jq -r '.token' /tmp/api_response.json)
    log INFO "Test user login successful."
    ((passed_tests++))
else
    log ERROR "Test user login failed."
    ((failed_tests++))
fi

echo -e "\n===== TESTING SETUP ====="
((total_tests++))
log STEP "Executing: GET /api/expert/areas"
if execute_curl GET "/api/expert/areas" "" "$admin_token" 200; then
    general_area_id=1 # Use first area (Art and Design)
    log INFO "Retrieved $(jq length /tmp/api_response.json) expert areas. Using generalArea ID: $general_area_id"
    ((passed_tests++))
else
    log ERROR "Failed to retrieve expert areas."
    ((failed_tests++))
fi

echo -e "\n===== TESTING EXPERT MANAGEMENT (Admin) ====="
((total_tests++))
log STEP "Executing: POST /api/experts"
expert_payload='{
    "name": "Test Expert '$TIMESTAMP'",
    "institution": "Test University '$TIMESTAMP'",
    "email": "expert'$TIMESTAMP'@example.com",
    "designation": "Professor",
    "isBahraini": true,
    "isAvailable": true,
    "rating": "5",
    "role": "evaluator",
    "employmentType": "academic",
    "generalArea": '$general_area_id',
    "specializedArea": "Software Engineering",
    "isTrained": true,
    "isPublished": true,
    "biography": "An expert created during automated testing."
}'
if execute_curl POST "/api/experts" "$expert_payload" "$admin_token" 201; then
    expert_internal_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Expert created successfully. ID: $expert_internal_id"
    ((passed_tests++))
else
    log ERROR "Failed to create expert."
    ((failed_tests++))
    expert_internal_id=""
fi

((total_tests++))
log STEP "Executing: POST /api/experts (No ID)"
no_id_expert_payload='{
    "name": "NoID Expert '$TIMESTAMP'",
    "institution": "Test University '$TIMESTAMP'",
    "email": "noid'$TIMESTAMP'@example.com",
    "generalArea": '$general_area_id'
}'
if execute_curl POST "/api/experts" "$no_id_expert_payload" "$admin_token" 201; then
    log INFO "Expert without expertId created successfully."
    ((passed_tests++))
else
    log ERROR "Failed to create expert without expertId."
    ((failed_tests++))
fi

((total_tests++))
log STEP "Executing: POST /api/experts (Invalid Payload)"
invalid_expert_payload='{
    "institution": "Test University",
    "generalArea": '$general_area_id'
}'
if execute_curl POST "/api/experts" "$invalid_expert_payload" "$admin_token" 400; then
    log INFO "Successfully rejected invalid expert payload."
    ((passed_tests++))
else
    log ERROR "Failed to reject invalid expert payload."
    ((failed_tests++))
fi

((total_tests++))
log STEP "Executing: POST /api/experts (Invalid General Area)"
invalid_area_expert_payload='{
    "name": "Invalid Area Expert '$TIMESTAMP'",
    "institution": "Test University",
    "email": "invalidarea'$TIMESTAMP'@example.com",
    "generalArea": -1
}'
if execute_curl POST "/api/experts" "$invalid_area_expert_payload" "$admin_token" 400; then
    log INFO "Successfully rejected invalid generalArea."
    ((passed_tests++))
else
    log ERROR "Failed to reject invalid generalArea."
    ((failed_tests++))
fi

((total_tests++))
log STEP "Executing: POST /api/experts (Invalid Email)"
invalid_email_expert_payload='{
    "name": "Invalid Email Expert '$TIMESTAMP'",
    "institution": "Test University",
    "email": "invalid-email",
    "generalArea": '$general_area_id'
}'
if execute_curl POST "/api/experts" "$invalid_email_expert_payload" "$admin_token" 400; then
    log INFO "Successfully rejected invalid email."
    ((passed_tests++))
else
    log ERROR "Failed to reject invalid email."
    ((failed_tests++))
fi

((total_tests++))
log STEP "Executing: GET /api/experts?limit=5"
if execute_curl GET "/api/experts?limit=5" "" "$admin_token" 200; then
    expert_count=$(jq length /tmp/api_response.json)
    log INFO "Listed $expert_count experts."
    ((passed_tests++))
else
    log ERROR "Failed to list experts."
    ((failed_tests++))
fi

if [[ -n "$expert_internal_id" ]]; then
    ((total_tests++))
    log STEP "Executing: GET /api/experts/$expert_internal_id"
    if execute_curl GET "/api/experts/$expert_internal_id" "" "$admin_token" 200; then
        log INFO "Retrieved expert details successfully."
        ((passed_tests++))
    else
        log ERROR "Failed to retrieve expert details."
        ((failed_tests++))
    fi

    ((total_tests++))
    log STEP "Executing: PUT /api/experts/$expert_internal_id"
    update_expert_payload='{
        "name": "Updated Expert '$TIMESTAMP'",
        "isAvailable": false
    }'
    if execute_curl PUT "/api/experts/$expert_internal_id" "$update_expert_payload" "$admin_token" 200; then
        log INFO "Expert updated successfully."
        ((passed_tests++))
    else
        log ERROR "Failed to update expert."
        ((failed_tests++))
    fi
else
    ((skipped_tests++))
    log SKIP "Skipping Get Expert Details (EXPERT_INTERNAL_ID not available)."
    ((skipped_tests++))
    log SKIP "Skipping Update Expert (EXPERT_INTERNAL_ID not available)."
fi

echo -e "\n===== TESTING EXPERT REQUEST MANAGEMENT ====="
((total_tests++))
log STEP "Executing: POST /api/expert-requests"
expert_request_payload='{
    "name": "Request Expert '$TIMESTAMP'",
    "designation": "Researcher",
    "institution": "Request University '$TIMESTAMP'",
    "isBahraini": false,
    "isAvailable": true,
    "rating": "4",
    "role": "reviewer",
    "employmentType": "freelance",
    "generalArea": '$general_area_id',
    "specializedArea": "Quantum Physics",
    "isTrained": false,
    "phone": "+9731111'$TIMESTAMP'",
    "email": "request'$TIMESTAMP'@example.com",
    "isPublished": false,
    "biography": "A researcher requesting addition."
}'
if execute_curl POST "/api/expert-requests" "$expert_request_payload" "$user_token" 201; then
    expert_request_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Expert request created successfully. ID: $expert_request_id"
    ((passed_tests++))
else
    log ERROR "Failed to create expert request."
    ((failed_tests++))
    expert_request_id=""
fi

((total_tests++))
log STEP "Executing: GET /api/expert-requests?limit=5"
if execute_curl GET "/api/expert-requests?limit=5" "" "$admin_token" 200; then
    log INFO "Listed expert requests successfully."
    ((passed_tests++))
else
    log ERROR "Failed to list expert requests."
    ((failed_tests++))
fi

if [[ -n "$expert_request_id" ]]; then
    ((total_tests++))
    log STEP "Executing: GET /api/expert-requests/$expert_request_id"
    if execute_curl GET "/api/expert-requests/$expert_request_id" "" "$admin_token" 200; then
        log INFO "Retrieved expert request details successfully."
        ((passed_tests++))
    else
        log ERROR "Failed to retrieve expert request details."
        ((failed_tests++))
    fi

    ((total_tests++))
    log STEP "Executing: PUT /api/expert-requests/$expert_request_id"
    update_request_payload='{"status":"approved"}'
    if execute_curl PUT "/api/expert-requests/$expert_request_id" "$update_request_payload" "$admin_token" 200; then
        log INFO "Expert request updated successfully."
        ((passed_tests++))
    else
        log ERROR "Failed to update expert request."
        ((failed_tests++))
    fi
else
    ((skipped_tests++))
    log SKIP "Skipping Get Expert Request Details (EXPERT_REQUEST_ID not available)."
    ((skipped_tests++))
    log SKIP "Skipping Update Expert Request (EXPERT_REQUEST_ID not available)."
fi

((total_tests++))
log STEP "Executing: POST /api/expert-requests (Reject Test)"
reject_request_payload='{
    "name": "Reject Request Expert '$TIMESTAMP'",
    "designation": "Researcher",
    "institution": "Reject University '$TIMESTAMP'",
    "email": "reject'$TIMESTAMP'@example.com",
    "generalArea": '$general_area_id'
}'
if execute_curl POST "/api/expert-requests" "$reject_request_payload" "$user_token" 201; then
    reject_request_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Expert request for rejection test created successfully. ID: $reject_request_id"
    ((passed_tests++))
else
    log ERROR "Failed to create expert request for rejection test."
    ((failed_tests++))
fi

echo -e "\n===== TESTING DOCUMENT MANAGEMENT ====="
if [[ -n "$expert_internal_id" ]]; then
    # Add document tests here if needed
    ((skipped_tests++))
    log SKIP "Document Management tests not implemented in this version."
else
    ((skipped_tests++))
    log SKIP "Skipping Document Management tests (EXPERT_INTERNAL_ID not available)."
fi

echo -e "\n===== TESTING STATISTICS ====="
((total_tests++))
log STEP "Executing: GET /api/statistics"
if execute_curl GET "/api/statistics" "" "$admin_token" 200; then
    log INFO "Retrieved system statistics successfully."
    ((passed_tests++))
else
    log ERROR "Failed to retrieve system statistics."
    ((failed_tests++))
fi

((total_tests++))
log STEP "Executing: GET /api/statistics/growth?months=6"
if execute_curl GET "/api/statistics/growth?months=6" "" "$admin_token" 200; then
    log INFO "Retrieved growth statistics successfully."
    ((passed_tests++))
else
    log ERROR "Failed to retrieve growth statistics."
    ((failed_tests++))
fi

((total_tests++))
log STEP "Executing: GET /api/statistics/nationality"
if execute_curl GET "/api/statistics/nationality" "" "$admin_token" 200; then
    log INFO "Retrieved nationality statistics successfully."
    ((passed_tests++))
else
    log ERROR "Failed to retrieve nationality statistics."
    ((failed_tests++))
fi

((total_tests++))
log STEP "Executing: GET /api/statistics/engagements"
if execute_curl GET "/api/statistics/engagements" "" "$admin_token" 200; then
    log INFO "Retrieved engagement statistics successfully."
    ((passed_tests++))
else
    log ERROR "Failed to retrieve engagement statistics."
    ((failed_tests++))
fi

echo -e "\n===== TESTING CLEANUP ====="
if [[ -n "$expert_internal_id" ]]; then
    ((total_tests++))
    log STEP "Executing: DELETE /api/experts/$expert_internal_id"
    if execute_curl DELETE "/api/experts/$expert_internal_id" "" "$admin_token" 200; then
        log INFO "Expert deleted successfully."
        ((passed_tests++))
    else
        log ERROR "Failed to delete expert."
        ((failed_tests++))
    fi
else
    ((skipped_tests++))
    log SKIP "Skipping Expert cleanup (EXPERT_INTERNAL_ID not available)."
fi

((total_tests++))
log STEP "Executing: DELETE /api/users/$user_id"
if execute_curl DELETE "/api/users/$user_id" "" "$admin_token" 200; then
    log INFO "Test user deleted successfully."
    ((passed_tests++))
else
    log ERROR "Failed to delete test user."
    ((failed_tests++))
fi

echo -e "\n===== TEST SUMMARY ====="
log INFO "Total tests run: $total_tests"
log SUCCESS "Passed: $passed_tests"
log ERROR "Failed: $failed_tests"
log SKIP "Skipped: $skipped_tests"

if [[ $failed_tests -gt 0 ]]; then
    log ERROR "Failed Tests:"
    # Note: Specific failed tests are logged inline above
fi

if [[ $skipped_tests -gt 0 ]]; then
    log SKIP "Skipped Tests:"
    log SKIP "- $([[ -n "$expert_internal_id" ]] || echo "Expert-related tests due to creation failure")"
fi

log INFO "Detailed logs available in: $LOG_FILE"

if [[ $failed_tests -gt 0 ]]; then
    log ERROR "===== Test Run Finished with Failures ====="
    exit 1
else
    log SUCCESS "===== Test Run Completed Successfully ====="
    exit 0
fi
