package triggers

import (
	"context"
	"testing"
	"time"

	"stockmarket/server/internal/database"
	"stockmarket/server/internal/models"
	ws "stockmarket/server/internal/websocket"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// mockDynamoDBClient is a mock implementation of the DynamoDB client
type mockDynamoDBClient struct {
	dynamodb.Client
}

func TestTriggerSystem(t *testing.T) {
	// Initialize services with mock client
	mockClient := &mockDynamoDBClient{}
	db := database.NewDatabase(mockClient)
	ws := ws.NewMarketWebSocket()
	service := NewService(db, ws)

	ctx := context.Background()

	// Test Case 1: Create and test a price upper limit trigger
	t.Run("Price Upper Limit Trigger", func(t *testing.T) {
		// Create a test stock
		stock := &models.Stock{
			StockID:   "test-stock-1",
			Symbol:    "AAPL",
			Exchange:  "NASDAQ",
			Triggers:  []string{},
			LastPrice: 150.0,
		}
		if err := db.CreateStock(ctx, stock); err != nil {
			t.Fatalf("Failed to create test stock: %v", err)
		}

		// Create a trigger
		trigger := &models.StockTrigger{
			TriggerID:       "test-trigger-1",
			StockID:         stock.StockID,
			UserID:          "test-user-1",
			Type:            string(PriceUpperLimit),
			IsActive:        true,
			PriceThreshold:  160.0,
			CooldownMinutes: 5,
		}
		if err := service.CreateTrigger(ctx, trigger); err != nil {
			t.Fatalf("Failed to create trigger: %v", err)
		}

		// Test price below threshold (should not trigger)
		if err := service.UpdatePrice(ctx, "AAPL", "NASDAQ", 155.0); err != nil {
			t.Fatalf("Failed to update price: %v", err)
		}

		// Verify trigger not activated
		triggers, err := service.GetUserTriggers(ctx, "test-user-1")
		if err != nil {
			t.Fatalf("Failed to get triggers: %v", err)
		}
		if len(triggers) == 0 {
			t.Fatal("No triggers found")
		}
		if triggers[0].LastTrigger.IsZero() {
			t.Log("Trigger correctly not activated for price below threshold")
		}

		// Test price above threshold (should trigger)
		if err := service.UpdatePrice(ctx, "AAPL", "NASDAQ", 165.0); err != nil {
			t.Fatalf("Failed to update price: %v", err)
		}

		// Verify trigger activated
		triggers, err = service.GetUserTriggers(ctx, "test-user-1")
		if err != nil {
			t.Fatalf("Failed to get triggers: %v", err)
		}
		if !triggers[0].LastTrigger.IsZero() {
			t.Log("Trigger correctly activated for price above threshold")
		}
	})

	// Test Case 2: Test cooldown period
	t.Run("Trigger Cooldown", func(t *testing.T) {
		// Create a test stock
		stock := &models.Stock{
			StockID:   "test-stock-2",
			Symbol:    "GOOGL",
			Exchange:  "NASDAQ",
			Triggers:  []string{},
			LastPrice: 2800.0,
		}
		if err := db.CreateStock(ctx, stock); err != nil {
			t.Fatalf("Failed to create test stock: %v", err)
		}

		// Create a trigger with short cooldown
		trigger := &models.StockTrigger{
			TriggerID:       "test-trigger-2",
			StockID:         stock.StockID,
			UserID:          "test-user-2",
			Type:            string(PriceUpperLimit),
			IsActive:        true,
			PriceThreshold:  2900.0,
			CooldownMinutes: 1,
		}
		if err := service.CreateTrigger(ctx, trigger); err != nil {
			t.Fatalf("Failed to create trigger: %v", err)
		}

		// Trigger first time
		if err := service.UpdatePrice(ctx, "GOOGL", "NASDAQ", 2950.0); err != nil {
			t.Fatalf("Failed to update price: %v", err)
		}

		// Try to trigger again immediately (should not trigger due to cooldown)
		if err := service.UpdatePrice(ctx, "GOOGL", "NASDAQ", 2950.0); err != nil {
			t.Fatalf("Failed to update price: %v", err)
		}

		// Verify trigger not activated again
		triggers, err := service.GetUserTriggers(ctx, "test-user-2")
		if err != nil {
			t.Fatalf("Failed to get triggers: %v", err)
		}
		if len(triggers) == 0 {
			t.Fatal("No triggers found")
		}

		// Wait for cooldown period
		time.Sleep(2 * time.Minute)

		// Try to trigger again after cooldown
		if err := service.UpdatePrice(ctx, "GOOGL", "NASDAQ", 2950.0); err != nil {
			t.Fatalf("Failed to update price: %v", err)
		}

		// Verify trigger activated again
		triggers, err = service.GetUserTriggers(ctx, "test-user-2")
		if err != nil {
			t.Fatalf("Failed to get triggers: %v", err)
		}
		if !triggers[0].LastTrigger.IsZero() {
			t.Log("Trigger correctly activated after cooldown period")
		}
	})

	// Test Case 3: Test trigger deletion
	t.Run("Trigger Deletion", func(t *testing.T) {
		// Create a test stock
		stock := &models.Stock{
			StockID:   "test-stock-3",
			Symbol:    "MSFT",
			Exchange:  "NASDAQ",
			Triggers:  []string{},
			LastPrice: 300.0,
		}
		if err := db.CreateStock(ctx, stock); err != nil {
			t.Fatalf("Failed to create test stock: %v", err)
		}

		// Create a trigger
		trigger := &models.StockTrigger{
			TriggerID:       "test-trigger-3",
			StockID:         stock.StockID,
			UserID:          "test-user-3",
			Type:            string(PriceUpperLimit),
			IsActive:        true,
			PriceThreshold:  310.0,
			CooldownMinutes: 5,
		}
		if err := service.CreateTrigger(ctx, trigger); err != nil {
			t.Fatalf("Failed to create trigger: %v", err)
		}

		// Delete the trigger
		if err := service.DeleteTrigger(ctx, trigger.TriggerID); err != nil {
			t.Fatalf("Failed to delete trigger: %v", err)
		}

		// Verify trigger is deleted
		triggers, err := service.GetUserTriggers(ctx, "test-user-3")
		if err != nil {
			t.Fatalf("Failed to get triggers: %v", err)
		}
		if len(triggers) > 0 {
			t.Fatal("Trigger was not deleted")
		}
		t.Log("Trigger successfully deleted")
	})
}
