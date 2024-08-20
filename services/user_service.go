package services

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"tugaspagik/models"
	"tugaspagik/repositories"

	"github.com/go-redis/redis"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (models.User, error)
	CreateUser(user models.User) error
	UpdateUser(user models.User) error
}

type userService struct {
	repo repositories.UserRepository
	rdb  *redis.Client
}

func NewUserService(repo repositories.UserRepository, rdb *redis.Client) UserService {
	return &userService{repo, rdb}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	var users []models.User

	// Cek apakah data ada di cache Redis
	cachedUsers, err := s.rdb.Get("all_users").Result()
	if err == redis.Nil { // Jika tidak ada di cache
		// Ambil dari database
		err := s.repo.FindAll(&users)
		if err != nil {
			return nil, err
		}

		// Simpan hasil ke Redis sebagai JSON
		userJSON, _ := json.Marshal(users)
		s.rdb.Set("all_users", userJSON, 10*time.Minute) // Cache selama 10 menit
		return users, nil
	} else if err != nil {
		return nil, err
	}

	// Jika ada di cache, ambil dari cache
	json.Unmarshal([]byte(cachedUsers), &users)
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
