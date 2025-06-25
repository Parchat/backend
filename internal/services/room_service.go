package services

import (
	"fmt"

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

// GetRoomMessages obtiene los mensajes de una sala con paginación
func (s *RoomService) GetRoomMessages(roomID string, limit int, cursor string) ([]models.MessageResponse, string, error) {
	return s.MessageRepo.GetRoomMessages(roomID, limit, cursor)
}

// GetRoomMessagesSimple obtiene los mensajes de una sala sin paginación
func (s *RoomService) GetRoomMessagesSimple(roomID string, limit int) ([]models.MessageResponse, error) {
	return s.MessageRepo.GetRoomMessagesSimple(roomID, limit)
}

// GetAllRooms obtiene todas las salas ordenadas por fecha de actualización
func (s *RoomService) GetAllRooms() ([]models.Room, error) {
	return s.RoomRepo.GetAllRooms()
}

// JoinRoom permite a un usuario unirse a una sala si tiene permiso
func (s *RoomService) JoinRoom(roomID string, userID string) error {
	// Verificar si el usuario puede unirse a la sala
	// canJoin := s.RoomRepo.CanJoinRoomWebSocket(roomID, userID)
	// if !canJoin {
	// 	return fmt.Errorf("user is not allowed to join this room")
	// }

	// Añadir usuario a la sala
	return s.RoomRepo.AddMemberToRoom(roomID, userID)
}

// IsUserAdminOrOwner checks if a user is an admin or owner of a room
func (s *RoomService) IsUserAdminOrOwner(roomID, userID string) (bool, error) {
	room, err := s.RoomRepo.GetRoom(roomID)
	if err != nil {
		return false, fmt.Errorf("error getting room: %v", err)
	}

	// Check if user is the owner
	if room.OwnerID == userID {
		return true, nil
	}

	// Check if user is an admin
	for _, adminID := range room.Admins {
		if adminID == userID {
			return true, nil
		}
	}

	return false, nil
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
