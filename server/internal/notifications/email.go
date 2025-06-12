package notifications

import (
	"context"
	"fmt"
	"time"

	"stockmarket/server/internal/config"
)

// EmailService handles email notifications
type EmailService struct {
	topicArn string
	config   *config.Config
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) (*EmailService, error) {
	if cfg.SNSTopicName == "" {
		return nil, fmt.Errorf("SNS topic name not configured")
	}

	// Create topic if it doesn't exist
	topicArn, err := CreateTopic(context.Background(), cfg.SNSTopicName)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %v", err)
	}

	return &EmailService{
		topicArn: topicArn,
		config:   cfg,
	}, nil
}

// SubscribeUser subscribes a user to email notifications
func (s *EmailService) SubscribeUser(ctx context.Context, email string) error {
	_, err := SubscribeEmail(ctx, s.topicArn, email)
	if err != nil {
		return fmt.Errorf("failed to subscribe user: %v", err)
	}
	return nil
}

// SendTriggerNotification sends a notification when a trigger is activated
func (s *EmailService) SendTriggerNotification(ctx context.Context, notification TriggerNotification) error {
	subject := fmt.Sprintf("Stock Alert: %s", notification.Symbol)
	message := fmt.Sprintf(
		"Your stock alert for %s has been triggered!\n\n"+
			"Symbol: %s\n"+
			"Current Price: $%.2f\n"+
			"Trigger Type: %s\n"+
			"Time: %s\n\n"+
			"This is an automated message from your Stock Market Alert System.",
		notification.Symbol,
		notification.Symbol,
		notification.Price,
		notification.TriggerType,
		time.Now().Format(time.RFC1123),
	)

	return PublishMessage(ctx, s.topicArn, subject, message)
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(ctx context.Context, notification WelcomeNotification) error {
	subject := "Welcome to Stock Market Alert System"
	message := fmt.Sprintf(
		"Welcome to the Stock Market Alert System!\n\n"+
			"Hello %s,\n\n"+
			"Thank you for joining our Stock Market Alert System. You can now:\n"+
			"- Set up stock price alerts\n"+
			"- Monitor volume spikes\n"+
			"- Track technical indicators\n"+
			"- Receive real-time notifications\n\n"+
			"Get started by setting up your first alert!\n\n"+
			"Best regards,\n"+
			"Stock Market Alert Team",
		notification.Username,
	)

	return PublishMessage(ctx, s.topicArn, subject, message)
}

// SendPriceAlert sends a price alert notification
func (s *EmailService) SendPriceAlert(ctx context.Context, symbol string, currentPrice, targetPrice float64, email string) error {
	subject := fmt.Sprintf("Price Alert: %s", symbol)
	message := fmt.Sprintf(
		"Price Alert for %s\n\n"+
			"Current Price: $%.2f\n"+
			"Target Price: $%.2f\n"+
			"Time: %s\n\n"+
			"This is an automated message from your Stock Market Alert System.",
		symbol,
		currentPrice,
		targetPrice,
		time.Now().Format(time.RFC1123),
	)

	return PublishMessage(ctx, s.topicArn, subject, message)
}

// SendVolumeAlert sends a volume alert notification
func (s *EmailService) SendVolumeAlert(ctx context.Context, symbol string, currentVolume, averageVolume float64, email string) error {
	subject := fmt.Sprintf("Volume Alert: %s", symbol)
	message := fmt.Sprintf(
		"Unusual Volume Alert for %s\n\n"+
			"Current Volume: %.0f\n"+
			"Average Volume: %.0f\n"+
			"Time: %s\n\n"+
			"This is an automated message from your Stock Market Alert System.",
		symbol,
		currentVolume,
		averageVolume,
		time.Now().Format(time.RFC1123),
	)

	return PublishMessage(ctx, s.topicArn, subject, message)
}
