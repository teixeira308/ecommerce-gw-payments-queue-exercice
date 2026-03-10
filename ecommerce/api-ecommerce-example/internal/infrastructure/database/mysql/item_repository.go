package mysql

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/internal/domain/repository"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type ItemRepository struct {
	DB *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{DB: db}
}

func (r *ItemRepository) Save(item *entity.Item) error {
	ctx := context.Background()
	// Check if the item already exists
	existingItem, err := r.FindByID(item.ID)
	if err != nil {
		return err
	}

	if existingItem != nil {
		// Item exists, perform an update
		query := `
			UPDATE items
			SET name = ?, price = ?
			WHERE id = ?
		`
		_, err := r.DB.ExecContext(ctx, query, item.Name, item.Price, item.ID)
		if err != nil {
			return err
		}
	} else {
		// Item does not exist, perform an insert
		query := `
			INSERT INTO items (id, name, price)
			VALUES (?, ?, ?)
		`
		_, err := r.DB.ExecContext(ctx, query, item.ID, item.Name, item.Price)
		if err != nil {
			return err
		}
	}

	return r.DB.QueryRowContext(ctx,
		"SELECT created_at FROM items WHERE id = ?",
		item.ID,
	).Scan(&item.CreatedAt)
}

func (r *ItemRepository) FindByID(id string) (*entity.Item, error) {
	ctx := context.Background()
	item := &entity.Item{}
	err := r.DB.QueryRowContext(ctx, "SELECT id, name, price, created_at FROM items WHERE id = ?", id).Scan(&item.ID, &item.Name, &item.Price, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) FindAll(page, limit int) ([]*entity.Item, error) {
	ctx := context.Background()
	offset := (page - 1) * limit
	query := "SELECT id, name, price, created_at FROM items LIMIT ? OFFSET ?"
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.Item

	for rows.Next() {
		i := &entity.Item{}
		if err := rows.Scan(&i.ID, &i.Name, &i.Price, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func (r *ItemRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := r.DB.ExecContext(ctx, "DELETE FROM items WHERE id = ?", id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
			return &repository.ErrConflict{Message: fmt.Sprintf("cannot delete item with ID %s because it is referenced by other records (e.g., in an order)", id)}
		}
		return err
	}
	return nil
}
