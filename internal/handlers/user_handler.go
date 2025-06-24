package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/services"
)

type UserHandler struct {
	UserService *services.UserService
	AuthService *services.AuthService
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(userService *services.UserService, authService *services.AuthService) *UserHandler {
	return &UserHandler{
		UserService: userService,
		AuthService: authService,
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

// EnsureUserExists verifica si el usuario existe en la base de datos, si no, lo crea con los datos de autenticación
//
//	@Summary		Asegura que el usuario exista en la base de datos
//	@Description	Verifica si el usuario autenticado existe en la base de datos, si no, lo crea con los datos de autenticación
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.User	"Datos del usuario"
//	@Failure		500	{string}	string		"Error interno del servidor"
//	@Router			/user/create [post]
func (h *UserHandler) EnsureUserExists(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Obtener los datos del usuario desde Firebase Auth
	authUser, err := h.AuthService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error getting user from auth: "+err.Error(), http.StatusInternalServerError)
		return
	}
	println("Auth User:", authUser.UID, authUser.Email, authUser.DisplayName)

	// Asegurar que el usuario exista en la base de datos
	user, err := h.UserService.EnsureUserExists(r.Context(), authUser)
	if err != nil {
		http.Error(w, "Error ensuring user exists: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Responder con los datos del usuario
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
