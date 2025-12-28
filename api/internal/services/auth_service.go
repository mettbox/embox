package services

import (
	"embox/internal/config"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/joho/godotenv/autoload"
)

type AuthService struct {
	config       *config.AuthConfig
	emailService *EmailService
}

type Claims struct {
	jwt.StandardClaims
	User string `json:"user"`
}

func NewAuthService(config *config.AuthConfig, emailService *EmailService) *AuthService {
	return &AuthService{config: config, emailService: emailService}
}

// === public functions ===

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
	loginToken, err := s.generateUserToken(payload, s.config.AccessSecret, s.config.AccessExpiration)
	if err != nil {
		return err
	}

	refreshToken, err := s.generateUserToken(payload, s.config.RefreshSecret, s.config.RefreshExpiration)
	if err != nil {
		return err
	}

	c.SetCookie("access_token", loginToken, s.config.AccessExpiration, "/", s.config.Domain, s.config.IsSecure, true)
	c.SetCookie("refresh_token", refreshToken, s.config.RefreshExpiration, "/", s.config.Domain, s.config.IsSecure, true)

	return nil
}

func (s *AuthService) ValidateAccessCookie(cookie string) (Claims, error) {
	claims, err := s.validateToken(cookie, s.config.AccessSecret)
	if err != nil {
		return Claims{}, err
	}
	return claims, nil
}

func (s *AuthService) ValidateRefreshCookie(cookie string) (Claims, error) {
	claims, err := s.validateToken(cookie, s.config.RefreshSecret)
	if err != nil {
		return Claims{}, err
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

// === private functions ===

func (s *AuthService) generateUserToken(user, secret string, expiration int) (string, error) {
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + int64(expiration), // Add expiration in seconds
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func (s *AuthService) validateToken(tokenString, secret string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}

	return *claims, nil
}
