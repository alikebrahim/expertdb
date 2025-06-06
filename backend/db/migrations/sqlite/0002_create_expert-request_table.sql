-- +goose Up
CREATE TABLE IF NOT EXISTS "expert_requests" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id TEXT,          -- Original ID if provided
    name TEXT NOT NULL,
    designation TEXT,
    institution TEXT,
    is_bahraini BOOLEAN,
    is_available BOOLEAN,
    rating TEXT,
    role TEXT,               -- Evaluator, Validator or both
    employment_type TEXT,    -- Academic, Employer or both
    general_area INTEGER,    -- Reference to expert_areas table
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_path TEXT,            -- Path to the CV file
    approval_document_path TEXT, -- Path to the approval document
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    biography TEXT,          -- Extended profile information
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

-- +goose Down
DROP INDEX IF EXISTS idx_expert_requests_created_by;
DROP INDEX IF EXISTS idx_expert_requests_general_area;
DROP INDEX IF EXISTS idx_expert_requests_created_at;
DROP INDEX IF EXISTS idx_expert_requests_status;
DROP TABLE IF EXISTS "expert_requests";