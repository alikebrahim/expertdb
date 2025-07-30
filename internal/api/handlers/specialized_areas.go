package handlers

import (
	"fmt"
	"net/http"
	"strings"
	
	"expertdb/internal/api/utils"
	"expertdb/internal/domain"
	"expertdb/internal/storage"
)

type SpecializedAreasHandler struct {
	store storage.Storage
}

func NewSpecializedAreasHandler(store storage.Storage) *SpecializedAreasHandler {
	return &SpecializedAreasHandler{store: store}
}

// HandleListSpecializedAreas returns all specialized areas with optional search functionality
func (h *SpecializedAreasHandler) HandleListSpecializedAreas(w http.ResponseWriter, r *http.Request) error {
	areas, err := h.store.ListSpecializedAreas()
	if err != nil {
		return fmt.Errorf("failed to retrieve specialized areas: %w", err)
	}
	
	// Optional: implement search filtering
	search := r.URL.Query().Get("search")
	if search != "" {
		filtered := []*domain.SpecializedArea{}
		searchLower := strings.ToLower(search)
		for _, area := range areas {
			if strings.Contains(strings.ToLower(area.Name), searchLower) {
				filtered = append(filtered, area)
			}
		}
		areas = filtered
	}
	
	return utils.RespondWithSuccess(w, "", areas)
}