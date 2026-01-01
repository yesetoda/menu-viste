package utils

import (
	"crypto/rand"
	"encoding/json"
	"math/big"

	"github.com/google/uuid"
)

// DerefString returns the string value of a pointer, or an empty string if nil
func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// DerefBool returns the bool value of a pointer, or false if nil
func DerefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// DerefInt32 returns the int32 value of a pointer, or 0 if nil
func DerefInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// DerefFloat64 returns the float64 value of a pointer, or 0.0 if nil
func DerefFloat64(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

// DerefUUID returns the uuid.UUID value of a pointer, or uuid.Nil if nil
func DerefUUID(u *uuid.UUID) uuid.UUID {
	if u == nil {
		return uuid.Nil
	}
	return *u
}

// ParseUUID parses a string into uuid.UUID
func ParseUUID(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return u
}

// UUIDToString converts uuid.UUID to string
func UUIDToString(u uuid.UUID) string {
	if u == uuid.Nil {
		return ""
	}
	return u.String()
}

// ToPGUUID is a legacy helper that now just returns the input if it's already uuid.UUID
// or converts string to uuid.UUID if needed.
func ToPGUUID(u interface{}) uuid.UUID {
	switch v := u.(type) {
	case uuid.UUID:
		return v
	case string:
		return ParseUUID(v)
	default:
		return uuid.Nil
	}
}

// FromPGUUID is a legacy helper that returns uuid.UUID for backward compatibility where needed
func FromPGUUID(pg uuid.UUID) uuid.UUID {
	return pg
}

// ToUUIDPtr converts uuid.UUID to *uuid.UUID
func ToUUIDPtr(u uuid.UUID) *uuid.UUID {
	if u == uuid.Nil {
		return nil
	}
	return &u
}

// UnmarshalJSON is a helper to unmarshal JSON bytes into a struct
func UnmarshalJSON(data []byte, v interface{}) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// GeneratePassword generates a random password
func GeneratePassword(length int) (string, error) {
	charset := lowercase + uppercase + digits + symbols

	password := make([]byte, length)
	for i := range password {
		// Generate random index
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}
