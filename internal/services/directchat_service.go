package services

import (
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

// GetDirectChatMessages obtiene los mensajes de un chat directo
func (s *DirectChatService) GetDirectChatMessages(directChatID string, limit int) ([]models.MessageResponse, error) {
	return s.MessageRepo.GetDirectChatMessages(directChatID, limit)
}

// FindOrCreateDirectChat encuentra un chat directo entre dos usuarios o lo crea si no existe
func (s *DirectChatService) FindOrCreateDirectChat(userID1, userID2 string) (*models.DirectChat, error) {
	return s.DirectChatRepo.FindOrCreateDirectChat(userID1, userID2)
}
