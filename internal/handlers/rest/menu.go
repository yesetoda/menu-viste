package rest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"menuvista/internal/models"
	"menuvista/internal/services/menu"
	"menuvista/internal/services/restaurant"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	errInvalidCategoryID = "Invalid category ID"
	errInvalidItemID     = "Invalid item ID"
	errUnauthorized      = "Unauthorized"
)

type MenuHandler struct {
	service           *menu.Service
	restaurantService *restaurant.Service
}

func NewMenuHandler(service *menu.Service, restaurantService *restaurant.Service) *MenuHandler {
	return &MenuHandler{
		service:           service,
		restaurantService: restaurantService,
	}
}

// Categories

// Categories

func (h *MenuHandler) CreateCategory(c *gin.Context) {
	log.Printf("[MenuHandler] CreateCategory request received")

	restaurantIDStr := c.Param("restaurant_id")
	if restaurantIDStr == "" {
		RespondError(c, http.StatusBadRequest, "Restaurant ID is required", "INVALID_INPUT")
		return
	}
	fmt.Println("Restaurant ID:", restaurantIDStr)
	restaurantID := utils.ParseUUID(restaurantIDStr)

	fmt.Println("Restaurant ID:", restaurantIDStr, restaurantID)
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[MenuHandler] CreateCategory bind error: %v", err)
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_INPUT")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	userID := userIDVal.(uuid.UUID)

	result, err := h.service.CreateCategory(c.Request.Context(), userID, restaurantID, req)
	if err != nil {
		log.Printf("[MenuHandler] CreateCategory service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[MenuHandler] Category created: %v", result.ID)
	RespondSuccess(c, http.StatusCreated, result, nil)
}

func (h *MenuHandler) ListCategories(c *gin.Context) {
	log.Printf("[MenuHandler] ListCategories request received")
	restaurantIDStr := c.Param("restaurant_id")
	pagination := ParsePaginationParams(c)
	if restaurantIDStr == "" {
		slug := c.Param("slug")
		if slug != "" {
			restaurant, err := h.restaurantService.GetRestaurantBySlug(c.Request.Context(), slug)
			if err != nil {
				RespondError(c, http.StatusNotFound, "Restaurant not found", "NOT_FOUND")
				return
			}
			restaurantIDStr = restaurant.ID.String()
		}
	}

	if restaurantIDStr == "" {
		RespondError(c, http.StatusBadRequest, "Restaurant ID is required", "INVALID_INPUT")
		return
	}

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "Invalid restaurant ID", "INVALID_INPUT")
		return
	}

	categories, meta, err := h.service.ListCategories(c.Request.Context(), restaurantID, pagination)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, categories, meta)
}

func (h *MenuHandler) UpdateCategory(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID := utils.ParseUUID(categoryIDStr)

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_INPUT")
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	result, err := h.service.UpdateCategory(c.Request.Context(), userID, categoryID, req)
	if err != nil {
		log.Printf("[MenuHandler] UpdateCategory service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[MenuHandler] Category updated: %v", result.ID)
	RespondSuccess(c, http.StatusOK, result, nil)
}

func (h *MenuHandler) DeleteCategory(c *gin.Context) {
	log.Printf("[MenuHandler] DeleteCategory request received")
	categoryIDStr := c.Param("category_id")
	categoryID := utils.ParseUUID(categoryIDStr)

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	err := h.service.DeleteCategory(c.Request.Context(), userID, categoryID)
	if err != nil {
		log.Printf("[MenuHandler] DeleteCategory service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[MenuHandler] Category deleted: %v", categoryID)
	RespondSuccess(c, http.StatusOK, gin.H{"message": "Category deleted"}, nil)
}

func (h *MenuHandler) ListItems(c *gin.Context) {
	log.Printf("[MenuHandler] ListItems request received")
	restaurantIDStr := c.Param("restaurant_id")
	if restaurantIDStr == "" {
		RespondError(c, http.StatusBadRequest, "Restaurant ID is required", "INVALID_INPUT")
		return
	}
	restaurantID := utils.ParseUUID(restaurantIDStr)

	categoryIDStr := c.Param("category_id")
	categoryID := utils.ParseUUID(categoryIDStr)

	pagination := ParsePaginationParams(c)

	items, meta, err := h.service.ListItems(c.Request.Context(), restaurantID, categoryID, pagination)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, items, meta)
}

// Items

func (h *MenuHandler) CreateItem(c *gin.Context) {
	log.Printf("[MenuHandler] CreateItem request received")
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		RespondError(c, http.StatusBadRequest, "Failed to parse form", "INVALID_INPUT")
		return
	}

	restaurantIDStr := c.PostForm("restaurant_id")
	restaurantID := utils.ParseUUID(restaurantIDStr)

	categoryIDStr := c.PostForm("category_id")
	categoryID := utils.ParseUUID(categoryIDStr)

	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	currency := c.PostForm("currency")
	isAvailableStr := c.PostForm("is_available")

	price, _ := strconv.ParseFloat(priceStr, 64)
	isAvailable, _ := strconv.ParseBool(isAvailableStr)

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	input := models.CreateMenuItemRequest{
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		Price:       price,
		Currency:    currency,
		IsAvailable: isAvailable,
		// TODO: Map other fields
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		// input.ImageFile = file
		// input.ImageName = header.Filename
		_ = header
	}

	result, err := h.service.CreateMenuItem(c.Request.Context(), userID, restaurantID, input)
	if err != nil {
		log.Printf("[MenuHandler] CreateItem service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[MenuHandler] Item created: %v", result.ID)
	RespondSuccess(c, http.StatusCreated, result, nil)
}

func (h *MenuHandler) UpdateItem(c *gin.Context) {
	itemIDStr := c.Param("item_id")
	itemID := utils.ParseUUID(itemIDStr)

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		RespondError(c, http.StatusBadRequest, "Failed to parse form", "INVALID_INPUT")
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	isAvailableStr := c.PostForm("is_available")

	price, _ := strconv.ParseFloat(priceStr, 64)
	isAvailable, _ := strconv.ParseBool(isAvailableStr)

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	input := models.UpdateMenuItemRequest{
		Name:        &name,
		Description: &description,
		Price:       &price,
		IsAvailable: &isAvailable,
		// TODO: Map other fields
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		// input.ImageFile = file
		// input.ImageName = header.Filename
		_ = header
	}

	result, err := h.service.UpdateMenuItem(c.Request.Context(), userID, itemID, input)
	if err != nil {
		log.Printf("[MenuHandler] UpdateItem service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[MenuHandler] Item updated: %v", result.ID)
	RespondSuccess(c, http.StatusOK, result, nil)
}

func (h *MenuHandler) DeleteItem(c *gin.Context) {
	itemIDStr := c.Param("item_id")
	itemID := utils.ParseUUID(itemIDStr)

	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	err := h.service.DeleteMenuItem(c.Request.Context(), userID, itemID)
	if err != nil {
		log.Printf("[MenuHandler] DeleteItem service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[MenuHandler] Item deleted: %v", itemID)
	RespondSuccess(c, http.StatusOK, gin.H{"message": "Item deleted"}, nil)
}
