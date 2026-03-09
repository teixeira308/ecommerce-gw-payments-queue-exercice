package service

import (
	"context"
	"ecommerce-api/internal/core/domain"
	"ecommerce-api/internal/core/port"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo port.OrderRepository
	itemRepo  port.ItemRepository // Assumindo que criaremos port.ItemRepository similarmente
	broker    port.MessageBroker
}

func NewOrderService(
	orderRepo port.OrderRepository,
	itemRepo port.ItemRepository,
	broker port.MessageBroker,
) port.OrderService {
	return &orderService{
		orderRepo: orderRepo,
		itemRepo:  itemRepo,
		broker:    broker,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, input port.CreateOrderInput) (*domain.Order, error) {
	var total float64
	orderID := uuid.NewString()
	orderItems := make([]*domain.OrderItem, 0, len(input.Items))

	for _, itemReq := range input.Items {
		item, err := s.itemRepo.FindByID(ctx, itemReq.ItemID)
		if err != nil {
			return nil, fmt.Errorf("item lookup failed: %w", err)
		}
		if item == nil {
			return nil, fmt.Errorf("item %s not found", itemReq.ItemID)
		}

		subtotal := float64(itemReq.Quantity) * item.Price
		total += subtotal

		orderItems = append(orderItems, &domain.OrderItem{
			ID:        uuid.NewString(),
			OrderID:   orderID,
			ItemID:    itemReq.ItemID,
			Quantity:  itemReq.Quantity,
			Subtotal:  subtotal,
			CreatedAt: time.Now(),
		})
	}

	order := &domain.Order{
		ID:        orderID,
		Items:     orderItems,
		Total:     total,
		Status:    domain.StatusPending,
		Method:    input.Method,
		CreatedAt: time.Now(),
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Notificar sistema de pagamento de forma assíncrona (Event Driven)
	go func() {
		event := map[string]interface{}{
			"order_id": order.ID,
			"amount":   order.Total,
			"method":   order.Method,
		}
		_ = s.broker.Publish(context.Background(), "payments.exchange", "payment.requested", event)
	}()

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}
	return order, nil
}

func (s *orderService) ListOrders(ctx context.Context, page, limit int) ([]*domain.Order, error) {
	return s.orderRepo.FindAll(ctx, page, limit)
}

func (s *orderService) HandlePaymentStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}

	order.Status = status
	return s.orderRepo.UpdateStatus(ctx, order)
}
