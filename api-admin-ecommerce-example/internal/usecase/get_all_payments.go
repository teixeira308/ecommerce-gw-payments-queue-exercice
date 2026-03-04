package usecase

import (
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
	return &GetAllPayments{
		Repo: repo,
	}
}

func (gap *GetAllPayments) Execute(input GetAllPaymentsInput) (*GetAllPaymentsOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}

	payments, err := gap.Repo.FindAll(input.Page, input.Limit)
	if err != nil {
		return nil, err
	}

	return &GetAllPaymentsOutput{Payments: payments}, nil
}
