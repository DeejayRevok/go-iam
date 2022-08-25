package database

import (
	"go-uaa/src/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDbRepository struct {
	db *gorm.DB
}

func (repo *UserDbRepository) Save(user user.User) error {
	result := repo.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user)
	return result.Error
}

func (repo *UserDbRepository) FindByID(userID uuid.UUID) (*user.User, error) {
	var foundUser user.User
	result := repo.db.Preload("Roles").Preload("Roles.Permissions").Where(user.User{ID: userID}).First(&foundUser)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (repo *UserDbRepository) FindByUsername(username string) (*user.User, error) {
	var foundUser user.User
	result := repo.db.Preload("Roles").Preload("Roles.Permissions").Where(user.User{Username: username}).First(&foundUser)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}

func (repo *UserDbRepository) FindByEmail(email string) (*user.User, error) {
	var foundUser user.User
	result := repo.db.Preload("Roles").Preload("Roles.Permissions").Where(user.User{Email: email}).First(&foundUser)
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
