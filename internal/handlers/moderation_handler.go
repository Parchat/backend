package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

// ModerationHandler handles the HTTP requests related to content moderation
type ModerationHandler struct {
	moderationService *services.ModerationService
	roomService       *services.RoomService
}

// NewModerationHandler creates a new instance of ModerationHandler
func NewModerationHandler(
	moderationService *services.ModerationService,
	roomService *services.RoomService,
) *ModerationHandler {
	return &ModerationHandler{
		moderationService: moderationService,
		roomService:       roomService,
	}
}

// ReportMessage handles the request to report an inappropriate message
//
//	@Summary		Report an inappropriate message
//	@Description	Reports a message as inappropriate in a chat room
//	@Tags			Moderation
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			roomId	path		string					true	"Room ID"
//	@Param			report	body		models.ReportRequest	true	"Report Request"
//	@Success		200		{string}	string					"Message reported successfully"
//	@Failure		400		{string}	string					"Invalid request"
//	@Failure		401		{string}	string					"Unauthorized"
//	@Failure		403		{string}	string					"Forbidden"
//	@Failure		404		{string}	string					"Not found"
//	@Failure		500		{string}	string					"Internal server error"
//	@Router			/chat/rooms/{roomId}/report [post]
func (h *ModerationHandler) ReportMessage(w http.ResponseWriter, r *http.Request) {
	// Get the room ID from the URL
	roomID := chi.URLParam(r, "roomId")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Parse the request body
	var reportReq models.ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&reportReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if reportReq.MessageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	// Call the service to report the message
	err := h.moderationService.ReportMessage(userID, roomID, reportReq.MessageID, reportReq.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message reported successfully"))
}

// GetBannedUsers handles the request to get all banned users in a room
//
//	@Summary		Get banned users in a room
//	@Description	Retrieves a list of users who have been banned in a chat room
//	@Tags			Moderation
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			roomId	path		string						true	"Room ID"
//	@Success		200		{object}	models.BannedUsersResponse	"Banned users"
//	@Failure		401		{string}	string						"Unauthorized"
//	@Failure		403		{string}	string						"Forbidden"
//	@Failure		404		{string}	string						"Not found"
//	@Failure		500		{string}	string						"Internal server error"
//	@Router			/chat/rooms/{roomId}/banned-users [get]
func (h *ModerationHandler) GetBannedUsers(w http.ResponseWriter, r *http.Request) {
	// Get the room ID from the URL
	roomID := chi.URLParam(r, "roomId")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Check if the user is an admin or owner of the room
	isAuthorized, err := h.roomService.IsUserAdminOrOwner(roomID, userID)
	if err != nil {
		http.Error(w, "Error checking user authorization", http.StatusInternalServerError)
		return
	}

	if !isAuthorized {
		http.Error(w, "Unauthorized: Only room admins or owner can view reported users", http.StatusForbidden)
		return
	}
	// Get banned users in the room
	bannedUsers, err := h.moderationService.GetBannedUsersInRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bannedUsers)
}

// ClearUserReports handles the request to clear all reports for a user in a room
//
//	@Summary		Clear reports for a user
//	@Description	Clears all reports for a specific user in a chat room
//	@Tags			Moderation
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			roomId			path		string						true	"Room ID"
//	@Param			clearRequest	body		models.ClearReportRequest	true	"Clear Report Request"
//	@Success		200				{string}	string						"Reports cleared successfully"
//	@Failure		400				{string}	string						"Invalid request"
//	@Failure		401				{string}	string						"Unauthorized"
//	@Failure		403				{string}	string						"Forbidden"
//	@Failure		404				{string}	string						"Not found"
//	@Failure		500				{string}	string						"Internal server error"
//	@Router			/chat/rooms/{roomId}/clear-reports [post]
func (h *ModerationHandler) ClearUserReports(w http.ResponseWriter, r *http.Request) {
	// Get the room ID from the URL
	roomID := chi.URLParam(r, "roomId")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	// Obtener el ID del usuario del contexto
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Check if the user is an admin or owner of the room
	isAuthorized, err := h.roomService.IsUserAdminOrOwner(roomID, userID)
	if err != nil {
		http.Error(w, "Error checking user authorization", http.StatusInternalServerError)
		return
	}

	if !isAuthorized {
		http.Error(w, "Unauthorized: Only room admins or owner can clear reports", http.StatusForbidden)
		return
	}

	// Parse the request body
	var clearReq models.ClearReportRequest
	if err := json.NewDecoder(r.Body).Decode(&clearReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if clearReq.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Call the service to clear the reports
	err = h.moderationService.ClearReportsForUser(roomID, clearReq.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reports cleared successfully"))
}
