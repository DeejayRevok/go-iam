package user

type EmailValidator interface {
	Validate(email string) error
}
