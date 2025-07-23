-- +goose Up

-- Create expert experience entries table
CREATE TABLE IF NOT EXISTS "expert_experience_entries" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    organization TEXT NOT NULL,
    position TEXT NOT NULL,
    start_date TEXT, -- Can be year only or year-month format
    end_date TEXT,   -- Can be year only, year-month format, or "Present"
    is_current BOOLEAN DEFAULT 0,
    country TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE
);

-- Create expert education entries table
CREATE TABLE IF NOT EXISTS "expert_education_entries" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    institution TEXT NOT NULL,
    degree TEXT NOT NULL,
    field_of_study TEXT,
    graduation_year TEXT, -- Can be year only or year-month format
    country TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE
);

-- Create expert request experience entries table
CREATE TABLE IF NOT EXISTS "expert_request_experience_entries" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_request_id INTEGER NOT NULL,
    organization TEXT NOT NULL,
    position TEXT NOT NULL,
    start_date TEXT, -- Can be year only or year-month format
    end_date TEXT,   -- Can be year only, year-month format, or "Present"
    is_current BOOLEAN DEFAULT 0,
    country TEXT,
    description TEXT,
    FOREIGN KEY (expert_request_id) REFERENCES expert_requests(id) ON DELETE CASCADE
);

-- Create expert request education entries table
CREATE TABLE IF NOT EXISTS "expert_request_education_entries" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_request_id INTEGER NOT NULL,
    institution TEXT NOT NULL,
    degree TEXT NOT NULL,
    field_of_study TEXT,
    graduation_year TEXT, -- Can be year only or year-month format
    country TEXT,
    description TEXT,
    FOREIGN KEY (expert_request_id) REFERENCES expert_requests(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_expert_experience_entries_expert_id ON expert_experience_entries(expert_id);
CREATE INDEX idx_expert_education_entries_expert_id ON expert_education_entries(expert_id);
CREATE INDEX idx_expert_request_experience_entries_request_id ON expert_request_experience_entries(expert_request_id);
CREATE INDEX idx_expert_request_education_entries_request_id ON expert_request_education_entries(expert_request_id);

-- Remove biography column from experts table
ALTER TABLE experts DROP COLUMN biography;

-- Remove biography column from expert_requests table  
ALTER TABLE expert_requests DROP COLUMN biography;

-- +goose Down

-- Recreate biography columns (though data will be lost)
ALTER TABLE experts ADD COLUMN biography TEXT;
ALTER TABLE expert_requests ADD COLUMN biography TEXT;

-- Drop indexes
DROP INDEX IF EXISTS idx_expert_request_education_entries_request_id;
DROP INDEX IF EXISTS idx_expert_request_experience_entries_request_id;
DROP INDEX IF EXISTS idx_expert_education_entries_expert_id;
DROP INDEX IF EXISTS idx_expert_experience_entries_expert_id;

-- Drop tables in reverse order
DROP TABLE IF EXISTS expert_request_education_entries;
DROP TABLE IF EXISTS expert_request_experience_entries;
DROP TABLE IF EXISTS expert_education_entries;
DROP TABLE IF EXISTS expert_experience_entries;