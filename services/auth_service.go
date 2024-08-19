package services

import (
	"errors"
	"tugaspagik/models"
	"tugaspagik/repositories"
)

type AuthService interface {
	Authenticate(username, password string) (models.User, error)
}

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Authenticate(username, password string) (models.User, error) {
	var user models.User
	err := s.repo.FindByID(username, &user) // Menyediakan variabel untuk menampung hasil
	if err != nil {
		return user, err
	}
	if user.Password != password {
		return user, errors.New("invalid credentials")
	}
	return user, nil
}
