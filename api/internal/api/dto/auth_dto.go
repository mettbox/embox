package dto

type AuthTokenRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthLoginRequest struct {
	Token string `json:"token" binding:"required"`
}
