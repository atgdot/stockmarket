package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"stockmarket/server/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var db *dynamodb.Client

// InitDynamoDB initializes the DynamoDB client and ensures table exists
func InitDynamoDB() error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	db = dynamodb.NewFromConfig(cfg)

	// Ensure table exists
	if err := ensureTableExists(); err != nil {
		return fmt.Errorf("failed to ensure table exists: %v", err)
	}

	return nil
}

// ensureTableExists creates the Users and Stocks tables if they don't exist
func ensureTableExists() error {
	// Create Users table
	if err := ensureUsersTableExists(); err != nil {
		return fmt.Errorf("failed to ensure Users table exists: %v", err)
	}

	// Create Stocks table
	if err := ensureStocksTableExists(); err != nil {
		return fmt.Errorf("failed to ensure Stocks table exists: %v", err)
	}

	return nil
}

// ensureUsersTableExists creates the Users table if it doesn't exist
func ensureUsersTableExists() error {
	tableName := os.Getenv("USERS_TABLE")
	if tableName == "" {
		return fmt.Errorf("USERS_TABLE environment variable not set")
	}

	// Check if table exists
	_, err := db.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err == nil {
		// Table exists
		return nil
	}

	// Create table
	_, err = db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("email"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return fmt.Errorf("failed to create Users table: %v", err)
	}

	// Wait for table to be active
	waiter := dynamodb.NewTableExistsWaiter(db)
	err = waiter.Wait(context.Background(),
		&dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
		2*time.Minute)
	if err != nil {
		return fmt.Errorf("timeout waiting for Users table creation: %v", err)
	}

	return nil
}

// ensureStocksTableExists creates the Stocks table if it doesn't exist
func ensureStocksTableExists() error {
	tableName := os.Getenv("STOCKS_TABLE")
	if tableName == "" {
		return fmt.Errorf("STOCKS_TABLE environment variable not set")
	}

	// Check if table exists
	_, err := db.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err == nil {
		// Table exists
		return nil
	}

	// Create table with composite key (user_id as partition key, stock_id as sort key)
	_, err = db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("user_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("stock_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("user_id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("stock_id"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return fmt.Errorf("failed to create Stocks table: %v", err)
	}

	// Wait for table to be active
	waiter := dynamodb.NewTableExistsWaiter(db)
	err = waiter.Wait(context.Background(),
		&dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
		2*time.Minute)
	if err != nil {
		return fmt.Errorf("timeout waiting for Stocks table creation: %v", err)
	}

	return nil
}

func SaveUser(ctx context.Context, user models.User) error {
	tableName := os.Getenv("USERS_TABLE")
	if tableName == "" {
		fmt.Println("DEBUG: USERS_TABLE environment variable is not set")
		return fmt.Errorf("USERS_TABLE environment variable not set")
	}

	fmt.Printf("DEBUG: Using table name: %s\n", tableName)

	// Marshal the user object to DynamoDB attribute map
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		fmt.Printf("DEBUG: Failed to marshal user data: %+v\n", err)
		return fmt.Errorf("failed to marshal user data: %v", err)
	}
	fmt.Printf("DEBUG: Marshalled user data: %+v\n", av)

	// Perform the PutItem operation to save the user into DynamoDB
	output, err := db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		fmt.Printf("DEBUG: Failed to save user to DynamoDB: %+v\n", err)
		return fmt.Errorf("failed to save user to DynamoDB: %v", err)
	}

	fmt.Printf("DEBUG: Successfully saved user. DynamoDB response: %+v\n", output)
	return nil
}

// GetUserByEmail retrieves a user from DynamoDB based on the email.
func GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	key := map[string]types.AttributeValue{
		"email": &types.AttributeValueMemberS{Value: email},
	}

	result, err := db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("USERS_TABLE")),
		Key:       key,
	})
	if err != nil || result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, fmt.Errorf("failed to unmarshal result: %v", err)
	}

	return user, nil
}
