package services

import (
	"context"

	"firebase.google.com/go/v4/auth"
	authP "github.com/Parchat/backend/internal/auth"
	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/repositories"
)

type UserService struct {
	UserRepo     *repositories.UserRepository
	FirebaseAuth *authP.FirebaseAuth
}

// UserToCreate represents the data required to create a user in Firebase Auth
type UserToCreate struct {
	Email    string
	Password string
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(userRepo *repositories.UserRepository, firebaseAuth *authP.FirebaseAuth) *UserService {
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

	//fmt.Println("User created successfully:", user)

	return nil
}

// CreateUserWithAuth crea un usuario en Firebase Authentication y Firestore
func (s *UserService) CreateUserWithAuth(password string, user *models.User) error {
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

	// Crear usuario en Firestore
	if err := s.UserRepo.CreateUser(user); err != nil {
		return err
	}

	return nil
}
