package handler

import (
	"fmt"
	"net/http"

	"stockmarket/server/internal/features/auth"

	"github.com/labstack/echo/v4"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUp handles user registration
func SignUp(c echo.Context) error {
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
	}

	// Create user using request context
	err := auth.CreateUser(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		// Log the error but don't send it to client
		fmt.Printf("[SignUp] Failed to create user: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User created successfully",
	})
}

// Login handles user authentication
func Login(c echo.Context) error {
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
	}

	// Validate user using request context
	token, err := auth.ValidateUser(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		// Log the error but don't expose details
		fmt.Printf("[Login] Authentication failed: %v\n", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
