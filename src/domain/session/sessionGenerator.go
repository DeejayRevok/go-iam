package session

import (
	"time"

	"github.com/google/uuid"
)

type SessionGenerator struct{}

func (*SessionGenerator) Generate(userID uuid.UUID) *Session {
	now := time.Now()
	expiration := now.Add(time.Hour * time.Duration(SessionExpirationHours))
	return &Session{
		UserID:    userID.String(),
		Exp:       expiration.Unix(),
		LastUsage: now.Unix(),
	}
}

func NewSessionGenerator() *SessionGenerator {
	return &SessionGenerator{}
}
