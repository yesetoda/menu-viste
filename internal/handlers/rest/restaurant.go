package rest

import (
	"fmt"
	"log"
	"net/http"

	"menuvista/internal/models"
	"menuvista/internal/services/restaurant"
	"menuvista/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RestaurantHandler struct {
	service *restaurant.Service
}

func NewRestaurantHandler(service *restaurant.Service) *RestaurantHandler {
	return &RestaurantHandler{service: service}
}

func (h *RestaurantHandler) CreateRestaurant(c *gin.Context) {
	log.Printf("[RestaurantHandler] CreateRestaurant request received")
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		log.Printf("[RestaurantHandler] CreateRestaurant form parse error: %v", err)
		RespondError(c, http.StatusBadRequest, "Failed to parse form", "INVALID_INPUT")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	ownerID := userIDVal.(uuid.UUID)

	slug := c.PostForm("slug")
	name := c.PostForm("name")
	description := c.PostForm("description")

	input := models.CreateRestaurantRequest{
		Slug:        slug,
		Name:        name,
		Description: description,
		// TODO: Map other fields from form
	}

	// Handle Logo
	file, header, err := c.Request.FormFile("logo")
	if err == nil {
		defer file.Close()
		// input.LogoFile = file
		// input.LogoName = header.Filename
		// TODO: Handle file upload in service or here
		_ = header
	}

	// Handle Cover
	coverFile, coverHeader, err := c.Request.FormFile("cover")
	if err == nil {
		defer coverFile.Close()
		// input.CoverFile = coverFile
		// input.CoverName = coverHeader.Filename
		_ = coverHeader
	}

	result, err := h.service.CreateRestaurant(c.Request.Context(), ownerID, input)
	if err != nil {
		log.Printf("[RestaurantHandler] CreateRestaurant service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[RestaurantHandler] Restaurant created: %v", result.ID)
	RespondSuccess(c, http.StatusCreated, result, nil)
}

func (h *RestaurantHandler) GetRestaurant(c *gin.Context) {
	log.Printf("[RestaurantHandler] GetRestaurant request received")
	slug := c.Param("slug")
	result, err := h.service.GetRestaurantBySlug(c.Request.Context(), slug)
	if err != nil {
		log.Printf("[RestaurantHandler] GetRestaurant service error: %v", err)
		RespondError(c, http.StatusNotFound, "Restaurant not found", "NOT_FOUND")
		return
	}
	RespondSuccess(c, http.StatusOK, result, nil)
}

func (h *RestaurantHandler) ListMyRestaurants(c *gin.Context) {
	log.Printf("[RestaurantHandler] ListMyRestaurants request received")
	userIDStr, exists := c.Get("user_id")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}

	pagination := ParsePaginationParams(c)
	filters, err := ParseFilterParams(c, "restaurants")
	if err != nil {
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_FILTER")
		return
	}

	uidUUID := userIDStr.(uuid.UUID)
	uid := uidUUID.String()
	filters.OwnerID = &uid

	restaurantFilters := models.RestaurantFilter{
		OwnerID: filters.OwnerID,
		Status:  filters.Status,
		Search:  filters.Search,
	}

	results, meta, err := h.service.ListRestaurantsWithFilters(c.Request.Context(), restaurantFilters, pagination)
	if err != nil {
		log.Printf("[RestaurantHandler] ListMyRestaurants service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, results, meta)
}

func (h *RestaurantHandler) ListRestaurants(c *gin.Context) {
	log.Printf("[RestaurantHandler] ListRestaurants request received")
	pagination := ParsePaginationParams(c)
	fmt.Println("[RestaurantHandler] ListRestaurants pagination:", pagination)
	filters, err := ParseFilterParams(c, "restaurants")
	if err != nil {
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_FILTER")
		return
	}

	restaurantFilters := models.RestaurantFilter{
		OwnerID: filters.OwnerID,
		// Status:  filters.Status,
		Search: filters.Search,
	}

	results, meta, err := h.service.ListRestaurantsWithFilters(c.Request.Context(), restaurantFilters, pagination)
	if err != nil {
		log.Printf("[RestaurantHandler] ListRestaurants service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, results, meta)
}

func (h *RestaurantHandler) UpdateRestaurant(c *gin.Context) {
	log.Printf("[RestaurantHandler] UpdateRestaurant request received")
	restaurantIDStr := c.Param("restaurant_id")
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		RespondError(c, http.StatusBadRequest, "Failed to parse form", "INVALID_INPUT")
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	website := c.PostForm("website")
	address := c.PostForm("address")
	city := c.PostForm("city")
	country := c.PostForm("country")
	isPublished := c.PostForm("is_published") == "true"

	restaurantID := utils.ToUUID(&restaurantIDStr)

	userIDVal, exists := c.Get("user_id")
	if !exists {
		RespondError(c, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
		return
	}
	ownerID := userIDVal.(uuid.UUID)

	// Check if admin
	role, _ := c.Get("role")
	isAdmin := role == "admin"

	input := models.UpdateRestaurantRequest{
		Name:        &name,
		Description: &description,
		Email:       &email,
		Phone:       &phone,
		Website:     &website,
		Address:     &address,
		City:        &city,
		Country:     &country,
		IsPublished: &isPublished,
		// TODO: Map other fields
	}

	file, header, err := c.Request.FormFile("logo")
	if err == nil {
		defer file.Close()
		// input.LogoFile = file
		// input.LogoName = header.Filename
		_ = header
	}

	coverFile, coverHeader, err := c.Request.FormFile("cover")
	if err == nil {
		defer coverFile.Close()
		// input.CoverFile = coverFile
		// input.CoverName = coverHeader.Filename
		_ = coverHeader
	}

	result, err := h.service.UpdateRestaurant(c.Request.Context(), restaurantID, ownerID, isAdmin, input)
	if err != nil {
		log.Printf("[RestaurantHandler] UpdateRestaurant service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[RestaurantHandler] Restaurant updated: %v", result.ID)
	RespondSuccess(c, http.StatusOK, result, nil)
}

func (h *RestaurantHandler) DeleteRestaurant(c *gin.Context) {
	log.Printf("[RestaurantHandler] DeleteRestaurant request received")
	restaurantIDStr := c.Param("restaurant_id")
	restaurantID := utils.ParseUUID(restaurantIDStr)

	userIDVal, _ := c.Get("user_id")
	ownerID := userIDVal.(uuid.UUID)
	err := h.service.DeleteRestaurant(c.Request.Context(), restaurantID, ownerID)
	if err != nil {
		log.Printf("[RestaurantHandler] DeleteRestaurant service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[RestaurantHandler] Restaurant deleted: %v", restaurantID)
	RespondSuccess(c, http.StatusOK, gin.H{"message": "Restaurant deleted"}, nil)
}
