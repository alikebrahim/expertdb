-- +goose Up
-- This migration attempts to map existing expert general areas to ISCED fields
-- It uses fuzzy matching based on general area descriptions

-- Map Business - Banking & Finance
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '04')
WHERE general_area LIKE '%Banking%' OR general_area LIKE '%Finance%';

-- Map Business - Management & Marketing
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '04')
WHERE general_area LIKE '%Management%' OR general_area LIKE '%Marketing%';

-- Map Business - Accounting & Audit
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '04')
WHERE general_area LIKE '%Accounting%' OR general_area LIKE '%Audit%';

-- Map Information Technology
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '06')
WHERE general_area LIKE '%Information Technology%';

-- Map Engineering fields
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '07')
WHERE general_area LIKE '%Engineering%';

-- Map Medical Science
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '09')
WHERE general_area LIKE '%Medical%' OR general_area LIKE '%Health%';

-- Map Science fields
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '05')
WHERE general_area LIKE '%Science%';

-- Map Education
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '01')
WHERE general_area LIKE '%Education%';

-- Map Law
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '04')
WHERE general_area LIKE '%Law%';

-- Map Arts and Design
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '02')
WHERE general_area LIKE '%Art%' OR general_area LIKE '%Design%';

-- Map Aviation and Transport
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '10')
WHERE general_area LIKE '%Aviation%' OR general_area LIKE '%Transport%';

-- Map Hospitality and Tourism
UPDATE experts 
SET isced_field_id = (SELECT id FROM isced_fields WHERE broad_code = '10')
WHERE general_area LIKE '%Hospitality%' OR general_area LIKE '%Tourism%';

-- Set default ISCED education level based on designation
-- (This is a simple heuristic - in reality you would need more information)
UPDATE experts 
SET isced_level_id = (SELECT id FROM isced_levels WHERE code = 'ISCED 6')  -- Bachelor's
WHERE designation = 'Mr.' OR designation = 'Ms.' OR designation = 'Mrs.';

UPDATE experts 
SET isced_level_id = (SELECT id FROM isced_levels WHERE code = 'ISCED 7')  -- Master's
WHERE designation = 'Dr.';

UPDATE experts 
SET isced_level_id = (SELECT id FROM isced_levels WHERE code = 'ISCED 8')  -- Doctoral
WHERE designation = 'Prof.' OR designation = 'Dr.';

-- +goose Down
-- Reset mapped ISCED fields and levels
UPDATE experts SET isced_field_id = NULL, isced_level_id = NULL;