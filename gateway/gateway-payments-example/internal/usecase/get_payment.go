package usecase

import (
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
	return &GetPayment{
		Repo: repo,
	}
}

func (gp *GetPayment) Execute(input GetPaymentInput) (*entity.Payment, error) {
	payment, err := gp.Repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
