package dto

import (
	"time"
)

type MediaUploadRequestDto struct {
	Type     string `form:"type"` // e.g "image/jpeg"
	FileName string `form:"fileName"`
	Date     string `form:"date"`
	Caption  string `form:"caption"`
}

type MediaResponseDto struct {
	Id          uint      `json:"id"`
	IsFavourite bool      `json:"isFavourite"`
	Caption     string    `json:"caption"`
	Date        string    `json:"date"` // yyyy-mm-dd
	Type        string    `json:"type"` // "Image", "Audio", "Video"
	CreatedAt   time.Time `json:"createdAt"`
}

/*
Why are pointers useful?
The use of pointers (string, bool) in your DTO is useful because
they allow you to distinguish between "not set" (nil) and "set but empty" ("" or false).
This is especially important for update operations where only specific fields need to be modified.
*/
type MediaUpdateRequestDto struct {
	ID      uint    `json:"id"`
	Date    *string `json:"date,omitempty"`    // Optional: yyyy-mm-dd oder nil
	Caption *string `json:"caption,omitempty"` // Optional: Text oder nil
}

type MediaUserResponseDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
