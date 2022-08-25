package permission

type PermissionRepository interface {
	Save(permission Permission) error
	FindByNames(permissionNames []string) ([]Permission, error)
}
