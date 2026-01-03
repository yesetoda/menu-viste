package rest

import (
	"log"
	"net/http"

	"menuvista/internal/models"
	"menuvista/internal/services/staff"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StaffHandler struct {
	service *staff.Service
}

func NewStaffHandler(service *staff.Service) *StaffHandler {
	return &StaffHandler{
		service: service,
	}
}

func (h *StaffHandler) ListStaff(c *gin.Context) {
	log.Printf("[StaffHandler] ListStaff request received")
	restaurantIDStr := c.Param("restaurant_id")
	if restaurantIDStr == "" {
		RespondError(c, http.StatusBadRequest, "Restaurant ID is required", "INVALID_INPUT")
		return
	}

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid restaurant ID", "INVALID_INPUT")
		return
	}

	pagination := ParsePaginationParams(c)

	staffList, meta, err := h.service.ListStaff(c.Request.Context(), restaurantID, pagination)
	if err != nil {
		log.Printf("[StaffHandler] ListStaff service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, staffList, meta)
}

func (h *StaffHandler) AddStaff(c *gin.Context) {
	var req models.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerID, exists := c.Get("owner_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req.OwnerID = ownerID.(uuid.UUID)

	user, err := h.service.CreateStaff(c, req.OwnerID, req.RestaurantID, models.CreateUserRequest{
		Email:    req.Email,
		FullName: req.FullName,
		Role:     models.RoleStaff,
		Phone:    req.Phone,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *StaffHandler) UpdateStaffStatus(c *gin.Context) {
	staffIDStr := c.Param("id")

	var req struct {
		IsActive     bool      `json:"is_active"`
		RestaurantID uuid.UUID `json:"restaurant_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStaffStatus(c, utils.ParseUUID(staffIDStr), req.RestaurantID, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

func (h *StaffHandler) RemoveStaff(c *gin.Context) {
	staffIDStr := c.Param("staff_id")

	// ownerID, _ := c.Get("owner_id") // Not used in service currently, but maybe should be for check

	restaurantIDStr := c.Param("restaurant_id")
	if restaurantIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id  param is required"})
		return
	}

	restaurantID := utils.ParseUUID(restaurantIDStr)

	if err := h.service.DeleteStaff(c, utils.ParseUUID(staffIDStr), restaurantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff removed"})
}
