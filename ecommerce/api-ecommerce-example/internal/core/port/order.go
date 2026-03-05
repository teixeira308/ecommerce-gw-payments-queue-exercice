package port

import (
	"context"
	"ecommerce-api/internal/core/domain"
)

type CreateOrderInput struct {
	Items  []OrderItemInput
	Method string
}

type OrderItemInput struct {
	ItemID   string
	Quantity int
}

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindAll(ctx context.Context, page, limit int) ([]*domain.Order, error)
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, order *domain.Order) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error)
	GetOrder(ctx context.Context, id string) (*domain.Order, error)
	ListOrders(ctx context.Context, page, limit int) ([]*domain.Order, error)
	HandlePaymentStatus(ctx context.Context, orderID string, status domain.OrderStatus) error
}

type MessageBroker interface {
	Publish(ctx context.Context, exchange, routingKey string, body interface{}) error
}
