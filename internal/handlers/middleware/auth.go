package middleware

import (
	"fmt"
	"log"
	"menuvista/internal/services/sms"
	"menuvista/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtSecret  []byte
	logger     *log.Logger
	smsService *sms.Service
}

func NewAuthMiddleware(secret string, logger *log.Logger, smsService *sms.Service) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:  []byte(secret),
		logger:     logger,
		smsService: smsService,
	}
}

	type AuthClaims struct {
		UserID       string `json:"user_id"`
		Role         string `json:"role"`
		OwnerID      string `json:"owner_id,omitempty"`
		RestaurantID string `json:"restaurant_id,omitempty"`
		SubStatus    string `json:"sub_status,omitempty"`
		SubEnd       int64  `json:"sub_end,omitempty"`
		jwt.RegisteredClaims
	}

func (am *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := am.extractToken(c)
		if err != nil {
			am.logger.Printf("[AuthMiddleware] Token extraction failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, err := am.parseAndValidateToken(tokenString)
		if err != nil {
			am.logger.Printf("[AuthMiddleware] Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		am.setContext(c, claims)
		c.Next()
	}
}

func (am *AuthMiddleware) extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("Authorization header must be Bearer {token}")
	}

	return parts[1], nil
}

func (am *AuthMiddleware) parseAndValidateToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return am.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (am *AuthMiddleware) setContext(c *gin.Context, claims *AuthClaims) {
	userID := utils.ParseUUID(claims.UserID)
	if userID == uuid.Nil {
		am.logger.Printf("[AuthMiddleware] Invalid user ID in token: %v", claims.UserID)
		return
	}

	c.Set("user_id", userID)
	c.Set("email", "")
	c.Set("role", claims.Role)
	if claims.Role == "owner" {
		c.Set("owner_id", userID)
	}
	if claims.OwnerID != "" {
		ownerID := utils.ParseUUID(claims.OwnerID)
		c.Set("owner_id", ownerID)
	}
	if claims.RestaurantID != "" {
		restaurantID := utils.ParseUUID(claims.RestaurantID)
		c.Set("restaurant_id", restaurantID)
	}

	if claims.SubStatus != "" {
		c.Set("sub_status", claims.SubStatus)
	}
	if claims.SubEnd > 0 {
		c.Set("sub_end", claims.SubEnd)

		// Check for expiry and notify/warn
		// Note: In a real app, we might want to rate limit this notification
		// or handle it more gracefully than just checking on every request.
		// For now, we'll just log it or set a flag.
		if time.Now().Unix() > claims.SubEnd {
			c.Set("sub_expired", true)

			// Send notification if expired
			// In a real app, we should check if we already sent this recently (e.g. via Redis)
			// For now, we'll just log and attempt to send
			if am.smsService != nil {
				// We need the user's phone number here.
				// Since it's not in the token, we can't send SMS without a DB lookup or adding phone to token.
				// For now, let's just log the intent.
				am.logger.Printf("[AuthMiddleware] Subscription expired for user %s. Triggering notification flow.", claims.UserID)
			}
		}
	}

	c.Set("auth_claims", claims)
}

func (am *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			log.Printf("[RequireRole] Role not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role type"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		log.Printf("[RequireRole] Access denied: user role %v, required one of %v", userRole, roles)
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
		c.Abort()
	}
}
