-- +goose Up
-- Create a table for expert areas (categories)
CREATE TABLE IF NOT EXISTS "expert_areas" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    parent_id INTEGER,       -- For hierarchical categorization (null for top-level)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Populate the expert_areas table with predefined specializations
INSERT INTO expert_areas (name) VALUES
    ("Art and Design"),
    ("Aviation"),
    ("Business"),
    ("Business - Accounting & Audit"),
    ("Business - Banking & Finance"),
    ("Business - Compliance"),
    ("Business - Economics"),
    ("Business - Insurance"),
    ("Business - Islamic Banking & Finance"),
    ("Business - Management & Marketing"),
    ("Business - Project Management"),
    ("Education"),
    ("Engineering"),
    ("Engineering - Architectural"),
    ("Engineering - Chemical"),
    ("Engineering - Civil"),
    ("Engineering - Electrical and Electronic"),
    ("Engineering - Mechanical"),
    ("English"),
    ("Health & Safety"),
    ("Hospitality and Tourism"),
    ("Information Technology"),
    ("Law"),
    ("Medical Science"),
    ("Quality Assurance"),
    ("Science"),
    ("Science - Biology"),
    ("Science - Chemistry"),
    ("Science - Environment"),
    ("Science - Mathematics"),
    ("Science - Physics"),
    ("Social Sciences"),
    ("Training");

-- +goose Down
DROP TABLE IF EXISTS "expert_areas";
