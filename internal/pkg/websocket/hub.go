package websocket

import (
	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/repositories"
)

// BroadcastMessage contiene la información para transmitir un mensaje
type BroadcastMessage struct {
	Message    WebSocketMessage
	RoomID     string // ID de la sala si es un mensaje de sala
	DirectChat string // ID del chat directo si es un mensaje directo
}

// Hub mantiene el conjunto de clientes activos y transmite mensajes a los clientes
type Hub struct {
	// Clientes registrados
	clients map[*Client]bool

	// Registra a un nuevo cliente
	Register chan *Client

	// Desregistra a un cliente
	Unregister chan *Client

	// Canal para transmitir mensajes a las salas de chat
	Broadcast chan BroadcastMessage

	// Canal para transmitir mensajes a chats directos
	BroadcastDirect chan BroadcastMessage

	// Repositorios
	messageRepo    *repositories.MessageRepository
	roomRepo       *repositories.RoomRepository
	directChatRepo *repositories.DirectChatRepository
	reportRepo     *repositories.ReportRepository

	// Firestore client
	firestoreClient *config.FirestoreClient
}

// NewHub inicializa un nuevo Hub
func NewHub(
	messageRepo *repositories.MessageRepository,
	roomRepo *repositories.RoomRepository,
	directChatRepo *repositories.DirectChatRepository,
	reportRepo *repositories.ReportRepository,
	client *config.FirestoreClient,
) *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Broadcast:       make(chan BroadcastMessage),
		BroadcastDirect: make(chan BroadcastMessage),
		messageRepo:     messageRepo,
		roomRepo:        roomRepo,
		directChatRepo:  directChatRepo,
		reportRepo:      reportRepo,
		firestoreClient: client,
	}
}

// Run comienza el hub, gestionando las conexiones de los clientes y los mensajes
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			// Difundir a todos los clientes que están en la sala
			for client := range h.clients {
				if client.IsInRoom(message.RoomID) {
					select {
					case client.send <- message.Message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		case message := <-h.BroadcastDirect:
			// Difundir a todos los clientes que están en el chat directo
			for client := range h.clients {
				if client.IsInDirectChat(message.DirectChat) {
					select {
					case client.send <- message.Message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
