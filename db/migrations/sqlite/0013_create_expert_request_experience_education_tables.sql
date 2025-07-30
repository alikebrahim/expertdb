-- +goose Up
-- Expert request experience table for tracking professional experience entries during request workflow
CREATE TABLE IF NOT EXISTS "expert_request_experience_entries" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_request_id INTEGER NOT NULL,
    organization TEXT NOT NULL,
    position TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT,
    is_current BOOLEAN DEFAULT FALSE,
    country TEXT,
    description TEXT,
    FOREIGN KEY (expert_request_id) REFERENCES expert_requests(id) ON DELETE CASCADE
);

-- Expert request education table for tracking educational background entries during request workflow
CREATE TABLE IF NOT EXISTS "expert_request_education_entries" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_request_id INTEGER NOT NULL,
    institution TEXT NOT NULL,
    degree TEXT NOT NULL,
    field_of_study TEXT,
    graduation_year TEXT NOT NULL,
    country TEXT,
    description TEXT,
    FOREIGN KEY (expert_request_id) REFERENCES expert_requests(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_expert_request_experience_entries_request_id ON expert_request_experience_entries(expert_request_id);
CREATE INDEX idx_expert_request_education_entries_request_id ON expert_request_education_entries(expert_request_id);

-- +goose Down
DROP INDEX IF EXISTS idx_expert_request_education_entries_request_id;
DROP INDEX IF EXISTS idx_expert_request_experience_entries_request_id;
DROP TABLE IF EXISTS "expert_request_education_entries";
DROP TABLE IF EXISTS "expert_request_experience_entries";