package routes

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"embox/internal/api/response"
	"embox/internal/config"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

func generateCSRFToken(secret string) (string, error) {
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	encodedNonce := base64.RawURLEncoding.EncodeToString(nonce)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(encodedNonce))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encodedNonce + "." + sig, nil
}

func RegisterCsrfRoutes(router *gin.Engine, config *config.CsrfConfig) {
	router.GET("/csrf-token", func(c *gin.Context) {
		token, err := generateCSRFToken(config.Secret)
		if err != nil {
			response.JSONError(c, 500, "Failed to generate CSRF token", err.Error())
			return
		}
		response.JSONSuccess(c, gin.H{
			"csrf_token": token,
		})
	})
}
