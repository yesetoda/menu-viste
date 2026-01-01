package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userID uuid.UUID, role string, ownerID *uuid.UUID, restaurantID *uuid.UUID, subStatus string, subEnd *time.Time) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.New().String()

	var err error
	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = UUIDToString(userID)
	atClaims["role"] = role
	if ownerID != nil {
		atClaims["owner_id"] = UUIDToString(*ownerID)
	}
	if restaurantID != nil {
		atClaims["restaurant_id"] = UUIDToString(*restaurantID)
	}
	if subStatus != "" {
		atClaims["sub_status"] = subStatus
	}
	if subEnd != nil {
		atClaims["sub_end"] = subEnd.Unix()
	}
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = UUIDToString(userID)
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return td, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("AUTH_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
