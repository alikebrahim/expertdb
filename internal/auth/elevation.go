package auth

import (
	"context"
	"net/http"

	"expertdb/internal/domain"
	"expertdb/internal/storage"
)

// HasPlannerAccess checks if user has planner access for a specific application
// This combines base role check with contextual elevation
func HasPlannerAccess(ctx context.Context, storage storage.Storage, applicationID int) (bool, error) {
	// Get user from context
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return false, domain.ErrUnauthorized
	}
	
	role, ok := GetUserRoleFromContext(ctx)
	if !ok {
		return false, domain.ErrUnauthorized
	}
	
	// Admin and super_user have planner access to everything
	if role == RoleAdmin || role == RoleSuperUser {
		return true, nil
	}
	
	// Check contextual planner assignment for regular users
	return storage.IsUserPlannerForApplication(int(userID), applicationID)
}

// HasManagerAccess checks if user has manager access for a specific application
func HasManagerAccess(ctx context.Context, storage storage.Storage, applicationID int) (bool, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return false, domain.ErrUnauthorized
	}
	
	role, ok := GetUserRoleFromContext(ctx)
	if !ok {
		return false, domain.ErrUnauthorized
	}
	
	// Admin and super_user have manager access to everything
	if role == RoleAdmin || role == RoleSuperUser {
		return true, nil
	}
	
	// Check contextual manager assignment for regular users
	return storage.IsUserManagerForApplication(int(userID), applicationID)
}

// RequirePlannerAccess middleware for application-specific planner access
func RequirePlannerAccess(storage storage.Storage, getApplicationID func(*http.Request) (int, error)) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			applicationID, err := getApplicationID(r)
			if err != nil {
				return domain.ErrBadRequest
			}
			
			hasAccess, err := HasPlannerAccess(r.Context(), storage, applicationID)
			if err != nil {
				return err
			}
			
			if !hasAccess {
				return domain.ErrForbidden
			}
			
			return next(w, r)
		}
	}
}

// RequireManagerAccess middleware for application-specific manager access
func RequireManagerAccess(storage storage.Storage, getApplicationID func(*http.Request) (int, error)) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			applicationID, err := getApplicationID(r)
			if err != nil {
				return domain.ErrBadRequest
			}
			
			hasAccess, err := HasManagerAccess(r.Context(), storage, applicationID)
			if err != nil {
				return err
			}
			
			if !hasAccess {
				return domain.ErrForbidden
			}
			
			return next(w, r)
		}
	}
}