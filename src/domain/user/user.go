package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"column:id;primaryKey"`
	Username     string    `gorm:"column:username"`
	Email        string    `gorm:"column:email"`
	Password     string    `gorm:"column:password"`
	RefreshToken string    `gorm:"column:refresh_token"`
	Superuser    bool      `gorm:"column:superuser"`
}
