package user

import "github.com/google/uuid"

type UserPasswordResetRepository interface {
	Save(userPasswordReset UserPasswordReset) error
	FindByUserID(userID uuid.UUID) (*UserPasswordReset, error)
	Delete(userPasswordReset UserPasswordReset) error
}
