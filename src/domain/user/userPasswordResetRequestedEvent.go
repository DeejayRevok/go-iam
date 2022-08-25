package user

type UserPasswordResetRequestedEvent struct {
	ResetToken string
	UserID     string
}

func UserPasswordResetRequestedEventFromMap(eventMap map[string]string) *UserPasswordResetRequestedEvent {
	return &UserPasswordResetRequestedEvent{
		ResetToken: eventMap["ResetToken"],
		UserID:     eventMap["UserID"],
	}
}
