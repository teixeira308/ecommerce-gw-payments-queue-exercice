package dto

import (
	"gateway-payments/internal/domain/entity"
	"time"
)

type CreatePaymentRequest struct {
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
	OrderID string  `json:"order_id"`
}

type UpdatePaymentRequest struct {
	Status string `json:"status"`
}

type PaymentResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func CreatePaymentResponse(payment *entity.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:        payment.ID,
		OrderID:   payment.OrderID,
		Method:    payment.Method,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
	}
}
