-- +goose Up
-- Expert edit requests table for tracking proposed changes to existing expert profiles
CREATE TABLE IF NOT EXISTS "expert_edit_requests" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,                    -- References experts(id) - the expert being edited
    
    -- Core profile fields (NULL = no change proposed)
    name TEXT,
    designation TEXT,
    institution TEXT,
    phone TEXT,
    email TEXT,
    is_bahraini BOOLEAN,
    is_available BOOLEAN,
    rating INTEGER CHECK (rating >= 0 AND rating <= 5),
    role TEXT,                                     -- Evaluator, Validator or both
    employment_type TEXT,                          -- Academic, Employer or both
    general_area INTEGER,                          -- Reference to expert_areas table
    specialized_area TEXT,                         -- Comma-separated specialized area IDs
    is_trained BOOLEAN,
    is_published BOOLEAN,
    biography TEXT,
    suggested_specialized_areas TEXT,              -- JSON array of user-suggested area names
    
    -- Document updates (NULL = no change proposed)
    new_cv_document_id INTEGER REFERENCES expert_documents(id),     -- Reference to updated CV document
    new_approval_document_id INTEGER REFERENCES expert_documents(id), -- Reference to updated approval document
    remove_cv BOOLEAN DEFAULT FALSE,              -- Flag to indicate CV should be removed
    remove_approval_document BOOLEAN DEFAULT FALSE, -- Flag to indicate approval document should be removed
    
    -- Change metadata
    change_summary TEXT NOT NULL,                 -- User-provided summary of changes
    change_reason TEXT NOT NULL,                  -- Reason for requesting the edit
    fields_changed TEXT NOT NULL,                 -- JSON array of field names being changed
    
    -- Status and workflow
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    rejection_reason TEXT,                        -- Reason for rejection if status is 'rejected'
    admin_notes TEXT,                            -- Internal notes for admin review
    
    -- Audit trail
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    applied_at TIMESTAMP,                        -- When the changes were applied to the expert
    created_by INTEGER NOT NULL,                -- References users(id) - who requested the edit
    reviewed_by INTEGER,                         -- References users(id) - who reviewed the request
    
    -- Constraints
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE,
    FOREIGN KEY (general_area) REFERENCES expert_areas(id),
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Experience entries for edit requests
CREATE TABLE IF NOT EXISTS "expert_edit_request_experience" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_edit_request_id INTEGER NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('add', 'update', 'delete')), -- What action to perform
    experience_id INTEGER,                       -- For update/delete actions, references original experience entry
    organization TEXT,
    position TEXT,
    start_date TEXT,
    end_date TEXT,
    is_current BOOLEAN,
    country TEXT,
    description TEXT,
    FOREIGN KEY (expert_edit_request_id) REFERENCES expert_edit_requests(id) ON DELETE CASCADE
);

-- Education entries for edit requests  
CREATE TABLE IF NOT EXISTS "expert_edit_request_education" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_edit_request_id INTEGER NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('add', 'update', 'delete')), -- What action to perform
    education_id INTEGER,                        -- For update/delete actions, references original education entry
    institution TEXT,
    degree TEXT,
    field_of_study TEXT,
    graduation_year TEXT,
    country TEXT,
    description TEXT,
    FOREIGN KEY (expert_edit_request_id) REFERENCES expert_edit_requests(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_expert_edit_requests_expert_id ON expert_edit_requests(expert_id);
CREATE INDEX idx_expert_edit_requests_status ON expert_edit_requests(status);
CREATE INDEX idx_expert_edit_requests_created_at ON expert_edit_requests(created_at);
CREATE INDEX idx_expert_edit_requests_created_by ON expert_edit_requests(created_by);
CREATE INDEX idx_expert_edit_requests_reviewed_by ON expert_edit_requests(reviewed_by);

CREATE INDEX idx_expert_edit_request_experience_request_id ON expert_edit_request_experience(expert_edit_request_id);
CREATE INDEX idx_expert_edit_request_education_request_id ON expert_edit_request_education(expert_edit_request_id);

-- +goose Down
DROP INDEX IF EXISTS idx_expert_edit_request_education_request_id;
DROP INDEX IF EXISTS idx_expert_edit_request_experience_request_id;
DROP INDEX IF EXISTS idx_expert_edit_requests_reviewed_by;
DROP INDEX IF EXISTS idx_expert_edit_requests_created_by;
DROP INDEX IF EXISTS idx_expert_edit_requests_created_at;
DROP INDEX IF EXISTS idx_expert_edit_requests_status;
DROP INDEX IF EXISTS idx_expert_edit_requests_expert_id;
DROP TABLE IF EXISTS "expert_edit_request_education";
DROP TABLE IF EXISTS "expert_edit_request_experience";
DROP TABLE IF EXISTS "expert_edit_requests";