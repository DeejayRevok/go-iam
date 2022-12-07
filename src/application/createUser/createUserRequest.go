package createUser

type CreateUserRequest struct {
	Username  string
	Email     string
	Password  string
	Superuser bool
}
