-- Create system_statistics table for caching frequently accessed statistics
CREATE TABLE system_statistics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stat_key TEXT NOT NULL,
    stat_value TEXT NOT NULL, -- JSON formatted statistics data
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(stat_key)
);

-- Add nationality tracking to experts table if not already present
ALTER TABLE experts ADD COLUMN nationality TEXT DEFAULT 'Bahraini' CHECK (nationality IN ('Bahraini', 'Non-Bahraini', 'Unknown'));

-- Create index for nationality filtering
CREATE INDEX IF NOT EXISTS idx_experts_nationality ON experts(nationality);