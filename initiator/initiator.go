package initiator

import (
	"context"
	"fmt"
	"log"
	"os"

	"menuvista/internal/glue"
	"menuvista/internal/handlers/middleware"
	"menuvista/internal/services/activity"
	"menuvista/internal/services/admin"
	"menuvista/internal/services/analytics"
	"menuvista/internal/services/auth"
	"menuvista/internal/services/email"
	"menuvista/internal/services/menu"
	"menuvista/internal/services/payment"
	"menuvista/internal/services/restaurant"
	"menuvista/internal/services/staff"
	"menuvista/internal/services/subscription"
	"menuvista/internal/storage/persistence"
	"menuvista/platform/cache"
	"menuvista/platform/core"
	"menuvista/platform/storage"

	"github.com/gin-gonic/gin"
)

func Init(ctx context.Context) *gin.Engine {
	// 1. Database
	dbPool, err := core.NewPostgresDB(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	queries := persistence.New(dbPool)

	// 2. Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		// Default or error? Let's assume localhost if not set for dev, or error
		redisURL = "redis://localhost:6379"
	}
	fmt.Println("this is the redis url", redisURL)
	redisClient, err := cache.NewRedisClient(ctx, redisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 3. R2 Storage
	r2Client, err := storage.NewR2Client(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize R2 client: %v", err)
	}

	// 4. Email Service
	resendAPIKey := os.Getenv("RESEND_API_KEY")
	if resendAPIKey == "" {
		log.Printf("WARNING: RESEND_API_KEY not set, email notifications will be disabled")
	}
	emailService := email.NewService(resendAPIKey, queries)

	// 5. Services
	paymentService := payment.NewService(queries)
	webhookService := payment.NewWebhookService(queries, emailService)
	authService := auth.NewService(queries, redisClient, emailService, paymentService)
	restaurantService := restaurant.NewService(queries, r2Client, emailService)
	menuService := menu.NewService(queries, r2Client)
	adminService := admin.NewService(queries)
	staffService := staff.NewStaffService(queries, emailService)
	activityService := activity.NewService(queries)
	analyticsService := analytics.NewService(queries)
	subscriptionService := subscription.NewService(queries)

	// Assuming cfg and logger are defined elsewhere or need to be added.
	// For now, I'll use the existing os.Getenv and log.New for the first two arguments
	// and add a placeholder for smsService as it's not defined in the current context.
	// You will need to define `smsService` and potentially `cfg` and `logger`
	// if you intend to use them as shown in your example snippet.
	// For a syntactically correct change based *only* on adding smsService,
	// I'll add it as a nil placeholder.
	// If `cfg.AuthSecret` and `logger` are intended to replace the existing arguments,
	// then `cfg` and `logger` must be defined.
	// Given the instruction "Pass smsService to NewAuthMiddleware", and the provided
	// "Code Edit" snippet which is syntactically broken, I will interpret it as
	// adding `smsService` as a new argument, keeping the existing ones.
	// If `cfg.AuthSecret` and `logger` are indeed meant to replace the existing arguments,
	// please provide the full context for `cfg` and `logger`.
	// For now, I'll assume `smsService` is the *only* new argument to be added.
	// Since `smsService` is not defined, I'll add a `nil` placeholder.
	authMiddleware := middleware.NewAuthMiddleware(os.Getenv("AUTH_SECRET"), log.New(os.Stdout, "[AuthMiddleware] ", log.LstdFlags), nil /* smsService */)

	// 5. Router
	router := glue.InitRouter(
		glue.Services{
			Auth:         authService,
			Restaurant:   restaurantService,
			Menu:         menuService,
			Admin:        adminService,
			Payment:      paymentService,
			Webhook:      webhookService,
			Staff:        staffService,
			Activity:     activityService,
			Analytics:    analyticsService,
			Subscription: subscriptionService,
		},
		authMiddleware,
	)

	return router
}
