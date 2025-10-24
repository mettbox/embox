package services

import (
	"embox/internal/config"
	"embox/pkg/jwt"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type AuthService struct {
	config       *config.AuthConfig
	emailService *EmailService
}

func NewAuthService(config *config.AuthConfig, emailService *EmailService) *AuthService {
	return &AuthService{config: config, emailService: emailService}
}

func (s *AuthService) GetAccessCookie(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return cookie, nil
}

func (s *AuthService) GetRefreshCookie(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return "", err
	}
	return cookie, nil
}

func (s *AuthService) SetCookie(c *gin.Context, payload string) error {
	loginToken, err := jwt.GenerateUserToken(payload, s.config.AccessSecret, s.config.AccessExpiration)
	if err != nil {
		return err
	}

	refreshToken, err := jwt.GenerateUserToken(payload, s.config.RefreshSecret, s.config.RefreshExpiration)
	if err != nil {
		return err
	}

	c.SetCookie("access_token", loginToken, s.config.AccessExpiration, "/", s.config.Domain, s.config.IsSecure, true)
	c.SetCookie("refresh_token", refreshToken, s.config.RefreshExpiration, "/", s.config.Domain, s.config.IsSecure, true)

	return nil
}

func (s *AuthService) ValidateAccessCookie(cookie string) (jwt.Claims, error) {
	claims, err := jwt.ValidateToken(cookie, s.config.AccessSecret)
	if err != nil {
		return jwt.Claims{}, err
	}
	return claims, nil
}

func (s *AuthService) ValidateRefreshCookie(cookie string) (jwt.Claims, error) {
	claims, err := jwt.ValidateToken(cookie, s.config.RefreshSecret)
	if err != nil {
		return jwt.Claims{}, err
	}
	return claims, nil
}

func (s *AuthService) UnsetCookie(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", s.config.Domain, s.config.IsSecure, true)
	c.SetCookie("refresh_token", "", -1, "/", s.config.Domain, s.config.IsSecure, true)
}

func (s *AuthService) SendLoginTokenEmail(c *gin.Context, email, name, token string) error {
	lang, _ := c.Get("lang")

	minutes := s.config.LoginTokenExpiration / 60
	body := fmt.Sprintf(s.config.LoginEmailTemplate, name, token, minutes)

	if err := s.emailService.SendEmail(lang, email, s.config.LoginEmailSubject, body); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetLoginTokenExpiration() time.Duration {
	return time.Duration(s.config.LoginTokenExpiration) * time.Second
}
