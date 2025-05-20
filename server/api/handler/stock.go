package handler

import (
	"net/http"
	"stockmarket/server/internal/database"
	"stockmarket/server/internal/features/stock"
	"strings"

	"github.com/labstack/echo/v4"
)

// SearchStockRequest represents a stock search request
type SearchStockRequest struct {
	Query string `json:"query"`
}

// AddStockRequest represents a request to add a stock to user's portfolio
type AddStockRequest struct {
	Symbol string `json:"symbol"`
}

// SearchStock handles stock search requests
func SearchStock(c echo.Context) error {
	var req SearchStockRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if req.Query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	results, err := stock.SearchStocks(req.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, results)
}

// FetchStockDetails handles detailed stock information requests
func FetchStockDetails(c echo.Context) error {
	symbol := c.QueryParam("symbol")
	if symbol == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Stock symbol is required",
		})
	}

	details, err := stock.FetchStockDetails(strings.ToUpper(symbol))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, details)
}

// AddStock handles adding a stock to user's portfolio
func AddStock(c echo.Context) error {
	var req AddStockRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if req.Symbol == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Stock symbol is required",
		})
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Get("user").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	// Fetch stock details including current price
	details, err := stock.FetchStockDetails(strings.ToUpper(req.Symbol))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Add stock to user's portfolio with current price
	err = database.AddStockToUser(c.Request().Context(), userID, details)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to add stock to portfolio",
		})
	}

	return c.JSON(http.StatusOK, details)
}

// GetUserStocks handles retrieving all stocks in user's portfolio
func GetUserStocks(c echo.Context) error {
	userID := c.Get("user").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	stocks, err := database.GetUserStocks(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user stocks",
		})
	}

	return c.JSON(http.StatusOK, stocks)
}

// RemoveStock handles removing a stock from user's portfolio
func RemoveStock(c echo.Context) error {
	stockID := c.Param("stockId")
	if stockID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Stock ID is required",
		})
	}

	userID := c.Get("user").(string)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	err := database.RemoveStockFromUser(c.Request().Context(), userID, stockID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove stock",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Stock removed successfully",
	})
}
