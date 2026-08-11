package middleware

import (
"sync"

	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	clients = make(map[string]*rate.Limiter)
	mu      sync.Mutex
)

func RateLimit(enabled bool, rps float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled {
			c.Next()
			return
		}

		ip := c.ClientIP()
		mu.Lock()
		limiter, ok := clients[ip]
		if !ok {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			clients[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.Header("Retry-After", "1")
			response.Error(c, errcode.ErrCooldown)
			c.Abort()
			return
		}

		c.Next()
	}
}
