package models

import "time"

type User struct {
	UserID       string    `dynamodbav:"user_id"`
	Email        string    `dynamodbav:"email"` // partition key
	PasswordHash string    `dynamodbav:"password_hash"`
	CreatedAt    time.Time `dynamodbav:"created_at"`
	UpdatedAt    time.Time `dynamodbav:"updated_at"`

	// Notification preferences
	NotificationPreferences struct {
		Email     bool   `dynamodbav:"email"`
		WebSocket bool   `dynamodbav:"websocket"`
		SMS       bool   `dynamodbav:"sms"`
		Phone     string `dynamodbav:"phone,omitempty"`
	} `dynamodbav:"notification_preferences"`

	// Active triggers count
	ActiveTriggers int `dynamodbav:"active_triggers"`
}
