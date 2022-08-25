package user

type PasswordResetTokenSender interface {
	Send(resetToken string, receiver *User) error
}
