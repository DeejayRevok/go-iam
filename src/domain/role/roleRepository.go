package role

import "github.com/google/uuid"

type RoleRepository interface {
	Save(role Role) error
	FindByIDs(roleIDs []uuid.UUID) ([]Role, error)
}
