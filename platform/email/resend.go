package email

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

type ResendClient struct {
	Client *resend.Client
	From   string
}

func NewResendClient() (*ResendClient, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("RESEND_FROM_EMAIL")

	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY is not set")
	}
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}

	client := resend.NewClient(apiKey)

	return &ResendClient{
		Client: client,
		From:   fromEmail,
	}, nil
}

func (r *ResendClient) SendEmail(to []string, subject string, htmlContent string) error {
	params := &resend.SendEmailRequest{
		From:    r.From,
		To:      to,
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := r.Client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
