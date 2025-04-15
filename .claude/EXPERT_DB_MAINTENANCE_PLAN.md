# ExpertDB Maintenance Plan

## Current Issues

Based on the test logs from `/home/alikebrahim/dev/expertdb/RUN.md`, we've identified several recurring issues that need to be addressed systematically.

### 1. NULL Value Handling in Database Operations

**Issues:**
- ✅ Fixed: `updated_at` NULL handling in `list_experts.go`
- ✅ Fixed: `is_available` NULL handling in `list_experts.go`
- ✅ Fixed: `specialized_area` NULL handling in `list_experts.go`
- ✅ Fixed: `specialized_area` and `biography` NULL handling in `expert_operations.go`
- ✅ Fixed: `last_login` NULL handling in `user_storage.go` (added sql.NullTime)
- ✅ Fixed: timestamp handling in `expert_request_operations.go` to check for empty strings
- ❓ Potential: Other database fields that might contain NULL values

**Root Cause:**
Go's SQL package doesn't support scanning NULL values directly into standard Go types. We need to use appropriate `sql.Null*` types for all potentially NULL fields.

### 2. Data Type Conversion in JSON Requests

**Issues:**
- ✅ Fixed: `generalArea` string-to-int64 conversion in expert request approval
- ❓ Potential: Similar issues in other API endpoints

**Root Cause:**
JSON doesn't distinguish between number types, and numeric fields are often encoded as strings in client requests. We need consistent type handling for all fields with specific type requirements.

### 3. Database Constraints

**Issues:**
- ✅ Fixed: `UNIQUE constraint failed: experts.expert_id` when creating experts
  - Added uniqueness check before creating experts
  - Implemented better error handling for constraint violations
  - Added mechanism to generate unique expert_ids
- ❓ Potential: Other constraint violations may occur

**Root Cause:**
The application wasn't handling database constraints appropriately, leading to 500 errors instead of meaningful client responses.

### 4. Test Script Issues

**Issues:**
- ✅ Fixed: `Cannot index array with string "name"` when getting user details
  - Improved test script to handle both array and object responses
  - Added detection of response format (array or object)
  - Added fallback values when fields don't exist
  - Added better skipping logic when IDs aren't available
- ❓ Potential: Other test script failures

**Root Cause:**
The test script expected certain response formats but received different ones, particularly with array vs. object responses.

## Action Plan

### 1. Review All Database Field Scanning Operations

**Approach:**
1. Identify all functions that retrieve data from the database (SQL query + scan operations)
2. For each string, bool, or similar non-nullable type in Go, check if the corresponding database field could be NULL
3. Replace with appropriate `sql.Null*` types (NullString, NullBool, NullInt64, etc.)
4. Add proper handling logic after scanning

**Files to Review:**
- `list_experts.go` - Already fixed some fields but needs specialized_area fixed
- `expert_operations.go` - Partially fixed, needs complete review
- `expert_request_operations.go` - May need similar NULL handling
- Any other files with direct database interaction

**Priority: High** - This is causing immediate failures in the API

### 2. Implement Consistent Type Conversion for API Requests

**Approach:**
1. Review all API handlers that accept JSON payloads
2. For fields that require specific types (especially numeric fields), ensure proper conversion
3. Standardize error handling for type conversion failures
4. Add validation to check data types before processing

**Files to Review:**
- `api.go` - Review all JSON handling functions
- `expert_request_operations.go` - Check type conversions
- `test_api.sh` - Ensure all JSON payloads use proper types (no quotes around numbers)

**Priority: Medium** - Fixed for some endpoints but may affect others

### 3. Enhance Error Handling for Database Constraints

**Approach:**
1. Identify common database constraints (UNIQUE, NOT NULL, FOREIGN KEY)
2. Add explicit error handling for constraint violations
3. Return meaningful error messages to API clients
4. For unique constraints, add checks before insertion attempts

**Implementation Tasks:**
- Add a function to check if an expert with a given expert_id already exists
- Return appropriate error codes (409 Conflict) rather than 500 for constraint violations
- Add logic to generate unique identifiers when needed

**Priority: Medium** - Affects user experience with misleading error messages

### 4. Fix Test Script Issues

**Approach:**
1. Fix the user details retrieval - it's trying to access a single user but getting an array
2. Update error handling in the test script to provide more context
3. Update reporting to display sent/received payloads.
4. Add proper cleanup between test runs to avoid unique constraint violations

**Implementation Tasks:**
- Fix the URL for getting user details (use proper ID)
- Enhance jq handling for arrays vs. objects
- Consider adding a cleanup/reset step at the beginning of tests

**Priority: medium** - Not blocking functionality but affects test reliability

## Implementation Order

1. ✅ Fix specialized_area NULL handling in list_experts.go
2. ✅ Fix NULL handling in expert_operations.go and user_storage.go 
3. ✅ Improve timestamp handling in expert_request_operations.go
4. ✅ Develop a mechanism for uniqueness constraint handling, especially for expert_id
   - Added uniqueness check before creating experts
   - Added unique ID generation with timestamp and random components
   - Improved error handling for constraint violations
5. ✅ Update test_api.sh to handle array vs. object responses correctly
   - Added format detection (array vs. object)
   - Improved error handling and skipping logic
   - Added fallback values for missing fields

## Notes

### Progress Update (2025-04-15)

We've systematically fixed NULL handling issues across the backend:

1. Fixed `specialized_area` NULL handling in `list_experts.go` by using sql.NullString
2. Fixed `specialized_area` and `biography` NULL handling in `expert_operations.go` using sql.NullString
3. Fixed `last_login` NULL handling in `user_storage.go` using sql.NullTime for all user-related functions
4. Improved timestamp parsing in `expert_request_operations.go` by checking for empty strings

We've also addressed the uniqueness constraint issue for expert_id:

1. Added `ExpertIDExists()` function to check if an expert_id already exists
2. Implemented `GenerateUniqueExpertID()` to create unique IDs with timestamp and random components
3. Updated expert creation to verify uniqueness before insertion
4. Improved error handling for constraint violations to return appropriate HTTP status codes (409 Conflict)
5. Applied these changes to both direct expert creation and expert creation via request approval

We've also fixed the test script to handle response format inconsistencies:

1. Added detection of response format (array vs. object)
2. Improved error handling for both formats
3. Added checks for empty or missing IDs
4. Added fallback values when expected fields are missing
5. Fixed specific issues with user, expert, and expert request detail retrieval

### Progress Update (2025-04-16)

Major improvements to the test script with advanced response handling:

1. Created universal `process_entity_response` function to consistently handle both array and object responses
2. Added detailed response summary showing request details and response data
3. Added sophisticated error handling and reporting for API responses
4. Implemented proper test status tracking and reporting
5. Added comprehensive test summary with pass/fail counts
6. Added detailed logging to `/tmp/api_test_results.log` for post-test analysis
7. Added timestamp-based IDs to test data to avoid conflicts

Also fixed general_area field handling:

1. Updated area field handling in `list_experts.go` to properly query INTEGER fields
2. Improved type conversion for general_area to ensure proper integer values
3. Fixed LIKE queries on general_area to use proper equality comparison for integers

Next steps:
1. Perform comprehensive testing with the updated code
2. Address any remaining edge cases or issues discovered during testing
3. Consider adding additional validation and error handling improvements
