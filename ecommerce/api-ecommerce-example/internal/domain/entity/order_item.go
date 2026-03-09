package entity

import "time"

type OrderItem struct {
	ID        string
	OrderID   string
	ItemID    string
	Quantity  int
	Subtotal  float64
	CreatedAt time.Time
}