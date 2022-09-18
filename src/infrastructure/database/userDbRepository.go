package database

import (
	"context"
	"go-uaa/src/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDbRepository struct {
	db *gorm.DB
}

func (repo *UserDbRepository) Save(ctx context.Context, user user.User) error {
	db := repo.db.WithContext(ctx)
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user)
	return result.Error
}

func (repo *UserDbRepository) FindByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	var foundUser user.User
	db := repo.db.WithContext(ctx)
	result := db.Preload("Roles").Preload("Roles.Permissions").Where(user.User{ID: userID}).First(&foundUser)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (repo *UserDbRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var foundUser user.User
	db := repo.db.WithContext(ctx)
	result := db.Preload("Roles").Preload("Roles.Permissions").Where(user.User{Username: username}).First(&foundUser)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (repo *UserDbRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var foundUser user.User
	db := repo.db.WithContext(ctx)
	result := db.Preload("Roles").Preload("Roles.Permissions").Where(user.User{Email: email}).First(&foundUser)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func NewUserDbRepository(db *gorm.DB) *UserDbRepository {
	repo := UserDbRepository{
		db: db,
	}
	return &repo
}
