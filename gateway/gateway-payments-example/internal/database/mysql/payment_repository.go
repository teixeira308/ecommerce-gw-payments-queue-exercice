package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gateway-payments/internal/domain/entity"
	domainErrors "gateway-payments/internal/domain/errors"
	"log"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Save(ctx context.Context, p *entity.Payment) error {
	log.Printf("[Repository] Saving payment with ID: %s", p.ID)
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
		log.Printf("[Repository Error] failed to save payment %s: %v", p.ID, err)
		return fmt.Errorf("error saving payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	log.Printf("[Repository] Finding payment by ID: %s", id)
	p := &entity.Payment{}
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Method, &p.Amount, &p.Status, &p.OrderID, &p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[Repository Info] payment not found: %s", id)
			return nil, domainErrors.ErrPaymentNotFound
		}
		log.Printf("[Repository Error] failed to find payment %s: %v", id, err)
		return nil, fmt.Errorf("error finding payment: %w", err)
	}
	return p, nil
}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID string) (*entity.Payment, error) {
	log.Printf("[Repository] Finding payment by OrderID: %s", orderID)
	p := &entity.Payment{}
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments WHERE order_id = ?`
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&p.ID, &p.Method, &p.Amount, &p.Status, &p.OrderID, &p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[Repository Info] payment not found for OrderID: %s", orderID)
			return nil, domainErrors.ErrPaymentNotFound
		}
		log.Printf("[Repository Error] failed to find payment for OrderID %s: %v", orderID, err)
		return nil, fmt.Errorf("error finding payment by order: %w", err)
	}
	return p, nil
}

func (r *paymentRepository) FindAll(ctx context.Context, page, limit int) ([]*entity.Payment, error) {
	offset := (page - 1) * limit
	log.Printf("[Repository] Listing payments (page: %d, limit: %d, offset: %d)", page, limit, offset)
	query := `SELECT id, method, amount, status, order_id, created_at FROM payments LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Printf("[Repository Error] failed to list payments: %v", err)
		return nil, fmt.Errorf("error listing payments: %w", err)
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		p := &entity.Payment{}
		if err := rows.Scan(&p.ID, &p.Method, &p.Amount, &p.Status, &p.OrderID, &p.CreatedAt); err != nil {
			log.Printf("[Repository Error] failed to scan payment row: %v", err)
			return nil, err
		}
		payments = append(payments, p)
	}
	log.Printf("[Repository] Successfully retrieved %d payments", len(payments))
	return payments, nil
}

func (r *paymentRepository) Delete(ctx context.Context, id string) error {
	log.Printf("[Repository] Deleting payment with ID: %s", id)
	query := `DELETE FROM payments WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("[Repository Error] failed to delete payment %s: %v", id, err)
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		log.Printf("[Repository Info] no rows affected when deleting payment %s", id)
		return domainErrors.ErrPaymentNotFound
	}
	log.Printf("[Repository] Successfully deleted payment with ID: %s", id)
	return nil
}
