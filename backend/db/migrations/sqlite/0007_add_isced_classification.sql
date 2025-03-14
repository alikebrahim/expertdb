-- +goose Up
-- Create ISCED education levels reference table
CREATE TABLE IF NOT EXISTS "isced_levels" (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    UNIQUE(code)
);

-- Create ISCED fields of education reference table
CREATE TABLE IF NOT EXISTS "isced_fields" (
    id INTEGER PRIMARY KEY,
    broad_code TEXT NOT NULL,
    broad_name TEXT NOT NULL,
    narrow_code TEXT,
    narrow_name TEXT,
    detailed_code TEXT,
    detailed_name TEXT,
    description TEXT,
    UNIQUE(broad_code, narrow_code, detailed_code)
);

-- Populate ISCED education levels (0-8)
INSERT INTO isced_levels (id, code, name, description) VALUES
(1, 'ISCED 0', 'Early childhood education', 'Pre-primary education'),
(2, 'ISCED 1', 'Primary education', 'Primary education or first stage of basic education'),
(3, 'ISCED 2', 'Lower secondary education', 'Lower secondary or second stage of basic education'),
(4, 'ISCED 3', 'Upper secondary education', 'Upper secondary education'),
(5, 'ISCED 4', 'Post-secondary non-tertiary education', 'Post-secondary non-tertiary education'),
(6, 'ISCED 5', 'Short-cycle tertiary education', 'First stage of tertiary education (short or medium duration)'),
(7, 'ISCED 6', 'Bachelor''s or equivalent level', 'Bachelor''s degree or equivalent'),
(8, 'ISCED 7', 'Master''s or equivalent level', 'Master''s degree or equivalent'),
(9, 'ISCED 8', 'Doctoral or equivalent level', 'Doctoral degree or equivalent');

-- Populate a few main ISCED broad fields as examples
INSERT INTO isced_fields (id, broad_code, broad_name, description) VALUES
(1, '00', 'Generic programmes and qualifications', 'General education programs'),
(2, '01', 'Education', 'Teacher training and education science'),
(3, '02', 'Arts and humanities', 'Arts, humanities, languages, etc.'),
(4, '03', 'Social sciences, journalism and information', 'Social and behavioral sciences, journalism, etc.'),
(5, '04', 'Business, administration and law', 'Business, management, law, etc.'),
(6, '05', 'Natural sciences, mathematics and statistics', 'Life sciences, physical sciences, mathematics, etc.'),
(7, '06', 'Information and Communication Technologies', 'Computer use, software, hardware, etc.'),
(8, '07', 'Engineering, manufacturing and construction', 'Engineering, manufacturing, architecture, etc.'),
(9, '08', 'Agriculture, forestry, fisheries and veterinary', 'Agriculture, forestry, fishery, veterinary'),
(10, '09', 'Health and welfare', 'Medicine, nursing, social services, etc.'),
(11, '10', 'Services', 'Personal services, transport, security, etc.');

-- +goose Down
-- Drop ISCED reference tables
DROP TABLE IF EXISTS isced_fields;
DROP TABLE IF EXISTS isced_levels;