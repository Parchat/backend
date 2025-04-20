package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Parchat/backend/internal/models"
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

// CreateUser crea un nuevo usuario
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
