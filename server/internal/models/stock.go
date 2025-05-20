package models

import "time"

type Stock struct {
	StockID     string    `dynamodbav:"stock_id"`
	UserID      string    `dynamodbav:"user_id"` // Foreign key to User
	Symbol      string    `dynamodbav:"symbol"`
	Name        string    `dynamodbav:"name"`
	Exchange    string    `dynamodbav:"exchange"`
	Currency    string    `dynamodbav:"currency"`
	Price       float64   `dynamodbav:"price"`      // Current price
	LastPrice   float64   `dynamodbav:"last_price"` // Previous price for calculating change
	AddedAt     time.Time `dynamodbav:"added_at"`
	LastUpdated time.Time `dynamodbav:"last_updated"`
}

// StockPrice represents the minimal stock information (symbol and price)
type StockPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
