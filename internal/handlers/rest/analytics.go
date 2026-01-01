package rest

import (
	"net/http"
	"time"

	"menuvista/internal/models"
	"menuvista/internal/services/analytics"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service *analytics.Service
}

func NewAnalyticsHandler(service *analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

func (h *AnalyticsHandler) TrackEvent(c *gin.Context) {
	var req models.CreateAnalyticsEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Capture IP and User Agent if not provided
	if req.IPAddress == "" {
		req.IPAddress = c.ClientIP()
	}
	if req.Browser == "" {
		req.Browser = c.Request.UserAgent()
	}

	if err := h.service.TrackEvent(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *AnalyticsHandler) GetAggregates(c *gin.Context) {
	restaurantIDStr := c.Param("restaurant_id")

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, 0, -7) // Default last 7 days
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}

	aggregates, err := h.service.GetAggregates(c, utils.ParseUUID(restaurantIDStr), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aggregates})
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	restaurantIDStr := c.Param("restaurant_id")

	// Default to last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if startStr := c.Query("start_date"); startStr != "" {
		if parsed, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = parsed
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if parsed, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = parsed
		}
	}

	overview, err := h.service.GetOverview(c, utils.ParseUUID(restaurantIDStr), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": overview})
}
