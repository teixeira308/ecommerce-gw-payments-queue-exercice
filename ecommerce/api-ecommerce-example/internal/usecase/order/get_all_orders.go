package order

import (
	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/internal/domain/repository"
)

type GetAllOrdersInput struct {
	Page  int
	Limit int
}

type GetAllOrdersOutput struct {
	Orders []*entity.Order
}

type GetAllOrders struct {
	Repo repository.OrderRepository
}

func NewGetAllOrdersUseCase(repo repository.OrderRepository) *GetAllOrders {
	return &GetAllOrders{
		Repo: repo,
	}
}

func (gap *GetAllOrders) Execute(input GetAllOrdersInput) (*GetAllOrdersOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}

	Orders, err := gap.Repo.FindAll(input.Page, input.Limit)
	if err != nil {
		return nil, err
	}

	return &GetAllOrdersOutput{Orders: Orders}, nil
}
