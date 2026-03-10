package order

import (
	"context"
	"ecommerce-api/internal/domain/event"
	"ecommerce-api/internal/infrastructure/broker"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentProcessedConsumer struct {
	Broker            *broker.RabbitMQClient
	UpdateOrderStatus *UpdateOrderStatus
}

func NewPaymentProcessedConsumer(broker *broker.RabbitMQClient, updateOrderStatus *UpdateOrderStatus) *PaymentProcessedConsumer {
	return &PaymentProcessedConsumer{
		Broker:            broker,
		UpdateOrderStatus: updateOrderStatus,
	}
}

func (c *PaymentProcessedConsumer) StartConsuming(queueName, consumerName string) {
	err := c.Broker.Consume(queueName, consumerName, c.HandleMessage)
	if err != nil {
		slog.Error("Failed to start consuming messages", "error", err)
		return
	}
	slog.Info("Started consuming messages", "queue", queueName)
}

func (c *PaymentProcessedConsumer) HandleMessage(d amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("Received a message", "routingKey", d.RoutingKey, "body", string(d.Body))

	var paymentProcessedEvent event.PaymentProcessed
	if err := json.Unmarshal(d.Body, &paymentProcessedEvent); err != nil {
		slog.Error("Error unmarshaling message", "error", err)
		d.Nack(false, false) // Nack, don't requeue
		return
	}

	var newStatus string
	if paymentProcessedEvent.Status == "APPROVED" {
		newStatus = "paid"
	} else if paymentProcessedEvent.Status == "REJECTED" {
		newStatus = "rejected"
	} else {
		slog.Warn("Unknown payment status received", "status", paymentProcessedEvent.Status, "orderID", paymentProcessedEvent.OrderID)
		d.Nack(false, false) // Nack, don't requeue
		return
	}

	err := c.UpdateOrderStatus.Execute(ctx, UpdateOrderStatusInput{
		ID:     paymentProcessedEvent.OrderID,
		Status: newStatus,
	})
	if err != nil {
		slog.Error("Error updating order status", "orderID", paymentProcessedEvent.OrderID, "error", err)
		d.Nack(false, false) // Nack, don't requeue
		return
	}

	slog.Info("Order status updated", "orderID", paymentProcessedEvent.OrderID, "status", newStatus)
	d.Ack(false) // Ack, message processed successfully
}
