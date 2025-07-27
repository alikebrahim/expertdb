package utils

import (
	"net/http"
	"strconv"
)

// Pagination holds pagination parameters
type Pagination struct {
	Limit  int
	Offset int
}

// PaginationResponse holds pagination metadata for responses
type PaginationResponse struct {
	TotalCount   int  `json:"totalCount"`
	TotalPages   int  `json:"totalPages"`
	CurrentPage  int  `json:"currentPage"`
	PageSize     int  `json:"pageSize"`
	HasNextPage  bool `json:"hasNextPage"`
	HasPrevPage  bool `json:"hasPrevPage"`
	HasMore      bool `json:"hasMore"`
}

// ParsePaginationParams extracts pagination parameters from request with defaults
func ParsePaginationParams(r *http.Request, defaultLimit int) Pagination {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}
	
	// Cap the limit to prevent abuse
	if limit > 1000 {
		limit = 1000
	}
	
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	
	return Pagination{Limit: limit, Offset: offset}
}

// BuildPaginationResponse creates pagination metadata for API responses
func BuildPaginationResponse(data interface{}, totalCount int, pagination Pagination) map[string]interface{} {
	pageSize := pagination.Limit
	currentPage := (pagination.Offset / pageSize) + 1
	totalPages := (totalCount + pageSize - 1) / pageSize
	
	if totalPages == 0 {
		totalPages = 1
	}
	
	hasNext := currentPage < totalPages
	hasPrev := currentPage > 1
	hasMore := pagination.Offset+pageSize < totalCount
	
	return map[string]interface{}{
		"data": data,
		"pagination": PaginationResponse{
			TotalCount:   totalCount,
			TotalPages:   totalPages,
			CurrentPage:  currentPage,
			PageSize:     pageSize,
			HasNextPage:  hasNext,
			HasPrevPage:  hasPrev,
			HasMore:      hasMore,
		},
	}
}

// BuildSimplePaginationResponse creates a simpler pagination response without total counts
func BuildSimplePaginationResponse(data interface{}, hasMore bool, pagination Pagination) map[string]interface{} {
	pageSize := pagination.Limit
	currentPage := (pagination.Offset / pageSize) + 1
	hasPrev := pagination.Offset > 0
	
	return map[string]interface{}{
		"data": data,
		"pagination": map[string]interface{}{
			"currentPage": currentPage,
			"pageSize":    pageSize,
			"hasMore":     hasMore,
			"hasPrevPage": hasPrev,
			"hasNextPage": hasMore,
			"offset":      pagination.Offset,
			"limit":       pagination.Limit,
		},
	}
}