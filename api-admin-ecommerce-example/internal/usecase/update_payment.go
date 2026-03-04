package usecase

import (
	"context"
	"errors"
	"fmt"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/infrastructure/broker"
	"log"
	"time"
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
	return &UpdatePayment{
		Repo:   repo,
		Broker: broker,
	}
}

func (up *UpdatePayment) Execute(ctx context.Context, input UpdatePaymentInput) error {
	payment, err := up.Repo.FindByID(input.ID)
	if err != nil {
		return errors.New("payment not found")
	}

	payment.Status = input.Status

	log.Print(payment)
	err = up.Repo.Save(payment)
	if err != nil {
		return err
	}

	if payment.Status == "APPROVED" || payment.Status == "REJECTED" {
		paymentProcessedEvent := event.PaymentProcessed{
			Event:       "payment.processed",
			OrderID:     payment.OrderID,
			Status:      payment.Status,
			ProcessedAt: time.Now(),
		}

		// Publica na fila para que o ecommerce-api receba e atualize o pedido
		err = up.Broker.Publish(ctx, "payments.exchange", "payment.processed", paymentProcessedEvent)
		if err != nil {
			return fmt.Errorf("error publishing payment.processed event: %w", err)
		}
		fmt.Printf("Status do pagamento %s atualizado e enviado para a fila: %s\n", payment.ID, payment.Status)
	}

	return nil
}
