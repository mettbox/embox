package routes

import (
	"embox/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAlbumRoutes(group *gin.RouterGroup, albumHandler *handlers.AlbumHandler) {
	group.POST("/", albumHandler.CreateAlbum)
	group.PUT("/:id", albumHandler.UpdateAlbum)
	group.PUT("/:id/cover", albumHandler.SetCover)
	group.GET("/", albumHandler.GetAlbumList)
	group.GET("/:id", albumHandler.GetAlbumByID)
	group.DELETE("/:id", albumHandler.DeleteAlbum)
	group.POST("/:id/media", albumHandler.AddMediaToAlbum)
	group.DELETE("/:id/media", albumHandler.RemoveMediaFromAlbum)
}
