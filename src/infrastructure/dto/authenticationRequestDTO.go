package dto

type AuthenticationRequestDTO struct {
	GrantType    string `json:"grant_type" validate:"required"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}
