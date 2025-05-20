package router

import (
	"stockmarket/server/api/handler"
	middleware "stockmarket/server/api/middleware"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func StartServer() {
	e := echo.New()

	// Global middleware
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORS())

	// Public routes
	e.POST("/signup", handler.SignUp)
	e.POST("/login", handler.Login)

	// Public stock routes
	e.POST("/api/stock/search", handler.SearchStock)
	e.GET("/api/stock/details", handler.FetchStockDetails)

	// Protected routes (require authentication)
	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware())

	// Stock management routes (fixed paths)
	api.POST("/stock/add", handler.AddStock)
	api.GET("/stock/list", handler.GetUserStocks)
	api.DELETE("/stock/:stockId", handler.RemoveStock)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
