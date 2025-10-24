package middleware

import (
	"embox/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware sets up CORS headers for the Gin router.
func CORSMiddleware(cfg *config.ServerConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.CorsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-xsrf-token", "X-XSRF-TOKEN"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// MaxAge specifies how long (in seconds) the results of a preflight request (OPTIONS)
		// can be cached by the browser. Here, 12 * time.Hour means the browser will cache
		// the CORS preflight response for 12 hours, reducing the number of OPTIONS requests.
		MaxAge: cfg.CorsMaxAge,
	})
}
