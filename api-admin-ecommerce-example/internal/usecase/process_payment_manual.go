package usecase

import (
	"context"
	"fmt"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/domain/service"
	"log"
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
	log.Printf("[ManualProcessUseCase] Manually processing payment: %s", paymentID)
	
	payment, err := uc.Repo.FindByID(ctx, paymentID)
	if err != nil {
		log.Printf("[ManualProcessUseCase] Error finding payment %s: %v", paymentID, err)
		return err
	}

	err = uc.PaymentGateway.ProcessPayment(
		payment.ID,
		payment.Amount,
		payment.Method,
		payment.OrderID,
	)

	if err != nil {
		log.Printf("[ManualProcessUseCase] Error processing payment %s with gateway: %v", paymentID, err)
		return fmt.Errorf("manual processing failed: %w", err)
	}

	log.Printf("[ManualProcessUseCase] Payment %s manually processed successfully", paymentID)
	return nil
}
