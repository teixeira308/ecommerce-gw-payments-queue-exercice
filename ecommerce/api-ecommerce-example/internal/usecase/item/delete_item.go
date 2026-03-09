package item

import (
	"ecommerce-api/internal/domain/repository"
)

type DeleteItemInput struct {
	ID string
}

type DeleteItem struct {
	Repo repository.ItemRepository
}

func NewDeleteItemUseCase(repo repository.ItemRepository) *DeleteItem {
	return &DeleteItem{
		Repo: repo,
	}
}

func (dp *DeleteItem) Execute(input DeleteItemInput) error {
	err := dp.Repo.Delete(input.ID)
	if err != nil {
		return err
	}
	return nil
}
