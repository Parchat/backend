package models

// User representa un usuario en la aplicación
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	PhotoURL  string `json:"photo_url,omitempty"`
	CreatedAt int64  `json:"created_at"`
}
