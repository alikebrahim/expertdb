# CSV Import Investigation (2025-03-12)

## Original Task

- Remove in-memory DB implementation from storage.go
- Move to goose-based migrations instead of automated schema initialization
- Update CSV import functionality to work with this new approach

## Implementation Changes Made

1. **Removed in-memory DB initialization from storage.go**
   - Removed the initSchema() function that auto-created tables
   - Simplified SQLiteStore to only handle database connection, not schema

2. **Updated server.go to verify schema exists**
   - Added verifyDatabaseSchema() function to check if required tables exist
   - Server now fails with clear error if database schema is incomplete

3. **Fixed migration files with schema issues**
   - Fixed 0005_create_foreign_keys.sql migration which had mismatched columns
   - Fixed 0012_create_statistics_table.sql to avoid duplicate column definitions

4. **Updated CSV import tool**
   - Modified import_csv/main.go to include all required columns in schema
   - Added biography field with default value "Expert BIO"
   - Set isced_level_id, isced_field_id, and original_request_id to NULL

## Current Issue

The CSV import tool fails with the following error:
```
Failed to import experts: error inserting expert 'Aalla Hajih': UNIQUE constraint failed: experts.expert_id
```

This happens when trying to import the first record after the header (E002 - Aalla Hajih).

## Investigation Steps Taken

1. **Verified database is empty**
   - Checked with: `sqlite3 ./db/sqlite/expertdb.sqlite "SELECT COUNT(*) FROM experts;"`
   - Result: 0 rows, confirmed the database is empty

2. **Checked for duplicate expert_ids in CSV**
   - Ran: `cut -d, -f1 ./experts.csv | sort | uniq -c | sort -rn | head -10`
   - Found only one entry for each expert_id, including E002
   - Found 5 malformed lines with quotes that might be causing issues

3. **Examined the specific problematic record**
   - Looked at E002 record: `grep -n "E002" ./experts.csv`
   - Found only one record with ID E002 (line 3)

4. **Verified database schema**
   - Checked table definition: `sqlite3 ./db/sqlite/expertdb.sqlite ".schema experts"`
   - Found mismatch between columns in schema and import script (fixed)

5. **Examined the migrations**
   - Checked all migrations to ensure they run in proper sequence
   - Successfully ran all migrations from 0001 through 0012

## Current State

- All migrations run successfully
- Database schema is correctly created
- CSV import script has been updated to match table structure
- Import still fails with UNIQUE constraint violation

## Code Changes

1. Updated storage.go:
   - Removed in-memory DB implementation
   - Removed automatic schema initialization

2. Updated server.go:
   - Added schema verification at startup

3. Fixed migrations:
   - Fixed column definitions in 0005_create_foreign_keys.sql
   - Fixed duplicate column addition in 0012_create_statistics_table.sql

4. Updated import_csv/main.go:
   - Added missing columns in INSERT statement (biography, isced_level_id, isced_field_id, original_request_id)
   - Set default biography to "Expert BIO"

## Potential Solutions to Try

1. **Transaction Issues**: 
   - Check if there might be a transaction issue causing partial commits
   - Try modifying the import tool to use smaller transactions

2. **Sqlite Database Corruption**:
   - Consider completely removing the database file
   - Retry with a fresh goose migration run

3. **Hidden Data**: 
   - Check if there might be hidden data from a previous import
   - Look for any migrations or processes that might be inserting data

4. **CSV Parsing Issues**:
   - Examine the CSV parsing logic for any issues with line breaks
   - Try with a simplified subset of the data

5. **UNIQUE Constraint Debug**:
   - Add logging to print exactly what values are being inserted
   - Modify the import tool to skip duplicate expert_ids instead of failing

6. **Inspect Raw Database**:
   - Use sqlite3 directly to examine if anything exists in the experts table despite COUNT(*) showing 0
   - Check if any triggers or views might be causing issues

## Commands to Resume Work

1. **Check database content**:
   ```
   sqlite3 ./db/sqlite/expertdb.sqlite "SELECT * FROM experts WHERE expert_id = 'E002';"
   ```

2. **Rebuild import tool**:
   ```
   cd /home/alikebrahim/dev/expertdb_grok/backend && go build -o import_csv ./cmd/import_csv
   ```

3. **Run migrations from scratch**:
   ```
   cd /home/alikebrahim/dev/expertdb_grok/backend && rm -f ./db/sqlite/expertdb.sqlite && goose -dir db/migrations/sqlite sqlite3 ./db/sqlite/expertdb.sqlite up
   ```

4. **Run import with the fixed tool**:
   ```
   cd /home/alikebrahim/dev/expertdb_grok/backend && ./import_csv -csv ./experts.csv -db ./db/sqlite/expertdb.sqlite
   ```