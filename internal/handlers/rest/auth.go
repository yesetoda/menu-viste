package rest

import (
	"fmt"
	"log"
	"net/http"

	"menuvista/internal/models"
	"menuvista/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	log.Printf("[AuthHandler] Register request received")
	var input models.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[AuthHandler] Register bind error: %v", err)
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_INPUT")
		return
	}

	response, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		log.Printf("[AuthHandler] Register service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	log.Printf("[AuthHandler] User registered: %v", response.User.ID)
	RespondSuccess(c, http.StatusCreated, response, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	log.Printf("[AuthHandler] Login request received")
	var input models.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[AuthHandler] Login bind error: %v", err)
		RespondError(c, http.StatusBadRequest, err.Error(), "INVALID_INPUT")
		return
	}

	response, err := h.service.Login(c.Request.Context(), input)
	fmt.Println("this is the input", input)
	fmt.Println("this is the response", response)
	if err != nil {
		if err == auth.ErrPaymentRequired {
			log.Printf("[AuthHandler] Payment required for user: %s", input.Email)
			RespondSuccess(c, http.StatusPaymentRequired, response, nil) // 402
			return
		}
		if err == auth.ErrSubscriptionInactive {
			log.Printf("[AuthHandler] Subscription inactive for user: %s", input.Email)
			RespondError(c, http.StatusForbidden, "Restaurant subscription is inactive. Please contact the owner.", "SUBSCRIPTION_INACTIVE") // 403
			return
		}

		log.Printf("[AuthHandler] Login service error: %v", err)
		RespondError(c, http.StatusUnauthorized, err.Error(), "UNAUTHORIZED")
		return
	}

	log.Printf("[AuthHandler] User logged in: %v", response.User.ID)
	RespondSuccess(c, http.StatusOK, response, nil)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	log.Printf("[AuthHandler] GetProfile request received")
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[AuthHandler] GetProfile service error: %v", err)
		RespondError(c, http.StatusInternalServerError, err.Error(), "INTERNAL_ERROR")
		return
	}

	RespondSuccess(c, http.StatusOK, user, nil)
}

func (h *AuthHandler) ActivateAccount(c *gin.Context) {
	log.Printf("[AuthHandler] ActivateAccount request received")
	token := c.Query("token")
	if token == "" {
		RespondError(c, http.StatusBadRequest, "Activation token is required", "MISSING_TOKEN")
		return
	}

	response, err := h.service.ActivateUser(c.Request.Context(), token)
	if err != nil {
		log.Printf("[AuthHandler] Activation failed: %v", err)
		RespondError(c, http.StatusBadRequest, err.Error(), "ACTIVATION_FAILED")
		return
	}

	RespondSuccess(c, http.StatusOK, response, nil)
}
