package security

import (
	"go-iam/src/domain/user"

	mail_checker "github.com/FGRibreau/mailchecker/v4/platform/go"
)

type MailCheckerEmailValidator struct{}

func (validator *MailCheckerEmailValidator) Validate(email string) error {
	if !mail_checker.IsValid(email) {
		return user.EmailNotValidError{Email: email}
	}
	return nil
}

func NewMailCheckerEmailValidator() *MailCheckerEmailValidator {
	validator := MailCheckerEmailValidator{}
	return &validator
}
