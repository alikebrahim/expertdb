-- +goose Up
-- Create system_statistics table for caching frequently accessed statistics
CREATE TABLE system_statistics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stat_key TEXT NOT NULL,
    stat_value TEXT NOT NULL, -- JSON formatted statistics data
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(stat_key)
);

-- Nationality column already exists in experts table from migration 0001
-- No need to add it again

-- Ensure index for nationality filtering exists
CREATE INDEX IF NOT EXISTS idx_experts_nationality ON experts(nationality);

-- +goose Down
-- Drop index first
DROP INDEX IF EXISTS idx_experts_nationality;

-- Drop table
DROP TABLE IF EXISTS system_statistics;