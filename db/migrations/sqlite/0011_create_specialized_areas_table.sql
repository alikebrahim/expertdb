-- +goose Up
CREATE TABLE specialized_areas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_specialized_areas_name ON specialized_areas(name);

-- +goose Down
DROP INDEX IF EXISTS idx_specialized_areas_name;
DROP TABLE IF EXISTS specialized_areas;