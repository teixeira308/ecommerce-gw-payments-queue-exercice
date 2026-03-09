package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/repositories"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo        repositories.PaymentRepository
	webhookURL  string
	autoApprove bool
}

func NewPaymentService(repo repositories.PaymentRepository, webhookURL string, autoApprove bool) *PaymentService {
	return &PaymentService{
		repo:        repo,
		webhookURL:  webhookURL,
		autoApprove: autoApprove,
	}
}

func (s *PaymentService) Create(ctx context.Context, id, orderID, method string, amount float64) (*entity.Payment, error) {
	existing, _ := s.repo.FindByOrderID(ctx, orderID)
	if existing != nil {
		return existing, nil
	}

	if id == "" {
		id = uuid.NewString()
	}
	payment := entity.NewPayment(id, orderID, method, amount)

	if s.autoApprove {
		payment.Approve()
	}

	if err := s.repo.Save(ctx, payment); err != nil {
		return nil, err
	}

	if payment.Status != entity.StatusPending {
		go s.sendWebhook(payment)
	}

	return payment, nil
}

func (s *PaymentService) GetByID(ctx context.Context, id string) (*entity.Payment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *PaymentService) List(ctx context.Context, page, limit int) ([]*entity.Payment, error) {
	return s.repo.FindAll(ctx, page, limit)
}

func (s *PaymentService) UpdateStatus(ctx context.Context, id, status string) error {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	p.Status = status
	if err := s.repo.Save(ctx, p); err != nil {
		return err
	}

	go s.sendWebhook(p)

	return nil
}

func (s *PaymentService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *PaymentService) sendWebhook(payment *entity.Payment) {
	if s.webhookURL == "" {
		return
	}

	url := fmt.Sprintf("%s/payments/%s", s.webhookURL, payment.ID)
	payload := map[string]interface{}{
		"event":        "payment.processed",
		"order_id":     payment.OrderID,
		"status":       payment.Status,
		"processed_at": time.Now().UTC(),
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Webhook Error] %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("[Webhook Error] received status %d from %s", resp.StatusCode, url)
	}
}
