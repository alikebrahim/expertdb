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
    general_area TEXT,
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_path TEXT,            -- Path to the CV file
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    status TEXT DEFAULT 'pending', -- pending, approved, rejected
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by INTEGER      -- References users(id)
);

-- Create indexes for tracking
CREATE INDEX idx_expert_requests_status ON expert_requests(status);
CREATE INDEX idx_expert_requests_created_at ON expert_requests(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_expert_requests_created_at;
DROP INDEX IF EXISTS idx_expert_requests_status;
DROP TABLE IF EXISTS "expert_requests";
