package models

import "time"

type User struct {
	UserID       string    `dynamodbav:"user_id"`
	Email        string    `dynamodbav:"email"` // partition key
	PasswordHash string    `dynamodbav:"password_hash"`
	CreatedAt    time.Time `dynamodbav:"created_at"`
}
