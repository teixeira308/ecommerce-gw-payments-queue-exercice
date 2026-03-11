package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/repository"
	"log"
)

type PaymentRepository struct {
	DB *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

// Save implements atomic UPSERT (Insert on Duplicate Key Update) for MySQL
func (r *PaymentRepository) Save(ctx context.Context, payment *entity.Payment) error {
	if payment.Status == "" {
		payment.Status = entity.StatusPending
	}

	log.Printf("[DB] Saving payment %s for order %s", payment.ID, payment.OrderID)

	query := `
		INSERT INTO payments (id, method, amount, status, order_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			method = VALUES(method),
			amount = VALUES(amount),
			status = VALUES(status),
			order_id = VALUES(order_id)
	`

	_, err := r.DB.ExecContext(ctx, query,
		payment.ID,
		payment.Method,
		payment.Amount,
		payment.Status,
		payment.OrderID,
		payment.CreatedAt,
	)

	if err != nil {
		log.Printf("[DB] Error saving payment %s: %v", payment.ID, err)
		return fmt.Errorf("failed to save payment: %w", err)
	}

	log.Printf("[DB] Payment %s saved successfully", payment.ID)
	return nil
}

func (r *PaymentRepository) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	log.Printf("[DB] Finding payment by ID: %s", id)
	payment := &entity.Payment{}
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments WHERE id = ?`
	
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.Method,
		&payment.Amount,
		&payment.Status,
		&payment.OrderID,
		&payment.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[DB] Payment not found: %s", id)
			return nil, repository.NewErrNotFound(id)
		}
		log.Printf("[DB] Error finding payment %s: %v", id, err)
		return nil, fmt.Errorf("database error finding payment: %w", err)
	}

	return payment, nil
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*entity.Payment, error) {
	log.Printf("[DB] Finding payment by Order ID: %s", orderID)
	payment := &entity.Payment{}
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments WHERE order_id = ?`
	
	err := r.DB.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.Method,
		&payment.Amount,
		&payment.Status,
		&payment.OrderID,
		&payment.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[DB] Payment not found for order: %s", orderID)
			return nil, repository.NewErrNotFound(orderID)
		}
		log.Printf("[DB] Error finding payment for order %s: %v", orderID, err)
		return nil, fmt.Errorf("database error finding payment by order: %w", err)
	}

	return payment, nil
}

func (r *PaymentRepository) FindAll(ctx context.Context, page, limit int) ([]*entity.Payment, error) {
	log.Printf("[DB] Finding all payments (page=%d, limit=%d)", page, limit)
	offset := (page - 1) * limit
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments LIMIT ? OFFSET ?`
	
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Printf("[DB] Error listing payments: %v", err)
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		p := &entity.Payment{}
		if err := rows.Scan(
			&p.ID,
			&p.Method,
			&p.Amount,
			&p.Status,
			&p.OrderID,
			&p.CreatedAt,
		); err != nil {
			log.Printf("[DB] Error scanning payment row: %v", err)
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, p)
	}

	log.Printf("[DB] Found %d payments", len(payments))
	return payments, nil
}

func (r *PaymentRepository) Delete(ctx context.Context, id string) error {
	log.Printf("[DB] Deleting payment: %s", id)
	query := `DELETE FROM payments WHERE id = ?`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("[DB] Error deleting payment %s: %v", id, err)
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		log.Printf("[DB] No payment found to delete: %s", id)
		return repository.NewErrNotFound(id)
	}

	log.Printf("[DB] Payment %s deleted successfully", id)
	return nil
}
