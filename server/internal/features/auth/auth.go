package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/atgdot/stockmarket/serverinternal/database"
	"github.com/atgdot/stockmarket/serverinternal/features/auth/jwt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"-"`
	CreatedAt    string `json:"createdAt"`
}

func CreateUser(ctx context.Context, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Email:        email,
		PasswordHash: string(hashed),
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	return database.SaveUser(ctx, user)
}

func ValidateUser(ctx context.Context, email, password string) (string, error) {
	user, err := database.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := jwt.GenerateJWT(email)
	if err != nil {
		return "", err
	}
	return token, nil
}
