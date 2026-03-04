package order

import (
	"context"
	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/internal/domain/event"
	"ecommerce-api/internal/domain/repository"
	"ecommerce-api/internal/infrastructure/broker"
	"ecommerce-api/internal/interface/dto"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type CreateOrder struct {
	OrderRepo repository.OrderRepository
	ItemRepo  repository.ItemRepository
	Broker    *broker.RabbitMQClient
}

func NewCreateOrderUseCase(orderRepo repository.OrderRepository, itemRepo repository.ItemRepository, broker *broker.RabbitMQClient) *CreateOrder {
	return &CreateOrder{
		OrderRepo: orderRepo,
		ItemRepo:  itemRepo,
		Broker:    broker,
	}
}

func (uc *CreateOrder) Execute(itemsRequest []dto.OrderItemRequest, method string) (*entity.Order, error) {
	orderItems := make([]*entity.OrderItem, len(itemsRequest))
	var total float64

	log.Print(orderItems)
	for i, itemReq := range itemsRequest {
		item, err := uc.ItemRepo.FindByID(itemReq.ItemID)
		if err != nil {
			return nil, fmt.Errorf("item with ID %s not found: %w", itemReq.ItemID, err)
		}
		if item == nil {
			return nil, fmt.Errorf("item with ID %s not found", itemReq.ItemID)
		}

		subtotal := float64(itemReq.Quantity) * item.Price
		total += subtotal

		orderItems[i] = &entity.OrderItem{
			ID:       uuid.NewString(),
			ItemID:   itemReq.ItemID,
			Quantity: itemReq.Quantity,
			Subtotal: subtotal,
		}
	}

	order := &entity.Order{
		ID:        uuid.NewString(),
		Items:     orderItems,
		Total:     total,
		Status:    "pending",
		Method:    method,
		CreatedAt: time.Now(),
	}
	log.Print("Pedido que chegou: ", order)

	for _, oi := range order.Items {
		oi.OrderID = order.ID
		oi.CreatedAt = order.CreatedAt
	}

	err := uc.OrderRepo.Save(order)
	if err != nil {
		return nil, err
	}

	// Publish payment.requested event
	paymentRequestedEvent := event.PaymentRequested{
		Event:       "payment.requested",
		OrderID:     order.ID,
		Amount:      order.Total,
		Currency:    "BRL", // Assuming BRL as currency as per spec
		RequestedAt: time.Now(),
	}

	log.Print("enviado para fila payments: ", paymentRequestedEvent)
	err = uc.Broker.Publish(context.Background(), "payments.exchange", "payment.requested", paymentRequestedEvent)
	if err != nil {
		// Log the error but don't fail the order creation, as per eventual consistency
		log.Printf("Error publishing payment.requested event for order %s: %v\n", order.ID, err)
	}

	return order, nil
}
