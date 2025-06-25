package repositories

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/models"
	"github.com/google/uuid"
)

// RoomRepository maneja las operaciones de base de datos para las salas
type RoomRepository struct {
	FirestoreClient *config.FirestoreClient
}

// NewRoomRepository crea una nueva instancia de RoomRepository
func NewRoomRepository(client *config.FirestoreClient) *RoomRepository {
	return &RoomRepository{
		FirestoreClient: client,
	}
}

// CreateRoom crea una nueva sala en Firestore
func (r *RoomRepository) CreateRoom(room *models.Room) error {
	ctx := context.Background()

	// Asignar ID si no tiene uno
	if room.ID == "" {
		room.ID = uuid.New().String()
	}

	// Asignar timestamps
	now := time.Now()
	room.CreatedAt = now
	room.UpdatedAt = now

	// Guarda la sala en Firestore
	_, err := r.FirestoreClient.Client.Collection("rooms").Doc(room.ID).Set(ctx, room)
	if err != nil {
		return err
	}

	return nil
}

// GetRoom obtiene una sala por ID
func (r *RoomRepository) GetRoom(roomID string) (*models.Room, error) {
	ctx := context.Background()

	docRef := r.FirestoreClient.Client.Collection("rooms").Doc(roomID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	var room models.Room
	if err := docSnap.DataTo(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

// UpdateLastMessage actualiza el último mensaje de una sala
func (r *RoomRepository) UpdateLastMessage(roomID string, message *models.Message) error {
	ctx := context.Background()

	// Usar firestore.Update correctamente
	_, err := r.FirestoreClient.Client.Collection("rooms").Doc(roomID).Update(ctx, []firestore.Update{
		{Path: "lastMessage", Value: message},
		{Path: "updatedAt", Value: time.Now()},
	})

	return err
}

// HasRoomAccess verifica si un usuario tiene acceso a una sala
func (r *RoomRepository) HasRoomAccess(room *models.Room, userID string) bool {
	// Verificar si el usuario es miembro
	for _, member := range room.Members {
		if member == userID {
			return true
		}
	}

	// Verificar si el usuario es admin
	for _, admin := range room.Admins {
		if admin == userID {
			return true
		}
	}

	// Verificar si el usuario es el propietario
	return room.OwnerID == userID
}

// CanJoinRoomWebSocket verifica si un usuario puede conectarse a una sala por WebSocket
func (r *RoomRepository) CanJoinRoomWebSocket(roomID string, userID string) bool {
	room, err := r.GetRoom(roomID)
	if err != nil {
		return false
	}

	// Si la sala no es privada, cualquiera puede unirse
	if !room.IsPrivate {
		return true
	}

	return r.HasRoomAccess(room, userID)
}

// CanTalkInRoom verifica si un usuario puede hablar en una sala
func (r *RoomRepository) CanTalkInRoomWebSocket(roomID string, userID string) bool {
	room, err := r.GetRoom(roomID)
	if err != nil {
		return false
	}

	return r.HasRoomAccess(room, userID)
}

// GetUserRooms obtiene todas las salas a las que pertenece un usuario
func (r *RoomRepository) GetUserRooms(userID string) ([]models.Room, error) {
	ctx := context.Background()

	var rooms []models.Room

	// Buscar salas donde el usuario sea miembro
	memberQuery := r.FirestoreClient.Client.Collection("rooms").Where("members", "array-contains", userID)
	memberDocs, err := memberQuery.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range memberDocs {
		var room models.Room
		if err := doc.DataTo(&room); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	// Buscar salas donde el usuario sea admin
	adminQuery := r.FirestoreClient.Client.Collection("rooms").Where("admins", "array-contains", userID)
	adminDocs, err := adminQuery.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range adminDocs {
		var room models.Room
		if err := doc.DataTo(&room); err != nil {
			continue // Ya podría estar en la lista
		}

		// Verificar que no esté duplicada (ya añadida como miembro)
		isDuplicate := false
		for _, existingRoom := range rooms {
			if existingRoom.ID == room.ID {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			rooms = append(rooms, room)
		}
	}

	// Buscar salas donde el usuario sea propietario
	ownerQuery := r.FirestoreClient.Client.Collection("rooms").Where("ownerId", "==", userID)
	ownerDocs, err := ownerQuery.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range ownerDocs {
		var room models.Room
		if err := doc.DataTo(&room); err != nil {
			continue
		}

		// Verificar que no esté duplicada
		isDuplicate := false
		for _, existingRoom := range rooms {
			if existingRoom.ID == room.ID {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			rooms = append(rooms, room)
		}
	}

	return rooms, nil
}

// GetAllRooms obtiene todas las salas ordenadas por fecha de actualización
func (r *RoomRepository) GetAllRooms() ([]models.Room, error) {
	ctx := context.Background()

	// Consultar todas las salas ordenadas por updatedAt descendente
	query := r.FirestoreClient.Client.Collection("rooms").
		//Where("isDeleted", "==", false).
		OrderBy("updatedAt", firestore.Desc)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var rooms []models.Room
	for _, doc := range docs {
		var room models.Room
		if err := doc.DataTo(&room); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// AddMemberToRoom añade un usuario como miembro de una sala
func (r *RoomRepository) AddMemberToRoom(roomID string, userID string) error {
	ctx := context.Background()

	// Obtener la sala
	room, err := r.GetRoom(roomID)
	if err != nil {
		return err
	}

	// Verificar si el usuario es miembro
	for _, member := range room.Members {
		if member == userID {
			return fmt.Errorf("user is already a member of the room")
		}
	}

	// Si la sala es privada, no puede unirse
	// if room.IsPrivate {
	// 	return fmt.Errorf("room is private, user cannot join")
	// }

	// Añadir el usuario como miembro y actualizar la fecha de modificación
	_, err = r.FirestoreClient.Client.Collection("rooms").Doc(roomID).Update(ctx, []firestore.Update{
		{Path: "members", Value: firestore.ArrayUnion(userID)},
		{Path: "updatedAt", Value: time.Now()},
	})

	return err
}
