package resetPassword

type ResetPasswordRequest struct {
	UserEmail   string
	ResetToken  string
	NewPassword string
}
