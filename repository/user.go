package repository

import (
	"github.com/marktrs/simple-todo/database"
	"github.com/marktrs/simple-todo/model"
)

// UserRepository - interface for user repository
type UserRepository interface {
	CreateUser(user *model.User) error
	GetByUsername(username string) (*model.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

// CreateUser add a new user record to the database
func (r *userRepository) CreateUser(user *model.User) error {
	db := database.DB
	return db.Create(&user).Error
}

// GetByUsername query a user by username from the database
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	db := database.DB
	var user *model.User

	if err := db.Where(&model.User{Username: username}).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
