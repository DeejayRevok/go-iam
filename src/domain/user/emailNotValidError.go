package user

import "fmt"

type EmailNotValidError struct {
	Email string
}

func (err EmailNotValidError) Error() string {
	return fmt.Sprintf("Email %s is not valid", err.Email)
}
