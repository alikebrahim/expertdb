-- +goose Up
-- Add foreign key constraints between tables

-- Add foreign key from experts to expert_requests
CREATE TABLE IF NOT EXISTS "experts_temp" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id TEXT UNIQUE,
    name TEXT NOT NULL,
    designation TEXT,
    institution TEXT,
    is_bahraini BOOLEAN,
    is_available BOOLEAN,
    rating TEXT,
    role TEXT,
    employment_type TEXT,
    general_area TEXT,
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_path TEXT,
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    original_request_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (original_request_id) REFERENCES expert_requests(id) ON DELETE SET NULL
);

-- Copy data from experts to experts_temp
INSERT INTO experts_temp 
SELECT * FROM experts;

-- Drop old table and rename new one
DROP TABLE experts;
ALTER TABLE experts_temp RENAME TO experts;

-- Recreate indexes
CREATE INDEX idx_experts_name ON experts(name);
CREATE INDEX idx_experts_general_area ON experts(general_area);
CREATE INDEX idx_experts_is_available ON experts(is_available);

-- Add foreign key from expert_requests to users
CREATE TABLE IF NOT EXISTS "expert_requests_temp" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id TEXT,
    name TEXT NOT NULL,
    designation TEXT,
    institution TEXT,
    is_bahraini BOOLEAN,
    is_available BOOLEAN,
    rating TEXT,
    role TEXT,
    employment_type TEXT,
    general_area TEXT,
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_path TEXT,
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by INTEGER,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Copy data from expert_requests to expert_requests_temp
INSERT INTO expert_requests_temp 
SELECT * FROM expert_requests;

-- Drop old table and rename new one
DROP TABLE expert_requests;
ALTER TABLE expert_requests_temp RENAME TO expert_requests;

-- Recreate indexes
CREATE INDEX idx_expert_requests_status ON expert_requests(status);
CREATE INDEX idx_expert_requests_created_at ON expert_requests(created_at);

-- +goose Down
-- No specific down migration needed, as the tables will be dropped 
-- by their original migration files