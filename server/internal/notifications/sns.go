package notifications

import (
	"context"
	"fmt"

	"stockmarket/server/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var snsClient *sns.Client

// InitSNS initializes the SNS client
func InitSNS(cfg *config.Config) error {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(aws.NewCredentialsCache(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     cfg.AWSAccessKeyID,
					SecretAccessKey: cfg.AWSSecretAccessKey,
				}, nil
			},
		))),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	snsClient = sns.NewFromConfig(awsCfg)
	return nil
}

// GetSNSClient returns the SNS client instance
func GetSNSClient() *sns.Client {
	return snsClient
}

// CreateTopic creates a new SNS topic
func CreateTopic(ctx context.Context, topicName string) (string, error) {
	result, err := snsClient.CreateTopic(ctx, &sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create topic: %v", err)
	}

	return *result.TopicArn, nil
}

// SubscribeEmail subscribes an email to a topic
func SubscribeEmail(ctx context.Context, topicArn, email string) (string, error) {
	result, err := snsClient.Subscribe(ctx, &sns.SubscribeInput{
		Protocol: aws.String("email"),
		TopicArn: aws.String(topicArn),
		Endpoint: aws.String(email),
	})
	if err != nil {
		return "", fmt.Errorf("failed to subscribe email: %v", err)
	}

	return *result.SubscriptionArn, nil
}

// PublishMessage publishes a message to a topic
func PublishMessage(ctx context.Context, topicArn, subject, message string) error {
	_, err := snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Subject:  aws.String(subject),
		Message:  aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

// Unsubscribe removes an email subscription
func Unsubscribe(ctx context.Context, subscriptionArn string) error {
	_, err := snsClient.Unsubscribe(ctx, &sns.UnsubscribeInput{
		SubscriptionArn: aws.String(subscriptionArn),
	})
	if err != nil {
		return fmt.Errorf("failed to unsubscribe: %v", err)
	}

	return nil
}

// ListSubscriptions lists all subscriptions for a topic
func ListSubscriptions(ctx context.Context, topicArn string) ([]string, error) {
	result, err := snsClient.ListSubscriptionsByTopic(ctx, &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %v", err)
	}

	var subscriptions []string
	for _, sub := range result.Subscriptions {
		if sub.Endpoint != nil {
			subscriptions = append(subscriptions, *sub.Endpoint)
		}
	}

	return subscriptions, nil
}
