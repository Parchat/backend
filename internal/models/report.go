package models

import "time"

// Report represents a user report for an inappropriate message
type Report struct {
	ID         string    `json:"id" firestore:"id"`
	MessageID  string    `json:"messageId" firestore:"messageId"`
	RoomID     string    `json:"roomId" firestore:"roomId"`
	ReportedID string    `json:"reportedId" firestore:"reportedId"` // ID of the user being reported
	ReporterID string    `json:"reporterId" firestore:"reporterId"` // ID of the user making the report
	Reason     string    `json:"reason" firestore:"reason"`
	CreatedAt  time.Time `json:"createdAt" firestore:"createdAt"`
}

// ReportRequest represents the request to report a message
type ReportRequest struct {
	MessageID string `json:"messageId"`
	Reason    string `json:"reason,omitempty"`
}

// BannedUserResponse represents a user who has been banned due to reports
type BannedUserResponse struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	ReportCount int    `json:"reportCount"`
}

// BannedUsersResponse represents a list of banned users in a room
type BannedUsersResponse struct {
	Users []BannedUserResponse `json:"users"`
}

// ClearReportRequest represents the request to clear reports for a user in a room
type ClearReportRequest struct {
	UserID string `json:"userId"`
}
