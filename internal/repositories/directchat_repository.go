package repositories

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/models"
	"github.com/google/uuid"
)

// DirectChatRepository maneja las operaciones de base de datos para los chats directos
type DirectChatRepository struct {
	FirestoreClient *config.FirestoreClient
	UserRepo        *UserRepository
}

// NewDirectChatRepository crea una nueva instancia de DirectChatRepository
func NewDirectChatRepository(client *config.FirestoreClient, userRepo *UserRepository) *DirectChatRepository {
	return &DirectChatRepository{
		FirestoreClient: client,
		UserRepo:        userRepo,
	}
}

// CreateDirectChat crea un nuevo chat directo entre dos usuarios
func (r *DirectChatRepository) CreateDirectChat(directChat *models.DirectChat) error {
	ctx := context.Background()

	// Asignar ID si no tiene uno
	if directChat.ID == "" {
		directChat.ID = uuid.New().String()
	}

	// Asignar timestamps
	now := time.Now()
	directChat.CreatedAt = now
	directChat.UpdatedAt = now

	// Guarda el chat directo en Firestore
	_, err := r.FirestoreClient.Client.Collection("directChats").Doc(directChat.ID).Set(ctx, directChat)
	if err != nil {
		return err
	}

	return nil
}

// GetDirectChat obtiene un chat directo por ID
func (r *DirectChatRepository) GetDirectChat(directChatID string) (*models.DirectChat, error) {
	ctx := context.Background()

	docRef := r.FirestoreClient.Client.Collection("directChats").Doc(directChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	var directChat models.DirectChat
	if err := docSnap.DataTo(&directChat); err != nil {
		return nil, err
	}

	return &directChat, nil
}

// UpdateLastMessage actualiza el último mensaje de un chat directo
func (r *DirectChatRepository) UpdateLastMessage(directChatID string, message *models.Message) error {
	ctx := context.Background()

	// Usar firestore.Update correctamente
	_, err := r.FirestoreClient.Client.Collection("directChats").Doc(directChatID).Update(ctx, []firestore.Update{
		{Path: "lastMessage", Value: message},
		{Path: "updatedAt", Value: time.Now()},
	})

	return err
}

// IsUserInDirectChat verifica si un usuario es parte de un chat directo
func (r *DirectChatRepository) IsUserInDirectChat(directChatID string, userID string) bool {
	chat, err := r.GetDirectChat(directChatID)
	if err != nil {
		return false
	}

	// Verificar si el usuario está en la lista de userIDs
	for _, id := range chat.UserIDs {
		if id == userID {
			return true
		}
	}

	return false
}

// GetUserDirectChats obtiene todos los chats directos de un usuario
func (r *DirectChatRepository) GetUserDirectChats(userID string) ([]models.DirectChat, error) {
	ctx := context.Background()

	var chats []models.DirectChat

	// Buscar chats donde el usuario sea parte
	query := r.FirestoreClient.Client.Collection("directChats").Where("userIds", "array-contains", userID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var chat models.DirectChat
		if err := doc.DataTo(&chat); err != nil {
			return nil, err
		}

		// Obtener los nombres actualizados de los usuarios
		chat.DisplayNames = make([]string, len(chat.UserIDs))
		for i, id := range chat.UserIDs {
			user, err := r.UserRepo.GetUserByID(ctx, id)
			if err == nil && user != nil {
				chat.DisplayNames[i] = user.DisplayName
			}
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

// FindOrCreateDirectChat encuentra un chat directo entre dos usuarios o lo crea si no existe
func (r *DirectChatRepository) FindOrCreateDirectChat(userID1, userID2 string) (*models.DirectChat, error) {
	// Primero intentamos encontrar un chat existente
	userChats, err := r.GetUserDirectChats(userID1)
	if err != nil {
		return nil, err
	}

	for _, chat := range userChats {
		if len(chat.UserIDs) == 2 {
			// Verificar si el otro usuario está en este chat
			if (chat.UserIDs[0] == userID1 && chat.UserIDs[1] == userID2) ||
				(chat.UserIDs[0] == userID2 && chat.UserIDs[1] == userID1) {
				return &chat, nil
			}
		}
	}

	// No se encontró, crear uno nuevo
	newChat := &models.DirectChat{
		ID:      uuid.New().String(),
		UserIDs: []string{userID1, userID2},
	}

	err = r.CreateDirectChat(newChat)
	if err != nil {
		return nil, err
	}

	return newChat, nil
}
