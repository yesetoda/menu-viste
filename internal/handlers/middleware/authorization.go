package middleware

import (
	"fmt"
	"log"
	"net/http"

	"menuvista/internal/models"
	"menuvista/internal/storage/persistence"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthorizationMiddleware struct {
	queries *persistence.Queries
	logger  *log.Logger
}

func NewAuthorizationMiddleware(queries *persistence.Queries, logger *log.Logger) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		queries: queries,
		logger:  logger,
	}
}

func (am *AuthorizationMiddleware) RequireRestaurantAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Admin has access to everything
		if userRole == string(models.RoleAdmin) {
			c.Next()
			return
		}

		// Get restaurant ID from param
		restaurantIDStr := c.Param("restaurant_id")
		if restaurantIDStr == "" {
			// Try "id" if "restaurant_id" is missing, assuming the resource ID is the restaurant ID
			// This depends on the route. For now let's be strict or allow configuration.
			// Or maybe we are accessing a sub-resource (e.g. /restaurants/:id/menu).
			// If we are accessing /menu-items/:id, we need to fetch the item to know the restaurant ID.
			// This middleware is best used on routes where restaurant_id is in the path.
			c.JSON(http.StatusBadRequest, gin.H{"error": "Restaurant ID is required"})
			c.Abort()
			return
		}

		targetRestaurantID := utils.ParseUUID(restaurantIDStr)
		if targetRestaurantID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			c.Abort()
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
			c.Abort()
			return
		}

		if userRole == string(models.RoleOwner) {
			// Check if owner owns this restaurant
			// We can check the token's owner_id if it's set, but that's the owner's ID, not the restaurant's.
			// Wait, if I am an owner, my ID is in the token.
			// I need to check if the restaurant belongs to me.
			// I can query the DB.

			// Optimization: If we trust the client to pass the correct restaurant_id, we still need to verify ownership.
			// We can use the persistence layer.

			restaurant, err := am.queries.GetRestaurantByID(c, targetRestaurantID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
				c.Abort()
				return
			}

			fmt.Println("Owner ID:", restaurant.OwnerID)
			fmt.Println("User ID:", userID)
			if restaurant.OwnerID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this restaurant"})
				c.Abort()
				return
			}

		} else if userRole == string(models.RoleStaff) {
			// Check if staff is assigned to this restaurant
			staffRestaurantIDVal, exists := c.Get("restaurant_id")
			if !exists {
				c.JSON(http.StatusForbidden, gin.H{"error": "Staff not assigned to any restaurant"})
				c.Abort()
				return
			}

			staffRestaurantID, ok := staffRestaurantIDVal.(uuid.UUID)
			if !ok || staffRestaurantID != targetRestaurantID {
				c.JSON(http.StatusForbidden, gin.H{"error": "You are not assigned to this restaurant"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
