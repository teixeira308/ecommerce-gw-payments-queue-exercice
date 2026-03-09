package item

import (
	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/internal/domain/repository"

	"github.com/google/uuid"
)

type CreateItem struct {
	Repo repository.ItemRepository
}

func NewCreateItemUseCase(repo repository.ItemRepository) *CreateItem {
	return &CreateItem{
		Repo: repo,
	}
}

func (uc *CreateItem) Execute(name string, price float64) (*entity.Item, error) {

	item := &entity.Item{
		ID:    uuid.NewString(),
		Name:  name,
		Price: price,
	}

	err := uc.Repo.Save(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
