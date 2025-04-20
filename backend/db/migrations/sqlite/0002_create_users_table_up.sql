-- +goose Up
-- Create users table with support for role hierarchy
CREATE TABLE IF NOT EXISTS "users" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('super_user', 'admin', 'scheduler', 'user')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Create indexes for email and role lookups
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- +goose Down
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS "users";