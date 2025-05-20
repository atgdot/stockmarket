package stock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"stockmarket/server/internal/cache"
)

// StockData represents the structure of stock data
type StockData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

// StockSearchResult represents a single stock search result
type StockSearchResult struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Exchange    string `json:"exchange"`
	Currency    string `json:"currency"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// StockDetails represents detailed stock information
type StockDetails struct {
	Symbol        string      `json:"symbol"`
	Name          string      `json:"name"`
	Exchange      string      `json:"exchange"`
	Currency      string      `json:"currency"`
	Price         json.Number `json:"price"`
	Change        json.Number `json:"change"`
	ChangePercent json.Number `json:"change_percent"`
	High          json.Number `json:"high"`
	Low           json.Number `json:"low"`
	Volume        json.Number `json:"volume"`
	LastUpdated   time.Time   `json:"last_updated"`
}

// SetStockCache sets the cache with error suppression
func SetStockCache(symbol string, data string) {
	// Attempt to cache but ignore errors
	cache.SetStockCache(symbol, data, 15*time.Second)
}

// GetStockCache gets from cache with error suppression
func GetStockCache(symbol string) (string, bool) {
	data, err := cache.GetStockCache(symbol)
	if err != nil {
		return "", false
	}
	return data, true
}

// FetchStockPrice fetches the stock price for a given symbol
func FetchStockPrice(symbol string) (*StockData, error) {
	// Try fetching from cache first (silently fail if cache unavailable)
	if cached, ok := GetStockCache(symbol); ok {
		price := parsePrice(cached)
		return &StockData{Symbol: symbol, Price: price}, nil
	}

	// Get API key from environment
	apiKey := os.Getenv("TWELVEDATA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TWELVEDATA_API_KEY not set")
	}

	// Make API request
	url := fmt.Sprintf("https://api.twelvedata.com/price?symbol=%s&apikey=%s", symbol, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}
	defer resp.Body.Close()

	// Check if the API response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	// Unmarshal the response body into a map
	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %v", err)
	}

	// Extract the price from the result
	priceStr := result["price"]
	if priceStr == "" {
		return nil, fmt.Errorf("price not found in API response")
	}

	// Cache the result (silently fail if cache unavailable)
	SetStockCache(symbol, priceStr)

	// Return the stock data
	return &StockData{
		Symbol: symbol,
		Price:  parsePrice(priceStr),
	}, nil
}

// parsePrice converts the price string from the API into a float64
func parsePrice(priceStr string) float64 {
	var price float64
	fmt.Sscanf(priceStr, "%f", &price)
	return price
}

// SearchStocks searches for stocks matching the query
func SearchStocks(query string) ([]StockSearchResult, error) {
	apiKey := os.Getenv("TWELVEDATA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TWELVEDATA_API_KEY not set")
	}

	// Make API request to search endpoint
	url := fmt.Sprintf("https://api.twelvedata.com/symbol_search?symbol=%s&apikey=%s", query, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	var result struct {
		Data []StockSearchResult `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %v", err)
	}

	return result.Data, nil
}

// FetchStockDetails fetches detailed information for a specific stock
func FetchStockDetails(symbol string) (*StockDetails, error) {
	// Try cache first
	if cached, err := cache.GetStockCache(symbol); err == nil {
		var details StockDetails
		if err := json.Unmarshal([]byte(cached), &details); err == nil {
			return &details, nil
		}
	}

	apiKey := os.Getenv("TWELVEDATA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TWELVEDATA_API_KEY not set")
	}

	// Make API request for quote endpoint
	url := fmt.Sprintf("https://api.twelvedata.com/quote?symbol=%s&apikey=%s", symbol, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	// Print the raw response for debugging
	fmt.Printf("Raw API Response: %s\n", string(body))

	// First unmarshal into a map to debug the response structure
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw API response: %v", err)
	}

	// Print the parsed response for debugging
	fmt.Printf("Parsed Response: %+v\n", rawResponse)

	// Check if there's an error message in the response
	if errMsg, ok := rawResponse["error"].(string); ok {
		return nil, fmt.Errorf("API returned error: %s", errMsg)
	}

	// Create a temporary struct to match the API response structure
	type APIResponse struct {
		Symbol        string `json:"symbol"`
		Name          string `json:"name"`
		Exchange      string `json:"exchange"`
		Currency      string `json:"currency"`
		Close         string `json:"close"`
		Change        string `json:"change"`
		PercentChange string `json:"percent_change"`
		High          string `json:"high"`
		Low           string `json:"low"`
		Volume        string `json:"volume"`
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response into struct: %v", err)
	}

	// Create StockDetails from the API response
	details := &StockDetails{
		Symbol:        apiResp.Symbol,
		Name:          apiResp.Name,
		Exchange:      apiResp.Exchange,
		Currency:      apiResp.Currency,
		Price:         json.Number(apiResp.Close),
		Change:        json.Number(apiResp.Change),
		ChangePercent: json.Number(apiResp.PercentChange),
		High:          json.Number(apiResp.High),
		Low:           json.Number(apiResp.Low),
		Volume:        json.Number(apiResp.Volume),
		LastUpdated:   time.Now(),
	}

	// Cache the result
	if cached, err := json.Marshal(details); err == nil {
		cache.SetStockCache(symbol, string(cached), 15*time.Second)
	}

	return details, nil
}
