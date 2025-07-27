package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	
	"expertdb/internal/domain"
	"expertdb/internal/logger"
)

// User context key type to avoid key collisions in context
type contextKey string

const (
	// UserClaimsContextKey is the key used to store user claims in context
	UserClaimsContextKey contextKey = "userClaims"
)

// UserClaims provides strong typing for JWT claims
type UserClaims struct {
	UserID    int64
	Name      string
	Email     string
	Role      string
	ExpiresAt int64
}

// GetUserClaimsFromContext extracts user claims from the request context
func GetUserClaimsFromContext(ctx context.Context) (map[string]interface{}, bool) {
	claims, ok := ctx.Value(UserClaimsContextKey).(map[string]interface{})
	return claims, ok
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	claims, ok := GetUserClaimsFromContext(ctx)
	if !ok {
		return 0, false
	}
	
	// Extract UserID
	if sub, ok := claims["sub"].(string); ok {
		id, err := strconv.ParseInt(sub, 10, 64)
		if err == nil {
			return id, true
		}
	}
	
	return 0, false
}

// GetUserRoleFromContext extracts the user role from the request context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	claims, ok := GetUserClaimsFromContext(ctx)
	if !ok {
		return "", false
	}
	
	role, ok := claims["role"].(string)
	return role, ok
}

// IsAdmin checks if the user in the context is an admin or higher
func IsAdmin(ctx context.Context) bool {
	role, ok := GetUserRoleFromContext(ctx)
	return ok && (role == RoleAdmin || role == RoleSuperUser)
}

// IsSuperUser checks if the user in the context is a super user
func IsSuperUser(ctx context.Context) bool {
	role, ok := GetUserRoleFromContext(ctx)
	return ok && role == RoleSuperUser
}

// GetUserIDFromRequest extracts user ID from request context with proper error handling
func GetUserIDFromRequest(r *http.Request) (int64, error) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		return 0, domain.ErrUnauthorized
	}
	return userID, nil
}

// GetUserRoleFromRequest extracts user role from request context with proper error handling
func GetUserRoleFromRequest(r *http.Request) (string, error) {
	role, ok := GetUserRoleFromContext(r.Context())
	if !ok {
		return "", domain.ErrUnauthorized
	}
	return role, nil
}


// SetUserClaimsInContext adds user claims to the request context
func SetUserClaimsInContext(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, UserClaimsContextKey, claims)
}

// ExtractTokenFromHeader extracts the JWT token from the Authorization header
func ExtractTokenFromHeader(r *http.Request) (string, error) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	
	// Split the header into parts and validate format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	
	// Return the token part
	return parts[1], nil
}

// HandlerFunc is the type for HTTP handlers that can return errors
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// RequireAuth is middleware that verifies a user is authenticated before allowing access
func RequireAuth(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log := logger.Get()
		
		// Extract the token from the Authorization header
		token, err := ExtractTokenFromHeader(r)
		if err != nil {
			log.Debug("Authentication failed: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			log.Debug("JWT verification failed: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Add user claims to request context for downstream handlers
		ctx := SetUserClaimsInContext(r.Context(), claims)
		
		// Pass control to the next handler with the updated context
		return next(w, r.WithContext(ctx))
	}
}

// RequireAdmin is middleware that ensures only admin users can access protected endpoints
func RequireAdmin(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log := logger.Get()
		
		// Extract the token from the Authorization header
		token, err := ExtractTokenFromHeader(r)
		if err != nil {
			log.Debug("Admin check failed - missing token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			log.Debug("Admin check failed - invalid token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Check if the user has admin role or higher
		role, ok := claims["role"].(string)
		if !ok || !HasRole(role, RoleAdmin) {
			// User is authenticated but doesn't have admin privileges
			log.Info("Forbidden access attempt by user with insufficient privileges (ID: %v, Role: %v) to %s", 
				claims["sub"], role, r.URL.Path)
			return domain.ErrForbidden
		}
		
		// Add user claims to request context for downstream handlers
		ctx := SetUserClaimsInContext(r.Context(), claims)
		
		// Pass control to the next handler with the updated context
		return next(w, r.WithContext(ctx))
	}
}

// RequireSuperUser is middleware that ensures only super users can access protected endpoints
func RequireSuperUser(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log := logger.Get()
		
		// Extract the token from the Authorization header
		token, err := ExtractTokenFromHeader(r)
		if err != nil {
			log.Debug("Super user check failed - missing token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			log.Debug("Super user check failed - invalid token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Check if the user has super_user role
		role, ok := claims["role"].(string)
		if !ok || role != RoleSuperUser {
			// User is authenticated but not a super user
			log.Info("Forbidden access attempt by non-super user (ID: %v, Role: %v) to %s", 
				claims["sub"], role, r.URL.Path)
			return domain.ErrForbidden
		}
		
		// Add user claims to request context for downstream handlers
		ctx := SetUserClaimsInContext(r.Context(), claims)
		
		// Pass control to the next handler with the updated context
		return next(w, r.WithContext(ctx))
	}
}

// RequireRole is a more flexible middleware that ensures users have at least the specified role
func RequireRole(minRole string, next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log := logger.Get()
		
		// Extract the token from the Authorization header
		token, err := ExtractTokenFromHeader(r)
		if err != nil {
			log.Debug("Role check failed - missing token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Verify the token and extract claims
		_, claims, err := VerifyJWT(token)
		if err != nil {
			log.Debug("Role check failed - invalid token: %v", err)
			return domain.ErrUnauthorized
		}
		
		// Check if the user has the required role or higher
		role, ok := claims["role"].(string)
		if !ok || !HasRole(role, minRole) {
			// User is authenticated but doesn't have sufficient privileges
			log.Info("Forbidden access attempt by user with insufficient privileges (ID: %v, Role: %v, Required: %v) to %s", 
				claims["sub"], role, minRole, r.URL.Path)
			return domain.ErrForbidden
		}
		
		// Add user claims to request context for downstream handlers
		ctx := SetUserClaimsInContext(r.Context(), claims)
		
		// Pass control to the next handler with the updated context
		return next(w, r.WithContext(ctx))
	}
}