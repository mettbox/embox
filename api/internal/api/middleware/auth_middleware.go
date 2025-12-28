package middleware

import (
	"embox/internal/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := authService.GetAccessCookie(c)
		var claims services.Claims
		if err == nil {
			claims, err = authService.ValidateAccessCookie(cookie)
		}

		if err == nil {
			c.Set("user", claims.User) // = user Email
			c.Next()
			return
		}

		c.Next()
	}
}
