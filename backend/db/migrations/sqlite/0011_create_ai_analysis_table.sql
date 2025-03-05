-- Create the ai_analysis_results table to store AI-generated content
CREATE TABLE ai_analysis_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER,
    document_id INTEGER,
    analysis_type TEXT NOT NULL, -- 'profile', 'isced_suggestion', 'skills_extraction', etc.
    result_data TEXT NOT NULL, -- JSON or text result from AI
    confidence_score REAL, -- Optional confidence score (0-1)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE,
    FOREIGN KEY (document_id) REFERENCES expert_documents(id) ON DELETE CASCADE
);

-- Create indexes for efficient queries
CREATE INDEX idx_ai_analysis_expert_id ON ai_analysis_results(expert_id);
CREATE INDEX idx_ai_analysis_document_id ON ai_analysis_results(document_id);
CREATE INDEX idx_ai_analysis_type ON ai_analysis_results(analysis_type);