package routes

import (
	"embox/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterFavouriteRoutes(group *gin.RouterGroup, favouriteHandler *handlers.FavouriteHandler) {
	group.GET("/user", favouriteHandler.GetUsersWithLatestFavourite)
	group.GET("/user/:id", favouriteHandler.GetFavouriteByUserID)
	group.POST("/", favouriteHandler.AddFavourites)
	group.DELETE("/", favouriteHandler.RemoveFavourites)
}
