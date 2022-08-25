package dto

type RequestResetPasswordDTO struct {
	Email string `json:"email" validate:"required"`
}
