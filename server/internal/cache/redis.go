package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379" // Default Redis address
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: os.Getenv("REDIS_PASSWORD"), // Empty string if not set
		DB:       0,
		// Add some reasonable timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	// Try to connect
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
		log.Println("Application will continue without caching")
		return err
	}

	log.Printf("Redis connected successfully to %s", redisHost)
	return nil
}

func SetStockCache(symbol string, price string, ttl time.Duration) {
	if RedisClient == nil {
		return // Silently skip if Redis is not available
	}

	err := RedisClient.Set(ctx, symbol, price, ttl).Err()
	if err != nil {
		log.Printf("Failed to cache %s: %v", symbol, err)
	}
}

func GetStockCache(symbol string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("cache not available")
	}
	return RedisClient.Get(ctx, symbol).Result()
}
