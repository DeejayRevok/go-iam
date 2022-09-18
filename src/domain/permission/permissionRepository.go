package permission

import "context"

type PermissionRepository interface {
	Save(ctx context.Context, permission Permission) error
	FindByNames(ctx context.Context, permissionNames []string) ([]Permission, error)
}
