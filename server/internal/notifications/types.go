package notifications

// NotificationType represents different types of notifications
type NotificationType string

const (
	// NotificationTypeTrigger represents a stock trigger notification
	NotificationTypeTrigger NotificationType = "TRIGGER"
	// NotificationTypeWelcome represents a welcome notification
	NotificationTypeWelcome NotificationType = "WELCOME"
	// NotificationTypePriceAlert represents a price alert notification
	NotificationTypePriceAlert NotificationType = "PRICE_ALERT"
	// NotificationTypeVolumeAlert represents a volume alert notification
	NotificationTypeVolumeAlert NotificationType = "VOLUME_ALERT"
)

// Notification represents a notification message
type Notification struct {
	Type    NotificationType
	Email   string
	Subject string
	Message string
	Data    map[string]interface{}
}

// TriggerNotification represents a stock trigger notification
type TriggerNotification struct {
	Symbol      string
	Price       float64
	TriggerType string
	UserID      string
	Email       string
}

// WelcomeNotification represents a welcome notification
type WelcomeNotification struct {
	Email    string
	Username string
}
