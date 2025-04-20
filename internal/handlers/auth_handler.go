package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Parchat/backend/internal/services"
)

// AuthHandler maneja las peticiones relacionadas con usuarios
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler crea una nueva instancia de AuthHandler
func NewAuthHandler(userService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: userService,
	}
}

// GetCurrentUser obtiene el usuario actual
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Obtener el usuario
	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error getting user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con los datos del usuario
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
