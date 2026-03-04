package mysql

import (
	"database/sql"
	"ecommerce-api/internal/domain/entity"
	"fmt"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback on error

	query := `
		INSERT INTO orders (id, total, status, method, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(query, order.ID, order.Total, order.Status, order.Method, order.CreatedAt)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO order_items (id, order_id, item_id, quantity, subtotal, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(itemQuery, item.ID, item.OrderID, item.ItemID, item.Quantity, item.Subtotal, item.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) SaveOrderItem(orderItem *entity.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, item_id, quantity, subtotal, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.DB.Exec(query, orderItem.ID, orderItem.OrderID, orderItem.ItemID, orderItem.Quantity, orderItem.Subtotal, orderItem.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) FindAll(page, limit int) ([]*entity.Order, error) {
	offset := (page - 1) * limit
	query := `SELECT id, total, status, method, created_at FROM orders order by created_at LIMIT ? OFFSET ?`
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying orders: %w", err)
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order := &entity.Order{}
		if err := rows.Scan(
			&order.ID,
			&order.Total,
			&order.Status,
			&order.Method,
			&order.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning order row: %w", err)
		}

		items, err := r.FindOrderItemsByOrderID(order.ID)
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

func (r *OrderRepository) FindByID(id string) (*entity.Order, error) {
	order := &entity.Order{}
	err := r.DB.QueryRow(`SELECT id, total, status, method, created_at FROM orders WHERE id = ?`, id).Scan(&order.ID, &order.Total, &order.Status, &order.Method, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}

	items, err := r.FindOrderItemsByOrderID(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func (r *OrderRepository) FindOrderItemsByOrderID(orderID string) ([]*entity.OrderItem, error) {
	query := `SELECT id, order_id, item_id, quantity, subtotal, created_at FROM order_items WHERE order_id = ?`
	rows, err := r.DB.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("error querying order items: %w", err)
	}
	defer rows.Close()

	orderItems := make([]*entity.OrderItem, 0)
	for rows.Next() {
		orderItem := &entity.OrderItem{}
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

func (r *OrderRepository) UpdateStatus(order *entity.Order) error {
	query := `UPDATE orders SET status = ? WHERE id = ?`
	_, err := r.DB.Exec(query, order.Status, order.ID)
	if err != nil {
		return err
	}
	return nil
}
