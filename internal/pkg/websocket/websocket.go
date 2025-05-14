package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Parchat/backend/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Tiempo permitido para escribir un mensaje al cliente
	writeWait = 10 * time.Second

	// Tiempo permitido para leer el próximo mensaje del cliente
	readWait = 60 * time.Second

	// Enviar pings al cliente con esta periodicidad
	pingPeriod = (readWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

// Tipo de mensaje para diferenciar entre tipos de eventos
type MessageType string

const (
	MessageTypeChatRoom       MessageType = "CHAT_ROOM"
	MessageTypeDirectChat     MessageType = "DIRECT_CHAT"
	MessageTypeJoinRoom       MessageType = "JOIN_ROOM"
	MessageTypeJoinDirectChat MessageType = "JOIN_DIRECT_CHAT"
	MessageTypeUserLeave      MessageType = "USER_LEAVE"
	MessageTypeError          MessageType = "ERROR"
	MessageTypeRoomCreated    MessageType = "ROOM_CREATED"
)

// WebSocketMessage representa el formato de mensaje que se intercambia entre cliente y servidor
type WebSocketMessage struct {
	Type      MessageType     `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

// Client representa un cliente de WebSocket
type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan WebSocketMessage
	userID     string
	rooms      map[string]bool // RoomIDs que el cliente está escuchando
	directChat map[string]bool // DirectChatIDs que el cliente está escuchando
}

// NewClient crea un nuevo cliente
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan WebSocketMessage, 256),
		userID:     userID,
		rooms:      make(map[string]bool),
		directChat: make(map[string]bool),
	}
}

// ReadPump bombea mensajes desde el WebSocket al hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(readWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(readWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Manejar diferentes tipos de mensajes
		switch wsMessage.Type {
		case MessageTypeChatRoom:
			var chatMsg models.Message
			if err := json.Unmarshal(wsMessage.Payload, &chatMsg); err != nil {
				log.Printf("Error unmarshaling chat message: %v", err)
				continue
			}

			// Verificar si el usuario es parte de la sala antes de enviar el mensaje
			if !c.hub.roomRepo.CanTalkInRoomWebSocket(chatMsg.RoomID, c.userID) {
				errMsg := "No permission to send messages to this room"
				errorPayload, _ := json.Marshal(errMsg)
				c.send <- WebSocketMessage{
					Type:      MessageTypeError,
					Payload:   errorPayload,
					Timestamp: time.Now(),
				}
				log.Printf("User %s attempted to send message to room %s without permission", c.userID, chatMsg.RoomID)
				continue
			}

			// Asignar ID y timestamps si no existen
			if chatMsg.ID == "" {
				chatMsg.ID = uuid.New().String()
			}
			now := time.Now()
			if chatMsg.CreatedAt.IsZero() {
				chatMsg.CreatedAt = now
			}
			chatMsg.UpdatedAt = now

			// Asegurarse que el userID es el correcto
			chatMsg.UserID = c.userID

			// Guardar el mensaje en Firestore
			err = c.hub.messageRepo.SaveMessage(&chatMsg)
			if err != nil {
				log.Printf("Error saving message: %v", err)
				continue
			}

			// Actualizar el último mensaje en la sala
			err = c.hub.roomRepo.UpdateLastMessage(chatMsg.RoomID, &chatMsg)
			if err != nil {
				log.Printf("Error updating last message: %v", err)
			}

			// Convertir el mensaje de vuelta a JSON para difundir
			payload, _ := json.Marshal(chatMsg)
			wsMessage.Payload = payload
			wsMessage.Timestamp = time.Now()

			// Transmitir mensaje a todos en la sala
			c.hub.Broadcast <- BroadcastMessage{
				Message: wsMessage,
				RoomID:  chatMsg.RoomID,
			}

		case MessageTypeDirectChat:
			var chatMsg models.Message
			if err := json.Unmarshal(wsMessage.Payload, &chatMsg); err != nil {
				log.Printf("Error unmarshaling direct chat message: %v", err)
				continue
			}

			// Verificar si el usuario es parte del chat directo antes de enviar el mensaje
			if !c.hub.directChatRepo.IsUserInDirectChat(chatMsg.RoomID, c.userID) {
				errMsg := "Not a member of this direct chat"
				errorPayload, _ := json.Marshal(errMsg)
				c.send <- WebSocketMessage{
					Type:      MessageTypeError,
					Payload:   errorPayload,
					Timestamp: time.Now(),
				}
				log.Printf("User %s attempted to send message to direct chat %s without being a member", c.userID, chatMsg.RoomID)
				continue
			}

			// Asignar ID y timestamps si no existen
			if chatMsg.ID == "" {
				chatMsg.ID = uuid.New().String()
			}
			now := time.Now()
			if chatMsg.CreatedAt.IsZero() {
				chatMsg.CreatedAt = now
			}
			chatMsg.UpdatedAt = now

			// Asegurarse que el userID es el correcto
			chatMsg.UserID = c.userID

			// Guardar el mensaje en Firestore
			err = c.hub.messageRepo.SaveDirectMessage(&chatMsg)
			if err != nil {
				log.Printf("Error saving direct message: %v", err)
				continue
			}

			// Actualizar el último mensaje en el chat directo
			err = c.hub.directChatRepo.UpdateLastMessage(chatMsg.RoomID, &chatMsg)
			if err != nil {
				log.Printf("Error updating last message in direct chat: %v", err)
			}

			// Convertir el mensaje de vuelta a JSON para difundir
			payload, _ := json.Marshal(chatMsg)
			wsMessage.Payload = payload
			wsMessage.Timestamp = time.Now()

			// Transmitir mensaje a todos en el chat directo
			c.hub.BroadcastDirect <- BroadcastMessage{
				Message:    wsMessage,
				DirectChat: chatMsg.RoomID,
			}

		case MessageTypeJoinRoom:
			var roomID string
			if err := json.Unmarshal(wsMessage.Payload, &roomID); err != nil {
				log.Printf("Error unmarshaling room ID: %v", err)
				continue
			}

			// Verificar si el usuario tiene permiso para unirse a la sala
			if c.hub.roomRepo.CanJoinRoomWebSocket(roomID, c.userID) {
				c.rooms[roomID] = true
				log.Printf("User %s joined room %s", c.userID, roomID)
			} else {
				errMsg := "No permission to join this room"
				errorPayload, _ := json.Marshal(errMsg)
				c.send <- WebSocketMessage{
					Type:      MessageTypeError,
					Payload:   errorPayload,
					Timestamp: time.Now(),
				}
			}

		// Manejar cuando un usuario quiere escuchar un chat directo
		case MessageTypeJoinDirectChat:
			var directChatID string
			if err := json.Unmarshal(wsMessage.Payload, &directChatID); err != nil {
				log.Printf("Error unmarshaling direct chat ID: %v", err)
				continue
			}

			// Verificar si el usuario es parte del chat directo
			if c.hub.directChatRepo.IsUserInDirectChat(directChatID, c.userID) {
				c.directChat[directChatID] = true
				log.Printf("User %s joined direct chat %s", c.userID, directChatID)
			} else {
				errMsg := "Not a member of this direct chat"
				errorPayload, _ := json.Marshal(errMsg)
				c.send <- WebSocketMessage{
					Type:      MessageTypeError,
					Payload:   errorPayload,
					Timestamp: time.Now(),
				}
			}
		}
	}
}

// WritePump bombea mensajes desde el hub al WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// El hub cerró el canal
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Serializar el mensaje WebSocketMessage a JSON
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			w.Write(messageBytes)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// IsInRoom comprueba si el cliente está en una sala específica
func (c *Client) IsInRoom(roomID string) bool {
	_, ok := c.rooms[roomID]
	return ok
}

// IsInDirectChat comprueba si el cliente está en un chat directo específico
func (c *Client) IsInDirectChat(directChatID string) bool {
	_, ok := c.directChat[directChatID]
	return ok
}
