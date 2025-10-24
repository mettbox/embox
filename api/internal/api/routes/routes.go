package routes

import (
	"embox/internal/api/handlers"
	"embox/internal/api/middleware"
	"embox/internal/config"
	"embox/internal/repositories"
	"embox/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter initializes the Gin router with all routes and middleware.
// It sets up the necessary handlers and applies middleware for logging, CORS, CSRF protection and authentication.
func Init(db *gorm.DB, apiConfig *config.ApiConfig) *gin.Engine {
	repos := repositories.Init(db)
	services := services.Init(apiConfig, repos)
	handlers := handlers.Init(apiConfig, services)

	gin.SetMode(apiConfig.Router.ReleaseMode)
	router := gin.New()

	// Recover middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Global middleware
	router.Use(middleware.LoggingMiddleware(apiConfig.Router))
	router.Use(middleware.LanguageMiddleware())
	router.Use(middleware.CORSMiddleware(apiConfig.Server))
	router.Use(middleware.CSRFMiddleware())
	router.Use(middleware.AuthMiddleware(services.Auth))

	// === Public routes ===

	RegisterHealthRoutes(router)
	RegisterCsrfRoutes(router, apiConfig.Csrf)
	RegisterAuthRoutes(router.Group("/auth"), handlers.Auth, apiConfig.Router.RateLimit)

	// === Protected routes ===

	userGroup := router.Group("/user")
	userGroup.Use(middleware.RequireAuthMiddleware())
	RegisterUserRoutes(userGroup, handlers.User)

	mediaGroup := router.Group("/media")
	mediaGroup.Use(middleware.RequireAuthMiddleware())
	RegisterMediaRoutes(mediaGroup, handlers.Media)

	favouriteGroup := router.Group("/favourite")
	favouriteGroup.Use(middleware.RequireAuthMiddleware())
	RegisterFavouriteRoutes(favouriteGroup, handlers.Favourite)

	albumGroup := router.Group("/album")
	albumGroup.Use(middleware.RequireAuthMiddleware())
	RegisterAlbumRoutes(albumGroup, handlers.Album)

	return router
}
