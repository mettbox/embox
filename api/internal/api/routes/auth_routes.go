package routes

import (
	"embox/internal/api/handlers"
	"embox/internal/api/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(group *gin.RouterGroup, authHandler *handlers.AuthHandler, rateLimit int) {
	group.POST("/token", middleware.RateLimitMiddleware(rateLimit, time.Minute), authHandler.GenerateToken)
	group.POST("/login", middleware.RateLimitMiddleware(rateLimit, time.Minute), authHandler.Login)
	group.POST("/refresh", middleware.RateLimitMiddleware(rateLimit, time.Minute), authHandler.Refresh)
	group.POST("/logout", authHandler.Logout)
}
