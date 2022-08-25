package user

import (
	"time"

	"github.com/google/uuid"
)

type UserPasswordReset struct {
	Token      string    `gorm:"column:token;primaryKey"`
	Expiration time.Time `gorm:"column:expiration"`
	UserID     uuid.UUID `gorm:"column:user_id"`
}

func (reset *UserPasswordReset) IsExpired() bool {
	now := time.Now()
	return now.After(reset.Expiration)
}
