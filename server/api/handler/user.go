package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/atgdot/stockmarket/serverinternal/features/auth"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(c echo.Context) error {
	var req AuthPayload
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}
	err := auth.CreateUser(context.TODO(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "User created")
}

func Login(c echo.Context) error {
	var req AuthPayload
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}
	token, err := auth.ValidateUser(context.TODO(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
