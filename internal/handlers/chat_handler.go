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

// GetRoomMessages obtiene los mensajes de una sala
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

	messages, err := h.RoomService.GetRoomMessages(roomID, limit)
	if err != nil {
		http.Error(w, "Error getting messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

// CreateDirectChat crea o encuentra un chat directo entre dos usuarios
func (h *ChatHandler) CreateDirectChat(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Obtener el ID del otro usuario
	var requestData struct {
		OtherUserID string `json:"otherUserId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if requestData.OtherUserID == "" {
		http.Error(w, "Other user ID is required", http.StatusBadRequest)
		return
	}

	chat, err := h.DirectChatService.FindOrCreateDirectChat(userID, requestData.OtherUserID)
	if err != nil {
		http.Error(w, "Error creating direct chat: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

// GetUserDirectChats obtiene todos los chats directos del usuario
func (h *ChatHandler) GetUserDirectChats(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	chats, err := h.DirectChatService.GetUserDirectChats(userID)
	if err != nil {
		http.Error(w, "Error getting direct chats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chats)
}

// GetDirectChatMessages obtiene los mensajes de un chat directo
func (h *ChatHandler) GetDirectChatMessages(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chatId")

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
