package repositories

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

// GetMessageByID retrieves a message by its ID from a specific room
func (r *MessageRepository) GetMessageByID(roomID, messageID string) (*models.Message, error) {
	ctx := context.Background()

	// Get the message from the room's messages collection
	doc, err := r.FirestoreClient.Client.
		Collection("rooms").Doc(roomID).
		Collection("messages").Doc(messageID).
		Get(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting message: %v", err)
	}

	var message models.Message
	if err := doc.DataTo(&message); err != nil {
		return nil, fmt.Errorf("error converting document to message: %v", err)
	}

	return &message, nil
}

func (r *MessageRepository) GetRoomMessages(roomID string, limit int, cursor string) ([]models.MessageResponse, string, error) {
	ctx := context.Background()

	var messages []models.Message
	var response []models.MessageResponse
	var nextCursor string

	// Crear la consulta base - Orden ASCENDENTE para paginación hacia atrás
	query := r.FirestoreClient.Client.
		Collection("rooms").Doc(roomID).
		Collection("messages").
		OrderBy("createdAt", firestore.Desc). // Cambiado a Asc
		Limit(limit)

	// Si hay un cursor, añadir la condición para empezar desde ese punto
	if cursor != "" {
		// Intentar convertir el cursor a un timestamp
		// Se espera que el cursor sea un timestamp en formato Unix (string)
		timestamp, err := strconv.ParseInt(cursor, 10, 64)
		if err == nil {
			// Convertir el timestamp a time.Time
			cursorTime := time.Unix(timestamp, 0)
			//fmt.Println("Cursor time:", cursorTime.String())

			// Usar EndBefore con el valor convertido
			query = query.StartAfter(cursorTime)
		}
	}

	// Ejecutar la consulta
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, "", fmt.Errorf("error obtaining messages: %v", err)
	}

	// Get messages and track unique user IDs
	userIDs := make(map[string]bool)

	for i, doc := range docs {
		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, "", fmt.Errorf("error decoding message: %v", err)
		}
		messages = append(messages, message)
		userIDs[message.UserID] = true

		// Guardar el último timestamp para el cursor de siguiente página
		if i == len(docs)-1 && len(docs) == limit {
			//fmt.Println("Last message timestamp:", message.CreatedAt.Unix())
			//fmt.Println("Content:", message.Content)
			nextCursor = strconv.FormatInt(message.CreatedAt.Unix(), 10)
		}
	}

	// Map to store user data to avoid duplicate fetches
	userDataCache := make(map[string]string) // userId -> displayName

	// Fetch user data for all unique userIds
	for userID := range userIDs {
		userDoc, err := r.FirestoreClient.Client.Collection("users").Doc(userID).Get(ctx)
		if err != nil {
			// Si hay error, continuamos pero sin el displayName
			continue
		}
		var user models.User
		if err := userDoc.DataTo(&user); err == nil {
			userDataCache[userID] = user.DisplayName
		}
	}

	// Construir la respuesta con los displayNames
	for _, message := range messages {
		msgResponse := models.MessageResponse{
			Message:     message,
			DisplayName: userDataCache[message.UserID], // Puede estar vacío si no se encontró
		}
		response = append(response, msgResponse)
	}

	return response, nextCursor, nil
}

// GetDirectChatMessagesSimple obtiene los mensajes de un chat directo sin paginación
func (r *MessageRepository) GetDirectChatMessagesSimple(directChatID string, limit int) ([]models.MessageResponse, error) {
	ctx := context.Background()

	var messages []models.Message
	var response []models.MessageResponse

	// Obtener mensajes en orden descendente (más recientes primero)
	messagesRef := r.FirestoreClient.Client.
		Collection("directChats").Doc(directChatID).
		Collection("messages").
		OrderBy("createdAt", firestore.Desc).
		Limit(limit)

	docs, err := messagesRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	// Get messages and track unique user IDs
	userIDs := make(map[string]bool)
	for _, doc := range docs {
		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
		userIDs[message.UserID] = true
	}

	// Map to store user data to avoid duplicate fetches
	userDataCache := make(map[string]string) // userId -> displayName

	// Fetch user data for all unique userIds
	for userID := range userIDs {
		userDoc, err := r.FirestoreClient.Client.Collection("users").Doc(userID).Get(ctx)
		if err == nil {
			var user models.User
			if err := userDoc.DataTo(&user); err == nil {
				userDataCache[userID] = user.DisplayName
			}
		}
	}

	// Crear respuestas con DisplayName
	var responseTemp []models.MessageResponse
	for _, message := range messages {
		msgMap := models.MessageResponse{
			Message: message,
		}

		// Add displayName if available in cache
		if displayName, exists := userDataCache[message.UserID]; exists {
			msgMap.DisplayName = displayName
		}

		responseTemp = append(responseTemp, msgMap)
	}

	// Invertir el orden para que queden en orden ascendente (más antiguos primero)
	for i := len(responseTemp) - 1; i >= 0; i-- {
		response = append(response, responseTemp[i])
	}

	return response, nil
}

// GetRoomMessagesSimple obtiene los mensajes de una sala sin paginación
func (r *MessageRepository) GetRoomMessagesSimple(roomID string, limit int) ([]models.MessageResponse, error) {
	ctx := context.Background()

	var messages []models.Message
	var response []models.MessageResponse

	// Obtener mensajes en orden descendente (más recientes primero)
	messagesRef := r.FirestoreClient.Client.
		Collection("rooms").Doc(roomID).
		Collection("messages").
		OrderBy("createdAt", firestore.Desc).
		Limit(limit)

	docs, err := messagesRef.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	// Get messages and track unique user IDs
	userIDs := make(map[string]bool)
	for _, doc := range docs {
		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
		userIDs[message.UserID] = true
	}

	// Map to store user data to avoid duplicate fetches
	userDataCache := make(map[string]string) // userId -> displayName

	// Fetch user data for all unique userIds
	for userID := range userIDs {
		userDoc, err := r.FirestoreClient.Client.Collection("users").Doc(userID).Get(ctx)
		if err != nil {
			// Si hay error, continuamos pero sin el displayName
			continue
		}
		var user models.User
		if err := userDoc.DataTo(&user); err == nil {
			userDataCache[userID] = user.DisplayName
		}
	}

	// Crear respuestas con DisplayName
	var responseTemp []models.MessageResponse
	for _, message := range messages {
		msgResponse := models.MessageResponse{
			Message:     message,
			DisplayName: userDataCache[message.UserID], // Puede estar vacío si no se encontró
		}
		responseTemp = append(responseTemp, msgResponse)
	}

	// Invertir el orden para que queden en orden ascendente (más antiguos primero)
	for i := len(responseTemp) - 1; i >= 0; i-- {
		response = append(response, responseTemp[i])
	}

	return response, nil
}
