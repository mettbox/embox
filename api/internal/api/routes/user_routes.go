package routes

import (
	"embox/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// Binds all admin user routes to the router
func RegisterUserRoutes(group *gin.RouterGroup, userHandler *handlers.UserHandler) {
	group.GET("/", userHandler.GetAllUsers)
	group.POST("/", userHandler.CreateUser)
	group.PUT("/:id", userHandler.UpdateUser)
	group.DELETE("/:id", userHandler.DeleteUser)
}
