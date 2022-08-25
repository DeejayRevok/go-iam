package user

import "github.com/google/uuid"

type UserRepository interface {
	Save(user User) error
	FindByID(userID uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
}
