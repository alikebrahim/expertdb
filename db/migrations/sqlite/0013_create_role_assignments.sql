-- +goose Up
-- Migration file for contextual role assignments
-- Creates tables for application-specific planner and manager assignments

-- Create application_planners table for planner assignments to specific applications
CREATE TABLE IF NOT EXISTS application_planners (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id INTEGER NOT NULL,        -- Reference to phase_applications table
    user_id INTEGER NOT NULL,               -- Reference to users table
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES phase_applications(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(application_id, user_id)         -- Ensure unique assignment per application-user pair
);

-- Create application_managers table for manager assignments to specific applications
CREATE TABLE IF NOT EXISTS application_managers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id INTEGER NOT NULL,        -- Reference to phase_applications table
    user_id INTEGER NOT NULL,               -- Reference to users table
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES phase_applications(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(application_id, user_id)         -- Ensure unique assignment per application-user pair
);

-- Create indexes for better performance
CREATE INDEX idx_application_planners_application_id ON application_planners(application_id);
CREATE INDEX idx_application_planners_user_id ON application_planners(user_id);
CREATE INDEX idx_application_managers_application_id ON application_managers(application_id);
CREATE INDEX idx_application_managers_user_id ON application_managers(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_application_managers_user_id;
DROP INDEX IF EXISTS idx_application_managers_application_id;
DROP INDEX IF EXISTS idx_application_planners_user_id;
DROP INDEX IF EXISTS idx_application_planners_application_id;
DROP TABLE IF EXISTS application_managers;
DROP TABLE IF EXISTS application_planners;