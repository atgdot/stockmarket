package triggers

import (
	"time"
)

// TriggerType represents different types of triggers
type TriggerType string

const (
	PriceUpperLimit    TriggerType = "PRICE_UPPER_LIMIT"
	PriceLowerLimit    TriggerType = "PRICE_LOWER_LIMIT"
	PriceChangePercent TriggerType = "PRICE_CHANGE_PERCENT"
	VolumeSpike        TriggerType = "VOLUME_SPIKE"
	RSI                TriggerType = "RSI"
	MACD               TriggerType = "MACD"
	BollingerBands     TriggerType = "BOLLINGER_BANDS"
	TimeBased          TriggerType = "TIME_BASED"
)

// TriggerEvaluation represents the result of evaluating a trigger
type TriggerEvaluation struct {
	TriggerID    string    `json:"trigger_id"`
	UserID       string    `json:"user_id"`
	Symbol       string    `json:"symbol"`
	Exchange     string    `json:"exchange"`
	Triggered    bool      `json:"triggered"`
	CurrentPrice float64   `json:"current_price"`
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message"`
}
