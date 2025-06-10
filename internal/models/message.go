package models

import "time"

// Message representa un mensaje enviado por un usuario
type Message struct {
	ID        string    `json:"id" firestore:"id"`
	Content   string    `json:"content" firestore:"content"`
	UserID    string    `json:"userId" firestore:"userId"`
	RoomID    string    `json:"roomId" firestore:"roomId"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
	IsDeleted bool      `json:"isDeleted" firestore:"isDeleted"`
}

// MessageResponse es la respuesta que incluye un mensaje y el nombre del usuario que lo envió
type MessageResponse struct {
	Message
	DisplayName string `json:"displayName,omitempty"`
}

// PaginatedMessagesResponse representa una respuesta paginada de mensajes
type PaginatedMessagesResponse struct {
	Messages   []MessageResponse `json:"messages"`
	NextCursor string            `json:"nextCursor,omitempty"`
	HasMore    bool              `json:"hasMore"`
}
