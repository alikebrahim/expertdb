-- +goose Up
-- Create a table for expert areas (categories)
CREATE TABLE IF NOT EXISTS "expert_areas" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Populate the expert_areas table with predefined specializations
-- Organized by category for logical grouping and efficient mapping
INSERT INTO expert_areas (name) VALUES
    -- Business Areas (IDs 1-11)
    ("Business"),
    ("Business - Accounting & Audit"),
    ("Business - Banking & Finance"),
    ("Business - Compliance"),
    ("Business - Economics"),
    ("Business - Entrepreneurship & Innovation"),
    ("Business - Human Resources"),
    ("Business - Insurance"),
    ("Business - Islamic Banking & Finance"),
    ("Business - Management & Marketing"),
    ("Business - Project Management"),
    
    -- Education & Training (IDs 12-13)
    ("Education"),
    ("Training"),
    
    -- Engineering Areas (IDs 14-24)
    ("Engineering"),
    ("Engineering - Architectural"),
    ("Engineering - Chemical"),
    ("Engineering - Civil"),
    ("Engineering - Electrical and Electronic"),
    ("Engineering - Environmental"),
    ("Engineering - Industrial"),
    ("Engineering - Marine"),
    ("Engineering - Mechanical"),
    ("Engineering - Petroleum & Gas"),
    ("Engineering - Software"),
    
    -- Information Technology (ID 25)
    ("Information Technology"),
    
    -- Science Areas (IDs 26-34)
    ("Science"),
    ("Science - Biology"),
    ("Science - Chemistry"),
    ("Science - Computer Science"),
    ("Science - Data Science & Analytics"),
    ("Science - Environment"),
    ("Science - Geology"),
    ("Science - Mathematics"),
    ("Science - Physics"),
    
    -- Medical Science (IDs 35-38)
    ("Medical Science"),
    ("Medical Science - Healthcare Management"),
    ("Medical Science - Pharmaceutical"),
    ("Medical Science - Public Health"),
    
    -- Legal & Compliance (ID 39)
    ("Law"),
    
    -- Arts, Design & Media (IDs 40-42)
    ("Art and Design"),
    ("English"),
    ("Media & Communications"),
    
    -- Specialized Fields (IDs 43-54)
    ("Agriculture & Food Security"),
    ("Aviation"),
    ("Cybersecurity"),
    ("Finance - Investment & Capital Markets"),
    ("Health & Safety"),
    ("Hospitality and Tourism"),
    ("Psychology & Behavioral Sciences"),
    ("Quality Assurance"),
    ("Renewable Energy & Sustainability"),
    ("Social Sciences");

-- +goose Down
DROP TABLE IF EXISTS "expert_areas";
