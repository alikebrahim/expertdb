#!/bin/bash

# Quick test script for expert filtering capabilities

# Configuration
API_BASE_URL="http://localhost:8080"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Login to get token
echo -e "${BLUE}Getting authentication token...${NC}"
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"email":"admin@expertdb.com","password":"adminpassword"}' $API_BASE_URL/api/auth/login | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo -e "${RED}Failed to get authentication token. Exiting.${NC}"
  exit 1
fi

echo -e "${GREEN}Successfully obtained token${NC}"

# Test each filter type
echo -e "\n${BLUE}Testing all expert filters:${NC}"

# Array of filter tests
declare -a TESTS=(
  "by_nationality=Bahraini&limit=5:Testing nationality filter"
  "by_general_area=1&limit=5:Testing general area filter"
  "by_specialized_area=Software&limit=5:Testing specialized area filter (partial match)"
  "by_employment_type=academic&limit=5:Testing employment type filter"
  "by_role=evaluator&limit=5:Testing role filter"
  "by_nationality=Bahraini&by_employment_type=academic&limit=5:Testing combined filters"
  "sort_by=name&sort_order=asc&limit=5:Testing sort by name ascending"
  "sort_by=institution&sort_order=desc&limit=5:Testing sort by institution descending"
  "sort_by=rating&limit=5:Testing sort by rating"
  "sort_by=created_at&sort_order=desc&limit=5:Testing sort by creation date (newest first)"
  "limit=5&offset=5:Testing pagination (page 2)"
)

# Run each test
for test in "${TESTS[@]}"; do
  # Split the test into query and description
  IFS=':' read -r query description <<< "$test"
  
  echo -e "\n${YELLOW}$description${NC}"
  echo -e "Query: $query"
  
  # Execute the test
  response=$(curl -s -H "Authorization: Bearer $TOKEN" "$API_BASE_URL/api/experts?$query")
  
  # Check if response is valid JSON
  if echo "$response" | jq . >/dev/null 2>&1; then
    # Check if we're using the new response format with pagination
    if echo "$response" | jq '.experts' >/dev/null 2>&1; then
      # Extract count and pagination info from new format
      count=$(echo "$response" | jq '.experts | length')
      total_count=$(echo "$response" | jq '.pagination.totalCount')
      current_page=$(echo "$response" | jq '.pagination.currentPage')
      total_pages=$(echo "$response" | jq '.pagination.totalPages')
      
      # Show results with pagination info
      echo -e "${GREEN}Request successful: Found $count experts (page $current_page of $total_pages, total: $total_count)${NC}"
      
      # Show pagination metadata
      echo -e "\nPagination metadata:"
      echo "$response" | jq '.pagination'
      
      # Show first result if any
      if [ "$count" -gt 0 ]; then
        echo -e "\nSample result:"
        echo "$response" | jq '.experts[0] | {id, name, nationality, role, employmentType, generalArea, generalAreaName, specializedArea}'
      fi
    else
      # Old format without pagination
      count=$(echo "$response" | jq '. | length')
      
      # Show results
      echo -e "${GREEN}Request successful: Found $count experts${NC}"
      
      # Show first result if any
      if [ "$count" -gt 0 ]; then
        echo -e "\nSample result:"
        echo "$response" | jq '.[0] | {id, name, nationality, role, employmentType, generalArea, generalAreaName, specializedArea}'
      fi
    fi
  else
    echo -e "${RED}Request failed:${NC}"
    echo "$response"
  fi
done

echo -e "\n${GREEN}Filter testing complete${NC}"