package item

import (
	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/internal/domain/repository"
)

type GetAllItemsInput struct {
	Page  int
	Limit int
}

type GetAllItemsOutput struct {
	Items []*entity.Item
}

type GetAllItems struct {
	Repo repository.ItemRepository
}

func NewGetAllItemsUseCase(repo repository.ItemRepository) *GetAllItems {
	return &GetAllItems{
		Repo: repo,
	}
}

func (gap *GetAllItems) Execute(input GetAllItemsInput) (*GetAllItemsOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}

	Items, err := gap.Repo.FindAll(input.Page, input.Limit)
	if err != nil {
		return nil, err
	}

	return &GetAllItemsOutput{Items: Items}, nil
}
