# Instructions for Refactoring Backend Code

These guidelines are designed to help you refactor the provided backend code in `backend/ai_service.go` and `backend/api.go` (and other relevant files) to enhance organization, readability, and maintainability. The focus is on adding clear, descriptive comments that outline the purpose and processing steps within each function, handler, or method, ensuring that every code block is explicitly linked to its role in the processing flow. Follow these steps carefully to meet the user's requirements.

## Objective
Refactor the backend code to:
- Improve readability by adding detailed, step-by-step comments.
- Organize code blocks logically under descriptive headers that reflect their purpose in the processing sequence.
- Ensure maintainability by making the intent and flow of each function/handler/method immediately clear to developers.
- Remove all AI integration implementations and related routes.
- Add developer notes using the format `// NOTE: MESSAGE`.

## Guidelines

### 1. Remove AI Integration
#### Objective
- Remove all AI-related functionality, as the application will no longer integrate with an AI service.
#### Steps
- **Delete AI Service File**: Remove `ai_service.go` entirely, as it contains the `AIService` struct and methods like `GenerateProfile`, `SuggestISCED`, `ExtractSkills`, and `SuggestExpertPanel`, which are no longer needed.
- **Remove AI Routes**: In `api.go`, remove the following AI-related routes from the `Run` method in `APIServer`:
  - `POST /api/ai/generate-profile`
  - `POST /api/ai/suggest-isced`
  - `POST /api/ai/extract-skills`
  - `POST /api/ai/suggest-panel`
- **Remove AI Handlers**: Delete the corresponding handler methods in `api.go`:
  - `handleGenerateProfile`
  - `handleSuggestISCED`
  - `handleExtractSkills`
  - `handleSuggestPanel`
- **Remove AI Storage Methods**: In `storage.go`, remove the AI-related methods from the `Storage` interface and their implementations in `SQLiteStore`:
  - `StoreAIAnalysisResult`
  - `SuggestISCED`
  - `ExtractSkills`
  - `GenerateProfile`
  - `SuggestExpertPanel`
- **Update Database Schema**: In `dbschema`, remove the `ai_analysis_results` table and its indexes, as it’s no longer needed:
  - Drop the table `ai_analysis_results`.
  - Drop the indexes `idx_ai_analysis_expert_id`, `idx_ai_analysis_document_id`, and `idx_ai_analysis_type`.
- **Remove AI Migration**: Delete the migration file `db/migrations/sqlite/0011_create_ai_analysis_table.sql`, which creates the `ai_analysis_results` table.
- **Remove AI Types**: In `types.go`, remove the `AIAnalysisResult` struct and `AIAnalysisRequest` struct, as they are no longer used.
- **Update APIServer**: In `api.go`, remove the `aiService` field from the `APIServer` struct and its initialization in `NewAPIServer`.
- **Update Configuration**: In `types.go` and `server.go`, remove the `AIServiceURL` field from the `Configuration` struct and its handling in `loadConfig`.
- **Developer Note**: After removing AI integration, add a note in `api.go` at the top of the file:
  // NOTE: AI integration has been removed as per requirements. Previous AI routes and services are no longer available.

### 2. Commenting and Documentation
#### Function-Level Comments
- **Purpose**: Every function, handler, or method must start with a comment block describing its purpose, inputs, outputs, and side effects.
- **Format**: Use a consistent comment style, such as:
  // handleGetExperts handles GET /api/experts requests.
  // Inputs: w (http.ResponseWriter) - the response writer, r (*http.Request) - the HTTP request with query parameters.
  // Outputs: error - an error if the request fails, otherwise nil (response written to w).
  // Side effects: Writes JSON response to the client.
- **Application**: Apply this to all methods (e.g., `handleGetExperts`, `CreateExpert`).

#### Processing Steps
- **Purpose**: Break down the logic within each function into distinct steps, each preceded by an inline comment.
- **Structure**: Group related operations under a single comment describing the step, followed by the code block. Use blank lines to separate steps visually.
- **Example**:
  // Step 1: Parse request body into CreateExpertRequest
  createReq := new(CreateExpertRequest)
  if err := json.NewDecoder(r.Body).Decode(createReq); err != nil {
      return WriteJson(w, http.StatusBadRequest, ApiError{Error: "Invalid request payload"})
  }

  // Step 2: Validate the request data
  if err := ValidateCreateExpertRequest(createReq); err != nil {
      return WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
  }
- **Application**: Refactor functions like `handleCreateExpert` and `CreateExpertRequest` to clearly delineate steps such as validation, data retrieval, processing, and response handling.

#### Handler Comments
- **Purpose**: For API handlers, document the HTTP method, endpoint, and overall processing flow.
- **Format**: Include a comment block at the start of each handler with this structure:
  // handleGetExperts handles GET /api/experts requests.
  // It retrieves a list of experts based on query parameters.
  // Flow:
  // 1. Parse query parameters for filters and pagination.
  // 2. Query the database with applied filters.
  // 3. Set pagination headers and return JSON response.
- **Application**: Update all handlers in `api.go` (e.g., `handleGetExpert`, `handleUploadDocument`) with this level of detail.
- **Developer Note**: After updating `handleExtractSkills` removal, add a note in `api.go` where the handler was removed:
  // NOTE: handleExtractSkills removed due to AI integration removal. Skills extraction is no longer supported.

### 3. Code Organization
#### Group Related Logic
- **Purpose**: Ensure code within a function is grouped by purpose (e.g., validation, data fetching, response preparation).
- **Approach**: Use blank lines to separate logical sections, aligning with the processing step comments.
- **Application**: In `handleGetExperts`, group all filter parsing together, followed by pagination parsing, then database querying.

#### Extract Helper Functions
- **Purpose**: Reduce complexity by extracting repetitive or distinct tasks into helper functions.
- **Guidance**: If a function exceeds 20-30 lines or performs multiple tasks (e.g., validation and processing), extract logic into a helper function with its own documentation.
- **Example**:
  // parseFilters extracts filter parameters from the request query.
  // Inputs: queryParams (url.Values) - the request query parameters.
  // Outputs: map[string]interface{} - the parsed filters.
  func parseFilters(queryParams url.Values) map[string]interface{} {
      filters := make(map[string]interface{})
      if name := queryParams.Get("name"); name != "" {
          filters["name"] = name
      }
      return filters
  }
- **Application**: Consider extracting filter parsing from `handleGetExperts` or input validation from `CreateExpertRequest`.
- **Developer Note**: After extracting helpers, add a note in the relevant file (e.g., `api.go`):
  // NOTE: Extracted parseFilters to simplify handleGetExperts. Review for additional opportunities to reduce handler complexity.

### 4. Error Handling
#### Consistent Error Messages
- **Purpose**: Standardize error formatting for consistency across similar operations.
- **Approach**: Use a pattern like `"failed to [action]: %w"` (e.g., `"failed to retrieve expert: %w"`).
- **Application**: Update error messages in `handleUpdateExpert` and `CreateExpert` to follow this pattern.

#### Logging
- **Purpose**: Enhance debugging by logging key events and errors.
- **Approach**: Add `logger.Info` for successful operations and `logger.Error` for failures, using appropriate log levels.
- **Application**: In `handleLogin`, log successful logins and failed attempts.
- **Developer Note**: After adding logging, add a note in `auth.go`:
  // NOTE: Added logging for login attempts to improve traceability. Consider adding more detailed logs for user actions.

### 5. Database Operations
#### SQL Queries
- **Purpose**: Clarify complex queries with comments explaining their purpose.
- **Approach**: Add a comment above each query describing what it retrieves or modifies.
- **Example**:
  // Retrieve a list of available experts ordered by name
  query := `
      SELECT * FROM experts
      WHERE is_available = 1
      ORDER BY name ASC
      LIMIT ? OFFSET ?
  `
- **Application**: Apply to queries in `ListExperts` and `ListExpertRequests`.

#### Transaction Management
- **Purpose**: Document transaction boundaries for clarity.
- **Approach**: Comment the start, rollback, and commit points in multi-step database operations.
- **Application**: In `DeleteExpert`, add comments around the transaction logic.
- **Developer Note**: After updating `DeleteExpert`, add a note in `expert_operations.go`:
  // NOTE: Added transaction comments for clarity. Ensure all multi-step DB operations follow this pattern.

### 6. API and Middleware
#### Middleware Comments
- **Purpose**: Explain middleware functionality and its impact on requests/responses.
- **Approach**: Add a comment block describing what the middleware does.
- **Example**:
  // requireAuth middleware ensures the request has a valid JWT token.
  // It verifies the token and adds user claims to the context if valid.
- **Application**: Update `requireAuth` and `requireAdmin` in `auth.go`.

#### Handler Flow
- **Purpose**: Outline the sequence of operations in API handlers.
- **Approach**: Use step comments to detail request parsing, business logic, and response steps.
- **Application**: Refactor `handleUploadDocument` to show steps like form parsing, validation, file saving, and response.

### 7. Type Definitions
#### Struct Comments
- **Purpose**: Describe the purpose of struct fields where not obvious.
- **Approach**: Add field-level comments in struct definitions.
- **Example**:
  type Expert struct {
      ID        int64  `json:"id"`
      Name      string `json:"name"` // Full name of the expert
      Rating    string `json:"rating"` // Expert's performance rating (e.g., "4.5")
  }
- **Application**: Update `Expert` and `Engagement` in `types.go`.
- **Developer Note**: After updating `types.go`, add a note:
  // NOTE: Added field comments to Expert struct for clarity. Review other structs for similar documentation needs.

### 8. Refactoring Suggestions
#### Duplicated Code
- **Purpose**: Eliminate redundancy by extracting shared logic.
- **Approach**: Identify repeated patterns (e.g., JSON encoding in handlers) and move to utilities.
- **Application**: Extract `WriteJson` calls into a helper if reused frequently.

#### Magic Numbers and Strings
- **Purpose**: Replace hardcoded values with named constants for clarity.
- **Approach**: Define constants with comments explaining their purpose.
- **Example**:
  const DefaultLimit = 10 // Default number of items per page
- **Application**: Replace `10` in `handleGetExperts` pagination with a constant.
- **Developer Note**: After replacing magic numbers, add a note in `api.go`:
  // NOTE: Replaced magic number 10 with DefaultLimit constant. Search for other hardcoded values to replace.

### 9. Input Validation
#### Objective
- **Purpose**: Secure the application by validating inputs.
- **Approach**: Comment validation logic clearly, explaining what’s checked.
- **Application**: Enhance `ValidateCreateExpertRequest` with step comments for each validation.
- **Developer Note**: After updating validation, add a note in `types.go`:
  // NOTE: Enhanced validation comments in ValidateCreateExpertRequest. Consider adding validation for other structs like ExpertRequest.

## Final Notes
- **Consistency**: Use the same comment style and code organization across all files.
- **Clarity**: Write comments that a new developer can understand without prior context.
- **Examples**: Use brief examples (like above) to clarify complex logic where needed.
- **Review**: After refactoring, verify each function/handler/method has a clear purpose and step-by-step breakdown, and ensure all AI-related code is removed.

Implement these guidelines across the provided files, focusing on `api.go`, `auth.go`, `storage.go`, and other relevant files after removing AI integration. If you encounter ambiguities or need clarification, note them for review.
