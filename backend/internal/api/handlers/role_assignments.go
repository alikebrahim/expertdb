package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"expertdb/internal/api/response"
	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"

	"github.com/gorilla/mux"
)

// RoleAssignmentHandler handles role assignment operations
type RoleAssignmentHandler struct {
	storage storage.Storage
}

// NewRoleAssignmentHandler creates a new role assignment handler
func NewRoleAssignmentHandler(storage storage.Storage) *RoleAssignmentHandler {
	return &RoleAssignmentHandler{
		storage: storage,
	}
}

// AssignPlannerApplications assigns a user as planner to multiple applications
// POST /api/users/{id}/planner-assignments
func (h *RoleAssignmentHandler) AssignPlannerApplications(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	vars := mux.Vars(r)
	
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	var req struct {
		ApplicationIDs []int `json:"application_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode request body", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	if len(req.ApplicationIDs) == 0 {
		response.Error(w, domain.ErrBadRequest)
		return
	}

	err = h.storage.AssignUserToPlannerApplications(int(userID), req.ApplicationIDs)
	if err != nil {
		log.Error("Failed to assign planner applications", "error", err, "userID", userID)
		response.Error(w, domain.ErrInternalServer)
		return
	}

	log.Info("User assigned as planner to applications", "userID", userID, "applicationCount", len(req.ApplicationIDs))
	response.Success(w, http.StatusOK, "Planner assignments updated successfully", map[string]interface{}{
		"user_id": userID,
		"assigned_applications": len(req.ApplicationIDs),
	})
}

// AssignManagerApplications assigns a user as manager to multiple applications
// POST /api/users/{id}/manager-assignments
func (h *RoleAssignmentHandler) AssignManagerApplications(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	vars := mux.Vars(r)
	
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	var req struct {
		ApplicationIDs []int `json:"application_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode request body", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	if len(req.ApplicationIDs) == 0 {
		response.Error(w, domain.ErrBadRequest)
		return
	}

	err = h.storage.AssignUserToManagerApplications(int(userID), req.ApplicationIDs)
	if err != nil {
		log.Error("Failed to assign manager applications", "error", err, "userID", userID)
		response.Error(w, domain.ErrInternalServer)
		return
	}

	log.Info("User assigned as manager to applications", "userID", userID, "applicationCount", len(req.ApplicationIDs))
	response.Success(w, http.StatusOK, "Manager assignments updated successfully", map[string]interface{}{
		"user_id": userID,
		"assigned_applications": len(req.ApplicationIDs),
	})
}

// RemovePlannerAssignments removes planner assignments for a user
// DELETE /api/users/{id}/planner-assignments
func (h *RoleAssignmentHandler) RemovePlannerAssignments(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	vars := mux.Vars(r)
	
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	var req struct {
		ApplicationIDs []int `json:"application_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode request body", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	if len(req.ApplicationIDs) == 0 {
		response.Error(w, domain.ErrBadRequest)
		return
	}

	err = h.storage.RemoveUserPlannerAssignments(int(userID), req.ApplicationIDs)
	if err != nil {
		log.Error("Failed to remove planner assignments", "error", err, "userID", userID)
		response.Error(w, domain.ErrInternalServer)
		return
	}

	log.Info("Planner assignments removed for user", "userID", userID, "applicationCount", len(req.ApplicationIDs))
	response.Success(w, http.StatusOK, "Planner assignments removed successfully", map[string]interface{}{
		"user_id": userID,
		"removed_applications": len(req.ApplicationIDs),
	})
}

// RemoveManagerAssignments removes manager assignments for a user
// DELETE /api/users/{id}/manager-assignments
func (h *RoleAssignmentHandler) RemoveManagerAssignments(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	vars := mux.Vars(r)
	
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	var req struct {
		ApplicationIDs []int `json:"application_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode request body", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	if len(req.ApplicationIDs) == 0 {
		response.Error(w, domain.ErrBadRequest)
		return
	}

	err = h.storage.RemoveUserManagerAssignments(int(userID), req.ApplicationIDs)
	if err != nil {
		log.Error("Failed to remove manager assignments", "error", err, "userID", userID)
		response.Error(w, domain.ErrInternalServer)
		return
	}

	log.Info("Manager assignments removed for user", "userID", userID, "applicationCount", len(req.ApplicationIDs))
	response.Success(w, http.StatusOK, "Manager assignments removed successfully", map[string]interface{}{
		"user_id": userID,
		"removed_applications": len(req.ApplicationIDs),
	})
}

// GetUserAssignments returns all planner and manager assignments for a user
// GET /api/users/{id}/assignments
func (h *RoleAssignmentHandler) GetUserAssignments(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	vars := mux.Vars(r)
	
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Error("Invalid user ID", "error", err)
		response.Error(w, domain.ErrBadRequest)
		return
	}

	plannerApps, err := h.storage.GetUserPlannerApplications(int(userID))
	if err != nil {
		log.Error("Failed to get planner applications", "error", err, "userID", userID)
		response.Error(w, domain.ErrInternalServer)
		return
	}

	managerApps, err := h.storage.GetUserManagerApplications(int(userID))
	if err != nil {
		log.Error("Failed to get manager applications", "error", err, "userID", userID)
		response.Error(w, domain.ErrInternalServer)
		return
	}

	response.Success(w, http.StatusOK, "User assignments retrieved successfully", map[string]interface{}{
		"user_id": userID,
		"planner_applications": plannerApps,
		"manager_applications": managerApps,
	})
}