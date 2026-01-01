package rest

import (
	"net/http"
	"strconv"

	"menuvista/internal/services/activity"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	service *activity.Service
}

func NewActivityHandler(service *activity.Service) *ActivityHandler {
	return &ActivityHandler{
		service: service,
	}
}

func (h *ActivityHandler) ListLogs(c *gin.Context) {
	restaurantIDStr := c.Param("restaurant_id")

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	logs, err := h.service.ListLogs(c, utils.ParseUUID(restaurantIDStr), int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func (h *ActivityHandler) GetActivityLogs(c *gin.Context) {
	h.ListLogs(c)
}
