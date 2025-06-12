package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port string

	// JWT configuration
	JWTSecret string

	// AWS configuration
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	// Database configuration
	UsersTable    string
	StocksTable   string
	TriggersTable string

	// Redis configuration
	RedisHost     string
	RedisPassword string

	// API Keys
	TwelveDataAPIKey string

	// SNS configuration
	SNSTopicName string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Port:               getEnvOrDefault("PORT", "8080"),
		JWTSecret:          getEnvOrDefault("JWT_SECRET", ""),
		AWSRegion:          getEnvOrDefault("AWS_REGION", "ap-south-1"),
		AWSAccessKeyID:     getEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
		UsersTable:         getEnvOrDefault("USERS_TABLE", "Users"),
		StocksTable:        getEnvOrDefault("STOCKS_TABLE", "Stocks"),
		TriggersTable:      getEnvOrDefault("TRIGGERS_TABLE", "Triggers"),
		RedisHost:          getEnvOrDefault("REDIS_HOST", "localhost:6379"),
		RedisPassword:      getEnvOrDefault("REDIS_PASSWORD", ""),
		TwelveDataAPIKey:   getEnvOrDefault("TWELVEDATA_API_KEY", ""),
		SNSTopicName:       getEnvOrDefault("SNS_TOPIC_NAME", "stock-market-alerts"),
	}

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if all required configuration values are set
func (c *Config) validate() error {
	required := map[string]string{
		"JWT_SECRET":            c.JWTSecret,
		"AWS_ACCESS_KEY_ID":     c.AWSAccessKeyID,
		"AWS_SECRET_ACCESS_KEY": c.AWSSecretAccessKey,
		"TWELVEDATA_API_KEY":    c.TwelveDataAPIKey,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("required environment variable %s is not set", name)
		}
	}

	return nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetRedisConfig returns Redis configuration
func (c *Config) GetRedisConfig() map[string]interface{} {
	return map[string]interface{}{
		"host":     c.RedisHost,
		"password": c.RedisPassword,
	}
}

// GetAWSConfig returns AWS configuration
func (c *Config) GetAWSConfig() map[string]string {
	return map[string]string{
		"region":          c.AWSRegion,
		"accessKeyID":     c.AWSAccessKeyID,
		"secretAccessKey": c.AWSSecretAccessKey,
	}
}
