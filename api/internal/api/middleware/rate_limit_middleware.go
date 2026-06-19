package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const clientTTL = 5 * time.Minute

type rateLimitClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients   = make(map[string]*rateLimitClient)
	mu        sync.Mutex
	cleanOnce sync.Once
)

func startCleanup() {
	cleanOnce.Do(func() {
		go func() {
			for {
				time.Sleep(clientTTL)
				mu.Lock()
				for ip, c := range clients {
					if time.Since(c.lastSeen) > clientTTL {
						delete(clients, ip)
					}
				}
				mu.Unlock()
			}
		}()
	})
}

func getLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	c, exists := clients[ip]
	if !exists {
		c = &rateLimitClient{limiter: rate.NewLimiter(r, b)}
		clients[ip] = c
	}
	c.lastSeen = time.Now()
	return c.limiter
}

func ResetRateLimiter() {
	mu.Lock()
	clients = make(map[string]*rateLimitClient)
	mu.Unlock()
}

func RateLimitMiddleware(maxRequests int, duration time.Duration) gin.HandlerFunc {
	startCleanup()
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
