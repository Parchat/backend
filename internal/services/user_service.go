package services

import (
	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/repositories"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(user *models.User) error {
	err := s.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}

	//fmt.Println("User created successfully:", user)

	return nil
}
