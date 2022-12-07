package dto

type AuthenticationRequestDTO struct {
	GrantType    string `json:"grant_type" validate:"required"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}
