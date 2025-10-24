package middleware

import (
	"embox/internal/services"
	"embox/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := authService.GetAccessCookie(c)
		var claims jwt.Claims
		if err == nil {
			claims, err = authService.ValidateAccessCookie(cookie)
		}

		if err == nil {
			c.Set("user", claims.User) // = user Email
			c.Next()
			return
		}

		// refreshCookie, err := authService.GetRefreshCookie(c)
		// if err == nil {
		// 	refreshClaims, err := authService.ValidateRefreshCookie(refreshCookie)
		// 	if err == nil {
		// 		_, _ = authService.SetCookie(c, refreshClaims.User)
		// 		c.Set("user", refreshClaims.User)
		// 	}
		// }

		c.Next()
	}
}
