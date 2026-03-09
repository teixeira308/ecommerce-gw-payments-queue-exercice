package order

import (
	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/internal/domain/repository"
)

type GetOrderInput struct {
	ID string
}

type GetOrder struct {
	Repo repository.OrderRepository
}

func NewGetOrderUseCase(repo repository.OrderRepository) *GetOrder {
	return &GetOrder{
		Repo: repo,
	}
}

func (gp *GetOrder) Execute(input GetOrderInput) (*entity.Order, error) {
	Item, err := gp.Repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return Item, nil
}
