Part 1: Error Messaging Report

Below is a detailed report identifying functions/methods in the ExpertDB backend where error messaging can be improved. The focus is on enhancing clarity and specificity for the small team (10-12 users) managing a modest database (up to 1200 entries), ensuring messages are actionable and user-friendly while maintaining simplicity, as security is handled organizationally and high load is not a concern.

1. File: internal/api/handlers/expert.go





Function/Method: HandleCreateExpert





Signature: func (h *ExpertHandler) HandleCreateExpert(w http.ResponseWriter, r *http.Request) error



Error Handling Stage:





JSON parsing error (line ~359).



Validation errors (line ~395).



Database insertion error (line ~404).



Current Behavior:





JSON parsing: Returns invalid request body: %w for any decoding error, lacking specifics (e.g., "invalid character").



Validation: Returns generic messages like invalid email format, without indicating the field or context.



Database: Returns failed to create expert: %w or an expert with this ID already exists for UNIQUE constraint violations, but other DB errors are vague.



Suggestion:





JSON Parsing: Include the specific JSON error (e.g., "unexpected EOF") to help developers debug malformed requests.



Validation: Collect all validation errors (e.g., invalid email, negative generalArea) and return them as a list for clarity.



Database: Differentiate between constraint violations (HTTP 409) and other errors (HTTP 500), providing field-specific messages (e.g., "expert_id already exists").



Proposed Code:

func (h *ExpertHandler) HandleCreateExpert(w http.ResponseWriter, r *http.Request) error {
    log := logger.Get()
    log.Debug("Processing POST /api/experts request")
    
    var expert domain.Expert
    if err := json.NewDecoder(r.Body).Decode(&expert); err != nil {
        log.Warn("Failed to parse expert creation request: %v", err)
        return writeJSON(w, http.StatusBadRequest, map[string]string{
            "error": fmt.Sprintf("Invalid JSON format: %v", err),
        })
    }
    
    // Validate fields
    errors := []string{}
    if expert.Name == "" {
        errors = append(errors, "name is required")
    }
    if expert.GeneralArea < 0 {
        errors = append(errors, "generalArea must be positive")
    }
    if expert.Email != "" && !isValidEmail(expert.Email) {
        errors = append(errors, fmt.Sprintf("invalid email format: %s", expert.Email))
    }
    if len(errors) > 0 {
        log.Warn("Expert creation validation failed: %v", errors)
        return writeJSON(w, http.StatusBadRequest, map[string][]string{
            "errors": errors,
        })
    }
    
    if expert.CreatedAt.IsZero() {
        expert.CreatedAt = time.Now()
    }
    
    log.Debug("Creating expert: %s, Institution: %s", expert.Name, expert.Institution)
    id, err := h.store.CreateExpert(&expert)
    if err != nil {
        log.Error("Failed to create expert in database: %v", err)
        if strings.Contains(err.Error(), "UNIQUE constraint failed: experts.expert_id") {
            return writeJSON(w, http.StatusConflict, map[string]string{
                "error": fmt.Sprintf("Expert ID %s already exists", expert.ExpertID),
            })
        }
        return writeJSON(w, http.StatusInternalServerError, map[string]string{
            "error": fmt.Sprintf("Database error creating expert: %v", err),
        })
    }
    
    log.Info("Expert created successfully with ID: %d", id)
    return writeJSON(w, http.StatusCreated, map[string]interface{}{
        "id":      id,
        "success": true,
        "message": "Expert created successfully",
    })
}

func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

2. File: internal/api/handlers/expert_request.go





Function/Method: HandleCreateExpertRequest





Signature: func (h *ExpertRequestHandler) HandleCreateExpertRequest(w http.ResponseWriter, r *http.Request) error



Error Handling Stage:





JSON parsing error (line ~1270).



Validation errors (line ~1306).



Database insertion error (line ~1338).



Current Behavior:





JSON parsing: Returns Invalid request payload for any JSON error, which is too generic.



Validation: Returns vague messages like name is required, without context for multiple errors.



Database: Returns failed to create expert request: %w, which doesn’t clarify the issue (e.g., constraint violation).



Suggestion:





JSON Parsing: Specify the JSON parsing issue (e.g., "missing closing brace") to aid debugging.



Validation: Aggregate all validation errors into a list, including field names and expected formats.



Database: Identify specific database errors (e.g., foreign key violations for general_area) and return appropriate HTTP codes (e.g., 400 for invalid references).



Proposed Code:

func (h *ExpertRequestHandler) HandleCreateExpertRequest(w http.ResponseWriter, r *http.Request) error {
    log := logger.Get()
    log.Debug("Processing POST /api/expert-requests request")
    
    var req domain.ExpertRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Warn("Failed to parse expert request: %v", err)
        return writeJSON(w, http.StatusBadRequest, map[string]string{
            "error": fmt.Sprintf("Invalid JSON format: %v", err),
        })
    }
    
    log.Debug("Validating expert request fields")
    errors := []string{}
    if req.Name == "" {
        errors = append(errors, "name is required")
    }
    if req.GeneralArea <= 0 {
        errors = append(errors, "generalArea must be a positive integer")
    }
    if req.Email != "" && !isValidEmail(req.Email) {
        errors = append(errors, fmt.Sprintf("invalid email format: %s", req.Email))
    }
    if len(errors) > 0 {
        log.Warn("Expert request validation failed: %v", errors)
        return writeJSON(w, http.StatusBadRequest, map[string][]string{
            "errors": errors,
        })
    }
    
    log.Debug("Setting default values for expert request")
    if req.CreatedAt.IsZero() {
        req.CreatedAt = time.Now()
    }
    if req.Status == "" {
        req.Status = "pending"
    }
    
    log.Debug("Creating expert request in database: %s, Institution: %s", req.Name, req.Institution)
    id, err := h.store.CreateExpertRequest(&req)
    if err != nil {
        log.Error("Failed to create expert request: %v", err)
        if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
            return writeJSON(w, http.StatusBadRequest, map[string]string{
                "error": "Invalid general_area: referenced expert area does not exist",
            })
        }
        return writeJSON(w, http.StatusInternalServerError, map[string]string{
            "error": fmt.Sprintf("Database error creating expert request: %v", err),
        })
    }
    
    log.Info("Expert request created successfully: ID: %d, Name: %s", id, req.Name)
    return writeJSON(w, http.StatusCreated, req)
}

func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

3. File: internal/storage/sqlite/expert.go





Function/Method: CreateExpert





Signature: func (s *SQLiteStore) CreateExpert(expert *domain.Expert) (int64, error)



Error Handling Stage:





Expert ID generation (line ~30).



Database insertion (line ~70).



Current Behavior:





ID generation: Returns failed to generate unique expert ID: %w if uniqueness checks fail, without indicating why (e.g., too many retries).



Database: Returns failed to create expert: %w, which wraps the raw SQLite error (e.g., "UNIQUE constraint failed"), but doesn’t clarify the conflicting field.



Suggestion:





ID Generation: Specify the failure reason (e.g., "maximum retry attempts exceeded") and suggest checking database state.



Database: Parse SQLite errors to return specific messages (e.g., "expert_id already exists" or "foreign key violation for general_area").



Proposed Code:

func (s *SQLiteStore) CreateExpert(expert *domain.Expert) (int64, error) {
    if expert.ExpertID == "" {
        var err error
        expert.ExpertID, err = s.GenerateUniqueExpertID()
        if err != nil {
            return 0, fmt.Errorf("failed to generate unique expert ID: %w", err)
        }
    } else {
        exists, err := s.ExpertIDExists(expert.ExpertID)
        if err != nil {
            return 0, fmt.Errorf("failed to check if expert ID exists: %w", err)
        }
        if exists {
            return 0, fmt.Errorf("expert ID %s already exists", expert.ExpertID)
        }
    }
    
    query := `
        INSERT INTO experts (
            expert_id, name, designation, institution, is_bahraini, is_available, rating,
            role, employment_type, general_area, specialized_area, is_trained,
            cv_path, phone, email, is_published, biography, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    if expert.CreatedAt.IsZero() {
        expert.CreatedAt = time.Now().UTC()
        expert.UpdatedAt = expert.CreatedAt
    }
    
    result, err := s.db.Exec(
        query,
        expert.ExpertID, expert.Name, expert.Designation, expert.Institution,
        expert.IsBahraini, expert.IsAvailable, expert.Rating,
        expert.Role, expert.EmploymentType, expert.GeneralArea, expert.SpecializedArea,
        expert.IsTrained, expert.CVPath, expert.Phone, expert.Email, expert.IsPublished,
        expert.Biography, expert.CreatedAt, expert.UpdatedAt,
    )
    
    if err != nil {
        if strings.Contains(err.Error(), "UNIQUE constraint failed: experts.expert_id") {
            return 0, fmt.Errorf("expert ID %s already exists", expert.ExpertID)
        }
        if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
            return 0, fmt.Errorf("invalid general_area: referenced expert area does not exist")
        }
        return 0, fmt.Errorf("failed to insert expert into database: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve new expert ID: %w", err)
    }
    
    return id, nil
}

4. File: internal/storage/sqlite/expert_request.go





Function/Method: CreateExpertRequest





Signature: func (s *SQLiteStore) CreateExpertRequest(req *domain.ExpertRequest) (int64, error)



Error Handling Stage:





Database insertion (line ~50).



Current Behavior:





Returns failed to create expert request: %w, wrapping the raw SQLite error without distinguishing between constraint violations, foreign key issues, or other failures.



Suggestion:





Parse database errors to provide specific messages (e.g., "invalid general_area reference") and ensure the caller can act on the information.



Proposed Code:

func (s *SQLiteStore) CreateExpertRequest(req *domain.ExpertRequest) (int64, error) {
    query := `
        INSERT INTO expert_requests (
            name, designation, institution, is_bahraini, is_available,
            rating, role, employment_type, general_area, specialized_area,
            is_trained, cv_path, phone, email, is_published, biography,
            status, created_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    if req.CreatedAt.IsZero() {
        req.CreatedAt = time.Now()
    }
    if req.Status == "" {
        req.Status = "pending"
    }
    
    designation := req.Designation
    if designation == "" {
        designation = ""
    }
    
    institution := req.Institution
    if institution == "" {
        institution = ""
    }
    
    var rating interface{} = nil
    if req.Rating != "" {
        rating = req.Rating
    }
    
    var specializedArea interface{} = nil
    if req.SpecializedArea != "" {
        specializedArea = req.SpecializedArea
    }
    
    var cvPath interface{} = nil
    if req.CVPath != "" {
        cvPath = req.CVPath
    }
    
    var biography interface{} = nil
    if req.Biography != "" {
        biography = req.Biography
    }
    
    result, err := s.db.Exec(
        query,
        req.Name, designation, institution,
        req.IsBahraini, req.IsAvailable, rating,
        req.Role, req.EmploymentType, req.GeneralArea,
        specializedArea, req.IsTrained, cvPath,
        req.Phone, req.Email, req.IsPublished, biography,
        req.Status, req.CreatedAt,
    )
    
    if err != nil {
        if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
            return 0, fmt.Errorf("invalid general_area: referenced expert area does not exist")
        }
        return 0, fmt.Errorf("failed to insert expert request into database: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve new expert request ID: %w", err)
    }
    
    req.ID = id
    return id, nil
}
