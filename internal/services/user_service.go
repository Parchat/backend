package services

import (
	"context"

	"github.com/Parchat/backend/internal/auth"
	"github.com/Parchat/backend/internal/models"
)

// UserService maneja la lógica de negocio relacionada con usuarios
type UserService struct {
	firebaseAuth *auth.FirebaseAuth
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(firebaseAuth *auth.FirebaseAuth) *UserService {
	return &UserService{
		firebaseAuth: firebaseAuth,
	}
}

// GetUserByID obtiene un usuario por su ID
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	// Obtener el usuario de Firebase
	firebaseUser, err := s.firebaseAuth.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convertir a nuestro modelo de usuario
	user := &models.User{
		ID:        firebaseUser.UID,
		Email:     firebaseUser.Email,
		Name:      firebaseUser.DisplayName,
		PhotoURL:  firebaseUser.PhotoURL,
		CreatedAt: firebaseUser.UserMetadata.CreationTimestamp,
	}

	return user, nil
}
