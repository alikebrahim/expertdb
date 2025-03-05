-- Create the expert_documents table for CV and certificate storage
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

-- Create indexes for efficient queries
CREATE INDEX idx_documents_expert_id ON expert_documents(expert_id);
CREATE INDEX idx_documents_type ON expert_documents(document_type);