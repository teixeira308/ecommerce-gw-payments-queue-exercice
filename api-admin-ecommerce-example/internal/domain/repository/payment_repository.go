package repository

import "gateway-payments/internal/domain/entity"

type PaymentRepository interface {
	Save(payment *entity.Payment) error
	FindByID(id string) (*entity.Payment, error)
	FindAll(page, limit int) ([]*entity.Payment, error)
	Delete(id string) error
	FindByOrderID(orderID string) (*entity.Payment, error)
}
