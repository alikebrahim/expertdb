#!/bin/bash

DB_PATH="./backend/db/sqlite/expertdb.sqlite"

if [ ! -f "$DB_PATH" ]; then
  echo "Error: Database file not found at $DB_PATH"
  echo "Make sure you're running this script from the project root directory"
  exit 1
fi

echo "Populating database with mock engagement data..."
sqlite3 "$DB_PATH" < populate_mock_data.sql

# Get counts for verification
ENGAGEMENT_COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM expert_engagements")
echo "Total engagements in database: $ENGAGEMENT_COUNT"

# Show distribution by type
echo "Engagement distribution by type:"
sqlite3 "$DB_PATH" "SELECT engagement_type, COUNT(*) FROM expert_engagements GROUP BY engagement_type"

# Show distribution by status
echo "Engagement distribution by status:"
sqlite3 "$DB_PATH" "SELECT status, COUNT(*) FROM expert_engagements GROUP BY status"

# Show statistics records
echo "Statistics records:"
sqlite3 "$DB_PATH" "SELECT stat_key, last_updated FROM system_statistics"

echo "Done!"