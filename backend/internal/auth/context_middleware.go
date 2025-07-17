package auth

import (
	"net/http"
	"strconv"

	"expertdb/internal/domain"
	"expertdb/internal/logger"
	"expertdb/internal/storage"

	"github.com/gorilla/mux"
)

// RequirePlannerForApplication is middleware that ensures the user has planner privileges for the specific application
func RequirePlannerForApplication(store storage.Storage, next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log := logger.Get()
		
		// First, verify authentication
		token, err := ExtractTokenFromHeader(r)
		if err != nil {
			log.Debug("Planner check failed - missing token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			log.Debug("Planner check failed - invalid token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Get user ID from claims
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			log.Debug("Planner check failed - invalid user ID in token")
			return domain.ErrUnauthorized
		}
		
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Debug("Planner check failed - invalid user ID format: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Check if user is admin or super_user (they can bypass application-specific checks)
		role, ok := claims["role"].(string)
		if ok && (role == RoleAdmin || role == RoleSuperUser) {
			// Add user claims to context
			ctx := SetUserClaimsInContext(r.Context(), claims)
			return next(w, r.WithContext(ctx))
		}
		
		// Extract application ID from URL
		vars := mux.Vars(r)
		appIDStr, exists := vars["app_id"]
		if !exists {
			log.Debug("Planner check failed - no application ID in URL")
			return domain.ErrBadRequest
		}
		
		appID, err := strconv.Atoi(appIDStr)
		if err != nil {
			log.Debug("Planner check failed - invalid application ID: %v", err)
			return domain.ErrBadRequest
		}
		
		// Check if user has planner privileges for this application
		hasAccess, err := store.IsUserPlannerForApplication(int(userID), appID)
		if err != nil {
			log.Error("Failed to check planner permissions", "error", err, "userID", userID, "appID", appID)
			return domain.ErrInternalServer
		}
		
		if !hasAccess {
			log.Info("Forbidden access attempt by user without planner privileges (ID: %v) to application %v", userID, appID)
			return domain.ErrForbidden
		}
		
		// Add user claims to context
		ctx := SetUserClaimsInContext(r.Context(), claims)
		return next(w, r.WithContext(ctx))
	}
}

// RequireManagerForApplication is middleware that ensures the user has manager privileges for the specific application
func RequireManagerForApplication(store storage.Storage, next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log := logger.Get()
		
		// First, verify authentication
		token, err := ExtractTokenFromHeader(r)
		if err != nil {
			log.Debug("Manager check failed - missing token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			log.Debug("Manager check failed - invalid token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Get user ID from claims
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			log.Debug("Manager check failed - invalid user ID in token")
			return domain.ErrUnauthorized
		}
		
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Debug("Manager check failed - invalid user ID format: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Check if user is admin or super_user (they can bypass application-specific checks)
		role, ok := claims["role"].(string)
		if ok && (role == RoleAdmin || role == RoleSuperUser) {
			// Add user claims to context
			ctx := SetUserClaimsInContext(r.Context(), claims)
			return next(w, r.WithContext(ctx))
		}
		
		// Extract application ID from URL
		vars := mux.Vars(r)
		appIDStr, exists := vars["app_id"]
		if !exists {
			log.Debug("Manager check failed - no application ID in URL")
			return domain.ErrBadRequest
		}
		
		appID, err := strconv.Atoi(appIDStr)
		if err != nil {
			log.Debug("Manager check failed - invalid application ID: %v", err)
			return domain.ErrBadRequest
		}
		
		// Check if user has manager privileges for this application
		hasAccess, err := store.IsUserManagerForApplication(int(userID), appID)
		if err != nil {
			log.Error("Failed to check manager permissions", "error", err, "userID", userID, "appID", appID)
			return domain.ErrInternalServer
		}
		
		if !hasAccess {
			log.Info("Forbidden access attempt by user without manager privileges (ID: %v) to application %v", userID, appID)
			return domain.ErrForbidden
		}
		
		// Add user claims to context
		ctx := SetUserClaimsInContext(r.Context(), claims)
		return next(w, r.WithContext(ctx))
	}
}

// RequireApplicationAccess is middleware that checks if a user has any access (planner or manager) to an application
func RequireApplicationAccess(store storage.Storage, next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log := logger.Get()
		
		// First, verify authentication
		token, err := ExtractTokenFromHeader(r)
		if err != nil {
			log.Debug("Application access check failed - missing token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			log.Debug("Application access check failed - invalid token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Get user ID from claims
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			log.Debug("Application access check failed - invalid user ID in token")
			return domain.ErrUnauthorized
		}
		
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Debug("Application access check failed - invalid user ID format: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Check if user is admin or super_user (they can bypass application-specific checks)
		role, ok := claims["role"].(string)
		if ok && (role == RoleAdmin || role == RoleSuperUser) {
			// Add user claims to context
			ctx := SetUserClaimsInContext(r.Context(), claims)
			return next(w, r.WithContext(ctx))
		}
		
		// Extract application ID from URL
		vars := mux.Vars(r)
		appIDStr, exists := vars["app_id"]
		if !exists {
			log.Debug("Application access check failed - no application ID in URL")
			return domain.ErrBadRequest
		}
		
		appID, err := strconv.Atoi(appIDStr)
		if err != nil {
			log.Debug("Application access check failed - invalid application ID: %v", err)
			return domain.ErrBadRequest
		}
		
		// Check if user has planner or manager privileges for this application
		isPlannerForApp, err := store.IsUserPlannerForApplication(int(userID), appID)
		if err != nil {
			log.Error("Failed to check planner permissions", "error", err, "userID", userID, "appID", appID)
			return domain.ErrInternalServer
		}
		
		isManagerForApp, err := store.IsUserManagerForApplication(int(userID), appID)
		if err != nil {
			log.Error("Failed to check manager permissions", "error", err, "userID", userID, "appID", appID)
			return domain.ErrInternalServer
		}
		
		if !isPlannerForApp && !isManagerForApp {
			log.Info("Forbidden access attempt by user without application privileges (ID: %v) to application %v", userID, appID)
			return domain.ErrForbidden
		}
		
		// Add user claims to context
		ctx := SetUserClaimsInContext(r.Context(), claims)
		return next(w, r.WithContext(ctx))
	}
}