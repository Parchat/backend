package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

// ChatHandler maneja las peticiones relacionadas con chats
type ChatHandler struct {
	RoomService       *services.RoomService
	DirectChatService *services.DirectChatService
}

// NewChatHandler crea una nueva instancia de ChatHandler
func NewChatHandler(roomService *services.RoomService, directChatService *services.DirectChatService) *ChatHandler {
	return &ChatHandler{
		RoomService:       roomService,
		DirectChatService: directChatService,
	}
}

// CreateRoom crea una nueva sala de chat
//
//	@Summary		Crea una nueva sala de chat
//	@Description	Crea una nueva sala de chat con el usuario actual como propietario
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			room	body		models.CreateRoomRequest	true	"Detalles de la sala"
//	@Success		201		{object}	models.Room					"Sala creada exitosamente"
//	@Failure		400		{string}	string						"Solicitud inválida"
//	@Failure		401		{string}	string						"No autorizado"
//	@Failure		500		{string}	string						"Error interno del servidor"
//	@Router			/chat/rooms [post]
func (h *ChatHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Asignar el creador como propietario
	room.OwnerID = userID

	if err := h.RoomService.CreateRoom(&room); err != nil {
		http.Error(w, "Error creating room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

// GetRoom obtiene una sala de chat por ID
//
//	@Summary		Obtiene una sala por ID
//	@Description	Devuelve los detalles de una sala específica
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			roomId	path		string		true	"ID de la sala"
//	@Success		200		{object}	models.Room	"Detalles de la sala"
//	@Failure		401		{string}	string		"No autorizado"
//	@Failure		404		{string}	string		"Sala no encontrada"
//	@Failure		500		{string}	string		"Error interno del servidor"
//	@Router			/chat/rooms/{roomId} [get]
func (h *ChatHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")

	room, err := h.RoomService.GetRoom(roomID)
	if err != nil {
		http.Error(w, "Error getting room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(room)
}

// GetUserRooms obtiene todas las salas a las que pertenece un usuario
//
//	@Summary		Obtiene las salas del usuario
//	@Description	Devuelve todas las salas a las que pertenece el usuario autenticado
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.Room	"Lista de salas"
//	@Failure		401	{string}	string		"No autorizado"
//	@Failure		500	{string}	string		"Error interno del servidor"
//	@Router			/chat/rooms/me [get]
func (h *ChatHandler) GetUserRooms(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	rooms, err := h.RoomService.GetUserRooms(userID)
	if err != nil {
		http.Error(w, "Error getting rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}

// GetRoomMessages obtiene los mensajes de una sala con paginación ordernada por fecha de creación descendente
//
//	@Summary		Obtiene mensajes de una sala
//	@Description	Devuelve los mensajes de una sala específica con soporte para paginación ordernada por fecha de creación descendente
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			roomId	path		string								true	"ID de la sala"
//	@Param			limit	query		int									false	"Límite de mensajes a obtener"			default(50)
//	@Param			cursor	query		string								false	"Cursor para paginación (timestamp)"	default("1747441934")
//	@Success		200		{object}	models.PaginatedMessagesResponse	"Mensajes paginados de la sala"
//	@Failure		401		{string}	string								"No autorizado"
//	@Failure		404		{string}	string								"Sala no encontrada"
//	@Failure		500		{string}	string								"Error interno del servidor"
//	@Router			/chat/rooms/{roomId}/messages/paginated [get]
func (h *ChatHandler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")

	limitStr := r.URL.Query().Get("limit")
	limit := 50 // valor por defecto

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	cursor := r.URL.Query().Get("cursor")

	messages, nextCursor, err := h.RoomService.GetRoomMessages(roomID, limit, cursor)
	if err != nil {
		http.Error(w, "Error getting messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.PaginatedMessagesResponse{
		Messages:   messages,
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
	}

	json.NewEncoder(w).Encode(response)
}

// CreateDirectChat crea o encuentra un chat directo entre dos usuarios
//
//	@Summary		Crea un chat directo
//	@Description	Crea o encuentra un chat directo entre el usuario autenticado y otro usuario
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			otherUserId	path		string				true	"ID del otro usuario"
//	@Success		201			{object}	models.DirectChat	"Chat directo creado o encontrado"
//	@Failure		400			{string}	string				"Solicitud inválida"
//	@Failure		401			{string}	string				"No autorizado"
//	@Failure		500			{string}	string				"Error interno del servidor"
//	@Router			/chat/direct/{otherUserId} [post]
func (h *ChatHandler) CreateDirectChat(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Obtener el ID del otro usuario desde el parámetro de la URL
	otherUserID := chi.URLParam(r, "otherUserId")
	if otherUserID == "" {
		http.Error(w, "Other user ID is required", http.StatusBadRequest)
		return
	}

	chat, err := h.DirectChatService.FindOrCreateDirectChat(userID, otherUserID)
	if err != nil {
		http.Error(w, "Error creating direct chat: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

// GetUserDirectChats obtiene todos los chats directos del usuario
//
//	@Summary		Obtiene chats directos
//	@Description	Devuelve todos los chats directos del usuario autenticado
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.DirectChat	"Lista de chats directos"
//	@Failure		401	{string}	string				"No autorizado"
//	@Failure		500	{string}	string				"Error interno del servidor"
//	@Router			/chat/direct/me [get]
func (h *ChatHandler) GetUserDirectChats(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	chats, err := h.DirectChatService.GetUserDirectChatsWithSenderNames(userID)
	if err != nil {
		http.Error(w, "Error getting direct chats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chats)
}

// GetDirectChatMessages obtiene los mensajes de un chat directo
//
//	@Summary		Obtiene mensajes de un chat directo
//	@Description	Devuelve los mensajes de un chat directo específico
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			chatId	path		string			true	"ID del chat directo"
//	@Param			limit	query		int				false	"Límite de mensajes a obtener"	default(50)
//	@Success		200		{array}		models.Message	"Lista de mensajes del chat directo"
//	@Failure		401		{string}	string			"No autorizado"
//	@Failure		403		{string}	string			"Acceso prohibido"
//	@Failure		404		{string}	string			"Chat no encontrado"
//	@Failure		500		{string}	string			"Error interno del servidor"
//	@Router			/chat/direct/{chatId}/messages [get]
func (h *ChatHandler) GetDirectChatMessages(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chatId")

	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Verificar si el usuario pertenece al chat directo
	if !h.DirectChatService.DirectChatRepo.IsUserInDirectChat(chatID, userID) {
		http.Error(w, "Unauthorized access to this chat", http.StatusForbidden)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50 // valor por defecto

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	messages, err := h.DirectChatService.GetDirectChatMessages(chatID, limit)
	if err != nil {
		http.Error(w, "Error getting messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

// GetAllRooms obtiene todas las salas ordenadas por updatedAt
//
//	@Summary		Obtiene todas las salas
//	@Description	Devuelve todas las salas ordenadas por fecha de actualización descendente
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.Room	"Lista de salas"
//	@Failure		401	{string}	string		"No autorizado"
//	@Failure		500	{string}	string		"Error interno del servidor"
//	@Router			/chat/rooms [get]
func (h *ChatHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.RoomService.GetAllRooms()
	if err != nil {
		http.Error(w, "Error getting rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}

// JoinRoom permite a un usuario unirse a una sala
//
//	@Summary		Unirse a una sala
//	@Description	Permite al usuario autenticado unirse a una sala específica
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			roomId	path		string	true	"ID de la sala"
//	@Success		200		{string}	string	"Usuario unido exitosamente"
//	@Failure		401		{string}	string	"No autorizado"
//	@Failure		403		{string}	string	"No permitido unirse a esta sala"
//	@Failure		404		{string}	string	"Sala no encontrada"
//	@Failure		409		{string}	string	"Usuario ya es miembro de la sala"
//	@Failure		500		{string}	string	"Error interno del servidor"
//	@Router			/chat/rooms/{roomId}/join [post]
func (h *ChatHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")

	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	err := h.RoomService.JoinRoom(roomID, userID)
	if err != nil {
		if err.Error() == "user is not allowed to join this room" {
			http.Error(w, "Error joining room: "+err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "user is already a member of the room" {
			http.Error(w, "Error joining room: "+err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Error joining room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Successfully joined the room",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetRoomMessagesSimple obtiene los mensajes de una sala sin paginación
//
//	@Summary		Obtiene mensajes de una sala (versión simple)
//	@Description	Devuelve los mensajes de una sala específica sin paginación
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			roomId	path		string					true	"ID de la sala"
//	@Param			limit	query		int						false	"Límite de mensajes a obtener"	default(50)
//	@Success		200		{array}		models.MessageResponse	"Lista de mensajes de la sala"
//	@Failure		401		{string}	string					"No autorizado"
//	@Failure		404		{string}	string					"Sala no encontrada"
//	@Failure		500		{string}	string					"Error interno del servidor"
//	@Router			/chat/rooms/{roomId}/messages [get]
func (h *ChatHandler) GetRoomMessagesSimple(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")

	limitStr := r.URL.Query().Get("limit")
	limit := 50 // valor por defecto

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	messages, err := h.RoomService.GetRoomMessagesSimple(roomID, limit)
	if err != nil {
		http.Error(w, "Error getting messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

// GetChat obtiene un chat directo por ID
//
//	@Summary		Obtiene un chat directo por ID
//	@Description	Devuelve los detalles de un chat directo específico
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			chatId	path		string				true	"ID del chat directo"
//	@Success		200		{object}	models.DirectChat	"Detalles del chat directo"
//	@Failure		401		{string}	string				"No autorizado"
//	@Failure		404		{string}	string				"Chat no encontrado"
//	@Failure		500		{string}	string				"Error interno del servidor"
//	@Router			/chat/direct/{chatId} [get]
func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chatId")

	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	chat, err := h.DirectChatService.GetDirectChatWithSenderName(chatID)
	if err != nil {
		http.Error(w, "Error getting chat: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar si el usuario tiene acceso a este chat
	hasAccess := false
	for _, id := range chat.UserIDs {
		if id == userID {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		http.Error(w, "Unauthorized access to this chat", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(chat)
}
