package usecase

import (
	"context"
	"gateway-payments/internal/domain/repository"
)

type DeletePaymentInput struct {
	ID string
}

type DeletePayment struct {
	Repo repository.PaymentRepository
}

func NewDeletePaymentUseCase(repo repository.PaymentRepository) *DeletePayment {
	return &DeletePayment{Repo: repo}
}

func (uc *DeletePayment) Execute(ctx context.Context, input DeletePaymentInput) error {
	return uc.Repo.Delete(ctx, input.ID)
}
