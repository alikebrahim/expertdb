#!/bin/bash

# Test script for integration between Phase Planning and Engagements
# This script tests the complete workflow from phase creation to engagement creation

set -e  # Exit on error
SERVER_URL="http://localhost:8080"
AUTH_TOKEN=""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Testing Phase Planning and Engagement Integration ===${NC}"

# Login as admin
echo -e "${BLUE}Logging in as admin...${NC}"
AUTH_TOKEN=$(curl -s -X POST ${SERVER_URL}/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@expertdb.com", "password": "adminpassword"}' | jq -r .token)

if [ -z "$AUTH_TOKEN" ]; then
  echo -e "${RED}Failed to login${NC}"
  exit 1
fi

echo -e "${GREEN}Logged in successfully${NC}"

# Get a valid planner ID
echo -e "${BLUE}Getting planner ID...${NC}"
PLANNER_ID=$(curl -s -X GET ${SERVER_URL}/api/users?role=planner \
  -H "Authorization: Bearer ${AUTH_TOKEN}" | jq -r '.[0].id')

if [ -z "$PLANNER_ID" ]; then
  echo -e "${BLUE}No existing planner, creating one...${NC}"
  PLANNER_ID=$(curl -s -X POST ${SERVER_URL}/api/users \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${AUTH_TOKEN}" \
    -d '{"name": "Test Planner", "email": "planner@expertdb.com", "password": "password123", "role": "planner"}' | jq -r .id)
fi

echo -e "${GREEN}Using planner ID: ${PLANNER_ID}${NC}"

# Create a new phase with both QP and IL applications
echo -e "${BLUE}Creating new phase with QP and IL applications...${NC}"
PHASE_RESPONSE=$(curl -s -X POST ${SERVER_URL}/api/phases \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${AUTH_TOKEN}" \
  -d '{
    "title": "Test Phase for QP and IL",
    "assignedPlannerId": '"$PLANNER_ID"',
    "status": "in_progress",
    "applications": [
      {
        "type": "QP",
        "institutionName": "Test University",
        "qualificationName": "Bachelor of Science"
      },
      {
        "type": "IL",
        "institutionName": "Test College",
        "qualificationName": "Diploma in Engineering"
      }
    ]
  }')

PHASE_ID=$(echo $PHASE_RESPONSE | jq -r .id)
echo -e "${GREEN}Created phase with ID: ${PHASE_ID}${NC}"

# Get application IDs
echo -e "${BLUE}Getting application IDs...${NC}"
APPLICATIONS=$(curl -s -X GET ${SERVER_URL}/api/phases/${PHASE_ID} \
  -H "Authorization: Bearer ${AUTH_TOKEN}" | jq -r .applications)

QP_APP_ID=$(echo $APPLICATIONS | jq -r '.[0].id')
IL_APP_ID=$(echo $APPLICATIONS | jq -r '.[1].id')

echo -e "${GREEN}QP Application ID: ${QP_APP_ID}${NC}"
echo -e "${GREEN}IL Application ID: ${IL_APP_ID}${NC}"

# Get some expert IDs
echo -e "${BLUE}Getting expert IDs...${NC}"
EXPERTS=$(curl -s -X GET ${SERVER_URL}/api/experts?limit=4 \
  -H "Authorization: Bearer ${AUTH_TOKEN}" | jq -r .experts)

EXPERT1_ID=$(echo $EXPERTS | jq -r '.[0].id')
EXPERT2_ID=$(echo $EXPERTS | jq -r '.[1].id')
EXPERT3_ID=$(echo $EXPERTS | jq -r '.[2].id')
EXPERT4_ID=$(echo $EXPERTS | jq -r '.[3].id')

if [ -z "$EXPERT1_ID" ] || [ -z "$EXPERT2_ID" ] || [ -z "$EXPERT3_ID" ] || [ -z "$EXPERT4_ID" ]; then
  echo -e "${RED}Not enough experts in the database. Please add at least 4 experts.${NC}"
  exit 1
fi

echo -e "${GREEN}Using experts: ${EXPERT1_ID}, ${EXPERT2_ID}, ${EXPERT3_ID}, ${EXPERT4_ID}${NC}"

# Assign experts to applications (simulating planner's role)
echo -e "${BLUE}Assigning experts to QP application...${NC}"
curl -s -X PUT ${SERVER_URL}/api/phases/${PHASE_ID}/applications/${QP_APP_ID} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${AUTH_TOKEN}" \
  -d '{
    "expert1": '"$EXPERT1_ID"',
    "expert2": '"$EXPERT2_ID"'
  }'

echo -e "${BLUE}Assigning experts to IL application...${NC}"
curl -s -X PUT ${SERVER_URL}/api/phases/${PHASE_ID}/applications/${IL_APP_ID} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${AUTH_TOKEN}" \
  -d '{
    "expert1": '"$EXPERT3_ID"',
    "expert2": '"$EXPERT4_ID"'
  }'

# Admin reviews and approves applications
echo -e "${BLUE}Admin approving QP application...${NC}"
curl -s -X PUT ${SERVER_URL}/api/phases/${PHASE_ID}/applications/${QP_APP_ID}/review \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${AUTH_TOKEN}" \
  -d '{
    "action": "approve"
  }'

echo -e "${BLUE}Admin approving IL application...${NC}"
curl -s -X PUT ${SERVER_URL}/api/phases/${PHASE_ID}/applications/${IL_APP_ID}/review \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${AUTH_TOKEN}" \
  -d '{
    "action": "approve"
  }'

# Verify engagements were created correctly
echo -e "${BLUE}Verifying engagements for Expert 1 (should be validator/QP)...${NC}"
ENGAGEMENTS=$(curl -s -X GET ${SERVER_URL}/api/expert-engagements?expertId=${EXPERT1_ID} \
  -H "Authorization: Bearer ${AUTH_TOKEN}")

ENGAGEMENT_TYPE=$(echo $ENGAGEMENTS | jq -r '.engagements[0].engagementType')
if [ "$ENGAGEMENT_TYPE" == "validator" ]; then
  echo -e "${GREEN}✅ Expert 1 correctly assigned as validator (QP)${NC}"
else
  echo -e "${RED}❌ Expert 1 not assigned as validator: $ENGAGEMENT_TYPE${NC}"
fi

echo -e "${BLUE}Verifying engagements for Expert 3 (should be evaluator/IL)...${NC}"
ENGAGEMENTS=$(curl -s -X GET ${SERVER_URL}/api/expert-engagements?expertId=${EXPERT3_ID} \
  -H "Authorization: Bearer ${AUTH_TOKEN}")

ENGAGEMENT_TYPE=$(echo $ENGAGEMENTS | jq -r '.engagements[0].engagementType')
if [ "$ENGAGEMENT_TYPE" == "evaluator" ]; then
  echo -e "${GREEN}✅ Expert 3 correctly assigned as evaluator (IL)${NC}"
else
  echo -e "${RED}❌ Expert 3 not assigned as evaluator: $ENGAGEMENT_TYPE${NC}"
fi

# Check statistics
echo -e "${BLUE}Checking engagement statistics for QP/IL mapping...${NC}"
STATS=$(curl -s -X GET ${SERVER_URL}/api/statistics/engagements \
  -H "Authorization: Bearer ${AUTH_TOKEN}")

echo -e "${GREEN}Engagement statistics: $(echo $STATS | jq -r .)${NC}"

echo -e "${GREEN}=== Phase Planning and Engagement Integration Tests Completed ===${NC}"