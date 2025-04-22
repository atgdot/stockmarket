package stock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/atgdot/stockmarket/serverinternal/cache"
)

type StockData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

func FetchStockPrice(symbol string) (*StockData, error) {
	// Try from cache first
	if cached, err := cache.GetStockCache(symbol); err == nil {
		price := parsePrice(cached)
		return &StockData{Symbol: symbol, Price: price}, nil
	}

	// If not cached, hit the API
	url := fmt.Sprintf("https://api.twelvedata.com/price?symbol=%s&apikey=%s", symbol, "YOUR_API_KEY")
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	priceStr := result["price"]
	cache.SetStockCache(symbol, priceStr, 15*time.Second) // Cache for 15s

	return &StockData{
		Symbol: symbol,
		Price:  parsePrice(priceStr),
	}, nil
}

func parsePrice(priceStr string) float64 {
	var price float64
	fmt.Sscanf(priceStr, "%f", &price)
	return price
}
