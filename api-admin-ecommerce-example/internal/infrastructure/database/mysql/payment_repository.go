package mysql

import (
	"database/sql"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/repository"
)

type PaymentRepository struct {
	DB *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) Save(payment *entity.Payment) error {
	if payment.Status == "" {
		payment.Status = entity.StatusPending
	}

	// Check if the payment already exists to decide between INSERT and UPDATE
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM payments WHERE id = ?)", payment.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if payment exists: %w", err)
	}

	if exists {
		query := `UPDATE payments SET method = ?, amount = ?, status = ?, order_id = ? WHERE id = ?`
		_, err := r.DB.Exec(
			query,
			payment.Method,
			payment.Amount,
			payment.Status,
			payment.OrderID,
			payment.ID,
		)
		if err != nil {
			return fmt.Errorf("error updating payment [%s]: %w", payment.ID, err)
		}
	} else {
		query := `INSERT INTO payments (id, method, amount, status, order_id, created_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err := r.DB.Exec(
			query,
			payment.ID,
			payment.Method,
			payment.Amount,
			payment.Status,
			payment.OrderID,
			payment.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("error persisting payment [%s]: %w", payment.ID, err)
		}
	}

	return nil
}

func (r *PaymentRepository) FindByID(id string) (*entity.Payment, error) {
	payment := &entity.Payment{}
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments WHERE id = ?`
	err := r.DB.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.Method,
		&payment.Amount,
		&payment.Status,
		&payment.OrderID,
		&payment.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment with ID %s not found", id)
		}
		return nil, fmt.Errorf("error finding payment by ID [%s]: %w", id, err)
	}

	return payment, nil
}

func (r *PaymentRepository) FindByOrderID(orderID string) (*entity.Payment, error) {
	payment := &entity.Payment{}
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments WHERE order_id = ?`
	err := r.DB.QueryRow(query, orderID).Scan(
		&payment.ID,
		&payment.Method,
		&payment.Amount,
		&payment.Status,
		&payment.OrderID,
		&payment.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &repository.ErrNotFound{Message: fmt.Sprintf("payment with order ID %s not found", orderID)}
		}
		return nil, fmt.Errorf("error finding payment by order ID [%s]: %w", orderID, err)
	}

	return payment, nil
}

func (r *PaymentRepository) FindAll(page, limit int) ([]*entity.Payment, error) {
	offset := (page - 1) * limit
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments LIMIT ? OFFSET ?`
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying payments: %w", err)
	}
	defer rows.Close()

	payments := make([]*entity.Payment, 0)
	for rows.Next() {
		payment := &entity.Payment{}
		if err := rows.Scan(
			&payment.ID,
			&payment.Method,
			&payment.Amount,
			&payment.Status,
			&payment.OrderID,
			&payment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning payment row: %w", err)
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return payments, nil
}

func (r *PaymentRepository) Delete(id string) error {
	query := `DELETE FROM payments WHERE id = ?`
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting payment [%s]: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected after delete: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment with ID %s not found for deletion", id)
	}

	return nil
}
