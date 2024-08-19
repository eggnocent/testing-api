package services

import (
	"errors"
	"strconv"
	"tugaspagik/models"
	"tugaspagik/repositories"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (models.User, error)
	CreateUser(user models.User) error
	UpdateUser(user models.User) error
	DeleteUser(id uint) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	users := []models.User{}
	err := s.repo.FindAll(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetUserByID(id string) (models.User, error) {
	var user models.User

	// Convert id from string to uint
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return user, errors.New("Invalid user ID")
	}

	err = s.repo.FindByID(uint(userID), &user)
	if err != nil {
		return user, errors.New("id tidak ditemukan")
	}
	return user, nil
}

func (s *userService) CreateUser(user models.User) error {
	return s.repo.Create(user)
}

func (s *userService) UpdateUser(user models.User) error {
	var existingUser models.User

	err := s.repo.FindByID(user.ID, &existingUser)
	if err != nil {
		return errors.New("id tidak ditemukan")
	}

	// Lakukan update jika pengguna ditemukan
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
