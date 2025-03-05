package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// ImportExperts reads the experts CSV file and imports the data into the database
func ImportExperts(db *sql.DB, csvFilePath string) error {
	// Open the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header row
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading CSV header: %w", err)
	}

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV records: %w", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Map to store unique general areas
	areaMap := make(map[string]int64)

	// First pass: Insert unique general areas
	for _, record := range records {
		generalArea := getFieldByHeaderName(header, record, "General Area")
		if generalArea != "" && generalArea != "No Rating" {
			// Check if this area already exists
			var areaID int64
			err := tx.QueryRow("SELECT id FROM expert_areas WHERE name = ?", generalArea).Scan(&areaID)
			if err == sql.ErrNoRows {
				// Insert new area
				result, err := tx.Exec(
					"INSERT INTO expert_areas (name, created_at) VALUES (?, ?)",
					generalArea, time.Now(),
				)
				if err != nil {
					return fmt.Errorf("error inserting general area '%s': %w", generalArea, err)
				}

				newID, err := result.LastInsertId()
				if err != nil {
					return fmt.Errorf("error getting last insert ID: %w", err)
				}
				areaMap[generalArea] = newID
			} else if err != nil {
				return fmt.Errorf("error checking for existing area: %w", err)
			} else {
				areaMap[generalArea] = areaID
			}
		}
	}

	// Create a map for specialized areas
	specializedAreaMap := make(map[string]int64)

	// Insert unique specialized areas
	for _, record := range records {
		specializedArea := getFieldByHeaderName(header, record, "Specialised Area")
		if specializedArea != "" {
			// Check if this area already exists
			var areaID int64
			err := tx.QueryRow("SELECT id FROM expert_areas WHERE name = ?", specializedArea).Scan(&areaID)
			if err == sql.ErrNoRows {
				// Insert new area
				result, err := tx.Exec(
					"INSERT INTO expert_areas (name, created_at) VALUES (?, ?)",
					specializedArea, time.Now(),
				)
				if err != nil {
					return fmt.Errorf("error inserting specialized area '%s': %w", specializedArea, err)
				}

				newID, err := result.LastInsertId()
				if err != nil {
					return fmt.Errorf("error getting last insert ID: %w", err)
				}
				specializedAreaMap[specializedArea] = newID
			} else if err != nil {
				return fmt.Errorf("error checking for existing area: %w", err)
			} else {
				specializedAreaMap[specializedArea] = areaID
			}
		}
	}

	// Second pass: Insert experts
	for _, record := range records {
		expertID := getFieldByHeaderName(header, record, "ID")
		name := getFieldByHeaderName(header, record, "Name")
		designation := getFieldByHeaderName(header, record, "Designation")
		institution := getFieldByHeaderName(header, record, "Institution")
		isBahraini := convertYesNoToBoolean(getFieldByHeaderName(header, record, "BH"))
		
		// Set nationality based on BH field
		nationality := "Unknown"
		if isBahraini {
			nationality = "Bahraini"
		} else {
			bhValue := getFieldByHeaderName(header, record, "BH")
			if strings.ToLower(bhValue) == "no" {
				nationality = "Non-Bahraini"
			}
		}
		
		isAvailable := convertYesNoToBoolean(getFieldByHeaderName(header, record, "Available"))
		rating := getFieldByHeaderName(header, record, "Rating")
		role := getFieldByHeaderName(header, record, "Validator/ Evaluator")
		employmentType := getFieldByHeaderName(header, record, "Academic/Employer")
		generalArea := getFieldByHeaderName(header, record, "General Area")
		specializedArea := getFieldByHeaderName(header, record, "Specialised Area")
		isTrained := convertYesNoToBoolean(getFieldByHeaderName(header, record, "Trained"))
		cvPath := "" // CV field is empty in the CSV
		phone := getFieldByHeaderName(header, record, "Phone")
		email := getFieldByHeaderName(header, record, "Email")
		isPublished := convertYesNoToBoolean(getFieldByHeaderName(header, record, "Published"))

		// Insert expert record
		result, err := tx.Exec(
			`INSERT INTO experts (
				expert_id, name, designation, institution, is_bahraini, nationality, is_available,
				rating, role, employment_type, general_area, specialized_area,
				is_trained, cv_path, phone, email, is_published, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			expertID, name, designation, institution, isBahraini, nationality, isAvailable,
			rating, role, employmentType, generalArea, specializedArea,
			isTrained, cvPath, phone, email, isPublished, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("error inserting expert '%s': %w", name, err)
		}

		expertDBID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting expert insert ID: %w", err)
		}

		// Insert expert general area specialization if applicable
		if areaID, ok := areaMap[generalArea]; ok && generalArea != "" {
			_, err = tx.Exec(
				"INSERT INTO expert_specializations (expert_id, area_id, created_at) VALUES (?, ?, ?)",
				expertDBID, areaID, time.Now(),
			)
			if err != nil {
				return fmt.Errorf("error inserting expert general area specialization: %w", err)
			}
		}

		// Insert expert specialized area if applicable
		if areaID, ok := specializedAreaMap[specializedArea]; ok && specializedArea != "" {
			_, err = tx.Exec(
				"INSERT INTO expert_specializations (expert_id, area_id, created_at) VALUES (?, ?, ?)",
				expertDBID, areaID, time.Now(),
			)
			if err != nil {
				return fmt.Errorf("error inserting expert specialized area: %w", err)
			}
		}
		
		// Map to ISCED classification
		err = mapExpertToISCED(tx, expertDBID, generalArea, specializedArea, designation)
		if err != nil {
			log.Printf("Warning: Failed to map expert %s to ISCED: %v", name, err)
			// Continue processing other experts even if ISCED mapping fails for some
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	log.Printf("Successfully imported %d experts", len(records))
	return nil
}

// mapExpertToISCED maps an expert to appropriate ISCED classifications
func mapExpertToISCED(tx *sql.Tx, expertID int64, generalArea, specializedArea, designation string) error {
	var iscedFieldID, iscedLevelID int64
	
	// Map general area to ISCED field
	generalAreaLower := strings.ToLower(generalArea)
	specializedAreaLower := strings.ToLower(specializedArea)
	
	// First try to find exact match by broad field name
	err := tx.QueryRow(`
		SELECT id FROM isced_fields 
		WHERE LOWER(broad_name) LIKE ? 
		LIMIT 1`, 
		"%"+generalAreaLower+"%").Scan(&iscedFieldID)
		
	// If no match, try to map based on keywords
	if err == sql.ErrNoRows {
		iscedFieldID, err = mapAreaToISCEDField(tx, generalAreaLower, specializedAreaLower)
		if err != nil {
			return fmt.Errorf("failed to map area to ISCED field: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error finding ISCED field: %w", err)
	}
	
	// Map designation to ISCED level
	iscedLevelID, err = mapDesignationToISCEDLevel(tx, designation)
	if err != nil {
		return fmt.Errorf("failed to map designation to ISCED level: %w", err)
	}
	
	// Update expert record with ISCED classifications
	_, err = tx.Exec(
		"UPDATE experts SET isced_field_id = ?, isced_level_id = ? WHERE id = ?",
		iscedFieldID, iscedLevelID, expertID,
	)
	if err != nil {
		return fmt.Errorf("error updating expert ISCED classifications: %w", err)
	}
	
	return nil
}

// mapAreaToISCEDField maps general and specialized areas to ISCED field based on keywords
func mapAreaToISCEDField(tx *sql.Tx, generalArea, specializedArea string) (int64, error) {
	// Default to "Generic programmes and qualifications" (00)
	var defaultFieldID int64
	err := tx.QueryRow(`SELECT id FROM isced_fields WHERE broad_code = '00' LIMIT 1`).Scan(&defaultFieldID)
	if err != nil {
		return 0, fmt.Errorf("error finding default ISCED field: %w", err)
	}
	
	// Map of keywords to ISCED field codes
	keywordMapping := map[string]string{
		// ISCED 00 - Generic programmes and qualifications
		"generic": "00",
		
		// ISCED 01 - Education
		"education": "01",
		"teaching": "01",
		"instructor": "01",
		"training": "01",
		
		// ISCED 02 - Arts and humanities
		"art": "02",
		"design": "02",
		"humanities": "02",
		"language": "02",
		"literature": "02",
		"philosophy": "02",
		"history": "02",
		"fashion": "02",
		"graphic": "02",
		"interior design": "02",
		"visual design": "02",
		
		// ISCED 03 - Social sciences, journalism and information
		"social science": "03",
		"journalism": "03",
		"media": "03",
		"communication": "03",
		"public relations": "03",
		"sociology": "03",
		"psychology": "03",
		"political": "03",
		
		// ISCED 04 - Business, administration and law
		"business": "04",
		"administration": "04",
		"management": "04",
		"marketing": "04",
		"finance": "04",
		"accounting": "04",
		"audit": "04",
		"banking": "04",
		"economics": "04",
		"human resource": "04",
		"law": "04",
		"legal": "04",
		"compliance": "04",
		"insurance": "04",
		"logistics": "04",
		"project management": "04",
		
		// ISCED 05 - Natural sciences, mathematics and statistics
		"science": "05",
		"mathematics": "05",
		"statistics": "05",
		"physics": "05",
		"chemistry": "05",
		"biology": "05",
		"environment": "05",
		
		// ISCED 06 - Information and Communication Technologies
		"information technology": "06",
		"computer": "06",
		"software": "06",
		"programming": "06",
		"database": "06",
		"network": "06",
		"web": "06",
		"multimedia": "06",
		"it": "06",
		"artificial intelligence": "06",
		
		// ISCED 07 - Engineering, manufacturing and construction
		"engineering": "07",
		"mechanical": "07",
		"electrical": "07",
		"electronic": "07",
		"civil": "07",
		"chemical": "07",
		"architectural": "07",
		"architecture": "07",
		"manufacturing": "07",
		"construction": "07",
		
		// ISCED 08 - Agriculture, forestry, fisheries and veterinary
		"agriculture": "08",
		"forestry": "08",
		"fisheries": "08",
		"veterinary": "08",
		
		// ISCED 09 - Health and welfare
		"health": "09",
		"medicine": "09",
		"medical": "09",
		"nursing": "09",
		"pharmacy": "09",
		"dental": "09",
		"radiology": "09",
		"physiotherapy": "09",
		"welfare": "09",
		
		// ISCED 10 - Services
		"services": "10",
		"hospitality": "10",
		"tourism": "10",
		"transportation": "10",
		"aviation": "10",
		"security": "10",
		"safety": "10",
	}
	
	// First check general area
	for keyword, code := range keywordMapping {
		if strings.Contains(generalArea, keyword) {
			var fieldID int64
			err := tx.QueryRow(`SELECT id FROM isced_fields WHERE broad_code = ? LIMIT 1`, code).Scan(&fieldID)
			if err == nil {
				return fieldID, nil
			}
		}
	}
	
	// Then check specialized area (higher priority)
	for keyword, code := range keywordMapping {
		if strings.Contains(specializedArea, keyword) {
			var fieldID int64
			err := tx.QueryRow(`SELECT id FROM isced_fields WHERE broad_code = ? LIMIT 1`, code).Scan(&fieldID)
			if err == nil {
				return fieldID, nil
			}
		}
	}
	
	// Return default if no match found
	return defaultFieldID, nil
}

// mapDesignationToISCEDLevel maps designation to ISCED level
func mapDesignationToISCEDLevel(tx *sql.Tx, designation string) (int64, error) {
	designationLower := strings.ToLower(designation)
	
	// Get all ISCED levels
	rows, err := tx.Query("SELECT id, code FROM isced_levels ORDER BY id")
	if err != nil {
		return 0, fmt.Errorf("error querying ISCED levels: %w", err)
	}
	defer rows.Close()
	
	// Store levels in a map
	levelMap := make(map[string]int64)
	var doctoralID, mastersID int64
	
	for rows.Next() {
		var id int64
		var code string
		if err := rows.Scan(&id, &code); err != nil {
			return 0, fmt.Errorf("error scanning ISCED level: %w", err)
		}
		
		levelMap[code] = id
		
		// Store common levels
		switch code {
		case "8": // Doctoral
			doctoralID = id
		case "7": // Master's
			mastersID = id
		case "6": // Bachelor's
			levelMap["6"] = id
		}
	}
	
	// Map based on designation
	if strings.Contains(designationLower, "dr.") || 
	   strings.Contains(designationLower, "prof.") || 
	   strings.Contains(designationLower, "professor") {
		return doctoralID, nil
	} else if strings.Contains(designationLower, "mr.") ||
              strings.Contains(designationLower, "mrs.") ||
              strings.Contains(designationLower, "ms.") ||
              strings.Contains(designationLower, "miss") {
		return mastersID, nil
	}
	
	// Default to master's level if not specified
	return mastersID, nil
}

// Helper function to get a field value by header name
func getFieldByHeaderName(header []string, record []string, headerName string) string {
	for i, h := range header {
		if h == headerName && i < len(record) {
			return strings.TrimSpace(record[i])
		}
	}
	return ""
}

// Helper function to convert "Yes"/"No" to boolean
func convertYesNoToBoolean(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "yes" || value == "true" || value == "y"
}