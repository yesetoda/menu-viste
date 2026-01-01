package rest

import (
	"net/http"

	"menuvista/internal/services/staff"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *staff.Service
}

func NewUserHandler(service *staff.Service) *StaffHandler {
	return &StaffHandler{
		service: service,
	}
}

func (h *StaffHandler) UpdateUser(c *gin.Context) {
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
