package menu

import (
	"context"
	"fmt"
	"log"
	"time"

	"menuvista/internal/models"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"
	"menuvista/platform/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries *persistence.Queries
	r2      *storage.R2Client
}

func NewService(queries *persistence.Queries, r2 *storage.R2Client) *Service {
	return &Service{
		queries: queries,
		r2:      r2,
	}
}

const (
	errCategoryNotFound = "category not found: %w"
	errMenuItemNotFound = "menu item not found: %w"
)

// Categories

// Categories

func (s *Service) CreateCategory(ctx context.Context, userID uuid.UUID, restaurantID uuid.UUID, input models.CreateCategoryRequest) (*models.Category, error) {
	user, err := s.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	log.Printf("[MenuService] Creating category: %s for restaurant: %v by user: %v", input.Name, restaurantID, userID)

	// Verify access
	if err := s.verifyAccess(ctx, user, restaurantID); err != nil {
		return nil, err
	}

	// Tier validation
	ownerID := user.ID
	if user.Role == models.RoleStaff {
		ownerID = *user.OwnerID
	}
	ownerIDStr := ownerID.String()
	sub, err := s.queries.GetSubscriptionByOwner(ctx, utils.ToUUID(&ownerIDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	// Check status and expiration
	now := time.Now()
	isSubActive := false
	if string(sub.Status) == string(models.SubscriptionStatusActive) {
		if sub.CurrentPeriodEnd.Time.After(now) {
			isSubActive = true
		}
	} else if string(sub.Status) == string(models.SubscriptionStatusTrialing) {
		if sub.TrialEnd.Valid && sub.TrialEnd.Time.After(now) {
			isSubActive = true
		}
	}

	if !isSubActive {
		return nil, fmt.Errorf("subscription is inactive or expired")
	}

	var features models.FeatureLimits
	if err := utils.UnmarshalJSON(sub.Features, &features); err != nil {
		log.Printf("[MenuService] Warning: Failed to unmarshal features: %v", err)
	}

	restaurantIDStr := restaurantID.String()
	existingCategories, err := s.queries.ListCategoriesByRestaurant(ctx, utils.ToUUID(&restaurantIDStr))
	if err == nil && len(existingCategories) >= features.MaxCategories {
		return nil, fmt.Errorf("category limit reached for your tier (%d)", features.MaxCategories)
	}

	userIDStr := user.ID.String()
	categoryRow, err := s.queries.CreateCategory(ctx, persistence.CreateCategoryParams{
		RestaurantID: utils.ToUUID(&restaurantIDStr),
		Name:         input.Name,
		Description:  pgtype.Text{String: input.Description, Valid: input.Description != ""},
		Icon:         pgtype.Text{String: input.Icon, Valid: input.Icon != ""},
		DisplayOrder: input.DisplayOrder,
		IsActive:     input.IsActive,
		CreatedBy:    utils.ToUUID(&userIDStr),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return s.mapToDomainCategory(categoryRow), nil
}

func (s *Service) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	idStr := id.String()
	row, err := s.queries.GetCategoryByID(ctx, utils.ToUUID(&idStr))
	if err != nil {
		return nil, fmt.Errorf(errCategoryNotFound, err)
	}
	return s.mapToDomainCategory(row), nil
}

func (s *Service) ListCategories(ctx context.Context, restaurantID uuid.UUID) ([]*models.Category, error) {
	restaurantIDStr := restaurantID.String()
	rows, err := s.queries.ListCategoriesByRestaurant(ctx, utils.ToUUID(&restaurantIDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	categories := make([]*models.Category, len(rows))
	for i, row := range rows {
		categories[i] = s.mapToDomainCategory(row)
	}
	return categories, nil
}

func (s *Service) UpdateCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID, input models.UpdateCategoryRequest) (*models.Category, error) {
	user, err := s.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	idStr := id.String()
	category, err := s.queries.GetCategoryByID(ctx, utils.ToUUID(&idStr))
	if err != nil {
		return nil, fmt.Errorf(errCategoryNotFound, err)
	}

	if err := s.verifyAccess(ctx, user, category.RestaurantID); err != nil {
		return nil, err
	}

	categoryRow, err := s.queries.UpdateCategory(ctx, persistence.UpdateCategoryParams{
		ID:           utils.ToUUID(&idStr),
		Name:         pgtype.Text{String: utils.DerefString(input.Name), Valid: input.Name != nil},
		Description:  pgtype.Text{String: utils.DerefString(input.Description), Valid: input.Description != nil},
		Icon:         pgtype.Text{String: utils.DerefString(input.Icon), Valid: input.Icon != nil},
		DisplayOrder: pgtype.Int4{Int32: utils.DerefInt32(input.DisplayOrder), Valid: input.DisplayOrder != nil},
		IsActive:     pgtype.Bool{Bool: utils.DerefBool(input.IsActive), Valid: input.IsActive != nil},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return s.mapToDomainCategory(categoryRow), nil
}

func (s *Service) DeleteCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	user, err := s.getUser(ctx, userID)
	if err != nil {
		return err
	}

	idStr := id.String()
	category, err := s.queries.GetCategoryByID(ctx, utils.ToUUID(&idStr))
	if err != nil {
		return fmt.Errorf(errCategoryNotFound, err)
	}

	if err := s.verifyAccess(ctx, user, category.RestaurantID); err != nil {
		return err
	}

	return s.queries.DeleteCategory(ctx, utils.ToUUID(&idStr))
}

// Items

func (s *Service) CreateMenuItem(ctx context.Context, userID uuid.UUID, restaurantID uuid.UUID, input models.CreateMenuItemRequest) (*models.MenuItem, error) {
	user, err := s.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	log.Printf("[MenuService] Creating item: %s for restaurant: %v by user: %v", input.Name, restaurantID, userID)

	// Verify access
	if err := s.verifyAccess(ctx, user, restaurantID); err != nil {
		return nil, err
	}

	// Tier validation
	ownerID := user.ID
	if user.Role == models.RoleStaff {
		ownerID = *user.OwnerID
	}
	ownerIDStr := ownerID.String()
	sub, err := s.queries.GetSubscriptionByOwner(ctx, utils.ToUUID(&ownerIDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	// Check status and expiration
	now := time.Now()
	isSubActive := false
	if string(sub.Status) == string(models.SubscriptionStatusActive) {
		if sub.CurrentPeriodEnd.Time.After(now) {
			isSubActive = true
		}
	} else if string(sub.Status) == string(models.SubscriptionStatusTrialing) {
		if sub.TrialEnd.Valid && sub.TrialEnd.Time.After(now) {
			isSubActive = true
		}
	}

	if !isSubActive {
		return nil, fmt.Errorf("subscription is inactive or expired")
	}

	var features models.FeatureLimits
	if err := utils.UnmarshalJSON(sub.Features, &features); err != nil {
		log.Printf("[MenuService] Warning: Failed to unmarshal features: %v", err)
	}

	restaurantIDStr := restaurantID.String()
	existingItems, err := s.queries.ListMenuItemsByRestaurant(ctx, utils.ToUUID(&restaurantIDStr))
	if err == nil && len(existingItems) >= features.MaxMenuItems {
		return nil, fmt.Errorf("menu item limit reached for your tier (%d)", features.MaxMenuItems)
	}

	categoryIDStr := input.CategoryID.String()
	userIDStr := user.ID.String()
	itemRow, err := s.queries.CreateMenuItem(ctx, persistence.CreateMenuItemParams{
		RestaurantID: utils.ToUUID(&restaurantIDStr),
		CategoryID:   utils.ToUUID(&categoryIDStr),
		Name:         input.Name,
		Description:  pgtype.Text{String: input.Description, Valid: input.Description != ""},
		Price:        utils.ToNumeric(input.Price),
		Currency:     input.Currency,
		Images:       input.Images,
		Allergens:    input.Allergens,
		DietaryTags:  input.DietaryTags,
		SpiceLevel:   pgtype.Int4{Int32: input.SpiceLevel, Valid: true},
		Calories:     pgtype.Int4{Int32: input.Calories, Valid: input.Calories != 0},
		IsAvailable:  input.IsAvailable,
		DisplayOrder: input.DisplayOrder,
		CreatedBy:    utils.ToUUID(&userIDStr),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create menu item: %w", err)
	}

	return s.mapToDomainMenuItem(itemRow), nil
}

func (s *Service) GetMenuItemByID(ctx context.Context, id uuid.UUID) (*models.MenuItem, error) {
	idStr := id.String()
	row, err := s.queries.GetMenuItemByID(ctx, utils.ToUUID(&idStr))
	if err != nil {
		return nil, fmt.Errorf(errMenuItemNotFound, err)
	}
	return s.mapToDomainMenuItem(row), nil
}

func (s *Service) ListMenuItems(ctx context.Context, restaurantID uuid.UUID) ([]*models.MenuItem, error) {
	restaurantIDStr := restaurantID.String()
	rows, err := s.queries.ListMenuItemsByRestaurant(ctx, utils.ToUUID(&restaurantIDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items: %w", err)
	}

	items := make([]*models.MenuItem, len(rows))
	for i, row := range rows {
		items[i] = s.mapToDomainMenuItem(row)
	}
	return items, nil
}

func (s *Service) UpdateMenuItem(ctx context.Context, userID uuid.UUID, id uuid.UUID, input models.UpdateMenuItemRequest) (*models.MenuItem, error) {
	user, err := s.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	idStr := id.String()
	item, err := s.queries.GetMenuItemByID(ctx, utils.ToUUID(&idStr))
	if err != nil {
		return nil, fmt.Errorf(errMenuItemNotFound, err)
	}

	if err := s.verifyAccess(ctx, user, item.RestaurantID); err != nil {
		return nil, err
	}

	categoryID := utils.DerefUUID(input.CategoryID)
	categoryIDStr := categoryID.String()
	catID := utils.ToUUID(&categoryIDStr)
	params := persistence.UpdateMenuItemParams{
		ID:           utils.ToUUID(&idStr),
		CategoryID:   &catID,
		Name:         pgtype.Text{String: utils.DerefString(input.Name), Valid: input.Name != nil},
		Description:  pgtype.Text{String: utils.DerefString(input.Description), Valid: input.Description != nil},
		Price:        utils.ToNumeric(utils.DerefFloat64(input.Price)),
		Images:       input.Images,
		Allergens:    input.Allergens,
		DietaryTags:  input.DietaryTags,
		SpiceLevel:   pgtype.Int4{Int32: utils.DerefInt32(input.SpiceLevel), Valid: input.SpiceLevel != nil},
		Calories:     pgtype.Int4{Int32: utils.DerefInt32(input.Calories), Valid: input.Calories != nil},
		IsAvailable:  pgtype.Bool{Bool: utils.DerefBool(input.IsAvailable), Valid: input.IsAvailable != nil},
		DisplayOrder: pgtype.Int4{Int32: utils.DerefInt32(input.DisplayOrder), Valid: input.DisplayOrder != nil},
	}

	itemRow, err := s.queries.UpdateMenuItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update menu item: %w", err)
	}

	return s.mapToDomainMenuItem(itemRow), nil
}

func (s *Service) DeleteMenuItem(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	user, err := s.getUser(ctx, userID)
	if err != nil {
		return err
	}

	idStr := id.String()
	item, err := s.queries.GetMenuItemByID(ctx, utils.ToUUID(&idStr))
	if err != nil {
		return fmt.Errorf("menu item not found: %w", err)
	}

	if err := s.verifyAccess(ctx, user, item.RestaurantID); err != nil {
		return err
	}

	return s.queries.DeleteMenuItem(ctx, utils.ToUUID(&idStr))
}

// Helpers

func (s *Service) getUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	userIDStr := userID.String()
	userRow, err := s.queries.GetUserByID(ctx, utils.ToUUID(&userIDStr))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	role := models.UserRole(userRow.Role)

	ownerID := userRow.OwnerID
	restaurantID := userRow.RestaurantID

	return &models.User{
		ID:           userID,
		Email:        userRow.Email,
		FullName:     userRow.FullName,
		Role:         role,
		OwnerID:      ownerID,
		RestaurantID: restaurantID,
	}, nil
}

func (s *Service) verifyAccess(ctx context.Context, user *models.User, restaurantID uuid.UUID) error {
	if user.Role == models.RoleAdmin {
		return nil
	}

	restaurantIDStr := restaurantID.String()
	restaurant, err := s.queries.GetRestaurantByID(ctx, utils.ToUUID(&restaurantIDStr))
	if err != nil {
		return fmt.Errorf("failed to fetch restaurant: %w", err)
	}

	if user.Role == models.RoleOwner {
		if restaurant.OwnerID != user.ID {
			return fmt.Errorf("unauthorized: you do not own this restaurant")
		}
	} else if user.Role == models.RoleStaff {
		if user.RestaurantID == nil || *user.RestaurantID != restaurantID {
			return fmt.Errorf("unauthorized: you are not assigned to this restaurant")
		}
	}

	return nil
}

func (s *Service) mapToDomainCategory(row persistence.Category) *models.Category {
	id := row.ID
	restaurantID := row.RestaurantID
	createdBy := row.CreatedBy

	return &models.Category{
		ID:           id,
		RestaurantID: restaurantID,
		Name:         row.Name,
		Description:  row.Description.String,
		Icon:         row.Icon.String,
		DisplayOrder: row.DisplayOrder,
		IsActive:     row.IsActive,
		CreatedBy:    createdBy,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

func (s *Service) mapToDomainMenuItem(row persistence.MenuItem) *models.MenuItem {
	id := row.ID
	restaurantID := row.RestaurantID
	categoryID := row.CategoryID
	createdBy := row.CreatedBy

	price, _ := row.Price.Float64Value()

	return &models.MenuItem{
		ID:           id,
		RestaurantID: restaurantID,
		CategoryID:   categoryID,
		Name:         row.Name,
		Description:  row.Description.String,
		Price:        price.Float64,
		Currency:     row.Currency,
		Images:       row.Images,
		Allergens:    row.Allergens,
		DietaryTags:  row.DietaryTags,
		SpiceLevel:   row.SpiceLevel.Int32,
		Calories:     row.Calories.Int32,
		IsAvailable:  row.IsAvailable,
		DisplayOrder: row.DisplayOrder,
		ViewCount:    row.ViewCount.Int32,
		CreatedBy:    createdBy,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}
