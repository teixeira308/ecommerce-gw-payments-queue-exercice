package dto

import (
	"ecommerce-api/internal/domain/entity"
	"time"
)

type CreateItemRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type UpdateItemRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ItemResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func ToItemResponse(i *entity.Item) *ItemResponse {
	return &ItemResponse{
		ID:        i.ID,
		Name:      i.Name,
		Price:     i.Price,
		CreatedAt: i.CreatedAt,
	}
}

func CreateItemResponse(i *entity.Item) *ItemResponse {
	return &ItemResponse{
		ID:        i.ID,
		Name:      i.Name,
		Price:     i.Price,
		CreatedAt: i.CreatedAt,
	}
}

func UpdateItemResponse(i *entity.Item) *ItemResponse {
	return &ItemResponse{
		ID:        i.ID,
		Name:      i.Name,
		Price:     i.Price,
		CreatedAt: i.CreatedAt,
	}
}
