package repositories

import (
	"tugaspagik/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(users *[]models.User) error
	FindByID(username string, user *models.User) error
	Create(user models.User) error
	Update(user models.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByID(username string, user *models.User) error {
	return r.db.Where("username = ?", username).First(user).Error
}

func (r *userRepository) FindAll(users *[]models.User) error {
	return r.db.Find(users).Error
}

func (r *userRepository) Create(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) Update(user models.User) error {
	return r.db.Save(&user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.User{}, id).Error
}
