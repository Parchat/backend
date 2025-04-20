package models

import "time"

// User representa un usuario en la aplicación
type User struct {
	UID         string `json:"uid" firestore:"uid"`
	Email       string `json:"email" firestore:"email"`
	DisplayName string `json:"displayName" firestore:"displayName"`
	PhotoURL    string `json:"photoUrl" firestore:"photoUrl"`
	Status      string `json:"status" firestore:"status"`
	LastSeen    string `json:"lastSeen" firestore:"lastSeen"`
	//BlockedUsers []string  `json:"blockedUsers" firestore:"blockedUsers"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
	IsDeleted bool      `json:"isDeleted" firestore:"isDeleted"`
}
