-- +goose Up
-- Expert edit history table for tracking all changes made to expert profiles
CREATE TABLE IF NOT EXISTS "expert_edit_history" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expert_id INTEGER NOT NULL,                    -- References experts(id) - the expert that was edited
    edited_by INTEGER NOT NULL,                    -- References users(id) - who made the edit
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the edit was made
    fields_changed TEXT NOT NULL,                  -- JSON array of field names that were changed
    old_values TEXT,                               -- JSON object of previous field values
    new_values TEXT,                               -- JSON object of new field values
    change_reason TEXT,                            -- Optional reason for the change
    
    -- Constraints
    FOREIGN KEY (expert_id) REFERENCES experts(id) ON DELETE CASCADE,
    FOREIGN KEY (edited_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes for performance
CREATE INDEX idx_expert_edit_history_expert_id ON expert_edit_history(expert_id);
CREATE INDEX idx_expert_edit_history_edited_by ON expert_edit_history(edited_by);
CREATE INDEX idx_expert_edit_history_edited_at ON expert_edit_history(edited_at);

-- +goose Down
DROP INDEX IF EXISTS idx_expert_edit_history_edited_at;
DROP INDEX IF EXISTS idx_expert_edit_history_edited_by;
DROP INDEX IF EXISTS idx_expert_edit_history_expert_id;
DROP TABLE IF EXISTS "expert_edit_history";