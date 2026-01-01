package email

import (
	"context"
	"fmt"
	"log"
	"strings"

	"menuvista/internal/models"
	"menuvista/internal/storage/persistence"

	"github.com/resend/resend-go/v2"
)

const senderEmail = "MenuVista <onboarding@resend.dev>"

type Service struct {
	client  *resend.Client
	queries *persistence.Queries
}

// NewService creates a new email service
func NewService(apiKey string, queries *persistence.Queries) *Service {
	client := resend.NewClient(apiKey)
	return &Service{
		client:  client,
		queries: queries,
	}
}

// SendAdminNotificationForNewOwner sends email to all admins when a new owner registers
func (s *Service) SendAdminNotificationForNewOwner(ctx context.Context, owner *persistence.User) error {
	log.Printf("[EmailService] Sending admin notification for new owner: %s", owner.Email)

	// Get all admin emails
	adminEmails, err := s.queries.GetAllAdminEmails(ctx)
	if err != nil {
		log.Printf("[EmailService] Failed to get admin emails: %v", err)
		return fmt.Errorf("failed to get admin emails: %w", err)
	}

	if len(adminEmails) == 0 {
		log.Printf("[EmailService] No active admins found to notify")
		return nil
	}

	// Prepare email content
	name := strings.Split(owner.FullName, " ")
	htmlContent := AdminNotificationNewOwnerTemplate(
		owner.Email,
		name[0],
		name[1],
	)

	// Send to all admins
	for _, adminEmail := range adminEmails {
		params := &resend.SendEmailRequest{
			From:    senderEmail,
			To:      []string{adminEmail},
			Subject: adminNotificationNewOwnerSubject,
			Html:    htmlContent,
		}

		_, err := s.client.Emails.Send(params)
		if err != nil {
			log.Printf("[EmailService] Failed to send email to admin %s: %v", adminEmail, err)
			// Continue sending to other admins even if one fails
			continue
		}
		log.Printf("[EmailService] Notification sent to admin: %s", adminEmail)
	}

	return nil
}

// SendOwnerApprovalEmail sends approval email to restaurant owner
func (s *Service) SendOwnerApprovalEmail(ctx context.Context, owner *persistence.User) error {
	log.Printf("[EmailService] Sending approval email to owner: %s", owner.Email)

	htmlContent := OwnerApprovalTemplate(owner.FullName)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{owner.Email},
		Subject: ownerApprovalSubject,
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send approval email: %v", err)
		return fmt.Errorf("failed to send approval email: %w", err)
	}

	log.Printf("[EmailService] Approval email sent successfully")
	return nil
}

// SendOwnerRejectionEmail sends rejection email to restaurant owner
func (s *Service) SendOwnerRejectionEmail(ctx context.Context, owner *persistence.User, reason string) error {
	log.Printf("[EmailService] Sending rejection email to owner: %s", owner.Email)

	htmlContent := OwnerRejectionTemplate(owner.FullName, reason)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{owner.Email},
		Subject: ownerRejectionSubject,
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send rejection email: %v", err)
		return fmt.Errorf("failed to send rejection email: %w", err)
	}

	log.Printf("[EmailService] Rejection email sent successfully")
	return nil
}

// SendRestaurantApprovalEmail sends approval email for restaurant
func (s *Service) SendRestaurantApprovalEmail(ctx context.Context, restaurant *persistence.Restaurant, owner *persistence.User) error {
	log.Printf("[EmailService] Sending restaurant approval email for: %s", restaurant.Name)

	htmlContent := RestaurantApprovalTemplate(owner.FullName, restaurant.Name)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{owner.Email},
		Subject: restaurantApprovalSubject,
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send restaurant approval email: %v", err)
		return fmt.Errorf("failed to send restaurant approval email: %w", err)
	}

	log.Printf("[EmailService] Restaurant approval email sent successfully")
	return nil
}

// SendRestaurantRejectionEmail sends rejection email for restaurant
func (s *Service) SendRestaurantRejectionEmail(ctx context.Context, restaurant *persistence.Restaurant, owner *persistence.User, reason string) error {
	log.Printf("[EmailService] Sending restaurant rejection email for: %s", restaurant.Name)

	htmlContent := RestaurantRejectionTemplate(owner.FullName, restaurant.Name, reason)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{owner.Email},
		Subject: restaurantRejectionSubject,
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send restaurant rejection email: %v", err)
		return fmt.Errorf("failed to send restaurant rejection email: %w", err)
	}

	log.Printf("[EmailService] Restaurant rejection email sent successfully")
	return nil
}

// SendWelcomeEmail sends a welcome email to new users
func (s *Service) SendWelcomeEmail(ctx context.Context, user *models.User) error {
	log.Printf("[EmailService] Sending welcome email to: %s", user.Email)

	// Simple welcome content for now
	htmlContent := fmt.Sprintf(`
		<h1>Welcome to MenuVista, %s!</h1>
		<p>We are excited to have you on board.</p>
		<p>Get started by creating your first restaurant.</p>
	`, user.FullName)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{user.Email},
		Subject: "Welcome to MenuVista!",
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send welcome email: %v", err)
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	log.Printf("[EmailService] Welcome email sent successfully")
	return nil
}

// SendPaymentSuccessEmail sends payment confirmation email
func (s *Service) SendPaymentSuccessEmail(ctx context.Context, email, firstName, invoiceNumber string, amount float64, currency string) error {
	log.Printf("[EmailService] Sending payment success email to: %s", email)

	htmlContent := PaymentSuccessTemplate(firstName, invoiceNumber, amount, currency)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{email},
		Subject: paymentSuccessSubject,
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send payment success email: %v", err)
		return fmt.Errorf("failed to send payment success email: %w", err)
	}

	log.Printf("[EmailService] Payment success email sent successfully")
	return nil
}

// SendPaymentFailedEmail sends payment failed email
func (s *Service) SendPaymentFailedEmail(ctx context.Context, email, firstName, updatePaymentURL string) error {
	log.Printf("[EmailService] Sending payment failed email to: %s", email)

	htmlContent := PaymentFailedTemplate(firstName, updatePaymentURL)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{email},
		Subject: paymentFailedSubject,
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send payment failed email: %v", err)
		return fmt.Errorf("failed to send payment failed email: %w", err)
	}

	log.Printf("[EmailService] Payment failed email sent successfully")
	return nil
}

// SendPaymentPendingEmail sends payment pending email
func (s *Service) SendPaymentPendingEmail(ctx context.Context, email, firstName string) error {
	log.Printf("[EmailService] Sending payment pending email to: %s", email)

	htmlContent := PaymentPendingTemplate(firstName)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{email},
		Subject: paymentPendingSubject,
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send payment pending email: %v", err)
		return fmt.Errorf("failed to send payment pending email: %w", err)
	}

	log.Printf("[EmailService] Payment pending email sent successfully")
	return nil
}

// SendStaffWelcomeEmail sends welcome email to new staff with credentials
func (s *Service) SendStaffWelcomeEmail(ctx context.Context, user *models.User, password string) error {
	log.Printf("[EmailService] Sending staff welcome email to: %s", user.Email)

	htmlContent := StaffWelcomeTemplate(user.FullName, user.Email, password)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{user.Email},
		Subject: "Welcome to MenuVista - Your Credentials",
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send staff welcome email: %v", err)
		return fmt.Errorf("failed to send staff welcome email: %w", err)
	}

	log.Printf("[EmailService] Staff welcome email sent successfully")
	return nil
}

// SendVerificationEmail sends verification email to new users
func (s *Service) SendVerificationEmail(ctx context.Context, user *models.User, token string) error {
	log.Printf("[EmailService] Sending verification email to: %s", user.Email)

	// Construct verification URL
	// TODO: Get base URL from config
	verificationURL := fmt.Sprintf("http://localhost:8080/api/v1/auth/activate?token=%s", token)
	expiresIn := "24 hours"

	htmlContent := VerificationEmailTemplate(user.FullName, verificationURL, expiresIn)

	params := &resend.SendEmailRequest{
		From:    senderEmail,
		To:      []string{user.Email},
		Subject: "Verify your email - MenuVista",
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[EmailService] Failed to send verification email: %v", err)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	log.Printf("[EmailService] Verification email sent successfully")
	return nil
}
