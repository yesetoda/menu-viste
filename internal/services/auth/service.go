package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"menuvista/internal/models"
	"menuvista/internal/services/payment"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"
	"menuvista/platform/cache"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	queries        *persistence.Queries
	redis          *cache.RedisClient
	emailService   EmailService
	paymentService PaymentService
}

// EmailService interface for sending emails
type EmailService interface {
	SendWelcomeEmail(ctx context.Context, user *models.User) error
	SendVerificationEmail(ctx context.Context, user *models.User, token string) error
}

// PaymentService interface for initiating payments
type PaymentService interface {
	InitiatePayment(ctx context.Context, input payment.InitiatePaymentInput) (*payment.InitiatePaymentResponse, error)
}

func NewService(queries *persistence.Queries, redis *cache.RedisClient, emailService EmailService, paymentService PaymentService) *Service {
	return &Service{
		queries:        queries,
		redis:          redis,
		emailService:   emailService,
		paymentService: paymentService,
	}
}

func (s *Service) Register(ctx context.Context, input models.CreateUserRequest) (*models.AuthResponse, error) {
	log.Printf("[AuthService] Registering user: %s with role: %s", input.Email, input.Role)

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user in DB
	userRow, err := s.queries.CreateUser(ctx, persistence.CreateUserParams{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		FullName:     input.FullName,
		Role:         persistence.UserRoleOwner,
		Phone:        pgtype.Text{String: input.Phone, Valid: input.Phone != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Immediately set active=false and email_verified=false
	userRow, err = s.queries.UpdateUser(ctx, persistence.UpdateUserParams{
		ID:            userRow.ID,
		IsActive:      pgtype.Bool{Bool: false, Valid: true},
		EmailVerified: pgtype.Bool{Bool: false, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set initial user state: %w", err)
	}

	domainUser := s.mapToDomainUser(userRow)

	// Generate activation token and send email
	token, err := s.GenerateActivationToken(ctx, domainUser.ID)
	if err != nil {
		log.Printf("[AuthService] Failed to generate activation token: %v", err)
		// Don't fail registration, user can request new token later
	} else {
		go func() {
			if err := s.emailService.SendVerificationEmail(context.Background(), domainUser, token); err != nil {
				log.Printf("[AuthService] Failed to send verification email: %v", err)
			}
		}()
	}

	// Return empty auth response or specific message
	// For now, returning user object but no tokens
	return &models.AuthResponse{
		User: *domainUser,
	}, nil
}

func (s *Service) GenerateActivationToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token := uuid.New().String()
	// Store in Redis: token -> userID
	err := s.redis.Set(ctx, redisKeyActivationToken+token, utils.UUIDToString(userID), 24*60*60) // 24 hours
	if err != nil {
		return "", err
	}
	// Store in Redis: userID -> token (to check if valid token exists)
	err = s.redis.Set(ctx, redisKeyUserActivation+utils.UUIDToString(userID), token, 24*60*60)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) ActivateUser(ctx context.Context, token string) (*models.AuthResponse, error) {
	userIDStr, err := s.redis.Get(ctx, redisKeyActivationToken+token)
	if err != nil {
		return nil, errors.New("invalid or expired activation token")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	// Update user in DB
	userRow, err := s.queries.UpdateUser(ctx, persistence.UpdateUserParams{
		ID:            userID,
		IsActive:      pgtype.Bool{Bool: true, Valid: true},
		EmailVerified: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	// Create Free Trial Subscription
	// 1. Get Free Plan
	freePlan, err := s.queries.GetSubscriptionPlanBySlug(ctx, "free-trial")
	if err != nil {
		log.Printf("[AuthService] Warning: Free plan not found, skipping trial creation: %v", err)
	} else {
		// 2. Create Subscription
		trialEnd := time.Now().AddDate(0, 0, 14) // 14 days trial
		_, err = s.queries.CreateSubscription(ctx, persistence.CreateSubscriptionParams{
			OwnerID:            userID,
			PlanID:             freePlan.ID,
			Status:             persistence.SubscriptionStatusTrialing,
			CurrentPeriodStart: pgtype.Timestamp{Time: time.Now(), Valid: true},
			CurrentPeriodEnd:   pgtype.Timestamp{Time: trialEnd, Valid: true},
			TrialEnd:           pgtype.Timestamp{Time: trialEnd, Valid: true},
		})
		if err != nil {
			log.Printf("[AuthService] Warning: Failed to create trial subscription: %v", err)
		}
	}

	// Delete tokens from Redis
	// s.redis.Del(ctx, redisKeyActivationToken+token) // Need Del method in RedisClient
	// s.redis.Del(ctx, redisKeyUserActivation+userIDStr)

	domainUser := s.mapToDomainUser(userRow)
	s.sendWelcomeEmailAsync(domainUser)

	return s.generateAuthResponse(ctx, userRow, "")
}

func (s *Service) ResendActivationEmail(ctx context.Context, email string) error {
	userRow, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if userRow.IsActive && userRow.EmailVerified {
		return errors.New("account already active")
	}

	// Check if valid token exists
	_, err = s.redis.Get(ctx, redisKeyUserActivation+utils.UUIDToString(userRow.ID))
	if err == nil {
		return errors.New("activation email already sent. please check your inbox")
	}

	token, err := s.GenerateActivationToken(ctx, userRow.ID)
	if err != nil {
		return err
	}

	domainUser := s.mapToDomainUser(userRow)
	go func() {
		if err := s.emailService.SendVerificationEmail(context.Background(), domainUser, token); err != nil {
			log.Printf("[AuthService] Failed to send verification email: %v", err)
		}
	}()

	return nil
}

func (s *Service) handleOwnerSubscription(ctx context.Context, user persistence.User, planSlug string) (string, error) {
	if planSlug == "" {
		planSlug = "free"
	}

	plan, err := s.queries.GetSubscriptionPlanBySlug(ctx, planSlug)
	if err != nil {
		log.Printf("[AuthService] Warning: Subscription plan '%s' not found: %v", planSlug, err)
		if planSlug != "free" {
			plan, _ = s.queries.GetSubscriptionPlanBySlug(ctx, "free")
		}
	}

	if plan.ID == (uuid.UUID{}) {
		return "", nil
	}

	status := persistence.SubscriptionStatusActive
	if plan.PriceMonthly > 0 {
		status = persistence.SubscriptionStatusIncomplete
	}

	sub, err := s.queries.CreateSubscription(ctx, persistence.CreateSubscriptionParams{
		OwnerID:            user.ID,
		PlanID:             plan.ID,
		Status:             status,
		CurrentPeriodStart: pgtype.Timestamp{Time: time.Now(), Valid: true},
		CurrentPeriodEnd:   pgtype.Timestamp{Time: time.Now().AddDate(0, 1, 0), Valid: true},
	})
	if err != nil {
		log.Printf("[AuthService] Warning: Failed to create initial subscription: %v", err)
	}

	if status == persistence.SubscriptionStatusIncomplete && s.paymentService != nil {
		payResp, err := s.paymentService.InitiatePayment(ctx, payment.InitiatePaymentInput{
			OwnerID:        user.ID,
			SubscriptionID: sub.ID,
			Plan:           plan.Slug,
			Email:          user.Email,
			Name:           user.FullName,
			Type:           "registration",
		})
		if err == nil {
			return payResp.CheckoutURL, nil
		}
		log.Printf("[AuthService] Warning: Failed to initiate payment: %v", err)
	}

	return "", nil
}

func (s *Service) sendWelcomeEmailAsync(user *models.User) {
	if s.emailService != nil {
		go func() {
			if err := s.emailService.SendWelcomeEmail(context.Background(), user); err != nil {
				log.Printf("[AuthService] Failed to send welcome email: %v", err)
			}
		}()
	}
}

func (s *Service) generateAuthResponse(ctx context.Context, user persistence.User, checkoutURL string) (*models.AuthResponse, error) {
	var subStatus string
	var subEnd *time.Time

	if user.Role == persistence.UserRoleOwner {
		sub, err := s.queries.GetActiveSubscriptionByOwner(ctx, user.ID)
		if err == nil {
			subStatus = string(sub.Status)
			if sub.CurrentPeriodEnd.Valid {
				t := sub.CurrentPeriodEnd.Time
				subEnd = &t
			}
		}
	}

	tokenDetails, err := utils.CreateToken(user.ID, string(user.Role), &user.OwnerID, &user.RestaurantID, subStatus, subEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}

	// Store token metadata in Redis
	_ = s.redis.Set(ctx, tokenDetails.AccessUUID, utils.UUIDToString(user.ID), int(tokenDetails.AtExpires-time.Now().Unix()))
	_ = s.redis.Set(ctx, tokenDetails.RefreshUUID, utils.UUIDToString(user.ID), int(tokenDetails.RtExpires-time.Now().Unix()))

	return &models.AuthResponse{
		User:         *s.mapToDomainUser(user),
		AccessToken:  tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshToken,
		CheckoutURL:  checkoutURL,
	}, nil
}

var (
	ErrPaymentRequired      = errors.New("payment required")
	ErrSubscriptionInactive = errors.New("subscription inactive")

	// Redis keys
	redisKeyActivationToken = "activation_token:"
	redisKeyUserActivation  = "user_activation:"
)

func (s *Service) Login(ctx context.Context, input models.LoginRequest) (*models.AuthResponse, error) {
	log.Printf("[AuthService] Login attempt for: %s", input.Email)

	userRow, err := s.queries.GetUserByEmail(ctx, input.Email)
	fmt.Println("this is the input", input)
	fmt.Println("this is the user row", userRow)
	fmt.Println("this is the error", err)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(input.Password, userRow.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if !userRow.IsActive {
		// Check if token expired and resend if needed
		// For now, just return error telling them to check email
		// Logic to resend if expired is in ResendActivationEmail which client can call
		// Or we can auto-resend here?
		// Requirement: "if the activation token expires the user can try to login but since they are not email verified and not active it and since the token is expired it should send a new activation email"

		_, err := s.redis.Get(ctx, redisKeyUserActivation+utils.UUIDToString(userRow.ID))
		if err != nil {
			// Token expired or not found, generate new one
			token, err := s.GenerateActivationToken(ctx, userRow.ID)
			if err == nil {
				domainUser := s.mapToDomainUser(userRow)
				go func() {
					_ = s.emailService.SendVerificationEmail(context.Background(), domainUser, token)
				}()
				return nil, errors.New("activation link expired. a new one has been sent to your email")
			}
		}

		return nil, errors.New("account is not activated. please check your email")
	}

	// Update last login
	_, _ = s.queries.UpdateUser(ctx, persistence.UpdateUserParams{
		ID:          userRow.ID,
		LastLoginAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	userID := userRow.ID
	ownerID := userRow.OwnerID
	restaurantID := userRow.RestaurantID

	var subStatus string
	var subEnd *time.Time
	var checkoutURL string

	// Check subscription status
	if userRow.Role == persistence.UserRoleOwner {
		sub, err := s.queries.GetActiveSubscriptionByOwner(ctx, userID)
		if err == nil {
			// Check if subscription is active or trialing
			isActive := sub.Status == persistence.SubscriptionStatusActive || sub.Status == persistence.SubscriptionStatusTrialing

			// If not active, check if we need to enforce payment
			if !isActive {
				// If it's incomplete or past_due, we should prompt for payment
				// Generate checkout URL
				// plan is not used directly in InitiatePaymentInput as it takes slug, but we need to check if err is nil
				_, err := s.queries.GetSubscriptionPlanBySlug(ctx, sub.PlanSlug)
				if err == nil && s.paymentService != nil {
					payResp, err := s.paymentService.InitiatePayment(ctx, payment.InitiatePaymentInput{
						OwnerID:        userID,
						SubscriptionID: sub.ID,
						Plan:           sub.PlanSlug,
						Email:          userRow.Email,
						Name:           userRow.FullName,
						Type:           "renewal", // or activation
					})
					if err == nil {
						checkoutURL = payResp.CheckoutURL
						return &models.AuthResponse{
							CheckoutURL: checkoutURL,
						}, ErrPaymentRequired
					}
				}

				return &models.AuthResponse{
					CheckoutURL: checkoutURL,
				}, ErrPaymentRequired
			}

			subStatus = string(sub.Status)
			if sub.CurrentPeriodEnd.Valid {
				t := sub.CurrentPeriodEnd.Time
				subEnd = &t
			}
		}
	} else if userRow.Role == persistence.UserRoleStaff {
		// Check owner's subscription
		if &ownerID != nil {
			sub, err := s.queries.GetActiveSubscriptionByOwner(ctx, ownerID)
			if err == nil {
				isActive := sub.Status == persistence.SubscriptionStatusActive || sub.Status == persistence.SubscriptionStatusTrialing
				if !isActive {
					return nil, ErrSubscriptionInactive
				}
			} else {
				// If owner has no subscription, staff shouldn't login?
				// Assuming yes.
				return nil, ErrSubscriptionInactive
			}
		}
	}

	tokenDetails, err := utils.CreateToken(userID, string(userRow.Role), &ownerID, &restaurantID, subStatus, subEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}

	// Store token metadata in Redis
	_ = s.redis.Set(ctx, tokenDetails.AccessUUID, utils.UUIDToString(userID), int(tokenDetails.AtExpires-time.Now().Unix()))
	_ = s.redis.Set(ctx, tokenDetails.RefreshUUID, utils.UUIDToString(userID), int(tokenDetails.RtExpires-time.Now().Unix()))

	return &models.AuthResponse{
		User:         *s.mapToDomainUser(userRow),
		AccessToken:  tokenDetails.AccessToken,
		RefreshToken: tokenDetails.RefreshToken,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	userRow, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return s.mapToDomainUser(userRow), nil
}

func (s *Service) mapToDomainUser(row persistence.User) *models.User {
	id := row.ID
	ownerID := row.OwnerID
	restaurantID := row.RestaurantID

	var lastLogin *time.Time
	if row.LastLoginAt.Valid {
		lastLogin = &row.LastLoginAt.Time
	}

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
		LastLoginAt:   lastLogin,
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}
