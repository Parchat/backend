package services

import (
	"context"

	"github.com/Parchat/backend/internal/auth"
	"github.com/Parchat/backend/internal/models"
)

// AuthService maneja la lógica de negocio relacionada con usuarios
type AuthService struct {
	firebaseAuth *auth.FirebaseAuth
}

// NewAuthService crea una nueva instancia de AuthService
func NewAuthService(firebaseAuth *auth.FirebaseAuth) *AuthService {
	return &AuthService{
		firebaseAuth: firebaseAuth,
	}
}

// GetUserByID obtiene un usuario por su ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	// Obtener el usuario de Firebase
	firebaseUser, err := s.firebaseAuth.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convertir a nuestro modelo de usuario
	user := &models.User{
		ID:          firebaseUser.UID,
		Email:       firebaseUser.Email,
		DisplayName: firebaseUser.DisplayName,
		PhotoURL:    firebaseUser.PhotoURL,
	}

	return user, nil
}
