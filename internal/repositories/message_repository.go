package repositories

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/models"
)

// MessageRepository maneja las operaciones de base de datos para los mensajes
type MessageRepository struct {
	FirestoreClient *config.FirestoreClient
}

// NewMessageRepository crea una nueva instancia de MessageRepository
func NewMessageRepository(client *config.FirestoreClient) *MessageRepository {
	return &MessageRepository{
		FirestoreClient: client,
	}
}

// SaveMessage guarda un mensaje de sala en Firestore
func (r *MessageRepository) SaveMessage(message *models.Message) error {
	ctx := context.Background()

	// Guardar el mensaje en la colección de mensajes de la sala
	_, err := r.FirestoreClient.Client.
		Collection("rooms").Doc(message.RoomID).
		Collection("messages").Doc(message.ID).
		Set(ctx, message)

	if err != nil {
		return err
	}

	return nil
}

// SaveDirectMessage guarda un mensaje de chat directo en Firestore
func (r *MessageRepository) SaveDirectMessage(message *models.Message) error {
	ctx := context.Background()

	// Guardar el mensaje en la colección de mensajes del chat directo
	_, err := r.FirestoreClient.Client.
		Collection("directChats").Doc(message.RoomID).
		Collection("messages").Doc(message.ID).
		Set(ctx, message)

	if err != nil {
		return err
	}

	return nil
}

// GetRoomMessages obtiene los mensajes de una sala
func (r *MessageRepository) GetRoomMessages(roomID string, limit int) ([]models.Message, error) {
	ctx := context.Background()

	var messages []models.Message

	messagesRef := r.FirestoreClient.Client.
		Collection("rooms").Doc(roomID).
		Collection("messages").
		OrderBy("createdAt", firestore.Desc).
		Limit(limit)

	docs, err := messagesRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetDirectChatMessages obtiene los mensajes de un chat directo
func (r *MessageRepository) GetDirectChatMessages(directChatID string, limit int) ([]models.Message, error) {
	ctx := context.Background()

	var messages []models.Message

	messagesRef := r.FirestoreClient.Client.
		Collection("directChats").Doc(directChatID).
		Collection("messages").
		OrderBy("createdAt", firestore.Desc).
		Limit(limit)

	docs, err := messagesRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
