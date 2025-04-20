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

// CreateUser maneja la creación de un nuevo usuario
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Crear el usuario
	user := &models.User{
		Email:       payload.Email,
		DisplayName: payload.DisplayName,
	}

	if err := h.UserService.CreateUserWithAuth(payload.Password, user); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
