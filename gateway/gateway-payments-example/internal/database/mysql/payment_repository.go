package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gateway-payments/internal/domain/entity"
	domainErrors "gateway-payments/internal/domain/errors"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Save(ctx context.Context, p *entity.Payment) error {
	query := `
		INSERT INTO payments (id, method, amount, status, order_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		method = VALUES(method),
		amount = VALUES(amount),
		status = VALUES(status),
		order_id = VALUES(order_id)
	`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Method, p.Amount, p.Status, p.OrderID, p.CreatedAt)
	if err != nil {
		return fmt.Errorf("error saving payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	p := &entity.Payment{}
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Method, &p.Amount, &p.Status, &p.OrderID, &p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("error finding payment: %w", err)
	}
	return p, nil
}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID string) (*entity.Payment, error) {
	p := &entity.Payment{}
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments WHERE order_id = ?`
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&p.ID, &p.Method, &p.Amount, &p.Status, &p.OrderID, &p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("error finding payment by order: %w", err)
	}
	return p, nil
}

func (r *paymentRepository) FindAll(ctx context.Context, page, limit int) ([]*entity.Payment, error) {
	offset := (page - 1) * limit
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing payments: %w", err)
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		p := &entity.Payment{}
		if err := rows.Scan(&p.ID, &p.Method, &p.Amount, &p.Status, &p.OrderID, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *paymentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM payments WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return domainErrors.ErrPaymentNotFound
	}
	return nil
}
