package dto

type ResetPasswordDTO struct {
	UserEmail   string `json:"user_email" validate:"required"`
	ResetToken  string `json:"reset_token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}
