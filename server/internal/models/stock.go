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
	Triggers    []string  `dynamodbav:"triggers"` // List of trigger IDs associated with this stock
}

// StockPrice represents the minimal stock information (symbol and price)
type StockPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

// StockTrigger represents a price trigger for a stock
type StockTrigger struct {
	TriggerID   string    `dynamodbav:"trigger_id"`
	StockID     string    `dynamodbav:"stock_id"` // Foreign key to Stock
	UserID      string    `dynamodbav:"user_id"`  // Foreign key to User
	Type        string    `dynamodbav:"type"`     // PRICE_UPPER, PRICE_LOWER, etc.
	IsActive    bool      `dynamodbav:"is_active"`
	CreatedAt   time.Time `dynamodbav:"created_at"`
	UpdatedAt   time.Time `dynamodbav:"updated_at"`
	LastTrigger time.Time `dynamodbav:"last_trigger"`

	// Trigger specific configurations
	PriceThreshold       float64  `dynamodbav:"price_threshold,omitempty"`
	VolumeMultiplier     float64  `dynamodbav:"volume_multiplier,omitempty"`
	NotificationChannels []string `dynamodbav:"notification_channels"`
	CooldownMinutes      int      `dynamodbav:"cooldown_minutes"`
}

// UserStockTriggers represents all triggers for a user's stocks
type UserStockTriggers struct {
	UserID    string         `dynamodbav:"user_id"`  // Partition key
	Triggers  []StockTrigger `dynamodbav:"triggers"` // List of triggers
	UpdatedAt time.Time      `dynamodbav:"updated_at"`
}
