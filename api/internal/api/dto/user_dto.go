package dto

type UserResponse struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	IsAdmin             bool   `json:"isAdmin"`
	HasPublicFavourites *bool  `json:"hasPublicFavourites,omitempty"`
}

type CreateUserRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	IsAdmin bool   `json:"isAdmin"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
}

type UpdateUserRequest struct {
	Name                *string `json:"name,omitempty"`
	Email               *string `json:"email,omitempty" binding:"omitempty,email"`
	IsAdmin             *bool   `json:"isAdmin,omitempty"`
	Subject             *string `json:"subject,omitempty"`
	Message             *string `json:"message,omitempty"`
	HasPublicFavourites *bool   `json:"hasPublicFavourites,omitempty"`
}

type UserWithLatestFavourite struct {
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Media struct {
		ID      uint   `json:"id"`
		Caption string `json:"caption"`
	} `json:"media"`
}
