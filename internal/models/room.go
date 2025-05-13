package models

import "time"

// Room representa una sala de chat
type Room struct {
	ID          string    `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	Description string    `json:"description" firestore:"description"`
	OwnerID     string    `json:"ownerId" firestore:"ownerId"`
	IsPrivate   bool      `json:"isPrivate" firestore:"isPrivate"`
	Members     []string  `json:"members" firestore:"members"`
	Admins      []string  `json:"admins" firestore:"admins"`
	LastMessage *Message  `json:"lastMessage,omitempty" firestore:"lastMessage,omitempty"`
	ImageURL    string    `json:"imageUrl" firestore:"imageUrl"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`
	IsDeleted   bool      `json:"isDeleted" firestore:"isDeleted"`
}
