package sqlite

import (
	"fmt"
	"time"
	
	"expertdb/internal/domain"
)

// GetStatistics retrieves system-wide statistics
func (s *SQLiteStore) GetStatistics() (*domain.Statistics, error) {
	stats := &domain.Statistics{
		LastUpdated: time.Now(),
	}
	
	// Get total experts count
	var totalExperts int
	err := s.db.QueryRow("SELECT COUNT(*) FROM experts").Scan(&totalExperts)
	if err != nil {
		return nil, fmt.Errorf("failed to count experts: %w", err)
	}
	stats.TotalExperts = totalExperts
	
	// Get active experts count
	var activeExperts int
	err = s.db.QueryRow("SELECT COUNT(*) FROM experts WHERE is_available = 1").Scan(&activeExperts)
	if err != nil {
		return nil, fmt.Errorf("failed to count active experts: %w", err)
	}
	stats.ActiveCount = activeExperts
	
	// Get Bahraini experts percentage
	bahrainiCount, _, err := s.GetExpertsByNationality()
	if err != nil {
		return nil, fmt.Errorf("failed to count experts by nationality: %w", err)
	}
	
	if totalExperts > 0 {
		stats.BahrainiPercentage = float64(bahrainiCount) / float64(totalExperts) * 100
	}
	
	// Get published experts count and ratio
	publishedCount, publishedRatio, err := s.GetPublishedExpertStats()
	if err != nil {
		return nil, fmt.Errorf("failed to count published experts: %w", err)
	}
	stats.PublishedCount = publishedCount
	stats.PublishedRatio = publishedRatio
	
	// Get top areas
	rows, err := s.db.Query(`
		SELECT ea.name as area_name, COUNT(*) as count
		FROM experts e
		JOIN expert_areas ea ON e.general_area = ea.id
		GROUP BY e.general_area
		ORDER BY count DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query top areas: %w", err)
	}
	defer rows.Close()
	
	// Initialize topAreas as empty slice to prevent null in JSON
	topAreas := []domain.AreaStat{}
	for rows.Next() {
		var area domain.AreaStat
		var count int
		if err := rows.Scan(&area.Name, &count); err != nil {
			return nil, fmt.Errorf("failed to scan area row: %w", err)
		}
		area.Count = count
		if totalExperts > 0 {
			area.Percentage = float64(count) / float64(totalExperts) * 100
		}
		topAreas = append(topAreas, area)
	}
	stats.TopAreas = topAreas
	
	// Get engagement statistics
	engagementStats, err := s.GetEngagementStatistics()
	if err != nil {
		// Initialize as empty array on error
		stats.EngagementsByType = []domain.AreaStat{}
	} else {
		stats.EngagementsByType = engagementStats
	}
	
	// Ensure empty array instead of null
	if stats.EngagementsByType == nil {
		stats.EngagementsByType = []domain.AreaStat{}
	}
	
	// Get yearly growth (instead of monthly growth)
	yearlyGrowth, err := s.GetExpertGrowthByYear(5) // Last 5 years
	if err != nil {
		// Initialize as empty array on error
		stats.YearlyGrowth = []domain.GrowthStat{}
	} else {
		stats.YearlyGrowth = yearlyGrowth
	}
	
	// Get most requested experts
	mostRequested := []domain.ExpertStat{} // Initialize as empty slice
	
	rows, err = s.db.Query(`
		SELECT e.expert_id, e.name, COUNT(eng.id) as request_count
		FROM experts e
		JOIN expert_engagements eng ON e.id = eng.expert_id
		GROUP BY e.id
		ORDER BY request_count DESC
		LIMIT 10
	`)
	if err != nil {
		// On error, just use empty array
		stats.MostRequestedExperts = mostRequested
	} else {
		defer rows.Close()
		
		for rows.Next() {
			var stat domain.ExpertStat
			if err := rows.Scan(&stat.ExpertID, &stat.Name, &stat.Count); err != nil {
				// On error, just use empty array
				stats.MostRequestedExperts = mostRequested
				break
			}
			mostRequested = append(mostRequested, stat)
		}
		stats.MostRequestedExperts = mostRequested
	}
	
	// Guarantee it's not null
	if stats.MostRequestedExperts == nil {
		stats.MostRequestedExperts = []domain.ExpertStat{}
	}
	
	return stats, nil
}

// UpdateStatistics updates the statistics in the database
func (s *SQLiteStore) UpdateStatistics(stats *domain.Statistics) error {
	// This is a placeholder method - in this implementation we calculate
	// statistics on-the-fly rather than storing them
	return nil
}

// GetExpertsByNationality retrieves counts of experts by nationality (Bahraini vs non-Bahraini)
func (s *SQLiteStore) GetExpertsByNationality() (int, int, error) {
	var bahrainiCount, nonBahrainiCount int
	
	// Count Bahraini experts
	err := s.db.QueryRow("SELECT COUNT(*) FROM experts WHERE is_bahraini = 1").Scan(&bahrainiCount)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count Bahraini experts: %w", err)
	}
	
	// Count non-Bahraini experts
	err = s.db.QueryRow("SELECT COUNT(*) FROM experts WHERE is_bahraini = 0").Scan(&nonBahrainiCount)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count non-Bahraini experts: %w", err)
	}
	
	return bahrainiCount, nonBahrainiCount, nil
}

// GetEngagementStatistics retrieves statistics about expert engagements
func (s *SQLiteStore) GetEngagementStatistics() ([]domain.AreaStat, error) {
	// Query to analyze engagement distribution by type - restrict to validator and evaluator types
	// and map them to QP and IL application types
	rows, err := s.db.Query(`
		SELECT 
			CASE 
				WHEN engagement_type = 'validator' THEN 'QP (Qualification Placement)'
				WHEN engagement_type = 'evaluator' THEN 'IL (Institutional Listing)'
				ELSE engagement_type 
			END as type_label,
			COUNT(*) as count, 
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count
		FROM expert_engagements
		WHERE engagement_type IN ('validator', 'evaluator')
		GROUP BY engagement_type
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query engagement statistics: %w", err)
	}
	defer rows.Close()
	
	// Initialize stats as empty slice to prevent null in JSON
	stats := []domain.AreaStat{}
	var totalEngagements int
	
	// First, collect all counts
	for rows.Next() {
		var stat domain.AreaStat
		var completedCount int
		if err := rows.Scan(&stat.Name, &stat.Count, &completedCount); err != nil {
			return nil, fmt.Errorf("failed to scan engagement row: %w", err)
		}
		totalEngagements += stat.Count
		stats = append(stats, stat)
	}
	
	// Calculate percentages
	if totalEngagements > 0 {
		for i := range stats {
			stats[i].Percentage = float64(stats[i].Count) / float64(totalEngagements) * 100
		}
	}
	
	return stats, nil
}

// GetExpertGrowthByMonth retrieves statistics about expert growth by month
func (s *SQLiteStore) GetExpertGrowthByMonth(months int) ([]domain.GrowthStat, error) {
	// Default to 12 months if not specified
	if months <= 0 {
		months = 12
	}
	
	// Query to analyze the growth pattern of experts over time
	rows, err := s.db.Query(`
		SELECT 
			strftime('%Y-%m', created_at) as month,
			COUNT(*) as count
		FROM experts
		WHERE created_at >= date('now', '-' || ? || ' months')
		GROUP BY month
		ORDER BY month
	`, months)
	if err != nil {
		return nil, fmt.Errorf("failed to query expert growth: %w", err)
	}
	defer rows.Close()
	
	// Initialize stats as empty slice to prevent null in JSON
	stats := []domain.GrowthStat{}
	var prevCount int
	
	// Process each month
	for rows.Next() {
		var stat domain.GrowthStat
		var monthStr string
		var count int
		if err := rows.Scan(&monthStr, &count); err != nil {
			return nil, fmt.Errorf("failed to scan growth stats row: %w", err)
		}
		
		stat.Period = monthStr
		stat.Count = count
		
		// Calculate growth rate (except for first month)
		if len(stats) > 0 && prevCount > 0 {
			stat.GrowthRate = (float64(count) - float64(prevCount)) / float64(prevCount) * 100
		}
		
		prevCount = count
		stats = append(stats, stat)
	}
	
	// If no data for some months in the range, fill with zeroes for continuity
	if len(stats) < months {
		// Generate a complete list of months
		endDate := time.Now()
		startDate := endDate.AddDate(0, -months, 0)
		
		filledStats := make([]domain.GrowthStat, 0, months)
		
		// Create a map of existing stats for lookup
		existingStats := make(map[string]domain.GrowthStat)
		for _, stat := range stats {
			existingStats[stat.Period] = stat
		}
		
		// Fill in all months
		for m := 0; m < months; m++ {
			currDate := startDate.AddDate(0, m, 0)
			monthStr := fmt.Sprintf("%04d-%02d", currDate.Year(), int(currDate.Month()))
			
			if stat, exists := existingStats[monthStr]; exists {
				filledStats = append(filledStats, stat)
			} else {
				// Add empty stat
				filledStats = append(filledStats, domain.GrowthStat{
					Period: monthStr,
					Count:  0,
				})
			}
		}
		
		// Recalculate growth rates with filled data
		for i := 1; i < len(filledStats); i++ {
			prevCount := filledStats[i-1].Count
			if prevCount > 0 {
				filledStats[i].GrowthRate = (float64(filledStats[i].Count) - float64(prevCount)) / float64(prevCount) * 100
			}
		}
		
		stats = filledStats
	}
	
	return stats, nil
}

// GetExpertGrowthByYear retrieves statistics about expert growth by year
func (s *SQLiteStore) GetExpertGrowthByYear(years int) ([]domain.GrowthStat, error) {
	// Default to 5 years if not specified
	if years <= 0 {
		years = 5
	}
	
	// Query to analyze the yearly growth pattern of experts
	rows, err := s.db.Query(`
		SELECT 
			strftime('%Y', created_at) as year,
			COUNT(*) as count
		FROM experts
		WHERE created_at >= date('now', '-' || ? || ' years')
		GROUP BY year
		ORDER BY year
	`, years)
	if err != nil {
		return nil, fmt.Errorf("failed to query expert yearly growth: %w", err)
	}
	defer rows.Close()
	
	// Initialize stats as empty slice to prevent null in JSON
	stats := []domain.GrowthStat{}
	var prevCount int
	
	// Process each year
	for rows.Next() {
		var stat domain.GrowthStat
		var yearStr string
		var count int
		if err := rows.Scan(&yearStr, &count); err != nil {
			return nil, fmt.Errorf("failed to scan yearly growth stats row: %w", err)
		}
		
		stat.Period = yearStr
		stat.Count = count
		
		// Calculate growth rate (except for first year)
		if len(stats) > 0 && prevCount > 0 {
			stat.GrowthRate = (float64(count) - float64(prevCount)) / float64(prevCount) * 100
		}
		
		prevCount = count
		stats = append(stats, stat)
	}
	
	// If no data for some years in the range, fill with zeroes for continuity
	if len(stats) < years {
		// Generate a complete list of years
		currentYear := time.Now().Year()
		startYear := currentYear - years + 1
		
		filledStats := make([]domain.GrowthStat, 0, years)
		
		// Create a map of existing stats for lookup
		existingStats := make(map[string]domain.GrowthStat)
		for _, stat := range stats {
			existingStats[stat.Period] = stat
		}
		
		// Fill in all years
		for y := 0; y < years; y++ {
			yearStr := fmt.Sprintf("%04d", startYear + y)
			
			if stat, exists := existingStats[yearStr]; exists {
				filledStats = append(filledStats, stat)
			} else {
				// Add empty stat
				filledStats = append(filledStats, domain.GrowthStat{
					Period: yearStr,
					Count:  0,
				})
			}
		}
		
		// Recalculate growth rates with filled data
		for i := 1; i < len(filledStats); i++ {
			prevCount := filledStats[i-1].Count
			if prevCount > 0 {
				filledStats[i].GrowthRate = (float64(filledStats[i].Count) - float64(prevCount)) / float64(prevCount) * 100
			}
		}
		
		stats = filledStats
	}
	
	return stats, nil
}

// GetPublishedExpertStats retrieves the count and percentage of published experts
func (s *SQLiteStore) GetPublishedExpertStats() (int, float64, error) {
	// Get published experts count
	var publishedCount int
	err := s.db.QueryRow("SELECT COUNT(*) FROM experts WHERE is_published = 1").Scan(&publishedCount)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count published experts: %w", err)
	}
	
	// Get total experts count for calculating ratio
	var totalExperts int
	err = s.db.QueryRow("SELECT COUNT(*) FROM experts").Scan(&totalExperts)
	if err != nil {
		return publishedCount, 0, fmt.Errorf("failed to count total experts: %w", err)
	}
	
	// Calculate published ratio, avoid division by zero
	var publishedRatio float64
	if totalExperts > 0 {
		publishedRatio = float64(publishedCount) / float64(totalExperts) * 100
	}
	
	return publishedCount, publishedRatio, nil
}

// GetAreaStatistics returns statistics for general areas and specialized areas
func (s *SQLiteStore) GetAreaStatistics() (map[string][]domain.AreaStat, error) {
	result := make(map[string][]domain.AreaStat)
	
	// Get total experts count for percentage calculations
	var totalExperts int
	err := s.db.QueryRow("SELECT COUNT(*) FROM experts").Scan(&totalExperts)
	if err != nil {
		return nil, fmt.Errorf("failed to count total experts: %w", err)
	}
	
	// Get general area statistics
	generalRows, err := s.db.Query(`
		SELECT ea.name as area_name, COUNT(*) as count
		FROM experts e
		JOIN expert_areas ea ON e.general_area = ea.id
		GROUP BY e.general_area
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query general area statistics: %w", err)
	}
	defer generalRows.Close()
	
	generalStats := []domain.AreaStat{}
	for generalRows.Next() {
		var stat domain.AreaStat
		var count int
		if err := generalRows.Scan(&stat.Name, &count); err != nil {
			return nil, fmt.Errorf("failed to scan general area row: %w", err)
		}
		stat.Count = count
		if totalExperts > 0 {
			stat.Percentage = float64(count) / float64(totalExperts) * 100
		}
		generalStats = append(generalStats, stat)
	}
	result["generalAreas"] = generalStats
	
	// Get specialized area statistics - top 5
	specializedRows, err := s.db.Query(`
		SELECT specialized_area as area_name, COUNT(*) as count
		FROM experts
		WHERE specialized_area != ''
		GROUP BY specialized_area
		ORDER BY count DESC
		LIMIT 5
	`)
	if err != nil {
		return result, fmt.Errorf("failed to query top specialized areas: %w", err)
	}
	defer specializedRows.Close()
	
	topSpecializedStats := []domain.AreaStat{}
	for specializedRows.Next() {
		var stat domain.AreaStat
		var count int
		if err := specializedRows.Scan(&stat.Name, &count); err != nil {
			return result, fmt.Errorf("failed to scan specialized area row: %w", err)
		}
		stat.Count = count
		if totalExperts > 0 {
			stat.Percentage = float64(count) / float64(totalExperts) * 100
		}
		topSpecializedStats = append(topSpecializedStats, stat)
	}
	result["topSpecializedAreas"] = topSpecializedStats
	
	// Get specialized area statistics - bottom 5
	bottomSpecializedRows, err := s.db.Query(`
		SELECT specialized_area as area_name, COUNT(*) as count
		FROM experts
		WHERE specialized_area != ''
		GROUP BY specialized_area
		ORDER BY count ASC
		LIMIT 5
	`)
	if err != nil {
		return result, fmt.Errorf("failed to query bottom specialized areas: %w", err)
	}
	defer bottomSpecializedRows.Close()
	
	bottomSpecializedStats := []domain.AreaStat{}
	for bottomSpecializedRows.Next() {
		var stat domain.AreaStat
		var count int
		if err := bottomSpecializedRows.Scan(&stat.Name, &count); err != nil {
			return result, fmt.Errorf("failed to scan bottom specialized area row: %w", err)
		}
		stat.Count = count
		if totalExperts > 0 {
			stat.Percentage = float64(count) / float64(totalExperts) * 100
		}
		bottomSpecializedStats = append(bottomSpecializedStats, stat)
	}
	result["bottomSpecializedAreas"] = bottomSpecializedStats
	
	return result, nil
}