package services

import (
	"context"

	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/repositories"
)

// DirectChatService maneja la lógica de negocio relacionada con chats directos
type DirectChatService struct {
	DirectChatRepo *repositories.DirectChatRepository
	MessageRepo    *repositories.MessageRepository
}

// NewDirectChatService crea una nueva instancia de DirectChatService
func NewDirectChatService(directChatRepo *repositories.DirectChatRepository, messageRepo *repositories.MessageRepository) *DirectChatService {
	return &DirectChatService{
		DirectChatRepo: directChatRepo,
		MessageRepo:    messageRepo,
	}
}

// CreateDirectChat crea un nuevo chat directo entre usuarios
func (s *DirectChatService) CreateDirectChat(directChat *models.DirectChat) error {
	return s.DirectChatRepo.CreateDirectChat(directChat)
}

// GetDirectChat obtiene un chat directo por su ID
func (s *DirectChatService) GetDirectChat(directChatID string) (*models.DirectChat, error) {
	return s.DirectChatRepo.GetDirectChat(directChatID)
}

// GetUserDirectChats obtiene todos los chats directos de un usuario
func (s *DirectChatService) GetUserDirectChats(userID string) ([]models.DirectChat, error) {
	return s.DirectChatRepo.GetUserDirectChats(userID)
}

// GetUserDirectChatsWithSenderNames obtiene todos los chats directos con nombres de remitentes
func (s *DirectChatService) GetUserDirectChatsWithSenderNames(userID string) ([]models.DirectChat, error) {
	chats, err := s.DirectChatRepo.GetUserDirectChats(userID)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client := s.DirectChatRepo.FirestoreClient.Client

	// Para cada chat, obtener el nombre del remitente del último mensaje
	for i := range chats {
		if chats[i].LastMessage != nil && chats[i].LastMessage.UserID != "" {
			userDoc, err := client.Collection("users").Doc(chats[i].LastMessage.UserID).Get(ctx)
			if err == nil {
				var user models.User
				if err := userDoc.DataTo(&user); err == nil {
					// Crear una copia del mensaje para añadir el displayName
					messageCopy := *chats[i].LastMessage
					// Agregar el displayName a la estructura Message (se serializa como JSON aunque no esté en la estructura)
					messageCopy.DisplayName = user.DisplayName
					// Reemplazar el mensaje original con la copia que incluye displayName
					chats[i].LastMessage = &messageCopy
				}
			}
		}
	}

	return chats, nil
}

// GetDirectChatMessages obtiene los mensajes de un chat directo
func (s *DirectChatService) GetDirectChatMessages(directChatID string, limit int) ([]models.MessageResponse, error) {
	return s.MessageRepo.GetDirectChatMessages(directChatID, limit)
}

// FindOrCreateDirectChat encuentra un chat directo entre dos usuarios o lo crea si no existe
func (s *DirectChatService) FindOrCreateDirectChat(userID1, userID2 string) (*models.DirectChat, error) {
	return s.DirectChatRepo.FindOrCreateDirectChat(userID1, userID2)
}

// FindOrCreateDirectChatWithSenderName encuentra o crea un chat directo e incluye el nombre del remitente
func (s *DirectChatService) FindOrCreateDirectChatWithSenderName(userID1, userID2 string) (*models.DirectChat, error) {
	chat, err := s.DirectChatRepo.FindOrCreateDirectChat(userID1, userID2)
	if err != nil {
		return nil, err
	}

	// Añadir el displayName al último mensaje si existe
	if chat.LastMessage != nil && chat.LastMessage.UserID != "" {
		ctx := context.Background()
		client := s.DirectChatRepo.FirestoreClient.Client

		userDoc, err := client.Collection("users").Doc(chat.LastMessage.UserID).Get(ctx)
		if err == nil {
			var user models.User
			if err := userDoc.DataTo(&user); err == nil {
				// Crear una copia del mensaje para añadir el displayName
				messageCopy := *chat.LastMessage
				// Agregar el displayName a la estructura Message
				messageCopy.DisplayName = user.DisplayName
				// Reemplazar el mensaje original con la copia
				chat.LastMessage = &messageCopy
			}
		}
	}

	return chat, nil
}

// GetDirectChatWithSenderName obtiene un chat directo por su ID e incluye el nombre del remitente del último mensaje
func (s *DirectChatService) GetDirectChatWithSenderName(directChatID string) (*models.DirectChat, error) {
	chat, err := s.DirectChatRepo.GetDirectChat(directChatID)
	if err != nil {
		return nil, err
	}

	// Añadir el displayName al último mensaje si existe
	if chat.LastMessage != nil && chat.LastMessage.UserID != "" {
		ctx := context.Background()
		client := s.DirectChatRepo.FirestoreClient.Client

		userDoc, err := client.Collection("users").Doc(chat.LastMessage.UserID).Get(ctx)
		if err == nil {
			var user models.User
			if err := userDoc.DataTo(&user); err == nil {
				// Crear una copia del mensaje para añadir el displayName
				messageCopy := *chat.LastMessage
				// Agregar el displayName a la estructura Message
				messageCopy.DisplayName = user.DisplayName
				// Reemplazar el mensaje original con la copia
				chat.LastMessage = &messageCopy
			}
		}
	}

	return chat, nil
}
