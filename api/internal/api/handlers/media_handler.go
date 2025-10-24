package handlers

import (
	"embox/internal/api/dto"
	"embox/internal/api/response"
	"embox/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	userService  *services.UserService
	mediaService *services.MediaService
}

func NewMediaHandler(userService *services.UserService, mediaService *services.MediaService) *MediaHandler {
	return &MediaHandler{userService, mediaService}
}

func (h *MediaHandler) GetMediaList(c *gin.Context) {
	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, 401, "Unauthorized", "")
		return
	}

	results, err := h.mediaService.GetMediaList(userEmail)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to fetch media", err.Error())
		return
	}

	response.JSONSuccess(c, results)
}

// Get media thumbnail as blob by ID
func (h *MediaHandler) GetMediaThumbnail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid media ID", err.Error())
		return
	}

	data, mimeType, err := h.mediaService.GetThumbnail(uint(id))
	if err != nil {
		response.JSONError(c, http.StatusNotFound, "File not found", err.Error())
		return
	}

	c.Data(200, mimeType, data)
}

// Get media original file from storage by ID
func (h *MediaHandler) GetMediaFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid media ID", err.Error())
		return
	}
	data, mimeType, err := h.mediaService.GetMediaFile(uint(id))
	if err != nil {
		response.JSONError(c, http.StatusNotFound, "File not found", err.Error())
		return
	}
	c.Data(200, mimeType, data)
}

// Upload one or multiple media files
func (h *MediaHandler) UploadMedia(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.JSONError(c, 400, "Invalid multipart form", err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.JSONError(c, 400, "No files uploaded", "")
		return
	}

	metaJson := c.PostForm("meta")
	var metaList []dto.MediaUploadRequestDto
	if err := json.Unmarshal([]byte(metaJson), &metaList); err != nil {
		response.JSONError(c, 400, "Invalid meta data", err.Error())
		return
	}
	if len(metaList) != len(files) {
		response.JSONError(c, 400, "Meta and files count mismatch", "")
		return
	}

	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, 401, "Unauthorized", "")
		return
	}

	uploaded, err := h.mediaService.UploadMedia(files, metaList, userEmail)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to upload media", err.Error())
		return
	}

	response.JSONSuccess(c, uploaded)
}

// Update one or multiple media items
func (h *MediaHandler) UpdateMedia(c *gin.Context) {
	userEmail, ok := GetContextUserEmail(c)
	if !ok {
		response.JSONError(c, http.StatusUnauthorized, "Unauthorized", "")
		return
	}

	var payload struct {
		Updates []dto.MediaUpdateRequestDto `json:"updates"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		response.JSONError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	updatedMedia, err := h.mediaService.UpdateMediaBatch(payload.Updates, userEmail)
	if err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to update media", err.Error())
		return
	}

	response.JSONSuccess(c, updatedMedia)
}

// Delete one or multiple media files (IDs per JSON/body)
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
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

	if err := h.mediaService.DeleteMedia(payload.IDs); err != nil {
		response.JSONError(c, http.StatusInternalServerError, "Failed to delete media", err.Error())
		return
	}

	response.JSONSuccess(c, gin.H{"message": "Media deleted successfully"})
}
