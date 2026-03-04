package usecase

import (
	"gateway-payments/internal/domain/repository"
)

type DeletePaymentInput struct {
	ID string
}

type DeletePayment struct {
	Repo repository.PaymentRepository
}

func NewDeletePaymentUseCase(repo repository.PaymentRepository) *DeletePayment {
	return &DeletePayment{
		Repo: repo,
	}
}

func (dp *DeletePayment) Execute(input DeletePaymentInput) error {
	err := dp.Repo.Delete(input.ID)
	if err != nil {
		return err
	}
	return nil
}
