package repositories

import (
	"context"
	"time"

	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/models"
)

type UserRepository struct {
	FirestoreClient *config.FirestoreClient
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository(client *config.FirestoreClient) *UserRepository {
	return &UserRepository{
		FirestoreClient: client,
	}
}

// CreateUser crea un nuevo usuario en la base de datos
func (r *UserRepository) CreateUser(user *models.User) error {
	ctx := context.Background()

	// Asignar timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Guardar el usuario en Firestore
	_, err := r.FirestoreClient.Client.Collection("users").Doc(user.UID).Set(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
