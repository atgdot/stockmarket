package database

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Database represents the database client
type Database struct {
	client *dynamodb.Client
}

// NewDatabase creates a new database client
func NewDatabase(client *dynamodb.Client) *Database {
	return &Database{
		client: client,
	}
}
