-- Create the expert_engagements table to track expert utilization
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

-- Create indexes for efficient queries
CREATE INDEX idx_engagements_expert_id ON expert_engagements(expert_id);
CREATE INDEX idx_engagements_status ON expert_engagements(status);
CREATE INDEX idx_engagements_date ON expert_engagements(start_date);