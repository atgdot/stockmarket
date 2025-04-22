package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/atgdot/stockmarket/server/internal/features/stock"
)

func FetchStock(c echo.Context) error {
	symbol := c.QueryParam("symbol")
	if symbol == "" {
		return c.JSON(http.StatusBadRequest, "Missing stock symbol")
	}

	data, err := stock.FetchStockPrice(symbol)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, data)
}
