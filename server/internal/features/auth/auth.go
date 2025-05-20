package auth

import (
	"context"
	"fmt"
	"time"

	"stockmarket/server/internal/database"
	"stockmarket/server/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
}

// CreateUser creates a new user in the database.
func CreateUser(ctx context.Context, email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	user := models.User{
		UserID:       uuid.New().String(),
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	return database.SaveUser(ctx, user)
}

// ValidateUser validates user credentials and returns a JWT token if valid
func ValidateUser(ctx context.Context, email, password string) (string, error) {
	user, err := database.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := GenerateJWT(email)
	if err != nil {
		return "", err
	}
	return token, nil
}
