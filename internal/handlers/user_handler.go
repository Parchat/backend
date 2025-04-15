package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Parchat/backend/internal/services"
)

// UserHandler maneja las peticiones relacionadas con usuarios
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// AuthStatus devuelve el estado de autenticación
func (h *UserHandler) AuthStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "Authentication service is running",
	})
}

// GetCurrentUser obtiene el usuario actual
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Obtener el usuario
	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error getting user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con los datos del usuario
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
