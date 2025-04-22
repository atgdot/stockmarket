package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/atgdot/stockmarket/serverapi/handler"
	"github.com/atgdot/stockmarket/serverapi/middleware/auth"
)

func StartServer() {
	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public routes
	e.POST("/signup", handler.SignUp)
	e.POST("/login", handler.Login)

	// Protected routes
	api := e.Group("/api")
	api.Use(auth.JWTMiddleware())

	api.GET("/me", handler.Me)
	api.POST("/stock/add", handler.AddStock)
	api.GET("/stock/fetch", handler.FetchStock)

	e.Logger.Fatal(e.Start(":8080"))
}
