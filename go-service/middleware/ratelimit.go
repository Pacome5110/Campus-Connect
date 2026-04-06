package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string]int
}

var limiter = &rateLimiter{requests: make(map[string]int)}

// Simple 10s window reset
func init() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			limiter.mu.Lock()
			limiter.requests = make(map[string]int)
			limiter.mu.Unlock()
		}
	}()
}

func RateLimiter(maxPerWindow int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter.mu.Lock()
		limiter.requests[ip]++
		count := limiter.requests[ip]
		limiter.mu.Unlock()

		if count > maxPerWindow {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests. Rate limit exceeded."})
			c.Abort()
			return
		}

		c.Next()
	}
}
