package database

import (
	"context"
	"fmt"
	"os"

	"github.com/atgdot/stockmarket/serverinternal/features/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

var db *dynamodb.Client

func InitDynamoDB() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}
	db = dynamodb.NewFromConfig(cfg)
}

func SaveUser(ctx context.Context, user auth.User) error {
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("USERS_TABLE")),
		Item:      av,
	})
	return err
}

func GetUserByEmail(ctx context.Context, email string) (auth.User, error) {
	var user auth.User

	result, err := db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("USERS_TABLE")),
		Key: map[string]dynamodb.AttributeValue{
			"email": &dynamodb.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil || result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	return user, err
}
