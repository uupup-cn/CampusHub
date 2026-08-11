package middleware

import (
	"github.com/campushub/chb-backend/internal/config"
	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}

		allowed := false
		for _, o := range cfg.AllowedOrigins {
			if o == origin || o == "*" {
				allowed = true
				break
			}
		}
		if allowed || len(cfg.AllowedOrigins) == 0 {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", joinStrings(cfg.AllowedMethods))
		c.Header("Access-Control-Allow-Headers", joinStrings(cfg.AllowedHeaders))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func joinStrings(s []string) string {
	if len(s) == 0 {
		return "*"
	}
	result := ""
	for i, v := range s {
		if i > 0 {
			result += ", "
		}
		result += v
	}
	return result
}
