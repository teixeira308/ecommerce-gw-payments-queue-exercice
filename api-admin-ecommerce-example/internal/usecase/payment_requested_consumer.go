package usecase

import (
	"context"
	"encoding/json"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/infrastructure/broker"
	"log"

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
		log.Fatalf("Failed to start consuming: %v", err)
	}
}

func (c *PaymentRequestedConsumer) HandleMessage(d amqp.Delivery) {
	var paymentRequested event.PaymentRequested
	err := json.Unmarshal(d.Body, &paymentRequested)
	if err != nil {
		log.Printf("Error unmarshaling payment requested: %v", err)
		d.Nack(false, false) // Rejects and does not requeue if payload is invalid
		return
	}

	log.Printf("Processing payment requested for order: %s", paymentRequested.OrderID)

	// We can use a background context or a specific context with timeout here
	_, err = c.CreatePayment.Execute(context.Background(), paymentRequested)
	if err != nil {
		log.Printf("Error creating payment: %v", err)
		d.Nack(false, true) // Requeues on processing error
		return
	}

	d.Ack(false)
}
