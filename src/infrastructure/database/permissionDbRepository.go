package database

import (
	"go-uaa/src/domain/permission"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionDbRepository struct {
	db *gorm.DB
}

func (repo *PermissionDbRepository) Save(permission permission.Permission) error {
	result := repo.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&permission)
	return result.Error
}

func (repo *PermissionDbRepository) FindByNames(permissionNames []string) ([]permission.Permission, error) {
	var foundPermissions []permission.Permission
	result := repo.db.Where("name IN ?", permissionNames).Find(&foundPermissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return foundPermissions, nil
}

func NewPermissionDbRepository(db *gorm.DB) *PermissionDbRepository {
	repo := PermissionDbRepository{
		db: db,
	}
	return &repo
}
