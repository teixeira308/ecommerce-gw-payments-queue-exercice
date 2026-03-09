package usecase

import (
	"context"
	"fmt"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/domain/service"
)

type ProcessPaymentManual struct {
	Repo           repository.PaymentRepository
	PaymentGateway service.PaymentGateway
}

func NewProcessPaymentManualUseCase(repo repository.PaymentRepository, paymentGateway service.PaymentGateway) *ProcessPaymentManual {
	return &ProcessPaymentManual{
		Repo:           repo,
		PaymentGateway: paymentGateway,
	}
}

func (uc *ProcessPaymentManual) Execute(ctx context.Context, paymentID string) error {
	payment, err := uc.Repo.FindByID(ctx, paymentID)
	if err != nil {
		return err
	}

	err = uc.PaymentGateway.ProcessPayment(
		payment.ID,
		payment.Amount,
		payment.Method,
		payment.OrderID,
	)

	if err != nil {
		return fmt.Errorf("manual processing failed: %w", err)
	}

	return nil
}
