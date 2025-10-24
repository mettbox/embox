package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var clients = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, exists := clients[ip]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		clients[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware(maxRequests int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip, rate.Every(duration/time.Duration(maxRequests)), maxRequests)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}
		c.Next()
	}
}
