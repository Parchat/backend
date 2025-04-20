package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Parchat/backend/internal/models"
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

// SignUpAndCreateUser maneja el registro de un nuevo usuario
func (h *AuthHandler) SignUpAndCreateUser(w http.ResponseWriter, r *http.Request) {
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

	if err := h.authService.SignUpAndCreateUser(payload.Password, user); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
