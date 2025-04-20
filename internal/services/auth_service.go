package services

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/models"
)

// AuthService maneja la lógica de negocio relacionada con usuarios
type AuthService struct {
	FirebaseAuth *config.FirebaseAuth
	UserService  *UserService
}

// NewAuthService crea una nueva instancia de AuthService
func NewAuthService(FirebaseAuth *config.FirebaseAuth, UserService *UserService) *AuthService {
	return &AuthService{
		FirebaseAuth: FirebaseAuth,
		UserService:  UserService,
	}
}

// GetUserByID obtiene un usuario por su ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	// Obtener el usuario de Firebase
	firebaseUser, err := s.FirebaseAuth.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convertir a nuestro modelo de usuario
	user := &models.User{
		UID:         firebaseUser.UID,
		Email:       firebaseUser.Email,
		DisplayName: firebaseUser.DisplayName,
		PhotoURL:    firebaseUser.PhotoURL,
	}

	return user, nil
}

// SignUpAndCreateUser crea un usuario en Firebase Authentication y luego lo guarda en Firestore
func (s *AuthService) SignUpAndCreateUser(password string, user *models.User) error {
	// Crear usuario en Firebase Authentication
	ctx := context.Background()

	params := (&auth.UserToCreate{}).
		DisplayName(user.DisplayName).
		Email(user.Email).
		Password(password)

	authUser, err := s.FirebaseAuth.Client.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	// Asignar UID de Firebase al usuario
	user.UID = authUser.UID

	// Usar el método CreateUser para guardar el usuario en Firestore
	if err := s.UserService.CreateUser(user); err != nil {
		return err
	}

	return nil
}
