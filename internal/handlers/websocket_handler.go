package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/Parchat/backend/internal/config"
	pws "github.com/Parchat/backend/internal/pkg/websocket"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir conexiones de cualquier origen en desarrollo
	},
}

// WebSocketHandler maneja las conexiones WebSocket
type WebSocketHandler struct {
	hub          *pws.Hub
	firebaseAuth *config.FirebaseAuth
}

// NewWebSocketHandler crea una nueva instancia de WebSocketHandler
func NewWebSocketHandler(hub *pws.Hub, firebaseAuth *config.FirebaseAuth) *WebSocketHandler {
	return &WebSocketHandler{
		hub:          hub,
		firebaseAuth: firebaseAuth,
	}
}

// HandleWebSocket maneja las conexiones WebSocket
//
//	@Summary		Conexión WebSocket para chat en tiempo real
//	@Description	Establece una conexión WebSocket para mensajería en tiempo real
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			token	query		string	true	"Firebase Auth Token"
//	@Success		101		{string}	string	"Switching Protocols a WebSocket"
//	@Failure		401		{string}	string	"No autorizado"
//	@Router			/chat/ws [get]
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extraer el token desde el parámetro de consulta
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token not provided", http.StatusUnauthorized)
		return
	}

	// Limpiar el token (por si viene con "Bearer ")
	token = strings.TrimPrefix(token, "Bearer ")

	// Verificar el token con Firebase
	authToken, err := h.firebaseAuth.VerifyIDToken(r.Context(), token)
	if err != nil {
		log.Printf("Error verifying token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extraer el userID del token verificado
	userID := authToken.UID
	if userID == "" {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	// Actualizar la conexión a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Crear un nuevo cliente
	client := pws.NewClient(h.hub, conn, userID)

	// Registrar cliente con el hub
	h.hub.Register <- client

	// Permitir recolección de basura del WebSocket cuando complete
	go client.WritePump()
	go client.ReadPump()
}
