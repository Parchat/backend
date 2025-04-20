package services

import (
	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/repositories"
)

type UserService struct {
	UserRepo     *repositories.UserRepository
	FirebaseAuth *config.FirebaseAuth
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(userRepo *repositories.UserRepository, firebaseAuth *config.FirebaseAuth) *UserService {
	return &UserService{
		UserRepo:     userRepo,
		FirebaseAuth: firebaseAuth,
	}
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(user *models.User) error {
	err := s.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}
