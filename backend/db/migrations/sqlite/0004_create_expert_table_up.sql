-- +goose Up
CREATE TABLE IF NOT EXISTS "experts" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id TEXT UNIQUE,  -- Original ID like "EXP-0001"
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
    approval_document_path TEXT, -- Path to the approval document
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,   -- Convert "Yes/No" to boolean
    biography TEXT,         -- Extended profile information
    original_request_id INTEGER, -- Foreign key referencing expert_requests
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (general_area) REFERENCES expert_areas(id),
    FOREIGN KEY (original_request_id) REFERENCES expert_requests(id) ON DELETE SET NULL
);

-- Create a sequence table for expert ID generation
CREATE TABLE IF NOT EXISTS expert_id_sequence (
    id INTEGER PRIMARY KEY CHECK (id = 1), -- Only one row allowed
    next_val INTEGER NOT NULL DEFAULT 1
);

-- Initialize the sequence with a single row
INSERT OR IGNORE INTO expert_id_sequence (id, next_val) VALUES (1, 1);

-- Create indexes for common search fields
CREATE INDEX idx_experts_name ON experts(name);
CREATE INDEX idx_experts_general_area ON experts(general_area);
CREATE INDEX idx_experts_is_available ON experts(is_available);
CREATE INDEX idx_experts_nationality ON experts(nationality);
CREATE INDEX idx_experts_expert_id ON experts(expert_id);
CREATE INDEX idx_experts_is_bahraini ON experts(is_bahraini);
CREATE INDEX idx_experts_specialized_area ON experts(specialized_area);
CREATE INDEX idx_experts_employment_type ON experts(employment_type);
CREATE INDEX idx_experts_role ON experts(role);

-- +goose Down
DROP INDEX IF EXISTS idx_experts_role;
DROP INDEX IF EXISTS idx_experts_employment_type;
DROP INDEX IF EXISTS idx_experts_specialized_area;
DROP INDEX IF EXISTS idx_experts_is_bahraini;
DROP INDEX IF EXISTS idx_experts_expert_id;
DROP INDEX IF EXISTS idx_experts_nationality;
DROP INDEX IF EXISTS idx_experts_is_available;
DROP INDEX IF EXISTS idx_experts_general_area;
DROP INDEX IF EXISTS idx_experts_name;
DROP TABLE IF EXISTS expert_id_sequence;
DROP TABLE IF EXISTS "experts";