package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, "Missing token")
			}

			// Check if the header starts with "Bearer "
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return c.JSON(http.StatusUnauthorized, "Invalid token format")
			}

			// Extract the token by removing the "Bearer " prefix
			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

			// Parse the token
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, "Invalid token")
			}

			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, "Token is not valid")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "Invalid token claims")
			}

			c.Set("user", claims["email"])
			return next(c)
		}
	}
}
