package services

import (
	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/repositories"
	"github.com/google/uuid"
)

// RoomService maneja la lógica de negocio relacionada con salas de chat
type RoomService struct {
	RoomRepo    *repositories.RoomRepository
	MessageRepo *repositories.MessageRepository
}

// NewRoomService crea una nueva instancia de RoomService
func NewRoomService(roomRepo *repositories.RoomRepository, messageRepo *repositories.MessageRepository) *RoomService {
	return &RoomService{
		RoomRepo:    roomRepo,
		MessageRepo: messageRepo,
	}
}

// CreateRoom crea una nueva sala de chat
func (s *RoomService) CreateRoom(room *models.Room) error {
	// Asignar ID si no tiene
	if room.ID == "" {
		room.ID = uuid.New().String()
	}

	// Garantizar que el creador es admin y miembro
	if !contains(room.Admins, room.OwnerID) {
		room.Admins = append(room.Admins, room.OwnerID)
	}

	if !contains(room.Members, room.OwnerID) {
		room.Members = append(room.Members, room.OwnerID)
	}

	return s.RoomRepo.CreateRoom(room)
}

// GetRoom obtiene una sala por su ID
func (s *RoomService) GetRoom(roomID string) (*models.Room, error) {
	return s.RoomRepo.GetRoom(roomID)
}

// GetUserRooms obtiene todas las salas a las que pertenece un usuario
func (s *RoomService) GetUserRooms(userID string) ([]models.Room, error) {
	return s.RoomRepo.GetUserRooms(userID)
}

// GetRoomMessages obtiene los mensajes de una sala
func (s *RoomService) GetRoomMessages(roomID string, limit int) ([]models.Message, error) {
	return s.MessageRepo.GetRoomMessages(roomID, limit)
}

// Helper para verificar si un slice contiene un valor
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
