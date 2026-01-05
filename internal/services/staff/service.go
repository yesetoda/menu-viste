package staff

import (
	"context"
	"fmt"
	"log"
	"time"

	"menuvista/internal/models"
	"menuvista/internal/services/email"

	// "menuvista/internal/services/sms"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"
	"menuvista/platform/storage"

	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries      *persistence.Queries
	r2           *storage.R2Client
	emailService *email.Service
	// smsService   *sms.Service
}

func NewStaffService(queries *persistence.Queries, r2 *storage.R2Client, emailService *email.Service) *Service {
	return &Service{
		queries:      queries,
		r2:           r2,
		emailService: emailService,
		// smsService:   smsService,
	}
}

func (s *Service) CreateStaff(ctx context.Context, ownerID uuid.UUID, restaurantID uuid.UUID, input models.CreateUserRequest) (*models.User, error) {
	log.Printf("[StaffService] Creating staff: %s for restaurant: %v", input.Email, restaurantID)

	// Tier validation
	sub, err := s.queries.GetActiveSubscriptionByOwner(ctx, ownerID) // Changed from GetSubscriptionByOwner to GetActiveSubscriptionByOwner and ownerID to *ownerID (corrected to ownerID for syntactic correctness)
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
		log.Printf("[StaffService] Warning: Failed to unmarshal features: %v", err)
	}

	existingStaff, err := s.queries.ListStaffByRestaurant(ctx, persistence.ListStaffByRestaurantParams{RestaurantID: restaurantID})
	if err == nil && len(existingStaff) >= features.MaxStaffAccounts {
		return nil, fmt.Errorf("staff account limit reached for your tier (%d)", features.MaxStaffAccounts)
	}
	password, err := utils.GeneratePassword(12)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userRow, err := s.queries.CreateUser(ctx, persistence.CreateUserParams{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		FullName:     input.FullName,
		Role:         persistence.UserRoleStaff,
		OwnerID:      ownerID,
		RestaurantID: restaurantID,
		Phone:        pgtype.Text{String: input.Phone, Valid: input.Phone != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create staff user: %w", err)
	}

	domainUser := s.mapToDomainUser(userRow)

	// Send notifications async
	go func() {
		// Email
		if err := s.emailService.SendStaffWelcomeEmail(context.Background(), domainUser, password); err != nil {
			log.Printf("[StaffService] Failed to send welcome email: %v", err)
		}

		// SMS (if phone exists)
		// if input.Phone != "" {
		// 	if err := s.smsService.SendCredentials(context.Background(), input.Phone, input.FullName, input.Email, password); err != nil {
		// 		log.Printf("[StaffService] Failed to send SMS: %v", err)
		// 	}
		// }
	}()

	return domainUser, nil
}

func (s *Service) ListStaff(ctx context.Context, restaurantID uuid.UUID, pagination models.PaginationParams) ([]*models.User, *models.Meta, error) {
	fmt.Printf("[StaffService] ListStaff: %v", restaurantID)
	rows, err := s.queries.ListStaffByRestaurant(ctx, persistence.ListStaffByRestaurantParams{
		RestaurantID: restaurantID,
		Limit:        int32(pagination.PageSize),
		Offset:       int32(pagination.Offset),
	})
	fmt.Printf("[StaffService] ListStaffByRestaurant: %v", rows)
	fmt.Printf("[StaffService] ListStaffByRestaurant: %v", err)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list staff: %w", err)
	}

	totalRecords, err := s.queries.CountStaffByRestaurant(ctx, restaurantID)
	if err != nil {
		log.Printf("[StaffService] Warning: Failed to count staff: %v", err)
	}

	staff := make([]*models.User, len(rows))
	for i, row := range rows {
		staff[i] = s.mapToDomainUser(row)
	}

	meta := models.CalculateMeta(pagination.Page, pagination.PageSize, int(totalRecords))

	return staff, meta, nil
}

func (s *Service) UpdateStaffStatus(ctx context.Context, staffID uuid.UUID, restaurantID uuid.UUID, isActive bool) error {
	return s.queries.UpdateStaffStatus(ctx, persistence.UpdateStaffStatusParams{
		ID:           staffID,
		RestaurantID: restaurantID,
		IsActive:     isActive,
	})
}

func (s *Service) DeleteStaff(ctx context.Context, staffID uuid.UUID, restaurantID uuid.UUID) error {
	fmt.Printf("[StaffService] Deleting staff: %v", staffID)
	fmt.Printf("[StaffService] Deleting staff: %v", restaurantID)
	return s.queries.DeleteStaff(ctx, persistence.DeleteStaffParams{
		ID:           staffID,
		RestaurantID: restaurantID,
	})
}

func (s *Service) mapToDomainUser(row persistence.User) *models.User {
	id := row.ID
	ownerID := row.OwnerID
	restaurantID := row.RestaurantID

	return &models.User{
		ID:            id,
		Email:         row.Email,
		FullName:      row.FullName,
		Role:          models.UserRole(row.Role),
		OwnerID:       &ownerID,
		RestaurantID:  &restaurantID,
		Phone:         row.Phone.String,
		AvatarURL:     row.AvatarUrl.String,
		EmailVerified: row.EmailVerified,
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}
func (s *Service) UpdateStaff(ctx context.Context, staffID uuid.UUID, restaurantID uuid.UUID, input models.UpdateUserRequest) (*models.User, error) {
	log.Printf("[StaffService] Updating staff: %v", staffID)

	params := persistence.UpdateUserParams{
		ID:       staffID,
		FullName: pgtype.Text{String: utils.DerefString(input.FullName), Valid: input.FullName != nil},
		Phone:    pgtype.Text{String: utils.DerefString(input.Phone), Valid: input.Phone != nil},
		IsActive: pgtype.Bool{Bool: utils.DerefBool(input.IsActive), Valid: input.IsActive != nil},
	}

	if input.Avatar != nil {
		url, err := s.uploadFile(ctx, fmt.Sprintf("users/%s/avatar", staffID.String()), input.Avatar)
		if err == nil {
			params.AvatarUrl = pgtype.Text{String: url, Valid: true}
		}
	}

	userRow, err := s.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	return s.mapToDomainUser(userRow), nil
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

	return s.r2.UploadFile(ctx, key, f, file.Header.Get("Content-Type"))
}
