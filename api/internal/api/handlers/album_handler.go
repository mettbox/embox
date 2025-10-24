package handlers

import (
	"embox/internal/api/dto"
	"embox/internal/api/response"
	"embox/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	albumService *services.AlbumService
}

func NewAlbumHandler(albumService *services.AlbumService) *AlbumHandler {
	return &AlbumHandler{albumService}
}

func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, 401, "Unauthorized", "")
		return
	}

	var albumRequest dto.CreateAlbumRequestDto

	if err := c.ShouldBindJSON(&albumRequest); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	albumDto, err := h.albumService.CreateAlbum(&albumRequest, userEmail)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to create album", err.Error())
		return
	}

	response.JSONSuccess(c, albumDto)
}

func (h *AlbumHandler) GetAlbumList(c *gin.Context) {
	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, 401, "Unauthorized", "")
		return
	}

	albumList, err := h.albumService.GetAlbumList(userEmail)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to retrieve album list", err.Error())
		return
	}

	response.JSONSuccess(c, albumList)
}

func (h *AlbumHandler) GetAlbumByID(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid album ID", err.Error())
		return
	}

	album, err := h.albumService.GetAlbumByID(uint(id))
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to retrieve album", err.Error())
		return
	}

	response.JSONSuccess(c, album)
}

func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid album ID", err.Error())
		return
	}

	var albumRequest dto.UpdateAlbumRequestDto
	if err := c.ShouldBindJSON(&albumRequest); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	album, err := h.albumService.UpdateAlbum(uint(id), &albumRequest)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to update album", err.Error())
		return
	}

	response.JSONSuccess(c, album)
}

func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid album ID", err.Error())
		return
	}

	err = h.albumService.DeleteAlbum(uint(id))
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to delete album", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"message": "Album deleted successfully"})
}

func (h *AlbumHandler) AddMediaToAlbum(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid album ID", err.Error())
		return
	}

	var payload struct {
		MediaIDs []uint `json:"mediaIds" binding:"required"`
		IsCover  bool   `json:"isCover"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}

	err = h.albumService.AddMediaToAlbum(uint(id), payload.MediaIDs, payload.IsCover)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to add media to album", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"message": "Media added to album successfully"})
}

func (h *AlbumHandler) RemoveMediaFromAlbum(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid album ID", err.Error())
		return
	}

	var payload struct {
		MediaIDs []uint `json:"mediaIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid payload", err.Error())
		return
	}

	err = h.albumService.RemoveMediaFromAlbum(uint(id), payload.MediaIDs)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to remove media from album", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"message": "Media removed from album successfully"})
}
