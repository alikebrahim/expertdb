-- +goose Up
-- Migration to update phase application types from "validation"/"evaluation" to "QP"/"IL"

-- Update existing data to use the new types
UPDATE phase_applications 
SET type = 'QP' 
WHERE type = 'validation';

UPDATE phase_applications 
SET type = 'IL' 
WHERE type = 'evaluation';

-- +goose Down
-- Revert back to old types
UPDATE phase_applications 
SET type = 'validation' 
WHERE type = 'QP';

UPDATE phase_applications 
SET type = 'evaluation' 
WHERE type = 'IL';