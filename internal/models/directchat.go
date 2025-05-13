package models

import "time"

// DirectChat representa un chat directo entre dos usuarios
type DirectChat struct {
	ID          string    `json:"id" firestore:"id"`
	UserIDs     []string  `json:"userIds" firestore:"userIds"`
	LastMessage *Message  `json:"lastMessage,omitempty" firestore:"lastMessage,omitempty"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`
	IsDeleted   bool      `json:"isDeleted" firestore:"isDeleted"`
}
