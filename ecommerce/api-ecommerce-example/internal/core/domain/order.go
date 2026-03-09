package domain

import (
	"errors"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidStatus = errors.New("invalid order status")
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusCanceled  OrderStatus = "canceled"
	StatusCompleted OrderStatus = "completed"
)

type Order struct {
	ID        string
	Items     []*OrderItem
	Total     float64
	Status    OrderStatus
	Method    string
	CreatedAt time.Time
}

type OrderItem struct {
	ID        string
	OrderID   string
	ItemID    string
	Quantity  int
	Subtotal  float64
	CreatedAt time.Time
}
