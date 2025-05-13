package handlers

import (
	"log"
	"net/http"

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
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Verificar el token desde el contexto (añadido por el middleware de autenticación)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
