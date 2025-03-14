CREATE TABLE goose_db_version (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                version_id INTEGER NOT NULL,
                is_applied INTEGER NOT NULL,
                tstamp TIMESTAMP DEFAULT (datetime('now'))
            );
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE IF NOT EXISTS "users" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,      -- admin, reviewer, standard
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);
CREATE TABLE IF NOT EXISTS "expert_areas" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    parent_id INTEGER,       -- For hierarchical categorization (null for top-level)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS "expert_specializations" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    area_id INTEGER NOT NULL,
    proficiency_level TEXT,  -- beginner, intermediate, expert
    years_experience INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE,
    FOREIGN KEY (area_id) REFERENCES expert_areas(id) ON DELETE CASCADE,
    UNIQUE(expert_id, area_id)
);
CREATE INDEX idx_expert_specializations_expert_id ON expert_specializations(expert_id);
CREATE INDEX idx_expert_specializations_area_id ON expert_specializations(area_id);
CREATE TABLE IF NOT EXISTS "isced_levels" (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    UNIQUE(code)
);
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
CREATE TABLE expert_documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    document_type TEXT NOT NULL, -- 'cv', 'certificate', 'publication', etc.
    filename TEXT NOT NULL,
    file_path TEXT NOT NULL,
    content_type TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE
);
CREATE INDEX idx_documents_expert_id ON expert_documents(expert_id);
CREATE INDEX idx_documents_type ON expert_documents(document_type);
CREATE TABLE expert_engagements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,
    engagement_type TEXT NOT NULL, -- 'evaluation', 'consultation', 'project', etc.
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    project_name TEXT,
    status TEXT NOT NULL, -- 'pending', 'active', 'completed', 'cancelled'
    feedback_score INTEGER, -- 1-5 rating
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE
);
CREATE INDEX idx_engagements_expert_id ON expert_engagements(expert_id);
CREATE INDEX idx_engagements_status ON expert_engagements(status);
CREATE INDEX idx_engagements_date ON expert_engagements(start_date);
CREATE TABLE ai_analysis (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER,
    document_id INTEGER,
    analysis_type TEXT NOT NULL, -- 'profile', 'isced_suggestion', 'skills_extraction', etc.
    analysis_result TEXT NOT NULL, -- JSON or text result from AI
    confidence_score REAL, -- Optional confidence score (0-1)
    model_used TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE,
    FOREIGN KEY (document_id) REFERENCES expert_documents(id) ON DELETE CASCADE
);
CREATE INDEX idx_ai_analysis_expert_id ON ai_analysis(expert_id);
CREATE INDEX idx_ai_analysis_document_id ON ai_analysis(document_id);
CREATE INDEX idx_ai_analysis_type ON ai_analysis(analysis_type);
CREATE TABLE system_statistics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stat_key TEXT NOT NULL,
    stat_value TEXT NOT NULL, -- JSON formatted statistics data
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(stat_key)
);
CREATE TABLE migration_versions (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				filename TEXT UNIQUE NOT NULL,
				applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
CREATE TABLE IF NOT EXISTS "experts" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id TEXT UNIQUE,
    name TEXT NOT NULL,
    designation TEXT,
    institution TEXT,
    is_bahraini BOOLEAN,
    nationality TEXT DEFAULT 'Bahraini' CHECK (nationality IN ('Bahraini', 'Non-Bahraini', 'Unknown')),
    is_available BOOLEAN,
    rating TEXT,
    role TEXT,
    employment_type TEXT,
    general_area TEXT,
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_path TEXT,
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    biography TEXT,
    isced_level_id INTEGER,
    isced_field_id INTEGER,
    original_request_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (original_request_id) REFERENCES expert_requests(id) ON DELETE SET NULL
);
CREATE INDEX idx_experts_name ON experts(name);
CREATE INDEX idx_experts_general_area ON experts(general_area);
CREATE INDEX idx_experts_is_available ON experts(is_available);
CREATE TABLE IF NOT EXISTS "expert_requests" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id TEXT,
    name TEXT NOT NULL,
    designation TEXT,
    institution TEXT,
    is_bahraini BOOLEAN,
    is_available BOOLEAN,
    rating TEXT,
    role TEXT,
    employment_type TEXT,
    general_area TEXT,
    specialized_area TEXT,
    is_trained BOOLEAN,
    cv_path TEXT,
    phone TEXT,
    email TEXT,
    is_published BOOLEAN,
    biography TEXT,
    status TEXT DEFAULT 'pending',
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by INTEGER,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_expert_requests_status ON expert_requests(status);
CREATE INDEX idx_expert_requests_created_at ON expert_requests(created_at);
