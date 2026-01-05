package restaurant

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"menuvista/internal/models"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"
	"menuvista/platform/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmailService interface {
	SendRestaurantApprovalEmail(ctx context.Context, restaurant *persistence.Restaurant, owner *persistence.User) error
	SendRestaurantRejectionEmail(ctx context.Context, restaurant *persistence.Restaurant, owner *persistence.User, reason string) error
}

type Service struct {
	queries      *persistence.Queries
	r2           *storage.R2Client
	emailService EmailService
}

func NewService(queries *persistence.Queries, r2 *storage.R2Client, emailService EmailService) *Service {
	return &Service{
		queries:      queries,
		r2:           r2,
		emailService: emailService,
	}
}

func (s *Service) CreateRestaurant(ctx context.Context, ownerID uuid.UUID, input models.CreateRestaurantRequest) (*models.Restaurant, error) {
	log.Printf("[RestaurantService] Creating restaurant: %s for owner: %v", input.Name, ownerID)

	// Tier validation
	ownerIDStr := ownerID.String()
	sub, err := s.queries.GetActiveSubscriptionByOwner(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	var features models.FeatureLimits
	if err := utils.UnmarshalJSON(sub.Features, &features); err != nil {
		log.Printf("[RestaurantService] Warning: Failed to unmarshal features: %v", err)
	}

	existingRestaurants, err := s.queries.ListRestaurantsByOwner(ctx, utils.ToUUID(&ownerIDStr))
	if err == nil && len(existingRestaurants) >= features.MaxRestaurants {
		return nil, fmt.Errorf("restaurant limit reached for your tier (%d)", features.MaxRestaurants)
	}

	var logoURL, coverURL string
	if input.Logo != nil {
		url, err := s.uploadFile(ctx, fmt.Sprintf("restaurants/%s/logo", input.Slug), input.Logo)
		if err != nil {
			log.Printf("[RestaurantService] Warning: Failed to upload logo: %v", err)
		} else {
			logoURL = url
		}
	}
	if input.CoverImage != nil {
		url, err := s.uploadFile(ctx, fmt.Sprintf("restaurants/%s/cover", input.Slug), input.CoverImage)
		if err != nil {
			log.Printf("[RestaurantService] Warning: Failed to upload cover: %v", err)
		} else {
			coverURL = url
		}
	}

	restaurantRow, err := s.queries.CreateRestaurant(ctx, persistence.CreateRestaurantParams{
		OwnerID:       utils.ToUUID(&ownerIDStr),
		Name:          input.Name,
		Slug:          input.Slug,
		Description:   pgtype.Text{String: input.Description, Valid: input.Description != ""},
		CuisineType:   pgtype.Text{String: input.CuisineType, Valid: input.CuisineType != ""},
		Phone:         pgtype.Text{String: input.Phone, Valid: input.Phone != ""},
		Email:         pgtype.Text{String: input.Email, Valid: input.Email != ""},
		Website:       pgtype.Text{String: input.Website, Valid: input.Website != ""},
		Address:       pgtype.Text{String: input.Address, Valid: input.Address != ""},
		City:          pgtype.Text{String: input.City, Valid: input.City != ""},
		Country:       pgtype.Text{String: input.Country, Valid: input.Country != ""},
		LogoUrl:       pgtype.Text{String: logoURL, Valid: logoURL != ""},
		CoverImageUrl: pgtype.Text{String: coverURL, Valid: coverURL != ""},
		ThemeSettings: input.ThemeSettings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create restaurant: %w", err)
	}

	return s.mapToDomainRestaurant(restaurantRow), nil
}

// func (s *Service) UpdateRestaurantStatus(ctx context.Context, id uuid.UUID, status persistence.RestaurantStatus, reason string) (*models.Restaurant, error) {
// 	idStr := id.String()
// 	restaurantRow, err := s.queries.UpdateRestaurantStatus(ctx, persistence.UpdateRestaurantStatusParams{
// 		ID:     utils.ToUUID(&idStr),
// 		Status: status,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update restaurant status: %w", err)
// 	}

// 	// Fetch owner to send email
// 	owner, err := s.queries.GetUserByID(ctx, restaurantRow.OwnerID)
// 	if err != nil {
// 		log.Printf("[RestaurantService] Warning: Failed to fetch owner for email notification: %v", err)
// 	} else {
// 		if status == persistence.RestaurantStatusApproved {
// 			if err := s.emailService.SendRestaurantApprovalEmail(ctx, &restaurantRow, &owner); err != nil {
// 				log.Printf("[RestaurantService] Warning: Failed to send approval email: %v", err)
// 			}
// 		} else if status == persistence.RestaurantStatusRejected {
// 			if err := s.emailService.SendRestaurantRejectionEmail(ctx, &restaurantRow, &owner, reason); err != nil {
// 				log.Printf("[RestaurantService] Warning: Failed to send rejection email: %v", err)
// 			}
// 		}
// 	}

// 	return s.mapToDomainRestaurant(restaurantRow), nil
// }

func (s *Service) GetRestaurantBySlug(ctx context.Context, slug string) (*models.Restaurant, error) {
	restaurantRow, err := s.queries.GetRestaurantBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("restaurant not found: %w", err)
	}
	return s.mapToDomainRestaurant(restaurantRow), nil
}

func (s *Service) GetRestaurantByID(ctx context.Context, id uuid.UUID) (*models.Restaurant, error) {
	idStr := id.String()
	restaurantRow, err := s.queries.GetRestaurantByID(ctx, utils.ToUUID(&idStr))
	if err != nil {
		return nil, fmt.Errorf("restaurant not found: %w", err)
	}
	return s.mapToDomainRestaurant(restaurantRow), nil
}

func (s *Service) ListRestaurantsByOwner(ctx context.Context, ownerID uuid.UUID) ([]*models.Restaurant, error) {
	ownerIDStr := ownerID.String()
	rows, err := s.queries.ListRestaurantsByOwner(ctx, utils.ToUUID(&ownerIDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to list restaurants: %w", err)
	}

	restaurants := make([]*models.Restaurant, len(rows))
	for i, row := range rows {
		restaurants[i] = s.mapToDomainRestaurant(row)
	}
	return restaurants, nil
}

func (s *Service) UpdateRestaurant(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, isAdmin bool, input models.UpdateRestaurantRequest) (*models.Restaurant, error) {
	log.Printf("[RestaurantService] Updating restaurant: %v", id)

	idStr := id.String()
	ownerIDStr := ownerID.String()
	params := persistence.UpdateRestaurantParams{
		ID:            utils.ToUUID(&idStr),
		OwnerID:       utils.ToUUID(&ownerIDStr),
		IsAdmin:       isAdmin,
		Name:          pgtype.Text{String: utils.DerefString(input.Name), Valid: input.Name != nil},
		Description:   pgtype.Text{String: utils.DerefString(input.Description), Valid: input.Description != nil},
		CuisineType:   pgtype.Text{String: utils.DerefString(input.CuisineType), Valid: input.CuisineType != nil},
		Phone:         pgtype.Text{String: utils.DerefString(input.Phone), Valid: input.Phone != nil},
		Email:         pgtype.Text{String: utils.DerefString(input.Email), Valid: input.Email != nil},
		Website:       pgtype.Text{String: utils.DerefString(input.Website), Valid: input.Website != nil},
		Address:       pgtype.Text{String: utils.DerefString(input.Address), Valid: input.Address != nil},
		City:          pgtype.Text{String: utils.DerefString(input.City), Valid: input.City != nil},
		Country:       pgtype.Text{String: utils.DerefString(input.Country), Valid: input.Country != nil},
		ThemeSettings: input.ThemeSettings,
		IsPublished:   pgtype.Bool{Bool: utils.DerefBool(input.IsPublished), Valid: input.IsPublished != nil},
	}

	if input.Logo != nil {
		url, err := s.uploadFile(ctx, fmt.Sprintf("restaurants/%s/logo", idStr), input.Logo)
		if err == nil {
			params.LogoUrl = pgtype.Text{String: url, Valid: true}
		}
	}
	if input.CoverImage != nil {
		url, err := s.uploadFile(ctx, fmt.Sprintf("restaurants/%s/cover", idStr), input.CoverImage)
		if err == nil {
			params.CoverImageUrl = pgtype.Text{String: url, Valid: true}
		}
	}

	restaurantRow, err := s.queries.UpdateRestaurant(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update restaurant: %w", err)
	}

	return s.mapToDomainRestaurant(restaurantRow), nil
}

func (s *Service) DeleteRestaurant(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	idStr := id.String()
	ownerIDStr := ownerID.String()
	return s.queries.DeleteRestaurant(ctx, persistence.DeleteRestaurantParams{
		ID:      utils.ToUUID(&idStr),
		OwnerID: utils.ToUUID(&ownerIDStr),
	})
}

func (s *Service) ListRestaurantsWithFilters(ctx context.Context, filters models.RestaurantFilter, pagination models.PaginationParams) ([]*models.Restaurant, *models.Meta, error) {
	fmt.Println("this is the pagination data", pagination)
	var ownerID uuid.UUID
	if filters.OwnerID != nil {
		ownerID = utils.ToUUID(filters.OwnerID)
	}
	// var status persistence.NullRestaurantStatus
	// if filters.Status != nil {
	// 	status = persistence.NullRestaurantStatus{
	// 		RestaurantStatus: persistence.RestaurantStatus(*filters.Status),
	// 		Valid:            true,
	// 	}
	// }

	var search pgtype.Text
	if filters.Search != nil {
		search = pgtype.Text{String: *filters.Search, Valid: true}
	}
	fmt.Println("search", search)
	fmt.Println("ownerID", ownerID)
	fmt.Println("filters", filters)
	fmt.Println("pagination", pagination)

	rows, err := s.queries.ListRestaurantsWithFilters(ctx, persistence.ListRestaurantsWithFiltersParams{
		OwnerID:     &ownerID,
		CuisineType: utils.ToText(filters.CuisineType),
		City:        utils.ToText(filters.City),
		Country:     utils.ToText(filters.Country),
		IsPublished: utils.ToBool(filters.IsPublished),
		Search:      search,
		Limit:       int32(pagination.PageSize),
		Offset:      int32(pagination.Offset),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list restaurants: %w", err)
	}

	totalRecords, err := s.queries.CountRestaurantsWithFilters(ctx, persistence.CountRestaurantsWithFiltersParams{
		OwnerID:     ownerID,
		CuisineType: utils.ToText(filters.CuisineType),
		City:        utils.ToText(filters.City),
		Country:     utils.ToText(filters.Country),
		IsPublished: utils.ToBool(filters.IsPublished),
		Search:      search,
	})
	if err != nil {
		log.Printf("[RestaurantService] Warning: Failed to count restaurants: %v", err)
	}

	restaurants := make([]*models.Restaurant, len(rows))
	for i, row := range rows {
		restaurants[i] = s.mapToDomainRestaurant(row)
	}

	fmt.Println("this is the pagination data", pagination)
	meta := models.CalculateMeta(pagination.Page, pagination.PageSize, int(totalRecords))

	return restaurants, meta, nil
}

func (s *Service) mapToDomainRestaurant(row persistence.Restaurant) *models.Restaurant {
	id := row.ID
	ownerID := row.OwnerID

	rankScore, _ := row.RankScore.Float64Value()

	return &models.Restaurant{
		ID:            id,
		OwnerID:       ownerID,
		Name:          row.Name,
		Slug:          row.Slug,
		Description:   row.Description.String,
		CuisineType:   row.CuisineType.String,
		Phone:         row.Phone.String,
		Email:         row.Email.String,
		Website:       row.Website.String,
		Address:       row.Address.String,
		City:          row.City.String,
		Country:       row.Country.String,
		LogoURL:       row.LogoUrl.String,
		CoverImageURL: row.CoverImageUrl.String,
		ThemeSettings: row.ThemeSettings,
		IsPublished:   row.IsPublished,
		// Status:        string(row.Status),
		ViewCount: row.ViewCount.Int32,
		RankScore: rankScore.Float64,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
func (s *Service) uploadFile(ctx context.Context, key string, file *multipart.FileHeader) (string, error) {
	if s.r2 == nil {
		return "", fmt.Errorf("R2 client not initialized")
	}

	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Use the original filename extension if possible
	// For now, we just use the provided key
	return s.r2.UploadFile(ctx, key, f, file.Header.Get("Content-Type"))
}
