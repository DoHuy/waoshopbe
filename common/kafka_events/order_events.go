package kafka_events

type OrderEventPayload struct {
	EventID       string  `json:"event_id"`
	EventType     string  `json:"event_type"` // e.g., "ORDER_PAID", "ORDER_CREATED"
	OrderID       uint64  `json:"order_id"`
	OrderNumber   string  `json:"order_number"`
	CustomerEmail string  `json:"customer_email"`
	CustomerName  string  `json:"customer_name"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`
	Timestamp     int64   `json:"timestamp"`
}
