package entity

import "time"

const (
	StatusPending  = "PENDING"
	StatusRejected = "REJECTED"
	StatusApproved = "APPROVED"
)

type Payment struct {
	ID        string
	OrderID   string
	Amount    float64 // Melhor usar float64 ou int (centavos) para cálculos
	Method    string
	Status    string
	CreatedAt time.Time
}

func NewPayment(id string, orderID string, amount float64, method string) *Payment {
	location := time.FixedZone("America/Sao_Paulo", -3*60*60)
	return &Payment{
		ID:        id,
		OrderID:   orderID,
		Amount:    amount,
		Method:    method,
		Status:    StatusPending,
		CreatedAt: time.Now().In(location),
	}
}

func UpdatePayment(id string, orderID string, amount float64, method string) *Payment {
	location := time.FixedZone("America/Sao_Paulo", -3*60*60)
	return &Payment{
		ID:        id,
		OrderID:   orderID,
		Amount:    amount,
		Method:    method,
		Status:    StatusPending,
		CreatedAt: time.Now().In(location),
	}
}
