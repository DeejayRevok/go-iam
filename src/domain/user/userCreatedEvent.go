package user

type UserCreatedEvent struct {
	ID        string
	Username  string
	Email     string
	Superuser bool
}

func (*UserCreatedEvent) EventName() string {
	return "event.user_created"
}
