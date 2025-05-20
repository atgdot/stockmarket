package main

import (
	"context"
	"log"
	"os"
	"time"

	"stockmarket/server/api/router"
	"stockmarket/server/internal/cache"
	"stockmarket/server/internal/database"
	"stockmarket/server/internal/features/stock"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment configs
	err := godotenv.Load("C:/Users/canud/OneDrive/Desktop/HOM work/stockmarket/server/configs/secrets.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check if the environment variables are loaded correctly
	log.Println("USERS_TABLE:", os.Getenv("USERS_TABLE"))

	// After loading the environment file
	log.Printf("JWT_SECRET loaded: %v", os.Getenv("JWT_SECRET") != "")

	// Initialize DynamoDB
	if err := database.InitDynamoDB(); err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}

	// Initialize Redis (non-fatal if it fails)
	if err := cache.InitRedis(); err != nil {
		log.Printf("Note: Application will run without caching. Redis error: %v", err)
	}

	// Fetch stock prices in the background for user portfolios
	go func() {
		ctx := context.Background()
		for {
			// Get unique stocks from all user portfolios
			stocks, err := database.GetAllUniqueStocks(ctx)
			if err != nil {
				log.Printf("Failed to fetch user stocks: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			// If no stocks in any portfolio, wait and try again
			if len(stocks) == 0 {
				time.Sleep(10 * time.Second)
				continue
			}

			// Update prices for all stocks in portfolios
			for _, s := range stocks {
				_, err := stock.FetchStockPrice(s.Symbol)
				if err != nil {
					log.Printf("Failed to fetch price for %s: %v", s.Symbol, err)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// Start HTTP server
	log.Println("Starting Stock Tracker Server...")
	router.StartServer()
}
