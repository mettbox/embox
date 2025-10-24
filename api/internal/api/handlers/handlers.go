package handlers

import (
	"embox/internal/config"
	"embox/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth      *AuthHandler
	User      *UserHandler
	Media     *MediaHandler
	Favourite *FavouriteHandler
	Album     *AlbumHandler
}

// Init initializes all handlers with the provided API configuration and services.
// It sets up the necessary dependencies for each handler and returns a Handlers struct.
func Init(apiConfig *config.ApiConfig, services *services.Services) *Handlers {
	return &Handlers{
		Auth:      NewAuthHandler(services.Auth, services.User),
		User:      NewUserHandler(services.User, services.Email),
		Media:     NewMediaHandler(services.User, services.Media),
		Favourite: NewFavouriteHandler(services.Favourite),
		Album:     NewAlbumHandler(services.Album),
	}
}

// GetUserEmail retrieves the authenticated user's email from the Gin context.
func GetContextUserEmail(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("user")
	if !exists {
		return "", false
	}
	emailStr, ok := userEmail.(string)
	if !ok {
		return "", false
	}
	return emailStr, true
}
