package usecase

import (
	"context"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/repository"
)

type GetPaymentInput struct {
	ID string
}

type GetPayment struct {
	Repo repository.PaymentRepository
}

func NewGetPaymentUseCase(repo repository.PaymentRepository) *GetPayment {
	return &GetPayment{Repo: repo}
}

func (uc *GetPayment) Execute(ctx context.Context, input GetPaymentInput) (*entity.Payment, error) {
	return uc.Repo.FindByID(ctx, input.ID)
}
