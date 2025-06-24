package services

import (
	"context"

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

// GetUserByID obtiene un usuario de la base de datos por su ID
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.UserRepo.GetUserByID(ctx, userID)
}

// EnsureUserExists verifica si el usuario existe en la base de datos, si no, lo crea
func (s *UserService) EnsureUserExists(ctx context.Context, authUser *models.User) (*models.User, error) {
	// Verificar si el usuario ya existe en la base de datos
	user, err := s.GetUserByID(ctx, authUser.UID)

	// Si hay un error pero NO es del tipo "no encontrado", retornamos el error
	// Si es un error de "no encontrado" o si no hay error pero user es nil, creamos el usuario
	if err == nil && user != nil {
		// Usuario encontrado, lo retornamos
		return user, nil
	}

	// Si llegamos aquí, o hubo un error de "usuario no encontrado" o user es nil,
	// en ambos casos queremos crear un nuevo usuario
	err = s.CreateUser(authUser)
	if err != nil {
		return nil, err
	}
	return authUser, nil
}
