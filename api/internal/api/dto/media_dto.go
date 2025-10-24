package dto

import (
	"time"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type MediaUploadRequestDto struct {
	Type        string  `form:"type"` // e.g "image/jpeg"
	FileName    string  `form:"fileName"`
	Date        string  `form:"date"`
	IsPublic    bool    `form:"isPublic"`
	Caption     string  `form:"caption"`
	Orientation string  `form:"orientation"`
	LocationLat float64 `form:"locationLat"`
	LocationLng float64 `form:"locationLng"`
}

type MediaResponseDto struct {
	Id          uint   `json:"id"`
	IsPublic    bool   `json:"isPublic"`
	IsFavourite bool   `json:"isFavourite"`
	Caption     string `json:"caption"`
	Date        string `json:"date"` // yyyy-mm-dd
	Type        string `json:"type"` // "Image", "Audio", "Video"
	// Location    Location  `json:"location"` // currently not uses in UI
	CreatedAt time.Time `json:"createdAt"`
}

/*
Why are pointers useful?
The use of pointers (string, bool) in your DTO is useful because
they allow you to distinguish between "not set" (nil) and "set but empty" ("" or false).
This is especially important for update operations where only specific fields need to be modified.
*/
type MediaUpdateRequestDto struct {
	ID       uint      `json:"id"`
	IsPublic *bool     `json:"isPublic,omitempty"` // Optional: true, false oder nil
	Date     *string   `json:"date,omitempty"`     // Optional: yyyy-mm-dd oder nil
	Caption  *string   `json:"caption,omitempty"`  // Optional: Text oder nil
	Location *Location `json:"location,omitempty"` // Optional: Location-Objekt oder nil
}

type MediaUserResponseDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
