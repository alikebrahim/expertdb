#!/bin/bash

# Database Schema Verification Script
# This script checks the database schema and verifies that all required migrations have been applied

DATABASE_PATH="./backend/db/sqlite/expertdb.sqlite"

if [ ! -f "$DATABASE_PATH" ]; then
    echo "Error: Database file not found at $DATABASE_PATH"
    exit 1
fi

echo "=== Checking Database Schema ==="
echo

# Check if migration_versions table exists
echo "Checking migration tracking table..."
MIGRATION_TABLE_COUNT=$(sqlite3 "$DATABASE_PATH" "SELECT count(*) FROM sqlite_master WHERE type='table' AND name='migration_versions';")

if [ "$MIGRATION_TABLE_COUNT" -eq "0" ]; then
    echo "Warning: migration_versions table does not exist. Migrations are not being tracked."
else
    echo "Migration tracking table exists."
    echo "Applied migrations:"
    sqlite3 "$DATABASE_PATH" "SELECT filename FROM migration_versions ORDER BY filename;" | sed 's/^/- /'
fi

echo

# Check expert_requests table schema
echo "Checking expert_requests table schema..."
echo "Columns in expert_requests table:"
sqlite3 "$DATABASE_PATH" "PRAGMA table_info(expert_requests);" | awk -F'|' '{print "- " $2 " (" $3 ")"}'

echo

# Check for specific columns that have been problematic
echo "Checking for specific columns..."
for COLUMN in "rejection_reason" "biography"; do
    COLUMN_EXISTS=$(sqlite3 "$DATABASE_PATH" "SELECT count(*) FROM pragma_table_info('expert_requests') WHERE name='$COLUMN';")
    if [ "$COLUMN_EXISTS" -eq "0" ]; then
        echo "Warning: Column '$COLUMN' does not exist in expert_requests table!"
    else
        echo "Column '$COLUMN' exists in expert_requests table."
    fi
done

echo

# Check for data in key tables
echo "=== Checking Table Data ==="
echo

tables=("users" "experts" "expert_requests" "isced_levels" "isced_fields")
for table in "${tables[@]}"; do
    count=$(sqlite3 "$DATABASE_PATH" "SELECT count(*) FROM $table;")
    echo "$table: $count records"
done

echo

# Verify migrations directory
echo "=== Checking Migration Files ==="
echo

ls -1 "./backend/db/migrations/sqlite/" | sort | while read -r migration; do
    if [ -n "$migration" ]; then
        applied=$(sqlite3 "$DATABASE_PATH" "SELECT count(*) FROM migration_versions WHERE filename='$migration';")
        if [ "$applied" -eq "0" ]; then
            echo "Warning: Migration '$migration' has not been applied!"
        else
            echo "Migration '$migration' has been applied."
        fi
    fi
done

echo
echo "=== Schema Verification Complete ==="