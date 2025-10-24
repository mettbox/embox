package handlers

import (
	"embox/internal/api/dto"
	"embox/internal/api/response"
	"embox/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
	userService *services.UserService
}

// NewUserHandler creates a new instance of AuthHandler
func NewAuthHandler(authService *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{authService, userService}
}

// Generate a token and send it via email
func (h *AuthHandler) GenerateToken(c *gin.Context) {
	var req dto.AuthTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.userService.GenerateToken(req.Email)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Failed to generate token", err.Error())
		return
	}

	if err := h.authService.SendLoginTokenEmail(c, user.Email, user.Name, user.Token); err != nil {
		fmt.Printf("Error sending email: %v", err)
		response.JSONError(c, http.StatusInternalServerError, "Failed to send email", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"message": "Token generated successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.AuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.userService.GetByToken(req.Token, h.authService.GetLoginTokenExpiration())
	if err != nil {
		if err.Error() == "token expired" {
			response.JSONError(c, http.StatusUnauthorized, "Token expired", err.Error())
			return
		}
		response.JSONError(c, http.StatusUnauthorized, "Failed to login", err.Error())
		return
	}

	err = h.authService.SetCookie(c, user.Email)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Error setting cookie", err.Error())
		return
	}

	response.JSONSuccess(c, user)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshCookie, err := h.authService.GetRefreshCookie(c)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "Refresh token not found", err.Error())
		return
	}
	claims, err := h.authService.ValidateRefreshCookie(refreshCookie)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "Invalid refresh token", err.Error())
		return
	}

	user, err := h.userService.GetUserByEmail(claims.User)
	if err != nil {
		response.JSONError(c, http.StatusUnauthorized, "User not found", err.Error())
		return
	}
	err = h.authService.SetCookie(c, user.Email)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Error setting cookie", err.Error())
		return
	}

	response.JSONSuccess(c, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.authService.UnsetCookie(c)
	response.JSONSuccess(c, gin.H{"message": "Logged out successfully"})
}
