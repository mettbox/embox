package routes

import (
	"embox/internal/api/response"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		response.JSONSuccess(c, gin.H{
			"message": "OK",
		})
	})
}
