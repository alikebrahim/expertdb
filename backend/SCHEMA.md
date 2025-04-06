# Database Schema Simplification: Implementation Report

## Migration Files Changes

### 1. 0001_create_expert_table_up.sql
**NOTES:**
- "-- NOTE: below indexing is not necessary anymore, edit in the drop section below accordingly" - regarding ISCED indexes

**ACTIONS:**
- Changed `general_area` column from TEXT to INTEGER
- Removed `isced_level_id` and `isced_field_id` columns
- Removed ISCED-related indexes
- Updated DROP statements to match the new structure

**AFFECTED BACKEND LOGIC:**
- [Expert struct] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [51-78]
- [CreateExpert method] - [/home/alikebrahim/dev/expertdb_grok/backend/expert_storage.go] - [9-61]
- [GetExpert method] - [/home/alikebrahim/dev/expertdb_grok/backend/expert_operations.go] - [29-146]
- [UpdateExpert method] - [/home/alikebrahim/dev/expertdb_grok/backend/expert_operations.go] - [163-222]
- [NewExpert function] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [231-259]

**RECOMMENDATIONS:**
- Update Expert struct to change GeneralArea from string to int64
- Remove ISCEDLevel and ISCEDField fields from Expert struct
- Add GeneralAreaName string field for frontend compatibility
- Update all methods that handle Expert creation/updates to work with the new integer field

### 2. 0002_create_expert-request_table.sql
**NOTES:**
- "-- Path to the CV file NOTE: This is better be replaced with expert_documents(id)" - regarding CV path

**ACTIONS:**
- Changed `general_area` column from TEXT to INTEGER
- Added an index for general_area column
- Updated DROP statements to include the new index

**AFFECTED BACKEND LOGIC:**
- [ExpertRequest struct] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [87-111]
- [CreateExpertRequest struct] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [10-23]
- [CreateExpertRequest handler] - [Not directly visible in provided files]

**RECOMMENDATIONS:**
- Update ExpertRequest struct to change GeneralArea from string to int64
- Add GeneralAreaName string field for frontend compatibility
- Update CreateExpertRequest struct to accept generalAreaId instead of generalArea text

### 3. 0004_create_expert_areas_table.sql
**NOTES:**
- "-- NOTE: These tables are to be removed" - regarding expert_specializations junction table

**ACTIONS:**
- Kept `expert_areas` table as instructed
- Removed `expert_specializations` junction table completely
- Added direct population of expert_areas table with 29 predefined specializations

**AFFECTED BACKEND LOGIC:**
- [Area struct] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [80-84]
- [GetExpert method] - [/home/alikebrahim/dev/expertdb_grok/backend/expert_operations.go] - [106-122]
- [DeleteExpert method] - [/home/alikebrahim/dev/expertdb_grok/backend/expert_operations.go] - [270-271]

**RECOMMENDATIONS:**
- Create new method GetExpertAreas() to fetch all expert areas
- Remove area mapping logic in GetExpert since there's no junction table anymore
- Add API endpoint for retrieving all expert areas

### 4. 0005_create_foreign_keys.sql
**NOTES:**
- "-- NOTE: The below tables appear to be redundant. Assess for removal."

**ACTIONS:**
- Modified `experts_temp` table to:
  - Change general_area to INTEGER
  - Remove isced_level_id and isced_field_id fields
  - Add foreign key constraint for general_area referencing expert_areas
- Modified `expert_requests_temp` table to:
  - Change general_area to INTEGER 
  - Add foreign key constraint for general_area referencing expert_areas
- Added data migration logic to transform text general_area to integer IDs
- Updated index recreation statements

**AFFECTED BACKEND LOGIC:**
- [SQLiteStore.CreateExpert] - [/home/alikebrahim/dev/expertdb_grok/backend/expert_storage.go] - [9-61]
- [SQLiteStore.UpdateExpert] - [/home/alikebrahim/dev/expertdb_grok/backend/expert_operations.go] - [163-222]
- [Expert creation from request] - [Various handlers]

**RECOMMENDATIONS:**
- Update query construction in CreateExpert and UpdateExpert to handle integer general_area
- Modify frontend form submission to send general_area as ID not text
- Add validation to ensure general_area ID exists in expert_areas table

### 5. Migration File Cleanup
**NOTES:**
- ISCED-related migration files (0007_add_isced_classification.sql and 0008_map_experts_to_isced.sql) were completely removed as they're no longer needed
- Remaining migration files were renumbered for sequential consistency
- Renamed files to better reflect their purpose

**ACTIONS:**
- Removed migration files:
  - 0007_add_isced_classification.sql
  - 0008_map_experts_to_isced.sql
- Renamed and resequenced migration files for consistency:
  - 0009_create_expert_documents_table.sql → 0006_create_expert_documents_table.sql
  - 0010_create_expert_engagements_table.sql → 0007_create_expert_engagements_table.sql
  - 0012_create_statistics_table.sql → 0008_create_statistics_table.sql
- Removed redundant migration file:
  - 0006_csv_import_script.sql (removed entirely as it duplicated expert_areas population in 0004)

**AFFECTED BACKEND LOGIC:**
- [ISCEDLevel struct] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [32-37]
- [ISCEDField struct] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [40-49]
- [Expert struct iscedLevel and iscedField fields] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [69-71]
- [ISCED retrieval in GetExpert] - [/home/alikebrahim/dev/expertdb_grok/backend/expert_operations.go] - [74-104]
- [Statistics struct ExpertsByISCEDField field] - [/home/alikebrahim/dev/expertdb_grok/backend/types.go] - [158]
- [Statistics calculation logic] - [Not directly visible in provided files]

**RECOMMENDATIONS:**
- Remove ISCEDLevel and ISCEDField struct definitions
- Remove ISCED-related fields from Expert struct
- Remove ISCED retrieval code in GetExpert method
- Update Statistics struct to remove ExpertsByISCEDField
- Remove any statistics calculation related to ISCED fields
- Update frontend dashboard components that displayed ISCED statistics

## Required Go Code Changes

### 1. types.go
- Update Expert struct:
  - Change GeneralArea from string to int64
  - Add GeneralAreaName string field
  - Remove ISCEDLevel and ISCEDField fields and their references
- Update ExpertRequest struct:
  - Change GeneralArea from string to int64
  - Add GeneralAreaName string field
- Update CreateExpertRequest struct:
  - Change GeneralArea from string to int64 (renamed to GeneralAreaID)
- Update Statistics struct:
  - Remove ExpertsByISCEDField field
- Remove ISCEDLevel and ISCEDField struct definitions

### 2. expert_storage.go
- Update CreateExpert method:
  - Remove ISCED fields from SQL query
  - Remove ISCED null handling logic

### 3. expert_operations.go
- Update GetExpert method:
  - Remove ISCED field processing (Lines 72-104)
  - Add join to expert_areas table to get area name
  - Remove expert_specializations related queries
- Update UpdateExpert method:
  - Remove ISCED fields from SQL query
  - Remove ISCED null handling logic
- Update DeleteExpert method:
  - Remove expert_specializations deletion logic

### 4. api.go
- Add endpoint for retrieving expert areas

### 5. New methods needed
- Add GetExpertAreas method to retrieve list of all expert areas

## Implementation Recommendations

1. Add new SQL methods:
```go
// GetExpertAreaByID retrieves a single expert area by ID
func (s *SQLiteStore) GetExpertAreaByID(id int64) (*Area, error) {
    var area Area
    err := s.db.QueryRow("SELECT id, name FROM expert_areas WHERE id = ?", id).Scan(&area.ID, &area.Name)
    if err != nil {
        return nil, err
    }
    return &area, nil
}

// GetExpertAreas retrieves all expert areas
func (s *SQLiteStore) GetExpertAreas() ([]Area, error) {
    rows, err := s.db.Query("SELECT id, name FROM expert_areas ORDER BY name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var areas []Area
    for rows.Next() {
        var area Area
        if err := rows.Scan(&area.ID, &area.Name); err != nil {
            return nil, err
        }
        areas = append(areas, area)
    }
    
    return areas, nil
}
```

2. Update Expert struct with backward compatibility:
```go
type Expert struct {
    // Other fields remained same
    GeneralArea     int64        `json:"generalAreaId"`
    GeneralAreaName string       `json:"generalArea"`    // For backward compatibility
    // ISCED fields removed
    // Other fields remained same
}
```

3. Modify SQL queries to join with expert_areas table:
```sql
SELECT e.*, ea.name as general_area_name
FROM experts e
JOIN expert_areas ea ON e.general_area = ea.id
WHERE e.id = ?
```

4. Add API endpoint for expert areas:
```go
router.GET("/api/expert/areas", s.listExpertAreasHandler)
```

This approach ensures a smooth transition while maintaining backward compatibility and introducing the simplified schema.
