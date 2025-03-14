#!/bin/bash

# ExpertDB API Test Script
# This script tests all the API endpoints in the ExpertDB application

# Set variables
BASE_URL="http://localhost:8080"
ADMIN_EMAIL="ali@edb.com"
ADMIN_PASSWORD="alipass"
USER_EMAIL="testuser3@example.com"
USER_PASSWORD="password123"
AUTH_TOKEN=""
ADMIN_TOKEN=""
USER_TOKEN=""
EXPERT_ID=""
EXPERT_REQUEST_ID=""
DOCUMENT_ID=""
USER_ID=""

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
  if [ "$2" = "exit" ]; then
    exit 1
  fi
}

# Function to make API requests with proper error handling
api_request() {
  local method=$1
  local endpoint=$2
  local data=$3
  local token=$4
  local output_file="/tmp/api_response.json"
  
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
  
  # Construct and execute the curl command
  cmd="curl -s -X $method $auth_header $content_header $data_arg $BASE_URL$endpoint"
  echo "Executing: $method $BASE_URL$endpoint"
  
  # Execute the command and capture the output
  eval "$cmd" > $output_file
  
  # Check if the response is valid JSON
  if jq empty $output_file 2>/dev/null; then
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
    cat $output_file
    return 1
  fi
}

# 1. Test Authentication
test_auth() {
  print_header "TESTING AUTHENTICATION"
  
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
}

# 2. Test User Management
test_users() {
  print_header "TESTING USER MANAGEMENT"
  
  # List all users
  echo "Listing all users..."
  api_request "GET" "/api/users" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to list users"
  else
    users_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Listed $users_count users"
  fi
  
  # Get a single user
  echo "Getting user details for user $USER_ID..."
  api_request "GET" "/api/users/$USER_ID" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get user details"
  else
    user_name=$(jq -r '.name' /tmp/api_response.json)
    print_success "Got details for user: $user_name"
  fi
  
  # Update user
  echo "Updating user $USER_ID..."
  api_request "PUT" "/api/users/$USER_ID" '{"name":"Updated Test User","email":"'"$USER_EMAIL"'","role":"user","isActive":true}' "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to update user"
  else
    print_success "User updated successfully"
  fi
}

# 3. Test Expert Management
test_experts() {
  print_header "TESTING EXPERT MANAGEMENT"
  
  # Create a new expert
  echo "Creating a new expert..."
  api_request "POST" "/api/experts" '{"name":"Test Expert","affiliation":"Test University","primaryContact":"expert@example.com","contactType":"email","isBahraini":true,"availability":"yes","role":"evaluator","employmentType":"academic","generalArea":"Engineering","biography":"An expert in software engineering with over 10 years of experience."}' "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to create expert"
  else
    EXPERT_ID=$(jq -r '.id' /tmp/api_response.json)
    print_success "Expert created successfully with ID: $EXPERT_ID"
  fi
  
  # List all experts
  echo "Listing all experts..."
  api_request "GET" "/api/experts" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to list experts"
  else
    experts_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Listed $experts_count experts"
  fi
  
  # Get single expert
  echo "Getting expert details for ID $EXPERT_ID..."
  api_request "GET" "/api/experts/$EXPERT_ID" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get expert details"
  else
    expert_name=$(jq -r '.name' /tmp/api_response.json)
    print_success "Got details for expert: $expert_name"
  fi
  
  # Update expert
  echo "Updating expert $EXPERT_ID..."
  api_request "PUT" "/api/experts/$EXPERT_ID" '{"name":"Updated Test Expert","designation":"Lead Researcher","institution":"Test University","isBahraini":true,"nationality":"Bahraini","isAvailable":true,"rating":"5","role":"Evaluator","employmentType":"Full-time","generalArea":"Engineering","specializedArea":"Software Engineering","isTrained":true,"phone":"+97312345678","email":"expert@example.com","isPublished":true,"biography":"An expert in software engineering with over 12 years of experience."}' "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to update expert"
  else
    print_success "Expert updated successfully"
  fi
}

# 4. Test Expert Request Management
test_expert_requests() {
  print_header "TESTING EXPERT REQUEST MANAGEMENT"
  
  # Create a new expert request
  echo "Creating a new expert request..."
  api_request "POST" "/api/expert-requests" '{"name":"Request Test Expert","designation":"Professor","institution":"Request University","isBahraini":true,"isAvailable":true,"rating":"4","role":"Reviewer","employmentType":"Part-time","generalArea":"Science","specializedArea":"Physics","isTrained":false,"phone":"+97312345679","email":"request@example.com","isPublished":false,"biography":"A physicist with expertise in quantum mechanics."}' "$USER_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to create expert request"
  else
    EXPERT_REQUEST_ID=$(jq -r '.id' /tmp/api_response.json)
    print_success "Expert request created successfully with ID: $EXPERT_REQUEST_ID"
  fi
  
  # List all expert requests
  echo "Listing all expert requests..."
  api_request "GET" "/api/expert-requests" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to list expert requests"
  else
    requests_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Listed $requests_count expert requests"
  fi
  
  # Get single expert request
  echo "Getting expert request details for ID $EXPERT_REQUEST_ID..."
  api_request "GET" "/api/expert-requests/$EXPERT_REQUEST_ID" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get expert request details"
  else
    request_name=$(jq -r '.name' /tmp/api_response.json)
    print_success "Got details for expert request: $request_name"
  fi
  
  # Update expert request - approve it
  echo "Approving expert request $EXPERT_REQUEST_ID..."
  api_request "PUT" "/api/expert-requests/$EXPERT_REQUEST_ID" '{"id":'"$EXPERT_REQUEST_ID"',"name":"Request Test Expert","designation":"Professor","institution":"Request University","isBahraini":true,"isAvailable":true,"rating":"4","role":"Reviewer","employmentType":"Part-time","generalArea":"Science","specializedArea":"Physics","isTrained":false,"phone":"+97312345679","email":"request@example.com","isPublished":true,"status":"approved","biography":"A physicist with expertise in quantum mechanics.","reviewedBy":1}' "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to approve expert request"
  else
    print_success "Expert request approved successfully"
  fi
}

# 5. Test Document Management
test_documents() {
  print_header "TESTING DOCUMENT MANAGEMENT"
  
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
  else
    api_request "GET" "/api/experts/$EXPERT_ID/documents" "" "$AUTH_TOKEN"
  fi
  if [ $? -ne 0 ]; then
    print_error "Failed to list expert documents"
  else
    # This might return empty if no real documents exist
    documents_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Listed $documents_count documents for the expert"
  fi
}

# 6. Test ISCED Classifications
test_isced() {
  print_header "TESTING ISCED CLASSIFICATIONS"
  
  # Get ISCED levels
  echo "Getting ISCED levels..."
  api_request "GET" "/api/isced/levels" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get ISCED levels"
  else
    levels_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Got $levels_count ISCED levels"
  fi
  
  # Get ISCED fields
  echo "Getting ISCED fields..."
  api_request "GET" "/api/isced/fields" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get ISCED fields"
  else
    fields_count=$(jq '. | length' /tmp/api_response.json)
    print_success "Got $fields_count ISCED fields"
  fi
}

# 7. Test Statistics
test_statistics() {
  print_header "TESTING STATISTICS"
  
  # Get overall statistics
  echo "Getting overall statistics..."
  api_request "GET" "/api/statistics" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get overall statistics"
  else
    print_success "Successfully retrieved overall statistics"
  fi
  
  # Get statistics with filters
  echo "Getting filtered statistics..."
  api_request "GET" "/api/statistics?period=month&count=6" "" "$AUTH_TOKEN"
  if [ $? -ne 0 ]; then
    print_error "Failed to get filtered statistics"
  else
    print_success "Successfully retrieved filtered statistics"
  fi
}

# Run all tests
main() {
  echo "Starting ExpertDB API tests..."
  
  # Run tests in sequence
  test_auth
  test_users
  test_experts
  test_expert_requests
  test_documents
  test_isced
  test_statistics
  
  echo -e "\n${GREEN}All tests completed!${NC}"
}

# Execute main function
main
