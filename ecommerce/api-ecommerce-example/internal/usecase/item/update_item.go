package item

import (
	"ecommerce-api/internal/domain/repository"
	"errors"
)

type UpdateItemInput struct {
	ID    string
	Price float64
	Name  string
}

type UpdateItem struct {
	Repo repository.ItemRepository
}

func NewUpdateItemUseCase(repo repository.ItemRepository) *UpdateItem {
	return &UpdateItem{
		Repo: repo,
	}
}

func (up *UpdateItem) Execute(input UpdateItemInput) error {
	Item, err := up.Repo.FindByID(input.ID)
	if err != nil {
		return errors.New("Item not found")
	}

	Item.Name = input.Name
	Item.Price = input.Price

	err = up.Repo.Save(Item) // Assuming Save also handles updates
	if err != nil {
		return err
	}

	return nil
}
