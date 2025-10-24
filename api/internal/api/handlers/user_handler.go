package handlers

import (
	"embox/internal/api/dto"
	"embox/internal/api/response"
	"embox/internal/models"
	"embox/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler contains the service for user operations
type UserHandler struct {
	userService  *services.UserService
	emailService *services.EmailService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *services.UserService, emailService *services.EmailService) *UserHandler {
	return &UserHandler{userService, emailService}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user := &models.User{
		Name:    req.Name,
		Email:   req.Email,
		IsAdmin: req.IsAdmin,
	}

	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, err.Error(), "")
		return
	}

	if req.Message != "" && req.Subject != "" {
		h.sendEmail(c, createdUser.Email, req.Subject, req.Message)
	}

	response.JSONCreated(c, createdUser)
}

// GetAllUsers returns all users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to retrieve users", err.Error())
		return
	}

	response.JSONSuccess(c, users)
}

// GetUserById returns a user by ID
func (h *UserHandler) GetUserById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.userService.GetUserById(id)
	if err != nil {
		if err.Error() == "user not found" {
			response.JSONError(c, http.StatusNotFound, "User not found", "")
		} else {
			response.JSONError(c, http.StatusInternalServerError, "Failed to retrieve user", err.Error())
		}
		return
	}

	response.JSONSuccess(c, user)
}

// GetUserByEmail returns a user by email
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "user not found" {
			response.JSONError(c, http.StatusNotFound, "User not found", "")
		} else {
			response.JSONError(c, http.StatusInternalServerError, "Failed to retrieve user", err.Error())
		}
		return
	}

	response.JSONSuccess(c, user)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	updatedUser, err := h.userService.UpdateUser(id, req)
	if err != nil {
		if err.Error() == "user not found" {
			response.JSONError(c, http.StatusNotFound, "User not found", "")
		} else {
			response.JSONError(c, http.StatusInternalServerError, "Failed to update user", err.Error())
		}
		return
	}

	if req.Message != nil && req.Subject != nil && *req.Message != "" && *req.Subject != "" {
		h.sendEmail(c, updatedUser.Email, *req.Subject, *req.Message)
	}

	response.JSONSuccess(c, updatedUser)
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	err = h.userService.DeleteUser(uuid)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) sendEmail(c *gin.Context, email, subject, body string) {
	lang, _ := c.Get("lang")
	if err := h.emailService.SendEmail(lang, email, subject, body); err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to send email", err.Error())
		return
	}
}
