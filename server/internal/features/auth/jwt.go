package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}
