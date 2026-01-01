package rest

import (
	"log"
	"net/http"
	"strconv"

	"menuvista/internal/models"
	"menuvista/internal/services/admin"
	"menuvista/internal/services/restaurant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	service           *admin.Service
	restaurantService *restaurant.Service
}

func NewAdminHandler(service *admin.Service, restaurantService *restaurant.Service) *AdminHandler {
	return &AdminHandler{
		service:           service,
		restaurantService: restaurantService,
	}
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	log.Printf("[AdminHandler] GetStats request received")
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		log.Printf("[AdminHandler] GetStats service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}
	RespondSuccess(c, http.StatusOK, stats, nil)
}

func (h *AdminHandler) GetRecentLogs(c *gin.Context) {
	log.Printf("[AdminHandler] GetRecentLogs request received")
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.service.GetRecentLogs(c.Request.Context(), int32(limit))
	if err != nil {
		log.Printf("[AdminHandler] GetRecentLogs service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}
	RespondSuccess(c, http.StatusOK, logs, nil)
}

func (h *AdminHandler) GetRestaurantDetails(c *gin.Context) {
	log.Printf("[AdminHandler] GetRestaurantDetails request received")
	idStr := c.Param("restaurant_id")
	var id uuid.UUID
	if err := id.Scan(idStr); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid restaurant ID", "INVALID_ID")
		return
	}

	details, err := h.service.GetRestaurantDetails(c.Request.Context(), id)
	if err != nil {
		log.Printf("[AdminHandler] GetRestaurantDetails service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}
	RespondSuccess(c, http.StatusOK, details, nil)
}

func (h *AdminHandler) GetRestaurants(c *gin.Context) {
	log.Printf("[AdminHandler] GetRestaurants request received")
	pagination := ParsePaginationParams(c)
	filters, err := ParseFilterParams(c, "restaurants")
	if err != nil {
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_FILTER")
		return
	}

	restaurantFilters := models.RestaurantFilter{
		OwnerID: filters.OwnerID,
		Status:  filters.Status,
		Search:  filters.Search,
	}

	results, meta, err := h.restaurantService.ListRestaurantsWithFilters(c.Request.Context(), restaurantFilters, pagination)
	if err != nil {
		log.Printf("[AdminHandler] GetRestaurants service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, results, meta)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	log.Printf("[AdminHandler] ListUsers request received")
	pagination := ParsePaginationParams(c)
	filters, err := ParseFilterParams(c, "users")
	if err != nil {
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_FILTER")
		return
	}

	users, err := h.service.ListUsersWithFilters(c.Request.Context(), filters, pagination)
	if err != nil {
		log.Printf("[AdminHandler] ListUsers service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	// Calculate meta (placeholder for total records)
	meta := models.CalculateMeta(pagination.Page, pagination.PageSize, len(users))
	RespondSuccess(c, http.StatusOK, users, meta)
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	log.Printf("[AdminHandler] UpdateUserStatus request received")
	userIDStr := c.Param("user_id")
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AdminHandler] UpdateUserStatus bind error: %v", err)
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_INPUT")
		return
	}

	var userID uuid.UUID
	if err := userID.Scan(userIDStr); err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid user ID", "INVALID_ID")
		return
	}

	user, err := h.service.UpdateUserStatus(c.Request.Context(), userID, req.IsActive)
	if err != nil {
		log.Printf("[AdminHandler] UpdateUserStatus service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}
	RespondSuccess(c, http.StatusOK, user, nil)
}
