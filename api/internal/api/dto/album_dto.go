package dto

import "time"

type CreateAlbumRequestDto struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description,omitempty"`
	IsPublic     bool   `json:"isPublic"`
	MediaIDs     []uint `json:"mediaIds,omitempty"`
	CoverMediaID uint   `json:"coverMediaId,omitempty"`
}

type UpdateAlbumRequestDto struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublic    *bool  `json:"isPublic,omitempty"`
}

type AlbumMediaResponseDto struct {
	Id          uint      `json:"id"`
	IsCover     bool      `json:"isCover"`
	IsPublic    bool      `json:"isPublic"`
	IsFavourite bool      `json:"isFavourite"`
	Caption     string    `json:"caption"`
	Date        string    `json:"date"` // yyyy-mm-dd
	Type        string    `json:"type"` // "Image", "Audio", "Video"
	CreatedAt   time.Time `json:"createdAt"`
}

type AlbumResponseDto struct {
	Id          uint                    `json:"id"`
	IsPublic    bool                    `json:"isPublic"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	MediaCount  int                     `json:"mediaCount"`
	Media       []AlbumMediaResponseDto `json:"media"`
}
