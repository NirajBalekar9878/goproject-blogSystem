package repositories

import (
	"errors"

	"cms-blog-api/config"
	"cms-blog-api/models"

	"gorm.io/gorm"
)

// UserRepository handles database operations for users.
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create inserts a new user into MySQL.
func (r *UserRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// FindByEmail finds a user by email address.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username.
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by ID.
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
