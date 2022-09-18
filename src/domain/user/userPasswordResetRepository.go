package user

import (
	"context"

	"github.com/google/uuid"
)

type UserPasswordResetRepository interface {
	Save(ctx context.Context, userPasswordReset UserPasswordReset) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*UserPasswordReset, error)
	Delete(ctx context.Context, userPasswordReset UserPasswordReset) error
}
