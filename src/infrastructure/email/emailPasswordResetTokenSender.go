package email

import (
	"errors"
	"fmt"
	"go-iam/src/domain/user"

	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailPasswordResetTokenSender struct {
	smtpClient *mail.SMTPClient
}

func (sender *EmailPasswordResetTokenSender) Send(resetToken string, receiver *user.User) error {
	if sender.smtpClient == nil {
		return errors.New("SMTP client not provided")
	}
	email := mail.NewMSG()

	email.SetFrom("Go iam <system@iam.com>")

	emailBody := sender.getResetPasswordMessage(resetToken)
	email.SetBody(mail.TextPlain, emailBody)
	email.SetSubject("iam reset password request")

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
