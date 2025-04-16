-- Populate expert_engagements table with mock data

-- Helper function for random range (using deterministic values for consistency)
WITH RECURSIVE range(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM range WHERE n < 100
),
-- Get all expert IDs from the database
expert_ids AS (
  SELECT id FROM experts ORDER BY id LIMIT 50
)

-- Insert engagements for each expert
INSERT INTO expert_engagements 
  (expert_id, engagement_type, start_date, end_date, project_name, status, feedback_score, notes, created_at)
SELECT 
  e.id AS expert_id,
  CASE (n % 5)
    WHEN 0 THEN 'evaluation'
    WHEN 1 THEN 'consultation'
    WHEN 2 THEN 'project'
    WHEN 3 THEN 'training'
    WHEN 4 THEN 'review'
  END AS engagement_type,
  date('now', '-' || (n % 12) || ' months', '-' || (n % 30) || ' days') AS start_date,
  CASE WHEN n % 10 > 3 -- 70% chance of having an end date
    THEN date('now', '-' || (n % 6) || ' months', '+' || (n % 15) || ' days')
    ELSE NULL
  END AS end_date,
  CASE WHEN n % 10 > 2 -- 80% chance of having a project name
    THEN 'Project ' || substr(hex(e.id), 1, 1) || '-' || n 
    ELSE NULL
  END AS project_name,
  CASE (n % 4)
    WHEN 0 THEN 'pending'
    WHEN 1 THEN 'active'
    WHEN 2 THEN 'completed'
    WHEN 3 THEN 'cancelled'
  END AS status,
  CASE 
    WHEN (n % 4) = 2 THEN (n % 5) + 1 -- Feedback score only for completed (1-5)
    ELSE NULL
  END AS feedback_score,
  CASE WHEN n % 10 > 5 -- 50% chance of having notes
    THEN 'Notes for engagement ' || n || ' with expert ID ' || e.id
    ELSE NULL
  END AS notes,
  date('now', '-' || (n % 6) || ' months', '-' || (n % 30) || ' days') AS created_at
FROM 
  expert_ids e
CROSS JOIN
  range
WHERE 
  (e.id + n) % 5 > 0 -- Creates 1-4 engagements per expert
LIMIT 
  200; -- Limit to 200 total engagements

-- Update statistics table to reflect changes
-- The system uses a key-value approach with JSON in the stat_value field
INSERT OR REPLACE INTO system_statistics 
  (stat_key, stat_value, last_updated)
VALUES 
  ('expert_counts', json_object(
    'bahraini_count', (SELECT COUNT(*) FROM experts WHERE is_bahraini = 1),
    'non_bahraini_count', (SELECT COUNT(*) FROM experts WHERE is_bahraini = 0),
    'total_experts', (SELECT COUNT(*) FROM experts),
    'active_experts', (SELECT COUNT(*) FROM experts WHERE is_available = 1)
  ), datetime('now')),
  ('engagement_stats', json_object(
    'total_engagements', (SELECT COUNT(*) FROM expert_engagements),
    'completed_engagements', (SELECT COUNT(*) FROM expert_engagements WHERE status = 'completed'),
    'active_engagements', (SELECT COUNT(*) FROM expert_engagements WHERE status = 'active')
  ), datetime('now'));
  
-- Print summary
SELECT 'Added ' || COUNT(*) || ' engagements to database' FROM expert_engagements;