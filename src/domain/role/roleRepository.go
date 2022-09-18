package role

import (
	"context"

	"github.com/google/uuid"
)

type RoleRepository interface {
	Save(ctx context.Context, role Role) error
	FindByIDs(ctx context.Context, roleIDs []uuid.UUID) ([]Role, error)
}
