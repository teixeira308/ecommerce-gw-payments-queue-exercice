package usecase

import (
	"context"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/domain/service"
)

type ProcessPaymentManual struct {
	Repo           repository.PaymentRepository
	paymentGateway service.PaymentGateway
}

func NewProcessPaymentManualUseCase(repo repository.PaymentRepository, paymentGateway service.PaymentGateway) *ProcessPaymentManual {
	return &ProcessPaymentManual{
		Repo:           repo,
		paymentGateway: paymentGateway,
	}
}

func (pm *ProcessPaymentManual) Execute(ctx context.Context, paymentID string) error {
	payment, err := pm.Repo.FindByID(paymentID)
	if err != nil {
		return fmt.Errorf("payment with ID %s not found: %w", paymentID, err)
	}

	if payment.Status != entity.StatusPending {
		return fmt.Errorf("only pending payments can be sent to gateway, current status: %s", payment.Status)
	}

	err = pm.paymentGateway.ProcessPayment(
		payment.ID,
		payment.Amount,
		payment.Method,
		payment.OrderID,
	)
	if err != nil {
		return fmt.Errorf("error processing payment with gateway: %w", err)
	}

	return nil
}
