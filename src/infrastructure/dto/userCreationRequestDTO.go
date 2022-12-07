package dto

type UserCreationRequestDTO struct {
	Username *string `json:"username" validate:"required"`
	Email    *string `json:"email" validate:"required"`
	Password *string `json:"password" validate:"required"`
}
