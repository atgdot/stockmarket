package database

import (
	"context"
	"fmt"
	"time"

	"stockmarket/server/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

// CreateTrigger creates a new trigger in DynamoDB
func (db *Database) CreateTrigger(ctx context.Context, trigger *models.StockTrigger) error {
	trigger.TriggerID = uuid.New().String()
	trigger.CreatedAt = time.Now()
	trigger.UpdatedAt = time.Now()
	trigger.LastTrigger = time.Time{} // Zero time for new triggers

	item, err := attributevalue.MarshalMap(trigger)
	if err != nil {
		return err
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Triggers"),
		Item:      item,
	})
	return err
}

// GetTriggersBySymbol gets all triggers for a specific stock symbol
func (db *Database) GetTriggersBySymbol(ctx context.Context, symbol, exchange string) ([]*models.StockTrigger, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Triggers"),
		IndexName:              aws.String("SymbolIndex"),
		KeyConditionExpression: aws.String("symbol = :symbol AND exchange = :exchange"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":symbol":   &types.AttributeValueMemberS{Value: symbol},
			":exchange": &types.AttributeValueMemberS{Value: exchange},
		},
	}

	result, err := db.client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var triggers []*models.StockTrigger
	err = attributevalue.UnmarshalListOfMaps(result.Items, &triggers)
	return triggers, err
}

// GetTriggersByUser gets all triggers for a specific user
func (db *Database) GetTriggersByUser(ctx context.Context, userID string) ([]*models.StockTrigger, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Triggers"),
		IndexName:              aws.String("UserIndex"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := db.client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var triggers []*models.StockTrigger
	err = attributevalue.UnmarshalListOfMaps(result.Items, &triggers)
	return triggers, err
}

// UpdateTrigger updates an existing trigger
func (db *Database) UpdateTrigger(ctx context.Context, trigger *models.StockTrigger) error {
	trigger.UpdatedAt = time.Now()

	item, err := attributevalue.MarshalMap(trigger)
	if err != nil {
		return err
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Triggers"),
		Item:      item,
	})
	return err
}

// DeleteTrigger deletes a trigger
func (db *Database) DeleteTrigger(ctx context.Context, triggerID string) error {
	_, err := db.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("Triggers"),
		Key: map[string]types.AttributeValue{
			"trigger_id": &types.AttributeValueMemberS{Value: triggerID},
		},
	})
	return err
}

// GetUserStockTriggers gets all triggers for a user's stocks
func (db *Database) GetUserStockTriggers(ctx context.Context, userID string) (*models.UserStockTriggers, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("UserStockTriggers"),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := db.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	var userTriggers models.UserStockTriggers
	err = attributevalue.UnmarshalMap(result.Item, &userTriggers)
	return &userTriggers, err
}

// // UpdateUserStockTriggers updates a user's stock triggers
// func (db *Database) UpdateUserStockTriggers(ctx context.Context, userTriggers *models.UserStockTriggers) error {
// 	userTriggers.LastUpdated = time.Now()

// 	item, err := attributevalue.MarshalMap(userTriggers)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
// 		TableName: aws.String("UserStockTriggers"),
// 		Item:      item,
// 	})
// 	return err
// }

// GetTrigger retrieves a trigger by its ID
func (db *Database) GetTrigger(ctx context.Context, triggerID string) (*models.StockTrigger, error) {
	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("Triggers"),
		Key: map[string]types.AttributeValue{
			"trigger_id": &types.AttributeValueMemberS{Value: triggerID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("trigger not found")
	}

	var trigger models.StockTrigger
	err = attributevalue.UnmarshalMap(result.Item, &trigger)
	if err != nil {
		return nil, err
	}

	return &trigger, nil
}
