package event

import "time"

type PaymentProcessed struct {
	Event       string    `json:"event"`
	OrderID     string    `json:"order_id"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}
