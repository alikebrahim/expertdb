package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AIService provides a client for communicating with the external AI service
type AIService struct {
	baseURL    string
	httpClient *http.Client
	storage    Storage // Added for accessing expert data when suggesting panels
}

// NewAIService creates a new AIService instance
func NewAIService(baseURL string, storage Storage) *AIService {
	return &AIService{
		baseURL: baseURL,
		storage: storage,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateProfile requests an AI-generated profile for an expert
// Implementation matches Storage.GenerateProfile(expertID int64) (*AIAnalysisResult, error)
func (s *AIService) GenerateProfile(expertID int64) (*AIAnalysisResult, error) {
	// First, retrieve the expert
	expert, err := s.storage.GetExpert(expertID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve expert data: %w", err)
	}
	
	// Prepare expert data to send to AI service
	data := map[string]interface{}{
		"name":            expert.Name,
		"designation":     expert.Designation,
		"institution":     expert.Institution,
		"generalArea":     expert.GeneralArea,
		"specializedArea": expert.SpecializedArea,
	}

	result, err := s.callAIService("generate-profile", data)
	if err != nil {
		return nil, err
	}

	// Create an analysis result record
	analysisResult := &AIAnalysisResult{
		ExpertID:        expertID,
		AnalysisType:    "profile_generation",
		AnalysisResult:  result,
		ResultData:      result,
		ConfidenceScore: 0.0, // AI service should provide this
		ModelUsed:       "placeholder",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Store the result in the database
	if err := s.storage.StoreAIAnalysisResult(analysisResult); err != nil {
		return nil, fmt.Errorf("failed to store AI analysis result: %w", err)
	}

	return analysisResult, nil
}

// SuggestISCED requests AI to suggest appropriate ISCED classifications
// Implementation matches Storage.SuggestISCED(expertID int64, input string) (*AIAnalysisResult, error)
func (s *AIService) SuggestISCED(expertID int64, input string) (*AIAnalysisResult, error) {
	// First, retrieve the expert if input is not provided
	var generalArea, specializedArea string
	
	if input == "" {
		expert, err := s.storage.GetExpert(expertID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve expert data: %w", err)
		}
		generalArea = expert.GeneralArea
		specializedArea = expert.SpecializedArea
	} else {
		// Parse input as needed or use it directly
		generalArea = input
		specializedArea = input
	}
	
	data := map[string]interface{}{
		"generalArea":     generalArea,
		"specializedArea": specializedArea,
		"expertID":        expertID,
	}

	result, err := s.callAIService("suggest-isced", data)
	if err != nil {
		return nil, err
	}

	analysisResult := &AIAnalysisResult{
		ExpertID:        expertID,
		AnalysisType:    "isced_suggestion",
		AnalysisResult:  result,
		ResultData:      result,
		ConfidenceScore: 0.85, // AI service should provide this
		ModelUsed:       "placeholder",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Store the result in the database
	if err := s.storage.StoreAIAnalysisResult(analysisResult); err != nil {
		return nil, fmt.Errorf("failed to store AI analysis result: %w", err)
	}

	return analysisResult, nil
}

// ExtractSkills requests AI to extract skills from a document or text
// Implementation matches Storage.ExtractSkills(expertID int64, input string) (*AIAnalysisResult, error)
func (s *AIService) ExtractSkills(expertID int64, input string) (*AIAnalysisResult, error) {
	// Use input as document text or extract text from document if needed
	documentText := input
	
	// If input is empty or looks like a file path, try to extract text from document
	if input == "" || (len(input) > 5 && input[0] == '/') {
		// Check if the expert has CV
		expert, err := s.storage.GetExpert(expertID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve expert data: %w", err)
		}
		
		if input == "" && expert.CVPath != "" {
			// Extract from CV
			documentText = "Sample CV text extracted from " + expert.CVPath
		}
	}
	
	data := map[string]interface{}{
		"text": documentText,
	}

	result, err := s.callAIService("extract-skills", data)
	if err != nil {
		return nil, err
	}

	analysisResult := &AIAnalysisResult{
		ExpertID:        expertID,
		AnalysisType:    "skills_extraction",
		AnalysisResult:  result,
		ResultData:      result,
		ConfidenceScore: 0.9, // AI service should provide this
		ModelUsed:       "placeholder",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Store the result in the database
	if err := s.storage.StoreAIAnalysisResult(analysisResult); err != nil {
		return nil, fmt.Errorf("failed to store AI analysis result: %w", err)
	}

	return analysisResult, nil
}

// SuggestExpertPanel suggests a panel of experts for a given project
// Implementation matches Storage.SuggestExpertPanel(request string, count int) ([]Expert, error)
func (s *AIService) SuggestExpertPanel(request string, count int) ([]Expert, error) {
	// Parse the request string to extract any relevant information
	// For example, parse for project name and ISCED field if included
	projectName := request
	var iscedFieldID int64 = 0
	
	// Get all relevant experts first
	filters := map[string]interface{}{
		"is_available": true,
		"is_published": true,
	}
	
	// Add ISCED field filter if extracted from request
	if iscedFieldID > 0 {
		filters["isced_field_id"] = iscedFieldID
	}
	
	// Get experts matching criteria
	experts, err := s.storage.ListExperts(filters, count*2, 0) // Get more than needed for AI to select from
	if err != nil {
		return nil, fmt.Errorf("failed to fetch experts: %w", err)
	}
	
	if len(experts) == 0 {
		return []Expert{}, nil // Return empty list if no experts found
	}
	
	// Create data for AI processing
	data := map[string]interface{}{
		"project_name":     projectName,
		"request":          request,
		"num_experts":      count,
		"available_experts": experts,
	}
	
	// Call AI service to rank the experts
	result, err := s.callAIService("suggest-panel", data)
	if err != nil {
		return nil, fmt.Errorf("AI service error: %w", err)
	}
	
	// In a real implementation, we would parse the AI response to get ranked experts
	// For now, just return the top N experts from our list
	
	// Create analysis result for tracking
	analysisResult := &AIAnalysisResult{
		AnalysisType:    "panel_suggestion",
		AnalysisResult:  result,
		ResultData:      result,
		ConfidenceScore: 0.85, // Higher confidence due to data-driven approach
		ModelUsed:       "placeholder",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Store the analysis for future reference
	if err := s.storage.StoreAIAnalysisResult(analysisResult); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: Failed to store AI analysis result: %v\n", err)
	}
	
	// Convert *Expert to Expert for the top N results
	selectedExperts := make([]Expert, 0, count)
	for i, expert := range experts {
		if i >= count {
			break
		}
		selectedExperts = append(selectedExperts, *expert)
	}
	
	return selectedExperts, nil
}

// callAIService is a helper method to make API calls to the AI service
func (s *AIService) callAIService(endpoint string, data map[string]interface{}) (string, error) {
	// In a real implementation, this would make an HTTP request to the AI service
	// For now, we'll simulate AI responses

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	// Construct URL
	url := fmt.Sprintf("%s/%s", s.baseURL, endpoint)

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request - in a real implementation
	// For now, simulate a response
	simulatedResponse := "This is a simulated AI response. In production, this would be replaced with actual calls to the AI service."

	// In production, you would:
	// resp, err := s.httpClient.Do(req)
	// if err != nil {
	//     return "", fmt.Errorf("failed to send request to AI service: %w", err)
	// }
	// defer resp.Body.Close()
	//
	// // Read response body
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	//     return "", fmt.Errorf("failed to read response: %w", err)
	// }
	//
	// return string(body), nil

	return simulatedResponse, nil
}