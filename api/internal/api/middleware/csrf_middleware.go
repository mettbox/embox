package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware is a middleware that checks for CSRF tokens in requests.
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for GET, HEAD, and OPTIONS methods
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		cookie, err := c.Cookie("XSRF-TOKEN")
		if err != nil || c.GetHeader("X-XSRF-TOKEN") != cookie {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
