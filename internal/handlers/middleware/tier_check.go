package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"menuvista/internal/models"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TierCheckMiddleware struct {
	queries *persistence.Queries
	logger  *log.Logger
}

func NewTierCheckMiddleware(queries *persistence.Queries, logger *log.Logger) *TierCheckMiddleware {
	return &TierCheckMiddleware{
		queries: queries,
		logger:  logger,
	}
}

func (tm *TierCheckMiddleware) CheckMenuLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ownerID, exists := c.Get("owner_id")
		if !exists {
			if role, _ := c.Get("role"); role == string(models.RoleAdmin) {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Owner context required for tier check"})
			c.Abort()
			return
		}

		ownerIDUUID := ownerID.(uuid.UUID)
		features, err := tm.getOwnerFeatures(c, ownerIDUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if features.MaxMenuItems > 0 {
			totalItems, err := tm.countOwnerMenuItems(c, ownerIDUUID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count menu items"})
				c.Abort()
				return
			}

			if totalItems >= features.MaxMenuItems {
				c.JSON(http.StatusForbidden, gin.H{
					"error": fmt.Sprintf("Menu item limit reached (%d/%d). Upgrade your plan.", totalItems, features.MaxMenuItems),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func (tm *TierCheckMiddleware) getOwnerFeatures(ctx context.Context, ownerID uuid.UUID) (models.FeatureLimits, error) {
	ownerIDStr := ownerID.String()
	subRow, err := tm.queries.GetActiveSubscriptionByOwner(ctx, utils.ToUUID(&ownerIDStr))
	if err != nil {
		tm.logger.Printf("[TierCheck] Failed to get subscription: %v", err)
		return models.FeatureLimits{}, fmt.Errorf("failed to check subscription")
	}

	// Check status and expiration
	now := time.Now()
	isActive := false

	if string(subRow.Status) == string(models.SubscriptionStatusActive) {
		if subRow.CurrentPeriodEnd.Time.After(now) {
			isActive = true
		}
	} else if string(subRow.Status) == string(models.SubscriptionStatusTrialing) {
		if subRow.TrialEnd.Valid && subRow.TrialEnd.Time.After(now) {
			isActive = true
		}
	}

	if !isActive {
		return models.FeatureLimits{}, fmt.Errorf("subscription is inactive or expired")
	}

	var features models.FeatureLimits
	if len(subRow.Features) > 0 {
		if err := json.Unmarshal(subRow.Features, &features); err != nil {
			tm.logger.Printf("[TierCheck] Failed to parse features: %v", err)
			return models.FeatureLimits{}, fmt.Errorf("failed to parse subscription features")
		}
	}
	return features, nil
}

func (tm *TierCheckMiddleware) countOwnerMenuItems(ctx context.Context, ownerID uuid.UUID) (int, error) {
	ownerIDStr := ownerID.String()
	restaurants, err := tm.queries.ListRestaurantsByOwner(ctx, utils.ToUUID(&ownerIDStr))
	if err != nil {
		return 0, err
	}

	var totalItems int
	for _, r := range restaurants {
		items, err := tm.queries.ListMenuItemsByRestaurant(ctx, r.ID)
		if err != nil {
			continue
		}
		totalItems += len(items)
	}
	return totalItems, nil
}

func (tm *TierCheckMiddleware) CheckStaffLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ownerID, exists := c.Get("owner_id")
		if !exists {
			if role, _ := c.Get("role"); role == string(models.RoleAdmin) {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Owner context required"})
			c.Abort()
			return
		}

		ownerIDUUID := ownerID.(uuid.UUID)
		features, err := tm.getOwnerFeatures(c, ownerIDUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if features.MaxStaffAccounts > 0 {

			staff, err := tm.queries.ListStaffByOwner(c, &ownerIDUUID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count staff"})
				c.Abort()
				return
			}

			if len(staff) >= features.MaxStaffAccounts {
				c.JSON(http.StatusForbidden, gin.H{
					"error": fmt.Sprintf("Staff limit reached (%d/%d). Upgrade your plan.", len(staff), features.MaxStaffAccounts),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
