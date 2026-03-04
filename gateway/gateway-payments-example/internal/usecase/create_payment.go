package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/domain/repository"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreatePayment struct {
	Repo        repository.PaymentRepository
	WebhookURL  string
	AutoApprove bool
}

func NewCreatePaymentUseCase(repo repository.PaymentRepository, webhookURL string, autoApprove bool) *CreatePayment {
	return &CreatePayment{
		Repo:        repo,
		WebhookURL:  webhookURL,
		AutoApprove: autoApprove,
	}
}

func (pc *CreatePayment) Execute(ctx context.Context, paymentRequested event.PaymentRequested) (*entity.Payment, error) {
	// Check idempotency
	existingPayment, err := pc.Repo.FindByOrderID(paymentRequested.OrderID)
	if err != nil && err.Error() != fmt.Sprintf("payment with order ID %s not found", paymentRequested.OrderID) {
		return nil, fmt.Errorf("error checking existing payment for order %s: %w", paymentRequested.OrderID, err)
	}
	if existingPayment != nil {
		fmt.Printf("Payment for order %s already exists, ignoring. Status: %s\n", paymentRequested.OrderID, existingPayment.Status)
		return existingPayment, nil
	}

	status := entity.StatusPending
	if pc.AutoApprove {
		status = entity.StatusApproved
	}

	// Persist payment record
	id := paymentRequested.ID
	if id == "" {
		id = uuid.NewString()
	}
	payment := entity.NewPayment(id, paymentRequested.OrderID, paymentRequested.Amount, "Credit Card")
	payment.Status = status

	log.Printf("Creating payment: %+v", payment)
	err = pc.Repo.Save(payment)
	if err != nil {
		return nil, fmt.Errorf("error saving payment: %w", err)
	}

	if payment.Status == entity.StatusApproved {
		go pc.sendWebhook(payment)
	}

	return payment, nil
}

func (pc *CreatePayment) sendWebhook(payment *entity.Payment) {
	if pc.WebhookURL == "" {
		log.Println("Webhook URL not configured, skipping...")
		return
	}

	url := fmt.Sprintf("%s/payments/%s", pc.WebhookURL, payment.ID)

	payload := event.PaymentProcessed{
		Event:       "payment.processed",
		OrderID:     payment.OrderID,
		Status:      payment.Status,
		ProcessedAt: time.Now(),
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating webhook request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending webhook: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("Webhook returned error status: %d", resp.StatusCode)
	}
}
