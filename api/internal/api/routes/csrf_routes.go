package routes

import (
	"crypto/rand"
	"embox/internal/api/response"
	"embox/internal/config"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

func generateCSRFToken() (string, error) {
	b := make([]byte, 32) // 256 Bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func RegisterCsrfRoutes(router *gin.Engine, config *config.CsrfConfig) {
	router.GET("/csrf-token", func(c *gin.Context) {
		token, err := generateCSRFToken()
		if err != nil {
			response.JSONError(c, 500, "Failed to generate CSRF token", err.Error())
			return
		}
		c.SetCookie("XSRF-TOKEN", token, config.MaxAge, "/", config.Domain, config.IsSecure, false)
		response.JSONSuccess(c, gin.H{
			"csrf_token": token,
		})
	})
}
