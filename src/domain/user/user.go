package user

import (
	"go-uaa/src/domain/role"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID   `gorm:"column:id;primaryKey"`
	Username         string      `gorm:"column:username"`
	Email            string      `gorm:"column:email"`
	Password         string      `gorm:"column:password"`
	Roles            []role.Role `gorm:"many2many:user_role"`
	RefreshToken     string      `gorm:"column:refresh_token"`
	Superuser        bool        `gorm:"column:superuser"`
	permissionsIndex map[string]struct{}
}

func (user *User) indexPermissions() {
	if user.permissionsIndex == nil {
		return
	}

	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			user.permissionsIndex[permission.Name] = struct{}{}
		}
	}
}

func (user *User) HasPermission(permission string) bool {
	if user.permissionsIndex == nil {
		user.indexPermissions()
	}

	_, hasPermission := user.permissionsIndex[permission]
	return hasPermission
}
