package services

import (
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
	err := s.repo.FindByID(id, &user) // Menyediakan variabel untuk menampung hasil
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userService) CreateUser(user models.User) error {
	return s.repo.Create(user)
}

func (s *userService) UpdateUser(user models.User) error {
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
