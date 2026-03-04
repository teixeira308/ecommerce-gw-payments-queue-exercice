package event

import "time"

type PaymentRequested struct {
	ID          string    `json:"id"`
	Event       string    `json:"event"`
	OrderID     string    `json:"order_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	RequestedAt time.Time `json:"requested_at"`
}
