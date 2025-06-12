package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"stockmarket/server/internal/features/stock"
	"stockmarket/server/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

// AddStockToUser adds a stock to a user's portfolio
func AddStockToUser(ctx context.Context, userID string, stockDetails *stock.StockDetails) error {
	tableName := os.Getenv("STOCKS_TABLE")
	if tableName == "" {
		return fmt.Errorf("STOCKS_TABLE environment variable not set")
	}

	// Convert price from json.Number to float64
	price, err := stockDetails.Price.Float64()
	if err != nil {
		return fmt.Errorf("failed to parse price: %v", err)
	}

	// Create the stock entry
	stock := models.Stock{
		StockID:     uuid.New().String(),
		UserID:      userID,
		Symbol:      stockDetails.Symbol,
		Name:        stockDetails.Name,
		Exchange:    stockDetails.Exchange,
		Currency:    stockDetails.Currency,
		Price:       price,
		LastPrice:   price, // Initially same as current price
		AddedAt:     time.Now(),
		LastUpdated: time.Now(),
	}

	// Marshal the stock object
	av, err := attributevalue.MarshalMap(stock)
	if err != nil {
		return fmt.Errorf("failed to marshal stock data: %v", err)
	}

	// Save to DynamoDB
	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to save stock: %v", err)
	}

	return nil
}

// GetUserStocks retrieves symbols of all stocks for a given user
func GetUserStocks(ctx context.Context, userID string) ([]models.Stock, error) {
	tableName := os.Getenv("STOCKS_TABLE")
	if tableName == "" {
		return nil, fmt.Errorf("STOCKS_TABLE environment variable not set")
	}

	// Query DynamoDB for all stocks with the given userID
	result, err := db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks: %v", err)
	}

	var stocks []models.Stock
	err = attributevalue.UnmarshalListOfMaps(result.Items, &stocks)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stocks: %v", err)
	}

	// Update current prices for all stocks
	for i := range stocks {
		details, err := stock.FetchStockDetails(stocks[i].Symbol)
		if err != nil {
			log.Printf("Failed to fetch price for %s: %v", stocks[i].Symbol, err)
			continue
		}

		// Update price information
		price, _ := details.Price.Float64()
		stocks[i].LastPrice = stocks[i].Price
		stocks[i].Price = price
		stocks[i].LastUpdated = time.Now()

		// Update the stock in DynamoDB with new price using a background context
		go func(s models.Stock) {
			// Create a new background context for the update
			updateCtx := context.Background()

			av, err := attributevalue.MarshalMap(s)
			if err != nil {
				log.Printf("Failed to marshal updated stock: %v", err)
				return
			}

			_, err = db.PutItem(updateCtx, &dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item:      av,
			})
			if err != nil {
				log.Printf("Failed to update stock price in DB: %v", err)
			}
		}(stocks[i])
	}

	return stocks, nil
}

// RemoveStockFromUser removes a stock from a user's portfolio
func RemoveStockFromUser(ctx context.Context, userID, stockID string) error {
	tableName := os.Getenv("STOCKS_TABLE")
	if tableName == "" {
		return fmt.Errorf("STOCKS_TABLE environment variable not set")
	}

	_, err := db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"user_id":  &types.AttributeValueMemberS{Value: userID},
			"stock_id": &types.AttributeValueMemberS{Value: stockID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete stock: %v", err)
	}

	return nil
}

// GetAllUniqueStocks retrieves all unique stocks across all user portfolios
func GetAllUniqueStocks(ctx context.Context) ([]models.Stock, error) {
	tableName := os.Getenv("STOCKS_TABLE")
	if tableName == "" {
		return nil, fmt.Errorf("STOCKS_TABLE environment variable not set")
	}

	// Scan the entire table
	result, err := db.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan stocks table: %v", err)
	}

	// Use a map to track unique stocks by symbol
	uniqueStocks := make(map[string]models.Stock)

	// Unmarshal items and collect unique stocks
	for _, item := range result.Items {
		var stock models.Stock
		if err := attributevalue.UnmarshalMap(item, &stock); err != nil {
			log.Printf("Warning: Failed to unmarshal stock: %v", err)
			continue
		}
		// Only keep the most recently added stock for each symbol
		existing, exists := uniqueStocks[stock.Symbol]
		if !exists || stock.AddedAt.After(existing.AddedAt) {
			uniqueStocks[stock.Symbol] = stock
		}
	}

	// Convert map to slice
	var stocks []models.Stock
	for _, stock := range uniqueStocks {
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetStock retrieves a stock by its ID
func (db *Database) GetStock(ctx context.Context, stockID string) (*models.Stock, error) {
	tableName := os.Getenv("STOCKS_TABLE")
	if tableName == "" {
		return nil, fmt.Errorf("STOCKS_TABLE environment variable not set")
	}

	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"stock_id": &types.AttributeValueMemberS{Value: stockID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get stock: %v", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("stock not found")
	}

	var stock models.Stock
	err = attributevalue.UnmarshalMap(result.Item, &stock)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stock: %v", err)
	}

	return &stock, nil
}

// UpdateStock updates a stock in the database
func (db *Database) UpdateStock(ctx context.Context, stock *models.Stock) error {
	tableName := os.Getenv("STOCKS_TABLE")
	if tableName == "" {
		return fmt.Errorf("STOCKS_TABLE environment variable not set")
	}

	stock.LastUpdated = time.Now()

	av, err := attributevalue.MarshalMap(stock)
	if err != nil {
		return fmt.Errorf("failed to marshal stock: %v", err)
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to update stock: %v", err)
	}

	return nil
}
