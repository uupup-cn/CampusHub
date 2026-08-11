package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID     = "discourse_user_id"
	ContextKeyUsername   = "username"
	ContextKeyTrustLevel = "trust_level"
	ContextKeyScopes     = "scopes"
	ContextKeyClientID   = "client_id"
)

// BearerAuth validates a Bearer token against the database.
// Use for routes that REQUIRE OAuth authentication.
func BearerAuth(appRepo *repository.AppRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			response.Error(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		at, err := appRepo.GetAccessToken(token)
		if err != nil {
			response.Error(c, errcode.ErrTokenInvalid)
			c.Abort()
			return
		}
		if at.AccessExpiresAt.Before(time.Now()) {
			response.Error(c, errcode.ErrTokenExpired)
			c.Abort()
			return
		}
		c.Set(ContextKeyUserID, at.DiscourseUserID)
		c.Set(ContextKeyClientID, at.ClientID)
		c.Set(ContextKeyScopes, scopesFromJSON(at.Scopes))
		c.Next()
	}
}

// OptionalAuth tries Bearer token first, then falls back to X-User-ID header.
// If Authorization header is present but token is invalid, returns 401 (does NOT fall back).
func OptionalAuth(appRepo *repository.AppRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			at, err := appRepo.GetAccessToken(token)
			if err != nil {
				response.Error(c, errcode.ErrTokenInvalid)
				c.Abort()
				return
			}
			if at.AccessExpiresAt.Before(time.Now()) {
				response.Error(c, errcode.ErrTokenExpired)
				c.Abort()
				return
			}
			c.Set(ContextKeyUserID, at.DiscourseUserID)
			c.Set(ContextKeyClientID, at.ClientID)
			c.Set(ContextKeyScopes, scopesFromJSON(at.Scopes))
			c.Next()
			return
		}
		// No Authorization header - fall back to X-User-ID (internal/plugin)
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr != "" {
			c.Set(ContextKeyUserID, parseUserID(userIDStr))
		}
		c.Next()
	}
}

// RequireScope checks that the OAuth token has the required scope.
// Internal calls (X-User-ID) bypass scope checks.
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get(ContextKeyClientID); !exists {
			c.Next()
			return
		}
		scopes, _ := c.Get(ContextKeyScopes)
		scopesStr := scopes.(string)
		if scopesStr == "" || !strings.Contains(scopesStr, scope) {
			response.Error(c, errcode.ErrScopeInsufficient)
			c.Abort()
			return
		}
		c.Next()
	}
}

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			response.Error(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		c.Set("api_key", key)
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Admin-Key")
		if key == "" {
			response.Error(c, errcode.ErrPermissionDenied)
			c.Abort()
			return
		}
		c.Set("admin_key", key)
		c.Next()
	}
}

func parseUserID(s string) int64 {
	var id int64
	_, err := fmt.Sscanf(s, "%d", &id)
	_ = err
	return id
}

func scopesFromJSON(scopeJSON string) string {
	return scopeJSON
}
