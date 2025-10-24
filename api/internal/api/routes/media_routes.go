package routes

import (
	"embox/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterMediaRoutes(group *gin.RouterGroup, mediaHandler *handlers.MediaHandler) {
	group.GET("/", mediaHandler.GetMediaList)
	group.GET("/:id/thumbnail", mediaHandler.GetMediaThumbnail)
	group.GET("/:id/file", mediaHandler.GetMediaFile)
	group.POST("/", mediaHandler.UploadMedia)
	group.PUT("/", mediaHandler.UpdateMedia)
	group.DELETE("/", mediaHandler.DeleteMedia)
}
