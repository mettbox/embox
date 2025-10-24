package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAuthMiddleware checks if the user is authenticated.
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
