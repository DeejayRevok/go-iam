package role

import (
	"go-uaa/src/domain/permission"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID               `gorm:"column:id;primaryKey"`
	Name        string                  `gorm:"column:name"`
	Permissions []permission.Permission `gorm:"many2many:role_permission"`
}
