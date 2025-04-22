package cache

import (
	"context"
	"os"
	"time"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"), // e.g. "localhost:6379"
		Password: "",                      // Set if you use a Redis password
		DB:       0,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("Redis connected.")
}

func SetStockCache(symbol string, price string, ttl time.Duration) {
	err := RedisClient.Set(ctx, symbol, price, ttl).Err()
	if err != nil {
		log.Printf("Failed to cache %s: %v", symbol, err)
	}
}

func GetStockCache(symbol string) (string, error) {
	return RedisClient.Get(ctx, symbol).Result()
}
