package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/domain/repository"
	"log"
	"net/http"
	"time"
)

type UpdatePaymentInput struct {
	ID     string
	Status string
}

type UpdatePayment struct {
	Repo       repository.PaymentRepository
	WebhookURL string
}

func NewUpdatePaymentUseCase(repo repository.PaymentRepository, webhookURL string) *UpdatePayment {
	return &UpdatePayment{
		Repo:       repo,
		WebhookURL: webhookURL,
	}
}

func (up *UpdatePayment) Execute(ctx context.Context, input UpdatePaymentInput) error {
	payment, err := up.Repo.FindByID(input.ID)
	if err != nil {
		return errors.New("payment not found")
	}

	// Atualiza o status
	payment.Status = input.Status

	log.Printf("Updating payment: %+v", payment)
	err = up.Repo.Save(payment)
	if err != nil {
		return err
	}

	// Envia webhook após atualização (aprovação manual ou rejeição)
	go up.sendWebhook(payment)

	return nil
}

func (up *UpdatePayment) sendWebhook(payment *entity.Payment) {
	if up.WebhookURL == "" {
		log.Println("Webhook URL not configured, skipping...")
		return
	}

	url := fmt.Sprintf("%s/payments/%s", up.WebhookURL, payment.ID)

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
