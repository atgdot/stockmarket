package main

import (
	"log"

	"github.com/atgdot/stockmarket/serverinternal/config"
	"github.com/atgdot/stockmarket/serverinternal/database"
	"github.com/atgdot/stockmarket/serverinternal/cache"
	"github.com/atgdot/stockmarket/serverapi/router"
)

func main() {
	// Load environment configs
	config.LoadEnv()

	// Initialize DynamoDB and Redis
	database.InitDynamoDB()
	cache.InitRedis()

	go func() {
	symbols := []string{"TCS", "AAPL", "MSFT"}
	for {
		for _, s := range symbols {
			_, _ = stock.FetchStockPrice(s) // ignore errors for now
		}
		time.Sleep(10 * time.Second)
	}
}()


	// Start HTTP server
	log.Println("Starting Stock Tracker Server...")
	router.StartServer()
}
