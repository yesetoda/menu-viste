package sms

import (
	"context"
	"fmt"
	"log"
)

type Service struct {
	// In a real implementation, this would hold the SMS provider client (e.g., Twilio, AWS SNS)
}

func NewService() *Service {
	return &Service{}
}

// SendCredentials sends the staff credentials via SMS
func (s *Service) SendCredentials(ctx context.Context, phone, name, email, password string) error {
	// Mock implementation: Log the SMS
	message := fmt.Sprintf("Hi %s, you have been added to MenuVista. Your login email is %s and password is %s. Please change it after login.", name, email, password)
	log.Printf("[SMSService] Sending SMS to %s: %s", phone, message)
	return nil
}

// SendSubscriptionEnded sends a subscription ended notification
func (s *Service) SendSubscriptionEnded(ctx context.Context, phone, name string) error {
	message := fmt.Sprintf("Hi %s, your MenuVista subscription has ended. Please renew to continue using premium features.", name)
	log.Printf("[SMSService] Sending SMS to %s: %s", phone, message)
	return nil
}
