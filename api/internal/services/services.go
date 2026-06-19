package services

import (
	"embox/internal/config"
	"embox/internal/repositories"
	"log"
)

type Services struct {
	Auth      *AuthService
	User      *UserService
	Email     *EmailService
	Media     *MediaService
	Favourite *FavouriteService
	Album     *AlbumService
}

// Init initializes all services with the provided API configuration and repositories.
// It sets up the necessary dependencies for each service and returns a Services struct.
func Init(apiConfig *config.ApiConfig, repos *repositories.Repositories) *Services {
	var storageService Storage
	if apiConfig.Storage.Adapter == "local" {
		log.Println("INFO: Using local storage adapter")
		storageService = NewLocalStorageAdapter(apiConfig.Storage.LocalDir)
	} else {
		storageService = NewStorageService(apiConfig.Storage)
	}
	emailService := NewEmailService(apiConfig.Email)
	userService := NewUserService(repos.User)
	authService := NewAuthService(apiConfig.Auth, emailService)
	mediaService := NewMediaService(storageService, repos.Media, repos.User)
	favouriteService := NewFavouriteService(repos.User, repos.Favourite)
	albumService := NewAlbumService(repos.User, repos.Album)

	return &Services{
		Auth:      authService,
		User:      userService,
		Email:     emailService,
		Media:     mediaService,
		Favourite: favouriteService,
		Album:     albumService,
	}
}
