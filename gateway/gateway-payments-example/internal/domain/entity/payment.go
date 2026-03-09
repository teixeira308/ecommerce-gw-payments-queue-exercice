package entity

import (
	"time"
)

const (
	StatusPending  = "PENDING"
	StatusRejected = "REJECTED"
	StatusApproved = "APPROVED"
)

type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPayment(id, orderID, method string, amount float64) *Payment {
	return &Payment{
		ID:        id,
		OrderID:   orderID,
		Amount:    amount,
		Method:    method,
		Status:    StatusPending,
		CreatedAt: time.Now().UTC(),
	}
}

func (p *Payment) Approve() {
	p.Status = StatusApproved
}

func (p *Payment) Reject() {
	p.Status = StatusRejected
}
