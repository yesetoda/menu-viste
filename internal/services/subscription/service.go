package subscription

import (
	"context"
	"fmt"
	"time"

	"menuvista/internal/models"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"

	"github.com/google/uuid"
)

type Service struct {
	queries *persistence.Queries
}

func NewService(queries *persistence.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) ListPlans(ctx context.Context) ([]*models.SubscriptionPlan, error) {
	rows, err := s.queries.ListSubscriptionPlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list plans: %w", err)
	}

	plans := make([]*models.SubscriptionPlan, len(rows))
	for i, row := range rows {
		plans[i] = s.mapToDomainPlan(row)
	}
	return plans, nil
}

func (s *Service) GetPlanBySlug(ctx context.Context, slug string) (*models.SubscriptionPlan, error) {
	row, err := s.queries.GetSubscriptionPlanBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}
	return s.mapToDomainPlan(row), nil
}

func (s *Service) GetSubscriptionByOwner(ctx context.Context, ownerID uuid.UUID) (*models.Subscription, error) {
	ownerIDStr := ownerID.String()
	row, err := s.queries.GetActiveSubscriptionByOwner(ctx, utils.ToUUID(&ownerIDStr))
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	return s.mapToDomainSubscription(row), nil
}

func (s *Service) GetSubscriptionDetails(ctx context.Context, ownerID uuid.UUID) (*models.SubscriptionDetailsResponse, error) {
	sub, err := s.GetSubscriptionByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	plan, err := s.GetPlanBySlug(ctx, sub.PlanSlug)
	if err != nil {
		return nil, err
	}

	var daysRemaining int
	now := time.Now()

	if sub.Status == models.SubscriptionStatusTrialing && sub.TrialEnd != nil {
		daysRemaining = int(sub.TrialEnd.Sub(now).Hours() / 24)
	} else {
		daysRemaining = int(sub.CurrentPeriodEnd.Sub(now).Hours() / 24)
	}

	if daysRemaining < 0 {
		daysRemaining = 0
	}

	return &models.SubscriptionDetailsResponse{
		PlanName:      plan.Name,
		PlanSlug:      plan.Slug,
		Price:         plan.PriceMonthly, // Assuming monthly for now
		Currency:      plan.Currency,
		Status:        sub.Status,
		StartDate:     sub.CurrentPeriodStart,
		EndDate:       sub.CurrentPeriodEnd,
		TrialEnd:      sub.TrialEnd,
		DaysRemaining: daysRemaining,
		Features:      plan.Features,
	}, nil
}

func (s *Service) mapToDomainPlan(row persistence.SubscriptionPlan) *models.SubscriptionPlan {

	var features models.FeatureLimits
	utils.UnmarshalJSON(row.Features, &features)

	return &models.SubscriptionPlan{
		ID:           row.ID,
		Name:         row.Name,
		Slug:         row.Slug,
		Description:  row.Description.String,
		PriceMonthly: row.PriceMonthly,
		PriceAnnual:  row.PriceAnnual.Int32,
		Currency:     row.Currency,
		Features:     features,
		DisplayOrder: row.DisplayOrder,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

func (s *Service) mapToDomainSubscription(row persistence.GetActiveSubscriptionByOwnerRow) *models.Subscription {
	id := row.ID
	ownerID := row.OwnerID
	planID := row.PlanID

	var features models.FeatureLimits
	utils.UnmarshalJSON(row.Features, &features)

	var trialEnd *time.Time
	if row.TrialEnd.Valid {
		trialEnd = &row.TrialEnd.Time
	}
	var cancelledAt *time.Time
	if row.CancelledAt.Valid {
		cancelledAt = &row.CancelledAt.Time
	}

	return &models.Subscription{
		ID:                 id,
		OwnerID:            ownerID,
		PlanID:             planID,
		Status:             models.SubscriptionStatus(row.Status),
		CurrentPeriodStart: row.CurrentPeriodStart.Time,
		CurrentPeriodEnd:   row.CurrentPeriodEnd.Time,
		TrialEnd:           trialEnd,
		CancelledAt:        cancelledAt,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
		PlanName:           row.PlanName,
		PlanSlug:           row.PlanSlug,
		Features:           features,
	}
}
