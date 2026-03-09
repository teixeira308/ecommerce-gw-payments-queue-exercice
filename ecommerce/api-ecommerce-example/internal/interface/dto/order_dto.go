// internal/interface/http/dto/order_dto.go

package dto

import "time"

type OrderItemRequest struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type CreateOrderRequest struct {
	Items  []OrderItemRequest `json:"items"`
	Method string             `json:"method"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

type OrderItemResponse struct {
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	Subtotal  float64   `json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateOrderResponse struct {
	ID        string              `json:"id"`
	Items     []OrderItemResponse `json:"items"`
	Total     float64             `json:"total"`
	Status    string              `json:"status"`
	Method    string              `json:"method"`
	CreatedAt time.Time           `json:"created_at"`
}

type GetOrderResponse struct {
	ID        string              `json:"id"`
	Items     []OrderItemResponse `json:"items"`
	Total     float64             `json:"total"`
	Status    string              `json:"status"`
	Method    string              `json:"method"`
	CreatedAt time.Time           `json:"created_at"`
}
