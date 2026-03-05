package usecase

import (
	"context"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/infrastructure/broker"
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
	return uc.Repo.Save(ctx, payment)
}
