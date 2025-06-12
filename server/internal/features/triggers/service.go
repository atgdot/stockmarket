package triggers

import (
	"context"
	"log"
	"sync"
	"time"

	"stockmarket/server/internal/database"
	"stockmarket/server/internal/models"
	ws "stockmarket/server/internal/websocket"
)

// Service handles all trigger-related operations
type Service struct {
	db         *database.Database
	ws         *ws.MarketWebSocket
	priceCache map[string]float64 // symbol:exchange -> price
	mu         sync.RWMutex
}

// NewService creates a new trigger service
func NewService(db *database.Database, ws *ws.MarketWebSocket) *Service {
	return &Service{
		db:         db,
		ws:         ws,
		priceCache: make(map[string]float64),
	}
}

// CreateTrigger creates a new trigger
func (s *Service) CreateTrigger(ctx context.Context, trigger *models.StockTrigger) error {
	// Create the trigger
	if err := s.db.CreateTrigger(ctx, trigger); err != nil {
		return err
	}

	// Update the stock's triggers list
	stock, err := s.db.GetStock(ctx, trigger.StockID)
	if err != nil {
		return err
	}

	stock.Triggers = append(stock.Triggers, trigger.TriggerID)
	return s.db.UpdateStock(ctx, stock)
}

// GetUserTriggers gets all triggers for a user
func (s *Service) GetUserTriggers(ctx context.Context, userID string) ([]*models.StockTrigger, error) {
	return s.db.GetTriggersByUser(ctx, userID)
}

// UpdatePrice updates the current price and evaluates triggers
func (s *Service) UpdatePrice(ctx context.Context, symbol, exchange string, price float64) error {
	// Update price cache
	key := symbol + ":" + exchange
	s.mu.Lock()
	s.priceCache[key] = price
	s.mu.Unlock()

	// Only evaluate triggers if market is open
	if !s.ws.IsMarketOpen(exchange) {
		return nil
	}

	// Get all triggers for this symbol
	triggers, err := s.db.GetTriggersBySymbol(ctx, symbol, exchange)
	if err != nil {
		return err
	}

	// Evaluate each trigger
	for _, trigger := range triggers {
		if !trigger.IsActive {
			continue
		}

		// Check cooldown period
		if time.Since(trigger.LastTrigger).Minutes() < float64(trigger.CooldownMinutes) {
			continue
		}

		// Get the associated stock to get the exchange
		stock, err := s.db.GetStock(ctx, trigger.StockID)
		if err != nil {
			log.Printf("Error getting stock for trigger %s: %v", trigger.TriggerID, err)
			continue
		}

		// Evaluate trigger conditions
		evaluation := s.evaluateTrigger(trigger, stock, price)
		if evaluation.Triggered {
			// Update last trigger time
			trigger.LastTrigger = time.Now()
			if err := s.db.UpdateTrigger(ctx, trigger); err != nil {
				log.Printf("Error updating trigger: %v", err)
				continue
			}

			// Send notification
			s.notifyTrigger(evaluation)
		}
	}

	return nil
}

// evaluateTrigger evaluates a single trigger against current price
func (s *Service) evaluateTrigger(trigger *models.StockTrigger, stock *models.Stock, price float64) TriggerEvaluation {
	evaluation := TriggerEvaluation{
		TriggerID:    trigger.TriggerID,
		UserID:       trigger.UserID,
		Symbol:       stock.Symbol,
		Exchange:     stock.Exchange,
		CurrentPrice: price,
		Timestamp:    time.Now(),
	}

	switch TriggerType(trigger.Type) {
	case PriceUpperLimit:
		if price >= trigger.PriceThreshold {
			evaluation.Triggered = true
			evaluation.Message = "Price exceeded upper limit"
		}
	case PriceLowerLimit:
		if price <= trigger.PriceThreshold {
			evaluation.Triggered = true
			evaluation.Message = "Price fell below lower limit"
		}
	case VolumeSpike:
		// This would need historical volume data
		if price > trigger.VolumeMultiplier {
			evaluation.Triggered = true
			evaluation.Message = "Unusual volume detected"
		}
	}

	return evaluation
}

// notifyTrigger sends notifications for triggered alerts
func (s *Service) notifyTrigger(evaluation TriggerEvaluation) {
	// Send websocket notification
	if err := s.ws.SendMessage(evaluation.UserID, evaluation); err != nil {
		log.Printf("Error sending websocket notification: %v", err)
	}

	// TODO: Send email notification
	// This would be implemented in a separate email service
}

// DeleteTrigger deletes a trigger
func (s *Service) DeleteTrigger(ctx context.Context, triggerID string) error {
	// Get the trigger to find its stock
	trigger, err := s.db.GetTrigger(ctx, triggerID)
	if err != nil {
		return err
	}

	// Remove trigger from stock's triggers list
	stock, err := s.db.GetStock(ctx, trigger.StockID)
	if err != nil {
		return err
	}

	// Remove trigger ID from stock's triggers
	for i, id := range stock.Triggers {
		if id == triggerID {
			stock.Triggers = append(stock.Triggers[:i], stock.Triggers[i+1:]...)
			break
		}
	}

	// Update stock
	if err := s.db.UpdateStock(ctx, stock); err != nil {
		return err
	}

	// Delete the trigger
	return s.db.DeleteTrigger(ctx, triggerID)
}
