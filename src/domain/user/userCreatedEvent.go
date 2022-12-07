package user

type UserCreatedEvent struct {
	ID        string
	Username  string
	Email     string
	Superuser bool
}
