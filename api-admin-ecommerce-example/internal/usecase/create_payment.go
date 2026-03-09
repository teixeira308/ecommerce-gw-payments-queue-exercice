package usecase

import (
	"context"
	"errors"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/domain/service"

	"github.com/google/uuid"
)

type CreatePayment struct {
	Repo              repository.PaymentRepository
	paymentGateway    service.PaymentGateway
	autoSendPayments  bool
}

func NewCreatePaymentUseCase(repo repository.PaymentRepository, paymentGateway service.PaymentGateway, autoSendPayments bool) *CreatePayment {
	return &CreatePayment{
		Repo:             repo,
		paymentGateway:   paymentGateway,
		autoSendPayments: autoSendPayments,
	}
}

func (pc *CreatePayment) Execute(ctx context.Context, paymentRequested event.PaymentRequested) (*entity.Payment, error) {
	// Check idempotency
	existingPayment, err := pc.Repo.FindByOrderID(ctx, paymentRequested.OrderID)
	if err != nil && !errors.Is(err, repository.ErrPaymentNotFound) {
		return nil, fmt.Errorf("error checking existing payment: %w", err)
	}

	if existingPayment != nil {
		if existingPayment.Status == entity.StatusPending && pc.autoSendPayments {
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

	err = pc.Repo.Save(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("error saving payment: %w", err)
	}

	if pc.autoSendPayments {
		err = pc.paymentGateway.ProcessPayment(
			payment.ID,
			payment.Amount,
			payment.Method,
			payment.OrderID,
		)
		if err != nil {
			return nil, fmt.Errorf("error processing with gateway: %w", err)
		}
	}

	return payment, nil
}
