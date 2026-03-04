package usecase

import (
	"context"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/domain/service"
	"log"

	"github.com/google/uuid"
)

type CreatePayment struct {
	Repo           repository.PaymentRepository
	paymentGateway service.PaymentGateway
}

func NewCreatePaymentUseCase(repo repository.PaymentRepository, paymentGateway service.PaymentGateway) *CreatePayment {
	return &CreatePayment{
		Repo:           repo,
		paymentGateway: paymentGateway,
	}
}

func (pc *CreatePayment) Execute(ctx context.Context, paymentRequested event.PaymentRequested) (*entity.Payment, error) {
	// Check idempotency
	existingPayment, err := pc.Repo.FindByOrderID(paymentRequested.OrderID)
	if err != nil && err.Error() != fmt.Sprintf("payment with order ID %s not found", paymentRequested.OrderID) {
		return nil, fmt.Errorf("error checking existing payment for order %s: %w", paymentRequested.OrderID, err)
	}

	if existingPayment != nil {
		fmt.Printf("Payment for order %s already exists. Status: %s\n", paymentRequested.OrderID, existingPayment.Status)

		if existingPayment.Status == entity.StatusPending {
			fmt.Println("Resending pending payment to gateway...")
			err := pc.paymentGateway.ProcessPayment(
				existingPayment.ID,
				existingPayment.Amount,
				existingPayment.Method,
				existingPayment.OrderID,
			)
			if err != nil {
				return nil, fmt.Errorf("error resending to gateway: %w", err)
			}
		}
		return existingPayment, nil
	}

	payment := entity.NewPayment(uuid.NewString(), paymentRequested.OrderID, paymentRequested.Amount, paymentRequested.Method)
	payment.Status = entity.StatusPending

	log.Printf("Saving initial payment: %+v", payment)
	err = pc.Repo.Save(payment)
	if err != nil {
		return nil, fmt.Errorf("error saving payment: %w", err)
	}

	// Always send to gateway
	err = pc.paymentGateway.ProcessPayment(
		payment.ID,
		payment.Amount,
		payment.Method,
		payment.OrderID,
	)
	if err != nil {
		return nil, fmt.Errorf("error processing payment with gateway: %w", err)
	}

	return payment, nil
}
