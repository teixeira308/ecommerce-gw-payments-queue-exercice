package usecase

import (
	"context"
	"encoding/json"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/infrastructure/broker"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentRequestedConsumer struct {
	Broker        *broker.RabbitMQClient
	CreatePayment *CreatePayment
}

func NewPaymentRequestedConsumer(broker *broker.RabbitMQClient, createPayment *CreatePayment) *PaymentRequestedConsumer {
	return &PaymentRequestedConsumer{
		Broker:        broker,
		CreatePayment: createPayment,
	}
}

func (c *PaymentRequestedConsumer) StartConsuming(queueName, consumerName string) {
	err := c.Broker.Consume(queueName, consumerName, c.HandleMessage)
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}
	log.Printf("Started consuming messages from queue: %s", queueName)
}

func (c *PaymentRequestedConsumer) HandleMessage(d amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Received a message from queue %s: %s", d.RoutingKey, d.Body)

	var paymentRequestedEvent event.PaymentRequested
	if err := json.Unmarshal(d.Body, &paymentRequestedEvent); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		d.Nack(false, false) // Nack, don't requeue
		return
	}

	_, err := c.CreatePayment.Execute(ctx, paymentRequestedEvent)
	if err != nil {
		log.Printf("Error creating payment for order %s: %v", paymentRequestedEvent.OrderID, err)
		d.Nack(false, false) // Nack, don't requeue
		return
	}

	log.Printf("Payment for order %s processed and event published", paymentRequestedEvent.OrderID)
	d.Ack(false) // Ack, message processed successfully
}
