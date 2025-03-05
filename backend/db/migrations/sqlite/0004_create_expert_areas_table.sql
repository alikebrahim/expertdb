-- +goose Up
-- Create a table for expert areas (categories)
CREATE TABLE IF NOT EXISTS "expert_areas" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    parent_id INTEGER,       -- For hierarchical categorization (null for top-level)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create a junction table for expert-to-area many-to-many relationship
CREATE TABLE IF NOT EXISTS "expert_specializations" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    area_id INTEGER NOT NULL,
    proficiency_level TEXT,  -- beginner, intermediate, expert
    years_experience INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE,
    FOREIGN KEY (area_id) REFERENCES expert_areas(id) ON DELETE CASCADE,
    UNIQUE(expert_id, area_id)
);

-- Create indexes for the junction table
CREATE INDEX idx_expert_specializations_expert_id ON expert_specializations(expert_id);
CREATE INDEX idx_expert_specializations_area_id ON expert_specializations(area_id);

-- +goose Down
DROP INDEX IF EXISTS idx_expert_specializations_area_id;
DROP INDEX IF EXISTS idx_expert_specializations_expert_id;
DROP TABLE IF EXISTS "expert_specializations";
DROP TABLE IF EXISTS "expert_areas";