-- +goose Up
CREATE TABLE IF NOT EXISTS "expert_requests" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    designation TEXT,
    affiliation TEXT,
    is_bahraini BOOLEAN,
    is_available BOOLEAN,
    role TEXT,               -- Evaluator, Validator or both
    employment_type TEXT,    -- Academic, Employer or both
    general_area INTEGER,    -- Reference to expert_areas table
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_document_id INTEGER REFERENCES expert_documents(id),
    approval_document_id INTEGER REFERENCES expert_documents(id),
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    suggested_specialized_areas TEXT, -- JSON array of user-suggested area names: ["Area Name 1", "Area Name 2"]
    status TEXT DEFAULT 'pending', -- pending, approved, rejected
    rejection_reason TEXT,   -- Reason for rejection if status is 'rejected'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by INTEGER,     -- References users(id)
    created_by INTEGER,      -- References users(id)
    FOREIGN KEY (general_area) REFERENCES expert_areas(id),
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes for tracking
CREATE INDEX idx_expert_requests_status ON expert_requests(status);
CREATE INDEX idx_expert_requests_created_at ON expert_requests(created_at);
CREATE INDEX idx_expert_requests_general_area ON expert_requests(general_area);
CREATE INDEX idx_expert_requests_created_by ON expert_requests(created_by);
CREATE INDEX idx_expert_requests_cv_document ON expert_requests(cv_document_id);
CREATE INDEX idx_expert_requests_approval_document ON expert_requests(approval_document_id);

-- +goose Down
DROP INDEX IF EXISTS idx_expert_requests_approval_document;
DROP INDEX IF EXISTS idx_expert_requests_cv_document;
DROP INDEX IF EXISTS idx_expert_requests_created_by;
DROP INDEX IF EXISTS idx_expert_requests_general_area;
DROP INDEX IF EXISTS idx_expert_requests_created_at;
DROP INDEX IF EXISTS idx_expert_requests_status;
DROP TABLE IF EXISTS "expert_requests";