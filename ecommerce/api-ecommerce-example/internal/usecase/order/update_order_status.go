package order

import (
	"context"
	"ecommerce-api/internal/domain/repository"
	"fmt"
)

type UpdateOrderStatusInput struct {
	ID     string
	Status string
}

type UpdateOrderStatus struct {
	Repo repository.OrderRepository
}

func NewUpdateOrderStatusUseCase(repo repository.OrderRepository) *UpdateOrderStatus {
	return &UpdateOrderStatus{
		Repo: repo,
	}
}

func (uc *UpdateOrderStatus) Execute(ctx context.Context, input UpdateOrderStatusInput) error {
	order, err := uc.Repo.FindByID(input.ID)
	if err != nil {
		return err
	}

	if order == nil {
		return &repository.ErrNotFound{Message: "order not found"}
	}

	order.Status = input.Status

	fmt.Print(order)
	return uc.Repo.UpdateStatus(order)
}
