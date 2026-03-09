package order

import (
	"context"
	"ecommerce-api/internal/domain/event"
	"ecommerce-api/internal/infrastructure/broker"
	"encoding/json"
	"log"
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
		log.Fatalf("Failed to start consuming messages: %v", err)
	}
	log.Printf("Started consuming messages from queue: %s", queueName)
}

func (c *PaymentProcessedConsumer) HandleMessage(d amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Received a message from queue %s: %s", d.RoutingKey, d.Body)

	var paymentProcessedEvent event.PaymentProcessed
	if err := json.Unmarshal(d.Body, &paymentProcessedEvent); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		d.Nack(false, false) // Nack, don't requeue
		return
	}

	var newStatus string
	if paymentProcessedEvent.Status == "APPROVED" {
		newStatus = "paid"
	} else if paymentProcessedEvent.Status == "REJECTED" {
		newStatus = "rejected"
	} else {
		log.Printf("Unknown payment status received: %s for order %s", paymentProcessedEvent.Status, paymentProcessedEvent.OrderID)
		d.Nack(false, false) // Nack, don't requeue
		return
	}

	err := c.UpdateOrderStatus.Execute(ctx, UpdateOrderStatusInput{
		ID:     paymentProcessedEvent.OrderID,
		Status: newStatus,
	})
	if err != nil {
		log.Printf("Error updating order status for order %s: %v", paymentProcessedEvent.OrderID, err)
		d.Nack(false, false) // Nack, don't requeue
		return
	}

	log.Printf("Order %s status updated to %s", paymentProcessedEvent.OrderID, newStatus)
	d.Ack(false) // Ack, message processed successfully
}
