package handlers

import (
	"embox/internal/api/response"
	"embox/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FavouriteHandler struct {
	favService *services.FavouriteService
}

func NewFavouriteHandler(favService *services.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{favService}
}

// GetUsersWithLatestFavourite retrieves all users with their latest favourite media
func (h *FavouriteHandler) GetUsersWithLatestFavourite(c *gin.Context) {
	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, 401, "Unauthorized", "")
		return
	}

	responseDtos, err := h.favService.GetUsersWithLatestFavourite(userEmail)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to retrieve users with latest favourites", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"users": responseDtos})
}

func (h *FavouriteHandler) GetFavouriteByUserID(c *gin.Context) {
	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, 401, "Unauthorized", "")
		return
	}

	favUserID := c.Param("id")
	if favUserID == "" {
		response.JSONError(c, http.StatusBadRequest, "Missing user ID", "")
		return
	}

	favouritesResponse, err := h.favService.GetFavouritesByUserID(favUserID, userEmail)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to retrieve favourites", err.Error())
		return
	}

	response.JSONSuccess(c, favouritesResponse)
}

// Add one or more media to the user's favourites
func (h *FavouriteHandler) AddFavourites(c *gin.Context) {
	var payload struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}
	if len(payload.IDs) == 0 {
		response.JSONError(c, http.StatusBadRequest, "No IDs provided", "")
		return
	}

	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, 401, "Unauthorized", "")
		return
	}

	if err := h.favService.AddFavourites(userEmail, payload.IDs); err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to add favourites", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"message": "Favourites added"})
}

// Delete one or more media from the user's favourites
func (h *FavouriteHandler) RemoveFavourites(c *gin.Context) {
	var payload struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}
	if len(payload.IDs) == 0 {
		response.JSONError(c, http.StatusBadRequest, "No IDs provided", "")
		return
	}

	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, 401, "Unauthorized", "")
		return
	}

	if err := h.favService.RemoveFavourites(userEmail, payload.IDs); err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to remove favourites", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"message": "Favourites removed"})
}
