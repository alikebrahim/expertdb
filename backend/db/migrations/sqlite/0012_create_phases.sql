-- Migration file for Phase Planning and Engagement System
-- Creates tables for phases and phase applications

-- Create phases table
CREATE TABLE IF NOT EXISTS phases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    phase_id TEXT NOT NULL UNIQUE,  -- Business identifier (e.g., "PH-2025-001")
    title TEXT NOT NULL,            -- Title/name of the phase
    assigned_scheduler_id INTEGER,  -- ID of the scheduler user assigned to this phase
    status TEXT NOT NULL,           -- Status: "draft", "in_progress", "completed", "cancelled"
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (assigned_scheduler_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create phase_applications table
CREATE TABLE IF NOT EXISTS phase_applications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    phase_id INTEGER NOT NULL,          -- Reference to phases table
    type TEXT NOT NULL,                 -- Type: "validation" or "evaluation"
    institution_name TEXT NOT NULL,     -- Name of the institution
    qualification_name TEXT NOT NULL,   -- Name of the qualification being reviewed
    expert_1 INTEGER,                   -- First expert ID (reference to experts table)
    expert_2 INTEGER,                   -- Second expert ID (reference to experts table)
    status TEXT NOT NULL,               -- Status: "pending", "assigned", "approved", "rejected"
    rejection_notes TEXT,               -- Notes for rejection (if status is "rejected")
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (phase_id) REFERENCES phases(id) ON DELETE CASCADE,
    FOREIGN KEY (expert_1) REFERENCES experts(id) ON DELETE SET NULL,
    FOREIGN KEY (expert_2) REFERENCES experts(id) ON DELETE SET NULL
);

-- Create indexes for better performance
CREATE INDEX idx_phases_phase_id ON phases(phase_id);
CREATE INDEX idx_phases_status ON phases(status);
CREATE INDEX idx_phases_assigned_scheduler ON phases(assigned_scheduler_id);
CREATE INDEX idx_phase_applications_phase_id ON phase_applications(phase_id);
CREATE INDEX idx_phase_applications_type ON phase_applications(type);
CREATE INDEX idx_phase_applications_status ON phase_applications(status);
CREATE INDEX idx_phase_applications_expert_1 ON phase_applications(expert_1);
CREATE INDEX idx_phase_applications_expert_2 ON phase_applications(expert_2);