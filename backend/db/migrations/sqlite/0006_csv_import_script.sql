-- +goose Up
-- This is a placeholder for the CSV import script
-- In real implementation, this would be a Go script that:
-- 1. Parses the CSV file
-- 2. Converts string values to appropriate data types
-- 3. Populates the expert_areas table with unique areas
-- 4. Inserts experts with proper relationships
-- 5. Creates appropriate entries in expert_specializations

-- Sample data for expert_areas (to be generated from CSV)
INSERT INTO expert_areas (name) VALUES 
('Business - Banking & Finance'),
('Business - Management & Marketing'),
('Business - Accounting & Audit'),
('Business - Islamic Banking & Finance'),
('Business - Economics'),
('Business - Project Management'),
('Business - Compliance'),
('Business - Insurance'),
('Information Technology'),
('Engineering - Electrical and Electronic'),
('Engineering - Mechanical'),
('Engineering - Chemical'),
('Engineering - Civil'),
('Engineering - Architectural'),
('Medical Science'),
('Science - Chemistry'),
('Science - Physics'),
('Science - Biology'),
('Science - Environment'),
('Science - Mathematics'),
('Art and Design'),
('Education'),
('English'),
('Law'),
('Aviation'),
('Health & Safety'),
('Quality Assurance'),
('Hospitality and Tourism'),
('Training');

-- +goose Down
-- Clear imported data (if needed)
DELETE FROM expert_specializations;
DELETE FROM experts WHERE expert_id LIKE 'E%';
DELETE FROM expert_areas;