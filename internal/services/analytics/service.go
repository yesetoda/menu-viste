package analytics

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (s *Service) TrackEvent(ctx context.Context, input models.CreateAnalyticsEventRequest) error {
	_, err := s.queries.CreateAnalyticsEvent(ctx, persistence.CreateAnalyticsEventParams{
		RestaurantID: input.RestaurantID,
		EventType:    input.EventType,
		VisitorID:    input.VisitorID,
		SessionID:    uuid.Nil, // Optional, can be added to request model if needed
		TargetID:     *utils.ToUUIDPtr(input.TargetID),
		IpAddress:    pgtype.Text{String: input.IPAddress, Valid: input.IPAddress != ""},
		DeviceType:   pgtype.Text{String: input.DeviceType, Valid: input.DeviceType != ""},
		Browser:      pgtype.Text{String: input.Browser, Valid: input.Browser != ""},
		Os:           pgtype.Text{String: input.OS, Valid: input.OS != ""},
		Country:      pgtype.Text{String: input.Country, Valid: input.Country != ""},
		City:         pgtype.Text{String: input.City, Valid: input.City != ""},
	})
	if err != nil {
		log.Printf("[AnalyticsService] Failed to track event: %v", err)
		return fmt.Errorf("failed to track event: %w", err)
	}

	// Update aggregates (simplified for now, ideally async)
	// For example, if event is "view_menu", increment view count
	if input.EventType == "view_restaurant" {
		if err := s.queries.IncrementRestaurantViewCount(ctx, input.RestaurantID); err != nil {
			log.Printf("[AnalyticsService] Failed to increment restaurant view count: %v", err)
		}
	} else if input.EventType == "view_item" && (input.TargetID != uuid.Nil) {
		if err := s.queries.IncrementMenuItemViewCount(ctx, input.TargetID); err != nil {
			log.Printf("[AnalyticsService] Failed to increment menu item view count: %v", err)
		}
	}

	// Also update daily aggregates
	now := time.Now()
	hour := int32(now.Hour())
	date := pgtype.Date{Time: now, Valid: true}

	_, err = s.queries.UpsertAnalyticsAggregate(ctx, persistence.UpsertAnalyticsAggregateParams{
		RestaurantID: input.RestaurantID,
		Date:         date,
		Hour:         pgtype.Int4{Int32: hour, Valid: true},
		MetricType:   input.EventType,
		TargetID:     *utils.ToUUIDPtr(input.TargetID),
		Value:        pgtype.Int4{Int32: 1, Valid: true},
	})
	if err != nil {
		log.Printf("[AnalyticsService] Failed to upsert aggregate: %v", err)
	}

	return nil
}

func (s *Service) GetAggregates(ctx context.Context, restaurantID uuid.UUID, startDate, endDate time.Time) ([]*models.AnalyticsAggregate, error) {
	rows, err := s.queries.GetAnalyticsAggregates(ctx, persistence.GetAnalyticsAggregatesParams{
		RestaurantID: restaurantID,
		Date:         pgtype.Date{Time: startDate, Valid: true},
		Date_2:       pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregates: %w", err)
	}

	aggregates := make([]*models.AnalyticsAggregate, len(rows))
	for i, row := range rows {
		aggregates[i] = s.mapToDomainAggregate(row)
	}
	return aggregates, nil
}

func (s *Service) GetOverview(ctx context.Context, restaurantID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	// This is a simplified overview. In a real app, this would be more complex queries.
	// For now, let's just return total views and some basic stats based on aggregates.

	aggregates, err := s.GetAggregates(ctx, restaurantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var totalViews int32
	var totalMenuViews int32

	for _, agg := range aggregates {
		if agg.MetricType == "view_restaurant" {
			totalViews += agg.Value
		} else if agg.MetricType == "view_item" {
			totalMenuViews += agg.Value
		}
	}

	return map[string]interface{}{
		"total_views":      totalViews,
		"total_menu_views": totalMenuViews,
		"period_start":     startDate,
		"period_end":       endDate,
	}, nil
}

func (s *Service) mapToDomainAggregate(row persistence.AnalyticsAggregate) *models.AnalyticsAggregate {
	id := row.ID
	restaurantID := row.RestaurantID

	var hour *int32
	if row.Hour.Valid {
		h := row.Hour.Int32
		hour = &h
	}

	return &models.AnalyticsAggregate{
		ID:           id,
		RestaurantID: restaurantID,
		Date:         row.Date.Time,
		Hour:         hour,
		MetricType:   row.MetricType,
		TargetID:     &row.TargetID,
		Value:        row.Value.Int32,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}
