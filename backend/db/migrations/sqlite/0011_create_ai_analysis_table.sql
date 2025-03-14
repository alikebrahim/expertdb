-- +goose Up
-- Create the ai_analysis table to store AI-generated content
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

-- Create indexes for efficient queries
CREATE INDEX idx_ai_analysis_expert_id ON ai_analysis(expert_id);
CREATE INDEX idx_ai_analysis_document_id ON ai_analysis(document_id);
CREATE INDEX idx_ai_analysis_type ON ai_analysis(analysis_type);

-- +goose Down
DROP INDEX IF EXISTS idx_ai_analysis_type;
DROP INDEX IF EXISTS idx_ai_analysis_document_id;
DROP INDEX IF EXISTS idx_ai_analysis_expert_id;
DROP TABLE IF EXISTS ai_analysis;