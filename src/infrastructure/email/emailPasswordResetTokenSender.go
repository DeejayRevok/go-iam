package email

import (
	"fmt"
	"go-uaa/src/domain/user"

	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailPasswordResetTokenSender struct {
	smtpClient *mail.SMTPClient
}

func (sender *EmailPasswordResetTokenSender) Send(resetToken string, receiver *user.User) error {
	email := mail.NewMSG()

	email.SetFrom("Go UAA <system@go-uaa.com>")

	emailBody := sender.getResetPasswordMessage(resetToken)
	email.SetBody(mail.TextPlain, emailBody)
	email.SetSubject("UAA reset password request")

	email.AddTo(receiver.Email)

	return email.Send(sender.smtpClient)
}

func (*EmailPasswordResetTokenSender) getResetPasswordMessage(resetToken string) string {
	return fmt.Sprintf("The code to reset your password is %s", resetToken)
}

func NewEmailPasswordResetTokenSender(smtpClient *mail.SMTPClient) *EmailPasswordResetTokenSender {
	return &EmailPasswordResetTokenSender{
		smtpClient: smtpClient,
	}
}
