package port

import (
	"context"
	"ecommerce-api/internal/core/domain"
)

type ItemRepository interface {
	Save(ctx context.Context, item *domain.Item) error
	FindByID(ctx context.Context, id string) (*domain.Item, error)
	FindAll(ctx context.Context, page, limit int) ([]*domain.Item, error)
	Delete(ctx context.Context, id string) error
}
