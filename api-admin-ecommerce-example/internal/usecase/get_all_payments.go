package usecase

import (
	"context"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/repository"
)

type GetAllPaymentsInput struct {
	Page  int
	Limit int
}

type GetAllPaymentsOutput struct {
	Payments []*entity.Payment
}

type GetAllPayments struct {
	Repo repository.PaymentRepository
}

func NewGetAllPaymentsUseCase(repo repository.PaymentRepository) *GetAllPayments {
	return &GetAllPayments{Repo: repo}
}

func (uc *GetAllPayments) Execute(ctx context.Context, input GetAllPaymentsInput) (*GetAllPaymentsOutput, error) {
	payments, err := uc.Repo.FindAll(ctx, input.Page, input.Limit)
	if err != nil {
		return nil, err
	}

	return &GetAllPaymentsOutput{Payments: payments}, nil
}
