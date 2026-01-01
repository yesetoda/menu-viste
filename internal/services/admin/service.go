package admin

import (
	"context"
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

func (s *Service) GetStats(ctx context.Context) (*persistence.GetAdminDashboardStatsRow, error) {
	log.Printf("[AdminService] Fetching admin stats")
	stats, err := s.queries.GetAdminDashboardStats(ctx)
	if err != nil {
		log.Printf("[AdminService] GetAdminDashboardStats error: %v", err)
		return nil, err
	}
	return &stats, nil
}

func (s *Service) GetRecentLogs(ctx context.Context, limit int32) ([]persistence.GetRecentAdminLogsRow, error) {
	log.Printf("[AdminService] Fetching recent admin logs (limit: %d)", limit)
	logs, err := s.queries.GetRecentAdminLogs(ctx, limit)
	if err != nil {
		log.Printf("[AdminService] GetRecentAdminLogs error: %v", err)
		return nil, err
	}
	return logs, nil
}

func (s *Service) GetRestaurantDetails(ctx context.Context, id uuid.UUID) (*persistence.GetRestaurantDetailsForAdminRow, error) {
	log.Printf("[AdminService] Fetching restaurant details for admin: %v", id)
	details, err := s.queries.GetRestaurantDetailsForAdmin(ctx, id)
	if err != nil {
		log.Printf("[AdminService] GetRestaurantDetailsForAdmin error: %v", err)
		return nil, err
	}
	return &details, nil
}

func (s *Service) ListUsers(ctx context.Context, limit, offset int32) ([]persistence.User, error) {
	log.Printf("[AdminService] Listing users limit=%d offset=%d", limit, offset)
	users, err := s.queries.ListUsers(ctx, persistence.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Printf("[AdminService] ListUsers error: %v", err)
		return nil, err
	}
	return users, nil
}

func (s *Service) ListUsersWithFilters(ctx context.Context, filters models.FilterParams, pagination models.PaginationParams) ([]persistence.User, error) {
	log.Printf("[AdminService] Listing users with filters")

	role, _ := utils.ToUserRole(filters.Role)
	users, err := s.queries.ListUsersWithFilters(ctx, persistence.ListUsersWithFiltersParams{
		Limit:    int32(pagination.PageSize),
		Offset:   int32(pagination.Offset),
		Email:    utils.ToText(filters.Email),
		Role:     persistence.NullUserRole{UserRole: role, Valid: filters.Role != nil && *filters.Role != ""},
		IsActive: utils.ToBool(filters.IsActive),
		Search:   utils.ToText(filters.Search),
	})
	if err != nil {
		log.Printf("[AdminService] ListUsersWithFilters error: %v", err)
		return nil, err
	}
	return users, nil
}

func (s *Service) UpdateUserStatus(ctx context.Context, userID uuid.UUID, isActive bool) (*persistence.User, error) {
	log.Printf("[AdminService] Updating user status: %v to %v", userID, isActive)
	user, err := s.queries.UpdateUser(ctx, persistence.UpdateUserParams{
		ID:       userID,
		IsActive: pgtype.Bool{Bool: isActive, Valid: true},
	})
	if err != nil {
		log.Printf("[AdminService] UpdateUserStatus error: %v", err)
		return nil, err
	}
	return &user, nil
}
