package handlers

import (
	"github.com/Parchat/backend/internal/services"
)

type UserHandler struct {
	UserService *services.UserService
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}
