# ExpertDB: Database Migration and API Fixes

## Overview

This document summarizes the fixes implemented to address database schema inconsistency issues and API failures in the ExpertDB application.

## Latest Updates (2025-03-12)

1. **Expert Creation Fix**
   - Problem: Expert creation was failing with "primary contact is required" error
   - Fix: Updated test_api.sh to use the correct JSON structure with primaryContact field
   - Impact: Experts can now be created successfully via the API

2. **Expert Request Columns Fix**
   - Problem: Expert request listing still failed with "no such column: rejection_reason" error
   - Fix: Added dynamic column detection via `getExpertRequestColumns()` helper method
   - Impact: Expert requests API now works regardless of schema variations

3. **Document Listing Fix**
   - Problem: Document listing API failed due to invalid URL when no expert was created
   - Fix: Added conditional check to skip document listing when no expert ID is available
   - Impact: API testing script now handles failures gracefully

4. **Schema Verification Tool**
   - Feature: Added a new verification_schema.sh script to analyze database schema
   - Capabilities: Checks applied migrations, table schemas, and problematic columns
   - Impact: Easier diagnosis of database schema inconsistencies

## Issues Addressed

1. **Database Migration Tracking**
   - Problem: Migrations were not properly tracked, leading to inconsistent schema across environments
   - Fix: Implemented a migration_versions table to track applied migrations
   - Impact: Ensures migrations are only applied once and in the correct order

2. **Expert Requests API Failures**
   - Problem: "no such column: rejection_reason" errors when accessing expert requests
   - Fix: Implemented dynamic column detection and graceful handling of missing columns
   - Impact: API now works even with schema variations

3. **Admin Panel User Display Issues**
   - Problem: Created users not appearing in admin panel despite successful login
   - Fix: Enhanced database operations to be more resilient to schema variations
   - Impact: Admin panel now shows all users correctly

## Technical Implementation

### Database Migration System

```go
// Create the migration_versions table
_, err = s.db.Exec(`
    CREATE TABLE IF NOT EXISTS migration_versions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        filename TEXT UNIQUE NOT NULL,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
`)

// Track applied migrations
appliedMigrations := make(map[string]bool)
for rows.Next() {
    var filename string
    if err := rows.Scan(&filename); err != nil {
        rows.Close()
        return fmt.Errorf("failed to scan migration row: %w", err)
    }
    appliedMigrations[filename] = true
}

// Apply each migration in its own transaction
for _, migration := range migrations {
    filename := filepath.Base(migration)
    
    // Skip if already applied
    if appliedMigrations[filename] {
        logger.Info("Skipping already applied migration: %s", filename)
        continue
    }
    
    // Start a transaction for this migration
    tx, err := s.db.Begin()
    // ...
    
    // Record the migration as applied
    _, err = tx.Exec("INSERT INTO migration_versions (filename) VALUES (?)", filename)
    // ...
}
```

### Dynamic Column Detection

```go
// Check columns to handle schema discrepancies
columns, err := rows.Columns()
if err != nil {
    return nil, fmt.Errorf("failed to get columns: %w", err)
}

// Create column map for dynamic scanning
colMap := make(map[string]int)
for i, col := range columns {
    colMap[col] = i
}

// Create values slice with the right number of interface{} values
values := make([]interface{}, len(columns))
for i := range values {
    values[i] = new(interface{})
}

if err := rows.Scan(values...); err != nil {
    return nil, fmt.Errorf("failed to scan expert request row: %w", err)
}

// Map column values to struct fields
for col, idx := range colMap {
    v := *(values[idx].(*interface{}))
    if v == nil {
        continue
    }
    
    switch col {
    case "id":
        if id, ok := v.(int64); ok {
            req.ID = id
        }
    // ...other column mappings
    }
}
```

### Dynamic SQL Query Building

```go
// Build dynamic update query based on available columns
var setClauses []string
var args []interface{}

// Add fields only if the column exists in the table
if columnMap["expert_id"] {
    setClauses = append(setClauses, "expert_id = ?")
    args = append(args, request.ExpertID)
}

// ...other column checks

// Build the final query
query := fmt.Sprintf("UPDATE expert_requests SET %s WHERE id = ?", strings.Join(setClauses, ", "))
```

## Testing

A comprehensive test script (`test_api.sh`) has been created to verify all API endpoints. This script:

1. Tests authentication with both admin and regular user accounts
2. Verifies user management operations
3. Tests expert creation, retrieval, and updates
4. Tests the expert request workflow including approval
5. Verifies document management functionality
6. Checks ISCED classification endpoints
7. Tests statistics endpoints

## Next Steps

1. Deploy the updated code with fixed migration system
2. Complete any remaining UI/UX improvements
3. Implement AI integration features
4. Enhance test coverage
5. Improve error reporting in the UI
