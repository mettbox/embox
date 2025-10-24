package dto

type UserWithLatestFavouriteDto struct {
	User struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		MediaCount int    `json:"media_count"`
	} `json:"user"`
	Media struct {
		ID      uint   `json:"id"`
		Caption string `json:"caption"`
		Type    string `json:"type"`
	} `json:"media"`
}

type FavouritesResponseDto struct {
	User  MediaUserResponseDto `json:"user"`
	Media []MediaResponseDto   `json:"media"`
}
