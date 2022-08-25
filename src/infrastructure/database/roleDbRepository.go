package database

import (
	"go-uaa/src/domain/role"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleDbRepository struct {
	db *gorm.DB
}

func (repo *RoleDbRepository) Save(role role.Role) error {
	result := repo.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&role)
	return result.Error
}

func (repo *RoleDbRepository) FindByIDs(roleIDs []uuid.UUID) ([]role.Role, error) {
	var foundRoles []role.Role
	result := repo.db.Where("id IN ?", roleIDs).Find(&foundRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return foundRoles, nil
}

func NewRoleDbRepository(db *gorm.DB) *RoleDbRepository {
	repo := RoleDbRepository{
		db: db,
	}
	return &repo
}
