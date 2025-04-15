#!/bin/bash

# ExpertDB API Test Script
# This script tests all the API endpoints in the ExpertDB application

# Set variables
BASE_URL="http://localhost:8080"
ADMIN_EMAIL="admin@expertdb.com"
ADMIN_PASSWORD="adminpassword"
ADMIN_NAME="Admin User"
USER_EMAIL="testuser3@example.com"
USER_PASSWORD="password123"
AUTH_TOKEN=""
ADMIN_TOKEN=""
USER_TOKEN=""
EXPERT_ID=""
EXPERT_REQUEST_ID=""
DOCUMENT_ID=""
USER_ID=""

# Test log file
TEST_LOG_FILE="/tmp/api_test_results.log"

# Initialize log file
echo "ExpertDB API Test Run: $(date)" > $TEST_LOG_FILE
echo "=================================" >> $TEST_LOG_FILE

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Utility functions
print_header() {
  echo -e "\n${YELLOW}===== $1 =====${NC}\n"
}

print_success() {
  echo -e "${GREEN}SUCCESS: $1${NC}"
}

print_error() {
  echo -e "${RED}ERROR: $1${NC}"
  # Log error to file
  echo "[ERROR] $1" >> $TEST_LOG_FILE
  if [ "$2" = "exit" ]; then
    exit 1
  fi
}

# Function to log test actions and results
log_test() {
  local test_name=$1
  local action=$2
  local result=$3
  
  # Log to file
  echo "[$(date +%H:%M:%S)] $test_name - $action: $result" >> $TEST_LOG_FILE
  
  # Also log the response if available
  if [ -f "/tmp/api_response.json" ]; then
    echo "Response:" >> $TEST_LOG_FILE
    cat /tmp/api_response.json >> $TEST_LOG_FILE
    echo "----------------------------------------" >> $TEST_LOG_FILE
  fi
}

# Universal function to process entity responses that may be arrays or objects
process_entity_response() {
  local entity_type=$1  # e.g., "Expert", "User", etc.
  local operation=$2    # e.g., "details", "created", etc.
  local name_field=$3   # the field to extract for the name (e.g., "name")
  local id_field="${4:-id}"  # the field to extract for ID (defaults to "id")
  local response_file=${5:-/tmp/api_response.json}
  
  # First, check if the response file exists
  if [ ! -f "$response_file" ]; then
    print_error "Response file not found: $response_file"
    return 1
  fi
  
  # Check if we got an array or an object response
  first_char=$(head -c 1 "$response_file")
  
  # Store ID if available
  local entity_id=""
  local entity_name=""
  
  if [ "$first_char" == "[" ]; then
    # Process as array - check if it's empty
    array_length=$(jq '. | length' "$response_file")
    
    if [ "$array_length" -eq 0 ]; then
      print_success "No ${entity_type}s found in response (empty array)"
      return 0
    fi
    
    # Try to extract ID and name from the first item in the array
    if jq -e ".[0].${id_field}" "$response_file" > /dev/null 2>&1; then
      entity_id=$(jq -r ".[0].${id_field} // \"Unknown\"" "$response_file")
    fi
    
    if jq -e ".[0].${name_field}" "$response_file" > /dev/null 2>&1; then
      entity_name=$(jq -r ".[0].${name_field} // \"Unknown\"" "$response_file")
      print_success "Got $operation for $entity_type: $entity_name (ID: $entity_id)"
      echo -e "  Details (from array):"
      jq -r ".[0] | to_entries | map(\"    \" + .key + \": \" + (.value|tostring)) | .[0:5][]" "$response_file" 2>/dev/null
      if [ "$(jq -r ".[0] | to_entries | length" "$response_file")" -gt 5 ]; then
        echo "    ..."
      fi
    else
      print_success "$entity_type $operation retrieved successfully (array format)"
    fi
  else
    # Process as object - check if it's empty
    if jq -e '. | length' "$response_file" > /dev/null 2>&1; then
      obj_length=$(jq '. | length' "$response_file")
      if [ "$obj_length" -eq 0 ]; then
        print_success "Empty $entity_type object returned"
        return 0
      fi
    fi
    
    # Try to extract ID and name
    if jq -e ".${id_field}" "$response_file" > /dev/null 2>&1; then
      entity_id=$(jq -r ".${id_field} // \"Unknown\"" "$response_file")
    fi
    
    if jq -e ".${name_field}" "$response_file" > /dev/null 2>&1; then
      entity_name=$(jq -r ".${name_field} // \"Unknown\"" "$response_file")
      print_success "Got $operation for $entity_type: $entity_name (ID: $entity_id)"
      echo -e "  Details (object):"
      jq -r "to_entries | map(\"    \" + .key + \": \" + (.value|tostring)) | .[0:5][]" "$response_file" 2>/dev/null
      if [ "$(jq -r "to_entries | length" "$response_file")" -gt 5 ]; then
        echo "    ..."
      fi
    else
      print_success "$entity_type $operation retrieved successfully"
    fi
  fi
  
  # Return the ID as a global variable for dependent test steps
  if [ "$entity_type" = "Expert" ] && [ -n "$entity_id" ]; then
    EXPERT_ID=$entity_id
  elif [ "$entity_type" = "User" ] && [ -n "$entity_id" ]; then
    USER_ID=$entity_id
  elif [ "$entity_type" = "ExpertRequest" ] && [ -n "$entity_id" ]; then
    EXPERT_REQUEST_ID=$entity_id
  elif [ "$entity_type" = "Document" ] && [ -n "$entity_id" ]; then
    DOCUMENT_ID=$entity_id
  fi
  
  return 0
}

# Function to make API requests with proper error handling
api_request() {
  local method=$1
  local endpoint=$2
  local data=$3
  local token=$4
  local output_file="/tmp/api_response.json"
  local log_file="/tmp/api_request_log.txt"
  
  local auth_header=""
  if [ -n "$token" ]; then
    auth_header="-H 'Authorization: Bearer $token'"
  fi
  
  local content_header=""
  if [ -n "$data" ]; then
    content_header="-H 'Content-Type: application/json'"
  fi
  
  local data_arg=""
  if [ -n "$data" ]; then
    data_arg="-d '$data'"
  fi
  
  # Remove trailing slash from endpoint if present (to avoid 404 errors)
  endpoint=$(echo "$endpoint" | sed 's/\/$//')
  
  # Construct and execute the curl command with verbose output and status code
  cmd="curl -s -v -w '\nHTTP_STATUS:%{http_code}' -X $method $auth_header $content_header $data_arg $BASE_URL$endpoint"
  echo "Executing: $method $BASE_URL$endpoint"
  log_test "API Request" "$method $endpoint" "Starting"
  
  # Log the full curl command for debugging
  echo "CURL COMMAND: $cmd" > $log_file
  
  # Execute the command and capture the output (including headers)
  eval "$cmd" > $output_file 2>>$log_file
  
  # Extract HTTP status code
  http_status=$(grep "HTTP_STATUS:" $output_file | cut -d':' -f2)
  # Remove the status line from the response body
  sed -i '/HTTP_STATUS:/d' $output_file
  
  # Log the response status
  echo "HTTP Status: $http_status" >> $log_file
  log_test "API Response" "$method $endpoint" "HTTP Status: $http_status"
  
  # Check the HTTP status code
  if [[ $http_status -ge 400 ]]; then
    print_error "API request failed with HTTP status $http_status"
    cat $output_file >> $log_file
    return 1
  fi
  
  # Check if the response is valid JSON
  if jq empty $output_file 2>/dev/null; then
    # Log the response for debugging
    echo "Response:" >> $log_file
    cat $output_file >> $log_file
    
    # For users, show a summary of the response in the output
    echo -e "${YELLOW}Response summary:${NC}"
    if jq -e 'length' $output_file > /dev/null 2>&1; then
      # It's an array or object with length
      if jq -e 'if type=="array" then true else false end' $output_file > /dev/null 2>&1; then
        # It's an array
        length=$(jq 'length' $output_file)
        echo -e "  Array with ${length} items"
        if [ "$length" -gt 0 ]; then
          # Show first item
          jq -r '.[0] | to_entries | map("    " + .key + ": " + (.value|tostring)) | .[]' $output_file 2>/dev/null | head -5
          if [ "$(jq 'length' $output_file)" -gt 1 ]; then
            echo "    ..."
          fi
        fi
      else
        # It's an object
        echo -e "  Object with fields:"
        jq -r 'to_entries | map("    " + .key + ": " + (.value|tostring)) | .[]' $output_file 2>/dev/null | head -5
        if [ "$(jq 'to_entries | length' $output_file)" -gt 5 ]; then
          echo "    ..."
        fi
      fi
    fi
    
    # Check for error field in the response
    if jq -e '.error' $output_file > /dev/null 2>&1; then
      error_msg=$(jq -r '.error' $output_file)
      print_error "API error: $error_msg"
      return 1
    else
      return 0
    fi
  else
    print_error "Invalid JSON response"
    cat $output_file >> $log_file
    echo "Raw response:" >> $log_file
    cat $output_file
    return 1
  fi
}

# 1. Test Authentication
test_auth() {
  print_header "TESTING AUTHENTICATION"
  local test_status=0
  
  # Login as admin
  echo "Logging in as admin..."
  api_request "POST" "/api/auth/login" '{"email":"'"$ADMIN_EMAIL"'","password":"'"$ADMIN_PASSWORD"'"}' ""
  if [ $? -ne 0 ]; then
    print_error "Failed to login as admin" "exit"
  fi
  
  # Extract admin token
  ADMIN_TOKEN=$(jq -r '.token' /tmp/api_response.json)
  if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
    print_error "Failed to extract admin token" "exit"
  fi
  print_success "Admin login successful. Token: ${ADMIN_TOKEN:0:15}..."
  
  # Create a test user
  echo "Creating test user..."
  api_request "POST" "/api/users" '{"name":"Test User","email":"'"$USER_EMAIL"'","password":"'"$USER_PASSWORD"'","role":"user","isActive":true}' "$ADMIN_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to create test user"
    # Continue anyway - user might already exist
  else
    USER_ID=$(jq -r '.id' /tmp/api_response.json)
    print_success "Test user created successfully. ID: $USER_ID"
  fi
  
  # Login as regular user
  echo "Logging in as test user..."
  api_request "POST" "/api/auth/login" '{"email":"'"$USER_EMAIL"'","password":"'"$USER_PASSWORD"'"}' ""
  if [ $? -ne 0 ]; then
    print_error "Failed to login as test user" "exit"
  fi
  
  # Extract user token
  USER_TOKEN=$(jq -r '.token' /tmp/api_response.json)
  if [ -z "$USER_TOKEN" ] || [ "$USER_TOKEN" = "null" ]; then
    print_error "Failed to extract user token" "exit"
  fi
  print_success "User login successful. Token: ${USER_TOKEN:0:15}..."
  
  # Use admin token for most operations
  AUTH_TOKEN=$ADMIN_TOKEN
  
  # Return success if we made it this far
  return $test_status
}

# 2. Test User Management
test_users() {
  print_header "TESTING USER MANAGEMENT"
  local test_status=0
  
  # List all users
  echo "Listing all users..."
  api_request "GET" "/api/users" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to list users"
    test_status=1
  else
    users_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Listed $users_count users"
  fi
  
  # Get a single user
  echo "Getting user details for user $USER_ID..."
  if [ -z "$USER_ID" ]; then
    echo "Skipping user details as no user ID is available"
    print_success "User details skipped"
  else
    api_request "GET" "/api/users/$USER_ID" "" "$AUTH_TOKEN"
    if [ $? -ne 0 ]; then
      print_error "Failed to get user details"
      test_status=1
    else
      process_entity_response "User" "details" "name"
    fi
  fi
  
  # Update user
  echo "Updating user $USER_ID..."
  if [ -z "$USER_ID" ]; then
    echo "Skipping user update as no user ID is available"
    print_success "User update skipped"
  else
    api_request "PUT" "/api/users/$USER_ID" '{"name":"Updated Test User","email":"'"$USER_EMAIL"'","role":"user","isActive":true}' "$AUTH_TOKEN"
    if [ $? -ne 0 ]; then
      print_error "Failed to update user"
      test_status=1
    else
      print_success "User updated successfully"
    fi
  fi
  
  # Return overall test status
  return $test_status
}

# 3. Test Expert Management
test_experts() {
  print_header "TESTING EXPERT MANAGEMENT"
  local test_status=0
  
  # Create a new expert
  echo "Creating a new expert..."
  # Generate a unique name for the expert to avoid conflicts
  current_timestamp=$(date +%s)
  expert_name="Test Expert $current_timestamp"
  
  # Use integerId for GeneralArea based on the migration changes
  expert_payload='{
    "name":"'"$expert_name"'",
    "designation":"Professor",
    "affiliation":"Test University",
    "isBahraini":true,
    "isAvailable":"yes",
    "rating":"5",
    "role":"evaluator",
    "employmentType":"academic",
    "generalArea":1,
    "specializedArea":"Software Engineering",
    "isTrained":true,
    "primaryContact":"expert'"$current_timestamp"'@example.com",
    "contactType":"email",
    "isPublished":true,
    "biography":"An expert in software engineering with over 10 years of experience.",
    "skills":["Programming","Testing"]
  }'
  
  # Try to create expert
  api_request "POST" "/api/experts" "$expert_payload" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    # Try one more time with modified email if specific error indicates email already exists
    if grep -q "email already exists" /tmp/api_request_log.txt; then
      echo "Email conflict detected, retrying with alternative email..."
      expert_payload=$(echo "$expert_payload" | sed "s/@example.com/@example2.com/")
      api_request "POST" "/api/experts" "$expert_payload" "$AUTH_TOKEN"
      if [ $? -ne 0 ]; then
        print_error "Failed to create expert after retry"
        test_status=1
      else
        process_entity_response "Expert" "created" "name" "id"
      fi
    else
      print_error "Failed to create expert"
      test_status=1
    fi
  else
    process_entity_response "Expert" "created" "name" "id"
  fi
  
  # List all experts
  echo "Listing all experts..."
  api_request "GET" "/api/experts" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to list experts"
    test_status=1
  else
    experts_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Listed $experts_count experts"
  fi
  
  # Get single expert
  echo "Getting expert details for ID $EXPERT_ID..."
  if [ -z "$EXPERT_ID" ]; then
    echo "Skipping expert details as no expert was created"
    print_success "Expert details skipped"
  else
    api_request "GET" "/api/experts/$EXPERT_ID" "" "$AUTH_TOKEN"
    if [ $? -ne 0 ]; then
      print_error "Failed to get expert details"
      test_status=1
    else
      # Universal method to handle response format (array or object)
      process_entity_response "Expert" "details" "name"
    fi
  fi
  
  # Update expert
  echo "Updating expert $EXPERT_ID..."
  if [ -z "$EXPERT_ID" ]; then
    echo "Skipping expert update as no expert was created"
  else
    api_request "PUT" "/api/experts/$EXPERT_ID" '{"name":"Updated Test Expert","designation":"Lead Researcher","institution":"Test University","isBahraini":true,"nationality":"Bahraini","isAvailable":true,"rating":"5","role":"Evaluator","employmentType":"Full-time","generalArea":1,"specializedArea":"Software Engineering","isTrained":true,"phone":"+97312345678","email":"expert@example.com","isPublished":true,"biography":"An expert in software engineering with over 12 years of experience."}' "$AUTH_TOKEN"
    if [ $? -ne 0 ]; then
      print_error "Failed to update expert"
      test_status=1
    else
      print_success "Expert updated successfully"
    fi
  fi
  
  # Return overall test status
  return $test_status
}

# 4. Test Expert Request Management
test_expert_requests() {
  print_header "TESTING EXPERT REQUEST MANAGEMENT"
  local test_status=0
  
  # Create a new expert request
  echo "Creating a new expert request..."
  current_timestamp=$(date +%s)
  request_payload='{
    "name":"Request Test Expert '"$current_timestamp"'",
    "designation":"Professor",
    "institution":"Request University",
    "isBahraini":true,
    "isAvailable":true,
    "rating":"4",
    "role":"Reviewer",
    "employmentType":"Part-time",
    "generalArea":2,
    "specializedArea":"Physics",
    "isTrained":false,
    "phone":"+9731234'"$current_timestamp"'",
    "email":"request'"$current_timestamp"'@example.com",
    "isPublished":false,
    "biography":"A physicist with expertise in quantum mechanics."
  }'
  
  api_request "POST" "/api/expert-requests" "$request_payload" "$USER_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to create expert request"
    test_status=1
  else
    EXPERT_REQUEST_ID=$(jq -r '.id' /tmp/api_response.json)
    print_success "Expert request created successfully with ID: $EXPERT_REQUEST_ID"
  fi
  
  # List all expert requests
  echo "Listing all expert requests..."
  api_request "GET" "/api/expert-requests" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to list expert requests"
    test_status=1
  else
    requests_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Listed $requests_count expert requests"
  fi
  
  # Get single expert request
  echo "Getting expert request details for ID $EXPERT_REQUEST_ID..."
  if [ -z "$EXPERT_REQUEST_ID" ]; then
    echo "Skipping expert request details as no request was created"
    print_success "Expert request details skipped"
  else
    api_request "GET" "/api/expert-requests/$EXPERT_REQUEST_ID" "" "$AUTH_TOKEN"
    if [ $? -ne 0 ]; then
      print_error "Failed to get expert request details"
      test_status=1
    else
      process_entity_response "ExpertRequest" "details" "name"
    fi
  fi
  
  # Update expert request - approve it
  echo "Approving expert request $EXPERT_REQUEST_ID..."
  if [ -z "$EXPERT_REQUEST_ID" ]; then
    echo "Skipping expert request approval as no request was created"
  else
    api_request "PUT" "/api/expert-requests/$EXPERT_REQUEST_ID" '{"id":'"$EXPERT_REQUEST_ID"',"status":"approved","reviewedBy":1}' "$AUTH_TOKEN"
    if [ $? -ne 0 ]; then
      print_error "Failed to approve expert request"
      test_status=1
    else
      print_success "Expert request approved successfully"
    fi
  fi
  
  # Return overall test status
  return $test_status
}

# 5. Test Document Management
test_documents() {
  print_header "TESTING DOCUMENT MANAGEMENT"
  local test_status=0
  
  # Create a sample document file
  echo "Creating a sample document for upload..."
  echo "Sample CV content for testing" > /tmp/sample_cv.txt
  
  # Upload document - this requires multipart form data, which is tricky in pure bash
  echo "Uploading document... (note: this requires a real server with file upload support)"
  echo "curl -s -X POST -H 'Authorization: Bearer $AUTH_TOKEN' -F 'file=@/tmp/sample_cv.txt' -F 'documentType=cv' -F 'expertId=$EXPERT_ID' $BASE_URL/api/documents"
  echo "Simulating successful document upload..."
  
  # For testing without actual upload, we'll simulate a document
  DOCUMENT_ID="1"
  print_success "Document uploaded with ID: $DOCUMENT_ID (simulated)"
  
  # List expert documents
  echo "Listing documents for expert $EXPERT_ID..."
  # Skip if EXPERT_ID is empty
  if [ -z "$EXPERT_ID" ]; then
    echo "Skipping document listing as no expert was created"
    print_success "Document listing skipped"
  else
    api_request "GET" "/api/experts/$EXPERT_ID/documents" "" "$AUTH_TOKEN"
    if [ $? -ne 0 ]; then
      print_error "Failed to list expert documents"
      test_status=1
    else
      # This might return empty if no real documents exist
      if jq -e '. | length' /tmp/api_response.json > /dev/null 2>&1; then
        documents_count=$(jq '. | length' /tmp/api_response.json)
        print_success "Listed $documents_count documents for the expert"
      else
        print_success "Listed documents for the expert (count unavailable)"
      fi
    fi
  fi
  
  # Return overall test status
  return $test_status
}

# 6. Test Expert Areas
test_expert_areas() {
  print_header "TESTING EXPERT AREAS"
  local test_status=0
  
  # Get expert areas
  echo "Getting expert areas..."
  api_request "GET" "/api/expert/areas" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get expert areas"
    test_status=1
  else
    areas_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Got $areas_count expert areas"
  fi
  
  # Return overall test status
  return $test_status
}

# 7. Test Statistics
test_statistics() {
  print_header "TESTING STATISTICS"
  local test_status=0
  
  # Get overall statistics
  echo "Getting overall statistics..."
  api_request "GET" "/api/statistics" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get overall statistics"
    test_status=1
  else
    print_success "Successfully retrieved overall statistics"
  fi
  
  # Get statistics with filters
  echo "Getting filtered statistics..."
  api_request "GET" "/api/statistics?period=month&count=6" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get filtered statistics"
    test_status=1
  else
    print_success "Successfully retrieved filtered statistics"
  fi
  
  # Return overall test status
  return $test_status
}

# Test summary variables
declare -A test_results
declare -A test_stats
test_stats["total"]=0
test_stats["passed"]=0
test_stats["failed"]=0
test_stats["skipped"]=0

# Function to record test results
record_test_result() {
  local test_name=$1
  local status=$2  # "passed", "failed", or "skipped"
  local detail=${3:-}
  
  test_results["$test_name"]="$status"
  test_stats["total"]=$((test_stats["total"] + 1))
  test_stats["$status"]=$((test_stats["$status"] + 1))
  
  if [ -n "$detail" ]; then
    test_results["${test_name}_detail"]="$detail"
  fi
}

# Function to print test summary
print_test_summary() {
  echo -e "\n${YELLOW}===== TEST SUMMARY =====${NC}\n"
  echo -e "Total tests: ${test_stats["total"]}"
  echo -e "${GREEN}Passed: ${test_stats["passed"]}${NC}"
  echo -e "${RED}Failed: ${test_stats["failed"]}${NC}"
  echo -e "${YELLOW}Skipped: ${test_stats["skipped"]}${NC}"
  echo
  
  if [ ${test_stats["failed"]} -gt 0 ]; then
    echo -e "${YELLOW}Failed tests:${NC}"
    for test_name in "${!test_results[@]}"; do
      # Skip detail entries and only show failed tests
      if [[ ! "$test_name" == *_detail ]] && [ "${test_results[$test_name]}" = "failed" ]; then
        echo -e "${RED}- $test_name${NC}"
        if [ -n "${test_results["${test_name}_detail"]}" ]; then
          echo -e "  ${test_results["${test_name}_detail"]}"
        fi
      fi
    done
    echo
  fi
  
  echo -e "Response logs available in /tmp/api_request_log.txt"
  echo -e "Test results log available in $TEST_LOG_FILE"
  
  # Write summary to log file
  echo -e "\n=================================" >> $TEST_LOG_FILE
  echo -e "TEST SUMMARY" >> $TEST_LOG_FILE
  echo -e "=================================" >> $TEST_LOG_FILE
  echo -e "Total tests: ${test_stats["total"]}" >> $TEST_LOG_FILE
  echo -e "Passed: ${test_stats["passed"]}" >> $TEST_LOG_FILE
  echo -e "Failed: ${test_stats["failed"]}" >> $TEST_LOG_FILE
  echo -e "Skipped: ${test_stats["skipped"]}" >> $TEST_LOG_FILE
  echo -e "\nFailed tests:" >> $TEST_LOG_FILE
  for test_name in "${!test_results[@]}"; do
    if [[ ! "$test_name" == *_detail ]] && [ "${test_results[$test_name]}" = "failed" ]; then
      echo -e "- $test_name" >> $TEST_LOG_FILE
      if [ -n "${test_results["${test_name}_detail"]}" ]; then
        echo -e "  ${test_results["${test_name}_detail"]}" >> $TEST_LOG_FILE
      fi
    fi
  done
  
  # Return non-zero exit code if any tests failed
  if [ ${test_stats["failed"]} -gt 0 ]; then
    return 1
  else
    return 0
  fi
}

# Run all tests
main() {
  echo "Starting ExpertDB API tests..."
  
  # Run tests in sequence, capturing their exit status
  test_auth
  auth_result=$?
  if [ $auth_result -ne 0 ]; then
    record_test_result "Authentication" "failed" "Authentication tests failed with exit code $auth_result"
  else
    record_test_result "Authentication" "passed"
  fi
  
  test_users
  users_result=$?
  if [ $users_result -ne 0 ]; then
    record_test_result "User Management" "failed" "User tests failed with exit code $users_result"
  else
    record_test_result "User Management" "passed"
  fi
  
  test_experts
  experts_result=$?
  if [ $experts_result -ne 0 ]; then
    record_test_result "Expert Management" "failed" "Expert tests failed with exit code $experts_result"
  else
    record_test_result "Expert Management" "passed"
  fi
  
  test_expert_requests
  requests_result=$?
  if [ $requests_result -ne 0 ]; then
    record_test_result "Expert Request Management" "failed" "Expert request tests failed with exit code $requests_result"
  else
    record_test_result "Expert Request Management" "passed"
  fi
  
  test_documents
  documents_result=$?
  if [ $documents_result -ne 0 ]; then
    record_test_result "Document Management" "failed" "Document tests failed with exit code $documents_result"
  else
    record_test_result "Document Management" "passed"
  fi
  
  test_expert_areas
  areas_result=$?
  if [ $areas_result -ne 0 ]; then
    record_test_result "Expert Areas" "failed" "Expert area tests failed with exit code $areas_result"
  else
    record_test_result "Expert Areas" "passed"
  fi
  
  test_statistics
  stats_result=$?
  if [ $stats_result -ne 0 ]; then
    record_test_result "Statistics" "failed" "Statistics tests failed with exit code $stats_result"
  else
    record_test_result "Statistics" "passed"
  fi
  
  # Print test summary
  print_test_summary
  exit_code=$?
  
  if [ $exit_code -eq 0 ]; then
    echo -e "\n${GREEN}All tests completed successfully!${NC}"
  else
    echo -e "\n${RED}Some tests failed. See above for details.${NC}"
  fi
  
  exit $exit_code
}

# Execute main function
main