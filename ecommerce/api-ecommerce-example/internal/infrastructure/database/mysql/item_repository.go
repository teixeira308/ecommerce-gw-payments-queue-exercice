package mysql

import (
	"database/sql"
	"ecommerce-api/internal/domain/entity"
)

type ItemRepository struct {
	DB *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{DB: db}
}
func (r *ItemRepository) Save(item *entity.Item) error {
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
		_, err := r.DB.Exec(query, item.Name, item.Price, item.ID)
		if err != nil {
			return err
		}
	} else {
		// Item does not exist, perform an insert
		query := `
			INSERT INTO items (id, name, price)
			VALUES (?, ?, ?)
		`
		_, err := r.DB.Exec(query, item.ID, item.Name, item.Price)
		if err != nil {
			return err
		}
	}

	// Re-fetch to ensure created_at is populated for new items or updated for existing ones if needed
	// (though created_at usually doesn't change on update, this ensures consistency if the entity is
	// expected to have it after a Save call, for updates it will just re-assign the same value)
	return r.DB.QueryRow(
		`SELECT created_at FROM items WHERE id = ?`,
		item.ID,
	).Scan(&item.CreatedAt)
}

func (r *ItemRepository) FindByID(id string) (*entity.Item, error) {
	item := &entity.Item{}
	err := r.DB.QueryRow(`SELECT id, name, price, created_at FROM items WHERE id = ?`, id).Scan(&item.ID, &item.Name, &item.Price, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) FindAll(page, limit int) ([]*entity.Item, error) {
	offset := (page - 1) * limit
	query := `SELECT id, name, price, created_at FROM items LIMIT ? OFFSET ?`
	rows, err := r.DB.Query(query, limit, offset)
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
	_, err := r.DB.Exec(`DELETE FROM items WHERE id = ?`, id)
	return err
}
