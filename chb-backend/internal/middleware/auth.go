package middleware

import (
	"strings"

	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID    = "discourse_user_id"
	ContextKeyUsername  = "username"
	ContextKeyTrustLevel = "trust_level"
	ContextKeyScopes    = "scopes"
	ContextKeyClientID  = "client_id"
)

func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.Next()
			return
		}
		c.Set(ContextKeyUserID, userID)
		c.Next()
	}
}

func BearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			response.Error(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		c.Set("access_token", token)
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

func ScopeCheck(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopes, exists := c.Get(ContextKeyScopes)
		if !exists {
			response.Error(c, errcode.ErrScopeInsufficient)
			c.Abort()
			return
		}
		scopesStr := scopes.(string)
		for _, required := range requiredScopes {
			if !strings.Contains(scopesStr, required) {
				response.Error(c, errcode.ErrScopeInsufficient)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
