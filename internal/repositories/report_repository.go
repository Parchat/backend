package repositories

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/Parchat/backend/internal/config"
	"github.com/Parchat/backend/internal/models"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// ReportRepository handles database operations for message reports
type ReportRepository struct {
	FirestoreClient *config.FirestoreClient
}

// NewReportRepository creates a new instance of ReportRepository
func NewReportRepository(client *config.FirestoreClient) *ReportRepository {
	return &ReportRepository{
		FirestoreClient: client,
	}
}

// CreateReport saves a new report in Firestore
func (r *ReportRepository) CreateReport(report *models.Report) error {
	ctx := context.Background()

	// Generate ID if not provided
	if report.ID == "" {
		report.ID = uuid.New().String()
	}

	// Save the report in the reports collection
	_, err := r.FirestoreClient.Client.
		Collection("reports").Doc(report.ID).
		Set(ctx, report)

	if err != nil {
		return fmt.Errorf("error creating report: %v", err)
	}

	return nil
}

// GetReportCountForUserInRoom gets the number of reports for a user in a specific room
func (r *ReportRepository) GetReportCountForUserInRoom(roomID, userID string) (int, error) {
	ctx := context.Background()

	// Query reports for this user in this room
	query := r.FirestoreClient.Client.
		Collection("reports").
		Where("roomId", "==", roomID).
		Where("reportedId", "==", userID)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return 0, fmt.Errorf("error getting reports: %v", err)
	}

	return len(docs), nil
}

// GetReportedUsersInRoom gets all users who have been reported in a room
func (r *ReportRepository) GetReportedUsersInRoom(roomID string) (map[string]int, error) {
	ctx := context.Background()
	reportedUsers := make(map[string]int)

	// Query all reports for this room
	query := r.FirestoreClient.Client.
		Collection("reports").
		Where("roomId", "==", roomID)

	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating reports: %v", err)
		}

		var report models.Report
		if err := doc.DataTo(&report); err != nil {
			log.Printf("Error converting document to report: %v", err)
			continue
		}

		// Increment report count for this user
		reportedUsers[report.ReportedID]++
	}

	return reportedUsers, nil
}

// DeleteReportsForUserInRoom deletes all reports for a user in a specific room
func (r *ReportRepository) DeleteReportsForUserInRoom(roomID, userID string) error {
	ctx := context.Background()

	// Query reports for this user in this room
	query := r.FirestoreClient.Client.
		Collection("reports").
		Where("roomId", "==", roomID).
		Where("reportedId", "==", userID)

	// Get all matching documents
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("error getting reports to delete: %v", err)
	}

	// Use a transaction instead of batch (as batch is deprecated)
	return r.FirestoreClient.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Delete each document in the transaction
		for _, doc := range docs {
			if err := tx.Delete(doc.Ref); err != nil {
				return fmt.Errorf("error adding delete operation to transaction: %v", err)
			}
		}
		return nil
	})
}

// HasUserReportedMessage checks if a user has already reported a specific message
func (r *ReportRepository) HasUserReportedMessage(reporterID, messageID string) (bool, error) {
	ctx := context.Background()

	// Query reports for this reporter and message
	query := r.FirestoreClient.Client.
		Collection("reports").
		Where("reporterId", "==", reporterID).
		Where("messageId", "==", messageID)

	// Check if any matching documents exist
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("error checking for existing report: %v", err)
	}

	return len(docs) > 0, nil
}

// UpdateRoomReportedUsers updates the reported users map in a room
func (r *ReportRepository) UpdateRoomReportedUsers(roomID string, reportedUsers map[string]int) error {
	ctx := context.Background()

	// Update the reportedUsers field in the room document
	_, err := r.FirestoreClient.Client.
		Collection("rooms").Doc(roomID).
		Update(ctx, []firestore.Update{
			{Path: "reportedUsers", Value: reportedUsers},
		})

	if err != nil {
		return fmt.Errorf("error updating room reported users: %v", err)
	}

	return nil
}
