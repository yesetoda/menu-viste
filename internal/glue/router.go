package glue

import (
	"log"
	"menuvista/internal/handlers/middleware"
	"menuvista/internal/handlers/rest"
	"menuvista/internal/services/activity"
	"menuvista/internal/services/admin"
	"menuvista/internal/services/analytics"
	"menuvista/internal/services/auth"
	"menuvista/internal/services/menu"
	"menuvista/internal/services/payment"
	"menuvista/internal/services/restaurant"
	"menuvista/internal/services/staff"
	"menuvista/internal/services/subscription"
	"menuvista/internal/services/user"

	"github.com/gin-gonic/gin"
)

type router struct {
	Logger         *log.Logger
	AuthMiddleware *middleware.AuthMiddleware
}

type Services struct {
	Auth         *auth.Service
	Restaurant   *restaurant.Service
	Menu         *menu.Service
	Admin        *admin.Service
	Payment      *payment.Service
	Webhook      *payment.WebhookService
	Staff        *staff.Service
	UserService  *user.UserService
	Activity     *activity.Service
	Analytics    *analytics.Service
	Subscription *subscription.Service
}

func InitRouter(
	services Services,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	r := gin.New()

	// Global Middleware
	r.Use(middleware.LoggingMiddleware())
	r.Use(gin.Recovery())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"version": "1.0.0",
		})
	})

	// Initialize handlers
	authH := rest.NewAuthHandler(services.Auth)
	restH := rest.NewRestaurantHandler(services.Restaurant)
	menuH := rest.NewMenuHandler(services.Menu)
	adminH := rest.NewAdminHandler(services.Admin, services.Restaurant)
	paymentH := rest.NewPaymentHandler(services.Payment)
	webhookH := rest.NewWebhookHandler(services.Webhook)
	staffH := rest.NewStaffHandler(services.Staff)
	// UserH := rest.NewUserHandler(services.Staff)
	activityH := rest.NewActivityHandler(services.Activity)

	analyticsH := rest.NewAnalyticsHandler(services.Analytics)
	subH := rest.NewSubscriptionHandler(services.Subscription)

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Public Routes
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
			auth.GET("/activate", authH.ActivateAccount)

			protectedAuth := auth.Group("")
			protectedAuth.Use(authMiddleware.AuthMiddleware())
			{
				protectedAuth.GET("/me", authH.GetProfile)
			}
		}

		subscription := api.Group("/subscription")
		subscription.Use(authMiddleware.AuthMiddleware())
		{
			subscription.GET("/me", subH.GetSubscriptionDetails)
		}

		restaurants := api.Group("/restaurants")
		{
			restaurants.GET("", restH.ListRestaurants)
			restaurants.GET("/:slug", restH.GetRestaurant)
			restaurants.GET("/:slug/categories", menuH.ListCategories)
			restaurants.GET("/:slug/categories/:category_id/items", menuH.ListItems)
		}

		payments := api.Group("/payment")
		{
			payments.POST("/chapa/webhook", webhookH.ChapaWebhook)
			payments.POST("/chapa/test", webhookH.ChapaWebhookTest)
		}
	}

	// Payment Return Routes (Public)
	r.GET("/payment/success", paymentH.PaymentSuccess)
	r.GET("/payment/cancel", paymentH.PaymentCancel)

	// Protected Routes
	protected := api.Group("")
	protected.Use(authMiddleware.AuthMiddleware())
	{
		// Restaurant Owner Routes
		owner := protected.Group("")
		owner.Use(authMiddleware.RequireRole("owner"))
		{
			myRestaurants := owner.Group("/my-restaurants")
			{
				myRestaurants.POST("", restH.CreateRestaurant)
				myRestaurants.GET("", restH.ListMyRestaurants)
				myRestaurants.PATCH("/:restaurant_id", restH.UpdateRestaurant)
				myRestaurants.DELETE("/:restaurant_id", restH.DeleteRestaurant)
			}

			categories := owner.Group("/my-restaurants/:restaurant_id/categories")
			{
				categories.POST("", menuH.CreateCategory)
				categories.GET("", menuH.ListCategories)
				categories.PATCH("/:category_id", menuH.UpdateCategory)
				categories.DELETE("/:category_id", menuH.DeleteCategory)
			}

			items := owner.Group("/my-restaurants/:restaurant_id/categories/:category_id/items")
			{
				items.POST("", menuH.CreateItem)
				items.GET("", menuH.ListItems)
				items.PATCH("/:item_id", menuH.UpdateItem)
				items.DELETE("/:item_id", menuH.DeleteItem)
			}

			staff := owner.Group("/my-restaurants/:restaurant_id/staff")
			{
				staff.POST("", staffH.AddStaff)
				staff.GET("", staffH.ListStaff)
				staff.DELETE("/:staff_id", staffH.RemoveStaff)
			}

			analytics := owner.Group("/my-restaurants/:restaurant_id/analytics")
			{
				analytics.GET("/overview", analyticsH.GetOverview)
			}

			activity := owner.Group("/my-restaurants/:restaurant_id/activity")
			{
				activity.GET("", activityH.GetActivityLogs)
			}

			payment := owner.Group("/payment")
			{
				payment.POST("/initiate", paymentH.InitiatePayment)
				// payment.POST("/renew", paymentH.RenewSubscription)
				// payment.POST("/upgrade", paymentH.UpgradeSubscription)
			}
		}

		// Admin Routes
		admin := protected.Group("/admin")
		admin.Use(authMiddleware.RequireRole("admin"))
		{
			admin.GET("/stats", adminH.GetStats)
			admin.GET("/logs", adminH.GetRecentLogs)
			admin.GET("/restaurants", adminH.GetRestaurants)
			admin.GET("/restaurants/:restaurant_id", adminH.GetRestaurantDetails)

			admin.GET("/users", adminH.ListUsers)
			admin.PATCH("/users/:user_id/status", adminH.UpdateUserStatus)

			admin.PATCH("/restaurants/:restaurant_id/status", restH.UpdateRestaurantStatus)
		}
	}

	return r
}
