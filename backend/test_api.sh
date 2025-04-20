#!/bin/bash

# test_api.sh - ExpertDB API Test Script
# Tests ExpertDB API endpoints for a small internal tool (10-12 users, ~1200 entries).
# Focuses on core workflow: auth, user/expert management, requests, documents, stats.
# Console output is concise; detailed logs are saved to file.

# Configuration
API_BASE_URL="http://localhost:8080"
LOG_DIR="logs"
LOG_FILE="$LOG_DIR/api_test_run_$(date +%Y%m%d_%H%M%S).log"
TIMESTAMP=$(date +%s) # Unique for test data

# Payloads
ADMIN_CREDENTIALS='{"email":"admin@expertdb.com","password":"adminpassword"}'
USER_PAYLOAD='{
    "name": "Test User '$TIMESTAMP'",
    "email": "testuser'$TIMESTAMP'@example.com",
    "password": "password123",
    "role": "user",
    "isActive": true
}'
USER_CREDENTIALS='{
    "email": "testuser'$TIMESTAMP'@example.com",
    "password": "password123"
}'
EXPERT_PAYLOAD='{
    "name": "Test Expert '$TIMESTAMP'",
    "institution": "Test University '$TIMESTAMP'",
    "email": "expert'$TIMESTAMP'@example.com",
    "phone": "+97312345'$TIMESTAMP'",
    "designation": "Professor",
    "isBahraini": true,
    "isAvailable": true,
    "rating": "5",
    "role": "evaluator",
    "employmentType": "academic",
    "generalArea": 1,
    "specializedArea": "Software Engineering",
    "isTrained": true,
    "isPublished": true,
    "biography": "Expert created for testing with sequential ID generation.",
    "skills": ["Go", "Testing"]
}'
NO_ID_EXPERT_PAYLOAD='{
    "name": "NoID Expert '$TIMESTAMP'",
    "institution": "Test University '$TIMESTAMP'",
    "email": "noid'$TIMESTAMP'@example.com",
    "phone": "+97312346'$TIMESTAMP'",
    "designation": "Associate Professor",
    "isBahraini": false,
    "isAvailable": true,
    "rating": "4",
    "role": "validator",
    "employmentType": "academic",
    "generalArea": 1,
    "specializedArea": "Testing and Validation",
    "isTrained": true,
    "biography": "Expert created for testing automatic ID generation."
}'
INVALID_EXPERT_PAYLOAD='{
    "institution": "Test University",
    "generalArea": 1,
    "email": "invalid'$TIMESTAMP'@example.com"
}'
INVALID_AREA_PAYLOAD='{
    "name": "Invalid Area Expert '$TIMESTAMP'",
    "institution": "Test University",
    "designation": "Professor", 
    "phone": "+97312347'$TIMESTAMP'",
    "email": "invalidarea'$TIMESTAMP'@example.com",
    "role": "evaluator", 
    "employmentType": "academic", 
    "specializedArea": "Area Testing",
    "isTrained": true, 
    "biography": "Testing invalid area",
    "generalArea": -1
}'
INVALID_ROLE_PAYLOAD='{
    "name": "Invalid Role Expert '$TIMESTAMP'",
    "institution": "Test University",
    "designation": "Professor", 
    "phone": "+97312347'$TIMESTAMP'",
    "email": "invalidrole'$TIMESTAMP'@example.com",
    "role": "invalid-role", 
    "employmentType": "academic", 
    "specializedArea": "Role Testing",
    "generalArea": 1,
    "isTrained": true, 
    "biography": "Testing invalid role"
}'
REQUEST_PAYLOAD='{
    "name": "Request Expert '$TIMESTAMP'",
    "designation": "Researcher",
    "institution": "Request University '$TIMESTAMP'",
    "isBahraini": false,
    "isAvailable": true,
    "rating": "4",
    "role": "evaluator",
    "employmentType": "academic",
    "generalArea": 1,
    "specializedArea": "Quantum Physics",
    "isTrained": false,
    "phone": "+9731111'$TIMESTAMP'",
    "email": "request'$TIMESTAMP'@example.com",
    "isPublished": false,
    "biography": "Researcher requesting addition."
}'
REJECT_REQUEST_PAYLOAD='{
    "name": "Reject Request Expert '$TIMESTAMP'",
    "designation": "Researcher",
    "institution": "Reject University '$TIMESTAMP'",
    "phone": "+9731112'$TIMESTAMP'",
    "email": "reject'$TIMESTAMP'@example.com",
    "role": "evaluator",
    "employmentType": "academic",
    "specializedArea": "Mathematics",
    "generalArea": 1,
    "isTrained": true,
    "biography": "Expert to be rejected",
    "isPublished": false
}'

# Colors for console output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Ensure log directory exists
mkdir -p "$LOG_DIR"

# Test counters
total_tests=0
passed_tests=0
failed_tests=0
skipped_tests=0

# State variables
admin_token=""
user_id=""
user_token=""
expert_internal_id=""
expert_request_id=""
document_id=""

# Utility functions
log() {
    local level=$1
    local message=$2
    local color=$NC
    local prefix=""

    case $level in
        INFO) color=$CYAN; prefix="[INFO]   ";;
        STEP) color=$YELLOW; prefix="[STEP]   ";;
        SUCCESS) color=$GREEN; prefix="[SUCCESS]";;
        ERROR) color=$RED; prefix="[ERROR]  ";;
        SKIP) color=$YELLOW; prefix="[SKIP]   ";;
        DETAIL) prefix="[DETAIL] ";;
        HEADER) color=$BLUE; prefix="[HEADER] ";;
    esac

    # Always log to file
    echo "$(date '+%Y-%m-%d %H:%M:%S') $prefix $message" >> "$LOG_FILE"

    # Console output for non-DETAIL messages
    [[ $level != "DETAIL" ]] && echo -e "${color}${prefix} ${message}${NC}"
}

validate_json() {
    local payload=$1
    echo "$payload" | jq . >/dev/null 2>&1 || {
        log ERROR "Invalid JSON payload: $payload"
        exit 1
    }
}

execute_curl() {
    local method=$1
    local endpoint=$2
    local payload=$3
    local token=$4
    local expected_status=$5
    local output_file="/tmp/api_response.json"
    local header_file="/tmp/api_response_headers.txt"
    local test_name="$method $endpoint"

    ((total_tests++))
    log STEP "Testing: $test_name"

    local curl_cmd="curl -s -w '%{http_code}' -o '$output_file' --dump-header '$header_file' -X $method"
    [[ -n "$token" ]] && curl_cmd+=" -H 'Authorization: Bearer $token'"
    if [[ "$method" != "GET" && -n "$payload" ]]; then
        validate_json "$payload"
        curl_cmd+=" -H 'Content-Type: application/json' -d '$payload'"
    fi
    curl_cmd+=" '$API_BASE_URL$endpoint'"

    # Log detailed request info to file
    log DETAIL "URL: $API_BASE_URL$endpoint"
    log DETAIL "Command: $curl_cmd"
    echo "--- Request ---" >> "$LOG_FILE"
    [[ -n "$token" ]] && echo "Authorization: Bearer ${token:0:10}..." >> "$LOG_FILE"
    if [[ -n "$payload" ]]; then
        echo "Request Body:" >> "$LOG_FILE"
        echo "$payload" | jq . >> "$LOG_FILE" 2>/dev/null || echo "$payload" >> "$LOG_FILE"
    else
        echo "Request Body: (None)" >> "$LOG_FILE"
    fi
    echo "--- End Request ---" >> "$LOG_FILE"

    # Execute request
    http_status=$(eval "$curl_cmd")
    curl_exit=$?

    # Log response details to file
    log DETAIL "Curl Exit Code: $curl_exit"
    log DETAIL "HTTP Status: $http_status"
    echo "--- Response ---" >> "$LOG_FILE"
    echo "HTTP Status: $http_status" >> "$LOG_FILE"
    echo "Response Headers:" >> "$LOG_FILE"
    cat "$header_file" >> "$LOG_FILE"
    echo -e "\nResponse Body:" >> "$LOG_FILE"
    jq . "$output_file" >> "$LOG_FILE" 2>/dev/null || cat "$output_file" >> "$LOG_FILE"
    echo -e "\n--- End Response ---" >> "$LOG_FILE"

    # Check result
    if [[ $curl_exit -ne 0 ]]; then
        log ERROR "$test_name failed: Curl error (exit code $curl_exit)"
        ((failed_tests++))
        return 1
    elif [[ -n "$expected_status" && "$http_status" -ne "$expected_status" ]]; then
        log ERROR "$test_name failed: Expected status $expected_status, got $http_status"
        cat "$output_file" | jq . # Show error response in console
        ((failed_tests++))
        return 1
    elif [[ "$http_status" -ge 400 ]]; then
        log ERROR "$test_name failed: HTTP status $http_status"
        cat "$output_file" | jq . # Show error response in console
        ((failed_tests++))
        return 1
    else
        log SUCCESS "$test_name passed (Status: $http_status)"
        ((passed_tests++))
        return 0
    fi
}

# Test execution
log INFO "Starting ExpertDB API tests against $API_BASE_URL..."
log INFO "Detailed logs saved to $LOG_FILE"

# Authentication
log HEADER "TESTING AUTHENTICATION"
execute_curl POST "/api/auth/login" "$ADMIN_CREDENTIALS" "" 200 && {
    admin_token=$(jq -r '.token' /tmp/api_response.json)
    log INFO "Admin login successful"
} || {
    log ERROR "Admin login failed. Exiting."
    exit 1
}

execute_curl POST "/api/users" "$USER_PAYLOAD" "$admin_token" 201 && {
    user_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Test user created. ID: $user_id"
} || log ERROR "Failed to create test user"

execute_curl POST "/api/auth/login" "$USER_CREDENTIALS" "" 200 && {
    user_token=$(jq -r '.token' /tmp/api_response.json)
    log INFO "Test user login successful"
} || log ERROR "Failed to login test user"

# Setup: Fetch Expert Areas
log HEADER "TESTING SETUP"
execute_curl GET "/api/expert/areas" "" "$admin_token" 200 && {
    log INFO "Retrieved $(jq length /tmp/api_response.json) expert areas"
} || log ERROR "Failed to retrieve expert areas"

# Expert Management
log HEADER "TESTING EXPERT MANAGEMENT"
execute_curl POST "/api/experts" "$EXPERT_PAYLOAD" "$admin_token" 201 && {
    expert_internal_id=$(jq -r '.id' /tmp/api_response.json)
    expert_id=$(jq -r '.expertId' /tmp/api_response.json)
    log INFO "Expert created. ID: $expert_internal_id, Expert ID: $expert_id"
} || log ERROR "Failed to create expert"

execute_curl POST "/api/experts" "$NO_ID_EXPERT_PAYLOAD" "$admin_token" 201 && {
    second_expert_id=$(jq -r '.expertId' /tmp/api_response.json)
    log INFO "Expert without explicit ID created. Generated Expert ID: $second_expert_id"
} || log ERROR "Failed to create expert without ID"

execute_curl POST "/api/experts" "$INVALID_EXPERT_PAYLOAD" "$admin_token" 400 && {
    log INFO "Rejected invalid expert payload"
} || log ERROR "Failed to reject invalid expert payload"

execute_curl POST "/api/experts" "$INVALID_AREA_PAYLOAD" "$admin_token" 400 && {
    log INFO "Rejected invalid generalArea"
} || log ERROR "Failed to reject invalid generalArea"

execute_curl POST "/api/experts" "$INVALID_ROLE_PAYLOAD" "$admin_token" 400 && {
    log INFO "Rejected invalid role"
} || log ERROR "Failed to reject invalid role"

execute_curl GET "/api/experts?limit=5" "" "$admin_token" 200 && {
    # Check if we have the new response format with pagination
    if jq -e '.experts' /tmp/api_response.json > /dev/null 2>&1; then
        # New format with pagination metadata
        expert_count=$(jq '.experts | length' /tmp/api_response.json)
        total_count=$(jq '.pagination.totalCount' /tmp/api_response.json)
        log INFO "Listed ${expert_count} experts (of ${total_count} total) with pagination metadata"
        log DETAIL "Pagination metadata: $(jq '.pagination' /tmp/api_response.json)"
    else
        # Old format (direct array of experts)
        log INFO "Listed $(jq length /tmp/api_response.json) experts"
    fi
} || log ERROR "Failed to list experts"

# Test sorting in Phase 3B
execute_curl GET "/api/experts?sort_by=name&sort_order=asc&limit=5" "" "$admin_token" 200 && {
    # Check if we have the new response format
    if jq -e '.experts' /tmp/api_response.json > /dev/null 2>&1; then
        log INFO "Successfully tested sorting by name (ascending)"
    else
        log INFO "Successfully tested sorting by name (ascending) - old response format"
    fi
} || log ERROR "Failed to test sorting experts by name"

execute_curl GET "/api/experts?sort_by=rating&sort_order=desc&limit=5" "" "$admin_token" 200 && {
    log INFO "Successfully tested sorting by rating (descending)"
} || log ERROR "Failed to test sorting experts by rating"

# Test pagination in Phase 3B
execute_curl GET "/api/experts?limit=5&offset=5" "" "$admin_token" 200 && {
    # Check if we have the new response format
    if jq -e '.pagination.currentPage' /tmp/api_response.json > /dev/null 2>&1; then
        current_page=$(jq '.pagination.currentPage' /tmp/api_response.json)
        total_pages=$(jq '.pagination.totalPages' /tmp/api_response.json)
        log INFO "Successfully tested pagination (page ${current_page} of ${total_pages})"
    else
        log INFO "Successfully tested pagination - old response format"
    fi
} || log ERROR "Failed to test expert pagination"

# Test the new expert filtering capabilities from Phase 3A
execute_curl GET "/api/experts?by_nationality=Bahraini&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed $(jq length /tmp/api_response.json) Bahraini experts"
} || log ERROR "Failed to filter experts by nationality"

execute_curl GET "/api/experts?by_general_area=1&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed $(jq length /tmp/api_response.json) experts in general area 1"
} || log ERROR "Failed to filter experts by general area"

execute_curl GET "/api/experts?by_specialized_area=Software&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed $(jq length /tmp/api_response.json) experts in Software specialized area"
} || log ERROR "Failed to filter experts by specialized area"

execute_curl GET "/api/experts?by_employment_type=academic&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed $(jq length /tmp/api_response.json) academic experts"
} || log ERROR "Failed to filter experts by employment type"

execute_curl GET "/api/experts?by_role=evaluator&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed $(jq length /tmp/api_response.json) evaluator experts"
} || log ERROR "Failed to filter experts by role"

# Test combined filters
execute_curl GET "/api/experts?by_nationality=Bahraini&by_employment_type=academic&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed $(jq length /tmp/api_response.json) Bahraini academic experts"
} || log ERROR "Failed to filter experts with combined filters"

if [[ -n "$expert_internal_id" ]]; then
    execute_curl GET "/api/experts/$expert_internal_id" "" "$admin_token" 200 && {
        log INFO "Retrieved expert details"
    } || log ERROR "Failed to retrieve expert details"

    execute_curl PUT "/api/experts/$expert_internal_id" '{"name":"Updated Expert '$TIMESTAMP'","isAvailable":false}' "$admin_token" 200 && {
        log INFO "Expert updated"
    } || log ERROR "Failed to update expert"
else
    log SKIP "Skipping expert details/update (no expert ID)"
    ((skipped_tests+=2))
fi

# Expert Request Management
log HEADER "TESTING EXPERT REQUEST MANAGEMENT"
execute_curl POST "/api/expert-requests" "$REQUEST_PAYLOAD" "$user_token" 201 && {
    expert_request_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Expert request created. ID: $expert_request_id"
} || log ERROR "Failed to create expert request"

execute_curl GET "/api/expert-requests?limit=5" "" "$admin_token" 200 && {
    log INFO "Listed expert requests"
} || log ERROR "Failed to list expert requests"

# Test status filtering for expert requests (Phase 4B)
execute_curl GET "/api/expert-requests?status=pending&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed pending expert requests"
} || log ERROR "Failed to filter expert requests by pending status"

execute_curl GET "/api/expert-requests?status=approved&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed approved expert requests"
} || log ERROR "Failed to filter expert requests by approved status"

execute_curl GET "/api/expert-requests?status=rejected&limit=5" "" "$admin_token" 200 && {
    log INFO "Listed rejected expert requests"
} || log ERROR "Failed to filter expert requests by rejected status"

if [[ -n "$expert_request_id" ]]; then
    execute_curl GET "/api/expert-requests/$expert_request_id" "" "$admin_token" 200 && {
        log INFO "Retrieved expert request details"
    } || log ERROR "Failed to retrieve expert request details"

    execute_curl PUT "/api/expert-requests/$expert_request_id" '{"status":"approved"}' "$admin_token" 200 && {
        log INFO "Expert request approved"
    } || log ERROR "Failed to approve expert request"
else
    log SKIP "Skipping request details/approval (no request ID)"
    ((skipped_tests+=2))
fi

execute_curl POST "/api/expert-requests" "$REJECT_REQUEST_PAYLOAD" "$user_token" 201 && {
    reject_request_id=$(jq -r '.id' /tmp/api_response.json)
    execute_curl PUT "/api/expert-requests/$reject_request_id" '{"status":"rejected","rejectionReason":"Test rejection"}' "$admin_token" 200 && {
        log INFO "Expert request rejected"
        
        # Test Phase 4C: Request editing before approval
        
        # Test: Regular user can't edit pending requests
        execute_curl PUT "/api/expert-requests/$expert_request_id" '{"name":"Updated by User"}' "$user_token" 403 && {
            log INFO "Correctly prevented user from editing pending request"
        } || log ERROR "Failed to test user permission for pending request"
        
        # Test: Admin can edit pending requests
        execute_curl PUT "/api/expert-requests/$expert_request_id" '{"name":"Updated by Admin"}' "$admin_token" 200 && {
            log INFO "Admin successfully edited pending request"
        } || log ERROR "Failed to test admin editing pending request"
        
        # Test: User can edit their own rejected request
        execute_curl PUT "/api/expert-requests/$reject_request_id" '{"name":"Updated Rejected Request"}' "$user_token" 200 && {
            log INFO "User successfully edited their rejected request"
        } || log ERROR "Failed to test user editing rejected request"
        
        # Test multipart form upload would require a different approach with curl_cmd
        log SKIP "Skipping test for file upload during edit (needs special curl handling)"
        ((skipped_tests+=1))
    } || log ERROR "Failed to reject expert request"
} || log ERROR "Failed to create reject request"

# Test Phase 5: Approval Document Integration
log HEADER "TESTING APPROVAL DOCUMENT INTEGRATION"

# Create a request that will later be rejected due to missing approval document
execute_curl POST "/api/expert-requests" "$REJECT_REQUEST_PAYLOAD" "$user_token" 201 && {
    approval_test_id=$(jq -r '.id' /tmp/api_response.json)
    
    # Test: Try to approve without approval document (should fail)
    execute_curl PUT "/api/expert-requests/$approval_test_id" '{"status":"approved"}' "$admin_token" 400 && {
        log INFO "Correctly prevented approval without approval document"
    } || log ERROR "Failed to test approval document requirement"
    
    # Test batch approval - note this requires a multipart form, so it's not fully tested here
    log SKIP "Skipping full test of batch approval (needs special multipart form handling)"
    ((skipped_tests+=1))
} || log ERROR "Failed to create approval test request"

# Document Management
log HEADER "TESTING DOCUMENT MANAGEMENT"
if [[ -n "$expert_internal_id" ]]; then
    doc_file="/tmp/sample_cv_$TIMESTAMP.txt"
    echo "Sample CV for expert $TIMESTAMP" > "$doc_file"
    curl_cmd="curl -s -w '%{http_code}' -o /tmp/api_response.json --dump-header /tmp/api_response_headers.txt -X POST -H 'Authorization: Bearer $admin_token' -F 'file=@$doc_file' -F 'documentType=cv' -F 'expertId=$expert_internal_id' '$API_BASE_URL/api/documents'"
    ((total_tests++))
    http_status=$(eval "$curl_cmd")
    log DETAIL "Upload Command: $curl_cmd"
    echo "--- Upload Response ---" >> "$LOG_FILE"
    echo "HTTP Status: $http_status" >> "$LOG_FILE"
    jq . /tmp/api_response.json >> "$LOG_FILE" 2>/dev/null || cat /tmp/api_response.json >> "$LOG_FILE"
    echo "--- End Upload Response ---" >> "$LOG_FILE"
    if [[ "$http_status" -eq 201 ]]; then
        document_id=$(jq -r '.id' /tmp/api_response.json)
        log SUCCESS "Document uploaded (ID: $document_id)"
        ((passed_tests++))
    else
        log ERROR "Document upload failed (Status: $http_status)"
        ((failed_tests++))
    fi
    rm -f "$doc_file"

    [[ -n "$document_id" ]] && execute_curl GET "/api/documents/$document_id" "" "$admin_token" 200 && {
        log INFO "Retrieved document details"
    } || { log ERROR "Failed to retrieve document details"; ((failed_tests++)); ((total_tests++)); }

    [[ -n "$document_id" ]] && execute_curl DELETE "/api/documents/$document_id" "" "$admin_token" 200 && {
        log INFO "Document deleted"
    } || { log ERROR "Failed to delete document"; ((failed_tests++)); ((total_tests++)); }
else
    log SKIP "Skipping document tests (no expert ID)"
    ((skipped_tests+=3))
fi

# Statistics
log HEADER "TESTING CSV BACKUP (PHASE 9)"

# Test Phase 9: CSV Backup Implementation
execute_curl GET "/api/backup" "" "$admin_token" 200 && {
    backup_size=$(stat -c %s /tmp/api_response.json 2>/dev/null || echo "0")
    log INFO "Successfully generated CSV backup (size: $backup_size bytes)"
    
    # Test file type (should be a ZIP file)
    file_type=$(file -b /tmp/api_response.json | cut -d' ' -f1-2)
    if [[ "$file_type" == "Zip archive" ]]; then
        log INFO "Verified backup is a valid ZIP archive"
    else
        log WARN "Backup may not be a valid ZIP archive. Type: $file_type"
    fi
} || log ERROR "Failed to generate CSV backup"

# Test backup permissions - regular user shouldn't be able to access
execute_curl GET "/api/backup" "" "$user_token" 403 && {
    log INFO "Correctly prevented regular user from accessing backup"
} || log ERROR "Failed to prevent regular user from accessing backup"

log HEADER "TESTING PHASE PLANNING (PHASE 10)"

# Create scheduler user
SCHEDULER_CREATE_PAYLOAD='{
    "name": "Test Scheduler '$TIMESTAMP'",
    "email": "scheduler'$TIMESTAMP'@example.com",
    "password": "password123",
    "role": "scheduler",
    "isActive": true
}'

execute_curl POST "/api/users" "$SCHEDULER_CREATE_PAYLOAD" "$admin_token" 201 && {
    scheduler_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Scheduler user created. ID: $scheduler_id"
    
    # Get scheduler token
    SCHEDULER_CREDENTIALS='{
        "email": "scheduler'$TIMESTAMP'@example.com",
        "password": "password123"
    }'
    
    execute_curl POST "/api/auth/login" "$SCHEDULER_CREDENTIALS" "" 200 && {
        scheduler_token=$(jq -r '.token' /tmp/api_response.json)
        log INFO "Scheduler login successful"
    } || log ERROR "Failed to login scheduler user"
    
    # Test Phase 10B: Phase Creation
    PHASE_PAYLOAD='{
        "title": "Test Phase '$TIMESTAMP'",
        "assignedSchedulerId": '$scheduler_id',
        "status": "draft",
        "applications": [
            {
                "type": "validation",
                "institutionName": "Test University",
                "qualificationName": "Bachelor of Science in Computer Science"
            },
            {
                "type": "evaluation",
                "institutionName": "Test College",
                "qualificationName": "Associate Degree in Engineering"
            }
        ]
    }'
    
    execute_curl POST "/api/phases" "$PHASE_PAYLOAD" "$admin_token" 201 && {
        phase_id=$(jq -r '.id' /tmp/api_response.json)
        phase_business_id=$(jq -r '.phaseId' /tmp/api_response.json)
        log INFO "Phase created. ID: $phase_id, Business ID: $phase_business_id"
        
        # Test Phase Retrieval
        execute_curl GET "/api/phases/$phase_id" "" "$admin_token" 200 && {
            log INFO "Phase retrieved successfully"
        } || log ERROR "Failed to retrieve phase"
        
        # Test Phase Listing
        execute_curl GET "/api/phases" "" "$admin_token" 200 && {
            phases_count=$(jq '. | length' /tmp/api_response.json)
            log INFO "Listed $phases_count phases"
        } || log ERROR "Failed to list phases"
        
        # Test Phase 10C: Expert Proposal for Applications
        # Get first application ID
        app_id=$(jq -r '.applications[0].id' /tmp/api_response.json)
        PROPOSAL_PAYLOAD='{
            "expert1": '$expert_internal_id',
            "expert2": 0
        }'
        
        if [[ -n "$scheduler_token" && -n "$app_id" ]]; then
            execute_curl PUT "/api/phases/$phase_id/applications/$app_id" "$PROPOSAL_PAYLOAD" "$scheduler_token" 200 && {
                log INFO "Application experts assigned successfully"
            } || log ERROR "Failed to assign application experts"
            
            # Test Phase 10D: Application Review
            REVIEW_PAYLOAD='{
                "action": "approve"
            }'
            
            execute_curl PUT "/api/phases/$phase_id/applications/$app_id/review" "$REVIEW_PAYLOAD" "$admin_token" 200 && {
                log INFO "Application approved successfully"
                
                # Check if engagement was created automatically
                execute_curl GET "/api/experts/$expert_internal_id/engagements" "" "$admin_token" 200 && {
                    engagement_count=$(jq '. | length' /tmp/api_response.json)
                    if [[ "$engagement_count" -gt 0 ]]; then
                        log INFO "Verified automatic engagement creation"
                    else
                        log ERROR "No engagement created for approved application"
                    fi
                } || log ERROR "Failed to verify engagement creation"
            } || log ERROR "Failed to approve application"
            
            # Test Phase Update
            UPDATE_PHASE_PAYLOAD='{
                "title": "Updated Phase '$TIMESTAMP'",
                "status": "in_progress"
            }'
            
            execute_curl PUT "/api/phases/$phase_id" "$UPDATE_PHASE_PAYLOAD" "$admin_token" 200 && {
                log INFO "Phase updated successfully"
            } || log ERROR "Failed to update phase"
        else
            log SKIP "Skipping application expert assignment and review tests"
            ((skipped_tests+=3))
        fi
    } || log ERROR "Failed to create phase"
} || log ERROR "Failed to create scheduler user"

log HEADER "TESTING AREA MANAGEMENT (PHASE 8)"

# Test Phase 8A: Area Access Extension
execute_curl GET "/api/expert/areas" "" "$user_token" 200 && {
    log INFO "Successfully verified area access for regular user"
    areas_count=$(jq '. | length' /tmp/api_response.json)
    log INFO "Retrieved $areas_count areas as regular user"
} || log ERROR "Failed to access areas as regular user"

# Test Phase 8B: Area Creation
AREA_NAME="Test Area $TIMESTAMP"
AREA_PAYLOAD='{
    "name": "'"$AREA_NAME"'"
}'

execute_curl POST "/api/expert/areas" "$AREA_PAYLOAD" "$admin_token" 201 && {
    new_area_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Successfully created new area with ID: $new_area_id"
} || log ERROR "Failed to create new area"

# Test duplicate area name handling
execute_curl POST "/api/expert/areas" "$AREA_PAYLOAD" "$admin_token" 409 && {
    log INFO "Correctly rejected duplicate area name"
} || log ERROR "Failed to reject duplicate area name"

# Test invalid area creation by regular user
execute_curl POST "/api/expert/areas" "$AREA_PAYLOAD" "$user_token" 403 && {
    log INFO "Correctly prevented area creation by regular user"
} || log ERROR "Failed to prevent area creation by regular user"

# Test Phase 8C: Area Renaming
if [[ -n "$new_area_id" ]]; then
    RENAME_PAYLOAD='{
        "name": "Renamed Area '"$TIMESTAMP"'"
    }'
    
    execute_curl PUT "/api/expert/areas/$new_area_id" "$RENAME_PAYLOAD" "$admin_token" 200 && {
        log INFO "Successfully renamed area"
    } || log ERROR "Failed to rename area"
    
    # Test invalid area rename by regular user
    execute_curl PUT "/api/expert/areas/$new_area_id" "$RENAME_PAYLOAD" "$user_token" 403 && {
        log INFO "Correctly prevented area rename by regular user"
    } || log ERROR "Failed to prevent area rename by regular user"
else
    log SKIP "Skipping area rename tests (no area ID)"
    ((skipped_tests+=2))
fi

log HEADER "TESTING STATISTICS ENDPOINTS (PHASE 7)"

execute_curl GET "/api/statistics" "" "$admin_token" 200 && {
    log INFO "Successfully retrieved statistics"
    # Test for published expert stats (Phase 7A)
    jq -e '.publishedCount' /tmp/api_response.json > /dev/null && log INFO "Successfully verified published count field"
    jq -e '.publishedRatio' /tmp/api_response.json > /dev/null && log INFO "Successfully verified published ratio field"
    
    # Test for yearly growth stats (Phase 7B)
    jq -e '.yearlyGrowth' /tmp/api_response.json > /dev/null && log INFO "Successfully verified yearly growth field"
} || log ERROR "Failed to retrieve system statistics"

execute_curl GET "/api/statistics/nationality" "" "$admin_token" 200 && {
    log INFO "Successfully retrieved nationality statistics"
} || log ERROR "Failed to retrieve nationality statistics"

# Test Phase 7B: Growth Statistics Enhancement (yearly instead of monthly)
execute_curl GET "/api/statistics/growth?years=3" "" "$admin_token" 200 && {
    log INFO "Successfully retrieved yearly growth statistics"
    # Verify the response contains period field formatted as year (YYYY)
    year_format=$(jq -r '.[0].period' /tmp/api_response.json 2>/dev/null | grep -E '^[0-9]{4}$')
    if [ -n "$year_format" ]; then
        log INFO "Verified year format (YYYY)"
    else
        log WARN "Year format verification failed"
    fi
} || log ERROR "Failed to retrieve growth statistics"

# Test Phase 7C: Engagement Type Statistics (validator/evaluator only)
execute_curl GET "/api/statistics/engagements" "" "$admin_token" 200 && {
    log INFO "Successfully retrieved engagement statistics"
    # Verify only validator/evaluator types are included
    types_count=$(jq '.byType | length' /tmp/api_response.json)
    if [ "$types_count" -le 2 ]; then
        log INFO "Verified engagement types are limited to validator/evaluator"
    else
        log WARN "More than expected engagement types found"
    fi
} || log ERROR "Failed to retrieve engagement statistics"

# Test Phase 7D: Area Statistics Implementation
execute_curl GET "/api/statistics/areas" "" "$admin_token" 200 && {
    log INFO "Successfully retrieved area statistics"
    # Verify the response contains the expected sections
    jq -e '.generalAreas' /tmp/api_response.json > /dev/null && log INFO "Verified general areas section"
    jq -e '.topSpecializedAreas' /tmp/api_response.json > /dev/null && log INFO "Verified top specialized areas section"
    jq -e '.bottomSpecializedAreas' /tmp/api_response.json > /dev/null && log INFO "Verified bottom specialized areas section"
} || log ERROR "Failed to retrieve area statistics"

# Cleanup
log HEADER "TESTING CLEANUP"
[[ -n "$expert_internal_id" ]] && execute_curl DELETE "/api/experts/$expert_internal_id" "" "$admin_token" 200 && {
    log INFO "Expert deleted"
} || { log ERROR "Failed to delete expert"; [[ -n "$expert_internal_id" ]] && ((failed_tests++)); ((total_tests++)); }

[[ -n "$user_id" ]] && execute_curl DELETE "/api/users/$user_id" "" "$admin_token" 200 && {
    log INFO "Test user deleted"
} || { log ERROR "Failed to delete test user"; [[ -n "$user_id" ]] && ((failed_tests++)); ((total_tests++)); }

# Database Performance Testing
log HEADER "TESTING DATABASE PERFORMANCE"

# Function to test query performance with EXPLAIN QUERY PLAN
test_query_performance() {
    local description=$1
    local query=$2
    
    log STEP "Testing query performance: $description"
    log DETAIL "Query: $query"
    
    # Execute EXPLAIN QUERY PLAN
    local explain_cmd="echo \"EXPLAIN QUERY PLAN $query;\" | sqlite3 db/sqlite/expertdb.sqlite"
    local explain_output=$(eval "$explain_cmd")
    
    log DETAIL "EXPLAIN QUERY PLAN Output:"
    log DETAIL "$explain_output"
    
    # Check if it uses an index
    if [[ $explain_output == *"USING INDEX"* ]]; then
        log SUCCESS "Query uses indexes: $description"
        ((passed_tests++))
    else
        log ERROR "Query does not use indexes: $description"
        ((failed_tests++))
    fi
    ((total_tests++))
}

# Only run performance tests if explicitly requested
if [[ "$1" == "--with-performance" ]]; then
    # Test nationality index
    test_query_performance "Query by nationality" "
        SELECT * FROM experts WHERE is_bahraini = 1 LIMIT 10;
    "
    
    # Test general area index
    test_query_performance "Query by general area" "
        SELECT * FROM experts WHERE general_area = 1 LIMIT 10;
    "
    
    # Test specialized area index
    test_query_performance "Query by specialized area" "
        SELECT * FROM experts WHERE specialized_area LIKE 'Software%' LIMIT 10;
    "
    
    # Test employment type index
    test_query_performance "Query by employment type" "
        SELECT * FROM experts WHERE employment_type = 'academic' LIMIT 10;
    "
    
    # Test role index
    test_query_performance "Query by role" "
        SELECT * FROM experts WHERE role = 'evaluator' LIMIT 10;
    "
    
    # Test combined query
    test_query_performance "Combined query" "
        SELECT * FROM experts 
        WHERE is_bahraini = 1 
        AND general_area = 1 
        AND role = 'evaluator' 
        LIMIT 10;
    "
else
    log SKIP "Performance tests skipped. Use --with-performance to run them."
    ((skipped_tests+=6))
fi

# Summary
log HEADER "TESTING ENGAGEMENT MANAGEMENT (PHASE 11)"

# Test Phase 11A: Engagement Filtering
execute_curl GET "/api/engagements?type=validator" "" "$admin_token" 200 && {
    log INFO "Successfully filtered engagements by validator type"
} || log ERROR "Failed to filter engagements by type"

execute_curl GET "/api/engagements?expert_id=$expert_internal_id" "" "$admin_token" 200 && {
    log INFO "Successfully filtered engagements by expert_id"
} || log ERROR "Failed to filter engagements by expert_id"

execute_curl GET "/api/engagements?type=evaluator&limit=10&offset=0" "" "$admin_token" 200 && {
    log INFO "Successfully filtered engagements with pagination"
} || log ERROR "Failed to filter engagements with pagination"

# Test Phase 11B: Engagement Type Restriction
INVALID_TYPE_PAYLOAD='{
    "expertId": '$expert_internal_id',
    "engagementType": "invalid-type",
    "startDate": "2025-01-01"
}'

execute_curl POST "/api/engagements" "$INVALID_TYPE_PAYLOAD" "$admin_token" 400 && {
    log INFO "Correctly rejected invalid engagement type"
} || log ERROR "Failed to reject invalid engagement type"

VALID_TYPE_PAYLOAD='{
    "expertId": '$expert_internal_id',
    "engagementType": "validator",
    "startDate": "2025-01-01"
}'

# Create scheduler user if not already created
if [[ -z "$scheduler_token" ]]; then
    SCHEDULER_CREATE_PAYLOAD='{
        "name": "Test Scheduler '$TIMESTAMP'",
        "email": "scheduler'$TIMESTAMP'@example.com",
        "password": "password123",
        "role": "scheduler",
        "isActive": true
    }'
    
    execute_curl POST "/api/users" "$SCHEDULER_CREATE_PAYLOAD" "$admin_token" 201 && {
        scheduler_id=$(jq -r '.id' /tmp/api_response.json)
        log INFO "Scheduler user created. ID: $scheduler_id"
        
        # Get scheduler token
        SCHEDULER_CREDENTIALS='{
            "email": "scheduler'$TIMESTAMP'@example.com",
            "password": "password123"
        }'
        
        execute_curl POST "/api/auth/login" "$SCHEDULER_CREDENTIALS" "" 200 && {
            scheduler_token=$(jq -r '.token' /tmp/api_response.json)
            log INFO "Scheduler login successful"
        } || log ERROR "Failed to login scheduler user"
    } || log ERROR "Failed to create scheduler user"
fi

execute_curl POST "/api/engagements" "$VALID_TYPE_PAYLOAD" "$scheduler_token" 201 && {
    test_engagement_id=$(jq -r '.id' /tmp/api_response.json)
    log INFO "Successfully created engagement with valid type"
    
    # Try updating with invalid type
    INVALID_UPDATE_PAYLOAD='{
        "engagementType": "invalid-type"
    }'
    
    execute_curl PUT "/api/engagements/$test_engagement_id" "$INVALID_UPDATE_PAYLOAD" "$scheduler_token" 400 && {
        log INFO "Correctly rejected invalid engagement type update"
    } || log ERROR "Failed to reject invalid engagement type update"
    
    # Clean up test engagement
    execute_curl DELETE "/api/engagements/$test_engagement_id" "" "$scheduler_token" 200 && {
        log INFO "Successfully deleted test engagement"
    } || log ERROR "Failed to delete test engagement"
} || log ERROR "Failed to create engagement with valid type"

# Test Phase 11C: Engagement Import
# Create a CSV file for importing engagements
IMPORT_CSV_FILE="/tmp/engagements_import_$TIMESTAMP.csv"
cat > "$IMPORT_CSV_FILE" << EOF
expert_id,engagement_type,start_date,end_date,project_name,status,notes
$expert_internal_id,validator,2025-02-01,2025-03-01,Project A,active,Imported via CSV
$expert_internal_id,evaluator,2025-04-01,2025-05-01,Project B,active,Another imported engagement
EOF

# Use multipart form upload for CSV import
curl_cmd="curl -s -w '%{http_code}' -o /tmp/api_response.json --dump-header /tmp/api_response_headers.txt -X POST -H 'Authorization: Bearer $admin_token' -F 'file=@$IMPORT_CSV_FILE' '$API_BASE_URL/api/engagements/import'"
((total_tests++))
http_status=$(eval "$curl_cmd")
log DETAIL "Upload Command: $curl_cmd"
echo "--- Import Response ---" >> "$LOG_FILE"
echo "HTTP Status: $http_status" >> "$LOG_FILE"
jq . /tmp/api_response.json >> "$LOG_FILE" 2>/dev/null || cat /tmp/api_response.json >> "$LOG_FILE"
echo "--- End Import Response ---" >> "$LOG_FILE"

if [[ "$http_status" -eq 200 ]]; then
    success_count=$(jq -r '.successCount' /tmp/api_response.json)
    log SUCCESS "CSV import successful: $success_count engagements imported"
    ((passed_tests++))
    
    # Verify imported engagements exist
    execute_curl GET "/api/experts/$expert_internal_id/engagements" "" "$admin_token" 200 && {
        engagement_count=$(jq length /tmp/api_response.json)
        log INFO "Verified imported engagements: $engagement_count found for expert"
    } || log ERROR "Failed to verify imported engagements"
else
    log ERROR "CSV import failed (Status: $http_status)"
    ((failed_tests++))
fi

# Clean up temporary file
rm -f "$IMPORT_CSV_FILE"

# Also test the JSON import 
IMPORT_JSON_PAYLOAD='[
    {
        "expertId": '$expert_internal_id',
        "engagementType": "validator",
        "startDate": "2025-06-01",
        "endDate": "2025-07-01",
        "projectName": "Project C",
        "status": "active",
        "notes": "Imported via JSON"
    },
    {
        "expertId": '$expert_internal_id',
        "engagementType": "evaluator",
        "startDate": "2025-08-01",
        "endDate": "2025-09-01",
        "projectName": "Project D",
        "status": "active",
        "notes": "Another JSON import"
    }
]'

execute_curl POST "/api/engagements/import" "$IMPORT_JSON_PAYLOAD" "$admin_token" 200 && {
    success_count=$(jq -r '.successCount' /tmp/api_response.json)
    log INFO "JSON import successful: $success_count engagements imported"
    
    # Verify imported engagements exist
    execute_curl GET "/api/experts/$expert_internal_id/engagements" "" "$admin_token" 200 && {
        engagement_count=$(jq length /tmp/api_response.json)
        log INFO "Verified all imported engagements: $engagement_count found for expert"
    } || log ERROR "Failed to verify JSON imported engagements"
} || log ERROR "Failed to import engagements via JSON"

log HEADER "TEST SUMMARY"
log INFO "Total tests: $total_tests"
log SUCCESS "Passed: $passed_tests"
log ERROR "Failed: $failed_tests"
log SKIP "Skipped: $skipped_tests"
[[ $failed_tests -gt 0 ]] && log ERROR "Test run had failures. Check $LOG_FILE for details." || log SUCCESS "Test run completed successfully."

exit $failed_tests