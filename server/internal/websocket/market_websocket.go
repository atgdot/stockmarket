package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MarketHours represents the trading hours for different exchanges
type MarketHours struct {
	StartTime time.Time
	EndTime   time.Time
	Timezone  string
}

// MarketWebSocket handles websocket connections during market hours
type MarketWebSocket struct {
	clients     map[string]*websocket.Conn // userID -> connection
	mu          sync.RWMutex
	upgrader    websocket.Upgrader
	marketHours map[string]MarketHours // exchange -> hours
}

// NewMarketWebSocket creates a new market websocket handler
func NewMarketWebSocket() *MarketWebSocket {
	return &MarketWebSocket{
		clients: make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: Implement proper origin checking
			},
		},
		marketHours: map[string]MarketHours{
			"NYSE": {
				StartTime: time.Date(0, 0, 0, 9, 30, 0, 0, time.UTC),
				EndTime:   time.Date(0, 0, 0, 16, 0, 0, 0, time.UTC),
				Timezone:  "America/New_York",
			},
			"NASDAQ": {
				StartTime: time.Date(0, 0, 0, 9, 30, 0, 0, time.UTC),
				EndTime:   time.Date(0, 0, 0, 16, 0, 0, 0, time.UTC),
				Timezone:  "America/New_York",
			},
			// Add other exchanges as needed
		},
	}
}

// IsMarketOpen checks if the market is currently open for a given exchange
func (m *MarketWebSocket) IsMarketOpen(exchange string) bool {
	hours, ok := m.marketHours[exchange]
	if !ok {
		return false
	}

	now := time.Now().UTC()

	// Check if it's a weekday
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}

	// Get current time in exchange's timezone
	loc, err := time.LoadLocation(hours.Timezone)
	if err != nil {
		log.Printf("Error loading timezone: %v", err)
		return false
	}

	currentTime := now.In(loc)
	marketStart := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
		hours.StartTime.Hour(), hours.StartTime.Minute(), 0, 0, loc)
	marketEnd := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
		hours.EndTime.Hour(), hours.EndTime.Minute(), 0, 0, loc)

	return currentTime.After(marketStart) && currentTime.Before(marketEnd)
}

// HandleConnection handles a new websocket connection
func (m *MarketWebSocket) HandleConnection(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Register the connection
	m.mu.Lock()
	m.clients[userID] = conn
	m.mu.Unlock()

	// Handle disconnection
	defer func() {
		m.mu.Lock()
		delete(m.clients, userID)
		m.mu.Unlock()
		conn.Close()
	}()

	// Keep connection alive and handle messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
	}
}

// BroadcastMessage sends a message to all connected clients
func (m *MarketWebSocket) BroadcastMessage(message interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for userID, conn := range m.clients {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to user %s: %v", userID, err)
		}
	}
}

// SendMessage sends a message to a specific user
func (m *MarketWebSocket) SendMessage(userID string, message interface{}) error {
	m.mu.RLock()
	conn, ok := m.clients[userID]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("user %s not connected", userID)
	}

	return conn.WriteJSON(message)
}
