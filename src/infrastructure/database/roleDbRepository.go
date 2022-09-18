package database

import (
	"context"
	"go-uaa/src/domain/role"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleDbRepository struct {
	db *gorm.DB
}

func (repo *RoleDbRepository) Save(ctx context.Context, role role.Role) error {
	db := repo.db.WithContext(ctx)
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&role)
	return result.Error
}

func (repo *RoleDbRepository) FindByIDs(ctx context.Context, roleIDs []uuid.UUID) ([]role.Role, error) {
	var foundRoles []role.Role
	db := repo.db.WithContext(ctx)
	result := db.Where("id IN ?", roleIDs).Find(&foundRoles)
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
