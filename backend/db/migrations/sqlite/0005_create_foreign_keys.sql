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
    nationality TEXT DEFAULT 'Bahraini' CHECK (nationality IN ('Bahraini', 'Non-Bahraini', 'Unknown')),
    is_available BOOLEAN,
    rating TEXT,
    role TEXT,
    employment_type TEXT,
    general_area INTEGER NOT NULL,
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_path TEXT,
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    biography TEXT,
    original_request_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (original_request_id) REFERENCES expert_requests(id) ON DELETE SET NULL,
    FOREIGN KEY (general_area) REFERENCES expert_areas(id)
);

-- Copy data from experts to experts_temp with transformed general_area
-- This assumes default to the first expert area (id=1) when converting from text
INSERT INTO experts_temp(id, expert_id, name, designation, institution, is_bahraini, 
                        nationality, is_available, rating, role, employment_type, 
                        general_area, specialized_area, is_trained, cv_path, phone, 
                        email, is_published, biography, original_request_id, 
                        created_at, updated_at)
SELECT id, expert_id, name, designation, institution, is_bahraini, 
       nationality, is_available, rating, role, employment_type, 
       (SELECT id FROM expert_areas WHERE expert_areas.name LIKE '%' || experts.general_area || '%' LIMIT 1) AS general_area,
       specialized_area, is_trained, cv_path, phone, 
       email, is_published, biography, original_request_id, 
       created_at, updated_at
FROM experts;

-- Drop old table and rename new one
DROP TABLE experts;
ALTER TABLE experts_temp RENAME TO experts;

-- Recreate indexes
CREATE INDEX idx_experts_name ON experts(name);
CREATE INDEX idx_experts_general_area ON experts(general_area);
CREATE INDEX idx_experts_is_available ON experts(is_available);
CREATE INDEX idx_experts_nationality ON experts(nationality);

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
    general_area INTEGER NOT NULL,
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_path TEXT,
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    biography TEXT,
    status TEXT DEFAULT 'pending',
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by INTEGER,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (general_area) REFERENCES expert_areas(id)
);

-- Copy data from expert_requests to expert_requests_temp with transformed general_area
-- This assumes default to the first expert area (id=1) when converting from text
INSERT INTO expert_requests_temp(id, expert_id, name, designation, institution, 
                               is_bahraini, is_available, rating, role, employment_type, 
                               general_area, specialized_area, is_trained, cv_path, 
                               phone, email, is_published, biography, status, 
                               rejection_reason, created_at, reviewed_at, reviewed_by)
SELECT id, expert_id, name, designation, institution, 
       is_bahraini, is_available, rating, role, employment_type, 
       (SELECT id FROM expert_areas WHERE expert_areas.name LIKE '%' || expert_requests.general_area || '%' LIMIT 1) AS general_area,
       specialized_area, is_trained, cv_path, 
       phone, email, is_published, biography, status, 
       rejection_reason, created_at, reviewed_at, reviewed_by
FROM expert_requests;

-- Drop old table and rename new one
DROP TABLE expert_requests;
ALTER TABLE expert_requests_temp RENAME TO expert_requests;

-- Recreate indexes
CREATE INDEX idx_expert_requests_status ON expert_requests(status);
CREATE INDEX idx_expert_requests_created_at ON expert_requests(created_at);
CREATE INDEX idx_expert_requests_general_area ON expert_requests(general_area);

-- +goose Down
-- No specific down migration needed, as the tables will be dropped 
-- by their original migration files
