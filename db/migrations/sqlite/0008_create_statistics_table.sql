-- +goose Up
-- Create system_statistics table for caching frequently accessed statistics
CREATE TABLE system_statistics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stat_key TEXT NOT NULL,
    stat_value TEXT NOT NULL, -- JSON formatted statistics data
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(stat_key)
);

-- Index for Bahraini filtering already exists from migration 0004
-- No additional indices needed

-- +goose Down
-- Drop table
DROP TABLE IF EXISTS system_statistics;