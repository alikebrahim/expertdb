-- +goose Up
-- Expert experience table for tracking professional experience entries
CREATE TABLE IF NOT EXISTS "expert_experience" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    organization TEXT NOT NULL,
    position TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT,
    is_current BOOLEAN DEFAULT FALSE,
    country TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE
);

-- Expert education table for tracking educational background entries
CREATE TABLE IF NOT EXISTS "expert_education" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    institution TEXT NOT NULL,
    degree TEXT NOT NULL,
    field_of_study TEXT,
    graduation_year TEXT NOT NULL,
    country TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_expert_experience_expert_id ON expert_experience(expert_id);
CREATE INDEX idx_expert_education_expert_id ON expert_education(expert_id);

-- +goose Down
DROP INDEX IF EXISTS idx_expert_education_expert_id;
DROP INDEX IF EXISTS idx_expert_experience_expert_id;
DROP TABLE IF EXISTS "expert_education";
DROP TABLE IF EXISTS "expert_experience";