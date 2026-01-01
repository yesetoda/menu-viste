package rest

import (
	"strconv"

	"menuvista/internal/models"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
)

// ParsePaginationParams extracts pagination parameters from query string
func ParsePaginationParams(c *gin.Context) models.PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	return models.NewPaginationParams(page, pageSize)
}

// ParseFilterParams extracts filter parameters from query string and validates them
func ParseFilterParams(c *gin.Context, entity string) (models.FilterParams, error) {
	fb := utils.NewFilterBuilder(entity)
	return fb.ValidateAndParse(c)
}

// RespondSuccess sends a standardized success response
func RespondSuccess(c *gin.Context, statusCode int, data interface{}, meta *models.Meta) {
	response := models.SuccessResponse{
		Success:    true,
		StatusCode: statusCode,
		Data:       data,
		Meta:       meta,
	}
	c.JSON(statusCode, response)
}

// RespondError sends a standardized error response
func RespondError(c *gin.Context, statusCode int, message string, code string) {
	response := models.ErrorResponse{
		Success:    false,
		StatusCode: statusCode,
		Error:      message,
		Code:       code,
	}
	c.JSON(statusCode, response)
}
