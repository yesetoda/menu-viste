package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"menuvista/internal/models"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries *persistence.Queries
}

func NewService(queries *persistence.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) LogActivity(ctx context.Context, input models.CreateActivityLogRequest) error {
	// Async logging is often better, but for simplicity we'll do sync for now or fire-and-forget in handler.
	// Here we just implement the core logic.

	beforeJSON, _ := json.Marshal(input.BeforeValue)
	afterJSON, _ := json.Marshal(input.AfterValue)

	_, err := s.queries.CreateActivityLog(ctx, persistence.CreateActivityLogParams{
		RestaurantID:   input.RestaurantID,
		UserID:         input.UserID,
		ActionType:     input.ActionType,
		ActionCategory: input.ActionCategory,
		Description:    pgtype.Text{String: input.Description, Valid: input.Description != ""},
		TargetType:     pgtype.Text{String: input.TargetType, Valid: input.TargetType != ""},
		TargetID:       *utils.ToUUIDPtr(input.TargetID),
		TargetName:     pgtype.Text{String: input.TargetName, Valid: input.TargetName != ""},
		BeforeValue:    beforeJSON,
		AfterValue:     afterJSON,
		IpAddress:      pgtype.Text{String: input.IPAddress, Valid: input.IPAddress != ""},
		UserAgent:      pgtype.Text{String: input.UserAgent, Valid: input.UserAgent != ""},
		DeviceType:     pgtype.Text{String: input.DeviceType, Valid: input.DeviceType != ""},
		Browser:        pgtype.Text{String: input.Browser, Valid: input.Browser != ""},
		Os:             pgtype.Text{String: input.OS, Valid: input.OS != ""},
		Success:        pgtype.Bool{Bool: input.Success, Valid: true},
	})

	if err != nil {
		log.Printf("[ActivityService] Failed to create activity log: %v", err)
		return fmt.Errorf("failed to log activity: %w", err)
	}

	return nil
}

func (s *Service) ListLogs(ctx context.Context, restaurantID uuid.UUID, limit, offset int32) ([]*models.ActivityLog, error) {
	rows, err := s.queries.ListActivityLogsByRestaurant(ctx, persistence.ListActivityLogsByRestaurantParams{
		RestaurantID: restaurantID,
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list activity logs: %w", err)
	}

	logs := make([]*models.ActivityLog, len(rows))
	for i, row := range rows {
		logs[i] = s.mapToDomainLog(row)
	}
	return logs, nil
}

func (s *Service) mapToDomainLog(row persistence.ListActivityLogsByRestaurantRow) *models.ActivityLog {
	id := row.ID
	restaurantID := row.RestaurantID
	userID := row.UserID

	var beforeVal, afterVal map[string]interface{}
	if len(row.BeforeValue) > 0 {
		json.Unmarshal(row.BeforeValue, &beforeVal)
	}
	if len(row.AfterValue) > 0 {
		json.Unmarshal(row.AfterValue, &afterVal)
	}

	return &models.ActivityLog{
		ID:             id,
		RestaurantID:   restaurantID,
		UserID:         userID,
		UserName:       row.UserName,
		UserEmail:      row.UserEmail,
		ActionType:     row.ActionType,
		ActionCategory: row.ActionCategory,
		Description:    row.Description.String,
		TargetType:     row.TargetType.String,
		TargetID:       &row.TargetID,
		TargetName:     row.TargetName.String,
		BeforeValue:    json.RawMessage(row.BeforeValue),
		AfterValue:     json.RawMessage(row.AfterValue),
		IPAddress:      row.IpAddress.String,
		UserAgent:      row.UserAgent.String,
		DeviceType:     row.DeviceType.String,
		Browser:        row.Browser.String,
		Os:             row.Os.String,
		Success:        row.Success.Bool,
		CreatedAt:      row.CreatedAt.Time,
	}
}
