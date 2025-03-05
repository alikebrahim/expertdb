package main

import (
	"context"
	"strconv"
	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a custom type for context keys to avoid conflicts
type contextKey string

// Define context keys
const (
	userContextKey contextKey = "user"
)

// UserClaims provides strong typing for JWT claims
type UserClaims struct {
	UserID    int64
	Name      string
	Email     string
	Role      string
	ExpiresAt int64
}

// setUserContext adds user claims to the request context
func setUserContext(ctx context.Context, claims jwt.MapClaims) context.Context {
	userClaims := mapToUserClaims(claims)
	return context.WithValue(ctx, userContextKey, userClaims)
}

// mapToUserClaims converts jwt.MapClaims to UserClaims
func mapToUserClaims(claims jwt.MapClaims) *UserClaims {
	userClaims := &UserClaims{}
	
	// Extract UserID
	if sub, ok := claims["sub"].(string); ok {
		id, err := strconv.ParseInt(sub, 10, 64)
		if err == nil {
			userClaims.UserID = id
		}
	}
	
	// Extract other fields
	if name, ok := claims["name"].(string); ok {
		userClaims.Name = name
	}
	
	if email, ok := claims["email"].(string); ok {
		userClaims.Email = email
	}
	
	if role, ok := claims["role"].(string); ok {
		userClaims.Role = role
	}
	
	if exp, ok := claims["exp"].(float64); ok {
		userClaims.ExpiresAt = int64(exp)
	}
	
	return userClaims
}

// getUserFromContext extracts user claims from the request context
func getUserFromContext(ctx context.Context) (*UserClaims, bool) {
	userClaims, ok := ctx.Value(userContextKey).(*UserClaims)
	return userClaims, ok
}

// getUserIDFromContext extracts the user ID from the request context
func getUserIDFromContext(ctx context.Context) (int64, bool) {
	claims, ok := getUserFromContext(ctx)
	if !ok {
		return 0, false
	}
	
	return claims.UserID, claims.UserID > 0
}

// getUserRoleFromContext extracts the user role from the request context
func getUserRoleFromContext(ctx context.Context) (string, bool) {
	claims, ok := getUserFromContext(ctx)
	if !ok {
		return "", false
	}
	
	return claims.Role, claims.Role != ""
}

// isAdmin checks if the user in the context is an admin
func isAdmin(ctx context.Context) bool {
	role, ok := getUserRoleFromContext(ctx)
	return ok && role == RoleAdmin
}