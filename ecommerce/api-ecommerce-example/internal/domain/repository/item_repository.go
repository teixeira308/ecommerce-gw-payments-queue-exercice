package repository

import "ecommerce-api/internal/domain/entity"

type ItemRepository interface {
	Save(*entity.Item) error
	FindByID(id string) (*entity.Item, error)
	FindAll(page, limit int) ([]*entity.Item, error)
	Delete(id string) error
}
