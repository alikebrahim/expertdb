# ExpertDB Instructions

## Database Recreation and CSV Import

The CSV import approach is a practical solution for this project. Here's how it works:

### Implementation

The system uses the `encoding/csv` package to import data from a CSV file (`experts.csv`) into the database. This process:

1. Reads the CSV file with headers
2. Processes each record in a transaction for data safety
3. Extracts expert information including newly added fields (role, employmentType, generalArea, cvPath, biography)
4. Creates records in the database with proper foreign key mappings

### Steps to Recreate Database from CSV

1. **Drop and recreate database**: 
   ```bash
   rm -f db/sqlite/expertdb.sqlite
   ```

2. **Run the application to initialize schema**: 
   ```bash
   go run .
   ```
   This will:
   - Create all necessary tables with the updated schema including the new biography field
   - Apply migrations to ensure the schema is up-to-date

3. **Import CSV data**:
   ```bash
   go run cmd/import_csv/main.go experts.csv
   ```

### Why This Approach Makes Sense

1. **Simplicity**: Direct recreation is much simpler than managing complex migration scripts for the development phase
2. **Small dataset**: The experts dataset is small enough to quickly reimport
3. **Development mode**: Since we're in development, frequent schema changes are expected
4. **Data consistency**: Ensures all data follows the latest schema patterns
5. **Quick validation**: Allows rapid testing of schema changes

## Feedback from Claude

I've reviewed and updated the codebase to handle the new fields in the `CreateExpertRequest` and `Expert` structs. Here's what I found and fixed:

1. Updated the `NewExpert` method to properly map all new fields from `CreateExpertRequest` to `Expert`:
   - Role
   - EmploymentType
   - GeneralArea
   - CVPath
   - Biography
   - IsBahraini

2. Added proper validation for new required fields in `ValidateCreateExpertRequest`:
   - Role is now required and validated against a list of valid values
   - EmploymentType is now required and validated against a list of valid values
   - GeneralArea is now required

3. Created a database migration to add the biography column to:
   - The `experts` table
   - The `expert_requests` table

4. Updated the in-memory database schema to include the biography field in both tables

5. The database can now be recreated from scratch with the updated schema.

Note: The CSV import functionality should work as expected with the new fields. The CSV file itself may need to be updated to include the biography field if you want to populate it during import.