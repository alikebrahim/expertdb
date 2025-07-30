-- +goose Up
CREATE TABLE IF NOT EXISTS "experts" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    designation TEXT,
    affiliation TEXT,
    is_bahraini BOOLEAN,    -- Convert "Yes/No" to boolean
    is_available BOOLEAN,   -- Convert "Yes/No" to boolean
    rating INTEGER DEFAULT 0 CHECK (rating >= 0 AND rating <= 5), -- Rating scale: 0=No Rating, 1-5=Performance rating
    role TEXT,              -- Evaluator, Validator or both
    employment_type TEXT,   -- Academic, Employer or both
    general_area INTEGER,   -- Reference to expert_areas table
    specialized_area TEXT,
    is_trained BOOLEAN,     -- Convert "Yes/No" to boolean
    cv_document_id INTEGER REFERENCES expert_documents(id),
    approval_document_id INTEGER REFERENCES expert_documents(id),
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,   -- Convert "Yes/No" to boolean
    biography TEXT,         -- Extended profile information
    original_request_id INTEGER, -- Foreign key referencing expert_requests
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    last_edited_by INTEGER, -- Foreign key referencing users(id) - who last edited this expert
    last_edited_at TIMESTAMP, -- When this expert was last edited
    FOREIGN KEY (general_area) REFERENCES expert_areas(id),
    FOREIGN KEY (original_request_id) REFERENCES expert_requests(id) ON DELETE SET NULL,
    FOREIGN KEY (last_edited_by) REFERENCES users(id) ON DELETE SET NULL
);


-- Create indexes for common search fields
CREATE INDEX idx_experts_name ON experts(name);
CREATE INDEX idx_experts_general_area ON experts(general_area);
CREATE INDEX idx_experts_is_available ON experts(is_available);
CREATE INDEX idx_experts_is_bahraini ON experts(is_bahraini);
CREATE INDEX idx_experts_specialized_area ON experts(specialized_area);
CREATE INDEX idx_experts_employment_type ON experts(employment_type);
CREATE INDEX idx_experts_role ON experts(role);
CREATE INDEX idx_experts_cv_document ON experts(cv_document_id);
CREATE INDEX idx_experts_approval_document ON experts(approval_document_id);
CREATE INDEX idx_experts_last_edited_by ON experts(last_edited_by);

-- +goose Down
DROP INDEX IF EXISTS idx_experts_last_edited_by;
DROP INDEX IF EXISTS idx_experts_approval_document;
DROP INDEX IF EXISTS idx_experts_cv_document;
DROP INDEX IF EXISTS idx_experts_role;
DROP INDEX IF EXISTS idx_experts_employment_type;
DROP INDEX IF EXISTS idx_experts_specialized_area;
DROP INDEX IF EXISTS idx_experts_is_bahraini;
DROP INDEX IF EXISTS idx_experts_is_available;
DROP INDEX IF EXISTS idx_experts_general_area;
DROP INDEX IF EXISTS idx_experts_name;
DROP TABLE IF EXISTS "experts";
