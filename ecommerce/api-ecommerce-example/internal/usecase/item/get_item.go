package item

import (
	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/internal/domain/repository"
)

type GetItemInput struct {
	ID string
}

type GetItem struct {
	Repo repository.ItemRepository
}

func NewGetItemUseCase(repo repository.ItemRepository) *GetItem {
	return &GetItem{
		Repo: repo,
	}
}

func (gp *GetItem) Execute(input GetItemInput) (*entity.Item, error) {
	Item, err := gp.Repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return Item, nil
}
