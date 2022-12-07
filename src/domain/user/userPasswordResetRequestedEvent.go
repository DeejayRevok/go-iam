package user

type UserPasswordResetRequestedEvent struct {
	ResetToken string
	UserID     string
}

func UserPasswordResetRequestedEventFromMap(eventMap map[string]interface{}) *UserPasswordResetRequestedEvent {
	return &UserPasswordResetRequestedEvent{
		ResetToken: eventMap["ResetToken"].(string),
		UserID:     eventMap["UserID"].(string),
	}
}
