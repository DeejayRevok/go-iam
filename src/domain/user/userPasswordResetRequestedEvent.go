package user

type UserPasswordResetRequestedEvent struct {
	ResetToken string
	UserID     string
}

func (*UserPasswordResetRequestedEvent) EventName() string {
	return "event.user_password_reset_requested"
}

func UserPasswordResetRequestedEventFromMap(eventMap map[string]interface{}) *UserPasswordResetRequestedEvent {
	return &UserPasswordResetRequestedEvent{
		ResetToken: eventMap["ResetToken"].(string),
		UserID:     eventMap["UserID"].(string),
	}
}
