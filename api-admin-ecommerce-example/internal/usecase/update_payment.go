package usecase

import (
	"context"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/infrastructure/broker"
	"log"
)

type UpdatePaymentInput struct {
	ID     string
	Status string
}

type UpdatePayment struct {
	Repo   repository.PaymentRepository
	Broker *broker.RabbitMQClient
}

func NewUpdatePaymentUseCase(repo repository.PaymentRepository, broker *broker.RabbitMQClient) *UpdatePayment {
	return &UpdatePayment{Repo: repo, Broker: broker}
}

func (uc *UpdatePayment) Execute(ctx context.Context, input UpdatePaymentInput) error {
	payment, err := uc.Repo.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}

	payment.Status = input.Status
	if err := uc.Repo.Save(ctx, payment); err != nil {
		return err
	}

	// Publish event to RabbitMQ for Ecommerce API
	event := map[string]interface{}{
		"payment_id": payment.ID,
		"order_id":   payment.OrderID,
		"status":     payment.Status,
	}

	err = uc.Broker.Publish(ctx, "payments.exchange", "payment.processed", event)
	if err != nil {
		log.Printf("Failed to publish payment processed event: %v", err)
	}

	return nil
}
