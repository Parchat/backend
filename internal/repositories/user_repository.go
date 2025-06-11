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

// GetUserByID obtiene un usuario de la base de datos por su ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	docRef := r.FirestoreClient.Client.Collection("users").Doc(userID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !docSnap.Exists() {
		return nil, nil // El usuario no existe
	}

	var user models.User
	if err := docSnap.DataTo(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
