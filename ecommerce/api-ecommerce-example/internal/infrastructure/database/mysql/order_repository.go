package mysql

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/core/domain"
	"fmt"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback on error

	query := `
		INSERT INTO orders (id, total, status, method, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, query, order.ID, order.Total, string(order.Status), order.Method, order.CreatedAt)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO order_items (id, order_id, item_id, quantity, subtotal, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err = tx.ExecContext(ctx, itemQuery, item.ID, item.OrderID, item.ItemID, item.Quantity, item.Subtotal, item.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) FindAll(ctx context.Context, page, limit int) ([]*domain.Order, error) {
	offset := (page - 1) * limit
	query := `SELECT id, total, status, method, created_at FROM orders order by created_at LIMIT ? OFFSET ?`
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying orders: %w", err)
	}
	defer rows.Close()

	orders := make([]*domain.Order, 0)
	for rows.Next() {
		order := &domain.Order{}
		var status string
		if err := rows.Scan(
			&order.ID,
			&order.Total,
			&status,
			&order.Method,
			&order.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning order row: %w", err)
		}
		order.Status = domain.OrderStatus(status)

		items, err := r.FindOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	order := &domain.Order{}
	var status string
	err := r.DB.QueryRowContext(ctx, `SELECT id, total, status, method, created_at FROM orders WHERE id = ?`, id).Scan(&order.ID, &order.Total, &status, &order.Method, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	order.Status = domain.OrderStatus(status)

	items, err := r.FindOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func (r *OrderRepository) FindOrderItemsByOrderID(ctx context.Context, orderID string) ([]*domain.OrderItem, error) {
	query := `SELECT id, order_id, item_id, quantity, subtotal, created_at FROM order_items WHERE order_id = ?`
	rows, err := r.DB.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("error querying order items: %w", err)
	}
	defer rows.Close()

	orderItems := make([]*domain.OrderItem, 0)
	for rows.Next() {
		orderItem := &domain.OrderItem{}
		if err := rows.Scan(
			&orderItem.ID,
			&orderItem.OrderID,
			&orderItem.ItemID,
			&orderItem.Quantity,
			&orderItem.Subtotal,
			&orderItem.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning order item row: %w", err)
		}
		orderItems = append(orderItems, orderItem)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during order items rows iteration: %w", err)
	}

	return orderItems, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, order *domain.Order) error {
	query := `UPDATE orders SET status = ? WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, string(order.Status), order.ID)
	if err != nil {
		return err
	}
	return nil
}
