package rest

import (
	"log"
	"net/http"

	"menuvista/internal/services/subscription"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service *subscription.Service
}

func NewSubscriptionHandler(service *subscription.Service) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) GetSubscriptionDetails(c *gin.Context) {
	log.Printf("[SubscriptionHandler] GetSubscriptionDetails request received")

	// Get user ID from context (set by auth middleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED")
		return
	}
	userID := userIDVal.(uuid.UUID)

	// Get role to determine if we need to fetch owner's subscription
	roleVal, _ := c.Get("role")
	role := roleVal.(string)

	var ownerID uuid.UUID
	if role == "owner" {
		ownerID = userID
	} else if role == "staff" {
		ownerIDVal, exists := c.Get("owner_id")
		if !exists || ownerIDVal == nil {
			RespondError(c, http.StatusBadRequest, "Staff member has no associated owner", "INVALID_STATE")
			return
		}
		// Handle pointer or value
		if ptr, ok := ownerIDVal.(*uuid.UUID); ok {
			ownerID = *ptr
		} else if val, ok := ownerIDVal.(uuid.UUID); ok {
			ownerID = val
		}
	} else {
		// Admin or other?
		RespondError(c, http.StatusForbidden, "Invalid role for subscription details", "FORBIDDEN")
		return
	}

	details, err := h.service.GetSubscriptionDetails(c.Request.Context(), ownerID)
	if err != nil {
		log.Printf("[SubscriptionHandler] Failed to get subscription details: %v", err)
		RespondError(c, http.StatusInternalServerError, "Failed to fetch subscription details", "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, details, nil)
}
