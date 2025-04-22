package auth

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")

			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, "Missing token")
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, "Invalid token")
			}

			claims := token.Claims.(jwt.MapClaims)
			c.Set("user", claims["email"])
			return next(c)
		}
	}
}
