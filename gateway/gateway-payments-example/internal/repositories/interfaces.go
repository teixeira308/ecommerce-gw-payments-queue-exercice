package repositories

import (
	"context"
	"gateway-payments/internal/domain/entity"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *entity.Payment) error
	FindByID(ctx context.Context, id string) (*entity.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*entity.Payment, error)
	FindAll(ctx context.Context, page, limit int) ([]*entity.Payment, error)
	Delete(ctx context.Context, id string) error
}
