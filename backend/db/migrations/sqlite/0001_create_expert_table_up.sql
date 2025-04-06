-- +goose Up
CREATE TABLE IF NOT EXISTS "experts" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id TEXT UNIQUE,  -- Original ID like "E001"
    name TEXT NOT NULL,
    designation TEXT,
    institution TEXT,
    is_bahraini BOOLEAN,    -- Convert "Yes/No" to boolean
    nationality TEXT DEFAULT 'Bahraini' CHECK (nationality IN ('Bahraini', 'Non-Bahraini', 'Unknown')),
    is_available BOOLEAN,   -- Convert "Yes/No" to boolean
    rating TEXT,
    role TEXT,              -- Evaluator, Validator or both
    employment_type TEXT,   -- Academic, Employer or both
    general_area INTEGER,   -- Reference to expert_areas table
    specialized_area TEXT,
    is_trained BOOLEAN,     -- Convert "Yes/No" to boolean
    cv_path TEXT,           -- Path to the CV file NOTE: This is better be replaced with expert_documents(id)
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,   -- Convert "Yes/No" to boolean
    biography TEXT,         -- Extended profile information
    original_request_id INTEGER, -- Foreign key referencing expert_requests
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create indexes for common search fields
CREATE INDEX idx_experts_name ON experts(name);
CREATE INDEX idx_experts_general_area ON experts(general_area);
CREATE INDEX idx_experts_is_available ON experts(is_available);
CREATE INDEX idx_experts_nationality ON experts(nationality);

-- +goose Down
DROP INDEX IF EXISTS idx_experts_nationality;
DROP INDEX IF EXISTS idx_experts_is_available;
DROP INDEX IF EXISTS idx_experts_general_area;
DROP INDEX IF EXISTS idx_experts_name;
DROP TABLE IF EXISTS "experts";
