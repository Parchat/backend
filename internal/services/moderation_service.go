package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Parchat/backend/internal/models"
	"github.com/Parchat/backend/internal/repositories"
	"github.com/google/uuid"
)

const (
	// MaxReportsBeforeBan is the threshold of reports a user can receive before being banned from a room
	MaxReportsBeforeBan = 3
)

// ModerationService handles operations related to content moderation and user reports
type ModerationService struct {
	reportRepo  *repositories.ReportRepository
	messageRepo *repositories.MessageRepository
	roomRepo    *repositories.RoomRepository
	userRepo    *repositories.UserRepository
}

// NewModerationService creates a new instance of ModerationService
func NewModerationService(
	reportRepo *repositories.ReportRepository,
	messageRepo *repositories.MessageRepository,
	roomRepo *repositories.RoomRepository,
	userRepo *repositories.UserRepository,
) *ModerationService {
	return &ModerationService{
		reportRepo:  reportRepo,
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
	}
}

// ReportMessage handles the reporting of an inappropriate message
func (s *ModerationService) ReportMessage(reporterID, roomID, messageID, reason string) error {
	// Validate that the message exists
	message, err := s.messageRepo.GetMessageByID(roomID, messageID)
	if err != nil {
		return fmt.Errorf("message not found: %v", err)
	}

	// Don't allow users to report their own messages
	if message.UserID == reporterID {
		return fmt.Errorf("users cannot report their own messages")
	}

	// Create a report record
	report := &models.Report{
		ID:         uuid.New().String(),
		MessageID:  messageID,
		RoomID:     roomID,
		ReportedID: message.UserID, // The user who sent the message
		ReporterID: reporterID,     // The user making the report
		Reason:     reason,
		CreatedAt:  time.Now(),
	}

	// Save the report
	if err := s.reportRepo.CreateReport(report); err != nil {
		return fmt.Errorf("failed to create report: %v", err)
	}

	// Get current reported users for the room
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		return fmt.Errorf("room not found: %v", err)
	}

	// Initialize reportedUsers map if it doesn't exist
	if room.ReportedUsers == nil {
		room.ReportedUsers = make(map[string]int)
	}

	// Increment report count for the user
	room.ReportedUsers[message.UserID]++

	// Update the room with the new reported users
	if err := s.reportRepo.UpdateRoomReportedUsers(roomID, room.ReportedUsers); err != nil {
		return fmt.Errorf("failed to update room reported users: %v", err)
	}

	return nil
}

// GetBannedUsersInRoom retrieves all users who have been banned in a room
func (s *ModerationService) GetBannedUsersInRoom(roomID string) (*models.BannedUsersResponse, error) {
	// Get the room to check if the room exists and to get the reported users
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		return nil, fmt.Errorf("room not found: %v", err)
	}

	// Initialize response
	response := &models.BannedUsersResponse{
		Users: []models.BannedUserResponse{},
	}

	// If no reported users, return empty list
	if len(room.ReportedUsers) == 0 {
		return response, nil
	}

	// Fetch user details for each reported user
	for userID, reportCount := range room.ReportedUsers {
		// Only include users who have reached or exceeded the threshold
		if reportCount >= MaxReportsBeforeBan {
			// Get user details
			ctx := context.Background()
			user, err := s.userRepo.GetUserByID(ctx, userID)
			if err != nil {
				// Skip this user if we can't get their details
				continue
			} // Add to response
			response.Users = append(response.Users, models.BannedUserResponse{
				UserID:      userID,
				DisplayName: user.DisplayName,
				ReportCount: reportCount,
			})
		}
	}

	return response, nil
}

// ClearReportsForUser clears all reports for a specific user in a room
func (s *ModerationService) ClearReportsForUser(roomID, userID string) error {
	// Delete the reports from the reports collection
	if err := s.reportRepo.DeleteReportsForUserInRoom(roomID, userID); err != nil {
		return fmt.Errorf("failed to delete reports: %v", err)
	}

	// Get the room to update the reported users map
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		return fmt.Errorf("room not found: %v", err)
	}

	// If no reported users, nothing to do
	if room.ReportedUsers == nil {
		return nil
	}

	// Remove the user from the reported users map
	delete(room.ReportedUsers, userID)

	// Update the room with the new reported users
	if err := s.reportRepo.UpdateRoomReportedUsers(roomID, room.ReportedUsers); err != nil {
		return fmt.Errorf("failed to update room reported users: %v", err)
	}

	return nil
}

// CanUserSendMessageInRoom checks if a user can send messages in a room based on report count
func (s *ModerationService) CanUserSendMessageInRoom(roomID, userID string) bool {
	// Get the room
	room, err := s.roomRepo.GetRoom(roomID)
	if err != nil {
		// If we can't get the room, default to allowing the message
		return true
	}

	// If no reported users or user not in reported users, they can send messages
	if room.ReportedUsers == nil {
		return true
	}

	reportCount, exists := room.ReportedUsers[userID]

	// If the user has fewer reports than the threshold or isn't reported, they can send messages
	return !exists || reportCount < MaxReportsBeforeBan
}
